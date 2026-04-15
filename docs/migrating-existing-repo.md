# Migrating an Existing Repo to agent-instruction

This guide walks through migrating a repository that already has hand-maintained `CLAUDE.md` or `AGENTS.md` files to use agent-instruction for centralized rule management.

## Overview

If your repo has existing AI instruction files, migrating means:

1. Analyzing your current content and grouping it into logical rule files
2. Running `agent-instruction init` to set up the `.agent-instruction/` directory
3. Creating rule JSON files from your existing content
4. Running `agent-instruction build` to regenerate the instruction files
5. Verifying the output matches your intent

The fastest way to do step 1 is to give an AI agent the analysis prompt below.

## The Migration Prompt

Copy this prompt and give it to an AI agent (e.g. Claude) in your repository. It will read your existing instruction files and output the JSON configuration you need.

---

```
I'm migrating this repository to use agent-instruction, a tool that manages CLAUDE.md and AGENTS.md
files from centralized JSON rule files stored in .agent-instruction/rules/.

Please analyze the existing AI instruction files in this repository (CLAUDE.md, AGENTS.md, or any
.claude/CLAUDE.md) and help me plan the migration.

## Step 1: Discover and read existing instruction files

Search for and read all of the following files if they exist:
- CLAUDE.md (repo root)
- AGENTS.md (repo root)
- .claude/CLAUDE.md
- Any CLAUDE.md or AGENTS.md files in subdirectories/packages

## Step 2: Analyze the content

Identify all distinct instructions and group them into logical categories. Common categories include:
- global (cross-cutting concerns, code style, general behavior)
- testing (test patterns, coverage requirements, test tooling)
- security (auth, input validation, secrets handling)
- Language-specific: golang, typescript, python, etc.
- Domain-specific: api-design, database, deployment, etc.

For a monorepo, also note which instructions are:
- Global (should apply to all packages)
- Package-specific (only relevant to one service or package)

## Step 3: Output the migration plan

For each rule file you recommend creating, output the complete JSON content.

### config.json

Output the recommended `.agent-instruction/config.json`. Set:
- "frameworks": include "claude" if CLAUDE.md exists, "agents" if AGENTS.md exists, or both
- "packages": use "auto" if this is a monorepo with subdirectory packages, or [] for single-package repos

Example format:
{
  "version": "1.0",
  "frameworks": ["claude"],
  "packages": []
}

### Rule files

For each rule category identified, output the complete JSON for
`.agent-instruction/rules/<category>.json`.

Each rule file follows this format:
{
  "title": "Human-readable title",
  "imports": ["other-rule-file"],   // optional, omit if not needed
  "instructions": [
    {
      "heading": "Section Heading",  // optional but recommended
      "rule": "The instruction text, preserved verbatim from the original"
    }
  ]
}

Guidelines:
- Preserve the original instruction text as closely as possible — don't rewrite or summarize
- global.json should have no imports (it is the base layer)
- Other rule files can import global if they extend it, but avoid circular imports
- Split instructions at natural category boundaries, not arbitrarily
- If an instruction spans multiple concerns, put it in the more specific file

### Package configs (monorepo only)

Before writing package configs, compare the instruction files across all packages to identify
shared rules — instructions that appear in two or more packages, even if worded slightly
differently. Shared rules belong in `.agent-instruction/rules/`, not repeated in each package.

For each instruction you encounter, classify it as one of:
- **Shared** — applies to all (or most) packages; goes in a repo-level rule file
- **Package-specific** — unique to one package; stays in that package's `agent-instruction.json`

Output a shared/package breakdown before writing the JSON:

  Shared rules → .agent-instruction/rules/global.json (or a named rule file):
    - "<instruction summary>"
    - "<instruction summary>"

  Package-specific:
    packages/api  → packages/api/agent-instruction.json
      - "<instruction summary>"
    packages/web  → packages/web/agent-instruction.json
      - "<instruction summary>"

Then, for each package that has package-specific instructions, output the complete JSON for
`<package-path>/agent-instruction.json`.

Each package config follows the same format as rule files and should:
- Import only the repo-level rule files it actually needs
- Contain only instructions unique to that package — never repeat rules already covered by imports

## Step 4: Summarize the migration steps

List the exact shell commands to run after creating the files, e.g.:

  agent-instruction init --non-interactive
  # (then write each JSON file)
  agent-instruction build
  git diff CLAUDE.md  # verify output matches original intent

Note any content from the original files that doesn't cleanly map to a rule (e.g. project
overviews, links, or prose sections) — the user may want to keep that as custom content
outside the managed section markers.
```

---

## After Running the Prompt

Once the agent produces the JSON output:

### 1. Initialize agent-instruction

```bash
agent-instruction init --non-interactive
```

This creates `.agent-instruction/config.json` and `.agent-instruction/rules/global.json` with placeholder content. You'll overwrite these with the agent's output.

### 2. Write the config and rule files

Replace the generated placeholders with the JSON the agent produced. For example:

```bash
# Overwrite config
vim .agent-instruction/config.json

# Overwrite global rules
vim .agent-instruction/rules/global.json

# Create additional rule files
vim .agent-instruction/rules/testing.json
vim .agent-instruction/rules/security.json
```

For monorepos, write any package configs too:

```bash
vim packages/api/agent-instruction.json
```

### 3. Build and verify

```bash
agent-instruction build --verbose
```

Compare the generated output to your original files:

```bash
diff CLAUDE.md.backup CLAUDE.md
```

`agent-instruction init` creates a `.backup` of any existing instruction files before overwriting them.

### 4. Preserve custom content

Any content you want to keep outside the managed section (e.g. a project overview header, links, or prose) should live **outside** the `<!-- BEGIN AGENT-INSTRUCTION -->` / `<!-- END AGENT-INSTRUCTION -->` markers. Edit the generated file directly — that content will be preserved on future rebuilds.

### 5. Commit

```bash
git add .agent-instruction/
git add CLAUDE.md AGENTS.md
# For monorepos:
git add **/agent-instruction.json **/CLAUDE.md **/AGENTS.md
git commit -m "Migrate to agent-instruction"
```

## Common Migration Patterns

### Flat CLAUDE.md with mixed content

A single file with 30–100 lines of mixed instructions maps well to 2–4 rule files:

```
global.json       ← code style, general behavior, commit conventions
testing.json      ← test patterns, coverage
security.json     ← auth, validation, secrets
```

### CLAUDE.md with a project overview header

Keep the overview as custom content above the managed section. The overview doesn't need to be a rule — it won't change frequently and doesn't need to be inherited by packages.

### Monorepo with one root CLAUDE.md duplicated per package

This is the primary pain point agent-instruction solves. Global rules go in `.agent-instruction/rules/global.json`. Package-specific additions go in each `<package>/agent-instruction.json`. After migrating, you'll delete the duplicated files and let the build generate them.

### Instructions written for a specific agent framework

If your `CLAUDE.md` contains Claude-specific behavioral instructions and your `AGENTS.md` contains different content for other agents, keep them as separate instruction sets and use the `frameworks` field to generate the right file for each. If the content is identical (most common case), a single rule set generates both.

## See Also

- [Configuration Reference](configuration.md)
- [Examples](examples.md)
- [init command](commands/init.md)
- [build command](commands/build.md)
