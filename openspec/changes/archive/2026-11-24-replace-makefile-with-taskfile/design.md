# Design

## Context

The repository root currently contains a Makefile with three independent
targets: `generate`, `test`, and `run`. The commands are simple wrappers around
the existing Go and templ tooling, so the migration can remain repository-local
and does not require changes to application code or build artifacts.

## Goals / Non-Goals

**Goals:**

- Provide a valid Taskfile v3 at the repository root.
- Keep the three existing developer workflows available under the same task
  names and with the same underlying commands.
- Make the README the source of truth for invoking the workflows.
- Allow the old Makefile to be removed without changing runtime behavior.

**Non-Goals:**

- Add task dependencies, caching, parallel execution, or environment-specific
  variants.
- Change Go module dependencies, templ generation output, test selection, or
  example-server behavior.
- Add a compatibility wrapper that keeps Make as a required tool.

## Decisions

### Use a root `Taskfile.yml` with Taskfile schema v3

The migration will use the conventional `Taskfile.yml` filename and declare
`version: '3'`. Each task will contain one `cmds` entry matching its current
Makefile command:

- `generate`: `templ generate ./...`
- `test`: `go test ./...`
- `run`: `go run ./cmd/example`

This keeps task discovery predictable and makes the migration a direct
replacement rather than a redesign of the developer workflow. A generated
Taskfile or a custom shell script was considered, but would add indirection
without providing a benefit for three stable commands.

### Remove the Makefile instead of maintaining dual entry points

The Makefile will be deleted after the Taskfile and documentation are updated.
Keeping both would obscure the supported workflow and allow the two definitions
to drift. The previous Make targets remain available as Task tasks with the same
names, so callers only need to change the launcher from `make` to `task`.

### Document installation and invocation in README.md

The README will add the Task installation prerequisite and show `task generate`
and `task run` in the example startup instructions, while documenting
`task test` for validation. The existing direct commands may remain as
implementation detail only where useful, but the supported entry point will be
the Taskfile.

## Risks / Trade-offs

- [Risk] Contributors without Task installed cannot use the documented
  shortcuts. → Document the prerequisite and retain the underlying Go/templ
  commands in the repository guidance where useful.
- [Risk] Taskfile syntax or version differences could make tasks fail to load.
  → Validate the file with the installed Task CLI (or an equivalent parser)
  before considering the migration complete, then run the generate and test
  tasks.
- [Risk] Removing Make may break undocumented automation invoking `make`.
  → Treat the removal as an intentional tooling migration, call it out in the
  proposal/documentation, and preserve task names and command bodies for a
  straightforward automation update.

## Migration Plan

1. Add `Taskfile.yml`, update README examples and prerequisites, and remove the
   Makefile in one change.
2. Verify task listing/parsing, run `task generate`, and run `task test`.
3. If rollback is required before downstream automation is updated, restore
   the Makefile from the previous revision and revert the README launcher
   changes; no application or data migration is needed.
