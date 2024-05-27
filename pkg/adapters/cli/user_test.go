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
			request:      `register user3`,
			hasErr:       false,
			wantResponse: "Add user3 successfully.\n",
		},
		{
			name:         "The [username] has already existed.",
			request:      `register user1`,
			hasErr:       true,
			wantResponse: "Error: The user1 has already existed.\n",
		},
	}

	fixture(t, testcase)
}
