package builder

import (
	"strings"
	"testing"

	"github.com/validkeys/agent-instruction/internal/files"
)

func TestWrapWithMarkers(t *testing.T) {
	tests := map[string]struct {
		content string
		want    string
	}{
		"empty content": {
			content: "",
			want:    "<!-- BEGIN AGENT-INSTRUCTION -->\n\n<!-- END AGENT-INSTRUCTION -->",
		},
		"simple content": {
			content: "# Rules\n\nAlways use error wrapping",
			want:    "<!-- BEGIN AGENT-INSTRUCTION -->\n# Rules\n\nAlways use error wrapping\n<!-- END AGENT-INSTRUCTION -->",
		},
		"multiline content": {
			content: "## Error Handling\n\nAlways use error wrapping\n\n## Testing\n\nWrite table-driven tests",
			want:    "<!-- BEGIN AGENT-INSTRUCTION -->\n## Error Handling\n\nAlways use error wrapping\n\n## Testing\n\nWrite table-driven tests\n<!-- END AGENT-INSTRUCTION -->",
		},
		"content with single trailing newline": {
			content: "Rule text\n",
			want:    "<!-- BEGIN AGENT-INSTRUCTION -->\nRule text\n\n<!-- END AGENT-INSTRUCTION -->",
		},
		"content with double trailing newline": {
			content: "Rule text\n\n",
			want:    "<!-- BEGIN AGENT-INSTRUCTION -->\nRule text\n\n<!-- END AGENT-INSTRUCTION -->",
		},
		"content without trailing newline": {
			content: "Rule text",
			want:    "<!-- BEGIN AGENT-INSTRUCTION -->\nRule text\n<!-- END AGENT-INSTRUCTION -->",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := WrapWithMarkers(tc.content)
			if got != tc.want {
				t.Errorf("WrapWithMarkers() mismatch\nGot:\n%q\n\nWant:\n%q", got, tc.want)
			}
		})
	}
}

func TestBuildManagedFile(t *testing.T) {
	tests := map[string]struct {
		generated string
		existing  *files.ManagedContent
		want      string
	}{
		"new file with no existing content": {
			generated: "## Rules\n\nAlways use error wrapping\n\n",
			existing:  nil,
			want:      "<!-- BEGIN AGENT-INSTRUCTION -->\n## Rules\n\nAlways use error wrapping\n\n<!-- END AGENT-INSTRUCTION -->",
		},
		"existing file with managed section only": {
			generated: "## Rules\n\nNew rule content\n\n",
			existing: &files.ManagedContent{
				Before:  "",
				Managed: "\n## Rules\n\nOld rule content\n\n",
				After:   "",
			},
			want: "<!-- BEGIN AGENT-INSTRUCTION -->\n## Rules\n\nNew rule content\n\n<!-- END AGENT-INSTRUCTION -->",
		},
		"existing file with content before managed section": {
			generated: "## Rules\n\nGenerated rules\n\n",
			existing: &files.ManagedContent{
				Before:  "# My Project\n\nCustom introduction text\n\n",
				Managed: "\n## Rules\n\nOld rules\n\n",
				After:   "",
			},
			want: "# My Project\n\nCustom introduction text\n\n<!-- BEGIN AGENT-INSTRUCTION -->\n## Rules\n\nGenerated rules\n\n<!-- END AGENT-INSTRUCTION -->",
		},
		"existing file with content after managed section": {
			generated: "## Generated Rules\n\n",
			existing: &files.ManagedContent{
				Before:  "",
				Managed: "\n## Old Rules\n\n",
				After:   "\n\n# Custom Footer\n\nUser notes here\n",
			},
			want: "<!-- BEGIN AGENT-INSTRUCTION -->\n## Generated Rules\n\n<!-- END AGENT-INSTRUCTION -->\n\n# Custom Footer\n\nUser notes here\n",
		},
		"existing file with content before and after": {
			generated: "## Rules\n\nGenerated content\n\n",
			existing: &files.ManagedContent{
				Before:  "# Header\n\nUser intro\n\n",
				Managed: "\n## Rules\n\nOld content\n\n",
				After:   "\n\n# Footer\n\nUser notes\n",
			},
			want: "# Header\n\nUser intro\n\n<!-- BEGIN AGENT-INSTRUCTION -->\n## Rules\n\nGenerated content\n\n<!-- END AGENT-INSTRUCTION -->\n\n# Footer\n\nUser notes\n",
		},
		"empty generated content with existing user content": {
			generated: "",
			existing: &files.ManagedContent{
				Before:  "# User Header\n\n",
				Managed: "\nOld managed content\n",
				After:   "\n# User Footer\n",
			},
			want: "# User Header\n\n<!-- BEGIN AGENT-INSTRUCTION -->\n\n<!-- END AGENT-INSTRUCTION -->\n# User Footer\n",
		},
		"preserves whitespace in user content": {
			generated: "Generated\n",
			existing: &files.ManagedContent{
				Before:  "Before   with   spaces\n\n",
				Managed: "\nOld\n",
				After:   "\n\nAfter   with   spaces",
			},
			want: "Before   with   spaces\n\n<!-- BEGIN AGENT-INSTRUCTION -->\nGenerated\n\n<!-- END AGENT-INSTRUCTION -->\n\nAfter   with   spaces",
		},
		"empty before and after sections": {
			generated: "New content\n",
			existing: &files.ManagedContent{
				Before:  "",
				Managed: "\nOld content\n",
				After:   "",
			},
			want: "<!-- BEGIN AGENT-INSTRUCTION -->\nNew content\n\n<!-- END AGENT-INSTRUCTION -->",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := BuildManagedFile(tc.generated, tc.existing)
			if got != tc.want {
				t.Errorf("BuildManagedFile() mismatch\nGot:\n%q\n\nWant:\n%q", got, tc.want)
			}
		})
	}
}

func TestBuildManagedFile_NoContentLoss(t *testing.T) {
	// Verify that user content is never lost when building managed files
	userBefore := "# My Custom Header\n\nThis is important user content.\n\n"
	userAfter := "\n\n# My Custom Footer\n\nMore user notes here.\n"
	generated := "## Generated Rules\n\nAlways test your code\n\n"

	existing := &files.ManagedContent{
		Before:  userBefore,
		Managed: "\n## Old Rules\n\nOld rule content\n\n",
		After:   userAfter,
	}

	result := BuildManagedFile(generated, existing)

	// Verify all user content is present
	if !strings.Contains(result, userBefore) {
		t.Error("user content before managed section was lost")
	}
	if !strings.Contains(result, userAfter) {
		t.Error("user content after managed section was lost")
	}
	if !strings.Contains(result, generated) {
		t.Error("generated content was not included")
	}

	// Verify old managed content is replaced
	if strings.Contains(result, "Old rule content") {
		t.Error("old managed content was not replaced")
	}

	// Verify markers are present
	if !strings.Contains(result, files.BeginMarker) {
		t.Error("begin marker missing")
	}
	if !strings.Contains(result, files.EndMarker) {
		t.Error("end marker missing")
	}
}

func TestBuildManagedFile_Integration(t *testing.T) {
	// Test realistic scenario: updating CLAUDE.md with new generated rules
	// while preserving user's custom introduction and notes

	userIntro := `# Project Documentation

This is my project. Please follow these guidelines.

`

	userNotes := `

## Additional Notes

These are my personal notes that should never be deleted.

- Note 1
- Note 2
`

	oldManaged := `
## Error Handling

Old error handling rules

`

	newGenerated := `## Error Handling

Always use fmt.Errorf with %w for error wrapping

## Testing

Write comprehensive table-driven tests

`

	existing := &files.ManagedContent{
		Before:  userIntro,
		Managed: oldManaged,
		After:   userNotes,
	}

	result := BuildManagedFile(newGenerated, existing)

	// Verify structure
	lines := strings.Split(result, "\n")
	if len(lines) < 10 {
		t.Errorf("expected at least 10 lines in result, got %d", len(lines))
	}

	// Verify content order
	introIdx := strings.Index(result, "This is my project")
	beginIdx := strings.Index(result, files.BeginMarker)
	generatedIdx := strings.Index(result, "Always use fmt.Errorf")
	endIdx := strings.Index(result, files.EndMarker)
	notesIdx := strings.Index(result, "Additional Notes")

	if introIdx == -1 || beginIdx == -1 || generatedIdx == -1 || endIdx == -1 || notesIdx == -1 {
		t.Fatal("missing expected content sections")
	}

	// Verify correct order: intro < begin < generated < end < notes
	if !(introIdx < beginIdx && beginIdx < generatedIdx && generatedIdx < endIdx && endIdx < notesIdx) {
		t.Error("content sections are not in correct order")
	}

	// Verify old content is gone
	if strings.Contains(result, "Old error handling rules") {
		t.Error("old managed content was not replaced")
	}
}
