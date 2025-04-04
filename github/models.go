package github

import "time"

// Repository represents a GitHub repository
type Repository struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	Description   string    `json:"description"`
	Private       bool      `json:"private"`
	Fork          bool      `json:"fork"`
	HTMLURL       string    `json:"html_url"`
	CloneURL      string    `json:"clone_url"`
	DefaultBranch string    `json:"default_branch"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateRepositoryRequest represents parameters for creating a repository
type CreateRepositoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private,omitempty"`
	AutoInit    bool   `json:"auto_init,omitempty"`
}

// FileContent represents the content of a file on GitHub
type FileContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int    `json:"size"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	DownloadURL string `json:"download_url"`
}

// CreateFileRequest represents parameters for creating a file
type CreateFileRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
	Branch  string `json:"branch,omitempty"`
}

// UpdateFileRequest represents parameters for updating a file
type UpdateFileRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
	SHA     string `json:"sha"`
	Branch  string `json:"branch,omitempty"`
}

// DeleteFileRequest represents parameters for deleting a file
type DeleteFileRequest struct {
	Message string `json:"message"`
	SHA     string `json:"sha"`
	Branch  string `json:"branch,omitempty"`
}

// Issue represents a GitHub issue
type Issue struct {
	ID        int64      `json:"id"`
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	State     string     `json:"state"`
	Body      string     `json:"body"`
	HTMLURL   string     `json:"html_url"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	Assignees []User     `json:"assignees"`
	Labels    []Label    `json:"labels"`
}

// PullRequestBranch represents a pull request branch
type PullRequestBranch struct {
	Ref  string     `json:"ref"`
	SHA  string     `json:"sha"`
	Repo Repository `json:"repo"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID        int64             `json:"id"`
	Number    int               `json:"number"`
	Title     string            `json:"title"`
	State     string            `json:"state"`
	Body      string            `json:"body"`
	HTMLURL   string            `json:"html_url"`
	Head      PullRequestBranch `json:"head"`
	Base      PullRequestBranch `json:"base"`
	Merged    bool              `json:"merged"`
	MergedAt  *time.Time        `json:"merged_at,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// User represents a GitHub user
type User struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
}

// Label represents a GitHub label
type Label struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

// Workflow represents a GitHub Actions workflow
type Workflow struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WorkflowRun represents a GitHub Actions workflow run
type WorkflowRun struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	Conclusion  string     `json:"conclusion"`
	HeadBranch  string     `json:"head_branch"`
	HeadSHA     string     `json:"head_sha"`
	RunNumber   int        `json:"run_number"`
	Event       string     `json:"event"`
	HTMLURL     string     `json:"html_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// ListIssuesOptions represents options for listing issues
type ListIssuesOptions struct {
	State   string   `json:"state,omitempty"`
	Labels  []string `json:"labels,omitempty"`
	Page    int      `json:"page,omitempty"`
	PerPage int      `json:"per_page,omitempty"`
}

// ListPullRequestsOptions represents options for listing pull requests
type ListPullRequestsOptions struct {
	State   string `json:"state,omitempty"`
	Head    string `json:"head,omitempty"`
	Base    string `json:"base,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}

// RepositorySearchResult represents a repository search result
type RepositorySearchResult struct {
	SearchResult
	Items []Repository `json:"items"`
}

// CreatePullRequestRequest represents parameters for creating a pull request
type CreatePullRequestRequest struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
	Head  string `json:"head"`
	Base  string `json:"base"`
	Draft bool   `json:"draft,omitempty"`
}

// PullRequestFile represents a file changed in a pull request
type PullRequestFile struct {
	SHA         string `json:"sha"`
	Filename    string `json:"filename"`
	Status      string `json:"status"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	Changes     int    `json:"changes"`
	BlobURL     string `json:"blob_url"`
	RawURL      string `json:"raw_url"`
	ContentsURL string `json:"contents_url"`
	Patch       string `json:"patch"`
}

// TriggerWorkflowRequest represents parameters for triggering a workflow
type TriggerWorkflowRequest struct {
	Ref    string                 `json:"ref"`
	Inputs map[string]interface{} `json:"inputs,omitempty"`
}
