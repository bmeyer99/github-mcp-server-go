package org

import (
	"testing"

	"github-mcp-server-go/internal/types"
)

// Additional test cases to ensure full coverage

func TestEdgeCases(t *testing.T) {
	orgsAPI := newMockOrgsAPI()
	teamsAPI := newMockTeamsAPI()
	service := New(orgsAPI, teamsAPI)

	t.Run("Empty Organization Operations", func(t *testing.T) {
		// List empty organizations
		orgs, err := service.ListOrganizations()
		if err != nil {
			t.Fatalf("Failed to list empty organizations: %v", err)
		}
		if len(orgs) != 0 {
			t.Errorf("Expected 0 organizations, got %d", len(orgs))
		}

		// Get non-existent organization
		_, err = service.GetOrganization("non-existent")
		if err == nil {
			t.Error("Expected error for non-existent organization")
		}
	})

	t.Run("Invalid Team Operations", func(t *testing.T) {
		// Create team with empty name
		emptyParams := &types.TeamParams{}
		_, err := service.CreateTeam("testorg", emptyParams)
		if err == nil {
			t.Error("Expected error for empty team name")
		}

		// Update team with invalid permission
		invalidParams := &types.TeamParams{
			Name:       "team",
			Permission: "invalid",
		}
		err = service.UpdateTeam("testorg", "team", invalidParams)
		if err == nil {
			t.Error("Expected error for invalid permission")
		}
	})

	t.Run("Complex Team Hierarchy", func(t *testing.T) {
		// Create multi-level team hierarchy
		params := &types.TeamParams{
			Name:        "root",
			Description: "Root team",
			Permission:  "admin",
		}
		root, err := service.CreateTeam("testorg", params)
		if err != nil {
			t.Fatalf("Failed to create root team: %v", err)
		}

		// Create multiple child teams
		childParams := []struct {
			name string
			perm string
		}{
			{"child1", "push"},
			{"child2", "pull"},
			{"child3", "maintain"},
		}

		for _, cp := range childParams {
			params := &types.TeamParams{
				Name:        cp.name,
				Description: "Child team",
				Permission:  cp.perm,
				ParentID:    root.ID,
			}
			_, err := service.CreateTeam("testorg", params)
			if err != nil {
				t.Errorf("Failed to create child team %s: %v", cp.name, err)
			}
		}

		// Verify nested teams
		nested, err := service.GetNestedTeams("testorg", root.Name)
		if err != nil {
			t.Fatalf("Failed to get nested teams: %v", err)
		}
		if len(nested) != len(childParams) {
			t.Errorf("Expected %d nested teams, got %d", len(childParams), len(nested))
		}
	})

	t.Run("Repository Access Control", func(t *testing.T) {
		// Add multiple repositories
		repos := []string{"repo1", "repo2", "repo3"}

		for _, repo := range repos {
			err := service.AddTeamRepo("testorg", "team", repo)
			if err != nil {
				t.Errorf("Failed to add repo %s: %v", repo, err)
			}
		}

		// List repositories
		repoList, err := service.ListTeamRepos("testorg", "team")
		if err != nil {
			t.Fatalf("Failed to list repos: %v", err)
		}
		if len(repoList) != len(repos) {
			t.Errorf("Expected %d repos, got %d", len(repos), len(repoList))
		}

		// Remove repositories
		for _, repo := range repos {
			err := service.RemoveTeamRepo("testorg", "team", repo)
			if err != nil {
				t.Errorf("Failed to remove repo %s: %v", repo, err)
			}
		}

		// Verify empty repo list
		repoList, err = service.ListTeamRepos("testorg", "team")
		if err != nil {
			t.Fatalf("Failed to list repos after removal: %v", err)
		}
		if len(repoList) != 0 {
			t.Errorf("Expected 0 repos after removal, got %d", len(repoList))
		}
	})

	t.Run("Member Role Operations", func(t *testing.T) {
		// Test member roles
		roles := []string{"", "invalid", "member", "admin"}
		for _, role := range roles {
			err := service.AddMember("testorg", "user", role)
			if role == "member" || role == "admin" {
				if err != nil {
					t.Errorf("Failed to add member with valid role %s: %v", role, err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for invalid role: %s", role)
				}
			}
		}
	})

	t.Run("Team Permission Operations", func(t *testing.T) {
		// Test team permissions
		permissions := []string{
			"",
			"invalid",
			"pull",
			"push",
			"admin",
			"maintain",
			"triage",
		}

		for _, perm := range permissions {
			params := &types.TeamParams{
				Name:       "team",
				Permission: perm,
			}
			_, err := service.CreateTeam("testorg", params)
			if perm == "pull" || perm == "push" || perm == "admin" || perm == "maintain" || perm == "triage" {
				if err != nil {
					t.Errorf("Failed to create team with valid permission %s: %v", perm, err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for invalid permission: %s", perm)
				}
			}
		}
	})

	t.Run("Error Conditions", func(t *testing.T) {
		cases := []struct {
			name string
			fn   func() error
		}{
			{
				name: "Add member with empty role",
				fn: func() error {
					return service.AddMember("org", "user", "")
				},
			},
			{
				name: "Get non-existent team",
				fn: func() error {
					_, err := service.GetTeam("org", "nonexistent")
					return err
				},
			},
			{
				name: "Update non-existent team",
				fn: func() error {
					return service.UpdateTeam("org", "nonexistent", &types.TeamParams{})
				},
			},
			{
				name: "Remove non-existent member",
				fn: func() error {
					return service.RemoveMember("org", "nonexistent")
				},
			},
			{
				name: "Delete non-existent team",
				fn: func() error {
					return service.DeleteTeam("org", "nonexistent")
				},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if err := tc.fn(); err == nil {
					t.Error("Expected error but got nil")
				}
			})
		}
	})
}
