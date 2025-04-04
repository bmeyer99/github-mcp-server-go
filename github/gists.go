package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github-mcp-server-go/gist"
)

// CreateGist creates a new gist
func (c *Client) CreateGist(ctx context.Context, description string, files map[string]gist.GistFile, public bool) (*gist.Gist, error) {
	url := "gists"

	payload := map[string]interface{}{
		"description": description,
		"public":      public,
		"files":       files,
	}

	req, err := c.newRequest(ctx, "POST", url, payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gistResponse gist.Gist
	if err := json.NewDecoder(resp.Body).Decode(&gistResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &gistResponse, nil
}

// ListGists returns all gists for the authenticated user
func (c *Client) ListGists(ctx context.Context) ([]*gist.Gist, error) {
	url := "gists"

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gists []*gist.Gist
	if err := json.NewDecoder(resp.Body).Decode(&gists); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return gists, nil
}

// GetGist retrieves a gist by ID
func (c *Client) GetGist(ctx context.Context, id string) (*gist.Gist, error) {
	url := fmt.Sprintf("gists/%s", id)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gistResponse gist.Gist
	if err := json.NewDecoder(resp.Body).Decode(&gistResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &gistResponse, nil
}

// UpdateGist updates an existing gist
func (c *Client) UpdateGist(ctx context.Context, id string, description string, files map[string]gist.GistFile) error {
	url := fmt.Sprintf("gists/%s", id)

	payload := map[string]interface{}{
		"description": description,
		"files":       files,
	}

	req, err := c.newRequest(ctx, "PATCH", url, payload)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteGist deletes a gist
func (c *Client) DeleteGist(ctx context.Context, id string) error {
	url := fmt.Sprintf("gists/%s", id)

	req, err := c.newRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ListGistCommits returns the commit history of a gist
func (c *Client) ListGistCommits(ctx context.Context, id string) ([]*GistCommit, error) {
	url := fmt.Sprintf("gists/%s/commits", id)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var commits []*GistCommit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return commits, nil
}

// GistCommit represents a commit in a gist's history
type GistCommit struct {
	URL          string    `json:"url"`
	Version      string    `json:"version"`
	User         User      `json:"user"`
	CommittedAt  time.Time `json:"committed_at"`
	ChangeStatus struct {
		Additions int `json:"additions"`
		Deletions int `json:"deletions"`
		Total     int `json:"total"`
	} `json:"change_status"`
}

// StarGist stars a gist
func (c *Client) StarGist(ctx context.Context, id string) error {
	url := fmt.Sprintf("gists/%s/star", id)

	req, err := c.newRequest(ctx, "PUT", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// UnstarGist unstars a gist
func (c *Client) UnstarGist(ctx context.Context, id string) error {
	url := fmt.Sprintf("gists/%s/star", id)

	req, err := c.newRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// IsGistStarred checks if a gist is starred
func (c *Client) IsGistStarred(ctx context.Context, id string) (bool, error) {
	url := fmt.Sprintf("gists/%s/star", id)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == 204, nil
}
