# Overview

Consolidate logging under the shared [`internal/logger`](https://github.com/container-registry/harbor-satellite/tree/main/internal/logger) package and remove the duplicated [`internal/groundcontrol/logger`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol/logger) package.

Satellite already uses the shared logger, while Ground Control primarily imports its component-local copy. Maintaining two implementations of the same audit, syslog, and OpenTelemetry logging behavior creates unnecessary duplication and allows fixes to diverge.

## Description

Both packages provide the same core audit logging model and transports through corresponding `audit.go`, `otel.go`, and `syslog.go` files and tests. The shared package is the appropriate canonical implementation because it serves repository-wide logging and additionally contains the common Zerolog context logger in [`internal/logger/logger.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/logger/logger.go).

The copies are functionally equivalent but are no longer byte-identical: they contain minor differences in error construction, OTLP response wording, formatting, and context-key placement. Keeping both packages means future fixes, tests, and behavior changes must be applied twice and can silently drift. Ground Control should consume the shared package, preserving `ComponentGroundControl` attribution and its existing audit behavior.

## Expected changes

- Replace `github.com/container-registry/harbor-satellite/internal/groundcontrol/logger` with `github.com/container-registry/harbor-satellite/internal/logger` in every current consumer:
  - [`internal/groundcontrol/server/audit_config_test.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/audit_config_test.go)
  - [`internal/groundcontrol/server/auth_handlers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/auth_handlers.go)
  - [`internal/groundcontrol/server/config_handlers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/config_handlers.go)
  - [`internal/groundcontrol/server/middleware.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/middleware.go)
  - [`internal/groundcontrol/server/satellite_handlers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/satellite_handlers.go)
  - [`internal/groundcontrol/server/server.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/server.go)
  - [`internal/groundcontrol/server/user_handlers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/user_handlers.go)
  - [`internal/env/utils.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/env/utils.go), which constructs Ground Control audit configurations outside the `internal/groundcontrol` tree
- Preserve the existing `auditlog` import alias where it keeps call sites readable; only the package source should change.
- Reconcile the minor implementation differences before deleting the local package. Prefer the shared behavior unless a Ground Control-specific difference is intentional and covered by a test.
- Move or merge any unique test coverage from [`internal/groundcontrol/logger`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol/logger) into [`internal/logger`](https://github.com/container-registry/harbor-satellite/tree/main/internal/logger).
- Remove the complete `internal/groundcontrol/logger` directory after all consumers use the shared package.
- Confirm Ground Control continues to emit audit events with `ComponentGroundControl`, including authentication, user, configuration, and satellite operations.
- Run shared logger tests, Ground Control server tests, environment configuration tests, lint, and a repository-wide search to ensure no imports of `internal/groundcontrol/logger` remain.
