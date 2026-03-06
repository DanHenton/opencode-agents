# OpenCode Agents

[![CI](https://github.com/opencode/opencode-agents/actions/workflows/ci.yml/badge.svg)](https://github.com/opencode/opencode-agents/actions/workflows/ci.yml)

`opencode-agents` is a fast, robust Go CLI tool and Terminal UI (TUI) designed to manage and sync external GitOps-style Markdown agent prompts into a deeply nested `opencode.json` configuration file.

It allows you to maintain your agent prompts as clean, version-controllable Markdown files (with YAML frontmatter) and easily inject them into your active OpenCode configuration, either interactively via a sleek TUI or automatically via CLI flags.

## Key Features

*   **Interactive TUI:** Easily toggle which agents are enabled or disabled using a sleek terminal interface powered by [Charmbracelet's Huh](https://github.com/charmbracelet/huh).
*   **GitOps-Friendly:** Define agents using standard Markdown files with YAML frontmatter. Keep your prompts organized and version-controlled.
*   **Safe JSON Merging:** Intelligently mutates the `opencode.json` file. It updates the `"agent"` key while safely preserving all other unknown top-level and nested fields (e.g., `server`, `skills`).
*   **Automation Ready:** Powerful CLI flags to bypass the TUI entirely, perfect for CI/CD pipelines or power-user scripts.

## Installation

### From Binaries
You can download pre-compiled binaries for Linux, macOS, and Windows from the [GitHub Releases](https://github.com/opencode/opencode-agents/releases) page.

### From Source (Using Go)
If you have Go installed (1.21+), you can install the latest version directly:

```bash
go install github.com/opencode/opencode-agents/cmd/opencode-agents@latest
```

## Usage & CLI Flags

By default, running `opencode-agents` without any flags will launch the interactive TUI, allowing you to select which agents to enable or disable.

```bash
opencode-agents
```

*(Note: The TUI will automatically pre-select any agents that are already enabled in your `opencode.json` file.)*

### CLI Flags

*   `--global`: Update the global configuration file located at `~/.config/opencode/opencode.json`.
*   `--local`: Update the local configuration file located at `./opencode.json`.
*   `--sync-all`: Bypass the interactive TUI and automatically sync and enable all discovered agents.
*   `--dir <path>`: Override the default source directory for agent Markdown files. Defaults to `$OPENCODE_AGENTS_DIR` or `~/.config/opencode-agents/agents`.

*Note: You cannot use both `--global` and `--local` flags simultaneously.*

## Agent File Format

Agents are defined using standard Markdown files (`.md`). The file uses YAML frontmatter to define metadata (which gets injected into the JSON) and the body of the Markdown file acts as the agent's actual `"prompt"`.

**Example: `code-reviewer.md`**

```markdown
---
name: code-reviewer
description: An expert code reviewer
model: gpt-4-turbo
temperature: 0.2
---
You are an expert Senior Software Engineer.
Please review the provided code for potential bugs, security vulnerabilities, and adherence to best practices.
Provide constructive feedback and suggest improvements.
```

*   **`name`**: (Optional) Overrides the agent's name in the JSON. If omitted, the filename (without `.md`) is used.
*   **`description`**: (Optional) A brief description of the agent, shown in the TUI.
*   **Other metadata**: Any other YAML keys (like `model`, `temperature`) are safely injected into the agent's JSON object.
*   **Markdown Body**: Everything below the `---` block becomes the agent's `"prompt"`.

## Configuration Paths

*   **Source Directory (Markdown files)**:
    1.  The path provided via the `--dir` flag.
    2.  The `$OPENCODE_AGENTS_DIR` environment variable.
    3.  Fallback: `~/.config/opencode-agents/agents`.
*   **Target JSON (`opencode.json`)**:
    *   With `--local`: `./opencode.json`
    *   With `--global`: `~/.config/opencode/opencode.json`

## Development & Contributing

To contribute to `opencode-agents` or run it locally from source:

1. Clone the repository:
   ```bash
   git clone https://github.com/opencode/opencode-agents.git
   cd opencode-agents
   ```

2. Run the CLI directly:
   ```bash
   go run ./cmd/opencode-agents --help
   ```

3. Format and tidy code before committing:
   ```bash
   go fmt ./...
   go mod tidy
   ```

4. Run the test suite:
   ```bash
   go test -v -race -cover ./...
   ```
