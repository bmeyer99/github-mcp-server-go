package org

import (
	"testing"
	"time"

	"github-mcp-server-go/types"
)

type mockTeamClient struct {
	teams   map[string][]*types.Team
	members map[string]map[string][]string            // org -> team -> members
	repos   map[string]map[string][]*types.Repository // org -> team -> repos
}

func newMockTeamClient() *mockTeamClient {
	return &mockTeamClient{
		teams:   make(map[string][]*types.Team),
		members: make(map[string]map[string][]string),
		repos:   make(map[string]map[string][]*types.Repository),
	}
}

func (m *mockTeamClient) ListTeams(org string) ([]*types.Team, error) {
	if teams, ok := m.teams[org]; ok {
		return teams, nil
	}
	return []*types.Team{}, nil
}

func (m *mockTeamClient) GetTeam(org string, team string) (*types.Team, error) {
	teams, ok := m.teams[org]
	if !ok {
		return nil, ErrNotFound
	}
	for _, t := range teams {
		if t.Slug == team {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockTeamClient) CreateTeam(org string, params *types.TeamParams) (*types.Team, error) {
	team := &types.Team{
		ID:          1,
		Name:        params.Name,
		Slug:        params.Name,
		Description: params.Description,
		Permission:  params.Permission,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if _, ok := m.teams[org]; !ok {
		m.teams[org] = []*types.Team{}
	}
	m.teams[org] = append(m.teams[org], team)
	return team, nil
}

func (m *mockTeamClient) UpdateTeam(org string, team string, params *types.TeamParams) error {
	t, err := m.GetTeam(org, team)
	if err != nil {
		return err
	}
	if params.Name != "" {
		t.Name = params.Name
	}
	if params.Description != "" {
		t.Description = params.Description
	}
	if params.Permission != "" {
		t.Permission = params.Permission
	}
	t.UpdatedAt = time.Now()
	return nil
}

func (m *mockTeamClient) DeleteTeam(org string, team string) error {
	teams, ok := m.teams[org]
	if !ok {
		return ErrNotFound
	}
	for i, t := range teams {
		if t.Slug == team {
			m.teams[org] = append(teams[:i], teams[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (m *mockTeamClient) AddTeamMember(org string, team string, username string) error {
	if _, ok := m.members[org]; !ok {
		m.members[org] = make(map[string][]string)
	}
	if _, ok := m.members[org][team]; !ok {
		m.members[org][team] = []string{}
	}
	m.members[org][team] = append(m.members[org][team], username)
	return nil
}

func (m *mockTeamClient) RemoveTeamMember(org string, team string, username string) error {
	if orgMembers, ok := m.members[org]; ok {
		if teamMembers, ok := orgMembers[team]; ok {
			for i, member := range teamMembers {
				if member == username {
					orgMembers[team] = append(teamMembers[:i], teamMembers[i+1:]...)
					return nil
				}
			}
		}
	}
	return ErrNotFound
}

func (m *mockTeamClient) ListTeamRepos(org string, team string) ([]*types.Repository, error) {
	if orgRepos, ok := m.repos[org]; ok {
		if teamRepos, ok := orgRepos[team]; ok {
			return teamRepos, nil
		}
	}
	return []*types.Repository{}, nil
}

func (m *mockTeamClient) AddTeamRepo(org string, team string, repo string) error {
	if _, ok := m.repos[org]; !ok {
		m.repos[org] = make(map[string][]*types.Repository)
	}
	if _, ok := m.repos[org][team]; !ok {
		m.repos[org][team] = []*types.Repository{}
	}
	repository := &types.Repository{
		ID:   1,
		Name: repo,
	}
	m.repos[org][team] = append(m.repos[org][team], repository)
	return nil
}

func (m *mockTeamClient) RemoveTeamRepo(org string, team string, repo string) error {
	if orgRepos, ok := m.repos[org]; ok {
		if teamRepos, ok := orgRepos[team]; ok {
			for i, r := range teamRepos {
				if r.Name == repo {
					orgRepos[team] = append(teamRepos[:i], teamRepos[i+1:]...)
					return nil
				}
			}
		}
	}
	return ErrNotFound
}

func TestTeamManager(t *testing.T) {
	client := newMockTeamClient()
	manager := NewTeamManager(client)

	t.Run("Create and Get Team", func(t *testing.T) {
		params := &types.TeamParams{
			Name:        "test-team",
			Description: "Test Team",
			Permission:  "push",
		}

		// Create team
		team, err := manager.Create("testorg", params)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		if team.Name != params.Name {
			t.Errorf("Expected team name %s, got %s", params.Name, team.Name)
		}

		// Get team
		fetchedTeam, err := manager.Get("testorg", "test-team")
		if err != nil {
			t.Fatalf("Failed to get team: %v", err)
		}
		if fetchedTeam.Name != team.Name {
			t.Errorf("Expected team name %s, got %s", team.Name, fetchedTeam.Name)
		}
	})

	t.Run("List Teams", func(t *testing.T) {
		teams, err := manager.List("testorg")
		if err != nil {
			t.Fatalf("Failed to list teams: %v", err)
		}
		if len(teams) != 1 {
			t.Errorf("Expected 1 team, got %d", len(teams))
		}
	})

	t.Run("Update Team", func(t *testing.T) {
		params := &types.TeamParams{
			Name:        "test-team",
			Description: "Updated Description",
			Permission:  "admin",
		}

		err := manager.Update("testorg", "test-team", params)
		if err != nil {
			t.Fatalf("Failed to update team: %v", err)
		}

		team, _ := manager.Get("testorg", "test-team")
		if team.Description != params.Description {
			t.Errorf("Expected description %s, got %s", params.Description, team.Description)
		}
		if team.Permission != params.Permission {
			t.Errorf("Expected permission %s, got %s", params.Permission, team.Permission)
		}
	})

	t.Run("Invalid Team Permission", func(t *testing.T) {
		params := &types.TeamParams{
			Name:       "invalid-team",
			Permission: "invalid",
		}

		_, err := manager.Create("testorg", params)
		if err == nil || err.Error() != "invalid team permission: invalid" {
			t.Errorf("Expected invalid team permission error, got %v", err)
		}
	})

	t.Run("Team Members", func(t *testing.T) {
		// Add member
		err := manager.AddMember("testorg", "test-team", "testuser")
		if err != nil {
			t.Fatalf("Failed to add team member: %v", err)
		}

		// Remove member
		err = manager.RemoveMember("testorg", "test-team", "testuser")
		if err != nil {
			t.Fatalf("Failed to remove team member: %v", err)
		}
	})

	t.Run("Team Repositories", func(t *testing.T) {
		// Add repository
		err := manager.AddRepo("testorg", "test-team", "test-repo")
		if err != nil {
			t.Fatalf("Failed to add team repository: %v", err)
		}

		// List repositories
		repos, err := manager.ListRepos("testorg", "test-team")
		if err != nil {
			t.Fatalf("Failed to list team repositories: %v", err)
		}
		if len(repos) != 1 {
			t.Errorf("Expected 1 repository, got %d", len(repos))
		}

		// Remove repository
		err = manager.RemoveRepo("testorg", "test-team", "test-repo")
		if err != nil {
			t.Fatalf("Failed to remove team repository: %v", err)
		}
	})
}
