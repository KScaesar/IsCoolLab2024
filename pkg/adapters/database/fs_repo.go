package database

import (
	"context"
	"errors"
	"fmt"

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
		// 只有單層結構, 有效能問題再使用樹狀遞回結構
		Preload("Root", "parent_id = ''"). // 取得 root 目錄本身
		Preload("Root.Folders").           // 取得 root 目錄的 dir
		Preload("Root.Files").             // 取得 root 目錄的 file
		Preload("Root.Folders.Files").     // 取得 dir 的 file
		Where("username = ?", username).
		Take(&fs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Error: The %v %w", username, app.ErrUserNotExists)
		}
		return nil, err
	}
	return &fs, nil
}

func (repo *FileSystemRepository) CreateFolder(ctx context.Context, fs *app.FileSystem) error {
	err := repo.db.WithContext(ctx).Table(FolderTable).
		Create(fs.Root.PersistentFolders).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *FileSystemRepository) DeleteFolder(ctx context.Context, fs *app.FileSystem) error {
	var ids []string
	for _, folder := range fs.Root.PersistentFolders {
		ids = append(ids, folder.Id)
	}

	err := repo.db.WithContext(ctx).Table(FolderTable).
		Delete(fs.Root.PersistentFolders, "id in (?)", ids).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *FileSystemRepository) UpdateFolder(ctx context.Context, fs *app.FileSystem) error {
	return nil
}

func (repo *FileSystemRepository) CreateFile(ctx context.Context, folder *app.Folder) error {
	err := repo.db.WithContext(ctx).Table(FileTable).
		Create(folder.PersistentFiles).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *FileSystemRepository) DeleteFile(ctx context.Context, folder *app.Folder) error {
	var ids []string
	for _, file := range folder.PersistentFiles {
		ids = append(ids, file.Id)
	}

	err := repo.db.WithContext(ctx).Table(FileTable).
		Delete(folder.PersistentFiles, "id in (?)", ids).Error
	if err != nil {
		return err
	}
	return nil
}
