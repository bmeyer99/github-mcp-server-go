package org

// ValidRole checks if a role is valid
func ValidRole(role string) bool {
	return role == "member" || role == "admin"
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
