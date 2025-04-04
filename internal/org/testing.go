package org

import (
	"fmt"
	"time"

	"github-mcp-server-go/internal/types"
	"sync"
)

// Mock implementations for testing

type mockOrgsAPI struct {
	mu       sync.RWMutex
	orgs     []*types.Organization
	members  map[string][]*types.Member
	settings map[string]*types.RepoSettings
}

type mockTeamsAPI struct {
	mu      sync.RWMutex
	teams   map[string][]*types.Team
	members map[string]map[string][]string            // org -> team -> members
	repos   map[string]map[string][]*types.Repository // org -> team -> repos
}

func newMockOrgsAPI() *mockOrgsAPI {
	return &mockOrgsAPI{
		orgs:     make([]*types.Organization, 0),
		members:  make(map[string][]*types.Member),
		settings: make(map[string]*types.RepoSettings),
	}
}

func newMockTeamsAPI() *mockTeamsAPI {
	return &mockTeamsAPI{
		teams:   make(map[string][]*types.Team),
		members: make(map[string]map[string][]string),
		repos:   make(map[string]map[string][]*types.Repository),
	}
}

// Thread-safe mock implementations

func (m *mockOrgsAPI) ListOrganizations() ([]*types.Organization, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.orgs, nil
}

func (m *mockOrgsAPI) GetOrganization(name string) (*types.Organization, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, org := range m.orgs {
		if org.Login == name {
			return org, nil
		}
	}
	return nil, fmt.Errorf("organization not found: %s", name)
}

func (m *mockOrgsAPI) UpdateSettings(name string, settings *types.RepoSettings) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings[name] = settings
	return nil
}

func (m *mockOrgsAPI) ListMembers(name string) ([]*types.Member, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.members[name], nil
}

func (m *mockOrgsAPI) AddMember(org, user string, role string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	member := &types.Member{
		ID:        1,
		Login:     user,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.members[org] = append(m.members[org], member)
	return nil
}

func (m *mockOrgsAPI) RemoveMember(org, user string) error {
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

func (m *mockTeamsAPI) ListTeams(org string) ([]*types.Team, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.teams[org], nil
}

func (m *mockTeamsAPI) GetTeam(org string, team string) (*types.Team, error) {
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

func (m *mockTeamsAPI) CreateTeam(org string, params *types.TeamParams) (*types.Team, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	team := &types.Team{
		ID:          1,
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

func (m *mockTeamsAPI) UpdateTeam(org string, team string, params *types.TeamParams) error {
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
			return nil
		}
	}
	return fmt.Errorf("team not found: %s", team)
}

func (m *mockTeamsAPI) DeleteTeam(org string, team string) error {
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

func (m *mockTeamsAPI) AddTeamMember(org string, team string, username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.members[org] == nil {
		m.members[org] = make(map[string][]string)
	}
	m.members[org][team] = append(m.members[org][team], username)
	return nil
}

func (m *mockTeamsAPI) RemoveTeamMember(org string, team string, username string) error {
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

func (m *mockTeamsAPI) ListTeamRepos(org string, team string) ([]*types.Repository, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if repos, ok := m.repos[org][team]; ok {
		return repos, nil
	}
	return []*types.Repository{}, nil
}

func (m *mockTeamsAPI) AddTeamRepo(org string, team string, repo string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.repos[org] == nil {
		m.repos[org] = make(map[string][]*types.Repository)
	}
	repository := &types.Repository{
		ID:   1,
		Name: repo,
	}
	m.repos[org][team] = append(m.repos[org][team], repository)
	return nil
}

func (m *mockTeamsAPI) RemoveTeamRepo(org string, team string, repo string) error {
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
