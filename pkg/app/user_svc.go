package app

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type UserService interface {
	Register(ctx context.Context, username string) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	QueryUserByName(ctx context.Context, username string) (*User, error)
}

func NewUserUseCase(userRepo UserRepository, fsRepo FileSystemRepository) *UserUseCase {
	return &UserUseCase{
		UserRepo: userRepo,
		FsRepo:   fsRepo,
		TimeNow:  time.Now,
	}
}

type UserUseCase struct {
	UserRepo UserRepository
	FsRepo   FileSystemRepository
	TimeNow  func() time.Time
}

func (uc *UserUseCase) Register(ctx context.Context, username string) error {
	user, err := newUser(username)
	if err != nil {
		return err
	}

	_, err = uc.UserRepo.QueryUserByName(ctx, user.Username)
	if err == nil {
		return fmt.Errorf("Error: The %v %w", user.Username, ErrUserExists)
	}

	if !errors.Is(err, ErrUserNotExists) {
		return err
	}

	err = uc.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	createdTime := uc.TimeNow()
	fs := newFileSystem(user.Username, createdTime)

	err = uc.FsRepo.CreateFileSystem(ctx, fs)
	if err != nil {
		return err
	}

	return nil
}
