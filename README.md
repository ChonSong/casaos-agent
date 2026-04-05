# CasaOS Agent CLI

> **Status:** Phase 1 — Scaffolded. Awaiting upstream research to wire real API calls.

Agent-native CLI for managing a CasaOS instance. Forked from [IceWhaleTech/CasaOS-CLI](https://github.com/IceWhaleTech/CasaOS-CLI) with machine-readable output, non-interactive flags, and webhook registration built in.

GitHub: https://github.com/ChonSong/casaos-agent

## Quick Start

```bash
# Build
make build

# Or install via Go
go install github.com/ChonSong/casaos-agent@latest

# List apps
casaos-agent app list --json

# Install an app (non-interactive)
casaos-agent app install homeassistant --yes --watch

# System health check
casaos-agent system health --json

# Register a webhook
casaos-agent webhook register https://agent.example.com/hooks/casaos \
  --event casaos:system:utilization \
  --event casaos:file:operate

# Watch system utilization in real-time
casaos-agent system utilization --watch
```

## Key Differences from upstream `CasaOS-CLI`

| Feature | CasaOS-CLI | CasaOS-Agent |
|---------|-----------|--------------|
| Output format | Human-formatted text | JSON (default with `--json`) |
| Interactive prompts | Yes (blocking) | Bypass with `--yes` |
| Streaming output | No | `--watch` flag |
| Webhook management | No | Built-in `webhook` command group |
| Event subscription | No | `event subscribe` |
| Error codes | No | Typed error codes in JSON |

## Command Overview

```
casaos-agent app          # App lifecycle (list, install, start, stop, restart, inspect, logs, update)
casaos-agent container    # Raw Docker container management
casaos-agent system       # System info, resources, health, logs, update, restart, reboot
casaos-agent storage      # Local storage management
casaos-agent webhook      # Register, list, test, deregister webhooks
casaos-agent event        # MessageBus interaction (list types, subscribe, publish)
casaos-agent gateway      # Gateway routes and status
```

## Configuration

CasaOS-Agent uses environment variables and a config file.

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CASAOS_URL` | `localhost:80` | CasaOS API root URL |
| `CASAOS_SOCKET` | `` | UNIX socket path (alternative to URL) |
| `CASAOS_TOKEN` | `` | Bearer token for API auth |
| `CASAOS_JSON` | `""` | Set to any value to force JSON output |
| `CASAOS_YES` | `""` | Set to any value to skip all confirmations |
| `CASAOS_OUTPUT` | `table` | Output format: `table`, `json`, `yaml` |

### Config File

`~/.config/casaos-agent/config.yaml` (optional):

```yaml
url: "localhost:80"
socket: ""
token: ""
timeout: 60
output:
  format: "json"
  force_json: true
  yes: true
```

## JSON Output Format

Every command returns a consistent envelope:

```json
{
  "ok": true,
  "command": "app list",
  "data": { ... },
  "timestamp": "2026-04-05T17:00:00Z"
}
```

On error:

```json
{
  "ok": false,
  "command": "app install",
  "error": {
    "code": "APP_INSTALL_FAILED",
    "message": "Docker compose validation error",
    "details": { ... }
  },
  "timestamp": "2026-04-05T17:00:00Z"
}
```

## Streaming (`--watch`)

Long-running operations support streaming JSON lines:

```json
{"type": "status", "message": "Pulling image nginx:latest"}
{"type": "progress", "current": 2, "total": 5}
{"type": "done", "duration_seconds": 47}
```

## Architecture

```
cmd/casaos-agent/main.go       # Entry point
internal/
  cli/           # Cobra command definitions
  client/        # HTTP / UNIX socket client (TODO)
  output/        # JSON response formatters + streaming
  config/        # Config loading from env + yaml
```

## Pair With

- **[casaos-webhook-emitter](https://github.com/ChonSong/casaos-webhook-emitter)** — Subscribes to CasaOS MessageBus and fans events out as HTTP webhooks to registered agent endpoints.

## License

MIT — See [LICENSE](./LICENSE) (same as upstream CasaOS-CLI)
