# Graph Traversal and Cycle Detection

**Purpose:** Import resolution algorithm with cycle detection
**Source:** Standard computer science algorithms
**Use cases:** Rule import resolution, dependency management

---

## Complete Import Resolver

```go
// ImportResolver handles rule import resolution with cycle detection
type ImportResolver struct {
    visited   map[string]bool
    pathStack []string
    baseDir   string
}

// NewImportResolver creates a new resolver
func NewImportResolver(baseDir string) *ImportResolver {
    return &ImportResolver{
        visited:   make(map[string]bool),
        pathStack: make([]string, 0),
        baseDir:   baseDir,
    }
}

// Resolve recursively resolves imports and returns merged instructions
func (r *ImportResolver) Resolve(rulePath string) ([]Instruction, error) {
    // Make path absolute
    absPath := r.resolveRelativePath(rulePath)

    // Check for cycle
    if r.inPathStack(absPath) {
        return nil, fmt.Errorf("import cycle detected: %s", r.formatCyclePath(absPath))
    }

    // Skip if already visited (but not a cycle)
    if r.visited[absPath] {
        return nil, nil
    }

    // Mark as visited and add to path stack
    r.visited[absPath] = true
    r.pathStack = append(r.pathStack, absPath)
    defer func() {
        // Remove from path stack when done
        r.pathStack = r.pathStack[:len(r.pathStack)-1]
    }()

    // Load rule file
    rule, err := loadRuleFile(absPath)
    if err != nil {
        return nil, err
    }

    // Collect all instructions (imports first, then local)
    var allInstructions []Instruction

    // Recursively resolve imports (depth-first)
    for _, importPath := range rule.Imports {
        importedInstructions, err := r.Resolve(importPath)
        if err != nil {
            return nil, err
        }
        allInstructions = append(allInstructions, importedInstructions...)
    }

    // Add local instructions after imports
    allInstructions = append(allInstructions, rule.Instructions...)

    return allInstructions, nil
}

// inPathStack checks if path is currently in the traversal path
func (r *ImportResolver) inPathStack(path string) bool {
    for _, p := range r.pathStack {
        if p == path {
            return true
        }
    }
    return false
}

// formatCyclePath creates a readable cycle error message
func (r *ImportResolver) formatCyclePath(cyclePath string) string {
    path := append(r.pathStack, cyclePath)
    return strings.Join(path, " → ")
}

// resolveRelativePath converts relative path to absolute based on baseDir
func (r *ImportResolver) resolveRelativePath(path string) string {
    if filepath.IsAbs(path) {
        return path
    }
    return filepath.Join(r.baseDir, path)
}
```

---

## Depth-First Search Pattern

```go
// DFS with visited set - prevents revisiting nodes
func dfs(node string, graph map[string][]string, visited map[string]bool, result *[]string) {
    if visited[node] {
        return // Already processed
    }

    visited[node] = true
    *result = append(*result, node)

    // Process neighbors
    for _, neighbor := range graph[node] {
        dfs(neighbor, graph, visited, result)
    }
}

// Usage
func traverseGraph(start string, graph map[string][]string) []string {
    visited := make(map[string]bool)
    result := make([]string, 0)
    dfs(start, graph, visited, &result)
    return result
}
```

---

## Cycle Detection with Path Stack

```go
// detectCycle uses DFS with path tracking
func detectCycle(node string, graph map[string][]string, visited, inPath map[string]bool) (bool, []string) {
    if inPath[node] {
        // Found cycle - node is in current path
        return true, []string{node}
    }

    if visited[node] {
        // Already fully explored, no cycle from here
        return false, nil
    }

    // Mark as visited and add to current path
    visited[node] = true
    inPath[node] = true
    defer func() {
        inPath[node] = false // Remove from path when backtracking
    }()

    // Check all neighbors
    for _, neighbor := range graph[node] {
        if hasCycle, cyclePath := detectCycle(neighbor, graph, visited, inPath); hasCycle {
            // Build cycle path
            return true, append(cyclePath, node)
        }
    }

    return false, nil
}

// Usage
func hasCycle(graph map[string][]string, start string) (bool, string) {
    visited := make(map[string]bool)
    inPath := make(map[string]bool)

    hasCycle, path := detectCycle(start, graph, visited, inPath)
    if hasCycle {
        // Reverse path for better readability
        for i := len(path)/2 - 1; i >= 0; i-- {
            j := len(path) - 1 - i
            path[i], path[j] = path[j], path[i]
        }
        return true, strings.Join(path, " → ")
    }

    return false, ""
}
```

---

## Topological Sort Pattern

```go
// TopologicalSort orders nodes so dependencies come before dependents
func TopologicalSort(graph map[string][]string) ([]string, error) {
    visited := make(map[string]bool)
    inPath := make(map[string]bool)
    result := make([]string, 0)

    var visit func(node string) error
    visit = func(node string) error {
        if inPath[node] {
            return fmt.Errorf("cycle detected at node: %s", node)
        }

        if visited[node] {
            return nil
        }

        visited[node] = true
        inPath[node] = true

        // Visit all dependencies first
        for _, dep := range graph[node] {
            if err := visit(dep); err != nil {
                return err
            }
        }

        inPath[node] = false
        result = append(result, node) // Add after dependencies

        return nil
    }

    // Visit all nodes
    for node := range graph {
        if !visited[node] {
            if err := visit(node); err != nil {
                return nil, err
            }
        }
    }

    return result, nil
}
```

---

## Complete Example: Import Graph

```go
// ImportGraph manages rule dependencies
type ImportGraph struct {
    nodes map[string]*RuleFile
    edges map[string][]string // node -> dependencies
}

// NewImportGraph creates a graph from rules directory
func NewImportGraph(rulesDir string) (*ImportGraph, error) {
    g := &ImportGraph{
        nodes: make(map[string]*RuleFile),
        edges: make(map[string][]string),
    }

    // Load all rule files
    entries, err := os.ReadDir(rulesDir)
    if err != nil {
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

        g.nodes[entry.Name()] = rule
        g.edges[entry.Name()] = rule.Imports
    }

    return g, nil
}

// DetectCycles finds all cycles in the graph
func (g *ImportGraph) DetectCycles() ([]string, error) {
    visited := make(map[string]bool)
    inPath := make(map[string]bool)
    cycles := make([]string, 0)

    var visit func(node string, path []string) error
    visit = func(node string, path []string) error {
        if inPath[node] {
            // Found cycle
            cyclePath := append(path, node)
            cycles = append(cycles, strings.Join(cyclePath, " → "))
            return nil // Continue to find all cycles
        }

        if visited[node] {
            return nil
        }

        visited[node] = true
        inPath[node] = true
        path = append(path, node)

        for _, dep := range g.edges[node] {
            if err := visit(dep, path); err != nil {
                return err
            }
        }

        inPath[node] = false
        return nil
    }

    for node := range g.nodes {
        if !visited[node] {
            if err := visit(node, []string{}); err != nil {
                return nil, err
            }
        }
    }

    if len(cycles) > 0 {
        return cycles, fmt.Errorf("found %d cycle(s)", len(cycles))
    }

    return nil, nil
}

// ResolveOrder returns rules in dependency order
func (g *ImportGraph) ResolveOrder() ([]string, error) {
    order, err := TopologicalSort(g.edges)
    if err != nil {
        return nil, fmt.Errorf("resolve import order: %w", err)
    }
    return order, nil
}
```

---

## Testing Cycle Detection

```go
func TestImportCycleDetection(t *testing.T) {
    tests := map[string]struct {
        rules    map[string]*RuleFile
        wantCycle bool
        cyclePath string
    }{
        "no cycle - linear chain": {
            rules: map[string]*RuleFile{
                "a.json": {Title: "A", Imports: []string{"b.json"}},
                "b.json": {Title: "B", Imports: []string{"c.json"}},
                "c.json": {Title: "C", Imports: []string{}},
            },
            wantCycle: false,
        },
        "simple cycle - A imports A": {
            rules: map[string]*RuleFile{
                "a.json": {Title: "A", Imports: []string{"a.json"}},
            },
            wantCycle: true,
            cyclePath: "a.json → a.json",
        },
        "cycle in chain - A→B→C→A": {
            rules: map[string]*RuleFile{
                "a.json": {Title: "A", Imports: []string{"b.json"}},
                "b.json": {Title: "B", Imports: []string{"c.json"}},
                "c.json": {Title: "C", Imports: []string{"a.json"}},
            },
            wantCycle: true,
            cyclePath: "a.json → b.json → c.json → a.json",
        },
        "diamond dependency - no cycle": {
            rules: map[string]*RuleFile{
                "a.json": {Title: "A", Imports: []string{"b.json", "c.json"}},
                "b.json": {Title: "B", Imports: []string{"d.json"}},
                "c.json": {Title: "C", Imports: []string{"d.json"}},
                "d.json": {Title: "D", Imports: []string{}},
            },
            wantCycle: false,
        },
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            // Setup test rules
            dir := setupTestRules(t, tc.rules)

            // Create import graph
            graph, err := NewImportGraph(dir)
            if err != nil {
                t.Fatalf("create graph: %v", err)
            }

            // Detect cycles
            cycles, err := graph.DetectCycles()

            if tc.wantCycle {
                if err == nil {
                    t.Fatal("expected cycle error, got nil")
                }
                if len(cycles) == 0 {
                    t.Fatal("expected cycles to be found")
                }
                if !strings.Contains(cycles[0], tc.cyclePath) {
                    t.Errorf("cycle path:\ngot:  %s\nwant: %s", cycles[0], tc.cyclePath)
                }
            } else {
                if err != nil {
                    t.Fatalf("unexpected error: %v", err)
                }
            }
        })
    }
}
```

---

## Performance Optimization

```go
// For large graphs, use optimized visited tracking
type OptimizedResolver struct {
    visited   map[string]bool
    inPath    map[string]bool
    resolved  map[string][]Instruction // Memoize results
    baseDir   string
}

func (r *OptimizedResolver) Resolve(rulePath string) ([]Instruction, error) {
    absPath := filepath.Join(r.baseDir, rulePath)

    // Check memo cache
    if result, ok := r.resolved[absPath]; ok {
        return result, nil
    }

    // Check for cycle
    if r.inPath[absPath] {
        return nil, fmt.Errorf("cycle: %s", absPath)
    }

    // Skip if already processed
    if r.visited[absPath] {
        return r.resolved[absPath], nil
    }

    // Mark and recurse
    r.visited[absPath] = true
    r.inPath[absPath] = true
    defer func() { r.inPath[absPath] = false }()

    // ... resolution logic ...

    // Cache result
    r.resolved[absPath] = result
    return result, nil
}
```

---

## Common Pitfalls

❌ **Bad: Only using visited set (misses cycles)**

```go
// This misses cycles because visited prevents re-entry
func badResolve(node string, visited map[string]bool) {
    if visited[node] {
        return // Can't distinguish between cycle and already-processed
    }
    visited[node] = true
    // ...
}
```

✅ **Good: Use both visited and path stack**

```go
func goodResolve(node string, visited, inPath map[string]bool) error {
    if inPath[node] {
        return fmt.Errorf("cycle") // Currently being processed = cycle
    }
    if visited[node] {
        return nil // Already fully processed = OK
    }
    // ...
}
```

---

## Key Principles

1. **Use visited set to track processed nodes** - Avoid redundant work
2. **Use path stack to detect cycles** - Track current traversal path
3. **Clear path stack on backtrack** - Use defer for cleanup
4. **Build informative cycle messages** - Show full cycle path
5. **Resolve imports depth-first** - Dependencies before dependents
6. **Memoize results for performance** - Cache resolved imports

---

## References

- DFS algorithm: https://en.wikipedia.org/wiki/Depth-first_search
- Cycle detection: https://en.wikipedia.org/wiki/Cycle_(graph_theory)
- Topological sort: https://en.wikipedia.org/wiki/Topological_sorting
- Project requirements: `/Users/kydavis/Sites/agent-instruction/docs/plan/001-initial-buildout/technical-requirements.yaml` (lines 442-474)
