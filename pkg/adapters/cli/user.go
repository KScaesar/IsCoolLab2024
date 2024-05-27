package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/KScaesar/IsCoolLab2024/pkg"
)

func registerUser() *cobra.Command {
	const prompt = "register [username]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "user", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.ErrOrStderr(), "%v-%v\n", len(args), args)
	}
	return command
}
