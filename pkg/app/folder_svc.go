package app

import (
	"context"
	"fmt"
)

type FolderService interface {
	CreateFolder(ctx context.Context, username string, params CreateFolderParams) error
	DeleteFolder(ctx context.Context, username string, params DeleteFolderParams) error
	ListFolders(ctx context.Context, username string, params ListFoldersParams) ([]ViewFolder, error)
	RenameFolder(ctx context.Context, username string, params RenameFolderParams) error
}

func NewFolderUseCase(fsRepo FileSystemRepository) *FolderUseCase {
	return &FolderUseCase{
		FsRepo: fsRepo,
	}
}

type FolderUseCase struct {
	FsRepo FileSystemRepository
}

func (uc *FolderUseCase) CreateFolder(ctx context.Context, username string, params CreateFolderParams) error {
	fs, err := uc.FsRepo.GetFileSystemByUsernameV2(ctx, username)
	if err != nil {
		return err
	}

	folder, err := fs.Root.CreateFolder(params)
	if err != nil {
		return err
	}

	err = uc.FsRepo.CreateFolder(ctx, folder)
	if err != nil {
		return err
	}

	return nil
}

func (uc *FolderUseCase) DeleteFolder(ctx context.Context, username string, params DeleteFolderParams) error {
	fs, err := uc.FsRepo.GetFileSystemByUsernameV2(ctx, username)
	if err != nil {
		return err
	}

	folder, err := fs.Root.DeleteFolder(params)
	if err != nil {
		return err
	}

	err = uc.FsRepo.DeleteFolder(ctx, folder)
	if err != nil {
		return err
	}

	return nil
}

func (uc *FolderUseCase) ListFolders(ctx context.Context, username string, params ListFoldersParams) ([]ViewFolder, error) {
	fs, err := uc.FsRepo.GetFileSystemByUsernameV2(ctx, username)
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

func (uc *FolderUseCase) RenameFolder(ctx context.Context, username string, params RenameFolderParams) error {
	fs, err := uc.FsRepo.GetFileSystemByUsernameV2(ctx, username)
	if err != nil {
		return err
	}

	folder, err := fs.Root.RenameFolder(params)
	if err != nil {
		return err
	}

	err = uc.FsRepo.UpdateFolder(ctx, folder)
	if err != nil {
		return err
	}

	return nil
}
