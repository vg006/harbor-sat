# Overview

Remove misleading unused tool-version variables and make the versions used by Task, CI, containers, releases, and contributor documentation authoritative and reproducible.

## Description

[`Taskfile.yml`](https://github.com/container-registry/harbor-satellite/blob/main/Taskfile.yml) declares `GO_VERSION`, `GOLANGCILINT_VERSION`, and `GORELEASER_VERSION`, but none of the current tasks consume these variables.

The actual tool selection happens elsewhere: Go is pinned independently in [`go.mod`](https://github.com/container-registry/harbor-satellite/blob/main/go.mod) and [`Dockerfile`](https://github.com/container-registry/harbor-satellite/blob/main/Dockerfile); lint tasks use a container digest from [`taskfiles/lint.yml`](https://github.com/container-registry/harbor-satellite/blob/main/taskfiles/lint.yml); and release tasks execute whichever `goreleaser` binary is available on `PATH`. [`taskfiles/README.md`](https://github.com/container-registry/harbor-satellite/blob/main/taskfiles/README.md) additionally documents golangci-lint `v2.0.2`, while the unused root variable claims `v2.12.2`.

These disconnected values appear to pin the toolchain without controlling it, making local and CI behavior harder to reproduce.

## Expected changes

- Inventory the effective Go, golangci-lint, GoReleaser, govulncheck, oapi-codegen, and SQLC versions used by local tasks, CI, Docker builds, and releases.
- Remove version variables that are only informational, or wire them into installation and verification tasks so they actually control the invoked tool.
- Keep one authoritative source for each tool version where practical, with dependent tasks and documentation referencing it rather than duplicating literal versions.
- Ensure the Go version remains aligned across [`go.mod`](https://github.com/container-registry/harbor-satellite/blob/main/go.mod), [`Dockerfile`](https://github.com/container-registry/harbor-satellite/blob/main/Dockerfile), CI setup, and contributor prerequisites.
- Make the golangci-lint container digest traceable to its intended release and align [`taskfiles/README.md`](https://github.com/container-registry/harbor-satellite/blob/main/taskfiles/README.md) with the actual local and CI execution path.
- Pin or verify GoReleaser before `snapshot` and `release` tasks instead of silently accepting any binary on `PATH`.
- Update contributor documentation with canonical installation commands and version checks.
- Add lightweight task preconditions or status checks that fail with a clear message when a required tool is missing or has an unsupported version.
- Verify local build, lint, vulnerability, generation, snapshot, and CI workflows use the documented versions with no unused version declarations remaining.
