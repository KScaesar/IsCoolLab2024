package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KScaesar/IsCoolLab2024/pkg"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	const success = 0
	if code != success {
		os.Exit(code)
	}
}

func setup() {
	sut = NewRootCommand()
}

func teardown() {
	sut = nil
}

// System Under Test
var sut *Command

func testCommand() *Command {
	return sut
	// return NewRootCommand()
}

func fixture(
	t *testing.T,
	root *Command,
	tests []struct {
		name         string
		request      string
		hasErr       bool
		wantResponse string
	},
) {
	// t.Parallel()

	spyStdout := &bytes.Buffer{}
	spyStderr := &bytes.Buffer{}
	root.SetOut(spyStdout)
	root.SetErr(spyStderr)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			root.SetArgs(pkg.CliParse(tt.request))
			root.Execute()

			var actualResponse string
			if !tt.hasErr {
				actualResponse = spyStdout.String()
			} else {
				actualResponse = spyStderr.String()
			}

			require.Equal(t, tt.wantResponse, actualResponse)

			spyStdout.Reset()
			spyStderr.Reset()
		})
	}
}
