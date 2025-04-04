// Package org provides organization management functionality
package org

import (
	"fmt"
	"github-mcp-server-go/internal/types"
)

const (
	roleMember = "member"
	roleAdmin  = "admin"
)

// Service provides organization and team management operations
type Service interface {
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

// impl is the internal implementation of the Service interface
type impl struct {
	orgAPI  OrgAPI
	teamAPI TeamAPI
}

// OrgAPI defines organization operations
type OrgAPI interface {
	ListOrganizations() ([]*types.Organization, error)
	GetOrganization(name string) (*types.Organization, error)
	UpdateSettings(name string, settings *types.RepoSettings) error
	ListMembers(name string) ([]*types.Member, error)
	AddMember(org, user string, role string) error
	RemoveMember(org, user string) error
}

// TeamAPI defines team operations
type TeamAPI interface {
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

// NewService creates a new organization service instance
func NewService(orgAPI OrgAPI, teamAPI TeamAPI) Service {
	return &impl{
		orgAPI:  orgAPI,
		teamAPI: teamAPI,
	}
}

// Organization operations implementation

func (s *impl) ListOrganizations() ([]*types.Organization, error) {
	return s.orgAPI.ListOrganizations()
}

func (s *impl) GetOrganization(name string) (*types.Organization, error) {
	return s.orgAPI.GetOrganization(name)
}

func (s *impl) UpdateSettings(name string, settings *types.RepoSettings) error {
	return s.orgAPI.UpdateSettings(name, settings)
}

func (s *impl) ListMembers(name string) ([]*types.Member, error) {
	return s.orgAPI.ListMembers(name)
}

func (s *impl) AddMember(org, user string, role string) error {
	if !validateRole(role) {
		return fmt.Errorf("invalid role: %s (must be '%s' or '%s')", role, roleMember, roleAdmin)
	}
	return s.orgAPI.AddMember(org, user, role)
}

func (s *impl) RemoveMember(org, user string) error {
	return s.orgAPI.RemoveMember(org, user)
}

// Team operations implementation

func (s *impl) ListTeams(org string) ([]*types.Team, error) {
	return s.teamAPI.ListTeams(org)
}

func (s *impl) GetTeam(org string, team string) (*types.Team, error) {
	return s.teamAPI.GetTeam(org, team)
}

func (s *impl) CreateTeam(org string, params *types.TeamParams) (*types.Team, error) {
	if params.Permission != "" && !validateTeamPermission(params.Permission) {
		return nil, fmt.Errorf("invalid team permission: %s", params.Permission)
	}
	return s.teamAPI.CreateTeam(org, params)
}

func (s *impl) UpdateTeam(org string, team string, params *types.TeamParams) error {
	if params.Permission != "" && !validateTeamPermission(params.Permission) {
		return fmt.Errorf("invalid team permission: %s", params.Permission)
	}
	return s.teamAPI.UpdateTeam(org, team, params)
}

func (s *impl) DeleteTeam(org string, team string) error {
	return s.teamAPI.DeleteTeam(org, team)
}

func (s *impl) AddTeamMember(org string, team string, username string) error {
	return s.teamAPI.AddTeamMember(org, team, username)
}

func (s *impl) RemoveTeamMember(org string, team string, username string) error {
	return s.teamAPI.RemoveTeamMember(org, team, username)
}

func (s *impl) ListTeamRepos(org string, team string) ([]*types.Repository, error) {
	return s.teamAPI.ListTeamRepos(org, team)
}

func (s *impl) AddTeamRepo(org string, team string, repo string) error {
	return s.teamAPI.AddTeamRepo(org, team, repo)
}

func (s *impl) RemoveTeamRepo(org string, team string, repo string) error {
	return s.teamAPI.RemoveTeamRepo(org, team, repo)
}

func (s *impl) GetNestedTeams(org string, team string) ([]*types.Team, error) {
	teams, err := s.ListTeams(org)
	if err != nil {
		return nil, err
	}

	var nestedTeams []*types.Team
	parentTeam, err := s.GetTeam(org, team)
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

// Validation helpers

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
