```
 ██████╗██████╗  ██████╗ ███████╗███████╗██████╗  █████╗ ████████╗██╗  ██╗
██╔════╝██╔══██╗██╔═══██╗██╔════╝██╔════╝██╔══██╗██╔══██╗╚══██╔══╝██║  ██║
██║     ██████╔╝██║   ██║███████╗███████╗██████╔╝███████║   ██║   ███████║
██║     ██╔══██╗██║   ██║╚════██║╚════██║██╔═══╝ ██╔══██║   ██║   ██╔══██║
╚██████╗██║  ██║╚██████╔╝███████║███████║██║     ██║  ██║   ██║   ██║  ██║
 ╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚══════╝╚═╝     ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝
                              C O M P A N I O N
```

# CrossPath Companion

Connect your local computer to [CrossPath](https://crosspath.im) — your private AI workspace.

The companion runs as a lightweight background daemon on your machine. It bridges your local files, terminal, browser, and MCP tools to the CrossPath web chat so the AI can act on your behalf locally.

---

## Quick Install

**Linux / macOS**
```bash
curl -fsSL https://raw.githubusercontent.com/originlabsai/crosspath-companion/main/scripts/install.sh | sh
```

**Windows (PowerShell)**
```powershell
irm https://raw.githubusercontent.com/originlabsai/crosspath-companion/main/scripts/install.ps1 | iex
```

---

## First-Time Setup

```bash
cross_companion login     # opens browser -> sign in to CrossPath -> saves token
cross_companion start     # starts the daemon and connects to CrossPath
cross_companion status    # verify it's connected
```

**Auto-start on login (optional)**
```bash
cross_companion service install
```

---

## Commands

| Command | Description |
|---|---|
| `cross_companion` | Launch interactive dashboard (default) |
| `cross_companion login` | Authenticate with CrossPath |
| `cross_companion logout` | Sign out and revoke this device |
| `cross_companion start` | Start daemon (headless, no UI) |
| `cross_companion status` | Show connection and auth status |
| `cross_companion list` | List configured MCP servers |
| `cross_companion add <name>` | Add a local MCP server |
| `cross_companion remove <name>` | Remove a local MCP server |
| `cross_companion service install` | Register as system service (auto-start) |
| `cross_companion service uninstall` | Remove system service |
| `cross_companion devices list` | List devices linked to your account |
| `cross_companion devices revoke <id>` | Revoke a device |

---

## What It Can Do

Once connected, CrossPath web chat gains access to your local environment:

| Tool | What it does |
|---|---|
| `read_file` | Read local files with line-range support |
| `write_file` | Create or overwrite local files |
| `string_replace` | Surgical in-place file edits |
| `list_directory` | Browse local folders |
| `find_files` | Recursive file search with glob patterns |
| `grep` | Search file contents |
| `execute_bash` | Run terminal commands |
| `run_background` | Start background processes (dev servers, watchers) |
| `browser` | Control a local browser via MCP |
| `http_request` | Make outbound HTTP calls from your machine |
| `get_device_info` | Report OS and hardware info |

You can also connect any third-party **MCP server** (filesystem, databases, VS Code, Obsidian, etc.) via `cross_companion add`.

---

## File Access Transparency

Every file and tool access is logged. View the full history in your CrossPath dashboard under **Settings > Companion > File Access Log**.

The log shows:
- Which file or command was accessed
- The operation (`read`, `write`, `exec`, `http`, ...)
- Timestamp and success/failure
- Which chat session triggered it

---

## Configuration

Config is stored at `~/.crosspath/mcp-config.yaml`.

Add a custom MCP server:
```bash
cross_companion add my-server \
  --command npx \
  --args @my-org/mcp-server@latest
```

---

## Building from Source

Requires Go 1.22+.

```bash
git clone https://github.com/originlabsai/crosspath-companion.git
cd crosspath-companion

# Build for current platform
go build -o cross_companion ./cmd/cross_companion

# Build all platforms
./build.sh prod
```

---

## Architecture

```
Your Computer                      CrossPath Cloud
--------------------------        --------------------------
cross_companion (daemon)  <--WSS-->  crosspath.im/mcp/connect
  |-- 11 built-in tools              |
  |-- local MCP servers              +-- web chat routes tool
  +-- file access log                    calls here
       ^ streamed to cloud
```

The companion connects to CrossPath over a persistent WebSocket. Tool calls flow from your browser session through the cloud, down to the companion, executed locally, and results returned — all in real time.

---

## Privacy & Security

- Your files never leave your machine unless the AI explicitly reads and returns them as part of a response
- All file accesses are logged locally at `~/.crosspath/access.log` and visible in your dashboard
- The companion only connects outbound — no inbound ports are opened
- Token refresh happens automatically; you only log in once per device
- Revoke any device at any time from `cross_companion devices revoke` or the CrossPath dashboard

---

## Links

- CrossPath: [https://crosspath.im](https://crosspath.im)
- GitHub: [https://github.com/originlabsai/crosspath-companion](https://github.com/originlabsai/crosspath-companion)
