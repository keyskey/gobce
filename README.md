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
go install github.com/keyskey/gobce/cmd/gobce@v0.1.0
```

## Binary Distribution

Prebuilt binaries can also be distributed through GitHub Releases.
Recommended archive set:
- darwin-arm64
- darwin-amd64
- linux-arm64
- linux-amd64

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

## Release Flow

This repository is configured to publish GitHub Releases via GoReleaser when a version tag is pushed.

```bash
git tag v0.1.0
git push origin v0.1.0
```

After push:
- GitHub Actions workflow `release` runs automatically
- GitHub Release is created for `v0.1.0`
- OS/arch archives and `checksums.txt` are uploaded

## Design

Design notes are available in:
- `docs/gobce.md` (concept and design notes)
- `docs/how-gobce-works.md` (beginner-friendly implementation walkthrough)