package main

// Recommendation represents a pair of developers who should work together
type Recommendation struct {
	A, B  string
	Count int
}

// Optimal pairing using brute-force for small N (minimize total pair count, each dev appears once)
func recommendPairsOptimal(devs []string, matrix map[Pair]int) []Recommendation {
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
				pa, pb := a, b
				if pa > pb {
					pa, pb = pb, pa
				}
				count := matrix[Pair{A: pa, B: pb}]
				pairs = append(pairs, Recommendation{A: pa, B: pb, Count: count})
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
