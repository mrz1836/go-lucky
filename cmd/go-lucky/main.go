// Package main provides the command-line interface for the NC Lucky for Life lottery analyzer.
//
// This cosmic-powered lottery analysis tool combines statistical analysis with celestial correlations
// to provide entertaining insights into lottery patterns. Whether you believe in the stars or statistics,
// this analyzer offers multiple perspectives on lottery data - all while maintaining a healthy skepticism
// about the truly random nature of lottery drawings.
//
// Features include frequency analysis, pattern detection, gap analysis, and whimsical cosmic correlations
// with moon phases, planetary positions, and other astronomical phenomena. Perfect for the mathematically
// curious lottery enthusiast who enjoys a good laugh at the universe's expense.
//
// Remember: Past performance does not predict future results, the house always wins, and the stars
// are not responsible for your gambling decisions! 🎲✨
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	ctx := context.Background()

	// Default configuration
	config := &AnalysisConfig{
		RecentWindow:     50,
		MinGapMultiplier: 1.5,
		ConfidenceLevel:  0.95,
		OutputMode:       "detailed",
		ExportFormat:     "console",
	}

	// Parse command line arguments
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "--simple":
				config.OutputMode = "simple"
			case "--statistical":
				config.OutputMode = "statistical"
			case "--cosmic":
				config.OutputMode = "cosmic"
			case "--export-json":
				config.ExportFormat = "json"
			case "--export-csv":
				config.ExportFormat = "csv"
			case "--recent":
				if i+1 < len(os.Args) {
					if val, err := strconv.Atoi(os.Args[i+1]); err == nil {
						config.RecentWindow = val
						i++
					}
				}
			case "--help":
				printHelp()
				return
			}
		}
	}

	// Create analyzer
	analyzer, err := NewAnalyzer(ctx, "../../data/lucky-numbers-history.csv", config)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Run analysis
	if err = analyzer.RunAnalysis(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error running analysis: %v\n", err)
		os.Exit(1)
	}

	// Export if requested
	if config.ExportFormat != "console" {
		filename := fmt.Sprintf("lottery_analysis_%s.%s",
			time.Now().Format("20060102_150405"),
			config.ExportFormat)
		if err = analyzer.ExportAnalysis(ctx, filename); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error exporting analysis: %v\n", err)
			os.Exit(1)
		}
		_, _ = fmt.Fprintf(os.Stdout, "\nAnalysis exported to: %s\n", filename)
	}
}

// printHelp displays usage information
func printHelp() {
	_, _ = fmt.Fprintln(os.Stdout, "NC Lucky for Life Lottery Analyzer")
	_, _ = fmt.Fprintln(os.Stdout, "==================================")
	_, _ = fmt.Fprintln(os.Stdout)
	_, _ = fmt.Fprintln(os.Stdout, "Usage: go run lottery_analyzer.go [options]")
	_, _ = fmt.Fprintln(os.Stdout)
	_, _ = fmt.Fprintln(os.Stdout, "Options:")
	_, _ = fmt.Fprintln(os.Stdout, "  --simple           Show simplified analysis")
	_, _ = fmt.Fprintln(os.Stdout, "  --statistical      Show detailed statistical analysis")
	_, _ = fmt.Fprintln(os.Stdout, "  --cosmic           Show cosmic correlation analysis")
	_, _ = fmt.Fprintln(os.Stdout, "  --export-json      Export results to JSON file")
	_, _ = fmt.Fprintln(os.Stdout, "  --export-csv       Export results to CSV file")
	_, _ = fmt.Fprintln(os.Stdout, "  --recent <n>       Set recent window size (default: 50)")
	_, _ = fmt.Fprintln(os.Stdout, "  --help             Show this help message")
	_, _ = fmt.Fprintln(os.Stdout)
	_, _ = fmt.Fprintln(os.Stdout, "Examples:")
	_, _ = fmt.Fprintln(os.Stdout, "  go run lottery_analyzer.go")
	_, _ = fmt.Fprintln(os.Stdout, "  go run lottery_analyzer.go --simple")
	_, _ = fmt.Fprintln(os.Stdout, "  go run lottery_analyzer.go --statistical --export-json")
	_, _ = fmt.Fprintln(os.Stdout, "  go run lottery_analyzer.go --recent 100")
}
