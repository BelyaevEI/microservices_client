package user

import (
	"github.com/BelyaevEI/microservices_auth/internal/repository"
	"github.com/BelyaevEI/microservices_auth/internal/service"
	"github.com/BelyaevEI/platform_common/pkg/db"
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
