package app

import (
	"context"
)

type FileSystemRepository interface {
	CreateFileSystem(ctx context.Context, fs *FileSystem) error
	GetFileSystemByUsername(ctx context.Context, username string) (*FileSystem, error)

	CreateFolder(ctx context.Context, fs *FileSystem) error
	DeleteFolder(ctx context.Context, fs *FileSystem) error
	UpdateFolder(ctx context.Context, fs *FileSystem) error

	CreateFile(ctx context.Context, folder *Folder) error
	DeleteFile(ctx context.Context, folder *Folder) error
}
