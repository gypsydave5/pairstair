package main

// OutputRenderer provides a unified interface for different output formats
type OutputRenderer interface {
	Render(matrix *Matrix, recencyMatrix *RecencyMatrix, devs []string, shortLabels map[string]string, emailToName map[string]string, strategy string) error
}

// CLIRenderer handles console output
type CLIRenderer struct{}

// HTMLRenderer handles HTML output
type HTMLRenderer struct{}

// Render outputs the matrix and recommendations to the console
func (r *CLIRenderer) Render(matrix *Matrix, recencyMatrix *RecencyMatrix, devs []string, shortLabels map[string]string, emailToName map[string]string, strategy string) error {
	PrintMatrixCLI(matrix, devs, shortLabels, emailToName)
	PrintRecommendationsCLI(matrix, recencyMatrix, devs, shortLabels, strategy)
	return nil
}

// Render outputs the matrix and recommendations as HTML
func (r *HTMLRenderer) Render(matrix *Matrix, recencyMatrix *RecencyMatrix, devs []string, shortLabels map[string]string, emailToName map[string]string, strategy string) error {
	return RenderHTMLAndOpen(matrix, devs, shortLabels, emailToName)
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
