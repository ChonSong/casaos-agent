package cli

import (
	"fmt"

	"github.com/chonSong/casaos-agent/internal/config"
	"github.com/spf13/cobra"
)

func NewContainerCmd(cfg *config.Config) *cobra.Command {
	container := &cobra.Command{
		Use:   "container",
		Short: "Raw Docker container management",
	}
	container.AddCommand(containerList(), containerInspect(), containerExec(), containerStats())
	return container
}

func containerList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all Docker containers (running and stopped)",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v2/docker/containers/json
			fmt.Println(`{"ok":true,"command":"container list","data":{"containers":[]}}`)
		},
	}
}

func containerInspect() *cobra.Command {
	return &cobra.Command{
		Use:   "inspect <id>",
		Short: "Full container state as JSON",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v2/docker/containers/{id}/json
			fmt.Println(`{"ok":true,"command":"container inspect","data":{}}`)
		},
	}
}

func containerExec() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec <id> <command...>",
		Short: "Run a command inside a running container",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: POST /v2/docker/containers/{id}/exec + start
			fmt.Println(`{"ok":true,"command":"container exec","data":{"exit_code":0}}`)
		},
	}
	cmd.Flags().Bool("detach", false, "Run command in background")
	return cmd
}

func containerStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats <id>",
		Short: "Live resource stats for a container",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v2/docker/containers/{id}/stats (stream)
			fmt.Println(`{"ok":true,"command":"container stats","data":{"cpu_percent":0,"memory_usage_bytes":0,"network_rx_bytes":0,"network_tx_bytes":0}}`)
		},
	}
	cmd.Flags().Bool("watch", false, "Stream stats continuously")
	return cmd
}
