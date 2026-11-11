package service

import (
	"context"
	"net/url"
	"qingmo-server/internal/dto"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/model"
	"qingmo-server/internal/repository"
	"strings"
	"unicode/utf8"
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
	if !req.UserName.Present && !req.AvatarURL.Present {
		return nil, errs.New(errs.KindInvalid, "provide userName or avatarURL")
	}

	update := repository.ProfileUpdate{SetAvatarURL: req.AvatarURL.Present}
	if req.UserName.Present {
		if req.UserName.Value == nil {
			return nil, errs.New(errs.KindInvalid, "userName cannot be null")
		}
		name := strings.TrimSpace(*req.UserName.Value)
		if name == "" || utf8.RuneCountInString(name) > 32 ||
			strings.ContainsRune(name, 0) {
			return nil, errs.New(
				errs.KindInvalid,
				"userName must be 1–32 characters without NUL",
			)
		}
		update.UserName = &name
	}
	if req.AvatarURL.Present && req.AvatarURL.Value != nil {
		avatar := strings.TrimSpace(*req.AvatarURL.Value)
		parsed, err := url.Parse(avatar)
		if err != nil ||
			(parsed.Scheme != "http" && parsed.Scheme != "https") ||
			parsed.Hostname() == "" ||
			strings.ContainsRune(avatar, 0) ||
			utf8.RuneCountInString(avatar) > 255 {
			return nil, errs.New(
				errs.KindInvalid,
				"avatarURL must be an HTTP(S) URL of at most 255 characters",
			)
		}
		update.AvatarURL = &avatar
	}

	if _, err := u.users.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}

	user, err := u.users.UpdateProfile(ctx, userID, update)
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
