package config

import (
	"os"
	"path/filepath"
	"testing"
)

// mockEnvProvider implements EnvProvider for testing
type mockEnvProvider struct {
	env map[string]string
}

func newMockEnvProvider() *mockEnvProvider {
	return &mockEnvProvider{
		env: make(map[string]string),
	}
}

func (p *mockEnvProvider) Get(key string) string {
	return p.env[key]
}

func (p *mockEnvProvider) Set(key, value string) error {
	p.env[key] = value
	return nil
}

func (p *mockEnvProvider) List() map[string]string {
	return p.env
}

func setupTestManager(t *testing.T) (*Manager, func()) {
	// Create temporary directories
	tmpDir, err := os.MkdirTemp("", "config-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create test directories
	globalConfigDir := filepath.Join(tmpDir, "global")
	localConfigDir := filepath.Join(tmpDir, "local")

	opts := ManagerOptions{
		ConfigDir:       globalConfigDir,
		LocalConfigPath: filepath.Join(localConfigDir, "config.json"),
		EnvProvider:     newMockEnvProvider(),
	}

	manager, err := NewManager(opts)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create Manager: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return manager, cleanup
}

func TestManager_SetGet(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	tests := []struct {
		name      string
		key       string
		value     interface{}
		scope     Scope
		wantError bool
	}{
		{
			name:      "global string",
			key:       "global_key",
			value:     "global value",
			scope:     ScopeGlobal,
			wantError: false,
		},
		{
			name:      "local string",
			key:       "local_key",
			value:     "local value",
			scope:     ScopeLocal,
			wantError: false,
		},
		{
			name:      "global number",
			key:       "global_num",
			value:     42,
			scope:     ScopeGlobal,
			wantError: false,
		},
		{
			name:      "local map",
			key:       "local_map",
			value:     map[string]interface{}{"nested": "value"},
			scope:     ScopeLocal,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set value
			err := manager.Set(tt.key, tt.value, tt.scope)
			if (err != nil) != tt.wantError {
				t.Errorf("Set() error = %v, wantError %v", err, tt.wantError)
				return
			}

			// Get value
			value, err := manager.Get(tt.key, tt.scope)
			if (err != nil) != tt.wantError {
				t.Errorf("Get() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if err == nil && !tt.wantError {
				// Compare values
				if !compareValues(value, tt.value) {
					t.Errorf("Value mismatch: got %v, want %v", value, tt.value)
				}
			}
		})
	}
}

func TestManager_ScopePrecedence(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	// Set up test values in different scopes
	testKey := "test_key"
	envValue := "env_value"
	globalValue := "global_value"
	localValue := "local_value"

	// Set values in different scopes
	envProvider := manager.env.(*mockEnvProvider)
	envProvider.Set(testKey, envValue)
	manager.Set(testKey, globalValue, ScopeGlobal)
	manager.Set(testKey, localValue, ScopeLocal)

	// Test precedence: env > local > global
	value, err := manager.Get(testKey, ScopeLocal)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if value != envValue {
		t.Errorf("Expected env value %q, got %q", envValue, value)
	}

	// Remove env var and test local precedence
	envProvider.env = make(map[string]string)
	value, err = manager.Get(testKey, ScopeLocal)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if value != localValue {
		t.Errorf("Expected local value %q, got %q", localValue, value)
	}
}

func TestManager_Delete(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	// Set up test values
	testKey := "test_key"
	globalValue := "global_value"
	localValue := "local_value"

	manager.Set(testKey, globalValue, ScopeGlobal)
	manager.Set(testKey, localValue, ScopeLocal)

	// Test delete from local scope
	if err := manager.Delete(testKey, ScopeLocal); err != nil {
		t.Errorf("Delete() local error = %v", err)
	}

	// Verify local value is deleted but global remains
	if value, err := manager.Get(testKey, ScopeLocal); err == nil {
		if value == localValue {
			t.Error("Local value should be deleted")
		}
	}

	// Test delete from global scope
	if err := manager.Delete(testKey, ScopeGlobal); err != nil {
		t.Errorf("Delete() global error = %v", err)
	}

	// Verify global value is deleted
	if _, err := manager.Get(testKey, ScopeGlobal); err == nil {
		t.Error("Global value should be deleted")
	}
}

func TestManager_ListAll(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	// Set up test values in different scopes
	globalValues := map[string]interface{}{
		"global1": "value1",
		"global2": 42,
	}

	localValues := map[string]interface{}{
		"local1": "value2",
		"local2": true,
	}

	envValues := map[string]string{
		"env1": "value3",
		"env2": "value4",
	}

	// Set values
	for k, v := range globalValues {
		manager.Set(k, v, ScopeGlobal)
	}

	for k, v := range localValues {
		manager.Set(k, v, ScopeLocal)
	}

	envProvider := manager.env.(*mockEnvProvider)
	for k, v := range envValues {
		envProvider.Set(k, v)
	}

	// Get all values
	allValues := manager.ListAll()

	// Verify all values are present
	for k, v := range globalValues {
		if !compareValues(allValues[k], v) {
			t.Errorf("Missing or incorrect global value for key %s", k)
		}
	}

	for k, v := range localValues {
		if !compareValues(allValues[k], v) {
			t.Errorf("Missing or incorrect local value for key %s", k)
		}
	}

	for k, v := range envValues {
		if !compareValues(allValues[k], v) {
			t.Errorf("Missing or incorrect env value for key %s", k)
		}
	}
}

// Helper function to compare values that might be of different types
func compareValues(a, b interface{}) bool {
	switch v := b.(type) {
	case string:
		if s, ok := a.(string); ok {
			return s == v
		}
	case int:
		if n, ok := a.(float64); ok {
			return int(n) == v
		}
	case bool:
		if b, ok := a.(bool); ok {
			return b == v
		}
	case map[string]interface{}:
		if m, ok := a.(map[string]interface{}); ok {
			if len(m) != len(v) {
				return false
			}
			for k, val := range v {
				if !compareValues(m[k], val) {
					return false
				}
			}
			return true
		}
	}
	return false
}
