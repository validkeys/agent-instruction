package rules

import (
	"errors"
	"fmt"
)

// Sentinel errors for rule validation
var (
	ErrTitleRequired        = errors.New("rule title is required")
	ErrInstructionsRequired = errors.New("must contain at least one instruction")
	ErrRuleTextRequired     = errors.New("rule text is required")
)

// RuleFile represents a rule file (.agent-instruction/rules/*.json)
type RuleFile struct {
	Title        string        `json:"title"`
	Instructions []Instruction `json:"instructions"`
	Imports      []string      `json:"imports,omitempty"`
}

// Instruction represents a single instruction rule
type Instruction struct {
	Heading    string      `json:"heading,omitempty"`
	Rule       string      `json:"rule"`
	References []Reference `json:"references,omitempty"`
}

// Reference represents a reference to another file or section
type Reference struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

// Validate checks rule file for required fields
func (r *RuleFile) Validate() error {
	if r.Title == "" {
		return ErrTitleRequired
	}

	if len(r.Instructions) == 0 {
		return ErrInstructionsRequired
	}

	for i, instr := range r.Instructions {
		if instr.Rule == "" {
			return fmt.Errorf("instruction %d: %w", i, ErrRuleTextRequired)
		}
	}

	return nil
}
