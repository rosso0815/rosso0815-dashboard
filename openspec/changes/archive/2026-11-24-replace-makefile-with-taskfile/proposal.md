# Proposal

## Why

The repository currently relies on a Makefile for three simple development
commands, while the project is intended to use Taskfile for task automation.
Replacing the Makefile with an equivalent Taskfile standardizes the workflow
without changing how the application is generated, tested, or run.

## What Changes

- Add a root `Taskfile.yml` with `generate`, `test`, and `run` tasks matching
  the current Makefile commands.
- Remove the root `Makefile`.
- Update the README to document Taskfile-based development commands and the
  required Task installation.
- Preserve the existing command behavior:
  `templ generate ./...`, `go test ./...`, and `go run ./cmd/example`.
- Add focused validation that the Taskfile parses and each task remains
  executable.

## Capabilities

### New Capabilities

None. This is a development-tooling and documentation migration.

### Modified Capabilities

None. No application behavior or existing OpenSpec capability requirements
change.

## Impact

- Affects the repository root task configuration and `README.md`.
- Removes the Make-based developer entry point and introduces Task as the
  development-task dependency.
- Does not change Go APIs, generated templates, routes, runtime behavior, or
  production dependencies.
