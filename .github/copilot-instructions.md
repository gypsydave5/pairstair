# Instructions to Copilot for writing Go (Golang) code.

We favour following idiomatic Go practices, using the standard library where possible, and writing clear, maintainable code. Here are some guidelines:

### Feature Development Approach

When the user requests a new feature, follow this collaborative approach:

1. **Interrogate the requirements**: Don't assume you understand the full scope. Ask clarifying questions about:
   - Specific behavior expected
   - Edge cases and error scenarios
   - User interface preferences (CLI flags, output format, etc.)
   - Performance or compatibility requirements

2. **Provide implementation options**: Present 2-3 different approaches with trade-offs:
   ```
   "I see a few ways to implement this:
   
   Option 1: Simple approach - [description with pros/cons]
   Option 2: Robust approach - [description with pros/cons]  
   Option 3: Full-featured approach - [description with pros/cons]
   
   Which direction would you prefer?"
   ```

3. **Confirm the approach**: Wait for the user to choose before proceeding with implementation.

4. **Follow TDD rigorously**: Write tests first, then implement (see detailed TDD guidelines below).

This collaborative approach ensures we build exactly what the user needs and avoids over-engineering or misunderstood requirements.

### Functions

Functions should be kept short, ideally no more than ten lines of code. They should do one thing and do it well. If a function is getting too long, consider breaking it down into smaller functions. The name of a function should reveal its purpose clearly.

### Structs

Collections of behaviour should be encapsulated in a named struct type. Access to the data should be mostly through methods.

Structs should for the most part be initialized using a public function starting with the word `New`. I.e. for a struct of type `Bob` the function should be called `NewBob`. Other constructor functions like this taking different arguments to construct the same struct are permissible and encouraged.

### Modules

Modules should be used to group related functionality. Each module should have a clear purpose and should not contain unrelated code. The module name should be descriptive of its functionality.

### Tests

Tests should be written for all public functions and methods. Use the `testing` package from the standard library. Tests should be placed in a file ending with `_test.go`. Each test function should start with `Test` followed by the name of the function being tested.

### Test-Driven Development and Change Management

When making changes to existing code, follow these critical practices:

1. **Always establish a baseline first**: Before making any changes, run all tests to ensure they pass. This provides a known good state to compare against.

2. **Run tests after each logical change**: After making modifications, immediately run tests to verify nothing was broken. If tests fail, assume the failure is related to your recent changes unless proven otherwise.

3. **Be extremely careful when updating test code**: When modifying function signatures that require test updates:
   - Update one test file at a time
   - Run tests after each file update
   - Pay special attention to copy-paste errors when updating multiple similar test calls
   - Avoid introducing new hardcoded values - use existing variables when possible

4. **When tests fail**: 
   - First check if your changes broke existing functionality
   - Look for typos or copy-paste errors in test modifications
   - Don't assume test failures are "pre-existing bugs" without verification
   - Use git to compare against the last known working state

5. **For function signature changes**: 
   - Search for all usages before making the change
   - Update all call sites systematically
   - Verify each update preserves the original intent
   - Use the compiler to help find missed updates

6. **When adding new features**:
   - Write tests for the new functionality first
   - Ensure existing tests still pass throughout development
   - Add integration tests to verify the feature works end-to-end

### True Test-Driven Development (TDD) Guidelines

**Critical: Always follow the Red-Green-Refactor cycle when adding new features**

#### When the user requests a new feature:

1. **Clarify requirements first**:
   - Ask detailed questions about the expected behavior
   - Provide multiple implementation options if appropriate
   - Discuss edge cases and error scenarios
   - Confirm the user interface (CLI flags, output format, etc.)

2. **Write acceptance tests FIRST**:
   - **Start with end-to-end acceptance tests** that verify the complete user experience
   - Use the actual binary and test real command-line behavior
   - Test the feature as users will actually interact with it
   - Example: For triple recommendations, write tests that check for `<->` notation in CLI output
   - Acceptance tests should fail initially and guide the implementation

3. **Write failing unit tests SECOND**:
   - Create comprehensive test cases that cover the main functionality
   - Include edge cases and error conditions
   - Test the actual production code interface, not test helper functions
   - Ensure tests fail for the right reason (missing functionality, not syntax errors)

3. **Write failing unit tests SECOND**:
   - Create comprehensive test cases that cover the main functionality
   - Include edge cases and error conditions
   - Test the actual production code interface, not test helper functions
   - Ensure tests fail for the right reason (missing functionality, not syntax errors)

4. **Make tests testable by design**:
   - Separate logic from external dependencies (file I/O, network, etc.)
   - Use dependency injection where appropriate
   - Create functions that take parameters instead of accessing globals
   - Example: `getVersionFromBuildInfo(info, hasInfo)` vs `getVersion()` that calls `debug.ReadBuildInfo()`

5. **Verify tests fail correctly**:
   - Run tests to confirm they fail as expected
   - Fix any compilation errors in tests
   - Ensure failure messages are clear and helpful

6. **Implement minimal production code**:
   - Write just enough code to make tests pass
   - Don't over-engineer or add features not covered by tests
   - Keep functions focused and single-purpose

7. **Verify tests pass**:
   - Run tests to confirm they now pass
   - Run all tests to ensure no regressions
   - **Trust passing tests - avoid redundant manual verification if acceptance tests exist**

8. **Refactor if needed**:
   - Clean up code while keeping tests green
   - Extract functions, improve naming, etc.
   - Re-run tests after each refactoring step

#### Common TDD Anti-Patterns to Avoid:

- **❌ Writing test code that duplicates production logic**: Tests should call actual production functions, not reimplementations
- **❌ Writing production code first**: Always write the test first to drive the design
- **❌ Writing tests that always pass**: Tests should fail initially and only pass after implementing the feature
- **❌ Testing implementation details**: Test behavior and interfaces, not internal implementation
- **❌ Skipping the refactor step**: Clean code is as important as working code

#### Feature Development Process:

1. **Requirements Gathering**:
   ```
   User: "Add version support"
   
   Copilot: "I'd like to clarify the requirements. Here are some options:
   
   Option 1: Simple version constant
   Option 2: Version constant + --version flag  
   Option 3: Smart version detection using build info + git tags
   
   Questions:
   - Should it show git commit hashes for development builds?
   - How should it handle dirty working directories?
   - What format do you prefer for the version output?"
   ```

2. **Test Design**:
   ```
   "Let me write acceptance tests first to define the expected user experience:
   - Test actual CLI commands and output format
   - Test the complete user workflow end-to-end
   - Verify the feature works as users will interact with it
   
   Then write unit tests for the implementation details:
   - Test version flag parsing
   - Test version detection with different build scenarios
   - Test fallback behavior when build info unavailable"
   ```

3. **Implementation**:
   ```
   "Now I'll implement the minimal code to make these tests pass,
   ensuring the production code is testable and focused."
   ```

### Acceptance Test Requirements

**Always write acceptance tests BEFORE implementing user-facing features**

#### When to Write Acceptance Tests:
- **CLI features**: New flags, commands, or output formats
- **User workflows**: Complete end-to-end functionality
- **Output changes**: Modified display formats or new information
- **Behavior changes**: Any change that affects how users interact with the tool

#### Acceptance Test Characteristics:
1. **Test the actual binary**: Use `go build` and execute the real program
2. **Test real command-line usage**: Use actual flags and arguments
3. **Verify complete output**: Check exact text, formatting, and behavior
4. **Test error scenarios**: Invalid inputs, missing files, etc.
5. **Use realistic test data**: Create git repositories with real commit history

#### Example Pattern for New Features:
```go
// Test for triple recommendations feature
func TestTripleRecommendations(t *testing.T) {
    // Setup: Create git repo with 5 developers (odd number)
    repoDir := setupGitRepoWithOddDevelopers(t)
    defer os.RemoveAll(repoDir)
    
    // Execute: Run pairstair with strategy
    cmd := exec.Command("./pairstair", "--strategy=least-paired")
    cmd.Dir = repoDir
    output, err := cmd.CombinedOutput()
    
    // Verify: Check for triple notation in output
    if !strings.Contains(string(output), "<->") {
        t.Errorf("Expected triple notation '<->' in output")
    }
    
    // Verify: Ensure exactly one triple and pairs for the rest
    lines := strings.Split(string(output), "\n")
    tripleCount := countTriples(lines)
    if tripleCount != 1 {
        t.Errorf("Expected exactly 1 triple, got %d", tripleCount)
    }
}
```

#### Integration with TDD Cycle:
1. **Red**: Write failing acceptance test that defines complete user experience
2. **Red**: Write failing unit tests for implementation details  
3. **Green**: Implement minimal code to make tests pass
4. **Green**: Verify acceptance test passes end-to-end
5. **Refactor**: Clean up while keeping all tests green

**Remember**: Acceptance tests guide the overall feature design, while unit tests guide the implementation details.

### Critical Process Guidelines for Code Changes

**Essential practices learned from major refactoring work to prevent errors and improve efficiency**

#### Comprehensive Testing Requirements

1. **Always run ALL test types before committing**: Even if acceptance tests pass, unit tests may still fail due to:
   - Type mismatches or import errors
   - Function signature changes not reflected in all test files
   - Copy-paste errors in test updates
   - Missing test updates after refactoring

2. **Test hierarchy and coverage**:
   - **Unit tests**: Test individual functions and methods in isolation
   - **Integration tests**: Test package interactions and data flow
   - **Acceptance tests**: Test the actual built binary end-to-end
   - **All must pass**: Don't rely on higher-level tests to catch lower-level issues

3. **Test after each logical change**: Run relevant tests immediately after:
   - Function signature changes
   - Type alias removal or addition
   - Package refactoring
   - Import statement updates

#### Build Artifact Management

1. **Clean up test artifacts immediately**: Test builds create artifacts that should never be committed:
   ```bash
   # Common test artifacts to clean up
   *.test           # Go test binaries (e.g., pairing.test, git.test)
   pairstair        # Main application binary
   coverage.out     # Coverage files
   ```

2. **Check for artifacts before committing**:
   - Use `git status` to verify no binaries are staged
   - Use `find . -name "*.test" -type f` to locate test artifacts
   - Remove with `rm` or `git rm` if accidentally staged
   - Update `.gitignore` if new artifact patterns emerge

3. **Automated artifact detection**: Add to pre-commit checks:
   ```bash
   # Check for common build artifacts
   if find . -name "*.test" -type f | grep -q .; then
       echo "Error: Test binaries found. Clean up before commit."
       exit 1
   fi
   ```

#### Function Signature Refactoring

**When changing function signatures that affect multiple files:**

1. **Identify all usage sites first**:
   - Use `grep_search` to find all function calls
   - Use `list_code_usages` for comprehensive reference finding
   - Include both production code and test files in the search

2. **Update systematically**:
   - Update production code first, then test files
   - Update one file at a time, testing after each
   - Pay special attention to test files with similar function calls
   - Avoid copy-paste errors by carefully checking each update

3. **Verify compiler catches all issues**:
   - Run `go build` to catch missing updates
   - Use compiler errors as a checklist for remaining work
   - Don't assume all issues are found by tests alone

#### Avoiding Redundant Manual Testing

**Critical principle: Trust your test suite and avoid duplicating what acceptance tests already verify**

1. **Leverage acceptance tests**: If comprehensive acceptance tests exist:
   - They provide end-to-end verification of CLI behavior
   - Manual binary testing becomes redundant and wastes time
   - Focus on ensuring all test types pass instead
   - **If acceptance tests pass, have confidence the functionality works as users will experience it**

2. **Don't duplicate test verification manually**:
   - ❌ **Avoid**: Running `go test ./...`, then building binary, then manually testing same functionality
   - ✅ **Do**: Run `go test ./...`, trust acceptance tests, move on
   - ❌ **Avoid**: Testing CLI flags manually when acceptance tests cover them
   - ✅ **Do**: Add missing acceptance test coverage if you feel manual testing is needed

3. **When manual testing is appropriate**:
   - New features not yet covered by acceptance tests (write the tests first!)
   - Complex edge cases that are hard to automate  
   - Performance or resource usage verification
   - User experience validation for new interfaces
   - Integration with external systems (browser opening, file operations)

4. **Acceptance test coverage should include**:
   - All CLI flags and combinations
   - Different input scenarios (various git repos)
   - Output format verification
   - Error handling and edge cases

**Remember**: Time spent on redundant manual verification is time not spent on valuable development work. Write comprehensive acceptance tests first, then trust them completely.

#### Responding to User Interventions

1. **When users redirect or halt work**:
   - Immediately stop the current approach
   - Ask clarifying questions about the new direction
   - Don't continue with planned steps that may no longer be relevant
   - Clean up any incomplete changes or artifacts

2. **Process improvement opportunities**:
   - Reflect on what led to the intervention
   - Update process guidelines to prevent similar issues
   - Document lessons learned for future reference
   - Adjust approach based on user feedback

3. **Communication during complex changes**:
   - Provide clear progress updates
   - Explain what each major step accomplishes
   - Ask for confirmation before proceeding with large changes
   - Highlight potential risks or breaking changes

#### Error Prevention and Recovery

1. **Establish baseline before major changes**:
   - Ensure all tests pass in current state
   - Create a git checkpoint with `git add . && git commit`
   - Document the known-good state for reference

2. **Incremental verification**:
   - Test after each logical change
   - Don't accumulate multiple changes before testing
   - Use `git add -p` to stage and test small chunks
   - Commit working increments regularly

3. **Recovery strategies**:
   - Use `git diff` to see what changed when tests fail
   - Use `git checkout` to revert problematic changes
   - Compare against last known working state
   - Fix one issue at a time rather than multiple at once

#### Code Quality During Refactoring

1. **Avoid shortcuts under pressure**:
   - Don't skip test updates to "save time"
   - Don't commit failing tests with plans to "fix later"
   - Don't leave build artifacts for "cleanup later"
   - Address issues immediately when discovered

2. **Maintain code consistency**:
   - Use domain types directly rather than aliases
   - Follow established patterns in the codebase
   - Update documentation alongside code changes
   - Ensure consistent naming and conventions

3. **Refactoring best practices**:
   - Remove deprecated code paths completely
   - Update all related documentation
   - Clean up unused imports and variables
   - Verify all TODOs and FIXMEs are addressed

**Remember**: These practices prevent the accumulation of technical debt and reduce the likelihood of subtle bugs. Time spent on proper process pays dividends in code quality and maintainability.

### Package Structure and Organization

**Follow consistent patterns for creating well-organized, maintainable packages**

When refactoring existing code or creating new functionality, use the established package structure pattern demonstrated by the `internal/update` package:

#### Package Creation Guidelines

1. **Use `internal/` for implementation packages**: All domain-specific packages should be placed under `internal/` to prevent external imports:
   ```
   internal/
   ├── update/     # Update notification logic
   ├── pairing/    # Pairing detection and analysis
   ├── team/       # Team file parsing and management  
   ├── git/        # Git log parsing and operations
   └── version/    # Version detection and reporting
   ```

2. **Single responsibility per package**: Each package should have one clear purpose:
   - `internal/update`: Check for new versions on GitHub
   - `internal/pairing`: Parse git logs for pairing information
   - `internal/team`: Handle team file operations
   - `internal/git`: Git repository operations
   - `internal/version`: Version detection from build info

3. **Minimal public APIs**: Export only what's necessary for external use:
   ```go
   // ✅ Good: Clean public API
   func CheckForUpdate(currentVersion string) string
   func IsNewerVersion(current, latest string) bool
   
   // ❌ Avoid: Exposing implementation details
   type release struct { ... } // Keep private
   func parseGitHubResponse(...) // Internal helper
   ```

#### Domain-Driven Type Organization

**Place types in packages that represent their primary domain**

When deciding where to define types, consider which domain they fundamentally belong to:

1. **Types belong in their primary domain package**: Define types where they conceptually fit best:
   ```go
   // ✅ Git-related types belong in the git package
   package git
   type Commit struct { ... }     // Committing is a git operation
   type Developer struct { ... }  // Parsed from git logs
   
   // ✅ Team-related types belong in the team package  
   package team
   type Team struct { ... }       // Team configuration and management
   
   // ❌ Avoid: Defining domain types in main package
   package main
   type Commit struct { ... }     // Should be in git package
   ```

2. **Use type aliases for backward compatibility during refactoring**:
   ```go
   // ✅ Transition approach: main package imports domain types
   package main
   import "github.com/project/internal/git"
   
   // Type alias for backward compatibility during refactoring
   type Commit = git.Commit
   type Developer = git.Developer
   ```

3. **Create conversion layers when needed**: During package refactoring, use conversion functions to bridge type systems:
   ```go
   // ✅ Temporary conversion during refactoring
   func convertCommit(gitCommit git.Commit) MainCommit {
       return MainCommit{
           Date: gitCommit.Date,
           Author: convertDeveloper(gitCommit.Author),
           // ...
       }
   }
   ```

4. **Migrate toward canonical domain types**: The goal is to use domain package types directly:
   ```go
   // ✅ End goal: main package uses domain types directly
   func ProcessCommits(commits []git.Commit) { ... }
   func AnalyzePairing(team team.Team, commits []git.Commit) { ... }
   ```

5. **Avoid cross-domain dependencies**: Keep packages focused on their domain:
   ```go
   // ✅ Good: git package doesn't depend on team concepts
   package git
   func GetCommits() []Commit { ... }
   
   // ✅ Good: higher-level coordination in main package
   package main 
   func main() {
       commits := git.GetCommits()
       team := team.LoadTeam()
       analysis := pairing.Analyze(team, commits)
       recommendation := pairing.RecommendPairs(pairs, strategy)
       output.Print(recommendation, format)
   }
   ```

This approach ensures that:
- **Types live where they conceptually belong**
- **Domain packages remain focused and cohesive**  
- **Refactoring can be done incrementally**
- **Dependencies flow in the right direction**

#### External Testing Pattern

**Always use external testing for packages to ensure proper encapsulation**

1. **Use `package name_test` pattern**: Test files should use external package naming:
   ```go
   // ✅ Correct external testing
   package update_test
   
   import (
       "testing"
       "github.com/gypsydave5/pairstair/internal/update"
   )
   
   // ❌ Avoid internal testing unless testing private functions
   package update
   ```

2. **Test only public interfaces**: External tests can only access exported functions and types:
   ```go
   // ✅ Test public API
   result := update.CheckForUpdate("v0.5.0")
   newer := update.IsNewerVersion("v0.5.0", "v0.6.0")
   
   // ❌ Cannot access private types (compiler prevents this)
   var r release // Error: undefined
   ```

3. **Use realistic test data**: Mock responses should use actual data formats:
   ```go
   // ✅ Use JSON strings for HTTP API responses
   mockResponse: `[
       {"tag_name": "v0.6.0", "draft": false},
       {"tag_name": "v0.5.0", "draft": false}
   ]`
   
   // ❌ Avoid exposing internal structs in tests
   mockResponse: []update.Release{...} // Error if Release is private
   ```

#### Dependency Injection for Testability

1. **Provide testable alternatives**: Export functions that accept dependencies for testing:
   ```go
   // ✅ Public function for production use
   func CheckForUpdate(currentVersion string) string {
       return CheckForUpdateWithURL(currentVersion, defaultURL)
   }
   
   // ✅ Testable function that accepts URL dependency
   func CheckForUpdateWithURL(currentVersion, url string) string {
       // Implementation can be tested with mock servers
   }
   ```

2. **Separate I/O from logic**: Keep business logic separate from external dependencies:
   ```go
   // ✅ Pure function for version comparison
   func IsNewerVersion(current, latest string) bool {
       // No I/O, easy to test
   }
   
   // ✅ I/O function that uses pure functions
   func CheckForUpdateWithURL(currentVersion, url string) string {
       // HTTP call, then use IsNewerVersion for logic
   }
   ```

#### Package Documentation and Examples

1. **Document package purpose clearly**: Each package should have clear documentation:
   ```go
   // Package update provides functionality for checking if newer versions
   // of the application are available on GitHub releases.
   //
   // The package performs silent HTTP requests to GitHub's API and compares
   // semantic versions to determine if an update notification should be shown.
   package update
   ```

2. **Provide usage examples**: Key functions should have examples:
   ```go
   // CheckForUpdate checks GitHub releases for a newer version.
   // Returns an empty string if no update is available or on any error.
   //
   // Example:
   //   message := CheckForUpdate("v0.5.0")
   //   if message != "" {
   //       fmt.Println(message)
   //   }
   func CheckForUpdate(currentVersion string) string
   ```

#### Refactoring Existing Code

When moving existing functionality into packages:

1. **Refactor incrementally**: Move one domain at a time, ensuring tests pass at each step
2. **Maintain backward compatibility**: Existing APIs should continue to work during transition
3. **Update imports systematically**: Use tools like `go mod tidy` and compiler errors to find all references
4. **Test after each move**: Run full test suite after each package extraction

#### Integration with Main Application

1. **Keep main package thin**: The main `pairstair.go` should primarily handle:
   - CLI flag parsing and validation
   - Coordinating calls to domain packages  
   - Output formatting and presentation
   - Error handling and user feedback

2. **Use packages for domain logic**: Move business logic into appropriate packages:
   ```go
   // ✅ Main handles CLI, packages handle logic
   func main() {
       // Parse flags...
       updateMsg := update.CheckForUpdate(getVersion())
       if updateMsg != "" {
           fmt.Fprintf(os.Stderr, "%s\n", updateMsg)
       }
       
       pairs := pairing.AnalyzeRepository(repoPath, timeWindow)
       recommendation := pairing.RecommendPairs(pairs, strategy)
       output.Print(recommendation, format)
   }
   ```

This package structure and testing approach ensures:
- **Clear separation of concerns** between different domains
- **Testable code** through external testing and dependency injection  
- **Maintainable APIs** with minimal surface area
- **Reliable refactoring** through comprehensive test coverage

### Git Commit Message Conventions

Use consistent prefixes for commit messages to indicate the type of change:

- **`-s-`** for structural changes that introduce no change in behavior (refactoring, documentation, formatting, test updates, etc.)
- **`-b-`** for behavioral changes that modify functionality (new features, bug fixes, API changes, etc.)

Examples:
```
-s- refactor matrix building logic for clarity
-s- update README with new installation instructions
-s- add tests for edge case handling
-b- add --strategy flag for least-recent pairing
-b- fix incorrect pair counting in team mode
-b- change default window from 1w to 2w
```

This convention helps reviewers and maintainers quickly understand the impact and scope of changes.

### Documentation Maintenance

Keep all documentation synchronized when adding features or making changes:

1. **When adding new CLI flags or features**:
   - Update `README.md` with new flag documentation and examples
   - Update `docs/pairstair.1.md` (man page source) with new options and examples
   - Regenerate the man page by running `cd docs && ./gen-man.sh`
   - Commit both the markdown source and generated man page

2. **Documentation files to update**:
   - `README.md`: Primary user documentation with examples and usage
   - `docs/pairstair.1.md`: Man page source (markdown format)
   - `docs/pairstair.1`: Generated man page (regenerate after editing source)
   - `.team.example`: Example team file showing current format

3. **Documentation standards**:
   - Include comprehensive examples for new features
   - Show flag combinations when applicable
   - Explain behavior clearly, especially edge cases
   - Update both basic and advanced usage examples

4. **Testing documentation**:
   - Verify examples work as documented
   - Test man page generation with `cd docs && ./gen-man.sh`
   - Ensure help output (`--help`) matches documented behavior

5. **Release process**:
   - Documentation updates warrant patch version bumps
   - Use `-s-` prefix for documentation-only commits
   - The CI/CD pipeline automatically generates man pages for releases

Remember: Users rely on both README and man page documentation. Keep them comprehensive and in sync.

### Documentation Accuracy and Verification

**Always verify functionality exists before documenting it**

When writing or updating documentation, follow these critical practices to ensure accuracy:

1. **Verify before documenting**: Before documenting any feature, command-line flag, or functionality:
   - Use `grep_search` or `semantic_search` to confirm the feature exists in the codebase
   - Check actual implementation files (e.g., `pairstair.go`, `recommend.go`, `print.go`)
   - Test the functionality locally if possible
   - Never assume features exist based on "what should be there"

2. **Cross-reference with implementation**: When documenting CLI options or strategies:
   - Check the actual switch/case statements in the code
   - Verify flag parsing logic in `pairstair.go`
   - Confirm strategy implementations in `recommend.go`
   - Match help text output with documented behavior

3. **Test examples and commands**: Before including code examples:
   - Verify commands actually work as documented
   - Test all flag combinations shown in examples
   - Ensure output samples match actual program output
   - Remove or correct any examples that don't work

4. **Distinguish planned vs. current features**: Clearly separate:
   - **Current features**: What exists in the codebase right now
   - **Planned features**: Future enhancements (document in `FEATURES.md` or "Planned Features" sections)
   - Use appropriate language: "PairStair supports..." vs. "PairStair will support..."

5. **When documenting new functionality**: After implementing a feature:
   - Document the actual implementation, not the ideal design
   - Include realistic examples with actual output
   - Update all relevant documentation files (README, man page, Jekyll site)
   - Test documentation examples against the implemented feature

6. **Regular documentation audits**: Periodically verify that:
   - All documented features still exist and work as described
   - Examples produce the documented output
   - No obsolete or incorrect information remains
   - Version numbers and capabilities are current

**Remember**: Documentation credibility depends on accuracy. Users trust that documented features actually work. Always verify before documenting.

### Development Scripts

PairStair includes two development helper scripts to streamline common workflows:

#### dev.sh - Development Helper

```bash
./dev.sh docs     # Regenerate man page from markdown source  
./dev.sh version  # Show comprehensive version information
```

**Use for**: Tasks that add value beyond standard Go commands. For basic operations, use `go test ./...`, `go build`, etc. directly.

#### release.sh - Release Automation  

```bash
./release.sh patch "Bug fixes"               # Auto-increment patch (0.7.2 -> 0.7.3)
./release.sh minor "New features"            # Auto-increment minor (0.7.2 -> 0.8.0)
./release.sh major "Breaking changes"        # Auto-increment major (0.7.2 -> 1.0.0)
./release.sh -v v2.0.0 "Custom version"      # Specify exact version when needed
```

**Features**:
- Automatic semantic version calculation from latest git tag
- Manual version override with -v flag for special cases
- Requires completely clean working directory (no uncommitted changes)
- Runs all tests and cleans up artifacts
- Creates annotated git tag and pushes to trigger CI/CD
- Non-interactive operation suitable for automation

**Safety**: Script aborts on any validation failure with helpful error messages.

**Release notes conventions**: Focus release notes on describing changes, not repeating version numbers. Git tags inherently contain version information, making version repetition redundant.

**Documentation**: Full details in `CONTRIBUTING.md` for both scripts.

### Git Repository Hygiene

Maintain a clean repository by following these practices:

1. **Never commit binary files**: Go binaries, executables, and compiled artifacts should never be committed to version control
   - The main binary (e.g., `pairstair`) should be listed in `.gitignore`
   - Test binaries (e.g., `*.test` files) should be cleaned up immediately after testing
   - Use `go build` to create binaries locally as needed
   - If you notice a binary has been accidentally committed, immediately remove it with `git rm <binary>` and update `.gitignore`

2. **Clean up test artifacts promptly**: Go tests create binary artifacts that must be removed:
   - Use `find . -name "*.test" -type f` to locate test binaries
   - Remove with `rm *.test` or similar commands
   - Check for artifacts before every commit with `git status`
   - Common artifacts: `pairing.test`, `git.test`, `output.test`, main binary name

3. **Review staged changes before committing**: Always check `git status` and `git diff --staged` before committing
   - Look for accidentally staged binary files, temporary files, or IDE artifacts
   - Ensure only intended source code changes are included
   - Pay special attention after running tests that create binaries

4. **Keep .gitignore comprehensive**: The `.gitignore` file should include:
   - IDE files and directories (`.idea/`, `.vscode/`, etc.)
   - Go binaries (the main executable name)
   - Test binaries (`*.test`)
   - Temporary files and build artifacts
   - Coverage files (`coverage.out`, `*.cover`)
   - OS-specific files (`.DS_Store`, `Thumbs.db`, etc.)

5. **Immediate correction**: If you notice binary files or other inappropriate content has been committed:
   - Remove it immediately with `git rm <file>`
   - Add appropriate entries to `.gitignore`
   - Commit the cleanup with a `-s-` prefix message
   - This prevents repository bloat and keeps the history clean

### Build and Release Documentation

**Maintain comprehensive documentation of the build and release process**

The PairStair project has a sophisticated build and release system that should be well-documented:

1. **CONTRIBUTING.md**: Primary developer documentation covering:
   - Development setup and local building
   - Version detection system and how it works
   - CI/CD pipeline details and version injection
   - Manual release process and checklist
   - Homebrew integration and tap repository
   - Testing procedures and TDD guidelines

2. **README.md Installation section**: User-focused documentation covering:
   - Multiple installation methods (Homebrew, Go install, manual download)
   - Platform support and requirements
   - Links to detailed build documentation

3. **External repository documentation**: The separate `gypsydave5/homebrew-pairstair` repository:
   - Maintains the Homebrew formula
   - Automatically updated by CI/CD pipeline
   - Requires manual intervention only if automation fails

4. **Key build concepts to document**:
   - **Version injection**: How `-ldflags` injects version into pre-built binaries
   - **Smart version detection**: Priority order from git tags to fallback constants
   - **Multi-platform builds**: Supported architectures and cross-compilation
   - **Automated releases**: Tag-triggered vs. manual dispatch workflows
   - **Homebrew integration**: Repository dispatch mechanism and formula updates

5. **When updating build documentation**:
   - Test all documented commands and procedures
   - Verify version numbers and URLs are current
   - Update both user-facing (README) and developer-facing (CONTRIBUTING) docs
   - Include examples of actual commands and expected outputs
   - Document any external dependencies or requirements

**Remember**: The build system is complex (version detection, multi-platform compilation, external Homebrew repository). Users and contributors need clear, accurate documentation to understand and work with the system effectively.

### Future Feature Ideas

For ideas about potential future enhancements to PairStair, see [FEATURES.md](../FEATURES.md). This is where new feature concepts should be documented and prioritized.
