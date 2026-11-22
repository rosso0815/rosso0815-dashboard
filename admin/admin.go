package admin

type NavItem struct {
	Label string
	Href  string
}

type TableColumn struct {
	Key   string
	Label string
}

type TableRow struct {
	ID     string
	Values map[string]string
}

type TableConfig struct {
	Columns    []TableColumn
	Rows       []TableRow
	EmptyText  string
	CreateHref string
	EditHref   func(string) string
	DeleteHref func(string) string
}

type FormField struct {
	Name        string
	Label       string
	Value       string
	Type        string
	Placeholder string
	Required    bool
}

type FormConfig struct {
	Title       string
	Action      string
	Method      string
	SubmitLabel string
	CancelHref  string
	Fields      []FormField
	Errors      map[string]string
}

type SubmissionResult struct {
	Success  bool
	Errors   map[string]string
	NotFound bool
}

type ResourceHandlers struct {
	List   func() ([]TableRow, error)
	Create func(map[string]string) SubmissionResult
	Update func(string, map[string]string) SubmissionResult
	Delete func(string) bool
}
