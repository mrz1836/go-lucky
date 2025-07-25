// Initialize Lucide icons
document.addEventListener('DOMContentLoaded', function() {
    lucide.createIcons();
    
    // Initialize interactive elements
    initializeInteractivity();
    
    // Initialize animations
    initializeAnimations();
    
    // Initialize tracking
    initializeTracking();
    
    // Initialize version display
    initializeVersionDisplay();
    
    // Initialize dynamic date
    updateCosmicDate();
    
    // Initialize copyright year
    updateCopyrightYear();
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
    // Animate numbers on page load for visible elements
    setTimeout(() => {
        const visibleStats = document.querySelectorAll('.stat-number');
        visibleStats.forEach(stat => {
            const rect = stat.getBoundingClientRect();
            if (rect.top < window.innerHeight && rect.bottom > 0) {
                // Don't animate the hero stats - they should show correct values
                if (!stat.closest('.hero-stats')) {
                    animateNumbers(stat);
                }
            }
        });
    }, 500);
}

// Animation utilities
function animateNumbers(element) {
    const text = element.textContent;
    
    // Skip animation for hero stats - they should show the correct static values
    if (element.closest('.hero-stats')) {
        return;
    }
    
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

// Version Display Functions
function initializeVersionDisplay() {
    console.log('🔍 Initializing GitHub version display...');
    fetchLatestVersion();
}

async function fetchLatestVersion() {
    const versionElement = document.getElementById('version-display');
    const versionLinkIcon = document.getElementById('version-link');
    
    if (!versionElement) {
        console.log('⚠️ Version display element not found');
        return;
    }
    
    try {
        // Add loading state
        versionElement.classList.add('loading');
        versionElement.textContent = 'Loading...';
        
        console.log('📡 Fetching latest release from GitHub API...');
        
        // Fetch latest release from GitHub API
        const response = await fetch('https://api.github.com/repos/mrz1836/go-lucky/releases/latest', {
            headers: {
                'Accept': 'application/vnd.github.v3+json',
                'User-Agent': 'Go-Lucky-Website'
            }
        });
        
        if (!response.ok) {
            throw new Error(`GitHub API responded with ${response.status}: ${response.statusText}`);
        }
        
        const release = await response.json();
        
        // Extract version info
        const version = release.tag_name || 'Unknown';
        const releaseUrl = release.html_url;
        const publishedAt = new Date(release.published_at);
        const isPrerelease = release.prerelease;
        
        // Update display
        versionElement.classList.remove('loading');
        versionElement.textContent = version;
        versionElement.title = `Released ${publishedAt.toLocaleDateString()}${isPrerelease ? ' (Pre-release)' : ''}`;
        
        // Update both the version number link and icon link to point to specific release
        if (releaseUrl) {
            versionElement.href = releaseUrl;
            versionElement.title = `View ${version} release notes`;
            
            if (versionLinkIcon) {
                versionLinkIcon.href = releaseUrl;
                versionLinkIcon.title = `View ${version} release notes`;
            }
        }
        
        // Add prerelease styling if applicable
        if (isPrerelease) {
            versionElement.style.background = 'rgba(245, 158, 11, 0.1)';
            versionElement.style.borderColor = 'rgba(245, 158, 11, 0.2)';
            versionElement.style.color = 'var(--warning)';
        }
        
        console.log(`✅ Version display updated: ${version}`);
        
        // Track version fetch success
        if (typeof gtag !== 'undefined') {
            gtag('event', 'version_fetch_success', {
                event_category: 'Technical',
                event_label: version,
                version: version,
                is_prerelease: isPrerelease,
                value: 1
            });
        }
        
    } catch (error) {
        console.error('❌ Failed to fetch version:', error);
        
        // Show error state
        versionElement.classList.remove('loading');
        versionElement.classList.add('error');
        versionElement.textContent = 'v1.0.0'; // Fallback version
        versionElement.title = 'Could not fetch latest version from GitHub';
        
        // Track version fetch error
        if (typeof gtag !== 'undefined') {
            gtag('event', 'version_fetch_error', {
                event_category: 'Technical',
                event_label: error.message,
                error_type: 'github_api_error',
                value: 0
            });
        }
    }
}

// Refresh version periodically (every 5 minutes) for long-running sessions
setInterval(() => {
    console.log('🔄 Refreshing version display...');
    fetchLatestVersion();
}, 5 * 60 * 1000);

// Google Analytics / GTM Tracking Functions
function initializeTracking() {
    console.log('🔍 Initializing tracking for external links and buttons...');
    
    // Track all external links
    trackExternalLinks();
    
    // Track specific button interactions
    trackButtonClicks();
    
    // Track scroll depth
    trackScrollDepth();
    
    // Track time on page milestones
    trackTimeOnPage();
    
    console.log('✅ Tracking initialized successfully');
}

function trackExternalLinks() {
    // Get all external links (links that go to different domains or specific external sites)
    const externalLinks = document.querySelectorAll('a[href^="http"], a[href^="https"], a[target="_blank"]');
    
    externalLinks.forEach(link => {
        const href = link.getAttribute('href');
        const text = link.textContent.trim();
        
        // Skip if it's an internal link to the same domain
        if (href && (href.includes(window.location.hostname) || href.startsWith('/'))) {
            return;
        }
        
        link.addEventListener('click', function(e) {
            const linkData = {
                event_category: 'External Link',
                event_label: href,
                link_text: text,
                link_url: href,
                link_domain: getDomainFromUrl(href),
                source_section: getSourceSection(this)
            };
            
            // Send to Google Analytics 4
            if (typeof gtag !== 'undefined') {
                gtag('event', 'click_external_link', {
                    event_category: linkData.event_category,
                    event_label: linkData.event_label,
                    link_text: linkData.link_text,
                    link_url: linkData.link_url,
                    link_domain: linkData.link_domain,
                    source_section: linkData.source_section,
                    value: 1
                });
            }
            
            // Send to Google Tag Manager (dataLayer)
            if (typeof dataLayer !== 'undefined') {
                dataLayer.push({
                    event: 'external_link_click',
                    link_category: linkData.event_category,
                    link_text: linkData.link_text,
                    link_url: linkData.link_url,
                    link_domain: linkData.link_domain,
                    source_section: linkData.source_section
                });
            }
            
            console.log('📊 External link tracked:', linkData);
        });
    });
    
    console.log(`🔗 Tracking ${externalLinks.length} external links`);
}

function trackButtonClicks() {
    // Track all buttons and CTA elements
    const buttons = document.querySelectorAll('.btn, button, .strategy-tab');
    
    buttons.forEach(button => {
        button.addEventListener('click', function(e) {
            const buttonText = this.textContent.trim();
            const buttonClass = this.className;
            const isExternal = this.getAttribute('href') && !this.getAttribute('href').startsWith('#');
            
            const buttonData = {
                event_category: 'Button Click',
                event_label: buttonText,
                button_type: getButtonType(this),
                button_class: buttonClass,
                source_section: getSourceSection(this),
                is_external: isExternal
            };
            
            // Send to Google Analytics 4
            if (typeof gtag !== 'undefined') {
                gtag('event', 'click_button', {
                    event_category: buttonData.event_category,
                    event_label: buttonData.event_label,
                    button_type: buttonData.button_type,
                    button_class: buttonData.button_class,
                    source_section: buttonData.source_section,
                    is_external: buttonData.is_external,
                    value: 1
                });
            }
            
            // Send to Google Tag Manager
            if (typeof dataLayer !== 'undefined') {
                dataLayer.push({
                    event: 'button_click',
                    button_category: buttonData.event_category,
                    button_text: buttonData.event_label,
                    button_type: buttonData.button_type,
                    button_class: buttonData.button_class,
                    source_section: buttonData.source_section,
                    is_external: buttonData.is_external
                });
            }
            
            console.log('🎯 Button click tracked:', buttonData);
        });
    });
    
    console.log(`🎯 Tracking ${buttons.length} buttons and interactive elements`);
}

function trackScrollDepth() {
    const scrollMilestones = [25, 50, 75, 90, 100];
    const reached = new Set();
    
    function checkScrollDepth() {
        const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        const docHeight = document.documentElement.scrollHeight - window.innerHeight;
        const scrollPercent = Math.round((scrollTop / docHeight) * 100);
        
        scrollMilestones.forEach(milestone => {
            if (scrollPercent >= milestone && !reached.has(milestone)) {
                reached.add(milestone);
                
                // Send to Google Analytics 4
                if (typeof gtag !== 'undefined') {
                    gtag('event', 'scroll_depth', {
                        event_category: 'Engagement',
                        event_label: `${milestone}%`,
                        scroll_depth: milestone,
                        value: milestone
                    });
                }
                
                // Send to Google Tag Manager
                if (typeof dataLayer !== 'undefined') {
                    dataLayer.push({
                        event: 'scroll_depth',
                        scroll_category: 'Engagement',
                        scroll_percentage: milestone,
                        scroll_label: `${milestone}%`
                    });
                }
                
                console.log(`📏 Scroll depth tracked: ${milestone}%`);
            }
        });
    }
    
    // Throttled scroll listener
    let scrollTimeout;
    window.addEventListener('scroll', function() {
        if (scrollTimeout) {
            clearTimeout(scrollTimeout);
        }
        scrollTimeout = setTimeout(checkScrollDepth, 100);
    });
}

function trackTimeOnPage() {
    const timeMilestones = [30, 60, 120, 300, 600]; // 30s, 1m, 2m, 5m, 10m
    const reached = new Set();
    const startTime = Date.now();
    
    function checkTimeOnPage() {
        const timeOnPage = Math.round((Date.now() - startTime) / 1000);
        
        timeMilestones.forEach(milestone => {
            if (timeOnPage >= milestone && !reached.has(milestone)) {
                reached.add(milestone);
                
                const minutes = Math.floor(milestone / 60);
                const seconds = milestone % 60;
                const timeLabel = minutes > 0 ? `${minutes}m${seconds > 0 ? ` ${seconds}s` : ''}` : `${seconds}s`;
                
                // Send to Google Analytics 4
                if (typeof gtag !== 'undefined') {
                    gtag('event', 'time_on_page', {
                        event_category: 'Engagement',
                        event_label: timeLabel,
                        time_seconds: milestone,
                        value: milestone
                    });
                }
                
                // Send to Google Tag Manager
                if (typeof dataLayer !== 'undefined') {
                    dataLayer.push({
                        event: 'time_on_page',
                        time_category: 'Engagement',
                        time_seconds: milestone,
                        time_label: timeLabel
                    });
                }
                
                console.log(`⏱️ Time on page tracked: ${timeLabel}`);
            }
        });
    }
    
    // Check time milestones every 10 seconds
    setInterval(checkTimeOnPage, 10000);
}

// Helper functions for tracking
function getDomainFromUrl(url) {
    try {
        return new URL(url).hostname;
    } catch (e) {
        return 'unknown';
    }
}

function getSourceSection(element) {
    // Determine which section of the page the element is in
    const section = element.closest('section, header, footer, nav');
    if (section) {
        // Try to get section class or ID
        if (section.className) {
            const classes = section.className.split(' ');
            const sectionClass = classes.find(cls => 
                ['hero', 'features', 'demo', 'usage', 'wisdom', 'disclaimer', 'footer'].includes(cls)
            );
            if (sectionClass) return sectionClass;
        }
        if (section.id) return section.id;
        return section.tagName.toLowerCase();
    }
    return 'unknown';
}

function getButtonType(button) {
    // Determine button type based on classes and context
    const classes = button.className.toLowerCase();
    
    if (classes.includes('btn-primary')) return 'primary';
    if (classes.includes('btn-secondary')) return 'secondary';
    if (classes.includes('strategy-tab')) return 'strategy-tab';
    if (button.tagName.toLowerCase() === 'button') return 'button';
    if (button.tagName.toLowerCase() === 'a') return 'link';
    
    return 'unknown';
}

// Track specific cosmic interactions
function trackCosmicInteraction(interactionType, details = {}) {
    const eventData = {
        event_category: 'Cosmic Interaction',
        event_label: interactionType,
        interaction_type: interactionType,
        ...details
    };
    
    // Send to Google Analytics 4
    if (typeof gtag !== 'undefined') {
        gtag('event', 'cosmic_interaction', eventData);
    }
    
    // Send to Google Tag Manager
    if (typeof dataLayer !== 'undefined') {
        dataLayer.push({
            event: 'cosmic_interaction',
            cosmic_category: eventData.event_category,
            cosmic_type: interactionType,
            ...details
        });
    }
    
    console.log('🌌 Cosmic interaction tracked:', eventData);
}

// Enhanced cosmic selection update with tracking
function updateCosmicSelection() {
    const cosmicNumbers = generateRandomCosmicNumbers();
    const luckyBall = generateRandomLuckyBall();
    
    const mainNumbers = document.querySelector('.cosmic-selection .main-numbers');
    const luckyBallElement = document.querySelector('.cosmic-selection .lucky');
    
    if (mainNumbers && luckyBallElement) {
        // Track the cosmic number generation
        trackCosmicInteraction('number_generation', {
            generated_numbers: cosmicNumbers.join('-'),
            lucky_ball: luckyBall,
            generation_type: 'automatic'
        });
        
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
    
    // Update the date as well
    updateCosmicDate();
}

// Update cosmic date to current browser time
function updateCosmicDate() {
    const cosmicDateElement = document.querySelector('.cosmic-date');
    if (cosmicDateElement) {
        const now = new Date();
        const options = { 
            year: 'numeric', 
            month: 'long', 
            day: 'numeric',
            timeZone: 'America/New_York' // Eastern Time for consistency
        };
        
        const formattedDate = now.toLocaleDateString('en-US', options);
        
        // Animate the date change
        cosmicDateElement.style.opacity = '0.7';
        setTimeout(() => {
            cosmicDateElement.textContent = formattedDate;
            cosmicDateElement.style.opacity = '1';
        }, 150);
        
        console.log(`🗓️ Cosmic date updated to: ${formattedDate}`);
    }
}

// Update cosmic date every minute to keep it current for long sessions
setInterval(updateCosmicDate, 60000);

// Update copyright year
function updateCopyrightYear() {
    const yearElement = document.getElementById('copyright-year');
    if (yearElement) {
        const currentYear = new Date().getFullYear();
        yearElement.textContent = currentYear;
        console.log(`© Copyright year updated to: ${currentYear}`);
    }
}