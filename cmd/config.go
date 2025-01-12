package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/RobKokochak/priori/internal/fileops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setPathCmd)
	configCmd.AddCommand(getPathCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure priori settings",
	Long:  `Configure various settings for priori, such as the path to the tasks file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		validCommands := []string{"set-path", "get-path"}
		if len(args) > 0 && !slices.Contains(validCommands, args[0]) {
			fullTask := "config " + strings.Join(args, " ")
			return cmd.Parent().RunE(cmd.Parent(), []string{fullTask})
		}

		return cmd.Help()
	},
}

var setPathCmd = &cobra.Command{
	Use:   "set-path [path]",
	Short: "Set the path to the tasks file",
	Long:  `Set the directory path where your tasks file will be stored`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tasksFilePath := args[0]
		if err := fileops.SetTasksFilePath(tasksFilePath); err != nil {
			return err
		}
		fmt.Printf("Successfully set tasks path to: %s\n", tasksFilePath+fileops.TASKS_FILENAME)
		return nil
	},
}

var getPathCmd = &cobra.Command{
	Use:   "get-path",
	Short: "Show the current tasks file path",
	Long:  `Display the current directory path where the tasks file is stored`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		tasksFilePath, err := fileops.GetTasksFilePath()
		if err != nil {
			return err
		}
		fmt.Printf("Current tasks path: %s\n", tasksFilePath)
		return nil
	},
}
