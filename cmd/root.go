package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "priori",
	Short: "A CLI task manager with dynamic priority setting and ordering via flags, built to be fast and efficient.",
	Long:  `Priori is a dynamic task manager.`,
}
