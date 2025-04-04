package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// SearchResult represents a generic search result
type SearchResult struct {
	TotalCount        int  `json:"total_count"`
	IncompleteResults bool `json:"incomplete_results"`
}

// CodeSearchResult represents a code search result
type CodeSearchResult struct {
	SearchResult
	Items []CodeSearchItem `json:"items"`
}

// CodeSearchItem represents a single code search result item
type CodeSearchItem struct {
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	SHA        string     `json:"sha"`
	URL        string     `json:"url"`
	HTMLURL    string     `json:"html_url"`
	Repository Repository `json:"repository"`
	Score      float64    `json:"score"`
}

// IssueSearchResult represents an issue search result
type IssueSearchResult struct {
	SearchResult
	Items []Issue `json:"items"`
}

// SearchCode searches for code in repositories
func (c *Client) SearchCode(ctx context.Context, query string, page, perPage int) (*CodeSearchResult, error) {
	url := "search/code"
	params := []string{
		"q=" + query,
		"page=" + strconv.Itoa(page),
		"per_page=" + strconv.Itoa(perPage),
	}
	url += "?" + params[0]
	for i := 1; i < len(params); i++ {
		url += "&" + params[i]
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

	var result CodeSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// SearchIssues searches for issues and pull requests
func (c *Client) SearchIssues(ctx context.Context, query string, page, perPage int) (*IssueSearchResult, error) {
	url := "search/issues"
	params := []string{
		"q=" + query,
		"page=" + strconv.Itoa(page),
		"per_page=" + strconv.Itoa(perPage),
	}
	url += "?" + params[0]
	for i := 1; i < len(params); i++ {
		url += "&" + params[i]
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

	var result IssueSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// SearchRepositories searches for repositories
func (c *Client) SearchRepositories(ctx context.Context, query string, page, perPage int) (*RepositorySearchResult, error) {
	url := "search/repositories"
	params := []string{
		"q=" + query,
		"page=" + strconv.Itoa(page),
		"per_page=" + strconv.Itoa(perPage),
	}
	url += "?" + params[0]
	for i := 1; i < len(params); i++ {
		url += "&" + params[i]
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

	var result RepositorySearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
