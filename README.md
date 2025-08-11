# Self-Updater

A Go-based self-updating server and launcher utility. The API server can update its binary at runtime using a signed release manifest.

## Features

- Self-updating server process with minimal downtime restart signals
- Signed release manifest generation and verification
- Multi-platform builds (Linux, macOS, Windows for amd64 and arm64)
- Simple CLI tool for signing artifacts

## Prerequisites

- Go 1.24+
- Git
- (Optional) Nix: `nix develop` for a consistent development shell
- (Optional) Direnv for automatic environment loading

## Getting Started

Clone the repository:

```bash
git clone https://github.com/danilevy1212/self-updater.git
cd self-updater
```

Enter a development shell (requires Nix):

```bash
nix develop
```

## Building

Use the Makefile to format, test, and build:

```bash
make all        # format, test, and build binaries
make build-api  # build only the API binaries
```

The binaries will be in the `bin/` directory:

- `bin/api-linux-amd64`
- `bin/api-windows-amd64.exe`
- `bin/api-darwin-amd64`
- `bin/api-linux-arm64`
- `bin/api-darwin-arm64`
- `bin/api-windows-arm64.exe`

## Configuration

The server can be configured via environment variables:

| Variable                | Default            | Description                                     |
| ----------------------- | ------------------ | ----------------------------------------------- |
| SERVER_PORT             | 3000               | Port for the API server to listen               |
| SERVER_IS_DEV           | false              | Enable development mode                         |
| UPDATER_IS_DEV          | false              | Enable updater development mode                 |
| UPDATER_CRON_SCHEDULE   | \* \* \* \* \*     | Cron schedule for updates                       |
| UPDATER_RUN_AT_BOOT     | true               | Run updater at boot time                        |
| LAUNCHER_IS_DEV         | false              | Enable launcher development mode                |
| LAUNCHER_SESSION_FOLDER | update-session     | Folder in temporary storage for update sessions |
| ARCHIVER_REPO           | self-updater       | GitHub repository for the release manifest      |
| ARCHIVER_OWNER          | your-org           | GitHub owner for the release manifest           |
| ARCHIVER_BASE_URL       | https://github.com | Base URL for the release manifest               |

## Usage

### Run the launcher

```bash
./bin/api-linux-amd64
```

### Run the server stand alone

```bash
./bin/api-linux-amd64 --server --current-session-dir /tmp/update-session/<unique-id>
```

- `--server`: run the API server and updater in the same process
- `--current-session-dir`: directory used to store and swap binaries

### Sign artifacts

Build the sign tool and use it to sign binaries or the manifest:

```bash
go run ./cmd/sign <private-key.pem> <file-to-sign>
# or, after building:
./bin/sign <private-key.pem> <file-to-sign>
```

## Release

To generate or update the signed release manifest (`internal/assets/release.json`):

```bash
export SIGNING_KEY_PEM="$(cat path/to/private.pem)"
export PUBLIC_KEY_PEM="$(cat path/to/public.pem)"
export ARCHIVER_BASE_URL="https://github.com/your-org/self-updater/releases/download"
export ARCHIVER_OWNER="your-org"
export ARCHIVER_REPO="self-updater"
make release
```

## Testing

```bash
make test
```

## Design

The self-updater is designed to allow a server process to update itself with minimal downtime. It uses a signed manifest to verify the integrity of updates and supports multiple platforms.

When executed in it's default mode (launcher mode), the binary will create a temporary directory for the update session. It will copy the current binary to this directory and start it in server mode. The server will run the updater in a separate goroutine, which will periodically check for updates based on the configured cron schedule.

If an update is available, the updater will download the new binary, verify it using the signed manifest, and signal to the launcher to restart with the new binary. The launcher will then swap the old binary with the new one and restart the server process.
