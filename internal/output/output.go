// Package output provides functionality for rendering pairing analysis results
// in different formats (CLI and HTML).
//
// The package provides a unified interface for different output formats,
// allowing the main application to render matrices and recommendations
// without being concerned with the specific output format details.
package output

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gypsydave5/pairstair/internal/git"
	"github.com/gypsydave5/pairstair/internal/pairing"
	"github.com/gypsydave5/pairstair/internal/recommend"
)

// OutputRenderer provides a unified interface for different output formats
type OutputRenderer interface {
	Render(matrix *pairing.Matrix, recencyMatrix *pairing.RecencyMatrix, developers []git.Developer, strategy string, recommendations []recommend.Recommendation) error
}

// CLIRenderer handles console output
type CLIRenderer struct{}

// HTMLRenderer handles HTML output
type HTMLRenderer struct {
	OpenInBrowser bool
}

// Render outputs the matrix and recommendations to the console
func (r *CLIRenderer) Render(matrix *pairing.Matrix, recencyMatrix *pairing.RecencyMatrix, developers []git.Developer, strategy string, recommendations []recommend.Recommendation) error {
	PrintMatrixCLI(matrix, developers)
	PrintRecommendationsCLI(recommendations, strategy)
	return nil
}

// Render outputs the matrix and recommendations as HTML
func (r *HTMLRenderer) Render(matrix *pairing.Matrix, recencyMatrix *pairing.RecencyMatrix, developers []git.Developer, strategy string, recommendations []recommend.Recommendation) error {
	if r.OpenInBrowser {
		return RenderHTMLAndOpen(matrix, developers, recommendations)
	} else {
		return RenderHTMLToWriter(os.Stdout, matrix, developers, recommendations)
	}
}

// NewRenderer creates the appropriate renderer based on output format
// This is kept for backward compatibility and defaults to not opening browser
func NewRenderer(outputFormat string) OutputRenderer {
	return NewRendererWithOpen(outputFormat, false)
}

// NewRendererWithOpen creates the appropriate renderer based on output format and open behavior
func NewRendererWithOpen(outputFormat string, openInBrowser bool) OutputRenderer {
	switch outputFormat {
	case "html":
		return &HTMLRenderer{OpenInBrowser: openInBrowser}
	default:
		return &CLIRenderer{}
	}
}

// PrintMatrixCLI prints the matrix and legend to the CLI
func PrintMatrixCLI(matrix *pairing.Matrix, developers []git.Developer) {
	fmt.Println("Legend:")
	for _, dev := range developers {
		fmt.Printf("  %-6s = %-20s %s\n", dev.AbbreviatedName, dev.DisplayName, dev.CanonicalEmail())
	}
	fmt.Println()

	fmt.Printf("%-8s", "")
	for _, dev := range developers {
		fmt.Printf("%-8s", dev.AbbreviatedName)
	}
	fmt.Println()
	for _, dev1 := range developers {
		fmt.Printf("%-8s", dev1.AbbreviatedName)
		for _, dev2 := range developers {
			if dev1.CanonicalEmail() == dev2.CanonicalEmail() {
				fmt.Printf("%-8s", "-")
				continue
			}
			fmt.Printf("%-8d", matrix.Count(dev1.CanonicalEmail(), dev2.CanonicalEmail()))
		}
		fmt.Println()
	}
}

// PrintRecommendationsCLI prints recommendations to the CLI
func PrintRecommendationsCLI(recommendations []recommend.Recommendation, strategy string) {
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
		if len(rec.B.EmailAddresses) == 0 {
			fmt.Printf("  %-6s (unpaired)\n", rec.A.AbbreviatedName)
		} else {
			if strategy == "least-recent" {
				if rec.HasPaired {
					if rec.DaysSince == 0 {
						fmt.Printf("  %-6s <-> %-6s : last paired today\n", rec.A.AbbreviatedName, rec.B.AbbreviatedName)
					} else if rec.DaysSince == 1 {
						fmt.Printf("  %-6s <-> %-6s : last paired 1 day ago\n", rec.A.AbbreviatedName, rec.B.AbbreviatedName)
					} else {
						fmt.Printf("  %-6s <-> %-6s : last paired %d days ago\n", rec.A.AbbreviatedName, rec.B.AbbreviatedName, rec.DaysSince)
					}
				} else {
					fmt.Printf("  %-6s <-> %-6s : never paired\n", rec.A.AbbreviatedName, rec.B.AbbreviatedName)
				}
			} else {
				fmt.Printf("  %-6s <-> %-6s : %d times\n", rec.A.AbbreviatedName, rec.B.AbbreviatedName, rec.Count)
			}
		}
	}
}

// RenderHTMLAndOpen renders HTML output and opens it in the default browser
func RenderHTMLAndOpen(matrix *pairing.Matrix, developers []git.Developer, recommendations []recommend.Recommendation) error {
	tmpfile, err := os.CreateTemp("", "pairstair-*.html")
	if err != nil {
		return err
	}
	defer tmpfile.Close()

	err = RenderHTMLToWriter(tmpfile, matrix, developers, recommendations)
	if err != nil {
		return err
	}

	// Open in default browser
	return openBrowser(tmpfile.Name())
}

// RenderHTMLToWriter renders HTML output to the provided io.Writer
// This is the testable version of HTML rendering that can write to any Writer
func RenderHTMLToWriter(w io.Writer, matrix *pairing.Matrix, developers []git.Developer, recommendations []recommend.Recommendation) error {
	html := renderHTML(matrix, developers, recommendations)
	_, err := w.Write([]byte(html))
	return err
}

// renderHTML generates HTML output for the matrix and recommendations
func renderHTML(matrix *pairing.Matrix, developers []git.Developer, recommendations []recommend.Recommendation) string {
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
	for _, dev := range developers {
		b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", dev.AbbreviatedName, dev.DisplayName, dev.CanonicalEmail()))
	}
	b.WriteString("</table>")

	// Matrix
	b.WriteString("<h2>Pair Matrix</h2><table><tr><th></th>")
	for _, dev := range developers {
		b.WriteString(fmt.Sprintf("<th>%s</th>", dev.AbbreviatedName))
	}
	b.WriteString("</tr>")
	for _, dev1 := range developers {
		b.WriteString(fmt.Sprintf("<tr><th>%s</th>", dev1.AbbreviatedName))
		for _, dev2 := range developers {
			if dev1.CanonicalEmail() == dev2.CanonicalEmail() {
				b.WriteString("<td>-</td>")
				continue
			}
			b.WriteString(fmt.Sprintf("<td>%d</td>", matrix.Count(dev1.CanonicalEmail(), dev2.CanonicalEmail())))
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
			if len(rec.B.EmailAddresses) == 0 {
				b.WriteString(fmt.Sprintf("<li><b>%s</b> (unpaired)</li>", rec.A.AbbreviatedName))
			} else {
				b.WriteString(fmt.Sprintf("<li><b>%s</b> &lt;-&gt; <b>%s</b> : %d times</li>", rec.A.AbbreviatedName, rec.B.AbbreviatedName, rec.Count))
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
