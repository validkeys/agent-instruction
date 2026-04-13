package commands

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// promptYesNo prompts user for a yes/no answer
func promptYesNo(cmd *cobra.Command, question string, defaultYes bool) bool {
	var defaultStr string
	if defaultYes {
		defaultStr = "Y/n"
	} else {
		defaultStr = "y/N"
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s (%s): ", question, defaultStr)

	reader := bufio.NewReader(cmd.InOrStdin())
	answer, err := reader.ReadString('\n')
	if err != nil {
		// On error, use default
		return defaultYes
	}

	answer = strings.TrimSpace(strings.ToLower(answer))

	// Empty answer uses default
	if answer == "" {
		return defaultYes
	}

	// Valid yes answers
	if answer == "y" || answer == "yes" {
		return true
	}

	// Valid no answers
	if answer == "n" || answer == "no" {
		return false
	}

	// Invalid answer uses default
	return defaultYes
}

// promptFrameworks prompts user to select frameworks
func promptFrameworks(cmd *cobra.Command) []string {
	fmt.Fprintf(cmd.OutOrStdout(), "Which AI frameworks do you want to support?\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  1) claude\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  2) agents\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  3) both\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Choice (default: 3): ")

	reader := bufio.NewReader(cmd.InOrStdin())
	answer, err := reader.ReadString('\n')
	if err != nil || strings.TrimSpace(answer) == "" {
		// Default to both
		return []string{"claude", "agents"}
	}

	answer = strings.TrimSpace(answer)
	switch answer {
	case "1":
		return []string{"claude"}
	case "2":
		return []string{"agents"}
	case "3", "":
		return []string{"claude", "agents"}
	default:
		// Invalid choice, use default
		fmt.Fprintf(cmd.OutOrStdout(), "Invalid choice, using default (both)\n")
		return []string{"claude", "agents"}
	}
}

// promptPackages prompts user for package discovery mode
func promptPackages(cmd *cobra.Command) []string {
	fmt.Fprintf(cmd.OutOrStdout(), "Package discovery mode?\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  1) auto - Automatically discover packages\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  2) manual - Specify package paths\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Choice (default: 1): ")

	reader := bufio.NewReader(cmd.InOrStdin())
	answer, err := reader.ReadString('\n')
	if err != nil || strings.TrimSpace(answer) == "" {
		// Default to auto
		return []string{"auto"}
	}

	answer = strings.TrimSpace(answer)
	switch answer {
	case "1", "":
		return []string{"auto"}
	case "2":
		// Prompt for manual package paths
		fmt.Fprintf(cmd.OutOrStdout(), "Enter comma-separated package paths (e.g., app,lib,services): ")
		paths, err := reader.ReadString('\n')
		if err != nil || strings.TrimSpace(paths) == "" {
			fmt.Fprintf(cmd.OutOrStdout(), "No packages specified, using auto\n")
			return []string{"auto"}
		}

		// Parse comma-separated paths
		parts := strings.Split(strings.TrimSpace(paths), ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}

		if len(result) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "No valid packages, using auto\n")
			return []string{"auto"}
		}

		return result
	default:
		// Invalid choice, use default
		fmt.Fprintf(cmd.OutOrStdout(), "Invalid choice, using default (auto)\n")
		return []string{"auto"}
	}
}
