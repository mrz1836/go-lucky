package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"time"
)

// Drawing represents a single lottery drawing with its results
type Drawing struct {
	Date      time.Time `json:"date"`
	Numbers   []int     `json:"numbers"`
	LuckyBall int       `json:"lucky_ball"`
	Index     int       `json:"index"` // Position in dataset (0 = most recent)
}

// NumberInfo contains comprehensive statistics for a single number
type NumberInfo struct {
	Number             int       `json:"number"`
	TotalFrequency     int       `json:"total_frequency"`
	RecentFrequency    int       `json:"recent_frequency"` // Last 50 drawings
	LastDrawnIndex     int       `json:"last_drawn_index"`
	LastDrawnDate      time.Time `json:"last_drawn_date"`
	GapsSinceDrawn     []int     `json:"gaps_since_drawn"`
	AverageGap         float64   `json:"average_gap"`
	StandardDeviation  float64   `json:"standard_deviation"`
	CurrentGap         int       `json:"current_gap"`
	ExpectedFrequency  float64   `json:"expected_frequency"`
	ChiSquareComponent float64   `json:"chi_square_component"`
}

// CombinationPattern represents a specific number combination and its frequency
type CombinationPattern struct {
	Numbers   []int  `json:"numbers"`
	Key       string `json:"key"`
	Frequency int    `json:"frequency"`
	LastSeen  int    `json:"last_seen"` // Index when last seen
}

// PatternStats tracks various pattern occurrences
type PatternStats struct {
	OddEvenPatterns map[string]int `json:"odd_even_patterns"`
	SumRanges       map[int]int    `json:"sum_ranges"`
	ConsecutiveCount int           `json:"consecutive_count"`
	DecadeDistribution map[int]int `json:"decade_distribution"`
}

// ScoredNumber represents a number with its calculated score and reasoning
type ScoredNumber struct {
	Number  int      `json:"number"`
	Score   float64  `json:"score"`
	Factors []string `json:"factors"`
}

// RecommendedSet represents a suggested number combination with metadata
type RecommendedSet struct {
	Numbers      []int    `json:"numbers"`
	LuckyBall    int      `json:"lucky_ball"`
	Strategy     string   `json:"strategy"`
	Confidence   float64  `json:"confidence"`
	Explanation  string   `json:"explanation"`
}

// AnalysisConfig holds configuration for analysis parameters
type AnalysisConfig struct {
	RecentWindow      int     `json:"recent_window"`       // How many drawings to consider "recent"
	MinGapMultiplier  float64 `json:"min_gap_multiplier"`  // Multiplier for "overdue" threshold
	ConfidenceLevel   float64 `json:"confidence_level"`    // Statistical confidence level
	OutputMode        string  `json:"output_mode"`         // "simple", "detailed", "statistical"
	ExportFormat      string  `json:"export_format"`       // "console", "csv", "json"
}

// Analyzer is the main lottery analysis engine
type Analyzer struct {
	config            *AnalysisConfig
	drawings          []Drawing
	mainNumbers       map[int]*NumberInfo
	luckyBalls        map[int]*NumberInfo
	pairPatterns      map[string]*CombinationPattern
	triplePatterns    map[string]*CombinationPattern
	quadPatterns      map[string]*CombinationPattern
	patternStats      *PatternStats
	chiSquareValue    float64
	randomnessScore   float64
	correlationEngine *CorrelationEngine
}

// NewAnalyzer creates a new analyzer instance with the given configuration
func NewAnalyzer(ctx context.Context, filename string, config *AnalysisConfig) (*Analyzer, error) {
	if config == nil {
		config = &AnalysisConfig{
			RecentWindow:     50,
			MinGapMultiplier: 1.5,
			ConfidenceLevel:  0.95,
			OutputMode:       "detailed",
			ExportFormat:     "console",
		}
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	analyzer := &Analyzer{
		config:         config,
		drawings:       make([]Drawing, 0),
		mainNumbers:    make(map[int]*NumberInfo),
		luckyBalls:     make(map[int]*NumberInfo),
		pairPatterns:   make(map[string]*CombinationPattern),
		triplePatterns: make(map[string]*CombinationPattern),
		quadPatterns:   make(map[string]*CombinationPattern),
		patternStats: &PatternStats{
			OddEvenPatterns:    make(map[string]int),
			SumRanges:          make(map[int]int),
			DecadeDistribution: make(map[int]int),
		},
	}

	// Initialize number tracking
	for i := 1; i <= 48; i++ {
		analyzer.mainNumbers[i] = &NumberInfo{
			Number:         i,
			GapsSinceDrawn: []int{},
			LastDrawnIndex: -1,
		}
	}
	for i := 1; i <= 18; i++ {
		analyzer.luckyBalls[i] = &NumberInfo{
			Number:         i,
			GapsSinceDrawn: []int{},
			LastDrawnIndex: -1,
		}
	}

	// Parse CSV data
	if err := analyzer.parseDrawings(ctx, records); err != nil {
		return nil, fmt.Errorf("failed to parse drawings: %w", err)
	}

	// Perform comprehensive analysis
	if err := analyzer.analyzeData(ctx); err != nil {
		return nil, fmt.Errorf("failed to analyze data: %w", err)
	}

	// Initialize correlation engine
	analyzer.correlationEngine = NewCorrelationEngine(analyzer)

	return analyzer, nil
}

// parseDrawings processes the CSV records into Drawing structs
func (a *Analyzer) parseDrawings(ctx context.Context, records [][]string) error {
	// Skip header row
	for i := 1; i < len(records); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if len(records[i]) < 7 || records[i][0] == "" {
			continue
		}

		date, err := time.Parse("01/02/2006", records[i][0])
		if err != nil {
			continue
		}

		drawing := Drawing{
			Date:    date,
			Numbers: make([]int, 5),
			Index:   len(a.drawings), // 0 = most recent
		}

		// Parse main numbers
		validDrawing := true
		for j := 1; j <= 5; j++ {
			num, err := strconv.Atoi(records[i][j])
			if err != nil {
				validDrawing = false
				break
			}
			drawing.Numbers[j-1] = num
		}

		if !validDrawing {
			continue
		}

		// Parse lucky ball
		luckyBall, err := strconv.Atoi(records[i][6])
		if err != nil {
			continue
		}
		drawing.LuckyBall = luckyBall

		a.drawings = append(a.drawings, drawing)
	}

	// Reverse to have newest first
	for i, j := 0, len(a.drawings)-1; i < j; i, j = i+1, j-1 {
		a.drawings[i], a.drawings[j] = a.drawings[j], a.drawings[i]
	}

	// Update indices after reversal
	for i := range a.drawings {
		a.drawings[i].Index = i
	}

	return nil
}

// analyzeData performs comprehensive analysis on the parsed drawings
func (a *Analyzer) analyzeData(ctx context.Context) error {
	// Process each drawing
	for idx, drawing := range a.drawings {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Track main numbers
		for _, num := range drawing.Numbers {
			if err := a.updateNumberInfo(a.mainNumbers[num], idx, drawing.Date); err != nil {
				return err
			}

			// Track recent frequency
			if idx < a.config.RecentWindow {
				a.mainNumbers[num].RecentFrequency++
			}
		}

		// Track lucky ball
		lbInfo := a.luckyBalls[drawing.LuckyBall]
		if err := a.updateNumberInfo(lbInfo, idx, drawing.Date); err != nil {
			return err
		}
		if idx < a.config.RecentWindow {
			lbInfo.RecentFrequency++
		}

		// Analyze patterns
		a.analyzeCombinations(drawing.Numbers, idx)
		a.analyzePatterns(drawing)
	}

	// Calculate statistical measures
	a.calculateStatistics()
	a.calculateChiSquare()

	return nil
}

// updateNumberInfo updates frequency and gap information for a number
func (a *Analyzer) updateNumberInfo(info *NumberInfo, idx int, date time.Time) error {
	info.TotalFrequency++

	if info.LastDrawnIndex != -1 {
		gap := idx - info.LastDrawnIndex
		info.GapsSinceDrawn = append(info.GapsSinceDrawn, gap)
	}

	info.LastDrawnIndex = idx
	info.LastDrawnDate = date

	return nil
}

// analyzeCombinations tracks pair, triple, and quad patterns
func (a *Analyzer) analyzeCombinations(numbers []int, drawingIndex int) {
	sorted := make([]int, len(numbers))
	copy(sorted, numbers)
	sort.Ints(sorted)

	// Pairs
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			key := fmt.Sprintf("%d-%d", sorted[i], sorted[j])
			if pattern, exists := a.pairPatterns[key]; exists {
				pattern.Frequency++
				pattern.LastSeen = drawingIndex
			} else {
				a.pairPatterns[key] = &CombinationPattern{
					Numbers:   []int{sorted[i], sorted[j]},
					Key:       key,
					Frequency: 1,
					LastSeen:  drawingIndex,
				}
			}
		}
	}

	// Triples
	for i := 0; i < len(sorted)-2; i++ {
		for j := i + 1; j < len(sorted)-1; j++ {
			for k := j + 1; k < len(sorted); k++ {
				key := fmt.Sprintf("%d-%d-%d", sorted[i], sorted[j], sorted[k])
				if pattern, exists := a.triplePatterns[key]; exists {
					pattern.Frequency++
					pattern.LastSeen = drawingIndex
				} else {
					a.triplePatterns[key] = &CombinationPattern{
						Numbers:   []int{sorted[i], sorted[j], sorted[k]},
						Key:       key,
						Frequency: 1,
						LastSeen:  drawingIndex,
					}
				}
			}
		}
	}

	// Quads
	for i := 0; i < len(sorted)-3; i++ {
		for j := i + 1; j < len(sorted)-2; j++ {
			for k := j + 1; k < len(sorted)-1; k++ {
				for l := k + 1; l < len(sorted); l++ {
					key := fmt.Sprintf("%d-%d-%d-%d", sorted[i], sorted[j], sorted[k], sorted[l])
					if pattern, exists := a.quadPatterns[key]; exists {
						pattern.Frequency++
						pattern.LastSeen = drawingIndex
					} else {
						a.quadPatterns[key] = &CombinationPattern{
							Numbers:   []int{sorted[i], sorted[j], sorted[k], sorted[l]},
							Key:       key,
							Frequency: 1,
							LastSeen:  drawingIndex,
						}
					}
				}
			}
		}
	}
}

// analyzePatterns tracks various statistical patterns
func (a *Analyzer) analyzePatterns(drawing Drawing) {
	odd, even := 0, 0
	sum := 0
	hasConsecutive := false

	sorted := make([]int, len(drawing.Numbers))
	copy(sorted, drawing.Numbers)
	sort.Ints(sorted)

	for i, num := range sorted {
		sum += num
		
		// Odd/Even count
		if num%2 == 1 {
			odd++
		} else {
			even++
		}

		// Decade distribution
		decade := (num - 1) / 10
		a.patternStats.DecadeDistribution[decade]++

		// Check for consecutive numbers
		if i > 0 && sorted[i]-sorted[i-1] == 1 {
			hasConsecutive = true
		}
	}

	// Record patterns
	oddEvenPattern := fmt.Sprintf("%dO-%dE", odd, even)
	a.patternStats.OddEvenPatterns[oddEvenPattern]++

	sumRange := (sum / 20) * 20
	a.patternStats.SumRanges[sumRange]++

	if hasConsecutive {
		a.patternStats.ConsecutiveCount++
	}
}

// calculateStatistics computes averages, standard deviations, and gaps
func (a *Analyzer) calculateStatistics() {
	// Calculate for main numbers
	for _, info := range a.mainNumbers {
		if len(info.GapsSinceDrawn) > 0 {
			// Average gap
			sum := 0
			for _, gap := range info.GapsSinceDrawn {
				sum += gap
			}
			info.AverageGap = float64(sum) / float64(len(info.GapsSinceDrawn))

			// Standard deviation
			var variance float64
			for _, gap := range info.GapsSinceDrawn {
				diff := float64(gap) - info.AverageGap
				variance += diff * diff
			}
			info.StandardDeviation = math.Sqrt(variance / float64(len(info.GapsSinceDrawn)))
		}
		info.CurrentGap = info.LastDrawnIndex

		// Expected frequency (assuming uniform distribution)
		info.ExpectedFrequency = float64(len(a.drawings)) * 5 / 48
	}

	// Calculate for lucky balls
	for _, info := range a.luckyBalls {
		if len(info.GapsSinceDrawn) > 0 {
			sum := 0
			for _, gap := range info.GapsSinceDrawn {
				sum += gap
			}
			info.AverageGap = float64(sum) / float64(len(info.GapsSinceDrawn))

			// Standard deviation
			var variance float64
			for _, gap := range info.GapsSinceDrawn {
				diff := float64(gap) - info.AverageGap
				variance += diff * diff
			}
			info.StandardDeviation = math.Sqrt(variance / float64(len(info.GapsSinceDrawn)))
		}
		info.CurrentGap = info.LastDrawnIndex
		info.ExpectedFrequency = float64(len(a.drawings)) / 18
	}
}

// calculateChiSquare performs chi-square test for randomness
func (a *Analyzer) calculateChiSquare() {
	var chiSquareMain, chiSquareLucky float64

	// Chi-square for main numbers
	for _, info := range a.mainNumbers {
		if info.ExpectedFrequency > 0 {
			diff := float64(info.TotalFrequency) - info.ExpectedFrequency
			info.ChiSquareComponent = (diff * diff) / info.ExpectedFrequency
			chiSquareMain += info.ChiSquareComponent
		}
	}

	// Chi-square for lucky balls
	for _, info := range a.luckyBalls {
		if info.ExpectedFrequency > 0 {
			diff := float64(info.TotalFrequency) - info.ExpectedFrequency
			info.ChiSquareComponent = (diff * diff) / info.ExpectedFrequency
			chiSquareLucky += info.ChiSquareComponent
		}
	}

	a.chiSquareValue = chiSquareMain + chiSquareLucky

	// Calculate randomness score (0-100, where 100 is perfectly random)
	// Using chi-square critical values for 95% confidence
	mainCritical := 64.001 // df=47, alpha=0.05
	luckyCritical := 27.587 // df=17, alpha=0.05

	mainRandomness := 100.0 * (1 - math.Min(chiSquareMain/mainCritical, 1))
	luckyRandomness := 100.0 * (1 - math.Min(chiSquareLucky/luckyCritical, 1))

	a.randomnessScore = (mainRandomness + luckyRandomness) / 2
}

// GetTopNumbers returns the most frequent numbers
func (a *Analyzer) GetTopNumbers(count int, recent bool) []*NumberInfo {
	numbers := make([]*NumberInfo, 0, len(a.mainNumbers))
	for _, info := range a.mainNumbers {
		numbers = append(numbers, info)
	}

	sort.Slice(numbers, func(i, j int) bool {
		if recent {
			return numbers[i].RecentFrequency > numbers[j].RecentFrequency
		}
		return numbers[i].TotalFrequency > numbers[j].TotalFrequency
	})

	if count > len(numbers) {
		count = len(numbers)
	}
	return numbers[:count]
}

// GetOverdueNumbers returns numbers that haven't been drawn recently
func (a *Analyzer) GetOverdueNumbers(count int) []*NumberInfo {
	overdue := make([]*NumberInfo, 0)
	
	for _, info := range a.mainNumbers {
		if info.AverageGap > 0 && float64(info.CurrentGap) > info.AverageGap*a.config.MinGapMultiplier {
			overdue = append(overdue, info)
		}
	}

	sort.Slice(overdue, func(i, j int) bool {
		ratioI := float64(overdue[i].CurrentGap) / overdue[i].AverageGap
		ratioJ := float64(overdue[j].CurrentGap) / overdue[j].AverageGap
		return ratioI > ratioJ
	})

	if count > len(overdue) {
		count = len(overdue)
	}
	return overdue[:count]
}

// GenerateRecommendations creates multiple number sets with different strategies
func (a *Analyzer) GenerateRecommendations(ctx context.Context, count int) ([]RecommendedSet, error) {
	recommendations := make([]RecommendedSet, 0, count)

	strategies := []string{"balanced", "hot", "overdue", "pattern", "frequency"}
	
	for i := 0; i < count && i < len(strategies); i++ {
		select {
		case <-ctx.Done():
			return recommendations, ctx.Err()
		default:
		}

		set, err := a.generateSetByStrategy(strategies[i])
		if err != nil {
			continue
		}
		recommendations = append(recommendations, set)
	}

	return recommendations, nil
}

// generateSetByStrategy creates a number set based on a specific strategy
func (a *Analyzer) generateSetByStrategy(strategy string) (RecommendedSet, error) {
	set := RecommendedSet{
		Strategy: strategy,
		Numbers:  make([]int, 0, 5),
	}

	// Score all numbers based on strategy
	scoredNumbers := a.scoreNumbersByStrategy(strategy)
	
	// Select 5 numbers ensuring no duplicates
	used := make(map[int]bool)
	for _, sn := range scoredNumbers {
		if !used[sn.Number] && len(set.Numbers) < 5 {
			set.Numbers = append(set.Numbers, sn.Number)
			used[sn.Number] = true
		}
	}

	sort.Ints(set.Numbers)

	// Select lucky ball
	luckyScores := a.scoreLuckyBalls()
	if len(luckyScores) > 0 {
		set.LuckyBall = luckyScores[0].Number
	}

	// Calculate confidence based on randomness score and strategy
	set.Confidence = a.calculateConfidence(strategy)
	set.Explanation = a.generateExplanation(strategy, set)

	return set, nil
}

// scoreNumbersByStrategy scores numbers based on the chosen strategy
func (a *Analyzer) scoreNumbersByStrategy(strategy string) []ScoredNumber {
	scoredNumbers := make([]ScoredNumber, 0, len(a.mainNumbers))

	for num, info := range a.mainNumbers {
		score := 0.0
		factors := []string{}

		switch strategy {
		case "balanced":
			// Mix of frequency, recency, and overdue
			freqScore := float64(info.TotalFrequency) / float64(len(a.drawings)) * 100
			score += freqScore

			if info.RecentFrequency > 3 {
				score += float64(info.RecentFrequency) * 10
				factors = append(factors, fmt.Sprintf("Hot-%d", info.RecentFrequency))
			}

			if info.AverageGap > 0 && float64(info.CurrentGap) > info.AverageGap*1.3 {
				overdueRatio := float64(info.CurrentGap) / info.AverageGap
				score += overdueRatio * 20
				factors = append(factors, fmt.Sprintf("Overdue-%.1fx", overdueRatio))
			}

		case "hot":
			// Focus on recent frequency
			score = float64(info.RecentFrequency) * 100
			if info.RecentFrequency > 0 {
				factors = append(factors, fmt.Sprintf("Recent-%d", info.RecentFrequency))
			}

		case "overdue":
			// Focus on numbers that haven't appeared recently
			if info.AverageGap > 0 {
				overdueRatio := float64(info.CurrentGap) / info.AverageGap
				score = overdueRatio * 100
				factors = append(factors, fmt.Sprintf("Gap-%d-days", info.CurrentGap))
			}

		case "pattern":
			// Look for numbers that appear in common patterns
			pairBonus := 0
			for _, pattern := range a.pairPatterns {
				for _, pNum := range pattern.Numbers {
					if pNum == num {
						pairBonus += pattern.Frequency
					}
				}
			}
			score = float64(pairBonus)
			if pairBonus > 50 {
				factors = append(factors, "StrongPairs")
			}

		case "frequency":
			// Pure frequency-based selection
			score = float64(info.TotalFrequency)
			factors = append(factors, fmt.Sprintf("Freq-%d", info.TotalFrequency))
		}

		scoredNumbers = append(scoredNumbers, ScoredNumber{
			Number:  num,
			Score:   score,
			Factors: factors,
		})
	}

	sort.Slice(scoredNumbers, func(i, j int) bool {
		return scoredNumbers[i].Score > scoredNumbers[j].Score
	})

	return scoredNumbers
}

// scoreLuckyBalls scores lucky ball numbers
func (a *Analyzer) scoreLuckyBalls() []ScoredNumber {
	scoredNumbers := make([]ScoredNumber, 0, len(a.luckyBalls))

	for num, info := range a.luckyBalls {
		score := float64(info.TotalFrequency) + float64(info.RecentFrequency*5)
		factors := []string{
			fmt.Sprintf("Total-%d", info.TotalFrequency),
			fmt.Sprintf("Recent-%d", info.RecentFrequency),
		}

		scoredNumbers = append(scoredNumbers, ScoredNumber{
			Number:  num,
			Score:   score,
			Factors: factors,
		})
	}

	sort.Slice(scoredNumbers, func(i, j int) bool {
		return scoredNumbers[i].Score > scoredNumbers[j].Score
	})

	return scoredNumbers
}

// calculateConfidence calculates confidence score for a strategy
func (a *Analyzer) calculateConfidence(strategy string) float64 {
	// Base confidence on randomness score
	baseConfidence := a.randomnessScore / 100.0

	// Adjust based on strategy
	switch strategy {
	case "balanced":
		return baseConfidence * 0.95 // Most reliable strategy
	case "hot":
		return baseConfidence * 0.85 // Recent trends may not continue
	case "overdue":
		return baseConfidence * 0.80 // Gambler's fallacy risk
	case "pattern":
		return baseConfidence * 0.75 // Patterns in random data are coincidental
	case "frequency":
		return baseConfidence * 0.90 // Long-term frequency is more stable
	default:
		return baseConfidence * 0.70
	}
}

// generateExplanation creates a human-readable explanation for the recommendation
func (a *Analyzer) generateExplanation(strategy string, set RecommendedSet) string {
	switch strategy {
	case "balanced":
		return "Combines hot numbers, overdue numbers, and frequency analysis for a well-rounded selection"
	case "hot":
		return "Focuses on numbers that have appeared frequently in recent drawings"
	case "overdue":
		return "Selects numbers that haven't appeared for longer than their average gap"
	case "pattern":
		return "Based on numbers that frequently appear together in winning combinations"
	case "frequency":
		return "Selects the most frequently drawn numbers throughout the entire history"
	default:
		return "Custom strategy based on statistical analysis"
	}
}

// ExportAnalysis exports the analysis results in the specified format
func (a *Analyzer) ExportAnalysis(ctx context.Context, filename string) error {
	switch a.config.ExportFormat {
	case "json":
		return a.exportJSON(ctx, filename)
	case "csv":
		return a.exportCSV(ctx, filename)
	default:
		return fmt.Errorf("unsupported export format: %s", a.config.ExportFormat)
	}
}

// exportJSON exports analysis results as JSON
func (a *Analyzer) exportJSON(ctx context.Context, filename string) error {
	data := map[string]interface{}{
		"metadata": map[string]interface{}{
			"total_drawings":   len(a.drawings),
			"date_range":       fmt.Sprintf("%s to %s", a.drawings[len(a.drawings)-1].Date.Format("01/02/2006"), a.drawings[0].Date.Format("01/02/2006")),
			"randomness_score": a.randomnessScore,
			"chi_square":       a.chiSquareValue,
		},
		"main_numbers": a.mainNumbers,
		"lucky_balls":  a.luckyBalls,
		"patterns": map[string]interface{}{
			"odd_even":    a.patternStats.OddEvenPatterns,
			"sum_ranges":  a.patternStats.SumRanges,
			"consecutive": a.patternStats.ConsecutiveCount,
		},
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// exportCSV exports analysis results as CSV
func (a *Analyzer) exportCSV(ctx context.Context, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	header := []string{"Number", "Total Frequency", "Recent Frequency", "Average Gap", "Current Gap", "Last Drawn", "Chi-Square Component"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Main numbers
	for i := 1; i <= 48; i++ {
		info := a.mainNumbers[i]
		record := []string{
			strconv.Itoa(info.Number),
			strconv.Itoa(info.TotalFrequency),
			strconv.Itoa(info.RecentFrequency),
			fmt.Sprintf("%.2f", info.AverageGap),
			strconv.Itoa(info.CurrentGap),
			info.LastDrawnDate.Format("01/02/2006"),
			fmt.Sprintf("%.4f", info.ChiSquareComponent),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// RunAnalysis performs complete analysis and outputs results
func (a *Analyzer) RunAnalysis(ctx context.Context) error {
	// Perform cosmic correlation analysis
	if err := a.correlationEngine.EnrichWithCosmicData(ctx); err != nil {
		fmt.Printf("Warning: Could not enrich with cosmic data: %v\n", err)
	}
	
	if err := a.correlationEngine.AnalyzeCorrelations(ctx); err != nil {
		fmt.Printf("Warning: Could not analyze correlations: %v\n", err)
	}

	switch a.config.OutputMode {
	case "simple":
		return a.printSimpleAnalysis(ctx)
	case "statistical":
		return a.printStatisticalAnalysis(ctx)
	case "cosmic":
		return a.printCosmicAnalysis(ctx)
	default:
		return a.printDetailedAnalysis(ctx)
	}
}

// printDetailedAnalysis outputs comprehensive analysis results
func (a *Analyzer) printDetailedAnalysis(ctx context.Context) error {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║        NC LUCKY FOR LIFE LOTTERY ANALYZER               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	
	// Metadata
	fmt.Printf("\nTotal Drawings Analyzed: %d\n", len(a.drawings))
	fmt.Printf("Date Range: %s to %s\n", 
		a.drawings[len(a.drawings)-1].Date.Format("01/02/2006"),
		a.drawings[0].Date.Format("01/02/2006"))
	fmt.Printf("Randomness Score: %.1f%% (100%% = perfectly random)\n", a.randomnessScore)
	fmt.Printf("Chi-Square Value: %.2f\n", a.chiSquareValue)

	// Frequency Analysis
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("                    FREQUENCY ANALYSIS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	hotNumbers := a.GetTopNumbers(10, true)
	fmt.Println("\nHOT NUMBERS (Last 50 Drawings):")
	for i, info := range hotNumbers {
		fmt.Printf("  %2d. Number %2d: %d times (%.1f%%) | Total: %d\n",
			i+1, info.Number, info.RecentFrequency,
			float64(info.RecentFrequency)/float64(a.config.RecentWindow)*100,
			info.TotalFrequency)
	}

	topNumbers := a.GetTopNumbers(10, false)
	fmt.Println("\nMOST FREQUENT (All Time):")
	for i, info := range topNumbers {
		deviation := float64(info.TotalFrequency) - info.ExpectedFrequency
		fmt.Printf("  %2d. Number %2d: %d times (%.1f%% deviation from expected)\n",
			i+1, info.Number, info.TotalFrequency,
			(deviation/info.ExpectedFrequency)*100)
	}

	// Overdue Analysis
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("                    OVERDUE ANALYSIS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	overdueNumbers := a.GetOverdueNumbers(10)
	fmt.Println("\nMOST OVERDUE NUMBERS:")
	for i, info := range overdueNumbers {
		overdueRatio := float64(info.CurrentGap) / info.AverageGap
		daysAgo := int(a.drawings[0].Date.Sub(info.LastDrawnDate).Hours() / 24)
		fmt.Printf("  %2d. Number %2d: Not drawn for %d drawings (%.1fx overdue) | %d days ago\n",
			i+1, info.Number, info.CurrentGap, overdueRatio, daysAgo)
	}

	// Pattern Analysis
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("                    PATTERN ANALYSIS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Odd/Even patterns
	fmt.Println("\nODD/EVEN DISTRIBUTION:")
	oddEvenSorted := make([]struct{ Pattern string; Count int }, 0)
	for pattern, count := range a.patternStats.OddEvenPatterns {
		oddEvenSorted = append(oddEvenSorted, struct{ Pattern string; Count int }{pattern, count})
	}
	sort.Slice(oddEvenSorted, func(i, j int) bool {
		return oddEvenSorted[i].Count > oddEvenSorted[j].Count
	})
	for i := 0; i < 3 && i < len(oddEvenSorted); i++ {
		percentage := float64(oddEvenSorted[i].Count) / float64(len(a.drawings)) * 100
		fmt.Printf("  %s: %d times (%.1f%%)\n", oddEvenSorted[i].Pattern, oddEvenSorted[i].Count, percentage)
	}

	consecutivePercent := float64(a.patternStats.ConsecutiveCount) / float64(len(a.drawings)) * 100
	fmt.Printf("\nConsecutive Numbers: %d drawings (%.1f%%)\n", a.patternStats.ConsecutiveCount, consecutivePercent)

	// Combination patterns
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("                 COMBINATION PATTERNS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Top pairs
	pairList := make([]*CombinationPattern, 0, len(a.pairPatterns))
	for _, pattern := range a.pairPatterns {
		pairList = append(pairList, pattern)
	}
	sort.Slice(pairList, func(i, j int) bool {
		return pairList[i].Frequency > pairList[j].Frequency
	})

	fmt.Println("\nTOP PAIRS:")
	for i := 0; i < 5 && i < len(pairList); i++ {
		fmt.Printf("  %s: %d times\n", pairList[i].Key, pairList[i].Frequency)
	}

	// Recommendations
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("                    RECOMMENDATIONS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	recommendations, err := a.GenerateRecommendations(ctx, 5)
	if err != nil {
		return fmt.Errorf("failed to generate recommendations: %w", err)
	}

	fmt.Println("\nRECOMMENDED NUMBER SETS:")
	for i, rec := range recommendations {
		fmt.Printf("\nSet %d - %s Strategy (%.1f%% confidence):\n", i+1, rec.Strategy, rec.Confidence*100)
		fmt.Printf("  Numbers: ")
		for j, num := range rec.Numbers {
			if j > 0 {
				fmt.Printf("-%02d", num)
			} else {
				fmt.Printf("%02d", num)
			}
		}
		fmt.Printf("  Lucky Ball: %d\n", rec.LuckyBall)
		fmt.Printf("  %s\n", rec.Explanation)
	}

	// Add cosmic correlation report
	fmt.Print(a.correlationEngine.GenerateCosmicReport())

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("                     DISCLAIMER")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("These recommendations are based on historical pattern analysis.")
	fmt.Printf("The randomness score of %.1f%% indicates the drawings are ", a.randomnessScore)
	if a.randomnessScore > 90 {
		fmt.Println("highly random.")
	} else if a.randomnessScore > 70 {
		fmt.Println("mostly random with minor deviations.")
	} else {
		fmt.Println("showing some non-random patterns.")
	}
	fmt.Println("Lottery drawings are designed to be random events.")
	fmt.Println("Past results do not influence future outcomes.")
	fmt.Println("Play responsibly!")

	return nil
}

// printSimpleAnalysis outputs a simplified analysis
func (a *Analyzer) printSimpleAnalysis(ctx context.Context) error {
	fmt.Println("LOTTERY ANALYSIS SUMMARY")
	fmt.Println("========================")
	
	fmt.Printf("Drawings analyzed: %d\n", len(a.drawings))
	fmt.Printf("Randomness: %.1f%%\n\n", a.randomnessScore)

	fmt.Println("TOP 5 HOT NUMBERS:")
	hotNumbers := a.GetTopNumbers(5, true)
	for _, info := range hotNumbers {
		fmt.Printf("  %2d (recent: %d times)\n", info.Number, info.RecentFrequency)
	}

	fmt.Println("\nTOP 5 OVERDUE:")
	overdueNumbers := a.GetOverdueNumbers(5)
	for _, info := range overdueNumbers {
		fmt.Printf("  %2d (gap: %d drawings)\n", info.Number, info.CurrentGap)
	}

	fmt.Println("\nQUICK PICKS:")
	recommendations, _ := a.GenerateRecommendations(ctx, 3)
	for i, rec := range recommendations {
		fmt.Printf("  Set %d: ", i+1)
		for j, num := range rec.Numbers {
			if j > 0 {
				fmt.Printf("-%02d", num)
			} else {
				fmt.Printf("%02d", num)
			}
		}
		fmt.Printf(" LB:%d\n", rec.LuckyBall)
	}

	// Add cosmic pick
	fmt.Println("\n🌌 COSMIC PICK:")
	cosmicNumbers := a.correlationEngine.PredictBasedOnCosmicConditions()
	fmt.Printf("  ")
	for i, num := range cosmicNumbers {
		if i > 0 {
			fmt.Printf("-%02d", num)
		} else {
			fmt.Printf("%02d", num)
		}
	}
	fmt.Printf(" LB:11\n")

	return nil
}

// printStatisticalAnalysis outputs detailed statistical analysis
func (a *Analyzer) printStatisticalAnalysis(ctx context.Context) error {
	fmt.Println("STATISTICAL ANALYSIS REPORT")
	fmt.Println("===========================")
	
	// Chi-square analysis
	fmt.Printf("\nChi-Square Test for Randomness:\n")
	fmt.Printf("  Total Chi-Square Value: %.4f\n", a.chiSquareValue)
	fmt.Printf("  Degrees of Freedom: %d (main) + %d (lucky)\n", 47, 17)
	fmt.Printf("  Randomness Score: %.2f%%\n", a.randomnessScore)
	
	// Distribution analysis
	fmt.Println("\nFrequency Distribution Analysis:")
	
	// Calculate mean and std dev for main numbers
	var sumFreq, sumSquaredDiff float64
	meanFreq := float64(len(a.drawings)) * 5 / 48
	
	for _, info := range a.mainNumbers {
		sumFreq += float64(info.TotalFrequency)
		diff := float64(info.TotalFrequency) - meanFreq
		sumSquaredDiff += diff * diff
	}
	
	stdDev := math.Sqrt(sumSquaredDiff / 48)
	fmt.Printf("  Expected frequency per number: %.2f\n", meanFreq)
	fmt.Printf("  Standard deviation: %.2f\n", stdDev)
	fmt.Printf("  Coefficient of variation: %.2f%%\n", (stdDev/meanFreq)*100)
	
	// Numbers outside normal range
	outsideCount := 0
	for _, info := range a.mainNumbers {
		if math.Abs(float64(info.TotalFrequency)-meanFreq) > 2*stdDev {
			outsideCount++
		}
	}
	fmt.Printf("  Numbers outside 2σ: %d (%.1f%%)\n", outsideCount, float64(outsideCount)/48*100)
	
	// Gap analysis
	fmt.Println("\nGap Analysis Statistics:")
	var totalGaps, minGap, maxGap int
	minGap = 999999
	
	for _, info := range a.mainNumbers {
		for _, gap := range info.GapsSinceDrawn {
			totalGaps++
			if gap < minGap {
				minGap = gap
			}
			if gap > maxGap {
				maxGap = gap
			}
		}
	}
	
	avgGap := float64(totalGaps) / float64(len(a.mainNumbers))
	fmt.Printf("  Average gap length: %.2f drawings\n", avgGap)
	fmt.Printf("  Minimum gap: %d drawings\n", minGap)
	fmt.Printf("  Maximum gap: %d drawings\n", maxGap)
	
	return nil
}

// printCosmicAnalysis outputs cosmic correlation analysis
func (a *Analyzer) printCosmicAnalysis(ctx context.Context) error {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║        COSMIC LOTTERY CORRELATION ANALYZER               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	
	fmt.Printf("\nTotal Drawings Analyzed: %d\n", len(a.drawings))
	fmt.Printf("Date Range: %s to %s\n", 
		a.drawings[len(a.drawings)-1].Date.Format("01/02/2006"),
		a.drawings[0].Date.Format("01/02/2006"))

	// Show cosmic correlation report
	fmt.Print(a.correlationEngine.GenerateCosmicReport())
	
	// Generate cosmic-influenced recommendations
	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("              🎯 COSMIC-INFLUENCED PREDICTIONS")
	fmt.Println("═══════════════════════════════════════════════════════════")
	
	fmt.Println("\nBased on current cosmic conditions:")
	cosmicNumbers := a.correlationEngine.PredictBasedOnCosmicConditions()
	fmt.Printf("\n🌟 Cosmic Selection: ")
	for i, num := range cosmicNumbers {
		if i > 0 {
			fmt.Printf("-%02d", num)
		} else {
			fmt.Printf("%02d", num)
		}
	}
	fmt.Printf("  Lucky Ball: 11\n")
	
	fmt.Println("\n📊 Combined Statistical + Cosmic Picks:")
	recommendations, _ := a.GenerateRecommendations(ctx, 3)
	for i, rec := range recommendations {
		fmt.Printf("\nSet %d - %s + Cosmic Alignment:\n", i+1, rec.Strategy)
		fmt.Printf("  Numbers: ")
		// Mix in some cosmic influence
		for j, num := range rec.Numbers {
			if j > 0 {
				fmt.Printf("-%02d", num)
			} else {
				fmt.Printf("%02d", num)
			}
		}
		fmt.Printf("  Lucky Ball: %d\n", rec.LuckyBall)
		fmt.Printf("  Confidence: %.1f%% (cosmic adjusted)\n", rec.Confidence*100*1.1)
	}
	
	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("                      COSMIC WISDOM")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("\n🌙 'As above, so below' - but lottery balls don't look up!")
	fmt.Println("☀️  The sun has witnessed every drawing, yet keeps its secrets.")
	fmt.Println("✨ Remember: The universe is under no obligation to make sense.")
	fmt.Println("🎲 ...or to make you wealthy!")
	
	return nil
}

// main is the entry point for the lottery analyzer
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
	analyzer, err := NewAnalyzer(ctx, "lucky-numbers-history.csv", config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Run analysis
	if err := analyzer.RunAnalysis(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error running analysis: %v\n", err)
		os.Exit(1)
	}

	// Export if requested
	if config.ExportFormat != "console" {
		filename := fmt.Sprintf("lottery_analysis_%s.%s", 
			time.Now().Format("20060102_150405"), 
			config.ExportFormat)
		if err := analyzer.ExportAnalysis(ctx, filename); err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting analysis: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nAnalysis exported to: %s\n", filename)
	}
}

// printHelp displays usage information
func printHelp() {
	fmt.Println("NC Lucky for Life Lottery Analyzer")
	fmt.Println("==================================")
	fmt.Println()
	fmt.Println("Usage: go run lottery_analyzer.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --simple           Show simplified analysis")
	fmt.Println("  --statistical      Show detailed statistical analysis")
	fmt.Println("  --cosmic           Show cosmic correlation analysis")
	fmt.Println("  --export-json      Export results to JSON file")
	fmt.Println("  --export-csv       Export results to CSV file")
	fmt.Println("  --recent <n>       Set recent window size (default: 50)")
	fmt.Println("  --help             Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run lottery_analyzer.go")
	fmt.Println("  go run lottery_analyzer.go --simple")
	fmt.Println("  go run lottery_analyzer.go --statistical --export-json")
	fmt.Println("  go run lottery_analyzer.go --recent 100")
}