package output_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gypsydave5/pairstair/internal/output"
	"github.com/gypsydave5/pairstair/internal/pairing"
)

func TestNewRenderer(t *testing.T) {
	tests := []struct {
		name           string
		outputFormat   string
		expectedType   string
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
	
	devs := []string{"alice@example.com", "bob@example.com", "charlie@example.com"}
	shortLabels := map[string]string{
		"alice@example.com":   "Alice",
		"bob@example.com":     "Bob",
		"charlie@example.com": "Char",
	}
	emailToName := map[string]string{
		"alice@example.com":   "Alice Smith",
		"bob@example.com":     "Bob Jones",
		"charlie@example.com": "Charlie Brown",
	}

	// Test that PrintMatrixCLI function exists and can be called
	t.Run("PrintMatrixCLI exists and callable", func(t *testing.T) {
		// This should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PrintMatrixCLI panicked: %v", r)
			}
		}()
		output.PrintMatrixCLI(matrix, devs, shortLabels, emailToName)
	})
}

func TestPrintRecommendationsCLI(t *testing.T) {
	shortLabels := map[string]string{
		"alice@example.com": "Alice",
		"bob@example.com":   "Bob",
	}

	tests := []struct {
		name            string
		recommendations []output.Recommendation
		strategy        string
	}{
		{
			name: "empty recommendations",
			recommendations: []output.Recommendation{},
			strategy: "least-paired",
		},
		{
			name: "single pair recommendation",
			recommendations: []output.Recommendation{
				{
					A: "alice@example.com", 
					B: "bob@example.com", 
					Count: 5,
				},
			},
			strategy: "least-paired",
		},
		{
			name: "unpaired developer",
			recommendations: []output.Recommendation{
				{
					A: "alice@example.com", 
					B: "", 
					Count: 0,
				},
			},
			strategy: "least-paired",
		},
		{
			name: "least-recent strategy",
			recommendations: []output.Recommendation{
				{
					A: "alice@example.com", 
					B: "bob@example.com", 
					Count: 5,
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
			output.PrintRecommendationsCLI(tt.recommendations, shortLabels, tt.strategy)
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

// Helper function to get type name for testing
func getTypeName(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
