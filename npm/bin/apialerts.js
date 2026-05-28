#!/usr/bin/env node
'use strict';

const { spawnSync } = require('child_process');
const { resolveBinaryPath } = require('..');

let binary;
try {
  binary = resolveBinaryPath();
} catch (err) {
  process.stderr.write(`${err.message}\n`);
  process.exit(1);
}

const result = spawnSync(binary, process.argv.slice(2), {
  stdio: 'inherit',
});

if (result.error) {
  process.stderr.write(`Failed to execute apialerts: ${result.error.message}\n`);
  process.exit(1);
}

process.exit(result.status ?? 1);
