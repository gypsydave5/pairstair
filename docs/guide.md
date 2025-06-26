---
layout: page
title: User Guide
permalink: /guide/
---

# User Guide

This guide covers all aspects of using PairStair to analyze and optimize developer pairing in your team.

## Basic Concepts

### Pairing Detection

PairStair detects pairs by analyzing git commit messages for **co-authored-by trailers**:

```
Some commit message

Co-authored-by: Alice Smith <alice@example.com>
Co-authored-by: Bob Jones <bob@example.com>
```

When git commits include these trailers, PairStair counts it as evidence that those developers worked together.

### Time Windows

PairStair analyzes commits within a configurable time window. The default is 1 week, but you can specify:

- **Days**: `1d`, `7d`, `30d`
- **Weeks**: `1w`, `2w`, `4w`
- **Months**: `1m`, `3m`, `6m`
- **Years**: `1y`

```bash
# Analyze the last 2 weeks
pairstair -window 2w

# Analyze the last 3 months
pairstair -window 3m
```

## Team Files

### Basic Team File

Create a `.team` file in your repository root to define your team:

```
Alice Smith <alice@example.com>
Bob Jones <bob@example.com>
Carol White <carol@example.com>
Dave Brown <dave@example.com>
```

This helps PairStair:
- Focus analysis on known team members
- Handle name/email variations
- Generate better recommendations

### Sub-Teams

Organize larger teams into sub-teams using sections:

```
[frontend]
Alice Smith <alice@example.com>
Bob Jones <bob@example.com>

[backend]
Carol White <carol@example.com>
Dave Brown <dave@example.com>

[mobile]
Eve Green <eve@example.com>
Frank Black <frank@example.com>
```

Analyze specific sub-teams with the `--team` flag:

```bash
# Analyze only the frontend team
pairstair --team frontend

# Analyze only the backend team
pairstair --team backend
```

## Output Formats

### Console Output (Default)

The default output shows a pair matrix and recommendations:

```
Pair Matrix (commits together):
        Alice  Bob  Carol  Dave
Alice     -    5     2     0
Bob       5    -     1     3
Carol     2    1     -     4
Dave      0    3     4     -

Recommendations:
1. Alice + Dave (never paired)
2. Bob + Carol (1 commit together)
```

### HTML Output

Generate a rich web interface with detailed analysis:

```bash
pairstair -output html
```

This creates an HTML file with:
- Interactive pair matrix
- Detailed statistics
- Visual charts
- Sortable tables

## Recommendation Strategies

PairStair offers different strategies for recommending pairs:

### Least Recent (Default)

Recommends pairs who haven't worked together recently:

```bash
pairstair -strategy least-recent
```

### Never Paired

Recommends developers who have never worked together:

```bash
pairstair -strategy never-paired
```

### Round Robin

Ensures all possible pairs get equal opportunities:

```bash
pairstair -strategy round-robin
```

## Advanced Usage

### Multiple Repositories

Analyze pairing across multiple repositories by running PairStair in each repo and combining results manually, or use shell scripting:

```bash
#!/bin/bash
for repo in repo1 repo2 repo3; do
    echo "=== $repo ==="
    cd $repo && pairstair
    cd ..
done
```

### Integration with CI/CD

Add pairing analysis to your CI pipeline:

```yaml
# .github/workflows/pairing.yml
name: Pairing Analysis
on:
  schedule:
    - cron: '0 9 * * MON'  # Every Monday at 9 AM

jobs:
  pairing:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0  # Full history needed
      - name: Install PairStair
        run: go install github.com/gypsydave5/pairstair@latest
      - name: Generate Report
        run: pairstair -output html > pairing-report.html
      - name: Upload Report
        uses: actions/upload-artifact@v3
        with:
          name: pairing-report
          path: pairing-report.html
```

### Custom Time Ranges

For more control over the analysis period, you can use git directly:

```bash
# Analyze commits from a specific date range
git log --since="2024-01-01" --until="2024-01-31" --pretty=fuller | pairstair
```

## Troubleshooting

### No Pairs Detected

If PairStair shows no pairs:

1. **Check co-authored-by trailers**: Ensure commits include proper co-author information
2. **Verify team file**: Make sure email addresses match git commits exactly
3. **Check time window**: Try a longer window (e.g., `-window 1m`)
4. **Repository history**: Ensure you have full git history (`git fetch --unshallow`)

### Incorrect Developer Names

If developers appear with multiple names/emails:

1. **Standardize in .team file**: List all variations for each developer
2. **Use git mailmap**: Create a `.mailmap` file to normalize names
3. **Configure git**: Ensure consistent `user.name` and `user.email` settings

### Performance Issues

For large repositories:

1. **Limit time window**: Use shorter windows (e.g., `-window 2w`)
2. **Use sub-teams**: Analyze smaller groups with `--team`
3. **Shallow clone**: Consider using `git clone --depth 100` for recent history

## Tips and Best Practices

### Establishing Pairing Culture

1. **Start small**: Begin with voluntary pairing sessions
2. **Rotate regularly**: Use PairStair recommendations weekly
3. **Track progress**: Regular analysis helps show improvement
4. **Celebrate success**: Acknowledge when pairing goals are met

### Effective Pairing Sessions

1. **Set clear goals**: Know what you want to accomplish
2. **Switch roles**: Driver and navigator should alternate
3. **Take breaks**: Pairing can be mentally intensive
4. **Document outcomes**: Update commit messages with co-authors

### Team Management

1. **Regular reviews**: Weekly pairing analysis and planning
2. **Skill sharing**: Pair experienced with junior developers
3. **Cross-team pairing**: Occasional pairing across sub-teams
4. **Knowledge transfer**: Use pairing for onboarding and mentoring
