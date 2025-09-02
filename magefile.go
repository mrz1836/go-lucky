//go:build mage

// Magefile for Go-Lucky Lottery Analyzer
package main

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
	binaryName   = "lottery-analyzer"
	binaryPath   = "./bin/lottery-analyzer"
	coverageFile = "coverage.out"
	testTimeout  = "30s"
)

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
	fmt.Println("🧹 Cleaning up...")

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

	fmt.Println("✅ Cleanup complete")
	return nil
}

// Build Commands

// (Build) Dev builds the development version of the analyzer
func (Build) Dev() error {
	fmt.Println("🔧 Building development version...")

	// Ensure bin directory exists
	if err := os.MkdirAll("bin", 0o750); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Build using magex
	return sh.RunV("magex", "build:dev")
}

// (Build) Default builds the default version
func (Build) Default() error {
	return sh.RunV("magex", "build:default")
}

// (Build) All builds all platforms
func (Build) All() error {
	return sh.RunV("magex", "build:all")
}

// Analysis Commands

// (Analysis) Full runs complete analysis with cosmic correlations (RECOMMENDED)
func (Analysis) Full() error {
	mg.Deps(Build{}.Dev)

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║        🌌 RUNNING FULL COSMIC LOTTERY ANALYSIS 🌌            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println("")

	return runAnalyzer("--cosmic")
}

// (Analysis) Simple runs simple analysis summary
func (Analysis) Simple() error {
	mg.Deps(Build{}.Dev)
	return runAnalyzer("--simple")
}

// (Analysis) Statistical runs detailed statistical analysis
func (Analysis) Statistical() error {
	mg.Deps(Build{}.Dev)
	return runAnalyzer("--statistical")
}

// (Analysis) Cosmic runs cosmic correlation analysis only
func (Analysis) Cosmic() error {
	mg.Deps(Build{}.Dev)
	return runAnalyzer("--cosmic")
}

// Export Commands

// (Export) JSON exports full analysis to JSON file
func (Export) JSON() error {
	mg.Deps(Build{}.Dev)

	fmt.Println("📊 Exporting analysis to JSON...")
	if err := runAnalyzer("--cosmic", "--export-json"); err != nil {
		return err
	}
	fmt.Println("✅ Export complete! Check lottery_analysis_*.json")
	return nil
}

// (Export) CSV exports analysis data to CSV file
func (Export) CSV() error {
	mg.Deps(Build{}.Dev)

	fmt.Println("📊 Exporting analysis to CSV...")
	if err := runAnalyzer("--cosmic", "--export-csv"); err != nil {
		return err
	}
	fmt.Println("✅ Export complete! Check lottery_analysis_*.csv")
	return nil
}

// Quick Analysis Commands

// (Quick) LuckyPicks generates 5 different analysis-based number sets
func (Quick) LuckyPicks() error {
	mg.Deps(Build{}.Dev)

	fmt.Println("🎰 Generating Lucky Picks...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

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
			fmt.Println(line)
			lineCount++
			if lineCount > 10 {
				break
			}
		}
	}

	fmt.Println("")
	fmt.Println("🌌 Cosmic Pick:")
	for _, line := range lines {
		if strings.Contains(line, "COSMIC PICK:") {
			fmt.Println(line)
			break
		}
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

// (Quick) HotNumbers shows current hot numbers
func (Quick) HotNumbers() error {
	mg.Deps(Build{}.Dev)

	fmt.Println("🔥 Current Hot Numbers:")
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
			fmt.Println(line)
			lineCount++
			if lineCount > 7 {
				break
			}
		}
	}

	return nil
}

// (Quick) Overdue shows most overdue numbers
func (Quick) Overdue() error {
	mg.Deps(Build{}.Dev)

	fmt.Println("⏰ Most Overdue Numbers:")
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
			fmt.Println(line)
			lineCount++
			if lineCount > 7 {
				break
			}
		}
	}

	return nil
}

// Fun Commands

// (Fun) CosmicWisdom displays cosmic lottery wisdom
func (Fun) CosmicWisdom() error {
	fmt.Println("")
	fmt.Println("✨ ═══════════════════════════════════════════════════════ ✨")
	fmt.Println("   🌙 The moon influences tides, not lottery numbers! 🌙")
	fmt.Println("   ☀️  Solar flares can't burn through randomness! ☀️")
	fmt.Println("   🌟 Every number has exactly 1/48 chance! 🌟")
	fmt.Println("   🎲 Play for fun, not for cosmic fortune! 🎲")
	fmt.Println("✨ ═══════════════════════════════════════════════════════ ✨")
	fmt.Println("")
	return nil
}

// (Fun) Fortune gets your lottery fortune
func (Fun) Fortune() error {
	mg.Deps(Build{}.Dev)

	fmt.Println("🔮 Your Lottery Fortune:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	output, err := runAnalyzerWithOutput("--simple")
	if err != nil {
		fmt.Println("The stars are silent today...")
	} else {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "COSMIC PICK:") {
				fmt.Println(line)
				break
			}
		}
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Remember: Fortune favors the prepared... wallet! 💸")
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
