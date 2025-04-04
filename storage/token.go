package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Token represents an authentication token with metadata
type Token struct {
	ID           string    `json:"id"`
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Scope        []string  `json:"scope,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SecureTokenStore defines the interface for token storage operations
type SecureTokenStore interface {
	// StoreToken saves a token to the store
	StoreToken(token *Token) error

	// GetToken retrieves a token by ID
	GetToken(id string) (*Token, error)

	// ListTokens returns all stored tokens
	ListTokens() ([]*Token, error)

	// DeleteToken removes a token from the store
	DeleteToken(id string) error
}

// FileSystemStore implements SecureTokenStore using the local filesystem
type FileSystemStore struct {
	BasePath string
	mu       sync.RWMutex
}

// NewFileSystemStore creates a new FileSystemStore instance
func NewFileSystemStore(basePath string) (*FileSystemStore, error) {
	if err := os.MkdirAll(basePath, 0700); err != nil {
		return nil, fmt.Errorf("failed to create token store directory: %w", err)
	}

	return &FileSystemStore{
		BasePath: basePath,
	}, nil
}

// tokenPath returns the full path for a token file
func (s *FileSystemStore) tokenPath(id string) string {
	return filepath.Join(s.BasePath, fmt.Sprintf("%s.json", id))
}

// StoreToken saves a token to the filesystem
func (s *FileSystemStore) StoreToken(token *Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update timestamps
	now := time.Now().UTC()
	if token.CreatedAt.IsZero() {
		token.CreatedAt = now
	}
	token.UpdatedAt = now

	// Marshal token to JSON
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Write to file
	if err := os.WriteFile(s.tokenPath(token.ID), data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// GetToken retrieves a token from the filesystem
func (s *FileSystemStore) GetToken(id string) (*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.tokenPath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("token not found: %s", id)
		}
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var token Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// ListTokens returns all tokens from the filesystem
func (s *FileSystemStore) ListTokens() ([]*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.BasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token directory: %w", err)
	}

	var tokens []*Token
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.BasePath, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read token file %s: %w", entry.Name(), err)
		}

		var token Token
		if err := json.Unmarshal(data, &token); err != nil {
			return nil, fmt.Errorf("failed to unmarshal token from %s: %w", entry.Name(), err)
		}

		tokens = append(tokens, &token)
	}

	return tokens, nil
}

// DeleteToken removes a token from the filesystem
func (s *FileSystemStore) DeleteToken(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.Remove(s.tokenPath(id)); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("token not found: %s", id)
		}
		return fmt.Errorf("failed to delete token file: %w", err)
	}

	return nil
}
