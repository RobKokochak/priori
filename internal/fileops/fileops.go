package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RobKokochak/priori/internal/models"
	"github.com/joho/godotenv"
)

// todo: allow user to set custom filename
const TASKS_FILENAME = "Tasks.md"

type Config struct {
	TasksFilePath string
}

var currentConfig Config

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

	tasksPath := os.Getenv("TASKS_FILE_PATH")
	currentConfig = Config{
		TasksFilePath: tasksPath,
	}
}

func getTasksFilePath() string {
	if currentConfig.TasksFilePath != "" {
		return filepath.Join(currentConfig.TasksFilePath, TASKS_FILENAME)
	}
	return getDefaultTasksFilePath()
}

func getDefaultTasksFilePath() string {
	homeDir, err := os.UserHomeDir()
	// if home isn't found, just use the current directory
	if err != nil {
		return TASKS_FILENAME
	}
	return filepath.Join(homeDir, TASKS_FILENAME)
}

// todo: allow user to set if task should go at top or bottom of section
func WriteTask(task string, priority models.Priority) error {
	filePath := getTasksFilePath()

	content, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading tasks.md: %w", err)
	}
	fileContent := strings.TrimLeft(string(content), "\n\r\t ")

	priorityHeadersOrdered := []struct {
		priority models.Priority
		header   string
	}{
		{models.HighPriority, "### High Priority"},
		{models.MediumPriority, "### Medium Priority"},
		{models.LowPriority, "### Low Priority"},
		{models.NoPriority, "### No Priority"},
	}

	targetHeader := ""
	for _, p := range priorityHeadersOrdered {
		if p.priority == priority {
			targetHeader = p.header
			break
		}
	}
	if targetHeader == "" {
		return fmt.Errorf("invalid priority")
	}

	if !strings.Contains(fileContent, targetHeader) {
		insertIndex := len(fileContent)
		for _, priorityHeader := range priorityHeadersOrdered {
			if priorityHeader.priority.LessThan(priority) {
				if idx := strings.Index(fileContent, priorityHeader.header); idx != -1 {
					insertIndex = idx
					break
				}
			}
		}

		beforeSection := fileContent[:insertIndex]
		afterSection := fileContent[insertIndex:]
		if insertIndex == 0 {
			fileContent = targetHeader + "\n" + afterSection
		} else {
			fileContent = strings.TrimRight(beforeSection, "\n") + "\n" + targetHeader + "\n" + afterSection
		}
	}

	lines := strings.Split(fileContent, "\n")
	var newLines []string
	foundSection := false
	taskAdded := false

	for i, line := range lines {
		newLines = append(newLines, line)
		if line == targetHeader {
			foundSection = true
			newLines = append(newLines, "- "+task)
			taskAdded = true
		} else if foundSection && !taskAdded && (i+1 == len(lines) || strings.HasPrefix(lines[i+1], "## ")) {
			newLines = append(newLines, "- "+task)
			taskAdded = true
		}
	}

	fileContent = strings.TrimRight(strings.Join(newLines, "\n"), "\n")
	err = os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to tasks.md: %w", err)
	}

	return nil
}
