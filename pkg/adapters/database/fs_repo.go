package database

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (repo *FileSystemRepository) GetFileSystemByUsernameV2(ctx context.Context, username string) (*app.FileSystem, error) {
	type Mapper struct {
		FsId              string     `gorm:"column:fs_id"`
		FolderId          string     `gorm:"column:folder_id"`
		ParentId          string     `gorm:"column:parent_id"`
		FolderName        string     `gorm:"column:folder_name"`
		FolderDescription string     `gorm:"column:folder_description"`
		FolderCreatedTime time.Time  `gorm:"column:folder_created_time"`
		FileId            *string    `gorm:"column:file_id"`
		FileName          *string    `gorm:"column:file_name"`
		FileDescription   *string    `gorm:"column:file_description"`
		FileCreatedTime   *time.Time `gorm:"column:file_created_time"`
	}
	var results []Mapper

	err := repo.db.WithContext(ctx).Raw(`
SELECT fs.id               AS fs_id,
       folder.id           AS folder_id,
       folder.parent_id    AS parent_id,
       folder.name         AS folder_name,
       folder.description  AS folder_description,
       folder.created_time AS folder_created_time,
       file.id             AS file_id,
       file.name           AS file_name,
       file.description    AS file_description,
       file.created_time   AS file_created_time
FROM file_systems fs
      LEFT JOIN
     folders folder ON folder.fs_id = fs.id
      LEFT JOIN
     files file ON file.folder_id = folder.id
WHERE fs.username = ?

UNION

SELECT fs.id                    AS fs_id,
       root_folder.id           AS folder_id,
       root_folder.parent_id    AS parent_id,
       root_folder.name         AS folder_name,
       root_folder.description  AS folder_description,
       root_folder.created_time AS folder_created_time,
       root_file.id             AS file_id,
       root_file.name           AS file_name,
       root_file.description    AS file_description,
       root_file.created_time   AS file_created_time
FROM file_systems fs
      LEFT JOIN
     folders root_folder ON root_folder.fs_id = fs.id AND root_folder.parent_id = ''
      LEFT JOIN
     files root_file ON root_file.folder_id = root_folder.id
WHERE fs.username = ?;`, username, username).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	folders := make(map[string]*app.Folder, len(results))

	for _, r := range results {
		folder, ok1 := folders[r.FolderId]
		if !ok1 {
			folder = &app.Folder{
				Id:             r.FolderId,
				ParentFolderId: r.ParentId,
				FsId:           r.FsId,
				Name:           r.FolderName,
				Description:    r.FolderDescription,
				CreatedTime:    r.FolderCreatedTime,
			}
			folders[r.FolderId] = folder
		}
		if r.FileId != nil {
			folder.Files = append(folder.Files, &app.File{
				Id:          *r.FileId,
				Name:        *r.FileName,
				FolderId:    r.FolderId,
				Foldername:  r.FolderName,
				Description: *r.FileDescription,
				CreatedTime: *r.FileCreatedTime,
			})
		}
	}

	var root *app.Folder
	for _, folder := range folders {
		if folder.ParentFolderId == "" {
			root = folder
			continue
		}
		parentFolder := folders[folder.ParentFolderId]
		parentFolder.Folders = append(parentFolder.Folders, folder)
	}

	if root == nil {
		return nil, fmt.Errorf("Error: The %v %w", username, app.ErrUserNotExists)
	}

	fs := &app.FileSystem{
		Id:       root.FsId,
		Username: username,
		Root:     *root,
	}
	return fs, nil
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
