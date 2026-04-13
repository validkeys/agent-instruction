package rules

import (
	"strings"
	"testing"
)

func TestDetectCycle(t *testing.T) {
	tests := map[string]struct {
		path      string
		pathStack []string
		wantErr   bool
		wantMsg   string
	}{
		"no cycle - path not in stack": {
			path:      "/rules/new.json",
			pathStack: []string{"/rules/a.json", "/rules/b.json"},
			wantErr:   false,
		},
		"direct cycle - same file": {
			path:      "/rules/a.json",
			pathStack: []string{"/rules/a.json"},
			wantErr:   true,
			wantMsg:   "/rules/a.json → /rules/a.json",
		},
		"indirect cycle - A→B→C→A": {
			path:      "/rules/a.json",
			pathStack: []string{"/rules/a.json", "/rules/b.json", "/rules/c.json"},
			wantErr:   true,
			wantMsg:   "/rules/a.json → /rules/b.json → /rules/c.json → /rules/a.json",
		},
		"cycle in middle of chain": {
			path:      "/rules/b.json",
			pathStack: []string{"/rules/a.json", "/rules/b.json", "/rules/c.json"},
			wantErr:   true,
			wantMsg:   "/rules/a.json → /rules/b.json → /rules/c.json → /rules/b.json",
		},
		"empty path stack - no cycle": {
			path:      "/rules/a.json",
			pathStack: []string{},
			wantErr:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := &ImportContext{
				visited:   make(map[string]bool),
				pathStack: tc.pathStack,
			}

			err := detectCycle(tc.path, ctx)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if tc.wantMsg != "" && !strings.Contains(err.Error(), tc.wantMsg) {
					t.Errorf("error message:\ngot:  %q\nwant: %q", err.Error(), tc.wantMsg)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestBuildCycleError(t *testing.T) {
	tests := map[string]struct {
		stack   []string
		newPath string
		want    string
	}{
		"simple self-reference": {
			stack:   []string{"/rules/a.json"},
			newPath: "/rules/a.json",
			want:    "import cycle: /rules/a.json → /rules/a.json",
		},
		"two-node cycle": {
			stack:   []string{"/rules/a.json", "/rules/b.json"},
			newPath: "/rules/a.json",
			want:    "import cycle: /rules/a.json → /rules/b.json → /rules/a.json",
		},
		"three-node cycle": {
			stack:   []string{"/rules/a.json", "/rules/b.json", "/rules/c.json"},
			newPath: "/rules/a.json",
			want:    "import cycle: /rules/a.json → /rules/b.json → /rules/c.json → /rules/a.json",
		},
		"cycle back to middle": {
			stack:   []string{"/rules/a.json", "/rules/b.json", "/rules/c.json"},
			newPath: "/rules/b.json",
			want:    "import cycle: /rules/a.json → /rules/b.json → /rules/c.json → /rules/b.json",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := buildCycleError(tc.stack, tc.newPath)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			got := err.Error()
			if got != tc.want {
				t.Errorf("error message:\ngot:  %q\nwant: %q", got, tc.want)
			}
		})
	}
}

func TestImportContext(t *testing.T) {
	t.Run("can initialize empty context", func(t *testing.T) {
		ctx := NewImportContext()

		if ctx.visited == nil {
			t.Error("visited map should be initialized")
		}

		if ctx.pathStack == nil {
			t.Error("pathStack should be initialized")
		}

		if len(ctx.visited) != 0 {
			t.Errorf("visited should be empty, got %d entries", len(ctx.visited))
		}

		if len(ctx.pathStack) != 0 {
			t.Errorf("pathStack should be empty, got %d entries", len(ctx.pathStack))
		}
	})

	t.Run("can push and pop from path stack", func(t *testing.T) {
		ctx := NewImportContext()

		ctx.pathStack = append(ctx.pathStack, "/rules/a.json")
		if len(ctx.pathStack) != 1 {
			t.Errorf("expected 1 item in stack, got %d", len(ctx.pathStack))
		}

		ctx.pathStack = append(ctx.pathStack, "/rules/b.json")
		if len(ctx.pathStack) != 2 {
			t.Errorf("expected 2 items in stack, got %d", len(ctx.pathStack))
		}

		// Pop
		ctx.pathStack = ctx.pathStack[:len(ctx.pathStack)-1]
		if len(ctx.pathStack) != 1 {
			t.Errorf("expected 1 item after pop, got %d", len(ctx.pathStack))
		}

		if ctx.pathStack[0] != "/rules/a.json" {
			t.Errorf("expected /rules/a.json, got %s", ctx.pathStack[0])
		}
	})
}
