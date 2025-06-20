package main

import (
	"bufio"
	"os"
	"strings"
)

func readTeamFile(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var team []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			team = append(team, line)
		}
	}
	return team, scanner.Err()
}

// Extract just the email part from "Name <email>"
func extractEmail(author string) string {
	start := strings.Index(author, "<")
	end := strings.Index(author, ">")
	if start >= 0 && end > start {
		return strings.ToLower(strings.TrimSpace(author[start+1 : end]))
	}
	return strings.ToLower(strings.TrimSpace(author))
}

// Extract just the name part from "Name <email>"
func extractName(author string) string {
	name := author
	if idx := strings.Index(author, "<"); idx > 0 {
		name = strings.TrimSpace(author[:idx])
	}
	return name
}
