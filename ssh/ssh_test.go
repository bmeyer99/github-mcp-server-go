package ssh

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestKeyGeneration(t *testing.T) {
	generator := NewKeyGenerator()

	tests := []struct {
		name    string
		keyType KeyType
	}{
		{
			name:    "ED25519 Key",
			keyType: KeyTypeED25519,
		},
		{
			name:    "RSA Key",
			keyType: KeyTypeRSA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := generator.Generate("test-key", tt.keyType)
			if err != nil {
				t.Fatalf("Failed to generate key: %v", err)
			}

			// Verify key properties
			if key.Type != tt.keyType {
				t.Errorf("Expected key type %s, got %s", tt.keyType, key.Type)
			}
			if key.Name != "test-key" {
				t.Errorf("Expected key name test-key, got %s", key.Name)
			}
			if key.PublicKey == "" {
				t.Error("Public key is empty")
			}
			if len(key.PrivateKey) == 0 {
				t.Error("Private key is empty")
			}
			if key.Fingerprint == "" {
				t.Error("Fingerprint is empty")
			}
		})
	}
}

func TestKeyEncryption(t *testing.T) {
	masterKey := GenerateRandomKey()
	encryption := NewKeyEncryption(masterKey)

	// Generate a test key
	generator := NewKeyGenerator()
	key, err := generator.Generate("test-key", KeyTypeED25519)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Save original private key
	originalKey := make([]byte, len(key.PrivateKey))
	copy(originalKey, key.PrivateKey)

	// Test encryption
	if err := encryption.EncryptKey(key); err != nil {
		t.Fatalf("Failed to encrypt key: %v", err)
	}

	// Verify key is encrypted (should be different)
	if bytes.Equal(key.PrivateKey, originalKey) {
		t.Error("Key was not encrypted")
	}

	// Test decryption
	if err := encryption.DecryptKey(key); err != nil {
		t.Fatalf("Failed to decrypt key: %v", err)
	}

	// Verify decrypted key matches original
	if !bytes.Equal(key.PrivateKey, originalKey) {
		t.Error("Decrypted key does not match original")
	}
}

func TestKeyStore(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "ssh-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create key store
	masterKey := GenerateRandomKey()
	store, err := NewKeyStore(tmpDir, masterKey)
	if err != nil {
		t.Fatalf("Failed to create key store: %v", err)
	}

	// Generate a test key
	generator := NewKeyGenerator()
	key, err := generator.Generate("test-key", KeyTypeED25519)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Test key save
	if err := store.Save(key); err != nil {
		t.Fatalf("Failed to save key: %v", err)
	}

	// Verify key file exists
	keyFile := filepath.Join(tmpDir, key.ID+".json")
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		t.Error("Key file not created")
	}

	// Test key load
	loadedKey, err := store.Load(key.ID)
	if err != nil {
		t.Fatalf("Failed to load key: %v", err)
	}

	// Verify loaded key matches original
	if loadedKey.ID != key.ID {
		t.Error("Loaded key ID does not match")
	}
	if loadedKey.Name != key.Name {
		t.Error("Loaded key name does not match")
	}
	if loadedKey.Type != key.Type {
		t.Error("Loaded key type does not match")
	}
	if loadedKey.PublicKey != key.PublicKey {
		t.Error("Loaded public key does not match")
	}
	if !bytes.Equal(loadedKey.PrivateKey, key.PrivateKey) {
		t.Error("Loaded private key does not match")
	}

	// Test key list
	keys, err := store.List()
	if err != nil {
		t.Fatalf("Failed to list keys: %v", err)
	}
	if len(keys) != 1 {
		t.Errorf("Expected 1 key, got %d", len(keys))
	}

	// Test key delete
	if err := store.Delete(key.ID); err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	// Verify key file is deleted
	if _, err := os.Stat(keyFile); !os.IsNotExist(err) {
		t.Error("Key file not deleted")
	}

	// Verify key is removed from store
	if _, err := store.Load(key.ID); err == nil {
		t.Error("Expected error loading deleted key")
	}
}

func TestKeyImport(t *testing.T) {
	generator := NewKeyGenerator()

	// Generate a key first
	key, err := generator.Generate("test-key", KeyTypeED25519)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Try to import the private key
	importedKey, err := generator.Import("imported-key", key.PrivateKey)
	if err != nil {
		t.Fatalf("Failed to import key: %v", err)
	}

	// Verify imported key
	if importedKey.Name != "imported-key" {
		t.Errorf("Expected name imported-key, got %s", importedKey.Name)
	}
	if importedKey.Type != key.Type {
		t.Errorf("Expected type %s, got %s", key.Type, importedKey.Type)
	}
	if importedKey.PublicKey == "" {
		t.Error("Public key is empty")
	}
	if len(importedKey.PrivateKey) == 0 {
		t.Error("Private key is empty")
	}
	if importedKey.Fingerprint == "" {
		t.Error("Fingerprint is empty")
	}
}
