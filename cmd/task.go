package cmd

import (
	"fmt"

	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.Flags().BoolP("help", "", false, "Help for the task command")
	taskCmd.Flags().BoolP("high", "h", false, "Marks the task as a high priority")
	taskCmd.Flags().BoolP("medium", "m", false, "Marks the task as a medium priority")
	taskCmd.Flags().BoolP("low", "l", false, "Marks the task as a low priority")
}

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
		priority, getPriorityErr := getPriority(cmd.Flags())
		if getPriorityErr != nil {
			return getPriorityErr
		}
		writeTaskErr := writeTask(task, priority)
		if writeTaskErr != nil {
			return writeTaskErr
		}
		fmt.Println("task added")
		return nil
	},
}

func getPriority(flags *pflag.FlagSet) (Priority, error) {
	highPriority, _ := flags.GetBool("high")
	mediumPriority, _ := flags.GetBool("medium")
	lowPriority, _ := flags.GetBool("low")

	var priority Priority
	priorityCount := 0

	if highPriority {
		priorityCount++
		priority = High
	}
	if mediumPriority {
		priorityCount++
		priority = Medium
	}
	if lowPriority {
		priorityCount++
		priority = Low
	}

	switch priorityCount {
	case 0:
		return None, nil
	case 1:
		return priority, nil
	default:
		return None, fmt.Errorf("only one priority flag can be set for a task.")
	}
}
