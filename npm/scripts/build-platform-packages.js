#!/usr/bin/env node
// Builds the six per-platform npm packages from the GoReleaser dist/ directory.
//
// Inputs (from project root):
//   dist/cli_<os>_<arch>(_<variant>)?/apialerts(.exe)?   GoReleaser binaries
//
// Outputs:
//   npm/platforms/<os>-<arch>/
//     package.json
//     bin/apialerts(.exe)
//
// Usage:
//   node npm/scripts/build-platform-packages.js [--version <semver>]
//
// If --version is omitted, the version is read from npm/package.json so the
// wrapper and platform packages always publish at the same version.
'use strict';

const fs = require('fs');
const path = require('path');

const ROOT = path.resolve(__dirname, '..', '..');
const NPM_DIR = path.resolve(__dirname, '..');
const DIST_DIR = path.join(ROOT, 'dist');
const OUT_DIR = path.join(NPM_DIR, 'platforms');

// GoReleaser builds these by default for the .goreleaser.yaml in this repo.
// Keep in sync with build matrix in .goreleaser.yaml and PACKAGES in index.js.
const TARGETS = [
  { os: 'darwin', arch: 'arm64', goDist: 'cli_darwin_arm64', binary: 'apialerts' },
  { os: 'darwin', arch: 'x64', goDist: 'cli_darwin_amd64_v1', binary: 'apialerts' },
  { os: 'linux', arch: 'arm64', goDist: 'cli_linux_arm64_v8.0', binary: 'apialerts' },
  { os: 'linux', arch: 'x64', goDist: 'cli_linux_amd64_v1', binary: 'apialerts' },
  { os: 'win32', arch: 'arm64', goDist: 'cli_windows_arm64', binary: 'apialerts.exe' },
  { os: 'win32', arch: 'x64', goDist: 'cli_windows_amd64_v1', binary: 'apialerts.exe' },
];

function readWrapperVersion() {
  const pkg = JSON.parse(
    fs.readFileSync(path.join(NPM_DIR, 'package.json'), 'utf8')
  );
  return pkg.version;
}

function parseArgs() {
  const args = process.argv.slice(2);
  let version = null;
  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--version' && args[i + 1]) {
      version = args[i + 1];
      i++;
    }
  }
  return { version: version ?? readWrapperVersion() };
}

function findBinary(target) {
  // GoReleaser's directory naming includes an arch variant suffix (_v1, _v8.0)
  // that can drift with releaser version. Tolerate either explicit form or
  // the simpler form by globbing.
  const exact = path.join(DIST_DIR, target.goDist, target.binary);
  if (fs.existsSync(exact)) return exact;

  const matches = fs
    .readdirSync(DIST_DIR)
    .filter((name) => {
      const prefix = `cli_${target.os === 'win32' ? 'windows' : target.os}_${target.arch === 'x64' ? 'amd64' : target.arch}`;
      return name.startsWith(prefix);
    })
    .map((name) => path.join(DIST_DIR, name, target.binary))
    .filter((p) => fs.existsSync(p));
  if (matches.length === 0) {
    throw new Error(
      `No GoReleaser output found for ${target.os}-${target.arch}. ` +
        `Expected dist/cli_${target.os === 'win32' ? 'windows' : target.os}_${target.arch === 'x64' ? 'amd64' : target.arch}*/.../${target.binary}. ` +
        `Run 'goreleaser release --snapshot --clean' or check .goreleaser.yaml.`
    );
  }
  return matches[0];
}

function writePlatformPackage(target, version) {
  const pkgName = `@apialerts/cli-${target.os}-${target.arch}`;
  const dir = path.join(OUT_DIR, `${target.os}-${target.arch}`);
  const binDir = path.join(dir, 'bin');
  fs.mkdirSync(binDir, { recursive: true });

  const sourceBinary = findBinary(target);
  const destBinary = path.join(binDir, target.binary);
  fs.copyFileSync(sourceBinary, destBinary);
  // Preserve executable bit on POSIX targets; Windows ignores chmod.
  if (target.os !== 'win32') {
    fs.chmodSync(destBinary, 0o755);
  }

  // npm's `os`/`cpu` fields cause non-matching platforms to skip install,
  // even though we list every package in optionalDependencies.
  const npmOs = target.os; // matches process.platform values
  const npmCpu = target.arch; // matches process.arch values

  const pkg = {
    name: pkgName,
    version,
    description: `API Alerts CLI binary for ${target.os}-${target.arch}`,
    license: 'MIT',
    author: 'API Alerts',
    homepage: 'https://apialerts.com',
    repository: {
      type: 'git',
      url: 'git+https://github.com/apialerts/cli.git',
    },
    bin: {
      apialerts: `bin/${target.binary}`,
    },
    files: ['bin'],
    os: [npmOs],
    cpu: [npmCpu],
  };
  fs.writeFileSync(
    path.join(dir, 'package.json'),
    JSON.stringify(pkg, null, 2) + '\n'
  );

  return { name: pkgName, dir };
}

function main() {
  const { version } = parseArgs();
  if (!fs.existsSync(DIST_DIR)) {
    console.error(
      `dist/ not found at ${DIST_DIR}. Run goreleaser first ` +
        `(e.g. 'goreleaser release --snapshot --clean').`
    );
    process.exit(1);
  }
  fs.rmSync(OUT_DIR, { recursive: true, force: true });

  const built = TARGETS.map((t) => writePlatformPackage(t, version));
  console.log(`Built ${built.length} platform packages at version ${version}:`);
  for (const b of built) {
    console.log(`  ${b.name}  ->  ${path.relative(ROOT, b.dir)}`);
  }
}

main();
