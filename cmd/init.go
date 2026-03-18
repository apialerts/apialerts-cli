package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/apialerts/cli/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up your API key",
	Long:  "Interactively prompt for your API key and save it to ~/.apialerts/config.json.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			return fmt.Errorf("no terminal detected — use: apialerts config --key \"your-api-key\"")
		}

		fmt.Print("Enter your API key: ")
		keyBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		key := strings.TrimSpace(string(keyBytes))
		if key == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		cfg.APIKey = key
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("API key saved: %s\n", maskAPIKey(key))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
