package server

import "github-mcp-server-go/protocol"

// configGetToolDef returns the definition for the config_get tool
func configGetToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "config_get",
		Description: "Get a configuration value",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"key": {
					Type:        "string",
					Description: "Configuration key to retrieve",
				},
				"scope": {
					Type:        "string",
					Description: "Configuration scope (global or local)",
					Enum:        []string{"global", "local"},
					Default:     "global",
				},
			},
			Required: []string{"key"},
		},
	}
}

// configSetToolDef returns the definition for the config_set tool
func configSetToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "config_set",
		Description: "Set a configuration value",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"key": {
					Type:        "string",
					Description: "Configuration key to set",
				},
				"value": {
					Type:        "string",
					Description: "Configuration value",
				},
				"scope": {
					Type:        "string",
					Description: "Configuration scope (global or local)",
					Enum:        []string{"global", "local"},
					Default:     "global",
				},
			},
			Required: []string{"key", "value"},
		},
	}
}

// configListToolDef returns the definition for the config_list tool
func configListToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "config_list",
		Description: "List configuration values",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"scope": {
					Type:        "string",
					Description: "Configuration scope (global, local, or all)",
					Enum:        []string{"global", "local", "all"},
					Default:     "global",
				},
			},
		},
	}
}

// configDeleteToolDef returns the definition for the config_delete tool
func configDeleteToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "config_delete",
		Description: "Delete a configuration value",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"key": {
					Type:        "string",
					Description: "Configuration key to delete",
				},
				"scope": {
					Type:        "string",
					Description: "Configuration scope (global or local)",
					Enum:        []string{"global", "local"},
					Default:     "global",
				},
			},
			Required: []string{"key"},
		},
	}
}

// aliasSetToolDef returns the definition for the alias_set tool
func aliasSetToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "alias_set",
		Description: "Create a new command alias",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"name": {
					Type:        "string",
					Description: "Alias name",
				},
				"command": {
					Type:        "string",
					Description: "Command to execute (use $1, $2, etc. for parameters)",
				},
				"description": {
					Type:        "string",
					Description: "Optional description of the alias",
				},
			},
			Required: []string{"name", "command"},
		},
	}
}

// aliasListToolDef returns the definition for the alias_list tool
func aliasListToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "alias_list",
		Description: "List all defined aliases",
		Schema: protocol.ToolSchema{
			Type:       "object",
			Properties: map[string]protocol.Property{},
		},
	}
}

// aliasDeleteToolDef returns the definition for the alias_delete tool
func aliasDeleteToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "alias_delete",
		Description: "Delete a command alias",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"name": {
					Type:        "string",
					Description: "Alias name to delete",
				},
			},
			Required: []string{"name"},
		},
	}
}

// aliasExpandToolDef returns the definition for the alias_expand tool
func aliasExpandToolDef() *protocol.Tool {
	return &protocol.Tool{
		Name:        "alias_expand",
		Description: "Expand an alias with provided arguments",
		Schema: protocol.ToolSchema{
			Type: "object",
			Properties: map[string]protocol.Property{
				"name": {
					Type:        "string",
					Description: "Alias name to expand",
				},
				"args": {
					Type:        "array",
					Description: "Arguments to replace in the alias command",
				},
			},
			Required: []string{"name"},
		},
	}
}
