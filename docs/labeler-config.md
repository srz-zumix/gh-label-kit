# Labeler Configuration

The `labeler` command uses a YAML configuration file (default: `.github/labeler.yml`) to define labeling rules. This configuration is compatible with [actions/labeler](https://github.com/actions/labeler) format, with additional support for `color` `description` and `codeowners` features.

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

```yaml
backend:
  - changed-files:
    - any-glob-to-any-file:
      - "api/**/*"
      - "server/**/*"
      - "**/*.go"
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

### Label Configuration

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
- The configuration is fully compatible with [actions/labeler](https://github.com/actions/labeler) with extensions for color and codeowners support
