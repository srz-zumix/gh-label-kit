# Labeler Configuration

The `labeler` command uses a YAML configuration file (default: `.github/labeler.yml`) to define labeling rules. This configuration is compatible with [actions/labeler](https://github.com/actions/labeler) format, with additional support for `author`, `color`, `description`, `codeowners`, and `all-files-to-any-glob` features.

## Compatibility with actions/labeler

gh-label-kit is fully compatible with [actions/labeler](https://github.com/actions/labeler) configuration files. This means you can:

1. **Use the same `.github/labeler.yml` file** for both actions/labeler and gh-label-kit
2. **Add gh-label-kit specific features** without breaking actions/labeler compatibility

When actions/labeler encounters gh-label-kit specific fields (like `author`, `color`, `description`, `codeowners`, or `all-files-to-any-glob`), it simply ignores them. This allows you to enhance your configuration for gh-label-kit without creating separate configuration files.

### Example: Shared Configuration

```yaml
# This configuration works with both actions/labeler and gh-label-kit

documentation:
  - changed-files:
    - any-glob-to-any-file: 'docs/**'
  - color: "0075ca"  # gh-label-kit only
  - description: "Documentation changes"  # gh-label-kit only

# gh-label-kit specific: all changed files must be source code
source-only:
  - all-files-to-any-glob:  # gh-label-kit only, ignored by actions/labeler
    - "src/**/*"
    - "*.go"
  - color: "28a745"

# gh-label-kit specific: match by author
bot-dependency-updates:
  - author: 'dependabot\[bot\]'  # gh-label-kit only, ignored by actions/labeler
  - changed-files:
    - any-glob-to-any-file: 'go.mod'
```

In the above example:

- **actions/labeler** will recognize and process the `changed-files` rules for all labels
- **gh-label-kit** will additionally process the `author`, `color`, `description`, and `all-files-to-any-glob` fields
- The configuration file remains valid for both tools

## Basic Structure

```yaml
Documentation:
- changed-files:
  - any-glob-to-any-file: 'docs/*'
```

## Configuration Options

### Changed File Matching

The labeler supports various file matching strategies:

- **any-glob-to-any-file**: Match if any changed file matches any of the provided patterns (default behavior)
- **any-glob-to-all-files**: Match if any pattern matches all changed files
- **all-globs-to-any-file**: Match if all patterns match at least one changed file
- **all-globs-to-all-files**: Match if all patterns match all changed files
- **all-files-to-any-glob**: Match if all changed files match at least one of the provided patterns

```yaml
backend:
  - changed-files:
    - any-glob-to-any-file:
      - "api/**/*"
      - "server/**/*"
      - "**/*.go"

# Ensure all changed files are either source code or documentation
source-or-docs:
  - changed-files:
    - all-files-to-any-glob:
      - "src/**/*"
      - "docs/**/*"
      - "*.md"
```

You can also use `all-files-to-any-glob` at the top level (same level as `color` and `description`) for cleaner configuration:

```yaml
# All changed files must be within specific directories
core-files-only:
  - all-files-to-any-glob:
    - "src/**/*"
    - "core/**/*"
  - color: "0366d6"  # Blue
  - description: "Changes to core files"
```

#### Glob Patterns

For the changed files options you provide a [path glob](https://github.com/bmatcuk/doublestar?tab=readme-ov-file#patterns)

#### Extended Glob Patterns (extglob) (**Experimental**)

In addition to standard glob patterns, the labeler supports extended glob patterns for more advanced matching:

- **!(pattern)**: Match anything except the pattern (negation)
- **?(pattern)**: Match zero or one occurrence of the pattern
- **+(pattern)**: Match one or more occurrences of the pattern
- ***(pattern)**: Match zero or more occurrences of the pattern
- **@(pattern)**: Match exactly one occurrence of the pattern

Extended glob patterns can be combined with standard doublestar (`**`) patterns and support multiple alternatives using the pipe (`|`) separator.

**Important**: Extended glob patterns follow bash shell extglob behavior. Inside extglob expressions, wildcards (`*` and `?`) can match across path separators (`/`). This differs from standard glob patterns where `*` does not cross directory boundaries.

For example:

- `!(test)` can match `dir/file` (shell behavior)
- `!(*test*)` will exclude any path containing "test", including `src/test/file.go`
- In standard glob, `*` only matches within a single directory level

##### Extended Glob Examples

```yaml
# Match files that are NOT markdown
no-markdown:
  - changed-files:
    - any-glob-to-any-file: "!(*.md)"

# Match only Go source files or Go modules
go-files-only:
  - changed-files:
    - any-glob-to-any-file: "@(*.go|go.mod|go.sum)"

# Match files that are NOT test files
no-tests:
  - changed-files:
    - any-glob-to-any-file: "!(**/*_test.go)"

# Match either source code or documentation
source-or-docs:
  - changed-files:
    - any-glob-to-any-file: "@(src/**|docs/**)"
```

### Branch Matching

Labels can be applied based on branch names:

```yaml
feature:
  - head-branch: 
    - "^feature/.*"
    - "feat/.*"

hotfix:
  - base-branch: 
    - "main"
    - "master"
```

#### Branch Name Patterns

For the branches you provide a [regexp](https://github.com/dlclark/regexp2) to match against the branch name.

### Author Matching

Labels can be applied based on PR author. This is useful for labeling bot PRs, team member PRs, or specific user PRs.

```yaml
# Match bot users
bot:
  - author: '.*\[bot\]$'

# Match specific user
user-pr:
  - author: 'username'

# Match multiple users
team-pr:
  - author:
    - 'user1'
    - 'user2'
```

#### Author Patterns

For authors, you can use:

- **Regular expressions**: Match against the PR author username using [regexp2](https://github.com/dlclark/regexp2)
- **Team membership**: `@org/team-slug` - Match if the author is a member of the specified team
- **Negated team membership**: `!@org/team-slug` - Match if the author is NOT a member of the specified team

##### Team Membership Examples

```yaml
# Match team members
internal-team:
  - author: '@myorg/developers'

# Match non-team members (external contributors)
external-contributor:
  - author: '!@myorg/developers'

# Combine with other conditions
team-feature:
  - all:
    - author: '@myorg/frontend-team'
    - head-branch: 'feature/.*'
```

##### Common Author Patterns

```yaml
# Match Dependabot PRs
dependabot:
  - author: 'dependabot\[bot\]'

# Match any bot PR
bot:
  - author: '.*\[bot\]$'

# Match non-bot PRs
human:
  - author: '^(?!.*\[bot\]$).*$'

# Match specific username prefix
team-prefix:
  - author: '^team-.*'
```

#### Color

You can specify colors for labels using the `color` property:

```yaml
bug:
  - changed-files:
    - any-glob-to-any-file: "**/*.{js,ts}"
  - color: "d73a4a"  # Red color for bug labels

enhancement:
  - head-branch: "feature/**"
  - color: "a2eeef"  # Light blue color for enhancement labels
```

#### Description

You can specify descriptions for labels using the `description` property:

```yaml
bug:
  - changed-files:
    - any-glob-to-any-file: "**/*.{js,ts}"
  - color: "d73a4a"  # Red color for bug labels
  - description: "Something isn't working"

enhancement:
  - head-branch: "feature/**"
  - color: "a2eeef"  # Light blue color for enhancement labels
  - description: "New feature or request"
```

### CODEOWNERS Support

You can specify reviewers for labels using the `codeowners` property:

```yaml
team-frontend:
  - changed-files:
    - any-glob-to-any-file: "**/*.{js,ts}"
  - codeowners:
    - "@org/frontend-team"
    - "@srz-zumix"
```

CODEOWNERS is a feature that sends review requests when label conditions are met.
You can specify when to send review requests using the `--review-request` option:

- **addto**: Send review requests only when labels are added
- **always**: Send review requests whenever label conditions are met
- **ready_for_review**: Send review requests when labels are added to non-draft PRs, and also when ready_for_review activity occurs if conditions are met
- **always_reviewable**: Send review requests when label conditions are met for non-draft PRs
- **never**: Disable the review request feature

For `ready_for_review` or `always_reviewable` options, you need to include the `ready_for_review` activity as a trigger.
When running outside of GitHub Actions environment, the `ready_for_review` activity processing will not be performed.

```yaml
  pull_request:
    types:
      - opened
      - synchronize
      - reopened
      - ready_for_review # for review-request: ready_for_review/always_reviewable
```

## Advanced Examples

### Multiple Conditions

```yaml
critical-bug:
  - changed-files:
    - any-glob-to-any-file: "src/core/**/*"
  - head-branch: "hotfix/**"
  - color: "b60205"  # Dark red
  - description: "Critical bug that needs immediate attention"
```

### Complex File Patterns

```yaml
config-change:
  - changed-files:
    - all-globs-to-any-file:
      - "*.json"
      - "*.yml"
      - "*.yaml"
      - ".github/**/*"
  - color: "fef2c0"  # Light yellow
  - description: "Configuration file changes"

# All changed files must be source code (no tests, no docs)
source-files-only:
  - all-files-to-any-glob:
    - "src/**/*.go"
    - "lib/**/*.go"
  - color: "fbca04"  # Yellow
  - description: "Source code changes only"
```

### Team-based Labeling with CODEOWNERS

```yaml
needs-review-security:
  - codeowners:
    - "@org/security-team"
  - changed-files:
    - any-glob-to-any-file:
      - "auth/**/*"
      - "security/**/*"
  - color: "d4c5f9"  # Light purple
  - description: "Needs review from security team"
```

## Sync Labels

When using the `--sync` flag, the labeler will remove labels that don't match any condition in the configuration file:

```sh
gh label-kit labeler 123 --sync
```

This ensures that only relevant labels based on the current configuration are applied to the PR.

## Notes

- Glob patterns follow standard glob syntax
- The configuration is fully compatible with [actions/labeler](https://github.com/actions/labeler)
- gh-label-kit specific features (`author`, `color`, `description`, `codeowners`, `all-files-to-any-glob` at top-level) are safely ignored by actions/labeler, allowing you to use a single configuration file for both tools
