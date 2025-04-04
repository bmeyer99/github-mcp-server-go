package ssh

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iterationCount = 100000
	keyLength      = 32
	saltLength     = 16
)

// KeyEncryption handles encryption/decryption of SSH keys
type KeyEncryption struct {
	masterKey []byte
}

// NewKeyEncryption creates a new KeyEncryption instance with a master key
func NewKeyEncryption(masterKey string) *KeyEncryption {
	// Generate a fixed-length key using PBKDF2
	salt := make([]byte, saltLength)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		panic(fmt.Sprintf("failed to generate salt: %v", err))
	}

	key := pbkdf2.Key([]byte(masterKey), salt, iterationCount, keyLength, sha256.New)

	return &KeyEncryption{
		masterKey: key,
	}
}

// Encrypt encrypts data using AES-GCM
func (k *KeyEncryption) Encrypt(data []byte) ([]byte, error) {
	// Create cipher block
	block, err := aes.NewCipher(k.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and seal
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return []byte(encoded), nil
}

// Decrypt decrypts data using AES-GCM
func (k *KeyEncryption) Decrypt(data []byte) ([]byte, error) {
	// Decode from base64
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create cipher block
	block, err := aes.NewCipher(k.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]

	// Decrypt and verify
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptKey encrypts an SSH key
func (k *KeyEncryption) EncryptKey(key *SSHKey) error {
	// Only encrypt private key
	encrypted, err := k.Encrypt(key.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}

	key.PrivateKey = encrypted
	return nil
}

// DecryptKey decrypts an SSH key
func (k *KeyEncryption) DecryptKey(key *SSHKey) error {
	// Only decrypt private key
	decrypted, err := k.Decrypt(key.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt private key: %w", err)
	}

	key.PrivateKey = decrypted
	return nil
}

// GenerateRandomKey generates a random master key
func GenerateRandomKey() string {
	key := make([]byte, keyLength)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(fmt.Sprintf("failed to generate random key: %v", err))
	}
	return base64.StdEncoding.EncodeToString(key)
}
