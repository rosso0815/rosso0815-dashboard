# Design

## Context

See `proposal.md` for motivation. The repository is currently empty, so the design establishes the module boundary and example layout without an existing implementation to preserve.

## Goals / Non-Goals

**Goals:**
- Provide composable templ components with typed Go configuration.
- Keep storage, routing, authentication, and validation under host-application control.
- Support server-rendered add, edit, and delete flows with no required JavaScript framework.
- Include a runnable example that uses in-memory representative data.
- Make the public API small enough to reuse across projects.

**Non-Goals:**
- A database abstraction or migration system.
- Authentication and authorization.
- A client-side SPA or framework-specific integration.
- Generic form generation for every Go field type.

## Decisions

### Separate reusable components from resource behavior
The library will provide templ layout, table, form, and navigation components plus small interfaces or handler contracts. Host applications supply row data, handlers, routes, and validation.

Alternatives considered: A generic reflection-based CRUD engine would reduce initial integration code but would couple the library to model tags and field behavior. Route-level components would offer more convenience but less reuse across frameworks.

### Use server-rendered templ pages and standard Go HTTP handlers
List, create, edit, and delete operations will submit standard HTTP requests and receive rendered templ responses or redirects. This works with the standard `net/http` ecosystem and avoids a required JavaScript runtime.

Alternatives considered: HTMX would improve partial updates but add a frontend dependency. A React or JSON API integration would expand the scope beyond templ.

### Keep example persistence in memory
The example application will implement the resource contracts over a thread-safe in-memory store. This demonstrates CRUD behavior without imposing a database dependency on library users.

Alternatives considered: SQLite would provide restart persistence but add a driver and schema setup to the example. PostgreSQL or a generic repository interface would require configuration and obscure the reusable integration.

### Define resource views with typed column and form configuration
Tables will use explicit column definitions mapping headers to values. CRUD forms will compose ordinary templ fields or host-provided field components, allowing projects to use their own validation and markup.

Alternatives considered: Automatic reflection-based forms would be harder to control and less portable across existing models.

### Organize the module by presentation and example boundaries
The reusable package will remain independent of `cmd/example`. The example will be a separate importable package layout so users can copy or adapt its wiring.

## Risks / Trade-offs

- [The component API may need refinement across real projects] → Keep the initial API host-neutral and add integration examples before expanding abstractions.
- [Server-rendered full-page responses are slower than partial updates] → Standard redirects and rendered pages keep behavior simple; partial templ rendering can be added later without changing storage contracts.
- [An in-memory example may be mistaken for durable storage] → Document the store as demonstration-only and keep it isolated from library packages.
- [No built-in authorization could expose host routes] → Explicitly leave authentication and authorization to the host and document route protection as integration responsibility.

## Migration Plan

This is a new greenfield module. Implement the reusable package and example together, verify both with tests and templ generation, then publish or consume as an ordinary Go module.

Rollback is limited to removing or reverting the initial module and example files because no existing persisted data or deployed API is involved.

## Open Questions

- No open questions. Project-specific CSS conventions and persistence engines will remain host choices.
