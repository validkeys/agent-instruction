package rules

import (
	"testing"
)

func TestMergeInstructions(t *testing.T) {
	tests := map[string]struct {
		sources [][]Instruction
		want    []Instruction
	}{
		"single source with multiple instructions": {
			sources: [][]Instruction{
				{
					{Heading: "Section 1", Rule: "Rule 1"},
					{Rule: "Rule 2"},
				},
			},
			want: []Instruction{
				{Heading: "Section 1", Rule: "Rule 1"},
				{Rule: "Rule 2"},
			},
		},
		"multiple sources merged in order": {
			sources: [][]Instruction{
				{
					{Heading: "Import A", Rule: "Rule from A"},
				},
				{
					{Heading: "Import B", Rule: "Rule from B"},
				},
				{
					{Heading: "Local", Rule: "Local rule"},
				},
			},
			want: []Instruction{
				{Heading: "Import A", Rule: "Rule from A"},
				{Heading: "Import B", Rule: "Rule from B"},
				{Heading: "Local", Rule: "Local rule"},
			},
		},
		"empty sources array": {
			sources: [][]Instruction{},
			want:    []Instruction{},
		},
		"sources with empty arrays": {
			sources: [][]Instruction{
				{},
				{
					{Rule: "Rule 1"},
				},
				{},
				{
					{Rule: "Rule 2"},
				},
			},
			want: []Instruction{
				{Rule: "Rule 1"},
				{Rule: "Rule 2"},
			},
		},
		"all empty sources": {
			sources: [][]Instruction{
				{},
				{},
				{},
			},
			want: []Instruction{},
		},
		"preserves heading hierarchy": {
			sources: [][]Instruction{
				{
					{Heading: "## Main", Rule: "Rule 1"},
					{Heading: "### Subsection", Rule: "Rule 2"},
				},
				{
					{Heading: "## Another", Rule: "Rule 3"},
				},
			},
			want: []Instruction{
				{Heading: "## Main", Rule: "Rule 1"},
				{Heading: "### Subsection", Rule: "Rule 2"},
				{Heading: "## Another", Rule: "Rule 3"},
			},
		},
		"preserves references": {
			sources: [][]Instruction{
				{
					{
						Rule: "Rule with reference",
						References: []Reference{
							{Title: "Doc", Path: "/docs/guide.md"},
						},
					},
				},
				{
					{
						Rule: "Another rule",
						References: []Reference{
							{Title: "API", Path: "/api/spec.json"},
						},
					},
				},
			},
			want: []Instruction{
				{
					Rule: "Rule with reference",
					References: []Reference{
						{Title: "Doc", Path: "/docs/guide.md"},
					},
				},
				{
					Rule: "Another rule",
					References: []Reference{
						{Title: "API", Path: "/api/spec.json"},
					},
				},
			},
		},
		"no deduplication - keeps all instructions": {
			sources: [][]Instruction{
				{
					{Rule: "Same rule"},
				},
				{
					{Rule: "Same rule"},
				},
			},
			want: []Instruction{
				{Rule: "Same rule"},
				{Rule: "Same rule"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := MergeInstructions(tc.sources)

			if len(got) != len(tc.want) {
				t.Fatalf("length mismatch: got %d, want %d", len(got), len(tc.want))
			}

			for i := range got {
				if got[i].Heading != tc.want[i].Heading {
					t.Errorf("instruction %d heading: got %q, want %q", i, got[i].Heading, tc.want[i].Heading)
				}

				if got[i].Rule != tc.want[i].Rule {
					t.Errorf("instruction %d rule: got %q, want %q", i, got[i].Rule, tc.want[i].Rule)
				}

				if len(got[i].References) != len(tc.want[i].References) {
					t.Errorf("instruction %d references length: got %d, want %d", i, len(got[i].References), len(tc.want[i].References))
					continue
				}

				for j := range got[i].References {
					if got[i].References[j] != tc.want[i].References[j] {
						t.Errorf("instruction %d reference %d: got %v, want %v", i, j, got[i].References[j], tc.want[i].References[j])
					}
				}
			}
		})
	}
}

func TestMergeInstructionsOrder(t *testing.T) {
	// Verify that order is strictly preserved
	source1 := []Instruction{
		{Rule: "A1"},
		{Rule: "A2"},
		{Rule: "A3"},
	}

	source2 := []Instruction{
		{Rule: "B1"},
		{Rule: "B2"},
	}

	source3 := []Instruction{
		{Rule: "C1"},
	}

	result := MergeInstructions([][]Instruction{source1, source2, source3})

	expected := []string{"A1", "A2", "A3", "B1", "B2", "C1"}

	if len(result) != len(expected) {
		t.Fatalf("expected %d instructions, got %d", len(expected), len(result))
	}

	for i, want := range expected {
		if result[i].Rule != want {
			t.Errorf("position %d: got %q, want %q", i, result[i].Rule, want)
		}
	}
}
