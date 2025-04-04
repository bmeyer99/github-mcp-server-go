package org

import (
	"testing"

	"github-mcp-server-go/internal/types"
)

func TestErrorConditions(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(orgAPI *mockOrgAPIConcurrent, teamAPI *mockTeamAPIConcurrent)
		test    func(s Service) error
		wantErr bool
	}{
		{
			name:  "Get non-existent organization",
			setup: func(orgAPI *mockOrgAPIConcurrent, teamAPI *mockTeamAPIConcurrent) {},
			test: func(s Service) error {
				_, err := s.GetOrganization("nonexistent")
				return err
			},
			wantErr: true,
		},
		{
			name:  "Add member with invalid role",
			setup: func(orgAPI *mockOrgAPIConcurrent, teamAPI *mockTeamAPIConcurrent) {},
			test: func(s Service) error {
				return s.AddMember("org", "user", "invalid-role")
			},
			wantErr: true,
		},
		{
			name:  "Create team with empty name",
			setup: func(orgAPI *mockOrgAPIConcurrent, teamAPI *mockTeamAPIConcurrent) {},
			test: func(s Service) error {
				_, err := s.CreateTeam("org", &types.TeamParams{})
				return err
			},
			wantErr: true,
		},
		{
			name:  "Create team with invalid permission",
			setup: func(orgAPI *mockOrgAPIConcurrent, teamAPI *mockTeamAPIConcurrent) {},
			test: func(s Service) error {
				_, err := s.CreateTeam("org", &types.TeamParams{
					Name:       "team",
					Permission: "invalid",
				})
				return err
			},
			wantErr: true,
		},
		{
			name:  "Get nested teams for non-existent team",
			setup: func(orgAPI *mockOrgAPIConcurrent, teamAPI *mockTeamAPIConcurrent) {},
			test: func(s Service) error {
				_, err := s.GetNestedTeams("org", "nonexistent")
				return err
			},
			wantErr: true,
		},
		{
			name:  "Remove non-existent team member",
			setup: func(orgAPI *mockOrgAPIConcurrent, teamAPI *mockTeamAPIConcurrent) {},
			test: func(s Service) error {
				return s.RemoveTeamMember("org", "team", "nonexistent")
			},
			wantErr: true,
		},
		{
			name:  "Remove non-existent team repository",
			setup: func(orgAPI *mockOrgAPIConcurrent, teamAPI *mockTeamAPIConcurrent) {},
			test: func(s Service) error {
				return s.RemoveTeamRepo("org", "team", "nonexistent")
			},
			wantErr: true,
		},
		{
			name:  "Update settings for non-existent organization",
			setup: func(orgAPI *mockOrgAPIConcurrent, teamAPI *mockTeamAPIConcurrent) {},
			test: func(s Service) error {
				return s.UpdateSettings("nonexistent", &types.RepoSettings{})
			},
			wantErr: true,
		},
		{
			name:  "List members for non-existent organization",
			setup: func(orgAPI *mockOrgAPIConcurrent, teamAPI *mockTeamAPIConcurrent) {},
			test: func(s Service) error {
				_, err := s.ListMembers("nonexistent")
				return err
			},
			wantErr: false, // ListMembers returns empty list for non-existent org
		},
		{
			name:  "List teams for non-existent organization",
			setup: func(orgAPI *mockOrgAPIConcurrent, teamAPI *mockTeamAPIConcurrent) {},
			test: func(s Service) error {
				_, err := s.ListTeams("nonexistent")
				return err
			},
			wantErr: false, // ListTeams returns empty list for non-existent org
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgAPI := newMockOrgAPIConcurrent()
			teamAPI := newMockTeamAPIConcurrent()
			tt.setup(orgAPI, teamAPI)
			service := New(orgAPI, teamAPI)

			err := tt.test(service)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
