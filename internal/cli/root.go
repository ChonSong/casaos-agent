package cli

import (
	"fmt"
	"os"

	"github.com/chonSong/casaos-agent/internal/config"
	"github.com/spf13/cobra"
)

var cfg *config.Config

// NewRootCmd returns the root cobra command
func NewRootCmd(c *config.Config) *cobra.Command {
	cfg = c
	root := &cobra.Command{
		Use:   "casaos-agent",
		Short: "Agent-native CasaOS CLI",
		Long: `CasaOS CLI designed for autonomous agent operation.
All commands return JSON output by default when --json is set or
CASAOS_JSON environment variable is set.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Respect --json flag
			if cmd.Flags().Changed("json") {
				cfg.Output.ForceJSON = true
				cfg.Output.Format = "json"
			}
			if cmd.Flags().Changed("yes") {
				cfg.Output.Yes = true
			}
			return nil
		},
	}

	// Global flags
	root.PersistentFlags().BoolVar(&cfg.Output.ForceJSON, "json", false, "Force JSON output")
	root.PersistentFlags().BoolVar(&cfg.Output.Watch, "watch", false, "Watch / stream output for long-running operations")
	root.PersistentFlags().BoolP("yes", "y", false, "Skip all confirmation prompts")
	root.PersistentFlags().StringVar(&cfg.GlobalURL, "url", cfg.GlobalURL, "CasaOS API root URL")
	root.PersistentFlags().StringVar(&cfg.SocketPath, "socket", cfg.SocketPath, "CasaOS UNIX socket path")
	root.PersistentFlags().StringVar(&cfg.Token, "token", cfg.Token, "API Authorization Bearer token")
	root.PersistentFlags().IntVar(&cfg.Timeout, "timeout", cfg.Timeout, "Request timeout in seconds")

	// Add command groups
	root.AddCommand(NewAppCmd(cfg))
	root.AddCommand(NewContainerCmd(cfg))
	root.AddCommand(NewSystemCmd(cfg))
	root.AddCommand(NewStorageCmd(cfg))
	root.AddCommand(NewWebhookCmd(cfg))
	root.AddCommand(NewEventCmd(cfg))
	root.AddCommand(NewGatewayCmd(cfg))
	root.AddCommand(NewVersionCmd())

	return root
}

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("casaos-agent v0.1.0")
		},
	}
}

func confirm(msg string) bool {
	if cfg.Output.Yes {
		return true
	}
	fmt.Printf("%s [y/N]: ", msg)
	var input string
	fmt.Scanln(&input)
	return input == "y" || input == "Y"
}

func exitJSON(cmd string, err error) {
	// Handled by root Execute context
	panic(err)
}

func requireURL() string {
	if cfg.GlobalURL == "" && cfg.SocketPath == "" {
		fmt.Fprintln(os.Stderr, "Error: must set either --url or --socket")
		os.Exit(1)
	}
	if cfg.SocketPath != "" {
		return "unix:" + cfg.SocketPath
	}
	return cfg.GlobalURL
}
