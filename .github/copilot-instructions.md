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

The public 
