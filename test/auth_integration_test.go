package test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github-mcp-server-go/protocol"
	"github-mcp-server-go/server"
)

func TestAuthIntegration(t *testing.T) {
	// Create temporary config directory
	tmpDir, err := os.MkdirTemp("", "auth-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create server instance with test config
	srv := server.New(server.Config{
		Token:     "dummy-token", // This is replaced by login
		ConfigDir: tmpDir,
	})

	// Initialize server
	initParams := &protocol.InitializeParams{
		ProtocolVersion: protocol.LatestProtocolVersion,
	}
	initRequest := protocol.NewRequest(1, "initialize", initParams)
	initResponse := srv.HandleRequest(context.Background(), initRequest)
	if initResponse.Error != nil {
		t.Fatalf("Failed to initialize server: %v", initResponse.Error)
	}

	// Test login
	loginArgs := map[string]interface{}{
		"token":  "ghp_test123",
		"scopes": []string{"repo", "user"},
	}
	loginRequest := protocol.NewRequest(2, "tools/call", &protocol.CallToolParams{
		Name:      "auth_login_token",
		Arguments: loginArgs,
	})
	loginResponse := srv.HandleRequest(context.Background(), loginRequest)
	if loginResponse.Error != nil {
		t.Fatalf("Failed to login: %v", loginResponse.Error)
	}

	// Verify token file exists
	tokenFile := filepath.Join(tmpDir, "auth", "pat.json")
	if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
		t.Error("Token file was not created")
	}

	// Test auth status
	statusRequest := protocol.NewRequest(3, "tools/call", &protocol.CallToolParams{
		Name:      "auth_status",
		Arguments: map[string]interface{}{},
	})
	statusResponse := srv.HandleRequest(context.Background(), statusRequest)
	if statusResponse.Error != nil {
		t.Fatalf("Failed to get auth status: %v", statusResponse.Error)
	}

	// Verify status shows authenticated
	var result protocol.CallToolResult
	if err := parseResponse(statusResponse, &result); err != nil {
		t.Fatalf("Failed to parse status response: %v", err)
	}
	if len(result.Content) == 0 {
		t.Fatal("Expected non-empty content in status result")
	}
	if result.Content[0].Error {
		t.Errorf("Unexpected error in status content: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, "Authenticated") {
		t.Errorf("Expected authenticated status, got: %s", result.Content[0].Text)
	}

	// Test logout
	logoutRequest := protocol.NewRequest(4, "tools/call", &protocol.CallToolParams{
		Name:      "auth_logout",
		Arguments: map[string]interface{}{},
	})
	logoutResponse := srv.HandleRequest(context.Background(), logoutRequest)
	if logoutResponse.Error != nil {
		t.Fatalf("Failed to logout: %v", logoutResponse.Error)
	}

	// Verify token file is removed
	if _, err := os.Stat(tokenFile); !os.IsNotExist(err) {
		t.Error("Token file still exists after logout")
	}

	// Verify status shows not authenticated
	statusResponse = srv.HandleRequest(context.Background(), statusRequest)
	if statusResponse.Error != nil {
		t.Fatalf("Failed to get auth status after logout: %v", statusResponse.Error)
	}
	if err := parseResponse(statusResponse, &result); err != nil {
		t.Fatalf("Failed to parse status response: %v", err)
	}
	if len(result.Content) == 0 {
		t.Fatal("Expected non-empty content in status result")
	}
	if !strings.Contains(result.Content[0].Text, "Not authenticated") {
		t.Errorf("Expected 'Not authenticated' status, got: %s", result.Content[0].Text)
	}
}

func parseResponse(response *protocol.Message, dest interface{}) error {
	if response.Result == nil {
		return fmt.Errorf("no result in response")
	}
	data, err := json.Marshal(response.Result)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}
