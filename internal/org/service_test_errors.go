package org

import (
	"fmt"
	"testing"

	"github-mcp-server-go/internal/types"
)

type errorOrgsAPI struct {
	shouldError bool
}

type errorTeamsAPI struct {
	shouldError bool
}

func (m *errorOrgsAPI) ListOrganizations() ([]*types.Organization, error) {
	if m.shouldError {
		return nil, fmt.Errorf("failed to list organizations")
	}
	return nil, nil
}

func (m *errorOrgsAPI) GetOrganization(name string) (*types.Organization, error) {
	if m.shouldError {
		return nil, fmt.Errorf("failed to get organization")
	}
	return nil, nil
}

func (m *errorOrgsAPI) UpdateSettings(name string, settings *types.RepoSettings) error {
	if m.shouldError {
		return fmt.Errorf("failed to update settings")
	}
	return nil
}

func (m *errorOrgsAPI) ListMembers(name string) ([]*types.Member, error) {
	if m.shouldError {
		return nil, fmt.Errorf("failed to list members")
	}
	return nil, nil
}

func (m *errorOrgsAPI) AddMember(org, user string, role string) error {
	if m.shouldError {
		return fmt.Errorf("failed to add member")
	}
	return nil
}

func (m *errorOrgsAPI) RemoveMember(org, user string) error {
	if m.shouldError {
		return fmt.Errorf("failed to remove member")
	}
	return nil
}

func (m *errorTeamsAPI) ListTeams(org string) ([]*types.Team, error) {
	if m.shouldError {
		return nil, fmt.Errorf("failed to list teams")
	}
	return nil, nil
}

func (m *errorTeamsAPI) GetTeam(org string, team string) (*types.Team, error) {
	if m.shouldError {
		return nil, fmt.Errorf("failed to get team")
	}
	return nil, nil
}

func (m *errorTeamsAPI) CreateTeam(org string, params *types.TeamParams) (*types.Team, error) {
	if m.shouldError {
		return nil, fmt.Errorf("failed to create team")
	}
	return nil, nil
}

func (m *errorTeamsAPI) UpdateTeam(org string, team string, params *types.TeamParams) error {
	if m.shouldError {
		return fmt.Errorf("failed to update team")
	}
	return nil
}

func (m *errorTeamsAPI) DeleteTeam(org string, team string) error {
	if m.shouldError {
		return fmt.Errorf("failed to delete team")
	}
	return nil
}

func (m *errorTeamsAPI) AddTeamMember(org string, team string, username string) error {
	if m.shouldError {
		return fmt.Errorf("failed to add team member")
	}
	return nil
}

func (m *errorTeamsAPI) RemoveTeamMember(org string, team string, username string) error {
	if m.shouldError {
		return fmt.Errorf("failed to remove team member")
	}
	return nil
}

func (m *errorTeamsAPI) ListTeamRepos(org string, team string) ([]*types.Repository, error) {
	if m.shouldError {
		return nil, fmt.Errorf("failed to list team repos")
	}
	return nil, nil
}

func (m *errorTeamsAPI) AddTeamRepo(org string, team string, repo string) error {
	if m.shouldError {
		return fmt.Errorf("failed to add team repo")
	}
	return nil
}

func (m *errorTeamsAPI) RemoveTeamRepo(org string, team string, repo string) error {
	if m.shouldError {
		return fmt.Errorf("failed to remove team repo")
	}
	return nil
}

func TestServiceErrors(t *testing.T) {
	tests := []struct {
		name        string
		orgsAPI     types.OrganizationAPI
		teamsAPI    types.TeamAPI
		operation   func(s Service) error
		expectError bool
	}{
		{
			name:        "ListOrganizations error",
			orgsAPI:     &errorOrgsAPI{shouldError: true},
			teamsAPI:    &errorTeamsAPI{},
			operation:   func(s Service) error { _, err := s.ListOrganizations(); return err },
			expectError: true,
		},
		{
			name:        "AddMember invalid role",
			orgsAPI:     &errorOrgsAPI{},
			teamsAPI:    &errorTeamsAPI{},
			operation:   func(s Service) error { return s.AddMember("org", "user", "invalid") },
			expectError: true,
		},
		{
			name:     "CreateTeam invalid permission",
			orgsAPI:  &errorOrgsAPI{},
			teamsAPI: &errorTeamsAPI{},
			operation: func(s Service) error {
				_, err := s.CreateTeam("org", &types.TeamParams{Permission: "invalid"})
				return err
			},
			expectError: true,
		},
		{
			name:        "GetNestedTeams list error",
			orgsAPI:     &errorOrgsAPI{},
			teamsAPI:    &errorTeamsAPI{shouldError: true},
			operation:   func(s Service) error { _, err := s.GetNestedTeams("org", "team"); return err },
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New(tt.orgsAPI, tt.teamsAPI)
			err := tt.operation(service)
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
