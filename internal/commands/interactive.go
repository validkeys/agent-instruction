package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// ListRuleFiles scans the rules directory and returns a sorted list of .json
// file names without the extension.
func ListRuleFiles(rulesDir string) ([]string, error) {
	entries, err := os.ReadDir(rulesDir)
	if err != nil {
		return nil, fmt.Errorf("read rules directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) == ".json" {
			// Remove .json extension
			baseName := strings.TrimSuffix(name, ".json")
			files = append(files, baseName)
		}
	}

	// Sort alphabetically for consistent ordering
	sort.Strings(files)

	return files, nil
}

// PromptRuleFile displays a numbered list of available rule files and prompts
// the user to select one. Returns the selected file name (without extension).
func PromptRuleFile(available []string, in io.Reader, out io.Writer) (string, error) {
	if len(available) == 0 {
		return "", fmt.Errorf("no rule files available")
	}

	// Display available files
	fmt.Fprintln(out, "\nAvailable rule files:")
	for i, file := range available {
		fmt.Fprintf(out, "  %d. %s\n", i+1, file)
	}
	fmt.Fprint(out, "\nSelect rule file (1-"+strconv.Itoa(len(available))+"): ")

	// Read user input
	reader := bufio.NewReader(in)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	// Parse selection
	input = strings.TrimSpace(input)
	selection, err := strconv.Atoi(input)
	if err != nil {
		return "", fmt.Errorf("invalid selection: must be a number")
	}

	// Validate range (1-indexed)
	if selection < 1 || selection > len(available) {
		return "", fmt.Errorf("invalid selection: must be between 1 and %d", len(available))
	}

	// Return selected file (convert to 0-indexed)
	return available[selection-1], nil
}
