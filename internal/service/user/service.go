package user

import (
	"github.com/BelyaevEI/microservices_auth/internal/client/db"
	"github.com/BelyaevEI/microservices_auth/internal/repository"
	"github.com/BelyaevEI/microservices_auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

// NewService creates a new user service.
func NewService(userRepository repository.UserRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
