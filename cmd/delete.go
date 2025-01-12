package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/RobKokochak/priori/internal/cmdops"
	"github.com/RobKokochak/priori/internal/fileops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func isValidIndex(index string) bool {
	num, err := strconv.Atoi(index)
	if err != nil {
		return false
	}
	return num >= 0
}

var deleteCmd = &cobra.Command{
	Use:   "delete [index]",
	Short: "Delete a task by its index in the priority section",
	Long:  `Delete a task by specifying its index within the priority section`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !isValidIndex(args[0]) {
			fullTask := "delete " + strings.Join(args, " ")
			return cmd.Parent().RunE(cmd.Parent(), []string{fullTask})
		}
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
