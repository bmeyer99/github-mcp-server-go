// Package org provides organization and team management functionality
package org

import (
	"fmt"
	"github-mcp-server-go/types"
)

// API provides access to organization and team management operations
type API struct {
	orgClient  Client
	teamClient TeamClient
}

// Client defines organization operations
type Client interface {
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

// Options provides configuration for creating a new API instance
type Options struct {
	OrgClient  Client
	TeamClient TeamClient
}

// New creates a new API instance
func New(opts Options) *API {
	return &API{
		orgClient:  opts.OrgClient,
		teamClient: opts.TeamClient,
	}
}
