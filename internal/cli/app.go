package cli

import (
	"fmt"

	"github.com/chonSong/casaos-agent/internal/config"
	"github.com/spf13/cobra"
)

func NewAppCmd(cfg *config.Config) *cobra.Command {
	app := &cobra.Command{
		Use:   "app",
		Short: "Compose app lifecycle management",
		Long:  `List, install, start, stop, restart, and remove compose apps from CasaOS.`,
	}
	app.AddCommand(appList(), appInspect(), appInstall(), appUninstall(), appStart(), appStop(), appRestart(), appLogs(), appUpdate(), appResources())
	return app
}

func appList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all installed apps",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: call AppManagement API → JSON
			fmt.Println(`{"ok":true,"command":"app list","data":{"apps":[]},"timestamp":"2026-04-05T00:00:00Z"}`)
		},
	}
	cmd.Flags().Bool("store", false, "List apps available in the CasaOS App Store")
	return cmd
}

func appInspect() *cobra.Command {
	return &cobra.Command{
		Use:   "inspect <name>",
		Short: "Show full app state as JSON",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v2/appmanagement/apps/{name}
			fmt.Println(`{"ok":true,"command":"app inspect","data":{}}`)
		},
	}
}

func appInstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <name-or-url>",
		Short: "Install an app from store or a docker-compose URL",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			watch, _ := cmd.Flags().GetBool("watch")
			_ = dryRun
			_ = watch
			// TODO: POST /v2/appmanagement/apps/install
			fmt.Println(`{"ok":true,"command":"app install","data":{"status":"installing"}}`)
		},
	}
	cmd.Flags().Bool("dry-run", false, "Validate compose file without installing")
	cmd.Flags().BoolP("watch", "w", false, "Stream installation progress")
	return cmd
}

func appUninstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall <name>",
		Short: "Remove an installed app",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			purge, _ := cmd.Flags().GetBool("purge")
			_ = purge
			// TODO: DELETE /v2/appmanagement/apps/{name}
			fmt.Println(`{"ok":true,"command":"app uninstall","data":{"status":"uninstalled"}}`)
		},
	}
	cmd.Flags().Bool("purge", false, "Remove app data volumes")
	return cmd
}

func appStart() *cobra.Command {
	return &cobra.Command{
		Use:   "start <name>",
		Short: "Start app container(s)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: POST /v2/appmanagement/apps/{name}/start
			fmt.Println(`{"ok":true,"command":"app start","data":{"status":"started"}}`)
		},
	}
}

func appStop() *cobra.Command {
	return &cobra.Command{
		Use:   "stop <name>",
		Short: "Stop app container(s)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: POST /v2/appmanagement/apps/{name}/stop
			fmt.Println(`{"ok":true,"command":"app stop","data":{"status":"stopped"}}`)
		},
	}
}

func appRestart() *cobra.Command {
	return &cobra.Command{
		Use:   "restart <name>",
		Short: "Restart app container(s)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: POST /v2/appmanagement/apps/{name}/restart
			fmt.Println(`{"ok":true,"command":"app restart","data":{"status":"restarting"}}`)
		},
	}
}

func appLogs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs <name>",
		Short: "Stream app container logs",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tail, _ := cmd.Flags().GetInt("tail")
			_ = tail
			// TODO: GET /v2/appmanagement/apps/{name}/logs (streaming)
			fmt.Println(`{"ok":true,"command":"app logs","data":{"logs":[]}}`)
		},
	}
	cmd.Flags().Int("tail", 100, "Number of recent lines to show")
	cmd.Flags().Bool("follow", false, "Follow log output")
	return cmd
}

func appUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <name>",
		Short: "Check for and apply app updates",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkOnly, _ := cmd.Flags().GetBool("check-only")
			_ = checkOnly
			// TODO: POST /v2/appmanagement/apps/{name}/update
			fmt.Println(`{"ok":true,"command":"app update","data":{}}`)
		},
	}
	cmd.Flags().Bool("check-only", false, "Only check for updates, don't apply")
	return cmd
}

func appResources() *cobra.Command {
	return &cobra.Command{
		Use:   "resources <name>",
		Short: "Show CPU, memory, and I/O usage for an app",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v2/appmanagement/apps/{name}/resources
			fmt.Println(`{"ok":true,"command":"app resources","data":{"cpu_percent":0,"memory_used_mb":0,"io_read_bytes":0,"io_write_bytes":0}}`)
		},
	}
}
