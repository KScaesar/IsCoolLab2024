package cli_test

import (
	"testing"
)

func TestNewRootCommand(t *testing.T) {
	testcase := []struct {
		name         string
		request      string
		hasErr       bool
		wantResponse string
	}{
		{
			// https://github.com/spf13/cobra/blob/v1.8.0/command.go#L1097-L1100
			name:         "unknown command",
			request:      `register-user user2`,
			hasErr:       true,
			wantResponse: "Error: Unrecognized command\n",
		},
		{
			// 為了控制 Error message, 開啟 SilenceErrors
			// 導致 Usage message 輸出到 stdout
			//
			// https://github.com/spf13/cobra/blob/v1.8.0/command.go#L1130-L1134
			// https://github.com/spf13/cobra/blob/v1.8.0/args.go#L97
			name:         "missing arg",
			request:      `rename-folder user4 folder15`,
			hasErr:       false,
			wantResponse: "rename-folder [username] [foldername] [new-folder-name]\n",
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
