# Contributing to OpenCode Agents

First off, thank you for considering contributing to OpenCode Agents! It's people like you that make OpenCode Agents such a great tool.

This project is participating in a hackathon, and we welcome all contributions that improve the CLI, TUI, or agent configurations.

## Development Environment Setup

### Prerequisites
- [Go](https://go.dev/doc/install) 1.24 or later.
- A terminal with support for ANSI colors and sequences (most modern terminals work).

### Building and Running
To run the application directly from source:
```bash
go run ./cmd/opencode-agents --local
```

To build a local binary:
```bash
go build -o bin/opencode-agents ./cmd/opencode-agents
```

## How to Contribute

1.  **Fork the repository** and create your branch from `main`.
2.  **Make your changes.** If you're adding a new feature, please include tests.
3.  **Run Tests and Formatting.** Ensure your code adheres to the project's standards:
    ```bash
    go fmt ./...
    go test ./...
    ```
4.  **Submit a Pull Request.** Provide a clear description of your changes and why they are being made.

## Coding Style & Conventions

- **Formatting:** Use standard `go fmt`.
- **Error Handling:** Check errors explicitly. Use `fmt.Errorf("context: %w", err)` for wrapping.
- **TUI/CLI:** We use the [Charmbracelet](https://charm.sh/) stack (`bubbletea`, `huh`, `lipgloss`). Please stick to the Elm architecture for Bubble Tea components.
- **Naming:** Use `camelCase` for internal and `PascalCase` for exported members.

## Reporting Issues
- Use the [Bug Report](https://github.com/opencode/opencode-agents/issues/new?template=bug_report.md) template for bugs.
- Use the [Feature Request](https://github.com/opencode/opencode-agents/issues/new?template=feature_request.md) template for suggestions.

Thank you for your support!
