package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds all command-line configuration
type Config struct {
	Window   string
	Output   string
	Strategy string
	Team     string
}

// parseFlags parses command-line flags and returns a Config
func parseFlags() *Config {
	config := &Config{}
	flag.StringVar(&config.Window, "window", "1w", "Time window to examine (e.g. 1d, 2w, 3m, 1y)")
	flag.StringVar(&config.Output, "output", "cli", "Output format: 'cli' (default) or 'html'")
	flag.StringVar(&config.Strategy, "strategy", "least-paired", "Recommendation strategy: 'least-paired' (default) or 'least-recent'")
	flag.StringVar(&config.Team, "team", "", "Sub-team to analyze (e.g. 'frontend', 'backend')")
	flag.Parse()
	return config
}

// exitOnError exits the program with an error message if err is not nil
func exitOnError(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
		os.Exit(1)
	}
}

func main() {
	config := parseFlags()

	wd, err := os.Getwd()
	exitOnError(err, "Error getting working directory")

	teamPath := filepath.Join(wd, ".team")
	team, err := NewTeamFromFile(teamPath, config.Team)
	useTeam := true
	if err != nil {
		if os.IsNotExist(err) {
			useTeam = false
		} else {
			exitOnError(err, "Error reading .team file")
		}
	}

	commits, err := getGitCommitsSince(config.Window)
	exitOnError(err, "Error getting git commits")

	matrix, pairRecency, devs, shortLabels, emailToName := BuildPairMatrix(team, commits, useTeam)

	renderer := NewRenderer(config.Output)
	err = renderer.Render(matrix, pairRecency, devs, shortLabels, emailToName, config.Strategy)
	exitOnError(err, "Error rendering output")
}
