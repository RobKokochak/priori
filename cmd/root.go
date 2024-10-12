package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "priori",
	Short: "A CLI task manager which features dynamic priority setting and ordering",
	Long:  `Priori is a dynamic task manager.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
