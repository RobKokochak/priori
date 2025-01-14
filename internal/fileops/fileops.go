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

func getConfigFilePath() string {
	homeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		fmt.Println("Error getting config: ", homeDirErr)
		return ""
	}
	configPath := filepath.Join(homeDir, ".config", "priori")
	if err := os.MkdirAll(configPath, 0755); err != nil {
		fmt.Println("Error creating config directory: ", err)
		return ""
	}
	return filepath.Join(configPath, configFileName)
}

func validatePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	path = strings.ReplaceAll(path, `\ `, ` `)

	if path == "/" {
		return "", fmt.Errorf("cannot use root directory (/)")
	}

	absPath, err := expandPath(path)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path: %w", err)
	}

	if absPath == "/" {
		return "", fmt.Errorf("cannot use root directory (/)")
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("directory does not exist: %s", absPath)
	}

	return absPath, nil
}

func HasValidPath() bool {
	configPath := getConfigFilePath()

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
		fmt.Print("Enter the location where you want to store Tasks.md: ")
		path, _ := reader.ReadString('\n')

		absPath, err := validatePath(path)
		if err != nil {
			fmt.Printf("Invalid path: %v\nPlease enter a valid directory path.\n", err)
			continue
		}

		return absPath
	}
}

func SetTasksFilePath(tasksFilePath string) error {
	absPath, err := validatePath(tasksFilePath)
	if err != nil {
		return err
	}
	configPath := getConfigFilePath()

	return os.WriteFile(configPath, []byte(absPath), 0644)
}

func GetTasksFilePath() (string, error) {
	configPath := getConfigFilePath()

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
		return fmt.Errorf("error writing to tasks file: %w", err)
	}

	return nil
}

func ReadTasks(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "Tasks file not found", nil
		}
		return "", fmt.Errorf("error retrieving tasks: %w", err)
	}

	if isTasksFileEmpty(string(content)) {
		return "Tasks file is empty", nil
	}

	return string(content), nil
}

func DeleteTask(filePath string, priority models.Priority, index string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading tasks file: %w", err)
	}

	if isTasksFileEmpty(string(content)) {
		return fmt.Errorf("Tasks file is empty")
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	var targetHeader string
	switch priority {
	case models.HighPriority:
		targetHeader = "### High Priority"
	case models.MediumPriority:
		targetHeader = "### Medium Priority"
	case models.LowPriority:
		targetHeader = "### Low Priority"
	case models.NoPriority:
		targetHeader = "### ~"
	default:
		return fmt.Errorf("invalid priority")
	}

	inTargetSection := false
	taskIndex := 0
	targetIndex := index
	sectionHasTasks := false
	taskFound := false

	// First pass: count tasks and check if target exists
	for _, line := range lines {
		if line == targetHeader {
			inTargetSection = true
			continue
		}

		if inTargetSection {
			if strings.HasPrefix(line, "### ") {
				inTargetSection = false
				continue
			}
			if strings.HasPrefix(line, "- ") {
				if fmt.Sprint(taskIndex) == targetIndex {
					taskFound = true
				} else {
					sectionHasTasks = true
				}
				taskIndex++
			}
		}
	}

	if !taskFound {
		return fmt.Errorf("task not found at index %s in %s priority section", targetIndex, priority)
	}

	// Reset for second pass
	inTargetSection = false
	taskIndex = 0

	// Second pass: build new content
	for _, line := range lines {
		if line == targetHeader {
			inTargetSection = true
			if sectionHasTasks {
				newLines = append(newLines, line)
			}
			continue
		}

		if inTargetSection {
			if strings.HasPrefix(line, "### ") {
				inTargetSection = false
				newLines = append(newLines, line)
				continue
			}
			if strings.HasPrefix(line, "- ") {
				if fmt.Sprint(taskIndex) != targetIndex {
					newLines = append(newLines, line)
				}
				taskIndex++
				continue
			}
		}

		if !inTargetSection {
			newLines = append(newLines, line)
		}
	}

	var cleanLines []string
	prevLineEmpty := false
	for _, line := range newLines {
		if line == "" {
			if !prevLineEmpty {
				cleanLines = append(cleanLines, line)
			}
			prevLineEmpty = true
		} else {
			cleanLines = append(cleanLines, line)
			prevLineEmpty = false
		}
	}

	for len(cleanLines) > 0 && cleanLines[len(cleanLines)-1] == "" {
		cleanLines = cleanLines[:len(cleanLines)-1]
	}

	// Write the modified content back to the file
	err = os.WriteFile(filePath, []byte(strings.Join(cleanLines, "\n")+"\n"), 0644)
	if err != nil {
		return fmt.Errorf("error writing to tasks file: %w", err)
	}

	return nil
}

func isTasksFileEmpty(content string) bool {
	if len(content) == 0 || len(strings.TrimSpace(string(content))) == 0 {
		return true
	}
	return false
}
