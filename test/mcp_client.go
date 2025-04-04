package test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// MCPRequest represents a JSON-RPC request to the MCP server
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// MCPResponse represents a JSON-RPC response from the MCP server
type MCPResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
}

// MCPClient represents a client for communicating with the MCP server
type MCPClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
}

// NewMCPClient creates a new MCP client
func NewMCPClient(token string) (*MCPClient, error) {
	// Find the binary in the project root
	binPath := "../github-mcp-server-go"
	cmd := exec.Command(binPath, "-debug")
	cmd.Env = append(os.Environ(), fmt.Sprintf("GITHUB_PERSONAL_ACCESS_TOKEN=%s", token))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start MCP server: %w", err)
	}

	return &MCPClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
	}, nil
}

// Call makes a request to the MCP server and returns the response
func (c *MCPClient) Call(method string, params interface{}) (*MCPResponse, error) {
	// Use 10 second timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create channels for response and error
	respCh := make(chan *MCPResponse, 1)
	errCh := make(chan error, 1)

	go func() {
		req := MCPRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  method,
			Params:  params,
		}

		reqBytes, err := json.Marshal(req)
		if err != nil {
			errCh <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		reqBytes = append(reqBytes, '\n')
		if _, err := c.stdin.Write(reqBytes); err != nil {
			errCh <- fmt.Errorf("failed to write request: %w", err)
			return
		}

		respBytes, err := c.stdout.ReadBytes('\n')
		if err != nil {
			errCh <- fmt.Errorf("failed to read response: %w", err)
			return
		}

		for bytes.HasPrefix(respBytes, []byte("2025/")) {
			respBytes, err = c.stdout.ReadBytes('\n')
			if err != nil {
				errCh <- fmt.Errorf("failed to read response: %w", err)
				return
			}
		}

		var resp MCPResponse
		if err := json.Unmarshal(respBytes, &resp); err != nil {
			errCh <- fmt.Errorf("failed to unmarshal response: %w", err)
			return
		}

		respCh <- &resp
	}()

	// Wait for either response, error, or timeout
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("request timed out after 10 seconds")
	case err := <-errCh:
		return nil, err
	case resp := <-respCh:
		return resp, nil
	}
}

// Close closes the MCP client and terminates the server
func (c *MCPClient) Close() error {
	if err := c.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill MCP server: %w", err)
	}
	return nil
}
