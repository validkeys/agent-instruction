## Agent Instruction

**Problem** It's very difficult to manage the CLAUDE.md/AGENTS.md files in a monorepo.

**IDEA**:

A command line utility that manages all of them for you.

**User Story**

- I open my monorepo and run agent-instruction init. This creates a .agent-instruction folder in the root of my current repo. This folder contains a rules/global.json. The init command asks me if I want to clear my current CLAUDE.md/AGENTS.md file or not. If I clear it, it wipes out the current content and adds two comment markers that it will later use to insert content between. The .agent-instruction folder probably also needs a config.json to determine what AI frameworks we are supports (whether to create CLAUDE.md or AGENTS.md)
- When I run agent-instruction build, the content between the comment markers in my root claude.md / agents.md file is updated with content generated from the array of rules in my global.json
- Often times as you're going through a process, you realize a rule that AI keeps messing up on that you forget to add to our claude.md file. There should be a command called agent-instruction add. This should then find all the files in the rules folder in the root .agent-instruction folder, and ask you which you want to add to. It could have agent-instruction list to list the rule files you have, then use agent-instruction add "Rule Content" --title "Optional Title for the rule" --rule="optional rule file name to add the rule to"
- This should have a skill that let's the AI know how to use agent-instruction
- Now in my different monorepo folders, instead of having claude.md and agents.md files, I should have agent-instruction.json at the root of the package. It would have a similar json as below where it could import one or more of our project-level rule files, but could also define it's own rule file. Then agent-instruction CLI command should also support adding a rule to a project's agent-instruction file.
- The "build" command should then parse the global agent-instruction config as well as all agent-browser.json files found in the repo, then generate the CLAUDE.md files and AGENTS.md files in the various folders in the repo

**JSON**

```json
{
  title: "Global agent instructions",
  instructions: [
    {
  heading: "Optional heading"
  rule: "Some markdown denoting a rule.",
  references: [
    { title: "", path: "" }
  ] // optional array of relative file references to be included
    }
  ],
  imports: ["./style-guide.json"] // if we want to also load the content from other agent-instruction files into this file
}
```

```

```
