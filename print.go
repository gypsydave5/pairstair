package main

import (
	"fmt"
	"math"
)

type Recommendation struct {
	A, B  string
	Count int
}

func PrintPairMatrix(matrix map[Pair]int, devs []string, shortLabels map[string]string, emailToName map[string]string) {
	fmt.Println("Legend:")
	for _, d := range devs {
		name := emailToName[d]
		if name == "" {
			name = d
		}
		fmt.Printf("  %-6s = %-20s %s\n", shortLabels[d], name, d)
	}
	fmt.Println()

	fmt.Printf("%-8s", "")
	for _, d := range devs {
		fmt.Printf("%-8s", shortLabels[d])
	}
	fmt.Println()
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
}

func PrintPairRecommendations(matrix map[Pair]int, devs []string, shortLabels map[string]string) {
	fmt.Println()
	fmt.Println("Pairing Recommendations (least-paired overall, optimal matching):")
	recommendations := recommendPairsOptimal(devs, matrix)
	for _, rec := range recommendations {
		labelA := shortLabels[rec.A]
		labelB := shortLabels[rec.B]
		fmt.Printf("  %-6s <-> %-6s : %d times\n", labelA, labelB, rec.Count)
	}
}

// Optimal pairing using brute-force for small N (minimize total pair count, each dev appears once)
func recommendPairsOptimal(devs []string, matrix map[Pair]int) []Recommendation {
	n := len(devs)
	if n < 2 {
		return nil
	}
	// Only even number of devs supported for perfect matching
	// If odd, leave one out (not paired)
	limit := n
	if n%2 != 0 {
		limit = n - 1
	}
	best := make([]Recommendation, 0)
	minSum := math.MaxInt
	perm := make([]string, n)
	copy(perm, devs)
	pairs := make([]Recommendation, 0, limit/2)
	used := make([]bool, n)

	var search func(pos int)
	search = func(pos int) {
		if pos == limit {
			// Evaluate this pairing
			pairs = pairs[:0]
			for i := 0; i < limit; i += 2 {
				a, b := perm[i], perm[i+1]
				pa, pb := a, b
				if pa > pb {
					pa, pb = pb, pa
				}
				count := matrix[Pair{A: pa, B: pb}]
				pairs = append(pairs, Recommendation{A: pa, B: pb, Count: count})
			}
			sum := 0
			for _, p := range pairs {
				sum += p.Count
			}
			if sum < minSum {
				minSum = sum
				best = append([]Recommendation(nil), pairs...)
			}
			return
		}
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
			search(pos + 2)
			perm[pos], perm[j] = perm[j], perm[pos]
			used[j] = false
		}
		used[pos] = false
	}

	search(0)
	// If odd, add the unpaired dev as a single
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
