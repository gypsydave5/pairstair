package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/gypsydave5/pairstair/internal/git"
	"github.com/gypsydave5/pairstair/internal/update"
)

// Version is the fallback version, overridden by build info when available
const Version = "0.5.0-dev"

// Commit represents a git commit with author and co-author information
// This uses the main package Developer type for compatibility with existing code
type Commit struct {
	Date      time.Time
	Author    Developer
	CoAuthors []Developer
}

// Conversion functions to bridge git package and main package types
func convertCommit(gc git.Commit) Commit {
	return Commit{
		Date:      gc.Date,
		Author:    convertDeveloper(gc.Author),
		CoAuthors: convertDevelopers(gc.CoAuthors),
	}
}

func convertDeveloper(gd git.Developer) Developer {
	return Developer{
		DisplayName:     gd.DisplayName,
		EmailAddresses:  gd.EmailAddresses,
		AbbreviatedName: gd.AbbreviatedName,
	}
}

func convertDevelopers(gds []git.Developer) []Developer {
	result := make([]Developer, len(gds))
	for i, gd := range gds {
		result[i] = convertDeveloper(gd)
	}
	return result
}

// getGitCommitsFromPackage wraps git.GetCommitsSince with type conversion
func getGitCommitsFromPackage(window string) ([]Commit, error) {
	gitCommits, err := git.GetCommitsSince(window)
	if err != nil {
		return nil, err
	}

	commits := make([]Commit, len(gitCommits))
	for i, gc := range gitCommits {
		commits[i] = convertCommit(gc)
	}
	return commits, nil
}

// getVersion returns the version string, preferring build info over the constant
func getVersion() string {
	info, ok := debug.ReadBuildInfo()
	return getVersionFromBuildInfo(info, ok)
}

// getVersionFromBuildInfo extracts version information from build info
// This function is separated to make it testable
func getVersionFromBuildInfo(info *debug.BuildInfo, hasInfo bool) string {
	if hasInfo && info != nil {
		// Check for git tag in VCS settings
		var revision, tag string
		var modified bool

		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.tag":
				tag = setting.Value
			case "vcs.revision":
				revision = setting.Value
			case "vcs.modified":
				modified = setting.Value == "true"
			}
		}

		// If we have a clean tag, use it
		if tag != "" && !modified {
			return tag
		}

		// If we have a tag but modified, show tag + dirty
		if tag != "" && modified {
			return tag + "-dirty"
		}

		// If we have a commit hash, show version + short hash
		if revision != "" {
			short := revision
			if len(revision) > 8 {
				short = revision[:8]
			}
			if modified {
				return fmt.Sprintf("%s+%s-dirty", Version, short)
			}
			return fmt.Sprintf("%s+%s", Version, short)
		}

		// Check if this was built as a module
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}

	// Fallback to compile-time constant
	return Version
}

// Config holds all command-line configuration
type Config struct {
	Window   string
	Output   string
	Strategy string
	Team     string
	Version  bool
}

// parseFlags parses command-line flags and returns a Config
func parseFlags() *Config {
	config := &Config{}
	flag.StringVar(&config.Window, "window", "1w", "Time window to examine (e.g. 1d, 2w, 3m, 1y)")
	flag.StringVar(&config.Output, "output", "cli", "Output format: 'cli' (default) or 'html'")
	flag.StringVar(&config.Strategy, "strategy", "least-paired", "Recommendation strategy: 'least-paired' (default) or 'least-recent'")
	flag.StringVar(&config.Team, "team", "", "Sub-team to analyze (e.g. 'frontend', 'backend')")
	flag.BoolVar(&config.Version, "version", false, "Show version information")
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

	// Check for updates (silent failure, no caching)
	if updateMessage := update.CheckForUpdate(getVersion()); updateMessage != "" {
		fmt.Fprintln(os.Stderr, updateMessage)
		fmt.Fprintln(os.Stderr, "")
	}

	if config.Version {
		fmt.Println(getVersion())
		return
	}

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

	commits, err := getGitCommitsFromPackage(config.Window)
	exitOnError(err, "Error getting git commits")

	matrix, pairRecency, devs, shortLabels, emailToName := BuildPairMatrix(team, commits, useTeam)

	renderer := NewRenderer(config.Output)
	err = renderer.Render(matrix, pairRecency, devs, shortLabels, emailToName, config.Strategy)
	exitOnError(err, "Error rendering output")
}
