package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	filename = "../tasks.md"
)

func main() {
	// Parse command-line flags
	today := flag.Bool("today", false, "Add task to today's list")
	medium := flag.Bool("m", false, "Set medium priority")
	top := flag.Bool("t", false, "Add task to the top of the list")
	flag.Parse()

	// Get the task from command-line arguments
	task := strings.Join(flag.Args(), " ")

	if task == "" {
		fmt.Println("Please provide a task.")
		return
	}

	// Determine the section
	section := "## Month"
	if *today {
		section = "## Today"
	}

	// Determine priority
	priority := "- [ ] "
	if *medium {
		priority = "- [ ] [M] "
	}

	// Create the task entry
	taskEntry := priority + task + "\n"

	// Read existing content
	content, err := os.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error reading file:", err)
		return
	}

	// Find the section and add the task
	lines := strings.Split(string(content), "\n")
	sectionFound := false
	newContent := ""

	for _, line := range lines {
		if strings.HasPrefix(line, section) {
			sectionFound = true
			newContent += line + "\n"
			if *top {
				newContent += taskEntry
			}
		} else if sectionFound && line == "" {
			if !*top {
				newContent += taskEntry
			}
			sectionFound = false
		}
		newContent += line + "\n"
	}

	if !sectionFound {
		newContent += "\n" + section + "\n" + taskEntry
	}

	// Write the updated content back to the file
	err = os.WriteFile(filename, []byte(newContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Task added successfully!")
}
