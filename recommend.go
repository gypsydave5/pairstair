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

// Optimal pairing using brute-force for small N (minimize total pair count, each dev appears once)
func recommendPairsOptimal(devs []string, matrix *Matrix) []Recommendation {
	n := len(devs)
	if n < 2 {
		return nil
	}

	// Set limit to the largest even number <= n
	limit := n - (n % 2)

	best := make([]Recommendation, 0)
	minSum := -1
	perm := make([]string, n)
	copy(perm, devs)
	pairs := make([]Recommendation, 0, limit/2)
	used := make([]bool, n)

	var search func(pos int)
	search = func(pos int) {
		// Base case: we've processed enough developers to make all possible pairs
		if pos >= limit {
			pairs = pairs[:0]
			// Create pairs from current permutation
			for i := 0; i < limit; i += 2 {
				a, b := perm[i], perm[i+1]
				count := matrix.Count(a, b)
				pairs = append(pairs, Recommendation{A: a, B: b, Count: count})
			}

			// Calculate total pairing count
			sum := 0
			for _, p := range pairs {
				sum += p.Count
			}

			// Keep the best (lowest) pairing count
			if minSum == -1 || sum < minSum {
				minSum = sum
				best = append([]Recommendation(nil), pairs...)
			}
			return
		}

		// Skip already used positions
		if used[pos] {
			search(pos + 1)
			return
		}

		used[pos] = true
		for j := pos + 1; j < n; j++ {
			if used[j] {
				continue
			}
			used[j] = true
			perm[pos], perm[j] = perm[j], perm[pos]

			// Make sure we don't go beyond array bounds
			if pos+2 < n {
				search(pos + 2)
			} else {
				search(n) // This will trigger the base case
			}

			perm[pos], perm[j] = perm[j], perm[pos]
			used[j] = false
		}
		used[pos] = false
	}

	search(0)

	// Handle odd number of developers - add the unpaired developer
	if n%2 != 0 {
		unpaired := ""
		usedMap := make(map[string]bool)
		for _, r := range best {
			usedMap[r.A] = true
			usedMap[r.B] = true
		}
		for _, d := range devs {
			if !usedMap[d] {
				unpaired = d
				break
			}
		}
		if unpaired != "" {
			best = append(best, Recommendation{A: unpaired, B: "", Count: 0})
		}
	}

	return best
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
