package cmd

import (
	"fmt"
	"strings"

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
		if len(args) > 0 {
			fullTask := "list " + strings.Join(args, " ")
			return cmd.Parent().RunE(cmd.Parent(), []string{fullTask})
		}

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
