package app

import (
	"time"

	"github.com/KScaesar/IsCoolLab2024/pkg"
)

// folder

type CreateFolderParams struct {
	Foldername  string `validate:"required,foldername"`
	Description string
	CreatedTime time.Time
}

type DeleteFolderParams struct {
	Foldername string `validate:"required,foldername"`
}

type ListFoldersParams struct {
	Sort *FileSystemSortParams
}

type RenameFolderParams struct {
	OldFolderName string `validate:"required,foldername"`
	NewFolderName string `validate:"required,foldername"`
}

// file

type CreateFileParams struct {
	Foldername  string `validate:"required,foldername"`
	Filename    string `validate:"required,filename"`
	Description string
	CreatedTime time.Time
}

type DeleteFileParams struct {
	Foldername string `validate:"required,foldername"`
	Filename   string `validate:"required,filename"`
}

type ListFilesParams struct {
	Foldername string `validate:"required,foldername"`
	Sort       *FileSystemSortParams
}

// sort

var (
	defaultFileSystemSortParams = &FileSystemSortParams{ByName: pkg.SortKind_Asc}
)

type FileSystemSortParams struct {
	ByName    pkg.SortKind `sort:"name"`
	ByCreated pkg.SortKind `sort:"created"`
}

func (p *FileSystemSortParams) Value() *FileSystemSortParams {
	if p.IsZero() {
		return defaultFileSystemSortParams
	}
	return p
}

func (p *FileSystemSortParams) IsZero() bool {
	return p == nil ||
		(p.ByName == "" && p.ByCreated == "")
}

func (p *FileSystemSortParams) Validate() (err error) {
	if p.IsZero() {
		return nil
	}
	return pkg.SortValidate(p)
}

// view model

func ToViewFolder(folder *Folder, username string) ViewFolder {
	return ViewFolder{
		Fodlername:  folder.Name,
		Description: folder.Description,
		CreatedTime: folder.CreatedTime,
		Username:    username,
	}
}

type ViewFolder struct {
	Fodlername  string
	Description string
	CreatedTime time.Time
	Username    string
}

func ToViewFile(file *File, username string) ViewFile {
	return ViewFile{
		Filename:    file.Name,
		Description: file.Description,
		CreatedTime: file.CreatedTime,
		Fodlername:  file.Foldername,
		Username:    username,
	}
}

type ViewFile struct {
	Filename    string
	Description string
	CreatedTime time.Time
	Fodlername  string
	Username    string
}
