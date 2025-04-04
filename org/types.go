// Package org provides organization and team management functionality
package org

import (
	"github-mcp-server-go/types"
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

// serviceConfig provides configuration for creating a new service
type serviceConfig struct {
	orgAPI  OrgAPI
	teamAPI TeamAPI
}

// service implements the Service interface
type service struct {
	orgAPI  OrgAPI
	teamAPI TeamAPI
}

// NewService creates a new organization service
func NewService(orgAPI OrgAPI, teamAPI TeamAPI) Service {
	return &service{
		orgAPI:  orgAPI,
		teamAPI: teamAPI,
	}
}
