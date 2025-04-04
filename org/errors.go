package org

import "errors"

var (
	// ErrNotFound indicates a requested resource was not found
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidRole indicates an invalid role was specified
	ErrInvalidRole = errors.New("invalid role: must be 'member' or 'admin'")

	// ErrInvalidPermission indicates an invalid permission was specified
	ErrInvalidPermission = errors.New("invalid permission: must be 'read', 'write', or 'admin'")

	// ErrInvalidTeamPermission indicates an invalid team permission was specified
	ErrInvalidTeamPermission = errors.New("invalid team permission: must be 'pull', 'push', 'admin', 'maintain', or 'triage'")
)

// ValidRole checks if a role is valid
func ValidRole(role string) bool {
	return role == "member" || role == "admin"
}

// ValidPermission checks if a permission is valid
func ValidPermission(permission string) bool {
	validPermissions := map[string]bool{
		"read":  true,
		"write": true,
		"admin": true,
	}
	return validPermissions[permission]
}

// ValidTeamPermission checks if a team permission is valid
func ValidTeamPermission(permission string) bool {
	validPermissions := map[string]bool{
		"pull":     true,
		"push":     true,
		"admin":    true,
		"maintain": true,
		"triage":   true,
	}
	return validPermissions[permission]
}
