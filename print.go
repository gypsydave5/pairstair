package main

import (
	"fmt"
	"sort"
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
	fmt.Println("Pairing Recommendations (least-paired first):")
	recommendations := recommendPairsUnique(devs, matrix)
	for _, rec := range recommendations {
		labelA := shortLabels[rec.A]
		labelB := shortLabels[rec.B]
		fmt.Printf("  %-6s <-> %-6s : %d times\n", labelA, labelB, rec.Count)
	}
}

// recommendPairsUnique returns a list of pairs such that each developer appears only once,
// optimizing for least-paired pairs (greedy matching).
func recommendPairsUnique(devs []string, matrix map[Pair]int) []Recommendation {
	type pairKey struct{ A, B string }
	pairCounts := make([]Recommendation, 0)
	used := make(map[string]bool)
	for i := 0; i < len(devs); i++ {
		for j := i + 1; j < len(devs); j++ {
			a, b := devs[i], devs[j]
			count := matrix[Pair{A: a, B: b}]
			pairCounts = append(pairCounts, Recommendation{A: a, B: b, Count: count})
		}
	}
	sort.Slice(pairCounts, func(i, j int) bool {
		if pairCounts[i].Count != pairCounts[j].Count {
			return pairCounts[i].Count < pairCounts[j].Count
		}
		if pairCounts[i].A != pairCounts[j].A {
			return pairCounts[i].A < pairCounts[j].A
		}
		return pairCounts[i].B < pairCounts[j].B
	})
	result := make([]Recommendation, 0)
	for _, rec := range pairCounts {
		if !used[rec.A] && !used[rec.B] {
			result = append(result, rec)
			used[rec.A] = true
			used[rec.B] = true
		}
	}
	return result
}
