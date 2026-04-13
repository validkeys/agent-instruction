package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <rule-content>",
		Short: "Add a new instruction rule",
		Long: `Adds a new instruction to a rule file.

This command appends a new instruction to an existing rule file in the
.agent-instruction/rules directory. You can specify the target file with
the --rule flag, or select interactively if omitted.`,
		Example: `  # Add rule with title to global rules
  agent-instruction add "Use explicit error handling" --title "Error Handling" --rule global

  # Add rule with interactive file selection
  agent-instruction add "Always validate input"

  # Add rule to specific rule file
  agent-instruction add "Test all edge cases" --rule testing`,
		Args: cobra.ExactArgs(1),
		RunE: runAdd,
	}

	cmd.Flags().String("title", "", "optional heading for the rule")
	cmd.Flags().String("rule", "", "target rule file (e.g., 'global', 'testing')")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	// Get rule content from args
	ruleContent := strings.TrimSpace(args[0])
	if ruleContent == "" {
		return fmt.Errorf("rule content cannot be empty")
	}

	// Get flags
	title, _ := cmd.Flags().GetString("title")
	ruleFile, _ := cmd.Flags().GetString("rule")

	// TODO: Implement interactive selection when rule file not specified
	// TODO: Implement service integration to add instruction to file
	// TODO: Display success message with next steps

	// Placeholder - command structure is valid
	fmt.Fprintf(cmd.OutOrStdout(), "Rule content: %s\n", ruleContent)
	if title != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Title: %s\n", title)
	}
	if ruleFile != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Target: %s\n", ruleFile)
	}

	return nil
}
