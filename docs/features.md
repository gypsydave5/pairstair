---
layout: page
title: Features
permalink: /features/
---

# Features

PairStair provides comprehensive tools for analyzing and optimizing developer pairing in your team.

## Current Features

### 📊 Pair Matrix Analysis

**Visual representation of team collaboration patterns**

- Shows how often each pair of developers has worked together
- Identifies collaboration gaps and strong partnerships
- Configurable time windows (days, weeks, months, years)
- Clear matrix format for quick understanding

```bash
pairstair -window 2w
```

### 👥 Team Management

**Flexible team organization and analysis**

- **Team Files**: Define your team members with `.team` files
- **Sub-Teams**: Organize large teams into focused groups
- **Email Normalization**: Handle multiple email addresses per developer
- **Selective Analysis**: Focus on specific sub-teams with `--team` flag

**Example `.team` file:**
```
[frontend]
Alice Smith <alice@company.com>
Bob Johnson <bob@company.com>

[backend]
Carol Williams <carol@company.com>
Dave Brown <dave@company.com>
```

### 🎯 Smart Recommendations

**Multiple strategies for optimal pair suggestions**

- **Least Recent**: Pairs who haven't worked together recently
- **Least Paired**: Pairs who have worked together the fewest times using optimal matching
- **Customizable**: Choose the strategy that fits your team's needs

```bash
pairstair -strategy least-recent
pairstair -strategy least-paired
```

### 📈 Rich HTML Output

**Detailed web interface for comprehensive analysis**

- Interactive pair matrix with drill-down capabilities
- Visual charts and graphs showing trends
- Detailed statistics and metrics
- Sortable and filterable data tables
- Professional reports for stakeholders

```bash
pairstair -output html > report.html
```

### ⏰ Flexible Time Windows

**Analyze any time period that matters to your team**

- **Days**: `1d`, `7d`, `30d`
- **Weeks**: `1w`, `2w`, `4w`, `8w`
- **Months**: `1m`, `3m`, `6m`, `12m`
- **Years**: `1y`, `2y`

Perfect for different analysis needs:
- Daily standups: `-window 1d`
- Sprint reviews: `-window 2w`
- Quarterly assessments: `-window 3m`

### 🔧 Git Integration

**Seamless integration with your existing workflow**

- Works with any git repository
- Analyzes co-authored-by trailers in commit messages
- No additional setup or configuration required
- Supports all git hosting platforms (GitHub, GitLab, Bitbucket, etc.)

### 🚀 Performance Optimized

**Fast analysis even for large repositories**

- Efficient git log parsing
- Minimal memory footprint
- Quick startup time
- Scales well with repository size and team size

## Planned Features

*These features are planned for future releases. See our [roadmap](https://github.com/gypsydave5/pairstair/projects) for more details.*

### 🔄 Advanced Pair Rotation

**Intelligent rotation management**

- **Round Robin Strategy**: Ensure all possible pairs get equal opportunities over time
- Detect current working pairs from recent commits
- Suggest optimal rotations while maintaining project continuity
- Balance skill sharing with project deadlines
- Integration with sprint planning tools

### 📅 Calendar Integration

**Factor in availability when making recommendations**

- Connect with Google Calendar, Outlook, or other calendar systems
- Consider time zones for distributed teams
- Respect PTO and availability preferences
- Schedule pairing sessions automatically

### 📊 Analytics Dashboard

**Comprehensive pairing insights over time**

- Web-based dashboard showing long-term trends
- Team health metrics and KPIs
- Skill distribution analysis
- Pairing effectiveness measurements

### 🎯 Skill-Based Pairing

**Match developers based on expertise and learning goals**

- Define skill matrices for team members
- Recommend pairs for knowledge transfer
- Balance mentor/mentee relationships
- Track skill development over time

### 🔗 Tool Integrations

**Connect with your existing development tools**

- **Project Management**: Jira, Asana, Trello integration
- **Communication**: Slack, Microsoft Teams notifications
- **Code Review**: GitHub, GitLab merge request analysis
- **CI/CD**: Jenkins, GitHub Actions workflow integration

### 📈 Advanced Analytics

**Deeper insights into team dynamics**

- **Trend Analysis**: How pairing patterns change over time
- **Correlation Analysis**: Link pairing to code quality metrics
- **Productivity Metrics**: Measure pairing effectiveness
- **Team Health Scores**: Overall collaboration indicators

### 🌐 Multi-Repository Support

**Analyze pairing across multiple projects**

- Aggregate data from multiple repositories
- Cross-project collaboration tracking
- Organization-wide pairing insights
- Microservices team coordination

### 📋 Custom Reporting

**Tailored reports for different stakeholders**

- **Manager Reports**: High-level team health summaries
- **Developer Reports**: Personal pairing progress
- **HR Reports**: Team integration and culture metrics
- **Custom Templates**: Configurable report formats

### 🔧 Advanced Configuration

**Fine-tune PairStair for your organization**

- **Custom Pairing Rules**: Define organization-specific policies
- **Ignore Patterns**: Exclude certain types of commits
- **Weighting Systems**: Prioritize different types of collaboration
- **Notification Preferences**: Customizable alert thresholds

### 📱 Mobile and Desktop Apps

**Access pairing insights anywhere**

- Native mobile apps for iOS and Android
- Desktop applications for Windows, macOS, Linux
- Offline analysis capabilities
- Synchronized data across devices

## Feature Comparison

| Feature | Current | Planned |
|---------|---------|---------|
| Basic pair analysis | ✅ | ✅ |
| Team file support | ✅ | ✅ |
| Sub-team filtering | ✅ | ✅ |
| HTML output | ✅ | ✅ |
| Multiple strategies | ✅ | ✅ |
| Calendar integration | ❌ | ✅ |
| Skill-based pairing | ❌ | ✅ |
| Multi-repo support | ❌ | ✅ |
| Analytics dashboard | ❌ | ✅ |
| Mobile apps | ❌ | ✅ |
| API access | ❌ | ✅ |
| Real-time notifications | ❌ | ✅ |

## Contributing Ideas

Have an idea for a new feature? We'd love to hear from you! Here are ways to contribute:

1. **Open an Issue**: Describe your feature idea on [GitHub Issues](https://github.com/gypsydave5/pairstair/issues)
2. **Join Discussions**: Participate in [GitHub Discussions](https://github.com/gypsydave5/pairstair/discussions)
3. **Submit a PR**: Implement a feature and submit a pull request
4. **Share Use Cases**: Tell us how you're using PairStair and what would help

### Feature Request Template

When requesting a feature, please include:

- **Problem**: What challenge are you trying to solve?
- **Solution**: How do you envision the feature working?
- **Alternatives**: What workarounds are you currently using?
- **Context**: What's your team size, structure, and workflow?

### Priority Guidelines

Features are prioritized based on:

1. **Impact**: How many users would benefit?
2. **Effort**: How complex is the implementation?
3. **Alignment**: Does it fit PairStair's core mission?
4. **Community**: How much interest from the community?

## Version History

### v1.2.0 (Current)
- ✅ Sub-team support with `--team` flag
- ✅ Enhanced `.team` file format with sections
- ✅ Improved HTML output styling
- ✅ Better error handling and validation

### v1.1.0
- ✅ HTML output format
- ✅ Multiple recommendation strategies
- ✅ Flexible time window parsing
- ✅ Performance improvements

### v1.0.0
- ✅ Basic pair matrix analysis
- ✅ Team file support
- ✅ Console output
- ✅ Git integration
- ✅ Cross-platform support

## API Reference

*Future feature: Programmatic access to PairStair functionality*

### REST API Endpoints

```
GET /api/v1/analysis?window=2w&team=frontend
POST /api/v1/recommendations
GET /api/v1/teams
PUT /api/v1/teams/{id}
```

### SDK Libraries

Planned SDKs for popular languages:

- **JavaScript/TypeScript**: `npm install @pairstair/sdk`
- **Python**: `pip install pairstair-sdk`
- **Go**: `go get github.com/gypsydave5/pairstair/sdk`
- **Java**: Maven/Gradle integration

### Webhook Integration

```bash
# Configure webhooks for real-time updates
pairstair webhook add https://your-app.com/pairing-update
```

## Enterprise Features

*Contact us for enterprise licensing and support*

### Advanced Security
- SSO integration (SAML, OAuth)
- Role-based access control
- Audit logging
- Data encryption at rest

### Scalability
- Multi-tenant architecture
- High availability deployment
- Performance monitoring
- Custom SLA agreements

### Support
- 24/7 technical support
- Dedicated customer success manager
- Custom training and onboarding
- Priority feature requests

---

*Want to see a feature implemented? [Let us know!](https://github.com/gypsydave5/pairstair/issues/new)*
