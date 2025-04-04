package org

import (
	"fmt"
	"github-mcp-server-go/types"
)

// Manager handles organization operations
type Manager struct {
	client Client
}

// Client defines the interface for organization operations
type Client interface {
	ListOrganizations() ([]*types.Organization, error)
	GetOrganization(name string) (*types.Organization, error)
	UpdateSettings(name string, settings *types.RepoSettings) error
	ListMembers(name string) ([]*types.Member, error)
	AddMember(org, user string, role string) error
	RemoveMember(org, user string) error
}

// TeamManager handles team operations
type TeamManager struct {
	client TeamClient
}

// TeamClient defines the interface for team operations
type TeamClient interface {
	ListTeams(org string) ([]*types.Team, error)
	GetTeam(org string, team string) (*types.Team, error)
	CreateTeam(org string, params *types.TeamParams) (*types.Team, error)
	UpdateTeam(org string, team string, params *types.TeamParams) error
	DeleteTeam(org string, team string) error
	AddTeamMember(org string, team string, username string) error
	RemoveTeamMember(org string, team string, username string) error
	ListTeamRepos(org string, team string) ([]*types.Repository, error)
	AddTeamRepo(org string, team string, repo string) error
	RemoveTeamRepo(org string, team string, repo string) error
}

// NewManager creates a new organization manager
func NewManager(client Client) *Manager {
	return &Manager{
		client: client,
	}
}

// NewTeamManager creates a new team manager
func NewTeamManager(client TeamClient) *TeamManager {
	return &TeamManager{
		client: client,
	}
}

// List returns all organizations for the authenticated user
func (m *Manager) List() ([]*types.Organization, error) {
	return m.client.ListOrganizations()
}

// Get retrieves an organization by name
func (m *Manager) Get(name string) (*types.Organization, error) {
	return m.client.GetOrganization(name)
}

// UpdateSettings updates organization repository settings
func (m *Manager) UpdateSettings(name string, settings *types.RepoSettings) error {
	return m.client.UpdateSettings(name, settings)
}

// ListMembers returns all members of an organization
func (m *Manager) ListMembers(name string) ([]*types.Member, error) {
	return m.client.ListMembers(name)
}

// AddMember adds a user to an organization
func (m *Manager) AddMember(org, user string, role string) error {
	if !ValidRole(role) {
		return fmt.Errorf("invalid role: %s (must be 'member' or 'admin')", role)
	}
	return m.client.AddMember(org, user, role)
}

// RemoveMember removes a user from an organization
func (m *Manager) RemoveMember(org, user string) error {
	return m.client.RemoveMember(org, user)
}

// List returns all teams in an organization
func (m *TeamManager) List(org string) ([]*types.Team, error) {
	return m.client.ListTeams(org)
}

// Get retrieves a team by name
func (m *TeamManager) Get(org string, team string) (*types.Team, error) {
	return m.client.GetTeam(org, team)
}

// Create creates a new team
func (m *TeamManager) Create(org string, params *types.TeamParams) (*types.Team, error) {
	// Validate team permission if set
	if params.Permission != "" && !ValidTeamPermission(params.Permission) {
		return nil, fmt.Errorf("invalid team permission: %s", params.Permission)
	}

	return m.client.CreateTeam(org, params)
}

// Update updates a team
func (m *TeamManager) Update(org string, team string, params *types.TeamParams) error {
	// Validate team permission if set
	if params.Permission != "" && !ValidTeamPermission(params.Permission) {
		return fmt.Errorf("invalid team permission: %s", params.Permission)
	}

	return m.client.UpdateTeam(org, team, params)
}

// Delete removes a team
func (m *TeamManager) Delete(org string, team string) error {
	return m.client.DeleteTeam(org, team)
}

// AddMember adds a user to a team
func (m *TeamManager) AddMember(org string, team string, username string) error {
	return m.client.AddTeamMember(org, team, username)
}

// RemoveMember removes a user from a team
func (m *TeamManager) RemoveMember(org string, team string, username string) error {
	return m.client.RemoveTeamMember(org, team, username)
}

// ListRepos returns all repositories a team has access to
func (m *TeamManager) ListRepos(org string, team string) ([]*types.Repository, error) {
	return m.client.ListTeamRepos(org, team)
}

// AddRepo adds a repository to a team
func (m *TeamManager) AddRepo(org string, team string, repo string) error {
	return m.client.AddTeamRepo(org, team, repo)
}

// RemoveRepo removes a repository from a team
func (m *TeamManager) RemoveRepo(org string, team string, repo string) error {
	return m.client.RemoveTeamRepo(org, team, repo)
}

// GetNestedTeams returns all nested teams for a parent team
func (m *TeamManager) GetNestedTeams(org string, team string) ([]*types.Team, error) {
	teams, err := m.List(org)
	if err != nil {
		return nil, err
	}

	var nestedTeams []*types.Team
	parentTeam, err := m.Get(org, team)
	if err != nil {
		return nil, err
	}

	for _, t := range teams {
		if t.Parent != nil && t.Parent.ID == parentTeam.ID {
			nestedTeams = append(nestedTeams, t)
		}
	}

	return nestedTeams, nil
}

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
