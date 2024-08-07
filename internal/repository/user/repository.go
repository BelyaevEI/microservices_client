package user

import (
	"github.com/BelyaevEI/microservices_auth/internal/client/db"
	"github.com/BelyaevEI/microservices_auth/internal/repository"
)

const (
	tableName = "user"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	passHashColumn  = "pass_hash"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

// NewRepository creates a new user repository.
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}
