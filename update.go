package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GitHubRelease represents a GitHub release from the API
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Draft   bool   `json:"draft"`
}

// checkForUpdate checks for a newer version and returns an update message if available
func checkForUpdate(currentVersion string) string {
	return checkForUpdateWithURL(currentVersion, "https://api.github.com/repos/gypsydave5/pairstair/releases")
}

// checkForUpdateWithURL checks for updates using a custom URL (for testing)
func checkForUpdateWithURL(currentVersion, url string) string {
	client := &http.Client{Timeout: 3 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return "" // Silent failure
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "" // Silent failure
	}

	var releases []GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "" // Silent failure
	}

	// Find the latest non-draft release
	var latestVersion string
	for _, release := range releases {
		if !release.Draft {
			latestVersion = release.TagName
			break
		}
	}

	if latestVersion == "" {
		return "" // No releases found
	}

	if isNewerVersion(currentVersion, latestVersion) {
		return fmt.Sprintf("A newer version of pairstair is available: %s (you have %s)", latestVersion, currentVersion)
	}

	return ""
}

// isNewerVersion compares two version strings and returns true if latest is newer than current
func isNewerVersion(current, latest string) bool {
	currentClean := cleanVersion(current)
	latestClean := cleanVersion(latest)

	// Simple semantic version comparison
	currentParts := parseVersion(currentClean)
	latestParts := parseVersion(latestClean)

	// Compare major.minor.patch
	for i := 0; i < 3; i++ {
		if latestParts[i] > currentParts[i] {
			return true
		}
		if latestParts[i] < currentParts[i] {
			return false
		}
	}

	return false // Versions are equal
}

// cleanVersion removes prefixes and suffixes to get core version
func cleanVersion(version string) string {
	// Remove 'v' prefix
	version = strings.TrimPrefix(version, "v")

	// Remove anything after '+' (commit hash)
	if idx := strings.Index(version, "+"); idx != -1 {
		version = version[:idx]
	}

	// Remove anything after '-' (pre-release info like -dev, -dirty)
	if idx := strings.Index(version, "-"); idx != -1 {
		version = version[:idx]
	}

	return version
}

// parseVersion parses a version string like "1.2.3" into [1, 2, 3]
func parseVersion(version string) [3]int {
	var parts [3]int
	segments := strings.Split(version, ".")

	for i := 0; i < 3 && i < len(segments); i++ {
		if num, err := strconv.Atoi(segments[i]); err == nil {
			parts[i] = num
		}
	}

	return parts
}
