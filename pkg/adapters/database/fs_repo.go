package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/KScaesar/IsCoolLab2024/pkg/app"
)

const (
	FileSystemTable = "file_systems"
	FolderTable     = "folders"
	FileTable       = "files"
)

func NewFileSystemRepository(db *gorm.DB) *FileSystemRepository {
	return &FileSystemRepository{db: db}
}

type FileSystemRepository struct {
	db *gorm.DB
}

func (repo *FileSystemRepository) CreateFileSystem(ctx context.Context, fs *app.FileSystem) error {
	err := repo.db.WithContext(ctx).Table(FileSystemTable).
		Create(fs).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *FileSystemRepository) GetFileSystemByUsername(ctx context.Context, username string) (*app.FileSystem, error) {
	// https://gorm.io/zh_CN/docs/preload.html#%E9%A2%84%E5%8A%A0%E8%BD%BD%E5%85%A8%E9%83%A8
	var fs app.FileSystem
	err := repo.db.WithContext(ctx).Table(FileSystemTable).
		Preload("Root", "name = ''").
		Preload("Root.Files").
		Preload("Root.Folders").
		Where("username = ?", username).
		Take(&fs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app.ErrUserNotExists
		}
		return nil, err
	}
	return &fs, nil
}

func (repo *FileSystemRepository) CreateFolder(ctx context.Context, fs *app.FileSystem) error {
	err := repo.db.WithContext(ctx).Table(FolderTable).
		Create(fs.Root.HelpedFolders).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *FileSystemRepository) DeleteFolder(ctx context.Context, fs *app.FileSystem) error {
	var ids []string
	for _, folder := range fs.Root.HelpedFolders {
		ids = append(ids, folder.Id)
	}

	err := repo.db.WithContext(ctx).Table(FolderTable).
		Delete(fs.Root.HelpedFolders, "id in (?)", ids).Error
	if err != nil {
		return err
	}
	return nil
}
