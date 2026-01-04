# ntmux

A powerful, declarative tmux session manager that simplifies the creation and management of complex tmux layouts through JSON or YAML configuration files.

## Overview

**ntmux** is a Go-based CLI tool that acts as an intelligent wrapper around tmux, allowing you to define your terminal workspace declaratively. Instead of manually creating sessions and windows with multiple commands, define your entire layout in a configuration file and apply it with a single command.

> **⚠️ Development Status**: This project is currently in active development. APIs, commands, and configuration formats may change in future releases.

### Key Features

- **Declarative Configuration**: Define sessions, windows, and startup commands in JSON or YAML
- **Idempotent Operations**: Safely re-run configurations without duplicating sessions
- **Template System**: Create reusable templates for different projects or workflows
- **Cross-Platform Support**: Works on macOS, Linux, and Windows (with appropriate shells)
- **Single Command Execution**: Optimized batching of tmux commands for better performance
- **JSON Schema Validation**: IDE autocomplete and validation support via JSON schema
- **Zero Configuration Start**: Works with sensible defaults out of the box
- **tmux Compatibility**: Full pass-through support for native tmux commands

## Requirements

- **Go 1.23.4+** (for building from source)
- **tmux** (must be installed and available in PATH)
- Unix-like environment (macOS, Linux, WSL, or Windows with PowerShell)

## Installation

### Using Go Install

```bash
go install github.com/coeeter/ntmux@latest
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/coeeter/ntmux.git
cd ntmux

# Install using Make
make install

# Or build to ./tmp/ntmux
make build
```

### Verify Installation

```bash
ntmux --help
```

## Quick Start

### 1. Create a Configuration File

Create a `ntmux.json` file in your project directory:

```json
{
  "$schema": "https://raw.githubusercontent.com/coeeter/ntmux/main/schema.json",
  "sessions": [
    {
      "name": "my-project",
      "dir": ".",
      "default": true,
      "windows": [
        {
          "name": "editor",
          "cmd": "nvim .",
          "default": true
        },
        {
          "name": "terminal"
        },
        {
          "name": "server",
          "cmd": "npm run dev"
        }
      ]
    }
  ]
}
```

### 2. Apply the Configuration

```bash
# If ntmux.json exists in current directory
ntmux

# Or explicitly specify the file
ntmux apply ntmux.json
```

### 3. You're Done!

ntmux will create your tmux session with all windows configured and attach you to the default session.

## Configuration

### File Formats

ntmux supports both JSON and YAML formats:

- `ntmux.json` (preferred)
- `ntmux.yaml` or `ntmux.yml`

### Configuration Schema

```json
{
  "$schema": "https://raw.githubusercontent.com/coeeter/ntmux/main/schema.json",
  "sessions": [
    {
      "name": "session-name", // Required: Session identifier
      "dir": "/path/to/directory", // Optional: Working directory (default: current directory)
      "default": true, // Optional: Attach to this session on startup
      "windows": [
        {
          "name": "window-name", // Required: Window identifier
          "dir": "/path/to/directory", // Optional: Window-specific directory
          "cmd": "command to run", // Optional: Command to execute when window opens
          "default": true // Optional: Select this window in the session
        }
      ]
    }
  ]
}
```

### Configuration Fields

#### Session Fields

| Field     | Type    | Required | Description                                       |
| --------- | ------- | -------- | ------------------------------------------------- |
| `name`    | string  | Yes      | Unique identifier for the tmux session            |
| `dir`     | string  | No       | Working directory (defaults to current directory) |
| `default` | boolean | No       | If true, attach to this session on startup        |
| `windows` | array   | Yes      | Array of window configurations                    |

#### Window Fields

| Field     | Type    | Required | Description                                     |
| --------- | ------- | -------- | ----------------------------------------------- |
| `name`    | string  | Yes      | Name of the window                              |
| `dir`     | string  | No       | Window-specific working directory               |
| `cmd`     | string  | No       | Command to execute when window opens            |
| `default` | boolean | No       | If true, select this window when session starts |

### Default Behavior

- If no session has `default: true`, the first session becomes the default
- If no window has `default: true`, the first window becomes the default
- Relative paths are resolved relative to the current working directory

## Usage

### Commands

#### Default Behavior

```bash
# If ntmux.json or ntmux.yaml exists in current directory
ntmux
```

When run without arguments, ntmux automatically searches for and applies `ntmux.json` or `ntmux.yaml` in the current directory.

#### Apply Configuration

```bash
# Apply default configuration file
ntmux apply

# Apply specific configuration file
ntmux apply path/to/template.json
ntmux apply path/to/template.yaml
```

#### Generate New Template

```bash
# Generate JSON template (default)
ntmux new-template

# Generate YAML template
ntmux new-template --format yaml
ntmux new-template -f yaml
```

The `new-template` command will:

1. Check for custom templates in `~/.config/ntmux/template.{json,yaml}`
2. Use custom template if found, otherwise use built-in defaults
3. Create the template file in the current directory

#### Pass-Through to tmux

```bash
# Any unrecognized command is passed to tmux
ntmux list-sessions
ntmux attach -t my-session
ntmux kill-session -t old-session
```

#### Help

```bash
# Show combined ntmux and tmux help
ntmux --help
ntmux -h
```

## Examples

### Basic Project Setup

```json
{
  "$schema": "https://raw.githubusercontent.com/coeeter/ntmux/main/schema.json",
  "sessions": [
    {
      "name": "my-app",
      "windows": [
        {
          "name": "editor",
          "cmd": "nvim .",
          "default": true
        },
        {
          "name": "terminal"
        }
      ]
    }
  ]
}
```

### Full-Stack Development Environment

```json
{
  "$schema": "https://raw.githubusercontent.com/coeeter/ntmux/main/schema.json",
  "sessions": [
    {
      "name": "fullstack",
      "default": true,
      "windows": [
        {
          "name": "backend",
          "dir": "./server",
          "cmd": "npm run dev",
          "default": true
        },
        {
          "name": "frontend",
          "dir": "./client",
          "cmd": "npm start"
        },
        {
          "name": "database",
          "cmd": "docker-compose up postgres"
        },
        {
          "name": "logs",
          "cmd": "tail -f ./logs/app.log"
        }
      ]
    }
  ]
}
```

### Multi-Project Workspace

```json
{
  "$schema": "https://raw.githubusercontent.com/coeeter/ntmux/main/schema.json",
  "sessions": [
    {
      "name": "api",
      "dir": "~/projects/api",
      "default": true,
      "windows": [
        {
          "name": "code",
          "cmd": "nvim .",
          "default": true
        },
        {
          "name": "server",
          "cmd": "go run main.go"
        }
      ]
    },
    {
      "name": "frontend",
      "dir": "~/projects/web",
      "windows": [
        {
          "name": "code",
          "cmd": "code ."
        },
        {
          "name": "dev",
          "cmd": "npm run dev"
        }
      ]
    }
  ]
}
```

### YAML Configuration Example

```yaml
sessions:
  - name: my-project
    dir: ~/projects/awesome-app
    default: true
    windows:
      - name: editor
        cmd: nvim .
        default: true
      - name: terminal
      - name: server
        cmd: npm run dev
```

## Advanced Usage

### Custom Template Location

Create a custom template that will be used by `ntmux new-template`:

```bash
# Create config directory
mkdir -p ~/.config/ntmux

# Create your custom template
cat > ~/.config/ntmux/template.json << 'EOF'
{
  "$schema": "https://raw.githubusercontent.com/coeeter/ntmux/main/schema.json",
  "sessions": [
    {
      "name": "project",
      "windows": [
        {
          "name": "nvim",
          "cmd": "nvim .",
          "default": true
        },
        {
          "name": "shell"
        },
        {
          "name": "git",
          "cmd": "git status"
        }
      ]
    }
  ]
}
EOF
```

Now `ntmux new-template` will use your custom template.

### Idempotent Sessions

ntmux checks if sessions already exist before creating them:

```bash
# First run: creates sessions
ntmux apply

# Second run: skips existing sessions (no duplicates)
ntmux apply
```

### Shell-Specific Commands

ntmux automatically detects your shell and formats commands appropriately:

- **Unix shells**: Commands wrapped with `shell -c 'cmd; exec shell'`
- **PowerShell**: Commands wrapped with `pwsh -NoExit -Command "& {cmd}"`
- **cmd.exe**: Commands wrapped with `cmd /K "cmd"`

### JSON Schema Benefits

The `$schema` field enables:

- IDE autocomplete for configuration fields
- Real-time validation in editors like VS Code
- Inline documentation for field descriptions

## Development

### Available Make Targets

```bash
make help              # Show available targets
make build             # Build binary to ./tmp/ntmux
make install           # Install to $GOPATH/bin/ntmux
make uninstall         # Remove from $GOPATH/bin
make test              # Run tests
make test-verbose      # Run tests with verbose output
make test-cover        # Generate coverage report
make clean             # Remove build artifacts
make run ARGS="args"   # Run with arguments
make fmt               # Format code with go fmt
make vet               # Run go vet
make deps              # Download dependencies
make tidy              # Tidy go.mod and go.sum
make check             # Run fmt, vet, and tests
make generate-schema   # Regenerate schema.json
```

## Contributing

Contributions are welcome! Here's how you can help:

### Reporting Issues

1. Check existing issues to avoid duplicates
2. Provide clear description of the problem
3. Include steps to reproduce
4. Share your configuration file (if applicable)
5. Include tmux and ntmux versions

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run quality checks (`make check`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Setup

```bash
# Clone your fork
git clone https://github.com/your-username/ntmux.git
cd ntmux

# Install dependencies
make deps

# Build and test
make build
make test

# Install locally for testing
make install
```

### Code Guidelines

- Follow Go best practices and idioms
- Add tests for new features
- Update documentation for user-facing changes
- Run `make check` before committing
- Keep commits focused and atomic

## Troubleshooting

### tmux not found

**Error**: `exec: "tmux": executable file not found in $PATH`

**Solution**: Install tmux using your package manager:

```bash
# macOS
brew install tmux

# Ubuntu/Debian
sudo apt-get install tmux

# Fedora
sudo dnf install tmux

# Arch Linux
sudo pacman -S tmux
```

### Session already exists

ntmux automatically skips existing sessions. To recreate a session:

```bash
# Kill the existing session first
tmux kill-session -t session-name

# Then apply your configuration
ntmux apply
```

### Commands not executing

Ensure your commands are properly quoted in JSON:

```json
{
  "cmd": "echo 'hello world'" // Correct
}
```

For complex commands, use shell scripts:

```json
{
  "cmd": "./scripts/setup.sh"
}
```

### IDE not recognizing schema

Ensure the `$schema` field is at the top of your configuration:

```json
{
  "$schema": "https://raw.githubusercontent.com/coeeter/ntmux/main/schema.json",
  "sessions": [...]
}
```

VS Code should automatically fetch the schema. For other editors, check their JSON schema configuration.

## License

This project is open source. Please check the repository for license details.

## Author

Created by [Coeeter](https://github.com/coeeter)

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [invopop/yaml](https://github.com/invopop/yaml) for YAML support
- JSON Schema generation via [invopop/jsonschema](https://github.com/invopop/jsonschema)

---
