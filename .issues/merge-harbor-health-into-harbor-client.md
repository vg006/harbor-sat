# Overview

Merge [`internal/groundcontrol/harborhealth`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol/harborhealth) into the existing [`internal/groundcontrol/harbor`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol/harbor) package and use the health client provided by the pinned Harbor Go SDK.

## Description

Ground Control already centralizes Harbor API access through [`internal/groundcontrol/harbor/client.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/harbor/client.go), which exposes the generated `HarborAPI` client from `github.com/goharbor/go-client` `v0.213.1`. That client includes a `Health` service and generated `GetHealth` operation, as documented by the [Harbor SDK health package](https://pkg.go.dev/github.com/goharbor/go-client@v0.213.1/pkg/sdk/v2.0/client/health).

The separate health package currently constructs its own HTTP client, calls `/api/v2.0/health`, and maintains local response models in [`internal/groundcontrol/harborhealth/check.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/harborhealth/check.go) and [`internal/groundcontrol/harborhealth/types.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/harborhealth/types.go). This duplicates URL, transport, endpoint, and response-model behavior already supplied by the SDK.

The package does contain useful Ground Control policy that must be preserved: skipping the startup check when configured, applying a bounded timeout, and ignoring selected optional Harbor components. These responsibilities fit within the existing `harbor` integration package and do not require a separate top-level package.

## Expected changes

- Move the health-check orchestration into `internal/groundcontrol/harbor`, for example as `harbor.CheckHealth` in a dedicated `health.go` file.
- Use `GetClient().Health.GetHealth` and the generated health request parameters and response models from the pinned Harbor SDK instead of issuing a handwritten `http.Client.Get` request.
- Preserve the existing startup policy:
  - Honor `SKIP_HARBOR_HEALTH_CHECK` and log when the check is skipped.
  - Apply the current five-second timeout through context or the generated request parameters.
  - Continue ignoring the configured optional components (`portal`, `trivy`, `registryctl`, and `jobservice`).
  - Return an error when a required Harbor component is unhealthy or the health request fails.
- Keep policy-specific filtering helpers inside the `harbor` package, but remove local wire-format types when the generated SDK models provide the required health fields.
- Update the only current importer, [`cmd/groundcontrol/server/main.go`](https://github.com/container-registry/harbor-satellite/blob/main/cmd/groundcontrol/server/main.go), from `internal/groundcontrol/harborhealth` to `internal/groundcontrol/harbor` and replace `harborhealth.CheckHealth()` with `harbor.CheckHealth()`. If the command-layout refactor lands first, apply the same change in `cmd/groundcontrold/main.go`.
- Remove the complete `internal/groundcontrol/harborhealth` directory after its policy and behavior are covered by the `harbor` package.
- Add focused tests for healthy responses, required-component failures, ignored optional components, request errors/timeouts, and the skip-health-check configuration.
- Run Ground Control package tests and a repository-wide search confirming that no `internal/groundcontrol/harborhealth` imports remain.
