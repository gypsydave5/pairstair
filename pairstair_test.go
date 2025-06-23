package main

import (
	"testing"
	"time"
)

func TestParseCoAuthors(t *testing.T) {
	body := `
Some commit message

Co-authored-by: Alice <alice@example.com>
Co-authored-by: Bob <bob@example.com>
`
	coauthors := parseCoAuthors(body)
	if len(coauthors) != 2 {
		t.Fatalf("expected 2 coauthors, got %d", len(coauthors))
	}
	if coauthors[0] != "Alice <alice@example.com>" {
		t.Errorf("unexpected coauthor: %s", coauthors[0])
	}
	if coauthors[1] != "Bob <bob@example.com>" {
		t.Errorf("unexpected coauthor: %s", coauthors[1])
	}
}

func TestMatrixLogic(t *testing.T) {
	commits := []Commit{
		{
			Date:      time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
			Author:    "Alice <alice@example.com>",
			CoAuthors: []string{"Bob <bob@example.com>"},
		},
		{
			Date:      time.Date(2024, 6, 1, 15, 0, 0, 0, time.UTC),
			Author:    "Bob <bob@example.com>",
			CoAuthors: []string{"Alice <alice@example.com>"},
		},
		{
			Date:      time.Date(2024, 6, 2, 10, 0, 0, 0, time.UTC),
			Author:    "Alice <alice@example.com>",
			CoAuthors: []string{"Carol <carol@example.com>"},
		},
	}

	matrix, _, _, _ := BuildPairMatrix(commits, []string{}, false)

	// Alice/Bob should have 1 (same day, only count once)
	a, b := "alice@example.com", "bob@example.com"
	if matrix[Pair{A: a, B: b}] != 1 && matrix[Pair{A: b, B: a}] != 1 {
		t.Errorf("expected Alice/Bob to have 1, got %d", matrix[Pair{A: a, B: b}])
	}
	// Alice/Carol should have 1
	c := "carol@example.com"
	if matrix[Pair{A: a, B: c}] != 1 && matrix[Pair{A: c, B: a}] != 1 {
		t.Errorf("expected Alice/Carol to have 1, got %d", matrix[Pair{A: a, B: c}])
	}
}
