package app

import (
	"context"
	"fmt"
	"time"
)

type FolderService interface {
	CreateFolder(ctx context.Context, username string, params CreateFolderParams) error
	DeleteFolder(ctx context.Context, username string, params DeleteFolderParams) error
	ListFolders(ctx context.Context, username string, params ListFoldersParams) ([]ViewFolder, error)
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
		return err
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
		return err
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

func (uc *FolderUseCase) ListFolders(ctx context.Context, username string, params ListFoldersParams) ([]ViewFolder, error) {
	fs, err := uc.FsRepo.GetFileSystemByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if len(fs.Root.Folders) == 0 {
		return nil, fmt.Errorf("Warning: The %v %w", username, ErrListFolderEmpty)
	}

	folders := fs.Root.ListFolders(params)
	response := make([]ViewFolder, len(folders))
	for i, folder := range folders {
		response[i] = ToViewFolder(folder, username)
	}

	return response, nil
}
