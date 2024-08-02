package service

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
)

// UserService represents a service for user entities.
type UserService interface {
	CreateUser(ctx context.Context, userCreate *model.UserCreate) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	UpdateUserByID(ctx context.Context, id int64, userUpdate *model.UserUpdate) error
	DeleteUserByID(ctx context.Context, id int64) error
}
