# PairStair

PairStair is a CLI tool for visualizing and optimizing software developer pairing. It analyzes your git history to show how often developers have paired (committed together) over a configurable time window, and recommends new pairs to maximize collaboration.

## Features

- Prints a "pair stair" matrix showing how often each pair of developers has worked together.
- Reads git commit authors and "Co-authored-by" trailers to detect pairs.
- Optionally restricts analysis to a team defined in a `.team` file.
- Supports configurable time windows (e.g., last week, last month).
- Provides pairing recommendations to encourage new collaborations.

## Installation

You need Go installed (version 1.18+ recommended).

### Homebrew

You can install PairStair using Homebrew with:

```
brew install gypsydave5/pairstair/pairstair
```

### Go install

Alternatively, you can install with Go:

```sh
go install github.com/gypsydave5/pairstair@latest
```

This will install the `pairstair` binary in your `$GOPATH/bin` or `$HOME/go/bin`.

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
  - `html`: Outputs the pairing data in HTML and opens it in a web browser.

#### `-strategy <strategy>`: Set the pairing recommendation strategy.

Options:
  - `least-paired` (default): Recommends pairs who have worked together the fewest times, using optimal matching to minimize total pair count.
  - `least-recent`: Recommends pairs who haven't worked together for the longest time, prioritizing pairs who have never collaborated.

Example:

```sh
pairstair -strategy least-recent
```

### The `.team` File

If you want to restrict the analysis to a specific team, create a `.team` file in your repository root. Each line should contain a developer's display name followed by their email address(es) in angle brackets.

**Multiple email addresses**: You can specify multiple email addresses for the same developer by separating them with commas (`,`) and enclosing each in angle brackets (`<>`).

Example `.team`:

```
Alice Example <alice@example.com>,<alice@gmail.com>
Bob Dev <bob@example.com>
Carol Tester <carol@example.com>,<carol@personal.com>,<carol@old-company.com>
```

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

## License

MIT
