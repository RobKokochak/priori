package fileops

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RobKokochak/priori/internal/models"
)

func TestSetTasksFilePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "priori-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Valid path",
			path:    filepath.Join(tmpDir, "tasks.md"),
			wantErr: false,
		},
		{
			name:    "Empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "Non-existent directory",
			path:    filepath.Join(tmpDir, "nonexistent", "tasks.md"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetTasksFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetTasksFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && currentConfig.TasksFilePath != tt.path {
				t.Errorf("SetTasksFilePath() didn't set the correct path, got = %v, want %v", currentConfig.TasksFilePath, tt.path)
			}
		})
	}
}

func TestGetTasksFilePath(t *testing.T) {
	originalConfig := currentConfig
	defer func() {
		currentConfig = originalConfig
	}()

	t.Run("Custom path set", func(t *testing.T) {
		customPath := "/custom/path/tasks.md"
		currentConfig.TasksFilePath = customPath
		if got := getTasksFilePath(); got != customPath {
			t.Errorf("getTasksFilePath() = %v, want %v", got, customPath)
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

	tasksFilePath := filepath.Join(tmpDir, "tasks.md")
	err = SetTasksFilePath(tasksFilePath)
	if err != nil {
		t.Fatalf("Failed to set tasks file path: %v", err)
	}

	tests := []struct {
		name     string
		task     string
		priority models.Priority
		wantErr  bool
		setup    func()
		verify   func(t *testing.T, content string)
	}{
		{
			name:     "Write high priority task",
			task:     "Test task 1",
			priority: models.HighPriority,
			wantErr:  false,
			setup: func() {
				os.Remove(tasksFilePath)
			},
			verify: func(t *testing.T, content string) {
				expected := "# Tasks\n### High Priority\n- Test task 1\n"
				if content != expected {
					t.Errorf("\nExpected content:\n%s\n\nGot:\n%s", expected, content)
				}
			},
		},
		{
			name:     "Write medium priority task",
			task:     "Test task 2",
			priority: models.MediumPriority,
			wantErr:  false,
			verify: func(t *testing.T, content string) {
				expected :=
					`# Tasks
### High Priority
- Test task 1
### Medium Priority
- Test task 2
`
				if content != expected {
					t.Errorf("\nExpected content:\n%s\n\nGot:\n%s", expected, content)
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
					`# Tasks
### High Priority
- Test task 1
### Medium Priority
- Test task 2
### Low Priority
- Test task 3
`
				if content != expected {
					t.Errorf("\nExpected content:\n%s\n\nGot:\n%s", expected, content)
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
					`# Tasks
### High Priority
- Test task 1
### Medium Priority
- Test task 4
- Test task 2
### Low Priority
- Test task 3
`
				if content != expected {
					t.Errorf("\nExpected content:\n%s\n\nGot:\n%s", expected, content)
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
					`# Tasks
### High Priority
- Test task 5
- Test task 1
### Medium Priority
- Test task 4
- Test task 2
### Low Priority
- Test task 3
`
				if content != expected {
					t.Errorf("\nExpected content:\n%s\n\nGot:\n%s", expected, content)
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
