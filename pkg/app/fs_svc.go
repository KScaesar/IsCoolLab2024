package app

import (
	"context"
)

type FileSystemRepository interface {
	CreateFileSystem(ctx context.Context, fs *FileSystem) error
}
