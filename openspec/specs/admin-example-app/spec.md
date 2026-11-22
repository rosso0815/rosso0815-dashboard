# admin-example-app Specification

## Purpose

Keeps a runnable example application aligned with the reusable admin library and demonstrates dashboard and CRUD integration.

## Requirements

### Requirement: Runnable library example
The repository SHALL include an example application that can be started locally and used with the reusable dashboard components.

#### Scenario: Example application starts
- **WHEN** a developer follows the repository run instructions
- **THEN** the example application serves a usable admin dashboard

### Requirement: CRUD example coverage
The example application SHALL demonstrate create, edit, and delete operations for a representative resource.

#### Scenario: Example resource managed
- **WHEN** a user adds, edits, or deletes a row in the example resource
- **THEN** the resource list reflects the operation

### Requirement: Example remains synchronized
Library changes SHALL keep the example application compiling and exercising the affected public integration points.

#### Scenario: Library API changes
- **WHEN** a library integration API changes
- **THEN** the example application is updated to use the new API before completion
