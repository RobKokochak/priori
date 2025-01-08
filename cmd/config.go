package cmd

import (
	"fmt"

	"github.com/RobKokochak/priori/internal/fileops"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure priori settings",
	Long:  `Configure various settings for priori, such as the path to the tasks file`,
}

var setPathCmd = &cobra.Command{
	Use:   "set-path [path]",
	Short: "Set the path to the tasks file",
	Long:  `Set the directory path where your tasks file will be stored`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tasksFilePath := args[0]
		if err := fileops.SaveTasksFilePath(tasksFilePath); err != nil {
			return err
		}
		fmt.Printf("Successfully set tasks path to: %s\n", tasksFilePath)
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

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setPathCmd)
	configCmd.AddCommand(getPathCmd)
}
