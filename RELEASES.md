# Release Process

1. Bump `Version` in `cmd/constants.go`.
2. PR to `main` and merge once tests pass.
3. Create a release on GitHub on `main` with a `v`-prefixed tag (e.g. `v1.3.1`).
4. The Publish workflow handles everything: GoReleaser builds and signs binaries, publishes to Homebrew, Scoop, apt, rpm, and pushes the winget manifest to the fork. The npm job substitutes the version from `cmd/constants.go` into `npm/package.json` and publishes `@apialerts/cli` plus the 6 platform sub-packages via OIDC. The workflow uses the `GORELEASER_TOKEN` repo secret (a PAT) to push to the tap/bucket/fork repos.
5. Manually open the winget PR from `apialerts/winget-pkgs:apialerts-<version>` to `microsoft/winget-pkgs:master` (until `pull_request.enabled: true` in `.goreleaser.yaml`).

`npm/package.json` ships with `0.0.0-dev` placeholders. CI substitutes the real version from `cmd/constants.go` at publish time. Do not change this.
