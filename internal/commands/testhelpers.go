package commands

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/validkeys/agent-instruction/internal/config"
	"github.com/validkeys/agent-instruction/internal/rules"
)

// executeCommand runs a Cobra command and captures output
func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return buf.String(), err
}

// setupTestRepo creates a temporary repository structure
func setupTestRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	agentDir := filepath.Join(dir, ".agent-instruction")

	if err := os.MkdirAll(filepath.Join(agentDir, "rules"), 0755); err != nil {
		t.Fatalf("create test repo: %v", err)
	}

	return dir
}

// createConfig writes a config.json file for testing
func createConfig(t *testing.T, baseDir string, cfg config.Config) {
	t.Helper()

	path := filepath.Join(baseDir, ".agent-instruction", "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}

	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

// createRuleFile writes a rule file for testing
func createRuleFile(t *testing.T, baseDir, filename string, rule *rules.RuleFile) {
	t.Helper()

	path := filepath.Join(baseDir, ".agent-instruction", "rules", filename)
	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		t.Fatalf("marshal rule: %v", err)
	}

	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write rule: %v", err)
	}
}

// assertFileExists fails if file doesn't exist
func assertFileExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("expected file to exist: %s", path)
	}
}

// assertFileNotExists fails if file exists
func assertFileNotExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected file to not exist: %s", path)
	}
}

// readFile reads file content and fails on error
func readFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %s: %v", path, err)
	}
	return string(data)
}

// loadConfig loads and parses config.json from directory
func loadConfig(t *testing.T, baseDir string) *config.Config {
	t.Helper()

	path := filepath.Join(baseDir, ".agent-instruction", "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("unmarshal config: %v", err)
	}

	return &cfg
}

// loadRuleFile loads and parses a rule file
func loadRuleFile(t *testing.T, baseDir, filename string) *rules.RuleFile {
	t.Helper()

	path := filepath.Join(baseDir, ".agent-instruction", "rules", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read rule file: %v", err)
	}

	var rule rules.RuleFile
	if err := json.Unmarshal(data, &rule); err != nil {
		t.Fatalf("unmarshal rule file: %v", err)
	}

	return &rule
}
