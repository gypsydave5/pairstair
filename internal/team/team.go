// Package team provides functionality for managing development teams and parsing
// team configuration files.
//
// The package handles team file parsing, sub-team organization, and developer
// identity management across multiple email addresses.
package team

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gypsydave5/pairstair/internal/git"
)

// Team represents a development team with member information and email mappings
type Team struct {
	team                []string
	developers          map[string]git.Developer
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

// GetDevelopers returns a slice of all developers in the team, sorted by canonical email
func (t Team) GetDevelopers() []git.Developer {
	var developers []git.Developer

	// Convert the map to a slice for consistent ordering
	var emails []string
	for email := range t.developers {
		emails = append(emails, email)
	}
	sort.Strings(emails)

	for _, email := range emails {
		developers = append(developers, t.developers[email])
	}

	return developers
}

// GetTeamMembers returns the original team member strings
func (t Team) GetTeamMembers() []string {
	return t.team
}

// NewTeamFromFile creates a Team from a team file, optionally filtering by sub-team
func NewTeamFromFile(filename string, subTeam string) (Team, error) {
	teamMembers, err := ReadTeamFile(filename, subTeam)
	if err != nil {
		return Team{}, err
	}

	return NewTeam(teamMembers)
}

// NewTeamFromDevelopers creates a Team from a slice of git.Developer objects
func NewTeamFromDevelopers(developers []git.Developer) Team {
	devMap := make(map[string]git.Developer)
	emailToName := make(map[string]string)
	emailToPrimaryEmail := make(map[string]string)
	var teamMembers []string

	for _, developer := range developers {
		if len(developer.EmailAddresses) == 0 {
			continue // Skip developers with no emails
		}

		// Associate all emails with this name and primary email
		for _, email := range developer.EmailAddresses {
			emailToName[email] = developer.DisplayName
			emailToPrimaryEmail[email] = developer.CanonicalEmail()
		}

		devMap[developer.CanonicalEmail()] = developer

		// Create a team member string for backward compatibility
		emailList := strings.Join(developer.EmailAddresses, ">,<")
		memberString := fmt.Sprintf("%s <%s>", developer.DisplayName, emailList)
		teamMembers = append(teamMembers, memberString)
	}

	return Team{
		team:                teamMembers,
		developers:          devMap,
		emailToName:         emailToName,
		emailToPrimaryEmail: emailToPrimaryEmail,
	}
}

// NewTeam creates a Team from a list of team member strings
func NewTeam(teamMembers []string) (Team, error) {
	developers := make(map[string]git.Developer)
	emailToName := make(map[string]string)
	emailToPrimaryEmail := make(map[string]string)

	for _, member := range teamMembers {
		developer := git.NewDeveloper(member)
		if len(developer.EmailAddresses) == 0 {
			continue // Skip invalid entries
		}

		// Associate all emails with this name and primary email
		for _, email := range developer.EmailAddresses {
			emailToName[email] = developer.DisplayName
			emailToPrimaryEmail[email] = developer.CanonicalEmail()
		}

		developers[developer.CanonicalEmail()] = developer
	}

	return Team{
		team:                teamMembers,
		developers:          developers,
		emailToName:         emailToName,
		emailToPrimaryEmail: emailToPrimaryEmail,
	}, nil
}

// ReadTeamFile reads and parses a team file, optionally filtering by sub-team
func ReadTeamFile(filename string, subTeam string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var teamMembers []string
	var currentSection string
	var inTargetSection bool

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if this is a section header [section_name]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			inTargetSection = (subTeam == "" || currentSection == subTeam)
			continue
		}

		// If no sub-team specified, include all lines not in sections
		// If sub-team specified, only include lines from that section
		if subTeam == "" {
			if currentSection == "" {
				teamMembers = append(teamMembers, line)
			}
		} else if inTargetSection {
			teamMembers = append(teamMembers, line)
		}
	}

	return teamMembers, scanner.Err()
}
