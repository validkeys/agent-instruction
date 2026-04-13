package rules

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRuleFileValidation(t *testing.T) {
	tests := map[string]struct {
		ruleFile RuleFile
		wantErr  bool
		errMsg   string
	}{
		"valid rule file": {
			ruleFile: RuleFile{
				Title: "Global Rules",
				Instructions: []Instruction{
					{Rule: "Always use error wrapping"},
				},
			},
			wantErr: false,
		},
		"valid rule file with heading and references": {
			ruleFile: RuleFile{
				Title: "API Rules",
				Instructions: []Instruction{
					{
						Heading: "Error Handling",
						Rule:    "Use fmt.Errorf with %w",
						References: []Reference{
							{Title: "Error Handling Guide", Path: "docs/errors.md"},
						},
					},
				},
			},
			wantErr: false,
		},
		"missing title": {
			ruleFile: RuleFile{
				Instructions: []Instruction{
					{Rule: "Some rule"},
				},
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		"missing instructions": {
			ruleFile: RuleFile{
				Title: "Rules",
			},
			wantErr: true,
			errMsg:  "at least one instruction",
		},
		"instruction missing rule text": {
			ruleFile: RuleFile{
				Title: "Rules",
				Instructions: []Instruction{
					{Heading: "Something"},
				},
			},
			wantErr: true,
			errMsg:  "rule text is required",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.ruleFile.Validate()

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr && tc.errMsg != "" {
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("error message %q does not contain %q", err.Error(), tc.errMsg)
				}
			}
		})
	}
}

func TestRuleFileMarshalUnmarshal(t *testing.T) {
	tests := map[string]struct {
		ruleFile RuleFile
	}{
		"basic rule file": {
			ruleFile: RuleFile{
				Title: "Global Rules",
				Instructions: []Instruction{
					{Rule: "Always validate input"},
				},
			},
		},
		"rule file with optional fields": {
			ruleFile: RuleFile{
				Title: "API Rules",
				Instructions: []Instruction{
					{
						Heading: "Validation",
						Rule:    "Check all inputs",
						References: []Reference{
							{Title: "Validation Guide", Path: "docs/validation.md"},
						},
					},
				},
				Imports: []string{"global.json"},
			},
		},
		"rule file with multiple instructions": {
			ruleFile: RuleFile{
				Title: "Multiple Rules",
				Instructions: []Instruction{
					{Rule: "Rule 1"},
					{Rule: "Rule 2"},
					{Rule: "Rule 3"},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tc.ruleFile)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			// Unmarshal back
			var decoded RuleFile
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			// Compare
			if decoded.Title != tc.ruleFile.Title {
				t.Errorf("title: got %q, want %q", decoded.Title, tc.ruleFile.Title)
			}

			if len(decoded.Instructions) != len(tc.ruleFile.Instructions) {
				t.Errorf("instructions count: got %d, want %d", len(decoded.Instructions), len(tc.ruleFile.Instructions))
			}
		})
	}
}

func TestInstructionWithOmitEmpty(t *testing.T) {
	tests := map[string]struct {
		instruction Instruction
		wantFields  []string // Fields that should appear in JSON
	}{
		"instruction with all fields": {
			instruction: Instruction{
				Heading: "Test",
				Rule:    "Test rule",
				References: []Reference{
					{Title: "Doc", Path: "path"},
				},
			},
			wantFields: []string{"heading", "rule", "references"},
		},
		"instruction with only rule": {
			instruction: Instruction{
				Rule: "Test rule",
			},
			wantFields: []string{"rule"}, // heading and references should be omitted
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			data, err := json.Marshal(tc.instruction)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			jsonStr := string(data)
			for _, field := range tc.wantFields {
				if !strings.Contains(jsonStr, field) {
					t.Errorf("expected field %q in JSON: %s", field, jsonStr)
				}
			}

			// Check omitempty works for optional fields
			if tc.instruction.Heading == "" && strings.Contains(jsonStr, "heading") {
				t.Errorf("empty heading should be omitted from JSON: %s", jsonStr)
			}
		})
	}
}
