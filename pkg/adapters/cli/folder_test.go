package cli_test

import (
	"testing"
)

func Test_createFolder(t *testing.T) {
	testcase := []struct {
		name         string
		request      string
		hasErr       bool
		wantResponse string
	}{
		{
			name:         "success",
			request:      "create-folder user1 folder4",
			hasErr:       false,
			wantResponse: "Create folder4 successfully.\n",
		},
		{
			name:         "with whitespace char",
			request:      `create-folder user2 folder2 "this-is-folder 2"`,
			hasErr:       false,
			wantResponse: "Create folder2 successfully.\n",
		},
		{
			name:         "The [foldername] has already existed.",
			request:      `create-folder user1 folder3`,
			hasErr:       true,
			wantResponse: "Error: The folder3 has already existed.\n",
		},
		{
			name:         "The [username] doesn't exist.",
			request:      `create-folder user4 folder1`,
			hasErr:       true,
			wantResponse: "Error: The user4 doesn't exist.\n",
		},
	}

	fixture(t, testcase)
}

func Test_deleteFolder(t *testing.T) {
	testcase := []struct {
		name         string
		request      string
		hasErr       bool
		wantResponse string
	}{
		{
			name:         "success",
			request:      `delete-folder user1 folder1`,
			hasErr:       false,
			wantResponse: "Delete folder1 successfully.\n",
		},
		{
			name:         "create same foldername after delete",
			request:      "create-folder user1 folder1",
			hasErr:       false,
			wantResponse: "Create folder1 successfully.\n",
		},
		{
			name:         "The [foldername] doesn't exist.",
			request:      `delete-folder user1 folder4`,
			hasErr:       true,
			wantResponse: "Error: The folder4 doesn't exist.\n",
		},
		{
			name:         "The [username] doesn't exist.",
			request:      `delete-folder user4 folder1`,
			hasErr:       true,
			wantResponse: "Error: The user4 doesn't exist.\n",
		},
	}

	fixture(t, testcase)
}

func Test_listFolders(t *testing.T) {
	testcase := []struct {
		name         string
		request      string
		hasErr       bool
		wantResponse string
	}{
		{
			name:    "by default",
			request: `list-folders user1`,
			hasErr:  false,
			wantResponse: `folder1 2024-05-27 23:00:03 user1
folder2 qa-folder 2024-05-27 23:00:01 user1
folder3 2024-05-27 23:00:02 user1
`,
		},
		{
			name:    "by created",
			request: `list-folders user1 --sort-created desc`,
			hasErr:  false,
			wantResponse: `folder1 2024-05-27 23:00:03 user1
folder3 2024-05-27 23:00:02 user1
folder2 qa-folder 2024-05-27 23:00:01 user1
`,
		},
		{
			name:    "by name",
			request: `list-folders user1 --sort-name desc`,
			hasErr:  false,
			wantResponse: `folder3 2024-05-27 23:00:02 user1
folder2 qa-folder 2024-05-27 23:00:01 user1
folder1 2024-05-27 23:00:03 user1
`,
		},
		{
			name:         "The [username] doesn't have any folders.",
			request:      `list-folders user2 --sort-name asc`,
			hasErr:       false,
			wantResponse: "Warning: The user2 doesn't have any folders.\n",
		},
	}

	fixture(t, testcase)
}

func Test_renameFolder(t *testing.T) {
	testcase := []struct {
		name         string
		request      string
		hasErr       bool
		wantResponse string
	}{
		{
			name:         "success",
			request:      `rename-folder user1 folder1 "folder isCool"`,
			hasErr:       false,
			wantResponse: "Rename folder1 to folder isCool successfully.\n",
		},
		{
			name:    "check rename for existing folder",
			request: `list-folders user1 --sort-created desc`,
			hasErr:  false,
			wantResponse: `folder isCool 2024-05-27 23:00:03 user1
folder3 2024-05-27 23:00:02 user1
folder2 qa-folder 2024-05-27 23:00:01 user1
`,
		},
		{
			name:    "check rename for existing file",
			request: `list-files user1 "folder isCool" --sort-created desc`,
			hasErr:  false,
			wantResponse: `file1 2024-05-27 23:00:03 folder isCool user1
file3 2024-05-27 23:00:02 folder isCool user1
file2 qa-file 2024-05-27 23:00:01 folder isCool user1
`,
		},
		{
			name:         "The [foldername] doesn't exist.",
			request:      `rename-folder user1 folder4 folder5`,
			hasErr:       true,
			wantResponse: "Error: The folder4 doesn't exist.\n",
		},
		{
			name:         "The [username] doesn't exist.",
			request:      `rename-folder user4 folder1 folder15`,
			hasErr:       true,
			wantResponse: "Error: The user4 doesn't exist.\n",
		},
	}

	fixture(t, testcase)
}
