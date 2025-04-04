package ssh

import (
	"context"
	"fmt"
	"github-mcp-server-go/github"
)

// GitHubKeyManager handles GitHub SSH key operations
type GitHubKeyManager struct {
	client *github.Client
	store  *KeyStore
}

// NewGitHubKeyManager creates a new GitHubKeyManager instance
func NewGitHubKeyManager(client *github.Client, store *KeyStore) *GitHubKeyManager {
	return &GitHubKeyManager{
		client: client,
		store:  store,
	}
}

// UploadKey uploads an SSH key to GitHub
func (m *GitHubKeyManager) UploadKey(ctx context.Context, keyID string, title string) error {
	// Load the key from storage
	key, err := m.store.Load(keyID)
	if err != nil {
		return fmt.Errorf("failed to load key: %w", err)
	}

	// Upload key to GitHub
	err = m.client.CreateSSHKey(ctx, title, key.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to upload key to GitHub: %w", err)
	}

	return nil
}

// ListGitHubKeys returns the list of SSH keys registered on GitHub
func (m *GitHubKeyManager) ListGitHubKeys(ctx context.Context) ([]github.SSHKey, error) {
	keys, err := m.client.ListSSHKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list GitHub SSH keys: %w", err)
	}

	return keys, nil
}

// DeleteGitHubKey removes an SSH key from GitHub
func (m *GitHubKeyManager) DeleteGitHubKey(ctx context.Context, keyID int64) error {
	err := m.client.DeleteSSHKey(ctx, keyID)
	if err != nil {
		return fmt.Errorf("failed to delete GitHub SSH key: %w", err)
	}

	return nil
}

// VerifyKey verifies that a key exists on GitHub
func (m *GitHubKeyManager) VerifyKey(ctx context.Context, keyID string) error {
	// Load the key from storage
	key, err := m.store.Load(keyID)
	if err != nil {
		return fmt.Errorf("failed to load key: %w", err)
	}

	// Get GitHub keys
	githubKeys, err := m.ListGitHubKeys(ctx)
	if err != nil {
		return err
	}

	// Check if our key exists
	keyExists := false
	for _, githubKey := range githubKeys {
		if githubKey.Key == key.PublicKey {
			keyExists = true
			break
		}
	}

	if !keyExists {
		return fmt.Errorf("key not found on GitHub")
	}

	return nil
}

// SyncKeys synchronizes local keys with GitHub
func (m *GitHubKeyManager) SyncKeys(ctx context.Context) error {
	// Get all local keys
	localKeys, err := m.store.List()
	if err != nil {
		return fmt.Errorf("failed to list local keys: %w", err)
	}

	// Get all GitHub keys
	githubKeys, err := m.ListGitHubKeys(ctx)
	if err != nil {
		return fmt.Errorf("failed to list GitHub keys: %w", err)
	}

	// Create map of GitHub keys by their public key
	githubKeyMap := make(map[string]*github.SSHKey)
	for i := range githubKeys {
		githubKeyMap[githubKeys[i].Key] = &githubKeys[i]
	}

	// Check each local key
	for _, localKey := range localKeys {
		// Load full key including private key
		fullKey, err := m.store.Load(localKey.ID)
		if err != nil {
			return fmt.Errorf("failed to load key %s: %w", localKey.ID, err)
		}

		// Check if key exists on GitHub
		if _, exists := githubKeyMap[fullKey.PublicKey]; !exists {
			// Upload missing key
			title := fmt.Sprintf("%s (synced)", fullKey.Name)
			if err := m.UploadKey(ctx, fullKey.ID, title); err != nil {
				return fmt.Errorf("failed to sync key %s: %w", fullKey.ID, err)
			}
		}
	}

	return nil
}
