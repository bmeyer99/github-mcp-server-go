package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// GetContent gets the content of a file
func (c *Client) GetContent(ctx context.Context, owner, repo, path, ref string) (*FileContent, error) {
	url := fmt.Sprintf("repos/%s/%s/contents/%s", owner, repo, path)
	if ref != "" {
		url += fmt.Sprintf("?ref=%s", ref)
	}

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var content FileContent
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Decode base64 content if present
	if content.Content != "" && content.Encoding == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(content.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to decode content: %w", err)
		}
		content.Content = string(decoded)
	}

	return &content, nil
}

// CreateFile creates a new file
func (c *Client) CreateFile(ctx context.Context, owner, repo, path string, req *CreateFileRequest) error {
	url := fmt.Sprintf("repos/%s/%s/contents/%s", owner, repo, path)

	// Encode content as base64
	req.Content = base64.StdEncoding.EncodeToString([]byte(req.Content))

	request, err := c.newRequest(ctx, "PUT", url, req)
	if err != nil {
		return err
	}

	resp, err := c.do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// UpdateFile updates an existing file
func (c *Client) UpdateFile(ctx context.Context, owner, repo, path string, req *UpdateFileRequest) error {
	url := fmt.Sprintf("repos/%s/%s/contents/%s", owner, repo, path)

	// Encode content as base64
	req.Content = base64.StdEncoding.EncodeToString([]byte(req.Content))

	request, err := c.newRequest(ctx, "PUT", url, req)
	if err != nil {
		return err
	}

	resp, err := c.do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteFile deletes a file
func (c *Client) DeleteFile(ctx context.Context, owner, repo, path string, req *DeleteFileRequest) error {
	url := fmt.Sprintf("repos/%s/%s/contents/%s", owner, repo, path)

	request, err := c.newRequest(ctx, "DELETE", url, req)
	if err != nil {
		return err
	}

	resp, err := c.do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
