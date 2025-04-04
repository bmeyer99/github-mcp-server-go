package ssh

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// KeyStore handles SSH key storage operations
type KeyStore struct {
	basePath string
	keyMap   map[string]*SSHKey
	mu       sync.RWMutex
	crypto   *KeyEncryption
}

// NewKeyStore creates a new KeyStore instance
func NewKeyStore(basePath string, masterKey string) (*KeyStore, error) {
	if err := os.MkdirAll(basePath, 0700); err != nil {
		return nil, fmt.Errorf("failed to create key store directory: %w", err)
	}

	store := &KeyStore{
		basePath: basePath,
		keyMap:   make(map[string]*SSHKey),
		crypto:   NewKeyEncryption(masterKey),
	}

	// Load existing keys
	if err := store.loadKeys(); err != nil {
		return nil, err
	}

	return store, nil
}

// Save stores an SSH key
func (s *KeyStore) Save(key *SSHKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Encrypt the private key before saving
	if err := s.crypto.EncryptKey(key); err != nil {
		return fmt.Errorf("failed to encrypt key: %w", err)
	}

	// Save to memory
	s.keyMap[key.ID] = key

	// Save to disk
	if err := s.saveKey(key); err != nil {
		delete(s.keyMap, key.ID)
		return err
	}

	return nil
}

// Load retrieves an SSH key by ID
func (s *KeyStore) Load(id string) (*SSHKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key, ok := s.keyMap[id]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", id)
	}

	// Create a copy of the key
	keyCopy := *key

	// Decrypt the private key
	if err := s.crypto.DecryptKey(&keyCopy); err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %w", err)
	}

	return &keyCopy, nil
}

// List returns all stored SSH keys
func (s *KeyStore) List() ([]*SSHKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]*SSHKey, 0, len(s.keyMap))
	for _, key := range s.keyMap {
		// Create a copy without the encrypted private key
		keyCopy := *key
		keyCopy.PrivateKey = nil
		keys = append(keys, &keyCopy)
	}

	return keys, nil
}

// Delete removes an SSH key
func (s *KeyStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.keyMap[id]; !ok {
		return fmt.Errorf("key not found: %s", id)
	}

	// Delete from disk first
	if err := s.deleteKey(id); err != nil {
		return err
	}

	// Delete from memory
	delete(s.keyMap, id)

	return nil
}

// loadKeys loads all keys from disk
func (s *KeyStore) loadKeys() error {
	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return fmt.Errorf("failed to read key store directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		keyID := strings.TrimSuffix(entry.Name(), ".json")
		key, err := s.loadKeyFromDisk(keyID)
		if err != nil {
			return fmt.Errorf("failed to load key %s: %w", keyID, err)
		}

		s.keyMap[key.ID] = key
	}

	return nil
}

// saveKey saves a key to disk
func (s *KeyStore) saveKey(key *SSHKey) error {
	data, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal key data: %w", err)
	}

	filename := filepath.Join(s.basePath, key.ID+".json")
	tempFile := filename + ".tmp"

	// Write to temporary file first
	if err := os.WriteFile(tempFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	// Rename to final filename (atomic operation)
	if err := os.Rename(tempFile, filename); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to save key file: %w", err)
	}

	return nil
}

// loadKeyFromDisk loads a key from disk
func (s *KeyStore) loadKeyFromDisk(id string) (*SSHKey, error) {
	filename := filepath.Join(s.basePath, id+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	var key SSHKey
	if err := json.Unmarshal(data, &key); err != nil {
		return nil, fmt.Errorf("failed to unmarshal key data: %w", err)
	}

	return &key, nil
}

// deleteKey removes a key file from disk
func (s *KeyStore) deleteKey(id string) error {
	filename := filepath.Join(s.basePath, id+".json")
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete key file: %w", err)
	}
	return nil
}
