# Go-Lucky Lottery Analyzer Makefile
# Advanced lottery analysis with cosmic correlations

# Include base make files
include .make/common.mk
include .make/go.mk

# Variables
BINARY_NAME=lottery-analyzer
MAIN_FILES=lottery_analyzer.go cosmic_correlator.go
TEST_TIMEOUT=30s
COVERAGE_FILE=coverage.out

# Default target
.DEFAULT_GOAL := help

# PHONY targets
.PHONY: clean full-analysis simple statistical cosmic export-json export-csv

##@ Analysis Commands

## full-analysis: 🌟 Run COMPLETE analysis with cosmic correlations (RECOMMENDED)
full-analysis: ## Run full analysis with cosmic correlations
	@magex build:dev
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║        🌌 RUNNING FULL COSMIC LOTTERY ANALYSIS 🌌            ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo ""
	@./bin/$(BINARY_NAME) --cosmic

## simple: Run simple analysis summary
simple: ## Run simple analysis summary
	@magex build:dev
	@./bin/$(BINARY_NAME) --simple

## statistical: Run detailed statistical analysis
statistical: ## Run detailed statistical analysis
	@magex build:dev
	@./bin/$(BINARY_NAME) --statistical

## cosmic: Run cosmic correlation analysis only
cosmic: ## Run cosmic correlation analysis only
	@magex build:dev
	@./bin/$(BINARY_NAME) --cosmic

##@ Export Commands

## export-json: Export full analysis to JSON file
export-json: ## Export analysis data to JSON file
	@magex build:dev
	@echo "📊 Exporting analysis to JSON..."
	@./bin/$(BINARY_NAME) --cosmic --export-json
	@echo "✅ Export complete! Check lottery_analysis_*.json"

## export-csv: Export analysis data to CSV file
export-csv: ## Export analysis data to CSV file
	@magex build:dev
	@echo "📊 Exporting analysis to CSV..."
	@./bin/$(BINARY_NAME) --cosmic --export-csv
	@echo "✅ Export complete! Check lottery_analysis_*.csv"


## clean: Clean build artifacts and generated files
clean: ## Clean up build artifacts and generated files
	@echo "🧹 Cleaning up..."
	@rm -rf bin/
	@rm -f $(BINARY_NAME)
	@rm -f $(COVERAGE_FILE) coverage.html coverage.txt coverage.out
	@rm -f lottery_analysis_*.json lottery_analysis_*.csv
	@rm -f test_*.csv debug_*.csv empty_*.csv invalid_*.csv
	@echo "✅ Cleanup complete"

## benchmark: Run performance benchmarks
benchmark: ## Run performance benchmarks
	@$(MAKE) bench

##@ Quick Analysis Sets

## lucky-picks: Generate 5 different analysis-based number sets
lucky-picks: ## Generate 5 different analysis-based number sets
	@magex build:dev
	@echo "🎰 Generating Lucky Picks..."
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@./bin/$(BINARY_NAME) --simple | grep -A 10 "QUICK PICKS:" || true
	@echo ""
	@echo "🌌 Cosmic Pick:"
	@./bin/$(BINARY_NAME) --simple | grep -A 1 "COSMIC PICK:" || true
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

## hot-numbers: Show current hot numbers
hot-numbers: ## Show current hot numbers
	@magex build:dev
	@echo "🔥 Current Hot Numbers:"
	@./bin/$(BINARY_NAME) --simple | grep -A 7 "TOP 5 HOT NUMBERS:" || true

## overdue: Show most overdue numbers
overdue: ## Show most overdue numbers
	@magex build:dev
	@echo "⏰ Most Overdue Numbers:"
	@./bin/$(BINARY_NAME) --simple | grep -A 7 "TOP 5 OVERDUE:" || true

# Special targets for fun
.PHONY: cosmic-wisdom fortune

## cosmic-wisdom: Display cosmic lottery wisdom
cosmic-wisdom: ## Display cosmic lottery wisdom
	@echo ""
	@echo "✨ ═══════════════════════════════════════════════════════ ✨"
	@echo "   🌙 The moon influences tides, not lottery numbers! 🌙"
	@echo "   ☀️  Solar flares can't burn through randomness! ☀️"
	@echo "   🌟 Every number has exactly 1/48 chance! 🌟"
	@echo "   🎲 Play for fun, not for cosmic fortune! 🎲"
	@echo "✨ ═══════════════════════════════════════════════════════ ✨"
	@echo ""

## fortune: Get your lottery fortune
fortune: ## Get your lottery fortune
	@magex build:dev
	@echo "🔮 Your Lottery Fortune:"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@./bin/$(BINARY_NAME) --simple | grep -A 1 "COSMIC PICK:" || echo "The stars are silent today..."
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "Remember: Fortune favors the prepared... wallet! 💸"
