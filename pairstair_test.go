package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gypsydave5/pairstair/internal/git"
	"github.com/gypsydave5/pairstair/internal/pairing"
	"github.com/gypsydave5/pairstair/internal/recommend"
	"github.com/gypsydave5/pairstair/internal/team"
)

func TestParseCoAuthors(t *testing.T) {
	body := `
Some commit message

Co-authored-by: Alice <alice@example.com>
Co-authored-by: Bob <bob@example.com>
`
	// Use the git package function directly now that types are aliases
	coauthors := git.ParseCoAuthors(body)

	if len(coauthors) != 2 {
		t.Fatalf("expected 2 coauthors, got %d", len(coauthors))
	}
	if coauthors[0].DisplayName != "Alice" {
		t.Errorf("unexpected coauthor name: %s", coauthors[0].DisplayName)
	}
	if coauthors[0].CanonicalEmail() != "alice@example.com" {
		t.Errorf("unexpected coauthor email: %s", coauthors[0].CanonicalEmail())
	}
	if coauthors[1].DisplayName != "Bob" {
		t.Errorf("unexpected coauthor name: %s", coauthors[1].DisplayName)
	}
	if coauthors[1].CanonicalEmail() != "bob@example.com" {
		t.Errorf("unexpected coauthor email: %s", coauthors[1].CanonicalEmail())
	}
}

func TestMatrixLogic(t *testing.T) {
	commits := []git.Commit{
		{
			Date:      time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
			Author:    git.NewDeveloper("Alice <alice@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Bob <bob@example.com>")},
		},
		{
			Date:      time.Date(2024, 6, 1, 15, 0, 0, 0, time.UTC),
			Author:    git.NewDeveloper("Bob <bob@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Alice <alice@example.com>")},
		},
		{
			Date:      time.Date(2024, 6, 2, 10, 0, 0, 0, time.UTC),
			Author:    git.NewDeveloper("Alice <alice@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Carol <carol@example.com>")},
		},
	}

	matrix, _, _ := pairing.BuildPairMatrix(team.Empty, commits, false)

	// Alice/Bob should have 1 (same day, only count once)
	a, b := "alice@example.com", "bob@example.com"
	if matrix.Count(a, b) != 1 {
		t.Errorf("expected Alice/Bob to have 1, got %d", matrix.Count(a, b))
	}
	// Alice/Carol should have 1
	c := "carol@example.com"
	if matrix.Count(a, c) != 1 {
		t.Errorf("expected Alice/Carol to have 1, got %d", matrix.Count(a, c))
	}
}

func TestMultipleEmailsInTeamFile(t *testing.T) {
	// Team with Alice having multiple email addresses
	developers := []git.Developer{
		git.NewDeveloper("Alice <alice@example.com>,<alice.work@company.com>"),
		git.NewDeveloper("Bob <bob@example.com>"),
	}
	teamObj := team.NewTeamFromDevelopers(developers)

	// Commits with Alice using different email addresses
	commits := []git.Commit{
		{
			Date:      time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
			Author:    git.NewDeveloper("Alice <alice@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Bob <bob@example.com>")},
		},
		{
			Date:      time.Date(2024, 6, 2, 10, 0, 0, 0, time.UTC),
			Author:    git.NewDeveloper("Alice <alice.work@company.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Bob <bob@example.com>")},
		},
	}

	matrix, _, developers := pairing.BuildPairMatrix(teamObj, commits, true)

	// We should only have 2 developers (Alice and Bob), not 3
	if len(developers) != 2 {
		t.Errorf("expected 2 developers, got %d: %v", len(developers), developers)
	}

	// Find Alice's canonical email (should be the first one in the team entry)
	var aliceEmail string
	for _, dev := range developers {
		if strings.Contains(dev.CanonicalEmail(), "alice") {
			aliceEmail = dev.CanonicalEmail()
			break
		}
	}

	if aliceEmail == "" {
		t.Fatalf("could not find Alice in developers list: %v", developers)
	}

	// Check that both commits are counted as Alice pairing with Bob
	bobEmail := "bob@example.com"

	// Should have 2 pairs (one from each day)
	if matrix.Count(aliceEmail, bobEmail) != 2 {
		t.Errorf("expected Alice/Bob pair to have count 2, got %d", matrix.Count(aliceEmail, bobEmail))
	}
}

func TestTeamFileCanonicalName(t *testing.T) {
	// Team with a canonical name
	developers := []git.Developer{
		git.NewDeveloper("Canonical Alice <alice@example.com>"),
		git.NewDeveloper("Bob <bob@example.com>"),
	}
	teamObj := team.NewTeamFromDevelopers(developers)

	// Commits with a different name for Alice
	commits := []git.Commit{
		{
			Date:      time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
			Author:    git.NewDeveloper("Different Alice <alice@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Bob <bob@example.com>")},
		},
	}

	// Build the matrix with useTeam=true
	_, _, developers = pairing.BuildPairMatrix(teamObj, commits, true)

	// Find Alice in the developers list
	var alice *git.Developer
	for _, dev := range developers {
		if dev.CanonicalEmail() == "alice@example.com" {
			alice = &dev
			break
		}
	}

	if alice == nil {
		t.Fatalf("could not find Alice in developers list: %v", developers)
	}

	// Check that Alice's name is the canonical one from the team file
	if alice.DisplayName != "Canonical Alice" {
		t.Errorf("expected name to be 'Canonical Alice', got '%s'", alice.DisplayName)
	}

	// Find Bob and verify his name is preserved
	var bob *git.Developer
	for _, dev := range developers {
		if dev.CanonicalEmail() == "bob@example.com" {
			bob = &dev
			break
		}
	}

	if bob == nil {
		t.Fatalf("could not find Bob in developers list: %v", developers)
	}

	if bob.DisplayName != "Bob" {
		t.Errorf("expected name to be 'Bob', got '%s'", bob.DisplayName)
	}
}

func TestMultipleAuthorsInCommit(t *testing.T) {
	// A single commit with three authors (one main author and two co-authors)
	commits := []git.Commit{
		{
			Date:   time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
			Author: git.NewDeveloper("Alice <alice@example.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Bob <bob@example.com>"),
				git.NewDeveloper("Carol <carol@example.com>"),
			},
		},
	}

	matrix, _, _ := pairing.BuildPairMatrix(team.Empty, commits, false)

	// With 3 authors, we should have 3 pairs: (Alice, Bob), (Alice, Carol), (Bob, Carol)
	if matrix.Len() != 3 {
		t.Errorf("expected 3 pairs in matrix, got %d", matrix.Len())
	}

	// Check each pair exists with count 1
	a, b, c := "alice@example.com", "bob@example.com", "carol@example.com"

	// Alice-Bob pair
	if matrix.Count(a, b) != 1 {
		t.Errorf("expected Alice/Bob pair to have count 1, got %d", matrix.Count(a, b))
	}

	// Alice-Carol pair
	if matrix.Count(a, c) != 1 {
		t.Errorf("expected Alice/Carol pair to have count 1, got %d", matrix.Count(a, c))
	}

	// Bob-Carol pair
	if matrix.Count(b, c) != 1 {
		t.Errorf("expected Bob/Carol pair to have count 1, got %d", matrix.Count(b, c))
	}
}

func TestComprehensivePairMatrix(t *testing.T) {
	// Create a large team with developers having multiple email addresses
	developers := []git.Developer{
		git.NewDeveloper("Alice Smith <alice@example.com>,<alice.smith@company.com>,<asmith@personal.net>"),
		git.NewDeveloper("Bob Jones <bob@example.com>,<bjones@company.com>"),
		git.NewDeveloper("Carol Davis <carol@example.com>,<cdavis@company.com>"),
		git.NewDeveloper("Dave Wilson <dave@example.com>"),
		git.NewDeveloper("Eve Brown <eve@example.com>,<ebrown@company.com>"),
		git.NewDeveloper("Frank Thomas <frank@example.com>"),
	}
	teamObj := team.NewTeamFromDevelopers(developers)

	// Create a comprehensive set of commits covering various scenarios
	now := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	day := 24 * time.Hour

	commits := []git.Commit{
		// Day 1: Alice pairs with Bob
		{
			Date:      now.Add(-14 * day),
			Author:    git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Bob Jones <bob@example.com>")},
		},
		// Day 2: Bob pairs with Carol
		{
			Date:      now.Add(-13 * day),
			Author:    git.NewDeveloper("Bob Jones <bjones@company.com>"), // Different email
			CoAuthors: []git.Developer{git.NewDeveloper("Carol Davis <cdavis@company.com>")},
		},
		// Day 3: Alice pairs with Carol and Dave (three-way pairing)
		{
			Date:   now.Add(-12 * day),
			Author: git.NewDeveloper("Alice Smith <alice.smith@company.com>"), // Different email
			CoAuthors: []git.Developer{
				git.NewDeveloper("Carol Davis <carol@example.com>"),
				git.NewDeveloper("Dave Wilson <dave@example.com>"),
			},
		},
		// Day 4: Eve pairs with Frank
		{
			Date:      now.Add(-11 * day),
			Author:    git.NewDeveloper("Eve Brown <ebrown@company.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Frank Thomas <frank@example.com>")},
		},
		// Day 5: Alice pairs with Eve
		{
			Date:      now.Add(-10 * day),
			Author:    git.NewDeveloper("Alice Smith <asmith@personal.net>"), // Different email
			CoAuthors: []git.Developer{git.NewDeveloper("Eve Brown <eve@example.com>")},
		},
		// Day 6: Dave pairs with Frank
		{
			Date:      now.Add(-9 * day),
			Author:    git.NewDeveloper("Dave Wilson <dave@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Frank Thomas <frank@example.com>")},
		},
		// Day 7: Bob pairs with Dave
		{
			Date:      now.Add(-8 * day),
			Author:    git.NewDeveloper("Bob Jones <bob@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Dave Wilson <dave@example.com>")},
		},
		// Day 7: Also, Carol pairs with Eve (same day, different pair)
		{
			Date:      now.Add(-8 * day),
			Author:    git.NewDeveloper("Carol Davis <carol@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Eve Brown <ebrown@company.com>")},
		},
		// Day 8: Alice pairs with Frank
		{
			Date:      now.Add(-7 * day),
			Author:    git.NewDeveloper("Frank Thomas <frank@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Alice Smith <alice@example.com>")},
		},
		// Day 9-10: No commits

		// Day 11: Multiple commits for the same pair on the same day (should count once)
		{
			Date:      now.Add(-4 * day),
			Author:    git.NewDeveloper("Bob Jones <bob@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Carol Davis <carol@example.com>")},
		},
		{
			Date:      now.Add(-4 * day),
			Author:    git.NewDeveloper("Carol Davis <cdavis@company.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Bob Jones <bjones@company.com>")},
		},
		// Day 12: External person not in the team (should be filtered out with useTeam=true)
		{
			Date:      now.Add(-3 * day),
			Author:    git.NewDeveloper("External Person <external@othercompany.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Alice Smith <alice@example.com>")},
		},
		// Day 13: Alice pairs with Bob again
		{
			Date:      now.Add(-2 * day),
			Author:    git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Bob Jones <bob@example.com>")},
		},
		// Day 14: All team members collaborate (large pairing)
		{
			Date:   now.Add(-1 * day),
			Author: git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Bob Jones <bob@example.com>"),
				git.NewDeveloper("Carol Davis <carol@example.com>"),
				git.NewDeveloper("Dave Wilson <dave@example.com>"),
				git.NewDeveloper("Eve Brown <eve@example.com>"),
				git.NewDeveloper("Frank Thomas <frank@example.com>"),
			},
		},
	}

	// Test with team information
	matrix, _, developers := pairing.BuildPairMatrix(teamObj, commits, true)

	// Check number of developers
	if len(developers) != 6 {
		t.Errorf("expected 6 developers, got %d: %v", len(developers), developers)
	}

	// Check expected pair counts
	expectedPairs := map[string]map[string]int{
		"alice@example.com": {
			"bob@example.com":   3, // Day 1, Day 13, Day 14
			"carol@example.com": 2, // Day 3, Day 14
			"dave@example.com":  2, // Day 3, Day 14
			"eve@example.com":   2, // Day 5, Day 14
			"frank@example.com": 2, // Day 8, Day 14
		},
		"bob@example.com": {
			"carol@example.com": 3, // Day 2, Day 11 (counts once), Day 14
			"dave@example.com":  2, // Day 7, Day 14
			"eve@example.com":   1, // Day 14
			"frank@example.com": 1, // Day 14
		},
		"carol@example.com": {
			"dave@example.com":  2, // Day 3 (three-way pairing), Day 14
			"eve@example.com":   2, // Day 8, Day 14
			"frank@example.com": 1, // Day 14
		},
		"dave@example.com": {
			"eve@example.com":   1, // Day 14
			"frank@example.com": 2, // Day 6, Day 14
		},
		"eve@example.com": {
			"frank@example.com": 2, // Day 4, Day 14
		},
	}

	for dev1, pairs := range expectedPairs {
		for dev2, expectedCount := range pairs {
			actualCount := matrix.Count(dev1, dev2)
			if actualCount != expectedCount {
				t.Errorf("pair %s/%s: expected count %d, got %d", dev1, dev2, expectedCount, actualCount)
			}
		}
	}

	// Verify developer names are properly extracted
	expectedNames := map[string]string{
		"alice@example.com": "Alice Smith",
		"bob@example.com":   "Bob Jones",
		"carol@example.com": "Carol Davis",
		"dave@example.com":  "Dave Wilson",
		"eve@example.com":   "Eve Brown",
		"frank@example.com": "Frank Thomas",
	}

	for _, dev := range developers {
		email := dev.CanonicalEmail()
		if expectedName, exists := expectedNames[email]; exists {
			if dev.DisplayName != expectedName {
				t.Errorf("developer %s: expected name %q, got %q", email, expectedName, dev.DisplayName)
			}
		}
	}

	// Now test without team information
	matrixNoTeam, _, developersNoTeam := pairing.BuildPairMatrix(team.Team{}, commits, false)

	// We expect more developers here because without team info, we don't consolidate alternate emails
	expectedNonTeamDevsCount := 12 // All unique email addresses appear as separate developers
	if len(developersNoTeam) != expectedNonTeamDevsCount {
		t.Errorf("expected %d developers with no team filter, got %d: %v",
			expectedNonTeamDevsCount, len(developersNoTeam), developersNoTeam)
	}

	// Check that external email has proper name
	var externalDev *git.Developer
	for _, dev := range developersNoTeam {
		if strings.Contains(dev.CanonicalEmail(), "external") {
			externalDev = &dev
			break
		}
	}

	if externalDev != nil {
		if externalDev.DisplayName != "External Person" {
			t.Errorf("external email %s should have name 'External Person', got %q", externalDev.CanonicalEmail(), externalDev.DisplayName)
		}
	} else {
		t.Error("external email not found in no-team developers list")
	}

	// Verify the Alice-External pair exists in the no-team matrix
	if externalDev != nil {
		aliceEmail := "alice@example.com"
		if matrixNoTeam.Count(aliceEmail, externalDev.CanonicalEmail()) != 1 {
			t.Errorf("expected Alice-External pair to have count 1 in no-team matrix")
		}
	}
}

func TestLeastRecentStrategy(t *testing.T) {
	// Create commits with different dates
	commits := []git.Commit{
		{
			Date:      time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC), // Most recent
			Author:    git.NewDeveloper("Alice <alice@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Bob <bob@example.com>")},
		},
		{
			Date:      time.Date(2024, 5, 15, 10, 0, 0, 0, time.UTC), // Less recent
			Author:    git.NewDeveloper("Alice <alice@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Carol <carol@example.com>")},
		},
		{
			Date:      time.Date(2024, 5, 10, 10, 0, 0, 0, time.UTC), // Least recent
			Author:    git.NewDeveloper("Bob <bob@example.com>"),
			CoAuthors: []git.Developer{git.NewDeveloper("Dave <dave@example.com>")},
		},
	}

	matrix, recencyMatrix, developers := pairing.BuildPairMatrix(team.Empty, commits, false)

	// Test recency tracking
	aliceEmail := "alice@example.com"
	bobEmail := "bob@example.com"
	carolEmail := "carol@example.com"

	// Check that recency data is correct
	lastPairedAB, hasDataAB := recencyMatrix.LastPaired(aliceEmail, bobEmail)
	if !hasDataAB {
		t.Error("expected Alice-Bob to have recency data")
	}
	expectedDateAB := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	if !lastPairedAB.Equal(expectedDateAB) {
		t.Errorf("expected Alice-Bob last paired on %v, got %v", expectedDateAB, lastPairedAB)
	}

	lastPairedAC, hasDataAC := recencyMatrix.LastPaired(aliceEmail, carolEmail)
	if !hasDataAC {
		t.Error("expected Alice-Carol to have recency data")
	}
	expectedDateAC := time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC)
	if !lastPairedAC.Equal(expectedDateAC) {
		t.Errorf("expected Alice-Carol last paired on %v, got %v", expectedDateAC, lastPairedAC)
	}

	// Test recommendations using least-recent strategy
	recommendations := recommend.GenerateRecommendations(developers, matrix, recencyMatrix, recommend.LeastRecent)

	// Should recommend pairs that haven't worked together or worked together longest ago
	if len(recommendations) < 2 {
		t.Errorf("expected at least 2 recommendations, got %d", len(recommendations))
	}

	// First recommendation should be for pairs that never worked together
	// or the least recently paired
	foundNeverPaired := false
	for _, rec := range recommendations {
		if !rec.HasPaired {
			foundNeverPaired = true
			break
		}
	}

	if !foundNeverPaired {
		// If all pairs have worked together, check that least recent is first
		firstRec := recommendations[0]
		if firstRec.DaysSince <= 0 {
			t.Error("expected first recommendation to have positive days since or be never paired")
		}
	}
}

func TestReadTeamFileWithSubTeams(t *testing.T) {
	// Create a temporary file with sub-team sections
	content := `Alice Example <alice@example.com>
Bob Dev <bob@example.com>

[frontend]
Carol Frontend <carol@example.com>
Dave UI <dave@example.com>

[backend]
Eve Backend <eve@example.com>
Frank API <frank@example.com>

[devops]
Grace Ops <grace@example.com>
`

	tempDir := t.TempDir()
	teamFile := filepath.Join(tempDir, ".team")
	err := ioutil.WriteFile(teamFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Test reading entire team (no sub-team specified)
	teamMembers, err := team.ReadTeamFile(teamFile, "")
	if err != nil {
		t.Fatalf("Failed to read team file: %v", err)
	}
	expected := []string{"Alice Example <alice@example.com>", "Bob Dev <bob@example.com>"}
	if len(teamMembers) != len(expected) {
		t.Fatalf("Expected %d members, got %d", len(expected), len(teamMembers))
	}
	for i, member := range expected {
		if teamMembers[i] != member {
			t.Errorf("Expected member %q, got %q", member, teamMembers[i])
		}
	}

	// Test reading frontend sub-team
	frontendTeam, err := team.ReadTeamFile(teamFile, "frontend")
	if err != nil {
		t.Fatalf("Failed to read frontend team: %v", err)
	}
	expectedFrontend := []string{"Carol Frontend <carol@example.com>", "Dave UI <dave@example.com>"}
	if len(frontendTeam) != len(expectedFrontend) {
		t.Fatalf("Expected %d frontend members, got %d", len(expectedFrontend), len(frontendTeam))
	}
	for i, member := range expectedFrontend {
		if frontendTeam[i] != member {
			t.Errorf("Expected frontend member %q, got %q", member, frontendTeam[i])
		}
	}

	// Test reading backend sub-team
	backendTeam, err := team.ReadTeamFile(teamFile, "backend")
	if err != nil {
		t.Fatalf("Failed to read backend team: %v", err)
	}
	expectedBackend := []string{"Eve Backend <eve@example.com>", "Frank API <frank@example.com>"}
	if len(backendTeam) != len(expectedBackend) {
		t.Fatalf("Expected %d backend members, got %d", len(expectedBackend), len(backendTeam))
	}
	for i, member := range expectedBackend {
		if backendTeam[i] != member {
			t.Errorf("Expected backend member %q, got %q", member, backendTeam[i])
		}
	}

	// Test reading non-existent sub-team
	nonExistentTeam, err := team.ReadTeamFile(teamFile, "nonexistent")
	if err != nil {
		t.Fatalf("Failed to read team file: %v", err)
	}
	if len(nonExistentTeam) != 0 {
		t.Errorf("Expected empty team for non-existent sub-team, got %d members", len(nonExistentTeam))
	}
}

func TestNewTeamFromFileWithSubTeam(t *testing.T) {
	content := `Alice Example <alice@example.com>
Bob Dev <bob@example.com>

[frontend]
Carol Frontend <carol@example.com>
Dave UI <dave@example.com>

[backend]
Eve Backend <eve@example.com>
Frank API <frank@example.com>
`

	tempDir := t.TempDir()
	teamFile := filepath.Join(tempDir, ".team")
	err := ioutil.WriteFile(teamFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Test creating team from frontend sub-team
	teamObj, err := team.NewTeamFromFile(teamFile, "frontend")
	if err != nil {
		t.Fatalf("Failed to create team from file: %v", err)
	}

	// Verify that only frontend developers are included
	if !teamObj.HasDeveloperByEmail("carol@example.com") {
		t.Error("Expected Carol to be in frontend team")
	}
	if !teamObj.HasDeveloperByEmail("dave@example.com") {
		t.Error("Expected Dave to be in frontend team")
	}
	if teamObj.HasDeveloperByEmail("eve@example.com") {
		t.Error("Expected Eve NOT to be in frontend team")
	}
	if teamObj.HasDeveloperByEmail("alice@example.com") {
		t.Error("Expected Alice NOT to be in frontend team")
	}
}

func TestDeveloperInMultipleSubTeams(t *testing.T) {
	content := `Alice Lead <alice@example.com>

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

	// Test frontend team - should include Bob and Carol
	frontendTeam, err := team.NewTeamFromFile(teamFile, "frontend")
	if err != nil {
		t.Fatalf("Failed to create frontend team: %v", err)
	}
	if !frontendTeam.HasDeveloperByEmail("bob@example.com") {
		t.Error("Expected Bob to be in frontend team")
	}
	if !frontendTeam.HasDeveloperByEmail("carol@example.com") {
		t.Error("Expected Carol to be in frontend team")
	}
	if frontendTeam.HasDeveloperByEmail("dave@example.com") {
		t.Error("Expected Dave NOT to be in frontend team")
	}

	// Test backend team - should include Bob and Dave
	backendTeam, err := team.NewTeamFromFile(teamFile, "backend")
	if err != nil {
		t.Fatalf("Failed to create backend team: %v", err)
	}
	if !backendTeam.HasDeveloperByEmail("bob@example.com") {
		t.Error("Expected Bob to be in backend team")
	}
	if !backendTeam.HasDeveloperByEmail("dave@example.com") {
		t.Error("Expected Dave to be in backend team")
	}
	if backendTeam.HasDeveloperByEmail("carol@example.com") {
		t.Error("Expected Carol NOT to be in backend team")
	}

	// Test main team - should only include Alice
	mainTeam, err := team.NewTeamFromFile(teamFile, "")
	if err != nil {
		t.Fatalf("Failed to create main team: %v", err)
	}
	if !mainTeam.HasDeveloperByEmail("alice@example.com") {
		t.Error("Expected Alice to be in main team")
	}
	if mainTeam.HasDeveloperByEmail("bob@example.com") {
		t.Error("Expected Bob NOT to be in main team (should be in sub-teams only)")
	}
}

func TestMainTeamWithDuplicatesInSubTeams(t *testing.T) {
	content := `Alice Main <alice@example.com>
Bob BothMainAndSub <bob@example.com>

[frontend]
Bob BothMainAndSub <bob@example.com>
Carol SubTeamOnly <carol@example.com>

[backend]  
Bob BothMainAndSub <bob@example.com>
Carol SubTeamOnly <carol@example.com>
Dave SubTeamOnly <dave@example.com>
`

	tempDir := t.TempDir()
	teamFile := filepath.Join(tempDir, ".team")
	err := ioutil.WriteFile(teamFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Test main team (no sub-team specified)
	mainTeam, err := team.NewTeamFromFile(teamFile, "")
	if err != nil {
		t.Fatalf("Failed to create main team: %v", err)
	}

	// Alice should be in main team (only listed in main)
	if !mainTeam.HasDeveloperByEmail("alice@example.com") {
		t.Error("Expected Alice to be in main team")
	}

	// Bob should be in main team (listed in both main and sub-teams)
	if !mainTeam.HasDeveloperByEmail("bob@example.com") {
		t.Error("Expected Bob to be in main team")
	}

	// Carol should NOT be in main team (only listed in sub-teams)
	if mainTeam.HasDeveloperByEmail("carol@example.com") {
		t.Error("Expected Carol NOT to be in main team (only in sub-teams)")
	}

	// Dave should NOT be in main team (only listed in sub-teams)
	if mainTeam.HasDeveloperByEmail("dave@example.com") {
		t.Error("Expected Dave NOT to be in main team (only in sub-teams)")
	}

	// Verify Bob appears only once (no duplication despite being in multiple sub-teams)
	emailToName, _ := mainTeam.GetEmailMappings()
	bobCount := 0
	for email := range emailToName {
		if email == "bob@example.com" {
			bobCount++
		}
	}
	if bobCount != 1 {
		t.Errorf("Expected Bob to appear exactly once in main team, got %d times", bobCount)
	}
}

func TestCoAuthorPairingDetection(t *testing.T) {
	// This test verifies that PairStair correctly identifies pairing from a real commit scenario
	// where the author email differs from the team file email but co-authored-by includes both devs

	// Create team with Ahmad and Tamara using specific emails
	// Include both Tamara emails to match the real-world scenario where team files
	// should list all email variations for each developer
	developers := []git.Developer{
		git.NewDeveloper("Ahmad Qurbanzada <ahmad.qurbanzada@springernature.com>"),
		git.NewDeveloper("Tamara Jordan <20561445+tamj0rd2@users.noreply.github.com>,<tamara.jordan@springernature.com>"),
	}
	teamObj := team.NewTeamFromDevelopers(developers)

	// Create a commit that matches the real scenario:
	// - Author: Tamara with work email (not in team file)
	// - Co-authored-by: Ahmad with work email (matches team file)
	// - Co-authored-by: Tamara with GitHub email (matches team file)
	commits := []git.Commit{
		{
			Date:   time.Date(2025, 6, 26, 16, 33, 50, 0, time.UTC),
			Author: git.NewDeveloper("Tamara Jordan <tamara.jordan@springernature.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Ahmad Qurbanzada <ahmad.qurbanzada@springernature.com>"),
				git.NewDeveloper("Tamara Jordan <20561445+tamj0rd2@users.noreply.github.com>"),
			},
		},
	}

	// Build pair matrix with team enabled
	matrix, _, developers := pairing.BuildPairMatrix(teamObj, commits, true)

	// Debug: print what we got
	t.Logf("Developers found: %v", developers)
	t.Logf("Matrix pairs: %d", matrix.Len())

	// Should have exactly 2 developers in the final result
	if len(developers) != 2 {
		t.Errorf("expected 2 developers, got %d: %v", len(developers), developers)
	}

	// Should have exactly 1 pair (Ahmad and Tamara)
	if matrix.Len() != 1 {
		t.Errorf("expected 1 pair in matrix, got %d", matrix.Len())
	}

	// Find Ahmad and Tamara's emails in the result
	var ahmadEmail, tamaraEmail string
	for _, dev := range developers {
		email := dev.CanonicalEmail()
		if strings.Contains(email, "ahmad") {
			ahmadEmail = email
		} else if strings.Contains(email, "tamj0rd2") || strings.Contains(email, "tamara") {
			tamaraEmail = email
		}
	}

	if ahmadEmail == "" {
		t.Fatalf("could not find Ahmad in developers list: %v", developers)
	}
	if tamaraEmail == "" {
		t.Fatalf("could not find Tamara in developers list: %v", developers)
	}

	// Verify that Ahmad and Tamara are counted as having paired once
	pairCount := matrix.Count(ahmadEmail, tamaraEmail)
	if pairCount != 1 {
		t.Errorf("expected Ahmad and Tamara to have paired 1 time, got %d", pairCount)
	}
}

// Tests for configuration and flag parsing
func TestConfigDefaults(t *testing.T) {
	// Test that Config has the expected default behavior
	config := &Config{}

	// The new Open field should default to false (stream to stdout)
	if config.Open {
		t.Error("Config.Open should default to false")
	}
}

func TestConfigOpenFlag(t *testing.T) {
	// Test that the -open flag changes the behavior of HTML output
	// This is a unit test for the configuration structure
	config := &Config{
		Output: "html",
		Open:   true,
	}

	if !config.Open {
		t.Error("Expected Open to be true when set")
	}

	if config.Output != "html" {
		t.Error("Expected Output to be html")
	}
}
