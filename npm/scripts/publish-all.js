#!/usr/bin/env node
// Publishes the wrapper and all built platform packages to npm. Designed to
// run after build-platform-packages.js inside a GitHub Actions job that has
// id-token: write (for OIDC trusted publishing) and a configured registry-url.
//
// Order: platform packages first, then the main wrapper. This guarantees the
// wrapper's optionalDependencies resolve immediately for anyone who hits the
// registry mid-publish.
//
// Usage:
//   node npm/scripts/publish-all.js [--dry-run] [--tag alpha]
'use strict';

const fs = require('fs');
const path = require('path');
const { spawnSync } = require('child_process');

const NPM_DIR = path.resolve(__dirname, '..');
const OUT_DIR = path.join(NPM_DIR, 'platforms');

function parseArgs() {
  const args = process.argv.slice(2);
  const flags = { dryRun: false, tag: null };
  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--dry-run') flags.dryRun = true;
    else if (args[i] === '--tag' && args[i + 1]) {
      flags.tag = args[i + 1];
      i++;
    }
  }
  return flags;
}

function isAlreadyPublished(name, version) {
  // `npm view <name>@<version> version` exits 0 and prints the version when it
  // exists. When the version doesn't exist on npm, exit code is non-zero (404).
  const result = spawnSync('npm', ['view', `${name}@${version}`, 'version'], {
    encoding: 'utf8',
  });
  return result.status === 0 && result.stdout.trim() === version;
}

function publish(dir, flags) {
  const pkg = JSON.parse(fs.readFileSync(path.join(dir, 'package.json'), 'utf8'));

  // Idempotent retry: if this exact version is already on npm, skip rather
  // than letting `npm publish` fail with "cannot publish over the previously
  // published version" partway through the loop. Lets us recover from a
  // mid-publish failure (e.g. trusted-publishing misconfig on one of the 7
  // packages) by fixing the bad one and re-running the workflow.
  if (!flags.dryRun && isAlreadyPublished(pkg.name, pkg.version)) {
    console.log(`Already published ${pkg.name}@${pkg.version}, skipping.`);
    return;
  }

  const args = ['publish', '--access', 'public', '--provenance'];
  if (flags.tag) args.push('--tag', flags.tag);
  if (flags.dryRun) args.push('--dry-run');

  console.log(
    `${flags.dryRun ? '[dry-run] ' : ''}npm publish ${pkg.name}@${pkg.version} in ${path.relative(process.cwd(), dir)}`
  );

  const result = spawnSync('npm', args, { cwd: dir, stdio: 'inherit' });
  if (result.status !== 0) {
    throw new Error(`npm publish failed for ${pkg.name} (exit ${result.status})`);
  }
}

function main() {
  const flags = parseArgs();
  if (!fs.existsSync(OUT_DIR)) {
    console.error(
      `npm/platforms/ not found. Run build-platform-packages.js first.`
    );
    process.exit(1);
  }

  const platformDirs = fs
    .readdirSync(OUT_DIR)
    .map((d) => path.join(OUT_DIR, d))
    .filter((d) => fs.statSync(d).isDirectory())
    .sort();

  for (const dir of platformDirs) {
    publish(dir, flags);
  }
  // Wrapper last so its optionalDependencies resolve straight away.
  publish(NPM_DIR, flags);
  console.log('All packages published.');
}

main();
