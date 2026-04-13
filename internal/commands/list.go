package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/validkeys/agent-instruction/internal/rules"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all instruction rules",
		Long: `Displays all rule files and their contents in a clear, readable format.

This command scans the .agent-instruction/rules directory and shows all
available rule files with their instruction counts. Use --verbose to see
the full content of each instruction.`,
		Example: `  # List all rules with summary
  agent-instruction list

  # List all rules with full content
  agent-instruction list --verbose`,
		Args: cobra.NoArgs,
		RunE: runList,
	}

	cmd.Flags().Bool("verbose", false, "show full rule content")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	// Get current directory
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Check if initialized
	rulesDir := filepath.Join(baseDir, ".agent-instruction", "rules")
	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		return fmt.Errorf("not initialized: run 'agent-instruction init' first")
	}

	// Get verbose flag
	verbose, _ := cmd.Flags().GetBool("verbose")

	// List rule files
	ruleFiles, err := ListRuleFiles(rulesDir)
	if err != nil {
		return fmt.Errorf("list rule files: %w", err)
	}

	if len(ruleFiles) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "No rule files found in %s\n", rulesDir)
		return nil
	}

	// Create service only when we have files to process
	configSvc := rules.NewFileConfigService(rulesDir)

	// Display each rule file
	for i, ruleFileName := range ruleFiles {
		ruleFilePath := filepath.Join(rulesDir, ruleFileName+".json")

		// Load rule file
		ruleFile, err := configSvc.LoadRuleFile(ruleFilePath)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Warning: failed to load %s: %v\n", ruleFileName+".json", err)
			continue
		}

		// Add spacing between files (except before first)
		if i > 0 {
			fmt.Fprintln(cmd.OutOrStdout())
		}

		// Display file header
		fmt.Fprintf(cmd.OutOrStdout(), "📄 %s\n", ruleFileName+".json")

		if verbose {
			// Show full content
			fmt.Fprintf(cmd.OutOrStdout(), "   Title: %s\n", ruleFile.Title)

			if len(ruleFile.Imports) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "   Imports: %d\n", len(ruleFile.Imports))
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\n")

			// Show each instruction
			for idx, instr := range ruleFile.Instructions {
				fmt.Fprintf(cmd.OutOrStdout(), "   [%d] ", idx+1)

				if instr.Heading != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n", instr.Heading)
					fmt.Fprintf(cmd.OutOrStdout(), "       %s\n", instr.Rule)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n", instr.Rule)
				}

				if len(instr.References) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "       References: %d\n", len(instr.References))
				}
			}
		} else {
			// Show summary
			fmt.Fprintf(cmd.OutOrStdout(), "   %d instruction(s)\n", len(ruleFile.Instructions))

			// Show instruction headings if available
			for _, instr := range ruleFile.Instructions {
				if instr.Heading != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "   - %s\n", instr.Heading)
				}
			}
		}
	}

	return nil
}
