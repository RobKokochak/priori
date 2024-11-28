package cmd

import (
	"fmt"

	"strings"

	"github.com/RobKokochak/priori/internal/cmdops"
	"github.com/RobKokochak/priori/internal/fileops"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.Flags().BoolP("help", "", false, "Help for the task command")
	taskCmd.Flags().BoolP("high", "h", false, "Marks the task as a high priority")
	taskCmd.Flags().BoolP("medium", "m", false, "Marks the task as a medium priority")
	taskCmd.Flags().BoolP("low", "l", false, "Marks the task as a low priority")
}

var taskCmd = &cobra.Command{
	Use:   `task "[description]"`,
	Short: "Create a new task",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || strings.TrimSpace(args[0]) == "" {
			return fmt.Errorf("Requires a task description")
		}
		if len(args) > 1 {
			return fmt.Errorf("Only one quote-wrapped task is allowed")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		task := args[0]
		priority, getPriorityErr := cmdops.GetPriority(cmd.Flags())

		if getPriorityErr != nil {
			return getPriorityErr
		}

		writeTaskErr := fileops.WriteTask(task, priority)
		if writeTaskErr != nil {
			return writeTaskErr
		}

		fmt.Println("Task added")
		return nil
	},
}
