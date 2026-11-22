package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func postForm(t *testing.T, mux http.Handler, path string, values url.Values) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestExampleCRUDFlows(t *testing.T) {
	store := NewUserStore()
	mux := NewMux(store)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/users", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Ada Lovelace") {
		t.Fatalf("list failed: status=%d", rec.Code)
	}

	rec = postForm(t, mux, "/admin/users/new", url.Values{"name": {"Katherine Johnson"}, "email": {"katherine@example.com"}})
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create status=%d", rec.Code)
	}
	if len(store.List()) != 3 {
		t.Fatalf("create did not add a user")
	}

	rec = postForm(t, mux, "/admin/users/new", url.Values{"name": {""}, "email": {"invalid"}})
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Name is required") {
		t.Fatalf("validation failed: status=%d", rec.Code)
	}

	rec = postForm(t, mux, "/admin/users/1/edit", url.Values{"name": {"Ada Byron"}, "email": {"ada@example.com"}})
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update status=%d", rec.Code)
	}
	user, ok := store.Get(1)
	if !ok || user.Name != "Ada Byron" {
		t.Fatal("update did not persist")
	}

	rec = postForm(t, mux, "/admin/users/1/delete", nil)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("delete status=%d", rec.Code)
	}
	if _, ok := store.Get(1); ok {
		t.Fatal("delete did not remove the user")
	}
}

func TestExampleCancelAndMissingEntry(t *testing.T) {
	store := NewUserStore()
	mux := NewMux(store)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/users/1/delete", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Confirm delete") ||
		!strings.Contains(rec.Body.String(), `href="/admin/users"`) ||
		!strings.Contains(rec.Body.String(), "btn-danger") {
		t.Fatal("delete confirmation or cancel link missing")
	}

	rec = postForm(t, mux, "/admin/users/999/edit", url.Values{"name": {"Nobody"}, "email": {"nobody@example.com"}})
	if rec.Code != http.StatusNotFound {
		t.Fatalf("missing edit status=%d", rec.Code)
	}
}
