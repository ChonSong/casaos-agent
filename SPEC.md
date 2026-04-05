# CasaOS Agent — Specification

> Forked from [IceWhaleTech/CasaOS-CLI](https://github.com/IceWhaleTech/CasaOS-CLI)  
> GitHub: https://github.com/ChonSong/casaos-agent

## Overview

**What it is:** A CLI tool for managing a CasaOS instance, designed for autonomous agent operation. Fork of the official `CasaOS-CLI` with machine-readable output, non-interactive operation, and webhook registration built in.

**What it replaces:** The existing `CasaOS-CLI` which outputs human-formatted text and has interactive confirmation prompts that block agents.

**What it pairs with:** `casaos-webhook-emitter` — a separate sidecar service that subscribes to CasaOS MessageBus events and fans them out as HTTP webhooks to registered agent endpoints.

---

## Goals

1. **Machine-readable output** — All commands return JSON. No human-formatted table output.
2. **Non-blocking operation** — No interactive prompts. All confirmation flags are `--yes` or default-allow.
3. **Webhook registration** — Agents can register webhook URLs against event types.
4. **Streaming output** — Long-running operations support `--watch` for real-time progress.
5. **Complete event coverage** — All CasaOS-MessageBus event types documented and subscribeable.

---

## Command Structure

```
casaos-agent [global flags] <command> <subcommand> [flags]
```

### Global Flags
```
--json              Force JSON output
--yaml              Force YAML output
--output format     Output format (json, yaml, table)
--socket path       Connect via UNIX socket
--url url           Connect via HTTP (default localhost:80)
--token string      Authorization Bearer token
--watch             Stream output for long-running operations
--timeout duration  Per-command timeout (default 60s)
--yes               Skip all confirmation prompts
```

### `casaos-agent app`
Compose app lifecycle (AppManagement v2 API: `/v2/apps/...`)

```
casaos-agent app list                          # GET /v2/apps/my (installed apps)
casaos-agent app list --store                  # GET /v2/apps/compose-app-store/info-list (catalog)
casaos-agent app inspect <name>                # GET /v2/apps/my/{name}
casaos-agent app install <name> [--dry-run] [--watch]   # POST /v2/apps/my/install
casaos-agent app uninstall <name> [--purge] [--yes]    # DELETE /v2/apps/my/{name}
casaos-agent app start <name>                  # POST /v2/apps/my/{name}/start
casaos-agent app stop <name>                   # POST /v2/apps/my/{name}/stop
casaos-agent app restart <name>                # POST /v2/apps/my/{name}/restart
casaos-agent app update <name> [--check-only]  # POST /v2/apps/my/{name}/update
casaos-agent app logs <name> [--tail N]        # GET /v2/apps/my/{name}/logs
casaos-agent app resources <name>              # GET /v2/apps/my/{name}/containers (stats)
casaos-agent app config <name>                # GET /v2/apps/my/{name}/settings
casaos-agent app ports <name>                 # GET /v2/apps/my/{name}/ports
```

### `casaos-agent container`
Raw Docker container management (AppManagement v1 API: `/v1/container`)

```
casaos-agent container list                    # GET /v1/container (raw Docker containers)
casaos-agent container inspect <id>          # GET /v1/container/{id}/json
casaos-agent container exec <id> <cmd...>    # POST /v1/container/{id}/exec
casaos-agent container logs <id> [--tail N]   # GET /v1/container/{id}/logs
casaos-agent container stats <id> [--watch]  # GET /v1/container/{id}/stats
casaos-agent container create <compose>       # POST /v1/container/compose
casaos-agent container recreate <id>          # POST /v2/apps/recreate-container/{id}
```

### `casaos-agent system`
System information and management (CasaOS main daemon: `/v1/sys/...`)

```
casaos-agent system info                      # GET /v1/sys/version/current
casaos-agent system resources                # GET /v1/sys/hardware (cpu, mem, disk, network)
casaos-agent system utilization [--watch]     # GET /v1/sys/utilization (live stats)
casaos-agent system health                    # GET /v2/health/services
casaos-agent system logs [--tail N] [--level debug|info|warn|error]
casaos-agent system update [--check-only]    # GET /v1/sys/version/check → POST /v1/sys/update
casaos-agent system restart                  # POST /v1/sys/stop
casaos-agent system reboot                   # POST /v1/sys/reboot
```

### `casaos-agent storage`
Local storage management (CasaOS LocalStorage service)

```
casaos-agent storage list                    # GET /v1/storage (mounts + disks)
casaos-agent storage info <mount>            # Space usage, inode counts
casaos-agent storage mount <device> <path>   # POST /v1/storage/mount
casaos-agent storage unmount <path> [--yes]  # DELETE /v1/storage/{mount}
```

### `casaos-agent webhook`
Webhook registration and management.

```
casaos-agent webhook register <url> [--event <name>]  # Register webhook
casaos-agent webhook list                     # List registered webhooks
casaos-agent webhook deregister <id|url>     # Remove a webhook
casaos-agent webhook test <id>               # Send test payload
casaos-agent webhook history <id>            # Recent delivery attempts
```

### `casaos-agent event`
MessageBus event interaction (read from CasaOS-MessageBus).

```
casaos-agent event list-types                # GET /v2/message_bus/event_type
casaos-agent event subscribe <event-name>    # WebSocket stream to stdout
casaos-agent event publish <name> [--data]   # POST /v2/message_bus/event/{source_id}/{name}
```

### `casaos-agent gateway`
CasaOS Gateway management.

```
casaos-agent gateway routes                  # List all gateway routes
casaos-agent gateway status                  # GET /v2/gateway/status
```

---

## Complete Event Catalog

All events from CasaOS-MessageBus. Source: `CasaOS-AppManagement/common/message.go` + `CasaOS/common/message.go`.

### App Lifecycle Events (source: `casaos-app-management`)

| Event | Description | Key Properties |
|-------|-------------|----------------|
| `app:install-begin` | App installation started | `app:name`, `app:icon` |
| `app:install-progress` | Installation progress update | `app:name`, `app:icon`, `app:progress`, `app:title` |
| `app:install-end` | Installation completed | `app:name`, `app:icon` |
| `app:install-error` | Installation failed | `app:name`, `app:icon`, `message` |
| `app:uninstall-begin` | Uninstall started | `app:name` |
| `app:uninstall-end` | Uninstall completed | `app:name` |
| `app:uninstall-error` | Uninstall failed | `app:name`, `message` |
| `app:update-begin` | Update started | — |
| `app:update-end` | Update completed | — |
| `app:update-error` | Update failed | `message` |
| `app:apply-changes-begin` | Config changes applied | — |
| `app:apply-changes-end` | Config changes completed | — |
| `app:apply-changes-error` | Config changes failed | `message` |
| `app:start-begin` | App start initiated | `app:name` |
| `app:start-end` | App started | `app:name` |
| `app:start-error` | App start failed | `app:name`, `message` |
| `app:stop-begin` | App stop initiated | `app:name` |
| `app:stop-end` | App stopped | `app:name` |
| `app:stop-error` | App stop failed | `app:name`, `message` |
| `app:restart-begin` | App restart initiated | `app:name` |
| `app:restart-end` | App restarted | `app:name` |
| `app:restart-error` | App restart failed | `app:name`, `message` |

### Docker Image Events (source: `casaos-app-management`)

| Event | Properties |
|-------|------------|
| `docker:image:pull-begin` | `app:name` |
| `docker:image:pull-progress` | `app:name`, `message` |
| `docker:image:pull-end` | `app:name`, `docker:image:updated` |
| `docker:image:pull-error` | `app:name`, `message` |
| `docker:image:remove-begin` | `app:name` |
| `docker:image:remove-end` | `app:name` |
| `docker:image:remove-error` | `app:name`, `message` |

### Docker Container Events (source: `casaos-app-management`)

| Event | Properties |
|-------|------------|
| `docker:container:create-begin` | `docker:container:name` |
| `docker:container:create-end` | `docker:container:id`, `docker:container:name` |
| `docker:container:create-error` | `docker:container:name`, `message` |
| `docker:container:start-begin` | `docker:container:id` |
| `docker:container:start-end` | `docker:container:id` |
| `docker:container:start-error` | `docker:container:id`, `message` |
| `docker:container:stop-begin` | `docker:container:id` |
| `docker:container:stop-end` | `docker:container:id` |
| `docker:container:stop-error` | `docker:container:id`, `message` |
| `docker:container:rename-begin` | `docker:container:id`, `docker:container:name` |
| `docker:container:rename-end` | `docker:container:id`, `docker:container:name` |
| `docker:container:rename-error` | `docker:container:id`, `docker:container:name`, `message` |
| `docker:container:remove-begin` | `docker:container:id` |
| `docker:container:remove-end` | `docker:container:id` |
| `docker:container:remove-error` | `docker:container:id`, `message` |

### System Events (source: `casaos`)

| Event | Description |
|-------|-------------|
| `casaos:system:utilization` | Periodic hardware utilization (CPU, RAM, disk, network) |
| `casaos:file:operate` | File operation (copy, move, delete) |
| `casaos:file:recover` | File recovery event |

---

## Architecture

```
cmd/casaos-agent/main.go
internal/
  cli/           # Cobra command definitions
  client/        # HTTP client (CasaOS API) + WebSocket client (MessageBus)
  output/        # JSON response envelope + streaming
  config/        # Env vars + YAML config
```

---

## API Endpoint Map

| Service | Base | Key Endpoints |
|---------|------|--------------|
| CasaOS main daemon | `:80` via Gateway | `GET /v1/sys/*`, `GET /v1/file/*`, `GET /v1/folder/*` |
| AppManagement | `:port` via Gateway | `GET/POST /v2/apps/*`, `GET/POST /v1/container/*` |
| MessageBus | `:port` via Gateway | `GET /v2/message_bus/event_type`, WebSocket `/v2/message_bus/subscribe/event` |
| LocalStorage | `:port` via Gateway | Storage management endpoints |
| Gateway | reverse proxy | Routes `/v1/`, `/v2/` to appropriate service |

All services register with the Gateway on startup. The CLI hits the Gateway at `localhost:80`.

---

## Output Format

### Success
```json
{
  "ok": true,
  "command": "app list",
  "data": { ... },
  "timestamp": "2026-04-05T17:00:00Z"
}
```

### Error
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

### Streaming (`--watch`)
```json
{"type": "status", "message": "Pulling image..."}
{"type": "progress", "current": 3, "total": 5}
{"type": "done", "duration_seconds": 47}
```

---

## Phase 2: What Needs Wiring (from research)

1. **HTTP client** — Replace `fmt.Println` stubs with real HTTP calls to Gateway (`localhost:80`)
2. **Auth** — Bearer token from `~/.config/casaos/auth/token`; pass via `--token`
3. **AppManagement port** — AppManagement runs on a dynamic port assigned at startup; registered in Gateway
4. **WebSocket for events** — Use `casaos-agent event subscribe` to test real WebSocket connection to MessageBus
5. **--watch streaming** — Wire progress events from AppManagement WebSocket to streaming JSON lines
6. **Interactive blockers** — Remove all `fmt.Scanln` and similar; `--yes` always succeeds

---

## File Structure

```
casaos-agent/
├── cmd/casaos-agent/main.go
├── internal/
│   ├── cli/
│   │   ├── root.go         # Root command + global flags
│   │   ├── app.go          # app lifecycle commands
│   │   ├── container.go    # docker container commands
│   │   ├── system.go       # system info commands
│   │   ├── storage.go      # storage commands
│   │   ├── webhook.go      # webhook registry
│   │   ├── event.go        # MessageBus interaction
│   │   └── gateway.go      # gateway commands
│   ├── client/
│   │   ├── http.go         # HTTP client for CasaOS API
│   │   └── websocket.go    # WebSocket client for MessageBus
│   ├── output/
│   │   └── response.go     # JSON envelope, streaming
│   └── config/
│       └── config.go       # Config loading
├── go.mod
├── Makefile
├── README.md
└── SPEC.md
```
