# gh-label-kit

asis-author-pr test

gh-label-kit is a set of gh extension commands for GitHub label management and auto-labeling.

## Installation

To install the tool, you can use the following command:

```sh
gh extension install srz-zumix/gh-label-kit
```

## Commands & Usage

---

### labeler: Auto-label PRs

```sh
gh label-kit labeler <pr-number...> [--repo <owner/repo>] [--config <path>] [--sync] [--dryrun] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>] [--name-only] [--ref <string>]
```

Automatically add or remove labels to GitHub Pull Requests based on changed files, branch name, and a YAML config file (default: .github/labeler.yml).
Supports glob/regex patterns and syncLabels option for label removal. This command behaves the same as [actions/labeler][labeler].

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

### repo copy: Copy labels between repositories

```sh
gh label-kit repo copy <dst-repository...> [--repo <owner/repo>] [--force]
```

Copy all labels from the source repository to the destination repositories. If a label already exists in the destination, it will be skipped unless --force is specified.

- --force/-f: Overwrite existing labels in the destination repository
- --repo/-R: Repository in the format 'owner/repo' (source repository)

---

### repo list: List labels

```sh
gh label-kit repo list [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

List all labels in the specified repository.

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

---

### repo sync: Sync label differences

```sh
gh label-kit repo sync <dst-repository...> [--repo <owner/repo>] [--force]
```

Sync all labels from the source repository to the destination repositories. If a label already exists in the destination, it will be updated if --force is specified.

- --force/-f: Overwrite existing labels in the destination repository
- --repo/-R: The repository in the format 'owner/repo' (source repository)

---

### runner list: List GitHub Actions runner labels

```sh
gh label-kit runner list [--repo <owner/repo>] [--owner <organization>] [--format <json>] [--jq <expression>] [--template <string>]
```

List all GitHub Actions runner labels in the specified repository.

- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --owner: Specify the organization name
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

---

### issue add: Add label(s) to issue

```sh
gh label-kit issue add <number> <label>... [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

Add one or more labels to a issue in the repository.

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

---

### issue clear: Remove all labels from issue

```sh
gh label-kit issue clear <number> [--repo <owner/repo>]
```

Remove all labels from a issue in the repository.

- --repo/-R: Repository in the format 'owner/repo'

---

### issue list: List labels for issue

```sh
gh label-kit issue list <number> [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

List all labels attached to a issue in the repository.

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

---

### issue remove: Remove label(s) from issue

```sh
gh label-kit issue remove <number> <label>... [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

Remove one or more labels from a issue in the repository.

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

---

### issue set: Set labels for issue (replace all)

```sh
gh label-kit issue set <number> <label>... [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

Set (replace) all labels for a issue in the repository.

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

---

### milestone list: List labels for milestone

```sh
gh label-kit milestone list <milestone> [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

List all labels attached to issues and PRs in the specified milestone.

- --color: Use color in diff output (always|never|auto, default: auto)
- --format: Output format (json)
- --jq: Filter JSON output using a jq expression
- --repo/-R: Repository in the format 'owner/repo'
- --template/-t: Format JSON output using a Go template

[labeler]: https://github.com/actions/labeler
