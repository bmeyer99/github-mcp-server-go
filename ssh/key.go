package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// KeyType represents the type of SSH key
type KeyType string

const (
	KeyTypeED25519 KeyType = "ed25519"
	KeyTypeRSA     KeyType = "rsa"
)

// SSHKey represents an SSH key pair with metadata
type SSHKey struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        KeyType   `json:"type"`
	PublicKey   string    `json:"public_key"`
	PrivateKey  []byte    `json:"private_key"`
	Fingerprint string    `json:"fingerprint"`
	Added       time.Time `json:"added"`
	LastUsed    time.Time `json:"last_used,omitempty"`
}

// SSHKeyManager handles SSH key operations
type SSHKeyManager interface {
	Generate(name string, keyType KeyType) (*SSHKey, error)
	Import(name string, privateKey []byte) (*SSHKey, error)
	List() ([]*SSHKey, error)
	Get(id string) (*SSHKey, error)
	Delete(id string) error
	UploadToGitHub(keyID string) error
}

// KeyGenerator handles SSH key generation
type KeyGenerator struct{}

// NewKeyGenerator creates a new KeyGenerator instance
func NewKeyGenerator() *KeyGenerator {
	return &KeyGenerator{}
}

// Generate creates a new SSH key pair
func (g *KeyGenerator) Generate(name string, keyType KeyType) (*SSHKey, error) {
	switch keyType {
	case KeyTypeED25519:
		return g.generateED25519(name)
	case KeyTypeRSA:
		return g.generateRSA(name)
	default:
		return nil, fmt.Errorf("unsupported key type: %s", keyType)
	}
}

// generateED25519 creates a new ED25519 key pair
func (g *KeyGenerator) generateED25519(name string) (*SSHKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ED25519 key: %w", err)
	}

	// Convert to SSH format
	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ED25519 public key: %w", err)
	}

	// Create PEM block for private key
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ED25519 private key: %w", err)
	}

	pemBlock := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: privKeyBytes,
	}

	// Create SSH key
	sshKey := &SSHKey{
		ID:          generateKeyID(),
		Name:        name,
		Type:        KeyTypeED25519,
		PublicKey:   string(ssh.MarshalAuthorizedKey(sshPubKey)),
		PrivateKey:  pem.EncodeToMemory(pemBlock),
		Fingerprint: ssh.FingerprintSHA256(sshPubKey),
		Added:       time.Now().UTC(),
	}

	return sshKey, nil
}

// generateRSA creates a new RSA key pair
func (g *KeyGenerator) generateRSA(name string) (*SSHKey, error) {
	// Generate 4096-bit RSA key
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Convert to SSH format
	sshPubKey, err := ssh.NewPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert RSA public key: %w", err)
	}

	// Create PEM block for private key
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	}

	// Create SSH key
	sshKey := &SSHKey{
		ID:          generateKeyID(),
		Name:        name,
		Type:        KeyTypeRSA,
		PublicKey:   string(ssh.MarshalAuthorizedKey(sshPubKey)),
		PrivateKey:  pem.EncodeToMemory(pemBlock),
		Fingerprint: ssh.FingerprintSHA256(sshPubKey),
		Added:       time.Now().UTC(),
	}

	return sshKey, nil
}

// Import imports an existing SSH private key
func (g *KeyGenerator) Import(name string, privateKeyData []byte) (*SSHKey, error) {
	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	var pubKey ssh.PublicKey
	var keyType KeyType

	// Parse private key based on type
	switch block.Type {
	case "OPENSSH PRIVATE KEY":
		privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ED25519 private key: %w", err)
		}
		if ed25519Key, ok := privKey.(ed25519.PrivateKey); ok {
			pubKey, err = ssh.NewPublicKey(ed25519Key.Public())
			if err != nil {
				return nil, fmt.Errorf("failed to convert ED25519 public key: %w", err)
			}
			keyType = KeyTypeED25519
		} else {
			return nil, fmt.Errorf("unsupported private key type")
		}

	case "RSA PRIVATE KEY":
		privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
		}
		pubKey, err = ssh.NewPublicKey(&privKey.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert RSA public key: %w", err)
		}
		keyType = KeyTypeRSA

	default:
		return nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}

	// Create SSH key
	sshKey := &SSHKey{
		ID:          generateKeyID(),
		Name:        name,
		Type:        keyType,
		PublicKey:   string(ssh.MarshalAuthorizedKey(pubKey)),
		PrivateKey:  privateKeyData,
		Fingerprint: ssh.FingerprintSHA256(pubKey),
		Added:       time.Now().UTC(),
	}

	return sshKey, nil
}

// generateKeyID creates a unique key ID
func generateKeyID() string {
	// Generate a random 16-byte ID
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("key_%x", b)
}
