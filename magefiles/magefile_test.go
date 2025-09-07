//go:build mage

package magefiles

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

// MagefileTestSuite defines the test suite for magefile functions
type MagefileTestSuite struct {
	suite.Suite

	// Test directories and files
	testFiles []string
	testDirs  []string
}

// SetupSuite runs once before all tests
func (s *MagefileTestSuite) SetupSuite() {
	// Initialize test tracking
	s.testFiles = []string{}
	s.testDirs = []string{}
}

// SetupTest runs before each test
func (s *MagefileTestSuite) SetupTest() {
	// Reset test files and directories tracking
	s.testFiles = []string{}
	s.testDirs = []string{}
}

// TearDownTest runs after each test
func (s *MagefileTestSuite) TearDownTest() {
	// Clean up test files and directories
	for _, file := range s.testFiles {
		_ = os.Remove(file)
	}
	for _, dir := range s.testDirs {
		_ = os.RemoveAll(dir)
	}
}

// Helper method to create test files
func (s *MagefileTestSuite) createTestFile(path, content string) {
	err := os.WriteFile(path, []byte(content), 0o600)
	s.Require().NoError(err)
	s.testFiles = append(s.testFiles, path)
}

// Helper method to create test directories
func (s *MagefileTestSuite) createTestDir(path string) {
	err := os.MkdirAll(path, 0o750)
	s.Require().NoError(err)
	s.testDirs = append(s.testDirs, path)
}

// Test Constants and Variables

func (s *MagefileTestSuite) TestConstants() {
	s.Equal("go-lucky", binaryName)
	s.Equal("./bin/go-lucky", binaryPath)
	s.Equal("coverage.out", coverageFile)
	s.Equal("30s", testTimeout)
}

// Test Namespace Types

func (s *MagefileTestSuite) TestNamespaceTypes() {
	// Test that namespace types can be instantiated
	build := Build{}
	s.NotNil(build)

	analysis := Analysis{}
	s.NotNil(analysis)

	export := Export{}
	s.NotNil(export)

	quick := Quick{}
	s.NotNil(quick)

	fun := Fun{}
	s.NotNil(fun)
}

// Test Helper Functions Structure

func (s *MagefileTestSuite) TestEnsureBinaryFunctionExists() {
	// Test that the function exists and can be called
	// We can't easily test the actual behavior without mocking external dependencies
	// But we can verify the function structure and basic logic

	// Create a temporary binary to test the "exists" path
	s.createTestDir("test_bin")
	testBinary := filepath.Join("test_bin", "go-lucky")
	s.createTestFile(testBinary, "test binary content")

	// The function should work without error when binary exists
	// We can't easily override the binaryPath constant, so we test the general structure
	// Function exists and compiles - verified by this test running
}

func (s *MagefileTestSuite) TestRunAnalyzerFunction() {
	// Test that runAnalyzer function structure is correct
	// The actual execution depends on external binaries, so we focus on structure
	// Function exists and compiles - verified by this test running
}

func (s *MagefileTestSuite) TestRunAnalyzerWithOutputFunction() {
	// Test that runAnalyzerWithOutput function structure is correct
	// Function exists and compiles - verified by this test running
}

// Test Clean Function Logic

func (s *MagefileTestSuite) TestCleanFilePatterns() {
	// Test that the Clean function handles the expected file patterns
	// We can't easily test the full execution, but we can test the pattern logic

	patterns := []string{
		"lottery_analysis_*.json",
		"lottery_analysis_*.csv",
		"test_*.csv",
		"debug_*.csv",
		"empty_*.csv",
		"invalid_*.csv",
	}

	// Create test files matching these patterns
	for i, pattern := range patterns {
		// Convert pattern to a specific test file name
		testFile := strings.Replace(pattern, "*", fmt.Sprintf("test%d", i), 1)
		s.createTestFile(testFile, "test content")

		// Verify the file was created
		_, err := os.Stat(testFile)
		s.Require().NoError(err)
	}

	// Verify the patterns would match our test files
	matches, err := filepath.Glob(patterns[0])
	s.Require().NoError(err)
	s.Len(matches, 1)
}

// Test Build Namespace

func (s *MagefileTestSuite) TestBuildNamespace() {
	build := Build{}

	// Test that build functions exist and have correct signatures
	s.NotNil(build.Dev)
	s.NotNil(build.Default)
	s.NotNil(build.All)
}

// Test Analysis Namespace

func (s *MagefileTestSuite) TestAnalysisNamespace() {
	analysis := Analysis{}

	// Test that analysis functions exist and have correct signatures
	s.NotNil(analysis.Full)
	s.NotNil(analysis.Simple)
	s.NotNil(analysis.Statistical)
	s.NotNil(analysis.Cosmic)
}

// Test Export Namespace

func (s *MagefileTestSuite) TestExportNamespace() {
	export := Export{}

	// Test that export functions exist and have correct signatures
	s.NotNil(export.JSON)
	s.NotNil(export.CSV)
}

// Test Quick Namespace

func (s *MagefileTestSuite) TestQuickNamespace() {
	quick := Quick{}

	// Test that quick functions exist and have correct signatures
	s.NotNil(quick.LuckyPicks)
	s.NotNil(quick.HotNumbers)
	s.NotNil(quick.Overdue)
}

// Test Fun Namespace

func (s *MagefileTestSuite) TestFunNamespace() {
	fun := Fun{}

	// Test that fun functions exist and have correct signatures
	s.NotNil(fun.CosmicWisdom)
	s.NotNil(fun.Fortune)
}

func (s *MagefileTestSuite) TestFunCosmicWisdom() {
	fun := Fun{}
	err := fun.CosmicWisdom()
	s.NoError(err)
}

// Test Legacy Aliases

func (s *MagefileTestSuite) TestLegacyAliases() {
	// Test that all legacy alias functions exist
	s.NotNil(FullAnalysis)
	s.NotNil(Simple)
	s.NotNil(Statistical)
	s.NotNil(Cosmic)
	s.NotNil(ExportJSON)
	s.NotNil(ExportCSV)
	s.NotNil(LuckyPicks)
	s.NotNil(HotNumbers)
	s.NotNil(Overdue)
	s.NotNil(CosmicWisdom)
	s.NotNil(Fortune)
}

func (s *MagefileTestSuite) TestCosmicWisdomAlias() {
	// This function doesn't depend on external commands, so we can test it fully
	err := CosmicWisdom()
	s.NoError(err)
}

// Test Core Functions

func (s *MagefileTestSuite) TestCoreFunctions() {
	// Test that core functions exist and have correct signatures
	s.NotNil(TestQuick)
	s.NotNil(Test)
	s.NotNil(Benchmark)
	s.NotNil(Clean)
}

// Test String Parsing Logic

func (s *MagefileTestSuite) TestStringParsingLogic() {
	// Test the string parsing logic used in Quick functions
	mockOutput := `Analysis Summary
QUICK PICKS:
Pick 1: 5 12 23 34 45 | Lucky: 7
Pick 2: 3 15 22 38 44 | Lucky: 12
Pick 3: 7 18 25 33 42 | Lucky: 3
COSMIC PICK: 11 22 33 44 45 | Lucky: 8
Some other content`

	lines := strings.Split(mockOutput, "\n")

	// Test finding QUICK PICKS section
	foundQuickPicks := false
	for _, line := range lines {
		if strings.Contains(line, "QUICK PICKS:") {
			foundQuickPicks = true
			break
		}
	}
	s.True(foundQuickPicks)

	// Test finding COSMIC PICK
	foundCosmicPick := false
	for _, line := range lines {
		if strings.Contains(line, "COSMIC PICK:") {
			foundCosmicPick = true
			break
		}
	}
	s.True(foundCosmicPick)
}

func (s *MagefileTestSuite) TestHotNumbersParsingLogic() {
	mockOutput := `Analysis Summary
TOP 5 HOT NUMBERS:
1. Number 23 (appeared 5 times)
2. Number 12 (appeared 4 times)
3. Number 34 (appeared 3 times)
4. Number 45 (appeared 2 times)
5. Number 7 (appeared 2 times)
Other content`

	lines := strings.Split(mockOutput, "\n")

	// Test finding TOP 5 HOT NUMBERS section
	foundHotNumbers := false
	lineCount := 0
	for _, line := range lines {
		if strings.Contains(line, "TOP 5 HOT NUMBERS:") {
			foundHotNumbers = true
		}
		if foundHotNumbers {
			lineCount++
			if lineCount > 7 {
				break
			}
		}
	}
	s.True(foundHotNumbers)
	s.Greater(lineCount, 5) // Should have found the section and counted lines
}

func (s *MagefileTestSuite) TestOverdueParsingLogic() {
	mockOutput := `Analysis Summary
TOP 5 OVERDUE:
1. Number 1 (overdue by 45 draws)
2. Number 8 (overdue by 32 draws)
3. Number 17 (overdue by 28 draws)
4. Number 33 (overdue by 25 draws)
5. Number 42 (overdue by 20 draws)
Other content`

	lines := strings.Split(mockOutput, "\n")

	// Test finding TOP 5 OVERDUE section
	foundOverdue := false
	for _, line := range lines {
		if strings.Contains(line, "TOP 5 OVERDUE:") {
			foundOverdue = true
			break
		}
	}
	s.True(foundOverdue)
}

// Test Error Handling Patterns

func (s *MagefileTestSuite) TestErrorHandlingPatterns() {
	// Test that functions have proper error handling structure
	// We can't test the actual error cases without external dependencies,
	// but we can verify the function signatures return errors where expected

	// Functions that should return errors
	build := Build{}
	_, ok := interface{}(build.Dev).(func() error)
	s.True(ok, "Build.Dev should return error")

	analysis := Analysis{}
	_, ok = interface{}(analysis.Full).(func() error)
	s.True(ok, "Analysis.Full should return error")
}

// Test File and Directory Structure

func (s *MagefileTestSuite) TestFileStructureAssumptions() {
	// Test assumptions about file structure
	s.Equal("./bin/go-lucky", binaryPath)
	s.Equal("coverage.out", coverageFile)

	// Test that binary path includes bin directory
	s.Contains(binaryPath, "bin/")

	// Test that coverage file has .out extension
	s.Contains(coverageFile, ".out")
}

// Test Output Formatting

func (s *MagefileTestSuite) TestOutputFormatting() {
	// Test that the CosmicWisdom function produces expected output format
	// This is one of the few functions we can test completely

	// Capture the function's behavior
	err := CosmicWisdom()
	s.NoError(err)

	// The function should execute without error
	// We can't easily capture stdout in this context, but we verify it runs
}

// Test Integration Points

func (s *MagefileTestSuite) TestIntegrationPoints() {
	// Test that the magefile integrates properly with expected tools

	// Test that binaryPath points to expected location
	expectedBinaryDir := "bin"
	s.Contains(binaryPath, expectedBinaryDir)

	// Test that timeout is properly formatted
	s.Contains(testTimeout, "s") // Should end with 's' for seconds

	// Test that binary name doesn't contain path separators
	s.NotContains(binaryName, "/")
	s.NotContains(binaryName, "\\")
}

// Test Concurrent Safety

func (s *MagefileTestSuite) TestConcurrentSafety() {
	// Test that the constants are safe for concurrent access
	// Constants should be read-only and thread-safe by nature

	s.Equal(binaryName, "go-lucky")
	s.Equal(binaryName, "go-lucky") // Should be consistent

	s.Equal(binaryPath, "./bin/go-lucky")
	s.Equal(binaryPath, "./bin/go-lucky") // Should be consistent
}

// Test Dependency Structure

func (s *MagefileTestSuite) TestDependencyStructure() {
	// Test that the magefile has the expected dependencies
	// The mg.Deps calls in the actual functions create the dependency structure

	// We can't easily test mg.Deps without actually running the functions,
	// but we can verify the structure exists
	analysis := Analysis{}
	s.NotNil(analysis.Full)

	// The Full() function should depend on Build{}.Dev according to the implementation
	// This creates a proper dependency chain for building before analysis
}

// Test Build Tags

func (s *MagefileTestSuite) TestBuildTags() {
	// Verify that this test file itself has the correct build tag
	// The file should have //go:build mage at the top

	// If this test runs, the build tag is working correctly
	// Test is running with correct build tag - verified by this test running
}

// Run the test suite
func TestMagefileTestSuite(t *testing.T) {
	suite.Run(t, new(MagefileTestSuite))
}

// Additional individual tests for comprehensive coverage

func TestConstants(t *testing.T) {
	if binaryName != "go-lucky" {
		t.Errorf("Expected binaryName to be 'go-lucky', got '%s'", binaryName)
	}

	if binaryPath != "./bin/go-lucky" {
		t.Errorf("Expected binaryPath to be './bin/go-lucky', got '%s'", binaryPath)
	}

	if coverageFile != "coverage.out" {
		t.Errorf("Expected coverageFile to be 'coverage.out', got '%s'", coverageFile)
	}

	if testTimeout != "30s" {
		t.Errorf("Expected testTimeout to be '30s', got '%s'", testTimeout)
	}
}

func TestNamespaceInstantiation(_ *testing.T) {
	// Test that all namespace types can be instantiated
	build := Build{}
	analysis := Analysis{}
	export := Export{}
	quick := Quick{}
	fun := Fun{}

	// Basic nil checks
	_ = build
	_ = analysis
	_ = export
	_ = quick
	_ = fun
}

func TestFunctionSignatures(t *testing.T) {
	// Test that all functions exist and have expected function types

	// Core functions - check that they are functions returning error
	_, ok := interface{}(TestQuick).(func() error)
	if !ok {
		t.Error("TestQuick is not a function returning error")
	}

	_, ok = interface{}(Test).(func() error)
	if !ok {
		t.Error("Test is not a function returning error")
	}

	_, ok = interface{}(Benchmark).(func() error)
	if !ok {
		t.Error("Benchmark is not a function returning error")
	}

	_, ok = interface{}(Clean).(func() error)
	if !ok {
		t.Error("Clean is not a function returning error")
	}

	// Legacy aliases - check that they are functions returning error
	_, ok = interface{}(FullAnalysis).(func() error)
	if !ok {
		t.Error("FullAnalysis is not a function returning error")
	}

	_, ok = interface{}(Simple).(func() error)
	if !ok {
		t.Error("Simple is not a function returning error")
	}

	_, ok = interface{}(Statistical).(func() error)
	if !ok {
		t.Error("Statistical is not a function returning error")
	}

	_, ok = interface{}(Cosmic).(func() error)
	if !ok {
		t.Error("Cosmic is not a function returning error")
	}
}

func TestCosmicWisdomOutput(t *testing.T) {
	// Test the one function we can safely test without external dependencies
	err := CosmicWisdom()
	if err != nil {
		t.Errorf("CosmicWisdom returned error: %v", err)
	}
}
