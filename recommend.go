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
	limit := n
	if n%2 != 0 {
		limit = n - 1
	}
	best := make([]Recommendation, 0)
	minSum := -1
	perm := make([]string, n)
	copy(perm, devs)
	pairs := make([]Recommendation, 0, limit/2)
	used := make([]bool, n)

	var search func(pos int)
	search = func(pos int) {
		if pos == limit {
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
			if minSum == -1 || sum < minSum {
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
