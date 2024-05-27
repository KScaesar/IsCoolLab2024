package cli_test

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
			name:         "success",
			request:      `register user2`,
			hasErr:       false,
			wantResponse: "Add user2 successfully.\n",
		},
		{
			name:         "The [username] has already existed.",
			request:      `register user2`,
			hasErr:       true,
			wantResponse: "Error: The user2 has already existed.\n",
		},
	}

	fixture(t, testcase)
}
