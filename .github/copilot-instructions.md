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
