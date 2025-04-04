package org

import (
	"testing"
	"time"

	"github-mcp-server-go/types"
)

type mockOrgClient struct {
	orgs    []*types.Organization
	members map[string][]*types.Member
}

func newMockOrgClient() *mockOrgClient {
	return &mockOrgClient{
		orgs:    make([]*types.Organization, 0),
		members: make(map[string][]*types.Member),
	}
}

func (m *mockOrgClient) ListOrganizations() ([]*types.Organization, error) {
	return m.orgs, nil
}

func (m *mockOrgClient) GetOrganization(name string) (*types.Organization, error) {
	for _, org := range m.orgs {
		if org.Login == name {
			return org, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockOrgClient) UpdateSettings(name string, settings *types.RepoSettings) error {
	org, err := m.GetOrganization(name)
	if err != nil {
		return err
	}
	org.DefaultRepoSettings = *settings
	return nil
}

func (m *mockOrgClient) ListMembers(name string) ([]*types.Member, error) {
	if members, ok := m.members[name]; ok {
		return members, nil
	}
	return []*types.Member{}, nil
}

func (m *mockOrgClient) AddMember(org, user string, role string) error {
	if !validateRole(role) { // Use package-level validator
		return ErrInvalidRole
	}
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

func (m *mockOrgClient) RemoveMember(org, user string) error {
	if members, ok := m.members[org]; ok {
		for i, member := range members {
			if member.Login == user {
				m.members[org] = append(members[:i], members[i+1:]...)
				return nil
			}
		}
	}
	return ErrNotFound
}

// Add TeamClient method implementations to mockOrgClient
func (m *mockOrgClient) ListTeams(org string) ([]*types.Team, error) {
	// Placeholder implementation
	return []*types.Team{}, nil
}

func (m *mockOrgClient) GetTeam(org string, team string) (*types.Team, error) {
	// Placeholder implementation
	return nil, ErrNotFound
}

func (m *mockOrgClient) CreateTeam(org string, params *types.TeamParams) (*types.Team, error) {
	// Placeholder implementation
	return &types.Team{ID: 1, Name: params.Name, Slug: params.Name}, nil
}

func (m *mockOrgClient) UpdateTeam(org string, team string, params *types.TeamParams) error {
	// Placeholder implementation
	return nil
}

func (m *mockOrgClient) DeleteTeam(org string, team string) error {
	// Placeholder implementation
	return nil
}

func (m *mockOrgClient) AddTeamMember(org string, team string, username string) error {
	// Placeholder implementation
	return nil
}

func (m *mockOrgClient) RemoveTeamMember(org string, team string, username string) error {
	// Placeholder implementation
	return nil
}

func (m *mockOrgClient) ListTeamRepos(org string, team string) ([]*types.Repository, error) {
	// Placeholder implementation
	return []*types.Repository{}, nil
}

func (m *mockOrgClient) AddTeamRepo(org string, team string, repo string) error {
	// Placeholder implementation
	return nil
}

func (m *mockOrgClient) RemoveTeamRepo(org string, team string, repo string) error {
	// Placeholder implementation
	return nil
}


// Renamed test function for clarity
func TestOrgManagerImplementation(t *testing.T) {
	client := newMockOrgClient()
	// Directly instantiate orgManager and assign to Manager interface type
	var manager Manager = &orgManager{
		orgClient:  client,
		teamClient: client,
	}

	t.Run("List Organizations", func(t *testing.T) {
		orgs, err := manager.ListOrganizations()
		if err != nil {
			t.Fatalf("Failed to list organizations: %v", err)
		}
		if len(orgs) != 0 {
			t.Errorf("Expected 0 organizations, got %d", len(orgs))
		}
	})

	t.Run("Add and Remove Member", func(t *testing.T) {
		// Test adding member
		err := manager.AddMember("testorg", "testuser", "member")
		if err != nil {
			t.Fatalf("Failed to add member: %v", err)
		}

		// Verify member was added
		members, err := manager.ListMembers("testorg")
		if err != nil {
			t.Fatalf("Failed to list members: %v", err)
		}
		if len(members) != 1 {
			t.Errorf("Expected 1 member, got %d", len(members))
		}
		if members[0].Login != "testuser" {
			t.Errorf("Expected member login 'testuser', got '%s'", members[0].Login)
		}

		// Test removing member
		err = manager.RemoveMember("testorg", "testuser")
		if err != nil {
			t.Fatalf("Failed to remove member: %v", err)
		}

		// Verify member was removed
		members, err = manager.ListMembers("testorg")
		if err != nil {
			t.Fatalf("Failed to list members: %v", err)
		}
		if len(members) != 0 {
			t.Errorf("Expected 0 members after removal, got %d", len(members))
		}
	})

	t.Run("Invalid Role", func(t *testing.T) {
		err := manager.AddMember("testorg", "testuser", "invalid")
		if err != ErrInvalidRole {
			t.Errorf("Expected ErrInvalidRole, got %v", err)
		}
	})
}
