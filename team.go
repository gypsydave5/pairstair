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

// HasDeveloper checks if the given email belongs to a developer on the team
func (t Team) HasDeveloper(email string) bool {
	_, ok := t.emailToPrimaryEmail[email]
	return ok
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
	emailToDeveloper := make(map[string]Developer)

	for _, member := range team {
		developer := NewDeveloper(member)
		if len(developer.EmailAddresses) == 0 {
			continue
		}

		// Associate all emails with this name and primary email
		for _, email := range developer.EmailAddresses {
			emailToName[email] = developer.DisplayName
			emailToPrimaryEmail[email] = developer.CanonicalEmail()
			emailToDeveloper[email] = developer
		}

		developers[developer.CanonicalEmail()] = developer
	}

	return Team{
		team:                team,
		developers:          developers,
		emailToName:         emailToName,
		emailToPrimaryEmail: emailToPrimaryEmail,
	}, nil
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
		return email
	}
	return strings.ToLower(strings.TrimSpace(author))
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
