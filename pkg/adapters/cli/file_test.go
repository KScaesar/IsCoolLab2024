package cli_test

import (
	"testing"
)

func Test_createFile(t *testing.T) {
	testcase := []struct {
		name         string
		request      string
		hasErr       bool
		wantResponse string
	}{
		{
			name:         "success",
			request:      "create-file user1 folder1 file4",
			hasErr:       false,
			wantResponse: "Create file4 in user1/folder1 successfully.\n",
		},
		{
			name:         "with whitespace char",
			request:      `create-file user1 folder1 file5 "this-is-file 5"`,
			hasErr:       false,
			wantResponse: "Create file5 in user1/folder1 successfully.\n",
		},
		{
			name:         "The [filename] has already existed.",
			request:      `create-file user1 folder1 file4`,
			hasErr:       true,
			wantResponse: "Error: The file4 has already existed.\n",
		},
		{
			name:         "The [username] doesn't exist.",
			request:      `create-file user4 folder1 file4`,
			hasErr:       true,
			wantResponse: "Error: The user4 doesn't exist.\n",
		},
		{
			name:         "The [foldername] doesn't exist.",
			request:      `create-file user1 folder5 file1`,
			hasErr:       true,
			wantResponse: "Error: The folder5 doesn't exist.\n",
		},
	}

	fixture(t, testcase)
}

func Test_deleteFile(t *testing.T) {
	testcase := []struct {
		name         string
		request      string
		hasErr       bool
		wantResponse string
	}{
		{
			name:         "success",
			request:      "delete-file user1 folder1 file1",
			hasErr:       false,
			wantResponse: "Delete file1 in user1/folder1 successfully.\n",
		},
		{
			name:         "The [filename] doesn't exist.",
			request:      `delete-file user1 folder1 file1`,
			hasErr:       true,
			wantResponse: "Error: The file1 doesn't exist.\n",
		},
		{
			name:         "The [username] doesn't exist.",
			request:      `delete-file user4 folder1 file4`,
			hasErr:       true,
			wantResponse: "Error: The user4 doesn't exist.\n",
		},
		{
			name:         "The [foldername] doesn't exist.",
			request:      `delete-file user1 folder5 file1`,
			hasErr:       true,
			wantResponse: "Error: The folder5 doesn't exist.\n",
		},
	}

	fixture(t, testcase)
}

func Test_listFiles(t *testing.T) {
	testcase := []struct {
		name         string
		request      string
		hasErr       bool
		wantResponse string
	}{
		{
			name:    "by default",
			request: `list-files user1 folder1`,
			hasErr:  false,
			wantResponse: `file1 2024-05-27 23:00:03 folder1 user1
file2 qa-file 2024-05-27 23:00:01 folder1 user1
file3 2024-05-27 23:00:02 folder1 user1
`,
		},
		{
			name:    "by created",
			request: `list-files user1 folder1 --sort-created desc`,
			hasErr:  false,
			wantResponse: `file1 2024-05-27 23:00:03 folder1 user1
file3 2024-05-27 23:00:02 folder1 user1
file2 qa-file 2024-05-27 23:00:01 folder1 user1
`,
		},
		{
			name:    "by name",
			request: `list-files user1 folder1 --sort-name desc`,
			hasErr:  false,
			wantResponse: `file3 2024-05-27 23:00:02 folder1 user1
file2 qa-file 2024-05-27 23:00:01 folder1 user1
file1 2024-05-27 23:00:03 folder1 user1
`,
		},
		{
			name:         "The folder is empty.",
			request:      `list-files user1 folder3 --sort-name asc`,
			hasErr:       false,
			wantResponse: "Warning: The folder is empty.\n",
		},
	}

	fixture(t, testcase)
}
