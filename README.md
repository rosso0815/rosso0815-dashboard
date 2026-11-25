# Golang Templ Admin Dashboard

A reusable server-rendered admin dashboard library for Go applications using
[templ](https://templ.guide).

## Run the example

Install [Task](https://taskfile.dev/docs/installation) and the templ CLI
before running the example. The example uses demonstration-only in-memory
persistence:

```sh
task generate
task run
```

Open `http://localhost:8080/admin/users`.

Run the complete test suite with:

```sh
task test
```

## Styling

The dashboard components are styled with Bootstrap 5.3. Bootstrap CSS is loaded
by the shared `admin.Layout` component. The components do not require Bootstrap
JavaScript or a JavaScript framework; standard templ markup, native HTML
elements (such as `<details>` for narrow-viewport navigation collapse), and Go
HTTP routes provide the interactions.

## Library integration

The reusable components live in `admin`. Import it from a host application,
provide navigation and table data, and keep persistence and validation in the
host:

```go
admin.Layout("Example", []admin.NavItem{
	{Label: "Users", Href: "/admin/users"},
}, "/admin/users", admin.Table(admin.TableConfig{
	Title:     "Users",
	Columns:   []admin.TableColumn{{Key: "name", Label: "Name"}},
	Rows: []admin.TableRow{{
		ID:     "1",
		Values: map[string]string{"name": "Ada"},
	}},
	EmptyText:   "No users found",
	CreateHref:  "/admin/users/new",
	CreateLabel: "Add User",
}))
```

The library does not require authentication, an ORM, or a JavaScript
framework. Host routes are responsible for authorization and persistence.

For full create/edit/delete route wiring backed by your own persistence and
validation, supply `admin.ResourceHandlers` to `admin.RegisterResource`:

```go
admin.RegisterResource(mux, admin.ResourceConfig{
	BasePath:      "/admin/users",
	Title:         "Users",
	SingularTitle: "User",
	Nav:           nav,
	Columns:       []admin.TableColumn{{Key: "name", Label: "Name"}},
	Fields:        []admin.FormField{{Name: "name", Label: "Name", Required: true}},
	Handlers:      handlers,
})
```

You can also compose `admin.ResourceForm`, `admin.DeleteConfirmation`, and
`admin.Unavailable` directly inside your own HTTP routes; the example under
`cmd/example` demonstrates both `admin.RegisterResource` and a standalone custom
page (`/admin/profile`).
