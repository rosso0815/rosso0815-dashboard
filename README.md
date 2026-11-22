# Golang Templ Admin Dashboard

A reusable server-rendered admin dashboard library for Go applications using
[templ](https://templ.guide).

## Run the example

The example uses demonstration-only in-memory persistence:

```sh
templ generate ./...
go run ./cmd/example
```

Open `http://localhost:8080/admin/users`.

## Styling

The dashboard components are styled with Bootstrap 5.3. Bootstrap CSS is loaded
by the shared `admin.Layout` component. The components do not require Bootstrap
JavaScript or a JavaScript framework; standard templ markup and Go HTTP routes
continue to provide the CRUD interactions.

## Library integration

The reusable components live in `admin`. Import it from a host application,
provide navigation and table data, and keep persistence and validation in the
host:

```go
admin.Layout("Example", []admin.NavItem{
	{Label: "Users", Href: "/admin/users"},
}, "/admin/users", admin.Table(admin.TableConfig{
	Columns: []admin.TableColumn{{Key: "name", Label: "Name"}},
	Rows: []admin.TableRow{{
		ID: "1",
		Values: map[string]string{"name": "Ada"},
	}},
	EmptyText: "No users found",
}))
```

The library does not require authentication, an ORM, or a JavaScript
framework. Host routes are responsible for authorization and persistence.

For create/edit/delete handlers, implement the host-owned
`admin.ResourceHandlers` contract or call the store operations from your HTTP
routes. Resource forms and delete confirmations are reusable templ components;
the example under `cmd/example` shows the complete route wiring.
