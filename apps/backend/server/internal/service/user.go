package service

import (
	"context"
	"qingmo-server/internal/dto"
	"qingmo-server/internal/repository"
)

type UserService struct {
	users *repository.UserPostgresRepository
}

func NewUserService(users *repository.UserPostgresRepository) *UserService {
	return &UserService{
		users: users,
	}
}

func (u *UserService) Profile(
	ctx context.Context,
	userID uint,
) (*dto.ProfileResponse, error) {
	user, err := u.users.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.ProfileResponse{
		UserName:  user.UserName,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}
