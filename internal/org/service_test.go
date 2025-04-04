package org

import (
	"testing"

	"github-mcp-server-go/internal/types"
)

func TestService(t *testing.T) {
	orgsAPI := newMockOrgsAPI()
	teamsAPI := newMockTeamsAPI()
	service := New(orgsAPI, teamsAPI)

	t.Run("Organization Operations", func(t *testing.T) {
		t.Run("Add and List Members", func(t *testing.T) {
			// Add member
			err := service.AddMember("testorg", "testuser", "member")
			if err != nil {
				t.Fatalf("Failed to add member: %v", err)
			}

			// Invalid role
			err = service.AddMember("testorg", "testuser2", "invalid")
			if err == nil {
				t.Error("Expected error for invalid role")
			}

			// List members
			members, err := service.ListMembers("testorg")
			if err != nil {
				t.Fatalf("Failed to list members: %v", err)
			}
			if len(members) != 1 {
				t.Errorf("Expected 1 member, got %d", len(members))
			}
			if members[0].Login != "testuser" {
				t.Errorf("Expected member testuser, got %s", members[0].Login)
			}
		})

		t.Run("Settings Management", func(t *testing.T) {
			settings := &types.RepoSettings{
				DefaultRepoPermission: "write",
				MembersCanCreateRepos: true,
				TwoFactorRequired:     true,
			}

			err := service.UpdateSettings("testorg", settings)
			if err != nil {
				t.Fatalf("Failed to update settings: %v", err)
			}

			if orgsAPI.settings["testorg"].DefaultRepoPermission != "write" {
				t.Error("Settings were not updated correctly")
			}
		})
	})

	t.Run("Team Operations", func(t *testing.T) {
		t.Run("Create and List Teams", func(t *testing.T) {
			params := &types.TeamParams{
				Name:        "testteam",
				Description: "Test Team",
				Permission:  "push",
			}

			team, err := service.CreateTeam("testorg", params)
			if err != nil {
				t.Fatalf("Failed to create team: %v", err)
			}
			if team.Name != params.Name {
				t.Errorf("Expected team name %s, got %s", params.Name, team.Name)
			}

			// Invalid permission
			params.Permission = "invalid"
			_, err = service.CreateTeam("testorg", params)
			if err == nil {
				t.Error("Expected error for invalid permission")
			}

			// List teams
			teams, err := service.ListTeams("testorg")
			if err != nil {
				t.Fatalf("Failed to list teams: %v", err)
			}
			if len(teams) != 1 {
				t.Errorf("Expected 1 team, got %d", len(teams))
			}
		})

		t.Run("Team Member Management", func(t *testing.T) {
			err := service.AddTeamMember("testorg", "testteam", "testuser")
			if err != nil {
				t.Fatalf("Failed to add team member: %v", err)
			}

			err = service.RemoveTeamMember("testorg", "testteam", "testuser")
			if err != nil {
				t.Fatalf("Failed to remove team member: %v", err)
			}
		})

		t.Run("Team Repository Management", func(t *testing.T) {
			err := service.AddTeamRepo("testorg", "testteam", "testrepo")
			if err != nil {
				t.Fatalf("Failed to add team repository: %v", err)
			}

			repos, err := service.ListTeamRepos("testorg", "testteam")
			if err != nil {
				t.Fatalf("Failed to list team repositories: %v", err)
			}
			if len(repos) != 1 {
				t.Errorf("Expected 1 repository, got %d", len(repos))
			}

			err = service.RemoveTeamRepo("testorg", "testteam", "testrepo")
			if err != nil {
				t.Fatalf("Failed to remove team repository: %v", err)
			}
		})

		t.Run("Nested Teams", func(t *testing.T) {
			// Create parent team
			parentParams := &types.TeamParams{
				Name:       "parent",
				Permission: "admin",
			}
			parent, err := service.CreateTeam("testorg", parentParams)
			if err != nil {
				t.Fatalf("Failed to create parent team: %v", err)
			}

			// Create child team
			childParams := &types.TeamParams{
				Name:       "child",
				Permission: "push",
				ParentID:   parent.ID,
			}
			child, err := service.CreateTeam("testorg", childParams)
			if err != nil {
				t.Fatalf("Failed to create child team: %v", err)
			}

			// Set parent reference for test
			child.Parent = parent

			// Get nested teams
			nested, err := service.GetNestedTeams("testorg", "parent")
			if err != nil {
				t.Fatalf("Failed to get nested teams: %v", err)
			}
			if len(nested) != 1 {
				t.Errorf("Expected 1 nested team, got %d", len(nested))
			}
			if nested[0].Name != "child" {
				t.Errorf("Expected child team, got %s", nested[0].Name)
			}
		})
	})
}
