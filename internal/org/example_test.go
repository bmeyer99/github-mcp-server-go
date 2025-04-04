package org_test

import (
	"fmt"
	"log"

	"github-mcp-server-go/internal/org"
	"github-mcp-server-go/internal/types"
)

func Example() {
	// Initialize organization and team APIs (implementation not shown)
	var orgAPI types.OrganizationAPI
	var teamAPI types.TeamAPI

	// Create a new organization service
	svc := org.New(orgAPI, teamAPI)

	// List organizations
	orgs, err := svc.ListOrganizations()
	if err != nil {
		log.Fatalf("Failed to list organizations: %v", err)
	}

	for _, org := range orgs {
		fmt.Printf("Organization: %s\n", org.Name)

		// List members of the organization
		members, err := svc.ListMembers(org.Login)
		if err != nil {
			log.Printf("Failed to list members for %s: %v", org.Login, err)
			continue
		}

		for _, member := range members {
			fmt.Printf("  Member: %s (Role: %s)\n", member.Login, member.Role)
		}

		// List teams in the organization
		teams, err := svc.ListTeams(org.Login)
		if err != nil {
			log.Printf("Failed to list teams for %s: %v", org.Login, err)
			continue
		}

		for _, team := range teams {
			fmt.Printf("  Team: %s (Permission: %s)\n", team.Name, team.Permission)
		}
	}
}

func ExampleService_CreateTeam() {
	// Initialize APIs (implementation not shown)
	var orgAPI types.OrganizationAPI
	var teamAPI types.TeamAPI

	svc := org.New(orgAPI, teamAPI)

	// Create a new team
	params := &types.TeamParams{
		Name:        "engineering",
		Description: "Engineering team",
		Permission:  "push",
	}

	team, err := svc.CreateTeam("my-org", params)
	if err != nil {
		log.Fatalf("Failed to create team: %v", err)
	}

	fmt.Printf("Created team %s with permission %s\n", team.Name, team.Permission)
}

func ExampleService_AddTeamMember() {
	// Initialize APIs (implementation not shown)
	var orgAPI types.OrganizationAPI
	var teamAPI types.TeamAPI

	svc := org.New(orgAPI, teamAPI)

	// Add a member to a team
	err := svc.AddTeamMember("my-org", "engineering", "jsmith")
	if err != nil {
		log.Fatalf("Failed to add team member: %v", err)
	}

	fmt.Println("Successfully added member to team")
}

func ExampleService_UpdateSettings() {
	// Initialize APIs (implementation not shown)
	var orgAPI types.OrganizationAPI
	var teamAPI types.TeamAPI

	svc := org.New(orgAPI, teamAPI)

	// Update organization settings
	settings := &types.RepoSettings{
		DefaultRepoPermission: "write",
		MembersCanCreateRepos: true,
		TwoFactorRequired:     true,
	}

	err := svc.UpdateSettings("my-org", settings)
	if err != nil {
		log.Fatalf("Failed to update settings: %v", err)
	}

	fmt.Println("Successfully updated organization settings")
}
