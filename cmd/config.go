package cmd

import (
	"fmt"

	"github.com/apialerts/cli/internal/config"
	"github.com/spf13/cobra"
)

var configKey string
var configServerURL string
var unsetKey bool

func maskAPIKey(key string) string {
	n := len(key)
	switch {
	case key == "":
		return "No API key configured."
	case n <= 3:
		return "***"
	case n <= 10:
		return key[:1] + "..." + key[n-1:]
	default:
		return key[:6] + "..." + key[n-4:]
	}
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the CLI",
	Long:  "Set your API key for authentication. The key is stored in ~/.apialerts/config.json.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if unsetKey {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			cfg.APIKey = ""
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to unset API key: %w", err)
			}
			fmt.Println("API key removed.")
			return nil
		}

		if cmd.Flags().Changed("server-url") {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			cfg.ServerURL = configServerURL
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save server URL: %w", err)
			}
			if configServerURL == "" {
				fmt.Println("Server URL reset to default.")
			} else {
				fmt.Printf("Server URL set to: %s\n", configServerURL)
			}
			return nil
		}

		if configKey != "" {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			cfg.APIKey = configKey
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Println("API key saved.")
			return nil
		}

		// Show current config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		if cfg.APIKey == "" {
			fmt.Println("No API key configured.")
			fmt.Println("Run: apialerts init")
		} else {
			fmt.Printf("API Key: %s\n", maskAPIKey(cfg.APIKey))
			if cfg.ServerURL != "" {
				fmt.Printf("Server URL: %s\n", cfg.ServerURL)
			}
		}
		return nil
	},
}

func init() {
	configCmd.Flags().StringVar(&configKey, "key", "", "Your API Alerts API key")
	configCmd.Flags().BoolVar(&unsetKey, "unset", false, "Remove the stored API key")
	configCmd.Flags().StringVar(&configServerURL, "server-url", "", "Override the API server URL")
	configCmd.Flags().MarkHidden("server-url")
	rootCmd.AddCommand(configCmd)
}
