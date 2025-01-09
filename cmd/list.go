package cmd

import (
	"fmt"

	"github.com/RobKokochak/priori/internal/fileops"
	"github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List all tasks",
    Long:  `Display all tasks organized by priority`,
    RunE: func(cmd *cobra.Command, args []string) error {
        filePath, err := fileops.GetTasksFilePath()
        if err != nil {
            return err
        }

        tasks, err := fileops.ReadTasks(filePath)
        if err != nil {
            return err
        }

        fmt.Println(tasks)
        return nil
    },
}
