package cmdops

import (
	"testing"

	"github.com/RobKokochak/priori/internal/models"
	"github.com/spf13/pflag"
)

func TestGetPriority(t *testing.T) {
	tests := []struct {
		name         string
		setupFlags   func(flags *pflag.FlagSet)
		wantPriority models.Priority
		wantErr      bool
		expectedErr  string
	}{
		{
			name: "No priority flags set",
			setupFlags: func(flags *pflag.FlagSet) {
				flags.Bool("high", false, "")
				flags.Bool("medium", false, "")
				flags.Bool("low", false, "")
			},
			wantPriority: models.NoPriority,
			wantErr:      false,
		},
		{
			name: "Multiple priority flags set",
			setupFlags: func(flags *pflag.FlagSet) {
				flags.Bool("high", true, "")
				flags.Bool("medium", true, "")
				flags.Bool("low", false, "")
			},
			wantPriority: models.NoPriority,
			wantErr:      true,
			expectedErr:  "only one priority flag can be set for a task.",
		},
		{
			name: "Single priority flag set",
			setupFlags: func(flags *pflag.FlagSet) {
				flags.Bool("high", false, "")
				flags.Bool("medium", false, "")
				flags.Bool("low", true, "")
			},
			wantPriority: models.LowPriority,
			wantErr:      false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
			test.setupFlags(flags)

			gotPriority, err := GetPriority(flags)

			if test.wantErr {
				if err == nil {
					t.Errorf("GetPriority() error = nil, wantErr %v", test.wantErr)
					return
				}
				if err.Error() != test.expectedErr {
					t.Errorf("GetPriority() error = %v, want %v", err.Error(), test.expectedErr)
					return
				}
			} else if err != nil {
				t.Errorf("GetPriority() unexpected error = %v", err)
				return
			}

			if gotPriority != test.wantPriority {
				t.Errorf("GetPriority() = %v, want %v", gotPriority, test.wantPriority)
			}
		})
	}
}
