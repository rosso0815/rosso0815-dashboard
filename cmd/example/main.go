package main

import (
	"fmt"
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

var nav = []admin.NavItem{{Label: "Users", Href: "/admin/users"}}

func NewMux(store *UserStore) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/users", func(w http.ResponseWriter, r *http.Request) {
		users := store.List()
		rows := make([]admin.TableRow, 0, len(users))
		for _, user := range users {
			rows = append(rows, admin.TableRow{
				ID: strconv.Itoa(user.ID),
				Values: map[string]string{
					"name":  user.Name,
					"email": user.Email,
				},
			})
		}
		table := admin.Table(admin.TableConfig{
			Columns: []admin.TableColumn{
				{Key: "name", Label: "Name"},
				{Key: "email", Label: "Email"},
			},
			Rows:       rows,
			EmptyText:  "No users found",
			CreateHref: "/admin/users/new",
			EditHref:   func(id string) string { return "/admin/users/" + id + "/edit" },
			DeleteHref: func(id string) string { return "/admin/users/" + id + "/delete" },
		})
		page := UserListPage(table)
		render(w, r, admin.Layout("Users", nav, r.URL.Path, page))
	})

	mux.HandleFunc("GET /admin/users/new", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, admin.Layout("Add User", nav, "/admin/users", admin.ResourceForm(userForm("", "", "/admin/users/new", nil))))
	})
	mux.HandleFunc("POST /admin/users/new", func(w http.ResponseWriter, r *http.Request) {
		name, email, errors := readUserForm(w, r)
		if len(errors) > 0 {
			render(w, r, admin.Layout("Add User", nav, "/admin/users", admin.ResourceForm(userForm(name, email, "/admin/users/new", errors))))
			return
		}
		store.Create(name, email)
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /admin/users/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		id, ok := routeID(w, r)
		if !ok {
			return
		}
		user, ok := store.Get(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		action := fmt.Sprintf("/admin/users/%d/edit", id)
		render(w, r, admin.Layout("Edit User", nav, "/admin/users", admin.ResourceForm(userForm(user.Name, user.Email, action, nil))))
	})
	mux.HandleFunc("POST /admin/users/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		id, ok := routeID(w, r)
		if !ok {
			return
		}
		name, email, errors := readUserForm(w, r)
		if len(errors) > 0 {
			action := fmt.Sprintf("/admin/users/%d/edit", id)
			render(w, r, admin.Layout("Edit User", nav, "/admin/users", admin.ResourceForm(userForm(name, email, action, errors))))
			return
		}
		if _, ok := store.Update(id, name, email); !ok {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /admin/users/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id, ok := routeID(w, r)
		if !ok {
			return
		}
		user, ok := store.Get(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		render(w, r, admin.Layout("Delete User", nav, "/admin/users", admin.DeleteConfirmation(admin.DeleteConfig{
			Title:      "Delete User",
			ItemLabel:  user.Name,
			Action:     fmt.Sprintf("/admin/users/%d/delete", id),
			CancelHref: "/admin/users",
		})))
	})
	mux.HandleFunc("POST /admin/users/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id, ok := routeID(w, r)
		if !ok {
			return
		}
		if !store.Delete(id) {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
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

func userForm(name, email, action string, errors map[string]string) admin.FormConfig {
	method := "post"
	return admin.FormConfig{
		Title:       "User",
		Action:      action,
		Method:      method,
		SubmitLabel: "Save",
		CancelHref:  "/admin/users",
		Errors:      errors,
		Fields: []admin.FormField{
			{Name: "name", Label: "Name", Value: name, Required: true},
			{Name: "email", Label: "Email", Value: email, Type: "email", Required: true},
		},
	}
}

func readUserForm(w http.ResponseWriter, r *http.Request) (string, string, map[string]string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return "", "", map[string]string{"form": "Invalid form submission"}
	}
	name := strings.TrimSpace(r.Form.Get("name"))
	email := strings.TrimSpace(r.Form.Get("email"))
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

func routeID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return 0, false
	}
	return id, true
}

func render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	store := NewUserStore()
	log.Println("Admin dashboard listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", NewMux(store)))
}
