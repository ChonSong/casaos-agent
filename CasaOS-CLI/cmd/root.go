/*
Copyright © 2022 IceWhaleTech

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-ini/ini"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
)

const (
	BasePathCasaOS = "v2/casaos"

	FlagDir     = "dir"
	FlagDryRun  = "dry-run"
	FlagFile    = "file"
	FlagForce   = "force"
	FlagRootURL = "root-url"

	GatewayPath = "/etc/casaos/gateway.ini"

	DefaultTimeout = 10 * time.Second
	RootGroupID    = "casaos-cli"
)

var (
	Version string
	Commit  string
	Date    string

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// Agent-native global flags
var (
	JSONOutput bool
	AutoYes    bool
	WatchMode  bool
)

// ResponseEnvelope is the standard JSON response envelope for agent-native output
type ResponseEnvelope struct {
	OK        bool        `json:"ok"`
	Command   string      `json:"command"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorBody `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// ErrorBody is the error payload
type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// StreamEvent is a single line in --watch streaming mode
type StreamEvent struct {
	Type    string      `json:"type"`
	Message string      `json:"message,omitempty"`
	Current int         `json:"current,omitempty"`
	Total   int         `json:"total,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "casaos-agent",
	Short: "Agent-native CasaOS CLI",
	Long: `CasaOS CLI designed for autonomous agent operation.

All commands return structured JSON output when --json is set.
Interactive prompts are bypassed when --yes is set.
Long-running operations support --watch for real-time streaming.

This is a fork of github.com/IceWhaleTech/CasaOS-CLI with:
  - Machine-readable JSON output (--json flag)
  - Non-interactive mode (--yes flag)
  - Streaming output (--watch flag)`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() int {
	err := rootCmd.Execute()
	if err != nil {
		if JSONOutput {
			JSONPrintResponse("root", nil, err)
		}
		return 1
	}
	return 0
}

func init() {
	url := ""

	rootCmd.PersistentFlags().StringP(FlagRootURL, "u", "", "root url of CasaOS API")

	// Agent-native flags
	rootCmd.PersistentFlags().BoolVarP(&JSONOutput, "json", "j", false, "Force JSON output for all commands")
	rootCmd.PersistentFlags().BoolVarP(&AutoYes, "yes", "y", false, "Skip all confirmation prompts")
	rootCmd.PersistentFlags().BoolVarP(&WatchMode, "watch", "w", false, "Stream output for long-running operations")

	if rootCmd.PersistentFlags().Changed(FlagRootURL) {
		url = rootCmd.PersistentFlags().Lookup(FlagRootURL).Value.String()
	} else {
		if _, err := os.Stat(GatewayPath); err == nil {
			cfgs, err := ini.Load(GatewayPath)
			if err != nil {
				log.Println("No gateway config found, use default root url")
			}

			port := cfgs.Section("gateway").Key("port").Value()
			if port != "" {
				url = fmt.Sprintf("localhost:%s", port)
			}
		}
	}

	if url == "" {
		url = "localhost:80"
	}

	rootCmd.PersistentFlags().Set(FlagRootURL, url)
	rootCmd.AddGroup(&cobra.Group{
		ID:    RootGroupID,
		Title: "Services",
	})
}

func trim(s string, l uint) string {
	if len(s) > int(l) {
		return s[:l] + "..."
	}
	return s
}

// JSONPrintResponse writes a structured JSON envelope to stdout
func JSONPrintResponse(command string, data interface{}, err error) {
	env := ResponseEnvelope{
		OK:        err == nil,
		Command:   command,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	if err != nil {
		env.Error = &ErrorBody{Code: "ERROR", Message: err.Error()}
	} else {
		env.Data = data
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(env)
}

// JSONPrintStream prints a streaming event line to stdout
func JSONPrintStream(t, message string) {
	ev := StreamEvent{
		Type:    t,
		Message: message,
	}
	data, _ := json.Marshal(ev)
	fmt.Println(string(data))
}

// JSONPrintStreamWithData prints a streaming event with data payload
func JSONPrintStreamWithData(t string, data interface{}) {
	ev := StreamEvent{
		Type: t,
		Data: data,
	}
	jsonData, _ := json.Marshal(ev)
	fmt.Println(string(jsonData))
}

// confirm returns true if the user approved, or if AutoYes is set
func confirm(msg string) bool {
	if AutoYes {
		return true
	}
	fmt.Printf("%s [y/N]: ", msg)
	var input string
	fmt.Scanln(&input)
	return input == "y" || input == "Y"
}
