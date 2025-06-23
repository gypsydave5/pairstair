package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func RenderHTMLAndOpen(matrix map[Pair]int, devs []string, shortLabels map[string]string, emailToName map[string]string) error {
	html := renderHTML(matrix, devs, shortLabels, emailToName)
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

func renderHTML(matrix map[Pair]int, devs []string, shortLabels map[string]string, emailToName map[string]string) string {
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
	b.WriteString("<div class=\"recommend\">")
	if len(devs) > 10 {
		b.WriteString("<h2>Pairing Recommendations</h2>")
		b.WriteString("<p>Skipping pairing recommendations - too many developers (> 10)</p>")
	} else {
		b.WriteString("<h2>Pairing Recommendations (least-paired overall, optimal matching)</h2><ul>")
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
		b.WriteString("</ul>")
	}
	b.WriteString("</div>")

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
