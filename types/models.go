// Package types provides the common type definitions for the application
package types

import "time"

// Organization represents a GitHub organization
type Organization struct {
	ID                  int64        `json:"id"`
	Login               string       `json:"login"`
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	HTMLURL             string       `json:"html_url"`
	AvatarURL           string       `json:"avatar_url"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
	TotalPrivateRepos   int          `json:"total_private_repos"`
	OwnedPrivateRepos   int          `json:"owned_private_repos"`   // Added from organization.go
	PublicRepos         int          `json:"public_repos"`          // Added from organization.go
	PublicGists         int          `json:"public_gists"`          // Added from organization.go
	PrivateGists        int          `json:"private_gists"`         // Added from organization.go
	Members             int          `json:"members"`               // Added from organization.go
	DefaultRepoSettings RepoSettings `json:"default_repo_settings"` // Renamed from Settings & matched organization.go
}

// Member represents a member of an organization
type Member struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	HTMLURL   string    `json:"html_url"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Team represents a GitHub team
type Team struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	Permission   string    `json:"permission"`
	Parent       *Team     `json:"parent,omitempty"`
	HTMLURL      string    `json:"html_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	MembersCount int       `json:"members_count"`
	ReposCount   int       `json:"repos_count"`
}

// TeamParams represents parameters for creating/updating a team
type TeamParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Permission  string `json:"permission,omitempty"`
	ParentID    int64  `json:"parent_team_id,omitempty"`
}

// Repository represents a GitHub repository
type Repository struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Private     bool      `json:"private"`
	HTMLURL     string    `json:"html_url"`
	CloneURL    string    `json:"clone_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RepoSettings represents organization repository settings
type RepoSettings struct {
	DefaultRepoPermission string `json:"default_repo_permission"`
	MembersCanCreateRepos bool   `json:"members_can_create_repos"`
	MembersCanCreatePages bool   `json:"members_can_create_pages"`
	MembersCanForkPrivate bool   `json:"members_can_fork_private"`
	TwoFactorRequired     bool   `json:"two_factor_required"`
}
