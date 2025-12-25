# mcp-local-context

A simple MCP (Model Context Protocol) server that provides prompts to AI assistants. This server helps ensure AI coding assistants have the right context when working with third-party packages by leveraging local module caches and custom prompts.

## Overview

`mcp-local-context` is an MCP server that provides prompts to AI assistants like Cursor, Claude, and GitHub Copilot. The primary use case is to provide systematic approaches for working with third-party packages by referencing local caches (like the Go module cache) rather than relying on potentially outdated documentation or assumptions.

## Why Use This?

When working with multiple AI tools, it's easier to configure one MCP server that provides your custom prompts everywhere, rather than adding the same prompts to each tool individually and keeping them synchronized. This centralizes your AI assistant context prompts and makes them reusable across different tools.

## Features

- **Golang Context Prompt**: Built-in prompt for working with third-party Go packages using the Go module cache
- **Custom Prompts**: Auto-discovery of custom prompt files from `~/.mcp-local-context/prompts/*.md`
- **Cross-platform**: Works on macOS, Linux, and Windows
- **Configurable**: Simple JSON configuration file
- **Extensible**: Easy to add new prompt providers (e.g., JavaScript, Python)

## Installation

```bash
go install github.com/svetlyi/mcp-local-context@latest
```

## Configuration

### Configuration File

Create a configuration file at `~/.mcp-local-context/config.json`:

```json
{
  "log_level": "info",
  "log_file": "~/mcp-local-context.log",
  "custom_prompt_dirs": ["~/custom-prompts", "/path/to/other/prompts"]
}
```

**Configuration Options**:
- `log_level`: Logging level (`debug`, `info`, `warn`, `error`). Default: `info`
- `log_file`: Path to log file (supports `~/` expansion). If not set, logs to a temporary file, depending on the OS.
- `custom_prompt_dirs`: Array of directories containing custom prompt files. The default `~/.mcp-local-context/prompts/` is always included

### Custom Prompts

Place custom prompt files (Markdown format) in `~/.mcp-local-context/prompts/`. Each `.md` file will be automatically discovered and made available as a prompt.

> **Important**: The first line of the file will be used as the prompt's description. If the first line is a markdown heading (starting with `#`), the heading markers will be automatically removed.

Example: `~/.mcp-local-context/prompts/my-custom-prompt.md`

```markdown
A description of my custom prompt.

# My Custom Prompt

This is my custom prompt that will be provided to AI assistants.
```

The prompt will be available as a prompt named `my-custom-prompt` with the description "A description of my custom prompt." (extracted from the first line).

## Usage

### Running the Server

The server communicates via stdio (standard input/output), which is the standard transport for MCP servers:

```bash
./bin/mcp-local-context
```

### Integration with AI Tools

For example, for Cursor IDE, add it to the settings:

```json
{
  "mcpServers": {
    "local-context": {
      "command": "/path/to/mcp-local-context/bin//"
    }
  }
}
```

If you installed it using `go install`, you can find the binary in your GO binary path, `echo $(go env GOPATH)/bin/mcp-local-context`.

## Available Prompts

### golang-context-rule

Provides a systematic approach for working with third-party Go packages by referencing the Go module cache. This prompt guides AI assistants to:

1. Identify the exact module version from `go.mod`
2. Locate the Go module cache
3. Explore the package structure
4. Use `go doc` to get documentation
5. Read the source code directly

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Related

- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
