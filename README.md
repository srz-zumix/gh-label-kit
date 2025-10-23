# gh-label-kit

[![CI/CD](https://github.com/srz-zumix/gh-label-kit/actions/workflows/ci.yaml/badge.svg)](https://github.com/srz-zumix/gh-label-kit/actions/workflows/ci.yaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

gh-label-kit is a powerful gh extension that simplifies GitHub label management. It provides a comprehensive set of commands to streamline your labeling workflow, from auto-labeling pull requests to syncing labels across repositories. This tool is designed to save you time and ensure consistency in your projects.

## Installation

To install the tool, you can use the following command:

```sh
gh extension install srz-zumix/gh-label-kit
```

## Table of Contents

- [Commands & Usage](#commands--usage)
  - [labeler](#labeler)
    - [labeler: Auto-label PRs](#labeler-auto-label-prs)
  - [repo](#repo)
    - [repo copy: Copy labels between repositories](#repo-copy-copy-labels-between-repositories)
    - [repo list: List labels](#repo-list-list-labels)
    - [repo sync: Sync label differences](#repo-sync-sync-label-differences)
  - [runner](#runner)
    - [runner list: List GitHub Actions runner labels](#runner-list-list-github-actions-runner-labels)
  - [issue](#issue)
    - [issue add: Add label(s) to issue](#issue-add-add-labels-to-issue)
    - [issue clear: Remove all labels from issue](#issue-clear-remove-all-labels-from-issue)
    - [issue list: List labels for issue](#issue-list-list-labels-for-issue)
    - [issue remove: Remove label(s) from issue](#issue-remove-remove-labels-from-issue)
    - [issue set: Set labels for issue (replace all)](#issue-set-set-labels-for-issue-replace-all)
  - [milestone](#milestone)
    - [milestone list: List labels for milestone](#milestone-list-list-labels-for-milestone)

## Commands & Usage

### labeler

#### labeler: Auto-label PRs

```sh
gh label-kit labeler <pr-number...> [--repo <owner/repo>] [--config <path>] [--sync] [--dryrun] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>] [--name-only] [--ref <string>]
```

Automatically add or remove labels to GitHub Pull Requests based on changed files, branch name, and a YAML config file (default: .github/labeler.yml).
Supports glob/regex patterns and syncLabels option for label removal. This command behaves the same as [actions/labeler][labeler].

**Example:**
```sh
# Auto-label PR #123 in the current repository
gh label-kit labeler 123

# Auto-label PR #123 in another repository with a dry run
gh label-kit labeler 123 --repo owner/repo --dryrun
```

- --color: Use color in diff output (auto|never|always, default: auto)
- --config: Path to labeler config YAML file (default: .github/labeler.yml)
- --dryrun/-n: Dry run: do not actually set labels
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --name-only: Output only team names
- --ref: Git reference (branch, tag, or commit SHA) to load config from repository
- --repo/-R: Target repository in the format 'owner/repo'
- --sync: Remove labels not matching any condition
- --template/-t: Format JSON output using a Go template

The `labeler` command uses a YAML configuration file to define labeling rules. The configuration format is compatible with [actions/labeler][labeler], with additional support for `color` and `codeowners` features.

For detailed configuration documentation, see [docs/labeler-config.md](docs/labeler-config.md).

---

### repo

#### repo copy: Copy labels between repositories

```sh
gh label-kit repo copy <dst-repository...> [--repo <owner/repo>] [--force]
```

Copy all labels from the source repository to the destination repositories. If a label already exists in the destination, it will be skipped unless --force is specified.

**Example:**
```sh
# Copy labels from the current repo to 'owner/repo2'
gh label-kit repo copy owner/repo2

# Copy labels from 'owner/repo1' to 'owner/repo2' and 'owner/repo3', overwriting existing labels
gh label-kit repo copy owner/repo2 owner/repo3 --repo owner/repo1 --force
```

- --force/-f: Overwrite existing labels in the destination repository
- --repo/-R: Repository in the format 'owner/repo' (source repository)

#### repo list: List labels

```sh
gh label-kit repo list [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

List all labels in the specified repository.

**Example:**
```sh
# List labels in the current repository
gh label-kit repo list

# List labels in 'owner/repo' in JSON format
gh label-kit repo list --repo owner/repo --format json
```

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

#### repo sync: Sync label differences

```sh
gh label-kit repo sync <dst-repository...> [--repo <owner/repo>] [--force]
```

Sync all labels from the source repository to the destination repositories. If a label already exists in the destination, it will be updated if --force is specified.

**Example:**
```sh
# Sync labels from the current repo to 'owner/repo2'
gh label-kit repo sync owner/repo2

# Sync labels from 'owner/repo1' to 'owner/repo2', overwriting existing labels
gh label-kit repo sync owner/repo2 --repo owner/repo1 --force
```

- --force/-f: Overwrite existing labels in the destination repository
- --repo/-R: The repository in the format 'owner/repo' (source repository)

---

### runner

#### runner list: List GitHub Actions runner labels

```sh
gh label-kit runner list [--repo <owner/repo>] [--owner <organization>] [--format <json>] [--jq <expression>] [--template <string>]
```

List all GitHub Actions runner labels in the specified repository.

**Example:**
```sh
# List runner labels for the current repository
gh label-kit runner list

# List runner labels for an organization
gh label-kit runner list --owner my-org
```

- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --owner: Specify the organization name
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

---

### issue

#### issue add: Add label(s) to issue

```sh
gh label-kit issue add <number> <label>... [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

Add one or more labels to a issue in the repository.

**Example:**
```sh
# Add the 'bug' and 'help wanted' labels to issue #456
gh label-kit issue add 456 bug "help wanted"
```

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

#### issue clear: Remove all labels from issue

```sh
gh label-kit issue clear <number> [--repo <owner/repo>]
```

Remove all labels from a issue in the repository.

**Example:**
```sh
# Remove all labels from issue #456
gh label-kit issue clear 456
```

- --repo/-R: Repository in the format 'owner/repo'

#### issue list: List labels for issue

```sh
gh label-kit issue list <number> [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

List all labels attached to a issue in the repository.

**Example:**
```sh
# List labels for issue #456
gh label-kit issue list 456
```

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

#### issue remove: Remove label(s) from issue

```sh
gh label-kit issue remove <number> <label>... [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

Remove one or more labels from a issue in the repository.

**Example:**
```sh
# Remove the 'bug' label from issue #456
gh label-kit issue remove 456 bug
```

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

#### issue set: Set labels for issue (replace all)

```sh
gh label-kit issue set <number> <label>... [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

Set (replace) all labels for a issue in the repository.

**Example:**
```sh
# Set the labels for issue #456 to 'enhancement' and 'v2.0'
gh label-kit issue set 456 enhancement v2.0
```

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

---

### milestone

#### milestone list: List labels for milestone

```sh
gh label-kit milestone list <milestone> [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

List all labels attached to issues and PRs in the specified milestone.

**Example:**
```sh
# List all labels for issues and PRs in the 'v1.0' milestone
gh label-kit milestone list v1.0
```

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

[labeler]: https://github.com/actions/labeler
