package cli_test

import (
	"bytes"
	_ "embed"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KScaesar/IsCoolLab2024/pkg"
	"github.com/KScaesar/IsCoolLab2024/pkg/adapters"
	"github.com/KScaesar/IsCoolLab2024/pkg/adapters/database"
	"github.com/KScaesar/IsCoolLab2024/pkg/inject"
)

//go:embed testdata.sql
var testdata string

// System Under Test
var sut *adapters.Infra

func TestMain(m *testing.M) {
	// setup()
	code := m.Run()
	// teardown()

	const success = 0
	if code != success {
		os.Exit(code)
	}
}

func setup() {
	conf := &database.GormConfing{
		// Dsn: "vFS.db",
		Dsn:     ":memory:",
		Migrate: true,
		// Debug:   true,
	}

	var err error
	sut, err = inject.NewInfra(conf)
	if err != nil {
		panic(err)
	}

	if !conf.Migrate {
		return
	}
	err = sut.Database.Exec(testdata).Error
	if err != nil {
		panic(err)
	}
}

func teardown() {
	sut.Cleanup()
	sut = nil
}

func fixture(
	t *testing.T,
	tests []struct {
		name         string
		request      string
		hasErr       bool
		wantResponse string
	},
) {
	// t.Parallel()
	setup()
	defer teardown()

	spyStdout := &bytes.Buffer{}
	spyStderr := &bytes.Buffer{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			root := inject.NewRootCommand(sut)
			root.SetOut(spyStdout)
			root.SetErr(spyStderr)
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
