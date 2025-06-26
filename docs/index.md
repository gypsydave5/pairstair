---
layout: home
---

# PairStair

**Visualize and optimize software developer pairing from git history**

PairStair is a CLI tool that analyzes your git commit history to show how often developers have paired (committed together) over a configurable time window, and recommends new pairs to maximize collaboration.

## Quick Start

### Installation

```bash
# Homebrew
brew install gypsydave5/pairstair/pairstair

# Go install
go install github.com/gypsydave5/pairstair@latest
```

### Basic Usage

```bash
# Analyze the last week
pairstair

# Analyze a specific time window
pairstair -window 4w

# Analyze a sub-team
pairstair -team frontend

# Export as HTML
pairstair -output html
```

## Features

- 📊 **Pair Matrix**: Visual representation of developer collaboration
- 👥 **Team Support**: Organize teams with `.team` files and sub-teams
- 🎯 **Smart Recommendations**: Multiple strategies for suggesting optimal pairs
- 📈 **HTML Output**: Rich web interface for detailed analysis
- ⏰ **Flexible Time Windows**: Analyze any time period
- 🔧 **Git Integration**: Works with any git repository

## Example Output

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

## Documentation

- **[Installation Guide](installation/)** - Get PairStair up and running
- **[User Guide](guide/)** - Complete guide to using PairStair
- **[Examples](examples/)** - Real-world usage scenarios and workflows
- **[Features](features/)** - Current and planned functionality

## Why PairStair?

### For Developers
- 🎯 **Focus**: Know who you should pair with next
- 📈 **Growth**: Track your collaboration and skill sharing
- 🤝 **Community**: Build stronger relationships with teammates

### For Team Leads
- 👀 **Visibility**: Understand team collaboration patterns
- 🎯 **Optimization**: Make data-driven pairing decisions
- 📊 **Metrics**: Track team health and knowledge sharing

### For Organizations
- 🚀 **Productivity**: Optimize team performance through better collaboration
- 🧠 **Knowledge**: Prevent knowledge silos and single points of failure
- 👥 **Culture**: Foster a collaborative development culture

## Getting Started

1. **Install** PairStair using your preferred method
2. **Navigate** to your git repository
3. **Create** a `.team` file with your team members
4. **Run** `pairstair` to see your first analysis
5. **Explore** the recommendations and start pairing!

## Community

- 🐛 **Report Issues**: [GitHub Issues](https://github.com/gypsydave5/pairstair/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/gypsydave5/pairstair/discussions)
- 📖 **Documentation**: You're reading it!
- 💬 **Support**: Community support via GitHub

## Contributing

PairStair is open source and welcomes contributions:

- 🔧 **Code**: Submit pull requests for bug fixes and features
- 📝 **Documentation**: Help improve guides and examples
- 🧪 **Testing**: Report bugs and test new features
- 💡 **Ideas**: Share your thoughts on improvements

Ready to optimize your team's pairing? [Get started with installation!](installation/)

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

## Learn More

- [Installation Guide](installation.html) - Multiple installation methods
- [User Guide](guide.html) - Complete documentation with examples
- [Features](features.html) - Current features and roadmap
- [Examples](examples.html) - Real-world usage scenarios
