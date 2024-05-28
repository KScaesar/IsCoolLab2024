package app

import (
	"context"
	"fmt"
	"time"
)

type FolderService interface {
	CreateFolder(ctx context.Context, username string, params CreateFolderParams) error
	DeleteFolder(ctx context.Context, username string, params DeleteFolderParams) error
}

func NewFolderUseCase(fsRepo FileSystemRepository) *FolderUseCase {
	return &FolderUseCase{
		FsRepo:  fsRepo,
		TimeNow: time.Now,
	}
}

type FolderUseCase struct {
	FsRepo  FileSystemRepository
	TimeNow func() time.Time
}

func (uc *FolderUseCase) CreateFolder(ctx context.Context, username string, params CreateFolderParams) error {
	fs, err := uc.FsRepo.GetFileSystemByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("Error: The %v %w", username, err)
	}

	params.createdTime = uc.TimeNow()
	err = fs.Root.CreateFolder(params)
	if err != nil {
		return err
	}

	err = uc.FsRepo.CreateFolder(ctx, fs)
	if err != nil {
		return err
	}

	return nil
}

func (uc *FolderUseCase) DeleteFolder(ctx context.Context, username string, params DeleteFolderParams) error {
	fs, err := uc.FsRepo.GetFileSystemByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("Error: The %v %w", username, err)
	}

	err = fs.Root.DeleteFolder(params)
	if err != nil {
		return err
	}

	err = uc.FsRepo.DeleteFolder(ctx, fs)
	if err != nil {
		return err
	}

	return nil
}
