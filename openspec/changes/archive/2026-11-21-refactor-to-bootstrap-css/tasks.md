# Tasks

## 1. Bootstrap Foundation

- [x] 1.1 Add a version-pinned Bootstrap 5.3 stylesheet link to `admin/layout.templ`, remove the obsolete inline custom stylesheet, and verify the generated layout contains the stylesheet reference and existing viewport metadata.
- [x] 1.2 Update the shared shell, sidebar, navigation, page header, card, and responsive layout markup to Bootstrap classes while preserving active-link semantics; verify `TestLayoutAndNavigation` covers the rendered shell and `aria-current` output.

## 2. Component Markup Migration

- [x] 2.1 Convert table, action, empty-state, form, validation-alert, and delete-confirmation markup in `admin/` to Bootstrap table, button, form, alert, spacing, and responsive utility classes; verify component rendering still includes headers, row values, actions, validation messages, and confirmation controls.
- [x] 2.2 Update `cmd/example/page.templ` to use the Bootstrap page-header and button classes without changing route links; verify the example page renders the expected users heading and add-user link.
- [x] 2.3 Update focused render assertions in `admin/admin_test.go` and `cmd/example/main_test.go` to check stable Bootstrap hooks plus existing semantic behavior; verify `go test ./admin ./cmd/example` passes.

## 3. Documentation and Integration

- [x] 3.1 Document the Bootstrap CSS dependency and the absence of a required Bootstrap JavaScript runtime in `README.md`; verify the documented example startup commands remain accurate.
- [x] 3.2 Run templ generation and the complete Go test suite, then inspect representative list, empty, validation, and delete-confirmation responses for correct Bootstrap markup; verify `templ generate ./... && go test ./...` succeeds.

## Workflow follow-up

- Archive the change after implementation and review are complete.
