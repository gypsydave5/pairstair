// Package pairing provides functionality for analyzing developer pairing patterns
// from git commit history.
//
// The package handles pair matrix construction, recency tracking, and developer
// label generation for visualization and analysis of collaboration patterns.
package pairing

import (
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

// CountByDeveloper returns the count of times a pair of developers has worked together
func (m *Matrix) CountByDeveloper(a, b git.Developer) int {
	return m.Count(a.CanonicalEmail(), b.CanonicalEmail())
}

// Add increments the count for a pair of developers
func (m *Matrix) Add(a, b string) {
	if a == b {
		return // Skip self-pairs
	}

	// Ensure consistent ordering
	if a > b {
		a, b = b, a
	}

	m.data[Pair{A: a, B: b}]++
}

// AddByDeveloper increments the count for a pair of developers
func (m *Matrix) AddByDeveloper(a, b git.Developer) {
	m.Add(a.CanonicalEmail(), b.CanonicalEmail())
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

// LastPairedByDeveloper returns the last time a pair of developers worked together
func (r *RecencyMatrix) LastPairedByDeveloper(a, b git.Developer) (time.Time, bool) {
	return r.LastPaired(a.CanonicalEmail(), b.CanonicalEmail())
}

// Record records the pairing time for two developers
func (r *RecencyMatrix) Record(a, b string, date time.Time) {
	if a == b {
		return // Skip self-pairs
	}

	// Ensure consistent ordering
	if a > b {
		a, b = b, a
	}

	pair := Pair{A: a, B: b}
	if existing, exists := r.data[pair]; !exists || date.After(existing) {
		r.data[pair] = date
	}
}

// RecordByDeveloper records the pairing time for two developers
func (r *RecencyMatrix) RecordByDeveloper(a, b git.Developer, date time.Time) {
	r.Record(a.CanonicalEmail(), b.CanonicalEmail(), date)
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
		for _, dev := range team.GetDevelopers() {
			primaryEmail := dev.CanonicalEmail()
			if _, ok := devsSet[primaryEmail]; !ok {
				emailToDevs[primaryEmail] = dev
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

