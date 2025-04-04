package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
)

const (
	apiBaseURL = "https://api.github.com/"
)

// Client represents a GitHub API client
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new GitHub API client
func NewClient(token string) *Client {
	return &Client{
		token:      token,
		baseURL:    apiBaseURL,
		httpClient: http.DefaultClient,
	}
}

// newRequest creates a new HTTP request with appropriate headers and base URL
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// do executes an HTTP request and returns the response
func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Check for error status codes
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		var errResp struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			return nil, fmt.Errorf("GitHub API error: %s", errResp.Message)
		}
		return nil, fmt.Errorf("GitHub API error: status code %d", resp.StatusCode)
	}

	return resp, nil
}

// buildURL builds a URL by joining the base URL and path components
func (c *Client) buildURL(pathComponents ...string) string {
	components := append([]string{c.baseURL}, pathComponents...)
	return path.Join(components...)
}

// GetRepository gets a repository by owner and name
func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*Repository, error) {
	url := fmt.Sprintf("repos/%s/%s", owner, repo)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repository Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &repository, nil
}

// ListRepositories lists repositories for the authenticated user
func (c *Client) ListRepositories(ctx context.Context, page, perPage int) ([]*Repository, error) {
	url := fmt.Sprintf("user/repos?page=%d&per_page=%d", page, perPage)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repositories []*Repository
	if err := json.NewDecoder(resp.Body).Decode(&repositories); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return repositories, nil
}

// CreateRepository creates a new repository
func (c *Client) CreateRepository(ctx context.Context, req *CreateRepositoryRequest) (*Repository, error) {
	request, err := c.newRequest(ctx, "POST", "user/repos", req)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repository Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &repository, nil
}
