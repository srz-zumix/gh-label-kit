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
`
	cfg, err := LoadConfigFromReader(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(cfg) != 2 {
		t.Errorf("expected 2 labels, got %d", len(cfg))
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
		matches := cfg[label]
		if len(matches) == 0 || len(matches[0].All) == 0 {
			t.Errorf("%s All not loaded", label)
			continue
		}
		var found bool
		for _, rule := range matches[0].All {
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
