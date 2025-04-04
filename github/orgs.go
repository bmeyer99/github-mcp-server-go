package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github-mcp-server-go/types"
)

// ListOrganizations lists organizations for the authenticated user
func (c *Client) ListOrganizations(ctx context.Context) ([]*types.Organization, error) {
	url := "user/orgs"

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orgs []*types.Organization
	if err := json.NewDecoder(resp.Body).Decode(&orgs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return orgs, nil
}

// GetOrganization gets an organization by name
func (c *Client) GetOrganization(ctx context.Context, name string) (*types.Organization, error) {
	url := fmt.Sprintf("orgs/%s", name)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var organization types.Organization
	if err := json.NewDecoder(resp.Body).Decode(&organization); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &organization, nil
}

// UpdateOrganizationSettings updates organization repository settings
func (c *Client) UpdateOrganizationSettings(ctx context.Context, name string, settings *types.RepoSettings) error {
	url := fmt.Sprintf("orgs/%s", name)

	req, err := c.newRequest(ctx, "PATCH", url, settings)
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

// ListOrganizationMembers lists members of an organization
func (c *Client) ListOrganizationMembers(ctx context.Context, name string) ([]*types.Member, error) {
	url := fmt.Sprintf("orgs/%s/members", name)

	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var members []*types.Member
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return members, nil
}

// AddOrganizationMember adds a user to an organization
func (c *Client) AddOrganizationMember(ctx context.Context, orgName, username, role string) error {
	url := fmt.Sprintf("orgs/%s/memberships/%s", orgName, username)

	payload := map[string]string{
		"role": role,
	}

	req, err := c.newRequest(ctx, "PUT", url, payload)
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

// RemoveOrganizationMember removes a user from an organization
func (c *Client) RemoveOrganizationMember(ctx context.Context, orgName, username string) error {
	url := fmt.Sprintf("orgs/%s/memberships/%s", orgName, username)

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
