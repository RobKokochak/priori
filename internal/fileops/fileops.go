package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RobKokochak/priori/internal/models"
)

// todo: allow user to set custom filename
const TASKS_FILENAME = "tasks.md"

func WriteTask(task string, priority models.Priority) error {
	filePath := getTasksFilePath()

	content, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading tasks.md: %w", err)
	}

	var fileContent string
	if len(content) == 0 || strings.TrimSpace(string(content)) == "" {
		fileContent = "# Tasks\n\n"
	} else {
		fileContent = string(content)
	}

	priorityHeadersOrdered := []struct {
		priority models.Priority
		header   string
	}{
		{models.HighPriority, "## High Priority"},
		{models.MediumPriority, "## Medium Priority"},
		{models.LowPriority, "## Low Priority"},
		{models.NoPriority, "## No Priority"},
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
		fileContent = beforeSection + targetHeader + "\n\n" + afterSection
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

	err = os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
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
