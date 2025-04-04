package config

import (
	"fmt"
	"testing"
)

// mockStore implements ConfigStore for testing
type mockStore struct {
	data map[string]interface{}
}

func newMockStore() *mockStore {
	return &mockStore{
		data: make(map[string]interface{}),
	}
}

func (s *mockStore) Get(key string) (interface{}, error) {
	if value, ok := s.data[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("key not found: %s", key)
}

func (s *mockStore) Set(key string, value interface{}) error {
	s.data[key] = value
	return nil
}

func (s *mockStore) Delete(key string) error {
	if _, ok := s.data[key]; !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	delete(s.data, key)
	return nil
}

func (s *mockStore) List() (map[string]interface{}, error) {
	return s.data, nil
}

func TestAliasManager_CreateAlias(t *testing.T) {
	store := newMockStore()
	manager := NewAliasManager(store)

	tests := []struct {
		name        string
		aliasName   string
		command     string
		description string
		wantErr     bool
	}{
		{
			name:        "valid alias",
			aliasName:   "test-alias",
			command:     "git commit -m $1",
			description: "Create a commit with message",
			wantErr:     false,
		},
		{
			name:        "empty name",
			aliasName:   "",
			command:     "git status",
			description: "Check git status",
			wantErr:     true,
		},
		{
			name:        "empty command",
			aliasName:   "status",
			command:     "",
			description: "Empty command",
			wantErr:     true,
		},
		{
			name:        "invalid name characters",
			aliasName:   "test@alias",
			command:     "git status",
			description: "Invalid characters",
			wantErr:     true,
		},
		{
			name:        "valid with parameters",
			aliasName:   "commit",
			command:     "git commit -m $1 -a $2",
			description: "Commit with message and flag",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.CreateAlias(tt.aliasName, tt.command, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify alias was created correctly
				alias, err := manager.GetAlias(tt.aliasName)
				if err != nil {
					t.Errorf("Failed to get created alias: %v", err)
					return
				}

				if alias.Name != tt.aliasName {
					t.Errorf("Alias name = %v, want %v", alias.Name, tt.aliasName)
				}
				if alias.Command != tt.command {
					t.Errorf("Alias command = %v, want %v", alias.Command, tt.command)
				}
				if alias.Description != tt.description {
					t.Errorf("Alias description = %v, want %v", alias.Description, tt.description)
				}
			}
		})
	}
}

func TestAliasManager_DeleteAlias(t *testing.T) {
	store := newMockStore()
	manager := NewAliasManager(store)

	// Create test alias
	aliasName := "test-alias"
	if err := manager.CreateAlias(aliasName, "test command", "test description"); err != nil {
		t.Fatalf("Failed to create test alias: %v", err)
	}

	// Test deletion
	if err := manager.DeleteAlias(aliasName); err != nil {
		t.Errorf("DeleteAlias() error = %v", err)
	}

	// Verify alias was deleted
	if _, err := manager.GetAlias(aliasName); err == nil {
		t.Error("Expected error getting deleted alias")
	}

	// Test deleting non-existent alias
	if err := manager.DeleteAlias("nonexistent"); err == nil {
		t.Error("Expected error deleting non-existent alias")
	}
}

func TestAliasManager_ListAliases(t *testing.T) {
	store := newMockStore()
	manager := NewAliasManager(store)

	// Create test aliases
	testAliases := []struct {
		name        string
		command     string
		description string
	}{
		{"alias1", "command1 $1", "description1"},
		{"alias2", "command2 $1 $2", "description2"},
		{"alias3", "command3", "description3"},
	}

	for _, a := range testAliases {
		if err := manager.CreateAlias(a.name, a.command, a.description); err != nil {
			t.Fatalf("Failed to create test alias: %v", err)
		}
	}

	// List aliases
	aliases, err := manager.ListAliases()
	if err != nil {
		t.Fatalf("ListAliases() error = %v", err)
	}

	if len(aliases) != len(testAliases) {
		t.Errorf("Expected %d aliases, got %d", len(testAliases), len(aliases))
	}

	// Verify each alias
	for _, want := range testAliases {
		found := false
		for _, got := range aliases {
			if got.Name == want.name {
				found = true
				if got.Command != want.command {
					t.Errorf("Alias command = %v, want %v", got.Command, want.command)
				}
				if got.Description != want.description {
					t.Errorf("Alias description = %v, want %v", got.Description, want.description)
				}
				break
			}
		}
		if !found {
			t.Errorf("Alias %s not found in list", want.name)
		}
	}
}

func TestAliasManager_ExpandAlias(t *testing.T) {
	store := newMockStore()
	manager := NewAliasManager(store)

	// Create test alias
	aliasName := "test-alias"
	command := "git commit -m $1 -a $2"
	if err := manager.CreateAlias(aliasName, command, "test description"); err != nil {
		t.Fatalf("Failed to create test alias: %v", err)
	}

	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "all parameters",
			args: []string{"message", "true"},
			want: "git commit -m message -a true",
		},
		{
			name: "partial parameters",
			args: []string{"message"},
			want: "git commit -m message -a",
		},
		{
			name: "extra parameters",
			args: []string{"message", "true", "extra"},
			want: "git commit -m message -a true",
		},
		{
			name: "no parameters",
			args: []string{},
			want: "git commit -m -a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expanded, err := manager.ExpandAlias(aliasName, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if expanded != tt.want {
				t.Errorf("ExpandAlias() = %v, want %v", expanded, tt.want)
			}
		})
	}
}

func TestAliasNameValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "test", false},
		{"valid with hyphen", "test-alias", false},
		{"valid with underscore", "test_alias", false},
		{"valid with numbers", "test123", false},
		{"empty", "", true},
		{"space", "test alias", true},
		{"special chars", "test@alias", true},
		{"dot", "test.alias", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAliasName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAliasName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
