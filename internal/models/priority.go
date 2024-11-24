package models

type Priority string

const (
	HighPriority   Priority = "High"
	MediumPriority Priority = "Medium"
	LowPriority    Priority = "Low"
	NoPriority     Priority = "None"
)

func (p Priority) LessThan(other Priority) bool {
	priorities := map[Priority]int{
		HighPriority:   3,
		MediumPriority: 2,
		LowPriority:    1,
		NoPriority:     0,
	}
	return priorities[p] < priorities[other]
}
