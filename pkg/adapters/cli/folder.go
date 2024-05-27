package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/KScaesar/IsCoolLab2024/pkg"
	"github.com/KScaesar/IsCoolLab2024/pkg/app"
)

func createFolder(svc app.FolderService) *cobra.Command {
	const prompt = "create-folder [username] [foldername] [description]?"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "folder", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	command.Args = cobra.RangeArgs(2, 3)
	command.Run = func(cmd *cobra.Command, args []string) {
		username := args[0]
		foldername := args[1]
		var description string
		if len(args) >= 3 {
			description = args[2]
		}
		req := app.CreateFolderParams{
			Foldername:  foldername,
			Description: description,
		}

		err := svc.CreateFolder(cmd.Context(), username, req)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Create %v successfully.\n", req.Foldername)
	}
	return command
}

func deleteFolder() *cobra.Command {
	const prompt = "delete-folder [username] [foldername]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "folder", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%v-%v\n", len(args), args)
	}
	return command
}

func listFolders() *cobra.Command {
	const prompt = "list-folders [username] [--sort-name|--sort-created] [asc|desc]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "folder", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	sortByName := command.Flags().String("sort-name", "asc", "sort by folder name [asc|desc]")
	sortByCreated := command.Flags().String("sort-created", "", "sort by created  [asc|desc]")
	command.MarkFlagsMutuallyExclusive("sort-name", "sort-created")
	_ = sortByName
	_ = sortByCreated

	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%v-%v-(%v)\n", len(args), args, cmd.Flags().NFlag())
	}
	return command
}

func renameFolder() *cobra.Command {
	const prompt = "rename-folder [username] [foldername] [new-folder-name]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "folder", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%v-%v\n", len(args), args)
	}
	return command
}
