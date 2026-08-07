# Overview

Group repository-wide internal packages under `internal/shared`, leaving only the component-owned [`internal/groundcontrol`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol) and [`internal/satellite`](https://github.com/container-registry/harbor-satellite/tree/main/internal/satellite) trees at the top level.

Use `shared` rather than `common` because it describes the ownership boundary without implying that unrelated helpers should be combined into a generic package. Each relocated package must retain its existing responsibility and package boundary.

## Description

The current [`internal`](https://github.com/container-registry/harbor-satellite/tree/main/internal) tree mixes component-specific code with packages shared by commands and components. Moving the shared packages under one namespace makes ownership clear and produces a predictable layout:

```text
internal/
├── groundcontrol/
├── satellite/
└── shared/
    ├── crypto/
    ├── env/
    ├── logger/
    ├── spiffe/
    └── utils/
```

This is an organizational refactor only. It must not merge these packages, broaden their APIs, or change runtime behavior.

## Expected changes

- Move the existing packages without changing their package names:
  - [`internal/crypto`](https://github.com/container-registry/harbor-satellite/tree/main/internal/crypto) to `internal/shared/crypto`.
  - [`internal/env`](https://github.com/container-registry/harbor-satellite/tree/main/internal/env) to `internal/shared/env`.
  - [`internal/logger`](https://github.com/container-registry/harbor-satellite/tree/main/internal/logger) to `internal/shared/logger`.
  - [`internal/spiffe`](https://github.com/container-registry/harbor-satellite/tree/main/internal/spiffe) to `internal/shared/spiffe`.
  - [`internal/utils`](https://github.com/container-registry/harbor-satellite/tree/main/internal/utils) to `internal/shared/utils`.
- Coordinate with the logger-consolidation work so [`internal/groundcontrol/logger`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol/logger) is removed and Ground Control imports `internal/shared/logger` rather than introducing another duplicate.
- Update the following active imports of `internal/crypto` to `internal/shared/crypto`:
  - [`internal/groundcontrol/auth/password.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/auth/password.go)
  - [`internal/groundcontrol/server/helpers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/helpers.go)
  - [`internal/groundcontrol/server/helpers_test.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/helpers_test.go)
  - [`internal/groundcontrol/server/middleware.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/middleware.go)
  - [`internal/groundcontrol/server/middleware_test.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/middleware_test.go)
  - [`internal/satellite/secure/config.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/secure/config.go)
  - [`internal/satellite/secure/config_test.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/secure/config_test.go)
  - [`pkg/config/manager.go`](https://github.com/container-registry/harbor-satellite/blob/main/pkg/config/manager.go)
- Update the following active imports of `internal/env` to `internal/shared/env`:
  - [`cmd/groundcontrol/server/main.go`](https://github.com/container-registry/harbor-satellite/blob/main/cmd/groundcontrol/server/main.go)
  - [`cmd/satellite/main.go`](https://github.com/container-registry/harbor-satellite/blob/main/cmd/satellite/main.go)
  - [`internal/groundcontrol/auth/policy.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/auth/policy.go)
  - [`internal/groundcontrol/auth/policy_test.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/auth/policy_test.go)
  - [`internal/groundcontrol/harbor/client.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/harbor/client.go)
  - [`internal/groundcontrol/harbor/robot.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/harbor/robot.go)
  - [`internal/groundcontrol/harbor/robot_test.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/harbor/robot_test.go)
  - [`internal/groundcontrol/harborhealth/check.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/harborhealth/check.go)
  - [`internal/groundcontrol/migrator/migrator.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/migrator/migrator.go)
  - [`internal/groundcontrol/server/audit_config_test.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/audit_config_test.go)
  - [`internal/groundcontrol/server/bootstrap.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/bootstrap.go)
  - [`internal/groundcontrol/server/config_handlers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/config_handlers.go)
  - [`internal/groundcontrol/server/group_handlers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/group_handlers.go)
  - [`internal/groundcontrol/server/helpers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/helpers.go)
  - [`internal/groundcontrol/server/satellite_handlers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/satellite_handlers.go)
  - [`internal/groundcontrol/server/satellite_handlers_test.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/satellite_handlers_test.go)
  - [`internal/groundcontrol/server/server.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/server.go)
  - [`internal/groundcontrol/spiffe/provider.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/spiffe/provider.go)
  - [`internal/groundcontrol/utils/helper.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/utils/helper.go)
- Update the following active imports of `internal/logger` to `internal/shared/logger`:
  - [`cmd/satellite/audit_config_test.go`](https://github.com/container-registry/harbor-satellite/blob/main/cmd/satellite/audit_config_test.go)
  - [`cmd/satellite/main.go`](https://github.com/container-registry/harbor-satellite/blob/main/cmd/satellite/main.go)
  - [`internal/satellite/satellite.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/satellite.go)
  - [`internal/satellite/state/catalog.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/catalog.go)
  - [`internal/satellite/state/catalog_test.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/catalog_test.go)
  - [`internal/satellite/state/direct_delivery.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/direct_delivery.go)
  - [`internal/satellite/state/registration_process.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/registration_process.go)
  - [`internal/satellite/state/replicator.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/replicator.go)
  - [`internal/satellite/state/report.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/report.go)
  - [`internal/satellite/state/reporting_process.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/reporting_process.go)
  - [`internal/satellite/state/spiffe_registration.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/spiffe_registration.go)
  - [`internal/satellite/state/state_process.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/state_process.go)
- As part of the coordinated logger consolidation, also replace `internal/groundcontrol/logger` imports with `internal/shared/logger` in:
  - [`internal/groundcontrol/server/audit_config_test.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/audit_config_test.go)
  - [`internal/groundcontrol/server/auth_handlers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/auth_handlers.go)
  - [`internal/groundcontrol/server/config_handlers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/config_handlers.go)
  - [`internal/groundcontrol/server/middleware.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/middleware.go)
  - [`internal/groundcontrol/server/satellite_handlers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/satellite_handlers.go)
  - [`internal/groundcontrol/server/server.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/server.go)
  - [`internal/groundcontrol/server/user_handlers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/server/user_handlers.go)
  - [`internal/env/utils.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/env/utils.go), at its new location under `internal/shared/env`
- Update the following active imports of `internal/spiffe` to `internal/shared/spiffe`:
  - [`internal/satellite/state/reporting_process.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/reporting_process.go)
  - [`internal/satellite/state/spiffe_registration.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/spiffe_registration.go)
- Update the following active imports of `internal/utils` to `internal/shared/utils`:
  - [`cmd/satellite/main.go`](https://github.com/container-registry/harbor-satellite/blob/main/cmd/satellite/main.go)
  - [`internal/logger/logger.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/logger/logger.go), at its new location under `internal/shared/logger`
  - [`internal/satellite/state/helpers.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/helpers.go)
  - [`internal/satellite/state/reporting_process.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/reporting_process.go)
  - [`internal/satellite/state/state_process.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/state/state_process.go)
- Update or remove stale commented import references in [`internal/satellite/container_runtime/host.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/container_runtime/host.go) and [`internal/satellite/container_runtime/read_config.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/satellite/container_runtime/read_config.go).
- Update non-Go references such as documentation, architecture diagrams, lint configuration, scripts, or tooling that mention the old package paths.
- Run formatting, lint, unit tests, integration tests, command builds, and a repository-wide search confirming that no old `internal/crypto`, `internal/env`, `internal/logger`, `internal/spiffe`, or `internal/utils` import paths remain.
