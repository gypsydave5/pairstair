package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// AcceptanceTest runs the actual pairstair binary and tests its behavior
func TestPairStairAcceptance(t *testing.T) {
	// Build the binary first
	binaryPath := buildPairStairBinary(t)
	defer os.Remove(binaryPath)

	tests := []struct {
		name         string
		setupRepo    func(t *testing.T, repoDir string)
		args         []string
		wantContains []string
		wantExitCode int
	}{
		{
			name: "version flag works",
			setupRepo: func(t *testing.T, repoDir string) {
				// No repo setup needed for version
			},
			args:         []string{"--version"},
			wantContains: []string{"0.6.0-dev"},
			wantExitCode: 0,
		},
		{
			name: "help flag works",
			setupRepo: func(t *testing.T, repoDir string) {
				// No repo setup needed for help
			},
			args:         []string{"--help"},
			wantContains: []string{"Usage of", "-window", "-strategy", "-team", "-output"},
			wantExitCode: 0,
		},
		{
			name: "basic pairing detection without team file",
			setupRepo: func(t *testing.T, repoDir string) {
				setupBasicPairingRepo(t, repoDir)
			},
			args: []string{"--window", "1y"},
			wantContains: []string{
				"alice@example.com",
				"bob@example.com",
				"Pairing Recommendations",
			},
			wantExitCode: 0,
		},
		{
			name: "pairing with team file",
			setupRepo: func(t *testing.T, repoDir string) {
				setupRepoWithTeamFile(t, repoDir)
			},
			args: []string{"--window", "1y"},
			wantContains: []string{
				"Alice Smith",
				"Bob Jones",
				"Legend:",
				"Pairing Recommendations",
			},
			wantExitCode: 0,
		},
		{
			name: "least-recent strategy",
			setupRepo: func(t *testing.T, repoDir string) {
				setupRepoWithTimestampedCommits(t, repoDir)
			},
			args: []string{"--strategy", "least-recent", "--window", "1y"},
			wantContains: []string{
				"least recent collaborations",
				"days ago",
			},
			wantExitCode: 0,
		},
		{
			name: "sub-team filtering",
			setupRepo: func(t *testing.T, repoDir string) {
				setupRepoWithSubTeams(t, repoDir)
			},
			args: []string{"--team", "frontend", "--window", "1y"},
			wantContains: []string{
				"carol@example.com",
				"dave@example.com",
			},
			wantExitCode: 0,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for this test
			testDir := t.TempDir()
			
			// Setup the repository if needed
			if tt.setupRepo != nil {
				tt.setupRepo(t, testDir)
			}

			// Run pairstair binary
			output, exitCode := runPairStair(t, binaryPath, testDir, tt.args)

			// Check exit code
			if exitCode != tt.wantExitCode {
				t.Errorf("expected exit code %d, got %d", tt.wantExitCode, exitCode)
			}

			// Check output contains expected strings
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output should contain %q, but got:\n%s", want, output)
				}
			}
		})
	}
}

// buildPairStairBinary builds the pairstair binary and returns its path
func buildPairStairBinary(t *testing.T) string {
	t.Helper()
	
	tempDir := t.TempDir()
	binaryPath := filepath.Join(tempDir, "pairstair")
	
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build pairstair binary: %v", err)
	}
	
	return binaryPath
}

// runPairStair runs the pairstair binary with given args in the specified directory
func runPairStair(t *testing.T, binaryPath, workDir string, args []string) (output string, exitCode int) {
	t.Helper()
	
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = workDir
	
	// Capture both stdout and stderr
	outputBytes, err := cmd.CombinedOutput()
	output = string(outputBytes)
	
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			t.Fatalf("failed to run pairstair: %v", err)
		}
	} else {
		exitCode = 0
	}
	
	return output, exitCode
}

// setupBasicPairingRepo creates a git repo with basic pairing commits
func setupBasicPairingRepo(t *testing.T, repoDir string) {
	t.Helper()
	
	// Initialize git repository
	runGitCommand(t, repoDir, "init")
	runGitCommand(t, repoDir, "config", "user.name", "Test User")
	runGitCommand(t, repoDir, "config", "user.email", "test@example.com")
	
	// Create initial commit
	writeFile(t, repoDir, "README.md", "# Test Project")
	runGitCommand(t, repoDir, "add", "README.md")
	runGitCommand(t, repoDir, "commit", "-m", "Initial commit")
	
	// Create commits with co-authored-by
	writeFile(t, repoDir, "feature1.txt", "Feature 1")
	runGitCommand(t, repoDir, "add", "feature1.txt")
	runGitCommand(t, repoDir, "commit", "-m", "Add feature 1\n\nCo-authored-by: Alice Smith <alice@example.com>\nCo-authored-by: Bob Jones <bob@example.com>")
	
	writeFile(t, repoDir, "feature2.txt", "Feature 2") 
	runGitCommand(t, repoDir, "add", "feature2.txt")
	runGitCommand(t, repoDir, "commit", "-m", "Add feature 2\n\nCo-authored-by: Alice Smith <alice@example.com>\nCo-authored-by: Carol Davis <carol@example.com>")
	
	writeFile(t, repoDir, "feature3.txt", "Feature 3")
	runGitCommand(t, repoDir, "add", "feature3.txt") 
	runGitCommand(t, repoDir, "commit", "-m", "Add feature 3\n\nCo-authored-by: Bob Jones <bob@example.com>\nCo-authored-by: Carol Davis <carol@example.com>")
}

// setupRepoWithTeamFile creates a repo with a .team file
func setupRepoWithTeamFile(t *testing.T, repoDir string) {
	t.Helper()
	
	// Setup basic repo first
	setupBasicPairingRepo(t, repoDir)
	
	// Add .team file
	teamContent := `Alice Smith <alice@example.com>
Bob Jones <bob@example.com>
Carol Davis <carol@example.com>
`
	writeFile(t, repoDir, ".team", teamContent)
}

// setupRepoWithTimestampedCommits creates commits at different times for recency testing
func setupRepoWithTimestampedCommits(t *testing.T, repoDir string) {
	t.Helper()
	
	// Initialize git repository
	runGitCommand(t, repoDir, "init")
	runGitCommand(t, repoDir, "config", "user.name", "Test User")
	runGitCommand(t, repoDir, "config", "user.email", "test@example.com")
	
	// Create initial commit
	writeFile(t, repoDir, "README.md", "# Test Project")
	runGitCommand(t, repoDir, "add", "README.md")
	runGitCommand(t, repoDir, "commit", "-m", "Initial commit")
	
	// Create commits with specific dates (using GIT_AUTHOR_DATE and GIT_COMMITTER_DATE)
	now := time.Now()
	
	// Recent commit (1 day ago) - Alice & Bob
	writeFile(t, repoDir, "recent.txt", "Recent work")
	runGitCommand(t, repoDir, "add", "recent.txt")
	runGitCommandWithDate(t, repoDir, now.AddDate(0, 0, -1), "commit", "-m", "Recent work\n\nCo-authored-by: Alice Smith <alice@example.com>\nCo-authored-by: Bob Jones <bob@example.com>")
	
	// Older commit (1 week ago) - Alice & Carol  
	writeFile(t, repoDir, "older.txt", "Older work")
	runGitCommand(t, repoDir, "add", "older.txt")
	runGitCommandWithDate(t, repoDir, now.AddDate(0, 0, -7), "commit", "-m", "Older work\n\nCo-authored-by: Alice Smith <alice@example.com>\nCo-authored-by: Carol Davis <carol@example.com>")
	
	// Very old commit (1 month ago) - Bob & Carol
	writeFile(t, repoDir, "oldest.txt", "Very old work")
	runGitCommand(t, repoDir, "add", "oldest.txt")
	runGitCommandWithDate(t, repoDir, now.AddDate(0, -1, 0), "commit", "-m", "Very old work\n\nCo-authored-by: Bob Jones <bob@example.com>\nCo-authored-by: Carol Davis <carol@example.com>")
}

// setupRepoWithSubTeams creates a repo with sub-teams in .team file
func setupRepoWithSubTeams(t *testing.T, repoDir string) {
	t.Helper()
	
	// Initialize git repository
	runGitCommand(t, repoDir, "init")
	runGitCommand(t, repoDir, "config", "user.name", "Test User")
	runGitCommand(t, repoDir, "config", "user.email", "test@example.com")
	
	// Create initial commit
	writeFile(t, repoDir, "README.md", "# Test Project")
	runGitCommand(t, repoDir, "add", "README.md")
	runGitCommand(t, repoDir, "commit", "-m", "Initial commit")
	
	// Create .team file with sub-teams
	teamContent := `Alice Lead <alice@example.com>
Bob Fullstack <bob@example.com>

[frontend]
Carol Frontend <carol@example.com>
Dave UI <dave@example.com>

[backend]
Eve Backend <eve@example.com>
Frank API <frank@example.com>
`
	writeFile(t, repoDir, ".team", teamContent)
	
	// Create commits involving different team members
	writeFile(t, repoDir, "frontend.txt", "Frontend work")
	runGitCommand(t, repoDir, "add", "frontend.txt")
	runGitCommand(t, repoDir, "commit", "-m", "Frontend work\n\nCo-authored-by: Carol Frontend <carol@example.com>\nCo-authored-by: Dave UI <dave@example.com>")
	
	writeFile(t, repoDir, "backend.txt", "Backend work")
	runGitCommand(t, repoDir, "add", "backend.txt")
	runGitCommand(t, repoDir, "commit", "-m", "Backend work\n\nCo-authored-by: Eve Backend <eve@example.com>\nCo-authored-by: Frank API <frank@example.com>")
}

// Helper functions for git operations and file writing

func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git command failed: git %v: %v", args, err)
	}
}

func runGitCommandWithDate(t *testing.T, dir string, date time.Time, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	
	// Set author and committer dates
	dateStr := date.Format(time.RFC3339)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE="+dateStr,
		"GIT_COMMITTER_DATE="+dateStr,
	)
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("git command with date failed: git %v: %v", args, err)
	}
}

func writeFile(t *testing.T, dir, filename, content string) {
	t.Helper()
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
}
