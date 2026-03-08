# OpenCode Agents

**Define, Version, and Sync your AI Personas.**

OpenCode Agents is a configuration manager that lets you define specialized AI personas as code, version control them with Git, and sync them instantly to your **OpenCode** environment.

## 💡 Why Personas?

Generic AI assistants are great, but complex engineering workflows often require specific roles and strict rules.

-   **Consistency:** Ensure your AI always follows your team's style guide, security protocols, or documentation standards.
-   **Specialization:** Configure **Reasoning** (Pro) models for complex architecture and **Speed** (Flash) models for quick syntax fixes.
-   **Workflow Control:** Create a chain of command—start with a **Planner** to scope the work, then switch to a **Builder** to execute it.

## 🤖 Workflow: Using Agents to Build Agents

The most powerful way to define a new persona is to use your existing AI toolkit to build it.

**Don't start from scratch.** Use **OpenCode** to scaffold your new agent definitions:

> "Create a new agent persona for **Security Auditing**. It should use a Reasoning model, be extremely paranoid about input validation, and strictly output findings in a markdown table."

Take the output, save it to your `agents/` directory, and you now have a repeatable, versioned expert for that task.

## ⚡️ Quick Start

### 1. Install

**Option A: Pre-compiled Binaries (Recommended)**

1. Download the latest release for your operating system and architecture from the **[Releases Page](https://github.com/DanHenton/opencode-agents/releases)**.
2. Extract the archive and place the `opencode-agents` binary into a directory that is included in your system's `PATH`.

> 💡 **What is the system PATH?**
> The `PATH` is an environment variable containing a list of directories that your operating system checks for executable files. Placing the binary in one of these directories (like `/usr/local/bin` on Linux/macOS) allows you to run the `opencode-agents` command from any folder in your terminal.

**Linux / macOS:**
```bash
# Extract the downloaded archive (replace filename with your specific version)
tar -xzf opencode-agents_Linux_x86_64.tar.gz

# Move the binary to a directory in your PATH
sudo mv opencode-agents /usr/local/bin/

# Verify the installation
opencode-agents --help
```

**Windows:**
1. Extract the downloaded `.zip` file.
2. Move the `opencode-agents.exe` binary to a permanent location (e.g., `C:\Program Files\OpenCodeAgents`).
3. Add that directory to your system's `Path` environment variable via the Windows Settings.
4. Open a new terminal and verify with `opencode-agents --help`.

**Option B: Build from Source**

If you have Go installed, you can build and install it directly:
```bash
go install github.com/opencode/opencode-agents/cmd/opencode-agents@latest
```

### 2. Run
```bash
opencode-agents
```
*(First run? We'll automatically generate a starter set of agents in `~/.config/opencode-agents/` for you.)*

### 3. Select
Use the interactive checklist to pick which agents you want active. This updates your local configuration (e.g., `opencode.json`), keeping your keys and server settings safe while swapping out system prompts.

## 📂 Configuration

Agents are defined in simple Markdown files. You can organize them however you like:


```text
~/.config/opencode-agents/
├── AGENTS.md           # Shared context & documentation
└── agents/
    ├── architect.md    # A "Reasoning" model for design
    ├── plan.md         # A read-only planning partner
    └── fix.md          # A "Speed" model for quick edits
```


```sh

cd ~/.config/opencode-agents/ && opencode .

```

Build a new agent persona, then import the agent into the target opencode.json file via:

```
opencode-agents
```

### Configuration Format

Each file uses YAML frontmatter for configuration and the body for the system prompt.

**Example: The Planner**
*Uses a **Reasoning** model to think before acting.*

```markdown
---
name: plan
description: Read-only pair programmer for planning
model: gemini-pro-1.5  # or o1, claude-3-opus, etc.
---
You are an expert Senior Software Engineer acting as a collaborative pair programmer.
**DEFAULT STATE:** PLAN MODE.

**CRITICAL CONSTRAINTS:**
- STRICTLY FORBIDDEN: Any file edits or modifications.
- NO CODE GENERATION: Do not generate functional code blocks.
- GOAL: Read files, understand the context, and produce a detailed Implementation Plan.
```

**Example: The Fixer**
*Uses a **Speed** model for instant results.*

```markdown
---
name: quick-fix
description: Fast syntax fixes and linting
model: gemini-flash-1.5 # or gpt-4o-mini, haiku, etc.
temperature: 0.1
---
You are a code cleaner. Fix syntax errors in the provided snippet.
Output ONLY the corrected code. No conversational filler.
```

## 🛠 Features

-   **Interactive TUI**: Easily toggle agents on/off with a terminal interface.
-   **GitOps Friendly**: Check your `~/.config/opencode-agents/` directory into Git to share personas with your team.
-   **Model Agnostic**: Configure any model ID supported by your backend in the frontmatter.

## License
MIT
