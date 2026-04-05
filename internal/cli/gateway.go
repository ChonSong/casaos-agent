package cli

import (
	"fmt"

	"github.com/chonSong/casaos-agent/internal/config"
	"github.com/spf13/cobra"
)

func NewGatewayCmd(cfg *config.Config) *cobra.Command {
	gateway := &cobra.Command{
		Use:   "gateway",
		Short: "CasaOS Gateway management",
	}
	gateway.AddCommand(gatewayRoutes(), gatewayStatus())
	return gateway
}

func gatewayRoutes() *cobra.Command {
	return &cobra.Command{
		Use:   "routes",
		Short: "List all gateway routes",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v2/gateway/routes
			fmt.Println(`{"ok":true,"command":"gateway routes","data":{"routes":[]}}`)
		},
	}
}

func gatewayStatus() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Gateway health and status",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v2/gateway/status
			fmt.Println(`{"ok":true,"command":"gateway status","data":{"healthy":true}}`)
		},
	}
}
