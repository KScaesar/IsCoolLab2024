package app

import (
	"time"

	"github.com/KScaesar/IsCoolLab2024/pkg"
)

// folder

type CreateFolderParams struct {
	Foldername  string `validate:"required,foldername"`
	Description string
	createdTime time.Time
}

func (p *CreateFolderParams) CreatedTime() time.Time {
	return p.createdTime
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
	folderId    string
	Foldername  string `validate:"required,foldername"`
	Filename    string `validate:"required,filename"`
	Description string
	CreatedTime time.Time
}

func (p *CreateFileParams) FolderId() string {
	return p.folderId
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

func ConvertFolderView(folder *Folder, username string) *FolderView {
	return &FolderView{
		Fodlername:  folder.Name,
		Description: folder.Description,
		CreatedTime: folder.CreatedTime,
		Username:    username,
	}
}

type FolderView struct {
	Fodlername  string
	Description string
	CreatedTime time.Time
	Username    string
}

func ConvertFileView(file *File, username string) *FileView {
	return &FileView{
		Filename:    file.Name,
		Description: file.Description,
		CreatedTime: file.CreatedTime,
		Fodlername:  file.Foldername,
		Username:    username,
	}
}

type FileView struct {
	Filename    string
	Description string
	CreatedTime time.Time
	Fodlername  string
	Username    string
}
