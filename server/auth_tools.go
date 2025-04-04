package server

import (
	"github-mcp-server-go/protocol"
)

func loginWithTokenToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "auth_login_token",
		Description: "Login with a GitHub Personal Access Token",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"token": {
					Type:        "string",
					Description: "GitHub Personal Access Token",
				},
				"scopes": {
					Type:        "array",
					Description: "Token scopes (optional, defaults to ['repo', 'user'])",
				},
			},
			Required: []string{"token"},
		},
	}
}

func logoutToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "auth_logout",
		Description: "Logout and remove stored authentication credentials",
		Schema: protocol.ToolSchema{
			Type:       "object",
			Properties: map[string]protocol.Property{},
		},
	}
}

func getAuthStatusToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "auth_status",
		Description: "Get current authentication status",
		Schema: protocol.ToolSchema{
			Type:       "object",
			Properties: map[string]protocol.Property{},
		},
	}
}
