package main

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-auth/internal/config"
	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/patyukin/mbs-auth/internal/server"
	"github.com/patyukin/mbs-auth/internal/usecase"
	desc "github.com/patyukin/mbs-auth/pkg/auth_v1"
	"github.com/patyukin/mbs-auth/pkg/dbconn"
	"github.com/patyukin/mbs-auth/pkg/migrator"
	"github.com/patyukin/mbs-auth/pkg/utils"
	_ "github.com/patyukin/mbs-auth/statik"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("failed to load config, error: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCServer.Port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	dbConn, err := dbconn.New(context.Background(), dbconn.PostgreSQLConfig(cfg.PostgreSQL))
	if err != nil {
		log.Fatal().Msgf("failed to connect to db: %v", err)
	}

	if err = migrator.UpMigrations(context.Background(), dbConn); err != nil {
		log.Fatal().Msgf("failed to up migrations: %v", err)
	}

	registry := db.New(dbConn)
	uc := usecase.New(registry)
	srv := server.New(uc)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthServiceServer(s, srv)

	log.Printf("server listening at %v", lis.Addr())

	wg := &sync.WaitGroup{}

	// GRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Info().Msgf("GRPC started on :%d", cfg.GRPCServer.Port)
		if err = s.Serve(lis); err != nil {
			log.Fatal().Msgf("failed to serve: %v", err)
		}
	}()

	// metrics server
	wg.Add(1)
	go func() {
		defer wg.Done()

		http.Handle("/metrics", promhttp.Handler())
		log.Info().Msgf("Prometheus metrics exposed on :%d/metrics", cfg.HttpServer.Port)
		if err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpServer.Port), nil); err != nil {
			log.Fatal().Msgf("Failed to serve Prometheus metrics: %v", err)
		}
	}()

	// swagger server
	wg.Add(1)
	go func() {
		defer wg.Done()

		statikFs, errStatikFs := fs.New()
		if errStatikFs != nil {
			log.Fatal().Msgf("failed to create statik fs: %v", errStatikFs)
		}

		mux := http.NewServeMux()
		mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
		mux.HandleFunc("/api.swagger.json", utils.ServeSwaggerFile("/api.swagger.json"))

		swaggerSrv := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.SwaggerServer.Port),
			Handler: mux,
		}

		log.Printf("swagger server started on port: %v", swaggerSrv.Addr)
		if err = swaggerSrv.ListenAndServe(); err != nil {
			log.Fatal().Msgf("failed to serve swagger: %v", err)
		}
	}()

	wg.Wait()
}
