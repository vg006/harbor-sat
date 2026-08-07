# Overview

Remove the obsolete go-swagger installation and generation workflow from the root Task configuration, leaving the OpenAPI specification and oapi-codegen tasks as the canonical Ground Control API generation path.

## Description

[`Taskfile.yml`](https://github.com/container-registry/harbor-satellite/blob/main/Taskfile.yml) still defines `GO_SWAGGER_VERSION`, `swagger:install`, and `swagger:generate`. The generation task runs from a removed `ground-control` directory and expects `meta.yml`, `swagger.yml`, and `swagger-flatten.yml`, none of which exist in the current repository layout. The task is therefore unusable.

Ground Control now keeps its OpenAPI source under [`spec/ground-control`](https://github.com/container-registry/harbor-satellite/tree/main/spec/ground-control) and generates server and client code through [`taskfiles/gen.yml`](https://github.com/container-registry/harbor-satellite/blob/main/taskfiles/gen.yml) using oapi-codegen. Retaining the old go-swagger workflow creates two apparent generation paths and leaves contributors unable to tell which one is authoritative.

## Expected changes

- Remove `GO_SWAGGER_VERSION`, `swagger:install`, and `swagger:generate` from [`Taskfile.yml`](https://github.com/container-registry/harbor-satellite/blob/main/Taskfile.yml).
- Remove any remaining go-swagger installation, cache, documentation, or CI references that are only used by the obsolete workflow.
- Keep [`spec/ground-control/openapi.yaml`](https://github.com/container-registry/harbor-satellite/blob/main/spec/ground-control/openapi.yaml) as the canonical Ground Control API specification, using its renamed `spec/groundcontrol` path if the naming refactor lands first.
- Keep the oapi-codegen server and client tasks in [`taskfiles/gen.yml`](https://github.com/container-registry/harbor-satellite/blob/main/taskfiles/gen.yml) as the only supported API code-generation workflow.
- Update contributor documentation so it exposes one canonical command for regenerating Ground Control API code.
- Verify the public generation tasks produce both committed generated files without requiring go-swagger, and confirm a repository-wide search finds no stale legacy Swagger task or tool references.
