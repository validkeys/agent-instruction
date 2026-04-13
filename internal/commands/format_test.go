package commands

import (
	"strings"
	"testing"
)

func TestFormatSuccess(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want string
	}{
		{
			name: "formats success message",
			msg:  "Operation completed",
			want: "✓ Operation completed",
		},
		{
			name: "handles empty message",
			msg:  "",
			want: "✓ ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSuccess(tt.msg)
			if got != tt.want {
				t.Errorf("FormatSuccess(%q) = %q, want %q", tt.msg, got, tt.want)
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want string
	}{
		{
			name: "formats error message",
			msg:  "Operation failed",
			want: "✗ Operation failed",
		},
		{
			name: "handles empty message",
			msg:  "",
			want: "✗ ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatError(tt.msg)
			if got != tt.want {
				t.Errorf("FormatError(%q) = %q, want %q", tt.msg, got, tt.want)
			}
		})
	}
}

func TestFormatHeading(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "formats heading",
			text: "Configuration",
			want: "\n=== Configuration ===\n",
		},
		{
			name: "handles empty text",
			text: "",
			want: "\n===  ===\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatHeading(tt.text)
			if got != tt.want {
				t.Errorf("FormatHeading(%q) = %q, want %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestFormatList(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  []string // Check for these substrings
	}{
		{
			name:  "formats list with items",
			items: []string{"Item 1", "Item 2", "Item 3"},
			want:  []string{"  - Item 1", "  - Item 2", "  - Item 3"},
		},
		{
			name:  "handles empty list",
			items: []string{},
			want:  []string{},
		},
		{
			name:  "handles single item",
			items: []string{"Only item"},
			want:  []string{"  - Only item"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatList(tt.items)

			if len(tt.items) == 0 {
				if got != "" {
					t.Errorf("FormatList([]) = %q, want empty string", got)
				}
				return
			}

			for _, wantItem := range tt.want {
				if !strings.Contains(got, wantItem) {
					t.Errorf("FormatList() missing %q in output:\n%s", wantItem, got)
				}
			}
		})
	}
}

func TestIndentText(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		spaces int
		want   string
	}{
		{
			name:   "indents single line with 2 spaces",
			text:   "Hello",
			spaces: 2,
			want:   "  Hello",
		},
		{
			name:   "indents multiple lines",
			text:   "Line 1\nLine 2\nLine 3",
			spaces: 4,
			want:   "    Line 1\n    Line 2\n    Line 3",
		},
		{
			name:   "handles zero spaces",
			text:   "No indent",
			spaces: 0,
			want:   "No indent",
		},
		{
			name:   "handles empty text",
			text:   "",
			spaces: 2,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IndentText(tt.text, tt.spaces)
			if got != tt.want {
				t.Errorf("IndentText(%q, %d) = %q, want %q", tt.text, tt.spaces, got, tt.want)
			}
		})
	}
}

func TestIsTerminal(t *testing.T) {
	// Note: This is difficult to test in a unit test since it depends on
	// the actual file descriptor. We just verify it doesn't crash.
	t.Run("doesn't crash", func(t *testing.T) {
		_ = isTerminal(1) // stdout
	})
}
