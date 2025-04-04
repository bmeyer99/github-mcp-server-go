# GitHub MCP Server Project Structure and Overview

## Directory Structure

```
github-mcp-server-go/
├── github/              # GitHub API client
│   ├── client.go        # HTTP client implementation
│   └── models.go        # GitHub API data models
├── protocol/            # MCP protocol implementation
│   └── protocol.go      # Protocol types and utilities
├── server/              # MCP server implementation
│   ├── server.go        # Core server implementation
│   ├── tools.go         # Tool implementations
│   └── tooldefs.go      # Tool definitions
├── transport/           # MCP transport implementation
│   └── transport.go     # Transport interfaces and implementations
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── main.go              # Entry point
├── LICENSE              # MIT license
└── README.md            # Documentation
```

## Component Overview

### Main (`main.go`)

The entry point for the application. It:
- Parses command line flags
- Sets up the logging system
- Creates a new server with the provided GitHub token
- Handles graceful shutdown via OS signals
- Starts the server with the appropriate transport

### Protocol (`protocol/protocol.go`)

Defines the core MCP protocol types and utilities:
- JSON-RPC 2.0 message structure
- MCP-specific message types
- Server capabilities (Tools, Resources, Prompts)
- Error codes and utility functions

### Transport (`transport/transport.go`)

Implements the communication layer between the MCP client and server:
- `Transport` interface for communication abstraction
- `StdioTransport` for stdin/stdout communication
- Placeholder for `SSETransport` for HTTP-based communication

### GitHub API Client (`github/client.go`)

Provides a high-level interface to the GitHub API:
- HTTP client with authentication
- Repository operations (get, list, create)
- Issue operations (get, list, create, close)
- Pull request operations (get, list, create, merge)
- File operations (get, create, update, delete)
- GitHub Actions operations (list workflows, list runs, trigger)
- Search operations (code, issues)

### Server Core (`server/server.go`)

Implements the core MCP server functionality:
- Message handling loop
- Request routing
- Tool registry
- Initialize handler
- List tools handler
- Call tool handler

### Tool Implementations (`server/tools.go`)

Implements the handlers for each tool:
- Repository tool handlers
- Issue tool handlers
- Pull request tool handlers
- GitHub Actions tool handlers
- File tool handlers
- Search tool handlers

### Tool Definitions (`server/tooldefs.go`)

Defines the schema for each tool:
- Property definitions
- Required fields
- Descriptions
- Default values
- Value constraints (enums, patterns)

## Flow of Execution

1. **Initialization**:
   - `main.go` parses command line arguments and environment variables
   - Server is created with configuration
   - GitHub client is initialized with provided token
   - Tools are registered with the server

2. **Message Handling**:
   - Server starts listening for messages via the configured transport
   - When a message is received, it is parsed as a JSON-RPC request
   - If the message is an initialize request, the server performs initialization
   - Otherwise, the server validates that it has been initialized
   - The server routes the request to the appropriate handler based on the method

3. **Tool Execution**:
   - For tools/list requests, the server returns a list of available tools
   - For tools/call requests, the server:
     - Validates the tool name and arguments
     - Calls the appropriate tool handler
     - Formats the result and sends it back

4. **GitHub API Interaction**:
   - Tool handlers use the GitHub client to interact with the GitHub API
   - The client makes HTTP requests with the provided token
   - Responses are parsed and returned to the tool handler
   - The tool handler formats the response as a tool result

5. **Response Handling**:
   - Results are formatted as JSON-RPC responses
   - Responses are sent back via the configured transport
   - Any errors are formatted as JSON-RPC error responses

## Authentication Flow

1. User provides GitHub Personal Access Token via command line flag or environment variable
2. Token is stored in the server configuration
3. GitHub client is initialized with the token
4. All API requests include the token in the Authorization header
5. GitHub API validates the token and permissions for each request

## Error Handling

- Protocol-level errors (e.g., parse errors, method not found) are returned as JSON-RPC error responses
- GitHub API errors are caught and formatted as tool result errors
- Transport errors are logged and may result in termination of the server

## Extension Points

The server is designed to be easily extended:

1. **Add New Tools**:
   - Add a new tool handler in `server/tools.go`
   - Register the handler in the appropriate register function
   - Add a tool definition in `server/tooldefs.go`

2. **Add New Transports**:
   - Implement the `Transport` interface in `transport/transport.go`
   - Add transport creation and selection logic in `main.go`

3. **Add New GitHub API Functionality**:
   - Add new methods to the GitHub client in `github/client.go`
   - Add new tool handlers that use these methods

4. **Add Resource Support**:
   - Implement resource-related protocol methods in `server/server.go`
   - Add resource handlers for GitHub resources
   - Update server capabilities to include resources