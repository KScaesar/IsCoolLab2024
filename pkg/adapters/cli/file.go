package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/KScaesar/IsCoolLab2024/pkg"
)

func createFile() *cobra.Command {
	const prompt = "create-file [username] [foldername] [filename] [description]?"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "file", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%v-%v\n", len(args), args)
	}
	return command
}

func deleteFile() *cobra.Command {
	const prompt = "delete-file [username] [foldername] [filename]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "file", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%v-%v\n", len(args), args)
	}
	return command
}

func listFiles() *cobra.Command {
	const prompt = "list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "file", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	sortByName := command.Flags().String("sort-name", "asc", "sort by file name [asc|desc]")
	sortByCreated := command.Flags().String("sort-created", "", "sort by created [asc|desc]")
	command.MarkFlagsMutuallyExclusive("sort-name", "sort-created")
	_ = sortByName
	_ = sortByCreated

	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%v-%v-(%v)\n", len(args), args, cmd.Flags().NFlag())
	}
	return command
}
