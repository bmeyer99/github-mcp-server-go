package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github-mcp-server-go/types"
)

// ListTeams lists teams in an organization
func (c *Client) ListTeams(ctx context.Context, orgName string) ([]*types.Team, error) {
	url := fmt.Sprintf("orgs/%s/teams", orgName)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var teams []*types.Team
	if err := json.NewDecoder(resp.Body).Decode(&teams); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return teams, nil
}

// CreateTeam creates a new team in an organization
func (c *Client) CreateTeam(ctx context.Context, orgName string, params *types.TeamParams) (*types.Team, error) {
	url := fmt.Sprintf("orgs/%s/teams", orgName)

	req, err := c.newRequest(ctx, "POST", url, params)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var team types.Team
	if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &team, nil
}

// GetTeam gets a team by slug
func (c *Client) GetTeam(ctx context.Context, orgName, teamSlug string) (*types.Team, error) {
	url := fmt.Sprintf("orgs/%s/teams/%s", orgName, teamSlug)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var team types.Team
	if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &team, nil
}

// UpdateTeam updates a team
func (c *Client) UpdateTeam(ctx context.Context, orgName, teamSlug string, params *types.TeamParams) error {
	url := fmt.Sprintf("orgs/%s/teams/%s", orgName, teamSlug)

	req, err := c.newRequest(ctx, "PATCH", url, params)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteTeam deletes a team
func (c *Client) DeleteTeam(ctx context.Context, orgName, teamSlug string) error {
	url := fmt.Sprintf("orgs/%s/teams/%s", orgName, teamSlug)

	req, err := c.newRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// AddTeamMember adds a user to a team
func (c *Client) AddTeamMember(ctx context.Context, orgName, teamSlug, username string) error {
	url := fmt.Sprintf("orgs/%s/teams/%s/memberships/%s", orgName, teamSlug, username)

	req, err := c.newRequest(ctx, "PUT", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// RemoveTeamMember removes a user from a team
func (c *Client) RemoveTeamMember(ctx context.Context, orgName, teamSlug, username string) error {
	url := fmt.Sprintf("orgs/%s/teams/%s/memberships/%s", orgName, teamSlug, username)

	req, err := c.newRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ListTeamRepos lists repositories a team has access to
func (c *Client) ListTeamRepos(ctx context.Context, orgName, teamSlug string) ([]*types.Repository, error) {
	url := fmt.Sprintf("orgs/%s/teams/%s/repos", orgName, teamSlug)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []*types.Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return repos, nil
}

// AddTeamRepo adds a repository to a team
func (c *Client) AddTeamRepo(ctx context.Context, orgName, teamSlug, repoName string) error {
	url := fmt.Sprintf("orgs/%s/teams/%s/repos/%s/%s", orgName, teamSlug, orgName, repoName)

	req, err := c.newRequest(ctx, "PUT", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// RemoveTeamRepo removes a repository from a team
func (c *Client) RemoveTeamRepo(ctx context.Context, orgName, teamSlug, repoName string) error {
	url := fmt.Sprintf("orgs/%s/teams/%s/repos/%s/%s", orgName, teamSlug, orgName, repoName)

	req, err := c.newRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
