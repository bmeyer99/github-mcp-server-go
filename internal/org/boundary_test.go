package org

import (
	"testing"

	"github-mcp-server-go/internal/types"
)

func TestBoundaryConditions(t *testing.T) {
	orgAPI := newMockOrgAPIConcurrent()
	teamAPI := newMockTeamAPIConcurrent()
	service := New(orgAPI, teamAPI)

	t.Run("Zero Value Parameters", func(t *testing.T) {
		// Test with zero-value/empty parameters
		tests := []struct {
			name    string
			test    func() error
			wantErr bool
		}{
			{
				name: "Empty organization name",
				test: func() error {
					_, err := service.GetOrganization("")
					return err
				},
				wantErr: true,
			},
			{
				name: "Nil settings",
				test: func() error {
					return service.UpdateSettings("org", nil)
				},
				wantErr: true,
			},
			{
				name: "Empty member role",
				test: func() error {
					return service.AddMember("org", "user", "")
				},
				wantErr: true,
			},
			{
				name: "Empty username",
				test: func() error {
					return service.AddMember("org", "", "member")
				},
				wantErr: true,
			},
			{
				name: "Nil team params",
				test: func() error {
					_, err := service.CreateTeam("org", nil)
					return err
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.test()
				if (err != nil) != tt.wantErr {
					t.Errorf("%s: got error = %v, wantErr = %v", tt.name, err, tt.wantErr)
				}
			})
		}
	})

	t.Run("Maximum Values", func(t *testing.T) {
		// Test with maximum allowed values
		longString := make([]byte, 256)
		for i := range longString {
			longString[i] = 'a'
		}
		longName := string(longString)

		tests := []struct {
			name    string
			test    func() error
			wantErr bool
		}{
			{
				name: "Very long organization name",
				test: func() error {
					_, err := service.GetOrganization(longName)
					return err
				},
				wantErr: true,
			},
			{
				name: "Very long team name",
				test: func() error {
					_, err := service.CreateTeam("org", &types.TeamParams{
						Name: longName,
					})
					return err
				},
				wantErr: true,
			},
			{
				name: "Very long username",
				test: func() error {
					return service.AddMember("org", longName, "member")
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.test()
				if (err != nil) != tt.wantErr {
					t.Errorf("%s: got error = %v, wantErr = %v", tt.name, err, tt.wantErr)
				}
			})
		}
	})

	t.Run("Special Characters", func(t *testing.T) {
		specialChars := []string{
			"test/org",
			"test\\org",
			"test:org",
			"test*org",
			"test?org",
			"test<org>",
			"test|org",
		}

		for _, char := range specialChars {
			t.Run("Organization name with "+char, func(t *testing.T) {
				_, err := service.GetOrganization(char)
				if err == nil {
					t.Errorf("Expected error for organization name with special character: %s", char)
				}
			})

			t.Run("Team name with "+char, func(t *testing.T) {
				_, err := service.CreateTeam("org", &types.TeamParams{Name: char})
				if err == nil {
					t.Errorf("Expected error for team name with special character: %s", char)
				}
			})

			t.Run("Username with "+char, func(t *testing.T) {
				err := service.AddMember("org", char, "member")
				if err == nil {
					t.Errorf("Expected error for username with special character: %s", char)
				}
			})
		}
	})
}
