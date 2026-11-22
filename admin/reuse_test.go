package admin_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/rosso0815/rosso0815-dashboard/admin"
)

func TestMinimalHostCanRenderLibrary(t *testing.T) {
	component := admin.Layout("Minimal Host", []admin.NavItem{{Label: "Users", Href: "/users"}}, "/users",
		admin.Table(admin.TableConfig{
			Columns: []admin.TableColumn{{Key: "name", Label: "Name"}},
			Rows: []admin.TableRow{{
				ID:     "1",
				Values: map[string]string{"name": "Host User"},
			}},
		}),
	)

	var out bytes.Buffer
	if err := component.Render(context.Background(), &out); err != nil {
		t.Fatalf("minimal host render failed: %v", err)
	}
	if out.Len() == 0 {
		t.Fatal("minimal host rendered empty output")
	}
}
