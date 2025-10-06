# gh-label-kit

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
gh labeler <pr-number...> [--repo <owner/repo>] [--config <path>] [--sync] [--dryrun] [--color <auto|always|never>]
```

Automatically add or remove labels to GitHub Pull Requests based on changed files, branch names, and a YAML config file (default: .github/labeler.yml).
This command behaves the same as [actions/labeler][labeler].

- Supports multiple PR numbers at once
- Supports glob/regex/negative lookahead patterns
- --sync: Remove labels that do not match any condition
- --dryrun/-n: Show results without actually changing labels
- --repo/-R: Target repository (default: current)
- --config: Path to config file (default: .github/labeler.yml)
- --color: Control color output

---

### repo copy: Bulk copy labels

```sh
gh repo copy <src> <dst...> [--force]
```

Copy labels from the src repository to one or more dst repositories.

- --force: Overwrite existing labels
- Multiple dst repositories supported

---

### repo list: List labels

```sh
gh repo list [--repo <owner/repo>]
```

Show all labels in the specified repository.

- --repo/-R: Target repository (default: current)

---

### repo sync: Sync label differences

```sh
gh repo sync <src> <dst...>
```

Synchronize label differences from src to one or more dst repositories.

---

### runner list: List GitHub Actions Runners

```sh
gh runner list [--repo <owner/repo>]
```

Show all GitHub Actions Runners for a repository or organization.

[labeler]: https://github.com/actions/labeler
