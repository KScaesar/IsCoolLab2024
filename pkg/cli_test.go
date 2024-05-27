package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliParse(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		// prompt
		{
			name: "prompt with flag",
			text: "list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]",
			want: []string{"list-files", "[username]", "[foldername]", "[--sort-name|--sort-created] [asc|desc]"},
		},
		{
			name: "prompt with dash but not flag # 1",
			text: "create-file [username] [folder--name] [filename]",
			want: []string{"create-file", "[username]", "[folder--name]", "[filename]"},
		},
		{
			name: "prompt with dash but not flag # 2",
			text: "list-files [username] [foldername] [--sort-created--time] [asc|desc]",
			want: []string{"list-files", "[username]", "[foldername]", "[--sort-created--time] [asc|desc]"},
		},
		{
			name: "prompt with optional",
			text: "create-file [username] [foldername] [filename] [description]?",
			want: []string{"create-file", "[username]", "[foldername]", "[filename]", "[description]?"},
		},

		// command
		{
			name: "command with flag",
			text: "list-files user1 folder1 --sort-name desc",
			want: []string{"list-files", "user1", "folder1", "--sort-name", "desc"},
		},
		{
			name: "command with dash but not flag",
			text: "create-file user-abc folder-abc config a--config-file",
			want: []string{"create-file", "user-abc", "folder-abc", "config", "a--config-file"},
		},
		{
			name: "command with flag which whitespace characters",
			text: `list-files "user 1" "New Folder" --filter "gopher book"`,
			want: []string{"list-files", "user 1", "New Folder", "--filter", "gopher book"},
		},
		{
			name: "command with flag which equal characters #1",
			text: `list-folders user1 --sort-name=asc`,
			want: []string{"list-folders", "user1", "--sort-name=asc"},
		},
		{
			name: "command with flag which equal characters #2",
			text: `list-files "user 1" "/keys" --filter =alkf=lakfj`,
			want: []string{"list-files", "user 1", "/keys", "--filter", "=alkf=lakfj"},
		},
		{
			name: "command with flag which whitespace and equal characters #1",
			text: `list-files user1 Folder1 --filter "gopher book" --sort-name=desc`,
			want: []string{"list-files", "user1", "Folder1", "--filter", "gopher book", "--sort-name=desc"},
		},
		{
			name: "command with flag which whitespace and equal characters #2",
			text: `list-files user1 Folder1 --filter="gopher book" --sort-name desc`,
			want: []string{"list-files", "user1", "Folder1", "--filter=gopher book", "--sort-name", "desc"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			prompts := CliParse(tt.text)
			assert.Equal(t, tt.want, prompts)
		})
	}
}
