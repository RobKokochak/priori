package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// todo: allow user to set custom filename
const TASKS_FILENAME = "tasks.md"

// writing tasks
func writeTask(task string) error {
	file, err := os.OpenFile(
		getTasksFilePath(),
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0644,
	)
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

func getTasksFilePath() string {
	// todo: allow user to set custom file path for tasks.md
	customPath := ""
	if customPath != "" {
		return customPath
	}
	return getDefaultTasksFilePath()
}

func getDefaultTasksFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return TASKS_FILENAME
	}
	return filepath.Join(homeDir, TASKS_FILENAME)
}
