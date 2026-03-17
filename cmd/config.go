package cmd

import (
	"fmt"

	"github.com/apialerts/cli/internal/config"
	"github.com/spf13/cobra"
)

var configKey string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the CLI",
	Long:  "Set your API key for authentication. The key is stored in ~/.apialerts/config.json.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if configKey == "" {
			// Show current config
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if cfg.APIKey == "" {
				fmt.Println("No API key configured.")
				fmt.Println("Run: apialerts config --key <your-api-key>")
			} else {
				masked := cfg.APIKey[:6] + "..." + cfg.APIKey[len(cfg.APIKey)-4:]
				fmt.Printf("API Key: %s\n", masked)
			}
			return nil
		}

		cfg := &config.CLIConfig{APIKey: configKey}
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println("API key saved.")
		return nil
	},
}

func init() {
	configCmd.Flags().StringVar(&configKey, "key", "", "Your API Alerts API key")
	rootCmd.AddCommand(configCmd)
}
