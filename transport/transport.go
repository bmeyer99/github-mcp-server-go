// github-mcp-server-go/transport/transport.go
package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

// Transport defines the interface for MCP communication channels
type Transport interface {
	// ReadMessage reads a message from the transport
	ReadMessage(ctx context.Context) ([]byte, error)

	// WriteMessage writes a message to the transport
	WriteMessage(ctx context.Context, message []byte) error

	// Close closes the transport
	Close() error
}

// StdioTransport implements Transport for stdin/stdout communication
type StdioTransport struct {
	reader  *bufio.Reader
	writer  io.Writer
	writeMu sync.Mutex
}

// NewStdioTransport creates a new StdioTransport
func NewStdioTransport() *StdioTransport {
	return &StdioTransport{
		reader: bufio.NewReader(os.Stdin),
		writer: os.Stdout,
	}
}

// ReadMessage reads a JSON-RPC message from stdin
func (t *StdioTransport) ReadMessage(ctx context.Context) ([]byte, error) {
	// Create a channel for the read operation
	messageCh := make(chan []byte)
	errCh := make(chan error)

	// Read in a separate goroutine so we can respect context cancellation
	go func() {
		line, err := t.reader.ReadBytes('\n')
		if err != nil {
			errCh <- fmt.Errorf("failed to read message: %w", err)
			return
		}
		messageCh <- line
	}()

	// Wait for either the message or context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case message := <-messageCh:
		return message, nil
	}
}

// WriteMessage writes a JSON-RPC message to stdout
func (t *StdioTransport) WriteMessage(ctx context.Context, message []byte) error {
	t.writeMu.Lock()
	defer t.writeMu.Unlock()

	// Create a channel for the write operation
	doneCh := make(chan error)

	// Write in a separate goroutine so we can respect context cancellation
	go func() {
		_, err := t.writer.Write(message)
		if err != nil {
			doneCh <- fmt.Errorf("failed to write message: %w", err)
			return
		}

		// Write a newline character to terminate the message
		_, err = t.writer.Write([]byte("\n"))
		if err != nil {
			doneCh <- fmt.Errorf("failed to write newline: %w", err)
			return
		}

		doneCh <- nil
	}()

	// Wait for either the write to complete or context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-doneCh:
		return err
	}
}

// Close closes the transport
func (t *StdioTransport) Close() error {
	// Nothing to close for stdio transport
	return nil
}

// SSETransport implements Transport for Server-Sent Events
// This would be implemented for HTTP-based MCP servers
type SSETransport struct {
	// Fields for SSE transport would go here
}

// NewSSETransport creates a new SSETransport
func NewSSETransport() *SSETransport {
	// Implementation would go here
	return &SSETransport{}
}

// ReadMessage reads a JSON-RPC message from the SSE connection
func (t *SSETransport) ReadMessage(ctx context.Context) ([]byte, error) {
	// Implementation would go here
	return nil, fmt.Errorf("SSE transport not implemented")
}

// WriteMessage writes a JSON-RPC message to the SSE connection
func (t *SSETransport) WriteMessage(ctx context.Context, message []byte) error {
	// Implementation would go here
	return fmt.Errorf("SSE transport not implemented")
}

// Close closes the SSE connection
func (t *SSETransport) Close() error {
	// Implementation would go here
	return fmt.Errorf("SSE transport not implemented")
}
