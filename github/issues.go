package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// CreateIssueRequest represents parameters for creating an issue
type CreateIssueRequest struct {
	Title     string   `json:"title"`
	Body      string   `json:"body,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Labels    []string `json:"labels,omitempty"`
}

// GetIssue gets an issue by number
func (c *Client) GetIssue(ctx context.Context, owner, repo string, number int) (*Issue, error) {
	url := fmt.Sprintf("repos/%s/%s/issues/%d", owner, repo, number)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &issue, nil
}

// ListIssues lists issues in a repository
func (c *Client) ListIssues(ctx context.Context, owner, repo string, opts *ListIssuesOptions) ([]Issue, error) {
	url := fmt.Sprintf("repos/%s/%s/issues", owner, repo)
	if opts != nil {
		params := make([]string, 0)
		if opts.State != "" {
			params = append(params, "state="+opts.State)
		}
		if len(opts.Labels) > 0 {
			params = append(params, "labels="+string(opts.Labels[0]))
		}
		if opts.Page > 0 {
			params = append(params, "page="+strconv.Itoa(opts.Page))
		}
		if opts.PerPage > 0 {
			params = append(params, "per_page="+strconv.Itoa(opts.PerPage))
		}
		if len(params) > 0 {
			url += "?" + params[0]
			for i := 1; i < len(params); i++ {
				url += "&" + params[i]
			}
		}
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

	var issues []Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return issues, nil
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(ctx context.Context, owner, repo string, req *CreateIssueRequest) (*Issue, error) {
	url := fmt.Sprintf("repos/%s/%s/issues", owner, repo)

	request, err := c.newRequest(ctx, "POST", url, req)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &issue, nil
}

// CloseIssue closes an issue
func (c *Client) CloseIssue(ctx context.Context, owner, repo string, number int) error {
	url := fmt.Sprintf("repos/%s/%s/issues/%d", owner, repo, number)

	update := map[string]string{
		"state": "closed",
	}

	request, err := c.newRequest(ctx, "PATCH", url, update)
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

// UpdateIssue updates an existing issue
func (c *Client) UpdateIssue(ctx context.Context, owner, repo string, number int, update map[string]interface{}) (*Issue, error) {
	url := fmt.Sprintf("repos/%s/%s/issues/%d", owner, repo, number)

	request, err := c.newRequest(ctx, "PATCH", url, update)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &issue, nil
}

// AddIssueComment adds a comment to an issue
func (c *Client) AddIssueComment(ctx context.Context, owner, repo string, number int, body string) error {
	url := fmt.Sprintf("repos/%s/%s/issues/%d/comments", owner, repo, number)

	comment := map[string]string{
		"body": body,
	}

	request, err := c.newRequest(ctx, "POST", url, comment)
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
