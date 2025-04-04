// Package org provides organization and team management functionality
package org

import (
	"fmt"
	"github-mcp-server-go/types"
)

const (
	roleMember = "member"
	roleAdmin  = "admin"
)

// Manager handles organization and team operations
type Manager interface {
	// Organization operations
	ListOrganizations() ([]*types.Organization, error)
	GetOrganization(name string) (*types.Organization, error)
	UpdateSettings(name string, settings *types.RepoSettings) error
	ListMembers(name string) ([]*types.Member, error)
	AddMember(org, user string, role string) error
	RemoveMember(org, user string) error

	// Team operations
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
	GetNestedTeams(org string, team string) ([]*types.Team, error)
}

// orgManager implements the Manager interface
type orgManager struct {
	orgClient  OrgClient
	teamClient TeamClient
}

// OrgClient defines organization operations
type OrgClient interface {
	ListOrganizations() ([]*types.Organization, error)
	GetOrganization(name string) (*types.Organization, error)
	UpdateSettings(name string, settings *types.RepoSettings) error
	ListMembers(name string) ([]*types.Member, error)
	AddMember(org, user string, role string) error
	RemoveMember(org, user string) error
}

// TeamClient defines team operations
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

// Config provides configuration for creating a new Manager
type Config struct {
	OrgClient  OrgClient
	TeamClient TeamClient
}

// New creates a new Manager instance
func New(cfg Config) Manager {
	return &orgManager{
		orgClient:  cfg.OrgClient,
		teamClient: cfg.TeamClient,
	}
}

// validate role and permission functions
func validateRole(role string) bool {
	return role == roleMember || role == roleAdmin
}

func validateTeamPermission(permission string) bool {
	validPermissions := map[string]bool{
		"pull":     true,
		"push":     true,
		"admin":    true,
		"maintain": true,
		"triage":   true,
	}
	return validPermissions[permission]
}

// Organization operations implementation

func (m *orgManager) ListOrganizations() ([]*types.Organization, error) {
	return m.orgClient.ListOrganizations()
}

func (m *orgManager) GetOrganization(name string) (*types.Organization, error) {
	return m.orgClient.GetOrganization(name)
}

func (m *orgManager) UpdateSettings(name string, settings *types.RepoSettings) error {
	return m.orgClient.UpdateSettings(name, settings)
}

func (m *orgManager) ListMembers(name string) ([]*types.Member, error) {
	return m.orgClient.ListMembers(name)
}

func (m *orgManager) AddMember(org, user string, role string) error {
	if !validateRole(role) {
		return fmt.Errorf("invalid role: %s (must be '%s' or '%s')", role, roleMember, roleAdmin)
	}
	return m.orgClient.AddMember(org, user, role)
}

func (m *orgManager) RemoveMember(org, user string) error {
	return m.orgClient.RemoveMember(org, user)
}

// Team operations implementation

func (m *orgManager) ListTeams(org string) ([]*types.Team, error) {
	return m.teamClient.ListTeams(org)
}

func (m *orgManager) GetTeam(org string, team string) (*types.Team, error) {
	return m.teamClient.GetTeam(org, team)
}

func (m *orgManager) CreateTeam(org string, params *types.TeamParams) (*types.Team, error) {
	if params.Permission != "" && !validateTeamPermission(params.Permission) {
		return nil, fmt.Errorf("invalid team permission: %s", params.Permission)
	}
	return m.teamClient.CreateTeam(org, params)
}

func (m *orgManager) UpdateTeam(org string, team string, params *types.TeamParams) error {
	if params.Permission != "" && !validateTeamPermission(params.Permission) {
		return fmt.Errorf("invalid team permission: %s", params.Permission)
	}
	return m.teamClient.UpdateTeam(org, team, params)
}

func (m *orgManager) DeleteTeam(org string, team string) error {
	return m.teamClient.DeleteTeam(org, team)
}

func (m *orgManager) AddTeamMember(org string, team string, username string) error {
	return m.teamClient.AddTeamMember(org, team, username)
}

func (m *orgManager) RemoveTeamMember(org string, team string, username string) error {
	return m.teamClient.RemoveTeamMember(org, team, username)
}

func (m *orgManager) ListTeamRepos(org string, team string) ([]*types.Repository, error) {
	return m.teamClient.ListTeamRepos(org, team)
}

func (m *orgManager) AddTeamRepo(org string, team string, repo string) error {
	return m.teamClient.AddTeamRepo(org, team, repo)
}

func (m *orgManager) RemoveTeamRepo(org string, team string, repo string) error {
	return m.teamClient.RemoveTeamRepo(org, team, repo)
}

func (m *orgManager) GetNestedTeams(org string, team string) ([]*types.Team, error) {
	teams, err := m.ListTeams(org)
	if err != nil {
		return nil, err
	}

	var nestedTeams []*types.Team
	parentTeam, err := m.GetTeam(org, team)
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
