// Package output provides functionality for rendering pairing analysis results
// in different formats (CLI and HTML).
//
// The package provides a unified interface for different output formats,
// allowing the main application to render matrices and recommendations
// without being concerned with the specific output format details.
package output

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/gypsydave5/pairstair/internal/pairing"
)

// Recommendation represents a pair of developers who should work together
type Recommendation struct {
	A, B         string
	Count        int
	LastPaired   time.Time
	DaysSince    int
	HasPaired    bool
}

// OutputRenderer provides a unified interface for different output formats
type OutputRenderer interface {
	Render(matrix *pairing.Matrix, recencyMatrix *pairing.RecencyMatrix, devs []string, shortLabels map[string]string, emailToName map[string]string, strategy string, recommendations []Recommendation) error
}

// CLIRenderer handles console output
type CLIRenderer struct{}

// HTMLRenderer handles HTML output
type HTMLRenderer struct{}

// Render outputs the matrix and recommendations to the console
func (r *CLIRenderer) Render(matrix *pairing.Matrix, recencyMatrix *pairing.RecencyMatrix, devs []string, shortLabels map[string]string, emailToName map[string]string, strategy string, recommendations []Recommendation) error {
	PrintMatrixCLI(matrix, devs, shortLabels, emailToName)
	PrintRecommendationsCLI(recommendations, shortLabels, strategy)
	return nil
}

// Render outputs the matrix and recommendations as HTML
func (r *HTMLRenderer) Render(matrix *pairing.Matrix, recencyMatrix *pairing.RecencyMatrix, devs []string, shortLabels map[string]string, emailToName map[string]string, strategy string, recommendations []Recommendation) error {
	return RenderHTMLAndOpen(matrix, devs, shortLabels, emailToName, recommendations)
}

// NewRenderer creates the appropriate renderer based on output format
func NewRenderer(outputFormat string) OutputRenderer {
	switch outputFormat {
	case "html":
		return &HTMLRenderer{}
	default:
		return &CLIRenderer{}
	}
}

// PrintMatrixCLI prints the matrix and legend to the CLI
func PrintMatrixCLI(matrix *pairing.Matrix, devs []string, shortLabels map[string]string, emailToName map[string]string) {
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
			fmt.Printf("%-8d", matrix.Count(d1, d2))
		}
		fmt.Println()
	}
}

// PrintRecommendationsCLI prints recommendations to the CLI
func PrintRecommendationsCLI(recommendations []Recommendation, shortLabels map[string]string, strategy string) {
	fmt.Println()
	if len(recommendations) == 0 {
		fmt.Println("Skipping pairing recommendations - too many developers (> 10)")
		return
	}

	switch strategy {
	case "least-recent":
		fmt.Println("Pairing Recommendations (least recent collaborations first):")
	default: // least-paired
		fmt.Println("Pairing Recommendations (least-paired overall, optimal matching):")
	}
	
	for _, rec := range recommendations {
		labelA := shortLabels[rec.A]
		labelB := shortLabels[rec.B]
		if rec.B == "" {
			fmt.Printf("  %-6s (unpaired)\n", labelA)
		} else {
			if strategy == "least-recent" {
				if rec.HasPaired {
					if rec.DaysSince == 0 {
						fmt.Printf("  %-6s <-> %-6s : last paired today\n", labelA, labelB)
					} else if rec.DaysSince == 1 {
						fmt.Printf("  %-6s <-> %-6s : last paired 1 day ago\n", labelA, labelB)
					} else {
						fmt.Printf("  %-6s <-> %-6s : last paired %d days ago\n", labelA, labelB, rec.DaysSince)
					}
				} else {
					fmt.Printf("  %-6s <-> %-6s : never paired\n", labelA, labelB)
				}
			} else {
				fmt.Printf("  %-6s <-> %-6s : %d times\n", labelA, labelB, rec.Count)
			}
		}
	}
}

// RenderHTMLAndOpen renders HTML output and opens it in the default browser
func RenderHTMLAndOpen(matrix *pairing.Matrix, devs []string, shortLabels map[string]string, emailToName map[string]string, recommendations []Recommendation) error {
	html := renderHTML(matrix, devs, shortLabels, emailToName, recommendations)
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

// renderHTML generates HTML output for the matrix and recommendations
func renderHTML(matrix *pairing.Matrix, devs []string, shortLabels map[string]string, emailToName map[string]string, recommendations []Recommendation) string {
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
			b.WriteString(fmt.Sprintf("<td>%d</td>", matrix.Count(d1, d2)))
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</table>")

	// Recommendations
	b.WriteString("<div class=\"recommend\">")
	if len(recommendations) == 0 {
		b.WriteString("<h2>Pairing Recommendations</h2>")
		b.WriteString("<p>Skipping pairing recommendations - too many developers (> 10)</p>")
	} else {
		b.WriteString("<h2>Pairing Recommendations (least-paired overall, optimal matching)</h2><ul>")
		for _, rec := range recommendations {
			labelA := shortLabels[rec.A]
			labelB := shortLabels[rec.B]
			if rec.B == "" {
				b.WriteString(fmt.Sprintf("<li><b>%s</b> (unpaired)</li>", labelA))
			} else {
				b.WriteString(fmt.Sprintf("<li><b>%s</b> &lt;-&gt; <b>%s</b> : %d times</li>", labelA, labelB, rec.Count))
			}
		}
		b.WriteString("</ul>")
	}
	b.WriteString("</div>")

	b.WriteString("</body></html>")
	return b.String()
}

// openBrowser opens the given file path in the default web browser
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

// RecommendPairsOptimal generates pairing recommendations using greedy approach 
// (minimize total pair count, each dev appears once)
func RecommendPairsOptimal(devs []string, matrix *pairing.Matrix) []Recommendation {
	if len(devs) < 2 {
		return nil
	}
	
	if len(devs) > 10 {
		return []Recommendation{} // Return empty list for too many developers
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

// RecommendPairsLeastRecent generates pairing recommendations based on least recent collaboration
func RecommendPairsLeastRecent(devs []string, matrix *pairing.Matrix, recencyMatrix *pairing.RecencyMatrix) []Recommendation {
	n := len(devs)
	if n < 2 {
		return nil
	}
	
	if n > 10 {
		return []Recommendation{} // Return empty list for too many developers
	}

	type pairWithRecency struct {
		pair     pairing.Pair
		lastTime time.Time
		hasData  bool
		count    int
	}

	var allPairs []pairWithRecency
	now := time.Now()

	// Generate all possible pairs
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			pair := pairing.Pair{A: devs[i], B: devs[j]}
			if devs[i] > devs[j] {
				pair = pairing.Pair{A: devs[j], B: devs[i]}
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
