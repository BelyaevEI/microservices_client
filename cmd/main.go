package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/BelyaevEI/microservices_client/internal/config"
	"github.com/BelyaevEI/microservices_client/internal/utils"
	desc "github.com/BelyaevEI/microservices_client/pkg/auth_v1"
)

var configPath string

type server struct {
	desc.UnimplementedAuthV1Server
	pool *pgxpool.Pool
}

func init() {
	configPath = os.Getenv("CONFIG_PATH")
}

func main() {

	ctx := context.Background()
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Create new user
func (s *server) CreateUser(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {

	if req.GetInfo().GetRole() == desc.Role_UNKNOWN {
		return nil, status.Error(codes.PermissionDenied, "role is unknown")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid password")
	}

	passHash, err := utils.GetHashPassword(req.GetPassword())
	if err != nil {
		log.Printf("create User error: %v", err)
		return nil, status.Error(codes.Internal, "calculate hash is failed")
	}

	builderInsert := sq.Insert("user").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "pass_hash", "role").
		Values(req.GetInfo().GetName(), req.GetInfo().GetEmail(), passHash, req.GetInfo().GetRole()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("create user builder is failed %v", err)
		return nil, status.Error(codes.Internal, "failed to build query")
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Printf("create user scan error: %v", err)
		return nil, status.Error(codes.Internal, "failed to insert user")
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

// Get user by id
func (s *server) GetUserByID(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("user").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("get user select failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to build query")
	}

	var (
		id, role    int64
		name, email string
		createdAt   time.Time
		updatedAt   time.Time
	)

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("get user query is failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to select users")
	}

	return &desc.GetResponse{
		User: &desc.User{
			Id: req.GetId(),
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  desc.Role(role),
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		},
	}, nil
}

// Update user
func (s *server) UpdateUserByID(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {

	if req.GetInfo().GetName() == nil && req.GetInfo().GetEmail() == nil && req.GetInfo().GetRole() == desc.Role_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "nothing update")
	}

	// Create builder for update user
	builderUpdate := sq.Update("user").
		PlaceholderFormat(sq.Dollar)

	if req.GetInfo().GetName() != nil {
		builderUpdate = builderUpdate.Set("name", req.GetInfo().GetName().GetValue())
	}

	if req.GetInfo().GetEmail() != nil {
		builderUpdate = builderUpdate.Set("email", req.GetInfo().GetEmail().GetValue())
	}
	if req.GetInfo().GetRole() != desc.Role_UNKNOWN {
		builderUpdate = builderUpdate.Set("role", int(req.GetInfo().GetRole()))
	}
	builderUpdate = builderUpdate.Set("updated_at", time.Now())

	builderUpdate = builderUpdate.Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("update user: %e", err)
		return nil, status.Error(codes.Internal, "failed to build query")
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return nil, nil
}

// Delete user
func (s *server) DeleteUserByID(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {

	builderDelete := sq.Delete("user").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("delete user: %v", err)
		return nil, status.Error(codes.Internal, "failed to build query")
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	log.Printf("deleted %d rows", res.RowsAffected())
	return nil, nil

}
