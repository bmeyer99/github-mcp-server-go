package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// GetPullRequest gets a pull request by number
func (c *Client) GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	url := fmt.Sprintf("repos/%s/%s/pulls/%d", owner, repo, number)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &pr, nil
}

// ListPullRequests lists pull requests in a repository
func (c *Client) ListPullRequests(ctx context.Context, owner, repo string, opts *ListPullRequestsOptions) ([]*PullRequest, error) {
	url := fmt.Sprintf("repos/%s/%s/pulls", owner, repo)
	if opts != nil {
		params := make([]string, 0)
		if opts.State != "" {
			params = append(params, "state="+opts.State)
		}
		if opts.Head != "" {
			params = append(params, "head="+opts.Head)
		}
		if opts.Base != "" {
			params = append(params, "base="+opts.Base)
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

	var prs []*PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return prs, nil
}

// CreatePullRequest creates a new pull request
func (c *Client) CreatePullRequest(ctx context.Context, owner, repo string, req *CreatePullRequestRequest) (*PullRequest, error) {
	url := fmt.Sprintf("repos/%s/%s/pulls", owner, repo)

	request, err := c.newRequest(ctx, "POST", url, req)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &pr, nil
}

// UpdatePullRequest updates an existing pull request
func (c *Client) UpdatePullRequest(ctx context.Context, owner, repo string, number int, update map[string]interface{}) (*PullRequest, error) {
	url := fmt.Sprintf("repos/%s/%s/pulls/%d", owner, repo, number)

	request, err := c.newRequest(ctx, "PATCH", url, update)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &pr, nil
}

// MergePullRequest merges a pull request
func (c *Client) MergePullRequest(ctx context.Context, owner, repo string, number int, message string, mergeMethod string) error {
	url := fmt.Sprintf("repos/%s/%s/pulls/%d/merge", owner, repo, number)

	merge := map[string]string{
		"merge_method": mergeMethod,
	}
	if message != "" {
		merge["commit_message"] = message
	}

	request, err := c.newRequest(ctx, "PUT", url, merge)
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

// GetPullRequestFiles gets the files changed in a pull request
func (c *Client) GetPullRequestFiles(ctx context.Context, owner, repo string, number int) ([]*PullRequestFile, error) {
	url := fmt.Sprintf("repos/%s/%s/pulls/%d/files", owner, repo, number)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var files []*PullRequestFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return files, nil
}
