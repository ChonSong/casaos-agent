package cli

import (
	"fmt"

	"github.com/chonSong/casaos-agent/internal/config"
	"github.com/spf13/cobra"
)

func NewStorageCmd(cfg *config.Config) *cobra.Command {
	storage := &cobra.Command{
		Use:   "storage",
		Short: "Local storage management",
	}
	storage.AddCommand(storageList(), storageInfo())
	return storage
}

func storageList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all mount points and disks",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v1/sys/storage (from CasaOS local-storage service)
			fmt.Println(`{"ok":true,"command":"storage list","data":{"storages":[]}}`)
		},
	}
}

func storageInfo() *cobra.Command {
	return &cobra.Command{
		Use:   "info <mount>",
		Short: "Show space usage and inode counts for a mount point",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: GET /v1/sys/storage/{mount}
			fmt.Println(`{"ok":true,"command":"storage info","data":{"mount":"","total_bytes":0,"used_bytes":0,"inodes_total":0,"inodes_used":0}}`)
		},
	}
}
