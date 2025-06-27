package team_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/gypsydave5/pairstair/internal/team"
)

func TestNewTeam(t *testing.T) {
	tests := []struct {
		name         string
		teamMembers  []string
		expectedDevs int
		checkEmail   string
		shouldHave   bool
	}{
		{
			name: "basic team creation",
			teamMembers: []string{
				"Alice Smith <alice@example.com>",
				"Bob Jones <bob@example.com>",
			},
			expectedDevs: 2,
			checkEmail:   "alice@example.com",
			shouldHave:   true,
		},
		{
			name: "team with multiple emails per developer",
			teamMembers: []string{
				"Alice Smith <alice@example.com>,<alice@company.com>",
				"Bob Jones <bob@example.com>",
			},
			expectedDevs: 2,
			checkEmail:   "alice@company.com",
			shouldHave:   true,
		},
		{
			name:         "empty team",
			teamMembers:  []string{},
			expectedDevs: 0,
			checkEmail:   "anyone@example.com",
			shouldHave:   false,
		},
		{
			name: "team with invalid entries",
			teamMembers: []string{
				"Alice Smith <alice@example.com>",
				"Invalid Entry Without Email",
				"Bob Jones <bob@example.com>",
			},
			expectedDevs: 2,
			checkEmail:   "alice@example.com",
			shouldHave:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team, err := team.NewTeam(tt.teamMembers)
			if err != nil {
				t.Fatalf("NewTeam() failed: %v", err)
			}

			// Check team size by examining email mappings
			emailToName, _ := team.GetEmailMappings()
			if len(emailToName) < tt.expectedDevs {
				t.Errorf("Expected at least %d email mappings, got %d", tt.expectedDevs, len(emailToName))
			}

			// Check specific email
			hasEmail := team.HasDeveloperByEmail(tt.checkEmail)
			if hasEmail != tt.shouldHave {
				t.Errorf("HasDeveloperByEmail(%q) = %v, expected %v", tt.checkEmail, hasEmail, tt.shouldHave)
			}
		})
	}
}

func TestTeamEmailMappings(t *testing.T) {
	teamMembers := []string{
		"Alice Smith <alice@example.com>,<alice@company.com>",
		"Bob Jones <bob@example.com>",
	}

	team, err := team.NewTeam(teamMembers)
	if err != nil {
		t.Fatalf("NewTeam() failed: %v", err)
	}

	emailToName, emailToPrimary := team.GetEmailMappings()

	// Check email to name mappings
	expectedNames := map[string]string{
		"alice@example.com": "Alice Smith",
		"alice@company.com": "Alice Smith",
		"bob@example.com":   "Bob Jones",
	}

	for email, expectedName := range expectedNames {
		if name, ok := emailToName[email]; !ok {
			t.Errorf("Missing email mapping for %q", email)
		} else if name != expectedName {
			t.Errorf("Email %q mapped to %q, expected %q", email, name, expectedName)
		}
	}

	// Check email to primary email mappings
	expectedPrimary := map[string]string{
		"alice@example.com": "alice@example.com", // First email is primary
		"alice@company.com": "alice@example.com", // Maps to primary
		"bob@example.com":   "bob@example.com",   // Single email, maps to itself
	}

	for email, expectedPrimary := range expectedPrimary {
		if primary, ok := emailToPrimary[email]; !ok {
			t.Errorf("Missing primary email mapping for %q", email)
		} else if primary != expectedPrimary {
			t.Errorf("Email %q maps to primary %q, expected %q", email, primary, expectedPrimary)
		}
	}
}

func TestNewTeamFromFile(t *testing.T) {
	// Create a temporary team file
	content := `Alice Smith <alice@example.com>
Bob Jones <bob@example.com>,<bob@company.com>
Carol White <carol@example.com>

[frontend]
Dave Brown <dave@example.com>
Eve Green <eve@example.com>

[backend]
Frank Black <frank@example.com>
Grace Gray <grace@example.com>
`

	tempDir := t.TempDir()
	teamFile := filepath.Join(tempDir, ".team")
	err := ioutil.WriteFile(teamFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	tests := []struct {
		name         string
		subTeam      string
		expectedDevs []string
		notExpected  []string
	}{
		{
			name:    "main team (no sub-team)",
			subTeam: "",
			expectedDevs: []string{
				"alice@example.com",
				"bob@example.com",
				"carol@example.com",
			},
			notExpected: []string{
				"dave@example.com",
				"eve@example.com",
			},
		},
		{
			name:    "frontend sub-team",
			subTeam: "frontend",
			expectedDevs: []string{
				"dave@example.com",
				"eve@example.com",
			},
			notExpected: []string{
				"alice@example.com",
				"frank@example.com",
			},
		},
		{
			name:    "backend sub-team",
			subTeam: "backend",
			expectedDevs: []string{
				"frank@example.com",
				"grace@example.com",
			},
			notExpected: []string{
				"alice@example.com",
				"dave@example.com",
			},
		},
		{
			name:        "non-existent sub-team",
			subTeam:     "nonexistent",
			expectedDevs: []string{},
			notExpected: []string{
				"alice@example.com",
				"dave@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team, err := team.NewTeamFromFile(teamFile, tt.subTeam)
			if err != nil {
				t.Fatalf("NewTeamFromFile() failed: %v", err)
			}

			// Check expected developers are present
			for _, email := range tt.expectedDevs {
				if !team.HasDeveloperByEmail(email) {
					t.Errorf("Expected developer %q to be in team", email)
				}
			}

			// Check unexpected developers are not present
			for _, email := range tt.notExpected {
				if team.HasDeveloperByEmail(email) {
					t.Errorf("Expected developer %q NOT to be in team", email)
				}
			}
		})
	}
}

func TestTeamFileWithMultipleEmailsPerDeveloper(t *testing.T) {
	content := `Alice Consolidated <alice@work.com>,<alice@personal.com>,<alice@old.com>
Bob Single <bob@example.com>
`

	tempDir := t.TempDir()
	teamFile := filepath.Join(tempDir, ".team")
	err := ioutil.WriteFile(teamFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	team, err := team.NewTeamFromFile(teamFile, "")
	if err != nil {
		t.Fatalf("NewTeamFromFile() failed: %v", err)
	}

	// All of Alice's emails should map to the same person
	aliceEmails := []string{"alice@work.com", "alice@personal.com", "alice@old.com"}
	for _, email := range aliceEmails {
		if !team.HasDeveloperByEmail(email) {
			t.Errorf("Expected Alice's email %q to be recognized", email)
		}
	}

	// Check that all Alice's emails map to the same primary email (first one)
	_, emailToPrimary := team.GetEmailMappings()
	expectedPrimary := "alice@work.com"
	for _, email := range aliceEmails {
		if primary := emailToPrimary[email]; primary != expectedPrimary {
			t.Errorf("Email %q maps to primary %q, expected %q", email, primary, expectedPrimary)
		}
	}
}

func TestTeamFileWithDuplicateSubTeamEntries(t *testing.T) {
	content := `Alice Lead <alice@example.com>
Bob Fullstack <bob@example.com>

[frontend]
Bob Fullstack <bob@example.com>
Carol Frontend <carol@example.com>

[backend]
Bob Fullstack <bob@example.com>
Dave Backend <dave@example.com>
`

	tempDir := t.TempDir()
	teamFile := filepath.Join(tempDir, ".team")
	err := ioutil.WriteFile(teamFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	tests := []struct {
		name         string
		subTeam      string
		shouldHaveBob bool
		otherExpected []string
	}{
		{
			name:          "main team includes Bob",
			subTeam:       "",
			shouldHaveBob: true,
			otherExpected: []string{"alice@example.com"},
		},
		{
			name:          "frontend team includes Bob",
			subTeam:       "frontend",
			shouldHaveBob: true,
			otherExpected: []string{"carol@example.com"},
		},
		{
			name:          "backend team includes Bob",
			subTeam:       "backend",
			shouldHaveBob: true,
			otherExpected: []string{"dave@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team, err := team.NewTeamFromFile(teamFile, tt.subTeam)
			if err != nil {
				t.Fatalf("NewTeamFromFile() failed: %v", err)
			}

			bobEmail := "bob@example.com"
			hasBob := team.HasDeveloperByEmail(bobEmail)
			if hasBob != tt.shouldHaveBob {
				t.Errorf("HasDeveloperByEmail(%q) = %v, expected %v", bobEmail, hasBob, tt.shouldHaveBob)
			}

			for _, email := range tt.otherExpected {
				if !team.HasDeveloperByEmail(email) {
					t.Errorf("Expected developer %q to be in %s team", email, tt.subTeam)
				}
			}
		})
	}
}

func TestReadTeamFileErrorHandling(t *testing.T) {
	// Test non-existent file
	_, err := team.NewTeamFromFile("/nonexistent/path/.team", "")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestTeamFileWithCommentsAndEmptyLines(t *testing.T) {
	content := `# This is a comment
Alice Smith <alice@example.com>

# Another comment
Bob Jones <bob@example.com>

# Comments in sections
[frontend]
# Frontend team members
Carol White <carol@example.com>

Dave Brown <dave@example.com>
`

	tempDir := t.TempDir()
	teamFile := filepath.Join(tempDir, ".team")
	err := ioutil.WriteFile(teamFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Test main team (should handle comments gracefully)
	team, err := team.NewTeamFromFile(teamFile, "")
	if err != nil {
		t.Fatalf("NewTeamFromFile() failed: %v", err)
	}

	expectedEmails := []string{"alice@example.com", "bob@example.com"}
	for _, email := range expectedEmails {
		if !team.HasDeveloperByEmail(email) {
			t.Errorf("Expected developer %q to be in main team", email)
		}
	}

	// Comments starting with # should be ignored
	if team.HasDeveloperByEmail("# This is a comment") {
		t.Error("Comments should be ignored, not treated as developers")
	}
}
