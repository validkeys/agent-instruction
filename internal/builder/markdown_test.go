package builder

import (
	"strings"
	"testing"

	"github.com/validkeys/agent-instruction/internal/rules"
)

func TestInstructionsToMarkdown(t *testing.T) {
	tests := map[string]struct {
		instructions []rules.Instruction
		want         string
	}{
		"empty instructions": {
			instructions: []rules.Instruction{},
			want:         "",
		},
		"simple instruction without heading": {
			instructions: []rules.Instruction{
				{Rule: "Always use error wrapping"},
			},
			want: "Always use error wrapping\n\n",
		},
		"instruction with heading": {
			instructions: []rules.Instruction{
				{
					Heading: "Error Handling",
					Rule:    "Always use error wrapping",
				},
			},
			want: "## Error Handling\n\nAlways use error wrapping\n\n",
		},
		"instruction with single reference": {
			instructions: []rules.Instruction{
				{
					Rule: "Follow the error handling guide",
					References: []rules.Reference{
						{Title: "Error Guide", Path: "/docs/errors.md"},
					},
				},
			},
			want: "Follow the error handling guide\n\n- [Error Guide](/docs/errors.md)\n\n",
		},
		"instruction with multiple references": {
			instructions: []rules.Instruction{
				{
					Rule: "Follow the style guides",
					References: []rules.Reference{
						{Title: "Go Style", Path: "/docs/go.md"},
						{Title: "Testing Guide", Path: "/docs/testing.md"},
					},
				},
			},
			want: "Follow the style guides\n\n- [Go Style](/docs/go.md)\n- [Testing Guide](/docs/testing.md)\n\n",
		},
		"instruction with heading and references": {
			instructions: []rules.Instruction{
				{
					Heading: "Testing",
					Rule:    "Write comprehensive tests",
					References: []rules.Reference{
						{Title: "Test Patterns", Path: "/docs/test-patterns.md"},
					},
				},
			},
			want: "## Testing\n\nWrite comprehensive tests\n\n- [Test Patterns](/docs/test-patterns.md)\n\n",
		},
		"multiple instructions mixed": {
			instructions: []rules.Instruction{
				{
					Heading: "Error Handling",
					Rule:    "Always use error wrapping",
				},
				{
					Rule: "Validate input at boundaries",
				},
				{
					Heading: "Testing",
					Rule:    "Write table-driven tests",
					References: []rules.Reference{
						{Title: "Test Guide", Path: "/docs/testing.md"},
					},
				},
			},
			want: "## Error Handling\n\nAlways use error wrapping\n\nValidate input at boundaries\n\n## Testing\n\nWrite table-driven tests\n\n- [Test Guide](/docs/testing.md)\n\n",
		},
		"multiline rule text": {
			instructions: []rules.Instruction{
				{
					Rule: "Always:\n- Use error wrapping\n- Add context\n- Return early",
				},
			},
			want: "Always:\n- Use error wrapping\n- Add context\n- Return early\n\n",
		},
		"rule with markdown formatting": {
			instructions: []rules.Instruction{
				{
					Rule: "Use `fmt.Errorf` with `%w` for wrapping",
				},
			},
			want: "Use `fmt.Errorf` with `%w` for wrapping\n\n",
		},
		"instruction with empty references array": {
			instructions: []rules.Instruction{
				{
					Rule:       "Simple rule",
					References: []rules.Reference{},
				},
			},
			want: "Simple rule\n\n",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := InstructionsToMarkdown(tc.instructions)
			if got != tc.want {
				t.Errorf("InstructionsToMarkdown() mismatch\nGot:\n%q\n\nWant:\n%q", got, tc.want)
			}
		})
	}
}

func TestFormatReference(t *testing.T) {
	tests := map[string]struct {
		ref  rules.Reference
		want string
	}{
		"simple reference": {
			ref:  rules.Reference{Title: "Guide", Path: "/docs/guide.md"},
			want: "- [Guide](/docs/guide.md)",
		},
		"reference with spaces in title": {
			ref:  rules.Reference{Title: "Error Handling Guide", Path: "/docs/errors.md"},
			want: "- [Error Handling Guide](/docs/errors.md)",
		},
		"reference with relative path": {
			ref:  rules.Reference{Title: "Local File", Path: "./local.md"},
			want: "- [Local File](./local.md)",
		},
		"reference with absolute path": {
			ref:  rules.Reference{Title: "Absolute", Path: "/absolute/path.md"},
			want: "- [Absolute](/absolute/path.md)",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := formatReference(tc.ref)
			if got != tc.want {
				t.Errorf("formatReference() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestInstructionsToMarkdown_Integration(t *testing.T) {
	// Test a realistic scenario with multiple instructions
	instructions := []rules.Instruction{
		{
			Heading: "Core Principles",
			Rule:    "Write simple, maintainable code that solves the problem at hand.",
		},
		{
			Rule: "Avoid over-engineering and premature optimization.",
		},
		{
			Heading: "Error Handling",
			Rule:    "Always handle errors explicitly. Use `fmt.Errorf` with `%w` to wrap errors with context.",
			References: []rules.Reference{
				{Title: "Error Handling Guide", Path: "/docs/style-anchors/error-handling.md"},
			},
		},
		{
			Heading: "Testing",
			Rule:    "Write table-driven tests using `map[string]struct` for test cases.",
			References: []rules.Reference{
				{Title: "Table-Driven Testing", Path: "/docs/style-anchors/table-driven-testing.md"},
				{Title: "Test Examples", Path: "/docs/examples/tests.md"},
			},
		},
	}

	got := InstructionsToMarkdown(instructions)

	// Verify structure
	if !strings.Contains(got, "## Core Principles") {
		t.Error("missing Core Principles heading")
	}
	if !strings.Contains(got, "## Error Handling") {
		t.Error("missing Error Handling heading")
	}
	if !strings.Contains(got, "## Testing") {
		t.Error("missing Testing heading")
	}

	// Verify rules
	if !strings.Contains(got, "Write simple, maintainable code") {
		t.Error("missing first rule text")
	}
	if !strings.Contains(got, "Avoid over-engineering") {
		t.Error("missing second rule text")
	}

	// Verify references
	if !strings.Contains(got, "[Error Handling Guide](/docs/style-anchors/error-handling.md)") {
		t.Error("missing error handling reference")
	}
	if !strings.Contains(got, "[Table-Driven Testing](/docs/style-anchors/table-driven-testing.md)") {
		t.Error("missing testing reference")
	}
	if !strings.Contains(got, "[Test Examples](/docs/examples/tests.md)") {
		t.Error("missing test examples reference")
	}

	// Verify proper spacing
	lines := strings.Split(got, "\n")
	if len(lines) < 10 {
		t.Errorf("expected at least 10 lines, got %d", len(lines))
	}
}
