% PAIRSTAIR(1) | User Commands

# NAME

pairstair - visualize and optimize software developer pairing from git history

# SYNOPSIS

**pairstair** [**-window** _window_] [**-output** _format_] [**-open**] [**-strategy** _strategy_] [**-team** _team_]

# DESCRIPTION

**pairstair** analyzes git commit history to show how often developers have paired (committed together) over a configurable time window, and recommends new pairs to maximize collaboration.

The tool scans commits in the specified time window, finds the author and any co-authors (from "Co-authored-by" trailers), and builds a matrix showing how many days each pair has worked together.

# OPTIONS

**-window** _window_  
:   Set the time window to analyze. Examples: `1d` (1 day), `2w` (2 weeks), `3m` (3 months), `1y` (1 year). Default: `1w`.

**-output** _format_  
:   Output format. Options: `cli` (default, prints to terminal), `html` (streams HTML to stdout).

**-open**  
:   Open HTML output in browser. Only applies when `-output html` is specified. Without this flag, HTML is streamed to stdout, allowing redirection to files or piping to other tools.

**-strategy** _strategy_  
:   Pairing recommendation strategy. Options: `least-paired` (default, recommends pairs who have worked together the fewest times), `least-recent` (recommends pairs who haven't worked together for the longest time).

**-team** _team_  
:   Specify a sub-team to analyze. When your `.team` file contains sub-teams (see below), analyze only that specific sub-team instead of the entire team.

# TEAM FILE

If a `.team` file is present in the working directory, only developers listed are included in the analysis. Each line should contain a developer's display name followed by their email address(es) in angle brackets. For developers who use multiple email addresses, separate them with commas and enclose each in angle brackets:

    Alice Example <alice@example.com>,<alice@gmail.com>
    Bob Dev <bob@example.com>
    Carol Tester <carol@example.com>,<carol@personal.com>

When multiple email addresses are specified for one developer, commits from any of those addresses will be attributed to the same person. This helps create an accurate pairing matrix even when developers use different email addresses.

## Sub-teams

You can organize your team into sub-teams using section headers in square brackets. When no `--team` flag is specified, only team members not in any sub-team section are analyzed:

    Alice Lead <alice@example.com>
    Bob Manager <bob@example.com>
    
    [frontend]
    Carol Frontend <carol@example.com>
    Dave UI <dave@example.com>
    
    [backend]
    Eve Backend <eve@example.com>
    Frank API <frank@example.com>

In this example, `pairstair` analyzes Alice and Bob, `pairstair --team=frontend` analyzes Carol and Dave, and `pairstair --team=backend` analyzes Eve and Frank.

If a developer needs to be in multiple sub-teams, duplicate their entry in each relevant section.

If no `.team` file exists, all authors from the git history are included.

# EXAMPLES

Analyze the last 4 weeks and show results in the terminal:

    pairstair -window 4w

Stream HTML results to stdout (can be redirected to a file):

    pairstair -output html > report.html

Open HTML results in your browser:

    pairstair -output html -open

Analyze only the frontend sub-team:

    pairstair -team frontend

Use least-recent strategy for recommendations:

    pairstair -strategy least-recent

Combine options to analyze backend team for the last month using least-recent strategy:

    pairstair -window 1m -team backend -strategy least-recent

# AUTHORS

Written by gypsydave5.

# SEE ALSO

**git-log**(1)
