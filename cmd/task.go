package cmd

import (
	"fmt"

	"strings"

	"github.com/spf13/cobra"
)

type Priority string

const (
	High   Priority = "high"
	Medium Priority = "medium"
	Low    Priority = "low"
	None   Priority = "none"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Add a task",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || strings.TrimSpace(args[0]) == "" {
			return fmt.Errorf("requires a task description")
		}
		if len(args) > 1 {
			return fmt.Errorf("Only one quote-wrapped task description is allowed")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		task := args[0]
		var priority Priority

		flags := cmd.Flags()
		highPriority, _ := flags.GetBool("high")
		mediumPriority, _ := flags.GetBool("medium")
		lowPriority, _ := flags.GetBool("low")

		if !highPriority && !mediumPriority && !lowPriority {
			priority = None
			fmt.Println("This task has no priority")
		} else if highPriority && !mediumPriority && !lowPriority {
			priority = High
			fmt.Println("This task has high priority")
		} else if !highPriority && mediumPriority && !lowPriority {
			priority = Medium
			fmt.Println("This task has medium priority")
		} else if !highPriority && !mediumPriority && lowPriority {
			priority = Low
			fmt.Println("This task has low priority")
		} else {
			return fmt.Errorf("only one priority flag can be set for a task.")
		}

		err := writeTask(task, priority)
		if err != nil {
			return fmt.Errorf("failed to write task: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.Flags().BoolP("help", "", false, "Help for the task command")
	taskCmd.Flags().BoolP("high", "h", false, "Marks the task as a high priority")
	taskCmd.Flags().BoolP("medium", "m", false, "Marks the task as a medium priority")
	taskCmd.Flags().BoolP("low", "l", false, "Marks the task as a low priority")
}
