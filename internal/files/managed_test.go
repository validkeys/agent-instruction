package files

import (
	"testing"
)

func TestParseManagedContent(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    *ManagedContent
		wantErr bool
	}{
		"content with managed section": {
			input: `# My Instructions

Some user content

<!-- BEGIN AGENT-INSTRUCTION -->
Generated content here
<!-- END AGENT-INSTRUCTION -->

More user content`,
			want: &ManagedContent{
				Before:  "# My Instructions\n\nSome user content\n\n",
				Managed: "\nGenerated content here\n",
				After:   "\n\nMore user content",
			},
			wantErr: false,
		},
		"content without markers": {
			input: `# My Instructions

Just user content`,
			want: &ManagedContent{
				Before:  "# My Instructions\n\nJust user content",
				Managed: "",
				After:   "",
			},
			wantErr: false,
		},
		"only begin marker": {
			input: `# My Instructions

<!-- BEGIN AGENT-INSTRUCTION -->
Content without end`,
			want:    nil,
			wantErr: true,
		},
		"only end marker": {
			input: `# My Instructions

Content without begin
<!-- END AGENT-INSTRUCTION -->`,
			want:    nil,
			wantErr: true,
		},
		"end marker before begin marker": {
			input: `# My Instructions

<!-- END AGENT-INSTRUCTION -->
Some content
<!-- BEGIN AGENT-INSTRUCTION -->`,
			want:    nil,
			wantErr: true,
		},
		"empty content": {
			input: "",
			want: &ManagedContent{
				Before:  "",
				Managed: "",
				After:   "",
			},
			wantErr: false,
		},
		"markers with empty managed section": {
			input: `<!-- BEGIN AGENT-INSTRUCTION -->
<!-- END AGENT-INSTRUCTION -->`,
			want: &ManagedContent{
				Before:  "",
				Managed: "\n",
				After:   "",
			},
			wantErr: false,
		},
		"preserves whitespace": {
			input: `  Some content

<!-- BEGIN AGENT-INSTRUCTION -->
  Indented content
<!-- END AGENT-INSTRUCTION -->

  More content  `,
			want: &ManagedContent{
				Before:  "  Some content\n\n",
				Managed: "\n  Indented content\n",
				After:   "\n\n  More content  ",
			},
			wantErr: false,
		},
		"multiple lines in managed section": {
			input: `User content

<!-- BEGIN AGENT-INSTRUCTION -->
Line 1
Line 2
Line 3
<!-- END AGENT-INSTRUCTION -->

After content`,
			want: &ManagedContent{
				Before:  "User content\n\n",
				Managed: "\nLine 1\nLine 2\nLine 3\n",
				After:   "\n\nAfter content",
			},
			wantErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseManagedContent(tc.input)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tc.wantErr {
				if got.Before != tc.want.Before {
					t.Errorf("Before mismatch:\ngot:  %q\nwant: %q", got.Before, tc.want.Before)
				}
				if got.Managed != tc.want.Managed {
					t.Errorf("Managed mismatch:\ngot:  %q\nwant: %q", got.Managed, tc.want.Managed)
				}
				if got.After != tc.want.After {
					t.Errorf("After mismatch:\ngot:  %q\nwant: %q", got.After, tc.want.After)
				}
			}
		})
	}
}

func TestHasManagedSection(t *testing.T) {
	tests := map[string]struct {
		input string
		want  bool
	}{
		"has both markers": {
			input: `<!-- BEGIN AGENT-INSTRUCTION -->
content
<!-- END AGENT-INSTRUCTION -->`,
			want: true,
		},
		"has only begin marker": {
			input: `<!-- BEGIN AGENT-INSTRUCTION -->
content`,
			want: false,
		},
		"has only end marker": {
			input: `content
<!-- END AGENT-INSTRUCTION -->`,
			want: false,
		},
		"has no markers": {
			input: "just plain content",
			want:  false,
		},
		"empty string": {
			input: "",
			want:  false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := HasManagedSection(tc.input)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestReplaceManagedSection(t *testing.T) {
	tests := map[string]struct {
		content    string
		newSection string
		want       string
		wantErr    bool
	}{
		"replaces existing managed section": {
			content: `User content

<!-- BEGIN AGENT-INSTRUCTION -->
Old generated content
<!-- END AGENT-INSTRUCTION -->

More user content`,
			newSection: "New generated content",
			want: `User content

<!-- BEGIN AGENT-INSTRUCTION -->
New generated content
<!-- END AGENT-INSTRUCTION -->

More user content`,
			wantErr: false,
		},
		"adds managed section when none exists": {
			content:    "Just user content",
			newSection: "Generated content",
			want: `Just user content

<!-- BEGIN AGENT-INSTRUCTION -->
Generated content
<!-- END AGENT-INSTRUCTION -->`,
			wantErr: false,
		},
		"only begin marker present": {
			content:    "<!-- BEGIN AGENT-INSTRUCTION -->\nContent",
			newSection: "New content",
			want:       "",
			wantErr:    true,
		},
		"only end marker present": {
			content:    "Content\n<!-- END AGENT-INSTRUCTION -->",
			newSection: "New content",
			want:       "",
			wantErr:    true,
		},
		"end before begin": {
			content:    "<!-- END AGENT-INSTRUCTION -->\n<!-- BEGIN AGENT-INSTRUCTION -->",
			newSection: "New content",
			want:       "",
			wantErr:    true,
		},
		"empty new section": {
			content: `User content

<!-- BEGIN AGENT-INSTRUCTION -->
Old content
<!-- END AGENT-INSTRUCTION -->

More content`,
			newSection: "",
			want: `User content

<!-- BEGIN AGENT-INSTRUCTION -->

<!-- END AGENT-INSTRUCTION -->

More content`,
			wantErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ReplaceManagedSection(tc.content, tc.newSection)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tc.wantErr && got != tc.want {
				t.Errorf("mismatch:\ngot:\n%s\n\nwant:\n%s", got, tc.want)
			}
		})
	}
}
