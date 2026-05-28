#!/usr/bin/env node
// Reads the version from cmd/constants.go, substitutes it into
// npm/package.json (replacing the 0.0.0-dev placeholders), and builds the
// six per-platform npm packages from the GoReleaser dist/ directory.
//
// Usage: node npm/scripts/build-platform-packages.js [--version <semver>]
'use strict';

const fs = require('fs');
const path = require('path');

const ROOT = path.resolve(__dirname, '..', '..');
const NPM_DIR = path.resolve(__dirname, '..');
const DIST_DIR = path.join(ROOT, 'dist');
const OUT_DIR = path.join(NPM_DIR, 'platforms');
const CONSTANTS_PATH = path.join(ROOT, 'cmd', 'constants.go');
const WRAPPER_PKG = path.join(NPM_DIR, 'package.json');

const TARGETS = [
  { os: 'darwin', arch: 'arm64', goDist: 'cli_darwin_arm64', binary: 'apialerts' },
  { os: 'darwin', arch: 'x64', goDist: 'cli_darwin_amd64_v1', binary: 'apialerts' },
  { os: 'linux', arch: 'arm64', goDist: 'cli_linux_arm64_v8.0', binary: 'apialerts' },
  { os: 'linux', arch: 'x64', goDist: 'cli_linux_amd64_v1', binary: 'apialerts' },
  { os: 'win32', arch: 'arm64', goDist: 'cli_windows_arm64', binary: 'apialerts.exe' },
  { os: 'win32', arch: 'x64', goDist: 'cli_windows_amd64_v1', binary: 'apialerts.exe' },
];

function readGoVersion() {
  const src = fs.readFileSync(CONSTANTS_PATH, 'utf8');
  const m = src.match(/Version\s*=\s*"([^"]+)"/);
  if (!m) {
    throw new Error(`Could not parse Version from ${CONSTANTS_PATH}`);
  }
  return m[1];
}

function syncWrapperVersion(version) {
  const pkg = JSON.parse(fs.readFileSync(WRAPPER_PKG, 'utf8'));
  pkg.version = version;
  if (pkg.optionalDependencies) {
    for (const key of Object.keys(pkg.optionalDependencies)) {
      pkg.optionalDependencies[key] = version;
    }
  }
  fs.writeFileSync(WRAPPER_PKG, JSON.stringify(pkg, null, 2) + '\n');
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
  return { version: version ?? readGoVersion() };
}

function findBinary(target) {
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
  if (target.os !== 'win32') {
    fs.chmodSync(destBinary, 0o755);
  }

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
    os: [target.os],
    cpu: [target.arch],
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

  syncWrapperVersion(version);
  console.log(`Synced npm/package.json (version + optionalDependencies) to ${version}`);

  const built = TARGETS.map((t) => writePlatformPackage(t, version));
  console.log(`Built ${built.length} platform packages at version ${version}:`);
  for (const b of built) {
    console.log(`  ${b.name}  ->  ${path.relative(ROOT, b.dir)}`);
  }
}

main();
