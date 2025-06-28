// Package recommend provides functionality for generating pairing recommendations
// based on different strategies and developer collaboration history.
//
// The package focuses purely on recommendation algorithms and logic,
// separating this concern from output formatting and display.
package recommend

import (
	"sort"
	"time"

	"github.com/gypsydave5/pairstair/internal/git"
	"github.com/gypsydave5/pairstair/internal/pairing"
)

// Recommendation represents a pairing recommendation for developers
type Recommendation struct {
	A, B       git.Developer
	Count      int
	LastPaired time.Time
	DaysSince  int
	HasPaired  bool
}

// Strategy represents a recommendation strategy
type Strategy string

const (
	LeastPaired Strategy = "least-paired"
	LeastRecent Strategy = "least-recent"
)

// GenerateRecommendations generates pairing recommendations using the specified strategy
func GenerateRecommendations(developers []git.Developer, matrix *pairing.Matrix, recencyMatrix *pairing.RecencyMatrix, strategy Strategy) []Recommendation {
	switch strategy {
	case LeastRecent:
		return generateLeastRecent(developers, matrix, recencyMatrix)
	default: // LeastPaired
		return generateLeastPaired(developers, matrix)
	}
}

// generateLeastPaired generates pairing recommendations using greedy approach
// (minimize total pair count, each dev appears once)
func generateLeastPaired(developers []git.Developer, matrix *pairing.Matrix) []Recommendation {
	if len(developers) < 2 {
		return nil
	}

	if len(developers) > 20 {
		return []Recommendation{} // Return empty list for too many developers
	}

	// Create all possible pairs with their counts
	type pairCandidate struct {
		devA, devB git.Developer
		count      int
	}

	var candidates []pairCandidate
	for i := 0; i < len(developers); i++ {
		for j := i + 1; j < len(developers); j++ {
			candidates = append(candidates, pairCandidate{
				devA:  developers[i],
				devB:  developers[j],
				count: matrix.Count(developers[i].CanonicalEmail(), developers[j].CanonicalEmail()),
			})
		}
	}

	// Sort by count (ascending - least paired first)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].count < candidates[j].count
	})

	// Greedily select pairs ensuring each dev appears only once
	used := make(map[string]bool)
	var recommendations []Recommendation

	for _, candidate := range candidates {
		emailA := candidate.devA.CanonicalEmail()
		emailB := candidate.devB.CanonicalEmail()
		if !used[emailA] && !used[emailB] {
			recommendations = append(recommendations, Recommendation{
				A:     candidate.devA,
				B:     candidate.devB,
				Count: candidate.count,
			})
			used[emailA] = true
			used[emailB] = true
		}
	}

	// Handle unpaired developer if odd number
	for _, dev := range developers {
		email := dev.CanonicalEmail()
		if !used[email] {
			recommendations = append(recommendations, Recommendation{
				A:     dev,
				B:     git.Developer{}, // Empty Developer object for unpaired
				Count: 0,
			})
			break
		}
	}

	return recommendations
}

// generateLeastRecent generates pairing recommendations based on least recent collaboration
func generateLeastRecent(developers []git.Developer, matrix *pairing.Matrix, recencyMatrix *pairing.RecencyMatrix) []Recommendation {
	n := len(developers)
	if n < 2 {
		return nil
	}

	if n > 20 {
		return []Recommendation{} // Return empty list for too many developers
	}

	type pairWithRecency struct {
		devA, devB git.Developer
		lastTime   time.Time
		hasData    bool
		count      int
	}

	var allPairs []pairWithRecency
	now := time.Now()

	// Generate all possible pairs
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			devA := developers[i]
			devB := developers[j]
			emailA := devA.CanonicalEmail()
			emailB := devB.CanonicalEmail()

			lastTime, hasData := recencyMatrix.LastPaired(emailA, emailB)
			count := matrix.Count(emailA, emailB)

			allPairs = append(allPairs, pairWithRecency{
				devA:     devA,
				devB:     devB,
				lastTime: lastTime,
				hasData:  hasData,
				count:    count,
			})
		}
	}

	// Sort pairs by recency (least recent first)
	sort.Slice(allPairs, func(i, j int) bool {
		// Pairs that have never worked together come first
		if !allPairs[i].hasData && allPairs[j].hasData {
			return true
		}
		if allPairs[i].hasData && !allPairs[j].hasData {
			return false
		}
		if !allPairs[i].hasData && !allPairs[j].hasData {
			return false // Both have no data, order doesn't matter
		}
		// Both have data, sort by oldest first
		return allPairs[i].lastTime.Before(allPairs[j].lastTime)
	})

	// Create recommendations using a greedy approach
	var recommendations []Recommendation
	used := make(map[string]bool)

	for _, pairData := range allPairs {
		emailA := pairData.devA.CanonicalEmail()
		emailB := pairData.devB.CanonicalEmail()
		if used[emailA] || used[emailB] {
			continue
		}

		daysSince := 0
		if pairData.hasData {
			daysSince = int(now.Sub(pairData.lastTime).Hours() / 24)
		} else {
			daysSince = -1 // Never paired
		}

		recommendations = append(recommendations, Recommendation{
			A:          pairData.devA,
			B:          pairData.devB,
			Count:      pairData.count,
			LastPaired: pairData.lastTime,
			DaysSince:  daysSince,
			HasPaired:  pairData.hasData,
		})

		used[emailA] = true
		used[emailB] = true
	}

	// Handle odd number of developers
	if n%2 != 0 {
		for _, dev := range developers {
			email := dev.CanonicalEmail()
			if !used[email] {
				recommendations = append(recommendations, Recommendation{
					A:         dev,
					B:         git.Developer{}, // Empty Developer object for unpaired
					Count:     0,
					DaysSince: 0,
					HasPaired: false,
				})
				break
			}
		}
	}

	return recommendations
}
