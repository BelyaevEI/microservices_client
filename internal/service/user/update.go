package user

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
)

func (s serv) UpdateUserByID(ctx context.Context, id int64, userUpdate *model.UserUpdate) error {

	err := s.userRepository.UpdateUserByID(ctx, id, userUpdate)
	if err != nil {
		return err
	}

	return nil
}
