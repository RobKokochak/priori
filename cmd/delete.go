package cmd

import (
	"fmt"
	"strings"

	"github.com/RobKokochak/priori/internal/cmdops"
	"github.com/RobKokochak/priori/internal/fileops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete [index]",
	Short: "Delete a task by its index in the priority section",
	Long:  `Delete a task by specifying its index within the priority section`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		priority, err := cmdops.GetPriority(cmd.Parent().PersistentFlags())
		if err != nil {
			return err
		}

		index := args[0]

		filePath, err := fileops.GetTasksFilePath()
		if err != nil {
			return err
		}

		if err := fileops.DeleteTask(filePath, priority, index); err != nil {
			if strings.Contains(err.Error(), "task not found") {
				fmt.Println(err.Error())
				return nil
			}
			return err
		}

		fmt.Println("Task deleted successfully")
		return nil
	},
}
