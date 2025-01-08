package fileops

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RobKokochak/priori/internal/models"
)

// todo: allow user to set custom filename
const TASKS_FILENAME = "/Tasks.md"

// todo: make config a json in .configs
const configFileName = "priori_config.txt"

func getPathToConfig() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory: ", err)
		return ""
	}
	return filepath.Join(homeDir, configFileName)
}

func HasValidPath() bool {
	configPath := getPathToConfig()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}

	path := strings.TrimSpace(string(content))

	absPath, err := expandPath(path)
	if err != nil {
		return false
	}

	// Verify the path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return false
	}

	return true
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(homeDir, path[1:])
	}
	return filepath.Abs(path)
}

func PromptForPath() string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Please enter the location where you want to store Tasks.md: ")
		path, _ := reader.ReadString('\n')

		path = strings.TrimSpace(path)

		if path == "/" {
			fmt.Println("Cannot use root directory (/). Please choose another location.")
			continue
		}

		absPath, err := expandPath(path)
		if err != nil {
			fmt.Println("Error getting absolute path: ", err)
			continue
		}

		if absPath == "/" {
			fmt.Println("Cannot use root directory (/). Please choose another location.")
			continue
		}

		if _, err := os.Stat(absPath); !os.IsNotExist(err) {
			return absPath
		}

		fmt.Println("Invalid path. Please enter a valid directory path.")
	}
}

func SaveTasksFilePath(tasksFilePath string) error {
	configPath := getPathToConfig()

	return os.WriteFile(configPath, []byte(tasksFilePath), 0644)
}

func GetTasksFilePath() (string, error) {
	configPath := getPathToConfig()

	content, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	// todo: allow user to set filename
	return strings.TrimSpace(string(content) + TASKS_FILENAME), nil
}

// todo: allow user to set if task should go at top or bottom of section
func WriteTask(task string, priority models.Priority, filePath string) error {

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
		{models.NoPriority, "### ~"},
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
