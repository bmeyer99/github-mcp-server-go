package server

import (
	"context"
	"fmt"
	"path/filepath"

	"github-mcp-server-go/auth"
	"github-mcp-server-go/protocol"
)

// registerAuthTools registers authentication-related tools
func (s *Server) registerAuthTools() {
	// Initialize auth tool
	authTool, err := auth.NewTool(filepath.Join(s.config.ConfigDir, "auth"))
	if err != nil {
		s.config.Logger.Printf("Failed to initialize auth tool: %v", err)
		return
	}
	s.authTool = authTool

	// Register auth tools
	s.tools["auth_login_token"] = s.handleLoginWithToken
	s.tools["auth_logout"] = s.handleLogout
	s.tools["auth_status"] = s.handleAuthStatus
}

// handleLoginWithToken handles the auth_login_token tool
func (s *Server) handleLoginWithToken(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	if s.authTool == nil {
		return nil, fmt.Errorf("auth tool not initialized")
	}

	// Call auth tool
	return s.authTool.LoginWithToken(ctx, args)
}

// handleLogout handles the auth_logout tool
func (s *Server) handleLogout(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	if s.authTool == nil {
		return nil, fmt.Errorf("auth tool not initialized")
	}

	// Call auth tool
	return s.authTool.Logout(ctx, args)
}

// handleAuthStatus handles the auth_status tool
func (s *Server) handleAuthStatus(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	if s.authTool == nil {
		return nil, fmt.Errorf("auth tool not initialized")
	}

	// Call auth tool
	return s.authTool.GetAuthStatus(ctx, args)
}
