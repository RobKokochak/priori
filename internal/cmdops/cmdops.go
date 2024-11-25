package cmdops

import (
	"fmt"

	"github.com/RobKokochak/priori/internal/models"
	"github.com/spf13/pflag"
)

func GetPriority(flags *pflag.FlagSet) (models.Priority, error) {
	highPriority, _ := flags.GetBool("high")
	mediumPriority, _ := flags.GetBool("medium")
	lowPriority, _ := flags.GetBool("low")

	var priority models.Priority
	priorityCount := 0

	if highPriority {
		priorityCount++
		priority = models.HighPriority
	}
	if mediumPriority {
		priorityCount++
		priority = models.MediumPriority
	}
	if lowPriority {
		priorityCount++
		priority = models.LowPriority
	}

	switch priorityCount {
	case 0:
		return models.NoPriority, nil
	case 1:
		return priority, nil
	default:
		return models.NoPriority, fmt.Errorf("only one priority flag can be set for a task.")
	}
}
