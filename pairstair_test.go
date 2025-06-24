package main

import (
	"strings"
	"testing"
	"time"
)

func TestParseCoAuthors(t *testing.T) {
	body := `
Some commit message

Co-authored-by: Alice <alice@example.com>
Co-authored-by: Bob <bob@example.com>
`
	coauthors := parseCoAuthors(body)
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
	commits := []Commit{
		{
			Date:      time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
			Author:    NewDeveloper("Alice <alice@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Bob <bob@example.com>")},
		},
		{
			Date:      time.Date(2024, 6, 1, 15, 0, 0, 0, time.UTC),
			Author:    NewDeveloper("Bob <bob@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Alice <alice@example.com>")},
		},
		{
			Date:      time.Date(2024, 6, 2, 10, 0, 0, 0, time.UTC),
			Author:    NewDeveloper("Alice <alice@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Carol <carol@example.com>")},
		},
	}

	emptyTeam, _ := NewTeam([]string{})
	matrix, _, _, _ := BuildPairMatrix(emptyTeam, commits, false)

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
	// Team file with Alice having multiple email addresses
	team, _ := NewTeam([]string{
		"Alice <alice@example.com>,<alice.work@company.com>",
		"Bob <bob@example.com>",
	})

	// Commits with Alice using different email addresses
	commits := []Commit{
		{
			Date:      time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
			Author:    NewDeveloper("Alice <alice@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Bob <bob@example.com>")},
		},
		{
			Date:      time.Date(2024, 6, 2, 10, 0, 0, 0, time.UTC),
			Author:    NewDeveloper("Alice <alice.work@company.com>"),
			CoAuthors: []Developer{NewDeveloper("Bob <bob@example.com>")},
		},
	}

	matrix, devs, _, _ := BuildPairMatrix(team, commits, true)

	// We should only have 2 developers (Alice and Bob), not 3
	if len(devs) != 2 {
		t.Errorf("expected 2 developers, got %d: %v", len(devs), devs)
	}

	// Find Alice's canonical email (should be the first one in the team entry)
	var aliceEmail string
	for _, dev := range devs {
		if strings.Contains(dev, "alice") {
			aliceEmail = dev
			break
		}
	}

	if aliceEmail == "" {
		t.Fatalf("could not find Alice in developers list: %v", devs)
	}

	// Check that both commits are counted as Alice pairing with Bob
	bobEmail := "bob@example.com"

	// Should have 2 pairs (one from each day)
	if matrix.Count(aliceEmail, bobEmail) != 2 {
		t.Errorf("expected Alice/Bob pair to have count 2, got %d", matrix.Count(aliceEmail, bobEmail))
	}
}

func TestTeamFileCanonicalName(t *testing.T) {
	// Team file with a canonical name
	team, _ := NewTeam([]string{
		"Canonical Alice <alice@example.com>",
		"Bob <bob@example.com>",
	})

	// Commits with a different name for Alice
	commits := []Commit{
		{
			Date:      time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
			Author:    NewDeveloper("Different Alice <alice@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Bob <bob@example.com>")},
		},
	}

	// Build the matrix with useTeam=true
	_, _, _, emailToName := BuildPairMatrix(team, commits, true)

	// Check that Alice's name is the canonical one from the team file
	aliceEmail := "alice@example.com"
	if name := emailToName[aliceEmail]; name != "Canonical Alice" {
		t.Errorf("expected name to be 'Canonical Alice', got '%s'", name)
	}

	// Also verify that Bob's name is preserved
	bobEmail := "bob@example.com"
	if name := emailToName[bobEmail]; name != "Bob" {
		t.Errorf("expected name to be 'Bob', got '%s'", name)
	}
}

func TestMultipleAuthorsInCommit(t *testing.T) {
	// A single commit with three authors (one main author and two co-authors)
	commits := []Commit{
		{
			Date:   time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
			Author: NewDeveloper("Alice <alice@example.com>"),
			CoAuthors: []Developer{
				NewDeveloper("Bob <bob@example.com>"),
				NewDeveloper("Carol <carol@example.com>"),
			},
		},
	}

	emptyTeam, _ := NewTeam([]string{})
	matrix, _, _, _ := BuildPairMatrix(emptyTeam, commits, false)

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
	team, err := NewTeam([]string{
		"Alice Smith <alice@example.com>,<alice.smith@company.com>,<asmith@personal.net>",
		"Bob Jones <bob@example.com>,<bjones@company.com>",
		"Carol Davis <carol@example.com>,<cdavis@company.com>",
		"Dave Wilson <dave@example.com>",
		"Eve Brown <eve@example.com>,<ebrown@company.com>",
		"Frank Thomas <frank@example.com>",
	})
	if err != nil {
		t.Fatalf("failed to create team: %v", err)
	}

	// Verify team setup
	if len(team.developers) != 6 {
		t.Errorf("expected 6 developers in team, got %d", len(team.developers))
	}

	// Create a comprehensive set of commits covering various scenarios
	now := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	day := 24 * time.Hour

	commits := []Commit{
		// Day 1: Alice pairs with Bob
		{
			Date:      now.Add(-14 * day),
			Author:    NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Bob Jones <bob@example.com>")},
		},
		// Day 2: Bob pairs with Carol
		{
			Date:      now.Add(-13 * day),
			Author:    NewDeveloper("Bob Jones <bjones@company.com>"), // Different email
			CoAuthors: []Developer{NewDeveloper("Carol Davis <cdavis@company.com>")},
		},
		// Day 3: Alice pairs with Carol and Dave (three-way pairing)
		{
			Date:   now.Add(-12 * day),
			Author: NewDeveloper("Alice Smith <alice.smith@company.com>"), // Different email
			CoAuthors: []Developer{
				NewDeveloper("Carol Davis <carol@example.com>"),
				NewDeveloper("Dave Wilson <dave@example.com>"),
			},
		},
		// Day 4: Eve pairs with Frank
		{
			Date:      now.Add(-11 * day),
			Author:    NewDeveloper("Eve Brown <ebrown@company.com>"),
			CoAuthors: []Developer{NewDeveloper("Frank Thomas <frank@example.com>")},
		},
		// Day 5: Alice pairs with Eve
		{
			Date:      now.Add(-10 * day),
			Author:    NewDeveloper("Alice Smith <asmith@personal.net>"), // Different email
			CoAuthors: []Developer{NewDeveloper("Eve Brown <eve@example.com>")},
		},
		// Day 6: Dave pairs with Frank
		{
			Date:      now.Add(-9 * day),
			Author:    NewDeveloper("Dave Wilson <dave@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Frank Thomas <frank@example.com>")},
		},
		// Day 7: Bob pairs with Dave
		{
			Date:      now.Add(-8 * day),
			Author:    NewDeveloper("Bob Jones <bob@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Dave Wilson <dave@example.com>")},
		},
		// Day 7: Also, Carol pairs with Eve (same day, different pair)
		{
			Date:      now.Add(-8 * day),
			Author:    NewDeveloper("Carol Davis <carol@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Eve Brown <ebrown@company.com>")},
		},
		// Day 8: Alice pairs with Frank
		{
			Date:      now.Add(-7 * day),
			Author:    NewDeveloper("Frank Thomas <frank@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Alice Smith <alice@example.com>")},
		},
		// Day 9-10: No commits

		// Day 11: Multiple commits for the same pair on the same day (should count once)
		{
			Date:      now.Add(-4 * day),
			Author:    NewDeveloper("Bob Jones <bob@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Carol Davis <carol@example.com>")},
		},
		{
			Date:      now.Add(-4 * day),
			Author:    NewDeveloper("Carol Davis <cdavis@company.com>"),
			CoAuthors: []Developer{NewDeveloper("Bob Jones <bjones@company.com>")},
		},
		// Day 12: External person not in the team (should be filtered out with useTeam=true)
		{
			Date:      now.Add(-3 * day),
			Author:    NewDeveloper("External Person <external@othercompany.com>"),
			CoAuthors: []Developer{NewDeveloper("Alice Smith <alice@example.com>")},
		},
		// Day 13: Alice pairs with Bob again
		{
			Date:      now.Add(-2 * day),
			Author:    NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []Developer{NewDeveloper("Bob Jones <bob@example.com>")},
		},
		// Day 14: All team members collaborate (large pairing)
		{
			Date:   now.Add(-1 * day),
			Author: NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []Developer{
				NewDeveloper("Bob Jones <bob@example.com>"),
				NewDeveloper("Carol Davis <carol@example.com>"),
				NewDeveloper("Dave Wilson <dave@example.com>"),
				NewDeveloper("Eve Brown <eve@example.com>"),
				NewDeveloper("Frank Thomas <frank@example.com>"),
			},
		},
	}

	// Test with team information
	matrix, devs, shortLabels, emailToName := BuildPairMatrix(team, commits, true)

	// Check number of developers
	if len(devs) != 6 {
		t.Errorf("expected 6 developers, got %d: %v", len(devs), devs)
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

	// Verify short labels are created for all developers
	if len(shortLabels) != 6 {
		t.Errorf("expected 6 short labels, got %d", len(shortLabels))
	}

	// Verify email to name mapping
	expectedNames := map[string]string{
		"alice@example.com": "Alice Smith",
		"bob@example.com":   "Bob Jones",
		"carol@example.com": "Carol Davis",
		"dave@example.com":  "Dave Wilson",
		"eve@example.com":   "Eve Brown",
		"frank@example.com": "Frank Thomas",
	}

	for email, expectedName := range expectedNames {
		if actualName := emailToName[email]; actualName != expectedName {
			t.Errorf("email %s: expected name %q, got %q", email, expectedName, actualName)
		}
	}

	// Now test without team information
	matrixNoTeam, devsNoTeam, shortLabelsNoTeam, emailToNameNoTeam := BuildPairMatrix(Team{}, commits, false)

	// We expect more developers here because without team info, we don't consolidate alternate emails
	expectedNonTeamDevsCount := 12 // All unique email addresses appear as separate developers
	if len(devsNoTeam) != expectedNonTeamDevsCount {
		t.Errorf("expected %d developers with no team filter, got %d: %v",
			expectedNonTeamDevsCount, len(devsNoTeam), devsNoTeam)
	}

	// Check that external email has a label
	var externalEmail string
	for _, dev := range devsNoTeam {
		if strings.Contains(dev, "external") {
			externalEmail = dev
			break
		}
	}

	if externalEmail != "" {
		if _, ok := shortLabelsNoTeam[externalEmail]; !ok {
			t.Errorf("external email %s should have a short label", externalEmail)
		}
		if name, ok := emailToNameNoTeam[externalEmail]; !ok || name != "External Person" {
			t.Errorf("external email %s should have name 'External Person', got %q", externalEmail, name)
		}
	} else {
		t.Error("external email not found in no-team developers list")
	}

	// Verify the Alice-External pair exists in the no-team matrix
	if externalEmail != "" {
		aliceEmail := "alice@example.com"
		if matrixNoTeam.Count(aliceEmail, externalEmail) != 1 {
			t.Errorf("expected Alice-External pair to have count 1 in no-team matrix")
		}
	}
}
