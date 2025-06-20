package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Pair struct {
	A, B string
}

type Recommendation struct {
	A, B  string
	Count int
}

func main() {
	window := flag.String("window", "1w", "Time window to examine (e.g. 1d, 2w, 3m, 1y)")
	flag.Parse()

	// Find .team in the current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
		os.Exit(1)
	}
	teamPath := filepath.Join(wd, ".team")
	team, err := readTeamFile(teamPath)
	useTeam := true
	if err != nil {
		if os.IsNotExist(err) {
			useTeam = false
		} else {
			fmt.Fprintf(os.Stderr, "Error reading .team file: %v\n", err)
			os.Exit(1)
		}
	}
	teamSet := make(map[string]struct{}, len(team))
	for _, t := range team {
		teamSet[t] = struct{}{}
	}

	commits, err := getGitCommitsSince(*window)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Build a mapping from full author string to canonical email (just the email part)
	authorToEmail := make(map[string]string)
	emailSet := make(map[string]struct{})

	// Build a mapping from canonical email to name (prefer name for abbreviation)
	emailToName := make(map[string]string)

	// Collect all emails from .team if present, else from commits
	if useTeam {
		for _, t := range team {
			email := extractEmail(t)
			name := extractName(t)
			authorToEmail[t] = email
			emailSet[email] = struct{}{}
			emailToName[email] = name
		}
	}

	datePairs := make(map[string]map[Pair]struct{})
	devsSet := make(map[string]struct{})

	for _, c := range commits {
		var devs []string
		if useTeam {
			// Only consider commits where the author is in the team
			if _, ok := teamSet[c.Author]; !ok {
				continue
			}
			// Only consider coauthors that are in the team
			filteredCoAuthors := make([]string, 0, len(c.CoAuthors))
			for _, ca := range c.CoAuthors {
				if _, ok := teamSet[ca]; ok {
					filteredCoAuthors = append(filteredCoAuthors, ca)
				}
			}
			devs = append([]string{c.Author}, filteredCoAuthors...)
		} else {
			// Use all authors and coauthors
			devs = append([]string{c.Author}, c.CoAuthors...)
		}
		// Map all devs to their canonical email
		emailMap := make(map[string]struct{})
		for _, d := range devs {
			email := extractEmail(d)
			name := extractName(d)
			authorToEmail[d] = email
			// Prefer the first name seen for this email
			if _, ok := emailToName[email]; !ok {
				emailToName[email] = name
			}
			emailMap[email] = struct{}{}
			devsSet[email] = struct{}{}
		}
		uniqueDevs := make([]string, 0, len(emailMap))
		for e := range emailMap {
			uniqueDevs = append(uniqueDevs, e)
		}
		if len(uniqueDevs) < 2 {
			continue // skip if not a pair
		}
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

	// Only use team members for the matrix if .team exists, else use all found devs
	devs := make([]string, 0, len(devsSet))
	for d := range devsSet {
		devs = append(devs, d)
	}
	if useTeam {
		// Ensure all team members are present, even if they didn't pair
		for _, t := range team {
			email := extractEmail(t)
			if _, ok := devsSet[email]; !ok {
				devs = append(devs, email)
			}
		}
	}
	sort.Strings(devs)

	// Build short labels for each developer, prefer name if available
	shortLabels := makeShortLabelsWithNames(devs, emailToName)

	matrix := make(map[Pair]int)
	for _, pairs := range datePairs {
		seen := make(map[Pair]struct{})
		for p := range pairs {
			if _, ok := seen[p]; !ok {
				matrix[p]++
				seen[p] = struct{}{}
			}
		}
	}

	// Print legend with name and email
	fmt.Println("Legend:")
	for _, d := range devs {
		name := emailToName[d]
		if name == "" {
			name = d
		}
		fmt.Printf("  %-6s = %-20s %s\n", shortLabels[d], name, d)
	}
	fmt.Println()

	// Print header with short labels
	fmt.Printf("%-8s", "")
	for _, d := range devs {
		fmt.Printf("%-8s", shortLabels[d])
	}
	fmt.Println()
	// Print rows
	for _, d1 := range devs {
		fmt.Printf("%-8s", shortLabels[d1])
		for _, d2 := range devs {
			if d1 == d2 {
				fmt.Printf("%-8s", "-")
				continue
			}
			a, b := d1, d2
			if a > b {
				a, b = b, a
			}
			fmt.Printf("%-8d", matrix[Pair{A: a, B: b}])
		}
		fmt.Println()
	}

	// --- Recommendations Section ---
	fmt.Println()
	fmt.Println("Pairing Recommendations (least-paired first):")
	recommendations := recommendPairsUnique(devs, matrix)
	for _, rec := range recommendations {
		labelA := shortLabels[rec.A]
		labelB := shortLabels[rec.B]
		fmt.Printf("  %-6s <-> %-6s : %d times\n", labelA, labelB, rec.Count)
	}
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

type Commit struct {
	Date      time.Time
	Author    string
	CoAuthors []string
}

func getGitCommitsSince(window string) ([]Commit, error) {
	// Validate window string (e.g., 1d, 2w, 3m, 1y)
	validWindow := regexp.MustCompile(`^\d+[dwmy]$`)
	if !validWindow.MatchString(window) {
		return nil, fmt.Errorf("invalid window format: %s", window)
	}
	sinceArg := windowToGitSince(window)
	cmd := exec.Command("git", "log", "--since="+sinceArg, "--pretty=format:%H%n%an <%ae>%n%ad%n%B%n==END==", "--date=iso")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	var commits []Commit
	var c Commit
	var bodyLines []string
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "==END==" {
			c.CoAuthors = parseCoAuthors(strings.Join(bodyLines, "\n"))
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
			c.Author = line
		case 2:
			t, err := time.Parse("2006-01-02 15:04:05 -0700", line)
			if err != nil {
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
	return commits, nil
}

// Convert window string to git --since argument
func windowToGitSince(window string) string {
	// git log accepts e.g. "1.day", "2.weeks", "3.months", "1.year"
	unitMap := map[byte]string{
		'd': "day",
		'w': "week",
		'm': "month",
		'y': "year",
	}
	n := window[:len(window)-1]
	unit := window[len(window)-1]
	if u, ok := unitMap[unit]; ok {
		return fmt.Sprintf("%s.%ss", n, u)
	}
	return window // fallback, should not happen due to validation
}

var coAuthorRe = regexp.MustCompile(`Co-authored-by:\s*(.+?)\s*<(.+?)>`)

func parseCoAuthors(body string) []string {
	var coauthors []string
	for _, line := range strings.Split(body, "\n") {
		m := coAuthorRe.FindStringSubmatch(line)
		if m != nil {
			coauthors = append(coauthors, fmt.Sprintf("%s <%s>", m[1], m[2]))
		}
	}
	return coauthors
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

// Generate short labels for each developer (e.g., initials or unique prefix)
func makeShortLabels(devs []string) map[string]string {
	labels := make(map[string]string, len(devs))
	used := make(map[string]struct{})
	for _, d := range devs {
		label := initialsFromEmail(d)
		origLabel := label
		i := 2
		for {
			if _, exists := used[label]; !exists {
				break
			}
			// If collision, append a number
			label = fmt.Sprintf("%s%d", origLabel, i)
			i++
		}
		labels[d] = label
		used[label] = struct{}{}
	}
	return labels
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
	at := strings.Index(email, "@")
	if at > 0 {
		parts := strings.Split(email[:at], ".")
		initials := ""
		for _, p := range parts {
			if len(p) > 0 {
				initials += strings.ToUpper(string(p[0]))
			}
		}
		if len(initials) == 0 {
			initials = strings.ToUpper(string(email[0]))
		}
		if len(initials) > 4 {
			initials = initials[:4]
		}
		return initials
	}
	if len(email) > 0 {
		return strings.ToUpper(string(email[0]))
	}
	return "??"
}

// Try to extract initials from "Name <email>"
func initialsFromAuthor(author string) string {
	name := author
	if idx := strings.Index(author, "<"); idx > 0 {
		name = strings.TrimSpace(author[:idx])
	}
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "??"
	}
	initials := ""
	for _, p := range parts {
		initials += strings.ToUpper(string(p[0]))
	}
	if len(initials) > 4 {
		initials = initials[:4]
	}
	return initials
}

func recommendPairs(devs []string, matrix map[Pair]int) []Recommendation {
	type pairKey struct{ A, B string }
	pairCounts := make([]Recommendation, 0)
	seen := make(map[pairKey]struct{})
	for i := 0; i < len(devs); i++ {
		for j := i + 1; j < len(devs); j++ {
			a, b := devs[i], devs[j]
			key := pairKey{A: a, B: b}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			count := matrix[Pair{A: a, B: b}]
			pairCounts = append(pairCounts, Recommendation{A: a, B: b, Count: count})
		}
	}
	// Sort by count ascending, then by names
	sort.Slice(pairCounts, func(i, j int) bool {
		if pairCounts[i].Count != pairCounts[j].Count {
			return pairCounts[i].Count < pairCounts[j].Count
		}
		if pairCounts[i].A != pairCounts[j].A {
			return pairCounts[i].A < pairCounts[j].A
		}
		return pairCounts[i].B < pairCounts[j].B
	})
	return pairCounts
}

// recommendPairsUnique returns a list of pairs such that each developer appears only once,
// optimizing for least-paired pairs (greedy matching).
func recommendPairsUnique(devs []string, matrix map[Pair]int) []Recommendation {
	type pairKey struct{ A, B string }
	pairCounts := make([]Recommendation, 0)
	used := make(map[string]bool)
	// Build all possible pairs with their counts
	for i := 0; i < len(devs); i++ {
		for j := i + 1; j < len(devs); j++ {
			a, b := devs[i], devs[j]
			count := matrix[Pair{A: a, B: b}]
			pairCounts = append(pairCounts, Recommendation{A: a, B: b, Count: count})
		}
	}
	// Sort by count ascending, then by names
	sort.Slice(pairCounts, func(i, j int) bool {
		if pairCounts[i].Count != pairCounts[j].Count {
			return pairCounts[i].Count < pairCounts[j].Count
		}
		if pairCounts[i].A != pairCounts[j].A {
			return pairCounts[i].A < pairCounts[j].A
		}
		return pairCounts[i].B < pairCounts[j].B
	})
	// Greedily pick pairs so each dev appears only once
	result := make([]Recommendation, 0)
	for _, rec := range pairCounts {
		if !used[rec.A] && !used[rec.B] {
			result = append(result, rec)
			used[rec.A] = true
			used[rec.B] = true
		}
	}
	return result
}
