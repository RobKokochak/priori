package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// todo: allow user to set custom filename
const TASKS_FILENAME = "tasks.md"

func writeTask(task string, priority Priority) error {
	filePath := getTasksFilePath()

	fileExists := false
	if _, err := os.Stat(filePath); err == nil {
		fileExists = true
	}

	file, err := os.OpenFile(
		filePath,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error opening tasks.md: %w", err)
	}
	defer file.Close()

	if !fileExists {
		_, err = file.WriteString("# Tasks\n\n")
		if err != nil {
			return fmt.Errorf("error writing heading to tasks.md: %w", err)
		}
	}

	// todo: place task in correct priority list

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
