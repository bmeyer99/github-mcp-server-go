package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Scope represents a configuration scope
type Scope string

const (
	ScopeGlobal Scope = "global"
	ScopeLocal  Scope = "local"
)

// EnvProvider defines the interface for environment variable operations
type EnvProvider interface {
	Get(key string) string
	Set(key, value string) error
	List() map[string]string
}

// ManagerOptions contains options for creating a new Manager
type ManagerOptions struct {
	ConfigDir       string
	LocalConfigPath string
	EnvProvider     EnvProvider
}

// Manager handles configuration operations across different scopes
type Manager struct {
	global ConfigStore
	local  ConfigStore
	env    EnvProvider
	mu     sync.RWMutex
}

// defaultEnvProvider implements EnvProvider using os environment
type defaultEnvProvider struct{}

func (p *defaultEnvProvider) Get(key string) string {
	return os.Getenv(key)
}

func (p *defaultEnvProvider) Set(key, value string) error {
	return os.Setenv(key, value)
}

func (p *defaultEnvProvider) List() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		if k, v, ok := splitEnv(e); ok {
			env[k] = v
		}
	}
	return env
}

// NewManager creates a new configuration manager
func NewManager(opts ManagerOptions) (*Manager, error) {
	if opts.ConfigDir == "" {
		opts.ConfigDir = filepath.Join(os.Getenv("HOME"), ".config", "github-mcp")
	}

	if opts.EnvProvider == nil {
		opts.EnvProvider = &defaultEnvProvider{}
	}

	// Create global config store
	global, err := NewFileStore(filepath.Join(opts.ConfigDir, "config.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to create global config store: %w", err)
	}

	// Create local config store if path is provided
	var local ConfigStore
	if opts.LocalConfigPath != "" {
		local, err = NewFileStore(opts.LocalConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create local config store: %w", err)
		}
	}

	return &Manager{
		global: global,
		local:  local,
		env:    opts.EnvProvider,
	}, nil
}

// Get retrieves a configuration value from the specified scope
func (m *Manager) Get(key string, scope Scope) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check environment variables first
	if envVal := m.env.Get(key); envVal != "" {
		return envVal, nil
	}

	// Check local config if available
	if scope == ScopeLocal && m.local != nil {
		if value, err := m.local.Get(key); err == nil {
			return value, nil
		}
	}

	// Fall back to global config
	return m.global.Get(key)
}

// Set stores a configuration value in the specified scope
func (m *Manager) Set(key string, value interface{}, scope Scope) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch scope {
	case ScopeLocal:
		if m.local == nil {
			return fmt.Errorf("local configuration is not available")
		}
		return m.local.Set(key, value)

	case ScopeGlobal:
		return m.global.Set(key, value)

	default:
		return fmt.Errorf("invalid scope: %s", scope)
	}
}

// Delete removes a configuration value from the specified scope
func (m *Manager) Delete(key string, scope Scope) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch scope {
	case ScopeLocal:
		if m.local == nil {
			return fmt.Errorf("local configuration is not available")
		}
		return m.local.Delete(key)

	case ScopeGlobal:
		return m.global.Delete(key)

	default:
		return fmt.Errorf("invalid scope: %s", scope)
	}
}

// List returns all configuration values from the specified scope
// GlobalStore returns the global configuration store
func (m *Manager) GlobalStore() ConfigStore {
	return m.global
}

func (m *Manager) List(scope Scope) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result map[string]interface{}
	var err error

	switch scope {
	case ScopeLocal:
		if m.local == nil {
			return nil, fmt.Errorf("local configuration is not available")
		}
		result, err = m.local.List()

	case ScopeGlobal:
		result, err = m.global.List()

	default:
		return nil, fmt.Errorf("invalid scope: %s", scope)
	}

	if err != nil {
		return nil, err
	}

	// Add environment variables to the result
	for k, v := range m.env.List() {
		result[k] = v
	}

	return result, nil
}

// ListAll returns all configuration values from all scopes
func (m *Manager) ListAll() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]interface{})

	// Add global config
	if global, err := m.global.List(); err == nil {
		for k, v := range global {
			result[k] = v
		}
	}

	// Add local config if available
	if m.local != nil {
		if local, err := m.local.List(); err == nil {
			for k, v := range local {
				result[k] = v
			}
		}
	}

	// Add environment variables
	for k, v := range m.env.List() {
		result[k] = v
	}

	return result
}

// Helper function to split environment variables
func splitEnv(env string) (key, value string, ok bool) {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return env[:i], env[i+1:], true
		}
	}
	return "", "", false
}
