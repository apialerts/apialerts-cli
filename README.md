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

Configure your API key once. The key is stored in `~/.apialerts/config.json`.

```bash
apialerts config --key your_api_key
```

Verify your configuration

```bash
apialerts config
```

## Send Events

Send an event with a message

```bash
apialerts send -m "Deploy completed"
```

Send an event with a name and title

```bash
apialerts send -e user.purchase -t "New Sale" -m "$49.99 from john@example.com" -c payments
```

### Optional Properties

You can optionally specify an event name, title, channel, tags, and a link.

```bash
apialerts send -m "Payment failed" -c payments -g billing,error -l https://dashboard.example.com
```

| Flag | Short | Description |
|------|-------|-------------|
| `--message` | `-m` | Event message (required) |
| `--event` | `-e` | Event name for routing (optional, e.g. `user.purchase`) |
| `--title` | `-t` | Event title (optional) |
| `--channel` | `-c` | Target channel (optional, uses default channel if not set) |
| `--tags` | `-g` | Comma-separated tags (optional) |
| `--link` | `-l` | Associated URL (optional) |
| `--key` | | API key override (optional, uses stored config if not set) |

### Override API Key

You can override the stored API key for a single request.

```bash
apialerts send -m "Hello World" --key other_api_key
```

## Test Connectivity

Send a test event to verify your API key and connection.

```bash
apialerts test
```

## CI/CD Examples

### GitHub Actions

```yaml
- name: Send deploy alert
  run: |
    apialerts send -m "Deployed ${{ github.sha }}" -c deployments -g ci,deploy --key ${{ secrets.APIALERTS_API_KEY }}
```

### Shell Script

```bash
#!/bin/bash
apialerts send -m "Backup completed" -c ops -g backup,cron
```
