[![Go Doc][godoc-image]][godoc-url]
[![Build Status][workflow-image]][workflow-url]
[![Test Coverage][codecov-image]][codecov-url]

# changelog

Changelog is a lightweight changelog generator for GitHub repositories.
Inspired by the popular [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator),
it takes a more focused approach to generating release notes.
It is designed to be **simple**, **fast**, and **dependency-free**.

## Why?

[github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator)
is great, battle-proven, and just works! So why build another tool?

This project is not a reinvention.
It is a practical Go alternative for teams that want easier integration and fewer moving parts.

For a long time, I used the Ruby gem to generate changelogs.
But when creating GitHub releases, I still needed extra tooling to run *github_changelog_generator* and extract only the new entries.

That setup introduced a few recurring problems:

  - Limited control over the gem version on developer machines
  - Extra effort to install and maintain a Ruby environment in CI and containers
  - Occasional breakages caused by changes in external dependencies

I wanted a Go implementation of this excellent idea (see https://github.com/github-changelog-generator/github-changelog-generator/issues/714):

  - A self-contained binary that requires no external tooling
  - A library that can be imported with standard Go module versioning

Starting fresh in Go also makes it possible to:

  - Simplify the user experience
  - Improve performance and efficiency
  - Remove legacy workarounds that are no longer needed
  - Build a clearer roadmap for future improvements

## Quick Start

### Install

```
brew install neatplatform/brew/changelog
```

For other platforms, you can download the binary from the [latest release](https://github.com/neatplatform/changelog/releases/latest).

### Examples

```bash
# Simply generate a changelog
changelog -access-token=$GITHUB_TOKEN

# Assign unreleased changes (changes without a tag) to a future tag that has not been yet created.
changelog -access-token=$GITHUB_TOKEN -future-tag v0.1.0
```

### Help

<details>
  <summary>changelog -help</summary>

```
  changelog is a simple command-line tool for generating changelogs based on issues and pull/merge requests.
  It expects one configured git remote repository.

  You can also have a changelog.yaml file in your repository for configuring how changelogs are generated.
  For more information, please see https://github.com/neatplatform/changelog

  Supported Remote Repositories:

    • GitHub (github.com)

  Usage: changelog [flags]

  Flags:

    -help                         Show the help text
    -version                      Print the version number

    -access-token                 The OAuth access token for making API calls
                                  The default value is read from the CHANGELOG_ACCESS_TOKEN environment variable

    -file                         The output file for the generated changelog (default: CHANGELOG.md)
    -base                         An optional file for appending the generated changelog to it 
                                  This option can only be used when generating the changelog for the first time
    -print                        Print the generated changelog to STDOUT (default: false)
                                  If this option is enabled, all logs will be disabled
    -verbose                      Show the verbosity logs (default: false)

    -from-tag                     Changelog will be generated for all changes after this tag (default: last tag on changelog)
    -to-tag                       Changelog will be generated for all changes before this tag (default: last git tag)
    -future-tag                   A future tag for all unreleased changes (changes after the last git tag) 
    -exclude-tags                 These tags will be excluded from changelog 
    -exclude-tags-regex           A POSIX-compliant regex for excluding certain tags from changelog 

    -issues-selection             Include closed issues in changelog (values: none|all|labeled) (default: all)
    -issues-include-labels        Include issues with these labels 
    -issues-exclude-labels        Exclude issues with these labels (default: duplicate,invalid,question,wontfix)
    -issues-grouping              Grouping style for issues (values: simple|milestone|label) (default: label)
    -issues-summary-labels        Labels for summary group (default: summary,release-summary)
    -issues-removed-labels        Labels for removed group (default: removed)
    -issues-breaking-labels       Labels for breaking group (default: breaking,backward-incompatible)
    -issues-deprecated-labels     Labels for deprecated group (default: deprecated)
    -issues-feature-labels        Labels for feature group (default: feature)
    -issues-enhancement-labels    Labels for enhancement group (default: enhancement)
    -issues-bug-labels            Labels for bug group (default: bug)
    -issues-security-labels       Labels for security group (default: security)

    -merges-selection             Include merged pull/merge requests in changelog (values: none|all|labeled) (default: all)
    -merges-branch                Include pull/merge requests merged into this branch (default: default remote branch)
    -merges-include-labels        Include merges with these labels 
    -merges-exclude-labels        Exclude merges with these labels 
    -merges-grouping              Grouping style for pull/merge requests (values: simple|milestone|label) (default: simple)
    -merges-summary-labels        Labels for summary group 
    -merges-removed-labels        Labels for removed group 
    -merges-breaking-labels       Labels for breaking group 
    -merges-deprecated-labels     Labels for deprecated group 
    -merges-feature-labels        Labels for feature group 
    -merges-enhancement-labels    Labels for enhancement group 
    -merges-bug-labels            Labels for bug group 
    -merges-security-labels       Labels for security group 

    -release-url                  An external release URL with the '{tag}' placeholder for the release tag

  Examples:

    changelog
    changelog -access-token=<your-access-token>
    changelog -access-token=<your-access-token> -base=HISTORY.md
    changelog -access-token=<your-access-token> -future-tag=v0.1.0
```
</details>

### Spec File

Add a configuration file to your repository to customize how changelogs are generated.

<details>
  <summary>changelog.yaml</summary>

```yaml
general:
  file: CHANGELOG.md
  base: HISTORY.md
  print: true
  verbose: false

tags:
  exclude: [ prerelease, candidate ]
  exclude-regex: (.*)-(alpha|beta)

issues:
  selection: labeled
  include-labels: [ breaking, bug, defect, deprecated, enhancement, feature, highlight, improvement, incompatible, privacy, removed, security, summary ]
  exclude-labels: [ documentation, duplicate, invalid, question, wontfix ]
  grouping: milestone
  summary-labels: [ summary, highlight ]
  removed-labels: [ removed ]
  breaking-labels: [ breaking, incompatible ]
  deprecated-labels: [ deprecated ]
  feature-labels: [ feature ]
  enhancement-labels: [ enhancement, improvement ]
  bug-labels: [ bug, defect ]
  security-labels: [ security, privacy ]

merges:
  selection: labeled
  branch: production
  include-labels: [ breaking, bug, defect, deprecated, enhancement, feature, highlight, improvement, incompatible, privacy, removed, security, summary ]
  exclude-labels: [ documentation, duplicate, invalid, question, wontfix ]
  grouping: label
  summary-labels: [ summary, highlight ]
  removed-labels: [ removed ]
  breaking-labels: [ breaking, incompatible ]
  deprecated-labels: [ deprecated ]
  feature-labels: [ feature ]
  enhancement-labels: [ enhancement, improvement ]
  bug-labels: [ bug, defect ]
  security-labels: [ security, privacy ]

content:
  release-url: https://storage.artifactory.com/project/releases/{tag}
```
</details>

## Features

  - Single, dependency-free, cross-platform binary
  - Generate changelogs from issues and pull/merge requests
  - Support for unreleased and draft releases
  - Filter tags by name or regex patterns
  - Organize entries by labels or milestones
  - Include/exclude items based on label filters

## Expected Behavior

When you run *changelog* inside a Git directory, the following steps occur:

  1. The remote repository is identified using the first *remote* name configured (SSH and HTTPS URLs are supported).
  2. The existing changelog file (if present) is scanned to identify which Git tags have already been documented.
  3. The candidate tags are filtered using the `exclude-tags` and `exclude-tags-regex` options (if specified).
  4. API calls are made to the remote platform (GitHub) to fetch all closed issues and merged pull requests.
  5. Issues are filtered based on the `selection`, `include-labels`, and `exclude-labels` settings.
  6. Pull requests are filtered based on the `selection`, `branch`, `include-labels`, and `exclude-labels` settings.
  7. Issues are organized using the `grouping` option (label, milestone, or simple).
  8. Pull requests are organized using the `grouping` option (label, milestone, or simple).
  9. The changelog is generated and written to the output file.

## Resources


[godoc-url]: https://pkg.go.dev/github.com/neatplatform/changelog
[godoc-image]: https://pkg.go.dev/badge/github.com/neatplatform/changelog
[workflow-url]: https://github.com/neatplatform/changelog/actions/workflows/go.yml
[workflow-image]: https://github.com/neatplatform/changelog/actions/workflows/go.yml/badge.svg
[codecov-url]: https://codecov.io/gh/neatplatform/changelog
[codecov-image]: https://codecov.io/gh/neatplatform/changelog/graph/badge.svg
