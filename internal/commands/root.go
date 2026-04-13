package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v1.0.0-dev"

// NewRootCmd creates and returns the root command
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "agent-instruction",
		Short: "Manage AI agent instructions in monorepos",
		Long: `agent-instruction is a CLI tool for managing CLAUDE.md and AGENTS.md files
across multiple packages in a monorepo.

It provides a declarative configuration system that builds instruction files
from modular rule sets, ensuring consistent AI agent behavior across your
entire project.`,
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			// When no subcommand is provided, show help
			return cmd.Help()
		},
	}

	// Set version template
	rootCmd.SetVersionTemplate(fmt.Sprintf("agent-instruction version %s\n", version))

	// Global flags can be added here
	// rootCmd.PersistentFlags().Bool("debug", false, "enable debug output")

	// Register subcommands
	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newBuildCmd())

	return rootCmd
}

// Execute runs the root command
func Execute() error {
	return NewRootCmd().Execute()
}
