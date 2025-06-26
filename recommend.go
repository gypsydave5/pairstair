package main

import (
	"sort"
	"time"
)

// Recommendation represents a pair of developers who should work together
type Recommendation struct {
	A, B         string
	Count        int
	LastPaired   time.Time
	DaysSince    int
	HasPaired    bool
}

// Optimal pairing using greedy approach (minimize total pair count, each dev appears once)
func recommendPairsOptimal(devs []string, matrix *Matrix) []Recommendation {
	if len(devs) < 2 {
		return nil
	}

	// Create all possible pairs with their counts
	type pairCandidate struct {
		devA, devB string
		count      int
	}

	var candidates []pairCandidate
	for i := 0; i < len(devs); i++ {
		for j := i + 1; j < len(devs); j++ {
			candidates = append(candidates, pairCandidate{
				devA:  devs[i],
				devB:  devs[j],
				count: matrix.Count(devs[i], devs[j]),
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
		if !used[candidate.devA] && !used[candidate.devB] {
			recommendations = append(recommendations, Recommendation{
				A:     candidate.devA,
				B:     candidate.devB,
				Count: candidate.count,
			})
			used[candidate.devA] = true
			used[candidate.devB] = true
		}
	}

	// Handle unpaired developer if odd number
	for _, dev := range devs {
		if !used[dev] {
			recommendations = append(recommendations, Recommendation{
				A:     dev,
				B:     "",
				Count: 0,
			})
			break
		}
	}

	return recommendations
}

// Recommend pairs based on least recent collaboration
func recommendPairsLeastRecent(devs []string, matrix *Matrix, recencyMatrix *RecencyMatrix) []Recommendation {
	n := len(devs)
	if n < 2 {
		return nil
	}

	type pairWithRecency struct {
		pair     Pair
		lastTime time.Time
		hasData  bool
		count    int
	}

	var allPairs []pairWithRecency
	now := time.Now()

	// Generate all possible pairs
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			pair := Pair{A: devs[i], B: devs[j]}
			if devs[i] > devs[j] {
				pair = Pair{A: devs[j], B: devs[i]}
			}

			lastTime, hasData := recencyMatrix.LastPaired(devs[i], devs[j])
			count := matrix.Count(devs[i], devs[j])

			allPairs = append(allPairs, pairWithRecency{
				pair:     pair,
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
		if used[pairData.pair.A] || used[pairData.pair.B] {
			continue
		}

		daysSince := 0
		if pairData.hasData {
			daysSince = int(now.Sub(pairData.lastTime).Hours() / 24)
		} else {
			daysSince = -1 // Never paired
		}

		recommendations = append(recommendations, Recommendation{
			A:          pairData.pair.A,
			B:          pairData.pair.B,
			Count:      pairData.count,
			LastPaired: pairData.lastTime,
			DaysSince:  daysSince,
			HasPaired:  pairData.hasData,
		})

		used[pairData.pair.A] = true
		used[pairData.pair.B] = true
	}

	// Handle odd number of developers
	if n%2 != 0 {
		for _, dev := range devs {
			if !used[dev] {
				recommendations = append(recommendations, Recommendation{
					A:         dev,
					B:         "",
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
