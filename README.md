# GitHub MCP Server

A Model Context Protocol (MCP) server for comprehensive GitHub CLI integration. This server implements all major GitHub CLI functionality, allowing AI assistants to interact with GitHub repositories, issues, pull requests, files, and more.

## Features

- **Full GitHub CLI Parity**: Implements the core functionality of the GitHub CLI
- **Secure**: Only performs actions within the permissions of your GitHub token
- **Efficient**: Written in Go for high performance and minimal resource usage
- **Standalone Binary**: No runtime dependencies required
- **Multi-Platform**: Works on macOS, Linux, and Windows

## Installation

### Download Binary

Download the pre-built binary for your platform from the [Releases](https://github.com/your-username/github-mcp-server-go/releases) page.

### Build from Source

```bash
# Clone the repository
git clone https://github.com/your-username/github-mcp-server-go.git
cd github-mcp-server-go

# Build the binary
go build -o github-mcp-server

# Move the binary to a directory in your PATH (optional)
sudo mv github-mcp-server /usr/local/bin/
```

## Usage

### GitHub Personal Access Token

Before using the GitHub MCP Server, you'll need a GitHub Personal Access Token with appropriate permissions:

1. Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. Create a new token with the `repo` scope (for full repository access) or `public_repo` for only public repository access
3. Copy the generated token for use with the MCP server

### Running the Server

You can run the GitHub MCP Server in two ways:

```bash
# Run with token provided as a command-line flag
./github-mcp-server -token YOUR_GITHUB_TOKEN

# Or run with token provided as an environment variable
export GITHUB_PERSONAL_ACCESS_TOKEN=YOUR_GITHUB_TOKEN
./github-mcp-server
```

### Integration with Claude Desktop

To use GitHub MCP Server with Claude Desktop:

1. Edit your Claude Desktop configuration file:
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Windows: `%APPDATA%\Claude\claude_desktop_config.json`
   - Linux: `~/.config/Claude/claude_desktop_config.json`

2. Add the GitHub MCP Server to the configuration:

```json
{
  "mcpServers": {
    "github": {
      "command": "/path/to/github-mcp-server",
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "YOUR_GITHUB_TOKEN"
      }
    }
  }
}
```

3. Restart Claude Desktop to apply the changes

## Supported Tools

GitHub MCP Server implements the following tool categories:

### Repository Management
- `get_repository`: Get repository details
- `list_repositories`: List user repositories
- `create_repository`: Create a new repository

### Issue Management
- `get_issue`: Get issue details
- `list_issues`: List repository issues
- `create_issue`: Create a new issue
- `close_issue`: Close an existing issue

### Pull Request Management
- `get_pull_request`: Get pull request details
- `list_pull_requests`: List repository pull requests
- `create_pull_request`: Create a new pull request
- `merge_pull_request`: Merge an existing pull request

### GitHub Actions
- `list_workflows`: List repository workflows
- `list_workflow_runs`: List workflow runs
- `trigger_workflow`: Trigger a workflow

### File Operations
- `get_file_content`: Get file content
- `create_file`: Create a new file
- `update_file`: Update an existing file
- `delete_file`: Delete a file

### Search Operations
- `search_code`: Search repositories for code
- `search_issues`: Search for issues and pull requests

## Examples

Here are some examples of how to use the tools with Claude:

### Create a New Repository

```
Please create a new GitHub repository named "awesome-project"
```

### Get Repository Information

```
Get information about the "cli/cli" repository
```

### Create an Issue

```
Create an issue in my "project-name" repository with the title "Fix button styling" and a description of the problem
```

### Search for Code

```
Search for examples of React components that use the useState hook in my repositories
```

## Security Considerations

- GitHub MCP Server only performs actions within the permissions of your GitHub Personal Access Token
- For maximum security, create a token with only the necessary permissions for your use case
- The server does not store your token, but transmits it with each request to the GitHub API
- Consider hosting the server locally rather than on a remote server

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- Inspired by the [GitHub CLI](https://github.com/cli/cli)
- Built with the [Model Context Protocol](https://modelcontextprotocol.io) specification