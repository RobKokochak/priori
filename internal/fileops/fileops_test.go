package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/RobKokochak/priori/internal/models"
)

func TestWriteTask(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "priori-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

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

			err := WriteTask(tt.task, tt.priority, tasksFilePath)
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
