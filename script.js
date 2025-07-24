// Initialize Lucide icons
document.addEventListener('DOMContentLoaded', function() {
    lucide.createIcons();
    
    // Initialize interactive elements
    initializeInteractivity();
    
    // Initialize animations
    initializeAnimations();
});

// Interactive elements
function initializeInteractivity() {
    // Strategy tab switching
    const strategyTabs = document.querySelectorAll('.strategy-tab');
    const strategySets = document.querySelectorAll('.strategy-set');
    
    strategyTabs.forEach(tab => {
        tab.addEventListener('click', function() {
            const strategy = this.dataset.strategy;
            
            // Remove active class from all tabs and sets
            strategyTabs.forEach(t => t.classList.remove('active'));
            strategySets.forEach(s => s.classList.remove('active'));
            
            // Add active class to clicked tab and corresponding set
            this.classList.add('active');
            const targetSet = document.querySelector(`.strategy-set[data-strategy="${strategy}"]`);
            if (targetSet) {
                targetSet.classList.add('active');
            }
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
                if (entry.target.classList.contains('stat-number')) {
                    animateNumbers(entry.target);
                } else if (entry.target.classList.contains('confidence-meter')) {
                    animateConfidenceBars(entry.target);
                }
                observer.unobserve(entry.target);
            }
        });
    }, observerOptions);
    
    // Observe stat numbers and confidence meters
    document.querySelectorAll('.stat-number').forEach(el => {
        observer.observe(el);
    });
    
    document.querySelectorAll('.confidence-meter').forEach(el => {
        observer.observe(el);
    });
    
    // Cosmic wisdom quote rotation
    initializeWisdomRotation();
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
    let ticking = false;
    
    function updateParallax() {
        const scrolled = window.pageYOffset;
        const parallax = document.querySelector('.hero-background');
        if (parallax) {
            const speed = scrolled * 0.5;
            parallax.style.transform = `translateY(${speed}px)`;
        }
        ticking = false;
    }
    
    function requestTick() {
        if (!ticking) {
            requestAnimationFrame(updateParallax);
            ticking = true;
        }
    }
    
    window.addEventListener('scroll', requestTick);
    
    // Update cosmic particles animation
    updateCosmicParticles();
    
    // Animate numbers on page load for visible elements
    setTimeout(() => {
        const visibleStats = document.querySelectorAll('.stat-number');
        visibleStats.forEach(stat => {
            const rect = stat.getBoundingClientRect();
            if (rect.top < window.innerHeight && rect.bottom > 0) {
                animateNumbers(stat);
            }
        });
    }, 500);
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
            formattedNumber = current.toFixed(0) + '%';
        } else if (text.includes('+')) {
            formattedNumber = Math.floor(current).toLocaleString() + '+';
        } else {
            formattedNumber = Math.floor(current).toLocaleString();
        }
        
        element.textContent = formattedNumber;
    }, duration / steps);
}

function animateConfidenceBars(meter) {
    const bars = meter.querySelectorAll('.confidence-fill');
    bars.forEach(bar => {
        const width = bar.style.width;
        bar.style.width = '0%';
        setTimeout(() => {
            bar.style.width = width;
        }, 100);
    });
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

function initializeWisdomRotation() {
    const quotes = document.querySelectorAll('.wisdom-quote');
    let currentQuote = 0;
    
    function rotateQuotes() {
        // Remove active class from current quote
        quotes[currentQuote].classList.remove('active');
        
        // Move to next quote
        currentQuote = (currentQuote + 1) % quotes.length;
        
        // Add active class to new quote
        quotes[currentQuote].classList.add('active');
    }
    
    // Rotate quotes every 4 seconds
    setInterval(rotateQuotes, 4000);
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

// Add some interactive number generation
function generateRandomCosmicNumbers() {
    const numbers = [];
    while (numbers.length < 5) {
        const num = Math.floor(Math.random() * 48) + 1;
        if (!numbers.includes(num)) {
            numbers.push(num);
        }
    }
    return numbers.sort((a, b) => a - b);
}

function generateRandomLuckyBall() {
    return Math.floor(Math.random() * 18) + 1;
}

// Update cosmic selection periodically (for demo purposes)
function updateCosmicSelection() {
    const cosmicNumbers = generateRandomCosmicNumbers();
    const luckyBall = generateRandomLuckyBall();
    
    const mainNumbers = document.querySelector('.cosmic-selection .main-numbers');
    const luckyBallElement = document.querySelector('.cosmic-selection .lucky');
    
    if (mainNumbers && luckyBallElement) {
        // Animate out
        mainNumbers.style.opacity = '0.5';
        luckyBallElement.style.opacity = '0.5';
        
        setTimeout(() => {
            // Update numbers
            const numberElements = mainNumbers.querySelectorAll('.number');
            cosmicNumbers.forEach((num, index) => {
                if (numberElements[index]) {
                    numberElements[index].textContent = num.toString().padStart(2, '0');
                }
            });
            
            luckyBallElement.textContent = luckyBall.toString();
            
            // Animate in
            mainNumbers.style.opacity = '1';
            luckyBallElement.style.opacity = '1';
        }, 300);
    }
}

// Update cosmic selection every 30 seconds for demo
setInterval(updateCosmicSelection, 30000);

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

Current cosmic conditions:
🌙 Moon Phase: New Moon (1% illuminated)
♌ Zodiac Sign: Leo
📅 Day: Thursday

🎯 Today's cosmic suggestion: New moon periods show balanced number distribution (1% illuminated).
A mix of high and low numbers may be favorable.
`);

// Export functions for potential external use
window.GoLucky = {
    animateNumbers,
    updateCosmicSelection,
    generateRandomCosmicNumbers,
    generateRandomLuckyBall
};