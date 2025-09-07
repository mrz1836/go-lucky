//go:build mage

// Package magefiles provides build automation tasks for the Go-Lucky Lottery Analyzer.
package magefiles

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Variables
const (
	binaryName   = "go-lucky"
	binaryPath   = "./bin/go-lucky"
	coverageFile = "coverage.out"
	testTimeout  = "30s"
)

// logInfo prints informational messages to stdout
func logInfo(msg string) {
	_, _ = fmt.Fprintln(os.Stdout, msg)
}

// Build namespace for build-related tasks
type Build mg.Namespace

// Analysis namespace for lottery analysis commands
type Analysis mg.Namespace

// Export namespace for export commands
type Export mg.Namespace

// Quick namespace for quick analysis commands
type Quick mg.Namespace

// Fun namespace for entertainment commands
type Fun mg.Namespace

// Helper function to ensure binary exists and is built
func ensureBinary() error {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return Build{}.Dev()
	}
	return nil
}

// Helper function to run the analyzer with given arguments
func runAnalyzer(args ...string) error {
	if err := ensureBinary(); err != nil {
		return err
	}
	return sh.RunV(binaryPath, args...)
}

// Helper function to run analyzer and capture output
func runAnalyzerWithOutput(args ...string) (string, error) {
	if err := ensureBinary(); err != nil {
		return "", err
	}
	return sh.Output(binaryPath, args...)
}

// TestQuick runs fast unit tests excluding performance tests
func TestQuick() error {
	return sh.RunV("go", "test", "-short", "./...")
}

// Test runs all tests
func Test() error {
	return sh.RunV("go", "test", "-timeout", testTimeout, "./...")
}

// Benchmark runs performance benchmarks
func Benchmark() error {
	return sh.RunV("go", "test", "-bench=.", "-benchmem", "./...")
}

// Clean removes build artifacts and generated files
func Clean() error {
	logInfo("🧹 Cleaning up...")

	// Remove directories
	dirs := []string{"bin/"}
	for _, dir := range dirs {
		if err := sh.Rm(dir); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", dir, err)
		}
	}

	// Remove specific files
	files := []string{
		binaryName,
		coverageFile,
		"coverage.html",
		"coverage.txt",
	}

	// Remove generated files with patterns
	patterns := []string{
		"lottery_analysis_*.json",
		"lottery_analysis_*.csv",
		"test_*.csv",
		"debug_*.csv",
		"empty_*.csv",
		"invalid_*.csv",
	}

	for _, file := range files {
		if err := sh.Rm(file); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", file, err)
		}
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("failed to glob %s: %w", pattern, err)
		}
		for _, match := range matches {
			if err := sh.Rm(match); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove %s: %w", match, err)
			}
		}
	}

	logInfo("✅ Cleanup complete")
	return nil
}

// Build Commands

// Dev builds the development version of the analyzer
func (Build) Dev() error {
	logInfo("🔧 Building development version...")

	// Ensure bin directory exists
	if err := os.MkdirAll("bin", 0o750); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Build using magex
	return sh.RunV("magex", "build:dev")
}

// Default builds the default version
func (Build) Default() error {
	return sh.RunV("magex", "build:default")
}

// All builds all platforms
func (Build) All() error {
	return sh.RunV("magex", "build:all")
}

// Analysis Commands

// Full runs complete analysis with cosmic correlations (RECOMMENDED)
func (Analysis) Full() error {
	mg.Deps(Build{}.Dev)

	logInfo("╔══════════════════════════════════════════════════════════════╗")
	logInfo("║        🌌 RUNNING FULL COSMIC LOTTERY ANALYSIS 🌌            ║")
	logInfo("╚══════════════════════════════════════════════════════════════╝")
	logInfo("")

	return runAnalyzer("--cosmic")
}

// Simple runs simple analysis summary
func (Analysis) Simple() error {
	mg.Deps(Build{}.Dev)
	return runAnalyzer("--simple")
}

// Statistical runs detailed statistical analysis
func (Analysis) Statistical() error {
	mg.Deps(Build{}.Dev)
	return runAnalyzer("--statistical")
}

// Cosmic runs cosmic correlation analysis only
func (Analysis) Cosmic() error {
	mg.Deps(Build{}.Dev)
	return runAnalyzer("--cosmic")
}

// Export Commands

// JSON exports full analysis to JSON file
func (Export) JSON() error {
	mg.Deps(Build{}.Dev)

	logInfo("📊 Exporting analysis to JSON...")
	if err := runAnalyzer("--cosmic", "--export-json"); err != nil {
		return err
	}
	logInfo("✅ Export complete! Check lottery_analysis_*.json")
	return nil
}

// CSV exports analysis data to CSV file
func (Export) CSV() error {
	mg.Deps(Build{}.Dev)

	logInfo("📊 Exporting analysis to CSV...")
	if err := runAnalyzer("--cosmic", "--export-csv"); err != nil {
		return err
	}
	logInfo("✅ Export complete! Check lottery_analysis_*.csv")
	return nil
}

// Quick Analysis Commands

// LuckyPicks generates 5 different analysis-based number sets
func (Quick) LuckyPicks() error {
	mg.Deps(Build{}.Dev)

	logInfo("🎰 Generating Lucky Picks...")
	logInfo("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	output, err := runAnalyzerWithOutput("--simple")
	if err != nil {
		return err
	}

	lines := strings.Split(output, "\n")
	printSection := false
	lineCount := 0

	for _, line := range lines {
		if strings.Contains(line, "QUICK PICKS:") {
			printSection = true
			lineCount = 0
		}
		if printSection {
			logInfo(line)
			lineCount++
			if lineCount > 10 {
				break
			}
		}
	}

	logInfo("")
	logInfo("🌌 Cosmic Pick:")
	for _, line := range lines {
		if strings.Contains(line, "COSMIC PICK:") {
			logInfo(line)
			break
		}
	}
	logInfo("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

// HotNumbers shows current hot numbers
func (Quick) HotNumbers() error {
	mg.Deps(Build{}.Dev)

	logInfo("🔥 Current Hot Numbers:")
	output, err := runAnalyzerWithOutput("--simple")
	if err != nil {
		return err
	}

	lines := strings.Split(output, "\n")
	printSection := false
	lineCount := 0

	for _, line := range lines {
		if strings.Contains(line, "TOP 5 HOT NUMBERS:") {
			printSection = true
			lineCount = 0
		}
		if printSection {
			logInfo(line)
			lineCount++
			if lineCount > 7 {
				break
			}
		}
	}

	return nil
}

// Overdue shows most overdue numbers
func (Quick) Overdue() error {
	mg.Deps(Build{}.Dev)

	logInfo("⏰ Most Overdue Numbers:")
	output, err := runAnalyzerWithOutput("--simple")
	if err != nil {
		return err
	}

	lines := strings.Split(output, "\n")
	printSection := false
	lineCount := 0

	for _, line := range lines {
		if strings.Contains(line, "TOP 5 OVERDUE:") {
			printSection = true
			lineCount = 0
		}
		if printSection {
			logInfo(line)
			lineCount++
			if lineCount > 7 {
				break
			}
		}
	}

	return nil
}

// Fun Commands

// CosmicWisdom displays cosmic lottery wisdom
func (Fun) CosmicWisdom() error {
	logInfo("")
	logInfo("✨ ═══════════════════════════════════════════════════════ ✨")
	logInfo("   🌙 The moon influences tides, not lottery numbers! 🌙")
	logInfo("   ☀️  Solar flares can't burn through randomness! ☀️")
	logInfo("   🌟 Every number has exactly 1/48 chance! 🌟")
	logInfo("   🎲 Play for fun, not for cosmic fortune! 🎲")
	logInfo("✨ ═══════════════════════════════════════════════════════ ✨")
	logInfo("")
	return nil
}

// Fortune gets your lottery fortune
func (Fun) Fortune() error {
	mg.Deps(Build{}.Dev)

	logInfo("🔮 Your Lottery Fortune:")
	logInfo("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	output, err := runAnalyzerWithOutput("--simple")
	if err != nil {
		logInfo("The stars are silent today...")
	} else {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "COSMIC PICK:") {
				logInfo(line)
				break
			}
		}
	}

	logInfo("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	logInfo("Remember: Fortune favors the prepared... wallet! 💸")
	return nil
}

// Legacy aliases for common commands (maintaining compatibility)

// FullAnalysis is an alias for Analysis.Full
func FullAnalysis() error {
	return Analysis{}.Full()
}

// Simple is an alias for Analysis.Simple
func Simple() error {
	return Analysis{}.Simple()
}

// Statistical is an alias for Analysis.Statistical
func Statistical() error {
	return Analysis{}.Statistical()
}

// Cosmic is an alias for Analysis.Cosmic
func Cosmic() error {
	return Analysis{}.Cosmic()
}

// ExportJSON is an alias for Export.JSON
func ExportJSON() error {
	return Export{}.JSON()
}

// ExportCSV is an alias for Export.CSV
func ExportCSV() error {
	return Export{}.CSV()
}

// LuckyPicks is an alias for Quick.LuckyPicks
func LuckyPicks() error {
	return Quick{}.LuckyPicks()
}

// HotNumbers is an alias for Quick.HotNumbers
func HotNumbers() error {
	return Quick{}.HotNumbers()
}

// Overdue is an alias for Quick.Overdue
func Overdue() error {
	return Quick{}.Overdue()
}

// CosmicWisdom is an alias for Fun.CosmicWisdom
func CosmicWisdom() error {
	return Fun{}.CosmicWisdom()
}

// Fortune is an alias for Fun.Fortune
func Fortune() error {
	return Fun{}.Fortune()
}
