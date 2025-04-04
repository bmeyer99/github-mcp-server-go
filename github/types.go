package github

import "time"

type (
	// Organization represents a GitHub organization
	Organization struct {
		ID                int64        `json:"id"`
		Login             string       `json:"login"`
		Name              string       `json:"name"`
		Description       string       `json:"description"`
		HTMLURL           string       `json:"html_url"`
		AvatarURL         string       `json:"avatar_url"`
		CreatedAt         time.Time    `json:"created_at"`
		UpdatedAt         time.Time    `json:"updated_at"`
		TotalPrivateRepos int          `json:"total_private_repos"`
		Settings          RepoSettings `json:"settings"`
	}

	// Member represents a member of an organization
	Member struct {
		ID        int64     `json:"id"`
		Login     string    `json:"login"`
		HTMLURL   string    `json:"html_url"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// Team represents a GitHub team
	Team struct {
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
	TeamParams struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		Permission  string `json:"permission,omitempty"`
		ParentID    int64  `json:"parent_team_id,omitempty"`
	}

	// Repository represents a GitHub repository
	Repository struct {
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
	RepoSettings struct {
		DefaultRepoPermission string `json:"default_repo_permission"`
		MembersCanCreateRepos bool   `json:"members_can_create_repos"`
		MembersCanCreatePages bool   `json:"members_can_create_pages"`
		MembersCanForkPrivate bool   `json:"members_can_fork_private"`
		TwoFactorRequired     bool   `json:"two_factor_required"`
	}
)
