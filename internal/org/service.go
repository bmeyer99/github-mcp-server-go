// Package org provides GitHub organization and team management functionality.
// It offers a high-level interface for managing organizations, their members,
// teams, and team repositories through a unified Service interface.
package org

import (
"fmt"
"github-mcp-server-go/internal/types"
)

const (
roleMember = "member"
roleAdmin  = "admin"
)

// Service defines a comprehensive interface for managing GitHub organizations
// and their teams. It provides methods for handling organization settings,
// member management, team operations, and repository access control.
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

// client provides the concrete implementation of the Service interface.
// It utilizes separate APIs for organization and team operations to
// maintain separation of concerns and facilitate testing.
type client struct {
orgs  types.OrganizationAPI
teams types.TeamAPI
}

// New creates a new organization service instance with the provided
// organization and team APIs. Both APIs must be non-nil.
func New(orgs types.OrganizationAPI, teams types.TeamAPI) Service {
return &client{
orgs:  orgs,
teams: teams,
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

// Organization operations

func (c *client) ListOrganizations() ([]*types.Organization, error) {
return c.orgs.ListOrganizations()
}

func (c *client) GetOrganization(name string) (*types.Organization, error) {
return c.orgs.GetOrganization(name)
}

func (c *client) UpdateSettings(name string, settings *types.RepoSettings) error {
return c.orgs.UpdateSettings(name, settings)
}

func (c *client) ListMembers(name string) ([]*types.Member, error) {
return c.orgs.ListMembers(name)
}

func (c *client) AddMember(org, user string, role string) error {
if !validateRole(role) {
return fmt.Errorf("invalid role: %s (must be '%s' or '%s')", role, roleMember, roleAdmin)
}
return c.orgs.AddMember(org, user, role)
}

func (c *client) RemoveMember(org, user string) error {
return c.orgs.RemoveMember(org, user)
}

// Team operations

func (c *client) ListTeams(org string) ([]*types.Team, error) {
return c.teams.ListTeams(org)
}

func (c *client) GetTeam(org string, team string) (*types.Team, error) {
return c.teams.GetTeam(org, team)
}

func (c *client) CreateTeam(org string, params *types.TeamParams) (*types.Team, error) {
if params == nil {
return nil, fmt.Errorf("team parameters cannot be nil")
}
if params.Permission != "" && !validateTeamPermission(params.Permission) {
return nil, fmt.Errorf("invalid team permission: %s", params.Permission)
}
return c.teams.CreateTeam(org, params)
}

func (c *client) UpdateTeam(org string, team string, params *types.TeamParams) error {
if params.Permission != "" && !validateTeamPermission(params.Permission) {
return fmt.Errorf("invalid team permission: %s", params.Permission)
}
return c.teams.UpdateTeam(org, team, params)
}

func (c *client) DeleteTeam(org string, team string) error {
return c.teams.DeleteTeam(org, team)
}

func (c *client) AddTeamMember(org string, team string, username string) error {
return c.teams.AddTeamMember(org, team, username)
}

func (c *client) RemoveTeamMember(org string, team string, username string) error {
return c.teams.RemoveTeamMember(org, team, username)
}

func (c *client) ListTeamRepos(org string, team string) ([]*types.Repository, error) {
return c.teams.ListTeamRepos(org, team)
}

func (c *client) AddTeamRepo(org string, team string, repo string) error {
return c.teams.AddTeamRepo(org, team, repo)
}

func (c *client) RemoveTeamRepo(org string, team string, repo string) error {
return c.teams.RemoveTeamRepo(org, team, repo)
}

func (c *client) GetNestedTeams(org string, team string) ([]*types.Team, error) {
teams, err := c.ListTeams(org)
if err != nil {
return nil, err
}

var nestedTeams []*types.Team
parentTeam, err := c.GetTeam(org, team)
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
