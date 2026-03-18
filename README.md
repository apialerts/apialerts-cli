# API Alerts • CLI

A command-line interface for [apialerts.com](https://apialerts.com). Send events from your terminal, scripts, and CI/CD pipelines.

## Installation

### Homebrew (macOS / Linux)

```bash
brew tap apialerts/tap
brew install --cask apialerts
```

### Scoop (Windows)

```bash
scoop bucket add apialerts https://github.com/apialerts/scoop-bucket
scoop install apialerts
```

### Go Install

```bash
go install github.com/apialerts/cli@latest
```

### Download Binary

Download the latest binary from the [Releases](https://github.com/apialerts/cli/releases) page.

## Setup

You'll need an API key from your workspace. After logging in to [apialerts.com](https://apialerts.com), navigate to your workspace and open the **API Keys** section. You can also find it in the mobile app under your workspace settings.

Your API key is stored locally in `~/.apialerts/config.json`.

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
| `--key` | | API key override (uses stored config if not set) |

## Test Connectivity

Send a test event to verify your API key and connection:

```bash
apialerts test
```

```
✓ Test event sent to My Workspace (general)
```

## Examples

### Claude Code

Because the CLI is installed on your machine, Claude Code can run it directly as part of any task. Just ask:

- "Refactor the auth module and send me an API Alert when you're done."
- "Run the full test suite and notify me via API Alerts with a summary of the results."
- "Migrate the database schema and send me an apialert if anything fails."

Claude will run `apialerts send` at the right moment — no extra configuration needed.

### Shell Script

```bash
#!/bin/bash
apialerts send -m "Backup completed" -c "ops" -g "backup,cron"
```
