# Go-Lucky Website

This is the showcase website for the Go-Lucky NC Lucky for Life Lottery Analyzer project.

## Overview

A modern, responsive static website that demonstrates the capabilities of the Go-Lucky lottery analysis tool. The site features:

- **Cutting-edge design** with cosmic-themed animations and data visualizations
- **Interactive charts** showing frequency distributions and cosmic correlations
- **Live demo dashboard** with mock data visualization
- **Responsive layout** that works on all devices
- **Performance optimized** with modern CSS and JavaScript

## Features

### Visual Design
- Dark cosmic theme with gradient backgrounds
- Animated particles and data grid effects
- Smooth transitions and hover effects
- Professional typography using Inter and JetBrains Mono fonts

### Interactive Elements
- Chart.js powered data visualizations
- Lucide icons for consistent iconography
- Smooth scrolling navigation
- Animated number counters
- Interactive demo controls

### Content Sections
1. **Hero Section** - Project introduction with key statistics
2. **Features** - Detailed breakdown of analysis capabilities
3. **Live Demo** - Interactive dashboard with charts and recommendations
4. **Usage** - Quick start guide with code examples
5. **Disclaimer** - Important educational notice
6. **Footer** - Links and project information

## Files

- `index.html` - Main HTML structure
- `styles.css` - Complete CSS styling with custom properties and responsive design
- `script.js` - JavaScript for interactivity, charts, and animations
- `README.md` - This documentation file

## Dependencies

The website uses the following external libraries via CDN:

- **Chart.js** - For data visualization charts
- **Lucide** - For consistent iconography
- **Google Fonts** - Inter and JetBrains Mono fonts

## Usage

Simply open `index.html` in a web browser or serve the files through any web server. The website is completely self-contained and doesn't require any build process.

### Local Development

```bash
# Navigate to the website directory
cd website

# Serve with Python (if available)
python -m http.server 8000

# Or use any other static file server
# Then visit http://localhost:8000
```

## Customization

The website uses CSS custom properties (variables) for easy theming. Key variables are defined in the `:root` selector in `styles.css`:

- Colors: Primary, secondary, accent, and cosmic theme colors
- Spacing: Consistent spacing scale
- Typography: Font families and sizes
- Animations: Transition durations and easing

## Performance

The website is optimized for performance with:

- Efficient CSS with minimal reflows
- Debounced scroll events
- Optimized animations using CSS transforms
- Lazy loading of chart data
- Minimal external dependencies

## Browser Support

The website supports all modern browsers including:

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Accessibility

The website includes accessibility features:

- Semantic HTML structure
- Proper heading hierarchy
- Focus indicators for keyboard navigation
- Alt text for visual elements
- High contrast color ratios
- Reduced motion support (respects user preferences)

## License

This website is part of the Go-Lucky project and follows the same licensing terms.