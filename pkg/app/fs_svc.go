package app

import (
	"context"
)

type FileSystemRepository interface {
	CreateFileSystem(ctx context.Context, fs *FileSystem) error
	GetFileSystemByUsername(ctx context.Context, username string) (*FileSystem, error)
	CreateFolder(ctx context.Context, fs *FileSystem) error
}
