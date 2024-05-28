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
		Where("username = ?", username).
		Preload("Root", "parent_id = ''"). // 取得 root 目錄本身
		Preload("Root.Folders").           // 取得 root 目錄的 dir
		Preload("Root.Folders.Files").     // 取得 dir 的 file
		Preload("Root.Files").             // 取得 root 目錄的 file
		Take(&fs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Error: The %v %w", username, app.ErrUserNotExists)
		}
		return nil, err
	}

	// 只有單層結構, 有 read 效能問題再使用樹狀結構
	//
	// SELECT * FROM `file_systems`
	// WHERE username = "user1"
	// LIMIT 1;
	//
	// SELECT *
	// FROM `folders`
	// WHERE `folders`.`fs_id` = "01HYXCC8AJ35Q5KKVACBGYDF5T" AND parent_id = '';
	//
	// SELECT *
	// FROM `folders`
	// WHERE `folders`.`parent_id` = "01HYXCC8AJ35Q5KKVACDEC38G7";
	//
	// SELECT *
	// FROM `files`
	// WHERE `files`.`folder_id` IN (
	//                              "01HYXCD1CD3VFFRYB9BWV19TM8",
	//                              "01HYXCD1CGB36V08CNRGJQMZHT",
	//                              "01HYXD0GV43XKBZ7Y1YDK7QDBQ"
	//                             );
	//
	// SELECT *
	// FROM `files`
	// WHERE `files`.`folder_id` = "01HYXCC8AJ35Q5KKVACDEC38G7";

	return &fs, nil
}

func (repo *FileSystemRepository) CreateFolder(ctx context.Context, folder *app.Folder) error {
	err := repo.db.WithContext(ctx).Table(FolderTable).
		Create(folder).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *FileSystemRepository) DeleteFolder(ctx context.Context, folder *app.Folder) error {
	err := repo.db.WithContext(ctx).Table(FolderTable).
		Delete(folder, "id = ?", folder.Id).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *FileSystemRepository) UpdateFolder(ctx context.Context, folder *app.Folder) error {
	err := repo.db.WithContext(ctx).Table(FolderTable).
		Omit("Files").
		Where("id = ?", folder.Id).
		Updates(folder.ByUpdate.StdMap()).Error
	if err != nil {
		return err
	}

	for _, file := range folder.Files {
		err := repo.db.WithContext(ctx).Table(FileTable).
			Where("id = ?", file.Id).
			Updates(file.ByUpdate.StdMap()).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *FileSystemRepository) CreateFile(ctx context.Context, file *app.File) error {
	err := repo.db.WithContext(ctx).Table(FileTable).
		Create(file).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *FileSystemRepository) DeleteFile(ctx context.Context, file *app.File) error {
	err := repo.db.WithContext(ctx).Table(FileTable).
		Delete(file, "id = ?", file.Id).Error
	if err != nil {
		return err
	}
	return nil
}
