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
JOIN folders folder ON folder.fs_id = fs.id AND fs.username = ? 
LEFT JOIN files file ON file.folder_id = folder.id;`, username).
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

func (repo *FileSystemRepository) GetFileSystemByUsernameV3(ctx context.Context, username string) (*app.FileSystem, error) {
	type Mapper struct {
		Id          string    `gorm:"column:id"`
		ParentID    string    `gorm:"column:parent_id"`
		FsId        string    `gorm:"column:fs_id"`
		Name        string    `gorm:"column:name"`
		Description string    `gorm:"column:description"`
		CreatedTime time.Time `gorm:"column:created_time"`
		Kind        string    `gorm:"column:kind"`
		Level       int       `gorm:"column:level"`
	}
	var rows []Mapper

	err := repo.db.WithContext(ctx).Raw(`
WITH RECURSIVE hierarchy AS (
 -- Anchor member: select the root nodes
 SELECT
  d.id AS id,
  parent_id,
  fs_id,
  d.name,
  d.description,
  d.created_time,
  'd' AS kind,
  0 AS level
 FROM file_systems fs
 JOIN folders d ON d.fs_id = fs.id AND fs.username = ?
 WHERE parent_id = ''

 UNION ALL
 -- Recursive member: select children of the current level
 SELECT
  f.id,
  f.folder_id AS parent_id,
  f.fs_id,
  f.name,
  f.description,
  f.created_time,
  'f' AS kind,
  h.level + 1 AS level
 FROM files f
 JOIN hierarchy h ON f.folder_id = h.id
 UNION ALL
 SELECT
  d.id,
  d.parent_id,
  d.fs_id,
  d.name,
  d.description,
  d.created_time,
  'd' AS kind,
  h.level + 1 AS level
 FROM folders d
 JOIN hierarchy h ON d.parent_id = h.id
)
SELECT * FROM hierarchy;`, username).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	var root *app.Folder
	if len(rows) > 0 && rows[0].Kind == "d" {
		root = &app.Folder{
			Id:             rows[0].Id,
			ParentFolderId: rows[0].ParentID,
			FsId:           rows[0].FsId,
			Name:           rows[0].Name,
			Description:    rows[0].Description,
			CreatedTime:    rows[0].CreatedTime,
		}
	}

	if root == nil {
		return nil, fmt.Errorf("Error: The %v %w", username, app.ErrUserNotExists)
	}

	folders := make(map[string]*app.Folder)
	folders[root.Id] = root
	level := 1
	for i := 1; i < len(rows); {
		for i < len(rows) && level == rows[i].Level {
			if rows[i].Kind == "d" {
				folder := &app.Folder{
					Id:             rows[i].Id,
					ParentFolderId: rows[i].ParentID,
					FsId:           rows[i].FsId,
					Name:           rows[i].Name,
					Description:    rows[i].Description,
					CreatedTime:    rows[i].CreatedTime,
				}
				folders[folder.Id] = folder
				parent := folders[folder.ParentFolderId]
				parent.Folders = append(parent.Folders, folder)
			}
			if rows[i].Kind == "f" {
				parent := folders[rows[i].ParentID]
				file := &app.File{
					Id:          rows[i].Id,
					FolderId:    rows[i].ParentID,
					FsId:        rows[i].FsId,
					Name:        rows[i].Name,
					Foldername:  parent.Name,
					Description: rows[i].Description,
					CreatedTime: rows[i].CreatedTime,
				}
				parent.Files = append(parent.Files, file)
			}
			i++
		}
		level++
	}

	fs := &app.FileSystem{
		Id:       root.FsId,
		Username: username,
		Root:     *root,
	}

	// [
	//  {
	//    "id": "01HYXCC8AJ35Q5KKVACDEC38G7",
	//    "parent_id": "",
	//    "fs_id": "01HYXCC8AJ35Q5KKVACBGYDF5T",
	//    "name": "/",
	//    "description": "",
	//    "created_time": "2024-05-27 23:00:00+08:00",
	//    "kind": "d",
	//    "level": 0
	//  },
	//  {
	//    "id": "01HYXCD1CD3VFFRYB9BWV19TM8",
	//    "parent_id": "01HYXCC8AJ35Q5KKVACDEC38G7",
	//    "fs_id": "01HYXCC8AJ35Q5KKVACBGYDF5T",
	//    "name": "folder isCool",
	//    "description": "",
	//    "created_time": "2024-05-27 23:00:03+08:00",
	//    "kind": "d",
	//    "level": 1
	//  },
	//  {
	//    "id": "01HYXCD1CGB36V08CNRGJQMZHT",
	//    "parent_id": "01HYXCC8AJ35Q5KKVACDEC38G7",
	//    "fs_id": "01HYXCC8AJ35Q5KKVACBGYDF5T",
	//    "name": "folder2",
	//    "description": "qa-folder",
	//    "created_time": "2024-05-27 23:00:01+08:00",
	//    "kind": "d",
	//    "level": 1
	//  },
	//  {
	//    "id": "01HYXD0GV43XKBZ7Y1YDK7QDBQ",
	//    "parent_id": "01HYXCC8AJ35Q5KKVACDEC38G7",
	//    "fs_id": "01HYXCC8AJ35Q5KKVACBGYDF5T",
	//    "name": "folder3",
	//    "description": "",
	//    "created_time": "2024-05-27 23:00:02+08:00",
	//    "kind": "d",
	//    "level": 1
	//  },
	//  {
	//    "id": "01HYYMFNZSFQ2FWPN1DYFTPADH",
	//    "parent_id": "01HYXCD1CD3VFFRYB9BWV19TM8",
	//    "fs_id": "01HYXCC8AJ35Q5KKVACBGYDF5T",
	//    "name": "file1",
	//    "description": "",
	//    "created_time": "2024-05-27 23:00:03+08:00",
	//    "kind": "f",
	//    "level": 2
	//  },
	//  {
	//    "id": "01HYYMMTX8F4D2BESDCAD2YXS5",
	//    "parent_id": "01HYXCD1CD3VFFRYB9BWV19TM8",
	//    "fs_id": "01HYXCC8AJ35Q5KKVACBGYDF5T",
	//    "name": "file2",
	//    "description": "qa-file",
	//    "created_time": "2024-05-27 23:00:01+08:00",
	//    "kind": "f",
	//    "level": 2
	//  },
	//  {
	//    "id": "01HYYMN2H854NWJJ32HJRCQKC0",
	//    "parent_id": "01HYXCD1CD3VFFRYB9BWV19TM8",
	//    "fs_id": "01HYXCC8AJ35Q5KKVACBGYDF5T",
	//    "name": "file3",
	//    "description": "",
	//    "created_time": "2024-05-27 23:00:02+08:00",
	//    "kind": "f",
	//    "level": 2
	//  }
	// ]
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
	db := repo.db.WithContext(ctx)

	err := db.Table(FolderTable).
		Delete(folder, "id = ?", folder.Id).Error
	if err != nil {
		return err
	}

	err = db.Table(FileTable).
		Delete(&app.File{}, "folder_id = ?", folder.Id).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *FileSystemRepository) UpdateFolder(ctx context.Context, folder *app.Folder) error {
	db := repo.db.WithContext(ctx)

	err := db.Table(FolderTable).
		Omit("Files").
		Where("id = ?", folder.Id).
		Updates(folder.ByUpdate.StdMap()).Error
	if err != nil {
		return err
	}

	if len(folder.Files) == 0 {
		return nil
	}

	file := folder.Files[0]
	err = db.Table(FileTable).
		Where("folder_id = ?", folder.Id).
		Updates(file.ByUpdate.StdMap()).Error
	if err != nil {
		return err
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
