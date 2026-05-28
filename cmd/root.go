package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// debugFlag is set by the persistent --debug flag on the root command and
// read by subcommands to enable verbose output.
var debugFlag bool

// jsonOutput is set by the persistent --json flag and switches command output
// (results and errors) to machine-readable JSON for scripting and agents.
var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "apialerts",
	Short: "API Alerts CLI - send events from your terminal",
	Long: `A command-line interface for apialerts.com. Send events from your terminal, scripts, and CI/CD pipelines.

Interactive setup:
  apialerts init
  apialerts send -m "Hello from the terminal"

Headless / agents / CI (no config file written):
  export APIALERTS_API_KEY="your-api-key"
  apialerts send --json -m "Hello"

API key resolution order: --key flag > APIALERTS_API_KEY env var > saved config file.
Pass --json on any command for parseable output ({"workspace","channel","warnings"} on stdout, {"error"} on stderr).`,
	Version: Version,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Enable debug output (request URL, integration, payload)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, `Emit JSON instead of human output (success: {"workspace","channel","warnings"} on stdout; error: {"error"} on stderr)`)
}

func Execute() {
	// Pre-scan os.Args for --json so we can intercept cobra's own error path
	// (unknown command, unknown flag, --version) BEFORE cobra parses flags.
	// This guarantees a single, consistent JSON-on-stderr error contract for
	// agents regardless of whether the failure is parse-time or runtime.
	preScanJSON()
	if jsonOutput {
		rootCmd.SetErr(io.Discard)
		rootCmd.SetVersionTemplate(`{"version":"{{.Version}}"}` + "\n")
	}

	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		if jsonOutput {
			payload, _ := json.Marshal(map[string]string{"error": err.Error()})
			fmt.Fprintln(os.Stderr, string(payload))
		} else {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
		os.Exit(1)
	}
}

// preScanJSON inspects os.Args to set jsonOutput before cobra parses flags.
// It accepts `--json`, `--json=true`, and `--json=false`.
func preScanJSON() {
	for _, a := range os.Args[1:] {
		switch a {
		case "--json", "--json=true":
			jsonOutput = true
		case "--json=false":
			jsonOutput = false
		}
	}
}
