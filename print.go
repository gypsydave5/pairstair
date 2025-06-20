package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Recommendation struct {
	A, B  string
	Count int
}

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

func RenderHTMLAndOpen(matrix map[Pair]int, devs []string, shortLabels map[string]string, emailToName map[string]string) error {
	html := RenderHTML(matrix, devs, shortLabels, emailToName)
	tmpfile, err := os.CreateTemp("", "pairstair-*.html")
	if err != nil {
		return err
	}
	defer tmpfile.Close()
	_, err = tmpfile.WriteString(html)
	if err != nil {
		return err
	}
	// Open in default browser
	return openBrowser(tmpfile.Name())
}

func RenderHTML(matrix map[Pair]int, devs []string, shortLabels map[string]string, emailToName map[string]string) string {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html><head><meta charset=\"utf-8\"><title>Pair Stair</title>")
	b.WriteString(`<style>
body { font-family: sans-serif; margin: 2em; }
table { border-collapse: collapse; }
th, td { border: 1px solid #ccc; padding: 0.5em 1em; text-align: center; }
th { background: #eee; }
.legend-table { margin-bottom: 2em; }
.recommend { margin-top: 2em; }
</style></head><body>`)
	b.WriteString("<h1>Pair Stair Matrix</h1>")

	// Legend
	b.WriteString("<h2>Legend</h2><table class=\"legend-table\"><tr><th>Initials</th><th>Name</th><th>Email</th></tr>")
	for _, d := range devs {
		name := emailToName[d]
		if name == "" {
			name = d
		}
		b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", shortLabels[d], name, d))
	}
	b.WriteString("</table>")

	// Matrix
	b.WriteString("<h2>Pair Matrix</h2><table><tr><th></th>")
	for _, d := range devs {
		b.WriteString(fmt.Sprintf("<th>%s</th>", shortLabels[d]))
	}
	b.WriteString("</tr>")
	for _, d1 := range devs {
		b.WriteString(fmt.Sprintf("<tr><th>%s</th>", shortLabels[d1]))
		for _, d2 := range devs {
			if d1 == d2 {
				b.WriteString("<td>-</td>")
				continue
			}
			a, b2 := d1, d2
			if a > b2 {
				a, b2 = b2, a
			}
			b.WriteString(fmt.Sprintf("<td>%d</td>", matrix[Pair{A: a, B: b2}]))
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</table>")

	// Recommendations
	b.WriteString("<div class=\"recommend\"><h2>Pairing Recommendations (least-paired overall, optimal matching)</h2><ul>")
	recommendations := recommendPairsOptimal(devs, matrix)
	for _, rec := range recommendations {
		labelA := shortLabels[rec.A]
		labelB := shortLabels[rec.B]
		if rec.B == "" {
			b.WriteString(fmt.Sprintf("<li><b>%s</b> (unpaired)</li>", labelA))
		} else {
			b.WriteString(fmt.Sprintf("<li><b>%s</b> &lt;-&gt; <b>%s</b> : %d times</li>", labelA, labelB, rec.Count))
		}
	}
	b.WriteString("</ul></div>")

	b.WriteString("</body></html>")
	return b.String()
}

func openBrowser(path string) error {
	url := path
	if !strings.HasPrefix(url, "file://") {
		url = "file://" + url
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
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
