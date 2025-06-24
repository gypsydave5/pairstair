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
