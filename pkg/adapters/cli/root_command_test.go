package cli

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
			name:         "unknown command",
			request:      `register-user user2`,
			hasErr:       true,
			wantResponse: "Error: Unrecognized command\n",
		},
	}

	fixture(t, testCommand(), testcase)
}
