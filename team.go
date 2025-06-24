package main

import (
	"bufio"
	"os"
	"strings"
)

// Maps all emails to their canonical email (first one listed for a developer)
var emailToCanonical = map[string]string{}

type Team struct {
	team                []string
	developers          map[string]Developer
	emailToName         map[string]string // Maps emails to display names
	emailToPrimaryEmail map[string]string // Maps all emails to their canonical/primary email
}

type Developer struct {
	DisplayName     string
	EmailAddresses  []string
	AbbreviatedName string
}

func NewTeamFromFile(filename string) (Team, error) {
	team, err := readTeamFile(filename)
	if err != nil {
		return Team{}, err
	}

	return NewTeam(team)
}

func NewTeam(team []string) (Team, error) {
	buildEmailMapping(team)
	developers := make(map[string]Developer)
	emailToName := make(map[string]string)
	emailToPrimaryEmail := make(map[string]string)

	for _, member := range team {
		name := extractName(member)
		emails := extractAllEmails(member)
		if len(emails) == 0 {
			continue
		}

		canonicalEmail := emails[0]

		// Associate all emails with this name and primary email
		for _, email := range emails {
			emailToName[email] = name
			emailToPrimaryEmail[email] = canonicalEmail
		}

		developers[canonicalEmail] = Developer{
			DisplayName:     name,
			EmailAddresses:  emails,
			AbbreviatedName: shortName(name),
		}
	}

	return Team{
		team:                team,
		developers:          developers,
		emailToName:         emailToName,
		emailToPrimaryEmail: emailToPrimaryEmail,
	}, nil
}

func shortName(name string) string {
	// Initials of all the words in a string
	words := strings.Fields(name)
	if len(words) == 0 {
		return "NAN"
	}

	initials := make([]string, len(words))

	for i, word := range words {
		if len(word) > 0 {
			initials[i] = strings.ToUpper(string(word[0]))
		} else {
			initials[i] = "."
		}
	}

	return strings.Join(initials, "")
}

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
		email := emails[0]
		// If we have a canonical mapping for this email, use it
		if canonical, ok := emailToCanonical[email]; ok {
			return canonical
		}
		return email
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

// Build the mapping from all emails to their canonical emails
func buildEmailMapping(team []string) {
	emailToCanonical = make(map[string]string)

	for _, member := range team {
		emails := extractAllEmails(member)
		if len(emails) == 0 {
			continue
		}

		// The first email is the canonical one
		canonical := emails[0]

		// Map all emails for this developer to the canonical one
		for _, email := range emails {
			emailToCanonical[email] = canonical
		}
	}
}
