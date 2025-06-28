# PairStair

PairStair is a CLI tool for visualizing and optimizing software developer pairing. It analyzes your git history to show how often developers have paired (committed together) over a configurable time window, and recommends new pairs to maximize collaboration.

## Features

- Prints a "pair stair" matrix showing how often each pair of developers has worked together.
- Reads git commit authors and "Co-authored-by" trailers to detect pairs.
- Optionally restricts analysis to a team defined in a `.team` file.
- Supports configurable time windows (e.g., last week, last month).
- Provides pairing recommendations to encourage new collaborations.

## Installation

PairStair can be installed through several methods:

### Homebrew (macOS)

The easiest way to install on macOS:

```bash
brew tap gypsydave5/pairstair
brew install pairstair
```

This installs a pre-built binary with proper version information.

### Go Install

Install directly from source (requires Go 1.21+):

```bash
go install github.com/gypsydave5/pairstair@latest
```

This will install the `pairstair` binary in your `$GOPATH/bin` or `$HOME/go/bin`.

### Manual Download

Download pre-built binaries from the [GitHub Releases page](https://github.com/gypsydave5/pairstair/releases). Available for:

- **Linux**: `amd64`, `arm64`
- **macOS**: `amd64` (Intel), `arm64` (Apple Silicon)  
- **Windows**: `amd64`

### Build from Source

For development or custom builds:

```bash
git clone https://github.com/gypsydave5/pairstair.git
cd pairstair
go build -o pairstair .
```

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for detailed build instructions and development setup.

## Usage

Run `pairstair` in the root of a git repository:

```sh
pairstair
```

### Options

#### `-window <window>`: Set the time window to analyze.

Examples:
  - `1d` (1 day)
  - `2w` (2 weeks)
  - `3m` (3 months)
  - `4y` (4 years)

Example:

```sh
pairstair -window 4w
```

#### `-output <type>`: Set the output format.

Options:
  - `cli` (default): Prints the pairing matrix on the command line.
  - `html`: Outputs the pairing data in HTML format to stdout (can be redirected to files).

#### `-open`: Open HTML output in browser.

When combined with `-output html`, opens the HTML results directly in your default web browser instead of streaming to stdout.

#### `-strategy <strategy>`: Set the pairing recommendation strategy.

Options:
  - `least-paired` (default): Recommends pairs who have worked together the fewest times, using optimal matching to minimize total pair count.
  - `least-recent`: Recommends pairs who haven't worked together for the longest time, prioritizing pairs who have never collaborated.

Example:

```sh
pairstair -strategy least-recent
```

Examples with HTML output:

```sh
# Stream HTML to a file
pairstair -output html > report.html

# Open HTML directly in browser
pairstair -output html -open
```

#### `-team <team>`: Specify a sub-team to analyze.

When your `.team` file contains sub-teams (see below), you can analyze just a specific sub-team instead of the entire team.

Example:

```sh
pairstair -team frontend
```

### The `.team` File

If you want to restrict the analysis to a specific team, create a `.team` file in your repository root. Each line should contain a developer's display name followed by their email address(es) in angle brackets.

**Multiple email addresses**: You can specify multiple email addresses for the same developer by separating them with commas (`,`) and enclosing each in angle brackets (`<>`).

#### Basic Team File

Example `.team`:

```
Alice Example <alice@example.com>,<alice@gmail.com>
Bob Dev <bob@example.com>
Carol Tester <carol@example.com>,<carol@personal.com>,<carol@old-company.com>
```

#### Sub-teams

You can organize your team into sub-teams using section headers in square brackets. When no `--team` flag is specified, only team members not in any sub-team section are analyzed.

Example `.team` with sub-teams:

```
Alice Lead <alice@example.com>
Bob Manager <bob@example.com>

[frontend]
Carol Frontend <carol@example.com>
Dave UI <dave@example.com>

[backend]
Eve Backend <eve@example.com>
Frank API <frank@example.com>

[devops]
Grace Ops <grace@example.com>
```

In this example:
- `pairstair` (no flags) analyzes Alice and Bob
- `pairstair --team=frontend` analyzes Carol and Dave
- `pairstair --team=backend` analyzes Eve and Frank
- `pairstair --team=devops` analyzes Grace only

**Multiple sub-teams**: If a developer needs to be in multiple sub-teams, you can duplicate their entry in each relevant section:

```
Alice Main <alice@example.com>
Bob BothMainAndSub <bob@example.com>

[frontend]
Bob BothMainAndSub <bob@example.com>
Carol SubTeamOnly <carol@example.com>

[backend]
Bob BothMainAndSub <bob@example.com>
Dave SubTeamOnly <dave@example.com>
```

In this case:
- `pairstair` (no flags) analyzes Alice and Bob only
- `pairstair --team=frontend` analyzes Bob and Carol
- `pairstair --team=backend` analyzes Bob and Dave
- Bob appears in all analyses, but Carol and Dave only appear in their respective sub-teams

If a developer has commits from different email addresses, they will be treated as the same person when calculating the pairing matrix.

If `.team` is not present, PairStair will use all authors found in the git history.

## How It Works

- Scans commits in the specified window.
- For each commit, finds the author and any co-authors.
- Groups developers by email address (so aliases are merged, and multiple email addresses in the `.team` file are combined).
- Builds a matrix showing how many days each pair has worked together.
- Prints a legend mapping short initials to developer names/emails.
- Prints pairing recommendations, suggesting pairs who have worked together the least (only if total number of developers is 10 or less).

## Example Output

### Default Strategy (least-paired)

```
Legend:
  AE     = Alice Example        alice@example.com
  BD     = Bob Dev              bob@example.com
  CT     = Carol Tester         carol@example.com

        AE      BD      CT
AE      -       2       1
BD      2       -       0
CT      1       0       -

Pairing Recommendations (least-paired overall, optimal matching):
  BD     <-> CT     : 0 times
```

### Least-Recent Strategy

```sh
pairstair -strategy least-recent
```

```
Legend:
  AE     = Alice Example        alice@example.com
  BD     = Bob Dev              bob@example.com
  CT     = Carol Tester         carol@example.com

        AE      BD      CT
AE      -       2       1
BD      2       -       0
CT      1       0       -

Pairing Recommendations (least recent collaborations first):
  BD     <-> CT     : never paired
  AE     <-> CT     : last paired 15 days ago
```

## Development and Contributing

PairStair welcomes contributions! Whether you're fixing bugs, adding features, or improving documentation, your help is appreciated.

### Quick Start

```bash
# Clone and build
git clone https://github.com/gypsydave5/pairstair.git
cd pairstair
go build -o pairstair .

# Run tests
go test -v ./...
```

### Key Resources

- **[`CONTRIBUTING.md`](CONTRIBUTING.md)** - Complete development guide including:
  - Build process and CI/CD pipeline
  - Release workflow and version management
  - Homebrew tap integration
  - Testing and TDD guidelines
  - Code style and conventions

- **[`.github/copilot-instructions.md`](.github/copilot-instructions.md)** - AI assistant guidelines and project conventions

### Release Process

PairStair uses an automated CI/CD pipeline that:

1. **Builds** binaries for multiple platforms with proper version injection
2. **Creates** GitHub releases with all artifacts  
3. **Updates** the Homebrew formula automatically via [`gypsydave5/homebrew-pairstair`](https://github.com/gypsydave5/homebrew-pairstair)

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for detailed release procedures.

## Contributing

Contributions are welcome! For development setup, testing guidelines, and the release process, see [CONTRIBUTING.md](CONTRIBUTING.md).

### Quick Development Reference

```bash
# Run tests
go test ./...

# Build locally
go build

# Development helpers
./dev.sh docs     # Regenerate man page
./dev.sh version  # Show version info

# Create a release (requires clean working directory)
./release.sh v0.x.y "Release notes"
```

## Documentation

For comprehensive documentation, examples, and guides, visit the **[PairStair Documentation Site](https://gypsydave5.github.io/pairstair/)**.

The documentation includes:

- **[Installation Guide](https://gypsydave5.github.io/pairstair/installation/)** - Multiple installation methods and setup
- **[User Guide](https://gypsydave5.github.io/pairstair/guide/)** - Complete guide to using all features
- **[Examples](https://gypsydave5.github.io/pairstair/examples/)** - Real-world scenarios and workflows
- **[Features](https://gypsydave5.github.io/pairstair/features/)** - Current and planned functionality

## License

MIT
