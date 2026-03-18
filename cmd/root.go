package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "apialerts",
	Short:   "API Alerts CLI — send events from your terminal",
	Long: `A command-line interface for apialerts.com. Send events from your terminal, scripts, and CI/CD pipelines.

Get started:
  apialerts init
  apialerts send -m "Hello from the terminal"`,
	Version: Version,
}

func Execute() {
	rootCmd.SilenceErrors = true
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
