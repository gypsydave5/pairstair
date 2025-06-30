package pairing_test

import (
	"testing"
	"time"

	"github.com/gypsydave5/pairstair/internal/git"
	"github.com/gypsydave5/pairstair/internal/pairing"
	"github.com/gypsydave5/pairstair/internal/team"
)

func TestMatrix(t *testing.T) {
	matrix := pairing.NewMatrix()
	
	// Test initial empty matrix
	if matrix.Len() != 0 {
		t.Errorf("Expected empty matrix length 0, got %d", matrix.Len())
	}
	
	// Test Count for non-existent pair
	if count := matrix.Count("alice@example.com", "bob@example.com"); count != 0 {
		t.Errorf("Expected count 0 for non-existent pair, got %d", count)
	}
	
	// Test self-pair returns 0
	if count := matrix.Count("alice@example.com", "alice@example.com"); count != 0 {
		t.Errorf("Expected count 0 for self-pair, got %d", count)
	}
}

func TestRecencyMatrix(t *testing.T) {
	recency := pairing.NewRecencyMatrix()
	
	// Test LastPaired for non-existent pair
	_, exists := recency.LastPaired("alice@example.com", "bob@example.com")
	if exists {
		t.Error("Expected false for non-existent pair in recency matrix")
	}
	
	// Test self-pair returns false
	_, exists = recency.LastPaired("alice@example.com", "alice@example.com")
	if exists {
		t.Error("Expected false for self-pair in recency matrix")
	}
}

func TestBuildPairMatrixEmptyCommits(t *testing.T) {
	commits := []git.Commit{}
	
	matrix, recencyMatrix, developers := pairing.BuildPairMatrix(team.Empty, commits, false)
	
	if matrix.Len() != 0 {
		t.Errorf("Expected empty matrix for no commits, got length %d", matrix.Len())
	}
	
	if len(developers) != 0 {
		t.Errorf("Expected no developers for no commits, got %d", len(developers))
	}
	
	// Test recency matrix is also empty
	_, exists := recencyMatrix.LastPaired("anyone", "else")
	if exists {
		t.Error("Expected empty recency matrix for no commits")
	}
}

func TestBuildPairMatrixSingleAuthor(t *testing.T) {
	commits := []git.Commit{
		{
			Date:      time.Now(),
			Author:    git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{}, // No co-authors
		},
	}
	
	matrix, _, developers := pairing.BuildPairMatrix(team.Empty, commits, false)
	
	// Single author commits should not create pairs
	if matrix.Len() != 0 {
		t.Errorf("Expected no pairs for single author commits, got %d", matrix.Len())
	}
	
	// But should include the developer
	if len(developers) != 1 {
		t.Errorf("Expected 1 developer, got %d", len(developers))
	}
	
	if developers[0].CanonicalEmail() != "alice@example.com" {
		t.Errorf("Expected alice@example.com, got %s", developers[0].CanonicalEmail())
	}
	
	// Should have abbreviated name for Alice
	if developers[0].AbbreviatedName == "" {
		t.Error("Expected non-empty abbreviated name for Alice")
	}
	
	// Should have display name
	if developers[0].DisplayName != "Alice Smith" {
		t.Errorf("Expected 'Alice Smith', got %s", developers[0].DisplayName)
	}
}

func TestBuildPairMatrixBasicPairing(t *testing.T) {
	commits := []git.Commit{
		{
			Date:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			Author: git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Bob Jones <bob@example.com>"),
			},
		},
	}
	
	matrix, recencyMatrix, developers := pairing.BuildPairMatrix(team.Empty, commits, false)
	
	// Should have one pair
	if matrix.Len() != 1 {
		t.Errorf("Expected 1 pair, got %d", matrix.Len())
	}
	
	// Check pair count
	count := matrix.Count("alice@example.com", "bob@example.com")
	if count != 1 {
		t.Errorf("Expected pair count 1, got %d", count)
	}
	
	// Check recency
	lastPaired, exists := recencyMatrix.LastPaired("alice@example.com", "bob@example.com")
	if !exists {
		t.Error("Expected recency data for Alice-Bob pair")
	}
	
	expectedDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC) // Should be date only
	if !lastPaired.Equal(expectedDate) {
		t.Errorf("Expected last paired %v, got %v", expectedDate, lastPaired)
	}
	
	// Should have both developers
	if len(developers) != 2 {
		t.Errorf("Expected 2 developers, got %d", len(developers))
	}
	
	// Check developers are sorted by email
	expectedEmails := []string{"alice@example.com", "bob@example.com"}
	for i, expectedEmail := range expectedEmails {
		if developers[i].CanonicalEmail() != expectedEmail {
			t.Errorf("Expected developer %s at index %d, got %s", expectedEmail, i, developers[i].CanonicalEmail())
		}
	}
	
	// Check developer names
	expectedNames := map[string]string{
		"alice@example.com": "Alice Smith",
		"bob@example.com":   "Bob Jones",
	}
	
	for _, dev := range developers {
		if expectedName, ok := expectedNames[dev.CanonicalEmail()]; ok {
			if dev.DisplayName != expectedName {
				t.Errorf("Expected name %s for %s, got %s", expectedName, dev.CanonicalEmail(), dev.DisplayName)
			}
		} else {
			t.Errorf("Unexpected developer email: %s", dev.CanonicalEmail())
		}
	}
}

func TestBuildPairMatrixWithTeam(t *testing.T) {
	developers := []git.Developer{
		git.NewDeveloper("Alice Smith <alice@example.com>"),
		git.NewDeveloper("Bob Jones <bob@example.com>"),
	}
	teamObj := team.NewTeamFromDevelopers(developers)
	
	commits := []git.Commit{
		{
			Date:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			Author: git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Bob Jones <bob@example.com>"),
			},
		},
		{
			Date:   time.Date(2024, 6, 2, 12, 0, 0, 0, time.UTC),
			Author: git.NewDeveloper("External Person <external@other.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Alice Smith <alice@example.com>"),
			},
		},
	}
	
	matrix, _, developers := pairing.BuildPairMatrix(teamObj, commits, true)
	
	// Should only include team members
	if len(developers) != 2 {
		t.Errorf("Expected 2 team members, got %d: %v", len(developers), developers)
	}
	
	// Should have Alice and Bob
	expectedEmails := []string{"alice@example.com", "bob@example.com"}
	for i, expectedEmail := range expectedEmails {
		if developers[i].CanonicalEmail() != expectedEmail {
			t.Errorf("Expected developer %s at index %d, got %s", expectedEmail, i, developers[i].CanonicalEmail())
		}
	}
	
	// Should have one pair (Alice-Bob from first commit)
	if matrix.Len() != 1 {
		t.Errorf("Expected 1 pair, got %d", matrix.Len())
	}
	
	// External person should be filtered out
	count := matrix.Count("alice@example.com", "external@other.com")
	if count != 0 {
		t.Errorf("Expected no pairing with external person, got count %d", count)
	}
}

func TestBuildPairMatrixMultipleEmailsPerDeveloper(t *testing.T) {
	developers := []git.Developer{
		git.NewDeveloper("Alice Smith <alice@example.com>,<alice@company.com>"),
		git.NewDeveloper("Bob Jones <bob@example.com>"),
	}
	teamObj := team.NewTeamFromDevelopers(developers)
	
	commits := []git.Commit{
		{
			Date:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			Author: git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Bob Jones <bob@example.com>"),
			},
		},
		{
			Date:   time.Date(2024, 6, 2, 12, 0, 0, 0, time.UTC),
			Author: git.NewDeveloper("Alice Smith <alice@company.com>"), // Different email
			CoAuthors: []git.Developer{
				git.NewDeveloper("Bob Jones <bob@example.com>"),
			},
		},
	}
	
	matrix, _, developers := pairing.BuildPairMatrix(teamObj, commits, true)
	
	// Should consolidate Alice's emails to primary
	if len(developers) != 2 {
		t.Errorf("Expected 2 developers, got %d: %v", len(developers), developers)
	}
	
	// Should have both commits count toward the same pair
	count := matrix.Count("alice@example.com", "bob@example.com")
	if count != 2 {
		t.Errorf("Expected pair count 2 (emails consolidated), got %d", count)
	}
}

func TestBuildPairMatrixThreeWayPairing(t *testing.T) {
	commits := []git.Commit{
		{
			Date:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			Author: git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Bob Jones <bob@example.com>"),
				git.NewDeveloper("Carol Davis <carol@example.com>"),
			},
		},
	}
	
	matrix, _, developers := pairing.BuildPairMatrix(team.Empty, commits, false)
	
	// Three-way pairing should create 3 pairs: A-B, A-C, B-C
	if matrix.Len() != 3 {
		t.Errorf("Expected 3 pairs for three-way pairing, got %d", matrix.Len())
	}
	
	// Check all pairs exist
	expectedPairs := [][]string{
		{"alice@example.com", "bob@example.com"},
		{"alice@example.com", "carol@example.com"},
		{"bob@example.com", "carol@example.com"},
	}
	
	for _, pair := range expectedPairs {
		count := matrix.Count(pair[0], pair[1])
		if count != 1 {
			t.Errorf("Expected count 1 for pair %s-%s, got %d", pair[0], pair[1], count)
		}
	}
	
	// Should have 3 developers
	if len(developers) != 3 {
		t.Errorf("Expected 3 developers, got %d", len(developers))
	}
}

func TestBuildPairMatrixSamePairMultipleDays(t *testing.T) {
	commits := []git.Commit{
		{
			Date:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			Author: git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Bob Jones <bob@example.com>"),
			},
		},
		{
			Date:   time.Date(2024, 6, 1, 15, 30, 0, 0, time.UTC), // Same day, different time
			Author: git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Bob Jones <bob@example.com>"),
			},
		},
		{
			Date:   time.Date(2024, 6, 2, 10, 0, 0, 0, time.UTC), // Different day
			Author: git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Bob Jones <bob@example.com>"),
			},
		},
	}
	
	matrix, recencyMatrix, _ := pairing.BuildPairMatrix(team.Empty, commits, false)
	
	// Should count as 2 separate pairing days
	count := matrix.Count("alice@example.com", "bob@example.com")
	if count != 2 {
		t.Errorf("Expected count 2 for pairs on different days, got %d", count)
	}
	
	// Recency should be the most recent date (June 2nd)
	lastPaired, exists := recencyMatrix.LastPaired("alice@example.com", "bob@example.com")
	if !exists {
		t.Error("Expected recency data")
	}
	
	expectedDate := time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC)
	if !lastPaired.Equal(expectedDate) {
		t.Errorf("Expected most recent date %v, got %v", expectedDate, lastPaired)
	}
}

func TestBuildPairMatrixConsistentPairOrdering(t *testing.T) {
	commits := []git.Commit{
		{
			Date:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			Author: git.NewDeveloper("Bob Jones <bob@example.com>"), // Bob first
			CoAuthors: []git.Developer{
				git.NewDeveloper("Alice Smith <alice@example.com>"), // Alice second
			},
		},
	}
	
	matrix, _, _ := pairing.BuildPairMatrix(team.Empty, commits, false)
	
	// Should work regardless of order in commit
	count1 := matrix.Count("alice@example.com", "bob@example.com")
	count2 := matrix.Count("bob@example.com", "alice@example.com")
	
	if count1 != 1 || count2 != 1 {
		t.Errorf("Expected consistent count 1 regardless of order, got %d and %d", count1, count2)
	}
	
	// Both should return the same value (pair ordering is normalized internally)
	if count1 != count2 {
		t.Errorf("Expected same count regardless of parameter order, got %d vs %d", count1, count2)
	}
}
