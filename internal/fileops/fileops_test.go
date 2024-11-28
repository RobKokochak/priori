package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/RobKokochak/priori/internal/models"
)

func TestGetTasksFilePath(t *testing.T) {
	originalConfig := currentConfig
	defer func() {
		currentConfig = originalConfig
	}()

	t.Run("Custom path set", func(t *testing.T) {
		customPath := "/custom/path"
		currentConfig.TasksFilePath = customPath
		if got := getTasksFilePath(); got != filepath.Join(customPath, TASKS_FILENAME) {
			t.Errorf("getTasksFilePath() = %v, want %v", got, filepath.Join(customPath, TASKS_FILENAME))
		}
	})

	t.Run("Default path", func(t *testing.T) {
		currentConfig.TasksFilePath = ""
		got := getTasksFilePath()
		homeDir, _ := os.UserHomeDir()
		expected := filepath.Join(homeDir, TASKS_FILENAME)
		if got != expected {
			t.Errorf("getTasksFilePath() = %v, want %v", got, expected)
		}
	})
}

func TestWriteTask(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "priori-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	defer os.RemoveAll(tmpDir)
	originalConfig := currentConfig
	defer func() {
		currentConfig = originalConfig
	}()

	// Set the tasks file path to the temporary directory
	currentConfig.TasksFilePath = tmpDir
	tasksFilePath := filepath.Join(tmpDir, TASKS_FILENAME)

	tests := []struct {
		name     string
		task     string
		priority models.Priority
		wantErr  bool
		setup    func()
		verify   func(t *testing.T, content string)
	}{
		{
			name:     "Write medium priority task",
			task:     "Test task 1",
			priority: models.MediumPriority,
			wantErr:  false,
			setup: func() {
				os.Remove(tasksFilePath)
			},
			verify: func(t *testing.T, content string) {
				expected := "### Medium Priority\n- Test task 1"
				if content != expected {
					t.Errorf("\nExpected content:\n%s\nGot:\n%s", expected, content)
				}
			},
		},
		{
			name:     "Write high priority task",
			task:     "Test task 2",
			priority: models.HighPriority,
			wantErr:  false,
			verify: func(t *testing.T, content string) {
				expected :=
					`### High Priority
- Test task 2
### Medium Priority
- Test task 1`
				if content != expected {
					fmt.Print("content----\n" + content)
					fmt.Print("\nexpected----\n" + expected)
					t.Errorf("\nExpected content:\n%s\nGot:\n%s", expected, content)
				}
			},
		},
		{
			name:     "Write low priority task",
			task:     "Test task 3",
			priority: models.LowPriority,
			wantErr:  false,
			verify: func(t *testing.T, content string) {
				expected :=
					`### High Priority
- Test task 2
### Medium Priority
- Test task 1
### Low Priority
- Test task 3`
				if content != expected {
					t.Errorf("\nExpected content:\n%s\nGot:\n%s", expected, content)
				}
			},
		},
		{
			name:     "Write 2nd medium priority task",
			task:     "Test task 4",
			priority: models.MediumPriority,
			wantErr:  false,
			verify: func(t *testing.T, content string) {
				expected :=
					`### High Priority
- Test task 2
### Medium Priority
- Test task 4
- Test task 1
### Low Priority
- Test task 3`
				if content != expected {
					t.Errorf("\nExpected content:\n%s\nGot:\n%s", expected, content)
				}
			},
		},
		{
			name:     "Write 2nd high priority task",
			task:     "Test task 5",
			priority: models.HighPriority,
			wantErr:  false,
			verify: func(t *testing.T, content string) {
				expected :=
					`### High Priority
- Test task 5
- Test task 2
### Medium Priority
- Test task 4
- Test task 1
### Low Priority
- Test task 3`
				if content != expected {
					t.Errorf("\nExpected content:\n%s\nGot:\n%s", expected, content)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			err := WriteTask(tt.task, tt.priority)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			content, err := os.ReadFile(tasksFilePath)
			if err != nil {
				t.Fatalf("Failed to read tasks file: %v", err)
			}

			if tt.verify != nil {
				tt.verify(t, string(content))
			}
		})
	}
}
