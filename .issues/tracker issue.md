# Overview

This is a meta issue used to track the Ground Control chore works, such as repository-structure and tooling cleanup as one parent issue with the following focused sub-issues. The order minimizes repeated edits to command paths, imports, specifications, and Task files.

## Issues

| Order | Issue | Number | Status | Depends on |
|---:|---|:---:|:---:|---|
| 1 | Refactor command layout and binary names | [#588](https://github.com/container-registry/harbor-satellite/issues/588) | Open | — |
| 2 | Remove obsolete go-swagger tasks | #590 | Open | #588 |
| 3 | Relocate SQLC sources and add generation tasks | #591 | Open | #588, #590 |
| 4 | Make tool-version pins authoritative | #592 | Open | #588, #590, #591 |
| 5 | Organize shared internal packages | #593 | Open | #588 |
| 6 | Consolidate the Ground Control logger | #594 | Open | #593 |
| 7 | Merge Harbor health into the Harbor client package | #595 | Open | #588, #593 |
