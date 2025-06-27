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
	emptyTeam, _ := team.NewTeam([]string{})
	commits := []git.Commit{}
	
	matrix, recencyMatrix, devs, shortLabels, emailToName := pairing.BuildPairMatrix(emptyTeam, commits, false)
	
	if matrix.Len() != 0 {
		t.Errorf("Expected empty matrix for no commits, got length %d", matrix.Len())
	}
	
	if len(devs) != 0 {
		t.Errorf("Expected no developers for no commits, got %d", len(devs))
	}
	
	if len(shortLabels) != 0 {
		t.Errorf("Expected no labels for no commits, got %d", len(shortLabels))
	}
	
	if len(emailToName) != 0 {
		t.Errorf("Expected no email mappings for no commits, got %d", len(emailToName))
	}
	
	// Test recency matrix is also empty
	_, exists := recencyMatrix.LastPaired("anyone", "else")
	if exists {
		t.Error("Expected empty recency matrix for no commits")
	}
}

func TestBuildPairMatrixSingleAuthor(t *testing.T) {
	emptyTeam, _ := team.NewTeam([]string{})
	commits := []git.Commit{
		{
			Date:      time.Now(),
			Author:    git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{}, // No co-authors
		},
	}
	
	matrix, _, devs, shortLabels, emailToName := pairing.BuildPairMatrix(emptyTeam, commits, false)
	
	// Single author commits should not create pairs
	if matrix.Len() != 0 {
		t.Errorf("Expected no pairs for single author commits, got %d", matrix.Len())
	}
	
	// But should include the developer
	if len(devs) != 1 {
		t.Errorf("Expected 1 developer, got %d", len(devs))
	}
	
	if devs[0] != "alice@example.com" {
		t.Errorf("Expected alice@example.com, got %s", devs[0])
	}
	
	// Should have short label for Alice
	if label, ok := shortLabels["alice@example.com"]; !ok {
		t.Error("Expected short label for Alice")
	} else if label == "" {
		t.Error("Expected non-empty short label for Alice")
	}
	
	// Should have email to name mapping
	if name, ok := emailToName["alice@example.com"]; !ok {
		t.Error("Expected email to name mapping for Alice")
	} else if name != "Alice Smith" {
		t.Errorf("Expected 'Alice Smith', got %s", name)
	}
}

func TestBuildPairMatrixBasicPairing(t *testing.T) {
	emptyTeam, _ := team.NewTeam([]string{})
	commits := []git.Commit{
		{
			Date:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			Author: git.NewDeveloper("Alice Smith <alice@example.com>"),
			CoAuthors: []git.Developer{
				git.NewDeveloper("Bob Jones <bob@example.com>"),
			},
		},
	}
	
	matrix, recencyMatrix, devs, shortLabels, emailToName := pairing.BuildPairMatrix(emptyTeam, commits, false)
	
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
	if len(devs) != 2 {
		t.Errorf("Expected 2 developers, got %d", len(devs))
	}
	
	// Check developers are sorted
	expectedDevs := []string{"alice@example.com", "bob@example.com"}
	for i, expectedDev := range expectedDevs {
		if devs[i] != expectedDev {
			t.Errorf("Expected developer %s at index %d, got %s", expectedDev, i, devs[i])
		}
	}
	
	// Check short labels
	if len(shortLabels) != 2 {
		t.Errorf("Expected 2 short labels, got %d", len(shortLabels))
	}
	
	// Check email mappings
	expectedMappings := map[string]string{
		"alice@example.com": "Alice Smith",
		"bob@example.com":   "Bob Jones",
	}
	
	for email, expectedName := range expectedMappings {
		if name, ok := emailToName[email]; !ok {
			t.Errorf("Missing email mapping for %s", email)
		} else if name != expectedName {
			t.Errorf("Expected name %s for %s, got %s", expectedName, email, name)
		}
	}
}

func TestBuildPairMatrixWithTeam(t *testing.T) {
	teamObj, err := team.NewTeam([]string{
		"Alice Smith <alice@example.com>",
		"Bob Jones <bob@example.com>",
	})
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}
	
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
	
	matrix, _, devs, _, _ := pairing.BuildPairMatrix(teamObj, commits, true)
	
	// Should only include team members
	if len(devs) != 2 {
		t.Errorf("Expected 2 team members, got %d: %v", len(devs), devs)
	}
	
	// Should have Alice and Bob
	expectedDevs := []string{"alice@example.com", "bob@example.com"}
	for i, expectedDev := range expectedDevs {
		if devs[i] != expectedDev {
			t.Errorf("Expected developer %s at index %d, got %s", expectedDev, i, devs[i])
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
	teamObj, err := team.NewTeam([]string{
		"Alice Smith <alice@example.com>,<alice@company.com>",
		"Bob Jones <bob@example.com>",
	})
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}
	
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
	
	matrix, _, devs, _, _ := pairing.BuildPairMatrix(teamObj, commits, true)
	
	// Should consolidate Alice's emails to primary
	if len(devs) != 2 {
		t.Errorf("Expected 2 developers, got %d: %v", len(devs), devs)
	}
	
	// Should have both commits count toward the same pair
	count := matrix.Count("alice@example.com", "bob@example.com")
	if count != 2 {
		t.Errorf("Expected pair count 2 (emails consolidated), got %d", count)
	}
}

func TestBuildPairMatrixThreeWayPairing(t *testing.T) {
	emptyTeam, _ := team.NewTeam([]string{})
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
	
	matrix, _, devs, _, _ := pairing.BuildPairMatrix(emptyTeam, commits, false)
	
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
	if len(devs) != 3 {
		t.Errorf("Expected 3 developers, got %d", len(devs))
	}
}

func TestBuildPairMatrixSamePairMultipleDays(t *testing.T) {
	emptyTeam, _ := team.NewTeam([]string{})
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
	
	matrix, recencyMatrix, _, _, _ := pairing.BuildPairMatrix(emptyTeam, commits, false)
	
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
	emptyTeam, _ := team.NewTeam([]string{})
	commits := []git.Commit{
		{
			Date:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			Author: git.NewDeveloper("Bob Jones <bob@example.com>"), // Bob first
			CoAuthors: []git.Developer{
				git.NewDeveloper("Alice Smith <alice@example.com>"), // Alice second
			},
		},
	}
	
	matrix, _, _, _, _ := pairing.BuildPairMatrix(emptyTeam, commits, false)
	
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
