# Overview

Refactor the Go command layout so `go install` produces clear binary names, and standardize Ground Control naming across repository paths, code, automation, and documentation.

Use the following command layout:

```text
cmd/
├── groundctl
├── groundcontrold
└── satellite
```

- `groundctl`: Ground Control administration client.
- `groundcontrold`: Ground Control server daemon. This is preferred over `groundsrv` because the conventional `d` suffix clearly identifies a daemon while preserving the product name.
- `satellite`: Satellite executable.

## Description

The current entry points are nested under [`cmd/groundcontrol/cli`](https://github.com/container-registry/harbor-satellite/tree/main/cmd/groundcontrol/cli) and [`cmd/groundcontrol/server`](https://github.com/container-registry/harbor-satellite/tree/main/cmd/groundcontrol/server). Since `go install` derives the executable name from the final directory, direct installation produces generic binaries named `cli` and `server`:

```shell
go install github.com/container-registry/harbor-satellite/cmd/groundcontrol/cli@<version>
go install github.com/container-registry/harbor-satellite/cmd/groundcontrol/server@<version>
```

The Task build system instead outputs `groundcontrol` for the client and `ground-control` for the server through [`taskfiles/build.yml`](https://github.com/container-registry/harbor-satellite/blob/main/taskfiles/build.yml). These names differ only by a hyphen and do not clearly distinguish the client from the server.

Naming is also inconsistent elsewhere. Code paths such as [`internal/groundcontrol`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol) use the closed compound form, while paths such as [`spec/ground-control`](https://github.com/container-registry/harbor-satellite/tree/main/spec/ground-control), [`examples/deploy/helm/ground-control`](https://github.com/container-registry/harbor-satellite/tree/main/examples/deploy/helm/ground-control), and [`.air.ground-control.toml`](https://github.com/container-registry/harbor-satellite/blob/main/.air.ground-control.toml) use hyphenated names. Tasks and build resources additionally mix `gc`, `ground-control`, and `groundcontrol`.

Adopt one convention:

- **Ground Control** in user-facing prose and product documentation.
- `groundcontrol` in paths, filenames, code identifiers, package names, tasks, build components, and generated-code conventions.
- `groundctl` only for the client executable.
- `groundcontrold` only for the server executable.

## Expected changes

- Move [`cmd/groundcontrol/cli/root.go`](https://github.com/container-registry/harbor-satellite/blob/main/cmd/groundcontrol/cli/root.go) to `cmd/groundctl/main.go`.
- Move [`cmd/groundcontrol/server/main.go`](https://github.com/container-registry/harbor-satellite/blob/main/cmd/groundcontrol/server/main.go) to `cmd/groundcontrold/main.go`.
- Update the client's Cobra name, help output, examples, and executable-specific configuration references to `groundctl`.
- Update server build and execution references to `groundcontrold`.
- Ensure direct installation produces the intended names:

  ```shell
  go install github.com/container-registry/harbor-satellite/cmd/groundctl@<version>
  go install github.com/container-registry/harbor-satellite/cmd/groundcontrold@<version>
  go install github.com/container-registry/harbor-satellite/cmd/satellite@<version>
  ```

- Rename repository-owned paths and files to the closed compound form, including:
  - `spec/ground-control` to `spec/groundcontrol`.
  - `examples/deploy/helm/ground-control` to `examples/deploy/helm/groundcontrol`.
  - `.air.ground-control.toml` to `.air.groundcontrol.toml`.
  - Other source and documentation filenames using `ground-control` or `ground_control`.
- Normalize task names, variables, artifact directories, component values, and repository-owned `gc` abbreviations in [`Taskfile.yml`](https://github.com/container-registry/harbor-satellite/blob/main/Taskfile.yml) and the files under [`taskfiles`](https://github.com/container-registry/harbor-satellite/tree/main/taskfiles).
- Update affected Docker build arguments, Compose files, generation configuration, scripts, tests, CI workflows, release automation, links, and documentation.
- Make local builds, cross-platform builds, release artifacts, and `go install` use the same canonical executable names.
- Review compatibility-sensitive public names—such as `GROUND_CONTROL_*` variables, `--ground-control-*` flags, image and service names, Helm metadata, API fields, and SPIFFE IDs—before renaming them. Preserve them or provide aliases and migration guidance where a rename would break existing deployments.
- Verify all three commands build successfully; run formatting, lint, generation, unit, E2E, container, and release checks; and confirm that any remaining `gc`, `ground-control`, or `ground_control` occurrences are intentional.
