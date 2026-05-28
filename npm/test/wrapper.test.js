'use strict';

const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { spawnSync } = require('node:child_process');

const { PACKAGES, platformKey, binaryName, resolveBinaryPath } = require('..');

test('PACKAGES covers every supported os/arch combination', () => {
  const expected = [
    'darwin-arm64',
    'darwin-x64',
    'linux-arm64',
    'linux-x64',
    'win32-arm64',
    'win32-x64',
  ];
  assert.deepEqual(Object.keys(PACKAGES).sort(), expected.sort());
  for (const key of expected) {
    assert.equal(PACKAGES[key], `@apialerts/cli-${key}`);
  }
});

test('platformKey reflects Node os.platform/os.arch', () => {
  assert.equal(platformKey(), `${os.platform()}-${os.arch()}`);
});

test('binaryName is apialerts.exe on win32, apialerts elsewhere', () => {
  const expected = os.platform() === 'win32' ? 'apialerts.exe' : 'apialerts';
  assert.equal(binaryName(), expected);
});

test('resolveBinaryPath throws a helpful error when sub-package is absent', () => {
  // In a fresh checkout the @apialerts/cli-* packages are not installed, so
  // require.resolve must fail and we expect the wrapped error message.
  assert.throws(() => resolveBinaryPath(), /Failed to locate @apialerts\/cli-/);
});

test('resolveBinaryPath returns the binary when a sub-package is installed', (t) => {
  // Simulate the per-platform package being installed under node_modules by
  // creating a fake one matching the current host. The shim should locate it.
  const key = platformKey();
  const pkgName = PACKAGES[key];
  // realpathSync to normalise the /private/var/folders vs /var/folders prefix
  // on macOS so require.resolve's output matches.
  const fakeRoot = fs.realpathSync(
    fs.mkdtempSync(path.join(os.tmpdir(), 'apialerts-cli-test-'))
  );
  t.after(() => fs.rmSync(fakeRoot, { recursive: true, force: true }));

  // node_modules/@apialerts/cli-<platform>/bin/<binary>
  const pkgDir = path.join(fakeRoot, 'node_modules', pkgName);
  const binDir = path.join(pkgDir, 'bin');
  fs.mkdirSync(binDir, { recursive: true });
  const binPath = path.join(binDir, binaryName());
  fs.writeFileSync(binPath, '#!/bin/sh\nexit 0\n', { mode: 0o755 });
  fs.writeFileSync(
    path.join(pkgDir, 'package.json'),
    JSON.stringify({ name: pkgName, version: '0.0.0-test', main: 'bin/' + binaryName() })
  );

  // Re-require the module with NODE_PATH pointing at our fake tree.
  const env = { ...process.env, NODE_PATH: path.join(fakeRoot, 'node_modules') };
  const probe = `
    require('module').Module._initPaths();
    const { resolveBinaryPath } = require(${JSON.stringify(path.resolve(__dirname, '..'))});
    process.stdout.write(resolveBinaryPath());
  `;
  const result = spawnSync(process.execPath, ['-e', probe], { env, encoding: 'utf8' });
  assert.equal(result.status, 0, `probe failed: stderr=${result.stderr}`);
  assert.equal(path.normalize(result.stdout), path.normalize(binPath));
});

test('bin/apialerts.js forwards exit code from the underlying binary', (t) => {
  // Stage a fake platform package under a temp dir and run the shim with
  // NODE_PATH pointing at it. The fake binary exits 42, which the shim should
  // surface as its own exit status.
  if (os.platform() === 'win32') {
    // The fake-binary shell script path doesn't apply on Windows runners;
    // skip rather than try to emulate cmd.exe semantics here.
    t.skip('shim exit-code test skipped on win32');
    return;
  }
  const key = platformKey();
  const pkgName = PACKAGES[key];
  const fakeRoot = fs.mkdtempSync(path.join(os.tmpdir(), 'apialerts-cli-shim-'));
  t.after(() => fs.rmSync(fakeRoot, { recursive: true, force: true }));

  const binDir = path.join(fakeRoot, 'node_modules', pkgName, 'bin');
  fs.mkdirSync(binDir, { recursive: true });
  const binPath = path.join(binDir, binaryName());
  fs.writeFileSync(binPath, '#!/bin/sh\nexit 42\n', { mode: 0o755 });
  fs.writeFileSync(
    path.join(fakeRoot, 'node_modules', pkgName, 'package.json'),
    JSON.stringify({ name: pkgName, version: '0.0.0-test' })
  );

  const shim = path.resolve(__dirname, '..', 'bin', 'apialerts.js');
  const env = { ...process.env, NODE_PATH: path.join(fakeRoot, 'node_modules') };
  const result = spawnSync(process.execPath, [shim, '--ignored'], { env });
  assert.equal(result.status, 42, `expected exit 42, got ${result.status}`);
});
