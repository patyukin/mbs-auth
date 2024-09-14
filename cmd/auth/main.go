package main

import (
	"fmt"
	"github.com/patyukin/mbs-auth/internal/config"
	desc "github.com/patyukin/mbs-auth/pkg/auth_v1"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"time"
)

type server struct {
	desc.UnimplementedAuthServiceServer
}

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

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}
