package labeler

import (
	"strings"
	"testing"
)

func TestLoadConfig_Simple(t *testing.T) {
	yamlContent := `
label-a:
  - any:
    - changed-files:
      - any-glob-to-any-file: "*.go"
label-b:
  - all:
    - changed-files:
      - any-glob-to-any-file: "*.md"
label-c:
  - changed-files:
    - any-glob-to-any-file: "*.txt"
`
	cfg, err := LoadConfigFromReader(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(cfg) != 3 {
		t.Errorf("expected 2 labels, got %d", len(cfg))
	}
	for _, label := range []string{"label-a", "label-b", "label-c"} {
		lc, ok := cfg[label]
		if !ok {
			t.Errorf("%s not found in config", label)
		}
		if len(lc.Matcher) != 1 {
			t.Errorf("%s should have exactly one match, got %d", label, len(lc.Matcher))
		}
		m := lc.Matcher[0]
		if len(m.Any) == 0 && len(m.All) == 0 {
			t.Errorf("%s should have either Any or All rules, got neither", label)
		}
	}
}

func TestLoadConfig_AnchorAlias(t *testing.T) {
	yamlContent := `
default-rule: &def
  - any:
      - changed-files:
          - any-glob-to-any-file: "*.go"
label-a: *def
label-b:
  - all:
      - changed-files:
          - any-glob-to-any-file: "*.md"
`
	cfg, err := LoadConfigFromReader(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(cfg) != 3 {
		t.Errorf("expected 3 labels, got %d", len(cfg))
	}
	if _, ok := cfg["label-a"]; !ok {
		t.Error("label-a not found")
	}
	if _, ok := cfg["label-b"]; !ok {
		t.Error("label-b not found")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/no/such/file.yml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadConfig_HeadBranchAnchor(t *testing.T) {
	yamlContent := `
ci:
  - all:
    - changed-files:
      - any-glob-to-any-file: ".github/workflows/*"
    - head-branch: &ignore_ci
      - "^(?!ci/github-actions/).*"
test:
  - all:
    - changed-files:
      - any-glob-to-any-file: "**/*_test.go"
    - head-branch:
      - *ignore_ci
`
	cfg, err := LoadConfigFromReader(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(cfg) != 2 {
		t.Errorf("expected 2 labels, got %d", len(cfg))
	}
	for _, label := range []string{"ci", "test"} {
		lc, ok := cfg[label]
		if !ok || len(lc.Matcher[0].All) == 0 {
			t.Errorf("%s All not loaded", label)
			continue
		}
		var found bool
		for _, rule := range lc.Matcher[0].All {
			flat := rule.GetHeadBranch()
			if len(flat) == 1 && flat[0] == "^(?!ci/github-actions/).*" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s head-branch does not contain exactly one expected pattern", label)
		}
	}
}

func TestLoadConfig_ColorKey(t *testing.T) {
	yamlContent := `
ci:
  - any:
      - changed-files:
          - any-glob-to-any-file: '.github/*'
  - color: '#7c0bb2'
labeler:
  - changed-files:
      - any-glob-to-any-file: 'labeler/*'
  - color: '#123456'
documentation:
  - changed-files:
      - any-glob-to-any-file:
          - 'docs/*'
          - README.md
  - color: '#abcdef'
`
	cfg, err := LoadConfigFromReader(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(cfg) != 3 {
		t.Errorf("expected 3 labels, got %d", len(cfg))
	}
	cases := []struct {
		label string
		want  string
	}{
		{"ci", "#7c0bb2"},
		{"labeler", "#123456"},
		{"documentation", "#abcdef"},
	}
	for _, c := range cases {
		lc, ok := cfg[c.label]
		if !ok {
			t.Errorf("%s not loaded", c.label)
			continue
		}
		color := lc.Color
		if color != c.want {
			t.Errorf("%s color = %q, want %q", c.label, color, c.want)
		}
	}
}

func TestLoadConfig_ColorOnly(t *testing.T) {
	yamlContent := `
ci:
  - color: '#7c0bb2'
`
	cfg, err := LoadConfigFromReader(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(cfg) != 1 {
		t.Errorf("expected 3 labels, got %d", len(cfg))
	}
	lc, ok := cfg["ci"]
	if !ok {
		t.Errorf("ci not loaded")
	}
	if len(lc.Matcher) != 0 {
		t.Errorf("expected no matchers for ci, got %d", len(lc.Matcher))
	}
}

func TestLoadConfig_Codeowners(t *testing.T) {
	yamlContent := `
backend:
  - changed-files:
      - any-glob-to-any-file:
          - "backend/**"
  - codeowners:
      - backend/owner1
      - backend/owner2
`
	cfg, err := LoadConfigFromReader(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(cfg) != 1 {
		t.Errorf("expected 1 label, got %d", len(cfg))
	}
	lc, ok := cfg["backend"]
	if !ok {
		t.Errorf("backend not loaded")
	}
	if len(lc.Matcher) != 1 {
		t.Errorf("expected 1 matcher for backend, got %d", len(lc.Matcher))
	}
	if len(lc.Codeowners) != 2 {
		t.Errorf("expected 2 codeowners entries, got %d", len(lc.Codeowners))
	}
}

func TestLoadConfig_DescriptionKey(t *testing.T) {
	yamlContent := `
ci:
  - any:
      - changed-files:
          - any-glob-to-any-file: '.github/*'
  - color: '#7c0bb2'
  - description: 'Continuous Integration'
labeler:
  - changed-files:
      - any-glob-to-any-file: 'labeler/*'
  - description: 'Auto-labeling functionality'
documentation:
  - changed-files:
      - any-glob-to-any-file:
          - 'docs/*'
          - README.md
  - color: '#abcdef'
  - description: 'Documentation updates'
`
	cfg, err := LoadConfigFromReader(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(cfg) != 3 {
		t.Errorf("expected 3 labels, got %d", len(cfg))
	}
	cases := []struct {
		label string
		want  string
	}{
		{"ci", "Continuous Integration"},
		{"labeler", "Auto-labeling functionality"},
		{"documentation", "Documentation updates"},
	}
	for _, c := range cases {
		lc, ok := cfg[c.label]
		if !ok {
			t.Errorf("%s not loaded", c.label)
			continue
		}
		description := lc.Description
		if description != c.want {
			t.Errorf("%s description = %q, want %q", c.label, description, c.want)
		}
	}
}

func TestLoadConfig_DescriptionOnly(t *testing.T) {
	yamlContent := `
ci:
  - description: 'CI/CD pipeline'
`
	cfg, err := LoadConfigFromReader(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(cfg) != 1 {
		t.Errorf("expected 1 labels, got %d", len(cfg))
	}
	lc, ok := cfg["ci"]
	if !ok {
		t.Errorf("ci not loaded")
	}
	if len(lc.Matcher) != 0 {
		t.Errorf("expected no matchers for ci, got %d", len(lc.Matcher))
	}
	if lc.Description != "CI/CD pipeline" {
		t.Errorf("ci description = %q, want %q", lc.Description, "CI/CD pipeline")
	}
}

func TestLoadConfig_MixedMatchersWithColor(t *testing.T) {
	// Test case: changed-files with all-files-to-any-glob and color in a single label
	// These should create two separate matchers with Any conditions
	yamlContent := `
ci:
  - any:
    - changed-files:
      - any-glob-to-all-files: '.github/**'
  - all-files-to-any-glob:
    - '.github/**'
    - 'zizmor.yml'
  - color: '#7c0bb2'
`
	cfg, err := LoadConfigFromReader(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(cfg) != 1 {
		t.Errorf("expected 1 label, got %d", len(cfg))
	}
	lc, ok := cfg["ci"]
	if !ok {
		t.Errorf("ci not loaded")
	}
	// Should create two separate matchers, each with Any rules
	if len(lc.Matcher) != 2 {
		t.Errorf("expected 2 matchers for ci, got %d", len(lc.Matcher))
	}
	if lc.Color != "#7c0bb2" {
		t.Errorf("ci color = %q, want %q", lc.Color, "#7c0bb2")
	}
	// Verify both matchers have Any rules
	for i, matcher := range lc.Matcher {
		if len(matcher.Any) == 0 {
			t.Errorf("matcher %d: expected at least 1 Any rule, got %d", i, len(matcher.Any))
		}
		if len(matcher.All) != 0 {
			t.Errorf("matcher %d: expected no All rules, got %d", i, len(matcher.All))
		}
	}
}

func TestLoadConfig_MixedMatchersUnderAny(t *testing.T) {
	// Test case: both changed-files and all-files-to-any-glob under any condition
	// These should be merged into a single matcher with one Any rule containing both conditions
	yamlContent := `
ci:
  - any:
    - changed-files:
      - any-glob-to-all-files: '.github/**'
    - all-files-to-any-glob:
      - '.github/**'
      - 'zizmor.yml'
  - color: '#7c0bb2'
`
	cfg, err := LoadConfigFromReader(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(cfg) != 1 {
		t.Errorf("expected 1 label, got %d", len(cfg))
	}
	lc, ok := cfg["ci"]
	if !ok {
		t.Errorf("ci not loaded")
	}
	// Should create one matcher with one Any rule containing both conditions
	if len(lc.Matcher) != 1 {
		t.Errorf("expected 1 matcher for ci, got %d", len(lc.Matcher))
	}
	if lc.Color != "#7c0bb2" {
		t.Errorf("ci color = %q, want %q", lc.Color, "#7c0bb2")
	}
	if len(lc.Matcher) > 0 {
		matcher := lc.Matcher[0]
		if len(matcher.Any) != 2 {
			t.Errorf("expected 2 Any rule, got %d", len(matcher.Any))
		}
		if len(matcher.All) != 0 {
			t.Errorf("expected no All rules, got %d", len(matcher.All))
		}
		// The Any rule should contain both changed-files and all-files-to-any-glob conditions
		for _, anyRule := range matcher.Any {
			hasChangedFiles := len(anyRule.ChangedFiles) > 0
			hasAllFilesToAnyGlob := len(anyRule.AllFilesToAnyGlob) > 0

			if !hasChangedFiles {
				t.Error("Any rule should contain changed-files condition")
			}
			if hasAllFilesToAnyGlob {
				t.Error("Any rule should contain all-files-to-any-glob condition")
			}
		}

		if len(matcher.Any[1].ChangedFiles[0].AllFilesToAnyGlob) == 0 {
			t.Error("Second Any rule should contain all-files-to-any-glob condition")
		}
	}
}
