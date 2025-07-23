package main

import (
	"context"
	"math"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// AnalyzerTestSuite defines the test suite for lottery analyzer
type AnalyzerTestSuite struct {
	suite.Suite
	ctx      context.Context
	analyzer *Analyzer
	testFile string
}

// SetupSuite runs once before all tests
func (s *AnalyzerTestSuite) SetupSuite() {
	s.ctx = context.Background()
	
	// Create a test CSV file
	s.testFile = "test_lottery_data.csv"
	content := `Date,Number 1,Number 2,Number 3,Number 4,Number 5,Lucky Ball
01/15/2024,5,12,23,34,45,7
01/12/2024,3,15,22,38,44,12
01/09/2024,5,18,23,35,42,7
01/06/2024,7,12,25,33,48,15
01/03/2024,2,11,23,34,41,3`
	
	err := os.WriteFile(s.testFile, []byte(content), 0644)
	require.NoError(s.T(), err)
}

// TearDownSuite runs once after all tests
func (s *AnalyzerTestSuite) TearDownSuite() {
	os.Remove(s.testFile)
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
	
	analyzer, err := NewAnalyzer(s.ctx, s.testFile, config)
	require.NoError(s.T(), err)
	s.analyzer = analyzer
}

// TestNewAnalyzerValidFile tests creating analyzer with valid file
func (s *AnalyzerTestSuite) TestNewAnalyzerValidFile() {
	assert.NotNil(s.T(), s.analyzer)
	assert.Equal(s.T(), 5, len(s.analyzer.drawings))
	assert.NotNil(s.T(), s.analyzer.mainNumbers)
	assert.NotNil(s.T(), s.analyzer.luckyBalls)
}

// TestNewAnalyzerInvalidFile tests creating analyzer with invalid file
func (s *AnalyzerTestSuite) TestNewAnalyzerInvalidFile() {
	_, err := NewAnalyzer(s.ctx, "nonexistent.csv", nil)
	assert.Error(s.T(), err)
}

// TestNewAnalyzerDefaultConfig tests default configuration
func (s *AnalyzerTestSuite) TestNewAnalyzerDefaultConfig() {
	analyzer, err := NewAnalyzer(s.ctx, s.testFile, nil)
	require.NoError(s.T(), err)
	
	assert.Equal(s.T(), 50, analyzer.config.RecentWindow)
	assert.Equal(s.T(), 1.5, analyzer.config.MinGapMultiplier)
	assert.Equal(s.T(), "detailed", analyzer.config.OutputMode)
}

// TestParseDrawings tests parsing CSV data
func (s *AnalyzerTestSuite) TestParseDrawings() {
	assert.Equal(s.T(), 5, len(s.analyzer.drawings))
	
	// Check first drawing (oldest after reversal - which is 01/03/2024)
	firstDrawing := s.analyzer.drawings[0]
	assert.Equal(s.T(), 2024, firstDrawing.Date.Year())
	assert.Equal(s.T(), time.January, firstDrawing.Date.Month())
	assert.Equal(s.T(), 3, firstDrawing.Date.Day())
	assert.Equal(s.T(), []int{2, 11, 23, 34, 41}, firstDrawing.Numbers)
	assert.Equal(s.T(), 3, firstDrawing.LuckyBall)
}

// TestNumberFrequencyTracking tests frequency calculation
func (s *AnalyzerTestSuite) TestNumberFrequencyTracking() {
	// Number 23 appears 3 times in test data
	info := s.analyzer.mainNumbers[23]
	assert.Equal(s.T(), 3, info.TotalFrequency)
	
	// Number 5 appears 2 times
	info5 := s.analyzer.mainNumbers[5]
	assert.Equal(s.T(), 2, info5.TotalFrequency)
	
	// Lucky ball 7 appears 2 times
	lb7 := s.analyzer.luckyBalls[7]
	assert.Equal(s.T(), 2, lb7.TotalFrequency)
}

// TestRecentFrequencyTracking tests recent frequency calculation
func (s *AnalyzerTestSuite) TestRecentFrequencyTracking() {
	// With recent window of 3, check recent frequencies
	info23 := s.analyzer.mainNumbers[23]
	assert.Equal(s.T(), 2, info23.RecentFrequency) // Appears in 2 of last 3
	
	info2 := s.analyzer.mainNumbers[2]
	assert.Equal(s.T(), 1, info2.RecentFrequency) // Appears once in last 3 (in 01/03/2024)
}

// TestGapCalculation tests gap tracking
func (s *AnalyzerTestSuite) TestGapCalculation() {
	// Number 5 appears at indices 2 and 4 (01/09 and 01/15)
	info5 := s.analyzer.mainNumbers[5]
	assert.Contains(s.T(), info5.GapsSinceDrawn, 2)
	assert.Equal(s.T(), 4, info5.CurrentGap) // Last appeared at index 4
}

// TestPatternAnalysis tests pattern detection
func (s *AnalyzerTestSuite) TestPatternAnalysis() {
	// Check odd/even patterns
	assert.Greater(s.T(), len(s.analyzer.patternStats.OddEvenPatterns), 0)
	
	// Check sum ranges
	assert.Greater(s.T(), len(s.analyzer.patternStats.SumRanges), 0)
}

// TestCombinationPatterns tests pair/triple/quad tracking
func (s *AnalyzerTestSuite) TestCombinationPatterns() {
	// Check that pairs are tracked
	assert.Greater(s.T(), len(s.analyzer.pairPatterns), 0)
	
	// Check specific pair (5-23 appears twice)
	pairKey := "5-23"
	pattern, exists := s.analyzer.pairPatterns[pairKey]
	assert.True(s.T(), exists)
	assert.Equal(s.T(), 2, pattern.Frequency)
}

// TestChiSquareCalculation tests statistical calculations
func (s *AnalyzerTestSuite) TestChiSquareCalculation() {
	assert.Greater(s.T(), s.analyzer.chiSquareValue, 0.0)
	assert.GreaterOrEqual(s.T(), s.analyzer.randomnessScore, 0.0)
	assert.LessOrEqual(s.T(), s.analyzer.randomnessScore, 100.0)
}

// TestGetTopNumbers tests retrieving top frequent numbers
func (s *AnalyzerTestSuite) TestGetTopNumbers() {
	topNumbers := s.analyzer.GetTopNumbers(5, false)
	assert.LessOrEqual(s.T(), len(topNumbers), 5)
	
	// Verify sorted by frequency
	for i := 1; i < len(topNumbers); i++ {
		assert.GreaterOrEqual(s.T(), topNumbers[i-1].TotalFrequency, topNumbers[i].TotalFrequency)
	}
}

// TestGetOverdueNumbers tests retrieving overdue numbers
func (s *AnalyzerTestSuite) TestGetOverdueNumbers() {
	overdueNumbers := s.analyzer.GetOverdueNumbers(10)
	
	// Verify sorted by overdue ratio
	for i := 1; i < len(overdueNumbers); i++ {
		ratioI := float64(overdueNumbers[i-1].CurrentGap) / overdueNumbers[i-1].AverageGap
		ratioJ := float64(overdueNumbers[i].CurrentGap) / overdueNumbers[i].AverageGap
		assert.GreaterOrEqual(s.T(), ratioI, ratioJ)
	}
}

// TestGenerateRecommendations tests recommendation generation
func (s *AnalyzerTestSuite) TestGenerateRecommendations() {
	recommendations, err := s.analyzer.GenerateRecommendations(s.ctx, 3)
	require.NoError(s.T(), err)
	
	assert.LessOrEqual(s.T(), len(recommendations), 3)
	
	for _, rec := range recommendations {
		assert.Equal(s.T(), 5, len(rec.Numbers))
		assert.Greater(s.T(), rec.LuckyBall, 0)
		assert.LessOrEqual(s.T(), rec.LuckyBall, 18)
		assert.NotEmpty(s.T(), rec.Strategy)
		assert.NotEmpty(s.T(), rec.Explanation)
		assert.Greater(s.T(), rec.Confidence, 0.0)
		assert.LessOrEqual(s.T(), rec.Confidence, 1.0)
	}
}

// TestScoreNumbersByStrategy tests different scoring strategies
func (s *AnalyzerTestSuite) TestScoreNumbersByStrategy() {
	strategies := []string{"balanced", "hot", "overdue", "pattern", "frequency"}
	
	for _, strategy := range strategies {
		scored := s.analyzer.scoreNumbersByStrategy(strategy)
		assert.Greater(s.T(), len(scored), 0)
		
		// Verify sorted by score
		for i := 1; i < len(scored); i++ {
			assert.GreaterOrEqual(s.T(), scored[i-1].Score, scored[i].Score)
		}
	}
}

// TestExportJSON tests JSON export functionality
func (s *AnalyzerTestSuite) TestExportJSON() {
	s.analyzer.config.ExportFormat = "json"
	testFile := "test_export.json"
	
	err := s.analyzer.ExportAnalysis(s.ctx, testFile)
	require.NoError(s.T(), err)
	
	// Verify file exists
	_, err = os.Stat(testFile)
	assert.NoError(s.T(), err)
	
	// Clean up
	os.Remove(testFile)
}

// TestExportCSV tests CSV export functionality
func (s *AnalyzerTestSuite) TestExportCSV() {
	s.analyzer.config.ExportFormat = "csv"
	testFile := "test_export.csv"
	
	err := s.analyzer.ExportAnalysis(s.ctx, testFile)
	require.NoError(s.T(), err)
	
	// Verify file exists
	_, err = os.Stat(testFile)
	assert.NoError(s.T(), err)
	
	// Clean up
	os.Remove(testFile)
}

// TestContextCancellation tests context cancellation handling
func (s *AnalyzerTestSuite) TestContextCancellation() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel() // Cancel immediately
	
	_, err := NewAnalyzer(ctx, s.testFile, nil)
	// Should still work as parsing is fast, but if it fails, should be context error
	if err != nil {
		assert.Contains(s.T(), err.Error(), "context canceled")
	}
}

// TestEmptyDataHandling tests handling of empty or invalid data
func (s *AnalyzerTestSuite) TestEmptyDataHandling() {
	// Create file with header only
	emptyFile := "empty_test.csv"
	content := `Date,Number 1,Number 2,Number 3,Number 4,Number 5,Lucky Ball`
	err := os.WriteFile(emptyFile, []byte(content), 0644)
	require.NoError(s.T(), err)
	defer os.Remove(emptyFile)
	
	analyzer, err := NewAnalyzer(s.ctx, emptyFile, nil)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 0, len(analyzer.drawings))
}

// TestInvalidDataHandling tests handling of malformed data
func (s *AnalyzerTestSuite) TestInvalidDataHandling() {
	// Create file with invalid data
	invalidFile := "invalid_test.csv"
	content := `Date,Number 1,Number 2,Number 3,Number 4,Number 5,Lucky Ball
invalid_date,5,12,23,34,45,7
01/12/2024,abc,15,22,38,44,12
01/09/2024,5,18,23,35,42,xyz`
	err := os.WriteFile(invalidFile, []byte(content), 0644)
	require.NoError(s.T(), err)
	defer os.Remove(invalidFile)
	
	analyzer, err := NewAnalyzer(s.ctx, invalidFile, nil)
	require.NoError(s.T(), err)
	// Should skip invalid rows
	assert.Equal(s.T(), 0, len(analyzer.drawings))
}

// TestNumberRangeValidation tests that numbers are within valid ranges
func (s *AnalyzerTestSuite) TestNumberRangeValidation() {
	for _, drawing := range s.analyzer.drawings {
		for _, num := range drawing.Numbers {
			assert.GreaterOrEqual(s.T(), num, 1)
			assert.LessOrEqual(s.T(), num, 48)
		}
		assert.GreaterOrEqual(s.T(), drawing.LuckyBall, 1)
		assert.LessOrEqual(s.T(), drawing.LuckyBall, 18)
	}
}

// TestStatisticalMeasures tests statistical calculations
func (s *AnalyzerTestSuite) TestStatisticalMeasures() {
	// Check that all numbers have expected frequency calculated
	for _, info := range s.analyzer.mainNumbers {
		assert.Greater(s.T(), info.ExpectedFrequency, 0.0)
	}
	
	// Check standard deviation is calculated for numbers with gaps
	for _, info := range s.analyzer.mainNumbers {
		if len(info.GapsSinceDrawn) > 0 {
			assert.GreaterOrEqual(s.T(), info.StandardDeviation, 0.0)
		}
	}
}

// TestCosmicCorrelationEngine tests the cosmic correlation functionality
func (s *AnalyzerTestSuite) TestCosmicCorrelationEngine() {
	// Test correlation engine initialization
	assert.NotNil(s.T(), s.analyzer.correlationEngine)
	
	// Test cosmic data enrichment
	err := s.analyzer.correlationEngine.EnrichWithCosmicData(s.ctx)
	require.NoError(s.T(), err)
	
	// Verify cosmic data was added
	assert.Greater(s.T(), len(s.analyzer.correlationEngine.cosmicData), 0)
	
	// Test correlation analysis
	err = s.analyzer.correlationEngine.AnalyzeCorrelations(s.ctx)
	require.NoError(s.T(), err)
	
	// Verify correlations were calculated
	assert.Greater(s.T(), len(s.analyzer.correlationEngine.correlationResults), 0)
}

// TestMoonPhaseCalculation tests moon phase calculations
func (s *AnalyzerTestSuite) TestMoonPhaseCalculation() {
	// Test known date - January 1, 2024
	testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	phase, illumination := s.analyzer.correlationEngine.calculateMoonPhase(testDate)
	
	// Verify phase is between 0 and 1
	assert.GreaterOrEqual(s.T(), phase, 0.0)
	assert.Less(s.T(), phase, 1.0)
	
	// Verify illumination is between 0 and 1
	assert.GreaterOrEqual(s.T(), illumination, 0.0)
	assert.LessOrEqual(s.T(), illumination, 1.0)
	
	// Test phase name
	phaseName := s.analyzer.correlationEngine.getMoonPhaseName(phase)
	assert.NotEmpty(s.T(), phaseName)
	assert.Contains(s.T(), []string{"New Moon", "Waxing Crescent", "First Quarter", "Waxing Gibbous", 
		"Full Moon", "Waning Gibbous", "Last Quarter", "Waning Crescent"}, phaseName)
}

// TestZodiacCalculation tests zodiac sign calculations
func (s *AnalyzerTestSuite) TestZodiacCalculation() {
	// Test known dates
	testCases := []struct {
		date     time.Time
		expected string
	}{
		{time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "Capricorn"},
		{time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), "Gemini"},
		{time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC), "Sagittarius"},
	}
	
	for _, tc := range testCases {
		zodiac := s.analyzer.correlationEngine.getZodiacSign(tc.date)
		assert.Equal(s.T(), tc.expected, zodiac)
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
		assert.Equal(s.T(), tc.expected, season)
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
		assert.True(s.T(), exists, "Missing position for %s", planet)
		assert.GreaterOrEqual(s.T(), position, 0.0)
		assert.Less(s.T(), position, 360.0)
	}
}

// TestCosmicPredictions tests cosmic-based number predictions
func (s *AnalyzerTestSuite) TestCosmicPredictions() {
	numbers := s.analyzer.correlationEngine.PredictBasedOnCosmicConditions()
	
	// Verify we get 5 numbers
	assert.Equal(s.T(), 5, len(numbers))
	
	// Verify all numbers are in valid range and unique
	used := make(map[int]bool)
	for _, num := range numbers {
		assert.GreaterOrEqual(s.T(), num, 1)
		assert.LessOrEqual(s.T(), num, 48)
		assert.False(s.T(), used[num], "Duplicate number: %d", num)
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
	assert.Greater(s.T(), corr, 0.99)
	// P-value should be very small for significant correlation
	if !math.IsNaN(pValue) {
		assert.LessOrEqual(s.T(), pValue, 0.05)
	}
	
	// Test with no correlation
	yRandom := []float64{5, 1, 8, 2, 9}
	corrRandom, pValueRandom := calculatePearsonCorrelation(x, yRandom)
	
	// Should be closer to 0
	assert.Less(s.T(), math.Abs(corrRandom), 0.9)
	assert.GreaterOrEqual(s.T(), pValueRandom, 0.0)
}

// TestCosmicReportGeneration tests cosmic report generation
func (s *AnalyzerTestSuite) TestCosmicReportGeneration() {
	// Ensure we have correlation results
	err := s.analyzer.correlationEngine.AnalyzeCorrelations(s.ctx)
	require.NoError(s.T(), err)
	
	report := s.analyzer.correlationEngine.GenerateCosmicReport()
	
	// Verify report contains expected sections
	assert.Contains(s.T(), report, "COSMIC CORRELATION ANALYSIS")
	assert.Contains(s.T(), report, "LUNAR CORRELATIONS")
	assert.Contains(s.T(), report, "CURRENT COSMIC CONDITIONS")
	assert.Contains(s.T(), report, "DISCLAIMER")
}

// TestCosmicAnalysisMode tests the cosmic analysis output mode
func (s *AnalyzerTestSuite) TestCosmicAnalysisMode() {
	// Change to cosmic mode
	s.analyzer.config.OutputMode = "cosmic"
	
	// Test that cosmic analysis runs without error
	err := s.analyzer.RunAnalysis(s.ctx)
	assert.NoError(s.T(), err)
}

// Run the test suite
func TestAnalyzerSuite(t *testing.T) {
	suite.Run(t, new(AnalyzerTestSuite))
}

// Benchmark tests
func BenchmarkAnalyzerCreation(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		NewAnalyzer(ctx, "lucky-numbers-history.csv", nil)
	}
}

func BenchmarkRecommendationGeneration(b *testing.B) {
	ctx := context.Background()
	analyzer, _ := NewAnalyzer(ctx, "lucky-numbers-history.csv", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.GenerateRecommendations(ctx, 5)
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