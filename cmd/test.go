package cmd

import (
	"fmt"
	"os"

	"github.com/apialerts/apialerts-go"
	"github.com/apialerts/cli/internal/config"
	"github.com/spf13/cobra"
)

var testKey string

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Send a test event",
	Long:  "Send a test event to verify your API key and connectivity.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		apiKey, err := resolveAPIKey(testKey, cfg)
		if err != nil {
			return err
		}

		apialerts.Configure(apiKey)
		apialerts.SetOverrides(IntegrationName, Version, cfg.ServerURL)
		if debugFlag {
			apialerts.SetDebug(true)
			fmt.Fprintf(os.Stderr, "[debug] integration: %s/%s\n", IntegrationName, Version)
			if cfg.ServerURL != "" {
				fmt.Fprintf(os.Stderr, "[debug] server URL: %s\n", cfg.ServerURL)
			}
		}

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

		printResult(result, "✓ Test event sent to ")
		return nil
	},
}

func init() {
	testCmd.Flags().StringVar(&testKey, "key", "", "API key override (instead of stored config)")
	rootCmd.AddCommand(testCmd)
}
