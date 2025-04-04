package gist

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type mockGistClient struct {
	gists map[string]*Gist
}

func newMockGistClient() *mockGistClient {
	return &mockGistClient{
		gists: make(map[string]*Gist),
	}
}

func (m *mockGistClient) Create(description string, files map[string]GistFile, public bool) (*Gist, error) {
	gist := &Gist{
		ID:          "gist123",
		Description: description,
		Public:      public,
		Files:       files,
		Owner:       "testuser",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		HTMLURL:     "https://gist.github.com/gist123",
	}
	m.gists[gist.ID] = gist
	return gist, nil
}

func (m *mockGistClient) List() ([]*Gist, error) {
	gists := make([]*Gist, 0, len(m.gists))
	for _, gist := range m.gists {
		gists = append(gists, gist)
	}
	return gists, nil
}

func (m *mockGistClient) Get(id string) (*Gist, error) {
	if gist, ok := m.gists[id]; ok {
		return gist, nil
	}
	return nil, fmt.Errorf("gist not found: %s", id)
}

func (m *mockGistClient) Update(id string, description string, files map[string]GistFile) error {
	if gist, ok := m.gists[id]; ok {
		gist.Description = description
		gist.Files = files
		gist.UpdatedAt = time.Now().UTC()
		return nil
	}
	return fmt.Errorf("gist not found: %s", id)
}

func (m *mockGistClient) Delete(id string) error {
	if _, ok := m.gists[id]; ok {
		delete(m.gists, id)
		return nil
	}
	return fmt.Errorf("gist not found: %s", id)
}

func TestGistManager_Create(t *testing.T) {
	client := newMockGistClient()
	manager := NewGistManager(client)

	files := map[string]string{
		"test.txt":  "test content",
		"test2.txt": "more content",
	}

	gist, err := manager.Create("Test Gist", files, true)
	if err != nil {
		t.Fatalf("Failed to create gist: %v", err)
	}

	if gist.Description != "Test Gist" {
		t.Errorf("Expected description 'Test Gist', got %s", gist.Description)
	}
	if !gist.Public {
		t.Error("Expected public gist")
	}
	if len(gist.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(gist.Files))
	}
}

func TestGistManager_List(t *testing.T) {
	client := newMockGistClient()
	manager := NewGistManager(client)

	// Create some test gists
	files := map[string]string{"test.txt": "content"}
	manager.Create("Gist 1", files, true)
	manager.Create("Gist 2", files, false)

	gists, err := manager.List()
	if err != nil {
		t.Fatalf("Failed to list gists: %v", err)
	}

	if len(gists) != 2 {
		t.Errorf("Expected 2 gists, got %d", len(gists))
	}
}

func TestGistManager_Clone(t *testing.T) {
	client := newMockGistClient()
	manager := NewGistManager(client)

	// Create a test gist
	files := map[string]string{
		"test.txt":  "test content",
		"test2.txt": "more content",
	}
	gist, err := manager.Create("Test Gist", files, true)
	if err != nil {
		t.Fatalf("Failed to create gist: %v", err)
	}

	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "gist-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone the gist
	err = manager.Clone(gist.ID, tmpDir)
	if err != nil {
		t.Fatalf("Failed to clone gist: %v", err)
	}

	// Verify files were created
	for name, file := range gist.Files {
		path := filepath.Join(tmpDir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", name, err)
			continue
		}
		if string(content) != file.Content {
			t.Errorf("File %s content mismatch: expected %q, got %q", name, file.Content, string(content))
		}
	}
}

func TestGistManager_CreateFromFiles(t *testing.T) {
	client := newMockGistClient()
	manager := NewGistManager(client)

	// Create temporary directory and files
	tmpDir, err := os.MkdirTemp("", "gist-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := map[string]string{
		"test1.txt": "test content 1",
		"test2.txt": "test content 2",
	}

	var filePaths []string
	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
		filePaths = append(filePaths, path)
	}

	// Create gist from files
	gist, err := manager.CreateFromFiles("Test Gist", true, filePaths...)
	if err != nil {
		t.Fatalf("Failed to create gist from files: %v", err)
	}

	// Verify gist contents
	if gist.Description != "Test Gist" {
		t.Errorf("Expected description 'Test Gist', got %s", gist.Description)
	}
	if !gist.Public {
		t.Error("Expected public gist")
	}
	if len(gist.Files) != len(files) {
		t.Errorf("Expected %d files, got %d", len(files), len(gist.Files))
	}

	for name, expectedContent := range files {
		if file, ok := gist.Files[filepath.Base(name)]; !ok {
			t.Errorf("Missing file %s in gist", name)
		} else if file.Content != expectedContent {
			t.Errorf("File %s content mismatch: expected %q, got %q", name, expectedContent, file.Content)
		}
	}
}
