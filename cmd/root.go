package cmd

import (
	"os"

	"fmt"

	"github.com/spf13/cobra"

	"strings"

	"github.com/RobKokochak/priori/internal/cmdops"
	"github.com/RobKokochak/priori/internal/fileops"
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("help", "", false, "Help for the task command")
	rootCmd.Flags().BoolP("high", "h", false, "Marks the task as a high priority")
	rootCmd.Flags().BoolP("medium", "m", false, "Marks the task as a medium priority")
	rootCmd.Flags().BoolP("low", "l", false, "Marks the task as a low priority")
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
