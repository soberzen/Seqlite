package repository

import (
	"context"
	"errors"
	"qingmo-server/internal/errs"
	"qingmo-server/internal/model"

	"gorm.io/gorm"
)

type UserPostgresRepository struct {
	db *gorm.DB
}

func NewUserPostgresRepository(db *gorm.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

func (u *UserPostgresRepository) ExistsByEmail(
	ctx context.Context,
	email string,
) (bool, error) {
	var count int64
	err := u.db.WithContext(ctx).
		Model(&model.Users{}).
		Where("email = ?", email).
		Count(&count).
		Error
	if err != nil {
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) {
			return false, err
		}
		return false, errs.Wrap(
			errs.KindInternal,
			"check user email exists",
			err,
		)
	}
	return count > 0, nil
}

func (u *UserPostgresRepository) Create(
	ctx context.Context,
	user *model.Users,
) (*model.Users, error) {
	if err := u.db.WithContext(ctx).
		Create(user).
		Error; err != nil {
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, errs.Wrap(
			errs.KindInternal,
			"create user failed",
			err,
		)
	}
	return user, nil
}

func (u *UserPostgresRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*model.Users, error) {
	var user model.Users
	if err := u.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error; err != nil {
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(errs.KindNotFound, "user not found", err)
		}
		return nil, errs.Wrap(
			errs.KindInternal,
			"get user by email failed",
			err,
		)
	}
	return &user, nil
}
