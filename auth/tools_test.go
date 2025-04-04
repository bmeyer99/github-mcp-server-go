package auth

import (
	"context"
	"os"
	"testing"
)

func setupTestTool(t *testing.T) (*Tool, func()) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "auth-tools-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create tool instance
	tool, err := NewTool(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create auth tool: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tool, cleanup
}

func TestLoginWithToken(t *testing.T) {
	tool, cleanup := setupTestTool(t)
	defer cleanup()

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name: "successful login",
			args: map[string]interface{}{
				"token": "ghp_test123",
				"scopes": []string{
					"repo",
					"user",
				},
			},
			wantErr: false,
		},
		{
			name: "missing token",
			args: map[string]interface{}{
				"scopes": []string{"repo"},
			},
			wantErr: true,
		},
		{
			name: "empty token",
			args: map[string]interface{}{
				"token": "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.LoginWithToken(context.Background(), tt.args)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("LoginWithToken() error = %v", err)
				return
			}

			if len(result.Content) == 0 {
				t.Error("Expected non-empty content in result")
			}

			// Verify token was stored
			token, err := tool.GetToken()
			if err != nil {
				t.Errorf("Failed to retrieve stored token: %v", err)
				return
			}

			if token.AccessToken != tt.args["token"] {
				t.Errorf("Expected token %v, got %v", tt.args["token"], token.AccessToken)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	tool, cleanup := setupTestTool(t)
	defer cleanup()

	// First login to create a token
	loginArgs := map[string]interface{}{
		"token": "ghp_test123",
	}
	if _, err := tool.LoginWithToken(context.Background(), loginArgs); err != nil {
		t.Fatalf("Failed to login for test setup: %v", err)
	}

	// Test logout
	result, err := tool.Logout(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Errorf("Logout() error = %v", err)
		return
	}

	if len(result.Content) == 0 {
		t.Error("Expected non-empty content in result")
	}

	// Verify token was removed
	if _, err := tool.GetToken(); err == nil {
		t.Error("Expected error getting token after logout, got nil")
	}
}

func TestGetAuthStatus(t *testing.T) {
	tool, cleanup := setupTestTool(t)
	defer cleanup()

	// Test status when not authenticated
	result, err := tool.GetAuthStatus(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Errorf("GetAuthStatus() error = %v", err)
		return
	}

	if len(result.Content) == 0 {
		t.Error("Expected non-empty content in result")
	}

	// Login and test authenticated status
	loginArgs := map[string]interface{}{
		"token": "ghp_test123",
		"scopes": []string{
			"repo",
			"user",
		},
	}
	if _, err := tool.LoginWithToken(context.Background(), loginArgs); err != nil {
		t.Fatalf("Failed to login for test setup: %v", err)
	}

	result, err = tool.GetAuthStatus(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Errorf("GetAuthStatus() error = %v", err)
		return
	}

	if len(result.Content) == 0 {
		t.Error("Expected non-empty content in result")
	}
}
