package service

import (
	"context"
	"qingmo-server/internal/dto"
	"qingmo-server/internal/model"
	"qingmo-server/internal/repository"
	"strings"
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

	return profileResponse(user), nil
}

func (u *UserService) UpdateProfile(
	ctx context.Context,
	userID uint,
	req dto.UpdateProfileRequest,
) (*dto.ProfileResponse, error) {

	if req.UserName != nil {
		name := strings.TrimSpace(*req.UserName)
		req.UserName = &name
	}
	if req.AvatarURL != nil {
		url := strings.TrimSpace(*req.AvatarURL)
		req.AvatarURL = &url
	}

	user, err := u.users.UpdateProfile(ctx, userID, req)
	if err != nil {
		return nil, err
	}
	return profileResponse(user), nil
}

func (u *UserService) PublicProfile(
	ctx context.Context,
	userID uint,
) (*dto.PublicProfileResponse, error) {
	user, err := u.users.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.PublicProfileResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		AvatarURL: user.AvatarURL,
	}, nil
}

func profileResponse(user *model.Users) *dto.ProfileResponse {
	return &dto.ProfileResponse{
		UserName:  user.UserName,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}
}
