# Proposal

## Why

Projects need a consistent Go-based admin interface without rebuilding sidebar navigation, table layouts, and create/edit/delete flows. This change establishes a reusable templ component library and keeps a working example application in sync so the library can be adopted and verified across projects.

## What Changes

- Add a reusable admin dashboard shell with sidebar navigation, page layout, table views, and responsive behavior.
- Add reusable resource CRUD components and server integration points for listing, creating, editing, and deleting entries.
- Provide an example application that exercises the reusable components with a representative resource.
- Keep the example application as an executable demonstration and verification fixture whenever the library changes.
- Define clear extension points for project-specific resources, storage, routes, validation, and presentation.

## Capabilities

### New Capabilities
- `admin-dashboard-shell`: Sidebar navigation, dashboard page layout, and reusable table presentation.
- `admin-resource-crud`: Reusable create, update, and delete flows for tabular resource entries.
- `admin-example-app`: An maintained example application demonstrating library usage and CRUD behavior.

### Modified Capabilities
- None.

## Impact

- New Go module and templ component packages.
- New example application and likely local development run configuration.
- Public component APIs and route integration contracts introduced for reuse across projects.
- Standard Go tests and templ generation/build tooling become part of the project workflow.
