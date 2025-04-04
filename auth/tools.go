package auth

import (
	"context"
	"fmt"
	"path/filepath"

	"github-mcp-server-go/protocol"
	"github-mcp-server-go/storage"
)

// Tool represents the authentication tools handler
type Tool struct {
	store storage.SecureTokenStore
}

// NewTool creates a new authentication tool instance
func NewTool(storagePath string) (*Tool, error) {
	store, err := storage.NewFileSystemStore(filepath.Join(storagePath, "auth"))
	if err != nil {
		return nil, fmt.Errorf("failed to create token store: %w", err)
	}

	return &Tool{
		store: store,
	}, nil
}

// LoginWithToken handles authentication using a personal access token
func (t *Tool) LoginWithToken(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	tokenValue, ok := args["token"].(string)
	if !ok || tokenValue == "" {
		return nil, fmt.Errorf("token argument is required")
	}

	scopes, _ := args["scopes"].([]string)
	if scopes == nil {
		scopes = []string{"repo", "user"}
	}

	token := &storage.Token{
		ID:          "pat", // Use a fixed ID for PAT
		AccessToken: tokenValue,
		TokenType:   "bearer",
		Scope:       scopes,
	}

	if err := t.store.StoreToken(token); err != nil {
		return nil, fmt.Errorf("failed to store token: %w", err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(fmt.Sprintf("Successfully authenticated with personal access token (scopes: %v)", scopes)),
		},
	}, nil
}

// Logout handles user logout by removing stored tokens
func (t *Tool) Logout(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	if err := t.store.DeleteToken("pat"); err != nil && !isNotFoundError(err) {
		return nil, fmt.Errorf("failed to remove token: %w", err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent("Successfully logged out"),
		},
	}, nil
}

// GetAuthStatus retrieves the current authentication status
func (t *Tool) GetAuthStatus(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	token, err := t.store.GetToken("pat")
	if err != nil {
		if isNotFoundError(err) {
			return &protocol.CallToolResult{
				Content: []protocol.Content{
					protocol.TextContent("Not authenticated"),
				},
			}, nil
		}
		return nil, fmt.Errorf("failed to get token status: %w", err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(fmt.Sprintf("Authenticated with %s token (scopes: %v)", token.TokenType, token.Scope)),
		},
	}, nil
}

// GetToken retrieves the current authentication token
func (t *Tool) GetToken() (*storage.Token, error) {
	return t.store.GetToken("pat")
}

// isNotFoundError checks if an error is a "not found" error
func isNotFoundError(err error) bool {
	return err != nil && err.Error() == "token not found: pat"
}
