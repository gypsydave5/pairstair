package main

import (
	"fmt"
	"sort"
	"strings"
)

type Pair struct {
	A, B string
}

func BuildPairMatrix(commits []Commit, team []string, useTeam bool) (map[Pair]int, []string, map[string]string, map[string]string) {
	authorToEmail := make(map[string]string)
	emailToName := make(map[string]string)
	teamSet := make(map[string]struct{}, len(team))
	for _, t := range team {
		teamSet[t] = struct{}{}
	}
	if useTeam {
		for _, t := range team {
			email := extractEmail(t)
			name := extractName(t)
			authorToEmail[t] = email
			emailToName[email] = name
		}
	}

	datePairs := make(map[string]map[Pair]struct{})
	devsSet := make(map[string]struct{})

	for _, c := range commits {
		var devs []string
		if useTeam {
			if _, ok := teamSet[c.Author]; !ok {
				continue
			}
			filteredCoAuthors := make([]string, 0, len(c.CoAuthors))
			for _, ca := range c.CoAuthors {
				if _, ok := teamSet[ca]; ok {
					filteredCoAuthors = append(filteredCoAuthors, ca)
				}
			}
			devs = append([]string{c.Author}, filteredCoAuthors...)
		} else {
			devs = append([]string{c.Author}, c.CoAuthors...)
		}
		emailMap := make(map[string]struct{})
		for _, d := range devs {
			email := extractEmail(d)
			name := extractName(d)
			authorToEmail[d] = email
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
			continue
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

	devs := make([]string, 0, len(devsSet))
	for d := range devsSet {
		devs = append(devs, d)
	}
	if useTeam {
		for _, t := range team {
			email := extractEmail(t)
			if _, ok := devsSet[email]; !ok {
				devs = append(devs, email)
			}
		}
	}
	sort.Strings(devs)

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
	return matrix, devs, shortLabels, emailToName
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
