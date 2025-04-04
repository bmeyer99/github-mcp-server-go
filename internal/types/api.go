package types

// OrganizationAPI defines organization operations
type OrganizationAPI interface {
	ListOrganizations() ([]*Organization, error)
	GetOrganization(name string) (*Organization, error)
	UpdateSettings(name string, settings *RepoSettings) error
	ListMembers(name string) ([]*Member, error)
	AddMember(org, user string, role string) error
	RemoveMember(org, user string) error
}

// TeamAPI defines team operations
type TeamAPI interface {
	ListTeams(org string) ([]*Team, error)
	GetTeam(org string, team string) (*Team, error)
	CreateTeam(org string, params *TeamParams) (*Team, error)
	UpdateTeam(org string, team string, params *TeamParams) error
	DeleteTeam(org string, team string) error
	AddTeamMember(org string, team string, username string) error
	RemoveTeamMember(org string, team string, username string) error
	ListTeamRepos(org string, team string) ([]*Repository, error)
	AddTeamRepo(org string, team string, repo string) error
	RemoveTeamRepo(org string, team string, repo string) error
}
