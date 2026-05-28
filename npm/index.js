'use strict';

const os = require('os');

// Map Node's platform+arch to the npm sub-package that ships the right binary.
// Keep in sync with optionalDependencies in package.json and the goreleaser
// build matrix (.goreleaser.yaml).
const PACKAGES = {
  'darwin-arm64': '@apialerts/cli-darwin-arm64',
  'darwin-x64': '@apialerts/cli-darwin-x64',
  'linux-arm64': '@apialerts/cli-linux-arm64',
  'linux-x64': '@apialerts/cli-linux-x64',
  'win32-arm64': '@apialerts/cli-win32-arm64',
  'win32-x64': '@apialerts/cli-win32-x64',
};

function platformKey() {
  return `${os.platform()}-${os.arch()}`;
}

function binaryName() {
  return os.platform() === 'win32' ? 'apialerts.exe' : 'apialerts';
}

function resolveBinaryPath() {
  const key = platformKey();
  const pkg = PACKAGES[key];
  if (!pkg) {
    throw new Error(
      `@apialerts/cli does not ship a binary for ${key}. ` +
        `Supported platforms: ${Object.keys(PACKAGES).join(', ')}.`
    );
  }

  try {
    return require.resolve(`${pkg}/bin/${binaryName()}`);
  } catch (err) {
    throw new Error(
      `Failed to locate ${pkg}. This usually means npm skipped the platform ` +
        `package during install (optionalDependencies are best-effort). ` +
        `Try: npm install ${pkg}. Underlying error: ${err.message}`
    );
  }
}

module.exports = { PACKAGES, platformKey, binaryName, resolveBinaryPath };
