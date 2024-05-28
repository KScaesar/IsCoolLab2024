package cli

import (
	"errors"
	"fmt"
	"time"

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
			CreatedTime: time.Now(),
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

func deleteFolder(svc app.FolderService) *cobra.Command {
	const prompt = "delete-folder [username] [foldername]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "folder", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	command.Args = cobra.ExactArgs(2)
	command.Run = func(cmd *cobra.Command, args []string) {
		username := args[0]
		req := app.DeleteFolderParams{
			Foldername: args[1],
		}

		err := svc.DeleteFolder(cmd.Context(), username, req)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Delete %v successfully.\n", req.Foldername)
	}
	return command
}

func listFolders(svc app.FolderService) *cobra.Command {
	const prompt = "list-folders [username] [--sort-name|--sort-created] [asc|desc]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "folder", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	sortByName := command.Flags().String("sort-name", "asc", "sort by folder name [asc|desc]")
	sortByCreated := command.Flags().String("sort-created", "", "sort by created  [asc|desc]")
	command.MarkFlagsMutuallyExclusive("sort-name", "sort-created")

	command.Args = cobra.ExactArgs(1)
	command.Run = func(cmd *cobra.Command, args []string) {
		username := args[0]
		req := app.ListFoldersParams{
			Sort: &app.FileSystemSortParams{
				ByName:    pkg.SortKind(*sortByName),
				ByCreated: pkg.SortKind(*sortByCreated),
			},
		}

		folders, err := svc.ListFolders(cmd.Context(), username, req)
		if err != nil {
			if errors.Is(err, app.ErrListFolderEmpty) {
				fmt.Fprintf(cmd.OutOrStdout(), "%v\n", err)
				return
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
			return
		}

		renderByText := func(folder *app.ViewFolder) {
			if folder.Description == "" {
				fmt.Fprintf(cmd.OutOrStdout(),
					"%v %v %v\n",
					folder.Fodlername,
					folder.CreatedTime.Format("2006-01-02 15:04:05"),
					folder.Username,
				)
				return
			}

			fmt.Fprintf(cmd.OutOrStdout(),
				"%v %v %v %v\n",
				folder.Fodlername,
				folder.Description,
				folder.CreatedTime.Format("2006-01-02 15:04:05"),
				folder.Username,
			)
		}

		for _, folder := range folders {
			renderByText(&folder)
		}
	}
	return command
}

func renameFolder(svc app.FolderService) *cobra.Command {
	const prompt = "rename-folder [username] [foldername] [new-folder-name]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "folder", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	command.Args = cobra.ExactArgs(3)
	command.Run = func(cmd *cobra.Command, args []string) {
		username := args[0]
		req := app.RenameFolderParams{
			OldFolderName: args[1],
			NewFolderName: args[2],
		}

		err := svc.RenameFolder(cmd.Context(), username, req)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(),
			"Rename %v to %v successfully.\n",
			req.OldFolderName,
			req.NewFolderName,
		)
	}
	return command
}
