package app

import (
	"context"
)

type FileSystemRepository interface {
	CreateFileSystem(ctx context.Context, fs *FileSystem) error
	GetFileSystemByUsername(ctx context.Context, username string) (*FileSystem, error)

	CreateFolder(ctx context.Context, folder *Folder) error
	DeleteFolder(ctx context.Context, folder *Folder) error
	UpdateFolder(ctx context.Context, folder *Folder) error

	CreateFile(ctx context.Context, file *File) error
	DeleteFile(ctx context.Context, file *File) error
}
