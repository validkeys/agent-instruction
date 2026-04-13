package commands

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// FormatSuccess formats a success message with a checkmark
func FormatSuccess(msg string) string {
	return "✓ " + msg
}

// FormatError formats an error message with an X mark
func FormatError(msg string) string {
	return "✗ " + msg
}

// FormatHeading formats a heading with visual separators
func FormatHeading(text string) string {
	return fmt.Sprintf("\n=== %s ===\n", text)
}

// FormatList formats a list of items with bullet points
func FormatList(items []string) string {
	if len(items) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, item := range items {
		sb.WriteString("  - ")
		sb.WriteString(item)
		sb.WriteString("\n")
	}

	// Remove trailing newline
	result := sb.String()
	return strings.TrimSuffix(result, "\n")
}

// IndentText indents all lines of text by the specified number of spaces
func IndentText(text string, spaces int) string {
	if text == "" || spaces == 0 {
		return text
	}

	indent := strings.Repeat(" ", spaces)
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}

	return strings.Join(lines, "\n")
}

// isTerminal checks if the given file descriptor is a terminal
func isTerminal(fd int) bool {
	return term.IsTerminal(fd)
}

// IsColorSupported checks if the output supports ANSI color codes.
// Uses os.Stdout and environment variables directly - tests verify
// behavior via the isTerminal() wrapper rather than mocking global state.
func IsColorSupported() bool {
	// Check if stdout is a terminal
	if !isTerminal(int(os.Stdout.Fd())) {
		return false
	}

	// Check if NO_COLOR environment variable is set
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check TERM environment variable
	termEnv := os.Getenv("TERM")
	if termEnv == "dumb" {
		return false
	}

	return true
}

// Color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// Colorize applies color to text if colors are supported
func Colorize(text string, color string) string {
	if !IsColorSupported() {
		return text
	}
	return color + text + colorReset
}

// ColorSuccess returns green colored text
func ColorSuccess(text string) string {
	return Colorize(text, colorGreen)
}

// ColorError returns red colored text
func ColorError(text string) string {
	return Colorize(text, colorRed)
}

// ColorWarning returns yellow colored text
func ColorWarning(text string) string {
	return Colorize(text, colorYellow)
}

// ColorInfo returns cyan colored text
func ColorInfo(text string) string {
	return Colorize(text, colorCyan)
}

// ColorHeading returns blue colored text
func ColorHeading(text string) string {
	return Colorize(text, colorBlue)
}
