/*
casaos-agent — Agent-native CasaOS CLI

Wraps CasaOS-CLI with:
  --json   Force structured JSON output for all commands
  --yes    Skip all confirmation prompts
  --watch  Stream output for long-running operations

This binary is the complete agent-native CLI. It replaces casaos-cli.
*/
package main

import (
	"os"

	"github.com/ChonSong/casaos-agent/CasaOS-CLI/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
