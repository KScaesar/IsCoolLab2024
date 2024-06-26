package app

import (
	"context"
)

type FileService interface {
	CreateFile(ctx context.Context, username string, params CreateFileParams) error
	DeleteFile(ctx context.Context, username string, params DeleteFileParams) error
	ListFiles(ctx context.Context, username string, params ListFilesParams) ([]ViewFile, error)
}

func NewFileUseCase(fsRepo FileSystemRepository) *FileUseCase {
	return &FileUseCase{
		FsRepo: fsRepo,
	}
}

type FileUseCase struct {
	FsRepo FileSystemRepository
}

func (uc *FileUseCase) CreateFile(ctx context.Context, username string, params CreateFileParams) error {
	fs, err := uc.FsRepo.GetFileSystemByUsernameV3(ctx, username)
	if err != nil {
		return err
	}

	file, err := fs.Root.CreateFile(params)
	if err != nil {
		return err
	}

	err = uc.FsRepo.CreateFile(ctx, file)
	if err != nil {
		return err
	}

	return nil
}

func (uc *FileUseCase) DeleteFile(ctx context.Context, username string, params DeleteFileParams) error {
	fs, err := uc.FsRepo.GetFileSystemByUsernameV3(ctx, username)
	if err != nil {
		return err
	}

	file, err := fs.Root.DeleteFile(params)
	if err != nil {
		return err
	}

	err = uc.FsRepo.DeleteFile(ctx, file)
	if err != nil {
		return err
	}

	return nil
}

func (uc *FileUseCase) ListFiles(ctx context.Context, username string, params ListFilesParams) ([]ViewFile, error) {
	fs, err := uc.FsRepo.GetFileSystemByUsernameV3(ctx, username)
	if err != nil {
		return nil, err
	}

	files, err := fs.Root.ListFiles(params)
	if err != nil {
		return nil, err
	}

	response := make([]ViewFile, len(files))
	for i, file := range files {
		response[i] = ToViewFile(file, username)
	}

	return response, nil
}
