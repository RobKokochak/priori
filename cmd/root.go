package cmd

import (
	"os"

	"fmt"

	"github.com/spf13/cobra"

	"strings"

	"github.com/RobKokochak/priori/internal/cmdops"
	"github.com/RobKokochak/priori/internal/fileops"
)

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.PersistentFlags().BoolP("help", "", false, "Help for the task command")
	rootCmd.PersistentFlags().BoolP("high", "h", false, "Marks the task as a high priority")
	rootCmd.PersistentFlags().BoolP("medium", "m", false, "Marks the task as a medium priority")
	rootCmd.PersistentFlags().BoolP("low", "l", false, "Marks the task as a low priority")

	rootCmd.Args = func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || strings.TrimSpace(args[0]) == "" {
			return fmt.Errorf("Requires a task description")
		}

		if len(args) > 1 && args[0] == "list" {
			cmd.SetArgs([]string{strings.Join(args, " ")})
			return nil
		}
		return nil
	}
}

var rootCmd = &cobra.Command{
	Use:   "priori [task]",
	Short: "A CLI task manager with dynamic priority setting and ordering via flags, built to be fast and efficient.",
	Long:  `Priori is a dynamic task manager.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !fileops.HasValidPath() {
			path := fileops.PromptForPath()
			fileops.SetTasksFilePath(path)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || strings.TrimSpace(args[0]) == "" {
			return fmt.Errorf("Requires a task description")
		}
		if len(args) == 1 && args[0] == "list" {
			return cobra.MinimumNArgs(0)(cmd, args)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		task := strings.Join(args, " ")
		priority, getPriorityErr := cmdops.GetPriority(cmd.Flags())

		if getPriorityErr != nil {
			return getPriorityErr
		}

		filePath, getTaskFilePathErr := fileops.GetTasksFilePath()
		if getTaskFilePathErr != nil {
			return getTaskFilePathErr
		}

		writeTaskErr := fileops.WriteTask(task, priority, filePath)
		if writeTaskErr != nil {
			return writeTaskErr
		}

		fmt.Println("Task added")
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
