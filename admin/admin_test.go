package admin

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestLayoutAndNavigation(t *testing.T) {
	component := Layout("Users", []NavItem{{Label: "Users", Href: "/admin/users"}}, "/admin/users", Table(TableConfig{
		Columns:   []TableColumn{{Key: "name", Label: "Name"}},
		Rows:      nil,
		EmptyText: "No users found",
	}))

	var out bytes.Buffer
	if err := component.Render(context.Background(), &out); err != nil {
		t.Fatalf("render layout: %v", err)
	}
	html := out.String()
	for _, want := range []string{"<aside", "nav-pills", "aria-current=\"page\"", "No users found", "viewport", "bootstrap@5.3.8"} {
		if !strings.Contains(html, want) {
			t.Errorf("rendered layout missing %q", want)
		}
	}
}

func TestTableRowsAndActions(t *testing.T) {
	component := Table(TableConfig{
		Columns: []TableColumn{{Key: "name", Label: "Name"}},
		Rows: []TableRow{{
			ID:     "1",
			Values: map[string]string{"name": "Ada"},
		}},
		EditHref:   func(id string) string { return "/users/" + id + "/edit" },
		DeleteHref: func(id string) string { return "/users/" + id + "/delete" },
	})

	var out bytes.Buffer
	if err := component.Render(context.Background(), &out); err != nil {
		t.Fatalf("render table: %v", err)
	}
	html := out.String()
	for _, want := range []string{"<th scope=\"col\">Name</th>", "Ada", "btn-outline-secondary", "btn-danger", "table-striped"} {
		if !strings.Contains(html, want) {
			t.Errorf("rendered table missing %q", want)
		}
	}
}

func TestSubmissionResultAndResourceHandlers(t *testing.T) {
	handlers := ResourceHandlers{
		Create: func(values map[string]string) SubmissionResult {
			if values["name"] == "" {
				return SubmissionResult{Errors: map[string]string{"name": "Name is required"}}
			}
			return SubmissionResult{Success: true}
		},
		Update: func(id string, values map[string]string) SubmissionResult {
			if id == "missing" {
				return SubmissionResult{NotFound: true}
			}
			return SubmissionResult{Success: true}
		},
		Delete: func(id string) bool { return id == "1" },
	}

	invalid := handlers.Create(map[string]string{})
	if invalid.Success || invalid.Errors["name"] == "" {
		t.Fatalf("expected validation result, got %+v", invalid)
	}
	missing := handlers.Update("missing", map[string]string{"name": "Ada"})
	if !missing.NotFound {
		t.Fatalf("expected not-found result, got %+v", missing)
	}
	if !handlers.Delete("1") || handlers.Delete("missing") {
		t.Fatal("delete handler did not preserve host outcomes")
	}
}
