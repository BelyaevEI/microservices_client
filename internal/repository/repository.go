package repository

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
)

// UserRepository represents a repository for user entities.
type UserRepository interface {
	CreateUser(ctx context.Context, userCreate *model.UserCreate) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	UpdateUserByID(ctx context.Context, id int64, userUpdate *model.UserUpdate) error
	DeleteUserByID(ctx context.Context, id int64) error
}
