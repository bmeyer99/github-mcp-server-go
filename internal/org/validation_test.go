package org

import "testing"

func TestValidateRole(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		isValid bool
	}{
		{
			name:    "Valid member role",
			role:    "member",
			isValid: true,
		},
		{
			name:    "Valid admin role",
			role:    "admin",
			isValid: true,
		},
		{
			name:    "Invalid empty role",
			role:    "",
			isValid: false,
		},
		{
			name:    "Invalid role",
			role:    "superuser",
			isValid: false,
		},
		{
			name:    "Case sensitivity check",
			role:    "MEMBER",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateRole(tt.role); got != tt.isValid {
				t.Errorf("validateRole(%q) = %v, want %v", tt.role, got, tt.isValid)
			}
		})
	}
}

func TestValidateTeamPermission(t *testing.T) {
	tests := []struct {
		name       string
		permission string
		isValid    bool
	}{
		{
			name:       "Valid pull permission",
			permission: "pull",
			isValid:    true,
		},
		{
			name:       "Valid push permission",
			permission: "push",
			isValid:    true,
		},
		{
			name:       "Valid admin permission",
			permission: "admin",
			isValid:    true,
		},
		{
			name:       "Valid maintain permission",
			permission: "maintain",
			isValid:    true,
		},
		{
			name:       "Valid triage permission",
			permission: "triage",
			isValid:    true,
		},
		{
			name:       "Invalid empty permission",
			permission: "",
			isValid:    false,
		},
		{
			name:       "Invalid permission",
			permission: "readwrite",
			isValid:    false,
		},
		{
			name:       "Case sensitivity check",
			permission: "PULL",
			isValid:    false,
		},
		{
			name:       "Invalid whitespace",
			permission: "pull ",
			isValid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateTeamPermission(tt.permission); got != tt.isValid {
				t.Errorf("validateTeamPermission(%q) = %v, want %v", tt.permission, got, tt.isValid)
			}
		})
	}
}
