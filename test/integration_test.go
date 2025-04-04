package test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGitHubMCPServer(t *testing.T) {
	// Check for required environment variables
	token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")
	if token == "" || username == "" {
		t.Skip("GITHUB_PERSONAL_ACCESS_TOKEN and GITHUB_USERNAME must be set")
	}

	// Create MCP client
	client, err := NewMCPClient(token)
	if err != nil {
		t.Fatalf("failed to create MCP client: %v", err)
	}
	defer client.Close()

	// Generate unique repository name
	repoName := fmt.Sprintf("test-repo-%d", time.Now().Unix())

	// Test cases
	tests := []struct {
		name   string
		method string
		params interface{}
		check  func(t *testing.T, result json.RawMessage)
	}{
		{
			name:   "Initialize",
			method: "initialize",
			params: map[string]interface{}{
				"protocolVersion": "1.0",
			},
			check: func(t *testing.T, result json.RawMessage) {
				var resp struct {
					ProtocolVersion string `json:"protocolVersion"`
				}
				if err := json.Unmarshal(result, &resp); err != nil {
					t.Errorf("failed to unmarshal result: %v", err)
					return
				}
				if resp.ProtocolVersion != "1.0" {
					t.Errorf("expected protocol version 1.0, got %s", resp.ProtocolVersion)
				}
			},
		},
		{
			name:   "Create Repository",
			method: "tools/call",
			params: map[string]interface{}{
				"name": "create_repository",
				"arguments": map[string]interface{}{
					"name":        repoName,
					"description": "Test repository for MCP server",
					"private":     true,
				},
			},
		},
		{
			name:   "Create File",
			method: "tools/call",
			params: map[string]interface{}{
				"name": "create_file",
				"arguments": map[string]interface{}{
					"owner":   username,
					"repo":    repoName,
					"path":    "README.md",
					"content": "# Test Repository\nThis is a test repository created by the GitHub MCP server.",
					"message": "Initial commit",
				},
			},
		},
		{
			name:   "Create Issue",
			method: "tools/call",
			params: map[string]interface{}{
				"name": "create_issue",
				"arguments": map[string]interface{}{
					"owner": username,
					"repo":  repoName,
					"title": "Test Issue",
					"body":  "This is a test issue created by the GitHub MCP server.",
				},
			},
		},
		{
			name:   "List Issues",
			method: "tools/call",
			params: map[string]interface{}{
				"name": "list_issues",
				"arguments": map[string]interface{}{
					"owner": username,
					"repo":  repoName,
				},
			},
			check: func(t *testing.T, result json.RawMessage) {
				// First unmarshal into a json.RawMessage to examine the structure
				var raw interface{}
				if err := json.Unmarshal(result, &raw); err != nil {
					t.Errorf("failed to unmarshal raw result: %v", err)
					return
				}

				// Log the actual response structure
				t.Logf("Raw response: %+v", raw)

				// Parse the content array
				var content struct {
					Content []struct {
						Text string `json:"text"`
					} `json:"content"`
				}
				if err := json.Unmarshal(result, &content); err != nil {
					t.Errorf("failed to unmarshal content: %v", err)
					return
				}

				// Parse the JSON array from the text content
				text := strings.TrimSpace(content.Content[0].Text)
				var issues []struct {
					Number    int    `json:"number"`
					Title     string `json:"title"`
					State     string `json:"state"`
					HTMLURL   string `json:"html_url"`
					CreatedAt string `json:"created_at"`
				}
				if err := json.Unmarshal([]byte(text), &issues); err != nil {
					t.Errorf("failed to unmarshal issues: %v, text was: %s", err, text)
					return
				}

				if len(issues) != 1 {
					t.Errorf("expected 1 issue, got %d", len(issues))
					return
				}

				if issues[0].Title != "Test Issue" || issues[0].State != "open" {
					t.Errorf("expected title 'Test Issue' and state 'open', got title '%s' and state '%s'",
						issues[0].Title, issues[0].State)
					return
				}
			},
		},
		{
			name:   "Search Code",
			method: "tools/call",
			params: map[string]interface{}{
				"name": "search_code",
				"arguments": map[string]interface{}{
					"query": "README " + username + "/" + repoName,
				},
			},
			check: func(t *testing.T, result json.RawMessage) {
				// Log the search response
				t.Logf("Search response: %s", string(result))

				var content struct {
					Content []struct {
						Text string `json:"text"`
					} `json:"content"`
				}
				if err := json.Unmarshal(result, &content); err != nil {
					t.Errorf("failed to unmarshal content: %v", err)
					return
				}

				if len(content.Content) == 0 {
					t.Error("no content returned from search")
					return
				}

				// Parse the search results
				text := strings.TrimSpace(content.Content[0].Text)
				if strings.Contains(text, "total_count\": 0") {
					t.Skipf("No search results found yet (likely still indexing): %s", text)
				} else if !strings.Contains(text, "README.md") {
					t.Errorf("expected to find README.md in search results, got: %s", text)
				}
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add delay before search test to allow GitHub to index content
			if tt.name == "Search Code" {
				t.Log("Waiting 30 seconds for GitHub to index repository content...")
				time.Sleep(30 * time.Second)
				t.Log("Note: If search fails, it may be because GitHub is still indexing the new content")
			}
			resp, err := client.Call(tt.method, tt.params)
			if err != nil {
				t.Fatalf("call failed: %v", err)
			}

			if tt.check != nil {
				tt.check(t, resp.Result)
			}
		})
	}
}
