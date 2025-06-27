package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gypsydave5/pairstair/internal/git"
)

type Pair struct {
	A, B string
}

type Matrix struct {
	data map[Pair]int
}

type RecencyMatrix struct {
	data map[Pair]time.Time
}

func NewMatrix() *Matrix {
	return &Matrix{data: make(map[Pair]int)}
}

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
func BuildPairMatrix(team Team, commits []Commit, useTeam bool) (*Matrix, *RecencyMatrix, []string, map[string]string, map[string]string) {
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
		var devsInCommit []Developer
		if useTeam {
			// When using team mode, include commits where any participant is a team member
			var teamMembers []Developer
			
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
			devsInCommit = append([]Developer{c.Author}, c.CoAuthors...)

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

	// Build list of developers
	devs := make([]string, 0, len(devsSet))
	for d := range devsSet {
		devs = append(devs, d)
	}

	// Add any team members not found in commits
	if useTeam {
		for _, tm := range team.team {
			emails := git.ExtractAllEmails(tm)
			if len(emails) > 0 {
				primaryEmail := emails[0]
				if _, ok := devsSet[primaryEmail]; !ok {
					devs = append(devs, primaryEmail)
				}
			}
		}
	}
	sort.Strings(devs)

	shortLabels := makeShortLabelsWithNames(devs, emailToName)

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
	return matrix, recencyMatrix, devs, shortLabels, emailToName
}

// Build short labels for each developer, prefer name if available
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

// Try to extract initials from an email address
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

// Try to extract initials from "Name <email>"
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

// Safe version of strings.Index that returns -1 if not found
func safeIndex(s, substr string) int {
	idx := strings.Index(s, substr)
	if idx == -1 {
		return -1
	}
	return idx
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

func fields(s string) []string {
	return strings.Fields(s)
}
