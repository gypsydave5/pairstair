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

### The `.team` File

If you want to restrict the analysis to a specific team, create a `.team` file in your repository root. Each line should be a git author string (e.g., `Alice Example <alice@example.com>`).

Example `.team`:

```
Alice Example <alice@example.com>
Bob Dev <bob@example.com>
Carol Tester <carol@example.com>
```

If `.team` is not present, PairStair will use all authors found in the git history.

## How It Works

- Scans commits in the specified window.
- For each commit, finds the author and any co-authors.
- Groups developers by email address (so aliases are merged).
- Builds a matrix showing how many days each pair has worked together.
- Prints a legend mapping short initials to developer names/emails.
- Prints pairing recommendations, suggesting pairs who have worked together the least.

## Example Output

```
Legend:
  AE     = Alice Example        alice@example.com
  BD     = Bob Dev              bob@example.com
  CT     = Carol Tester         carol@example.com

        AE      BD      CT
AE      -       2       1
BD      2       -       0
CT      1       0       -

Pairing Recommendations (least-paired first):
  BD     <-> CT     : 0 times
```

## License

MIT
