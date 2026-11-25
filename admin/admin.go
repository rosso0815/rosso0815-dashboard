package admin

// NavItem defines a single navigation link in the dashboard sidebar.
type NavItem struct {
	Label string
	Href  string
}

// TableColumn maps a row value key to a table column header.
type TableColumn struct {
	Key   string
	Label string
}

// TableRow represents one resource record rendered in a Table.
type TableRow struct {
	ID     string
	Values map[string]string
}

// TableConfig configures a reusable resource table and optional header actions.
type TableConfig struct {
	Title       string
	Columns     []TableColumn
	Rows        []TableRow
	EmptyText   string
	CreateHref  string
	CreateLabel string
	EditHref    func(string) string
	DeleteHref  func(string) string
}

// FormField defines a single input rendered by ResourceForm.
type FormField struct {
	Name        string
	Label       string
	Value       string
	Type        string
	Placeholder string
	Required    bool
}

// FormConfig configures a reusable create or edit form.
type FormConfig struct {
	Title       string
	Action      string
	Method      string
	SubmitLabel string
	CancelHref  string
	Fields      []FormField
	Errors      map[string]string
}

// DeleteConfig configures a delete confirmation page.
type DeleteConfig struct {
	Title      string
	ItemLabel  string
	Action     string
	CancelHref string
}

// UnavailableConfig configures the missing-resource state rendered inside Layout.
type UnavailableConfig struct {
	Title     string
	Message   string
	BackHref  string
	BackLabel string
}

// SubmissionResult reports the outcome of a host create or update handler.
type SubmissionResult struct {
	Success  bool
	Errors   map[string]string
	NotFound bool
}

// ResourceHandlers defines host-owned persistence and validation callbacks.
type ResourceHandlers struct {
	List   func() ([]TableRow, error)
	Get    func(string) (TableRow, bool)
	Create func(map[string]string) SubmissionResult
	Update func(string, map[string]string) SubmissionResult
	Delete func(string) bool
}
