package admin

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLayoutAndNavigation(t *testing.T) {
	component := Layout("Users", []NavItem{
		{Label: "Users", Href: "/admin/users"},
		{Label: "Profile", Href: "/admin/profile"},
	}, "/admin/users", Table(TableConfig{
		Columns:   []TableColumn{{Key: "name", Label: "Name"}},
		Rows:      nil,
		EmptyText: "No users found",
	}))

	var out bytes.Buffer
	if err := component.Render(context.Background(), &out); err != nil {
		t.Fatalf("render layout: %v", err)
	}
	html := out.String()
	for _, want := range []string{
		"<aside",
		"nav-pills",
		`aria-current="page"`,
		"No users found",
		"viewport",
		"bootstrap@5.3.8",
		"<details",
		"d-none d-md-block",
		"text-dark",
		`>Users</a>`,
		`>Profile</a>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("rendered layout missing %q", want)
		}
	}
	if strings.Contains(html, "text-white") {
		t.Error("rendered layout unexpectedly contains text-white on light sidebar")
	}
	if strings.Contains(html, "Users 1") {
		t.Error("rendered layout contains debug suffix in active nav label")
	}
}

func TestTableRowsAndActions(t *testing.T) {
	component := Table(TableConfig{
		Title:       "Users",
		CreateHref:  "/users/new",
		CreateLabel: "Add User",
		Columns:     []TableColumn{{Key: "name", Label: "Name"}},
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
	for _, want := range []string{
		`<h1 class="h2 mb-0">Users</h1>`,
		`href="/users/new"`,
		"Add User",
		"table-responsive",
		`<th scope="col">Name</th>`,
		"Ada",
		"btn-outline-secondary",
		"btn-danger",
		"table-striped",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("rendered table missing %q", want)
		}
	}
}

func TestResourceFormDefaultsAndDeterministicErrors(t *testing.T) {
	component := ResourceForm(FormConfig{
		Title:  "User",
		Action: "/users/new",
		Fields: []FormField{
			{Name: "name", Label: "Name", Value: "", Required: true},
			{Name: "email", Label: "Email", Value: "bad", Type: "email"},
		},
		Errors: map[string]string{
			"name":  "Name is required",
			"email": "Email is invalid",
		},
	})

	var out bytes.Buffer
	if err := component.Render(context.Background(), &out); err != nil {
		t.Fatalf("render form: %v", err)
	}
	html := out.String()
	for _, want := range []string{
		`method="post"`,
		`action="/users/new"`,
		"is-invalid",
		"invalid-feedback",
		`aria-invalid`,
		"Name is required",
		"Email is invalid",
		">Save</button>",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("rendered form missing %q", want)
		}
	}

	// Alert error messages must appear in deterministic alphabetical key order (email before name).
	emailIdx := strings.Index(html, "<div>Email is invalid</div>")
	nameIdx := strings.Index(html, "<div>Name is required</div>")
	if emailIdx == -1 || nameIdx == -1 || emailIdx >= nameIdx {
		t.Fatalf("expected deterministic error alert order (email before name), got emailIdx=%d nameIdx=%d", emailIdx, nameIdx)
	}
}

func TestDeleteConfirmationDefaultsAndUnavailable(t *testing.T) {
	var delBuf bytes.Buffer
	if err := DeleteConfirmation(DeleteConfig{
		ItemLabel: "Ada Lovelace",
		Action:    "/users/1/delete",
	}).Render(context.Background(), &delBuf); err != nil {
		t.Fatalf("render delete: %v", err)
	}
	delHTML := delBuf.String()
	if !strings.Contains(delHTML, "Confirm Delete") || !strings.Contains(delHTML, "Ada Lovelace") {
		t.Fatalf("unexpected delete confirmation output: %s", delHTML)
	}
	if strings.Contains(delHTML, ">Cancel</a>") {
		t.Fatal("expected cancel link to be omitted when CancelHref is empty")
	}

	var unavailBuf bytes.Buffer
	if err := Unavailable(UnavailableConfig{
		BackHref:  "/users",
		BackLabel: "Back to Users",
	}).Render(context.Background(), &unavailBuf); err != nil {
		t.Fatalf("render unavailable: %v", err)
	}
	unavailHTML := unavailBuf.String()
	for _, want := range []string{"Entry Unavailable", "no longer available", `href="/users"`, "Back to Users"} {
		if !strings.Contains(unavailHTML, want) {
			t.Errorf("rendered unavailable missing %q", want)
		}
	}
}

func TestSubmissionResultAndResourceHandlers(t *testing.T) {
	handlers := ResourceHandlers{
		Get: func(id string) (TableRow, bool) {
			if id == "1" {
				return TableRow{ID: "1", Values: map[string]string{"name": "Ada"}}, true
			}
			return TableRow{}, false
		},
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

	row, ok := handlers.Get("1")
	if !ok || row.Values["name"] != "Ada" {
		t.Fatalf("expected row from Get handler, got %+v ok=%v", row, ok)
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

func TestRegisterResourceLifecycle(t *testing.T) {
	items := map[string]string{"1": "Widget"}
	mux := http.NewServeMux()
	// Intentionally omit Get to verify fallback to List works.
	RegisterResource(mux, ResourceConfig{
		BasePath:      "/admin/items",
		Title:         "Items",
		SingularTitle: "Item",
		Nav:           []NavItem{{Label: "Items", Href: "/admin/items"}},
		Columns:       []TableColumn{{Key: "name", Label: "Name"}},
		EmptyText:     "No items",
		CreateLabel:   "New Item",
		Fields:        []FormField{{Name: "name", Label: "Name", Required: true}},
		Handlers: ResourceHandlers{
			List: func() ([]TableRow, error) {
				rows := make([]TableRow, 0, len(items))
				for id, name := range items {
					rows = append(rows, TableRow{ID: id, Values: map[string]string{"name": name}})
				}
				return rows, nil
			},
			Create: func(values map[string]string) SubmissionResult {
				if values["name"] == "" {
					return SubmissionResult{Errors: map[string]string{"name": "Name is required"}}
				}
				items["2"] = values["name"]
				return SubmissionResult{Success: true}
			},
			Update: func(id string, values map[string]string) SubmissionResult {
				if _, ok := items[id]; !ok {
					return SubmissionResult{NotFound: true}
				}
				if values["name"] == "" {
					return SubmissionResult{Errors: map[string]string{"name": "Name is required"}}
				}
				items[id] = values["name"]
				return SubmissionResult{Success: true}
			},
			Delete: func(id string) bool {
				if _, ok := items[id]; !ok {
					return false
				}
				delete(items, id)
				return true
			},
		},
	})

	// List
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/items", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Widget") || !strings.Contains(rec.Body.String(), "New Item") {
		t.Fatalf("list failed: status=%d", rec.Code)
	}

	// Create form & validation error
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/items/new", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Add Item") {
		t.Fatalf("new form failed: status=%d", rec.Code)
	}

	form := url.Values{"name": {""}}
	req := httptest.NewRequest(http.MethodPost, "/admin/items/new", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Name is required") {
		t.Fatalf("expected validation error, status=%d", rec.Code)
	}

	// Create valid
	form = url.Values{"name": {"Gadget"}}
	req = httptest.NewRequest(http.MethodPost, "/admin/items/new", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther || items["2"] != "Gadget" {
		t.Fatalf("create failed: status=%d items=%v", rec.Code, items)
	}

	// Edit GET pre-fills value via List fallback
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/items/1/edit", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `value="Widget"`) {
		t.Fatalf("edit form failed: status=%d", rec.Code)
	}

	// Edit missing GET returns 404 Unavailable inside Layout
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/items/missing/edit", nil))
	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "Entry Unavailable") {
		t.Fatalf("missing edit GET failed: status=%d", rec.Code)
	}

	// Delete GET & POST
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/items/1/delete", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Widget") {
		t.Fatalf("delete confirmation failed: status=%d", rec.Code)
	}

	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/admin/items/1/delete", nil))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("delete POST failed: status=%d", rec.Code)
	}
	if _, ok := items["1"]; ok {
		t.Fatal("expected item 1 to be deleted")
	}

	// Delete missing POST returns 404 Unavailable
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/admin/items/1/delete", nil))
	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "Entry Unavailable") {
		t.Fatalf("missing delete POST failed: status=%d", rec.Code)
	}
}

func TestRegisterResourceListError(t *testing.T) {
	mux := http.NewServeMux()
	RegisterResource(mux, ResourceConfig{
		BasePath: "/admin/broken",
		Title:    "Broken",
		Handlers: ResourceHandlers{
			List: func() ([]TableRow, error) {
				return nil, errors.New("database offline")
			},
		},
	})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/broken", nil))
	if rec.Code != http.StatusInternalServerError || !strings.Contains(rec.Body.String(), "database offline") {
		t.Fatalf("expected 500 error page, got status=%d body=%s", rec.Code, rec.Body.String())
	}
}
