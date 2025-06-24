package main

import "strings"

type Developer struct {
	DisplayName     string
	EmailAddresses  []string
	AbbreviatedName string
}

func NewDeveloper(entry string) Developer {
	name := extractName(entry)
	emails := extractAllEmails(entry)
	if len(emails) == 0 {
		return Developer{}
	}

	return Developer{
		DisplayName:     name,
		EmailAddresses:  emails,
		AbbreviatedName: shortName(name),
	}
}

func (d Developer) CanonicalEmail() string {
	return d.EmailAddresses[0]
}

// Extract all emails from "Name <email1>,<email2>,<email3>"
func extractAllEmails(author string) []string {
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

// Extract just the name part from "Name <email>"
func extractName(author string) string {
	name := author
	if idx := strings.Index(author, "<"); idx > 0 {
		name = strings.TrimSpace(author[:idx])
	}
	return name
}

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
