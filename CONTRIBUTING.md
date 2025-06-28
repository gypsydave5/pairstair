# Contributing to PairStair

Thank you for your interest in contributing to PairStair! This document provides guidelines for development, testing, and the release process.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- pandoc (for building man pages)

### Building Locally

```bash
# Clone the repository
git clone https://github.com/gypsydave5/pairstair.git
cd pairstair

# Build the binary
go build -o pairstair .

# Run tests
go test -v ./...

# Build with version injection (simulates release builds)
go build -ldflags "-X main.Version=v1.0.0-local" -o pairstair .
```

### Version Detection System

PairStair uses a sophisticated version detection system:

1. **Development builds**: Show git information (tags, commit hashes, dirty state)
2. **Release builds**: Show the injected version from CI/CD pipeline
3. **Fallback**: Use the `Version` constant in `pairstair.go`

The version detection priority (highest to lowest):
1. Clean git tag (e.g., `v0.5.1`)
2. Dirty git tag (e.g., `v0.5.1-dirty`)
3. Version constant + commit hash (e.g., `0.5.0-dev+abc12345`)
4. Module version from `go.mod`
5. Fallback to `Version` constant

## Testing

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific test
go test -v -run TestCoAuthorPairingDetection
```

### Test-Driven Development

We follow strict TDD practices:

1. **Write failing tests first** for any new functionality
2. **Make tests pass** with minimal implementation
3. **Refactor** while keeping tests green
4. **Verify** existing tests still pass

See `.github/copilot-instructions.md` for detailed TDD guidelines.

## Development Scripts

PairStair includes helper scripts to streamline common development tasks:

### dev.sh - Development Helper

Simple helper script for common development tasks:

```bash
# Show version information (builds temporarily if needed)
./dev.sh version

# Regenerate man page documentation
./dev.sh docs

# Show help
./dev.sh
```

**Use standard Go commands for basic operations**:
- `go test ./...` (instead of a test script)
- `go build` (instead of a build script)
- Manual cleanup as needed

### release.sh - Release Automation

Comprehensive script that automates the entire release process with intelligent version calculation:

```bash
# Semantic version releases (recommended)
./release.sh patch                    # v0.7.2 -> v0.7.3
./release.sh minor "New HTML feature" # v0.7.2 -> v0.8.0  
./release.sh major "Breaking changes" # v0.7.2 -> v1.0.0

# Explicit version (when needed)
./release.sh -v v2.0.0 "Complete rewrite"

# Show help
./release.sh
```

**What the release script does**:
1. Calculates next version from latest git tag (or uses specified version)
2. Validates version format and checks for existing tags
3. Checks working directory is completely clean (no uncommitted changes)
4. Runs all tests to ensure they pass
5. Pushes any unpushed commits to origin
6. Creates annotated git tag with release notes
7. Pushes the tag to trigger CI/CD pipeline

**Version calculation**:
- **patch**: Bug fixes and small improvements (0.7.2 -> 0.7.3)
- **minor**: New features, backward compatible (0.7.2 -> 0.8.0)
- **major**: Breaking changes (0.7.2 -> 1.0.0)
- **-v flag**: Specify exact version when automatic calculation isn't suitable

**Requirements**:
- Clean working directory (no staged or unstaged changes)
- All tests must pass
- Valid semantic version format

**Safety features**:
- Aborts if uncommitted changes exist
- Aborts if version tag already exists
- Provides helpful error messages and suggestions

## Build and Release Process

### CI/CD Pipeline

The project uses GitHub Actions for automated building and releasing:

**File**: `.github/workflows/release.yml`

#### Trigger Conditions

1. **Tag push**: Pushing a tag like `v1.0.0` triggers a release
2. **Manual dispatch**: Can be triggered manually with a custom version

#### Build Process

1. **Test Stage**: Run all tests across the codebase
2. **Build Stage**: 
   - Extract version from git tag or manual input
   - Build binaries for multiple platforms using version injection:
     ```bash
     go build -ldflags "-X main.Version=${VERSION}" -o pairstair .
     ```
   - Supported platforms: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`
3. **Man Page Stage**: Generate man page using pandoc
4. **Release Stage**: Create GitHub release with all artifacts

#### Version Injection

The CI/CD pipeline injects the correct version into binaries using Go's `-ldflags`:

```bash
# Extract version from tag (e.g., v0.5.1)
VERSION="${GITHUB_REF#refs/tags/}"

# Inject version into binary
go build -ldflags "-X main.Version=${VERSION}" -o pairstair .
```

This ensures that pre-built binaries (Homebrew packages, GitHub releases) show the correct version instead of the development fallback.

### Manual Release Process

1. **Ensure all tests pass**:
   ```bash
   go test -v ./...
   ```

2. **Update version constant** (optional, for development reference):
   ```go
   // In pairstair.go
   const Version = "0.6.0-dev"  // Next development version
   ```

3. **Commit and push changes**:
   ```bash
   git add .
   git commit -m "-s- prepare for v0.6.0 release"
   git push origin master
   ```

4. **Create and push release tag**:
   ```bash
   git tag v0.6.0 -m "Release v0.6.0

   Brief description of changes in this release.
   
   Features:
   - New feature descriptions
   
   Fixes:
   - Bug fix descriptions
   
   Breaking Changes:
   - Any breaking changes"
   
   git push origin v0.6.0
   ```

5. **Monitor CI/CD pipeline**: Check GitHub Actions for build status

6. **Verify release**: Once complete, check that:
   - GitHub release was created with all binary artifacts
   - Homebrew formula was updated (see below)
   - Version information is correct in built binaries

### Homebrew Integration

PairStair has a separate Homebrew tap repository for macOS distribution:

**Repository**: `gypsydave5/homebrew-pairstair`

#### How It Works

1. **Automatic trigger**: When a GitHub release is created, the CI/CD pipeline triggers an update to the Homebrew formula
2. **Repository dispatch**: Uses `peter-evans/repository-dispatch` action to notify the Homebrew repository
3. **Formula update**: The Homebrew repository automatically updates the formula with:
   - New version number
   - Updated download URLs
   - New checksums for the release artifacts

#### Manual Homebrew Updates

If the automatic update fails, manually update the Homebrew formula:

1. Clone the Homebrew tap repository:
   ```bash
   git clone https://github.com/gypsydave5/homebrew-pairstair.git
   ```

2. Update the formula in `Formula/pairstair.rb`:
   - Update version number
   - Update download URL
   - Update SHA256 checksums

3. Test the formula:
   ```bash
   brew install --build-from-source ./Formula/pairstair.rb
   brew test pairstair
   ```

4. Commit and push changes

#### Installation Methods

Users can install PairStair through multiple methods:

```bash
# Homebrew (macOS)
brew tap gypsydave5/pairstair
brew install pairstair

# Go install (any platform)
go install github.com/gypsydave5/pairstair@latest

# Manual download (any platform)
# Download from GitHub Releases page
```

## Documentation

### Documentation Files

- `README.md`: Primary user documentation
- `docs/pairstair.1.md`: Man page source (Markdown)
- `docs/pairstair.1`: Generated man page
- `CONTRIBUTING.md`: Development and contribution guidelines
- `.github/copilot-instructions.md`: AI assistant guidelines and project conventions

### Updating Documentation

When adding features or making changes:

1. **Update README.md** with user-facing changes
2. **Update man page source** in `docs/pairstair.1.md`
3. **Regenerate man page**:
   ```bash
   cd docs
   ./gen-man.sh
   ```
4. **Commit both source and generated files**

### Documentation Testing

Before committing documentation changes:

1. **Verify examples work**:
   ```bash
   # Test documented commands
   ./pairstair --help
   ./pairstair --version
   ./pairstair .team
   ```

2. **Test man page generation**:
   ```bash
   cd docs
   ./gen-man.sh
   man ./pairstair.1  # Verify formatting
   ```

3. **Check help output matches documentation**

## Commit Message Conventions

Use consistent prefixes to indicate change types:

- **`-s-`**: Structural changes (refactoring, documentation, tests, CI/CD)
- **`-b-`**: Behavioral changes (features, bug fixes, API changes)

Examples:
```
-s- update CI/CD pipeline to support new platforms
-s- add tests for edge case in team file parsing
-s- update README with new installation instructions
-b- add --strategy flag for pairing recommendations
-b- fix incorrect handling of co-authored commits
```

## Code Style and Standards

- **Follow Go conventions**: Use `gofmt`, follow standard Go naming conventions
- **Keep functions short**: Ideally 10 lines or less
- **Write comprehensive tests**: Test all public functions and edge cases
- **Document public APIs**: Use Go doc comments for exported functions
- **Separate concerns**: Use dependency injection for testability

## Getting Help

- **Issues**: Report bugs or request features on GitHub Issues
- **Discussions**: Use GitHub Discussions for questions and general discussion
- **Development**: Follow the TDD guidelines in `.github/copilot-instructions.md`

## Release Checklist

- [ ] All tests pass locally
- [ ] Documentation updated
- [ ] Man page regenerated
- [ ] Version constant updated (if desired)
- [ ] Changes committed and pushed
- [ ] Release tag created and pushed
- [ ] CI/CD pipeline completes successfully
- [ ] GitHub release created with artifacts
- [ ] Homebrew formula updated
- [ ] Installation methods tested
