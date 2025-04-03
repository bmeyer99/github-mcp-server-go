// github-mcp-server-go/server/tooldefs.go
package server

import (
	"github.com/your-username/github-mcp-server-go/protocol"
)

// getToolDefinition returns the definition for a tool by name
func getToolDefinition(name string) *protocol.Tool {
	switch name {
	// Repository tools
	case "get_repository":
		return getRepositoryToolDef()
	case "list_repositories":
		return listRepositoriesToolDef()
	case "create_repository":
		return createRepositoryToolDef()
	
	// Issue tools
	case "get_issue":
		return getIssueToolDef()
	case "list_issues":
		return listIssuesToolDef()
	case "create_issue":
		return createIssueToolDef()
	case "close_issue":
		return closeIssueToolDef()
	
	// Pull request tools
	case "get_pull_request":
		return getPullRequestToolDef()
	case "list_pull_requests":
		return listPullRequestsToolDef()
	case "create_pull_request":
		return createPullRequestToolDef()
	case "merge_pull_request":
		return mergePullRequestToolDef()
	
	// GitHub Actions tools
	case "list_workflows":
		return listWorkflowsToolDef()
	case "list_workflow_runs":
		return listWorkflowRunsToolDef()
	case "trigger_workflow":
		return triggerWorkflowToolDef()
	
	// File tools
	case "get_file_content":
		return getFileContentToolDef()
	case "create_file":
		return createFileToolDef()
	case "update_file":
		return updateFileToolDef()
	case "delete_file":
		return deleteFileToolDef()
	
	// Search tools
	case "search_code":
		return searchCodeToolDef()
	case "search_issues":
		return searchIssuesToolDef()
	
	default:
		return nil
	}
}

// ============== Repository Tool Definitions ==============

func getRepositoryToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "get_repository",
		Description: "Get a repository by owner and name",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

func listRepositoriesToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "list_repositories",
		Description: "List repositories for the authenticated user",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"page": {
					Type:        "number",
					Description: "Page number (1-based)",
					Default:     1,
				},
				"per_page": {
					Type:        "number",
					Description: "Number of results per page (max 100)",
					Default:     30,
				},
			},
		},
	}
}

func createRepositoryToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "create_repository",
		Description: "Create a new repository",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"name": {
					Type:        "string",
					Description: "Repository name",
				},
				"description": {
					Type:        "string",
					Description: "Repository description",
				},
				"private": {
					Type:        "boolean",
					Description: "Whether the repository is private",
					Default:     false,
				},
			},
			Required: []string{"name"},
		},
	}
}

// ============== Issue Tool Definitions ==============

func getIssueToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "get_issue",
		Description: "Get an issue by owner, repo, and number",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"number": {
					Type:        "number",
					Description: "Issue number",
				},
			},
			Required: []string{"owner", "repo", "number"},
		},
	}
}

func listIssuesToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "list_issues",
		Description: "List issues in a repository",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"state": {
					Type:        "string",
					Description: "Issue state (open, closed, all)",
					Enum:        []string{"open", "closed", "all"},
					Default:     "open",
				},
				"labels": {
					Type:        "array",
					Description: "Issue labels",
				},
				"page": {
					Type:        "number",
					Description: "Page number (1-based)",
					Default:     1,
				},
				"per_page": {
					Type:        "number",
					Description: "Number of results per page (max 100)",
					Default:     30,
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

func createIssueToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "create_issue",
		Description: "Create a new issue",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"title": {
					Type:        "string",
					Description: "Issue title",
				},
				"body": {
					Type:        "string",
					Description: "Issue body",
				},
				"assignees": {
					Type:        "array",
					Description: "Issue assignees (usernames)",
				},
				"labels": {
					Type:        "array",
					Description: "Issue labels",
				},
			},
			Required: []string{"owner", "repo", "title"},
		},
	}
}

func closeIssueToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "close_issue",
		Description: "Close an issue",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"number": {
					Type:        "number",
					Description: "Issue number",
				},
			},
			Required: []string{"owner", "repo", "number"},
		},
	}
}

// ============== Pull Request Tool Definitions ==============

func getPullRequestToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "get_pull_request",
		Description: "Get a pull request by owner, repo, and number",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"number": {
					Type:        "number",
					Description: "Pull request number",
				},
			},
			Required: []string{"owner", "repo", "number"},
		},
	}
}

func listPullRequestsToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "list_pull_requests",
		Description: "List pull requests in a repository",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"state": {
					Type:        "string",
					Description: "Pull request state (open, closed, all)",
					Enum:        []string{"open", "closed", "all"},
					Default:     "open",
				},
				"head": {
					Type:        "string",
					Description: "Filter by head branch (e.g., 'username:branch-name')",
				},
				"base": {
					Type:        "string",
					Description: "Filter by base branch (e.g., 'main')",
				},
				"page": {
					Type:        "number",
					Description: "Page number (1-based)",
					Default:     1,
				},
				"per_page": {
					Type:        "number",
					Description: "Number of results per page (max 100)",
					Default:     30,
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

func createPullRequestToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "create_pull_request",
		Description: "Create a new pull request",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"title": {
					Type:        "string",
					Description: "Pull request title",
				},
				"body": {
					Type:        "string",
					Description: "Pull request body",
				},
				"head": {
					Type:        "string",
					Description: "Head branch (e.g., 'username:branch-name' or just 'branch-name')",
				},
				"base": {
					Type:        "string",
					Description: "Base branch (e.g., 'main')",
				},
				"draft": {
					Type:        "boolean",
					Description: "Whether the pull request is a draft",
					Default:     false,
				},
			},
			Required: []string{"owner", "repo", "title", "head", "base"},
		},
	}
}

func mergePullRequestToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "merge_pull_request",
		Description: "Merge a pull request",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"number": {
					Type:        "number",
					Description: "Pull request number",
				},
				"commit_title": {
					Type:        "string",
					Description: "Title for the merge commit",
				},
				"commit_message": {
					Type:        "string",
					Description: "Message for the merge commit",
				},
				"merge_method": {
					Type:        "string",
					Description: "Merge method (merge, squash, rebase)",
					Enum:        []string{"merge", "squash", "rebase"},
					Default:     "merge",
				},
			},
			Required: []string{"owner", "repo", "number"},
		},
	}
}

// ============== GitHub Actions Tool Definitions ==============

func listWorkflowsToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "list_workflows",
		Description: "List workflows in a repository",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"page": {
					Type:        "number",
					Description: "Page number (1-based)",
					Default:     1,
				},
				"per_page": {
					Type:        "number",
					Description: "Number of results per page (max 100)",
					Default:     30,
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

func listWorkflowRunsToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "list_workflow_runs",
		Description: "List workflow runs for a repository workflow",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"workflow_id": {
					Type:        "number",
					Description: "Workflow ID",
				},
				"page": {
					Type:        "number",
					Description: "Page number (1-based)",
					Default:     1,
				},
				"per_page": {
					Type:        "number",
					Description: "Number of results per page (max 100)",
					Default:     30,
				},
			},
			Required: []string{"owner", "repo", "workflow_id"},
		},
	}
}

func triggerWorkflowToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "trigger_workflow",
		Description: "Trigger a workflow run",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"workflow_id": {
					Type:        "number",
					Description: "Workflow ID",
				},
				"ref": {
					Type:        "string",
					Description: "Git reference (branch, tag, or SHA)",
				},
				"inputs": {
					Type:        "object",
					Description: "Workflow inputs",
				},
			},
			Required: []string{"owner", "repo", "workflow_id", "ref"},
		},
	}
}

// ============== File Tool Definitions ==============

func getFileContentToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "get_file_content",
		Description: "Get the content of a file in a repository",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"path": {
					Type:        "string",
					Description: "File path in the repository",
				},
				"ref": {
					Type:        "string",
					Description: "Git reference (branch, tag, or SHA)",
				},
			},
			Required: []string{"owner", "repo", "path"},
		},
	}
}

func createFileToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "create_file",
		Description: "Create a file in a repository",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"path": {
					Type:        "string",
					Description: "File path in the repository",
				},
				"content": {
					Type:        "string",
					Description: "File content",
				},
				"message": {
					Type:        "string",
					Description: "Commit message",
				},
				"branch": {
					Type:        "string",
					Description: "Branch to commit to (default: repository's default branch)",
				},
			},
			Required: []string{"owner", "repo", "path", "content", "message"},
		},
	}
}

func updateFileToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "update_file",
		Description: "Update a file in a repository",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"path": {
					Type:        "string",
					Description: "File path in the repository",
				},
				"content": {
					Type:        "string",
					Description: "New file content",
				},
				"message": {
					Type:        "string",
					Description: "Commit message",
				},
				"sha": {
					Type:        "string",
					Description: "File SHA (required for updates)",
				},
				"branch": {
					Type:        "string",
					Description: "Branch to commit to (default: repository's default branch)",
				},
			},
			Required: []string{"owner", "repo", "path", "content", "message", "sha"},
		},
	}
}

func deleteFileToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "delete_file",
		Description: "Delete a file in a repository",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"path": {
					Type:        "string",
					Description: "File path in the repository",
				},
				"message": {
					Type:        "string",
					Description: "Commit message",
				},
				"sha": {
					Type:        "string",
					Description: "File SHA (required for deletions)",
				},
				"branch": {
					Type:        "string",
					Description: "Branch to commit to (default: repository's default branch)",
				},
			},
			Required: []string{"owner", "repo", "path", "message", "sha"},
		},
	}
}

// ============== Search Tool Definitions ==============

func searchCodeToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "search_code",
		Description: "Search for code in repositories",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"query": {
					Type:        "string",
					Description: "Search query",
				},
				"page": {
					Type:        "number",
					Description: "Page number (1-based)",
					Default:     1,
				},
				"per_page": {
					Type:        "number",
					Description: "Number of results per page (max 100)",
					Default:     30,
				},
			},
			Required: []string{"query"},
		},
	}
}

func searchIssuesToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "search_issues",
		Description: "Search for issues and pull requests",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"query": {
					Type:        "string",
					Description: "Search query",
				},
				"page": {
					Type:        "number",
					Description: "Page number (1-based)",
					Default:     1,
				},
				"per_page": {
					Type:        "number",
					Description: "Number of results per page (max 100)",
					Default:     30,
				},
			},
			Required: []string{"query"},
		},
	}
}