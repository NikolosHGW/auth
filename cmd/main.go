package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/NikolosHGW/auth/internal/entity"
	"github.com/NikolosHGW/auth/internal/infrastructure/config"
	"github.com/NikolosHGW/auth/internal/infrastructure/config/env"
	"github.com/NikolosHGW/auth/internal/infrastructure/db/repository"
	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userRepo interface {
	Create(context.Context, *entity.User) (*entity.User, error)
	GetByID(context.Context, int64) (*entity.User, error)
	UpdateByID(context.Context, *entity.User) error
	DeleteByID(context.Context, int64) error
}

type userServer struct {
	userpb.UserV1Server
	userRepo userRepo
}

const grpcPort = 3200

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	pgConfig, err := env.NewPG()
	if err != nil {
		log.Fatal(err)
	}

	pgxCon, err := pgx.Connect(context.Background(), pgConfig.GetDatabaseDSN())
	if err != nil {
		log.Fatal(err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	userpb.RegisterUserV1Server(s, &userServer{userRepo: repository.NewUser(pgxCon)})

	fmt.Println("Сервер gRPC начал работу")
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

func (s *userServer) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.CreateResponse, error) {
	user, err := s.userRepo.Create(
		ctx,
		&entity.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
			Role:     int(req.Role),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка создании: %w", err)
	}
	fmt.Println(req.Name, req.Email, req.Password, req.PasswordConfirm, req.Role)

	return &userpb.CreateResponse{
		Id: user.ID,
	}, nil
}

func (s *userServer) Get(ctx context.Context, req *userpb.GetRequest) (*userpb.GetResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения: %w", err)
	}
	fmt.Println(req.Id)

	createdAt := timestamppb.New(user.CreatedAt)

	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return &userpb.GetResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (s *userServer) Update(ctx context.Context, req *userpb.UpdateRequest) (*emptypb.Empty, error) {
	err := s.userRepo.UpdateByID(
		ctx,
		&entity.User{ID: req.Id, Name: req.Name.Value, Email: req.Email.Value},
	)
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
