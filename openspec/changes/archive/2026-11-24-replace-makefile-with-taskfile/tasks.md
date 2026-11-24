# Tasks

## 1. Taskfile Migration

- [x] 1.1 Add a root `Taskfile.yml` using Taskfile schema v3 with `generate`,
  `test`, and `run` tasks that execute the existing `templ generate ./...`,
  `go test ./...`, and `go run ./cmd/example` commands; verify `task --list`
  loads the file and lists all three tasks.
- [x] 1.2 Remove the root `Makefile` after the equivalent Taskfile tasks are
  present; verify the repository no longer presents Make as the supported task
  runner.

## 2. Developer Documentation

- [x] 2.1 Update `README.md` to document installing Task and invoking
  `task generate`, `task test`, and `task run`; verify every documented command
  matches a Taskfile task and the example startup flow remains accurate.

## 3. Integration Verification

- [x] 3.1 Run `task generate` and verify templ generation succeeds without
  changing the expected generated sources.
- [x] 3.2 Run `task test` and verify the complete Go test suite passes.
- [x] 3.3 Inspect the final diff to verify only the Make-to-Taskfile tooling
  migration and its documentation are included.

## Workflow follow-up

- Archive the change after implementation and review are complete.
