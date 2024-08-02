package user

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
)

func (s serv) CreateUser(ctx context.Context, userCreate *model.UserCreate) (int64, error) {

	id, err := s.userRepository.CreateUser(ctx, userCreate)
	if err != nil {
		return 0, err
	}

	return id, nil
}
