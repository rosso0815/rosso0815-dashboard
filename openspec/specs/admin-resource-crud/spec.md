# admin-resource-crud Specification

## Purpose

Defines reusable add, edit, and delete behavior for resource entries managed through the admin dashboard.

## Requirements

### Requirement: Create resource entries
The resource interface SHALL allow a host application to present and submit an entry creation form.

#### Scenario: Valid entry submitted
- **WHEN** a user submits a complete and valid create form
- **THEN** the host application creates the entry and refreshes the resource list

#### Scenario: Invalid entry submitted
- **WHEN** a user submits an entry that fails host validation
- **THEN** the interface remains visible and presents the validation errors

### Requirement: Update resource entries
The resource interface SHALL allow a host application to load and submit changes for an existing entry.

#### Scenario: Existing entry edited
- **WHEN** a user opens an entry and submits valid changed values
- **THEN** the host application updates the entry and refreshes the resource list

#### Scenario: Entry no longer exists
- **WHEN** an edit request targets an entry that cannot be found
- **THEN** the interface reports that the entry is unavailable

### Requirement: Delete resource entries
The resource interface SHALL require confirmation before deleting an entry.

#### Scenario: Delete confirmed
- **WHEN** a user confirms deletion of an existing entry
- **THEN** the host application deletes the entry and removes it from the resource list

#### Scenario: Delete cancelled
- **WHEN** a user cancels deletion
- **THEN** the entry remains unchanged

### Requirement: Host-owned persistence and validation
The reusable components SHALL expose operations for resource persistence and validation without requiring a specific database or form framework.

#### Scenario: Host supplies resource handlers
- **WHEN** the host application provides list, create, update, and delete handlers
- **THEN** the reusable interface invokes those handlers for the corresponding operations
