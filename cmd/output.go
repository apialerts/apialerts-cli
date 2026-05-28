package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/apialerts/apialerts-go"
	"github.com/apialerts/cli/internal/config"
)

// resolveAPIKey returns the API key to use, checking sources in priority order:
// explicit --key flag, APIALERTS_API_KEY env var, then the saved config file.
func resolveAPIKey(flagKey string, cfg *config.CLIConfig) (string, error) {
	if flagKey != "" {
		return flagKey, nil
	}
	if v := os.Getenv("APIALERTS_API_KEY"); v != "" {
		return v, nil
	}
	if cfg.APIKey != "" {
		return cfg.APIKey, nil
	}
	return "", fmt.Errorf("no API key configured - run: apialerts init, set APIALERTS_API_KEY, or pass --key")
}

// jsonSuccess is the success-path JSON shape emitted by --json. Field order
// matches the documented schema: workspace, channel, warnings.
type jsonSuccess struct {
	Workspace string   `json:"workspace"`
	Channel   string   `json:"channel"`
	Warnings  []string `json:"warnings"`
}

// printResult writes the successful send result to stdout. In --json mode it
// emits a single-line JSON object; otherwise a human-readable line plus any
// warnings.
func printResult(result *apialerts.Result, humanPrefix string) {
	if jsonOutput {
		warnings := result.Warnings
		if warnings == nil {
			warnings = []string{}
		}
		payload, _ := json.Marshal(jsonSuccess{
			Workspace: result.Workspace,
			Channel:   result.Channel,
			Warnings:  warnings,
		})
		fmt.Println(string(payload))
		return
	}
	fmt.Printf("%s%s (%s)\n", humanPrefix, result.Workspace, result.Channel)
	for _, w := range result.Warnings {
		fmt.Printf("! Warning: %s\n", w)
	}
}
