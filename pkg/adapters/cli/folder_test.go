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
			name:         "",
			request:      "create-folder user1 folder1",
			hasErr:       false,
			wantResponse: "2-[user1 folder1]\n",
		},
		{
			name:         "",
			request:      `create-folder user1 folder2 "this-is-folder 2"`,
			hasErr:       false,
			wantResponse: "3-[user1 folder2 this-is-folder 2]\n",
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
