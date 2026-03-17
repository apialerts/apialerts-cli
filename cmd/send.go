package cmd

import (
	"fmt"
	"strings"

	"github.com/apialerts/cli/internal/config"
	"github.com/apialerts/apialerts-go"
	"github.com/spf13/cobra"
)

var (
	sendEvent   string
	sendTitle   string
	sendMessage string
	sendChannel string
	sendTags    string
	sendLink    string
	sendKey     string
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an event",
	Long:  "Send an event to API Alerts. Requires a message at minimum.",
	Example: `  apialerts send -m "Deploy completed"
  apialerts send -e user.purchase -t "New Sale" -m "$49.99 from john@example.com"
  apialerts send -m "Payment failed" -c payments -g billing,error
  apialerts send -m "Build passed" -l https://ci.example.com/build/123`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if sendMessage == "" {
			return fmt.Errorf("message is required — use -m \"your message\"")
		}

		// Resolve API key: flag > config file
		apiKey := sendKey
		if apiKey == "" {
			key, err := config.GetAPIKey()
			if err != nil {
				return err
			}
			apiKey = key
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

		// Configure and send
		apialerts.ConfigureWithConfig(apiKey, apialerts.Config{Debug: true})
		apialerts.SetIntegration(IntegrationName)

		event := apialerts.Event{
			Event:   sendEvent,
			Title:   sendTitle,
			Message: sendMessage,
			Channel: sendChannel,
			Tags:    tags,
			Link:    sendLink,
		}

		if err := apialerts.SendAsync(event); err != nil {
			return fmt.Errorf("failed to send: %w", err)
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
	sendCmd.Flags().StringVar(&sendKey, "key", "", "API key override (instead of stored config)")
	rootCmd.AddCommand(sendCmd)
}
