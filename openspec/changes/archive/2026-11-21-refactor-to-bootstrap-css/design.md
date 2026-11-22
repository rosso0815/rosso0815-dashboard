# Design

## Context

See `proposal.md` for motivation. The current UI is server-rendered templ markup with one inline stylesheet in `admin/layout.templ`; component templates share custom semantic classes for the shell, navigation, cards, tables, forms, actions, and alerts. The example app consumes these components and has render tests that assert selected HTML output.

## Goals / Non-Goals

**Goals:**

- Use Bootstrap 5 utility and component classes for the existing dashboard shell and resource screens.
- Keep the reusable Go and templ APIs stable so host applications do not need route or configuration changes.
- Preserve responsive navigation, table overflow, form validation feedback, empty states, and CRUD interactions.
- Keep styling centralized in the shared layout and avoid duplicating framework setup across pages.

**Non-Goals:**

- Adding Bootstrap JavaScript, a JavaScript framework, or interactive behavior that does not exist today.
- Introducing an asset bundler, custom theme system, or host-configurable stylesheet API.
- Changing persistence, routing, authentication, validation, or resource handler contracts.

## Decisions

### Load a pinned Bootstrap stylesheet from the shared layout

Add a version-pinned Bootstrap 5.3 CSS link in `admin/layout.templ`, keeping the dependency centralized for every rendered page. The current repository has no frontend asset pipeline, so a CDN stylesheet avoids adding build tooling or generated assets while retaining the server-rendered integration model.

Alternatives considered: Vendoring Bootstrap would improve offline determinism but substantially increase repository size and require an asset-serving convention. An unpinned CDN URL would simplify updates but make rendered output change without a planned dependency update.

### Map existing layout responsibilities to Bootstrap classes

Replace custom classes with Bootstrap grid, spacing, flex, navigation, card, table, form-control, button, alert, and responsive utility classes. Keep small semantic classes only where the existing component structure needs a stable hook not provided by Bootstrap, and preserve the current markup semantics and templ control flow.

Alternatives considered: Retaining the custom stylesheet alongside Bootstrap would create competing rules and make the migration incomplete. Rewriting the page structure around Bootstrap examples would risk changing host-visible markup and component behavior unnecessarily.

### Keep the sidebar responsive with Bootstrap utilities

Use Bootstrap’s responsive display and sizing utilities to preserve a desktop sidebar plus main content and a narrow-screen presentation where the sidebar remains compact without hiding content. Do not add a collapse toggle because that would introduce new behavior and require Bootstrap JavaScript.

Alternatives considered: Bootstrap’s collapse component would provide a richer mobile navigation, but it changes behavior and adds a JavaScript dependency outside this refactor’s scope. A custom media-query implementation would preserve the current result but defeat the purpose of consolidating layout styling.

### Update assertions around stable behavior and framework markup

Adjust component and example tests to assert stable Bootstrap markers and existing semantics, such as the viewport metadata, `aria-current`, table headers, validation text, and Bootstrap class hooks. Avoid asserting the full stylesheet or incidental class ordering so future Bootstrap-compatible markup changes remain low-cost.

Alternatives considered: Removing all HTML assertions would reduce maintenance but would no longer protect the rendered integration contract. Snapshotting complete pages would be brittle for a class-name migration.

## Risks / Trade-offs

- [CDN availability affects styling in offline or restricted environments] → Pin the Bootstrap version, document the external stylesheet choice, and keep the HTML functional without relying on Bootstrap JavaScript.
- [Bootstrap defaults may not exactly match the current visual design] → Preserve layout semantics and add only minimal scoped custom CSS where a direct utility mapping cannot retain required responsive behavior.
- [Class names become part of the rendered integration surface] → Test stable framework hooks and preserve public Go APIs; treat visual class changes as an intentional presentation migration.

## Migration Plan

1. Replace the shared inline CSS and update all affected templ markup to Bootstrap classes.
2. Update the example page markup and render assertions.
3. Run templ generation and the Go test suite, then inspect representative rendered HTML for desktop, narrow, empty, validation, and confirmation states.

Rollback is a source-only revert: restore the current inline stylesheet and custom classes. No data migration, API migration, or deployment sequencing is required.

## Open Questions

None.
