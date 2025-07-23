# Go-Lucky Lottery Analyzer Makefile
# Advanced lottery analysis with cosmic correlations

# Variables
BINARY_NAME=lottery-analyzer
MAIN_FILES=lottery_analyzer.go cosmic_correlator.go
TEST_TIMEOUT=30s
COVERAGE_FILE=coverage.out

# Default target
.DEFAULT_GOAL := help

# PHONY targets
.PHONY: help build test clean full-analysis simple statistical cosmic export-json export-csv coverage lint

## help: Display this help message
help:
	@echo "Go-Lucky Lottery Analyzer - Available Commands:"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""

##@ Analysis Commands

## full-analysis: 🌟 Run COMPLETE analysis with cosmic correlations (RECOMMENDED)
full-analysis: build
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║        🌌 RUNNING FULL COSMIC LOTTERY ANALYSIS 🌌             ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo ""
	@./$(BINARY_NAME) --cosmic

## simple: Run simple analysis summary
simple: build
	@./$(BINARY_NAME) --simple

## statistical: Run detailed statistical analysis
statistical: build
	@./$(BINARY_NAME) --statistical

## cosmic: Run cosmic correlation analysis only
cosmic: build
	@./$(BINARY_NAME) --cosmic

##@ Export Commands

## export-json: Export full analysis to JSON file
export-json: build
	@echo "📊 Exporting analysis to JSON..."
	@./$(BINARY_NAME) --cosmic --export-json
	@echo "✅ Export complete! Check lottery_analysis_*.json"

## export-csv: Export analysis data to CSV file
export-csv: build
	@echo "📊 Exporting analysis to CSV..."
	@./$(BINARY_NAME) --cosmic --export-csv
	@echo "✅ Export complete! Check lottery_analysis_*.csv"

##@ Development Commands

## build: Build the lottery analyzer binary
build:
	@echo "🔨 Building lottery analyzer..."
	@go build -o $(BINARY_NAME) $(MAIN_FILES)
	@echo "✅ Build complete: $(BINARY_NAME)"

## test: Run all tests
test:
	@echo "🧪 Running tests..."
	@go test -v -timeout=$(TEST_TIMEOUT) ./...

## coverage: Run tests with coverage report
coverage:
	@echo "📊 Running tests with coverage..."
	@go test -coverprofile=$(COVERAGE_FILE) -v ./...
	@go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

## lint: Run linters
lint:
	@echo "🔍 Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		go vet ./...; \
		go fmt ./...; \
	fi

## clean: Clean build artifacts and generated files
clean:
	@echo "🧹 Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f $(COVERAGE_FILE) coverage.html
	@rm -f lottery_analysis_*.json lottery_analysis_*.csv
	@rm -f test_*.csv debug_*.csv empty_*.csv invalid_*.csv
	@echo "✅ Cleanup complete"

##@ Quick Analysis Sets

## lucky-picks: Generate 5 different analysis-based number sets
lucky-picks: build
	@echo "🎰 Generating Lucky Picks..."
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@./$(BINARY_NAME) --simple | grep -A 10 "QUICK PICKS:" || true
	@echo ""
	@echo "🌌 Cosmic Pick:"
	@./$(BINARY_NAME) --simple | grep -A 1 "COSMIC PICK:" || true
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

## hot-numbers: Show current hot numbers
hot-numbers: build
	@echo "🔥 Current Hot Numbers:"
	@./$(BINARY_NAME) --simple | grep -A 7 "TOP 5 HOT NUMBERS:" || true

## overdue: Show most overdue numbers
overdue: build
	@echo "⏰ Most Overdue Numbers:"
	@./$(BINARY_NAME) --simple | grep -A 7 "TOP 5 OVERDUE:" || true

##@ Utility Commands

## install-deps: Install/update Go dependencies
install-deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

## update: Update dependencies to latest versions
update:
	@echo "🔄 Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✅ Dependencies updated"

## benchmark: Run performance benchmarks
benchmark:
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem -run=^$

# Special targets for fun
.PHONY: cosmic-wisdom fortune

## cosmic-wisdom: Display cosmic lottery wisdom
cosmic-wisdom:
	@echo ""
	@echo "✨ ═══════════════════════════════════════════════════════ ✨"
	@echo "   🌙 The moon influences tides, not lottery numbers! 🌙"
	@echo "   ☀️  Solar flares can't burn through randomness! ☀️"
	@echo "   🌟 Every number has exactly 1/48 chance! 🌟"
	@echo "   🎲 Play for fun, not for cosmic fortune! 🎲"
	@echo "✨ ═══════════════════════════════════════════════════════ ✨"
	@echo ""

## fortune: Get your lottery fortune
fortune: build
	@echo "🔮 Your Lottery Fortune:"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@./$(BINARY_NAME) --simple | grep -A 1 "COSMIC PICK:" || echo "The stars are silent today..."
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "Remember: Fortune favors the prepared... wallet! 💸"