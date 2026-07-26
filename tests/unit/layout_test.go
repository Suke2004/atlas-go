package unit

import (
	"bytes"
	"context"
	"strings"
	"testing"

	layouts "github.com/Suke2004/atlas-go/web/templates/layout"
)

func TestLayout_SidebarRendering(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer

	err := layouts.Sidebar("/", "Alex Mercer").Render(ctx, &buf)
	if err != nil {
		t.Fatalf("failed to render Sidebar: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "Dashboard") {
		t.Errorf("expected Sidebar to contain 'Dashboard', got: %s", html)
	}
	if !strings.Contains(html, "Alex Mercer") {
		t.Errorf("expected Sidebar to contain username 'Alex Mercer', got: %s", html)
	}
}

func TestLayout_TopbarRendering(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer

	err := layouts.Topbar("Projects", "Alex Mercer").Render(ctx, &buf)
	if err != nil {
		t.Fatalf("failed to render Topbar: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "Projects") {
		t.Errorf("expected Topbar to contain title 'Projects', got: %s", html)
	}
	if !strings.Contains(html, "Search Atlas...") {
		t.Errorf("expected Topbar to contain global search trigger, got: %s", html)
	}
}

func TestLayout_ToastRendering(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer

	err := layouts.Toast("Project created successfully", layouts.ToastSuccess).Render(ctx, &buf)
	if err != nil {
		t.Fatalf("failed to render Toast: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "Project created successfully") {
		t.Errorf("expected Toast to contain message, got: %s", html)
	}
	if !strings.Contains(html, "emerald") {
		t.Errorf("expected ToastSuccess to contain emerald style, got: %s", html)
	}
}
