#!/usr/bin/env node
'use strict';

const { spawn } = require('child_process');
const { createHash } = require('crypto');
const fs = require('fs');
const fsp = require('fs/promises');
const https = require('https');
const os = require('os');
const path = require('path');
const zlib = require('zlib');
const pkg = require('../package.json');

const PACKAGE_NAME = pkg.name;
const VERSION = pkg.version;
const DEFAULT_REPOSITORY = 'TencentCloudCommunity/mcp-server';
const DEFAULT_RELEASE_TAG = `postgres-mcp-server-v${VERSION}`;
const CHECKSUM_FILE = 'checksums.txt';
const DEFAULT_CACHE_DIR = path.join(os.homedir(), '.cache', PACKAGE_NAME, VERSION, `${platformName()}-${archName()}`);

const HELP_TEXT = `${PACKAGE_NAME} v${VERSION}\n\n` +
  '用法:\n' +
  '  npx -y postgres-mcp-server@latest\n' +
  '  npx -y postgres-mcp-server@latest --transport sse --env-file .env\n' +
  '  npx -y ./postgres-mcp-server-1.0.3.tgz --release-repository <owner/repo>\n\n' +
  '参数:\n' +
  '  --transport <stdio|streamable-http|sse>  覆盖 MCP_TRANSPORT，默认 stdio\n' +
  '  --env-file <path>                        启动前按 key=value 语义加载 .env\n' +
  '  --binary-path <path>                     直接使用本地二进制，跳过下载\n' +
  '  --release-repository <owner/repo>        覆盖 GitHub Release 仓库，默认官方仓库\n' +
  '  --release-tag <tag>                      覆盖 GitHub Release tag，默认随版本推导\n' +
  '  --help                                   显示帮助\n' +
  '  --version                                输出版本\n\n' +
  '环境变量兼容:\n' +
  '  TRANSPORT -> MCP_TRANSPORT\n' +
  '  PORT -> MCP_SERVER_PORT\n' +
  '  TENCENTCLOUD_SECRET_ID / KEY 仍可直接复用\n' +
  '  POSTGRES_MCP_BINARY_PATH 可指定本地二进制路径\n' +
  '  POSTGRES_MCP_RELEASE_REPOSITORY 可改为 owner/repo\n' +
  '  POSTGRES_MCP_RELEASE_TAG 可覆盖默认 tag';

main().catch((error) => {
  console.error(`[${PACKAGE_NAME}] ${error.message}`);
  process.exit(1);
});

async function main() {
  const options = parseArgs(process.argv.slice(2));
  if (options.help) {
    process.stdout.write(`${HELP_TEXT}\n`);
    return;
  }
  if (options.version) {
    process.stdout.write(`${VERSION}\n`);
    return;
  }

  const env = await buildChildEnv(options);
  const binaryPath = await resolveBinaryPath(options);
  await ensureExecutable(binaryPath);
  await runServer(binaryPath, env);
}

function parseArgs(argv) {
  const options = {
    transport: undefined,
    envFile: undefined,
    binaryPath: undefined,
    releaseRepository: undefined,
    releaseTag: undefined,
    help: false,
    version: false,
  };

  for (let i = 0; i < argv.length; i += 1) {
    const arg = argv[i];
    switch (arg) {
      case '--transport':
        options.transport = requireValue(argv, ++i, '--transport');
        break;
      case '--env-file':
        options.envFile = requireValue(argv, ++i, '--env-file');
        break;
      case '--binary-path':
        options.binaryPath = requireValue(argv, ++i, '--binary-path');
        break;
      case '--release-repository':
        options.releaseRepository = requireValue(argv, ++i, '--release-repository');
        break;
      case '--release-tag':
        options.releaseTag = requireValue(argv, ++i, '--release-tag');
        break;
      case '--help':
      case '-h':
        options.help = true;
        break;
      case '--version':
      case '-v':
        options.version = true;
        break;
      default:
        throw new Error(`unknown argument: ${arg}`);
    }
  }

  return options;
}

function requireValue(argv, index, flag) {
  const value = argv[index];
  if (!value || value.startsWith('-')) {
    throw new Error(`${flag} requires a value`);
  }
  return value;
}

async function buildChildEnv(options) {
  const env = { ...process.env };

  if (options.envFile) {
    await loadEnvFile(path.resolve(options.envFile), env);
  }

  if (options.transport) {
    env.MCP_TRANSPORT = options.transport;
  }
  if (!env.MCP_TRANSPORT && env.TRANSPORT) {
    env.MCP_TRANSPORT = env.TRANSPORT;
  }
  if (!env.MCP_TRANSPORT) {
    env.MCP_TRANSPORT = 'stdio';
  }
  if (!env.MCP_SERVER_PORT && env.PORT) {
    env.MCP_SERVER_PORT = env.PORT;
  }

  return env;
}

async function loadEnvFile(filePath, targetEnv) {
  const content = await fsp.readFile(filePath, 'utf8');
  for (const rawLine of content.split(/\r?\n/)) {
    let line = rawLine.trim();
    if (!line || line.startsWith('#')) {
      continue;
    }
    if (line.startsWith('export ')) {
      line = line.slice('export '.length).trim();
    }
    const index = line.indexOf('=');
    if (index <= 0) {
      continue;
    }
    const key = line.slice(0, index).trim();
    if (!/^[A-Za-z_][A-Za-z0-9_]*$/.test(key)) {
      throw new Error(`invalid env key in ${filePath}: ${key}`);
    }
    if (Object.prototype.hasOwnProperty.call(targetEnv, key)) {
      continue;
    }
    let value = line.slice(index + 1).trim();
    if ((value.startsWith('"') && value.endsWith('"')) || (value.startsWith("'") && value.endsWith("'"))) {
      value = value.slice(1, -1);
    }
    targetEnv[key] = value;
  }
}

async function resolveBinaryPath(options) {
  const explicitBinary = options.binaryPath || process.env.POSTGRES_MCP_BINARY_PATH;
  if (explicitBinary) {
    const resolved = path.resolve(explicitBinary);
    await assertFileExists(resolved, 'configured binary path does not exist');
    return resolved;
  }

  const releaseSource = resolveReleaseSource(options);
  const cacheDir = process.env.POSTGRES_MCP_CACHE_DIR ? path.resolve(process.env.POSTGRES_MCP_CACHE_DIR) : DEFAULT_CACHE_DIR;
  const binaryPath = path.join(cacheDir, binaryFileName());
  if (await pathExists(binaryPath)) {
    return binaryPath;
  }

  await downloadBinary(binaryPath, cacheDir, releaseSource);
  return binaryPath;
}

async function downloadBinary(binaryPath, cacheDir, releaseSource) {
  await fsp.mkdir(cacheDir, { recursive: true });
  const checksumUrl = releaseAssetUrl(CHECKSUM_FILE, releaseSource);
  const checksums = await downloadText(checksumUrl);
  const assetName = releaseAssetName();
  const expectedChecksum = parseChecksum(checksums, assetName);
  if (!expectedChecksum) {
    throw new Error(`checksum entry not found for ${assetName}`);
  }

  const compressed = await downloadBuffer(releaseAssetUrl(assetName, releaseSource));
  const actualChecksum = sha256(compressed);
  if (actualChecksum !== expectedChecksum) {
    throw new Error(`checksum mismatch for ${assetName}`);
  }

  const extracted = gunzip(compressed);
  const tempPath = `${binaryPath}.${process.pid}.tmp`;
  await fsp.writeFile(tempPath, extracted, { mode: isWindows() ? undefined : 0o755 });
  if (!isWindows()) {
    await fsp.chmod(tempPath, 0o755);
  }
  await fsp.rename(tempPath, binaryPath).catch(async (error) => {
    await safeUnlink(tempPath);
    if (error && error.code === 'EEXIST' && (await pathExists(binaryPath))) {
      return;
    }
    throw error;
  });
}

function resolveReleaseSource(options) {
  const repository = (options.releaseRepository || process.env.POSTGRES_MCP_RELEASE_REPOSITORY || DEFAULT_REPOSITORY).trim();
  const tag = (options.releaseTag || process.env.POSTGRES_MCP_RELEASE_TAG || DEFAULT_RELEASE_TAG).trim();
  const parts = repository.split('/').filter(Boolean);
  if (parts.length !== 2 || !parts.every((part) => /^[A-Za-z0-9._-]+$/.test(part))) {
    throw new Error(`invalid release repository: ${repository}; expected owner/repo`);
  }
  if (!tag) {
    throw new Error('release tag cannot be empty');
  }
  return {
    owner: parts[0],
    repo: parts[1],
    tag,
  };
}

function releaseAssetUrl(fileName, releaseSource) {
  return `https://github.com/${encodeURIComponent(releaseSource.owner)}/${encodeURIComponent(releaseSource.repo)}/releases/download/${encodeURIComponent(releaseSource.tag)}/${encodeURIComponent(fileName)}`;
}

function releaseAssetName() {
  const suffix = isWindows() ? '.exe.gz' : '.gz';
  return `postgres-server_${VERSION}_${platformName()}_${archName()}${suffix}`;
}

function binaryFileName() {
  return isWindows() ? 'postgres-server.exe' : 'postgres-server';
}

function platformName() {
  switch (process.platform) {
    case 'darwin':
      return 'darwin';
    case 'linux':
      return 'linux';
    case 'win32':
      return 'windows';
    default:
      throw new Error(`unsupported platform: ${process.platform}`);
  }
}

function archName() {
  switch (process.arch) {
    case 'x64':
      return 'amd64';
    case 'arm64':
      return 'arm64';
    default:
      throw new Error(`unsupported architecture: ${process.arch}`);
  }
}

function isWindows() {
  return process.platform === 'win32';
}

async function assertFileExists(filePath, message) {
  try {
    const stat = await fsp.stat(filePath);
    if (!stat.isFile()) {
      throw new Error(message);
    }
  } catch (error) {
    if (error && error.code === 'ENOENT') {
      throw new Error(`${message}: ${filePath}`);
    }
    throw error;
  }
}

async function ensureExecutable(filePath) {
  if (isWindows()) {
    return;
  }
  try {
    await fsp.chmod(filePath, 0o755);
  } catch (error) {
    throw new Error(`failed to mark binary executable: ${filePath}: ${error.message}`);
  }
}

async function runServer(binaryPath, env) {
  await new Promise((resolve, reject) => {
    const child = spawn(binaryPath, [], {
      stdio: 'inherit',
      env,
    });

    const forwardSignal = (signal) => {
      if (!child.killed) {
        child.kill(signal);
      }
    };

    process.on('SIGINT', forwardSignal);
    process.on('SIGTERM', forwardSignal);

    child.on('error', (error) => {
      reject(new Error(`failed to start ${binaryPath}: ${error.message}`));
    });
    child.on('exit', (code, signal) => {
      process.removeListener('SIGINT', forwardSignal);
      process.removeListener('SIGTERM', forwardSignal);
      if (signal) {
        process.kill(process.pid, signal);
        return;
      }
      process.exitCode = code == null ? 1 : code;
      resolve();
    });
  });
}

function parseChecksum(content, fileName) {
  for (const line of content.split(/\r?\n/)) {
    const trimmed = line.trim();
    if (!trimmed) {
      continue;
    }
    const parts = trimmed.split(/\s+/);
    if (parts.length < 2) {
      continue;
    }
    const maybeFile = parts[parts.length - 1].replace(/^\*/, '');
    if (maybeFile === fileName) {
      return parts[0].toLowerCase();
    }
  }
  return '';
}

function sha256(buffer) {
  return createHash('sha256').update(buffer).digest('hex');
}

function gunzip(buffer) {
  try {
    return zlib.gunzipSync(buffer);
  } catch (error) {
    throw new Error(`failed to decompress release asset: ${error.message}`);
  }
}

async function pathExists(filePath) {
  try {
    await fsp.access(filePath, fs.constants.F_OK);
    return true;
  } catch (_) {
    return false;
  }
}

async function safeUnlink(filePath) {
  try {
    await fsp.unlink(filePath);
  } catch (_) {
    // ignore cleanup errors
  }
}

async function downloadText(url) {
  const buffer = await downloadBuffer(url);
  return buffer.toString('utf8');
}

async function downloadBuffer(url, redirectsLeft = 5) {
  if (redirectsLeft < 0) {
    throw new Error(`too many redirects while downloading ${url}`);
  }

  return new Promise((resolve, reject) => {
    const request = https.get(url, {
      headers: {
        'User-Agent': `${PACKAGE_NAME}/${VERSION}`,
        'Accept-Encoding': 'identity',
      },
    }, (response) => {
      const statusCode = response.statusCode || 0;
      if ([301, 302, 303, 307, 308].includes(statusCode) && response.headers.location) {
        response.resume();
        const redirectUrl = new URL(response.headers.location, url).toString();
        downloadBuffer(redirectUrl, redirectsLeft - 1).then(resolve, reject);
        return;
      }
      if (statusCode < 200 || statusCode >= 300) {
        response.resume();
        reject(new Error(`download failed: ${url} -> HTTP ${statusCode}`));
        return;
      }

      const chunks = [];
      response.on('data', (chunk) => chunks.push(chunk));
      response.on('end', () => resolve(Buffer.concat(chunks)));
      response.on('error', reject);
    });

    request.on('error', reject);
  });
}
