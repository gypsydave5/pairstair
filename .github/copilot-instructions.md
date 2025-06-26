# Instructions to Copilot for writing Go (Golang) code.

We favour following idiomatic Go practices, using the standard library where possible, and writing clear, maintainable code. Here are some guidelines:

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

### Error Investigation

When investigating test failures or bugs:
- Use git to check what changed recently
- Add temporary debug logging to understand data flow
- Verify assumptions about data with explicit checks
- Remove debug code once the issue is resolved

These practices help maintain code quality and prevent regressions during development.

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

### Git Repository Hygiene

Maintain a clean repository by following these practices:

1. **Never commit binary files**: Go binaries, executables, and compiled artifacts should never be committed to version control
   - The main binary (e.g., `pairstair`) should be listed in `.gitignore`
   - Use `go build` to create binaries locally as needed
   - If you notice a binary has been accidentally committed, immediately remove it with `git rm <binary>` and update `.gitignore`

2. **Review staged changes before committing**: Always check `git status` and `git diff --staged` before committing
   - Look for accidentally staged binary files, temporary files, or IDE artifacts
   - Ensure only intended source code changes are included

3. **Keep .gitignore comprehensive**: The `.gitignore` file should include:
   - IDE files and directories (`.idea/`, `.vscode/`, etc.)
   - Go binaries (the main executable name)
   - Temporary files and build artifacts
   - OS-specific files (`.DS_Store`, `Thumbs.db`, etc.)

4. **Immediate correction**: If you notice binary files or other inappropriate content has been committed:
   - Remove it immediately with `git rm <file>`
   - Add appropriate entries to `.gitignore`
   - Commit the cleanup with a `-s-` prefix message
   - This prevents repository bloat and keeps the history clean

### Future Feature Ideas

For ideas about potential future enhancements to PairStair, see [FEATURES.md](../FEATURES.md). This is where new feature concepts should be documented and prioritized.
