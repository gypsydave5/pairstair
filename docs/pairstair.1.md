% PAIRSTAIR(1) | User Commands

# NAME

pairstair - visualize and optimize software developer pairing from git history

# SYNOPSIS

**pairstair** [**-window** _window_] [**-output** _format_]

# DESCRIPTION

**pairstair** analyzes git commit history to show how often developers have paired (committed together) over a configurable time window, and recommends new pairs to maximize collaboration.

The tool scans commits in the specified time window, finds the author and any co-authors (from "Co-authored-by" trailers), and builds a matrix showing how many days each pair has worked together.

# OPTIONS

**-window** _window_  
:   Set the time window to analyze. Examples: `1d` (1 day), `2w` (2 weeks), `3m` (3 months), `1y` (1 year). Default: `1w`.

**-output** _format_  
:   Output format. Options: `cli` (default, prints to terminal), `html` (opens results in browser).

# TEAM FILE

If a `.team` file is present in the working directory, only developers listed are included in the analysis. Each line should contain a developer's display name followed by their email address(es) in angle brackets. For developers who use multiple email addresses, separate them with commas and enclose each in angle brackets:

    Alice Example <alice@example.com>,<alice@gmail.com>
    Bob Dev <bob@example.com>
    Carol Tester <carol@example.com>,<carol@personal.com>

When multiple email addresses are specified for one developer, commits from any of those addresses will be attributed to the same person. This helps create an accurate pairing matrix even when developers use different email addresses.

If no `.team` file exists, all authors from the git history are included.

# EXAMPLES

Analyze the last 4 weeks and show results in the terminal:

    pairstair -window 4w

Show results as HTML in your browser:

    pairstair -output html

# AUTHORS

Written by gypsydave5.

# SEE ALSO

**git-log**(1)
