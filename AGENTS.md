# OpenCode Agents Repository Guidelines

Welcome to the `opencode-agents` codebase! This document provides instructions and guidelines for agentic coding agents (like yourself) operating within this repository. 

This is a Go-based CLI and Terminal User Interface (TUI) application.

## 1. Build, Lint, and Test Commands

### Running the Application
To run the CLI application locally from source:
```bash
go run ./cmd/opencode-agents [flags]
```
Example: `go run ./cmd/opencode-agents --local` to sync agents locally.

### Building the Application
To build the binary into a `bin` directory:
```bash
go build -o bin/opencode-agents ./cmd/opencode-agents
```

### Testing
To run all tests in the repository:
```bash
go test ./...
```
To run tests with verbosity, race detection, and coverage:
```bash
go test -v -race -cover ./...
```

**Running a Single Test:**
To run a specific test (e.g., `TestLoadAgents` in the `manager` package):
```bash
go test -run ^TestLoadAgents$ ./internal/manager -v
```
To run a specific subtest:
```bash
go test -run ^TestLoadAgents$/^MySubtest$ ./internal/manager -v
```

### Linting and Formatting
To format the code (always run this before writing or committing code):
```bash
go fmt ./...
```
To tidy up module dependencies:
```bash
go mod tidy
```
If `golangci-lint` is available in the environment, you can run it via:
```bash
golangci-lint run
```

## 2. Code Style and Conventions

### Project Structure
- `cmd/opencode-agents/`: Contains the main application entry point (`main.go`). Keep this lightweight.
- `internal/`: Contains private application logic (`cli`, `config`, `manager`). Code here cannot be imported by other repositories.
- `agents/`: A directory containing the markdown files used for agent configurations.
- `go.mod` / `go.sum`: Dependency tracking. The module name is `github.com/opencode/opencode-agents`.

### Go Code Style
- **Formatting:** Code MUST be formatted with standard Go tooling (`go fmt`).
- **Naming Conventions:**
  - Use `camelCase` for unexported variables, constants, and functions.
  - Use `PascalCase` for exported variables, functions, and types (Structs, Interfaces).
  - Variable names should be concise but descriptive. Examples: `err` for errors, `req` for requests, `m` for a method receiver of a `Manager` type.
  - Interfaces should generally end in `-er` if they define an action (e.g., `Reader`, `Writer`).

### Error Handling
- **Explicit Checks:** Always check for errors immediately. Use the standard idiom: `if err != nil { return err }`.
- **Error Wrapping:** When returning errors up the stack, wrap them with contextual information using `fmt.Errorf("failed to do X: %w", err)`. This allows callers to use `errors.Is` or `errors.As`.
- **Fail Fast in Main:** In `main.go`, errors should be printed to `os.Stderr`, followed by `os.Exit(1)`. Inside packages, always return the error to the caller.

### Imports
Group imports into logically separated blocks (handled automatically by `goimports`):
1. Standard library packages (e.g., `fmt`, `os`, `path/filepath`).
2. Third-party packages (e.g., `github.com/charmbracelet/huh`, `github.com/adrg/frontmatter`).
3. Internal project packages (e.g., `github.com/opencode/opencode-agents/internal/manager`).

### TUI and CLI Libraries
This project relies heavily on [Charmbracelet](https://charm.sh/) libraries:
- **`huh`**: Used for building forms and interactive prompts.
- **`bubbletea`**: Used for more complex, stateful terminal user interfaces. When using Bubble Tea, adhere strictly to the Elm architecture (`Model`, `Init`, `Update`, `View`).
- **`lipgloss`**: Used for styling terminal output (colors, borders, layouts).
- **Standard `flag` package**: Used in `main.go` for basic command-line argument parsing.

### File Parsing
- The project reads Markdown files with YAML frontmatter to configure agents.
- We use `github.com/adrg/frontmatter` for parsing. Handle metadata structures defensively, checking types (e.g., `if val, ok := metadata["name"].(string); ok { ... }`).

## 3. General Agent Behavior

- **Read Before Write:** Always use the `Read` or `Glob` tool to examine a file or directory structure before modifying it. Do not assume file contents.
- **Idiomatic Edits:** When adding new features or fixing bugs, mimic the existing style and architecture of the neighboring code.
- **Small, Focused Changes:** Make targeted edits rather than rewriting entire files, unless an entire file rewrite is explicitly requested by the user.
- **Verification Loop:** After making Go code modifications, **you must run `go build ./...`** to ensure there are no compilation errors. If tests exist for the package you modified, run them.
- **Do not commit changes** unless explicitly asked to do so by the user.
- **Keep it Simple:** Prefer standard library solutions over adding new third-party dependencies unless absolutely necessary.
