package main

import (
	"bytes"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/a-h/templ"
	"github.com/rosso0815/rosso0815-dashboard/admin"
)

type User struct {
	ID    int
	Name  string
	Email string
}

// UserStore provides demonstration-only thread-safe in-memory persistence.
type UserStore struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: map[int]User{
			1: {ID: 1, Name: "Ada Lovelace", Email: "ada@example.com"},
			2: {ID: 2, Name: "Grace Hopper", Email: "grace@example.com"},
		},
		nextID: 3,
	}
}

func (s *UserStore) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })
	return users
}

func (s *UserStore) Get(id int) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[id]
	return user, ok
}

func (s *UserStore) Create(name, email string) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	user := User{ID: s.nextID, Name: name, Email: email}
	s.users[user.ID] = user
	s.nextID++
	return user
}

func (s *UserStore) Update(id int, name, email string) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return User{}, false
	}
	user := User{ID: id, Name: name, Email: email}
	s.users[id] = user
	return user, true
}

func (s *UserStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return false
	}
	delete(s.users, id)
	return true
}

var nav = []admin.NavItem{
	{Label: "Users", Href: "/admin/users"},
	{Label: "Profile", Href: "/admin/profile"},
}

func NewMux(store *UserStore) *http.ServeMux {
	mux := http.NewServeMux()

	admin.RegisterResource(mux, admin.ResourceConfig{
		BasePath:      "/admin/users",
		Title:         "Users",
		SingularTitle: "User",
		Nav:           nav,
		CurrentPath:   "/admin/users",
		Columns: []admin.TableColumn{
			{Key: "name", Label: "Name"},
			{Key: "email", Label: "Email"},
		},
		EmptyText:   "No users found",
		CreateLabel: "Add User",
		Fields: []admin.FormField{
			{Name: "name", Label: "Name", Required: true},
			{Name: "email", Label: "Email", Type: "email", Required: true},
		},
		ItemLabel: func(row admin.TableRow) string {
			return row.Values["name"]
		},
		Handlers: admin.ResourceHandlers{
			List: func() ([]admin.TableRow, error) {
				users := store.List()
				rows := make([]admin.TableRow, 0, len(users))
				for _, user := range users {
					rows = append(rows, userToRow(user))
				}
				return rows, nil
			},
			Get: func(rawID string) (admin.TableRow, bool) {
				id, ok := parseID(rawID)
				if !ok {
					return admin.TableRow{}, false
				}
				user, ok := store.Get(id)
				if !ok {
					return admin.TableRow{}, false
				}
				return userToRow(user), true
			},
			Create: func(values map[string]string) admin.SubmissionResult {
				name, email, errors := validateUserValues(values)
				if len(errors) > 0 {
					return admin.SubmissionResult{Errors: errors}
				}
				store.Create(name, email)
				return admin.SubmissionResult{Success: true}
			},
			Update: func(rawID string, values map[string]string) admin.SubmissionResult {
				id, ok := parseID(rawID)
				if !ok {
					return admin.SubmissionResult{NotFound: true}
				}
				name, email, errors := validateUserValues(values)
				if len(errors) > 0 {
					return admin.SubmissionResult{Errors: errors}
				}
				if _, ok := store.Update(id, name, email); !ok {
					return admin.SubmissionResult{NotFound: true}
				}
				return admin.SubmissionResult{Success: true}
			},
			Delete: func(rawID string) bool {
				id, ok := parseID(rawID)
				if !ok {
					return false
				}
				return store.Delete(id)
			},
		},
	})

	mux.HandleFunc("GET /admin/profile", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, admin.Layout("Profile", nav, "/admin/profile", ProfilePage()))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
			return
		}
		http.NotFound(w, r)
	})

	return mux
}

func userToRow(user User) admin.TableRow {
	return admin.TableRow{
		ID: strconv.Itoa(user.ID),
		Values: map[string]string{
			"name":  user.Name,
			"email": user.Email,
		},
	}
}

func parseID(raw string) (int, bool) {
	id, err := strconv.Atoi(raw)
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

func validateUserValues(values map[string]string) (string, string, map[string]string) {
	name := strings.TrimSpace(values["name"])
	email := strings.TrimSpace(values["email"])
	errors := map[string]string{}
	if name == "" {
		errors["name"] = "Name is required"
	}
	if email == "" {
		errors["email"] = "Email is required"
	} else if !strings.Contains(email, "@") {
		errors["email"] = "Email must contain @"
	}
	return name, email, errors
}

func render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	var buf bytes.Buffer
	if err := component.Render(r.Context(), &buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}

func main() {
	store := NewUserStore()
	log.Println("Admin dashboard listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", NewMux(store)))
}
