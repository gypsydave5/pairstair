package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	window := flag.String("window", "1w", "Time window to examine (e.g. 1d, 2w, 3m, 1y)")
	output := flag.String("output", "cli", "Output format: 'cli' (default) or 'html'")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
		os.Exit(1)
	}
	teamPath := filepath.Join(wd, ".team")
	team, err := NewTeamFromFile(teamPath)
	useTeam := true
	if err != nil {
		if os.IsNotExist(err) {
			useTeam = false
		} else {
			fmt.Fprintf(os.Stderr, "Error reading .team file: %v\n", err)
			os.Exit(1)
		}
	}

	commits, err := getGitCommitsSince(*window)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	matrix, devs, shortLabels, emailToName := BuildPairMatrix(commits, team, useTeam)

	if *output == "html" {
		err := RenderHTMLAndOpen(matrix, devs, shortLabels, emailToName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering HTML: %v\n", err)
			os.Exit(1)
		}
	} else {
		PrintMatrixCLI(matrix, devs, shortLabels, emailToName)
		PrintRecommendationsCLI(matrix, devs, shortLabels)
	}
}
