package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/KScaesar/IsCoolLab2024/pkg"
	"github.com/KScaesar/IsCoolLab2024/pkg/app"
)

func registerUser(svc app.UserService) *cobra.Command {
	const prompt = "register [username]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "user", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	command.Args = cobra.ExactArgs(1)
	command.Run = func(cmd *cobra.Command, args []string) {
		username := args[0]

		err := svc.Register(cmd.Context(), username, time.Now())
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Add %v successfully.\n", username)
	}
	return command
}
