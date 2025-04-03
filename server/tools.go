// github-mcp-server-go/server/tools.go
package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/your-username/github-mcp-server-go/protocol"
	"github.com/your-username/github-mcp-server-go/github"
)

// registerRepositoryTools registers repository-related tools
func (s *Server) registerRepositoryTools() {
	// Get repository
	s.tools["get_repository"] = s.handleGetRepository
	
	// List repositories
	s.tools["list_repositories"] = s.handleListRepositories
	
	// Create repository
	s.tools["create_repository"] = s.handleCreateRepository
}

// registerIssueTools registers issue-related tools
func (s *Server) registerIssueTools() {
	// Get issue
	s.tools["get_issue"] = s.handleGetIssue
	
	// List issues
	s.tools["list_issues"] = s.handleListIssues
	
	// Create issue
	s.tools["create_issue"] = s.handleCreateIssue
	
	// Close issue
	s.tools["close_issue"] = s.handleCloseIssue
}

// registerPullRequestTools registers pull request-related tools
func (s *Server) registerPullRequestTools() {
	// Get pull request
	s.tools["get_pull_request"] = s.handleGetPullRequest
	
	// List pull requests
	s.tools["list_pull_requests"] = s.handleListPullRequests
	
	// Create pull request
	s.tools["create_pull_request"] = s.handleCreatePullRequest
	
	// Merge pull request
	s.tools["merge_pull_request"] = s.handleMergePullRequest
}

// registerActionsTools registers GitHub Actions-related tools
func (s *Server) registerActionsTools() {
	// List workflows
	s.tools["list_workflows"] = s.handleListWorkflows
	
	// List workflow runs
	s.tools["list_workflow_runs"] = s.handleListWorkflowRuns
	
	// Trigger workflow
	s.tools["trigger_workflow"] = s.handleTriggerWorkflow
}

// registerFileTools registers file-related tools
func (s *Server) registerFileTools() {
	// Get file content
	s.tools["get_file_content"] = s.handleGetFileContent
	
	// Create file
	s.tools["create_file"] = s.handleCreateFile
	
	// Update file
	s.tools["update_file"] = s.handleUpdateFile
	
	// Delete file
	s.tools["delete_file"] = s.handleDeleteFile
}

// registerSearchTools registers search-related tools
func (s *Server) registerSearchTools() {
	// Search code
	s.tools["search_code"] = s.handleSearchCode
	
	// Search issues
	s.tools["search_issues"] = s.handleSearchIssues
}

// ============== Repository Tool Handlers ==============

// handleGetRepository handles the get_repository tool
func (s *Server) handleGetRepository(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	// Get repository
	repository, err := s.client.GetRepository(ctx, owner, repo)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to get repository: %v", err)),
			},
		}, nil
	}
	
	// Format result
	result := fmt.Sprintf(`{
	"name": "%s",
	"full_name": "%s",
	"description": "%s",
	"private": %v,
	"default_branch": "%s",
	"html_url": "%s",
	"clone_url": "%s"
}`, repository.Name, repository.FullName, repository.Description, repository.Private, repository.DefaultBranch, repository.HTMLURL, repository.CloneURL)
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// handleListRepositories handles the list_repositories tool
func (s *Server) handleListRepositories(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Parse page and per_page arguments
	page := 1
	if pageArg, ok := args["page"].(float64); ok {
		page = int(pageArg)
	}
	
	perPage := 30
	if perPageArg, ok := args["per_page"].(float64); ok {
		perPage = int(perPageArg)
	}
	
	// List repositories
	repositories, err := s.client.ListRepositories(ctx, page, perPage)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to list repositories: %v", err)),
			},
		}, nil
	}
	
	// Format result
	result := "[\n"
	for i, repo := range repositories {
		result += fmt.Sprintf(`  {
    "name": "%s",
    "full_name": "%s",
    "description": "%s",
    "private": %v,
    "default_branch": "%s",
    "html_url": "%s"
  }`, repo.Name, repo.FullName, repo.Description, repo.Private, repo.DefaultBranch, repo.HTMLURL)
		
		if i < len(repositories)-1 {
			result += ",\n"
		} else {
			result += "\n"
		}
	}
	result += "]"
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// handleCreateRepository handles the create_repository tool
func (s *Server) handleCreateRepository(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name is required and must be a string")
	}
	
	// Parse optional arguments
	description := ""
	if descArg, ok := args["description"].(string); ok {
		description = descArg
	}
	
	private := false
	if privateArg, ok := args["private"].(bool); ok {
		private = privateArg
	}
	
	// Create repository
	req := &github.CreateRepositoryRequest{
		Name:        name,
		Description: description,
		Private:     private,
	}
	
	repository, err := s.client.CreateRepository(ctx, req)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to create repository: %v", err)),
			},
		}, nil
	}
	
	// Format result
	result := fmt.Sprintf(`{
	"name": "%s",
	"full_name": "%s",
	"description": "%s",
	"private": %v,
	"default_branch": "%s",
	"html_url": "%s",
	"clone_url": "%s"
}`, repository.Name, repository.FullName, repository.Description, repository.Private, repository.DefaultBranch, repository.HTMLURL, repository.CloneURL)
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// ============== File Tool Handlers ==============

// handleGetFileContent handles the get_file_content tool
func (s *Server) handleGetFileContent(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	path, ok := args["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("path is required and must be a string")
	}
	
	// Parse optional arguments
	ref := ""
	if refArg, ok := args["ref"].(string); ok {
		ref = refArg
	}
	
	// Get file content
	content, err := s.client.GetContent(ctx, owner, repo, path, ref)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to get file content: %v", err)),
			},
		}, nil
	}
	
	// Decode content if it's base64 encoded
	decodedContent := content.Content
	if content.Encoding == "base64" {
		bytes, err := base64.StdEncoding.DecodeString(content.Content)
		if err != nil {
			return &protocol.CallToolResult{
				Content: []protocol.Content{
					protocol.ErrorContent(fmt.Sprintf("Failed to decode content: %v", err)),
				},
			}, nil
		}
		decodedContent = string(bytes)
	}
	
	// Format result
	result := fmt.Sprintf(`{
	"name": "%s",
	"path": "%s",
	"sha": "%s",
	"size": %d,
	"type": "%s",
	"content": %s
}`, content.Name, content.Path, content.SHA, content.Size, content.Type, strconv.Quote(decodedContent))
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// handleCreateFile handles the create_file tool
func (s *Server) handleCreateFile(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	path, ok := args["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("path is required and must be a string")
	}
	
	content, ok := args["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content is required and must be a string")
	}
	
	message, ok := args["message"].(string)
	if !ok || message == "" {
		return nil, fmt.Errorf("commit message is required and must be a string")
	}
	
	// Parse optional arguments
	branch := ""
	if branchArg, ok := args["branch"].(string); ok {
		branch = branchArg
	}
	
	// Create file
	req := &github.CreateFileRequest{
		Message: message,
		Content: base64.StdEncoding.EncodeToString([]byte(content)),
		Branch:  branch,
	}
	
	err := s.client.CreateFile(ctx, owner, repo, path, req)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to create file: %v", err)),
			},
		}, nil
	}
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(fmt.Sprintf("File '%s' created successfully in repository '%s/%s'", path, owner, repo)),
		},
	}, nil
}

// handleUpdateFile handles the update_file tool
func (s *Server) handleUpdateFile(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	path, ok := args["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("path is required and must be a string")
	}
	
	content, ok := args["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content is required and must be a string")
	}
	
	message, ok := args["message"].(string)
	if !ok || message == "" {
		return nil, fmt.Errorf("commit message is required and must be a string")
	}
	
	sha, ok := args["sha"].(string)
	if !ok || sha == "" {
		return nil, fmt.Errorf("sha is required and must be a string")
	}
	
	// Parse optional arguments
	branch := ""
	if branchArg, ok := args["branch"].(string); ok {
		branch = branchArg
	}
	
	// Update file
	req := &github.UpdateFileRequest{
		Message: message,
		Content: base64.StdEncoding.EncodeToString([]byte(content)),
		SHA:     sha,
		Branch:  branch,
	}
	
	err := s.client.UpdateFile(ctx, owner, repo, path, req)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to update file: %v", err)),
			},
		}, nil
	}
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(fmt.Sprintf("File '%s' updated successfully in repository '%s/%s'", path, owner, repo)),
		},
	}, nil
}

// handleDeleteFile handles the delete_file tool
func (s *Server) handleDeleteFile(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	path, ok := args["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("path is required and must be a string")
	}
	
	message, ok := args["message"].(string)
	if !ok || message == "" {
		return nil, fmt.Errorf("commit message is required and must be a string")
	}
	
	sha, ok := args["sha"].(string)
	if !ok || sha == "" {
		return nil, fmt.Errorf("sha is required and must be a string")
	}
	
	// Parse optional arguments
	branch := ""
	if branchArg, ok := args["branch"].(string); ok {
		branch = branchArg
	}
	
	// Delete file
	req := &github.DeleteFileRequest{
		Message: message,
		SHA:     sha,
		Branch:  branch,
	}
	
	err := s.client.DeleteFile(ctx, owner, repo, path, req)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to delete file: %v", err)),
			},
		}, nil
	}
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(fmt.Sprintf("File '%s' deleted successfully from repository '%s/%s'", path, owner, repo)),
		},
	}, nil
}

// ============== Issue Tool Handlers ==============

// handleGetIssue handles the get_issue tool
func (s *Server) handleGetIssue(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	// Parse number argument
	var number int
	switch n := args["number"].(type) {
	case float64:
		number = int(n)
	case int:
		number = n
	case string:
		var err error
		number, err = strconv.Atoi(n)
		if err != nil {
			return nil, fmt.Errorf("number must be an integer: %v", err)
		}
	default:
		return nil, fmt.Errorf("number is required and must be an integer")
	}
	
	// Get issue
	issue, err := s.client.GetIssue(ctx, owner, repo, number)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to get issue: %v", err)),
			},
		}, nil
	}
	
	// Format result
	result := fmt.Sprintf(`{
	"number": %d,
	"title": "%s",
	"state": "%s",
	"html_url": "%s",
	"created_at": "%s",
	"updated_at": "%s",
	"body": %s
}`, issue.Number, issue.Title, issue.State, issue.HTMLURL, issue.CreatedAt, issue.UpdatedAt, strconv.Quote(issue.Body))
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// handleListIssues handles the list_issues tool
func (s *Server) handleListIssues(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	// Parse optional arguments
	state := "open"
	if stateArg, ok := args["state"].(string); ok && (stateArg == "open" || stateArg == "closed" || stateArg == "all") {
		state = stateArg
	}
	
	page := 1
	if pageArg, ok := args["page"].(float64); ok {
		page = int(pageArg)
	}
	
	perPage := 30
	if perPageArg, ok := args["per_page"].(float64); ok {
		perPage = int(perPageArg)
	}
	
	// Parse labels argument
	var labels []string
	if labelsArg, ok := args["labels"]; ok {
		switch l := labelsArg.(type) {
		case string:
			if l != "" {
				labels = []string{l}
			}
		case []interface{}:
			for _, item := range l {
				if str, ok := item.(string); ok && str != "" {
					labels = append(labels, str)
				}
			}
		}
	}
	
	// Create options
	opts := &github.ListIssuesOptions{
		State:   state,
		Labels:  labels,
		Page:    page,
		PerPage: perPage,
	}
	
	// List issues
	issues, err := s.client.ListIssues(ctx, owner, repo, opts)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to list issues: %v", err)),
			},
		}, nil
	}
	
	// Format result
	result := "[\n"
	for i, issue := range issues {
		result += fmt.Sprintf(`  {
    "number": %d,
    "title": "%s",
    "state": "%s",
    "html_url": "%s",
    "created_at": "%s"
  }`, issue.Number, issue.Title, issue.State, issue.HTMLURL, issue.CreatedAt)
		
		if i < len(issues)-1 {
			result += ",\n"
		} else {
			result += "\n"
		}
	}
	result += "]"
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// handleCreateIssue handles the create_issue tool
func (s *Server) handleCreateIssue(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	title, ok := args["title"].(string)
	if !ok || title == "" {
		return nil, fmt.Errorf("title is required and must be a string")
	}
	
	// Parse optional arguments
	body := ""
	if bodyArg, ok := args["body"].(string); ok {
		body = bodyArg
	}
	
	// Parse assignees argument
	var assignees []string
	if assigneesArg, ok := args["assignees"]; ok {
		switch a := assigneesArg.(type) {
		case string:
			if a != "" {
				assignees = []string{a}
			}
		case []interface{}:
			for _, item := range a {
				if str, ok := item.(string); ok && str != "" {
					assignees = append(assignees, str)
				}
			}
		}
	}
	
	// Parse labels argument
	var labels []string
	if labelsArg, ok := args["labels"]; ok {
		switch l := labelsArg.(type) {
		case string:
			if l != "" {
				labels = []string{l}
			}
		case []interface{}:
			for _, item := range l {
				if str, ok := item.(string); ok && str != "" {
					labels = append(labels, str)
				}
			}
		}
	}
	
	// Create request
	req := &github.CreateIssueRequest{
		Title:     title,
		Body:      body,
		Assignees: assignees,
		Labels:    labels,
	}
	
	// Create issue
	issue, err := s.client.CreateIssue(ctx, owner, repo, req)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to create issue: %v", err)),
			},
		}, nil
	}
	
	// Format result
	result := fmt.Sprintf(`{
	"number": %d,
	"title": "%s",
	"html_url": "%s"
}`, issue.Number, issue.Title, issue.HTMLURL)
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// handleCloseIssue handles the close_issue tool
func (s *Server) handleCloseIssue(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// This would be implemented to close an issue
	// Not implemented in this example
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.ErrorContent("Not implemented"),
		},
	}, nil
}

// ============== Pull Request Tool Handlers ==============

// handleGetPullRequest handles the get_pull_request tool
func (s *Server) handleGetPullRequest(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	// Parse number argument
	var number int
	switch n := args["number"].(type) {
	case float64:
		number = int(n)
	case int:
		number = n
	case string:
		var err error
		number, err = strconv.Atoi(n)
		if err != nil {
			return nil, fmt.Errorf("number must be an integer: %v", err)
		}
	default:
		return nil, fmt.Errorf("number is required and must be an integer")
	}
	
	// Get pull request
	pr, err := s.client.GetPullRequest(ctx, owner, repo, number)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to get pull request: %v", err)),
			},
		}, nil
	}
	
	// Format result
	result := fmt.Sprintf(`{
	"number": %d,
	"title": "%s",
	"state": "%s",
	"html_url": "%s",
	"created_at": "%s",
	"updated_at": "%s",
	"merged_at": "%s",
	"head": {
		"ref": "%s",
		"sha": "%s",
		"repo": {
			"name": "%s",
			"full_name": "%s"
		}
	},
	"base": {
		"ref": "%s",
		"sha": "%s",
		"repo": {
			"name": "%s",
			"full_name": "%s"
		}
	},
	"body": %s
}`, pr.Number, pr.Title, pr.State, pr.HTMLURL, pr.CreatedAt, pr.UpdatedAt, pr.MergedAt,
		pr.Head.Ref, pr.Head.SHA, pr.Head.Repo.Name, pr.Head.Repo.FullName,
		pr.Base.Ref, pr.Base.SHA, pr.Base.Repo.Name, pr.Base.Repo.FullName,
		strconv.Quote(pr.Body))
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// handleListPullRequests handles the list_pull_requests tool
func (s *Server) handleListPullRequests(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	// Parse optional arguments
	state := "open"
	if stateArg, ok := args["state"].(string); ok && (stateArg == "open" || stateArg == "closed" || stateArg == "all") {
		state = stateArg
	}
	
	page := 1
	if pageArg, ok := args["page"].(float64); ok {
		page = int(pageArg)
	}
	
	perPage := 30
	if perPageArg, ok := args["per_page"].(float64); ok {
		perPage = int(perPageArg)
	}
	
	head := ""
	if headArg, ok := args["head"].(string); ok {
		head = headArg
	}
	
	base := ""
	if baseArg, ok := args["base"].(string); ok {
		base = baseArg
	}
	
	// Create options
	opts := &github.ListPullRequestsOptions{
		State:     state,
		Head:      head,
		Base:      base,
		Page:      page,
		PerPage:   perPage,
	}
	
	// List pull requests
	prs, err := s.client.ListPullRequests(ctx, owner, repo, opts)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to list pull requests: %v", err)),
			},
		}, nil
	}
	
	// Format result
	result := "[\n"
	for i, pr := range prs {
		result += fmt.Sprintf(`  {
    "number": %d,
    "title": "%s",
    "state": "%s",
    "html_url": "%s",
    "created_at": "%s",
    "head": {
      "ref": "%s"
    },
    "base": {
      "ref": "%s"
    }
  }`, pr.Number, pr.Title, pr.State, pr.HTMLURL, pr.CreatedAt, pr.Head.Ref, pr.Base.Ref)
		
		if i < len(prs)-1 {
			result += ",\n"
		} else {
			result += "\n"
		}
	}
	result += "]"
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// handleCreatePullRequest handles the create_pull_request tool
func (s *Server) handleCreatePullRequest(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	title, ok := args["title"].(string)
	if !ok || title == "" {
		return nil, fmt.Errorf("title is required and must be a string")
	}
	
	head, ok := args["head"].(string)
	if !ok || head == "" {
		return nil, fmt.Errorf("head branch is required and must be a string")
	}
	
	base, ok := args["base"].(string)
	if !ok || base == "" {
		return nil, fmt.Errorf("base branch is required and must be a string")
	}
	
	// Parse optional arguments
	body := ""
	if bodyArg, ok := args["body"].(string); ok {
		body = bodyArg
	}
	
	draft := false
	if draftArg, ok := args["draft"].(bool); ok {
		draft = draftArg
	}
	
	// Create request
	req := &github.CreatePullRequestRequest{
		Title: title,
		Body:  body,
		Head:  head,
		Base:  base,
		Draft: draft,
	}
	
	// Create pull request
	pr, err := s.client.CreatePullRequest(ctx, owner, repo, req)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to create pull request: %v", err)),
			},
		}, nil
	}
	
	// Format result
	result := fmt.Sprintf(`{
	"number": %d,
	"title": "%s",
	"html_url": "%s"
}`, pr.Number, pr.Title, pr.HTMLURL)
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// handleMergePullRequest handles the merge_pull_request tool
func (s *Server) handleMergePullRequest(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// This would be implemented to merge a pull request
	// Not implemented in this example
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.ErrorContent("Not implemented"),
		},
	}, nil
}

// ============== GitHub Actions Tool Handlers ==============

// handleListWorkflows handles the list_workflows tool
func (s *Server) handleListWorkflows(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	// Parse optional arguments
	page := 1
	if pageArg, ok := args["page"].(float64); ok {
		page = int(pageArg)
	}
	
	perPage := 30
	if perPageArg, ok := args["per_page"].(float64); ok {
		perPage = int(perPageArg)
	}
	
	// List workflows
	workflows, err := s.client.ListWorkflows(ctx, owner, repo, page, perPage)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to list workflows: %v", err)),
			},
		}, nil
	}
	
	// Format result
	result := "[\n"
	for i, workflow := range workflows {
		result += fmt.Sprintf(`  {
    "id": %d,
    "name": "%s",
    "path": "%s",
    "state": "%s"
  }`, workflow.ID, workflow.Name, workflow.Path, workflow.State)
		
		if i < len(workflows)-1 {
			result += ",\n"
		} else {
			result += "\n"
		}
	}
	result += "]"
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// handleListWorkflowRuns handles the list_workflow_runs tool
func (s *Server) handleListWorkflowRuns(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	// Parse workflow_id argument
	var workflowID int64
	switch w := args["workflow_id"].(type) {
	case float64:
		workflowID = int64(w)
	case int64:
		workflowID = w
	case int:
		workflowID = int64(w)
	case string:
		var err error
		id, err := strconv.ParseInt(w, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("workflow_id must be an integer: %v", err)
		}
		workflowID = id
	default:
		return nil, fmt.Errorf("workflow_id is required and must be an integer")
	}
	
	// Parse optional arguments
	page := 1
	if pageArg, ok := args["page"].(float64); ok {
		page = int(pageArg)
	}
	
	perPage := 30
	if perPageArg, ok := args["per_page"].(float64); ok {
		perPage = int(perPageArg)
	}
	
	// List workflow runs
	runs, err := s.client.ListWorkflowRuns(ctx, owner, repo, workflowID, page, perPage)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to list workflow runs: %v", err)),
			},
		}, nil
	}
	
	// Format result
	result := "[\n"
	for i, run := range runs {
		result += fmt.Sprintf(`  {
    "id": %d,
    "name": "%s",
    "head_branch": "%s",
    "head_sha": "%s",
    "status": "%s",
    "conclusion": "%s",
    "html_url": "%s",
    "created_at": "%s",
    "updated_at": "%s"
  }`, run.ID, run.Name, run.HeadBranch, run.HeadSHA, run.Status, run.Conclusion, run.HTMLURL, run.CreatedAt, run.UpdatedAt)
		
		if i < len(runs)-1 {
			result += ",\n"
		} else {
			result += "\n"
		}
	}
	result += "]"
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result),
		},
	}, nil
}

// handleTriggerWorkflow handles the trigger_workflow tool
func (s *Server) handleTriggerWorkflow(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	owner, ok := args["owner"].(string)
	if !ok || owner == "" {
		return nil, fmt.Errorf("owner is required and must be a string")
	}
	
	repo, ok := args["repo"].(string)
	if !ok || repo == "" {
		return nil, fmt.Errorf("repo is required and must be a string")
	}
	
	// Parse workflow_id argument
	var workflowID int64
	switch w := args["workflow_id"].(type) {
	case float64:
		workflowID = int64(w)
	case int64:
		workflowID = w
	case int:
		workflowID = int64(w)
	case string:
		var err error
		id, err := strconv.ParseInt(w, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("workflow_id must be an integer: %v", err)
		}
		workflowID = id
	default:
		return nil, fmt.Errorf("workflow_id is required and must be an integer")
	}
	
	ref, ok := args["ref"].(string)
	if !ok || ref == "" {
		return nil, fmt.Errorf("ref is required and must be a string")
	}
	
	// Parse inputs argument
	inputs := make(map[string]interface{})
	if inputsArg, ok := args["inputs"].(map[string]interface{}); ok {
		inputs = inputsArg
	}
	
	// Create request
	req := &github.TriggerWorkflowRequest{
		Ref:    ref,
		Inputs: inputs,
	}
	
	// Trigger workflow
	err := s.client.TriggerWorkflow(ctx, owner, repo, workflowID, req)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to trigger workflow: %v", err)),
			},
		}, nil
	}
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(fmt.Sprintf("Workflow '%d' triggered successfully in repository '%s/%s' on ref '%s'", workflowID, owner, repo, ref)),
		},
	}, nil
}

// ============== Search Tool Handlers ==============

// handleSearchCode handles the search_code tool
func (s *Server) handleSearchCode(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query is required and must be a string")
	}
	
	// Parse optional arguments
	page := 1
	if pageArg, ok := args["page"].(float64); ok {
		page = int(pageArg)
	}
	
	perPage := 30
	if perPageArg, ok := args["per_page"].(float64); ok {
		perPage = int(perPageArg)
	}
	
	// Search code
	result, err := s.client.SearchCode(ctx, query, page, perPage)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to search code: %v", err)),
			},
		}, nil
	}
	
	// Format result
	response := fmt.Sprintf(`{
	"total_count": %d,
	"items": [`, result.TotalCount)
	
	for i, item := range result.Items {
		response += fmt.Sprintf(`
		{
			"name": "%s",
			"path": "%s",
			"html_url": "%s",
			"repository": {
				"name": "%s",
				"full_name": "%s"
			}
		}`, item.Name, item.Path, item.HTMLURL, item.Repository.Name, item.Repository.FullName)
		
		if i < len(result.Items)-1 {
			response += ","
		}
	}
	
	response += `
	]
}`
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(response),
		},
	}, nil
}

// handleSearchIssues handles the search_issues tool
func (s *Server) handleSearchIssues(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	// Validate arguments
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query is required and must be a string")
	}
	
	// Parse optional arguments
	page := 1
	if pageArg, ok := args["page"].(float64); ok {
		page = int(pageArg)
	}
	
	perPage := 30
	if perPageArg, ok := args["per_page"].(float64); ok {
		perPage = int(perPageArg)
	}
	
	// Search issues
	result, err := s.client.SearchIssues(ctx, query, page, perPage)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to search issues: %v", err)),
			},
		}, nil
	}
	
	// Format result
	response := fmt.Sprintf(`{
	"total_count": %d,
	"items": [`, result.TotalCount)
	
	for i, item := range result.Items {
		response += fmt.Sprintf(`
		{
			"number": %d,
			"title": "%s",
			"state": "%s",
			"html_url": "%s",
			"created_at": "%s"
		}`, item.Number, item.Title, item.State, item.HTMLURL, item.CreatedAt)
		
		if i < len(result.Items)-1 {
			response += ","
		}
	}
	
	response += `
	]
}`
	
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(response),
		},
	}, nil
}