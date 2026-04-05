# CasaOS Agent CLI

> **Status:** Phase 2 — Fork complete, real CasaOS-CLI wired with agent-native flags.  
> GitHub: https://github.com/ChonSong/casaos-agent

Agent-native CLI for managing a CasaOS instance. Fork of [IceWhaleTech/CasaOS-CLI](https://github.com/IceWhaleTech/CasaOS-CLI) with machine-readable output, non-interactive flags, and streaming support built in.

## Quick Start

```bash
# Build
cd CasaOS-CLI && go build -o bin/casaos-agent .

# Or install
go install github.com/ChonSong/casaos-agent@latest

# List installed apps (JSON output)
casaos-agent --json app-management list apps

# Install app (non-interactive)
casaos-agent --yes app-management install nginx --file docker-compose.yml

# System health check
casaos-agent --json healthcheck services

# Subscribe to MessageBus events as JSON
casaos-agent --json message-bus subscribe websocket \
  --source-id casaos-app-management
```

## Key Differences from upstream `casaos-cli`

| Feature | `casaos-cli` | `casaos-agent` |
|---------|-------------|----------------|
| Output format | Human-formatted tables | **Structured JSON** (with `--json`) |
| Interactive prompts | Yes (blocking) | **Bypass with `--yes`** |
| Streaming output | No | **`--watch` flag** |
| Error codes | No | Typed error codes in JSON envelope |

## Agent-Native Flags

```
--json, -j     Force structured JSON output for all commands
--yes, -y      Skip all confirmation prompts (auto-confirm)
--watch, -w    Stream output for long-running operations
--url, -u      CasaOS API root URL (default: localhost:80)
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
    "code": "ERROR",
    "message": "404 Not Found - is the casaos-app-management service running?"
  },
  "timestamp": "2026-04-05T17:00:00Z"
}
```

## Command Tree

```
casaos-agent app-management    # Compose app lifecycle
  list apps                    # GET /v2/apps/my
  list app-stores             # GET /v2/apps/app-stores
  search                      # Search compose app store
  install <app>              # POST /v2/apps/my/install
  uninstall <appid>           # DELETE /v2/apps/my/{id}
  start <appid>              # POST /v2/apps/my/{id}/start
  stop <appid>               # POST /v2/apps/my/{id}/stop
  restart <appid>            # POST /v2/apps/my/{id}/restart
  logs <appid>               # GET /v2/apps/my/{id}/logs
  show <appid>               # GET /v2/apps/my/{id}

casaos-agent message-bus      # MessageBus interaction
  list event-types           # GET /v2/message_bus/event_type
  subscribe websocket         # WS /v2/message_bus/subscribe/event/{source_id}
  trigger <name>             # POST /v2/message_bus/action/{source_id}/{name}

casaos-agent local-storage   # Storage management
casaos-agent healthcheck     # Service health
casaos-agent gateway         # Gateway routes
casaos-agent user            # User management
```

## Complete Event Catalog

When subscribing to the MessageBus (`message-bus subscribe websocket`), CasaOS emits:

**App lifecycle** (source: `casaos-app-management`):
`app:install-{begin,progress,end,error}`, `app:uninstall-{begin,end,error}`, `app:update-{begin,end,error}`, `app:apply-changes-{begin,end,error}`, `app:{start,stop,restart}-{begin,end,error}`

**Docker image** (source: `casaos-app-management`):
`docker:image:pull-{begin,progress,end,error}`, `docker:image:remove-{begin,end,error}`

**Docker container** (source: `casaos-app-management`):
`docker:container:{create,start,stop,rename,remove}-{begin,end,error}`

**System** (source: `casaos`):
`casaos:system:utilization`, `casaos:file:operate`, `casaos:file:recover`

## Architecture

```
casaos-agent/
├── CasaOS-CLI/           # Fork of upstream CasaOS-CLI
│   ├── main.go           # Entry point → cmd.Execute()
│   ├── cmd/              # All Cobra commands (forked + modified)
│   │   ├── root.go       # Root cmd + agent-native flags + JSONPrintResponse
│   │   ├── appManagement*.go  # App lifecycle commands
│   │   ├── messageBus*.go     # MessageBus commands
│   │   └── ...
│   ├── codegen/          # Generated API clients from OpenAPI specs
│   │   ├── app_management/
│   │   ├── message_bus/
│   │   ├── casaos/
│   │   ├── local_storage/
│   │   └── user_service/
│   └── go.mod
├── go.mod                # Top-level module
└── README.md
```

## Pair With

- **[casaos-webhook-emitter](https://github.com/ChonSong/casaos-webhook-emitter)** — Subscribes to CasaOS MessageBus and fans events out as HTTP webhooks to registered agent endpoints.
