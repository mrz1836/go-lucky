# Fuzz Testing Guide for Go-Lucky

This document describes the fuzz testing infrastructure for the Go-Lucky lottery analyzer.

## Overview

Fuzz testing helps discover edge cases and potential security vulnerabilities by providing random, malformed, or unexpected input to functions. The Go-Lucky project includes comprehensive fuzz tests for both the lottery analyzer and cosmic correlator components.

## Running Fuzz Tests

### Run All Fuzz Tests
```bash
magex test:fuzz
```

### Run Specific Fuzz Test
```bash
# Run for 30 seconds (default)
go test -fuzz=FuzzNewAnalyzer

# Run for specific duration
go test -fuzz=FuzzNewAnalyzer -fuzztime=2m

# Run with specific number of workers
go test -fuzz=FuzzNewAnalyzer -parallel=8
```

### Run with Coverage
```bash
go test -fuzz=FuzzNewAnalyzer -cover -coverprofile=fuzz_coverage.out
go tool cover -html=fuzz_coverage.out
```

## Fuzz Test Coverage

### Lottery Analyzer Fuzzing (`lottery_analyzer_fuzz_test.go`)

1. **FuzzNewAnalyzer**
   - Tests CSV parsing with malformed input
   - Validates date parsing, number ranges, and data integrity
   - Tests various CSV formats and edge cases

2. **FuzzValidateFilePath**
   - Tests path traversal attempts
   - Validates handling of special characters and null bytes
   - Tests Windows/Unix path formats

3. **FuzzParseDrawings**
   - Tests CSV record parsing with invalid data
   - Validates number range enforcement
   - Tests empty and malformed records

4. **FuzzAnalysisConfig**
   - Tests configuration with extreme values
   - Validates defaults and boundary conditions
   - Tests invalid output modes and formats

5. **FuzzExportAnalysis**
   - Tests export with various file paths
   - Validates format handling (JSON/CSV)
   - Tests error conditions

6. **FuzzScoreNumbersByStrategy**
   - Tests all scoring strategies with invalid inputs
   - Validates score calculations and sorting
   - Tests unknown strategies

7. **FuzzGenerateSetByStrategy**
   - Tests recommendation generation
   - Validates uniqueness constraints
   - Tests number range enforcement

### Cosmic Correlator Fuzzing (`cosmic_correlator_fuzz_test.go`)

1. **FuzzCalculateMoonPhase**
   - Tests with extreme dates (past and future)
   - Validates phase and illumination ranges
   - Tests date overflow conditions

2. **FuzzGetZodiacSign**
   - Tests zodiac boundaries
   - Validates all months and edge dates
   - Tests leap year handling

3. **FuzzCalculatePearsonCorrelation**
   - Tests with various data patterns
   - Handles NaN and Inf values
   - Tests empty and mismatched arrays

4. **FuzzPlanetaryPositions**
   - Tests position calculations across time
   - Validates angle ranges (0-360°)
   - Tests extreme dates

5. **FuzzCosmicDataEnrichment**
   - Tests with various drawing counts
   - Validates date range handling
   - Tests memory efficiency with large datasets

6. **FuzzPredictBasedOnCosmicConditions**
   - Tests prediction consistency
   - Validates number uniqueness
   - Tests with pre-populated cosmic data

## Seed Corpus

The `testdata/fuzz/` directory contains seed inputs for fuzz tests:

- `FuzzNewAnalyzer/`: Sample CSV files with edge cases
- `FuzzValidateFilePath/`: Path traversal attempts

## Security Considerations

The fuzz tests specifically target:

1. **Input Validation**
   - CSV injection attempts
   - Path traversal vulnerabilities
   - Buffer overflow attempts with large inputs

2. **Data Integrity**
   - Number range validation (1-48 for main, 1-18 for lucky)
   - Date parsing and validation
   - Statistical calculation overflow

3. **Resource Exhaustion**
   - Large file handling
   - Memory allocation with extreme values
   - CPU-intensive calculations

## Adding New Fuzz Tests

When adding new functionality, create corresponding fuzz tests:

```go
func FuzzNewFeature(f *testing.F) {
    // Add seed corpus
    f.Add("normal input")
    f.Add("edge case")
    
    f.Fuzz(func(t *testing.T, input string) {
        // Call function with fuzz input
        result, err := NewFeature(input)
        
        // Validate invariants
        if err == nil {
            // Check result properties
        }
    })
}
```

## Best Practices

1. **Always validate output** - Don't just check for panics
2. **Add meaningful seed corpus** - Include known edge cases
3. **Test invariants** - Verify expected properties hold
4. **Handle special values** - Test NaN, Inf, nil, empty
5. **Check boundaries** - Test limits and overflows

## Continuous Fuzzing

Consider integrating with OSS-Fuzz or running nightly fuzz jobs:

```bash
# Run all fuzz tests for extended period
for test in $(go test -list Fuzz); do
    go test -fuzz=$test -fuzztime=10m
done
```

## Troubleshooting

### Out of Memory
- Reduce `-parallel` flag
- Add memory limits to test execution
- Check for memory leaks in test data generation

### Slow Tests
- Use `-fuzzminimizetime=10s` to speed up minimization
- Profile with `-cpuprofile=cpu.prof`
- Check for expensive validations in hot paths

### False Positives
- Validate that fuzz input meets function preconditions
- Skip invalid inputs early with `t.Skip()`
- Document expected input constraints