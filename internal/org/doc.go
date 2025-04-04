/*
Package org provides a high-level interface for managing GitHub organizations and teams.

This package implements a service layer that simplifies common organization management
tasks such as managing members, teams, and repository access. It provides a unified
interface that handles validation and implements best practices.

Basic usage:

	svc := org.New(orgAPI, teamAPI)

	// List organizations
	orgs, err := svc.ListOrganizations()

	// Add a member
	err = svc.AddMember("my-org", "username", "member")

	// Create and manage teams
	team, err := svc.CreateTeam("my-org", &types.TeamParams{
		Name:        "engineering",
		Description: "Engineering team",
		Permission:  "push",
	})

The package enforces proper role and permission validation, handles nested team
relationships, and provides a consistent interface for both organization and team
operations.

Organization operations include:
  - Managing organization settings
  - Adding/removing members
  - Setting member roles

Team operations include:
  - Creating/updating/deleting teams
  - Managing team membership
  - Controlling repository access
  - Managing nested team hierarchies

All operations use the interfaces defined in the types package, making it easy
to provide alternative implementations or mock clients for testing.
*/
package org
