package app

import (
	"errors"
	"testing"
	"time"

	"github.com/KScaesar/IsCoolLab2024/pkg"
)

func testFileSystem() *FileSystem {
	createdTime := pkg.MockTimeNow("2024-05-26T12:00:00+08:00")()
	fs := newFileSystem("caesar", createdTime)

	// folder
	err := fs.Root.CreateFolder(CreateFolderParams{
		FsId:        fs.Id,
		Foldername:  "/home",
		CreatedTime: createdTime,
	})
	if err != nil {
		panic(err)
	}
	fs.Root.CreateFolder(CreateFolderParams{
		FsId:        fs.Id,
		Foldername:  "/etc",
		CreatedTime: createdTime.Add(2 * time.Second),
	})
	fs.Root.CreateFolder(CreateFolderParams{
		FsId:        fs.Id,
		Foldername:  "/tmp",
		CreatedTime: createdTime.Add(time.Second),
	})

	// file
	err = fs.Root.CreateFile(CreateFileParams{
		FsId:        fs.Id,
		Foldername:  "/home",
		Filename:    "dev.conf",
		CreatedTime: createdTime,
	})
	if err != nil {
		panic(err)
	}
	fs.Root.CreateFile(CreateFileParams{
		FsId:        fs.Id,
		Foldername:  "/home",
		Filename:    "prod.conf",
		CreatedTime: createdTime.Add(2 * time.Second),
	})
	fs.Root.CreateFile(CreateFileParams{
		FsId:        fs.Id,
		Foldername:  "/home",
		Filename:    "qa.conf",
		CreatedTime: createdTime.Add(time.Second),
	})

	return fs
}

func TestFolder_CreateFolder(t *testing.T) {
	fs := testFileSystem()

	tests := []struct {
		name    string
		params  CreateFolderParams
		wantErr error
		assert  func(t *testing.T)
	}{
		{
			name: "success",
			params: CreateFolderParams{
				Foldername: "keys",
			},
			wantErr: nil,
			assert: func(t *testing.T) {
				n := len(fs.Root.Folders)
				want := 4
				if n != want {
					t.Errorf("CreateFolder() len=%v, want=%v", n, want)
				}
			},
		},
		{
			name: "The [foldername] has already existed.",
			params: CreateFolderParams{
				Foldername: "/home",
			},
			wantErr: ErrFolderExists,
		},
		{
			name: "The [foldername] contain invalid chars.",
			params: CreateFolderParams{
				Foldername: "/home#1",
			},
			wantErr: ErrInvalidParams,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := fs.Root.CreateFolder(tt.params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateFolder() error=%v, want=%v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert(t)
			}
		})
	}
}

func TestFolder_DeleteFolder(t *testing.T) {
	fs := testFileSystem()

	tests := []struct {
		name    string
		params  DeleteFolderParams
		wantErr error
		assert  func(t *testing.T)
	}{
		{
			name: "success",
			params: DeleteFolderParams{
				Foldername: "/home",
			},
			wantErr: nil,
			assert: func(t *testing.T) {
				n := len(fs.Root.Folders)
				want := 2
				if n != want {
					t.Errorf("DeleteFolder() len=%v, want=%v", n, want)
					return
				}
			},
		},
		{
			name: "The [foldername] doesn't exist.",
			params: DeleteFolderParams{
				Foldername: "home3",
			},
			wantErr: ErrFolderNotExists,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := fs.Root.DeleteFolder(tt.params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteFolder() error=%v, want=%v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert(t)
			}
		})
	}
}

func TestFolder_ListFolders(t *testing.T) {
	fs := testFileSystem()

	tests := []struct {
		name    string
		params  ListFoldersParams
		wantErr error
		assert  func(t *testing.T, folders []*Folder)
	}{
		{
			name: "by default",
			params: ListFoldersParams{
				Sort: nil,
			},
			wantErr: nil,
			assert: func(t *testing.T, folders []*Folder) {
				want := []string{"/etc", "/home", "/tmp"}
				for i, folder := range folders {
					if folder.Name != want[i] {
						t.Errorf("ListFolders() folder=%v, want=%v", folder.Name, want[i])
					}
				}
			},
		},
		{
			name: "by name",
			params: ListFoldersParams{
				Sort: &FileSystemSortParams{ByName: pkg.SortKind_Desc},
			},
			wantErr: nil,
			assert: func(t *testing.T, folders []*Folder) {
				want := []string{"/tmp", "/home", "/etc"}
				for i, folder := range folders {
					if folder.Name != want[i] {
						t.Errorf("ListFolders() folder=%v, want=%v", folder.Name, want[i])
					}
				}
			},
		},
		{
			name: "by createdTime",
			params: ListFoldersParams{
				Sort: &FileSystemSortParams{ByCreated: pkg.SortKind_Desc},
			},
			wantErr: nil,
			assert: func(t *testing.T, folders []*Folder) {
				want := []string{"/etc", "/tmp", "/home"}
				for i, folder := range folders {
					if folder.Name != want[i] {
						t.Errorf("ListFolders() folder=%v, want=%v", folder.Name, want[i])
					}
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			folders, err := fs.Root.ListFolders(tt.params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ListFolders() error=%v, want=%v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert(t, folders)
			}
		})
	}
}

func TestFolder_RenameFolder(t *testing.T) {
	fs := testFileSystem()

	tests := []struct {
		name    string
		params  RenameFolderParams
		wantErr error
		assert  func(t *testing.T, folder *Folder)
	}{
		{
			name: "success",
			params: RenameFolderParams{
				OldFolderName: "/home",
				NewFolderName: "home2",
			},
			wantErr: nil,
			assert: func(t *testing.T, folder *Folder) {
				name := folder.Name
				want := "home2"
				if name != want {
					t.Errorf("for folder: name=%v, want=%v", name, want)
					return
				}

				for _, file := range folder.Files {
					folderame := file.Foldername
					if name != want {
						t.Errorf("for file: name=%v, want=%v", folderame, want)
						return
					}
				}
			},
		},
		{
			name: "OldName doesn't exist",
			params: RenameFolderParams{
				OldFolderName: "home1",
				NewFolderName: "home3",
			},
			wantErr: ErrFolderNotExists,
		},
		{
			name: "NewName exist",
			params: RenameFolderParams{
				OldFolderName: "/home",
				NewFolderName: "/etc",
			},
			wantErr: ErrFolderExists,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := fs.Root.RenameFolder(tt.params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("RenameFolder() error=%v, want=%v", err, tt.wantErr)
			}

			folder, _ := fs.Root.findFolder(tt.params.NewFolderName)
			if tt.assert != nil {
				tt.assert(t, folder)
			}
		})
	}
}

func TestFolder_CreateFile(t *testing.T) {
	fs := testFileSystem()

	tests := []struct {
		name    string
		params  CreateFileParams
		wantErr error
		assert  func(t *testing.T, folder *Folder)
	}{
		{
			name: "success",
			params: CreateFileParams{
				Foldername: "/etc",
				Filename:   "nginx.conf",
			},
			wantErr: nil,
			assert: func(t *testing.T, folder *Folder) {
				n := len(folder.Files)
				want := 1
				if n != want {
					t.Errorf("CreateFile() len=%v, want=%v", n, want)
				}
			},
		},
		{
			name: "The [foldername] doesn't exist.",
			params: CreateFileParams{
				Foldername: "app",
				Filename:   "nginx.conf",
			},
			wantErr: ErrFolderNotExists,
		},
		{
			name: "The [filename] has already existed.",
			params: CreateFileParams{
				Foldername: "/home",
				Filename:   "dev.conf",
			},
			wantErr: ErrFileExists,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := fs.Root.CreateFile(tt.params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateFile() error=%v, want=%v", err, tt.wantErr)
			}

			folder, _ := fs.Root.findFolder(tt.params.Foldername)
			if tt.assert != nil {
				tt.assert(t, folder)
			}
		})
	}
}

func TestFolder_DeleteFile(t *testing.T) {
	fs := testFileSystem()

	tests := []struct {
		name    string
		params  DeleteFileParams
		wantErr error
		assert  func(t *testing.T, folder *Folder)
	}{
		{
			name: "success",
			params: DeleteFileParams{
				Foldername: "/home",
				Filename:   "qa.conf",
			},
			wantErr: nil,
			assert: func(t *testing.T, folder *Folder) {
				n := len(folder.Files)
				want := 2
				if n != want {
					t.Errorf("DeleteFile() len=%v, want=%v", n, want)
				}
			},
		},
		{
			name: "The [foldername] doesn't exist.",
			params: DeleteFileParams{
				Foldername: "app",
				Filename:   "qa.conf",
			},
			wantErr: ErrFolderNotExists,
		},
		{
			name: "The [filename] doesn't exist.",
			params: DeleteFileParams{
				Foldername: "/home",
				Filename:   "qa.key",
			},
			wantErr: ErrFileNotExists,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := fs.Root.DeleteFile(tt.params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteFile() error=%v, want=%v", err, tt.wantErr)
			}

			folder, _ := fs.Root.findFolder(tt.params.Foldername)
			if tt.assert != nil {
				tt.assert(t, folder)
			}
		})
	}
}

func TestFolder_ListFiles(t *testing.T) {
	fs := testFileSystem()

	tests := []struct {
		name    string
		params  ListFilesParams
		wantErr error
		assert  func(t *testing.T, files []*File)
	}{
		{
			name: "by default",
			params: ListFilesParams{
				Foldername: "/home",
				Sort:       &FileSystemSortParams{},
			},
			wantErr: nil,
			assert: func(t *testing.T, files []*File) {
				want := []string{"dev.conf", "prod.conf", "qa.conf"}
				for i, file := range files {
					if file.Name != want[i] {
						t.Errorf("ListFiles() file=%v, want=%v", file.Name, want[i])
					}
				}
			},
		},
		{
			name: "by name",
			params: ListFilesParams{
				Foldername: "/home",
				Sort:       &FileSystemSortParams{ByName: pkg.SortKind_Desc},
			},
			wantErr: nil,
			assert: func(t *testing.T, files []*File) {
				want := []string{"qa.conf", "prod.conf", "dev.conf"}
				for i, file := range files {
					if file.Name != want[i] {
						t.Errorf("ListFiles() file=%v, want=%v", file.Name, want[i])
					}
				}
			},
		},
		{
			name: "by created",
			params: ListFilesParams{
				Foldername: "/home",
				Sort:       &FileSystemSortParams{ByCreated: pkg.SortKind_Desc},
			},
			wantErr: nil,
			assert: func(t *testing.T, files []*File) {
				want := []string{"prod.conf", "qa.conf", "dev.conf"}
				for i, file := range files {
					if file.Name != want[i] {
						t.Errorf("ListFiles() file=%v, want=%v", file.Name, want[i])
					}
				}
			},
		},
		{
			name: "The folder is empty.",
			params: ListFilesParams{
				Foldername: "/etc",
				Sort:       &FileSystemSortParams{},
			},
			wantErr: ErrListFileEmpty,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			files, err := fs.Root.ListFiles(tt.params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ListFiles() error=%v, want=%v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert(t, files)
			}
		})
	}
}
