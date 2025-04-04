package server

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github-mcp-server-go/config"
	"github-mcp-server-go/protocol"
)

// registerConfigTools registers configuration-related tools
func (s *Server) registerConfigTools() {
	// Initialize managers if not already done
	if s.configManager == nil {
		manager, err := config.NewManager(config.ManagerOptions{
			ConfigDir: filepath.Join(s.config.ConfigDir, "config"),
		})
		if err != nil {
			s.config.Logger.Printf("Failed to initialize config manager: %v", err)
			return
		}
		s.configManager = manager

		// Initialize alias manager
		s.aliasManager = config.NewAliasManager(manager.GlobalStore())
	}

	// Register config tools
	s.tools["config_get"] = s.handleConfigGet
	s.tools["config_set"] = s.handleConfigSet
	s.tools["config_list"] = s.handleConfigList
	s.tools["config_delete"] = s.handleConfigDelete

	// Register alias tools
	s.tools["alias_set"] = s.handleAliasSet
	s.tools["alias_list"] = s.handleAliasList
	s.tools["alias_delete"] = s.handleAliasDelete
	s.tools["alias_expand"] = s.handleAliasExpand
}

// handleConfigGet handles the config_get tool
func (s *Server) handleConfigGet(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("key is required and must be a string")
	}

	scope := config.ScopeGlobal
	if scopeStr, ok := args["scope"].(string); ok {
		scope = config.Scope(strings.ToLower(scopeStr))
	}

	value, err := s.configManager.Get(key, scope)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to get config: %v", err)),
			},
		}, nil
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(fmt.Sprintf("%v", value)),
		},
	}, nil
}

// handleConfigSet handles the config_set tool
func (s *Server) handleConfigSet(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("key is required and must be a string")
	}

	value, ok := args["value"]
	if !ok {
		return nil, fmt.Errorf("value is required")
	}

	scope := config.ScopeGlobal
	if scopeStr, ok := args["scope"].(string); ok {
		scope = config.Scope(strings.ToLower(scopeStr))
	}

	if err := s.configManager.Set(key, value, scope); err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to set config: %v", err)),
			},
		}, nil
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(fmt.Sprintf("Successfully set %s = %v in %s scope", key, value, scope)),
		},
	}, nil
}

// handleConfigList handles the config_list tool
func (s *Server) handleConfigList(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	scope := config.ScopeGlobal
	if scopeStr, ok := args["scope"].(string); ok {
		scope = config.Scope(strings.ToLower(scopeStr))
	}

	var configs map[string]interface{}
	var err error

	if scope == "all" {
		configs = s.configManager.ListAll()
	} else {
		configs, err = s.configManager.List(scope)
		if err != nil {
			return &protocol.CallToolResult{
				Content: []protocol.Content{
					protocol.ErrorContent(fmt.Sprintf("Failed to list configs: %v", err)),
				},
			}, nil
		}
	}

	// Format output
	var result strings.Builder
	result.WriteString("Configuration values:\n")
	for k, v := range configs {
		result.WriteString(fmt.Sprintf("%s = %v\n", k, v))
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result.String()),
		},
	}, nil
}

// handleConfigDelete handles the config_delete tool
func (s *Server) handleConfigDelete(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("key is required and must be a string")
	}

	scope := config.ScopeGlobal
	if scopeStr, ok := args["scope"].(string); ok {
		scope = config.Scope(strings.ToLower(scopeStr))
	}

	if err := s.configManager.Delete(key, scope); err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to delete config: %v", err)),
			},
		}, nil
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(fmt.Sprintf("Successfully deleted %s from %s scope", key, scope)),
		},
	}, nil
}

// handleAliasSet handles the alias_set tool
func (s *Server) handleAliasSet(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name is required and must be a string")
	}

	command, ok := args["command"].(string)
	if !ok || command == "" {
		return nil, fmt.Errorf("command is required and must be a string")
	}

	description := ""
	if desc, ok := args["description"].(string); ok {
		description = desc
	}

	if s.aliasManager == nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent("Alias manager not initialized"),
			},
		}, nil
	}

	if err := s.aliasManager.CreateAlias(name, command, description); err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to set alias: %v", err)),
			},
		}, nil
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(fmt.Sprintf("Successfully created alias %s", name)),
		},
	}, nil
}

// handleAliasList handles the alias_list tool
func (s *Server) handleAliasList(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	if s.aliasManager == nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent("Alias manager not initialized"),
			},
		}, nil
	}

	aliases, err := s.aliasManager.ListAliases()
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to list aliases: %v", err)),
			},
		}, nil
	}

	var result strings.Builder
	result.WriteString("Defined aliases:\n")
	for _, alias := range aliases {
		result.WriteString(fmt.Sprintf("%s: %s", alias.Name, alias.Command))
		if alias.Description != "" {
			result.WriteString(fmt.Sprintf(" (%s)", alias.Description))
		}
		result.WriteString("\n")
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(result.String()),
		},
	}, nil
}

// handleAliasDelete handles the alias_delete tool
func (s *Server) handleAliasDelete(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name is required and must be a string")
	}

	if s.aliasManager == nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent("Alias manager not initialized"),
			},
		}, nil
	}

	if err := s.aliasManager.DeleteAlias(name); err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to delete alias: %v", err)),
			},
		}, nil
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(fmt.Sprintf("Successfully deleted alias %s", name)),
		},
	}, nil
}

// handleAliasExpand handles the alias_expand tool
func (s *Server) handleAliasExpand(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name is required and must be a string")
	}

	argsArray := []string{}
	if argsInterface, ok := args["args"].([]interface{}); ok {
		for _, arg := range argsInterface {
			if str, ok := arg.(string); ok {
				argsArray = append(argsArray, str)
			}
		}
	}

	if s.aliasManager == nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent("Alias manager not initialized"),
			},
		}, nil
	}

	alias, err := s.aliasManager.GetAlias(name)
	if err != nil {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.ErrorContent(fmt.Sprintf("Failed to get alias: %v", err)),
			},
		}, nil
	}

	expanded := alias.Command
	for i, arg := range argsArray {
		placeholder := fmt.Sprintf("$%d", i+1)
		expanded = strings.ReplaceAll(expanded, placeholder, arg)
	}

	// Replace any remaining numbered parameters with empty string
	for i := len(argsArray) + 1; i <= 9; i++ {
		placeholder := fmt.Sprintf("$%d", i)
		expanded = strings.ReplaceAll(expanded, placeholder, "")
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent(strings.TrimSpace(expanded)),
		},
	}, nil
}
