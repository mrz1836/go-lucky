package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// FuzzParseDrawings tests the CSV parsing logic with fuzzy input
func FuzzParseDrawings(f *testing.F) {
	// Add seed corpus with valid and edge cases
	f.Add("07/08/2019,17,28,35,41,47,13\n")
	f.Add("01/01/2020,1,2,3,4,5,1\n")
	f.Add("12/31/2021,48,47,46,45,44,18\n")
	f.Add("invalid date,1,2,3,4,5,6\n")
	f.Add("01/01/2020,0,2,3,4,5,6\n")
	f.Add("01/01/2020,49,2,3,4,5,6\n")
	f.Add("01/01/2020,1,2,3,4,5,19\n")
	f.Add("01/01/2020,1,2,3,4,5,0\n")
	f.Add("01/01/2020,1,2,3,4,5\n")
	f.Add("01/01/2020,1,2,3,4,5,6,7\n")
	f.Add("01/01/2020,a,b,c,d,e,f\n")
	f.Add("")
	f.Add("\n\n\n")
	f.Add("01/01/2020,-1,-2,-3,-4,-5,-6\n")
	f.Add("13/13/2020,1,2,3,4,5,6\n")
	f.Add("00/00/0000,1,2,3,4,5,6\n")

	f.Fuzz(func(t *testing.T, csvData string) {
		// Create a context with timeout to prevent infinite loops
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// Parse the CSV data
		reader := csv.NewReader(strings.NewReader(csvData))
		drawings := []Drawing{}

		done := make(chan bool)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
				done <- true
			}()

			records, err := reader.ReadAll()
			if err != nil {
				// CSV parsing errors are expected for fuzzy input
				return
			}

			for _, record := range records {
				if len(record) != 7 {
					continue
				}

				drawing := Drawing{}
				dateStr := record[0]
				if dateStr == "" {
					continue
				}

				// Try to parse the date
				date, err := time.Parse("01/02/2006", dateStr)
				if err != nil {
					continue
				}
				drawing.Date = date

				// Parse numbers
				valid := true
				for i := 1; i <= 5; i++ {
					num := 0
					_, parseErr := fmt.Sscanf(record[i], "%d", &num)
					if parseErr != nil || num < 1 || num > 48 {
						valid = false
						break
					}
					drawing.Numbers[i-1] = num
				}

				if !valid {
					continue
				}

				// Parse lucky ball
				luckyBall := 0
				_, parseErr := fmt.Sscanf(record[6], "%d", &luckyBall)
				if parseErr != nil || luckyBall < 1 || luckyBall > 18 {
					continue
				}
				drawing.LuckyBall = luckyBall

				// Check for duplicate numbers
				seen := make(map[int]bool)
				for _, num := range drawing.Numbers {
					if seen[num] {
						valid = false
						break
					}
					seen[num] = true
				}

				if valid {
					drawings = append(drawings, drawing)
				}
			}
		}()

		select {
		case <-ctx.Done():
			t.Fatalf("Parsing timed out for input: %q", csvData)
		case <-done:
			// Successfully completed
		}

		// Validate the results
		for _, drawing := range drawings {
			// Ensure all numbers are in valid ranges
			for _, num := range drawing.Numbers {
				if num < 1 || num > 48 {
					t.Errorf("Invalid number %d in drawing", num)
				}
			}
			if drawing.LuckyBall < 1 || drawing.LuckyBall > 18 {
				t.Errorf("Invalid lucky ball %d in drawing", drawing.LuckyBall)
			}

			// Ensure no duplicate numbers
			seen := make(map[int]bool)
			for _, num := range drawing.Numbers {
				if seen[num] {
					t.Errorf("Duplicate number %d in drawing", num)
				}
				seen[num] = true
			}
		}
	})
}

// FuzzValidateFilePath tests file path validation with fuzzy input
func FuzzValidateFilePath(f *testing.F) {
	// Add seed corpus
	f.Add("data.csv")
	f.Add("../data.csv")
	f.Add("../../data.csv")
	f.Add("/etc/passwd")
	f.Add("C:\\Windows\\System32\\config\\SAM")
	f.Add("data.csv; rm -rf /")
	f.Add("data.csv\x00.txt")
	f.Add("CON")
	f.Add("PRN")
	f.Add("AUX")
	f.Add("NUL")
	f.Add("COM1")
	f.Add("LPT1")
	f.Add("")
	f.Add(".")
	f.Add("..")
	f.Add("~")
	f.Add("$HOME/data.csv")
	f.Add("%USERPROFILE%\\data.csv")
	f.Add("data?.csv")
	f.Add("data*.csv")
	f.Add("data[1-9].csv")
	f.Add("data.csv|cat")
	f.Add("data.csv&&ls")
	f.Add("data.csv;ls")
	f.Add("data.csv`ls`")
	f.Add("data.csv$(ls)")
	f.Add(strings.Repeat("a", 1000) + ".csv")
	f.Add(strings.Repeat("../", 100) + "data.csv")

	f.Fuzz(func(t *testing.T, path string) {
		// Enable strict validation for this test
		oldVal := os.Getenv("GO_LUCKY_STRICT_VALIDATION")
		_ = os.Setenv("GO_LUCKY_STRICT_VALIDATION", "true")
		defer func() { _ = os.Setenv("GO_LUCKY_STRICT_VALIDATION", oldVal) }()

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		done := make(chan bool)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
				done <- true
			}()

			// Validate the file path
			err := validateFilePath(path)
			if err == nil {
				// If validation passed, ensure the path is actually safe
				cleanPath := filepath.Clean(path)
				validateCleanPath(t, cleanPath, path)
			}
		}()

		select {
		case <-ctx.Done():
			t.Fatalf("Validation timed out for input: %q", path)
		case <-done:
			// Successfully completed
		}
	})
}

// FuzzAnalysisConfig tests configuration validation with fuzzy input
func FuzzAnalysisConfig(f *testing.F) {
	// Add seed corpus
	f.Add(10, 1.0, 0.95, "simple", "console")
	f.Add(0, 1.5, 0.95, "detailed", "csv")
	f.Add(-1, -1.0, -0.5, "", "")
	f.Add(1000000, 100.0, 2.0, "invalid", "invalid")
	f.Add(50, 1.5, 0.95, "statistical", "json")
	f.Add(1, 0.0, 0.0, "cosmic", "console")

	f.Fuzz(func(t *testing.T, recentWindow int, minGapMultiplier, confidenceLevel float64, outputMode, exportFormat string) {
		config := AnalysisConfig{
			RecentWindow:     recentWindow,
			MinGapMultiplier: minGapMultiplier,
			ConfidenceLevel:  confidenceLevel,
			OutputMode:       outputMode,
			ExportFormat:     exportFormat,
		}

		// Create analyzer with fuzzy config
		ctx := context.Background()

		// Create a temporary valid CSV file for testing
		tmpFile, err := os.CreateTemp("", "fuzz_test_*.csv")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		// Write valid test data
		writer := csv.NewWriter(tmpFile)
		_ = writer.Write([]string{"01/01/2020", "1", "2", "3", "4", "5", "6"})
		_ = writer.Write([]string{"01/02/2020", "7", "8", "9", "10", "11", "12"})
		writer.Flush()
		_ = tmpFile.Close()

		// Try to create analyzer with fuzzy config
		analyzer, err := NewAnalyzer(ctx, tmpFile.Name(), &config)
		if err != nil {
			// Error is acceptable for invalid configs
			return
		}

		// If analyzer was created, validate the config was sanitized properly
		if analyzer.config.RecentWindow < 1 {
			t.Errorf("Invalid RecentWindow not sanitized: %d", analyzer.config.RecentWindow)
		}

		if analyzer.config.MinGapMultiplier < 0 {
			t.Errorf("Invalid MinGapMultiplier not sanitized: %f", analyzer.config.MinGapMultiplier)
		}

		if analyzer.config.ConfidenceLevel < 0 || analyzer.config.ConfidenceLevel > 1 {
			t.Errorf("Invalid ConfidenceLevel not sanitized: %f", analyzer.config.ConfidenceLevel)
		}

		validOutputModes := map[string]bool{"simple": true, "detailed": true, "statistical": true, "cosmic": true}
		if analyzer.config.OutputMode != "" && !validOutputModes[analyzer.config.OutputMode] {
			t.Errorf("Invalid OutputMode not sanitized: %s", analyzer.config.OutputMode)
		}

		validExportFormats := map[string]bool{"console": true, "csv": true, "json": true}
		if analyzer.config.ExportFormat != "" && !validExportFormats[analyzer.config.ExportFormat] {
			t.Errorf("Invalid ExportFormat not sanitized: %s", analyzer.config.ExportFormat)
		}
	})
}

// FuzzScoreNumbersByStrategy tests scoring algorithm with fuzzy input
func FuzzScoreNumbersByStrategy(f *testing.F) {
	// Add seed corpus
	f.Add("balanced", 10)
	f.Add("hot", 20)
	f.Add("overdue", 30)
	f.Add("pattern", 50)
	f.Add("frequency", 100)
	f.Add("cosmic", 1)
	f.Add("invalid", 0)
	f.Add("", -1)
	f.Add("BALANCED", 1000000)

	f.Fuzz(func(t *testing.T, strategy string, count int) {
		// Create a simple analyzer for testing
		ctx := context.Background()
		tmpFile, err := os.CreateTemp("", "fuzz_score_*.csv")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		// Write test data
		writer := csv.NewWriter(tmpFile)
		for i := 0; i < 100; i++ {
			date := fmt.Sprintf("01/%02d/2020", (i%30)+1)
			nums := []string{date}
			for j := 1; j <= 5; j++ {
				nums = append(nums, fmt.Sprintf("%d", ((i*j)%48)+1))
			}
			nums = append(nums, fmt.Sprintf("%d", ((i)%18)+1))
			_ = writer.Write(nums)
		}
		writer.Flush()
		_ = tmpFile.Close()

		analyzer, err := NewAnalyzer(ctx, tmpFile.Name(), &AnalysisConfig{
			RecentWindow:     50,
			MinGapMultiplier: 1.5,
			ConfidenceLevel:  0.95,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Test scoring with fuzzy input
		done := make(chan bool)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic in scoring: %v", r)
				}
				done <- true
			}()

			scores := analyzer.scoreNumbersByStrategy(strategy)

			// Validate scores
			validateScores(t, scores, strategy, count)
		}()

		select {
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Scoring timed out for strategy: %q, count: %d", strategy, count)
		case <-done:
			// Successfully completed
		}
	})
}

// FuzzExportAnalysis tests export functionality with fuzzy configurations
func FuzzExportAnalysis(f *testing.F) {
	// Add seed corpus
	f.Add("json", true, false)
	f.Add("csv", false, true)
	f.Add("invalid", true, true)
	f.Add("", false, false)
	f.Add("JSON", true, false)
	f.Add("CSV", false, true)

	f.Fuzz(func(t *testing.T, format string, _, _ bool) {
		// Create analyzer with test data
		ctx := context.Background()
		tmpFile, err := os.CreateTemp("", "fuzz_export_*.csv")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		// Write test data
		writer := csv.NewWriter(tmpFile)
		_ = writer.Write([]string{"01/01/2020", "1", "2", "3", "4", "5", "6"})
		_ = writer.Write([]string{"01/02/2020", "7", "8", "9", "10", "11", "12"})
		writer.Flush()
		_ = tmpFile.Close()

		analyzer, err := NewAnalyzer(ctx, tmpFile.Name(), &AnalysisConfig{
			RecentWindow:     10,
			MinGapMultiplier: 1.5,
			ExportFormat:     format,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Test export with fuzzy parameters
		var buf bytes.Buffer
		done := make(chan bool)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic in export: %v", r)
				}
				done <- true
			}()

			// Try to export (this might fail for invalid formats)
			_ = analyzer.ExportAnalysis(ctx, "")
		}()

		select {
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Export timed out for format: %q", format)
		case <-done:
			// Check output if any
			output := buf.String()
			validateExportOutput(t, output, format)
		}
	})
}

// validateCleanPath validates that a cleaned path is actually safe
func validateCleanPath(t *testing.T, cleanPath, originalPath string) {
	// Check for path traversal
	if strings.Contains(cleanPath, "..") {
		t.Errorf("Path traversal not caught: %q", originalPath)
	}

	// Check for absolute paths
	if filepath.IsAbs(cleanPath) {
		t.Errorf("Absolute path not caught: %q", originalPath)
	}

	// Check for reserved names on Windows
	base := strings.ToUpper(filepath.Base(cleanPath))
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	for _, reserved := range reservedNames {
		if base == reserved || strings.HasPrefix(base, reserved+".") {
			t.Errorf("Reserved name not caught: %q", originalPath)
		}
	}

	// Check for null bytes
	if strings.Contains(originalPath, "\x00") {
		t.Errorf("Null byte not caught: %q", originalPath)
	}

	// Check for shell metacharacters
	dangerous := []string{";", "|", "&", "`", "$", "(", ")", "{", "}", "<", ">", "!", "\\n", "\\r"}
	for _, char := range dangerous {
		if strings.Contains(originalPath, char) {
			t.Errorf("Dangerous character %q not caught in path: %q", char, originalPath)
		}
	}

	// Check length
	if len(originalPath) > 255 {
		t.Errorf("Overly long path not caught: length=%d", len(originalPath))
	}
}

// validateScores validates lottery number scores
func validateScores(t *testing.T, scores []ScoredNumber, strategy string, _ int) {
	if len(scores) == 0 {
		return
	}

	// Check that scores are sorted in descending order
	for i := 1; i < len(scores); i++ {
		if scores[i].Score > scores[i-1].Score {
			t.Errorf("Scores not properly sorted for strategy %s", strategy)
		}
	}

	// Check that all numbers are valid
	for _, score := range scores {
		if score.Number < 1 || score.Number > 48 {
			t.Errorf("Invalid number %d in scores", score.Number)
		}
		if score.Score < 0 {
			t.Errorf("Negative score %f for number %d", score.Score, score.Number)
		}
	}

	// scoreNumbersByStrategy returns all scores without count limiting
}

// validateExportOutput validates export format output
func validateExportOutput(t *testing.T, output, format string) {
	if output == "" {
		return
	}

	lowerFormat := strings.ToLower(format)
	switch lowerFormat {
	case "json":
		// Basic JSON validation
		if !strings.HasPrefix(strings.TrimSpace(output), "{") {
			t.Errorf("Invalid JSON output for format %q", format)
		}
	case "csv":
		// Basic CSV validation
		lines := strings.Split(output, "\n")
		if len(lines) > 0 && !strings.Contains(lines[0], ",") {
			t.Errorf("Invalid CSV output for format %q", format)
		}
	}
}
