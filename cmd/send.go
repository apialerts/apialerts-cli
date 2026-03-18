package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apialerts/apialerts-go"
	"github.com/apialerts/cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	sendEvent   string
	sendTitle   string
	sendMessage string
	sendChannel string
	sendTags    string
	sendLink    string
	sendData    string
	sendKey     string
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an event",
	Long: `Send an event to API Alerts. Requires a message at minimum.

Properties:
  -m  message   The notification message (required)
  -e  event     Event name for filtering/routing (e.g. user.purchase, deploy.success)
  -t  title     Short title displayed above the message
  -c  channel   Target channel (uses your default channel if not set)
  -g  tags      Comma-separated tags for filtering (e.g. billing,error)
  -l  link      URL attached to the notification
  -d  data      JSON object with additional event data (e.g. '{"user":"john","plan":"pro"}')
      key       API key override (uses stored config if not set)`,
	Example: `  apialerts send -m "Deploy completed"
  apialerts send -e "user.purchase" -t "New Sale" -m "$49.99 from john@example.com" -c "payments"
  apialerts send -m "Payment failed" -c "payments" -g "billing,error"
  apialerts send -m "Build passed" -l "https://ci.example.com/build/123"
  apialerts send -e "user.signup" -m "New user registered" -d '{"plan":"pro","source":"organic"}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if sendMessage == "" {
			return fmt.Errorf("message is required — use -m \"your message\"")
		}

		// Load config once
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Resolve API key: flag > config file
		apiKey := sendKey
		if apiKey == "" {
			if cfg.APIKey == "" {
				return fmt.Errorf("no API key configured — run: apialerts init")
			}
			apiKey = cfg.APIKey
		}

		// Parse tags
		var tags []string
		if sendTags != "" {
			for _, t := range strings.Split(sendTags, ",") {
				trimmed := strings.TrimSpace(t)
				if trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		}

		// Parse data JSON
		var data map[string]any
		if sendData != "" {
			if err := json.Unmarshal([]byte(sendData), &data); err != nil {
				return fmt.Errorf("invalid JSON for --data: %w", err)
			}
		}

		// Configure and send
		apialerts.Configure(apiKey)
		apialerts.SetOverrides(IntegrationName, Version, cfg.ServerURL)

		event := apialerts.Event{
			Event:   sendEvent,
			Title:   sendTitle,
			Message: sendMessage,
			Channel: sendChannel,
			Tags:    tags,
			Link:    sendLink,
			Data:    data,
		}

		result, err := apialerts.SendAsync(event)
		if err != nil {
			return fmt.Errorf("failed to send: %w", err)
		}

		fmt.Printf("✓ Alert sent to %s (%s)\n", result.Workspace, result.Channel)
		for _, w := range result.Warnings {
			fmt.Printf("! Warning: %s\n", w)
		}

		return nil
	},
}

func init() {
	sendCmd.Flags().StringVarP(&sendEvent, "event", "e", "", "Event name (e.g. user.purchase)")
	sendCmd.Flags().StringVarP(&sendTitle, "title", "t", "", "Event title")
	sendCmd.Flags().StringVarP(&sendMessage, "message", "m", "", "Event message (required)")
	sendCmd.Flags().StringVarP(&sendChannel, "channel", "c", "", "Target channel")
	sendCmd.Flags().StringVarP(&sendTags, "tags", "g", "", "Comma-separated tags")
	sendCmd.Flags().StringVarP(&sendLink, "link", "l", "", "Associated URL")
	sendCmd.Flags().StringVarP(&sendData, "data", "d", "", "JSON object with additional event data")
	sendCmd.Flags().StringVar(&sendKey, "key", "", "API key override (instead of stored config)")
	rootCmd.AddCommand(sendCmd)
}
