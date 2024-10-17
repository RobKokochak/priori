package cmd

import (
	"fmt"
	"os"
)

func WriteTaskToFile(task string) error {
	file, err := os.OpenFile("task-output/tasks.md", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening tasks.md: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString("- " + task + "\n")
	if err != nil {
		return fmt.Errorf("error writing to tasks.md: %w", err)
	}

	return nil
}
