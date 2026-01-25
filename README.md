# mcp-local-context

## What is it and why use this?

An MCP server that helps AI assistants work with GoLang third-party packages **using your local Go module cache**. Think of it as a **free, local alternative to [Context7](https://context7.com)**.

- **Free and unlimited** - uses your existing GoLang cache, no API limits or costs
- **Works with any module** - supports public and private modules already in your cache
- **Universal compatibility** - works with all AI tools (Cursor, Claude, GitHub Copilot, etc.)
- **Customizable** - add your own prompts and instructions for AI agents

## Quick Start

```bash
# Install
go install github.com/svetlyi/mcp-local-context@latest

# Find binary path
echo $(go env GOPATH)/bin/mcp-local-context
```

Add to your IDE (e.g., Cursor): Settings → MCP Servers → Add:

```json
{
  "mcpServers": {
    "local-context": {
      "command": "/path-to-the-binary/bin/mcp-local-context"
    }
  }
}
```

> The server communicates via stdio (standard input/output).

## Features

- **Built-in & Custom Prompts**: Built-in Go context prompt with automatic discovery of custom prompt files from `~/.mcp-local-context/prompts/*.md`
- **Configurable & Extensible**: JSON configuration file and easily add new prompt providers (e.g., JavaScript, Python)

## Configuration (optional)

Create `~/.mcp-local-context/config.json` (optional):

```json
{
  "log_level": "info",
  "log_file": "~/mcp-local-context.log",
  "custom_prompt_dirs": ["~/custom-prompts", "/path/to/other/prompts"]
}
```

**Options**:
- `log_level`: `debug`, `info`, `warn`, `error` (default: `info`)
- `log_file`: Path to log file (supports `~/` expansion). Default: OS temp file
- `custom_prompt_dirs`: Additional directories for custom prompts. The `~/.mcp-local-context/prompts/` directory is always included

### Custom Prompts

Place Markdown files in `~/.mcp-local-context/prompts/`. Each `.md` file is automatically discovered and loaded.

#### Prompt Configuration

You can configure prompts using key-value pairs at the start of the file. The configuration section ends with a blank line.

**Supported configuration keys**:
- `title`: Custom title/description for the prompt
- `lang`: Language identifier (e.g., `go`, `javascript`, `python`)

**Example with configuration**: `~/.mcp-local-context/prompts/my-custom-prompt.md`

```markdown
title: My Custom Prompt Title
lang: go

# My Custom Prompt

This is my custom prompt that will be provided to AI assistants.
```

**Example without configuration**: `~/.mcp-local-context/prompts/simple-prompt.md`

```markdown
A description of my custom prompt.

# My Custom Prompt

This is my custom prompt that will be provided to AI assistants.
```

> **Note**: If `title` is not specified, the first line is used as the prompt description. If the first line is a markdown heading (`#`), the heading markers are automatically removed.

The prompt will be available as a prompt named `my-custom-prompt` (derived from the filename) with the configured title or extracted description.

## Usage

### How to use

You can use the MCP in two ways:

**Method 1: Directly reference a specific prompt**

Manually invoke a specific prompt by name. This gives you precise control over which instructions are used.

For example, in Cursor, you can write:

```prompt
Create something in Go, use the /local_context/golang-context-rule prompt
```

**Method 2: Ask the AI to use the MCP tools**

Ask the AI to use the MCP's automatic language detection tools. The AI will call `list_supported_languages` and then `get_context_instructions` with the appropriate language, which automatically provides the right context instructions.

For example:

```prompt
Create something in Go, use the local-context MCP to get context instructions for Go
```

or more generally:

```prompt
When working with third-party packages, use the local-context MCP to get the appropriate context instructions
```

## Available Prompts

### golang-context-rule

Systematic approach for working with third-party Go packages using the Go module cache. Guides AI assistants to:

1. Identify the exact module version from `go.mod`
2. Locate the Go module cache
3. Explore the package structure
4. Use `go doc` for documentation
5. Read source code directly

## Future work

* **Module indexing**: Index existing modules to find appropriate functionality more efficiently based on the task
* **More languages**: Add support for other languages, such as JavaScript, etc.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

🤝 Contributions are welcome! Please feel free to submit a Pull Request.

## Related

🔗 [Model Context Protocol Specification](https://modelcontextprotocol.io/)
