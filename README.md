# API Alerts • CLI

A command-line interface for [apialerts.com](https://apialerts.com). Send events from your terminal, scripts, and CI/CD pipelines.

## For AI agents and automation

```bash
export APIALERTS_API_KEY="your-api-key"
npx @apialerts/cli send --json -m "Task complete" -c agents
```

No install step, no config file, parseable JSON output, works in any container with Node. See [Machine-readable output](#machine-readable-output) for the JSON schema.

## Installation

Pick the row that matches your situation:

| Your situation | Run this |
|---|---|
| **AI agent / CI runner / container** (recommended) | `npx @apialerts/cli` + set `APIALERTS_API_KEY` |
| macOS | `brew install --cask apialerts` |
| Debian / Ubuntu | apt repo - see below |
| Fedora / RHEL / CentOS | dnf repo - see below |
| Windows 10 / 11 | `winget install APIAlerts.CLI` |
| Windows + Scoop | Scoop - see below |
| Node already installed | `npx @apialerts/cli` or `npm install -g @apialerts/cli` |
| Go developer | `go install github.com/apialerts/cli@latest` |
| Direct download / air-gapped | [Releases page](https://github.com/apialerts/cli/releases) |

### Detailed instructions

#### Homebrew (macOS / Linux)

```bash
brew tap apialerts/tap
brew install --cask apialerts
```

#### apt (Debian / Ubuntu and derivatives)

```bash
curl -fsSL https://apt.apialerts.com/key.gpg | sudo gpg --dearmor -o /usr/share/keyrings/apialerts.gpg
echo "deb [signed-by=/usr/share/keyrings/apialerts.gpg] https://apt.apialerts.com stable main" | sudo tee /etc/apt/sources.list.d/apialerts.list
sudo apt update && sudo apt install apialerts
```

#### dnf (Fedora / RHEL / CentOS)

```bash
sudo rpm --import https://rpm.apialerts.com/key.gpg
sudo tee /etc/yum.repos.d/apialerts.repo <<EOF
[apialerts]
name=API Alerts
baseurl=https://rpm.apialerts.com
enabled=1
gpgcheck=1
gpgkey=https://rpm.apialerts.com/key.gpg
EOF
sudo dnf install apialerts
```

#### winget (Windows)

```powershell
winget install APIAlerts.CLI
```

#### Scoop (Windows)

```bash
scoop bucket add apialerts https://github.com/apialerts/scoop-bucket
scoop install apialerts
```

#### npm / npx (any platform with Node 18+)

```bash
# One-off, zero install (npx fetches and caches on demand)
npx @apialerts/cli send -m "Deploy complete"

# Global install
npm install -g @apialerts/cli
apialerts send -m "Deploy complete"

# Project devDependency
npm install -D @apialerts/cli
npx apialerts send -m "Deploy complete"
```

Install once with npm, run anywhere with npx. Same binary as Homebrew/winget/apt, shipped through the package manager your stack probably already uses.

#### Go Install

```bash
go install github.com/apialerts/cli@latest
```

#### Download Binary

Download the latest binary from the [Releases](https://github.com/apialerts/cli/releases) page.

## Setup

You'll need an API key from your workspace. After logging in to [apialerts.com](https://apialerts.com), navigate to your workspace and open the **API Keys** section. You can also find it in the mobile app under your workspace settings.

The CLI resolves the API key in this order: `--key` flag, then `APIALERTS_API_KEY` env var, then the saved config file at `~/.apialerts/config.json`.

### Interactive (recommended)

```bash
apialerts init
```

You will be prompted to paste your API key (input is hidden):

```
Enter your API key:
API key saved: abcdef...wxyz
```

### Non-interactive (CI/CD or scripts)

```bash
apialerts config --key "your-api-key"
```

### Environment variable (ephemeral runners, AI agents)

```bash
export APIALERTS_API_KEY="your-api-key"
apialerts send -m "Build complete"
```

The env var is the recommended path for ephemeral environments (Docker containers, sandboxed agents, GitHub Actions runners) where writing to `~/.apialerts/config.json` doesn't persist or isn't appropriate.

### View your current key

```bash
apialerts config
```

```
API Key: abcdef...wxyz
```

### Remove your key

```bash
apialerts config --unset
```

## Send Events

```bash
apialerts send -m "Deploy completed"
```

```bash
apialerts send -e "user.purchase" -t "New Sale" -m "$49.99 from john@example.com" -c "payments"
```

```bash
apialerts send -e "user.signup" -m "New user registered" -d '{"plan":"pro","source":"organic"}'
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--message` | `-m` | Event message **(required)** |
| `--event` | `-e` | Event name for routing (e.g. `user.purchase`) |
| `--title` | `-t` | Event title |
| `--channel` | `-c` | Target channel (uses your default channel if not set) |
| `--tags` | `-g` | Comma-separated tags (e.g. `billing,error`) |
| `--link` | `-l` | Associated URL |
| `--data` | `-d` | JSON object with additional event data (e.g. `'{"plan":"pro"}'`) |
| `--key` | | API key override (highest priority, beats env var and config) |
| `--debug` | | Print request URL, integration, and payload to stderr before sending |
| `--json` | | Emit machine-readable JSON instead of human-readable output |

## Machine-readable output

Pass `--json` on any command to switch output to a single-line JSON object. Successful sends print to stdout, errors print to stderr (and exit 1). This is intended for shell scripts, CI pipelines, and AI agents that need to parse the result.

```bash
apialerts send --json -m "Deploy complete"
# stdout: {"workspace":"Acme Corp","channel":"general","warnings":[]}

apialerts send --json -m "hi" --key bad-key
# stderr: {"error":"failed to send: unauthorized, check your API key"}
# exit code: 1
```

`--debug` output (when enabled) still goes to stderr as `[debug] ...` lines and does not interfere with the JSON payload on stdout.

## Test Connectivity

Send a test event to verify your API key and connection:

```bash
apialerts test
```

```
✓ Test event sent to My Workspace (general)
```

`test` also accepts `--key` to verify a specific key without writing it to config, and `--debug` for verbose output.

## Examples

### Claude Code

Because the CLI is installed on your machine, Claude Code can run it directly as part of any task. Just ask:

- "Refactor the auth module and send me an API Alert when you're done."
- "Run the full test suite and notify me via API Alerts with a summary of the results."
- "Migrate the database schema and send me an apialert if anything fails."

Claude will run `apialerts send` at the right moment - no extra configuration needed.

For agents running in ephemeral environments, set `APIALERTS_API_KEY` in the environment and use `--json` to get a parseable result back:

```bash
APIALERTS_API_KEY="..." apialerts send --json -m "Task complete" -c agents
```

### Shell Script

```bash
#!/bin/bash
apialerts send -m "Backup completed" -c "ops" -g "backup,cron"
```
