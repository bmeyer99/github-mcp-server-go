// github-mcp-server-go/protocol/protocol.go
package protocol

import (
	"encoding/json"
	"fmt"
)

// Protocol version constants
const (
	LatestProtocolVersion = "1.0"
)

// JSON-RPC constants
const (
	JSONRPCVersion = "2.0"
)

// Error codes
const (
	// JSON-RPC 2.0 error codes
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603

	// MCP-specific error codes
	NotInitialized        = -32001
	AlreadyInitialized    = -32002
	ResourceNotFound      = -32003
	InvalidResource       = -32004
	ToolNotFound          = -32005
	InvalidTool           = -32006
	PromptNotFound        = -32007
	InvalidPrompt         = -32008
	CapabilityNotSupported = -32009
)

// Message represents a JSON-RPC 2.0 message
type Message struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// NewRequest creates a new JSON-RPC request
func NewRequest(id interface{}, method string, params interface{}) *Message {
	return &Message{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  method,
		Params:  params,
	}
}

// NewResponse creates a new JSON-RPC response
func NewResponse(id interface{}, result interface{}) *Message {
	return &Message{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Result:  result,
	}
}

// NewErrorResponse creates a new JSON-RPC error response
func NewErrorResponse(id interface{}, code int, message string, data interface{}) *Message {
	return &Message{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

// Error represents a JSON-RPC 2.0 error
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ServerCapabilities represents the capabilities supported by the server
type ServerCapabilities struct {
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Prompts   *PromptsCapability   `json:"prompts,omitempty"`
}

// ResourcesCapability represents the resources capability
type ResourcesCapability struct {
	// Additional resource capability fields would go here
}

// ToolsCapability represents the tools capability
type ToolsCapability struct {
	// Additional tool capability fields would go here
}

// PromptsCapability represents the prompts capability
type PromptsCapability struct {
	// Additional prompt capability fields would go here
}

// ============== Initialization ==============

// InitializeParams represents the parameters for the initialize method
type InitializeParams struct {
	ProtocolVersion string `json:"protocolVersion"`
}

// InitializeResult represents the result of the initialize method
type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
	Capabilities    ServerCapabilities `json:"capabilities"`
}

// ServerInfo represents information about the server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url,omitempty"`
}

// ============== Tools ==============

// ListToolsParams represents the parameters for the tools/list method
type ListToolsParams struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// ListToolsResult represents the result of the tools/list method
type ListToolsResult struct {
	Tools      []Tool  `json:"tools"`
	NextCursor *string `json:"nextCursor,omitempty"`
}

// Tool represents a tool
type Tool struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Schema      ToolSchema   `json:"schema"`
}

// ToolSchema represents the schema for a tool
type ToolSchema struct {
	Type       string               `json:"type"`
	Properties map[string]Property  `json:"properties"`
	Required   []string             `json:"required,omitempty"`
}

// Property represents a property in a tool schema
type Property struct {
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Pattern     string      `json:"pattern,omitempty"`
	Format      string      `json:"format,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

// CallToolParams represents the parameters for the tools/call method
type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// CallToolResult represents the result of the tools/call method
type CallToolResult struct {
	Content []Content `json:"content"`
}

// Content represents content in a tool result
type Content struct {
	Type  string `json:"type"`
	Text  string `json:"text,omitempty"`
	Error bool   `json:"isError,omitempty"`
}

// TextContent creates a new text content
func TextContent(text string) Content {
	return Content{
		Type: "text",
		Text: text,
	}
}

// ErrorContent creates a new error content
func ErrorContent(text string) Content {
	return Content{
		Type:  "text",
		Text:  text,
		Error: true,
	}
}