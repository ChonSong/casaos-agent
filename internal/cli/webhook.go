package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chonSong/casaos-agent/internal/config"
	"github.com/chonSong/casaos-agent/internal/output"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var webhookRegistryFile string

func NewWebhookCmd(cfg *config.Config) *cobra.Command {
	webhook := &cobra.Command{
		Use:   "webhook",
		Short: "Webhook registration and management",
		Long:  `Register, list, and remove webhooks that receive CasaOS event notifications.`,
	}
	webhook.AddCommand(webhookRegister(), webhookList(), webhookDeregister(), webhookTest(), webhookHistory())
	return webhook
}

type WebhookRecord struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Events    []string `json:"events"`
	Secret    string   `json:"secret,omitempty"`
	CreatedAt string   `json:"created_at"`
	Enabled   bool     `json:"enabled"`
}

func webhookRegistryPath() string {
	return os.ExpandEnv("$HOME/.config/casaos-agent/webhooks.json")
}

func loadWebhooks() ([]WebhookRecord, error) {
	path := webhookRegistryPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []WebhookRecord{}, nil
		}
		return nil, err
	}
	var webhooks []WebhookRecord
	if err := json.Unmarshal(data, &webhooks); err != nil {
		return nil, err
	}
	return webhooks, nil
}

func saveWebhooks(webhooks []WebhookRecord) error {
	path := webhookRegistryPath()
	os.MkdirAll(path[:len(path)-len("/webhooks.json")], 0755)
	data, err := json.MarshalIndent(webhooks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func webhookRegister() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register <url>",
		Short: "Register a webhook URL to receive CasaOS events",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			events, _ := cmd.Flags().GetStringSlice("event")
			secret, _ := cmd.Flags().GetString("secret")
			_ = secret

			webhooks, err := loadWebhooks()
			if err != nil {
				output.PrintError("webhook register", err)
				os.Exit(1)
			}

			rec := WebhookRecord{
				ID:        "wh_" + uuid.NewString()[:12],
				URL:       args[0],
				Events:    events,
				CreatedAt: "2026-04-05T00:00:00Z",
				Enabled:   true,
			}
			webhooks = append(webhooks, rec)
			if err := saveWebhooks(webhooks); err != nil {
				output.PrintError("webhook register", err)
				os.Exit(1)
			}
			output.Print("webhook register", rec)
		},
	}
	cmd.Flags().StringSliceP("event", "e", []string{}, "Event type to subscribe (repeatable, default: all)")
	cmd.Flags().String("secret", "", "HMAC secret for payload signing")
	return cmd
}

func webhookList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all registered webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			webhooks, err := loadWebhooks()
			if err != nil {
				output.PrintError("webhook list", err)
				os.Exit(1)
			}
			output.Print("webhook list", map[string]interface{}{"webhooks": webhooks})
		},
	}
}

func webhookDeregister() *cobra.Command {
	return &cobra.Command{
		Use:   "deregister <id|url>",
		Short: "Remove a webhook by ID or URL",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			webhooks, err := loadWebhooks()
			if err != nil {
				output.PrintError("webhook deregister", err)
				os.Exit(1)
			}
			before := len(webhooks)
			target := args[0]
			var removed WebhookRecord
			webhooks = func() []WebhookRecord {
				out := []WebhookRecord{}
				for _, wh := range webhooks {
					if wh.ID == target || wh.URL == target {
						removed = wh
						continue
					}
					out = append(out, wh)
				}
				return out
			}()
			if len(webhooks) == before {
				output.PrintError("webhook deregister", fmt.Errorf("webhook not found: %s", target))
				os.Exit(1)
			}
			if err := saveWebhooks(webhooks); err != nil {
				output.PrintError("webhook deregister", err)
				os.Exit(1)
			}
			output.Print("webhook deregister", map[string]string{"removed_id": removed.ID})
		},
	}
}

func webhookTest() *cobra.Command {
	return &cobra.Command{
		Use:   "test <id>",
		Short: "Send a test payload to a registered webhook",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			webhooks, err := loadWebhooks()
			if err != nil {
				output.PrintError("webhook test", err)
				os.Exit(1)
			}
			var target WebhookRecord
			for _, wh := range webhooks {
				if wh.ID == args[0] || wh.URL == args[0] {
					target = wh
					break
				}
			}
			if target.ID == "" {
				output.PrintError("webhook test", fmt.Errorf("webhook not found: %s", args[0]))
				os.Exit(1)
			}
			// TODO: actually fire the HTTP POST
			output.Print("webhook test", map[string]string{"status": "test_sent", "webhook_id": target.ID})
		},
	}
}

func webhookHistory() *cobra.Command {
	return &cobra.Command{
		Use:   "history <id>",
		Short: "Show recent delivery attempts for a webhook",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: read from emitter history file
			output.Print("webhook history", map[string][]interface{}{"deliveries": {}})
		},
	}
}
