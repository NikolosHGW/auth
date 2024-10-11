package main

import (
	"context"
	"fmt"
	"log"
	"net"

	userServer "github.com/NikolosHGW/auth/internal/api/user"
	"github.com/NikolosHGW/auth/internal/infrastructure/config"
	"github.com/NikolosHGW/auth/internal/infrastructure/config/env"
	userRepo "github.com/NikolosHGW/auth/internal/infrastructure/db/repository/user"
	userService "github.com/NikolosHGW/auth/internal/service/user"
	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pgConfig, err := env.NewPG()
	if err != nil {
		log.Fatal(err)
	}

	grpcConfig, err := env.NewGRPC()
	if err != nil {
		log.Fatal(err)
	}

	pgxCon, err := pgxpool.New(context.Background(), pgConfig.GetDatabaseDSN())
	if err != nil {
		log.Fatal(err)
	}

	listen, err := net.Listen("tcp", grpcConfig.GetRunAddress())
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	reflection.Register(s)

	userRepo := userRepo.NewUser(pgxCon)
	userService := userService.NewService(userRepo)

	userpb.RegisterUserV1Server(s, userServer.NewImplementation(userService))

	fmt.Println("Сервер gRPC начал работу")
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
