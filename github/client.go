// github-mcp-server-go/github/client.go
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// GitHub API base URL
	apiBaseURL = "https://api.github.com"
	
	// Default user agent
	defaultUserAgent = "github-mcp-server-go"
	
	// Default timeout
	defaultTimeout = 30 * time.Second
)

// Client represents a GitHub API client
type Client struct {
	// HTTP client
	httpClient *http.Client
	
	// Personal access token
	token string
	
	// User agent
	userAgent string
}

// NewClient creates a new GitHub API client
func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		token:     token,
		userAgent: defaultUserAgent,
	}
}

// do performs an HTTP request
func (c *Client) do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	// Create request URL
	url := apiBaseURL + path
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	// Set content type if body is provided
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	
	// Check response status
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		
		// Parse error response
		var ghErr Error
		if err := json.NewDecoder(resp.Body).Decode(&ghErr); err != nil {
			return nil, fmt.Errorf("failed to parse error response: %w", err)
		}
		
		return nil, &ghErr
	}
	
	return resp, nil
}

// Error represents a GitHub API error
type Error struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	Status           int    `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("GitHub API error: %s", e.Message)
}

// ============== Repository Operations ==============

// Repository represents a GitHub repository
type Repository struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	Fork        bool   `json:"fork"`
	HTMLURL     string `json:"html_url"`
	SSHURL      string `json:"ssh_url"`
	CloneURL    string `json:"clone_url"`
	GitURL      string `json:"git_url"`
	DefaultBranch string `json:"default_branch"`
	// Add more fields as needed
}

// GetRepository retrieves a repository by owner and name
func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*Repository, error) {
	path := fmt.Sprintf("/repos/%s/%s", owner, repo)
	
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var repository Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, fmt.Errorf("failed to decode repository: %w", err)
	}
	
	return &repository, nil
}

// ListRepositories lists repositories for the authenticated user
func (c *Client) ListRepositories(ctx context.Context, page, perPage int) ([]Repository, error) {
	path := fmt.Sprintf("/user/repos?page=%d&per_page=%d", page, perPage)
	
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var repositories []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repositories); err != nil {
		return nil, fmt.Errorf("failed to decode repositories: %w", err)
	}
	
	return repositories, nil
}

// CreateRepositoryRequest represents a request to create a repository
type CreateRepositoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private,omitempty"`
	// Add more fields as needed
}

// CreateRepository creates a new repository
func (c *Client) CreateRepository(ctx context.Context, req *CreateRepositoryRequest) (*Repository, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	resp, err := c.do(ctx, http.MethodPost, "/user/repos", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var repository Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, fmt.Errorf("failed to decode repository: %w", err)
	}
	
	return &repository, nil
}

// ============== Issue Operations ==============

// Issue represents a GitHub issue
type Issue struct {
	ID        int64  `json:"id"`
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	HTMLURL   string `json:"html_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ClosedAt  string `json:"closed_at"`
	// Add more fields as needed
}

// GetIssue retrieves an issue by owner, repo, and number
func (c *Client) GetIssue(ctx context.Context, owner, repo string, number int) (*Issue, error) {
	path := fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, number)
	
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode issue: %w", err)
	}
	
	return &issue, nil
}

// ListIssuesOptions represents options for listing issues
type ListIssuesOptions struct {
	State     string
	Labels    []string
	Assignee  string
	Creator   string
	Mentioned string
	Sort      string
	Direction string
	Since     string
	Page      int
	PerPage   int
}

// ListIssues lists issues for a repository
func (c *Client) ListIssues(ctx context.Context, owner, repo string, opts *ListIssuesOptions) ([]Issue, error) {
	// Build query string
	query := fmt.Sprintf("/repos/%s/%s/issues?page=%d&per_page=%d", 
		owner, repo, opts.Page, opts.PerPage)
	
	if opts.State != "" {
		query += "&state=" + opts.State
	}
	
	if len(opts.Labels) > 0 {
		query += "&labels=" + strings.Join(opts.Labels, ",")
	}
	
	// Add other options as needed
	
	resp, err := c.do(ctx, http.MethodGet, query, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var issues []Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, fmt.Errorf("failed to decode issues: %w", err)
	}
	
	return issues, nil
}

// CreateIssueRequest represents a request to create an issue
type CreateIssueRequest struct {
	Title     string   `json:"title"`
	Body      string   `json:"body,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Labels    []string `json:"labels,omitempty"`
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(ctx context.Context, owner, repo string, req *CreateIssueRequest) (*Issue, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	path := fmt.Sprintf("/repos/%s/%s/issues", owner, repo)
	resp, err := c.do(ctx, http.MethodPost, path, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode issue: %w", err)
	}
	
	return &issue, nil
}

// ============== Pull Request Operations ==============

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID        int64  `json:"id"`
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	HTMLURL   string `json:"html_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ClosedAt  string `json:"closed_at"`
	MergedAt  string `json:"merged_at"`
	Head      struct {
		Ref  string `json:"ref"`
		SHA  string `json:"sha"`
		Repo struct {
			ID       int64  `json:"id"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
		} `json:"repo"`
	} `json:"head"`
	Base struct {
		Ref  string `json:"ref"`
		SHA  string `json:"sha"`
		Repo struct {
			ID       int64  `json:"id"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
		} `json:"repo"`
	} `json:"base"`
	// Add more fields as needed
}

// GetPullRequest retrieves a pull request by owner, repo, and number
func (c *Client) GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	path := fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, number)
	
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode pull request: %w", err)
	}
	
	return &pr, nil
}

// ListPullRequestsOptions represents options for listing pull requests
type ListPullRequestsOptions struct {
	State     string
	Head      string
	Base      string
	Sort      string
	Direction string
	Page      int
	PerPage   int
}

// ListPullRequests lists pull requests for a repository
func (c *Client) ListPullRequests(ctx context.Context, owner, repo string, opts *ListPullRequestsOptions) ([]PullRequest, error) {
	// Build query string
	query := fmt.Sprintf("/repos/%s/%s/pulls?page=%d&per_page=%d", 
		owner, repo, opts.Page, opts.PerPage)
	
	if opts.State != "" {
		query += "&state=" + opts.State
	}
	
	if opts.Head != "" {
		query += "&head=" + opts.Head
	}
	
	if opts.Base != "" {
		query += "&base=" + opts.Base
	}
	
	// Add other options as needed
	
	resp, err := c.do(ctx, http.MethodGet, query, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var prs []PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, fmt.Errorf("failed to decode pull requests: %w", err)
	}
	
	return prs, nil
}

// CreatePullRequestRequest represents a request to create a pull request
type CreatePullRequestRequest struct {
	Title               string `json:"title"`
	Body                string `json:"body,omitempty"`
	Head                string `json:"head"`
	Base                string `json:"base"`
	MaintainerCanModify bool   `json:"maintainer_can_modify,omitempty"`
	Draft               bool   `json:"draft,omitempty"`
}

// CreatePullRequest creates a new pull request
func (c *Client) CreatePullRequest(ctx context.Context, owner, repo string, req *CreatePullRequestRequest) (*PullRequest, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	path := fmt.Sprintf("/repos/%s/%s/pulls", owner, repo)
	resp, err := c.do(ctx, http.MethodPost, path, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode pull request: %w", err)
	}
	
	return &pr, nil
}

// ============== File Operations ==============

// Content represents a file or directory in a repository
type Content struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	DownloadURL string `json:"download_url"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
}

// GetContent retrieves content from a repository
func (c *Client) GetContent(ctx context.Context, owner, repo, path, ref string) (*Content, error) {
	query := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)
	
	if ref != "" {
		query += "?ref=" + ref
	}
	
	resp, err := c.do(ctx, http.MethodGet, query, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var content Content
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return nil, fmt.Errorf("failed to decode content: %w", err)
	}
	
	return &content, nil
}

// CreateFileRequest represents a request to create a file
type CreateFileRequest struct {
	Message   string `json:"message"`
	Content   string `json:"content"`
	Branch    string `json:"branch,omitempty"`
	Committer struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"committer,omitempty"`
	Author struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"author,omitempty"`
}

// CreateFile creates a file in a repository
func (c *Client) CreateFile(ctx context.Context, owner, repo, path string, req *CreateFileRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	query := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)
	resp, err := c.do(ctx, http.MethodPut, query, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

// UpdateFileRequest represents a request to update a file
type UpdateFileRequest struct {
	Message   string `json:"message"`
	Content   string `json:"content"`
	SHA       string `json:"sha"`
	Branch    string `json:"branch,omitempty"`
	Committer struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"committer,omitempty"`
	Author struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"author,omitempty"`
}

// UpdateFile updates a file in a repository
func (c *Client) UpdateFile(ctx context.Context, owner, repo, path string, req *UpdateFileRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	query := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)
	resp, err := c.do(ctx, http.MethodPut, query, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

// DeleteFileRequest represents a request to delete a file
type DeleteFileRequest struct {
	Message   string `json:"message"`
	SHA       string `json:"sha"`
	Branch    string `json:"branch,omitempty"`
	Committer struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"committer,omitempty"`
	Author struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"author,omitempty"`
}

// DeleteFile deletes a file in a repository
func (c *Client) DeleteFile(ctx context.Context, owner, repo, path string, req *DeleteFileRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	query := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)
	resp, err := c.do(ctx, http.MethodDelete, query, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

// ============== Branch Operations ==============

// Branch represents a GitHub branch
type Branch struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	Protected bool `json:"protected"`
}

// GetBranch retrieves a branch by owner, repo, and branch name
func (c *Client) GetBranch(ctx context.Context, owner, repo, branch string) (*Branch, error) {
	path := fmt.Sprintf("/repos/%s/%s/branches/%s", owner, repo, branch)
	
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var b Branch
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		return nil, fmt.Errorf("failed to decode branch: %w", err)
	}
	
	return &b, nil
}

// ListBranches lists branches for a repository
func (c *Client) ListBranches(ctx context.Context, owner, repo string, page, perPage int) ([]Branch, error) {
	path := fmt.Sprintf("/repos/%s/%s/branches?page=%d&per_page=%d", owner, repo, page, perPage)
	
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var branches []Branch
	if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
		return nil, fmt.Errorf("failed to decode branches: %w", err)
	}
	
	return branches, nil
}

// ============== Workflow/Actions Operations ==============

// Workflow represents a GitHub Actions workflow
type Workflow struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Path  string `json:"path"`
	State string `json:"state"`
}

// ListWorkflows lists workflows for a repository
func (c *Client) ListWorkflows(ctx context.Context, owner, repo string, page, perPage int) ([]Workflow, error) {
	path := fmt.Sprintf("/repos/%s/%s/actions/workflows?page=%d&per_page=%d", owner, repo, page, perPage)
	
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var response struct {
		TotalCount int        `json:"total_count"`
		Workflows  []Workflow `json:"workflows"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode workflows: %w", err)
	}
	
	return response.Workflows, nil
}

// WorkflowRun represents a GitHub Actions workflow run
type WorkflowRun struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	HeadBranch string `json:"head_branch"`
	HeadSHA    string `json:"head_sha"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	URL        string `json:"url"`
	HTMLURL    string `json:"html_url"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ListWorkflowRuns lists workflow runs for a repository
func (c *Client) ListWorkflowRuns(ctx context.Context, owner, repo string, workflowID int64, page, perPage int) ([]WorkflowRun, error) {
	path := fmt.Sprintf("/repos/%s/%s/actions/workflows/%d/runs?page=%d&per_page=%d", 
		owner, repo, workflowID, page, perPage)
	
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var response struct {
		TotalCount  int           `json:"total_count"`
		WorkflowRuns []WorkflowRun `json:"workflow_runs"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode workflow runs: %w", err)
	}
	
	return response.WorkflowRuns, nil
}

// TriggerWorkflowRequest represents a request to trigger a workflow
type TriggerWorkflowRequest struct {
	Ref        string                 `json:"ref"`
	Inputs     map[string]interface{} `json:"inputs,omitempty"`
}

// TriggerWorkflow triggers a workflow
func (c *Client) TriggerWorkflow(ctx context.Context, owner, repo string, workflowID int64, req *TriggerWorkflowRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	path := fmt.Sprintf("/repos/%s/%s/actions/workflows/%d/dispatches", owner, repo, workflowID)
	resp, err := c.do(ctx, http.MethodPost, path, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

// ============== Search Operations ==============

// SearchCodeResult represents a GitHub code search result
type SearchCodeResult struct {
	TotalCount int   `json:"total_count"`
	Items      []struct {
		Name       string `json:"name"`
		Path       string `json:"path"`
		SHA        string `json:"sha"`
		URL        string `json:"url"`
		HTMLURL    string `json:"html_url"`
		Repository struct {
			ID       int64  `json:"id"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
		} `json:"repository"`
	} `json:"items"`
}

// SearchCode searches code in repositories
func (c *Client) SearchCode(ctx context.Context, query string, page, perPage int) (*SearchCodeResult, error) {
	path := fmt.Sprintf("/search/code?q=%s&page=%d&per_page=%d", query, page, perPage)
	
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result SearchCodeResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search result: %w", err)
	}
	
	return &result, nil
}

// SearchIssuesResult represents a GitHub issues search result
type SearchIssuesResult struct {
	TotalCount int     `json:"total_count"`
	Items      []Issue `json:"items"`
}

// SearchIssues searches issues and pull requests
func (c *Client) SearchIssues(ctx context.Context, query string, page, perPage int) (*SearchIssuesResult, error) {
	path := fmt.Sprintf("/search/issues?q=%s&page=%d&per_page=%d", query, page, perPage)
	
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result SearchIssuesResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search result: %w", err)
	}
	
	return &result, nil
}