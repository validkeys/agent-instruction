package builder

import (
	"strings"

	"github.com/validkeys/agent-instruction/internal/files"
)

// WrapWithMarkers wraps content with HTML comment markers for managed sections.
// The markers indicate the boundaries of auto-generated content that can be
// safely replaced without affecting user-written content outside the markers.
func WrapWithMarkers(content string) string {
	var sb strings.Builder

	sb.WriteString(files.BeginMarker)
	sb.WriteString("\n")
	sb.WriteString(content)

	// Ensure proper spacing before end marker
	// Content should end with \n\n for spacing from END marker
	if strings.HasSuffix(content, "\n\n") {
		// Already has double newline - don't add more
	} else if strings.HasSuffix(content, "\n") {
		// Has single newline - add one more for spacing
		sb.WriteString("\n")
	} else {
		// No trailing newline - add one
		sb.WriteString("\n")
	}

	sb.WriteString(files.EndMarker)

	return sb.String()
}

// BuildManagedFile combines generated content with existing user content.
// It wraps the generated content with markers and preserves any user content
// that exists before or after the managed section.
//
// If existing is nil (new file), returns just the wrapped generated content.
// If existing has user content before or after the managed section, that
// content is preserved in the output.
func BuildManagedFile(generated string, existing *files.ManagedContent) string {
	wrapped := WrapWithMarkers(generated)

	// New file - no existing content
	if existing == nil {
		return wrapped
	}

	// Existing file - combine user content with new managed section
	var sb strings.Builder

	// Add content before managed section
	if existing.Before != "" {
		sb.WriteString(existing.Before)
	}

	// Add wrapped generated content
	sb.WriteString(wrapped)

	// Add content after managed section
	if existing.After != "" {
		sb.WriteString(existing.After)
	}

	return sb.String()
}
