package cmd

import (
	"fmt"

	"strings"

	"github.com/spf13/cobra"
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
		flags := cmd.Flags()

		highPriority, _ := flags.GetBool("high")
		mediumPriority, _ := flags.GetBool("medium")
		lowPriority, _ := flags.GetBool("low")

		fmt.Println("Task:", args[0])

		priorityCount := 0
		if highPriority {
			priorityCount++
		}
		if mediumPriority {
			priorityCount++
		}
		if lowPriority {
			priorityCount++
		}
		switch {
		case priorityCount > 1:
			return fmt.Errorf("Only one priority flag can be set at a time.")
		case highPriority:
			fmt.Println("This is a high priority task")
		case mediumPriority:
			fmt.Println("This is a medium priority task")
		case lowPriority:
			fmt.Println("This is a low priority task")
		default:
			fmt.Println("No priority set (default: medium priority)")
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
