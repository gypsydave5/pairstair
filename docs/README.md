# PairStair Documentation Site

This directory contains the Jekyll site for PairStair documentation, hosted on GitHub Pages.

## Site Structure

- `_config.yml` - Jekyll configuration
- `_layouts/` - Custom page layouts
- `index.md` - Home page
- `installation.md` - Installation guide
- `guide.md` - Complete user guide
- `examples.md` - Usage examples and scenarios
- `features.md` - Current and planned features

## Local Development

To run the site locally:

```bash
# Install Jekyll (one time setup)
gem install bundler jekyll

# Navigate to docs directory
cd docs/

# Create Gemfile if it doesn't exist
cat > Gemfile << 'EOF'
source "https://rubygems.org"
gem "github-pages", group: :jekyll_plugins
gem "webrick", "~> 1.7"
EOF

# Install dependencies
bundle install

# Serve the site locally
bundle exec jekyll serve

# Open http://localhost:4000/pairstair in your browser
```

## Publishing

The site is automatically published to GitHub Pages when changes are pushed to the `main` branch. GitHub Actions handles the Jekyll build and deployment.

## Content Guidelines

### Writing Style
- Use clear, concise language
- Include practical examples
- Assume readers are developers but may be new to pairing
- Use consistent terminology throughout

### Code Examples
- Always test code examples before publishing
- Use realistic data in examples
- Include both command and expected output
- Explain what each example demonstrates

### Navigation
- Keep the main navigation simple (4-5 top-level pages)
- Use clear, descriptive page titles
- Link between related sections
- Include "back to top" links on long pages

### Updates
When adding new features to PairStair:

1. Update the relevant documentation pages
2. Add examples if applicable
3. Update the features page
4. Consider if new navigation items are needed
5. Test the site locally before committing

## File Organization

```
docs/
├── _config.yml          # Jekyll configuration
├── _layouts/
│   └── home.html        # Custom home page layout
├── index.md             # Home page content
├── installation.md      # Installation instructions
├── guide.md            # Complete user guide
├── examples.md         # Usage examples
├── features.md         # Feature descriptions
├── gen-man.sh          # Man page generator script
├── pairstair.1         # Generated man page
├── pairstair.1.md      # Man page source
└── README.md           # This file
```

## Theme Customization

The site uses the `minima` theme with some customizations:

- Custom home page layout with feature cards
- Enhanced styling for code blocks and examples
- Responsive design for mobile devices
- Consistent color scheme and typography

To modify the theme:

1. Override theme files by creating files in the same path structure
2. Add custom CSS to `_sass/` directory
3. Modify layouts in `_layouts/` directory
4. Test changes locally before committing

## SEO and Performance

The site includes:

- SEO tags and metadata
- Sitemap generation
- Social media integration
- Fast loading times
- Mobile-friendly responsive design

## Maintenance

Regular maintenance tasks:

- Review and update documentation when PairStair features change
- Check for broken links
- Ensure examples remain current and accurate
- Update version numbers and release information
- Monitor site performance and accessibility
