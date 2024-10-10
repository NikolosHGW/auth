package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/NikolosHGW/auth/internal/infrastructure/config"
	"github.com/NikolosHGW/auth/internal/infrastructure/config/env"
	"github.com/NikolosHGW/auth/internal/infrastructure/db/repository"
	"github.com/NikolosHGW/auth/internal/infrastructure/db/repository/user"
	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userServer struct {
	userpb.UserV1Server
	userRepo repository.UserRepository
}

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
	userpb.RegisterUserV1Server(s, &userServer{userRepo: user.NewUser(pgxCon)})

	fmt.Println("Сервер gRPC начал работу")
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

func (s *userServer) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.CreateResponse, error) {
	id, err := s.userRepo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка создании: %w", err)
	}
	fmt.Println("успешно: ", req.Name, req.Email, req.Password, req.PasswordConfirm, req.Role)

	return &userpb.CreateResponse{
		Id: id,
	}, nil
}

func (s *userServer) Get(ctx context.Context, req *userpb.GetRequest) (*userpb.GetResponse, error) {
	response, err := s.userRepo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения: %w", err)
	}
	fmt.Println("успешно: ", req.Id)

	return response, nil
}

func (s *userServer) Update(ctx context.Context, req *userpb.UpdateRequest) (*emptypb.Empty, error) {
	err := s.userRepo.UpdateByID(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении: %w", err)
	}
	fmt.Println(req.Id, req.Name, req.Email)

	return &emptypb.Empty{}, nil
}

func (s *userServer) Delete(ctx context.Context, req *userpb.DeleteRequest) (*emptypb.Empty, error) {
	err := s.userRepo.DeleteByID(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при удалении: %w", err)
	}
	fmt.Println(req.Id)

	return &emptypb.Empty{}, nil
}
