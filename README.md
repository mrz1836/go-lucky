# Go-Lucky: NC Lucky for Life Lottery Analyzer 🌌

A sophisticated statistical analysis tool for NC Lucky for Life lottery data that combines traditional statistical methods with cosmic correlation analysis. The analyzer verifies randomness, identifies patterns, and generates number recommendations using multiple strategies including astronomical and environmental factors.

**⚠️ IMPORTANT**: This tool is for educational and entertainment purposes only. Lottery drawings are random events, and no analysis can predict future outcomes or improve odds.

## Features

### 🎯 Core Statistical Analysis
- **Frequency Analysis**: Track how often each number appears overall and recently
- **Gap Analysis**: Calculate average gaps between number appearances
- **Pattern Detection**: Identify odd/even distributions, consecutive numbers, sum ranges
- **Combination Tracking**: Monitor frequently appearing pairs, triples, and quads
- **Statistical Verification**: Chi-square test to verify randomness

### 🌌 Cosmic Correlation Analysis (NEW!)
- **Moon Phase Correlations**: Analyze number patterns during different lunar phases
- **Solar Activity Analysis**: Correlate solar wind, geomagnetic activity with drawings
- **Weather Pattern Analysis**: Historical weather data correlation with number selection
- **Planetary Position Tracking**: Track planetary influences and Mercury retrograde effects
- **Zodiac and Seasonal Analysis**: Temporal pattern analysis based on astronomical events
- **Cosmic-Influenced Predictions**: Generate number sets based on current cosmic conditions

### 📊 Advanced Features
- **Multi-Strategy Recommendations**: Generate number sets using different approaches
- **Confidence Scoring**: Statistical confidence levels for each recommendation
- **Export Capabilities**: Save results to JSON or CSV formats
- **Configurable Analysis**: Adjust recent window, gap thresholds, and more
- **Comprehensive Testing**: Full test suite with benchmarks
- **Make Commands**: Easy-to-use Makefile for all analysis types

## Installation

```bash
# Clone the repository
git clone https://github.com/mrz1836/go-lucky.git
cd go-lucky

# Install dependencies
make install-deps

# Run the full cosmic analysis (RECOMMENDED)
make full-analysis
```

## Quick Start

### Using Make Commands (Recommended)

```bash
# Run COMPLETE analysis with cosmic correlations
make full-analysis

# Quick summary with cosmic picks
make simple

# Show current hot numbers
make hot-numbers

# Generate multiple number sets
make lucky-picks

# Export full analysis to JSON
make export-json

# See all available commands
make help
```

### Direct Go Commands

```bash
# Run detailed analysis with cosmic correlations (default)
go run lottery_analyzer.go cosmic_correlator.go

# Simple summary view
go run lottery_analyzer.go cosmic_correlator.go --simple

# Cosmic correlation analysis only
go run lottery_analyzer.go cosmic_correlator.go --cosmic

# Statistical deep dive
go run lottery_analyzer.go cosmic_correlator.go --statistical
```

## Analysis Modes

### 🌟 Full Analysis (`make full-analysis`)
The complete experience including:
- Traditional statistical analysis
- Cosmic correlation analysis
- Moon phase patterns
- Solar activity correlations
- Weather pattern analysis
- Current cosmic conditions
- Multiple prediction strategies

Sample output:
```
🌌 COSMIC CORRELATION ANALYSIS 🌌
═══════════════════════════════════════════════════════════════════

🌙 LUNAR CORRELATIONS
─────────────────────
• Full Moon Lucky Numbers:
  Number 7 appears 2.9% more frequently during Full Moon

☀️ SOLAR ACTIVITY CORRELATIONS
─────────────────────────────
• Solar Wind vs High Numbers:
  No significant correlation detected

🔮 CURRENT COSMIC CONDITIONS
──────────────────────────
Date: July 23, 2025
Moon Phase: New Moon (4% illuminated)
Zodiac Sign: Leo

🌟 Cosmic Selection: 45-22-01-23-42  Lucky Ball: 11
```

### 📊 Statistical Mode (`--statistical`)
Detailed mathematical analysis including:
- Chi-square randomness testing
- Frequency distribution analysis
- Gap analysis statistics
- P-values and significance testing

### 🎯 Simple Mode (`--simple`)
Quick overview with:
- Top hot numbers
- Most overdue numbers
- Quick picks
- Cosmic pick of the day

## Understanding the Output

### Randomness Score
- **90-100%**: Highly random (expected for fair lottery)
- **70-90%**: Mostly random with minor deviations
- **Below 70%**: Showing non-random patterns (investigate further)

### Cosmic Correlations
- **P-values < 0.05**: Statistically significant (but likely coincidental)
- **P-values > 0.1**: No meaningful correlation
- All cosmic correlations are within expected random variation

### Number Categories
- **Hot Numbers**: Frequently appearing in recent drawings
- **Cold Numbers**: Haven't appeared recently
- **Overdue Numbers**: Haven't appeared for longer than their average gap
- **Cosmic Numbers**: Generated based on current astronomical conditions

### Recommendation Strategies
1. **Balanced**: Mix of frequency, recency, and overdue factors
2. **Hot**: Focus on recently frequent numbers
3. **Overdue**: Numbers that haven't appeared beyond average gap
4. **Pattern**: Numbers that frequently appear together
5. **Frequency**: Pure historical frequency approach
6. **Cosmic**: Based on current moon phase, zodiac, and planetary positions

## NC Lucky for Life Rules

- Pick 5 numbers from 1-48
- Pick 1 Lucky Ball from 1-18
- Drawings held every day
- Overall odds of winning any prize: 1 in 7.77
- Jackpot odds: 1 in 30,821,472

### Prize Structure
| Match | Prize | Odds |
|-------|-------|------|
| 5 + Lucky Ball | $1,000/day for life | 1 in 30,821,472 |
| 5 numbers | $25,000/year for life | 1 in 1,813,028 |
| 4 + Lucky Ball | $5,000 | 1 in 143,356 |
| 4 numbers | $200 | 1 in 8,433 |
| 3 + Lucky Ball | $150 | 1 in 3,413 |
| 3 numbers | $20 | 1 in 201 |
| 2 + Lucky Ball | $25 | 1 in 250 |
| 2 numbers | $3 | 1 in 15 |
| 1 + Lucky Ball | $6 | 1 in 50 |
| 0 + Lucky Ball | $4 | 1 in 32 |

## Development

### Project Structure
```
go-lucky/
├── lottery_analyzer.go      # Main analyzer implementation
├── cosmic_correlator.go     # Cosmic correlation engine
├── lottery_analyzer_test.go # Comprehensive test suite
├── lucky-numbers-history.csv # Historical drawing data
├── Makefile                 # Build and analysis commands
├── CLAUDE.md               # AI assistant context
├── tech-conventions.md     # Go development guidelines
├── LOTTERY_ANALYSIS_INSIGHTS.md # Deep analysis insights
└── README.md              # This file
```

### Testing

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run benchmarks
make benchmark

# Test specific functionality
go test -run TestAnalyzerSuite/TestCosmicCorrelations -v
```

### Local Development

```bash
# Serve the website using Node.js
npx serve . -p 8000

# Then visit http://localhost:8000
```

### Adding Features
1. Follow conventions in `tech-conventions.md`
2. Use context-first design
3. Add comprehensive tests
4. Update documentation
5. Ensure randomness verification isn't compromised

## Make Commands Reference

### Analysis Commands
- `make full-analysis` - 🌟 Complete cosmic + statistical analysis (RECOMMENDED)
- `make simple` - Quick summary with hot numbers and picks
- `make statistical` - Detailed mathematical analysis
- `make cosmic` - Cosmic correlations only
- `make lucky-picks` - Generate multiple number sets
- `make hot-numbers` - Show current hot numbers
- `make overdue` - Show most overdue numbers

### Export Commands
- `make export-json` - Export full analysis to JSON
- `make export-csv` - Export analysis data to CSV

### Development Commands
- `make build` - Build the analyzer binary
- `make test` - Run all tests
- `make coverage` - Generate test coverage report
- `make lint` - Run code linters
- `make clean` - Clean up generated files
- `make install-deps` - Install/update dependencies

### Fun Commands
- `make cosmic-wisdom` - Display cosmic lottery wisdom
- `make fortune` - Get your cosmic lottery fortune

## Mathematical Insights

### Why Cosmic Patterns Don't Predict
- Each drawing is an independent random event
- Mechanical randomization ensures fairness
- Past results and cosmic events have zero influence on future draws
- Patterns in historical data are statistical noise, not meaningful correlations

### What This Tool Actually Does
- ✅ Verifies lottery fairness through statistical tests
- ✅ Provides entertainment through pattern analysis
- ✅ Educates about probability, statistics, and correlation vs causation
- ✅ Demonstrates how to analyze random data scientifically
- ❌ Does NOT predict future numbers
- ❌ Does NOT improve winning odds
- ❌ Cosmic correlations are NOT predictive

### The Cosmic Correlation Findings
Our analysis shows:
- **Moon phases**: No significant correlation with number patterns
- **Solar activity**: No meaningful impact on lottery outcomes
- **Weather patterns**: No correlation with number selection
- **Planetary positions**: No detectable influence
- **All correlations**: Within expected random variation

This confirms the lottery operates as designed - as a truly random system unaffected by external cosmic forces.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Follow Go conventions and add tests
4. Commit your changes (`git commit -m 'feat: add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## Performance

- Analyzes 2000+ drawings in ~1 second
- Cosmic data enrichment adds ~2 seconds
- Memory usage: ~50MB for full dataset
- Correlation analysis: ~500ms additional

## License

This project is for educational purposes. Use at your own risk.

## Disclaimer

**CRITICAL UNDERSTANDING**: This tool demonstrates sophisticated statistical and correlation analysis of random data. It includes cosmic correlation analysis to show how even astronomical events have no meaningful relationship with lottery outcomes.

The cosmic correlation features are designed to:
1. **Educate** about correlation vs causation
2. **Demonstrate** proper statistical analysis techniques
3. **Entertain** with interesting but meaningless patterns
4. **Prove** that lottery drawings are truly random

**The tool cannot and will not**:
- Predict lottery numbers
- Improve your odds of winning
- Find meaningful cosmic influences on random events

Remember: In a fair lottery, every number combination has exactly the same probability of being drawn, regardless of moon phases, solar activity, weather patterns, or any other external factors.

*"The lottery is a tax on people who are bad at math, but it's also a fascinating demonstration of true randomness in action."*

---

*Play responsibly, if at all. The only guaranteed way to not lose money on the lottery is to not play.*