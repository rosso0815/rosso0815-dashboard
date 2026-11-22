# Spec Delta

## Purpose

Provides reusable templ components for a responsive admin dashboard shell, sidebar navigation, and consistent resource tables.

## ADDED Requirements

### Requirement: Dashboard shell layout
The dashboard SHALL provide a main content area and a sidebar that adapts between desktop and narrow viewports.

#### Scenario: Desktop viewport
- **WHEN** the dashboard is rendered on a desktop viewport
- **THEN** the sidebar and main content are shown side by side

#### Scenario: Narrow viewport
- **WHEN** the dashboard is rendered on a narrow viewport
- **THEN** the sidebar can be collapsed without hiding the primary page content

### Requirement: Configurable sidebar navigation
The dashboard SHALL render navigation items supplied by the host application.

#### Scenario: Navigation item supplied
- **WHEN** the host application provides a navigation item with a label and destination
- **THEN** the sidebar displays it as a selectable link

#### Scenario: Active navigation item
- **WHEN** a current route matches a navigation destination
- **THEN** the corresponding navigation item is marked active

### Requirement: Reusable table presentation
The dashboard SHALL render tabular resource data from column and row configuration supplied by the host application.

#### Scenario: Table columns rendered
- **WHEN** a resource table is configured with columns
- **THEN** the table renders one header and one value for each configured column

#### Scenario: No rows available
- **WHEN** a resource table has no rows
- **THEN** the table displays an explicit empty state
