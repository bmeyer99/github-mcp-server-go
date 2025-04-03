// github-mcp-server-go/main.go
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/your-username/github-mcp-server-go/server"
	"github.com/your-username/github-mcp-server-go/transport"
)

func main() {
	// Parse command line flags
	tokenFlag := flag.String("token", "", "GitHub Personal Access Token")
	debugFlag := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Check for token in environment variable if not provided via flag
	token := *tokenFlag
	if token == "" {
		token = os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
		if token == "" {
			log.Fatal("GitHub token must be provided either via -token flag or GITHUB_PERSONAL_ACCESS_TOKEN environment variable")
		}
	}

	// Setup logger
	logger := log.New(os.Stderr, "", log.LstdFlags)
	if *debugFlag {
		logger.Println("Debug logging enabled")
	}

	// Create server
	logger.Println("Initializing GitHub MCP server")
	srv := server.New(server.Config{
		Token:  token,
		Logger: logger,
		Debug:  *debugFlag,
	})

	// Create transport
	stdioTransport := transport.NewStdioTransport()

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Println("Received termination signal, shutting down...")
		cancel()
	}()

	// Start the server
	logger.Println("Starting GitHub MCP server")
	if err := srv.Serve(ctx, stdioTransport); err != nil {
		logger.Fatalf("Server error: %v", err)
	}

	logger.Println("Server shut down successfully")
}