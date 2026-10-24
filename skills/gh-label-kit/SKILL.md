---
name: gh-label-kit
description: GitHub CLI extension (gh label-kit) for GitHub label management and auto-labeling — auto-labeling pull requests from an actions/labeler compatible YAML config (changed files, branch, author and team membership, with label color/description and CODEOWNERS review requests), copying and syncing repository labels, adding/removing/setting/clearing labels on issues, pull requests and discussions, searching issues and discussions by label, and listing the labels of a milestone or of GitHub Actions self-hosted runners.
---

# gh-label-kit

Reference for gh-label-kit — a GitHub CLI extension for GitHub label management
and auto-labeling: labeling pull requests from a labeler config, keeping
repository labels in sync, editing the labels of issues, pull requests and
discussions, and listing labels by milestone or runner.

Version: 0.18.0

## Prerequisites

### Installation

```bash
gh extension install srz-zumix/gh-label-kit
```

### Authentication

```bash
gh auth login
```

Reading labels requires read access to the repository. Adding, removing or
syncing labels and requesting reviewers requires write (triage or higher)
access. Listing organization runners with `runner list --owner` requires
organization owner permission.

## CLI Structure

```text
gh label-kit                 # Root command
├── discussion               # Labels of discussions
│   ├── add                   # Add labels to a discussion
│   ├── clear                 # Remove all labels from a discussion
│   ├── list                  # List the labels of a discussion
│   ├── remove                # Remove labels from a discussion
│   ├── search                # Search discussions by query and labels
│   └── set                   # Replace all labels of a discussion
├── issue (alias: pr)        # Labels of issues and pull requests
│   ├── add                   # Add labels to an issue
│   ├── clear                 # Remove all labels from an issue
│   ├── list                  # List the labels of an issue
│   ├── remove                # Remove labels from an issue
│   ├── search                # Search issues by query and labels
│   └── set                   # Replace all labels of an issue
├── labeler                  # Auto-label pull requests from a config file
├── milestone                # Labels of milestones
│   └── list                  # List labels used by a milestone's issues and PRs
├── repo                     # Repository labels
│   ├── copy                  # Copy labels to other repositories
│   ├── list                  # List the labels of a repository
│   └── sync                  # Make other repositories match the source labels
├── runner                   # Runner labels
│   └── list                  # List self-hosted runners and their labels
└── skills                   # Manage the bundled agent skills
```

## Target Selection

- `-R`, `--repo [HOST/]OWNER/REPO` selects the repository. When omitted, the
  repository of the current directory is used.
- The `<number>` argument of `issue` and `discussion` subcommands and the
  `<milestone>` argument of `milestone list` also accept a GitHub URL. The
  repository is then taken from the URL.
- `pr` is an alias of `issue`. A pull request is an issue in the GitHub API, so
  `gh label-kit pr add 42 bug` labels pull request #42. With `pr search`, the
  `is:pr` qualifier is added to the query automatically.
- `issue search` and `discussion search` only add `repo:` / `org:` qualifiers
  when `--repo` / `--owner` is passed explicitly. Without them the search is not
  limited to the current repository.
- `runner list` targets the repository runners, or the organization runners when
  `--owner` is given without `--repo`.
- `repo copy` and `repo sync` read labels from `--repo` (source) and write them
  to every `<dst-repository>` argument. The source and destinations must be on
  the same host.

## Commands

### discussion add

Add one or more labels to a discussion in the repository.

```bash
gh label-kit discussion add <number> <label>... [--repo [HOST/]OWNER/REPO] \
  [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-R`, `--repo` | current repository | Select a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### discussion clear

Remove all labels from a discussion in the repository.

```bash
gh label-kit discussion clear <number> [--repo [HOST/]OWNER/REPO]
```

| Option | Default | Description |
| --- | --- | --- |
| `-R`, `--repo` | current repository | Select a repository |

### discussion list

List all labels attached to a discussion in the repository.

```bash
gh label-kit discussion list <number> [--repo [HOST/]OWNER/REPO] \
  [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-R`, `--repo` | current repository | Select a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### discussion remove

Remove one or more labels from a discussion in the repository.

```bash
gh label-kit discussion remove <number> <label>... [--repo [HOST/]OWNER/REPO] \
  [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-R`, `--repo` | current repository | Select a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### discussion search

Search discussions using a GitHub search query. Every `--label` is appended as a
`label:"NAME"` qualifier (all of them must match), `--owner` as `org:` and
`--repo` as `repo:`.

```bash
gh label-kit discussion search [QUERY...] [--repo [HOST/]OWNER/REPO] [--owner ORG] \
  [--label LABEL]... [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-l`, `--label` | - | Filter by label; repeatable or comma separated |
| `--owner` | - | Limit the search to an organization |
| `-R`, `--repo` | - | Limit the search to a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### discussion set

Set (replace) all labels of a discussion in the repository. Labels not listed are
removed.

```bash
gh label-kit discussion set <number> <label>... [--repo [HOST/]OWNER/REPO] \
  [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-R`, `--repo` | current repository | Select a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### issue add

Add one or more labels to an issue or pull request (`pr add`) in the repository.

```bash
gh label-kit issue add <number> <label>... [--repo [HOST/]OWNER/REPO] \
  [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-R`, `--repo` | current repository | Select a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### issue clear

Remove all labels from an issue or pull request in the repository.

```bash
gh label-kit issue clear <number> [--repo [HOST/]OWNER/REPO]
```

| Option | Default | Description |
| --- | --- | --- |
| `-R`, `--repo` | current repository | Select a repository |

### issue list

List all labels attached to an issue or pull request in the repository.

```bash
gh label-kit issue list <number> [--repo [HOST/]OWNER/REPO] \
  [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-R`, `--repo` | current repository | Select a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### issue remove

Remove one or more labels from an issue or pull request in the repository.

```bash
gh label-kit issue remove <number> <label>... [--repo [HOST/]OWNER/REPO] \
  [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-R`, `--repo` | current repository | Select a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### issue search

Search issues and pull requests using a GitHub search query. Every `--label` is
appended as a `label:"NAME"` qualifier (all of them must match), `--owner` as
`org:` and `--repo` as `repo:`. Invoked as `pr search`, `is:pr` is added.

```bash
gh label-kit issue search [QUERY...] [--repo [HOST/]OWNER/REPO] [--owner ORG] \
  [--label LABEL]... [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-l`, `--label` | - | Filter by label; repeatable or comma separated |
| `--owner` | - | Limit the search to an organization |
| `-R`, `--repo` | - | Limit the search to a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### issue set

Set (replace) all labels of an issue or pull request in the repository. Labels not
listed are removed.

```bash
gh label-kit issue set <number> <label>... [--repo [HOST/]OWNER/REPO] \
  [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-R`, `--repo` | current repository | Select a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### labeler

Automatically add or remove labels on pull requests based on changed files,
branch names, the PR author (including team membership) and a YAML config file.
The config format is compatible with
[actions/labeler](https://github.com/actions/labeler), with additional `author`,
`all-files-to-any-glob`, `color`, `description` and `codeowners` keys and
extended glob (extglob) patterns.

```bash
gh label-kit labeler <pr-number>... [--repo [HOST/]OWNER/REPO] [--config PATH|URL|OWNER/REPO[/PATH]@REF] \
  [--ref REF] [--skip-local-config] [--strict] [--no-hidden] [--sync] [--dryrun] \
  [--review-request addto|always|ready_for_review|always_reviewable|never|none] \
  [--name-only] [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--config` | `.github/labeler.yml` | Config location: a local path, a path in the repository, a GitHub URL (`https://github.com/OWNER/REPO[/tree/REF\|/blob/REF/PATH]`) or the actions format `OWNER/REPO[/PATH]@REF` |
| `-n`, `--dryrun` | `false` | Log the label changes and review requests without applying them |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `--name-only` | `false` | Print only the label names, comma separated |
| `--no-hidden` | `false` | Exclude hidden files (starting with `.`) from glob matching |
| `--ref` | `$GITHUB_SHA`, else the PR head branch | Git reference to load the config from the repository |
| `-R`, `--repo` | current repository | Select a repository |
| `--review-request` | `addto` | When to request reviews from the `codeowners` of matched labels (see below) |
| `--skip-local-config` | `false` | Ignore a local config file and load it from the repository |
| `--strict` | `false` | Treat unknown config fields as errors instead of warnings |
| `--sync` | `false` | Remove labels that no longer match any condition |
| `-t`, `--template` | - | Format JSON output using a Go template |

Config resolution: when `--config` points to an existing local file it is used
(unless `--skip-local-config`); a local file that fails to parse is an error and
does not fall back. Otherwise the config is read from the repository at `--ref`,
then `$GITHUB_SHA`, then — when exactly one PR number is given — the PR head
branch.

`--review-request` modes:

- `addto` — request reviews only for labels newly added by this run.
- `always` — request reviews for every label whose conditions match.
- `ready_for_review` — like `addto` for non-draft PRs; on a `ready_for_review`
  event in GitHub Actions, behaves like `always`.
- `always_reviewable` — like `always`, but skips draft PRs.
- `never` / `none` — never request reviews.

Inside GitHub Actions the command sets the `new-labels` and `all-labels` step
outputs. Labels on the PR that have a `color` or `description` in the config
are updated to match it.

Minimal config example:

```yaml
documentation:
  - changed-files:
    - any-glob-to-any-file: 'docs/**'
  - color: "0075ca"
  - description: "Documentation changes"

dependencies:
  - author: 'dependabot\[bot\]'

frontend:
  - changed-files:
    - any-glob-to-any-file: '**/*.{js,ts}'
  - codeowners:
    - "@my-org/frontend-team"
```

Author patterns are regular expressions, `@ORG/TEAM` (member of the team) or
`!@ORG/TEAM` (not a member). Full reference:
[docs/labeler-config.md](https://github.com/srz-zumix/gh-label-kit/blob/main/docs/labeler-config.md).

### milestone list

List all labels attached to the issues and pull requests in a milestone. The
milestone is given by its number or its URL.

```bash
gh label-kit milestone list <milestone> [--repo [HOST/]OWNER/REPO] \
  [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-R`, `--repo` | current repository | Select a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### repo copy

Copy all labels from the source repository to the destination repositories.
Labels that already exist in a destination are skipped unless `--force` is
given. Labels are never deleted.

```bash
gh label-kit repo copy <dst-repository>... [--repo [HOST/]OWNER/REPO] [--force]
```

| Option | Default | Description |
| --- | --- | --- |
| `-f`, `--force` | `false` | Overwrite the color and description of labels that already exist in the destination |
| `-R`, `--repo` | current repository | Source repository |

### repo list

List all labels in the repository.

```bash
gh label-kit repo list [--repo [HOST/]OWNER/REPO] \
  [--color always|never|auto] [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--color` | `auto` | Use color in output: `always`, `never` or `auto` |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `-R`, `--repo` | current repository | Select a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### repo sync

Make the labels of each destination repository match the source repository:
missing labels are created, labels that differ are updated, and labels that do
not exist in the source are deleted. Without `--force`, only destination-only
labels that are not used by any issue or pull request are deleted; with
`--force`, they are deleted even when in use.

```bash
gh label-kit repo sync <dst-repository>... [--repo [HOST/]OWNER/REPO] [--force]
```

| Option | Default | Description |
| --- | --- | --- |
| `-f`, `--force` | `false` | Also delete destination-only labels that are still in use |
| `-R`, `--repo` | current repository | Source repository |

`repo sync` is destructive. Prefer `repo copy` when the destination labels must
be kept.

### runner list

List the GitHub Actions self-hosted runners and their labels (`ID`, `NAME`, `OS`,
`STATUS`, `BUSY`, `LABELS`) of a repository, or of an organization when
`--owner` is given without `--repo`.

```bash
gh label-kit runner list [--repo [HOST/]OWNER/REPO | --owner ORG] \
  [--format json] [--jq EXPRESSION] [--template TEMPLATE]
```

| Option | Default | Description |
| --- | --- | --- |
| `--format` | - | Output format: `json`. Table output is used when not specified |
| `-q`, `--jq` | - | Filter JSON output using a jq expression |
| `--owner` | - | Select an organization |
| `-R`, `--repo` | current repository | Select a repository |
| `-t`, `--template` | - | Format JSON output using a Go template |

### skills

Install and manage the agent skills bundled with gh-label-kit.

```bash
gh label-kit skills list
gh label-kit skills install|update|reinstall|uninstall|status [-scope user|repo] [-prefix DIR] [-dry-run] [-force]
```

| Option | Default | Description |
| --- | --- | --- |
| `-dry-run` | `false` | Print what would happen without making changes |
| `-force` | `false` | Overwrite unmanaged skills |
| `-prefix` | - | Skill installation directory (overrides `-scope`) |
| `-scope` | `user` | `user` (`~/.agents/skills`) or `repo` (`<repo-root>/.agents/skills`) |

## Global Options

| Option | Default | Description |
| --- | --- | --- |
| `--http-timeout` | `30s` | Timeout for each GitHub API request |
| `-L`, `--log-level` | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `--read-only` | `false` | Run in read-only mode (blocks create, update and delete API calls) |

## Common Workflows

### Preview and apply auto-labeling to a pull request

```bash
# Check what would change first
gh label-kit labeler 123 --dryrun

# Apply, removing labels that no longer match
gh label-kit labeler 123 --sync
```

### Auto-label with a shared config from another repository

```bash
gh label-kit labeler 123 --config my-org/shared-config/.github/labeler.yml@main
```

### Auto-label in GitHub Actions

This example is for pull requests from branches in the same repository. For
pull requests from forks, the `pull_request` event receives a read-only token,
so applying labels and requesting reviewers fails.

```yaml
on:
  pull_request:
    types: [opened, synchronize, reopened, ready_for_review]

jobs:
  labeler:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write
    steps:
      - uses: actions/checkout@v4
      - run: gh extension install srz-zumix/gh-label-kit
        env:
          GH_TOKEN: ${{ github.token }}
      - run: gh label-kit labeler ${{ github.event.pull_request.number }} --sync --review-request ready_for_review
        env:
          GH_TOKEN: ${{ github.token }}
```

### Share a label set across repositories

```bash
# Add missing labels only
gh label-kit repo copy my-org/repo-a my-org/repo-b -R my-org/label-template

# Also overwrite the color and description of existing labels
gh label-kit repo copy my-org/repo-a -R my-org/label-template --force

# Make repo-a exactly match the template (deletes unused extra labels)
gh label-kit repo sync my-org/repo-a -R my-org/label-template
```

### Edit labels of issues and pull requests

```bash
gh label-kit issue add 42 bug priority:high
gh label-kit pr remove 128 wip
gh label-kit issue set 42 bug
gh label-kit issue clear 42
gh label-kit discussion add 7 question
```

### Find issues and pull requests with labels

```bash
# Open bugs with high priority in a repository
gh label-kit issue search is:open -R my-org/my-repo -l bug -l priority:high

# Pull requests across an organization
gh label-kit pr search is:open --owner my-org -l needs-review
```

### Summarize the labels of a milestone

```bash
gh label-kit milestone list 3 -R my-org/my-repo
```

### Check which runner labels are available

```bash
gh label-kit runner list -R my-org/my-repo
gh label-kit runner list --owner my-org --format json --jq '[.[].labels[].name] | unique'
```

## Troubleshooting

| Symptom | Cause / Resolution |
| --- | --- |
| `failed to load local config` | A local `--config` file exists but is invalid. Fix the YAML; the command does not fall back to the repository config in this case. |
| `failed to load config from repository` | The config was not found at the resolved ref. Pass `--ref`, or use `OWNER/REPO/PATH@REF` / a GitHub URL in `--config`. |
| Labeler warns about unknown fields | The config contains keys not understood by gh-label-kit. Fix the typo, or pass `--strict` to fail on them. |
| `labeler` does not remove stale labels | Labels are only removed with `--sync`. |
| No review requests are sent | `--review-request addto` only requests for newly added labels, draft PRs are skipped by `ready_for_review` / `always_reviewable`, and labels need a `codeowners` entry. |
| `ready_for_review` behaves like `addto` locally | The `ready_for_review` event is only detected inside GitHub Actions. |
| Search returns results from other repositories | `issue search` / `discussion search` only scope the query with `--repo` / `--owner` or `repo:` / `org:` in the query. |
| `source and destination repositories must be on the same host` | `repo copy` / `repo sync` cannot cross hosts. Run the command against each host separately. |
| `repo sync` left some extra labels | Without `--force`, labels still used by an issue or pull request are kept. |
| `unable to parse milestone number` | `milestone list` takes the milestone number or URL, not its title. |
| `runner list --owner` returns 403 | Listing organization runners requires organization owner permission. |
| Write commands fail with a read-only error | `--read-only` is set; remove it to apply changes. |
