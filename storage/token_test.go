package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestStore(t *testing.T) (*FileSystemStore, func()) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "token-store-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	store, err := NewFileSystemStore(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create FileSystemStore: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return store, cleanup
}

func TestFileSystemStore_StoreAndGetToken(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Create test token
	token := &Token{
		ID:          "test-token",
		AccessToken: "ghp_test123",
		TokenType:   "bearer",
		Scope:       []string{"repo", "user"},
		ExpiresAt:   time.Now().Add(time.Hour).UTC(),
	}

	// Store token
	if err := store.StoreToken(token); err != nil {
		t.Fatalf("StoreToken failed: %v", err)
	}

	// Retrieve token
	retrieved, err := store.GetToken(token.ID)
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	// Verify token data
	if retrieved.ID != token.ID {
		t.Errorf("Expected token ID %s, got %s", token.ID, retrieved.ID)
	}
	if retrieved.AccessToken != token.AccessToken {
		t.Errorf("Expected access token %s, got %s", token.AccessToken, retrieved.AccessToken)
	}
	if retrieved.TokenType != token.TokenType {
		t.Errorf("Expected token type %s, got %s", token.TokenType, retrieved.TokenType)
	}
}

func TestFileSystemStore_ListTokens(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Create multiple test tokens
	tokens := []*Token{
		{
			ID:          "token1",
			AccessToken: "ghp_test1",
			TokenType:   "bearer",
		},
		{
			ID:          "token2",
			AccessToken: "ghp_test2",
			TokenType:   "bearer",
		},
	}

	// Store tokens
	for _, token := range tokens {
		if err := store.StoreToken(token); err != nil {
			t.Fatalf("StoreToken failed: %v", err)
		}
	}

	// List tokens
	listed, err := store.ListTokens()
	if err != nil {
		t.Fatalf("ListTokens failed: %v", err)
	}

	// Verify token count
	if len(listed) != len(tokens) {
		t.Errorf("Expected %d tokens, got %d", len(tokens), len(listed))
	}

	// Verify each token exists
	for _, token := range tokens {
		found := false
		for _, listedToken := range listed {
			if listedToken.ID == token.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Token %s not found in listed tokens", token.ID)
		}
	}
}

func TestFileSystemStore_DeleteToken(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Create and store a token
	token := &Token{
		ID:          "token-to-delete",
		AccessToken: "ghp_test123",
		TokenType:   "bearer",
	}

	if err := store.StoreToken(token); err != nil {
		t.Fatalf("StoreToken failed: %v", err)
	}

	// Delete token
	if err := store.DeleteToken(token.ID); err != nil {
		t.Fatalf("DeleteToken failed: %v", err)
	}

	// Verify token is deleted
	_, err := store.GetToken(token.ID)
	if err == nil {
		t.Error("Expected error getting deleted token, got nil")
	}

	// Verify token file is deleted
	if _, err := os.Stat(filepath.Join(store.BasePath, token.ID+".json")); !os.IsNotExist(err) {
		t.Error("Token file still exists after deletion")
	}
}

func TestFileSystemStore_TokenNotFound(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Try to get non-existent token
	_, err := store.GetToken("non-existent-token")
	if err == nil {
		t.Error("Expected error getting non-existent token, got nil")
	}

	// Try to delete non-existent token
	err = store.DeleteToken("non-existent-token")
	if err == nil {
		t.Error("Expected error deleting non-existent token, got nil")
	}
}
