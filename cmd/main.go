package main

import (
	"context"
	"fmt"
	"log"
	"net"

	desc "github.com/BelyaevEI/microservices_client/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedAuthV1Server
}

// Create new user
func (s *server) CreateUser(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("User email: %s", req.GetInfo().GetEmail())
	return &desc.CreateResponse{
		Id: int64(gofakeit.Number(0, 1000)),
	}, nil
}

// Get user by id
func (s *server) GetUserByID(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())
	return &desc.GetResponse{
		User: &desc.User{
			Id: req.Id,
			Info: &desc.UserInfo{
				Name:  gofakeit.BeerName(),
				Email: gofakeit.Email(),
				Role:  1,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

// Update user
func (s *server) UpdateUserByID(_ context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("User email: %s", req.GetInfo().GetEmail())
	return nil, nil
}

// Delete user
func (s *server) DeleteUserByID(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("User id: %d", req.GetId())
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("listen is failed: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("serve is failed: %v", err)
	}
}
