# gobce

gobce (Go Branch Coverage Estimator) is a standalone CLI that estimates C1 branch coverage for Go projects.
It combines `go test -coverprofile` output with lightweight AST/CFG analysis to produce:
- statement coverage
- estimated C1 branch coverage
- uncovered branch findings

> Note: gobce reports **estimated C1**, not exact branch coverage, because Go coverprofiles are statement/block-oriented.

## Setup

Requirements:
- Go 1.26+
- A Go project where `go test ./... -coverprofile coverage.out` works

## Install

Install via `go install`:

```bash
go install github.com/keyskey/gobce/cmd/gobce@latest
```

For CI and team-wide reproducibility, pin a version tag:

```bash
go install github.com/keyskey/gobce/cmd/gobce@v0.2.0
```

## Binary Distribution

Prebuilt binaries can also be distributed through GitHub Releases.
Recommended archive set:
- darwin-arm64
- darwin-amd64
- linux-arm64
- linux-amd64

Install from a prebuilt archive (example: `v0.2.0`, `darwin-arm64`):

```bash
curl -fsSL "https://raw.githubusercontent.com/keyskey/gobce/main/scripts/install.sh" | sh
```

Install a specific version:

```bash
VERSION=v0.2.0 curl -fsSL "https://raw.githubusercontent.com/keyskey/gobce/main/scripts/install.sh" | sh
```

Install to a custom directory:

```bash
INSTALL_DIR="$HOME/.local/bin" curl -fsSL "https://raw.githubusercontent.com/keyskey/gobce/main/scripts/install.sh" | sh
```

## CLI Usage

```bash
go test ./... -coverprofile coverage.out
go run ./cmd/gobce analyze --coverprofile coverage.out --format json
```

Write JSON result to file (while keeping stdout output):

```bash
go run ./cmd/gobce analyze --coverprofile coverage.out --format json --output gobce-result.json
```

Build an executable:

```bash
go build -o gobce ./cmd/gobce
./gobce analyze --coverprofile coverage.out --format json
```

## GitHub Actions (CI)

Minimal workflow example:

```yaml
name: test-and-gobce

on:
  pull_request:
  push:
    branches: [main]

jobs:
  coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.26.x"

      - name: Generate coverage profile
        run: go test ./... -coverprofile coverage.out

      - name: Run gobce
        run: go run github.com/keyskey/gobce/cmd/gobce@latest analyze --coverprofile coverage.out --format json
```

If you want to fail CI by threshold, combine JSON output with `jq`:

```bash
go run github.com/keyskey/gobce/cmd/gobce@latest analyze --coverprofile coverage.out --format json --output gobce.json
jq -e '.estimatedBranchCoverage >= 70' gobce.json
```

For stable CI behavior, prefer pinned versions (`@vX.Y.Z`) over `@latest`.

## Versioning Policy

`gobce` follows [Semantic Versioning](https://semver.org/).

During the pre-1.0 phase, this project is versioned as `0.x` while coverage logic is validated in a real project.

- `PATCH` (`0.1.0` -> `0.1.1`): bug fixes only
- `MINOR` (`0.1.1` -> `0.2.0`): new features, behavior changes, and potentially breaking changes
- `MAJOR` (`1.x`): used after the interface and output become stable

Please assume that interfaces and output formats may change between `0.x` minor versions.

Version tags use the `vX.Y.Z` format (for example, `v0.2.0`).

## Release Flow

This repository is configured to publish GitHub Releases via GoReleaser when a version tag is pushed.
Release notes are generated automatically from git commits between tags using GoReleaser changelog support.

You can create the next SemVer tag with:

```bash
make tag-patch
make tag-minor
make tag-major
# or
make tag-next TYPE=patch
make tag-next TYPE=minor
make tag-next TYPE=major
```

```bash
git tag v0.2.0
git push origin v0.2.0
```

After push:
- GitHub Actions workflow `release` runs automatically
- GitHub Release is created for `v0.2.0`
- Release notes are generated from commit history
- OS/arch archives and `checksums.txt` are uploaded

To keep generated release notes clean, prefer commit prefixes such as `feat:`, `fix:`, `refactor:`, and `perf:`.

## Design

Design notes are available in:
- `docs/gobce.md` (concept and design notes)
- `docs/how-gobce-works.md` (beginner-friendly implementation walkthrough)