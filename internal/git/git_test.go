package git_test

import (
	"strings"
	"testing"

	"github.com/gypsydave5/pairstair/internal/git"
)

func TestParseCoAuthors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []git.Developer
	}{
		{
			name:  "single co-author",
			input: "Some commit message\n\nCo-authored-by: Alice Smith <alice@example.com>",
			expected: []git.Developer{
				git.NewDeveloper("Alice Smith <alice@example.com>"),
			},
		},
		{
			name:  "multiple co-authors",
			input: "Some commit message\n\nCo-authored-by: Alice Smith <alice@example.com>\nCo-authored-by: Bob Jones <bob@example.com>",
			expected: []git.Developer{
				git.NewDeveloper("Alice Smith <alice@example.com>"),
				git.NewDeveloper("Bob Jones <bob@example.com>"),
			},
		},
		{
			name:     "no co-authors",
			input:    "Some commit message with no co-authors",
			expected: []git.Developer{},
		},
		{
			name:  "co-authors with extra whitespace",
			input: "Some commit message\n\nCo-authored-by:  Alice Smith   <alice@example.com>  \nCo-authored-by:\tBob Jones\t<bob@example.com>",
			expected: []git.Developer{
				git.NewDeveloper("Alice Smith <alice@example.com>"),
				git.NewDeveloper("Bob Jones <bob@example.com>"),
			},
		},
		{
			name:  "mixed content with co-authors",
			input: "Fix bug in parser\n\nThis fixes the issue where the parser would fail.\n\nCo-authored-by: Alice Smith <alice@example.com>\nSome other text\nCo-authored-by: Bob Jones <bob@example.com>",
			expected: []git.Developer{
				git.NewDeveloper("Alice Smith <alice@example.com>"),
				git.NewDeveloper("Bob Jones <bob@example.com>"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should call the public API from the git package
			result := git.ParseCoAuthors(tt.input)
			
			if len(result) != len(tt.expected) {
				t.Errorf("ParseCoAuthors() returned %d co-authors, expected %d", len(result), len(tt.expected))
				return
			}
			
			for i, expected := range tt.expected {
				if i >= len(result) {
					t.Errorf("Missing co-author at index %d", i)
					continue
				}
				
				// Compare the canonical representation
				if result[i].CanonicalEmail() != expected.EmailAddresses[0] {
					t.Errorf("Co-author %d: got email %q, expected %q", i, result[i].CanonicalEmail(), expected.EmailAddresses[0])
				}
				
				if result[i].DisplayName != expected.DisplayName {
					t.Errorf("Co-author %d: got name %q, expected %q", i, result[i].DisplayName, expected.DisplayName)
				}
			}
		})
	}
}

func TestWindowToGitSince(t *testing.T) {
	tests := []struct {
		name     string
		window   string
		expected string
	}{
		{
			name:     "days",
			window:   "7d",
			expected: "7.days",
		},
		{
			name:     "weeks",
			window:   "2w",
			expected: "2.weeks",
		},
		{
			name:     "months",
			window:   "3m",
			expected: "3.months",
		},
		{
			name:     "years",
			window:   "1y",
			expected: "1.years",
		},
		{
			name:     "single digit",
			window:   "1d",
			expected: "1.days",
		},
		{
			name:     "multi digit",
			window:   "30d",
			expected: "30.days",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := git.WindowToGitSince(tt.window)
			if result != tt.expected {
				t.Errorf("WindowToGitSince(%q) = %q, expected %q", tt.window, result, tt.expected)
			}
		})
	}
}

func TestValidateWindow(t *testing.T) {
	tests := []struct {
		name    string
		window  string
		wantErr bool
	}{
		{
			name:    "valid days",
			window:  "7d",
			wantErr: false,
		},
		{
			name:    "valid weeks",
			window:  "2w",
			wantErr: false,
		},
		{
			name:    "valid months",
			window:  "3m",
			wantErr: false,
		},
		{
			name:    "valid years",
			window:  "1y",
			wantErr: false,
		},
		{
			name:    "invalid format - no number",
			window:  "d",
			wantErr: true,
		},
		{
			name:    "invalid format - no unit",
			window:  "7",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong unit",
			window:  "7x",
			wantErr: true,
		},
		{
			name:    "invalid format - multiple units",
			window:  "7dm",
			wantErr: true,
		},
		{
			name:    "empty string",
			window:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := git.ValidateWindow(tt.window)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWindow(%q) error = %v, wantErr %v", tt.window, err, tt.wantErr)
			}
		})
	}
}

func TestGetCommitsSince_Integration(t *testing.T) {
	// This is more of an integration test - we'll test with actual git commands
	// but we need to make it work in the test environment
	
	tests := []struct {
		name      string
		window    string
		expectErr bool
	}{
		{
			name:      "valid window format",
			window:    "1w",
			expectErr: false, // Should not error on window format validation
		},
		{
			name:      "invalid window format",
			window:    "invalid",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the public API
			_, err := git.GetCommitsSince(tt.window)
			
			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.expectErr && err != nil {
				// For valid window formats, we might still get git errors if not in a repo
				// but we should not get window validation errors
				if strings.Contains(err.Error(), "invalid window format") {
					t.Errorf("Got window validation error for valid window: %v", err)
				}
				// Git command errors are acceptable in test environment
			}
		})
	}
}

func TestGetCommitsInPath_WithMockData(t *testing.T) {
	// Test the testable function that accepts a git command runner
	// This will allow us to test the git parsing logic without actual git commands
	
	// Mock git log output
	mockGitOutput := `abc123
Alice Smith <alice@example.com>
2024-01-15 10:30:00 -0800
Add new feature

Co-authored-by: Bob Jones <bob@example.com>
==END==
def456
Carol White <carol@example.com>
2024-01-14 14:22:00 -0800
Fix bug in parser

==END==`

	// Test that we can parse the mock output correctly
	// This tests the parsing logic separately from git command execution
	result := git.ParseGitLogOutput(mockGitOutput)
	
	if len(result) != 2 {
		t.Fatalf("Expected 2 commits, got %d", len(result))
	}
	
	// Test first commit
	commit1 := result[0]
	if commit1.Author.DisplayName != "Alice Smith" {
		t.Errorf("First commit author: got %q, expected %q", commit1.Author.DisplayName, "Alice Smith")
	}
	
	if len(commit1.CoAuthors) != 1 {
		t.Errorf("First commit co-authors: got %d, expected 1", len(commit1.CoAuthors))
	} else if commit1.CoAuthors[0].DisplayName != "Bob Jones" {
		t.Errorf("First commit co-author: got %q, expected %q", commit1.CoAuthors[0].DisplayName, "Bob Jones")
	}
	
	// Test second commit
	commit2 := result[1]
	if commit2.Author.DisplayName != "Carol White" {
		t.Errorf("Second commit author: got %q, expected %q", commit2.Author.DisplayName, "Carol White")
	}
	
	if len(commit2.CoAuthors) != 0 {
		t.Errorf("Second commit co-authors: got %d, expected 0", len(commit2.CoAuthors))
	}
}
