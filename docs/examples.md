---
layout: page
title: Examples
permalink: /examples/
---

# Examples

This page provides practical examples of using PairStair in different scenarios and team configurations.

## Quick Start Examples

### Basic Analysis

Analyze pairing for the last week (default):

```bash
cd your-project
pairstair
```

**Output:**
```
Pair Matrix (commits together):
        Alice  Bob  Carol
Alice     -    3     1
Bob       3    -     2
Carol     1    2     -

Recommendations:
1. Alice + Carol (1 commit together)
2. Bob + Carol (2 commits together)
```

### Different Time Windows

```bash
# Last 2 weeks
pairstair -window 2w

# Last month
pairstair -window 1m

# Last 3 months
pairstair -window 3m
```

### HTML Output

Generate a rich web report:

```bash
pairstair -output html > pairing-report.html
open pairing-report.html  # macOS
```

## Team File Examples

### Simple Team

**`.team` file:**
```
Alice Smith <alice@company.com>
Bob Johnson <bob@company.com>
Carol Williams <carol@company.com>
Dave Brown <dave@company.com>
```

**Usage:**
```bash
pairstair -window 2w
```

### Multiple Email Addresses

Handle developers with multiple email addresses:

**`.team` file:**
```
Alice Smith <alice@company.com>
Alice Smith <alice.smith@company.com>
Alice Smith <a.smith@company.com>

Bob Johnson <bob@company.com>
Bob Johnson <bjohnson@company.com>
```

### Sub-Teams Organization

**`.team` file:**
```
[frontend]
Alice Smith <alice@company.com>
Bob Johnson <bob@company.com>
Carol Williams <carol@company.com>

[backend]
Dave Brown <dave@company.com>
Eve Davis <eve@company.com>
Frank Wilson <frank@company.com>

[mobile]
Grace Miller <grace@company.com>
Henry Garcia <henry@company.com>

[devops]
Ivy Rodriguez <ivy@company.com>
Jack Martinez <jack@company.com>
```

**Usage:**
```bash
# Analyze specific sub-teams
pairstair --team frontend
pairstair --team backend -window 3w
pairstair --team mobile -output html
```

## Real-World Scenarios

### Scenario 1: New Team Assessment

**Situation:** You've just joined a team and want to understand current pairing patterns.

**Commands:**
```bash
# Get overall picture for the last month
pairstair -window 1m

# Generate detailed HTML report
pairstair -window 1m -output html > team-analysis.html
```

**Sample Output:**
```
Pair Matrix (commits together):
        Alice  Bob  Carol  Dave  Eve
Alice     -    12    3     0    2
Bob      12     -    8     1    4
Carol     3     8    -     5    0
Dave      0     1    5     -    9
Eve       2     4    0     9    -

Recommendations:
1. Alice + Dave (never paired)
2. Carol + Eve (never paired)
3. Alice + Carol (3 commits together)
```

**Insights:**
- Alice and Bob pair frequently (12 commits)
- Dave and Eve have a strong pairing relationship (9 commits)
- Alice has never paired with Dave
- Carol has never paired with Eve

### Scenario 2: Sub-Team Rotation

**Situation:** You want to optimize pairing within your frontend team.

**Team File:**
```
[frontend]
Sarah Connor <sarah@company.com>
John Doe <john@company.com>
Jane Smith <jane@company.com>
Mike Wilson <mike@company.com>
```

**Commands:**
```bash
# Current frontend pairing (last 2 weeks)
pairstair --team frontend -window 2w

# Check if anyone hasn't paired recently
pairstair --team frontend -window 1w -strategy least-recent
```

**Sample Output:**
```
Frontend Team - Pair Matrix (commits together):
        Sarah  John  Jane  Mike
Sarah     -     5     2     0
John      5     -     1     3
Jane      2     1     -     0
Mike      0     3     0     -

Recommendations:
1. Sarah + Mike (never paired)
2. Jane + Mike (never paired)
3. John + Jane (1 commit together)
```

### Scenario 3: Cross-Team Collaboration

**Situation:** You want to encourage collaboration between frontend and backend teams.

**Commands:**
```bash
# Analyze all teams to see cross-team pairing
pairstair -window 1m

# Focus on backend team
pairstair --team backend -window 1m

# Compare with frontend
pairstair --team frontend -window 1m
```

**Strategy:** Use the insights to plan cross-team pairing sessions or shared projects.

### Scenario 4: Onboarding New Developer

**Situation:** A new developer "Alex" has joined and you want to track their integration.

**Updated `.team` file:**
```
[team]
Alice Smith <alice@company.com>
Bob Johnson <bob@company.com>
Carol Williams <carol@company.com>
Dave Brown <dave@company.com>
Alex Martinez <alex@company.com>  # New developer
```

**Commands:**
```bash
# Check Alex's pairing progress weekly
pairstair -window 1w

# Monthly review of integration
pairstair -window 1m -output html
```

**Tracking Progress:**
- Week 1: Alex pairs with Alice (mentor)
- Week 2: Alex pairs with Bob and Carol
- Week 3: Alex works independently but pairs with Dave
- Month review: Alex has paired with all team members

### Scenario 5: Remote Team Coordination

**Situation:** Distributed team needs structured pairing recommendations.

**Workflow:**
```bash
# Monday morning: Check last week's pairing
pairstair -window 1w

# Plan this week's pairs based on recommendations
pairstair -strategy least-recent

# Friday review: Generate weekly report
pairstair -window 1w -output html > weekly-report-$(date +%Y-%m-%d).html
```

**Automation Script:**
```bash
#!/bin/bash
# weekly-pairing-report.sh

DATE=$(date +%Y-%m-%d)
REPORT_DIR="pairing-reports"

mkdir -p $REPORT_DIR

echo "Generating weekly pairing report for $DATE"
pairstair -window 1w -output html > "$REPORT_DIR/pairing-report-$DATE.html"

echo "Last week's pairing summary:"
pairstair -window 1w

echo "This week's recommendations:"
pairstair -strategy least-recent
```

## Integration Examples

### GitHub Actions Workflow

**`.github/workflows/pairing-analysis.yml`:**
```yaml
name: Weekly Pairing Analysis

on:
  schedule:
    - cron: '0 9 * * MON'  # Every Monday at 9 AM
  workflow_dispatch:  # Allow manual trigger

jobs:
  analyze-pairing:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
        with:
          fetch-depth: 0  # Need full history
          
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.21'
          
      - name: Install PairStair
        run: go install github.com/gypsydave5/pairstair@latest
        
      - name: Generate Pairing Report
        run: |
          echo "## Weekly Pairing Analysis" > pairing-report.md
          echo "Generated on $(date)" >> pairing-report.md
          echo "" >> pairing-report.md
          echo "### Last Week" >> pairing-report.md
          echo "\`\`\`" >> pairing-report.md
          pairstair -window 1w >> pairing-report.md
          echo "\`\`\`" >> pairing-report.md
          echo "" >> pairing-report.md
          echo "### Recommendations" >> pairing-report.md
          echo "\`\`\`" >> pairing-report.md
          pairstair -strategy least-recent >> pairing-report.md
          echo "\`\`\`" >> pairing-report.md
          
      - name: Create Issue with Report
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            const report = fs.readFileSync('pairing-report.md', 'utf8');
            
            await github.rest.issues.create({
              owner: context.repo.owner,
              repo: context.repo.repo,
              title: `Weekly Pairing Analysis - ${new Date().toISOString().split('T')[0]}`,
              body: report,
              labels: ['pairing', 'weekly-report']
            });
```

### Slack Integration

**Bash script for Slack notifications:**
```bash
#!/bin/bash
# slack-pairing-notify.sh

SLACK_WEBHOOK_URL="your-slack-webhook-url"

# Generate pairing recommendations
RECOMMENDATIONS=$(pairstair -strategy least-recent | tail -n +5)

# Format for Slack
MESSAGE=$(cat << EOF
{
  "text": "Weekly Pairing Recommendations",
  "blocks": [
    {
      "type": "header",
      "text": {
        "type": "plain_text",
        "text": "🤝 Weekly Pairing Recommendations"
      }
    },
    {
      "type": "section",
      "text": {
        "type": "mrkdwn",
        "text": "\`\`\`$RECOMMENDATIONS\`\`\`"
      }
    },
    {
      "type": "context",
      "elements": [
        {
          "type": "mrkdwn",
          "text": "Generated by PairStair • $(date)"
        }
      ]
    }
  ]
}
EOF
)

# Send to Slack
curl -X POST -H 'Content-type: application/json' \
  --data "$MESSAGE" \
  "$SLACK_WEBHOOK_URL"
```

### Git Hooks Integration

**`.git/hooks/post-commit`:**
```bash
#!/bin/bash
# Auto-generate pairing report after commits with co-authors

# Check if this commit has co-authors
if git log -1 --pretty=format:"%B" | grep -q "Co-authored-by:"; then
    echo "Co-authored commit detected. Updating pairing analysis..."
    pairstair -window 1w > .pairing-summary.txt
    echo "Pairing summary updated in .pairing-summary.txt"
fi
```

## Output Format Examples

### Console Output

```
$ pairstair -window 2w

Pair Matrix (commits together):
        Alice  Bob  Carol  Dave
Alice     -    8     3     1
Bob       8    -     5     2
Carol     3    5     -     4
Dave      1    2     4     -

Total commits analyzed: 42
Commits with pairs: 15 (35.7%)
Unique pairs formed: 6

Recommendations (least recent):
1. Alice + Dave (1 commit together, 12 days ago)
2. Alice + Carol (3 commits together, 8 days ago)
3. Bob + Dave (2 commits together, 5 days ago)
```

### HTML Output Features

When using `-output html`, you get:

- **Interactive Matrix**: Click cells to see commit details
- **Time Series Charts**: Pairing trends over time
- **Developer Statistics**: Individual pairing frequency
- **Filtering Options**: Filter by date ranges, developers
- **Export Functions**: Save data as CSV or JSON

### JSON Output (Future Feature)

```bash
# Future feature example
pairstair -output json -window 1m
```

Would output:
```json
{
  "analysis_period": {
    "start": "2024-01-01T00:00:00Z",
    "end": "2024-01-31T23:59:59Z",
    "window": "1m"
  },
  "matrix": {
    "Alice": {"Bob": 8, "Carol": 3, "Dave": 1},
    "Bob": {"Alice": 8, "Carol": 5, "Dave": 2},
    "Carol": {"Alice": 3, "Bob": 5, "Dave": 4},
    "Dave": {"Alice": 1, "Bob": 2, "Carol": 4}
  },
  "recommendations": [
    {"pair": ["Alice", "Dave"], "score": 1, "last_paired": "2024-01-15"},
    {"pair": ["Alice", "Carol"], "score": 3, "last_paired": "2024-01-20"}
  ],
  "statistics": {
    "total_commits": 42,
    "paired_commits": 15,
    "pairing_percentage": 35.7,
    "unique_pairs": 6
  }
}
```

These examples should help you get started with PairStair and adapt it to your team's specific needs and workflows.
