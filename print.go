package main

import (
	"fmt"
)

// Print the matrix and legend to the CLI
func PrintMatrixCLI(matrix map[Pair]int, devs []string, shortLabels map[string]string, emailToName map[string]string) {
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

// Print recommendations to the CLI
func PrintRecommendationsCLI(matrix map[Pair]int, devs []string, shortLabels map[string]string) {
	fmt.Println()
	if len(devs) > 10 {
		fmt.Println("Skipping pairing recommendations - too many developers (> 10)")
		return
	}

	fmt.Println("Pairing Recommendations (least-paired overall, optimal matching):")
	recommendations := recommendPairsOptimal(devs, matrix)
	for _, rec := range recommendations {
		labelA := shortLabels[rec.A]
		labelB := shortLabels[rec.B]
		if rec.B == "" {
			fmt.Printf("  %-6s (unpaired)\n", labelA)
		} else {
			fmt.Printf("  %-6s <-> %-6s : %d times\n", labelA, labelB, rec.Count)
		}
	}
}
