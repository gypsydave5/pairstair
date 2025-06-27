// Package git provides functionality for parsing git repositories and extracting
// commit information for pairing analysis.
//
// The package handles git log parsing, co-author detection, and time window
// validation for analyzing developer collaboration patterns.
package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Developer represents a developer extracted from git commits.
// This structure matches the main package's Developer type.
type Developer struct {
	DisplayName     string
	EmailAddresses  []string
	AbbreviatedName string
}

// CanonicalEmail returns the primary email address for the developer
func (d Developer) CanonicalEmail() string {
	if len(d.EmailAddresses) == 0 {
		return ""
	}
	return d.EmailAddresses[0]
}

// NewDeveloper creates a Developer from a "Name <email>" string
// This is the public constructor for Developer instances
func NewDeveloper(entry string) Developer {
	return newDeveloper(entry)
}

// Commit represents a git commit with author and co-author information
type Commit struct {
	Date      time.Time
	Author    Developer
	CoAuthors []Developer
}

// GetCommitsSince retrieves git commits from the current repository within the specified time window
func GetCommitsSince(window string) ([]Commit, error) {
	if err := ValidateWindow(window); err != nil {
		return nil, err
	}
	
	sinceArg := WindowToGitSince(window)
	cmd := exec.Command("git", "log", "--since="+sinceArg, "--pretty=format:%H%n%an <%ae>%n%ad%n%B%n==END==", "--date=iso")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	return ParseGitLogOutput(string(out)), nil
}

// ParseGitLogOutput parses the output from git log command and returns commits
// This function is exported to allow testing with mock data
func ParseGitLogOutput(output string) []Commit {
	scanner := bufio.NewScanner(bytes.NewReader([]byte(output)))
	var commits []Commit
	var c Commit
	var bodyLines []string
	lineNum := 0
	
	for scanner.Scan() {
		line := scanner.Text()
		if line == "==END==" {
			c.CoAuthors = ParseCoAuthors(strings.Join(bodyLines, "\n"))
			commits = append(commits, c)
			c = Commit{}
			bodyLines = nil
			lineNum = 0
			continue
		}
		
		switch lineNum {
		case 0:
			// commit hash, skip
		case 1:
			c.Author = newDeveloper(line)
		case 2:
			t, err := time.Parse("2006-01-02 15:04:05 -0700", line)
			if err != nil {
				// Try alternative format if first parse fails
				t, err = time.Parse("2006-01-02 15:04:05 -0700", line)
				if err != nil {
					break
				}
			}
			c.Date = t
		default:
			bodyLines = append(bodyLines, line)
		}
		lineNum++
	}
	
	return commits
}

// ParseCoAuthors extracts co-author information from a commit message body
func ParseCoAuthors(body string) []Developer {
	var coAuthors []Developer
	coAuthorRe := regexp.MustCompile(`Co-authored-by:\s*(.+?)\s*<(.+?)>`)
	
	for _, line := range strings.Split(body, "\n") {
		matches := coAuthorRe.FindStringSubmatch(line)
		if matches != nil && len(matches) >= 3 {
			authorString := fmt.Sprintf("%s <%s>", matches[1], matches[2])
			coAuthors = append(coAuthors, newDeveloper(authorString))
		}
	}
	
	return coAuthors
}

// WindowToGitSince converts a time window string (e.g., "2w", "1m") to git's --since format
func WindowToGitSince(window string) string {
	unitMap := map[byte]string{
		'd': "day",
		'w': "week",
		'm': "month",
		'y': "year",
	}
	
	if len(window) < 2 {
		return window
	}
	
	n := window[:len(window)-1]
	unit := window[len(window)-1]
	
	if u, ok := unitMap[unit]; ok {
		return fmt.Sprintf("%s.%ss", n, u)
	}
	
	return window
}

// ValidateWindow checks if a time window string is in valid format (e.g., "2w", "1m", "7d")
func ValidateWindow(window string) error {
	validWindow := regexp.MustCompile(`^\d+[dwmy]$`)
	if !validWindow.MatchString(window) {
		return fmt.Errorf("invalid window format: %s", window)
	}
	return nil
}

// newDeveloper creates a developer from a "Name <email>" string
// This is internal to the git package
func newDeveloper(entry string) Developer {
	name := extractName(entry)
	emails := ExtractAllEmails(entry)
	
	if len(emails) == 0 {
		return Developer{}
	}
	
	return Developer{
		DisplayName:     name,
		EmailAddresses:  emails,
		AbbreviatedName: shortName(name),
	}
}

// extractName extracts the name part from "Name <email>" format
func extractName(author string) string {
	if idx := strings.Index(author, "<"); idx >= 0 {
		return strings.TrimSpace(author[:idx])
	}
	return strings.TrimSpace(author)
}

// ExtractAllEmails extracts all email addresses from the author string
// This function is exported for use by other packages
func ExtractAllEmails(author string) []string {
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

// shortName creates an abbreviated name from a full name
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
