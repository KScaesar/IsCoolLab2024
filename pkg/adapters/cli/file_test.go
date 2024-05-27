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
			name:         "",
			request:      `create-file user1 folder1 file1 this-is-file1`,
			hasErr:       false,
			wantResponse: "4-[user1 folder1 file1 this-is-file1]\n",
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
			name:         "",
			request:      `delete-file user1 folder1 file1`,
			hasErr:       false,
			wantResponse: "3-[user1 folder1 file1]\n",
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
			name:         "",
			request:      `list-files user1 folder1 --sort-name desc`,
			hasErr:       false,
			wantResponse: "2-[user1 folder1]-(1)\n",
		},
		{
			name:         "unknown flag",
			request:      `list-files user1 folder1 --sort-createdTime desc`,
			hasErr:       true,
			wantResponse: "list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]\n",
		},
	}

	fixture(t, testcase)
}
