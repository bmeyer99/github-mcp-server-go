package org

import (
	"sync"
	"testing"

	"github-mcp-server-go/internal/types"
)

func TestConcurrentOperations(t *testing.T) {
	orgsAPI := newMockOrgsAPI()
	teamsAPI := newMockTeamsAPI()
	service := New(orgsAPI, teamsAPI)

	t.Run("Concurrent Member Operations", func(t *testing.T) {
		var wg sync.WaitGroup
		users := []string{"user1", "user2", "user3", "user4", "user5"}
		const org = "testorg"

		// Add members concurrently
		wg.Add(len(users))
		for _, user := range users {
			go func(username string) {
				defer wg.Done()
				err := service.AddMember(org, username, roleMember)
				if err != nil {
					t.Errorf("Failed to add member %s: %v", username, err)
				}
			}(user)
		}
		wg.Wait()

		// Verify all members were added
		members, err := service.ListMembers(org)
		if err != nil {
			t.Fatalf("Failed to list members: %v", err)
		}
		if len(members) != len(users) {
			t.Errorf("Expected %d members, got %d", len(users), len(members))
		}
	})

	t.Run("Concurrent Team Operations", func(t *testing.T) {
		var wg sync.WaitGroup
		teams := []string{"team1", "team2", "team3"}
		const org = "testorg"

		// Create teams concurrently
		wg.Add(len(teams))
		for _, team := range teams {
			go func(teamName string) {
				defer wg.Done()
				params := &types.TeamParams{
					Name:        teamName,
					Description: "Test team",
					Permission:  "push",
				}
				_, err := service.CreateTeam(org, params)
				if err != nil {
					t.Errorf("Failed to create team %s: %v", teamName, err)
				}
			}(team)
		}
		wg.Wait()

		// Verify all teams were created
		teamList, err := service.ListTeams(org)
		if err != nil {
			t.Fatalf("Failed to list teams: %v", err)
		}
		if len(teamList) != len(teams) {
			t.Errorf("Expected %d teams, got %d", len(teams), len(teamList))
		}
	})

	t.Run("Concurrent Team Member Operations", func(t *testing.T) {
		const org = "testorg"
		const team = "team1"
		users := []string{"user1", "user2", "user3"}

		var wg sync.WaitGroup
		wg.Add(len(users))

		// Add team members concurrently
		for _, user := range users {
			go func(username string) {
				defer wg.Done()
				err := service.AddTeamMember(org, team, username)
				if err != nil {
					t.Errorf("Failed to add team member %s: %v", username, err)
				}
			}(user)
		}
		wg.Wait()

		// Remove team members concurrently
		wg.Add(len(users))
		for _, user := range users {
			go func(username string) {
				defer wg.Done()
				err := service.RemoveTeamMember(org, team, username)
				if err != nil {
					t.Errorf("Failed to remove team member %s: %v", username, err)
				}
			}(user)
		}
		wg.Wait()
	})

	t.Run("Concurrent Repository Operations", func(t *testing.T) {
		const org = "testorg"
		const team = "team1"
		repos := []string{"repo1", "repo2", "repo3"}

		var wg sync.WaitGroup
		wg.Add(len(repos))

		// Add repositories concurrently
		for _, repo := range repos {
			go func(repoName string) {
				defer wg.Done()
				err := service.AddTeamRepo(org, team, repoName)
				if err != nil {
					t.Errorf("Failed to add team repo %s: %v", repoName, err)
				}
			}(repo)
		}
		wg.Wait()

		// Verify repositories were added
		repoList, err := service.ListTeamRepos(org, team)
		if err != nil {
			t.Fatalf("Failed to list team repos: %v", err)
		}
		if len(repoList) != len(repos) {
			t.Errorf("Expected %d repos, got %d", len(repos), len(repoList))
		}
	})
}
