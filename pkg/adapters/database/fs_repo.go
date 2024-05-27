package database

import (
	"context"

	"gorm.io/gorm"

	"github.com/KScaesar/IsCoolLab2024/pkg/app"
)

const (
	FileSystemTable = "file_systems"
	FileNodeTable   = "file_nodes"
)

func NewFileSystemRepository(db *gorm.DB) *FileSystemRepository {
	return &FileSystemRepository{db: db}
}

type FileSystemRepository struct {
	db *gorm.DB
}

func (repo *FileSystemRepository) CreateFileSystem(ctx context.Context, fs *app.FileSystem) error {
	err := repo.db.WithContext(ctx).Table(FileSystemTable).Create(fs).Error
	if err != nil {
		return err
	}
	return nil
}
