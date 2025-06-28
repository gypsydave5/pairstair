package recommend_test

import (
	"testing"

	"github.com/gypsydave5/pairstair/internal/git"
	"github.com/gypsydave5/pairstair/internal/pairing"
	"github.com/gypsydave5/pairstair/internal/recommend"
)

func TestGenerateRecommendations_LeastPaired(t *testing.T) {
	// Create test developers
	alice := git.NewDeveloper("Alice Smith <alice@example.com>")
	bob := git.NewDeveloper("Bob Jones <bob@example.com>")
	carol := git.NewDeveloper("Carol Davis <carol@example.com>")
	developers := []git.Developer{alice, bob, carol}

	// Create test matrix with known pairing counts
	matrix := pairing.NewMatrix()
	recencyMatrix := pairing.NewRecencyMatrix()

	// Generate recommendations
	recommendations := recommend.GenerateRecommendations(developers, matrix, recencyMatrix, recommend.LeastPaired)

	// Should get recommendations for 3 developers
	if len(recommendations) == 0 {
		t.Error("Expected recommendations, got none")
	}

	// Verify no empty recommendations
	for i, rec := range recommendations {
		if len(rec.A.EmailAddresses) == 0 {
			t.Errorf("Recommendation %d has empty A field", i)
		}
	}
}

func TestGenerateRecommendations_LeastRecent(t *testing.T) {
	// Create test developers
	alice := git.NewDeveloper("Alice Smith <alice@example.com>")
	bob := git.NewDeveloper("Bob Jones <bob@example.com>")
	developers := []git.Developer{alice, bob}

	// Create test matrices
	matrix := pairing.NewMatrix()
	recencyMatrix := pairing.NewRecencyMatrix()

	// Generate recommendations
	recommendations := recommend.GenerateRecommendations(developers, matrix, recencyMatrix, recommend.LeastRecent)

	// Should get recommendations for 2 developers
	if len(recommendations) != 1 {
		t.Errorf("Expected 1 recommendation, got %d", len(recommendations))
	}

	if len(recommendations) > 0 {
		rec := recommendations[0]
		if len(rec.A.EmailAddresses) == 0 || len(rec.B.EmailAddresses) == 0 {
			t.Error("Expected both A and B to be populated")
		}
	}
}

func TestGenerateRecommendations_EmptyDevelopers(t *testing.T) {
	matrix := pairing.NewMatrix()
	recencyMatrix := pairing.NewRecencyMatrix()
	
	recommendations := recommend.GenerateRecommendations([]git.Developer{}, matrix, recencyMatrix, recommend.LeastPaired)
	
	if recommendations != nil {
		t.Errorf("Expected nil for empty developers, got %v", recommendations)
	}
}

func TestGenerateRecommendations_SingleDeveloper(t *testing.T) {
	alice := git.NewDeveloper("Alice Smith <alice@example.com>")
	developers := []git.Developer{alice}
	
	matrix := pairing.NewMatrix()
	recencyMatrix := pairing.NewRecencyMatrix()
	
	recommendations := recommend.GenerateRecommendations(developers, matrix, recencyMatrix, recommend.LeastPaired)
	
	if recommendations != nil {
		t.Errorf("Expected nil for single developer, got %v", recommendations)
	}
}
