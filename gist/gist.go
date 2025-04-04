package gist

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Gist represents a GitHub gist
type Gist struct {
	ID          string              `json:"id"`
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]GistFile `json:"files"`
	Owner       string              `json:"owner"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	HTMLURL     string              `json:"html_url"`
}

// GistFile represents a file within a gist
type GistFile struct {
	Content  string `json:"content"`
	Language string `json:"language"`
	Size     int    `json:"size"`
	RawURL   string `json:"raw_url"`
}

// GistManager handles gist operations
type GistManager struct {
	client GistClient
}

// GistClient defines the interface for gist operations
type GistClient interface {
	Create(description string, files map[string]GistFile, public bool) (*Gist, error)
	List() ([]*Gist, error)
	Get(id string) (*Gist, error)
	Update(id string, description string, files map[string]GistFile) error
	Delete(id string) error
}

// NewGistManager creates a new GistManager instance
func NewGistManager(client GistClient) *GistManager {
	return &GistManager{
		client: client,
	}
}

// Create creates a new gist
func (m *GistManager) Create(description string, files map[string]string, public bool) (*Gist, error) {
	// Convert file map to GistFile map
	gistFiles := make(map[string]GistFile)
	for name, content := range files {
		gistFiles[name] = GistFile{
			Content: content,
		}
	}

	return m.client.Create(description, gistFiles, public)
}

// List returns all gists for the authenticated user
func (m *GistManager) List() ([]*Gist, error) {
	return m.client.List()
}

// Get retrieves a gist by ID
func (m *GistManager) Get(id string) (*Gist, error) {
	return m.client.Get(id)
}

// Update updates an existing gist
func (m *GistManager) Update(id string, description string, files map[string]string) error {
	gistFiles := make(map[string]GistFile)
	for name, content := range files {
		gistFiles[name] = GistFile{
			Content: content,
		}
	}

	return m.client.Update(id, description, gistFiles)
}

// Delete removes a gist
func (m *GistManager) Delete(id string) error {
	return m.client.Delete(id)
}

// Clone clones a gist to a local directory
func (m *GistManager) Clone(id string, directory string) error {
	// Get the gist
	gist, err := m.Get(id)
	if err != nil {
		return fmt.Errorf("failed to get gist: %w", err)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write each file
	for name, file := range gist.Files {
		filePath := filepath.Join(directory, name)
		if err := os.WriteFile(filePath, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", name, err)
		}
	}

	return nil
}

// CreateFromFiles creates a gist from local files
func (m *GistManager) CreateFromFiles(description string, public bool, filePaths ...string) (*Gist, error) {
	files := make(map[string]string)

	for _, path := range filePaths {
		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Use base filename
		name := filepath.Base(path)
		files[name] = string(content)
	}

	return m.Create(description, files, public)
}
