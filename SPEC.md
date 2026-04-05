# CasaOS Agent — Specification

> Forked from [IceWhaleTech/CasaOS-CLI](https://github.com/IceWhaleTech/CasaOS-CLI)  
> GitHub: https://github.com/ChonSong/casaos-agent

## Overview

**What it is:** A CLI tool for managing a CasaOS instance, designed for autonomous agent operation. Fork of the official `CasaOS-CLI` with machine-readable output, non-interactive operation, and webhook registration built in.

**What it replaces:** The existing `CasaOS-CLI` which outputs human-formatted text and has interactive confirmation prompts that block agents.

**What it pairs with:** `casaos-webhook-emitter` — a separate sidecar service that subscribes to CasaOS MessageBus events and fans them out as HTTP webhooks to registered agent endpoints.

---

## Goals

1. **Machine-readable output** — All commands that query state return JSON. No human-formatted table output unless explicitly requested.
2. **Non-blocking operation** — No interactive prompts. All confirmation flags are `--yes` or default-allow.
3. **Webhook registration** — Agents can register webhook URLs against event types without leaving the CLI.
4. **Streaming output** — Long-running operations (app install, container start) support `--watch` and `--json` streaming output so agents can poll for completion.
5. **UNIX socket transport** — CLI can connect via UNIX socket (`/run/casaos-agent.sock`) in addition to HTTP, for agent-native local communication.

---

## Command Structure (Planned)

```
casaos-agent [global flags] <command> <subcommand> [flags]
```

### Global Flags
```
--json              Force JSON output for all commands
--yaml              Force YAML output
--output format     Alias for --json/--yaml
--socket path       Connect via UNIX socket (default /run/casaos-agent.sock)
--url url           Connect via HTTP (default localhost:80)
--token string      API token / Authorization Bearer
--watch             Stream output for long-running operations
--timeout duration  Per-command timeout (default 60s)
--yes               Skip all confirmation prompts (auto-confirm)
```

### Command Tree

#### `casaos-agent app`
Compose app lifecycle management.

```
casaos-agent app list                        # List all installed apps
casaos-agent app list --store                # List available store apps
casaos-agent app install <name|url> [--dry-run] [--watch]   # Install from store or compose URL
casaos-agent app uninstall <name> [--yes] [--purge]         # Remove app and optionally volumes
casaos-agent app start <name>                              # Start app container(s)
casaos-agent app stop <name>                               # Stop app container(s)
casaos-agent app restart <name>                           # Restart app container(s)
casaos-agent app inspect <name>                            # Full app state as JSON
casaos-agent app logs <name> [--tail N] [--since duration] # Stream logs as JSON lines
casaos-agent app update <name> [--check-only]              # Check / apply updates
casaos-agent app resources <name>                          # CPU, memory, I/O usage
```

#### `casaos-agent container`
Raw Docker container management.

```
casaos-agent container list                  # List all Docker containers
casaos-agent container inspect <id>          # Full container state as JSON
casaos-agent container exec <id> <cmd>       # Run command inside container
casaos-agent container stats <id>            # Live resource stats as JSON
```

#### `casaos-agent system`
System information and health.

```
casaos-agent system info                     # Hostname, OS, kernel, uptime
casaos-agent system resources                # CPU, memory, disk, network — full JSON
casaos-agent system utilization              # Live utilization (same as /v1/sys/utilization)
casaos-agent system health                   # Service health for all casaos-* services
casaos-agent system logs [--tail N] [--level info|warn|error]
casaos-agent system update [--check-only] [--version x.y.z]
casaos-agent system restart                  # Restart CasaOS service
casaos-agent system reboot                   # Reboot host
```

#### `casaos-agent storage`
Local storage management.

```
casaos-agent storage list                    # All mount points and disks
casaos-agent storage info <mount>            # Space usage, inode counts
casaos-agent storage mount <device> <path>
casaos-agent storage unmount <path> [--yes]
```

#### `casaos-agent webhook`
Webhook registration and management.

```
casaos-agent webhook register <url> [--event <event-name>] [--event <event-name>]  # Register a webhook
casaos-agent webhook list                     # List registered webhooks
casaos-agent webhook deregister <id|url>     # Remove a webhook
casaos-agent webhook test <id>               # Send a test payload to the webhook
casaos-agent webhook history <id>            # Recent delivery attempts and status
```

#### `casaos-agent event`
Event bus interaction (reads from CasaOS-MessageBus).

```
casaos-agent event list-types                # All registered event types in MessageBus
casaos-agent event subscribe <event-name>    # Subscribe to an event type (stdout streaming)
casaos-agent event publish <name> [--data <json>]  # Publish an event (for testing)
```

#### `casaos-agent gateway`
CasaOS Gateway management.

```
casaos-agent gateway routes                  # List all gateway routes
casaos-agent gateway status                  # Gateway health
```

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
    "message": "Docker compose failed: network my-net not found",
    "details": { ... }
  },
  "timestamp": "2026-04-05T17:00:00Z"
}
```

### Stream (--watch)
Each line is a JSON object:
```json
{"type": "status", "message": "Pulling image..."}
{"type": "status", "message": "Creating network..."}
{"type": "progress", "current": 3, "total": 5}
{"type": "done", "duration_seconds": 47}
```

---

## Event Catalog

*(Populated by researcher from CasaOS-MessageBus OpenAPI spec)*

Expected event types from CasaOS-MessageBus:

| Event Name | Source | Description |
|-----------|--------|-------------|
| `casaos:system:utilization` | CasaOS daemon | Periodic hardware utilization |
| `casaos:file:operate` | CasaOS daemon | File operations (copy, move, delete) |
| `casaos:file:recover` | CasaOS daemon | File recovery events |
| *(to be cataloged)* | AppManagement | App install/start/stop/uninstall |
| *(to be cataloged)* | AppManagement | Container state changes |

---

## Authentication

- **Default:** No auth for local CLI use (UNIX socket + same host = trusted)
- **Remote:** `--token` flag sets `Authorization: Bearer <token>` header
- **MessageBus:** Token-based auth; token stored in CasaOS config at `~/.config/casaos/auth/token`

---

## Non-Functional Requirements

1. **Single binary** — No external dependencies, fully static Go build
2. **Cross-platform** — Linux (amd64 + arm64), macOS
3. **Shell completions** — bash, zsh, fish available
4. **Man pages** — Generated via cobra doc
5. **Test coverage** — All command handlers have unit tests; integration tests hit a real CasaOS instance (or mock)

---

## File Structure

```
casaos-agent/
├── cmd/
│   └── casaos-agent/
│       └── main.go
├── internal/
│   ├── cli/           # Cobra command definitions
│   ├── client/        # HTTP/UNIX socket client wrappers
│   ├── output/        # JSON/YAML output formatters
│   └── config/        # Config file loading
├── CasaOS-CLI/        # Upstream fork (git submodule or bare clone)
│   └── ...            # Original source, minimally modified
├── Makefile
├── README.md
└── SPEC.md
```

---

## Relationship to Upstream

- Forked from `IceWhaleTech/CasaOS-CLI` commit `HEAD` as of 2026-04-05
- Additions: `--json`/`--yaml` flags, `--watch` streaming, `--yes` confirm bypass, `webhook` command group
- Modifications: Replace all `fmt.Printf` table output with structured JSON
- Upstream merge strategy: Rebase onto upstream main; conflicts expected in `cmd/` tree for any upstream command changes

---

## Out of Scope (Phase 1)

- Native agent protocol (MCP, etc.) — agents use CLI subprocess or HTTP
- Webhook emitter sidecar (separate repo: `casaos-webhook-emitter`)
- Installation/installer modifications
- CasaOS-UI changes
