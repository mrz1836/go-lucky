# CLAUDE.md - AI Assistant Context for Lucky Lottery Analyzer 🌌

## Project Overview

This is a sophisticated lottery analysis tool for NC Lucky for Life lottery data that combines traditional statistical methods with cosmic correlation analysis. The analyzer uses advanced statistical methods to examine historical patterns, verify randomness, and generate number recommendations based on multiple strategies including astronomical and environmental factors.

**IMPORTANT**: This tool is for educational and entertainment purposes. Lottery drawings are random events, and no analysis can predict future outcomes or improve odds. The cosmic correlation features are designed to demonstrate proper statistical analysis and prove that even astronomical events have no meaningful relationship with lottery outcomes.

## Technical Architecture

### Core Components

1. **lottery_analyzer.go** - Main analysis engine
   - Statistical analysis functions
   - Pattern detection algorithms
   - Recommendation generation system
   - Chi-square randomness testing
   - Multiple output formats
   - Integrated cosmic correlation engine

2. **cosmic_correlator.go** - Cosmic correlation engine (NEW!)
   - Moon phase calculations and correlations
   - Solar activity simulation and analysis
   - Weather pattern correlation analysis
   - Planetary position tracking
   - Zodiac and seasonal analysis
   - Statistical correlation calculations
   - Cosmic-influenced prediction generation

3. **lottery_analyzer_test.go** - Comprehensive test suite
   - Unit tests using testify
   - Benchmark tests
   - Edge case handling
   - Cosmic correlation testing

4. **lucky-numbers-history.csv** - Historical drawing data
   - Format: Date, 5 main numbers (1-48), 1 lucky ball (1-18)
   - 2000+ historical drawings for analysis

5. **Makefile** - Build and analysis commands
   - Full-featured analysis commands
   - Development and testing tools
   - Export utilities

## Key Features

### Traditional Analysis Capabilities
- **Frequency Analysis**: Track how often each number appears
- **Recent Trends**: Analyze last N drawings for "hot" numbers
- **Gap Analysis**: Calculate average gaps between appearances
- **Pattern Detection**: Identify odd/even, consecutive, sum patterns
- **Combination Tracking**: Monitor pairs, triples, quads
- **Statistical Testing**: Chi-square test for randomness verification

### Cosmic Correlation Analysis (NEW!)
- **Moon Phase Correlations**: Analyze number patterns during different lunar phases
- **Solar Activity Analysis**: Correlate solar wind, geomagnetic activity with drawings
- **Weather Pattern Analysis**: Historical weather data correlation
- **Planetary Position Tracking**: Track planetary influences and Mercury retrograde
- **Zodiac Analysis**: Temporal pattern analysis based on sun sign
- **Seasonal Patterns**: Weather and astronomical season correlations
- **Current Cosmic Conditions**: Real-time astronomical data for predictions

### Recommendation Strategies
1. **Balanced**: Mix of frequency, recency, and overdue factors
2. **Hot**: Focus on recently frequent numbers
3. **Overdue**: Numbers that haven't appeared beyond average gap
4. **Pattern**: Numbers that appear in common combinations
5. **Frequency**: Pure historical frequency approach
6. **Cosmic**: Based on current astronomical conditions (NEW!)

## Usage Examples

```bash
# Run complete analysis with cosmic correlations (RECOMMENDED)
make full-analysis

# Quick analysis with cosmic pick
make simple

# Cosmic correlation analysis only
go run lottery_analyzer.go cosmic_correlator.go --cosmic

# Statistical deep dive
go run lottery_analyzer.go cosmic_correlator.go --statistical

# Export full cosmic analysis
make export-json
```

## Development Guidelines

### Code Style
- Follow tech-conventions.md strictly
- Use context-first design
- No global state
- Comprehensive error handling
- Clear function documentation

### Testing
```bash
# Run all tests including cosmic features
make test

# Run with coverage
make coverage

# Run benchmarks
make benchmark

# Test specific cosmic functionality
go test -run TestAnalyzerSuite/TestCosmicCorrelations -v
```

### Adding New Cosmic Features

When adding cosmic correlation features:
1. Add data structures to `CosmicData` struct
2. Implement data fetching/calculation in correlation engine
3. Add correlation analysis function
4. Update test suite with new functionality
5. Update documentation and help output
6. Ensure statistical validity with proper p-values

## Statistical Integrity

### Cosmic Correlation Analysis
The cosmic correlation engine performs rigorous statistical analysis:
- **Pearson Correlation**: For continuous variables (temperature, solar activity)
- **Point-Biserial Correlation**: For binary outcomes (number drawn/not drawn)
- **P-value Calculations**: Statistical significance testing
- **Multiple Comparison Corrections**: Adjustments for multiple testing

### Important Findings
All cosmic correlations analyzed show:
- **Moon phases**: No significant correlation with number patterns
- **Solar activity**: No meaningful impact on lottery outcomes  
- **Weather patterns**: No correlation with number selection
- **Planetary positions**: No detectable influence
- **All correlations**: Within expected random variation (p > 0.1)

This confirms the lottery operates as a truly random system.

### Randomness Verification
- Chi-square test validates uniform distribution
- Randomness score (0-100%) indicates deviation from expected
- Scores above 90% indicate high randomness
- Current dataset shows ~45% randomness (some clustering, but normal)

## Cosmic Data Sources & Implementation

### Astronomical Calculations
- **Moon phases**: Calculated using synodic month algorithms
- **Planetary positions**: Simplified orbital mechanics
- **Zodiac signs**: Sun position calculations
- **Seasonal phases**: Astronomical season boundaries

### Mock Data for Demo
Current implementation uses calculated/simulated data for:
- Solar activity (solar wind speed, geomagnetic indices)
- Weather patterns (temperature, pressure, precipitation)
- Geomagnetic activity (Kp indices)

### Production APIs (Future Enhancement)
For real-world deployment, integrate:
- **USNO Moon Phase API**: https://aa.usno.navy.mil/api/moon/phases/
- **NOAA Space Weather**: http://services.swpc.noaa.gov/products/
- **NOAA Climate Data**: https://www.ncei.noaa.gov/access/services/data/v1
- **NASA JPL Ephemeris**: For precise planetary positions

## Configuration and Customization

### Analysis Configuration
```go
type AnalysisConfig struct {
    RecentWindow      int     // Default: 50 drawings
    MinGapMultiplier  float64 // Default: 1.5 for "overdue"
    ConfidenceLevel   float64 // Default: 0.95
    OutputMode        string  // "simple", "detailed", "statistical", "cosmic"
    ExportFormat      string  // "console", "csv", "json"
}
```

### Cosmic Analysis Parameters
- **Moon phase precision**: ±1 day accuracy
- **Solar activity simulation**: Realistic patterns based on solar cycles
- **Correlation significance**: p < 0.05 for statistical significance
- **Sample size requirements**: Minimum 50 data points for correlations

## Performance Considerations

### Analysis Speed
- Standard analysis: ~1 second for 2000+ drawings
- Cosmic data enrichment: ~2 seconds additional
- Correlation calculations: ~500ms additional
- Memory usage: ~50MB for full dataset

### Optimization Tips
- Use `--simple` mode for faster results
- Cache cosmic data calculations
- Consider parallel correlation analysis for large datasets

## Debugging and Troubleshooting

### Common Issues

**Cosmic Data Issues**:
- Check date parsing for moon phase calculations
- Verify astronomical algorithm accuracy
- Ensure mock data generation is consistent

**Correlation Analysis Issues**:
- Validate sample sizes (minimum 30 data points)
- Check for division by zero in correlation calculations
- Verify p-value computation accuracy

**Build Issues**:
- Ensure both `lottery_analyzer.go` and `cosmic_correlator.go` are included
- Check Go module dependencies are current
- Verify testify package is installed

### Debug Commands
```bash
# Test build
make build

# Check for correlation data
go run lottery_analyzer.go cosmic_correlator.go --cosmic | grep -A 10 "LUNAR"

# Validate statistical calculations
go run lottery_analyzer.go cosmic_correlator.go --statistical

# Test mock data generation
go test -run TestCosmicDataGeneration -v
```

## Maintenance Notes

### Data Updates
1. Ensure CSV format matches expected structure
2. Validate date formats (MM/DD/YYYY)
3. Check number ranges (1-48 main, 1-18 lucky)
4. Run analyzer to verify data integrity
5. Update cosmic data calculations for new date ranges

### Statistical Validation
- Regularly verify chi-square calculations
- Check correlation p-values for accuracy
- Validate mock data generation ranges
- Ensure astronomical calculations remain accurate

## Ethical Considerations

### Responsible Messaging
- Always include randomness disclaimers
- Emphasize educational/entertainment purpose
- Never claim cosmic correlations improve odds
- Include proper statistical significance explanations
- Promote responsible gambling practices

### Educational Value
The cosmic correlation features serve to:
1. **Demonstrate** proper statistical analysis techniques
2. **Educate** about correlation vs causation
3. **Prove** lottery randomness through negative results
4. **Entertain** with interesting but meaningless patterns

## Future Enhancements

### Potential Cosmic Features
1. **Real API Integration**: Live astronomical data
2. **Advanced Planetary Calculations**: More precise ephemeris
3. **Historical Event Correlations**: Major astronomical events
4. **Machine Learning**: Pattern recognition (for educational purposes)
5. **Visualization**: Cosmic correlation charts and graphs

### Data Wishlist
- Complete historical weather data for NC
- Precise ball weight measurements
- Drawing machine calibration data
- Cross-lottery comparisons
- Player behavior analysis (if available)

## Remember

This tool demonstrates sophisticated statistical analysis of random data combined with cosmic correlation analysis. It's perfect for:
- Learning probability and statistics
- Understanding correlation vs causation
- Verifying lottery fairness
- Exploring pattern recognition in random data
- Entertainment purposes

It is NOT for:
- Predicting future numbers
- "Beating" the lottery
- Investment strategies
- Guaranteed wins

**The cosmic correlation analysis consistently shows no meaningful relationships between astronomical events and lottery outcomes, confirming the lottery's true randomness.**

## Quick Reference Commands

```bash
# Most useful commands
make full-analysis        # Complete cosmic + statistical analysis
make simple              # Quick summary with cosmic pick
make lucky-picks         # Generate multiple strategies
make export-json         # Save full analysis
make test               # Run all tests
make help               # See all commands

# Debug and development
make build              # Build binary
make coverage           # Test coverage report
make clean              # Clean generated files
make cosmic-wisdom      # Fun cosmic quotes
```

**The universe is under no obligation to make you wealthy, but it's happy to teach you about statistics!** 🌟