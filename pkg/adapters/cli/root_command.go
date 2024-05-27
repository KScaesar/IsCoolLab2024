package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KScaesar/IsCoolLab2024/pkg/app"
)

func NewRootCommand(svc *app.Service) *Command {
	root := &cobra.Command{
		Use:                "vFS",
		Short:              "A Simple Virtual File System",
		DisableSuggestions: true,
		SilenceErrors:      true,
	}
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		fmt.Fprintf(cmd.ErrOrStderr(), cmd.UsageString())
		return nil
	})

	// user
	root.AddCommand(registerUser(svc.UserService))

	// folder
	root.AddCommand(createFolder())
	root.AddCommand(deleteFolder())
	root.AddCommand(listFolders())
	root.AddCommand(renameFolder())

	// file
	root.AddCommand(createFile())
	root.AddCommand(deleteFile())
	root.AddCommand(listFiles())

	return &Command{root}
}

type Command struct {
	*cobra.Command
}

func (c *Command) Execute() {
	err := c.Command.Execute()
	if err != nil {
		if strings.Contains(err.Error(), "unknown command") {
			fmt.Fprintf(c.ErrOrStderr(), "Error: Unrecognized command\n")
			return
		}
	}
}
