# JSON Config Handling

**Purpose:** JSON configuration parsing and validation
**Source:** Go encoding/json standard library and project data model
**Use cases:** Config service, rule service, all JSON file operations

---

## Struct Definitions with JSON Tags

```go
// Config represents .agent-instruction/config.json
type Config struct {
    Version    string   `json:"version"`
    Packages   []string `json:"packages"`
    Frameworks []string `json:"frameworks"`
}

// RuleFile represents a rule file (.agent-instruction/rules/*.json)
type RuleFile struct {
    Title        string        `json:"title"`
    Instructions []Instruction `json:"instructions"`
    Imports      []string      `json:"imports,omitempty"`
}

// Instruction represents a single instruction rule
type Instruction struct {
    Heading    string      `json:"heading,omitempty"`
    Rule       string      `json:"rule"`
    References []Reference `json:"references,omitempty"`
}

// Reference represents a reference to another file or section
type Reference struct {
    Title string `json:"title"`
    Path  string `json:"path"`
}
```

---

## Loading JSON Files

```go
// LoadConfig reads and validates config.json
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("config not found: %s (run 'agent-instruction init')", path)
        }
        return nil, fmt.Errorf("read config: %w", err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config JSON: %w", err)
    }

    if err := validateConfig(&cfg); err != nil {
        return nil, fmt.Errorf("validate config: %w", err)
    }

    return &cfg, nil
}

// LoadRuleFile reads and validates a rule file
func LoadRuleFile(path string) (*RuleFile, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read rule file: %w", err)
    }

    var rule RuleFile
    if err := json.Unmarshal(data, &rule); err != nil {
        return nil, fmt.Errorf("parse rule JSON in %s: %w", path, err)
    }

    if err := validateRuleFile(path, &rule); err != nil {
        return nil, err
    }

    return &rule, nil
}
```

---

## Validation Functions

```go
// validateConfig checks config for required fields and valid values
func validateConfig(cfg *Config) error {
    if cfg.Version == "" {
        return fmt.Errorf("version is required")
    }

    if len(cfg.Frameworks) == 0 {
        return fmt.Errorf("at least one framework is required")
    }

    validFrameworks := map[string]bool{
        "claude": true,
        "agents": true,
    }

    for _, fw := range cfg.Frameworks {
        if !validFrameworks[fw] {
            return fmt.Errorf("invalid framework: %s (must be 'claude' or 'agents')", fw)
        }
    }

    return nil
}

// validateRuleFile checks rule file for required fields
func validateRuleFile(path string, rule *RuleFile) error {
    if rule.Title == "" {
        return fmt.Errorf("rule file %s: title is required", path)
    }

    if len(rule.Instructions) == 0 {
        return fmt.Errorf("rule file %s: must contain at least one instruction", path)
    }

    for i, instr := range rule.Instructions {
        if instr.Rule == "" {
            return fmt.Errorf("rule file %s: instruction %d: rule text is required", path, i)
        }
    }

    return nil
}
```

---

## Saving JSON Files

```go
// SaveConfig writes config to disk with proper formatting
func SaveConfig(path string, cfg *Config) error {
    // Validate before saving
    if err := validateConfig(cfg); err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }

    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return fmt.Errorf("marshal config: %w", err)
    }

    // Add trailing newline for POSIX compliance
    data = append(data, '\n')

    if err := AtomicWrite(path, data, 0644); err != nil {
        return fmt.Errorf("write config: %w", err)
    }

    return nil
}

// SaveRuleFile writes rule file to disk
func SaveRuleFile(path string, rule *RuleFile) error {
    if err := validateRuleFile(path, rule); err != nil {
        return err
    }

    data, err := json.MarshalIndent(rule, "", "  ")
    if err != nil {
        return fmt.Errorf("marshal rule file: %w", err)
    }

    data = append(data, '\n')

    if err := AtomicWrite(path, data, 0644); err != nil {
        return fmt.Errorf("write rule file: %w", err)
    }

    return nil
}
```

---

## Optional Fields with omitempty

✅ **Good: Use omitempty for optional fields**

```go
type Instruction struct {
    Heading    string      `json:"heading,omitempty"`    // Optional
    Rule       string      `json:"rule"`                 // Required
    References []Reference `json:"references,omitempty"` // Optional
}

// Result: If Heading is empty or References is nil, they won't appear in JSON
```

❌ **Bad: Don't use omitempty for required fields**

```go
type Config struct {
    Version string `json:"version,omitempty"` // Bad: version should always be present
}
```

---

## Nested Struct Composition

```go
// Good: Compose complex structures from smaller pieces
type ProjectConfig struct {
    Config                // Embedded struct
    Metadata   Metadata   `json:"metadata"`
    Repository Repository `json:"repository"`
}

type Metadata struct {
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Repository struct {
    Type string `json:"type"` // "git", "monorepo"
    Root string `json:"root"`
}
```

---

## Handling Default Values

```go
// SetDefaults applies default values to config
func (cfg *Config) SetDefaults() {
    if cfg.Version == "" {
        cfg.Version = "1.0"
    }

    if len(cfg.Frameworks) == 0 {
        cfg.Frameworks = []string{"claude"}
    }

    if cfg.Packages == nil {
        cfg.Packages = []string{}
    }
}

// LoadConfigWithDefaults loads config and applies defaults
func LoadConfigWithDefaults(path string) (*Config, error) {
    cfg, err := LoadConfig(path)
    if err != nil {
        return nil, err
    }

    cfg.SetDefaults()
    return cfg, nil
}
```

---

## JSON Merge Pattern

```go
// MergeConfigs merges override config into base config
func MergeConfigs(base, override *Config) *Config {
    merged := *base // Copy base

    if override.Version != "" {
        merged.Version = override.Version
    }

    if len(override.Frameworks) > 0 {
        merged.Frameworks = override.Frameworks
    }

    // Append packages (don't replace)
    merged.Packages = append(merged.Packages, override.Packages...)

    return &merged
}
```

---

## Pretty Printing for Debugging

```go
// PrettyPrint outputs formatted JSON for debugging
func PrettyPrint(v interface{}) string {
    data, err := json.MarshalIndent(v, "", "  ")
    if err != nil {
        return fmt.Sprintf("error: %v", err)
    }
    return string(data)
}

// Usage in tests or debugging
func debugConfig(cfg *Config) {
    fmt.Println("Config:", PrettyPrint(cfg))
}
```

---

## Error Handling Patterns

✅ **Good: Specific error messages for JSON parsing**

```go
func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config file: %w", err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        // Check for syntax errors
        if syntaxErr, ok := err.(*json.SyntaxError); ok {
            return nil, fmt.Errorf("JSON syntax error at byte %d: %w", syntaxErr.Offset, err)
        }
        return nil, fmt.Errorf("parse JSON: %w", err)
    }

    return &cfg, nil
}
```

❌ **Bad: Generic error message**

```go
func loadConfig(path string) (*Config, error) {
    data, _ := os.ReadFile(path)
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err // User doesn't know what went wrong
    }
    return &cfg, nil
}
```

---

## Testing JSON Marshaling

```go
func TestConfigMarshalUnmarshal(t *testing.T) {
    original := &Config{
        Version:    "1.0",
        Frameworks: []string{"claude", "agents"},
        Packages:   []string{"api", "web"},
    }

    // Marshal to JSON
    data, err := json.Marshal(original)
    if err != nil {
        t.Fatalf("marshal: %v", err)
    }

    // Unmarshal back
    var decoded Config
    if err := json.Unmarshal(data, &decoded); err != nil {
        t.Fatalf("unmarshal: %v", err)
    }

    // Compare
    if decoded.Version != original.Version {
        t.Errorf("version: got %q, want %q", decoded.Version, original.Version)
    }

    if len(decoded.Frameworks) != len(original.Frameworks) {
        t.Errorf("frameworks count: got %d, want %d", len(decoded.Frameworks), len(original.Frameworks))
    }
}
```

---

## Common Patterns

```go
// Pattern: Load all rule files from directory
func LoadAllRules(rulesDir string) (map[string]*RuleFile, error) {
    rules := make(map[string]*RuleFile)

    entries, err := os.ReadDir(rulesDir)
    if err != nil {
        if os.IsNotExist(err) {
            return rules, nil // Empty map, not an error
        }
        return nil, fmt.Errorf("read rules directory: %w", err)
    }

    for _, entry := range entries {
        if !strings.HasSuffix(entry.Name(), ".json") {
            continue
        }

        path := filepath.Join(rulesDir, entry.Name())
        rule, err := LoadRuleFile(path)
        if err != nil {
            return nil, fmt.Errorf("load %s: %w", entry.Name(), err)
        }

        rules[entry.Name()] = rule
    }

    return rules, nil
}

// Pattern: Update specific field in config
func UpdateConfigFrameworks(configPath string, frameworks []string) error {
    cfg, err := LoadConfig(configPath)
    if err != nil {
        return err
    }

    cfg.Frameworks = frameworks

    if err := SaveConfig(configPath, cfg); err != nil {
        return err
    }

    return nil
}
```

---

## Key Principles

1. **Always validate after unmarshaling** - JSON can be syntactically valid but semantically wrong
2. **Use `omitempty` for optional fields** - Cleaner JSON output
3. **Marshal with indent for config files** - Human-readable
4. **Add trailing newline when writing** - POSIX compliance
5. **Wrap errors with context** - Include filename in error messages
6. **Validate before marshaling** - Catch errors early

---

## References

- Go encoding/json: https://pkg.go.dev/encoding/json
- JSON and Go: https://go.dev/blog/json
- Project data model: `/Users/kydavis/Sites/agent-instruction/docs/plan/001-initial-buildout/technical-requirements.yaml` (lines 101-143)
