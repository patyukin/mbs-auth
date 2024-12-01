package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/patyukin/mbs-auth/internal/cacher"
	"github.com/patyukin/mbs-auth/internal/config"
	"github.com/patyukin/mbs-auth/internal/cronjob"
	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/patyukin/mbs-auth/internal/metrics"
	"github.com/patyukin/mbs-auth/internal/server"
	"github.com/patyukin/mbs-auth/internal/usecase"
	"github.com/patyukin/mbs-pkg/pkg/dbconn"
	"github.com/patyukin/mbs-pkg/pkg/kafka"
	"github.com/patyukin/mbs-pkg/pkg/migrator"
	"github.com/patyukin/mbs-pkg/pkg/mux_server"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/patyukin/mbs-pkg/pkg/rabbitmq"
	"github.com/patyukin/mbs-pkg/pkg/tracing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/reflection"
)

const ServiceName = "AuthService"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("failed to load config, error: %v", err)
	}

	if err = metrics.Init(); err != nil {
		log.Fatal().Msgf("failed to init metrics: %v", err)
	}

	_, closer, err := tracing.InitJaeger(fmt.Sprintf(cfg.TracerHost), ServiceName)
	if err != nil {
		log.Fatal().Msgf("failed to initialize tracer: %v", err)
	}

	defer closer()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCServer.Port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	dbConn, err := dbconn.New(ctx, cfg.PostgreSQLDSN)
	if err != nil {
		log.Fatal().Msgf("failed to connect to db: %v", err)
	}

	if err = migrator.UpMigrations(ctx, dbConn); err != nil {
		log.Fatal().Msgf("failed to up migrations: %v", err)
	}

	rbt, err := rabbitmq.New(cfg.RabbitMQUrl, rabbitmq.Exchange)
	if err != nil {
		log.Fatal().Msgf("failed to create rabbit producer: %v", err)
	}

	err = rbt.BindQueueToExchange(
		rabbitmq.Exchange,
		rabbitmq.TelegramMessageQueue,
		[]string{rabbitmq.TelegramMessageRouteKey},
	)
	if err != nil {
		log.Fatal().Msgf("failed to bind NotifyAuthQueue to exchange with - NotifySignUpConfirmCodeRouteKey: %v", err)
	}

	kfk, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroup, cfg.Kafka.Topics)
	if err != nil {
		log.Fatal().Msgf("failed to create kafka consumer, err: %v", err)
	}

	chr, err := cacher.New(ctx, cfg.RedisDSN)
	if err != nil {
		log.Fatal().Msgf("failed to create redis cacher: %v", err)
	}

	registry := db.New(dbConn)
	uc := usecase.New(registry, rbt, chr, cfg)
	srv := server.New(uc)

	// grpc server
	s := server.NewGRPCServer(cfg)
	reflection.Register(s)
	desc.RegisterAuthServiceServer(s, srv)
	grpcPrometheus.Register(s)

	// http server
	muxServer := mux_server.New()

	errCh := make(chan error)

	// cron job
	cj := cronjob.New(uc)
	uc.RemoveNotRegisteredUsers(ctx) // temp
	go func() {
		if err = cj.Run(ctx); err != nil {
			log.Error().Msgf("failed adding cron job, err: %v", err)
			errCh <- err
		}
	}()

	// run consumer
	go func() {
		if err = kfk.ProcessMessages(ctx, uc.RegistrationSolutionProcess); err != nil {
			log.Error().Msgf("failed to process messages: %v", err)
			errCh <- err
		}
	}()

	// GRPC server
	go func() {
		log.Info().Msgf("GRPC started on :%d", cfg.GRPCServer.Port)
		if err = s.Serve(lis); err != nil {
			log.Error().Msgf("failed to serve: %v", err)
			errCh <- err
		}
	}()

	// metrics + pprof server
	go func() {
		if err = muxServer.Run(cfg.HttpServer.Port); err != nil {
			log.Error().Msgf("Failed to serve Prometheus metrics: %v", err)
			errCh <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err = <-errCh:
		log.Error().Msgf("Failed to run, err: %v", err)
	case res := <-sigChan:
		if res == syscall.SIGINT || res == syscall.SIGTERM {
			log.Info().Msg("Signal received")
		} else if res == syscall.SIGHUP {
			log.Info().Msg("Signal received")
		}
	}

	log.Info().Msg("Shutting Down")

	// stop servers
	s.GracefulStop()
	if err = muxServer.Shutdown(ctx); err != nil {
		log.Error().Msgf("failed to shutdown http server: %v", err)
	}

	if err = dbConn.Close(); err != nil {
		log.Error().Msgf("failed db connection close: %s", err.Error())
	}

	if err = chr.Close(); err != nil {
		log.Error().Msgf("failed redis connection close: %s", err.Error())
	}

	cj.Stop()
}
