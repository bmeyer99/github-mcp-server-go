package config

import (
	"fmt"
	"strings"
	"time"
)

// Alias represents a command alias
type Alias struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Description string            `json:"description,omitempty"`
	Params      map[string]string `json:"params,omitempty"`
	Created     time.Time         `json:"created"`
	Updated     time.Time         `json:"updated"`
}

// AliasManager handles command alias operations
type AliasManager struct {
	store ConfigStore
}

// NewAliasManager creates a new alias manager
func NewAliasManager(store ConfigStore) *AliasManager {
	return &AliasManager{
		store: store,
	}
}

// CreateAlias creates a new command alias
func (m *AliasManager) CreateAlias(name, command, description string) error {
	// Validate name and command
	if err := validateAliasName(name); err != nil {
		return err
	}
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Extract parameters from command
	params := extractParameters(command)

	// Create alias
	alias := &Alias{
		Name:        name,
		Command:     command,
		Description: description,
		Params:      params,
		Created:     time.Now().UTC(),
		Updated:     time.Now().UTC(),
	}

	// Store alias
	return m.store.Set("alias:"+name, alias)
}

// GetAlias retrieves an alias by name
func (m *AliasManager) GetAlias(name string) (*Alias, error) {
	value, err := m.store.Get("alias:" + name)
	if err != nil {
		return nil, err
	}

	// Convert stored value to Alias
	alias, ok := value.(*Alias)
	if !ok {
		return nil, fmt.Errorf("invalid alias data")
	}

	return alias, nil
}

// DeleteAlias removes an alias
func (m *AliasManager) DeleteAlias(name string) error {
	return m.store.Delete("alias:" + name)
}

// ListAliases returns all defined aliases
func (m *AliasManager) ListAliases() ([]*Alias, error) {
	values, err := m.store.List()
	if err != nil {
		return nil, err
	}

	var aliases []*Alias
	for key, value := range values {
		if !strings.HasPrefix(key, "alias:") {
			continue
		}

		if alias, ok := value.(*Alias); ok {
			aliases = append(aliases, alias)
		}
	}

	return aliases, nil
}

// ExpandAlias expands an alias with the provided arguments
func (m *AliasManager) ExpandAlias(name string, args []string) (string, error) {
	// Get alias
	alias, err := m.GetAlias(name)
	if err != nil {
		return "", err
	}

	// Replace parameters in command
	command := alias.Command
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		command = strings.ReplaceAll(command, placeholder, arg)
	}

	// Replace any remaining numbered parameters with empty string
	for i := len(args) + 1; i <= 9; i++ {
		placeholder := fmt.Sprintf("$%d", i)
		command = strings.ReplaceAll(command, placeholder, "")
	}

	return strings.TrimSpace(command), nil
}

// validateAliasName checks if the alias name is valid
func validateAliasName(name string) error {
	if name == "" {
		return fmt.Errorf("alias name cannot be empty")
	}

	if !isValidAliasName(name) {
		return fmt.Errorf("invalid alias name: must contain only letters, numbers, hyphens, and underscores")
	}

	return nil
}

// isValidAliasName checks if a name contains only valid characters
func isValidAliasName(name string) bool {
	for _, c := range name {
		if !isAlphanumeric(c) && c != '-' && c != '_' {
			return false
		}
	}
	return true
}

// isAlphanumeric checks if a character is a letter or number
func isAlphanumeric(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9')
}

// extractParameters gets parameter placeholders from a command
func extractParameters(command string) map[string]string {
	params := make(map[string]string)
	parts := strings.Fields(command)

	for _, part := range parts {
		if strings.HasPrefix(part, "$") && len(part) > 1 {
			num := strings.TrimPrefix(part, "$")
			if num >= "1" && num <= "9" {
				params[num] = ""
			}
		}
	}

	return params
}
