package main

import (
	"context"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	// Export formats
	exportFormatJSON = "json"
	exportFormatCSV  = "csv"
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
	s.analyzer.config.ExportFormat = exportFormatJSON
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
	s.analyzer.config.ExportFormat = exportFormatCSV
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
	s.analyzer.config.ExportFormat = exportFormatJSON
	err = s.analyzer.ExportAnalysis(ctx, "/nonexistent/directory/test.json")
	s.Require().Error(err)

	// Test read-only directory (simulate permission error)
	s.analyzer.config.ExportFormat = exportFormatCSV
	err = s.analyzer.ExportAnalysis(ctx, "/test_readonly.csv")
	s.Error(err)
}

// TestExportComprehensiveErrors tests comprehensive export error scenarios
func (s *AnalyzerTestSuite) TestExportComprehensiveErrors() {
	ctx := context.Background()

	// Store original environment for file validation
	originalEnv := os.Getenv("GO_LUCKY_STRICT_VALIDATION")
	defer func() {
		_ = os.Setenv("GO_LUCKY_STRICT_VALIDATION", originalEnv)
	}()

	// Enable strict validation to test file path validation
	_ = os.Setenv("GO_LUCKY_STRICT_VALIDATION", "true")

	// Test JSON export with invalid file path
	s.analyzer.config.ExportFormat = exportFormatJSON
	err := s.analyzer.ExportAnalysis(ctx, "")
	s.Require().Error(err)
	s.Contains(err.Error(), "invalid file path")

	// Test CSV export with invalid file path
	s.analyzer.config.ExportFormat = exportFormatCSV
	err = s.analyzer.ExportAnalysis(ctx, "../invalid.csv")
	s.Require().Error(err)
	s.Contains(err.Error(), "invalid file path")

	// Test JSON export with dangerous filename
	s.analyzer.config.ExportFormat = exportFormatJSON
	err = s.analyzer.ExportAnalysis(ctx, "test;rm.json")
	s.Require().Error(err)
	s.Contains(err.Error(), "invalid file path")

	// Test CSV export with dangerous filename
	s.analyzer.config.ExportFormat = exportFormatCSV
	err = s.analyzer.ExportAnalysis(ctx, "test|rm.csv")
	s.Require().Error(err)
	s.Contains(err.Error(), "invalid file path")

	// Test with null byte in filename
	s.analyzer.config.ExportFormat = exportFormatJSON
	err = s.analyzer.ExportAnalysis(ctx, "test\x00.json")
	s.Require().Error(err)
	s.Contains(err.Error(), "invalid file path")

	// Test with Windows reserved name
	s.analyzer.config.ExportFormat = exportFormatCSV
	err = s.analyzer.ExportAnalysis(ctx, "CON.csv")
	s.Require().Error(err)
	s.Contains(err.Error(), "invalid file path")

	// Test exportJSON and exportCSV directly with permission errors
	// This would test file creation failures
	s.analyzer.config.ExportFormat = exportFormatJSON
	err = s.analyzer.ExportAnalysis(ctx, "/root/test.json") // Should fail with permission error
	s.Require().Error(err)

	s.analyzer.config.ExportFormat = exportFormatCSV
	err = s.analyzer.ExportAnalysis(ctx, "/root/test.csv") // Should fail with permission error
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

// TestFileValidationErrors tests comprehensive file validation error scenarios
func (s *AnalyzerTestSuite) TestFileValidationErrors() {
	// Test empty path
	err := validateFilePath("")
	s.Require().Error(err)
	s.Equal(ErrInvalidFilePath, err)

	// Store original environment
	originalEnv := os.Getenv("GO_LUCKY_STRICT_VALIDATION")
	defer func() {
		_ = os.Setenv("GO_LUCKY_STRICT_VALIDATION", originalEnv)
	}()

	// Enable strict validation for these tests
	_ = os.Setenv("GO_LUCKY_STRICT_VALIDATION", "true")

	testCases := []struct {
		name     string
		filename string
		hasError bool
	}{
		{"Empty path", "", true},
		{"Valid filename", "test.csv", false},
		{"Too long filename", strings.Repeat("a", 256), true},
		{"Null byte", "test\x00.csv", true},
		{"Semicolon", "test;rm.csv", true},
		{"Pipe", "test|rm.csv", true},
		{"Ampersand", "test&rm.csv", true},
		{"Backtick", "test`rm.csv", true},
		{"Dollar sign", "test$rm.csv", true},
		{"Parentheses", "test(rm).csv", true},
		{"Curly braces", "test{rm}.csv", true},
		{"Angle brackets", "test<rm>.csv", true},
		{"Exclamation", "test!rm.csv", true},
		{"Newline", "test\\nrm.csv", true},
		{"Carriage return", "test\\rrm.csv", true},
		{"Path traversal", "../../../etc/passwd", true},
		{"Absolute path", "/tmp/test.csv", true},
		{"Windows reserved CON", "CON", true},
		{"Windows reserved PRN", "PRN.txt", true},
		{"Windows reserved AUX", "aux.csv", true},
		{"Windows reserved COM1", "COM1.dat", true},
		{"Windows reserved LPT1", "LPT1.log", true},
		{"Valid with subdir", "subdir/test.csv", false},
	}

	for _, tc := range testCases {
		validateErr := validateFilePath(tc.filename)
		if tc.hasError {
			s.Require().Error(validateErr, "Test case: %s should return error", tc.name)
			s.Equal(ErrInvalidFilePath, validateErr, "Test case: %s", tc.name)
		} else {
			s.Require().NoError(validateErr, "Test case: %s should not return error", tc.name)
		}
	}

	// Test validateFilePathStrict directly
	err = validateFilePathStrict("valid.csv", "valid.csv")
	s.Require().NoError(err)

	err = validateFilePathStrict(strings.Repeat("x", 256), strings.Repeat("x", 256))
	s.Require().Error(err)
	s.Equal(ErrInvalidFilePath, err)
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
	dateKey := today.Format(dateFormatISO)

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

// TestNetworkFailureHandling tests handling of network failures in cosmic data fetching
func (s *AnalyzerTestSuite) TestNetworkFailureHandling() {
	// Test that fetchMoonPhaseData handles errors gracefully
	ctx := context.Background()

	// Test with normal context
	err := s.analyzer.correlationEngine.fetchMoonPhaseData(ctx, 2024)
	s.Require().NoError(err) // Current implementation should not fail

	// Test with canceled context
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()
	err = s.analyzer.correlationEngine.fetchMoonPhaseData(canceledCtx, 2024)
	s.Require().NoError(err) // Current implementation ignores context cancellation

	// Test with timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
	defer cancel()
	time.Sleep(2 * time.Nanosecond) // Ensure timeout
	err = s.analyzer.correlationEngine.fetchMoonPhaseData(timeoutCtx, 2024)
	s.Require().NoError(err) // Current implementation should still work as it doesn't make network requests

	// Test error handling in EnrichWithCosmicData when fetchMoonPhaseData fails
	// This tests the warning path when fetch fails
	originalStderr := os.Stderr
	defer func() { os.Stderr = originalStderr }()

	// Capture stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Create a correlation engine that would fail on network requests
	testEngine := NewCorrelationEngine(s.analyzer)
	err = testEngine.EnrichWithCosmicData(ctx)
	s.Require().NoError(err) // Should complete despite warnings

	_ = w.Close()
	os.Stderr = originalStderr

	// Read captured output (this tests the warning path)
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	if n > 0 {
		output := string(buf[:n])
		// Current implementation shouldn't produce warnings since it doesn't fail
		s.NotContains(output, "Warning: Could not fetch moon data")
	}
}

// TestHTTPClientConfiguration tests HTTP client setup for cosmic data fetching
func (s *AnalyzerTestSuite) TestHTTPClientConfiguration() {
	engine := NewCorrelationEngine(s.analyzer)

	// Verify HTTP client is configured
	s.NotNil(engine.client)
	s.Equal(30*time.Second, engine.client.Timeout)
}

// TestCosmicDataEnrichmentErrors tests error scenarios during cosmic data enrichment
func (s *AnalyzerTestSuite) TestCosmicDataEnrichmentErrors() {
	ctx := context.Background()

	// Test with empty analyzer
	emptyAnalyzer := &Analyzer{
		drawings: []Drawing{},
	}
	emptyEngine := NewCorrelationEngine(emptyAnalyzer)
	err := emptyEngine.EnrichWithCosmicData(ctx)
	s.Require().NoError(err) // Should handle empty data gracefully

	// Test with analyzer containing drawings from multiple years
	multiYearAnalyzer := &Analyzer{
		drawings: []Drawing{
			{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			{Date: time.Date(2021, 6, 15, 0, 0, 0, 0, time.UTC)},
			{Date: time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)},
			{Date: time.Date(2023, 3, 20, 0, 0, 0, 0, time.UTC)},
			{Date: time.Date(2024, 8, 10, 0, 0, 0, 0, time.UTC)},
		},
	}
	multiYearEngine := NewCorrelationEngine(multiYearAnalyzer)
	err = multiYearEngine.EnrichWithCosmicData(ctx)
	s.Require().NoError(err)

	// Verify cosmic data was created for all drawing dates (may include additional calendar dates)
	s.GreaterOrEqual(len(multiYearEngine.cosmicData), 5)

	// Verify each date has cosmic data
	for _, drawing := range multiYearAnalyzer.drawings {
		dateKey := drawing.Date.Format(dateFormatISO)
		cosmic, exists := multiYearEngine.cosmicData[dateKey]
		s.True(exists, "Missing cosmic data for date %s", dateKey)
		s.NotNil(cosmic)
		s.Equal(drawing.Date, cosmic.Date)
	}
}

// TestMockDataForDemo tests the mock data generation for demonstration
func (s *AnalyzerTestSuite) TestMockDataForDemo() {
	cosmic := &CosmicData{
		Date: time.Now(),
	}

	// Call calculateAstronomicalData first to set up planetary positions
	s.analyzer.correlationEngine.calculateAstronomicalData(cosmic)

	// Test that addMockDataForDemo populates all fields
	s.analyzer.correlationEngine.addMockDataForDemo(cosmic)

	// Verify solar data is populated
	s.NotNil(cosmic.SolarActivity)
	s.Greater(cosmic.SolarActivity.SolarWindSpeed, 0.0)
	s.Greater(cosmic.SolarActivity.F107Index, 0.0)

	// Verify weather data is populated
	s.NotNil(cosmic.WeatherData)
	s.Greater(cosmic.WeatherData.Temperature, -50.0) // Basic range check

	// Verify planetary positions are populated
	s.NotEmpty(cosmic.PlanetaryPositions)
	s.Contains(cosmic.PlanetaryPositions, "Mercury")
	s.Contains(cosmic.PlanetaryPositions, "Venus")
	s.Contains(cosmic.PlanetaryPositions, "Mars")
	s.Contains(cosmic.PlanetaryPositions, "Jupiter")
	s.Contains(cosmic.PlanetaryPositions, "Saturn")
}

// TestCosmicCorrelationEdgeCases tests edge cases in cosmic correlation analysis
func (s *AnalyzerTestSuite) TestCosmicCorrelationEdgeCases() {
	ctx := context.Background()

	// Test with analyzer that has no drawings
	emptyAnalyzer := &Analyzer{
		drawings:     []Drawing{},
		mainNumbers:  make(map[int]*NumberInfo),
		luckyBalls:   make(map[int]*NumberInfo),
		pairPatterns: make(map[string]*CombinationPattern),
		patternStats: &PatternStats{},
	}
	emptyEngine := NewCorrelationEngine(emptyAnalyzer)

	// Test EnrichWithCosmicData with empty data
	err := emptyEngine.EnrichWithCosmicData(ctx)
	s.Require().NoError(err) // Should handle empty data gracefully

	// Test AnalyzeCorrelations with empty data
	err = emptyEngine.AnalyzeCorrelations(ctx)
	s.Require().NoError(err) // Should handle empty data gracefully

	// Test cosmic prediction with no historical data
	numbers := emptyEngine.PredictBasedOnCosmicConditions()
	s.Len(numbers, 5) // Should still generate 5 numbers

	// Test that all numbers are in valid range
	for _, num := range numbers {
		s.GreaterOrEqual(num, 1)
		s.LessOrEqual(num, 48)
	}

	// Test with minimal data
	minimalAnalyzer := &Analyzer{
		drawings: []Drawing{
			{
				Date:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Numbers:   []int{1, 2, 3, 4, 5},
				LuckyBall: 1,
			},
		},
		mainNumbers:  make(map[int]*NumberInfo),
		luckyBalls:   make(map[int]*NumberInfo),
		pairPatterns: make(map[string]*CombinationPattern),
		patternStats: &PatternStats{
			OddEvenPatterns:    make(map[string]int),
			SumRanges:          make(map[int]int),
			ConsecutiveCount:   0,
			DecadeDistribution: make(map[int]int),
		},
	}

	// Initialize number info for minimal analyzer
	for i := 1; i <= 48; i++ {
		minimalAnalyzer.mainNumbers[i] = &NumberInfo{
			Number:            i,
			TotalFrequency:    0,
			RecentFrequency:   0,
			GapsSinceDrawn:    []int{},
			CurrentGap:        1,
			AverageGap:        1.0,
			ExpectedFrequency: 1.0,
		}
	}
	for i := 1; i <= 18; i++ {
		minimalAnalyzer.luckyBalls[i] = &NumberInfo{
			Number:            i,
			TotalFrequency:    0,
			RecentFrequency:   0,
			GapsSinceDrawn:    []int{},
			CurrentGap:        1,
			AverageGap:        1.0,
			ExpectedFrequency: 1.0,
		}
	}

	minimalEngine := NewCorrelationEngine(minimalAnalyzer)

	// Test with minimal data
	err = minimalEngine.EnrichWithCosmicData(ctx)
	s.Require().NoError(err)

	err = minimalEngine.AnalyzeCorrelations(ctx)
	s.Require().NoError(err)

	// Test report generation with minimal data
	report := minimalEngine.GenerateCosmicReport()
	s.NotEmpty(report)
	s.Contains(report, "COSMIC CORRELATION ANALYSIS")
}

// TestCosmicWeatherCorrelations tests weather data correlation edge cases
func (s *AnalyzerTestSuite) TestCosmicWeatherCorrelations() {
	ctx := context.Background()

	// Test with diverse weather conditions
	// Create an analyzer with specific drawings to test weather correlations
	weatherAnalyzer := &Analyzer{
		drawings:     make([]Drawing, 0),
		mainNumbers:  make(map[int]*NumberInfo),
		luckyBalls:   make(map[int]*NumberInfo),
		pairPatterns: make(map[string]*CombinationPattern),
		patternStats: &PatternStats{
			OddEvenPatterns:    make(map[string]int),
			SumRanges:          make(map[int]int),
			ConsecutiveCount:   0,
			DecadeDistribution: make(map[int]int),
		},
	}

	// Add drawings with specific patterns
	dates := []time.Time{
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
	}

	for i, date := range dates {
		drawing := Drawing{
			Date:      date,
			Numbers:   []int{i*5 + 1, i*5 + 2, i*5 + 3, i*5 + 4, i*5 + 5},
			LuckyBall: i + 1,
		}
		weatherAnalyzer.drawings = append(weatherAnalyzer.drawings, drawing)
	}

	// Initialize number info
	for i := 1; i <= 48; i++ {
		weatherAnalyzer.mainNumbers[i] = &NumberInfo{
			Number:            i,
			TotalFrequency:    1,
			RecentFrequency:   1,
			GapsSinceDrawn:    []int{1, 2, 3},
			CurrentGap:        1,
			AverageGap:        2.0,
			ExpectedFrequency: 2.0,
			StandardDeviation: 1.0,
		}
	}
	for i := 1; i <= 18; i++ {
		weatherAnalyzer.luckyBalls[i] = &NumberInfo{
			Number:            i,
			TotalFrequency:    1,
			RecentFrequency:   1,
			GapsSinceDrawn:    []int{1, 2},
			CurrentGap:        1,
			AverageGap:        1.5,
			ExpectedFrequency: 1.5,
			StandardDeviation: 0.5,
		}
	}

	weatherEngine := NewCorrelationEngine(weatherAnalyzer)

	// Test enrichment and analysis
	err := weatherEngine.EnrichWithCosmicData(ctx)
	s.Require().NoError(err)

	err = weatherEngine.AnalyzeCorrelations(ctx)
	s.Require().NoError(err)

	// Verify correlations were calculated
	s.NotEmpty(weatherEngine.correlationResults)

	// Test report generation
	report := weatherEngine.GenerateCosmicReport()
	s.NotEmpty(report)
	s.Contains(report, "WEATHER CORRELATIONS")
}

// TestAstronomicalDataCalculation tests astronomical data calculation edge cases
func (s *AnalyzerTestSuite) TestAstronomicalDataCalculation() {
	// Test calculateAstronomicalData with various dates
	testDates := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),   // New Year
		time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),  // Summer Solstice
		time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC), // Winter Solstice
		time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),  // Spring Equinox
		time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),  // Fall Equinox
	}

	for _, date := range testDates {
		cosmic := &CosmicData{Date: date}
		s.analyzer.correlationEngine.calculateAstronomicalData(cosmic)

		// Verify all astronomical fields are populated
		s.NotEmpty(cosmic.ZodiacSign)
		s.NotEmpty(cosmic.DayOfWeek)
		s.NotEmpty(cosmic.SeasonalPhase)
		s.NotEmpty(cosmic.PlanetaryPositions)
		s.GreaterOrEqual(cosmic.GeomagneticIndex, 0.0)

		// Verify planetary positions are within valid range (0-360 degrees)
		for planet, position := range cosmic.PlanetaryPositions {
			s.GreaterOrEqual(position, 0.0, "Planet %s position should be >= 0", planet)
			s.Less(position, 360.0, "Planet %s position should be < 360", planet)
		}
	}
}

// TestCosmicReportGenerationEdgeCases tests edge cases in report generation
func (s *AnalyzerTestSuite) TestCosmicReportGenerationEdgeCases() {
	// Test report generation with no correlation results
	emptyEngine := NewCorrelationEngine(&Analyzer{
		drawings:     []Drawing{},
		mainNumbers:  make(map[int]*NumberInfo),
		luckyBalls:   make(map[int]*NumberInfo),
		pairPatterns: make(map[string]*CombinationPattern),
		patternStats: &PatternStats{},
	})

	// Generate report without any data
	report := emptyEngine.GenerateCosmicReport()
	s.NotEmpty(report)
	s.Contains(report, "COSMIC CORRELATION ANALYSIS")
	s.Contains(report, "DISCLAIMER")

	// Test with some correlation results
	emptyEngine.correlationResults = []CorrelationResult{
		{
			Factor:       "Moon Phase",
			SubFactor:    "Full Moon",
			Correlation:  0.75,
			PValue:       0.01,
			SampleSize:   100,
			Significance: "High",
		},
		{
			Factor:       "Solar Activity",
			SubFactor:    "Solar Wind Speed",
			Correlation:  -0.45,
			PValue:       0.03,
			SampleSize:   100,
			Significance: "Moderate",
		},
		{
			Factor:       "Weather",
			SubFactor:    "Temperature",
			Correlation:  0.25,
			PValue:       0.15,
			SampleSize:   100,
			Significance: "None",
		},
	}

	// Generate report with mock correlation results
	detailedReport := emptyEngine.GenerateCosmicReport()
	s.NotEmpty(detailedReport)
	s.Contains(detailedReport, "Moon Phase")
	s.Contains(strings.ToUpper(detailedReport), strings.ToUpper("Solar Activity"))
	s.Contains(strings.ToUpper(detailedReport), strings.ToUpper("Weather"))
}

// TestStatisticalAnalysisEdgeCases tests edge cases in statistical analysis
func (s *AnalyzerTestSuite) TestStatisticalAnalysisEdgeCases() {
	// Test gap analysis with no gaps
	analyzer := &Analyzer{
		drawings: []Drawing{
			{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Numbers: []int{1, 2, 3, 4, 5}, LuckyBall: 1},
		},
		mainNumbers:  make(map[int]*NumberInfo),
		luckyBalls:   make(map[int]*NumberInfo),
		pairPatterns: make(map[string]*CombinationPattern),
		patternStats: &PatternStats{
			OddEvenPatterns:    make(map[string]int),
			SumRanges:          make(map[int]int),
			ConsecutiveCount:   0,
			DecadeDistribution: make(map[int]int),
		},
	}

	// Test with numbers that have no gaps
	for i := 1; i <= 48; i++ {
		analyzer.mainNumbers[i] = &NumberInfo{
			Number:            i,
			TotalFrequency:    0,
			RecentFrequency:   0,
			GapsSinceDrawn:    []int{}, // No gaps
			CurrentGap:        0,
			AverageGap:        0.0,
			ExpectedFrequency: 0.0,
			StandardDeviation: 0.0,
		}
	}

	// Test statistical analysis output with edge case data
	ctx := context.Background()
	analyzer.config = &AnalysisConfig{OutputMode: "statistical"}
	err := analyzer.printStatisticalAnalysis(ctx)
	s.Require().NoError(err)

	// Test with numbers that have extreme gaps
	analyzer.mainNumbers[1].GapsSinceDrawn = []int{1, 100, 200, 500, 1000}
	analyzer.mainNumbers[1].CurrentGap = 999999
	analyzer.mainNumbers[1].AverageGap = 360.2
	analyzer.mainNumbers[1].StandardDeviation = 415.8

	err = analyzer.printStatisticalAnalysis(ctx)
	s.NoError(err)
}

// TestChiSquareCalculationEdgeCases tests chi-square calculation edge cases
func (s *AnalyzerTestSuite) TestChiSquareCalculationEdgeCases() {
	// Test calculateChiSquare with extreme values
	testCases := []struct {
		name           string
		frequencies    map[int]int
		expectedResult bool
	}{
		{
			name: "All zeros",
			frequencies: map[int]int{
				1: 0, 2: 0, 3: 0, 4: 0, 5: 0,
			},
			expectedResult: true, // Should handle gracefully
		},
		{
			name: "Single high frequency",
			frequencies: map[int]int{
				1: 1000, 2: 0, 3: 0, 4: 0, 5: 0,
			},
			expectedResult: true,
		},
		{
			name: "Equal frequencies",
			frequencies: map[int]int{
				1: 10, 2: 10, 3: 10, 4: 10, 5: 10,
			},
			expectedResult: true,
		},
		{
			name:           "Empty frequencies",
			frequencies:    map[int]int{},
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		// Create test analyzer with specific frequencies
		testAnalyzer := &Analyzer{
			drawings:     []Drawing{},
			mainNumbers:  make(map[int]*NumberInfo),
			luckyBalls:   make(map[int]*NumberInfo),
			pairPatterns: make(map[string]*CombinationPattern),
			patternStats: &PatternStats{},
		}

		// Set up frequencies
		for num, freq := range tc.frequencies {
			testAnalyzer.mainNumbers[num] = &NumberInfo{
				Number:         num,
				TotalFrequency: freq,
			}
		}

		// Test chi-square calculation
		testAnalyzer.calculateChiSquare()

		// Verify chi-square value is valid (not NaN or negative)
		s.False(math.IsNaN(testAnalyzer.chiSquareValue), "Chi-square should not be NaN for %s", tc.name)
		s.GreaterOrEqual(testAnalyzer.chiSquareValue, 0.0, "Chi-square should be non-negative for %s", tc.name)

		// Verify randomness score is valid
		s.False(math.IsNaN(testAnalyzer.randomnessScore), "Randomness score should not be NaN for %s", tc.name)
		s.GreaterOrEqual(testAnalyzer.randomnessScore, 0.0, "Randomness score should be non-negative for %s", tc.name)
		s.LessOrEqual(testAnalyzer.randomnessScore, 100.0, "Randomness score should be <= 100 for %s", tc.name)
	}
}

// TestCorrelationCalculationEdgeCases tests Pearson correlation calculation edge cases
func (s *AnalyzerTestSuite) TestCorrelationCalculationEdgeCases() {
	testCases := []struct {
		name           string
		x              []float64
		y              []float64
		expectError    bool
		expectedCorr   float64
		expectedPValue float64
	}{
		{
			name:        "Empty arrays",
			x:           []float64{},
			y:           []float64{},
			expectError: false, // Should handle gracefully
		},
		{
			name:        "Single value arrays",
			x:           []float64{1.0},
			y:           []float64{2.0},
			expectError: false, // Should handle gracefully
		},
		{
			name:         "Perfect positive correlation",
			x:            []float64{1, 2, 3, 4, 5},
			y:            []float64{2, 4, 6, 8, 10},
			expectError:  false,
			expectedCorr: 1.0,
		},
		{
			name:         "Perfect negative correlation",
			x:            []float64{1, 2, 3, 4, 5},
			y:            []float64{10, 8, 6, 4, 2},
			expectError:  false,
			expectedCorr: -1.0,
		},
		{
			name:        "All zeros",
			x:           []float64{0, 0, 0, 0, 0},
			y:           []float64{0, 0, 0, 0, 0},
			expectError: false, // Should handle gracefully
		},
		{
			name:        "Different lengths",
			x:           []float64{1, 2, 3},
			y:           []float64{1, 2, 3, 4, 5},
			expectError: false, // Should handle gracefully
		},
		{
			name:        "Constant x values",
			x:           []float64{5, 5, 5, 5, 5},
			y:           []float64{1, 2, 3, 4, 5},
			expectError: false, // Should handle gracefully (NaN is expected)
		},
		{
			name:        "Constant y values",
			x:           []float64{1, 2, 3, 4, 5},
			y:           []float64{7, 7, 7, 7, 7},
			expectError: false, // Should handle gracefully (NaN is expected)
		},
	}

	for _, tc := range testCases {
		corr, pValue := calculatePearsonCorrelation(tc.x, tc.y)

		switch tc.name {
		case "Perfect positive correlation":
			s.Greater(corr, 0.99, "Correlation should be close to 1.0 for %s", tc.name)
		case "Perfect negative correlation":
			s.Less(corr, -0.99, "Correlation should be close to -1.0 for %s", tc.name)
		}

		// P-value should be a valid number or NaN
		s.False(math.IsInf(pValue, 0), "P-value should not be infinite for %s", tc.name)

		// For edge cases, correlation might be NaN which is acceptable
		if !math.IsNaN(corr) {
			s.GreaterOrEqual(corr, -1.0, "Correlation should be >= -1.0 for %s", tc.name)
			s.LessOrEqual(corr, 1.0, "Correlation should be <= 1.0 for %s", tc.name)
		}
	}
}

// TestNumberAnalysisEdgeCases tests number analysis edge cases
func (s *AnalyzerTestSuite) TestNumberAnalysisEdgeCases() {
	// Test with extreme number distributions
	extremeAnalyzer := &Analyzer{
		drawings: []Drawing{
			// All same numbers
			{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Numbers: []int{1, 1, 1, 1, 1}, LuckyBall: 1},
			{Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Numbers: []int{48, 48, 48, 48, 48}, LuckyBall: 18},
		},
		mainNumbers:  make(map[int]*NumberInfo),
		luckyBalls:   make(map[int]*NumberInfo),
		pairPatterns: make(map[string]*CombinationPattern),
		patternStats: &PatternStats{
			OddEvenPatterns:    make(map[string]int),
			SumRanges:          make(map[int]int),
			ConsecutiveCount:   0,
			DecadeDistribution: make(map[int]int),
		},
		config: &AnalysisConfig{
			RecentWindow:     3,
			MinGapMultiplier: 1.5,
			ConfidenceLevel:  0.95,
			OutputMode:       "detailed",
		},
	}

	// Initialize all numbers
	for i := 1; i <= 48; i++ {
		extremeAnalyzer.mainNumbers[i] = &NumberInfo{
			Number:            i,
			TotalFrequency:    0,
			RecentFrequency:   0,
			GapsSinceDrawn:    []int{},
			CurrentGap:        0,
			AverageGap:        0.0,
			ExpectedFrequency: 0.0,
			StandardDeviation: 0.0,
		}
	}
	for i := 1; i <= 18; i++ {
		extremeAnalyzer.luckyBalls[i] = &NumberInfo{
			Number:            i,
			TotalFrequency:    0,
			RecentFrequency:   0,
			GapsSinceDrawn:    []int{},
			CurrentGap:        0,
			AverageGap:        0.0,
			ExpectedFrequency: 0.0,
			StandardDeviation: 0.0,
		}
	}

	// Test with extreme distributions
	ctx := context.Background()

	// Test GetTopNumbers with edge cases
	topNumbers := extremeAnalyzer.GetTopNumbers(0, false)
	s.Empty(topNumbers)

	topNumbers = extremeAnalyzer.GetTopNumbers(100, false)
	s.LessOrEqual(len(topNumbers), 48)

	// Test GetOverdueNumbers with edge cases
	overdueNumbers := extremeAnalyzer.GetOverdueNumbers(0)
	s.Empty(overdueNumbers)

	overdueNumbers = extremeAnalyzer.GetOverdueNumbers(100)
	s.LessOrEqual(len(overdueNumbers), 48)

	// Test GenerateRecommendations with edge cases
	recommendations, err := extremeAnalyzer.GenerateRecommendations(ctx, 0)
	s.Require().NoError(err)
	s.Empty(recommendations)

	recommendations, err = extremeAnalyzer.GenerateRecommendations(ctx, 1)
	s.Require().NoError(err)
	s.Len(recommendations, 1)
}

// TestPrintAnalysisEdgeCases tests print analysis functions with edge cases
func (s *AnalyzerTestSuite) TestPrintAnalysisEdgeCases() {
	ctx := context.Background()

	// Create analyzer with minimal data for testing edge cases in print functions
	minimalAnalyzer := &Analyzer{
		drawings: []Drawing{
			{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Numbers: []int{1, 2, 3, 4, 5}, LuckyBall: 1},
		},
		mainNumbers:  make(map[int]*NumberInfo),
		luckyBalls:   make(map[int]*NumberInfo),
		pairPatterns: make(map[string]*CombinationPattern),
		patternStats: &PatternStats{
			OddEvenPatterns:    make(map[string]int),
			SumRanges:          make(map[int]int),
			ConsecutiveCount:   0,
			DecadeDistribution: make(map[int]int),
		},
		config:            &AnalysisConfig{OutputMode: "detailed"},
		correlationEngine: NewCorrelationEngine(nil),
	}

	// Initialize minimal number info to prevent panics
	for i := 1; i <= 48; i++ {
		minimalAnalyzer.mainNumbers[i] = &NumberInfo{
			Number:            i,
			TotalFrequency:    0,
			RecentFrequency:   0,
			GapsSinceDrawn:    []int{},
			CurrentGap:        1,
			AverageGap:        1.0,
			ExpectedFrequency: 1.0,
			StandardDeviation: 0.0,
		}
	}
	for i := 1; i <= 18; i++ {
		minimalAnalyzer.luckyBalls[i] = &NumberInfo{
			Number:            i,
			TotalFrequency:    0,
			RecentFrequency:   0,
			GapsSinceDrawn:    []int{},
			CurrentGap:        1,
			AverageGap:        1.0,
			ExpectedFrequency: 1.0,
			StandardDeviation: 0.0,
		}
	}

	// Set correlation engine analyzer reference to avoid nil pointer
	minimalAnalyzer.correlationEngine.analyzer = minimalAnalyzer

	// Test print functions with minimal data
	err := minimalAnalyzer.printSimpleAnalysis(ctx)
	s.Require().NoError(err)

	err = minimalAnalyzer.printDetailedAnalysis(ctx)
	s.Require().NoError(err)

	err = minimalAnalyzer.printStatisticalAnalysis(ctx)
	s.Require().NoError(err)

	err = minimalAnalyzer.printCosmicAnalysis(ctx)
	s.NoError(err)
}

// Run the test suite
func TestAnalyzerSuite(t *testing.T) {
	suite.Run(t, new(AnalyzerTestSuite))
}

// Benchmark tests
func BenchmarkAnalyzerCreation(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, _ = NewAnalyzer(ctx, "../../data/lucky-numbers-history.csv", nil) // ignore error in benchmark
	}
}

func BenchmarkRecommendationGeneration(b *testing.B) {
	ctx := context.Background()
	analyzer, err := NewAnalyzer(ctx, "../../data/lucky-numbers-history.csv", nil)
	if err != nil {
		b.Skip("Skipping benchmark: CSV file not available")
	}
	if analyzer == nil {
		b.Skip("Skipping benchmark: analyzer is nil")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = analyzer.GenerateRecommendations(ctx, 5) // ignore error in benchmark
	}
}

func BenchmarkPatternAnalysis(b *testing.B) {
	ctx := context.Background()
	analyzer, err := NewAnalyzer(ctx, "../../data/lucky-numbers-history.csv", nil)
	if err != nil {
		b.Skip("Skipping benchmark: CSV file not available")
	}
	if analyzer == nil {
		b.Skip("Skipping benchmark: analyzer is nil")
	}

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
		_, _ = NewAnalyzer(ctx, "../../data/lucky-numbers-history.csv", nil)
	}
}

func BenchmarkCosmicCorrelations(b *testing.B) {
	ctx := context.Background()
	analyzer, err := NewAnalyzer(ctx, "../../data/lucky-numbers-history.csv", nil)
	if err != nil {
		b.Skip("Skipping benchmark: CSV file not available")
	}
	if analyzer == nil {
		b.Skip("Skipping benchmark: analyzer is nil")
	}

	correlationEngine := NewCorrelationEngine(analyzer)
	_ = correlationEngine.EnrichWithCosmicData(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = correlationEngine.AnalyzeCorrelations(ctx)
	}
}

func BenchmarkReportGeneration(b *testing.B) {
	ctx := context.Background()
	analyzer, err := NewAnalyzer(ctx, "../../data/lucky-numbers-history.csv", nil)
	if err != nil {
		b.Skip("Skipping benchmark: CSV file not available")
	}
	if analyzer == nil {
		b.Skip("Skipping benchmark: analyzer is nil")
	}

	correlationEngine := NewCorrelationEngine(analyzer)
	_ = correlationEngine.EnrichWithCosmicData(ctx)
	_ = correlationEngine.AnalyzeCorrelations(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = correlationEngine.GenerateCosmicReport()
	}
}

// TestCLIArgumentParsing tests command-line argument parsing logic
func (s *AnalyzerTestSuite) TestCLIArgumentParsing() {
	// Test parsing CLI arguments into config
	testCases := []struct {
		name           string
		args           []string
		expectedMode   string
		expectedFormat string
		expectedWindow int
	}{
		{
			name:           "Simple mode",
			args:           []string{"program", "--simple"},
			expectedMode:   "simple",
			expectedFormat: "console",
			expectedWindow: 50,
		},
		{
			name:           "Statistical mode",
			args:           []string{"program", "--statistical"},
			expectedMode:   "statistical",
			expectedFormat: "console",
			expectedWindow: 50,
		},
		{
			name:           "Cosmic mode",
			args:           []string{"program", "--cosmic"},
			expectedMode:   "cosmic",
			expectedFormat: "console",
			expectedWindow: 50,
		},
		{
			name:           "JSON export",
			args:           []string{"program", "--export-json"},
			expectedMode:   "detailed",
			expectedFormat: exportFormatJSON,
			expectedWindow: 50,
		},
		{
			name:           "CSV export",
			args:           []string{"program", "--export-csv"},
			expectedMode:   "detailed",
			expectedFormat: exportFormatCSV,
			expectedWindow: 50,
		},
		{
			name:           "Recent window",
			args:           []string{"program", "--recent", "25"},
			expectedMode:   "detailed",
			expectedFormat: "console",
			expectedWindow: 25,
		},
		{
			name:           "Combined flags",
			args:           []string{"program", "--cosmic", "--export-json", "--recent", "30"},
			expectedMode:   "cosmic",
			expectedFormat: exportFormatJSON,
			expectedWindow: 30,
		},
	}

	for _, tc := range testCases {
		config := parseCLIArgs(tc.args)
		s.Equal(tc.expectedMode, config.OutputMode, "Test case: %s", tc.name)
		s.Equal(tc.expectedFormat, config.ExportFormat, "Test case: %s", tc.name)
		s.Equal(tc.expectedWindow, config.RecentWindow, "Test case: %s", tc.name)
	}
}

// Test help functionality
func (s *AnalyzerTestSuite) TestPrintHelp() {
	// Test that printHelp function runs without error
	s.NotPanics(func() {
		printHelp()
	})
}

// TestCLIEdgeCases tests edge cases in CLI argument parsing
func (s *AnalyzerTestSuite) TestCLIEdgeCases() {
	// Test invalid recent window value
	config := parseCLIArgs([]string{"program", "--recent", "invalid"})
	s.Equal(50, config.RecentWindow) // Should remain default

	// Test recent without value
	config = parseCLIArgs([]string{"program", "--recent"})
	s.Equal(50, config.RecentWindow) // Should remain default

	// Test unknown flags are ignored
	config = parseCLIArgs([]string{"program", "--unknown-flag", "--simple"})
	s.Equal("simple", config.OutputMode) // Should still process known flags

	// Test empty args
	config = parseCLIArgs([]string{"program"})
	s.Equal("detailed", config.OutputMode) // Should use defaults
	s.Equal("console", config.ExportFormat)
	s.Equal(50, config.RecentWindow)
}

// parseCLIArgs extracts the CLI parsing logic for testing
func parseCLIArgs(args []string) *AnalysisConfig {
	// Default configuration
	config := &AnalysisConfig{
		RecentWindow:     50,
		MinGapMultiplier: 1.5,
		ConfidenceLevel:  0.95,
		OutputMode:       "detailed",
		ExportFormat:     "console",
	}

	// Parse command line arguments
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--simple":
				config.OutputMode = "simple"
			case "--statistical":
				config.OutputMode = "statistical"
			case "--cosmic":
				config.OutputMode = "cosmic"
			case "--export-json":
				config.ExportFormat = exportFormatJSON
			case "--export-csv":
				config.ExportFormat = exportFormatCSV
			case "--recent":
				if i+1 < len(args) {
					if val, err := strconv.Atoi(args[i+1]); err == nil {
						config.RecentWindow = val
						i++
					}
				}
			}
		}
	}

	return config
}
