package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	GlobalURL  string
	SocketPath string
	Token      string
	Timeout    int
	Output     OutputConfig
}

type OutputConfig struct {
	Format     string // "json", "yaml", "table"
	ForceJSON  bool
	Watch      bool
	Yes        bool
}

func Load() *Config {
	cfg := &Config{
		GlobalURL:  getEnv("CASAOS_URL", "localhost:80"),
		SocketPath: getEnv("CASAOS_SOCKET", ""),
		Token:      getEnv("CASAOS_TOKEN", ""),
		Timeout:    60,
		Output: OutputConfig{
			Format:    getEnv("CASAOS_OUTPUT", "table"),
			ForceJSON: os.Getenv("CASAOS_JSON") != "",
			Watch:     false,
			Yes:       os.Getenv("CASAOS_YES") != "",
		},
	}

	// Allow CLI flags to override via env-check in each command
	return cfg
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func ConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "casaos-agent")
	}
	return filepath.Join(os.Getenv("HOME"), ".config", "casaos-agent")
}
