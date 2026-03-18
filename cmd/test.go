package cmd

import (
	"fmt"

	"github.com/apialerts/apialerts-go"
	"github.com/apialerts/cli/internal/config"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Send a test event",
	Long:  "Send a test event to verify your API key and connectivity.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		if cfg.APIKey == "" {
			return fmt.Errorf("no API key configured — run: apialerts init")
		}

		apialerts.Configure(cfg.APIKey)
		apialerts.SetOverrides(IntegrationName, Version, cfg.ServerURL)

		event := apialerts.Event{
			Event:   "cli.test",
			Title:   "CLI Test Event",
			Message: "Test event from API Alerts CLI",
			Tags:    []string{"test", "cli"},
		}

		result, err := apialerts.SendAsync(event)
		if err != nil {
			return fmt.Errorf("test failed: %w", err)
		}

		fmt.Printf("✓ Test event sent to %s (%s)\n", result.Workspace, result.Channel)
		for _, w := range result.Warnings {
			fmt.Printf("! Warning: %s\n", w)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
