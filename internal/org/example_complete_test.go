package org_test

import (
	"fmt"
	"log"

	"github-mcp-server-go/internal/org"
	"github-mcp-server-go/internal/types"
)

func Example_complete() {
	// This example demonstrates the complete workflow of managing
	// an organization, teams, and repositories

	// Initialize mock API clients for the example
	orgAPI := &mockOrgAPI{}
	teamAPI := &mockTeamAPI{}
	svc := org.New(orgAPI, teamAPI)

	// 1. Get organization details
	org, err := svc.GetOrganization("my-org")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Update organization settings
	settings := &types.RepoSettings{
		DefaultRepoPermission: "write",
		MembersCanCreateRepos: true,
		TwoFactorRequired:     true,
	}
	if err := svc.UpdateSettings(org.Login, settings); err != nil {
		log.Fatal(err)
	}

	// 3. Create a team hierarchy
	engineeringParams := &types.TeamParams{
		Name:        "engineering",
		Description: "Engineering department",
		Permission:  "admin",
	}
	engineering, err := svc.CreateTeam(org.Login, engineeringParams)
	if err != nil {
		log.Fatal(err)
	}

	// Create sub-teams
	teamParams := &types.TeamParams{
		Name:        "backend",
		Description: "Backend development team",
		Permission:  "push",
		ParentID:    engineering.ID,
	}
	if _, err := svc.CreateTeam(org.Login, teamParams); err != nil {
		log.Fatal(err)
	}

	// 4. Add team members
	if err := svc.AddTeamMember(org.Login, engineering.Name, "tech-lead"); err != nil {
		log.Fatal(err)
	}

	// 5. Grant repository access
	if err := svc.AddTeamRepo(org.Login, engineering.Name, "api-service"); err != nil {
		log.Fatal(err)
	}

	// 6. List nested teams
	teams, err := svc.GetNestedTeams(org.Login, engineering.Name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Engineering sub-teams:\n")
	for _, team := range teams {
		fmt.Printf("- %s (%s)\n", team.Name, team.Permission)
	}

	// Output:
	// Engineering sub-teams:
	// - backend (push)
}

// Mock implementations for example
type mockOrgAPI struct {
	org *types.Organization
}

func newMockOrgAPI() *mockOrgAPI {
	return &mockOrgAPI{
		org: &types.Organization{
			ID:    1,
			Login: "my-org",
			Name:  "My Organization",
		},
	}
}

func (m *mockOrgAPI) GetOrganization(name string) (*types.Organization, error) {
	return m.org, nil
}

func (m *mockOrgAPI) UpdateSettings(name string, settings *types.RepoSettings) error {
	m.org.Settings = *settings
	return nil
}

func (m *mockOrgAPI) ListOrganizations() ([]*types.Organization, error) {
	return []*types.Organization{m.org}, nil
}

func (m *mockOrgAPI) ListMembers(name string) ([]*types.Member, error) {
	return nil, nil
}

func (m *mockOrgAPI) AddMember(org, user string, role string) error {
	return nil
}

func (m *mockOrgAPI) RemoveMember(org, user string) error {
	return nil
}

type mockTeamAPI struct {
	parentTeam *types.Team
	childTeam  *types.Team
}

func newMockTeamAPI() *mockTeamAPI {
	return &mockTeamAPI{}
}

func (m *mockTeamAPI) CreateTeam(org string, params *types.TeamParams) (*types.Team, error) {
	if m.parentTeam == nil {
		m.parentTeam = &types.Team{
			ID:         1,
			Name:       params.Name,
			Permission: params.Permission,
		}
		return m.parentTeam, nil
	}

	m.childTeam = &types.Team{
		ID:         2,
		Name:       params.Name,
		Permission: params.Permission,
		Parent:     m.parentTeam,
	}
	return m.childTeam, nil
}

func (m *mockTeamAPI) ListTeams(org string) ([]*types.Team, error) {
	teams := []*types.Team{m.parentTeam}
	if m.childTeam != nil {
		teams = append(teams, m.childTeam)
	}
	return teams, nil
}

func (m *mockTeamAPI) GetTeam(org string, team string) (*types.Team, error) {
	if team == m.parentTeam.Name {
		return m.parentTeam, nil
	}
	return m.childTeam, nil
}

func (m *mockTeamAPI) UpdateTeam(org string, team string, params *types.TeamParams) error {
	return nil
}

func (m *mockTeamAPI) DeleteTeam(org string, team string) error {
	return nil
}

func (m *mockTeamAPI) AddTeamMember(org string, team string, username string) error {
	return nil
}

func (m *mockTeamAPI) RemoveTeamMember(org string, team string, username string) error {
	return nil
}

func (m *mockTeamAPI) ListTeamRepos(org string, team string) ([]*types.Repository, error) {
	return []*types.Repository{{ID: 1, Name: "api-service"}}, nil
}

func (m *mockTeamAPI) AddTeamRepo(org string, team string, repo string) error {
	return nil
}

func (m *mockTeamAPI) RemoveTeamRepo(org string, team string, repo string) error {
	return nil
}

func init() {
	// Initialize the mock APIs with test data
	orgAPI := newMockOrgAPI()
	teamAPI := newMockTeamAPI()

	// Pre-create some test data
	_ = org.New(orgAPI, teamAPI)
}
