package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// ConfigStore defines the interface for configuration storage
type ConfigStore interface {
	// Get retrieves a configuration value by key
	Get(key string) (interface{}, error)

	// Set stores a configuration value by key
	Set(key string, value interface{}) error

	// Delete removes a configuration value by key
	Delete(key string) error

	// List returns all configuration values
	List() (map[string]interface{}, error)
}

// FileStore implements ConfigStore using a JSON file
type FileStore struct {
	path  string
	cache map[string]interface{}
	mu    sync.RWMutex
}

// NewFileStore creates a new file-based configuration store
func NewFileStore(path string) (*FileStore, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	store := &FileStore{
		path:  path,
		cache: make(map[string]interface{}),
	}

	// Load existing config if it exists
	if err := store.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return store, nil
}

// load reads the configuration file into memory
func (s *FileStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.cache)
}

// save writes the in-memory configuration to disk
func (s *FileStore) save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s.cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to temporary file first
	tmpFile := s.path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Rename temporary file to actual config file (atomic operation)
	if err := os.Rename(tmpFile, s.path); err != nil {
		os.Remove(tmpFile) // Clean up temp file if rename fails
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// Get retrieves a configuration value by key
func (s *FileStore) Get(key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.cache[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	return value, nil
}

// Set stores a configuration value by key
func (s *FileStore) Set(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update cache
	s.cache[key] = value

	// Save to disk
	if err := s.save(); err != nil {
		delete(s.cache, key) // Revert cache if save fails
		return err
	}

	return nil
}

// Delete removes a configuration value by key
func (s *FileStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if key exists
	if _, ok := s.cache[key]; !ok {
		return fmt.Errorf("key not found: %s", key)
	}

	// Remove from cache
	delete(s.cache, key)

	// Save to disk
	if err := s.save(); err != nil {
		return fmt.Errorf("failed to save config after delete: %w", err)
	}

	return nil
}

// List returns all configuration values
func (s *FileStore) List() (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy of the cache to avoid external modifications
	result := make(map[string]interface{}, len(s.cache))
	for k, v := range s.cache {
		result[k] = v
	}

	return result, nil
}
