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

func main() {
	config := parseFlags()

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	teamPath := filepath.Join(wd, ".team")
	team, err := NewTeamFromFile(teamPath, config.Team)
	useTeam := true
	if err != nil {
		if os.IsNotExist(err) {
			useTeam = false
		} else {
			fmt.Fprintf(os.Stderr, "Error reading .team file: %v\n", err)
			os.Exit(1)
		}
	}

	commits, err := getGitCommitsSince(config.Window)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	matrix, pairRecency, devs, shortLabels, emailToName := BuildPairMatrix(team, commits, useTeam)

	if config.Output == "html" {
		err := RenderHTMLAndOpen(matrix, devs, shortLabels, emailToName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering HTML: %v\n", err)
			os.Exit(1)
		}
	} else {
		PrintMatrixCLI(matrix, devs, shortLabels, emailToName)
		PrintRecommendationsCLI(matrix, pairRecency, devs, shortLabels, config.Strategy)
	}
}
