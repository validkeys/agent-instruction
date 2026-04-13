package files

import (
	"fmt"
	"strings"
)

const (
	BeginMarker = "<!-- BEGIN AGENT-INSTRUCTION -->"
	EndMarker   = "<!-- END AGENT-INSTRUCTION -->"
)

// ManagedContent represents parsed content with managed sections
type ManagedContent struct {
	Before  string // Content before managed section
	Managed string // Content within managed section
	After   string // Content after managed section
}

// ParseManagedContent parses content and extracts managed sections
func ParseManagedContent(content string) (*ManagedContent, error) {
	beginIdx := strings.Index(content, BeginMarker)
	endIdx := strings.Index(content, EndMarker)

	// No markers present - all content is "before"
	if beginIdx == -1 && endIdx == -1 {
		return &ManagedContent{
			Before:  content,
			Managed: "",
			After:   "",
		}, nil
	}

	// Only one marker present
	if beginIdx == -1 {
		return nil, fmt.Errorf("malformed managed section: end marker found without begin marker")
	}
	if endIdx == -1 {
		return nil, fmt.Errorf("malformed managed section: begin marker found without end marker")
	}

	// End marker before begin marker
	if beginIdx > endIdx {
		return nil, fmt.Errorf("malformed managed section: end marker appears before begin marker")
	}

	// Extract sections
	before := content[:beginIdx]
	managedStart := beginIdx + len(BeginMarker)
	managed := content[managedStart:endIdx]
	after := content[endIdx+len(EndMarker):]

	return &ManagedContent{
		Before:  before,
		Managed: managed,
		After:   after,
	}, nil
}

// HasManagedSection checks if content has valid managed section markers
func HasManagedSection(content string) bool {
	beginIdx := strings.Index(content, BeginMarker)
	endIdx := strings.Index(content, EndMarker)

	// Both markers must be present and in correct order
	return beginIdx != -1 && endIdx != -1 && beginIdx < endIdx
}

// ReplaceManagedSection replaces or adds managed section in content
func ReplaceManagedSection(content string, newSection string) (string, error) {
	beginIdx := strings.Index(content, BeginMarker)
	endIdx := strings.Index(content, EndMarker)

	// No managed section - append to end
	if beginIdx == -1 && endIdx == -1 {
		return content + "\n\n" + BeginMarker + "\n" + newSection + "\n" + EndMarker, nil
	}

	// Only one marker present
	if beginIdx == -1 {
		return "", fmt.Errorf("malformed managed section: end marker found without begin marker")
	}
	if endIdx == -1 {
		return "", fmt.Errorf("malformed managed section: begin marker found without end marker")
	}

	// End marker before begin marker
	if beginIdx > endIdx {
		return "", fmt.Errorf("malformed managed section: end marker appears before begin marker")
	}

	// Replace content between markers
	before := content[:beginIdx+len(BeginMarker)]
	after := content[endIdx:]
	return before + "\n" + newSection + "\n" + after, nil
}
