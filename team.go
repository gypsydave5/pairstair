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

// Extract just the email part from "Name <email>" or the first email from "Name <email1>,<email2>"
func extractEmail(author string) string {
	emails := extractAllEmails(author)
	if len(emails) > 0 {
		return emails[0]
	}
	return strings.ToLower(strings.TrimSpace(author))
}

// Extract all emails from "Name <email1>,<email2>,<email3>"
func extractAllEmails(author string) []string {
	var emails []string

	// Find all email parts between < and >
	parts := strings.Split(author, "<")
	for i := 1; i < len(parts); i++ {
		if idx := strings.Index(parts[i], ">"); idx >= 0 {
			email := strings.TrimSpace(parts[i][:idx])
			if email != "" {
				emails = append(emails, strings.ToLower(email))
			}
		}
	}

	if len(emails) == 0 {
		email := strings.ToLower(strings.TrimSpace(author))
		if email != "" {
			emails = append(emails, email)
		}
	}

	return emails
}

// Extract just the name part from "Name <email>"
func extractName(author string) string {
	name := author
	if idx := strings.Index(author, "<"); idx > 0 {
		name = strings.TrimSpace(author[:idx])
	}
	return name
}
