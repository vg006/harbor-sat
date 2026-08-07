# Overview

Move the Ground Control SQLC configuration, query definitions, and database schema into the Ground Control specification tree, and add reproducible Task targets for generating its database package.

The root [`sqlc.yaml`](https://github.com/container-registry/harbor-satellite/blob/main/sqlc.yaml) and SQL sources under [`internal/groundcontrol/sql`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol/sql) are specific to Ground Control and belong with its other source specifications under [`spec/ground-control`](https://github.com/container-registry/harbor-satellite/tree/main/spec/ground-control).

## Description

The SQLC configuration consumes Ground Control queries and schema definitions from [`internal/groundcontrol/sql`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol/sql) and writes generated Go code to [`internal/groundcontrol/database`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol/database). The query files are SQLC generation inputs rather than Go implementation code, so keeping them under `internal` obscures their role. Co-locating the configuration, queries, and schema under the specification tree makes their ownership and generation relationship explicit.

The schema directory has an additional runtime responsibility: [`internal/groundcontrol/migrator/migrator.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/migrator/migrator.go) loads it as Goose migrations during local execution, and [`Dockerfile`](https://github.com/container-registry/harbor-satellite/blob/main/Dockerfile) copies it to `/migrations` for container execution. Moving the schema therefore requires updating both SQLC generation and migration packaging/runtime paths.

The generated database package identifies SQLC `v1.31.1` in [`internal/groundcontrol/database/db.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/database/db.go), but [`Taskfile.yml`](https://github.com/container-registry/harbor-satellite/blob/main/Taskfile.yml) and [`taskfiles/gen.yml`](https://github.com/container-registry/harbor-satellite/blob/main/taskfiles/gen.yml) provide no SQLC installation or generation task. Contributors must install and invoke SQLC manually, with no project-managed version or canonical command. This can produce inconsistent generated output and allows query or schema changes to be merged without verifying that committed database code is current.

## Expected changes

- Move the complete Ground Control SQL source layout into the specification tree:
  - [`sqlc.yaml`](https://github.com/container-registry/harbor-satellite/blob/main/sqlc.yaml) to `spec/ground-control/sqlc.yaml`.
  - [`internal/groundcontrol/sql/queries`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol/sql/queries) to `spec/ground-control/sql/queries`.
  - [`internal/groundcontrol/sql/schema`](https://github.com/container-registry/harbor-satellite/tree/main/internal/groundcontrol/sql/schema) to `spec/ground-control/sql/schema`.
  - If the repository-wide naming refactor lands first, use `spec/groundcontrol` as the canonical base path instead.
- Update the relocated SQLC configuration so its query and schema inputs resolve within the specification tree and its generated output remains `internal/groundcontrol/database`.
- Update [`internal/groundcontrol/migrator/migrator.go`](https://github.com/container-registry/harbor-satellite/blob/main/internal/groundcontrol/migrator/migrator.go) to use the relocated schema directory as its local-development migration fallback.
- Update [`Dockerfile`](https://github.com/container-registry/harbor-satellite/blob/main/Dockerfile) to copy the relocated schema into `/migrations`, preserving container migration behavior.
- Add a pinned SQLC version variable to [`Taskfile.yml`](https://github.com/container-registry/harbor-satellite/blob/main/Taskfile.yml), initially matching the version recorded by the committed generated files unless an intentional upgrade is included.
- Add an internal SQLC installation task to [`taskfiles/gen.yml`](https://github.com/container-registry/harbor-satellite/blob/main/taskfiles/gen.yml). Install the pinned binary into the project `bin` directory and verify its version before reuse.
- Add a Ground Control database generation task that invokes SQLC with the relocated configuration explicitly, rather than depending on the caller's working directory or SQLC's default config discovery.
- Expose the database generator through a public task in [`Taskfile.yml`](https://github.com/container-registry/harbor-satellite/blob/main/Taskfile.yml), following the repository's final closed-compound naming convention, for example `generate:groundcontrol-database` with a concise alias.
- Include database generation in the aggregate Ground Control generation task alongside OpenAPI server and client generation, while retaining a standalone task for database-only changes.
- Declare the relocated SQL queries, migrations, and SQLC configuration as task sources and the generated database files as outputs where Task source tracking is reliable.
- Update contributor documentation with the canonical generation command and clarify that files under `internal/groundcontrol/database` are generated and must not be edited manually.
- Update migration documentation and all remaining references to the old `internal/groundcontrol/sql` paths.
- Add a CI verification step that runs SQLC generation and fails when it leaves a Git diff, ensuring committed database code remains synchronized with queries, migrations, configuration, and the pinned SQLC version.
- Verify SQLC generation, local Ground Control migrations, and container startup so both generation-time and runtime consumers of the relocated schema are covered.
