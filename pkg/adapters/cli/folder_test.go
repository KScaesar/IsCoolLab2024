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
			name:         "",
			request:      `delete-folder user1 folder1`,
			hasErr:       false,
			wantResponse: "2-[user1 folder1]\n",
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
			name:         "",
			request:      `list-folders user1 --sort-name asc`,
			hasErr:       false,
			wantResponse: "1-[user1]-(1)\n",
		},
		{
			name:         "unknown flag",
			request:      `list-folders user1 --sort-filename asc`,
			hasErr:       true,
			wantResponse: "list-folders [username] [--sort-name|--sort-created] [asc|desc]\n",
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
			name:         "",
			request:      `rename-folder user3 /home/doc1 "/home/doc2 go"`,
			hasErr:       false,
			wantResponse: "3-[user3 /home/doc1 /home/doc2 go]\n",
		},
	}

	fixture(t, testcase)
}
