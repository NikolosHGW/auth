package main

import (
	"context"
	"fmt"
	"log"
	"net"

	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userServer struct {
	userpb.UserV1Server
}

const grpcPort = 3200

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	userpb.RegisterUserV1Server(s, &userServer{})

	fmt.Println("Сервер gRPC начал работу")
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

func (s *userServer) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.CreateResponse, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	fmt.Println(req.Name, req.Email, req.Password, req.PasswordConfirm, req.Role)

	return &userpb.CreateResponse{
		Id: 1,
	}, nil
}

func (s *userServer) Get(ctx context.Context, req *userpb.GetRequest) (*userpb.GetResponse, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	fmt.Println(req.Id)

	return &userpb.GetResponse{
		Id:        req.Id,
		Name:      "Вася",
		Email:     "qwe@mailfake.com",
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}, nil
}

func (s *userServer) Update(ctx context.Context, req *userpb.UpdateRequest) (*emptypb.Empty, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	fmt.Println(req.Id, req.Name, req.Email)

	return &emptypb.Empty{}, nil
}

func (s *userServer) Delete(ctx context.Context, req *userpb.DeleteRequest) (*emptypb.Empty, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	fmt.Println(req.Id)

	return &emptypb.Empty{}, nil
}
