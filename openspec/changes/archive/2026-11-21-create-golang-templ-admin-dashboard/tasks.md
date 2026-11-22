# Tasks

## 1. Module and Tooling Setup

- [x] 1.1 Create the Go module and package directories for reusable components and the example application; verify `go list ./...` succeeds
- [x] 1.2 Add templ generation/build tooling and repository run instructions; verify templ generation succeeds from a clean checkout
- [x] 1.3 Add baseline Go tests for package initialization and module layout; verify `go test ./...` passes

## 2. Dashboard Shell

- [x] 2.1 Implement configurable navigation, active-route, and responsive sidebar templ components; verify component tests cover desktop and narrow rendering inputs
- [x] 2.2 Implement the reusable page shell and table components, including empty-state rendering; verify templ tests cover headers, rows, and no-row output
- [x] 2.3 Document sidebar and table integration for host applications; verify the README example compiles or is covered by an example package test

## 3. Resource CRUD Components

- [x] 3.1 Define host-owned resource handler contracts and typed resource form/table configuration; verify Go tests cover success and validation result mapping
- [x] 3.2 Implement create and edit templ pages plus server-facing form submission helpers; verify tests cover valid input, invalid input, and missing-entry behavior
- [x] 3.3 Implement confirmed delete behavior and test cancellation and successful deletion outcomes
- [x] 3.4 Document route and handler integration without prescribing storage or authentication; verify the example application demonstrates every documented contract

## 4. Example Application

- [x] 4.1 Create a runnable example server with a representative resource and thread-safe in-memory persistence; verify `go run ./cmd/example` starts successfully
- [x] 4.2 Wire example list, create, edit, and delete routes to the reusable components; verify integration tests cover create, update, delete, cancel, and validation flows
- [x] 4.3 Add example navigation, sidebar, table, form, and empty-state screens; verify the example renders each state through HTTP tests
- [x] 4.4 Add example startup instructions and identify in-memory persistence as demonstration-only; verify the instructions start the app from a clean checkout

## 5. Integration Verification

- [x] 5.1 Run templ generation, the full Go test suite, and the example HTTP integration suite; verify no generated or manual files are stale
- [x] 5.2 Review public APIs and documentation for host reuse across projects; verify a second minimal host package can import and render the library components

## Workflow follow-up

- Archive the change after implementation and review are complete.
