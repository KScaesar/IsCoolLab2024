package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/KScaesar/IsCoolLab2024/pkg"
	"github.com/KScaesar/IsCoolLab2024/pkg/app"
)

func createFile(svc app.FileService) *cobra.Command {
	const prompt = "create-file [username] [foldername] [filename] [description]?"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "file", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	command.Args = cobra.RangeArgs(3, 4)
	command.Run = func(cmd *cobra.Command, args []string) {
		username := args[0]
		foldername := args[1]
		filename := args[2]
		var description string
		if len(args) >= 4 {
			description = args[3]
		}
		req := app.CreateFileParams{
			Foldername:  foldername,
			Filename:    filename,
			Description: description,
			CreatedTime: time.Now(),
		}

		err := svc.CreateFile(cmd.Context(), username, req)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(),
			"Create %v in %v/%v successfully.\n",
			req.Filename,
			username,
			req.Foldername,
		)
	}
	return command
}

func deleteFile(svc app.FileService) *cobra.Command {
	const prompt = "delete-file [username] [foldername] [filename]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "file", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	command.Args = cobra.ExactArgs(3)
	command.Run = func(cmd *cobra.Command, args []string) {
		username := args[0]
		req := app.DeleteFileParams{
			Foldername: args[1],
			Filename:   args[2],
		}

		err := svc.DeleteFile(cmd.Context(), username, req)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(),
			"Delete %v in %v/%v successfully.\n",
			req.Filename,
			username,
			req.Foldername,
		)
	}
	return command
}

func listFiles(svc app.FileService) *cobra.Command {
	const prompt = "list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]"

	command := &cobra.Command{
		Use: prompt,
	}
	pkg.CliSetUsage(command, "file", prompt)
	pkg.CliSetActivePrompt(command, prompt)

	sortByName := command.Flags().String("sort-name", "asc", "sort by file name [asc|desc]")
	sortByCreated := command.Flags().String("sort-created", "", "sort by created [asc|desc]")
	command.MarkFlagsMutuallyExclusive("sort-name", "sort-created")

	command.Args = cobra.ExactArgs(2)
	command.Run = func(cmd *cobra.Command, args []string) {
		username := args[0]
		foldername := args[1]
		req := app.ListFilesParams{
			Foldername: foldername,
			Sort: &app.FileSystemSortParams{
				ByName:    pkg.SortKind(*sortByName),
				ByCreated: pkg.SortKind(*sortByCreated),
			},
		}

		files, err := svc.ListFiles(cmd.Context(), username, req)
		if err != nil {
			if errors.Is(err, app.ErrListFileEmpty) {
				fmt.Fprintf(cmd.OutOrStdout(), "%v\n", err)
				return
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "%v", err)
			return
		}

		renderByText := func(file *app.ViewFile) {
			if file.Description == "" {
				fmt.Fprintf(cmd.OutOrStdout(),
					"%v %v %v %v\n",
					file.Filename,
					file.CreatedTime.Format("2006-01-02 15:04:05"),
					file.Fodlername,
					file.Username,
				)
				return
			}

			fmt.Fprintf(cmd.OutOrStdout(),
				"%v %v %v %v %v\n",
				file.Filename,
				file.Description,
				file.CreatedTime.Format("2006-01-02 15:04:05"),
				file.Fodlername,
				file.Username,
			)
		}

		for _, folder := range files {
			renderByText(&folder)
		}
	}
	return command
}
