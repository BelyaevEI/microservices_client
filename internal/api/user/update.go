package user

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/converter"
	desc "github.com/BelyaevEI/microservices_auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Update updates a user.
func (i *Implementation) UpdateUserByID(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.UpdateUserByID(ctx, req.GetId(), converter.ToUserUpdateFromDesc(req))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
