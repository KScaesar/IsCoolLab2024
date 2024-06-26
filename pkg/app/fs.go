package app

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/KScaesar/IsCoolLab2024/pkg"
)

func newFileSystem(username string, createdTime time.Time) *FileSystem {
	fsId := pkg.NewUlid()
	return &FileSystem{
		Id:       fsId,
		Username: username,
		Root:     newRootFolder(fsId, createdTime),
	}
}

type FileSystem struct {
	Id       string `gorm:"column:id;type:char(26);not null;primaryKey"`
	Username string `gorm:"column:username;type:varchar(64);not null;uniqueIndex"`
	Root     Folder `gorm:"foreignKey:fs_id"`
}

func newRootFolder(fsId string, createdTime time.Time) Folder {
	return Folder{
		Id:          pkg.NewUlid(),
		FsId:        fsId,
		Name:        "/",
		CreatedTime: createdTime,
	}
}

func newFolder(parentId, fsId string, params CreateFolderParams) (*Folder, error) {
	err := validateFoldername(params.Foldername)
	if err != nil {
		return nil, err
	}

	return &Folder{
		Id:             pkg.NewUlid(),
		ParentFolderId: parentId,
		FsId:           fsId,
		Name:           params.Foldername,
		Description:    params.Description,
		CreatedTime:    params.CreatedTime,
	}, nil
}

type Folder struct {
	Id             string    `gorm:"column:id;type:char(26);not null;primaryKey"`
	ParentFolderId string    `gorm:"column:parent_id;type:char(26);not null;index"`
	FsId           string    `gorm:"column:fs_id;type:char(26);not null;index"`
	Name           string    `gorm:"column:name;type:varchar(256);not null"`
	Description    string    `gorm:"column:description;type:varchar(1024);not null"`
	CreatedTime    time.Time `gorm:"column:created_time;not null"`
	Files          []*File   `gorm:"foreignKey:folder_id"`
	Folders        []*Folder `gorm:"foreignKey:parent_id"`

	ByUpdate pkg.MapData `gorm:"-"`
}

func (dir *Folder) CreateFolder(params CreateFolderParams) (*Folder, error) {
	_, err := dir.findFolder(params.Foldername)
	if err == nil {
		return nil, fmt.Errorf("Error: The %v %w", params.Foldername, ErrFolderExists)
	}

	if !errors.Is(err, ErrFolderNotExists) {
		return nil, err
	}

	folder, err := newFolder(dir.Id, dir.FsId, params)
	if err != nil {
		return nil, err
	}

	dir.Folders = append(dir.Folders, folder)
	return folder, nil
}

func (dir *Folder) DeleteFolder(params DeleteFolderParams) (*Folder, error) {
	targetFolder, err := dir.findFolder(params.Foldername)
	if err != nil {
		return nil, err
	}

	for i, folder := range dir.Folders {
		if folder == targetFolder {
			dir.Folders = append(dir.Folders[:i], dir.Folders[i+1:]...)
			return folder, nil
		}
	}
	return nil, nil
}

func (dir *Folder) findFolder(foldername string) (*Folder, error) {
	if dir.Name == foldername {
		return dir, nil
	}
	for _, folder := range dir.Folders {
		if strings.EqualFold(folder.Name, foldername) {
			return folder, nil
		}
	}
	return nil, fmt.Errorf("Error: The %v %w", foldername, ErrFolderNotExists)
}

func (dir *Folder) ListFolders(params ListFoldersParams) []*Folder {
	pkg.SortTraversalParams(params.Sort.Value(), func(key string, value pkg.SortKind) {
		sort.Slice(dir.Folders, func(i, j int) bool {
			switch key {
			case "name":
				if value == pkg.SortKind_Asc {
					return dir.Folders[i].Name < dir.Folders[j].Name
				}
				if value == pkg.SortKind_Desc {
					return dir.Folders[i].Name > dir.Folders[j].Name
				}

			case "created":
				if value == pkg.SortKind_Asc {
					return dir.Folders[i].CreatedTime.Second() < dir.Folders[j].CreatedTime.Second()
				}
				if value == pkg.SortKind_Desc {
					return dir.Folders[i].CreatedTime.Second() > dir.Folders[j].CreatedTime.Second()
				}
			}

			return false
		})
	})

	return dir.Folders
}

func (dir *Folder) RenameFolder(params RenameFolderParams) (*Folder, error) {
	err := validateFoldername(params.NewFolderName)
	if err != nil {
		return nil, fmt.Errorf("Error: The %v %w", params.NewFolderName, err)
	}

	_, err = dir.findFolder(params.NewFolderName)
	if err == nil {
		return nil, fmt.Errorf("Error: The %v %w", params.NewFolderName, ErrFolderExists)
	}

	if !errors.Is(err, ErrFolderNotExists) {
		return nil, err
	}

	folder, err := dir.findFolder(params.OldFolderName)
	if err != nil {
		return nil, err
	}

	folder.Name = params.NewFolderName
	folder.ByUpdate.MustOk().Set("name", folder.Name)
	for _, file := range folder.Files {
		file.Foldername = params.NewFolderName
		file.ByUpdate.MustOk().Set("foldername", file.Foldername)
	}
	return folder, nil
}

func (dir *Folder) CreateFile(params CreateFileParams) (*File, error) {
	folder, err := dir.findFolder(params.Foldername)
	if err != nil {
		return nil, err
	}

	for _, file := range folder.Files {
		if file.Name == params.Filename {
			return nil, fmt.Errorf("Error: The %v %w", params.Filename, ErrFileExists)
		}
	}

	file, err := newFile(folder.Id, folder.FsId, params)
	if err != nil {
		return nil, err
	}

	folder.Files = append(folder.Files, file)
	return file, nil
}

func (dir *Folder) DeleteFile(params DeleteFileParams) (*File, error) {
	folder, err := dir.findFolder(params.Foldername)
	if err != nil {
		return nil, err
	}

	for i, file := range folder.Files {
		if file.Name == params.Filename {
			folder.Files = append(folder.Files[:i], folder.Files[i+1:]...)
			return file, nil
		}
	}

	return nil, fmt.Errorf("Error: The %v %w", params.Filename, ErrFileNotExists)
}

func (dir *Folder) ListFiles(params ListFilesParams) ([]*File, error) {
	folder, err := dir.findFolder(params.Foldername)
	if err != nil {
		return nil, err
	}

	if len(folder.Files) == 0 {
		return nil, ErrListFileEmpty
	}

	pkg.SortTraversalParams(params.Sort.Value(), func(key string, value pkg.SortKind) {
		sort.Slice(folder.Files, func(i, j int) bool {
			switch key {
			case "name":
				if value == pkg.SortKind_Asc {
					return folder.Files[i].Name < folder.Files[j].Name
				}
				if value == pkg.SortKind_Desc {
					return folder.Files[i].Name > folder.Files[j].Name
				}

			case "created":
				if value == pkg.SortKind_Asc {
					return folder.Files[i].CreatedTime.Second() < folder.Files[j].CreatedTime.Second()
				}
				if value == pkg.SortKind_Desc {
					return folder.Files[i].CreatedTime.Second() > folder.Files[j].CreatedTime.Second()
				}
			}

			return false
		})
	})

	return folder.Files, nil
}

func newFile(folderId, fsId string, params CreateFileParams) (*File, error) {
	err := validateFilename(params.Filename)
	if err != nil {
		return nil, err
	}

	return &File{
		Id:          pkg.NewUlid(),
		FolderId:    folderId,
		FsId:        fsId,
		Name:        params.Filename,
		Foldername:  params.Foldername,
		Description: params.Description,
		CreatedTime: params.CreatedTime,
		ByUpdate:    nil,
	}, nil
}

type File struct {
	Id          string    `gorm:"column:id;type:char(26);not null;primaryKey"`
	FolderId    string    `gorm:"column:folder_id;type:char(26);not null;index"`
	FsId        string    `gorm:"column:fs_id;type:char(26);not null;index"`
	Name        string    `gorm:"column:name;type:varchar(256);not null"`
	Foldername  string    `gorm:"column:foldername;type:varchar(256);not null"`
	Description string    `gorm:"column:description;type:varchar(1024);not null"`
	CreatedTime time.Time `gorm:"column:created_time;not null"`

	ByUpdate pkg.MapData `gorm:"-"`
}

// validate

func validateFoldername(foldername string) error {
	if len(foldername) > 256 {
		return fmt.Errorf("Error: The %v %w", foldername, ErrInvalidParams)
	}
	for _, char := range foldername {
		if !(unicode.IsLetter(char) || unicode.IsNumber(char) || char == '_' || char == '-' || char == '/' || char == ' ') {
			return fmt.Errorf("Error: The %v %w", foldername, ErrInvalidParams)
		}
	}
	return nil
}

func validateFilename(filename string) error {
	if len(filename) > 256 {
		return fmt.Errorf("Error: The %v %w", filename, ErrInvalidParams)
	}
	for _, char := range filename {
		if !(unicode.IsLetter(char) || unicode.IsNumber(char) || char == '_' || char == '-' || char == '.' || char == ' ') {
			return fmt.Errorf("Error: The %v %w", filename, ErrInvalidParams)
		}
	}
	return nil
}
