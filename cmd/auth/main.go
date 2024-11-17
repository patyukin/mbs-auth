package main

import (
	"context"
	"fmt"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/patyukin/mbs-auth/internal/cacher"
	"github.com/patyukin/mbs-auth/internal/config"
	"github.com/patyukin/mbs-auth/internal/cronjob"
	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/patyukin/mbs-auth/internal/metrics"
	"github.com/patyukin/mbs-auth/internal/server"
	"github.com/patyukin/mbs-auth/internal/usecase"
	"github.com/patyukin/mbs-pkg/pkg/dbconn"
	"github.com/patyukin/mbs-pkg/pkg/migrator"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/patyukin/mbs-pkg/pkg/rabbitmq"
	"github.com/patyukin/mbs-pkg/pkg/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	_, closer, err := tracing.InitJaeger(fmt.Sprintf("localhost:6831"), ServiceName)
	if err != nil {
		log.Fatal().Msgf("failed to initialize tracer: %v", err)
	}

	defer closer()

	log.Info().Msg("Jaeger connected")

	log.Info().Msg("Opentracing connected")
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
		rabbitmq.NotifyAuthQueue,
		[]string{rabbitmq.NotifySignUpConfirmCodeRouteKey},
	)
	if err != nil {
		log.Fatal().Msgf("failed to bind NotifyAuthQueue to exchange with - NotifySignUpConfirmCodeRouteKey: %v", err)
	}

	err = rbt.BindQueueToExchange(
		rabbitmq.Exchange,
		rabbitmq.AuthNotifyQueue,
		[]string{rabbitmq.AuthSignInConfirmCodeRouteKey, rabbitmq.AuthSignUpResultMessageRouteKey},
	)
	if err != nil {
		log.Fatal().Msgf(
			"failed to bind AuthNotifyQueue to exchange with - AuthSignInConfirmCodeRouteKey, "+
				"AuthSignUpResultMessageRouteKey: %v", err,
		)
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

	log.Printf("server listening at %v", lis.Addr())

	errCh := make(chan error)

	// cron job
	cj := cronjob.New(uc)
	uc.RemoveNotRegisteredUsers(ctx)
	go func() {
		if err = cj.Run(ctx); err != nil {
			log.Error().Msgf("failed adding cron job, err: %v", err)
			errCh <- err
		}
	}()

	// run consumer
	go func() {
		if err = rbt.Consume(ctx, rabbitmq.NotifyAuthQueue, uc.NotifyConsumeHandler); err != nil {
			log.Error().Msgf("failed to consume from queue: %v", err)
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
		http.Handle("/metrics", promhttp.Handler())
		log.Info().Msgf("Prometheus metrics exposed on :%d/metrics", cfg.HttpServer.Port)
		if err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpServer.Port), nil); err != nil {
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

	// stop server
	s.GracefulStop()

	if err = dbConn.Close(); err != nil {
		log.Error().Msgf("failed db connection close: %s", err.Error())
	}

	if err = chr.Close(); err != nil {
		log.Error().Msgf("failed redis connection close: %s", err.Error())
	}

	cj.Stop()

	//if err = tracerProvider.Shutdown(ctx); err != nil {
	//	log.Error().Msgf("Error shutting down tracer provider: %v", err)
	//}
}
