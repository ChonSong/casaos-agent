package cli

import (
	"fmt"

	"github.com/chonSong/casaos-agent/internal/config"
	"github.com/spf13/cobra"
)

func NewSystemCmd(cfg *config.Config) *cobra.Command {
	system := &cobra.Command{
		Use:   "system",
		Short: "System information and management",
	}
	system.AddCommand(systemInfo(), systemResources(), systemUtilization(), systemHealth(), systemLogs(), systemUpdate(), systemRestart(), systemReboot())
	return system
}

func systemInfo() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Host information: hostname, OS, kernel, uptime",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v1/sys/info
			fmt.Println(`{"ok":true,"command":"system info","data":{"hostname":"casa","os":"Linux","kernel":"6.8.0","uptime_seconds":0}}`)
		},
	}
}

func systemResources() *cobra.Command {
	return &cobra.Command{
		Use:   "resources",
		Short: "Full resource inventory: CPU, memory, disk, network",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v1/sys/hardware
			fmt.Println(`{"ok":true,"command":"system resources","data":{"cpu":{"cores":4,"percent":0},"memory":{"total_bytes":0,"used_bytes":0},"disk":[],"network":{}}}`)
		},
	}
}

func systemUtilization() *cobra.Command {
	return &cobra.Command{
		Use:   "utilization",
		Short: "Live system utilization (CPU, memory, disk I/O, network)",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v1/sys/utilization
			fmt.Println(`{"ok":true,"command":"system utilization","data":{"cpu_percent":0,"memory_percent":0,"disk_io_percent":0,"network_rx_bps":0,"network_tx_bps":0}}`)
		},
	}
}

func systemHealth() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Health status of all casaos-* services",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v2/health/services
			fmt.Println(`{"ok":true,"command":"system health","data":{"services":{"casaos":{"healthy":true},"casaos-gateway":{"healthy":true},"casaos-app-management":{"healthy":true}}}}`)
		},
	}
}

func systemLogs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "CasaOS service logs",
		Run: func(cmd *cobra.Command, args []string) {
			tail, _ := cmd.Flags().GetInt("tail")
			level, _ := cmd.Flags().GetString("level")
			_ = tail
			_ = level
			// TODO: GET /v2/health/logs
			fmt.Println(`{"ok":true,"command":"system logs","data":{"logs":[]}}`)
		},
	}
	cmd.Flags().Int("tail", 100, "Number of recent log lines")
	cmd.Flags().String("level", "info", "Minimum log level: debug, info, warn, error")
	return cmd
}

func systemUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Check for or apply CasaOS system updates",
		Run: func(cmd *cobra.Command, args []string) {
			checkOnly, _ := cmd.Flags().GetBool("check-only")
			version, _ := cmd.Flags().GetString("version")
			_ = checkOnly
			_ = version
			// TODO: GET /v1/sys/version/check + POST /v1/sys/update
			fmt.Println(`{"ok":true,"command":"system update","data":{"update_available":false}}`)
		},
	}
	cmd.Flags().Bool("check-only", false, "Check for updates without applying")
	cmd.Flags().String("version", "", "Target version to install")
	return cmd
}

func systemRestart() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart the CasaOS service (not the host)",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: POST /v1/sys/stop (triggers restart)
			fmt.Println(`{"ok":true,"command":"system restart","data":{"status":"restarting"}}`)
		},
	}
}

func systemReboot() *cobra.Command {
	return &cobra.Command{
		Use:   "reboot",
		Short: "Reboot the host system",
		PreRun: func(cmd *cobra.Command, args []string) {
			if !confirm("Reboot host system?") {
				fmt.Println("Aborted.")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: POST /v1/sys/reboot (if endpoint exists)
			fmt.Println(`{"ok":true,"command":"system reboot","data":{"status":"rebooting"}}`)
		},
	}
}
