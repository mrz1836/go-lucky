// Initialize Lucide icons
document.addEventListener('DOMContentLoaded', function() {
    lucide.createIcons();
    
    // Initialize charts
    initializeCharts();
    
    // Initialize interactive elements
    initializeInteractivity();
    
    // Initialize animations
    initializeAnimations();
});

// Chart initialization and management
function initializeCharts() {
    // Frequency Distribution Chart
    const frequencyCtx = document.getElementById('frequencyChart');
    if (frequencyCtx) {
        const frequencyChart = new Chart(frequencyCtx, {
            type: 'bar',
            data: {
                labels: generateNumberLabels(1, 48),
                datasets: [{
                    label: 'Frequency',
                    data: generateMockFrequencyData(),
                    backgroundColor: createGradientArray(48, '#6366f1', '#8b5cf6'),
                    borderColor: '#6366f1',
                    borderWidth: 1,
                    borderRadius: 4,
                    borderSkipped: false,
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        backgroundColor: 'rgba(15, 15, 35, 0.9)',
                        titleColor: '#ffffff',
                        bodyColor: '#a1a1aa',
                        borderColor: '#6366f1',
                        borderWidth: 1,
                        cornerRadius: 8,
                        displayColors: false,
                        callbacks: {
                            title: function(context) {
                                return `Number ${context[0].label}`;
                            },
                            label: function(context) {
                                return `Drawn ${context.parsed.y} times`;
                            }
                        }
                    }
                },
                scales: {
                    x: {
                        grid: {
                            display: false
                        },
                        ticks: {
                            color: '#71717a',
                            font: {
                                size: 10
                            },
                            maxTicksLimit: 12
                        }
                    },
                    y: {
                        grid: {
                            color: 'rgba(255, 255, 255, 0.1)'
                        },
                        ticks: {
                            color: '#71717a',
                            font: {
                                size: 10
                            }
                        }
                    }
                },
                interaction: {
                    intersect: false,
                    mode: 'index'
                }
            }
        });
        
        // Store chart reference for updates
        window.frequencyChart = frequencyChart;
    }
    
    // Cosmic Correlations Chart
    const cosmicCtx = document.getElementById('cosmicChart');
    if (cosmicCtx) {
        const cosmicChart = new Chart(cosmicCtx, {
            type: 'line',
            data: {
                labels: ['New Moon', 'Waxing Crescent', 'First Quarter', 'Waxing Gibbous', 'Full Moon', 'Waning Gibbous', 'Last Quarter', 'Waning Crescent'],
                datasets: [{
                    label: 'Correlation Strength',
                    data: [0.012, -0.008, 0.023, 0.015, 0.031, -0.005, 0.018, -0.012],
                    borderColor: '#7c3aed',
                    backgroundColor: 'rgba(124, 58, 237, 0.1)',
                    borderWidth: 2,
                    fill: true,
                    tension: 0.4,
                    pointBackgroundColor: '#7c3aed',
                    pointBorderColor: '#ffffff',
                    pointBorderWidth: 2,
                    pointRadius: 4,
                    pointHoverRadius: 6
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        backgroundColor: 'rgba(15, 15, 35, 0.9)',
                        titleColor: '#ffffff',
                        bodyColor: '#a1a1aa',
                        borderColor: '#7c3aed',
                        borderWidth: 1,
                        cornerRadius: 8,
                        displayColors: false,
                        callbacks: {
                            title: function(context) {
                                return context[0].label;
                            },
                            label: function(context) {
                                const value = context.parsed.y;
                                const strength = Math.abs(value) > 0.02 ? 'Moderate' : 'Weak';
                                return `Correlation: ${value.toFixed(3)} (${strength})`;
                            }
                        }
                    }
                },
                scales: {
                    x: {
                        grid: {
                            display: false
                        },
                        ticks: {
                            color: '#71717a',
                            font: {
                                size: 10
                            },
                            maxRotation: 45
                        }
                    },
                    y: {
                        grid: {
                            color: 'rgba(255, 255, 255, 0.1)'
                        },
                        ticks: {
                            color: '#71717a',
                            font: {
                                size: 10
                            },
                            callback: function(value) {
                                return value.toFixed(3);
                            }
                        }
                    }
                },
                interaction: {
                    intersect: false,
                    mode: 'index'
                }
            }
        });
    }
}

// Interactive elements
function initializeInteractivity() {
    // Demo chart controls
    const demoButtons = document.querySelectorAll('.demo-btn');
    demoButtons.forEach(button => {
        button.addEventListener('click', function() {
            // Remove active class from all buttons in the same group
            const group = this.closest('.demo-controls');
            group.querySelectorAll('.demo-btn').forEach(btn => btn.classList.remove('active'));
            
            // Add active class to clicked button
            this.classList.add('active');
            
            // Update chart based on selection
            const chartType = this.dataset.chart;
            updateFrequencyChart(chartType);
        });
    });
    
    // Smooth scrolling for anchor links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });
    
    // Animate numbers on scroll
    const observerOptions = {
        threshold: 0.5,
        rootMargin: '0px 0px -100px 0px'
    };
    
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                animateNumbers(entry.target);
                observer.unobserve(entry.target);
            }
        });
    }, observerOptions);
    
    // Observe stat numbers
    document.querySelectorAll('.stat-number, .score-value, .gauge-value').forEach(el => {
        observer.observe(el);
    });
    
    // Animate confidence bar
    const confidenceBars = document.querySelectorAll('.confidence-fill');
    confidenceBars.forEach(bar => {
        observer.observe(bar.closest('.confidence-meter'));
    });
}

// Animation functions
function initializeAnimations() {
    // Stagger animation for feature cards
    const featureCards = document.querySelectorAll('.feature-card');
    featureCards.forEach((card, index) => {
        card.style.animationDelay = `${index * 0.1}s`;
        card.classList.add('animate-fade-in-up');
    });
    
    // Parallax effect for hero background
    window.addEventListener('scroll', () => {
        const scrolled = window.pageYOffset;
        const parallax = document.querySelector('.hero-background');
        if (parallax) {
            const speed = scrolled * 0.5;
            parallax.style.transform = `translateY(${speed}px)`;
        }
    });
    
    // Update cosmic particles animation
    updateCosmicParticles();
}

// Chart update functions
function updateFrequencyChart(type) {
    if (!window.frequencyChart) return;
    
    let newData;
    let newColors;
    
    if (type === 'recent') {
        newData = generateMockRecentData();
        newColors = createGradientArray(48, '#06b6d4', '#8b5cf6');
    } else {
        newData = generateMockFrequencyData();
        newColors = createGradientArray(48, '#6366f1', '#8b5cf6');
    }
    
    window.frequencyChart.data.datasets[0].data = newData;
    window.frequencyChart.data.datasets[0].backgroundColor = newColors;
    window.frequencyChart.update('active');
}

// Data generation functions
function generateNumberLabels(start, end) {
    const labels = [];
    for (let i = start; i <= end; i++) {
        labels.push(i.toString());
    }
    return labels;
}

function generateMockFrequencyData() {
    const data = [];
    for (let i = 0; i < 48; i++) {
        // Generate realistic lottery frequency data with some variation
        const baseFreq = 85; // Average frequency
        const variation = Math.random() * 30 - 15; // ±15 variation
        data.push(Math.max(60, Math.floor(baseFreq + variation)));
    }
    return data;
}

function generateMockRecentData() {
    const data = [];
    for (let i = 0; i < 48; i++) {
        // Generate recent frequency data (lower numbers, more variation)
        const baseFreq = 3; // Average recent frequency
        const variation = Math.random() * 4 - 2; // ±2 variation
        data.push(Math.max(0, Math.floor(baseFreq + variation)));
    }
    return data;
}

function createGradientArray(length, startColor, endColor) {
    const colors = [];
    for (let i = 0; i < length; i++) {
        const ratio = i / (length - 1);
        colors.push(interpolateColor(startColor, endColor, ratio));
    }
    return colors;
}

function interpolateColor(color1, color2, ratio) {
    // Simple color interpolation
    const hex1 = color1.replace('#', '');
    const hex2 = color2.replace('#', '');
    
    const r1 = parseInt(hex1.substr(0, 2), 16);
    const g1 = parseInt(hex1.substr(2, 2), 16);
    const b1 = parseInt(hex1.substr(4, 2), 16);
    
    const r2 = parseInt(hex2.substr(0, 2), 16);
    const g2 = parseInt(hex2.substr(2, 2), 16);
    const b2 = parseInt(hex2.substr(4, 2), 16);
    
    const r = Math.round(r1 + (r2 - r1) * ratio);
    const g = Math.round(g1 + (g2 - g1) * ratio);
    const b = Math.round(b1 + (b2 - b1) * ratio);
    
    return `rgb(${r}, ${g}, ${b})`;
}

// Animation utilities
function animateNumbers(element) {
    const text = element.textContent;
    const number = parseFloat(text.replace(/[^\d.]/g, ''));
    
    if (isNaN(number)) return;
    
    const duration = 2000;
    const steps = 60;
    const increment = number / steps;
    let current = 0;
    let step = 0;
    
    const timer = setInterval(() => {
        current += increment;
        step++;
        
        if (step >= steps) {
            current = number;
            clearInterval(timer);
        }
        
        // Format the number based on original text
        let formattedNumber;
        if (text.includes('%')) {
            formattedNumber = current.toFixed(1) + '%';
        } else if (text.includes('+')) {
            formattedNumber = Math.floor(current).toLocaleString() + '+';
        } else {
            formattedNumber = Math.floor(current).toLocaleString();
        }
        
        element.textContent = formattedNumber;
    }, duration / steps);
}

function updateCosmicParticles() {
    const particles = document.querySelector('.cosmic-particles');
    if (!particles) return;
    
    // Add dynamic cosmic effects based on time
    const now = new Date();
    const seconds = now.getSeconds();
    const intensity = (Math.sin(seconds * 0.1) + 1) * 0.5; // 0 to 1
    
    particles.style.opacity = 0.3 + (intensity * 0.4);
    
    // Update every second
    setTimeout(updateCosmicParticles, 1000);
}

// Utility functions
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Performance optimization
const debouncedScroll = debounce(() => {
    // Handle scroll events here if needed
}, 16); // ~60fps

window.addEventListener('scroll', debouncedScroll);

// Error handling for charts
Chart.defaults.plugins.tooltip.enabled = true;
Chart.defaults.plugins.tooltip.external = function(context) {
    // Custom tooltip handling if needed
};

// Accessibility improvements
document.addEventListener('keydown', function(e) {
    // Handle keyboard navigation
    if (e.key === 'Tab') {
        document.body.classList.add('keyboard-navigation');
    }
});

document.addEventListener('mousedown', function() {
    document.body.classList.remove('keyboard-navigation');
});

// Console easter egg
console.log(`
🌌 Go-Lucky Lottery Analyzer
═══════════════════════════════════════════════════════════════════

Welcome to the cosmic correlation analysis tool!

The universe is under no obligation to make you wealthy,
but it's happy to teach you about statistics! 🌟

Remember: This tool is for educational and entertainment purposes only.
Lottery drawings are random events, and no analysis can predict future outcomes.

Check out the source code: https://github.com/mrz1836/go-lucky
`);

// Export functions for potential external use
window.GoLucky = {
    updateFrequencyChart,
    animateNumbers,
    generateMockFrequencyData,
    generateMockRecentData
};