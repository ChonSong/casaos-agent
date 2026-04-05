package cli

import (
	"fmt"

	"github.com/chonSong/casaos-agent/internal/config"
	"github.com/spf13/cobra"
)

func NewEventCmd(cfg *config.Config) *cobra.Command {
	event := &cobra.Command{
		Use:   "event",
		Short: "Interact with the CasaOS MessageBus",
	}
	event.AddCommand(eventListTypes(), eventSubscribe(), eventPublish())
	return event
}

func eventListTypes() *cobra.Command {
	return &cobra.Command{
		Use:   "list-types",
		Short: "List all event types registered in the MessageBus",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v2/message_bus/event_type
			fmt.Println(`{"ok":true,"command":"event list-types","data":{"event_types":[]}}`)
		},
	}
}

func eventSubscribe() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribe <event-name>",
		Short: "Subscribe to an event type and stream events to stdout",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Connect to MessageBus WebSocket, filter by event name
			fmt.Println(`{"ok":true,"command":"event subscribe","data":{"status":"subscribed","event":"` + args[0] + `"}}`)
		},
	}
	return cmd
}

func eventPublish() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish <name>",
		Short: "Publish an event to the MessageBus (for testing)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			data, _ := cmd.Flags().GetString("data")
			_ = data
			// TODO: POST /v2/message_bus/event
			fmt.Println(`{"ok":true,"command":"event publish","data":{"status":"published"}}`)
		},
	}
	cmd.Flags().String("data", "{}", "JSON data payload")
	return cmd
}
