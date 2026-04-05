package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/chonSong/casaos-agent/internal/cli"
	"github.com/chonSong/casaos-agent/internal/config"
	"github.com/chonSong/casaos-agent/internal/output"
)

const (
	Version = "0.1.0"
)

func main() {
	cfg := config.Load()

	root := cli.NewRootCmd(cfg)
	if err := root.Execute(); err != nil {
		if cfg.Output.Format == "json" || cfg.Output.ForceJSON {
			resp := output.ErrorResponse("root", err)
			json.NewEncoder(os.Stderr).Encode(resp)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
