package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
	"github.com/validkeys/agent-instruction/internal/builder"
	"github.com/validkeys/agent-instruction/internal/config"
	"github.com/validkeys/agent-instruction/internal/files"
	"github.com/validkeys/agent-instruction/internal/rules"
)

func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build instruction files from rules",
		Long: `Generates CLAUDE.md and/or AGENTS.md files for all packages in the monorepo.

This command discovers packages (either automatically or from the packages list
in config.json), composes instructions from global and package-level rules,
and generates instruction files for each configured framework.`,
		Example: `  # Build all packages
  agent-instruction build

  # Dry-run to see what would be generated
  agent-instruction build --dry-run

  # Verbose output with progress
  agent-instruction build --verbose

  # Disable parallel processing
  agent-instruction build --no-parallel`,
		Args: cobra.NoArgs,
		RunE: runBuild,
	}

	cmd.Flags().Bool("dry-run", false, "preview changes without writing files")
	cmd.Flags().Bool("verbose", false, "show detailed progress output")
	cmd.Flags().Bool("no-parallel", false, "disable parallel package processing")

	return cmd
}

func runBuild(cmd *cobra.Command, args []string) error {
	// Get current directory
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Check if initialized
	agentDir := filepath.Join(baseDir, ".agent-instruction")
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		return fmt.Errorf("not initialized: run 'agent-instruction init' first")
	}

	// Get flags
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	verbose, _ := cmd.Flags().GetBool("verbose")
	noParallel, _ := cmd.Flags().GetBool("no-parallel")

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "DRY RUN: Preview mode - no files will be written\n\n")
	}

	// Load config
	configPath := filepath.Join(agentDir, "config.json")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Discover packages
	if verbose {
		fmt.Fprintf(cmd.OutOrStdout(), "Discovering packages...\n")
	}

	packages, err := builder.DiscoverPackages(cfg, baseDir)
	if err != nil {
		return fmt.Errorf("discover packages: %w", err)
	}

	if len(packages) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "✓ No packages found\n")
		return nil
	}

	if verbose {
		fmt.Fprintf(cmd.OutOrStdout(), "Found %d package(s):\n", len(packages))
		for _, pkg := range packages {
			relPath, _ := filepath.Rel(baseDir, pkg)
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", relPath)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "\n")
	}

	// Create services
	rulesDir := filepath.Join(agentDir, "rules")
	configSvc := rules.NewFileConfigService(rulesDir)
	ruleSvc := rules.NewRuleService(configSvc)

	// Build packages
	start := time.Now()
	var successCount, errorCount atomic.Int32

	// Process function for each package
	processPackage := func(ctx context.Context, pkgPath string) error {
		relPath, _ := filepath.Rel(baseDir, pkgPath)

		if verbose {
			fmt.Fprintf(cmd.OutOrStdout(), "Processing %s...\n", relPath)
		}

		// Compose instructions
		globalRulesPath := filepath.Join(rulesDir, "global.json")
		packageConfigPath := filepath.Join(pkgPath, "agent-instruction.json")

		instructions, err := builder.ComposeInstructions(globalRulesPath, packageConfigPath, ruleSvc)
		if err != nil {
			return fmt.Errorf("compose instructions for %s: %w", relPath, err)
		}

		// Generate file for each framework
		for _, framework := range cfg.Frameworks {
			var filename string
			switch framework {
			case "claude":
				filename = "CLAUDE.md"
			case "agents":
				filename = "AGENTS.md"
			default:
				return fmt.Errorf("unknown framework: %s", framework)
			}

			outputPath := filepath.Join(pkgPath, filename)

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "  Would generate: %s\n", filepath.Join(relPath, filename))
				continue
			}

			// Generate file content
			generatedContent := builder.InstructionsToMarkdown(instructions)

			// Read existing file if it exists
			var existingContent *files.ManagedContent
			if existingData, err := os.ReadFile(outputPath); err == nil {
				existingContent, err = files.ParseManagedContent(string(existingData))
				if err != nil {
					relPath, _ := filepath.Rel(baseDir, outputPath)
					return fmt.Errorf("parse existing file %s: %w", relPath, err)
				}
			}

			// Build file with managed sections
			finalContent := builder.BuildManagedFile(generatedContent, existingContent)

			// Ensure output directory exists
			outputDir := filepath.Dir(outputPath)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("create directory: %w", err)
			}

			// Write to file
			relPath, _ := filepath.Rel(baseDir, outputPath)
			if err := os.WriteFile(outputPath, []byte(finalContent), 0644); err != nil {
				return fmt.Errorf("write %s: %w", relPath, err)
			}

			if verbose {
				fmt.Fprintf(cmd.OutOrStdout(), "  ✓ Generated %s\n", filename)
			}
		}

		successCount.Add(1)
		return nil
	}

	// Execute build (parallel or sequential)
	ctx := cmd.Context()

	if noParallel || len(packages) == 1 {
		// Sequential processing
		for _, pkg := range packages {
			if err := processPackage(ctx, pkg); err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
				errorCount.Add(1)
			}
		}
	} else {
		// Parallel processing
		if verbose {
			fmt.Fprintf(cmd.OutOrStdout(), "Building packages in parallel...\n\n")
		}

		if err := builder.ProcessPackagesParallel(ctx, packages, processPackage); err != nil {
			// ProcessPackagesParallel returns on first error
			errorCount.Add(1)
			fmt.Fprintf(cmd.OutOrStderr(), "Build failed: %v\n", err)
		}
	}

	elapsed := time.Since(start)

	// Print summary
	fmt.Fprintf(cmd.OutOrStdout(), "\n")
	success := int(successCount.Load())
	errors := int(errorCount.Load())

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "✓ DRY RUN: Checked %d package(s) in %v\n", len(packages), elapsed.Round(time.Millisecond))
	} else if errors == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "✓ Successfully processed %d package(s) in %v\n", success, elapsed.Round(time.Millisecond))
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "⚠ Processed %d package(s) with %d error(s) in %v\n", success, errors, elapsed.Round(time.Millisecond))
		return fmt.Errorf("build completed with errors")
	}

	return nil
}
