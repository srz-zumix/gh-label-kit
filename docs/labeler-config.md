# Labeler Configuration

The `labeler` command uses a YAML configuration file (default: `.github/labeler.yml`) to define labeling rules. This configuration is compatible with [actions/labeler](https://github.com/actions/labeler) format, with additional support for `color` and `codeowners` features.

## Basic Structure

```yaml
Documentation:
- changed-files:
  - any-glob-to-any-file: 'docs/*'
```

## Configuration Options

### File Matching

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

### Extended Glob Patterns (extglob) (**Experimental**)

In addition to standard glob patterns, the labeler supports extended glob patterns for more advanced matching:

- **!(pattern)**: Match anything except the pattern (negation)
- **?(pattern)**: Match zero or one occurrence of the pattern
- **+(pattern)**: Match one or more occurrences of the pattern
- ***(pattern)**: Match zero or more occurrences of the pattern
- **@(pattern)**: Match exactly one occurrence of the pattern

Extended glob patterns can be combined with standard doublestar (`**`) patterns and support multiple alternatives using the pipe (`|`) separator.

#### Extended Glob Examples

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
    - "feature/**"
    - "feat/**"

hotfix:
  - base-branch: 
    - "main"
    - "master"
```

### Color Support

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

### Description Support

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
