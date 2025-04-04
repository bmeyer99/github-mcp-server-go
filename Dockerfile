# github-mcp-server-go/Dockerfile
# Build stage
FROM golang:1.21-alpine AS build

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Fix import paths if needed (for Docker build context)
RUN grep -r "github.com/your-username/github-mcp-server-go" --include="*.go" . | \
    xargs -r sed -i 's|github.com/your-username/github-mcp-server-go|github-mcp-server-go|g'

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o github-mcp-server .

# Final stage
FROM alpine:3.19

WORKDIR /root/

# Copy the binary from the build stage
COPY --from=build /app/github-mcp-server .

# Set entrypoint
ENTRYPOINT ["./github-mcp-server"]

# Document that the container expects the GitHub token as an environment variable
ENV GITHUB_PERSONAL_ACCESS_TOKEN=""

# Document that the container exposes stdio for MCP communication
LABEL org.opencontainers.image.description="GitHub MCP Server for comprehensive GitHub CLI integration"
LABEL org.opencontainers.image.source="https://github.com/your-username/github-mcp-server-go"
LABEL org.opencontainers.image.licenses="MIT"