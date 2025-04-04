package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupTestStore(t *testing.T) (*FileStore, func()) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	store, err := NewFileStore(filepath.Join(tmpDir, "config.json"))
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create FileStore: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return store, cleanup
}

func TestFileStore_SetGet(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Test string value
	if err := store.Set("string_key", "test value"); err != nil {
		t.Errorf("Failed to set string value: %v", err)
	}

	value, err := store.Get("string_key")
	if err != nil {
		t.Errorf("Failed to get string value: %v", err)
	}
	if str, ok := value.(string); !ok || str != "test value" {
		t.Errorf("Expected string value 'test value', got %v", value)
	}

	// Test number value
	if err := store.Set("number_key", 42); err != nil {
		t.Errorf("Failed to set number value: %v", err)
	}

	value, err = store.Get("number_key")
	if err != nil {
		t.Errorf("Failed to get number value: %v", err)
	}
	if num, ok := value.(float64); !ok || num != 42 {
		t.Errorf("Expected number value 42, got %v", value)
	}

	// Test map value
	mapVal := map[string]interface{}{
		"nested": "value",
		"number": 123,
	}
	if err := store.Set("map_key", mapVal); err != nil {
		t.Errorf("Failed to set map value: %v", err)
	}

	value, err = store.Get("map_key")
	if err != nil {
		t.Errorf("Failed to get map value: %v", err)
	}
	if m, ok := value.(map[string]interface{}); !ok {
		t.Errorf("Expected map value, got %T", value)
	} else if m["nested"] != "value" || m["number"].(float64) != 123 {
		t.Errorf("Map value mismatch, got %v", m)
	}
}

func TestFileStore_Delete(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Set and then delete a value
	if err := store.Set("test_key", "test value"); err != nil {
		t.Errorf("Failed to set value: %v", err)
	}

	if err := store.Delete("test_key"); err != nil {
		t.Errorf("Failed to delete value: %v", err)
	}

	// Verify value is deleted
	if _, err := store.Get("test_key"); err == nil {
		t.Error("Expected error getting deleted key")
	}

	// Try to delete non-existent key
	if err := store.Delete("nonexistent"); err == nil {
		t.Error("Expected error deleting non-existent key")
	}
}

func TestFileStore_List(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Set multiple values
	testData := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": map[string]interface{}{"nested": "value"},
	}

	for k, v := range testData {
		if err := store.Set(k, v); err != nil {
			t.Errorf("Failed to set value for %s: %v", k, err)
		}
	}

	// List all values
	values, err := store.List()
	if err != nil {
		t.Errorf("Failed to list values: %v", err)
	}

	// Verify all values are present
	if len(values) != len(testData) {
		t.Errorf("Expected %d values, got %d", len(testData), len(values))
	}

	for k, v := range testData {
		if value, ok := values[k]; !ok {
			t.Errorf("Missing key %s in listed values", k)
		} else {
			// Convert both values to JSON for comparison
			expectedJSON, _ := json.Marshal(v)
			actualJSON, _ := json.Marshal(value)
			if string(expectedJSON) != string(actualJSON) {
				t.Errorf("Value mismatch for key %s: expected %v, got %v", k, string(expectedJSON), string(actualJSON))
			}
		}
	}
}

func TestFileStore_Persistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config-persistence-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")

	// Create store and set some values
	store, err := NewFileStore(configPath)
	if err != nil {
		t.Fatalf("Failed to create first store: %v", err)
	}

	if err := store.Set("test_key", "test value"); err != nil {
		t.Errorf("Failed to set initial value: %v", err)
	}

	// Create new store instance and verify values persist
	store2, err := NewFileStore(configPath)
	if err != nil {
		t.Fatalf("Failed to create second store: %v", err)
	}

	value, err := store2.Get("test_key")
	if err != nil {
		t.Errorf("Failed to get value from second store: %v", err)
	}
	if str, ok := value.(string); !ok || str != "test value" {
		t.Errorf("Expected value 'test value', got %v", value)
	}
}
