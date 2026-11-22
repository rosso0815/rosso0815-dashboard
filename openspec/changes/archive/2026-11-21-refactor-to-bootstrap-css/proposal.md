# Proposal

## Why

The admin components currently embed a large custom CSS block for layout, navigation, tables, forms, and controls. Adopting Bootstrap will reduce bespoke styling, provide familiar responsive utility classes, and make the rendered components easier for host applications to extend while preserving the existing server-rendered behavior.

## What Changes

- Replace the inline custom stylesheet in the admin layout with Bootstrap CSS loaded through the layout.
- Refactor dashboard, navigation, table, form, delete-confirmation, and example-page markup to use Bootstrap classes.
- Preserve the existing component APIs, routes, HTML semantics, active-navigation behavior, empty states, form validation, and CRUD flows.
- Keep the example application and tests aligned with the updated class names and Bootstrap-backed markup.
- Do not add a JavaScript framework, change persistence, or introduce a required client-side interaction model.

## Capabilities

### New Capabilities

None. This is a presentation-only refactor.

### Modified Capabilities

None. Existing dashboard, CRUD, and example-app requirements remain unchanged.

## Impact

- Affects the templ files under `admin/` and the example page under `cmd/example/`.
- Updates rendered HTML class names and the source of styling, but does not change public Go configuration types or HTTP route contracts.
- Adds Bootstrap CSS as a frontend dependency delivered by the shared layout; the implementation should choose a stable, documented loading strategy consistent with the project’s server-rendered model.
- Requires focused updates to render assertions and a full templ-generation/test verification pass.
