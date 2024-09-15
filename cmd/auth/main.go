package main

import (
	"fmt"
	"github.com/patyukin/mbs-auth/internal/config"
	"github.com/patyukin/mbs-auth/internal/server"
	"github.com/patyukin/mbs-auth/internal/usecase"
	desc "github.com/patyukin/mbs-auth/pkg/auth_v1"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
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

	uc := usecase.New()
	srv := server.New(uc)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthServiceServer(s, srv)

	log.Printf("server listening at %v", lis.Addr())

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		statikFs, errStatikFs := fs.New()
		if err != nil {
			log.Fatal().Msgf("failed to create statik fs: %v", errStatikFs)
		}

		mux := http.NewServeMux()
		mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
		mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

		swaggerSrv := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.HttpServer.Port),
			Handler: mux,
		}

		log.Printf("start swagger server: %v", swaggerSrv.Addr)

		if err = swaggerSrv.ListenAndServe(); err != nil {
			log.Fatal().Msgf("failed to serve swagger: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err = s.Serve(lis); err != nil {
			log.Fatal().Msgf("failed to serve: %v", err)
		}
	}()

	wg.Wait()
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}
