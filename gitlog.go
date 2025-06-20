package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Commit struct {
	Date      time.Time
	Author    string
	CoAuthors []string
}

func getGitCommitsSince(window string) ([]Commit, error) {
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

func windowToGitSince(window string) string {
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
	return window
}
