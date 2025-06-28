#!/bin/sh
set -e

# PairStair Release Script
# 
# Automates the complete release workflow with intelligent version calculation.
# Requires a completely clean working directory and passing tests.
#
# Features:
#   - Semantic version calculation (major/minor/patch)
#   - Manual version override with -v flag
#   - Clean working directory enforcement
#   - Automated testing and cleanup
#   - Git tag creation and pushing
#   - Non-interactive operation with command-line release notes
#
# Usage: ./release.sh <major|minor|patch|(-v version)> [release_notes]
# Documentation: See CONTRIBUTING.md for full release process details

# Automates the common release workflow for creating and pushing new versions

usage() {
    echo "Usage: $0 <release_type|version> [release_notes]"
    echo ""
    echo "Creates a new release by calculating the next version or using a specific version."
    echo ""
    echo "Release Types (recommended):"
    echo "  major          Increment major version (e.g., v0.7.2 -> v1.0.0)"
    echo "  minor          Increment minor version (e.g., v0.7.2 -> v0.8.0)"  
    echo "  patch          Increment patch version (e.g., v0.7.2 -> v0.7.3)"
    echo ""
    echo "Or specify exact version:"
    echo "  -v <version>   Use specific version number (e.g., -v v1.0.0)"
    echo ""
    echo "Arguments:"
    echo "  release_notes  Optional release notes (if not provided, uses default)"
    echo ""
    echo "Examples:"
    echo "  $0 patch                              # v0.7.2 -> v0.7.3"
    echo "  $0 minor 'Add new HTML streaming'     # v0.7.2 -> v0.8.0"
    echo "  $0 major 'Breaking API changes'      # v0.7.2 -> v1.0.0"
    echo "  $0 -v v2.0.0 'Complete rewrite'      # Use specific version"
    echo ""
    echo "The script will:"
    echo "  1. Calculate next version from latest git tag (or use specified version)"
    echo "  2. Verify working directory is completely clean"
    echo "  3. Run all tests to ensure they pass"
    echo "  4. Push any unpushed commits to origin"
    echo "  5. Create and push annotated git tag"
    echo ""
    exit 1
}

# Check arguments
if [ $# -lt 1 ] || [ $# -gt 2 ]; then
    echo "Error: Invalid number of arguments"
    usage
fi

RELEASE_TYPE_OR_VERSION="$1"
RELEASE_NOTES="$2"

# Function to calculate next version
calculate_next_version() {
    local release_type="$1"
    local current_version="$2"
    
    # Remove 'v' prefix and split version into parts
    local version_number=$(echo "$current_version" | sed 's/^v//')
    local major=$(echo "$version_number" | cut -d. -f1)
    local minor=$(echo "$version_number" | cut -d. -f2)
    local patch=$(echo "$version_number" | cut -d. -f3)
    
    case "$release_type" in
        "major")
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        "minor")
            minor=$((minor + 1))
            patch=0
            ;;
        "patch")
            patch=$((patch + 1))
            ;;
        *)
            echo "Error: Invalid release type '$release_type'. Use major, minor, or patch."
            exit 1
            ;;
    esac
    
    echo "v${major}.${minor}.${patch}"
}

# Determine version to use
if [ "$RELEASE_TYPE_OR_VERSION" = "-v" ]; then
    # Version flag mode - expect version as next argument
    if [ $# -lt 2 ]; then
        echo "Error: -v flag requires a version number"
        usage
    fi
    VERSION="$2"
    RELEASE_NOTES="$3"
    
    # Validate version format
    if ! echo "$VERSION" | grep -q '^v[0-9]\+\.[0-9]\+\.[0-9]\+$'; then
        echo "Error: Version must be in format vX.Y.Z (e.g., v0.7.1)"
        exit 1
    fi
    
    echo "🚀 Starting release process for specified version $VERSION"
elif echo "$RELEASE_TYPE_OR_VERSION" | grep -q '^v[0-9]\+\.[0-9]\+\.[0-9]\+$'; then
    # Legacy mode - full version provided directly
    VERSION="$RELEASE_TYPE_OR_VERSION"
    echo "🚀 Starting release process for specified version $VERSION"
else
    # Release type mode - calculate next version
    RELEASE_TYPE="$RELEASE_TYPE_OR_VERSION"
    
    # Check if we're in a git repository first
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        echo "Error: Not in a git repository"
        exit 1
    fi
    
    # Get current version from latest tag
    CURRENT_VERSION=$(git tag --sort=-version:refname | head -1 2>/dev/null)
    if [ -z "$CURRENT_VERSION" ]; then
        echo "Error: No existing version tags found. Use -v flag to specify initial version."
        echo "Example: $0 -v v0.1.0 'Initial release'"
        exit 1
    fi
    
    VERSION=$(calculate_next_version "$RELEASE_TYPE" "$CURRENT_VERSION")
    echo "🚀 Starting release process: $CURRENT_VERSION -> $VERSION ($RELEASE_TYPE)"
fi

# Continue with original validation logic
echo ""

# Check if we're in a git repository
if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo "Error: Not in a git repository"
    exit 1
fi

# Check if version tag already exists
if git tag -l | grep -q "^$VERSION$"; then
    echo "Error: Version tag $VERSION already exists"
    exit 1
fi

# Verify we're on master branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "master" ]; then
    echo "Warning: You are on branch '$CURRENT_BRANCH', not 'master'"
    echo "Continue anyway? (y/N)"
    read -r response
    if [ "$response" != "y" ] && [ "$response" != "Y" ]; then
        echo "Aborted"
        exit 1
    fi
fi

# Check for uncommitted changes (require completely clean working directory)
if git diff --quiet && git diff --cached --quiet; then
    echo "✓ Working directory is clean"
else
    echo "Error: Working directory has uncommitted changes. Please commit or stash them first."
    echo ""
    echo "All changes:"
    git status --porcelain
    echo ""
    echo "Tip: Use 'git add . && git commit -m \"prepare for release\"' to commit changes"
    exit 1
fi

# Run all tests
echo ""
echo "🧪 Running tests..."
if ! go test ./...; then
    echo "Error: Tests failed. Please fix them before releasing."
    exit 1
fi
echo "✓ All tests passed"

# Clean up any test artifacts
echo ""
echo "🧹 Cleaning up test artifacts..."
find . -name "*.test" -type f -delete 2>/dev/null || true
find . -name "pairstair" -type f -not -path "./pairstair" -delete 2>/dev/null || true
echo "✓ Test artifacts cleaned"

# Push commits to origin
echo ""
echo "⬆️  Pushing commits to origin..."
git push origin "$CURRENT_BRANCH"
echo "✓ Commits pushed"

# Get the latest tag for release notes context
PREVIOUS_TAG=$(git tag --sort=-version:refname | head -1 2>/dev/null || echo "")
if [ -n "$PREVIOUS_TAG" ]; then
    echo ""
    echo "📋 Previous version: $PREVIOUS_TAG"
    echo "📋 Commits since $PREVIOUS_TAG:"
    git log --oneline "$PREVIOUS_TAG"..HEAD | head -10
fi

# Create annotated tag with release notes
echo ""
echo "🏷️  Creating release tag $VERSION..."

# Use provided release notes or default
if [ -n "$RELEASE_NOTES" ]; then
    TAG_MESSAGE="$RELEASE_NOTES"
    echo "Using provided release notes"
else
    # Default tag message if none provided
    TAG_MESSAGE="$VERSION

Release $VERSION with latest changes and improvements."
    echo "Using default release notes"
fi

git tag -a "$VERSION" -m "$TAG_MESSAGE"
echo "✓ Tag created"

# Push the tag
echo ""
echo "🏷️  Pushing tag to origin..."
git push origin "$VERSION"
echo "✓ Tag pushed"

echo ""
echo "🎉 Release $VERSION completed successfully!"
echo ""
echo "The CI/CD pipeline will now:"
echo "  • Build binaries for multiple platforms"
echo "  • Create GitHub release with notes"
echo "  • Update Homebrew tap repository"
echo ""
echo "Monitor the release progress at:"
echo "  https://github.com/gypsydave5/pairstair/actions"
echo "  https://github.com/gypsydave5/pairstair/releases"
