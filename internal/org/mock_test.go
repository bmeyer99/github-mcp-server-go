package org

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github-mcp-server-go/internal/types"
)

// Mock test helpers for concurrent testing
type mockOrgAPIConcurrent struct {
	mu      sync.RWMutex
	orgs    []*types.Organization
	members map[string][]*types.Member
}

func newMockOrgAPIConcurrent() *mockOrgAPIConcurrent {
	return &mockOrgAPIConcurrent{
		orgs:    make([]*types.Organization, 0),
		members: make(map[string][]*types.Member),
	}
}

// Organization API methods
func (m *mockOrgAPIConcurrent) ListOrganizations() ([]*types.Organization, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]*types.Organization{}, m.orgs...), nil
}

func (m *mockOrgAPIConcurrent) GetOrganization(name string) (*types.Organization, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, org := range m.orgs {
		if org.Login == name {
			return org, nil
		}
	}
	return nil, fmt.Errorf("organization not found: %s", name)
}

func (m *mockOrgAPIConcurrent) UpdateSettings(name string, settings *types.RepoSettings) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, org := range m.orgs {
		if org.Login == name {
			org.Settings = *settings
			return nil
		}
	}
	return fmt.Errorf("organization not found: %s", name)
}

func (m *mockOrgAPIConcurrent) ListMembers(name string) ([]*types.Member, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]*types.Member{}, m.members[name]...), nil
}

func (m *mockOrgAPIConcurrent) AddMember(org, user string, role string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	member := &types.Member{
		ID:        time.Now().UnixNano(),
		Login:     user,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.members[org] = append(m.members[org], member)
	return nil
}

func (m *mockOrgAPIConcurrent) RemoveMember(org, user string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	members := m.members[org]
	for i, member := range members {
		if member.Login == user {
			m.members[org] = append(members[:i], members[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("member not found: %s", user)
}

type mockTeamAPIConcurrent struct {
	mu      sync.RWMutex
	teams   map[string][]*types.Team
	members map[string]map[string][]string
	repos   map[string]map[string][]*types.Repository
}

func newMockTeamAPIConcurrent() *mockTeamAPIConcurrent {
	return &mockTeamAPIConcurrent{
		teams:   make(map[string][]*types.Team),
		members: make(map[string]map[string][]string),
		repos:   make(map[string]map[string][]*types.Repository),
	}
}

// Team API methods
func (m *mockTeamAPIConcurrent) ListTeams(org string) ([]*types.Team, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]*types.Team{}, m.teams[org]...), nil
}

func (m *mockTeamAPIConcurrent) GetTeam(org string, team string) (*types.Team, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	teams := m.teams[org]
	for _, t := range teams {
		if t.Slug == team {
			return t, nil
		}
	}
	return nil, fmt.Errorf("team not found: %s", team)
}

func (m *mockTeamAPIConcurrent) CreateTeam(org string, params *types.TeamParams) (*types.Team, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	team := &types.Team{
		ID:          time.Now().UnixNano(),
		Name:        params.Name,
		Slug:        params.Name,
		Description: params.Description,
		Permission:  params.Permission,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if m.teams[org] == nil {
		m.teams[org] = []*types.Team{}
	}
	m.teams[org] = append(m.teams[org], team)
	return team, nil
}

func (m *mockTeamAPIConcurrent) UpdateTeam(org string, team string, params *types.TeamParams) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	teams := m.teams[org]
	for _, t := range teams {
		if t.Slug == team {
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
	}
	return fmt.Errorf("team not found: %s", team)
}

func (m *mockTeamAPIConcurrent) DeleteTeam(org string, team string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	teams := m.teams[org]
	for i, t := range teams {
		if t.Slug == team {
			m.teams[org] = append(teams[:i], teams[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("team not found: %s", team)
}

func (m *mockTeamAPIConcurrent) AddTeamMember(org string, team string, username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.members[org] == nil {
		m.members[org] = make(map[string][]string)
	}
	m.members[org][team] = append(m.members[org][team], username)
	return nil
}

func (m *mockTeamAPIConcurrent) RemoveTeamMember(org string, team string, username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if members, ok := m.members[org][team]; ok {
		for i, member := range members {
			if member == username {
				m.members[org][team] = append(members[:i], members[i+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("member not found: %s", username)
}

func (m *mockTeamAPIConcurrent) ListTeamRepos(org string, team string) ([]*types.Repository, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.repos[org] == nil {
		return []*types.Repository{}, nil
	}
	return append([]*types.Repository{}, m.repos[org][team]...), nil
}

func (m *mockTeamAPIConcurrent) AddTeamRepo(org string, team string, repo string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.repos[org] == nil {
		m.repos[org] = make(map[string][]*types.Repository)
	}
	repository := &types.Repository{
		ID:   time.Now().UnixNano(),
		Name: repo,
	}
	m.repos[org][team] = append(m.repos[org][team], repository)
	return nil
}

func (m *mockTeamAPIConcurrent) RemoveTeamRepo(org string, team string, repo string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if repos, ok := m.repos[org][team]; ok {
		for i, r := range repos {
			if r.Name == repo {
				m.repos[org][team] = append(repos[:i], repos[i+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("repository not found: %s", repo)
}

func TestConcurrentSafety(t *testing.T) {
	orgAPI := newMockOrgAPIConcurrent()
	teamAPI := newMockTeamAPIConcurrent()
	service := New(orgAPI, teamAPI)

	const numGoroutines = 10
	var wg sync.WaitGroup

	// Test concurrent organization operations
	t.Run("Concurrent Organization Operations", func(t *testing.T) {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()

				org := "testorg"
				user := fmt.Sprintf("user%d", id)

				err := service.AddMember(org, user, roleMember)
				if err != nil {
					t.Errorf("Failed to add member: %v", err)
				}

				members, err := service.ListMembers(org)
				if err != nil {
					t.Errorf("Failed to list members: %v", err)
				} else if len(members) < 1 {
					t.Errorf("Expected at least 1 member, got %d", len(members))
				}

				err = service.RemoveMember(org, user)
				if err != nil {
					t.Errorf("Failed to remove member: %v", err)
				}
			}(i)
		}
		wg.Wait()
	})

	// Test concurrent team operations
	t.Run("Concurrent Team Operations", func(t *testing.T) {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()

				org := "testorg"
				teamName := fmt.Sprintf("team%d", id)

				params := &types.TeamParams{
					Name:        teamName,
					Description: "Test team",
					Permission:  "push",
				}
				team, err := service.CreateTeam(org, params)
				if err != nil {
					t.Errorf("Failed to create team: %v", err)
					return
				}

				user := fmt.Sprintf("user%d", id)
				err = service.AddTeamMember(org, team.Name, user)
				if err != nil {
					t.Errorf("Failed to add team member: %v", err)
				}

				repo := fmt.Sprintf("repo%d", id)
				err = service.AddTeamRepo(org, team.Name, repo)
				if err != nil {
					t.Errorf("Failed to add team repo: %v", err)
				}

				repos, err := service.ListTeamRepos(org, team.Name)
				if err != nil {
					t.Errorf("Failed to list team repos: %v", err)
				} else if len(repos) < 1 {
					t.Errorf("Expected at least 1 repo, got %d", len(repos))
				}

				err = service.RemoveTeamRepo(org, team.Name, repo)
				if err != nil {
					t.Errorf("Failed to remove team repo: %v", err)
				}
			}(i)
		}
		wg.Wait()
	})
}
