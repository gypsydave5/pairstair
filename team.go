package main

import (
	"bufio"
	"os"
	"strings"
)

type Team struct {
	team                []string
	developers          map[string]Developer
	emailToName         map[string]string // Maps emails to display names
	emailToPrimaryEmail map[string]string // Maps all emails to their canonical/primary email
}

// HasDeveloperByEmail checks if the given email belongs to a developer on the team
func (t Team) HasDeveloperByEmail(email string) bool {
	_, ok := t.emailToPrimaryEmail[email]
	return ok
}

// GetEmailMappings returns the email-to-name and email-to-primary-email mappings
func (t Team) GetEmailMappings() (map[string]string, map[string]string) {
	return t.emailToName, t.emailToPrimaryEmail
}

func NewTeamFromFile(filename string) (Team, error) {
	team, err := readTeamFile(filename)
	if err != nil {
		return Team{}, err
	}

	return NewTeam(team)
}

func NewTeam(team []string) (Team, error) {
	developers := make(map[string]Developer)
	emailToName := make(map[string]string)
	emailToPrimaryEmail := make(map[string]string)

	for _, member := range team {
		developer := NewDeveloper(member)
		if len(developer.EmailAddresses) == 0 {
			continue
		}

		// Associate all emails with this name and primary email
		for _, email := range developer.EmailAddresses {
			emailToName[email] = developer.DisplayName
			emailToPrimaryEmail[email] = developer.CanonicalEmail()
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
