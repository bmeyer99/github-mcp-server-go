// github-mcp-server-go/server/server.go
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/your-username/github-mcp-server-go/protocol"
	"github.com/your-username/github-mcp-server-go/transport"
	"github.com/your-username/github-mcp-server-go/github"
)

// Config represents the server configuration
type Config struct {
	// GitHub Personal Access Token
	Token string
	
	// Logger for server logs
	Logger *log.Logger
	
	// Debug mode
	Debug bool
}

// Server represents an MCP server
type Server struct {
	config      Config
	initialized bool
	client      *github.Client
	mu          sync.Mutex
	// registry of supported operations
	tools       map[string]ToolHandler
}

// ToolHandler is a function that handles a tool call
type ToolHandler func(context.Context, map[string]interface{}) (*protocol.CallToolResult, error)

// New creates a new MCP server
func New(config Config) *Server {
	return &Server{
		config: config,
		tools:  make(map[string]ToolHandler),
	}
}

// Serve starts the server with the given transport
func (s *Server) Serve(ctx context.Context, t transport.Transport) error {
	// Initialize GitHub client
	s.client = github.NewClient(s.config.Token)
	
	// Register tools
	s.registerTools()
	
	// Main message handling loop
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Read message
			msg, err := t.ReadMessage(ctx)
			if err != nil {
				s.config.Logger.Printf("Error reading message: %v", err)
				continue
			}

			// Parse message
			var request protocol.Message
			if err := json.Unmarshal(msg, &request); err != nil {
				s.config.Logger.Printf("Error parsing message: %v", err)
				response := protocol.NewErrorResponse(nil, protocol.ParseError, "Invalid JSON", nil)
				if err := s.sendResponse(ctx, t, response); err != nil {
					s.config.Logger.Printf("Error sending response: %v", err)
				}
				continue
			}

			// Handle message
			go func() {
				response := s.handleRequest(ctx, &request)
				if err := s.sendResponse(ctx, t, response); err != nil {
					s.config.Logger.Printf("Error sending response: %v", err)
				}
			}()
		}
	}
}

// sendResponse sends a response message
func (s *Server) sendResponse(ctx context.Context, t transport.Transport, response *protocol.Message) error {
	// Debug log response
	if s.config.Debug {
		respJSON, _ := json.Marshal(response)
		s.config.Logger.Printf("Response: %s", string(respJSON))
	}

	// Marshal response
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// Send response
	return t.WriteMessage(ctx, responseJSON)
}

// handleRequest handles a request message
func (s *Server) handleRequest(ctx context.Context, request *protocol.Message) *protocol.Message {
	// Debug log request
	if s.config.Debug {
		reqJSON, _ := json.Marshal(request)
		s.config.Logger.Printf("Request: %s", string(reqJSON))
	}

	// Check if it's a valid request
	if request.Method == "" {
		return protocol.NewErrorResponse(request.ID, protocol.InvalidRequest, "Missing method", nil)
	}

	// Initialize request
	if request.Method == "initialize" {
		return s.handleInitialize(ctx, request)
	}

	// Check if server is initialized
	if !s.initialized {
		return protocol.NewErrorResponse(request.ID, protocol.NotInitialized, "Server not initialized", nil)
	}

	// Handle method
	switch request.Method {
	case "tools/list":
		return s.handleListTools(ctx, request)
	case "tools/call":
		return s.handleCallTool(ctx, request)
	default:
		return protocol.NewErrorResponse(request.ID, protocol.MethodNotFound, "Method not found", nil)
	}
}

// handleInitialize handles an initialize request
func (s *Server) handleInitialize(ctx context.Context, request *protocol.Message) *protocol.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if already initialized
	if s.initialized {
		return protocol.NewErrorResponse(request.ID, protocol.AlreadyInitialized, "Server already initialized", nil)
	}

	// Parse parameters
	var params protocol.InitializeParams
	if err := parseParams(request.Params, &params); err != nil {
		return protocol.NewErrorResponse(request.ID, protocol.InvalidParams, err.Error(), nil)
	}

	// Set server as initialized
	s.initialized = true

	// Create result
	result := protocol.InitializeResult{
		ProtocolVersion: protocol.LatestProtocolVersion,
		ServerInfo: protocol.ServerInfo{
			Name:    "GitHub MCP Server",
			Version: "1.0.0",
			URL:     "https://github.com/your-username/github-mcp-server-go",
		},
		Capabilities: protocol.ServerCapabilities{
			Tools: &protocol.ToolsCapability{},
		},
	}

	return protocol.NewResponse(request.ID, result)
}

// handleListTools handles a tools/list request
func (s *Server) handleListTools(ctx context.Context, request *protocol.Message) *protocol.Message {
	// Parse parameters
	var params protocol.ListToolsParams
	if err := parseParams(request.Params, &params); err != nil {
		return protocol.NewErrorResponse(request.ID, protocol.InvalidParams, err.Error(), nil)
	}

	// Get list of tools
	tools := make([]protocol.Tool, 0, len(s.tools))
	for name, _ := range s.tools {
		// Get tool definition from registry
		tool := getToolDefinition(name)
		if tool != nil {
			tools = append(tools, *tool)
		}
	}

	// Create result
	result := protocol.ListToolsResult{
		Tools: tools,
		// No pagination for now
		NextCursor: nil,
	}

	return protocol.NewResponse(request.ID, result)
}

// handleCallTool handles a tools/call request
func (s *Server) handleCallTool(ctx context.Context, request *protocol.Message) *protocol.Message {
	// Parse parameters
	var params protocol.CallToolParams
	if err := parseParams(request.Params, &params); err != nil {
		return protocol.NewErrorResponse(request.ID, protocol.InvalidParams, err.Error(), nil)
	}

	// Find tool handler
	handler, ok := s.tools[params.Name]
	if !ok {
		return protocol.NewErrorResponse(request.ID, protocol.ToolNotFound, "Tool not found", nil)
	}

	// Call tool
	result, err := handler(ctx, params.Arguments)
	if err != nil {
		return protocol.NewErrorResponse(request.ID, protocol.InternalError, err.Error(), nil)
	}

	return protocol.NewResponse(request.ID, result)
}

// registerTools registers all supported tools
func (s *Server) registerTools() {
	// Register repository tools
	s.registerRepositoryTools()
	
	// Register issue tools
	s.registerIssueTools()
	
	// Register pull request tools
	s.registerPullRequestTools()
	
	// Register GitHub Actions tools
	s.registerActionsTools()
	
	// Register file tools
	s.registerFileTools()
	
	// Register search tools
	s.registerSearchTools()
}

// Helper function to parse request parameters
func parseParams(params interface{}, dest interface{}) error {
	// Marshal and unmarshal to convert to the correct type
	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal params: %w", err)
	}

	return nil
}

// Placeholder for getting tool definitions - these would be defined elsewhere
func getToolDefinition(name string) *protocol.Tool {
	// This would be implemented to return actual tool definitions
	return nil
}