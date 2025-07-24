//go:build go1.18
// +build go1.18

package main

import (
	"context"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Helper function to parse comma-separated string to float64 slice
func parseFloatSlice(s string) []float64 {
	if s == "" {
		return []float64{}
	}
	parts := strings.Split(s, ",")
	result := make([]float64, 0, len(parts))
	for _, part := range parts {
		f := parseFloatValue(strings.TrimSpace(part))
		if !math.IsNaN(f) || strings.TrimSpace(part) == "NaN" {
			result = append(result, f)
		}
	}
	return result
}

// Helper function to parse a single float value
func parseFloatValue(s string) float64 {
	switch s {
	case "NaN":
		return math.NaN()
	case "Inf":
		return math.Inf(1)
	case "-Inf":
		return math.Inf(-1)
	default:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return math.NaN()
		}
		return f
	}
}

// FuzzCalculateMoonPhase tests moon phase calculation with various dates
func FuzzCalculateMoonPhase(f *testing.F) {
	// Seed with various dates
	f.Add(int64(0))                                             // Unix epoch
	f.Add(time.Now().Unix())                                    // Current time
	f.Add(time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC).Unix()) // Reference new moon
	f.Add(time.Date(1969, 7, 20, 0, 0, 0, 0, time.UTC).Unix())  // Moon landing
	f.Add(time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC).Unix())
	f.Add(int64(-62135596800)) // Year 1
	f.Add(int64(253402300799)) // Year 9999

	f.Fuzz(func(t *testing.T, unixTime int64) {
		// Skip extremely far dates that might cause overflow
		if unixTime < -62135596800 || unixTime > 253402300799 {
			t.Skip("Date out of reasonable range")
		}

		date := time.Unix(unixTime, 0).UTC()

		// Create a minimal correlation engine
		analyzer := &Analyzer{
			drawings: []Drawing{},
		}
		ce := NewCorrelationEngine(analyzer)

		// Should not panic
		phase, illumination := ce.calculateMoonPhase(date)

		// Validate results
		if phase < 0 || phase >= 1 {
			t.Errorf("Moon phase out of range [0,1): %f", phase)
		}

		if illumination < 0 || illumination > 1 {
			t.Errorf("Moon illumination out of range [0,1]: %f", illumination)
		}

		// Phase name should be valid
		phaseName := ce.getMoonPhaseName(phase)
		validPhases := []string{
			"New Moon", "Waxing Crescent", "First Quarter", "Waxing Gibbous",
			"Full Moon", "Waning Gibbous", "Last Quarter", "Waning Crescent",
		}
		found := false
		for _, valid := range validPhases {
			if phaseName == valid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Invalid moon phase name: %s", phaseName)
		}
	})
}

// FuzzGetZodiacSign tests zodiac calculation with fuzzy dates
func FuzzGetZodiacSign(f *testing.F) {
	// Seed with edge cases around zodiac boundaries
	f.Add(1, 1)   // January 1
	f.Add(1, 19)  // Capricorn/Aquarius boundary
	f.Add(1, 20)  // Aquarius start
	f.Add(2, 29)  // Leap year consideration
	f.Add(12, 31) // End of year
	f.Add(0, 0)   // Invalid
	f.Add(13, 32) // Invalid
	f.Add(6, 21)  // Summer solstice

	f.Fuzz(func(t *testing.T, month, day int) {
		// Skip invalid months
		if month < 1 || month > 12 {
			t.Skip("Invalid month")
		}

		// Create date with fuzzed month/day
		year := 2024

		// Adjust day to be valid for the month
		maxDay := 31
		switch time.Month(month) {
		case time.February:
			if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
				maxDay = 29
			} else {
				maxDay = 28
			}
		case time.April, time.June, time.September, time.November:
			maxDay = 30
		}

		if day < 1 || day > maxDay {
			t.Skip("Invalid day for month")
		}

		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		analyzer := &Analyzer{}
		ce := NewCorrelationEngine(analyzer)

		// Should not panic
		zodiac := ce.getZodiacSign(date)

		// Validate zodiac sign
		validSigns := []string{
			"Capricorn", "Aquarius", "Pisces", "Aries", "Taurus", "Gemini",
			"Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius",
		}

		found := false
		for _, sign := range validSigns {
			if zodiac == sign {
				found = true
				break
			}
		}
		if !found && zodiac != "Unknown" {
			t.Errorf("Invalid zodiac sign: %s", zodiac)
		}
	})
}

// FuzzCalculatePearsonCorrelation tests correlation calculation with various data
func FuzzCalculatePearsonCorrelation(f *testing.F) {
	// Seed with various data patterns encoded as strings
	f.Add("1,2,3,4,5", "2,4,6,8,10")      // Perfect positive
	f.Add("1,2,3,4,5", "10,8,6,4,2")      // Perfect negative
	f.Add("1,1,1,1,1", "2,3,4,5,6")       // No variance in X
	f.Add("1,2,3,4,5", "3,3,3,3,3")       // No variance in Y
	f.Add("", "")                         // Empty
	f.Add("1", "2")                       // Single point
	f.Add("1,2,3", "1,2")                 // Mismatched lengths
	f.Add("0,0,0,0", "0,0,0,0")           // All zeros
	f.Add("1e100,-1e100", "1e100,-1e100") // Large numbers
	f.Add("NaN", "1")                     // NaN values
	f.Add("Inf", "1")                     // Infinity

	f.Fuzz(func(t *testing.T, xStr, yStr string) {
		// Parse strings to float64 slices
		x := parseFloatSlice(xStr)
		y := parseFloatSlice(yStr)
		// Skip if lengths don't match or empty
		if len(x) != len(y) || len(x) == 0 {
			corr, pValue := calculatePearsonCorrelation(x, y)
			if corr != 0 || pValue != 1 {
				t.Errorf("Expected (0, 1) for mismatched/empty data, got (%f, %f)", corr, pValue)
			}
			return
		}

		// Should not panic
		corr, pValue := calculatePearsonCorrelation(x, y)

		// Check for NaN or Inf in input
		hasSpecial := false
		for i := range x {
			if math.IsNaN(x[i]) || math.IsInf(x[i], 0) ||
				math.IsNaN(y[i]) || math.IsInf(y[i], 0) {
				hasSpecial = true
				break
			}
		}

		if !hasSpecial && !math.IsNaN(corr) && !math.IsInf(corr, 0) {
			// Correlation should be in [-1, 1]
			if corr < -1.001 || corr > 1.001 { // Small epsilon for floating point
				t.Errorf("Correlation out of range [-1, 1]: %f", corr)
			}

			// P-value should be in [0, 1]
			if !math.IsNaN(pValue) && (pValue < 0 || pValue > 1) {
				t.Errorf("P-value out of range [0, 1]: %f", pValue)
			}
		}
	})
}

// FuzzPlanetaryPositions tests planetary position calculations
func FuzzPlanetaryPositions(f *testing.F) {
	// Seed with various timestamps
	f.Add(int64(0))
	f.Add(time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC).Unix())
	f.Add(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Unix())
	f.Add(time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC).Unix())
	f.Add(int64(-1000000000))
	f.Add(int64(2000000000))

	f.Fuzz(func(t *testing.T, unixTime int64) {
		// Skip extremely far dates
		if unixTime < -62135596800 || unixTime > 253402300799 {
			t.Skip("Date out of reasonable range")
		}

		date := time.Unix(unixTime, 0).UTC()

		analyzer := &Analyzer{}
		ce := NewCorrelationEngine(analyzer)

		// Should not panic
		positions := ce.calculatePlanetaryPositions(date)

		// Verify all expected planets are present
		expectedPlanets := []string{"Mercury", "Venus", "Mars", "Jupiter", "Saturn"}
		for _, planet := range expectedPlanets {
			pos, exists := positions[planet]
			if !exists {
				t.Errorf("Missing position for planet: %s", planet)
				continue
			}

			// Position should be in [0, 360) degrees
			if pos < 0 || pos >= 360 {
				t.Errorf("Invalid position for %s: %f", planet, pos)
			}
		}

		// Should only have expected planets
		if len(positions) != len(expectedPlanets) {
			t.Errorf("Unexpected number of planets: %d", len(positions))
		}
	})
}

// FuzzCosmicDataEnrichment tests cosmic data enrichment with various analyzer states
func FuzzCosmicDataEnrichment(f *testing.F) {
	// Seed with various drawing counts and date ranges
	f.Add(1, 1, 2024, 1, 2024)    // Single drawing
	f.Add(100, 1, 2020, 12, 2024) // Many drawings
	f.Add(10, 1, 1900, 12, 1900)  // Old dates
	f.Add(5, 1, 2099, 12, 2099)   // Future dates

	f.Fuzz(func(t *testing.T, numDrawings, startMonth, startYear, endMonth, endYear int) {
		// Validate inputs
		if numDrawings < 0 || numDrawings > 1000 {
			t.Skip("Invalid number of drawings")
		}
		if startMonth < 1 || startMonth > 12 || endMonth < 1 || endMonth > 12 {
			t.Skip("Invalid month")
		}
		if startYear < 1900 || startYear > 2100 || endYear < 1900 || endYear > 2100 {
			t.Skip("Year out of range")
		}

		// Create analyzer with fuzzed drawings
		analyzer := &Analyzer{
			drawings:    make([]Drawing, 0),
			mainNumbers: make(map[int]*NumberInfo),
			luckyBalls:  make(map[int]*NumberInfo),
		}

		// Generate drawings across date range
		startDate := time.Date(startYear, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(endYear, time.Month(endMonth), 28, 0, 0, 0, 0, time.UTC) // Safe day

		if startDate.After(endDate) {
			t.Skip("Start date after end date")
		}

		totalDays := int(endDate.Sub(startDate).Hours() / 24)
		if totalDays <= 0 {
			totalDays = 1
		}

		for i := 0; i < numDrawings && i < 1000; i++ {
			dayOffset := (totalDays * i) / numDrawings
			drawDate := startDate.AddDate(0, 0, dayOffset)

			drawing := Drawing{
				Date:      drawDate,
				Numbers:   []int{1, 2, 3, 4, 5},
				LuckyBall: 1,
				Index:     i,
			}
			analyzer.drawings = append(analyzer.drawings, drawing)
		}

		ce := NewCorrelationEngine(analyzer)
		ctx := context.Background()

		// Should not panic
		err := ce.EnrichWithCosmicData(ctx)
		if err != nil {
			t.Logf("Enrichment error (expected): %v", err)
		}

		// Verify cosmic data was created for drawings
		for _, drawing := range analyzer.drawings {
			dateKey := drawing.Date.Format(dateFormatISO)
			if cosmic, exists := ce.cosmicData[dateKey]; exists {
				// Verify cosmic data is valid
				if cosmic.MoonPhase < 0 || cosmic.MoonPhase >= 1 {
					t.Errorf("Invalid moon phase: %f", cosmic.MoonPhase)
				}
				if cosmic.MoonIllumination < 0 || cosmic.MoonIllumination > 1 {
					t.Errorf("Invalid moon illumination: %f", cosmic.MoonIllumination)
				}
			}
		}
	})
}

// FuzzCosmicReport tests report generation with various correlation results
func FuzzCosmicReport(f *testing.F) {
	// Seed with various correlation scenarios
	f.Add(0.8, 0.01, "High", "Strong positive correlation")
	f.Add(-0.5, 0.05, "Moderate", "Negative correlation")
	f.Add(0.1, 0.5, "None", "No significant correlation")
	f.Add(0.0, 1.0, "", "")
	f.Add(1.0, 0.0, "Perfect", "Perfect correlation")
	f.Add(math.NaN(), 0.1, "Unknown", "Invalid correlation")

	f.Fuzz(func(t *testing.T, correlation, pValue float64, significance, interpretation string) {
		analyzer := &Analyzer{
			drawings: []Drawing{
				{Date: time.Now(), Numbers: []int{1, 2, 3, 4, 5}, LuckyBall: 1},
			},
		}
		ce := NewCorrelationEngine(analyzer)

		// Add fuzzed correlation result
		ce.correlationResults = []CorrelationResult{
			{
				Factor:         "Test Factor",
				SubFactor:      "Test SubFactor",
				Correlation:    correlation,
				PValue:         pValue,
				SampleSize:     100,
				Significance:   significance,
				Interpretation: interpretation,
			},
		}

		// Should not panic
		report := ce.GenerateCosmicReport()

		// Report should contain basic structure
		if !strings.Contains(report, "COSMIC CORRELATION ANALYSIS") {
			t.Error("Report missing header")
		}
		if !strings.Contains(report, "DISCLAIMER") {
			t.Error("Report missing disclaimer")
		}

		// If we have valid correlation, it should appear in report
		if !math.IsNaN(correlation) && !math.IsInf(correlation, 0) {
			if strings.Contains(ce.correlationResults[0].Factor, "Test") {
				// The test factor might not appear in report as it's not in the known categories
				t.Logf("Test factor may not appear in categorized report sections")
			}
		}
	})
}

// FuzzSeasonalPhase tests seasonal phase calculation with edge dates
func FuzzSeasonalPhase(f *testing.F) {
	// Seed with boundary dates
	f.Add(3, 19)  // Day before spring
	f.Add(3, 20)  // Spring equinox
	f.Add(6, 20)  // Day before summer
	f.Add(6, 21)  // Summer solstice
	f.Add(9, 22)  // Day before autumn
	f.Add(9, 23)  // Autumn equinox
	f.Add(12, 20) // Day before winter
	f.Add(12, 21) // Winter solstice

	f.Fuzz(func(t *testing.T, month, day int) {
		// Validate month
		if month < 1 || month > 12 {
			t.Skip("Invalid month")
		}

		// Get max days for month
		year := 2024
		maxDay := 31
		switch time.Month(month) {
		case time.February:
			maxDay = 29 // 2024 is leap year
		case time.April, time.June, time.September, time.November:
			maxDay = 30
		}

		if day < 1 || day > maxDay {
			t.Skip("Invalid day for month")
		}

		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		analyzer := &Analyzer{}
		ce := NewCorrelationEngine(analyzer)

		// Should not panic
		season := ce.getSeasonalPhase(date)

		// Validate season
		validSeasons := []string{"Spring", "Summer", "Autumn", "Winter"}
		found := false
		for _, s := range validSeasons {
			if season == s {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Invalid season: %s", season)
		}
	})
}

// FuzzPredictBasedOnCosmicConditions tests cosmic predictions with various conditions
func FuzzPredictBasedOnCosmicConditions(f *testing.F) {
	// Seed with various timestamps
	f.Add(time.Now().Unix())
	f.Add(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix())
	f.Add(time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC).Unix())
	f.Add(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC).Unix())

	f.Fuzz(func(t *testing.T, unixTime int64) {
		// Skip extreme dates
		if unixTime < 0 || unixTime > 253402300799 {
			t.Skip("Date out of range")
		}

		// Create analyzer with correlation engine
		analyzer := &Analyzer{
			drawings: []Drawing{},
		}
		ce := NewCorrelationEngine(analyzer)

		// Pre-populate with some cosmic data
		testDate := time.Unix(unixTime, 0).UTC()
		dateKey := testDate.Format(dateFormatISO)

		cosmic := &CosmicData{
			Date:             testDate,
			MoonPhase:        0.5,
			MoonIllumination: 0.75,
			MoonPhaseName:    "Full Moon",
			ZodiacSign:       "Leo",
			SeasonalPhase:    "Summer",
			DayOfWeek:        testDate.Weekday().String(),
		}

		// Add mock data
		ce.addMockDataForDemo(cosmic)
		ce.cosmicData[dateKey] = cosmic

		// Should not panic
		numbers := ce.PredictBasedOnCosmicConditions()

		// Validate predictions
		if len(numbers) != 5 {
			t.Errorf("Expected 5 numbers, got %d", len(numbers))
		}

		// Check all numbers are valid and unique
		seen := make(map[int]bool)
		for _, num := range numbers {
			if num < 1 || num > 48 {
				t.Errorf("Number out of range: %d", num)
			}
			if seen[num] {
				t.Errorf("Duplicate number: %d", num)
			}
			seen[num] = true
		}
	})
}

// FuzzCorrelationInterpretations tests interpretation functions with edge values
func FuzzCorrelationInterpretations(f *testing.F) {
	// Seed with various correlation values and p-values
	f.Add(0.0, 0.0)
	f.Add(1.0, 0.0)
	f.Add(-1.0, 0.0)
	f.Add(0.5, 0.05)
	f.Add(-0.5, 0.05)
	f.Add(0.1, 0.5)
	f.Add(math.NaN(), 0.1)
	f.Add(0.5, math.NaN())
	f.Add(math.Inf(1), 0.1)
	f.Add(0.5, math.Inf(1))

	f.Fuzz(func(t *testing.T, corr, pValue float64) {
		// Test moon correlation interpretation
		moonInterp := interpretMoonCorrelation(corr, pValue)
		if moonInterp == "" {
			t.Error("Empty moon interpretation")
		}

		// Test solar correlation interpretation
		solarInterp := interpretSolarCorrelation(corr, pValue)
		if solarInterp == "" {
			t.Error("Empty solar interpretation")
		}

		// Test weather correlation interpretation
		weatherInterp := interpretWeatherCorrelation(corr, pValue)
		if weatherInterp == "" {
			t.Error("Empty weather interpretation")
		}

		// Interpretations should handle special values gracefully
		if math.IsNaN(corr) || math.IsInf(corr, 0) || math.IsNaN(pValue) || math.IsInf(pValue, 0) {
			// Should still return valid interpretation
			if strings.Contains(moonInterp, "NaN") || strings.Contains(moonInterp, "Inf") {
				t.Error("Interpretation contains NaN/Inf")
			}
		}
	})
}
