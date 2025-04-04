package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// SSHKey represents a GitHub SSH key
type SSHKey struct {
	ID        int64  `json:"id"`
	Key       string `json:"key"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	Verified  bool   `json:"verified"`
	ReadOnly  bool   `json:"read_only"`
}

// CreateSSHKey adds a new SSH key to GitHub
func (c *Client) CreateSSHKey(ctx context.Context, title, key string) error {
	url := "user/keys"

	payload := map[string]string{
		"title": title,
		"key":   key,
	}

	req, err := c.newRequest(ctx, "POST", url, payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("failed to create SSH key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// ListSSHKeys returns all SSH keys for the authenticated user
func (c *Client) ListSSHKeys(ctx context.Context) ([]SSHKey, error) {
	url := "user/keys"

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list SSH keys: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var keys []SSHKey
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return keys, nil
}

// DeleteSSHKey removes an SSH key from GitHub
func (c *Client) DeleteSSHKey(ctx context.Context, keyID int64) error {
	url := fmt.Sprintf("user/keys/%s", strconv.FormatInt(keyID, 10))

	req, err := c.newRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("failed to delete SSH key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetSSHKey retrieves a specific SSH key by ID
func (c *Client) GetSSHKey(ctx context.Context, keyID int64) (*SSHKey, error) {
	url := fmt.Sprintf("user/keys/%s", strconv.FormatInt(keyID, 10))

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get SSH key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var key SSHKey
	if err := json.NewDecoder(resp.Body).Decode(&key); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &key, nil
}
