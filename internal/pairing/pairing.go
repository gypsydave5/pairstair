// Package pairing provides functionality for analyzing developer pairing patterns
// from git commit history.
//
// The package handles pair matrix construction, recency tracking, and developer
// label generation for visualization and analysis of collaboration patterns.
package pairing

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gypsydave5/pairstair/internal/git"
	"github.com/gypsydave5/pairstair/internal/team"
)

// Pair represents a pair of developers identified by their email addresses
type Pair struct {
	A, B string
}

// Matrix tracks how many times each pair of developers has worked together
type Matrix struct {
	data map[Pair]int
}

// RecencyMatrix tracks when each pair of developers last worked together
type RecencyMatrix struct {
	data map[Pair]time.Time
}

// NewMatrix creates a new empty pairing matrix
func NewMatrix() *Matrix {
	return &Matrix{data: make(map[Pair]int)}
}

// NewRecencyMatrix creates a new empty recency matrix
func NewRecencyMatrix() *RecencyMatrix {
	return &RecencyMatrix{data: make(map[Pair]time.Time)}
}

// Count returns the count of times a pair has worked together
func (m *Matrix) Count(a, b string) int {
	if a == b {
		return 0 // Self-pairs always have count 0
	}

	// Ensure consistent ordering
	if a > b {
		a, b = b, a
	}

	return m.data[Pair{A: a, B: b}]
}

// LastPaired returns the last time a pair worked together
func (r *RecencyMatrix) LastPaired(a, b string) (time.Time, bool) {
	if a == b {
		return time.Time{}, false // Self-pairs don't have a last paired date
	}

	// Ensure consistent ordering
	if a > b {
		a, b = b, a
	}

	lastTime, exists := r.data[Pair{A: a, B: b}]
	return lastTime, exists
}

// Len returns the number of pairs in the matrix
func (m *Matrix) Len() int {
	return len(m.data)
}

// BuildPairMatrix constructs a pair matrix from the commits and team data
func BuildPairMatrix(team team.Team, commits []git.Commit, useTeam bool) (*Matrix, *RecencyMatrix, []git.Developer) {
	// Maps to track emails and names
	emailToName := make(map[string]string)
	emailToPrimaryEmail := make(map[string]string)

	// Process team file
	if useTeam {
		// Use the pre-calculated maps from the team
		emailToName, emailToPrimaryEmail = team.GetEmailMappings()
	}

	datePairs := make(map[string]map[Pair]struct{})
	devsSet := make(map[string]struct{})

	for _, c := range commits {
		var devsInCommit []git.Developer
		if useTeam {
			// When using team mode, include commits where any participant is a team member
			var teamMembers []git.Developer
			
			// Check if author is in team
			authorEmail := c.Author.CanonicalEmail()
			if team.HasDeveloperByEmail(authorEmail) {
				teamMembers = append(teamMembers, c.Author)
			}

			// Filter co-authors to only include team members
			for _, ca := range c.CoAuthors {
				coAuthorEmail := ca.CanonicalEmail()
				if team.HasDeveloperByEmail(coAuthorEmail) {
					teamMembers = append(teamMembers, ca)
				}
			}
			
			// Skip commits where no participants are team members
			if len(teamMembers) == 0 {
				continue
			}
			
			devsInCommit = teamMembers
		} else {
			devsInCommit = append([]git.Developer{c.Author}, c.CoAuthors...)

			for _, d := range devsInCommit {
				email := d.CanonicalEmail()
				if _, ok := emailToName[email]; !ok {
					emailToName[email] = d.DisplayName
				}
			}
		}

		// Create a set of unique developers (by primary email)
		emailMap := make(map[string]struct{})
		for _, d := range devsInCommit {
			email := d.CanonicalEmail()

			// If using team, map to primary email
			if useTeam {
				if primaryEmail, ok := emailToPrimaryEmail[email]; ok {
					emailMap[primaryEmail] = struct{}{}
				} else {
					// Not using team or email not in team
					emailMap[email] = struct{}{}
				}
			} else {
				// When not using team, each email is its own developer
				// We don't try to consolidate different emails for the same person
				emailMap[email] = struct{}{}
			}

			// Track all developers we've seen
			if useTeam {
				if primaryEmail, ok := emailToPrimaryEmail[email]; ok {
					devsSet[primaryEmail] = struct{}{}
				} else {
					devsSet[email] = struct{}{}
				}
			} else {
				devsSet[email] = struct{}{}
			}
		}

		uniqueDevs := make([]string, 0, len(emailMap))
		for e := range emailMap {
			uniqueDevs = append(uniqueDevs, e)
		}
		if len(uniqueDevs) < 2 {
			continue
		}

		// Create pairs for this date
		sort.Strings(uniqueDevs)
		date := c.Date.Format("2006-01-02")
		if _, ok := datePairs[date]; !ok {
			datePairs[date] = make(map[Pair]struct{})
		}
		for i := 0; i < len(uniqueDevs); i++ {
			for j := i + 1; j < len(uniqueDevs); j++ {
				p := Pair{A: uniqueDevs[i], B: uniqueDevs[j]}
				datePairs[date][p] = struct{}{}
			}
		}
	}

	// Build list of developers as Developer objects
	emailToDevs := make(map[string]git.Developer)
	
	// First, add developers from commits
	for email := range devsSet {
		// Try to find the developer info from the commits or team
		var dev git.Developer
		if useTeam {
			// For team mode, use team information
			if name, exists := emailToName[email]; exists {
				// Get all emails for this developer from team
				allEmails := []string{email}
				for teamEmail, primaryEmail := range emailToPrimaryEmail {
					if primaryEmail == email && teamEmail != email {
						allEmails = append(allEmails, teamEmail)
					}
				}
				dev = git.Developer{
					DisplayName:     name,
					EmailAddresses:  allEmails,
					AbbreviatedName: makeAbbreviatedName(name),
				}
			} else {
				// Fallback: create from email
				dev = git.NewDeveloper(email)
			}
		} else {
			// For non-team mode, use the display name we captured from commits
			if name, exists := emailToName[email]; exists && name != "" {
				dev = git.Developer{
					DisplayName:     name,
					EmailAddresses:  []string{email},
					AbbreviatedName: makeAbbreviatedName(name),
				}
			} else {
				// Fallback: create from email
				dev = git.NewDeveloper(email)
			}
		}
		emailToDevs[email] = dev
	}

	// Add any team members not found in commits
	if useTeam {
		for _, tm := range team.GetTeamMembers() {
			emails := git.ExtractAllEmails(tm)
			if len(emails) > 0 {
				primaryEmail := emails[0]
				if _, ok := devsSet[primaryEmail]; !ok {
					// Extract name from team member string
					name := extractNameFromTeamMember(tm)
					dev := git.Developer{
						DisplayName:     name,
						EmailAddresses:  emails,
						AbbreviatedName: makeAbbreviatedName(name),
					}
					emailToDevs[primaryEmail] = dev
				}
			}
		}
	}

	// Convert to sorted slice
	var devEmails []string
	for email := range emailToDevs {
		devEmails = append(devEmails, email)
	}
	sort.Strings(devEmails)

	devs := make([]git.Developer, len(devEmails))
	for i, email := range devEmails {
		devs[i] = emailToDevs[email]
	}

	// Build final matrix and recency matrix
	matrix := NewMatrix()
	recencyMatrix := NewRecencyMatrix()
	
	// Sort dates to process in chronological order
	var sortedDates []string
	for date := range datePairs {
		sortedDates = append(sortedDates, date)
	}
	sort.Strings(sortedDates)
	
	for _, date := range sortedDates {
		pairs := datePairs[date]
		seen := make(map[Pair]struct{})
		for p := range pairs {
			if _, ok := seen[p]; !ok {
				matrix.data[p]++
				// Parse the date and update recency
				if commitDate, err := time.Parse("2006-01-02", date); err == nil {
					recencyMatrix.data[p] = commitDate
				}
				seen[p] = struct{}{}
			}
		}
	}
	return matrix, recencyMatrix, devs
}

// makeAbbreviatedName creates initials from a full name, similar to the git package's shortName
func makeAbbreviatedName(name string) string {
	if name == "" {
		return "??"
	}
	
	words := strings.Fields(name)
	if len(words) == 0 {
		return "??"
	}

	initials := make([]string, len(words))
	for i, word := range words {
		if len(word) > 0 {
			initials[i] = strings.ToUpper(string(word[0]))
		} else {
			initials[i] = "?"
		}
	}

	return strings.Join(initials, "")
}

// extractNameFromTeamMember extracts the display name from a team member string like "Alice Smith <alice@example.com>"
func extractNameFromTeamMember(member string) string {
	if idx := strings.Index(member, "<"); idx >= 0 {
		return strings.TrimSpace(member[:idx])
	}
	return strings.TrimSpace(member)
}

// makeShortLabelsWithNames builds short labels for each developer, preferring name if available
func makeShortLabelsWithNames(devs []string, emailToName map[string]string) map[string]string {
	labels := make(map[string]string, len(devs))
	used := make(map[string]struct{})
	for _, d := range devs {
		name := emailToName[d]
		var label string
		if name != "" && name != d {
			label = initialsFromAuthor(name)
		} else {
			label = initialsFromEmail(d)
		}
		origLabel := label
		i := 2
		for {
			if _, exists := used[label]; !exists {
				break
			}
			label = fmt.Sprintf("%s%d", origLabel, i)
			i++
		}
		labels[d] = label
		used[label] = struct{}{}
	}
	return labels
}

// initialsFromEmail tries to extract initials from an email address
func initialsFromEmail(email string) string {
	at := -1
	for i, c := range email {
		if c == '@' {
			at = i
			break
		}
	}
	if at > 0 {
		parts := splitDot(email[:at])
		initials := ""
		for _, p := range parts {
			if len(p) > 0 {
				initials += string(p[0])
			}
		}
		if len(initials) == 0 && len(email) > 0 {
			initials = string(email[0])
		}
		if len(initials) > 4 {
			initials = initials[:4]
		}
		return initials
	}
	if len(email) > 0 {
		return string(email[0])
	}
	return "??"
}

// splitDot splits a string on '.' characters
func splitDot(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}

// initialsFromAuthor tries to extract initials from "Name <email>" format
func initialsFromAuthor(author string) string {
	name := author
	if idx := safeIndex(author, "<"); idx > 0 {
		name = trimSpace(author[:idx])
	}
	parts := fields(name)
	if len(parts) == 0 {
		return "??"
	}
	initials := ""
	for _, p := range parts {
		initials += string(p[0])
	}
	if len(initials) > 4 {
		initials = initials[:4]
	}
	return initials
}

// safeIndex is a safe version of strings.Index that returns -1 if not found
func safeIndex(s, substr string) int {
	idx := strings.Index(s, substr)
	if idx == -1 {
		return -1
	}
	return idx
}

// trimSpace trims whitespace from a string
func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

// fields splits a string into fields
func fields(s string) []string {
	return strings.Fields(s)
}
