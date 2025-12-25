# mcp-local-context

A simple MCP (Model Context Protocol) server that provides prompts to AI assistants. This server helps ensure AI coding assistants have the right context when working with third-party packages by leveraging local module caches and custom rules.

## Overview

`mcp-local-context` is an MCP server that provides prompts to AI assistants like Cursor, Claude, and GitHub Copilot. The primary use case is to provide systematic approaches for working with third-party packages by referencing local caches (like the Go module cache) rather than relying on potentially outdated documentation or assumptions.

## Why Use This?

When working with multiple AI tools, it's easier to configure one MCP server that provides your custom rules and prompts everywhere, rather than adding the same rules to each tool individually and keeping them synchronized. This centralizes your AI assistant context rules and makes them reusable across different tools.

## Features

- **Golang Context Rule**: Built-in prompt for working with third-party Go packages using the Go module cache
- **Custom Rules**: Auto-discovery of custom rule files from `~/.mcp-local-context/rules/*.md`
- **Cross-platform**: Works on macOS, Linux, and Windows
- **Configurable**: Simple JSON configuration file
- **Extensible**: Easy to add new prompt providers (e.g., JavaScript, Python)

## Installation

1. Clone this repository:
```bash
git clone https://github.com/svetlyi/mcp-local-context.git
cd mcp-local-context
```

2. Build the server:
```bash
make build
```

Or manually:
```bash
go build -o bin/mcp-local-context ./cmd/server
```

## Configuration

### Configuration File

Create a configuration file at `~/.mcp-local-context/config.json`:

```json
{
  "address": "localhost",
  "port": 8080,
  "log_level": "info"
}
```

**Note**: Currently, the server uses stdio transport (standard for MCP servers), so the address and port settings are reserved for future HTTP transport support.

### Custom Rules

Place custom rule files (Markdown format) in `~/.mcp-local-context/rules/`. Each `.md` file will be automatically discovered and made available as a prompt.

Example: `~/.mcp-local-context/rules/my-custom-rule.md`

```markdown
# My Custom Rule

This is my custom rule that will be provided to AI assistants.
```

The rule will be available as a prompt named `my-custom-rule`.

## Usage

### Running the Server

The server communicates via stdio (standard input/output), which is the standard transport for MCP servers:

```bash
./bin/mcp-local-context
```

### Integration with AI Tools

#### Cursor

Add to your Cursor settings:

```json
{
  "mcpServers": {
    "local-context": {
      "command": "/path/to/mcp-local-context/bin/mcp-local-context"
    }
  }
}
```

#### Claude Desktop

Add to your Claude Desktop configuration:

```json
{
  "mcpServers": {
    "local-context": {
      "command": "/path/to/mcp-local-context/bin/mcp-local-context"
    }
  }
}
```

## Available Prompts

### golang-context-rule

Provides a systematic approach for working with third-party Go packages by referencing the Go module cache. This prompt guides AI assistants to:

1. Identify the exact module version from `go.mod`
2. Locate the Go module cache
3. Explore the package structure
4. Use `go doc` to get documentation
5. Read the source code directly

See [adr/golang-context-rule.md](adr/golang-context-rule.md) for the full content.

## Development

### Project Structure

```
mcp-local-context/
├── cmd/
│   └── server/
│       └── main.go              # MCP server entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading
│   ├── prompts/
│   │   ├── provider.go          # Prompt provider interface
│   │   └── golang.go            # Golang prompt implementation
│   └── rules/
│       └── loader.go            # Rule discovery and loading
├── .github/
│   └── workflows/
│       └── ci.yml               # GitHub Actions pipeline
├── Makefile                     # Build, test, lint commands
└── README.md                    # This file
```

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Linting

```bash
make lint
```

### Running Locally

```bash
make run
```

### Cleaning

```bash
make clean
```

## MCP Protocol

This server implements the Model Context Protocol (MCP) using JSON-RPC 2.0 over stdio. It supports the following methods:

- `initialize`: Initialize the MCP connection
- `prompts/list`: List all available prompts
- `prompts/get`: Get a specific prompt by name

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Related

- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- Blog post: [Improve AI Context: Use Your Go Module Cache](adr/golang-context-rule.md)

