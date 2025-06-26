# Future Feature Ideas for PairStair

This document captures ideas for potential future enhancements to PairStair. Feel free to add new ideas as inspiration strikes.

## Feature Ideas

- **Odd number of devs in recommender**: Handle this by recommending triples
- **Pair rotation**: Infer which developers are currently working together on a feature based on the most recent commits; suggest a rotation of the developers ensuring that developers do change their current pair
- **Ignoring Commits**: Introduce a way to ignore certain commits from the analysis, for example, by using a `.pairignore` file with commit message patterns (e.g., `chore:`, `docs:`) or a command-line flag
- **Flexible Date Ranges**: Instead of just a relative `--window`, allow for absolute date ranges with `--since="YYYY-MM-DD"` and `--until="YYYY-MM-DD"` flags for more precise analysis
- **Data Export Formats**: Add support for machine-readable output formats like JSON or CSV (`--output=json`). This would allow the data to be used in other tools or custom dashboards
- **Specify Repository Path**: Add a flag to specify the path to the git repository to analyze, rather than always using the current working directory. This would make the tool more flexible for scripting
- **Use relative .team file**: Look for team file at the git repository root, not in the current working directory
- **Team rotation suggestions**: Recommend optimal team rotations based on pairing history
- **Pairing analytics dashboard**: Web interface showing detailed pairing trends over time
- **Integration with project management tools**: Import team data from Jira, GitHub teams, etc.
- **Configurable pairing goals**: Set targets for minimum pairing frequency per developer
- **Historical trend analysis**: Show how pairing patterns change over time periods
- **Slack/Teams integration**: Post pairing recommendations to team channels
- **Git hooks integration**: Automatically track pairing without co-authored-by trailers
- **Multiple repository analysis**: Aggregate pairing data across multiple repos
- **Skill-based pairing**: Consider developer skills/expertise when recommending pairs
- **Calendar integration**: Factor in availability when making recommendations
- **Export formats**: CSV, JSON export for external analysis tools
- **Custom pairing rules**: Define organization-specific pairing policies

## Implementation Guidelines

When implementing new features:

1. Follow the established patterns in the codebase
2. Add comprehensive tests and documentation (README + man page)
3. Consider backward compatibility with existing `.team` files
4. Use semantic versioning (minor bump for new features, patch for fixes)
5. Update examples and documentation to showcase new capabilities

## How to Contribute Ideas

To add a new feature idea, simply add it to the list above with a brief description of the functionality.
