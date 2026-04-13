package builder

import (
	"fmt"
	"strings"

	"github.com/validkeys/agent-instruction/internal/rules"
)

// InstructionsToMarkdown converts a slice of instructions to formatted markdown.
// It handles headings, rules, and references with proper spacing.
func InstructionsToMarkdown(instructions []rules.Instruction) string {
	if len(instructions) == 0 {
		return ""
	}

	var sb strings.Builder

	for _, instr := range instructions {
		// Add heading if present
		if instr.Heading != "" {
			sb.WriteString("## ")
			sb.WriteString(instr.Heading)
			sb.WriteString("\n\n")
		}

		// Add rule text
		sb.WriteString(instr.Rule)
		sb.WriteString("\n\n")

		// Add references if present
		if len(instr.References) > 0 {
			for _, ref := range instr.References {
				sb.WriteString(formatReference(ref))
				sb.WriteString("\n")
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// formatReference formats a reference as a markdown link in a bullet list.
func formatReference(ref rules.Reference) string {
	return fmt.Sprintf("- [%s](%s)", ref.Title, ref.Path)
}
