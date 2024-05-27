package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/KScaesar/IsCoolLab2024/pkg/app"
)

const (
	UserTable = "users"
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

type UserRepository struct {
	db *gorm.DB
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *app.User) error {
	err := repo.db.WithContext(ctx).Table(UserTable).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) QueryUserByName(ctx context.Context, username string) (*app.User, error) {
	var user app.User
	err := repo.db.WithContext(ctx).Table(UserTable).
		Where("username = ?", username).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app.ErrUserNotExists
		}
		return nil, err
	}
	return &user, nil
}
