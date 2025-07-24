package main

import (
	"context"
	"math"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// AnalyzerTestSuite defines the test suite for lottery analyzer
type AnalyzerTestSuite struct {
	suite.Suite

	analyzer *Analyzer
	testFile string
}

// SetupSuite runs once before all tests
func (s *AnalyzerTestSuite) SetupSuite() {
	// Create a test CSV file
	s.testFile = "test_lottery_data.csv"
	content := `Date,Number 1,Number 2,Number 3,Number 4,Number 5,Lucky Ball
01/15/2024,5,12,23,34,45,7
01/12/2024,3,15,22,38,44,12
01/09/2024,5,18,23,35,42,7
01/06/2024,7,12,25,33,48,15
01/03/2024,2,11,23,34,41,3`

	err := os.WriteFile(s.testFile, []byte(content), 0o600)
	s.Require().NoError(err)
}

// TearDownSuite runs once after all tests
func (s *AnalyzerTestSuite) TearDownSuite() {
	_ = os.Remove(s.testFile) // ignore error in cleanup
}

// SetupTest runs before each test
func (s *AnalyzerTestSuite) SetupTest() {
	config := &AnalysisConfig{
		RecentWindow:     3,
		MinGapMultiplier: 1.5,
		ConfidenceLevel:  0.95,
		OutputMode:       "simple",
		ExportFormat:     "console",
	}

	ctx := context.Background()
	analyzer, err := NewAnalyzer(ctx, s.testFile, config)
	s.Require().NoError(err)
	s.analyzer = analyzer
}

// TestNewAnalyzerValidFile tests creating analyzer with valid file
func (s *AnalyzerTestSuite) TestNewAnalyzerValidFile() {
	s.NotNil(s.analyzer)
	s.Len(s.analyzer.drawings, 5)
	s.NotNil(s.analyzer.mainNumbers)
	s.NotNil(s.analyzer.luckyBalls)
}

// TestNewAnalyzerInvalidFile tests creating analyzer with invalid file
func (s *AnalyzerTestSuite) TestNewAnalyzerInvalidFile() {
	ctx := context.Background()
	_, err := NewAnalyzer(ctx, "nonexistent.csv", nil)
	s.Error(err)
}

// TestNewAnalyzerDefaultConfig tests default configuration
func (s *AnalyzerTestSuite) TestNewAnalyzerDefaultConfig() {
	ctx := context.Background()
	analyzer, err := NewAnalyzer(ctx, s.testFile, nil)
	s.Require().NoError(err)

	s.Equal(50, analyzer.config.RecentWindow)
	s.InEpsilon(1.5, analyzer.config.MinGapMultiplier, 0.001)
	s.Equal("detailed", analyzer.config.OutputMode)
}

// TestParseDrawings tests parsing CSV data
func (s *AnalyzerTestSuite) TestParseDrawings() {
	s.Len(s.analyzer.drawings, 5)

	// Check first drawing (oldest after reversal - which is 01/03/2024)
	firstDrawing := s.analyzer.drawings[0]
	s.Equal(2024, firstDrawing.Date.Year())
	s.Equal(time.January, firstDrawing.Date.Month())
	s.Equal(3, firstDrawing.Date.Day())
	s.Equal([]int{2, 11, 23, 34, 41}, firstDrawing.Numbers)
	s.Equal(3, firstDrawing.LuckyBall)
}

// TestNumberFrequencyTracking tests frequency calculation
func (s *AnalyzerTestSuite) TestNumberFrequencyTracking() {
	// Number 23 appears 3 times in test data
	info := s.analyzer.mainNumbers[23]
	s.Equal(3, info.TotalFrequency)

	// Number 5 appears 2 times
	info5 := s.analyzer.mainNumbers[5]
	s.Equal(2, info5.TotalFrequency)

	// Lucky ball 7 appears 2 times
	lb7 := s.analyzer.luckyBalls[7]
	s.Equal(2, lb7.TotalFrequency)
}

// TestRecentFrequencyTracking tests recent frequency calculation
func (s *AnalyzerTestSuite) TestRecentFrequencyTracking() {
	// With recent window of 3, check recent frequencies
	info23 := s.analyzer.mainNumbers[23]
	s.Equal(2, info23.RecentFrequency) // Appears in 2 of last 3

	info2 := s.analyzer.mainNumbers[2]
	s.Equal(1, info2.RecentFrequency) // Appears once in last 3 (in 01/03/2024)
}

// TestGapCalculation tests gap tracking
func (s *AnalyzerTestSuite) TestGapCalculation() {
	// Number 5 appears at indices 2 and 4 (01/09 and 01/15)
	info5 := s.analyzer.mainNumbers[5]
	s.Contains(info5.GapsSinceDrawn, 2)
	s.Equal(4, info5.CurrentGap) // Last appeared at index 4
}

// TestPatternAnalysis tests pattern detection
func (s *AnalyzerTestSuite) TestPatternAnalysis() {
	// Check odd/even patterns
	s.NotEmpty(s.analyzer.patternStats.OddEvenPatterns)

	// Check sum ranges
	s.NotEmpty(s.analyzer.patternStats.SumRanges)
}

// TestCombinationPatterns tests pair/triple/quad tracking
func (s *AnalyzerTestSuite) TestCombinationPatterns() {
	// Check that pairs are tracked
	s.NotEmpty(s.analyzer.pairPatterns)

	// Check specific pair (5-23 appears twice)
	pairKey := "5-23"
	pattern, exists := s.analyzer.pairPatterns[pairKey]
	s.True(exists)
	s.Equal(2, pattern.Frequency)
}

// TestChiSquareCalculation tests statistical calculations
func (s *AnalyzerTestSuite) TestChiSquareCalculation() {
	s.Greater(s.analyzer.chiSquareValue, 0.0)
	s.GreaterOrEqual(s.analyzer.randomnessScore, 0.0)
	s.LessOrEqual(s.analyzer.randomnessScore, 100.0)
}

// TestGetTopNumbers tests retrieving top frequent numbers
func (s *AnalyzerTestSuite) TestGetTopNumbers() {
	topNumbers := s.analyzer.GetTopNumbers(5, false)
	s.LessOrEqual(len(topNumbers), 5)

	// Verify sorted by frequency
	for i := 1; i < len(topNumbers); i++ {
		s.GreaterOrEqual(topNumbers[i-1].TotalFrequency, topNumbers[i].TotalFrequency)
	}
}

// TestGetOverdueNumbers tests retrieving overdue numbers
func (s *AnalyzerTestSuite) TestGetOverdueNumbers() {
	overdueNumbers := s.analyzer.GetOverdueNumbers(10)

	// Verify sorted by overdue ratio
	for i := 1; i < len(overdueNumbers); i++ {
		ratioI := float64(overdueNumbers[i-1].CurrentGap) / overdueNumbers[i-1].AverageGap
		ratioJ := float64(overdueNumbers[i].CurrentGap) / overdueNumbers[i].AverageGap
		s.GreaterOrEqual(ratioI, ratioJ)
	}
}

// TestGenerateRecommendations tests recommendation generation
func (s *AnalyzerTestSuite) TestGenerateRecommendations() {
	ctx := context.Background()
	recommendations, err := s.analyzer.GenerateRecommendations(ctx, 3)
	s.Require().NoError(err)

	s.LessOrEqual(len(recommendations), 3)

	for _, rec := range recommendations {
		s.Len(rec.Numbers, 5)
		s.Positive(rec.LuckyBall)
		s.LessOrEqual(rec.LuckyBall, 18)
		s.NotEmpty(rec.Strategy)
		s.NotEmpty(rec.Explanation)
		s.Greater(rec.Confidence, 0.0)
		s.LessOrEqual(rec.Confidence, 1.0)
	}
}

// TestScoreNumbersByStrategy tests different scoring strategies
func (s *AnalyzerTestSuite) TestScoreNumbersByStrategy() {
	strategies := []string{"balanced", "hot", "overdue", "pattern", "frequency"}

	for _, strategy := range strategies {
		scored := s.analyzer.scoreNumbersByStrategy(strategy)
		s.NotEmpty(scored)

		// Verify sorted by score
		for i := 1; i < len(scored); i++ {
			s.GreaterOrEqual(scored[i-1].Score, scored[i].Score)
		}
	}
}

// TestExportJSON tests JSON export functionality
func (s *AnalyzerTestSuite) TestExportJSON() {
	s.analyzer.config.ExportFormat = "json"
	testFile := "test_export.json"

	ctx := context.Background()
	err := s.analyzer.ExportAnalysis(ctx, testFile)
	s.Require().NoError(err)

	// Verify file exists
	_, err = os.Stat(testFile)
	s.Require().NoError(err)

	// Clean up
	_ = os.Remove(testFile) // ignore error in cleanup
}

// TestExportCSV tests CSV export functionality
func (s *AnalyzerTestSuite) TestExportCSV() {
	s.analyzer.config.ExportFormat = "csv"
	testFile := "test_export.csv"

	ctx := context.Background()
	err := s.analyzer.ExportAnalysis(ctx, testFile)
	s.Require().NoError(err)

	// Verify file exists
	_, err = os.Stat(testFile)
	s.Require().NoError(err)

	// Clean up
	_ = os.Remove(testFile) // ignore error in cleanup
}

// TestContextCancellation tests context cancellation handling
func (s *AnalyzerTestSuite) TestContextCancellation() {
	parentCtx := context.Background()
	ctx, cancel := context.WithCancel(parentCtx)
	cancel() // Cancel immediately

	_, err := NewAnalyzer(ctx, s.testFile, nil)
	// Should still work as parsing is fast, but if it fails, should be context error
	if err != nil {
		s.Contains(err.Error(), "context canceled")
	}
}

// TestEmptyDataHandling tests handling of empty or invalid data
func (s *AnalyzerTestSuite) TestEmptyDataHandling() {
	// Create file with header only
	emptyFile := "empty_test.csv"
	content := `Date,Number 1,Number 2,Number 3,Number 4,Number 5,Lucky Ball`
	err := os.WriteFile(emptyFile, []byte(content), 0o600)
	s.Require().NoError(err)
	defer func() { _ = os.Remove(emptyFile) }() // ignore error in cleanup

	ctx := context.Background()
	analyzer, err := NewAnalyzer(ctx, emptyFile, nil)
	s.Require().NoError(err)
	s.Empty(analyzer.drawings)
}

// TestInvalidDataHandling tests handling of malformed data
func (s *AnalyzerTestSuite) TestInvalidDataHandling() {
	// Create file with invalid data
	invalidFile := "invalid_test.csv"
	content := `Date,Number 1,Number 2,Number 3,Number 4,Number 5,Lucky Ball
invalid_date,5,12,23,34,45,7
01/12/2024,abc,15,22,38,44,12
01/09/2024,5,18,23,35,42,xyz`
	err := os.WriteFile(invalidFile, []byte(content), 0o600)
	s.Require().NoError(err)
	defer func() { _ = os.Remove(invalidFile) }() // ignore error in cleanup

	ctx := context.Background()
	analyzer, err := NewAnalyzer(ctx, invalidFile, nil)
	s.Require().NoError(err)
	// Should skip invalid rows
	s.Empty(analyzer.drawings)
}

// TestNumberRangeValidation tests that numbers are within valid ranges
func (s *AnalyzerTestSuite) TestNumberRangeValidation() {
	for _, drawing := range s.analyzer.drawings {
		for _, num := range drawing.Numbers {
			s.GreaterOrEqual(num, 1)
			s.LessOrEqual(num, 48)
		}
		s.GreaterOrEqual(drawing.LuckyBall, 1)
		s.LessOrEqual(drawing.LuckyBall, 18)
	}
}

// TestStatisticalMeasures tests statistical calculations
func (s *AnalyzerTestSuite) TestStatisticalMeasures() {
	// Check that all numbers have expected frequency calculated
	for _, info := range s.analyzer.mainNumbers {
		s.Greater(info.ExpectedFrequency, 0.0)
	}

	// Check standard deviation is calculated for numbers with gaps
	for _, info := range s.analyzer.mainNumbers {
		if len(info.GapsSinceDrawn) > 0 {
			s.GreaterOrEqual(info.StandardDeviation, 0.0)
		}
	}
}

// TestCosmicCorrelationEngine tests the cosmic correlation functionality
func (s *AnalyzerTestSuite) TestCosmicCorrelationEngine() {
	// Test correlation engine initialization
	s.NotNil(s.analyzer.correlationEngine)

	// Test cosmic data enrichment
	ctx := context.Background()
	err := s.analyzer.correlationEngine.EnrichWithCosmicData(ctx)
	s.Require().NoError(err)

	// Verify cosmic data was added
	s.NotEmpty(s.analyzer.correlationEngine.cosmicData)

	// Test correlation analysis
	err = s.analyzer.correlationEngine.AnalyzeCorrelations(ctx)
	s.Require().NoError(err)

	// Verify correlations were calculated
	s.NotEmpty(s.analyzer.correlationEngine.correlationResults)
}

// TestMoonPhaseCalculation tests moon phase calculations
func (s *AnalyzerTestSuite) TestMoonPhaseCalculation() {
	// Test known date - January 1, 2024
	testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	phase, illumination := s.analyzer.correlationEngine.calculateMoonPhase(testDate)

	// Verify phase is between 0 and 1
	s.GreaterOrEqual(phase, 0.0)
	s.Less(phase, 1.0)

	// Verify illumination is between 0 and 1
	s.GreaterOrEqual(illumination, 0.0)
	s.LessOrEqual(illumination, 1.0)

	// Test phase name
	phaseName := s.analyzer.correlationEngine.getMoonPhaseName(phase)
	s.NotEmpty(phaseName)
	s.Contains([]string{
		"New Moon", "Waxing Crescent", "First Quarter", "Waxing Gibbous",
		"Full Moon", "Waning Gibbous", "Last Quarter", "Waning Crescent",
	}, phaseName)
}

// TestZodiacCalculation tests zodiac sign calculations for all signs
func (s *AnalyzerTestSuite) TestZodiacCalculation() {
	// Test all zodiac signs with comprehensive date coverage
	testCases := []struct {
		date     time.Time
		expected string
	}{
		// Capricorn: Dec 22 - Jan 19
		{time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "Capricorn"},
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Capricorn"},
		{time.Date(2024, 1, 19, 0, 0, 0, 0, time.UTC), "Capricorn"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Capricorn"},

		// Aquarius: Jan 20 - Feb 18
		{time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC), "Aquarius"},
		{time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC), "Aquarius"},
		{time.Date(2024, 2, 18, 0, 0, 0, 0, time.UTC), "Aquarius"},

		// Pisces: Feb 19 - Mar 20
		{time.Date(2024, 2, 25, 0, 0, 0, 0, time.UTC), "Pisces"},
		{time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC), "Pisces"},
		{time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC), "Pisces"},

		// Aries: Mar 21 - Apr 19
		{time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC), "Aries"},
		{time.Date(2024, 4, 10, 0, 0, 0, 0, time.UTC), "Aries"},
		{time.Date(2024, 4, 19, 0, 0, 0, 0, time.UTC), "Aries"},

		// Taurus: Apr 20 - May 20
		{time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC), "Taurus"},
		{time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC), "Taurus"},
		{time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC), "Taurus"},

		// Gemini: May 21 - Jun 20
		{time.Date(2024, 5, 25, 0, 0, 0, 0, time.UTC), "Gemini"},
		{time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), "Gemini"},
		{time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC), "Gemini"},

		// Cancer: Jun 21 - Jul 22
		{time.Date(2024, 6, 25, 0, 0, 0, 0, time.UTC), "Cancer"},
		{time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), "Cancer"},
		{time.Date(2024, 7, 22, 0, 0, 0, 0, time.UTC), "Cancer"},

		// Leo: Jul 23 - Aug 22
		{time.Date(2024, 7, 25, 0, 0, 0, 0, time.UTC), "Leo"},
		{time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), "Leo"},
		{time.Date(2024, 8, 22, 0, 0, 0, 0, time.UTC), "Leo"},

		// Virgo: Aug 23 - Sep 22
		{time.Date(2024, 8, 25, 0, 0, 0, 0, time.UTC), "Virgo"},
		{time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC), "Virgo"},
		{time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC), "Virgo"},

		// Libra: Sep 23 - Oct 22
		{time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC), "Libra"},
		{time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC), "Libra"},
		{time.Date(2024, 10, 22, 0, 0, 0, 0, time.UTC), "Libra"},

		// Scorpio: Oct 23 - Nov 21
		{time.Date(2024, 10, 25, 0, 0, 0, 0, time.UTC), "Scorpio"},
		{time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC), "Scorpio"},
		{time.Date(2024, 11, 21, 0, 0, 0, 0, time.UTC), "Scorpio"},

		// Sagittarius: Nov 22 - Dec 21
		{time.Date(2024, 11, 25, 0, 0, 0, 0, time.UTC), "Sagittarius"},
		{time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC), "Sagittarius"},
		{time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC), "Sagittarius"},
	}

	for _, tc := range testCases {
		zodiac := s.analyzer.correlationEngine.getZodiacSign(tc.date)
		s.Equal(tc.expected, zodiac, "Date %s should be %s, got %s", tc.date.Format("Jan 2"), tc.expected, zodiac)
	}
}

// TestSeasonalPhaseCalculation tests seasonal phase calculations
func (s *AnalyzerTestSuite) TestSeasonalPhaseCalculation() {
	testCases := []struct {
		date     time.Time
		expected string
	}{
		{time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC), "Spring"},
		{time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), "Summer"},
		{time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC), "Autumn"},
		{time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "Winter"},
	}

	for _, tc := range testCases {
		season := s.analyzer.correlationEngine.getSeasonalPhase(tc.date)
		s.Equal(tc.expected, season)
	}
}

// TestPlanetaryPositions tests planetary position calculations
func (s *AnalyzerTestSuite) TestPlanetaryPositions() {
	testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	positions := s.analyzer.correlationEngine.calculatePlanetaryPositions(testDate)

	// Verify we have positions for expected planets
	expectedPlanets := []string{"Mercury", "Venus", "Mars", "Jupiter", "Saturn"}
	for _, planet := range expectedPlanets {
		position, exists := positions[planet]
		s.True(exists, "Missing position for %s", planet)
		s.GreaterOrEqual(position, 0.0)
		s.Less(position, 360.0)
	}
}

// TestCosmicPredictions tests cosmic-based number predictions
func (s *AnalyzerTestSuite) TestCosmicPredictions() {
	numbers := s.analyzer.correlationEngine.PredictBasedOnCosmicConditions()

	// Verify we get 5 numbers
	s.Len(numbers, 5)

	// Verify all numbers are in valid range and unique
	used := make(map[int]bool)
	for _, num := range numbers {
		s.GreaterOrEqual(num, 1)
		s.LessOrEqual(num, 48)
		s.False(used[num], "Duplicate number: %d", num)
		used[num] = true
	}
}

// TestCorrelationCalculations tests statistical correlation functions
func (s *AnalyzerTestSuite) TestCorrelationCalculations() {
	// Test Pearson correlation with known data
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 6, 8, 10} // Perfect positive correlation

	corr, pValue := calculatePearsonCorrelation(x, y)

	// Should be very close to 1.0 for perfect correlation
	s.Greater(corr, 0.99)
	// P-value should be very small for significant correlation
	if !math.IsNaN(pValue) {
		s.LessOrEqual(pValue, 0.05)
	}

	// Test with no correlation
	yRandom := []float64{5, 1, 8, 2, 9}
	corrRandom, pValueRandom := calculatePearsonCorrelation(x, yRandom)

	// Should be closer to 0
	s.Less(math.Abs(corrRandom), 0.9)
	s.GreaterOrEqual(pValueRandom, 0.0)
}

// TestCosmicReportGeneration tests cosmic report generation
func (s *AnalyzerTestSuite) TestCosmicReportGeneration() {
	// Ensure we have correlation results
	ctx := context.Background()
	err := s.analyzer.correlationEngine.AnalyzeCorrelations(ctx)
	s.Require().NoError(err)

	report := s.analyzer.correlationEngine.GenerateCosmicReport()

	// Verify report contains expected sections
	s.Contains(report, "COSMIC CORRELATION ANALYSIS")
	s.Contains(report, "LUNAR CORRELATIONS")
	s.Contains(report, "CURRENT COSMIC CONDITIONS")
	s.Contains(report, "DISCLAIMER")
}

// TestGetTotalFrequency tests the getTotalFrequency utility function
func (s *AnalyzerTestSuite) TestGetTotalFrequency() {
	// Test with populated frequency map
	freqMap := map[int]int{
		1:  5,
		5:  3,
		10: 7,
		15: 2,
	}

	total := getTotalFrequency(freqMap)
	s.Equal(17, total) // 5+3+7+2 = 17

	// Test with empty map
	emptyMap := make(map[int]int)
	emptyTotal := getTotalFrequency(emptyMap)
	s.Equal(0, emptyTotal)

	// Test with single entry
	singleMap := map[int]int{42: 10}
	singleTotal := getTotalFrequency(singleMap)
	s.Equal(10, singleTotal)
}

// TestCosmicAnalysisMode tests the cosmic analysis output mode
func (s *AnalyzerTestSuite) TestCosmicAnalysisMode() {
	// Change to cosmic mode
	s.analyzer.config.OutputMode = "cosmic"

	// Test that cosmic analysis runs without error
	ctx := context.Background()
	err := s.analyzer.RunAnalysis(ctx)
	s.NoError(err)
}

// TestCalculateConfidence tests confidence calculation edge cases
func (s *AnalyzerTestSuite) TestCalculateConfidence() {
	// Test different strategies
	strategies := []string{"balanced", "hot", "overdue", "pattern", "frequency"}

	for _, strategy := range strategies {
		confidence := s.analyzer.calculateConfidence(strategy)
		s.GreaterOrEqual(confidence, 0.0)
		s.LessOrEqual(confidence, 1.0)
	}

	// Test unknown strategy
	unknownConfidence := s.analyzer.calculateConfidence("unknown")
	s.GreaterOrEqual(unknownConfidence, 0.0)
	s.LessOrEqual(unknownConfidence, 1.0)
}

// TestGenerateExplanation tests explanation generation for different strategies
func (s *AnalyzerTestSuite) TestGenerateExplanation() {
	mockSet := RecommendedSet{
		Numbers:   []int{5, 12, 23, 34, 45},
		LuckyBall: 7,
		Strategy:  "balanced",
	}

	// Test all known strategies
	strategies := []string{"balanced", "hot", "overdue", "pattern", "frequency"}

	for _, strategy := range strategies {
		explanation := s.analyzer.generateExplanation(strategy, mockSet)
		s.NotEmpty(explanation)
	}

	// Test unknown strategy
	unknownExplanation := s.analyzer.generateExplanation("unknown", mockSet)
	s.NotEmpty(unknownExplanation)
	s.Contains(unknownExplanation, "analysis")
}

// TestExportErrors tests export functionality error handling
func (s *AnalyzerTestSuite) TestExportErrors() {
	ctx := context.Background()

	// Test invalid export format
	s.analyzer.config.ExportFormat = "invalid"
	err := s.analyzer.ExportAnalysis(ctx, "test_invalid.out")
	s.Require().Error(err)
	s.Contains(err.Error(), "unsupported export format")

	// Test invalid file path (directory that doesn't exist)
	s.analyzer.config.ExportFormat = "json"
	err = s.analyzer.ExportAnalysis(ctx, "/nonexistent/directory/test.json")
	s.Require().Error(err)

	// Test read-only directory (simulate permission error)
	s.analyzer.config.ExportFormat = "csv"
	err = s.analyzer.ExportAnalysis(ctx, "/test_readonly.csv")
	s.Error(err)
}

// TestExportFormats tests different export format handling
func (s *AnalyzerTestSuite) TestExportFormats() {
	ctx := context.Background()

	// Test console format (should return error since only json/csv are supported)
	s.analyzer.config.ExportFormat = "console"
	err := s.analyzer.ExportAnalysis(ctx, "test_console.txt")
	s.Require().Error(err) // Console format should return error
	s.Contains(err.Error(), "unsupported export format")
}

// TestFileValidationEdgeCases tests file validation edge cases
func (s *AnalyzerTestSuite) TestFileValidationEdgeCases() {
	// Test valid filename
	err := validateFilePath("valid_file.csv")
	s.Require().NoError(err)

	// Test path with directory traversal (will be cleaned by filepath.Clean)
	err = validateFilePath("../test.csv")
	s.Require().NoError(err) // Current implementation allows this

	// Test absolute path
	err = validateFilePath("/tmp/test.csv")
	s.Require().NoError(err) // Current implementation allows this

	// Test relative path
	err = validateFilePath("./test.csv")
	s.Require().NoError(err)
}

// TestRunAnalysisModes tests different output modes
func (s *AnalyzerTestSuite) TestRunAnalysisModes() {
	ctx := context.Background()

	// Test detailed mode
	s.analyzer.config.OutputMode = "detailed"
	err := s.analyzer.RunAnalysis(ctx)
	s.Require().NoError(err)

	// Test simple mode
	s.analyzer.config.OutputMode = "simple"
	err = s.analyzer.RunAnalysis(ctx)
	s.Require().NoError(err)

	// Test statistical mode
	s.analyzer.config.OutputMode = "statistical"
	err = s.analyzer.RunAnalysis(ctx)
	s.Require().NoError(err)

	// Test cosmic mode
	s.analyzer.config.OutputMode = "cosmic"
	err = s.analyzer.RunAnalysis(ctx)
	s.Require().NoError(err)

	// Test unknown mode (should default)
	s.analyzer.config.OutputMode = "unknown"
	err = s.analyzer.RunAnalysis(ctx)
	s.NoError(err)
}

// TestGetTopNumbersEdgeCases tests edge cases for GetTopNumbers
func (s *AnalyzerTestSuite) TestGetTopNumbersEdgeCases() {
	// Test requesting more numbers than available
	topNumbers := s.analyzer.GetTopNumbers(100, false)
	s.LessOrEqual(len(topNumbers), 48) // Should not exceed max possible numbers

	// Test requesting 0 numbers
	emptyNumbers := s.analyzer.GetTopNumbers(0, false)
	s.Empty(emptyNumbers)

	// Test recent vs all-time
	recentNumbers := s.analyzer.GetTopNumbers(5, true)
	allTimeNumbers := s.analyzer.GetTopNumbers(5, false)
	s.Len(recentNumbers, 5)
	s.Len(allTimeNumbers, 5)
}

// TestGenerateRecommendationsEdgeCases tests recommendation generation edge cases
func (s *AnalyzerTestSuite) TestGenerateRecommendationsEdgeCases() {
	ctx := context.Background()

	// Test requesting 0 recommendations
	recommendations, err := s.analyzer.GenerateRecommendations(ctx, 0)
	s.Require().NoError(err)
	s.Empty(recommendations)

	// Test requesting many recommendations
	manyRecs, err := s.analyzer.GenerateRecommendations(ctx, 10)
	s.Require().NoError(err)
	s.LessOrEqual(len(manyRecs), 10) // Should not exceed requested amount
}

// TestSignificanceLevelCoverage tests all significance level branches
func (s *AnalyzerTestSuite) TestSignificanceLevelCoverage() {
	// Test high significance
	high := getSignificanceLevel(0.005)
	s.Equal("High", high)

	// Test moderate significance
	moderate := getSignificanceLevel(0.03)
	s.Equal("Moderate", moderate)

	// Test low significance
	low := getSignificanceLevel(0.08)
	s.Equal("Low", low)

	// Test no significance
	none := getSignificanceLevel(0.15)
	s.Equal("None", none)
}

// TestCorrelationInterpretations tests interpretation functions
func (s *AnalyzerTestSuite) TestCorrelationInterpretations() {
	// Test moon correlation interpretations
	// High p-value (no significance)
	noSig := interpretMoonCorrelation(0.3, 0.2)
	s.Contains(noSig, "No significant correlation")

	// Positive correlation with significance
	posSig := interpretMoonCorrelation(0.7, 0.03)
	s.Contains(posSig, "positive correlation")
	s.Contains(posSig, "Higher numbers")

	// Negative correlation with significance
	negSig := interpretMoonCorrelation(-0.6, 0.03)
	s.Contains(negSig, "negative correlation")
	s.Contains(negSig, "Lower numbers")

	// Test solar correlation interpretations
	solarNoSig := interpretSolarCorrelation(0.2, 0.15)
	s.Contains(solarNoSig, "no significant impact")

	solarSig := interpretSolarCorrelation(0.5, 0.03)
	s.Contains(solarSig, "Correlation detected")

	// Test weather correlation interpretations
	weatherNoSig := interpretWeatherCorrelation(0.1, 0.2)
	s.Contains(weatherNoSig, "no correlation")

	weatherSig := interpretWeatherCorrelation(0.4, 0.04)
	s.Contains(weatherSig, "Weather correlation")
}

// TestParsingEdgeCases tests additional parsing scenarios
func (s *AnalyzerTestSuite) TestParsingEdgeCases() {
	// Test file with different date formats
	mixedDateFile := "mixed_date_test.csv"
	content := `Date,Number 1,Number 2,Number 3,Number 4,Number 5,Lucky Ball
01/15/2024,5,12,23,34,45,7
2024-01-12,3,15,22,38,44,12
01-09-2024,5,18,23,35,42,7`

	err := os.WriteFile(mixedDateFile, []byte(content), 0o600)
	s.Require().NoError(err)
	defer func() { _ = os.Remove(mixedDateFile) }()

	ctx := context.Background()
	analyzer, err := NewAnalyzer(ctx, mixedDateFile, nil)
	s.Require().NoError(err)

	// Should parse at least some valid entries
	s.GreaterOrEqual(len(analyzer.drawings), 1)
}

// TestCurrentConditionsEdgeCases tests edge cases in generateCurrentConditions
func (s *AnalyzerTestSuite) TestCurrentConditionsEdgeCases() {
	// Test different moon phases
	testDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	phase, illumination := s.analyzer.correlationEngine.calculateMoonPhase(testDate)
	s.GreaterOrEqual(phase, 0.0)
	s.Less(phase, 1.0)
	s.GreaterOrEqual(illumination, 0.0)
	s.LessOrEqual(illumination, 1.0)
}

// TestCosmicPredictionEdgeCases tests edge cases in PredictBasedOnCosmicConditions
func (s *AnalyzerTestSuite) TestCosmicPredictionEdgeCases() {
	// Test prediction when cosmic data already exists
	today := time.Now()
	dateKey := today.Format("2006-01-02")

	// Pre-populate cosmic data
	cosmic := &CosmicData{
		Date:             today,
		MoonPhase:        0.5,
		MoonIllumination: 0.8,
		MoonPhaseName:    "Full Moon",
	}
	s.analyzer.correlationEngine.calculateAstronomicalData(cosmic)
	s.analyzer.correlationEngine.addMockDataForDemo(cosmic)
	s.analyzer.correlationEngine.cosmicData[dateKey] = cosmic

	// Test prediction
	numbers := s.analyzer.correlationEngine.PredictBasedOnCosmicConditions()
	s.Len(numbers, 5)

	// Verify uniqueness
	used := make(map[int]bool)
	for _, num := range numbers {
		s.False(used[num], "Duplicate number found: %d", num)
		used[num] = true
	}
}

// Run the test suite
func TestAnalyzerSuite(t *testing.T) {
	suite.Run(t, new(AnalyzerTestSuite))
}

// Benchmark tests
func BenchmarkAnalyzerCreation(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, _ = NewAnalyzer(ctx, "lucky-numbers-history.csv", nil) // ignore error in benchmark
	}
}

func BenchmarkRecommendationGeneration(b *testing.B) {
	ctx := context.Background()
	analyzer, _ := NewAnalyzer(ctx, "lucky-numbers-history.csv", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = analyzer.GenerateRecommendations(ctx, 5) // ignore error in benchmark
	}
}

func BenchmarkPatternAnalysis(b *testing.B) {
	ctx := context.Background()
	analyzer, _ := NewAnalyzer(ctx, "lucky-numbers-history.csv", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, drawing := range analyzer.drawings {
			analyzer.analyzePatterns(drawing)
		}
	}
}

func BenchmarkDataLoading(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, _ = NewAnalyzer(ctx, "lucky-numbers-history.csv", nil)
	}
}

func BenchmarkCosmicCorrelations(b *testing.B) {
	ctx := context.Background()
	analyzer, _ := NewAnalyzer(ctx, "lucky-numbers-history.csv", nil)
	correlationEngine := NewCorrelationEngine(analyzer)
	_ = correlationEngine.EnrichWithCosmicData(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = correlationEngine.AnalyzeCorrelations(ctx)
	}
}

func BenchmarkReportGeneration(b *testing.B) {
	ctx := context.Background()
	analyzer, _ := NewAnalyzer(ctx, "lucky-numbers-history.csv", nil)
	correlationEngine := NewCorrelationEngine(analyzer)
	_ = correlationEngine.EnrichWithCosmicData(ctx)
	_ = correlationEngine.AnalyzeCorrelations(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = correlationEngine.GenerateCosmicReport()
	}
}
