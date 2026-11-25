package admin

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

// ResourceConfig wires host-owned ResourceHandlers into standard server-rendered CRUD routes.
type ResourceConfig struct {
	BasePath      string
	Title         string
	SingularTitle string
	Nav           []NavItem
	CurrentPath   string
	Columns       []TableColumn
	EmptyText     string
	CreateLabel   string
	Fields        []FormField
	ItemLabel     func(TableRow) string
	Handlers      ResourceHandlers
}

// RegisterResource mounts list, create, edit, and delete routes onto mux using config.Handlers.
func RegisterResource(mux *http.ServeMux, config ResourceConfig) {
	base := "/" + strings.Trim(config.BasePath, "/")
	if base == "/" {
		base = "/resource"
	}
	currentPath := config.CurrentPath
	if currentPath == "" {
		currentPath = base
	}
	singular := config.SingularTitle
	if singular == "" {
		singular = config.Title
	}
	if singular == "" {
		singular = "Entry"
	}
	listTitle := config.Title
	if listTitle == "" {
		listTitle = singular
	}

	mux.HandleFunc("GET "+base, func(w http.ResponseWriter, r *http.Request) {
		if config.Handlers.List == nil {
			renderStatus(w, r, http.StatusNotFound, Layout("Unavailable", config.Nav, currentPath, Unavailable(UnavailableConfig{
				Title:   listTitle + " unavailable",
				Message: "No list handler is configured for this resource.",
			})))
			return
		}
		rows, err := config.Handlers.List()
		if err != nil {
			renderStatus(w, r, http.StatusInternalServerError, Layout(listTitle, config.Nav, currentPath, Unavailable(UnavailableConfig{
				Title:   listTitle,
				Message: "Unable to load entries: " + err.Error(),
			})))
			return
		}

		tableCfg := TableConfig{
			Title:       listTitle,
			Columns:     config.Columns,
			Rows:        rows,
			EmptyText:   config.EmptyText,
			CreateLabel: config.CreateLabel,
		}
		if config.Handlers.Create != nil {
			tableCfg.CreateHref = base + "/new"
		}
		if config.Handlers.Update != nil {
			tableCfg.EditHref = func(id string) string { return base + "/" + id + "/edit" }
		}
		if config.Handlers.Delete != nil {
			tableCfg.DeleteHref = func(id string) string { return base + "/" + id + "/delete" }
		}

		renderStatus(w, r, http.StatusOK, Layout(listTitle, config.Nav, currentPath, Table(tableCfg)))
	})

	if config.Handlers.Create != nil {
		createTitle := "Add " + singular
		createAction := base + "/new"

		mux.HandleFunc("GET "+createAction, func(w http.ResponseWriter, r *http.Request) {
			formCfg := buildFormConfig(createTitle, createAction, base, config.Fields, nil, nil)
			renderStatus(w, r, http.StatusOK, Layout(createTitle, config.Nav, currentPath, ResourceForm(formCfg)))
		})

		mux.HandleFunc("POST "+createAction, func(w http.ResponseWriter, r *http.Request) {
			values, ok := parseResourceForm(r, config.Fields)
			if !ok {
				formCfg := buildFormConfig(createTitle, createAction, base, config.Fields, values, map[string]string{"form": "Invalid form submission"})
				renderStatus(w, r, http.StatusBadRequest, Layout(createTitle, config.Nav, currentPath, ResourceForm(formCfg)))
				return
			}
			res := config.Handlers.Create(values)
			if !res.Success {
				formCfg := buildFormConfig(createTitle, createAction, base, config.Fields, values, res.Errors)
				renderStatus(w, r, http.StatusOK, Layout(createTitle, config.Nav, currentPath, ResourceForm(formCfg)))
				return
			}
			http.Redirect(w, r, base, http.StatusSeeOther)
		})
	}

	if config.Handlers.Update != nil {
		editTitle := "Edit " + singular

		mux.HandleFunc("GET "+base+"/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			row, found := findResourceRow(config.Handlers, id)
			if !found {
				renderUnavailable(w, r, config, currentPath, listTitle)
				return
			}
			action := base + "/" + id + "/edit"
			formCfg := buildFormConfig(editTitle, action, base, config.Fields, row.Values, nil)
			renderStatus(w, r, http.StatusOK, Layout(editTitle, config.Nav, currentPath, ResourceForm(formCfg)))
		})

		mux.HandleFunc("POST "+base+"/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			action := base + "/" + id + "/edit"
			values, ok := parseResourceForm(r, config.Fields)
			if !ok {
				formCfg := buildFormConfig(editTitle, action, base, config.Fields, values, map[string]string{"form": "Invalid form submission"})
				renderStatus(w, r, http.StatusBadRequest, Layout(editTitle, config.Nav, currentPath, ResourceForm(formCfg)))
				return
			}
			res := config.Handlers.Update(id, values)
			if res.NotFound {
				renderUnavailable(w, r, config, currentPath, listTitle)
				return
			}
			if !res.Success {
				formCfg := buildFormConfig(editTitle, action, base, config.Fields, values, res.Errors)
				renderStatus(w, r, http.StatusOK, Layout(editTitle, config.Nav, currentPath, ResourceForm(formCfg)))
				return
			}
			http.Redirect(w, r, base, http.StatusSeeOther)
		})
	}

	if config.Handlers.Delete != nil {
		deleteTitle := "Delete " + singular

		mux.HandleFunc("GET "+base+"/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			row, found := findResourceRow(config.Handlers, id)
			if !found {
				renderUnavailable(w, r, config, currentPath, listTitle)
				return
			}
			label := resolveItemLabel(config, row)
			renderStatus(w, r, http.StatusOK, Layout(deleteTitle, config.Nav, currentPath, DeleteConfirmation(DeleteConfig{
				Title:      deleteTitle,
				ItemLabel:  label,
				Action:     base + "/" + id + "/delete",
				CancelHref: base,
			})))
		})

		mux.HandleFunc("POST "+base+"/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			if !config.Handlers.Delete(id) {
				renderUnavailable(w, r, config, currentPath, listTitle)
				return
			}
			http.Redirect(w, r, base, http.StatusSeeOther)
		})
	}
}

func findResourceRow(handlers ResourceHandlers, id string) (TableRow, bool) {
	if handlers.Get != nil {
		return handlers.Get(id)
	}
	if handlers.List == nil {
		return TableRow{}, false
	}
	rows, err := handlers.List()
	if err != nil {
		return TableRow{}, false
	}
	for _, row := range rows {
		if row.ID == id {
			return row, true
		}
	}
	return TableRow{}, false
}

func resolveItemLabel(config ResourceConfig, row TableRow) string {
	if config.ItemLabel != nil {
		if label := config.ItemLabel(row); label != "" {
			return label
		}
	}
	for _, col := range config.Columns {
		if val := row.Values[col.Key]; val != "" {
			return val
		}
	}
	return row.ID
}

func buildFormConfig(title, action, cancelHref string, templateFields []FormField, values map[string]string, errors map[string]string) FormConfig {
	fields := make([]FormField, len(templateFields))
	for i, f := range templateFields {
		if values != nil {
			if v, ok := values[f.Name]; ok {
				f.Value = v
			}
		}
		fields[i] = f
	}
	return FormConfig{
		Title:       title,
		Action:      action,
		Method:      "post",
		SubmitLabel: "Save",
		CancelHref:  cancelHref,
		Fields:      fields,
		Errors:      errors,
	}
}

func parseResourceForm(r *http.Request, fields []FormField) (map[string]string, bool) {
	if err := r.ParseForm(); err != nil {
		return map[string]string{}, false
	}
	values := make(map[string]string, len(fields))
	for _, f := range fields {
		values[f.Name] = strings.TrimSpace(r.Form.Get(f.Name))
	}
	return values, true
}

func renderUnavailable(w http.ResponseWriter, r *http.Request, config ResourceConfig, currentPath, listTitle string) {
	base := "/" + strings.Trim(config.BasePath, "/")
	if base == "/" {
		base = "/resource"
	}
	renderStatus(w, r, http.StatusNotFound, Layout("Entry Unavailable", config.Nav, currentPath, Unavailable(UnavailableConfig{
		Title:     "Entry Unavailable",
		Message:   "The requested entry could not be found or was deleted.",
		BackHref:  base,
		BackLabel: "Back to " + listTitle,
	})))
}

func renderStatus(w http.ResponseWriter, r *http.Request, status int, component templ.Component) {
	var buf bytes.Buffer
	if err := component.Render(r.Context(), &buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}
