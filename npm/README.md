# @apialerts/cli

The [API Alerts](https://apialerts.com) CLI, distributed via npm. Wraps the native Go binary so you can install and run it with `npm` or `npx`, no Homebrew/Scoop/apt setup required.

## Install

```bash
# Global install
npm install -g @apialerts/cli
apialerts send -m "Deploy complete"

# Local devDependency
npm install -D @apialerts/cli
npx apialerts send -m "Deploy complete"

# Zero-install (downloads on demand, caches for re-use)
npx @apialerts/cli send -m "Deploy complete"
```

## How it works

This package is a small Node shim. On install, npm picks the right platform-specific sub-package via `optionalDependencies` (`@apialerts/cli-darwin-arm64`, `@apialerts/cli-linux-x64`, etc.) and the shim execs the native binary from it. The same binary you'd get from Homebrew or a GitHub release.

Supported platforms:

| OS | Architectures |
|---|---|
| macOS | arm64, x64 |
| Linux | arm64, x64 |
| Windows | arm64, x64 |

## Authentication

Set `APIALERTS_API_KEY` in your environment - this is the recommended path for any non-interactive use including CI, containers, and AI agents:

```bash
export APIALERTS_API_KEY="your-api-key"
npx apialerts send --json -m "Build complete"
```

For interactive setup, run `apialerts init` (writes to `~/.apialerts/config.json`).

Full docs: <https://apialerts.com/docs/tools/cli>

## Why a separate npm package?

This is the same CLI shipped via Homebrew, Scoop, apt, dnf, and direct download. npm distribution exists so JavaScript/TypeScript projects can pull it as a `devDependency` and AI agents in Node-based environments can `npx` it without a system install. The Go source lives at <https://github.com/apialerts/cli>.

The unscoped `apialerts` package on npm is the [JavaScript SDK](https://www.npmjs.com/package/apialerts), a different thing - that's a library you `import` in code.
