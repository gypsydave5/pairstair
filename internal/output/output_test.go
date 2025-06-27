package output_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gypsydave5/pairstair/internal/git"
	"github.com/gypsydave5/pairstair/internal/output"
	"github.com/gypsydave5/pairstair/internal/pairing"
)

func TestNewRenderer(t *testing.T) {
	tests := []struct {
		name         string
		outputFormat string
		expectedType string
	}{
		{
			name:         "CLI renderer for default format",
			outputFormat: "cli",
			expectedType: "*output.CLIRenderer",
		},
		{
			name:         "CLI renderer for empty format",
			outputFormat: "",
			expectedType: "*output.CLIRenderer",
		},
		{
			name:         "HTML renderer for html format",
			outputFormat: "html",
			expectedType: "*output.HTMLRenderer",
		},
		{
			name:         "CLI renderer for unknown format",
			outputFormat: "unknown",
			expectedType: "*output.CLIRenderer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := output.NewRenderer(tt.outputFormat)
			if renderer == nil {
				t.Fatal("Expected renderer, got nil")
			}

			// Check the type matches expected
			rendererType := getTypeName(renderer)
			if rendererType != tt.expectedType {
				t.Errorf("Expected renderer type %s, got %s", tt.expectedType, rendererType)
			}
		})
	}
}

func TestPrintMatrixCLI(t *testing.T) {
	// Create test data - just test with empty matrix
	matrix := pairing.NewMatrix()

	developers := []git.Developer{
		git.NewDeveloper("Alice Smith <alice@example.com>"),
		git.NewDeveloper("Bob Jones <bob@example.com>"),
		git.NewDeveloper("Charlie Brown <charlie@example.com>"),
	}

	// Test that PrintMatrixCLI function exists and can be called
	t.Run("PrintMatrixCLI exists and callable", func(t *testing.T) {
		// This should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PrintMatrixCLI panicked: %v", r)
			}
		}()
		output.PrintMatrixCLI(matrix, developers)
	})
}

func TestPrintRecommendationsCLI(t *testing.T) {
	tests := []struct {
		name            string
		recommendations []output.Recommendation
		strategy        string
	}{
		{
			name:            "empty recommendations",
			recommendations: []output.Recommendation{},
			strategy:        "least-paired",
		},
		{
			name: "single pair recommendation",
			recommendations: []output.Recommendation{
				{
					A:     "alice@example.com",
					B:     "bob@example.com",
					Count: 5,
				},
			},
			strategy: "least-paired",
		},
		{
			name: "unpaired developer",
			recommendations: []output.Recommendation{
				{
					A:     "alice@example.com",
					B:     "",
					Count: 0,
				},
			},
			strategy: "least-paired",
		},
		{
			name: "least-recent strategy",
			recommendations: []output.Recommendation{
				{
					A:         "alice@example.com",
					B:         "bob@example.com",
					Count:     5,
					DaysSince: 3,
					HasPaired: true,
				},
			},
			strategy: "least-recent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PrintRecommendationsCLI panicked: %v", r)
				}
			}()
			output.PrintRecommendationsCLI(tt.recommendations, tt.strategy)
		})
	}
}

func TestRecommendation(t *testing.T) {
	// Test the Recommendation struct
	rec := output.Recommendation{
		A:          "alice@example.com",
		B:          "bob@example.com",
		Count:      5,
		LastPaired: time.Now(),
		DaysSince:  3,
		HasPaired:  true,
	}

	if rec.A != "alice@example.com" {
		t.Errorf("Expected A to be 'alice@example.com', got %s", rec.A)
	}
	if rec.B != "bob@example.com" {
		t.Errorf("Expected B to be 'bob@example.com', got %s", rec.B)
	}
	if rec.Count != 5 {
		t.Errorf("Expected Count to be 5, got %d", rec.Count)
	}
	if rec.DaysSince != 3 {
		t.Errorf("Expected DaysSince to be 3, got %d", rec.DaysSince)
	}
	if !rec.HasPaired {
		t.Errorf("Expected HasPaired to be true, got %t", rec.HasPaired)
	}
}

func TestRenderHTMLToWriter(t *testing.T) {
	// Create test data using the existing test pattern
	alice := git.NewDeveloper("Alice Smith <alice@example.com>")
	bob := git.NewDeveloper("Bob Jones <bob@example.com>")
	developers := []git.Developer{alice, bob}

	// Create an empty matrix (we'll populate it via BuildPairMatrix in the future)
	matrix := pairing.NewMatrix()

	recommendations := []output.Recommendation{
		{
			A:         alice.AbbreviatedName,
			B:         bob.AbbreviatedName,
			Count:     2,
			HasPaired: true,
		},
	}

	// Test rendering to a string builder
	var result strings.Builder
	err := output.RenderHTMLToWriter(&result, matrix, developers, recommendations)
	if err != nil {
		t.Fatalf("RenderHTMLToWriter failed: %v", err)
	}

	htmlOutput := result.String()

	// Verify the HTML contains expected structural elements
	expectedContents := []string{
		"<!DOCTYPE html>",
		"<title>Pair Stair</title>",
		"<h1>Pair Stair Matrix</h1>",
		"<h2>Legend</h2>",
		alice.AbbreviatedName,
		alice.DisplayName,
		alice.CanonicalEmail(),
		bob.AbbreviatedName,
		bob.DisplayName,
		bob.CanonicalEmail(),
		"<h2>Pair Matrix</h2>",
		"<h2>Pairing Recommendations",
		"</html>",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(htmlOutput, expected) {
			t.Errorf("HTML output should contain %q, but got:\n%s", expected, htmlOutput)
		}
	}
}

func TestRenderHTMLToWriter_EmptyRecommendations(t *testing.T) {
	// Test with empty recommendations (too many developers case)
	alice := git.NewDeveloper("Alice Smith <alice@example.com>")
	developers := []git.Developer{alice}
	matrix := pairing.NewMatrix()
	recommendations := []output.Recommendation{}

	var result strings.Builder
	err := output.RenderHTMLToWriter(&result, matrix, developers, recommendations)
	if err != nil {
		t.Fatalf("RenderHTMLToWriter failed: %v", err)
	}

	htmlOutput := result.String()
	if !strings.Contains(htmlOutput, "too many developers") {
		t.Error("HTML output should mention too many developers when recommendations are empty")
	}
}

// Helper function to get type name for testing
func getTypeName(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
