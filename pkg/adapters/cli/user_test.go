package cli

import (
	"testing"
)

func Test_registerUser(t *testing.T) {
	testcase := []struct {
		name         string
		request      string
		hasErr       bool
		wantResponse string
	}{
		{
			name:         "",
			request:      `register user2`,
			hasErr:       true,
			wantResponse: "1-[user2]\n",
		},
	}

	fixture(t, testCommand(), testcase)
}
