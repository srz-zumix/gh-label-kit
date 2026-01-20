package labeler

import (
	"maps"
	"slices"

	"gopkg.in/yaml.v3"
)

type LabelerConfig map[string]LabelerLabelConfig
type LabelerLabelConfig struct {
	Matcher     []LabelerMatch
	Color       string
	Description string
	Codeowners  []string
}

type LabelerMatch struct {
	Any []LabelerRule
	All []LabelerRule
}

// LabelerConfig represents the YAML config for labeler (v5 compatible, supports per-label color key)
type labelerYamlConfig map[string][]labelerYamlMatch

// LabelerMatch supports per-label color key (actions/labeler v5 style)
type labelerYamlMatch struct {
	Any                []LabelerRule      `yaml:"any,omitempty"`
	All                []LabelerRule      `yaml:"all,omitempty"`
	ChangedFiles       []ChangedFilesRule `yaml:"changed-files,omitempty"`
	AllFilesToAnyGlob  StringOrSlice      `yaml:"all-files-to-any-glob,omitempty"`
	BaseBranch         StringOrSliceRaw   `yaml:"base-branch,omitempty"`
	HeadBranch         StringOrSliceRaw   `yaml:"head-branch,omitempty"`
	Author             StringOrSliceRaw   `yaml:"author,omitempty"`
	Color              string             `yaml:"color,omitempty"`
	Description        string             `yaml:"description,omitempty"`
	Codeowners         StringOrSlice      `yaml:"codeowners,omitempty"`
}

type LabelerRule struct {
	ChangedFiles      []ChangedFilesRule `yaml:"changed-files,omitempty"`
	AllFilesToAnyGlob StringOrSlice      `yaml:"all-files-to-any-glob,omitempty"`
	BaseBranch        StringOrSliceRaw   `yaml:"base-branch,omitempty"`
	HeadBranch        StringOrSliceRaw   `yaml:"head-branch,omitempty"`
	Author            StringOrSliceRaw   `yaml:"author,omitempty"`
}

type ChangedFilesRule struct {
	AnyGlobToAnyFile   StringOrSlice `yaml:"any-glob-to-any-file,omitempty"`
	AnyGlobToAllFiles  StringOrSlice `yaml:"any-glob-to-all-files,omitempty"`
	AllGlobsToAnyFile  StringOrSlice `yaml:"all-globs-to-any-file,omitempty"`
	AllGlobsToAllFiles StringOrSlice `yaml:"all-globs-to-all-files,omitempty"`
	AllFilesToAnyGlob  StringOrSlice `yaml:"all-files-to-any-glob,omitempty"`
}

type StringOrSliceRaw any

type StringOrSlice []string

func (s *StringOrSlice) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		var single string
		if err := value.Decode(&single); err != nil {
			return err
		}
		*s = []string{single}
		return nil
	case yaml.SequenceNode:
		var result []string
		for _, elem := range value.Content {
			switch elem.Kind {
			case yaml.ScalarNode:
				var v string
				if err := elem.Decode(&v); err != nil {
					return err
				}
				result = append(result, v)
			case yaml.SequenceNode:
				var inner StringOrSlice
				if err := elem.Decode(&inner); err != nil {
					return err
				}
				result = append(result, inner...)
			}
		}
		*s = result
		return nil
	default:
		return nil
	}
}

func flattenStringOrSliceRaw(v any) []string {
	switch vv := v.(type) {
	case nil:
		return nil
	case string:
		return []string{vv}
	case []string:
		return vv
	case []any:
		var result []string
		for _, e := range vv {
			result = append(result, flattenStringOrSliceRaw(e)...)
		}
		return result
	}
	return nil
}

func (m *labelerYamlMatch) GetBaseBranch() []string {
	return flattenStringOrSliceRaw(m.BaseBranch)
}
func (m *labelerYamlMatch) GetHeadBranch() []string {
	return flattenStringOrSliceRaw(m.HeadBranch)
}
func (m *labelerYamlMatch) GetAuthor() []string {
	return flattenStringOrSliceRaw(m.Author)
}
func (r *LabelerRule) GetBaseBranch() []string {
	return flattenStringOrSliceRaw(r.BaseBranch)
}
func (r *LabelerRule) GetHeadBranch() []string {
	return flattenStringOrSliceRaw(r.HeadBranch)
}
func (r *LabelerRule) GetAuthor() []string {
	return flattenStringOrSliceRaw(r.Author)
}

// ColorOfLabel returns the color string for a label (if any), allowing for color-only elements in the config.
func colorOfLabel(matches []labelerYamlMatch) string {
	for _, m := range matches {
		if m.Color != "" {
			return m.Color
		}
	}
	return ""
}

func descriptionOfLabel(matches []labelerYamlMatch) string {
	for _, m := range matches {
		if m.Description != "" {
			return m.Description
		}
	}
	return ""
}

func codeownersOfLabel(matches []labelerYamlMatch) []string {
	ownerSet := make(map[string]struct{})
	for _, m := range matches {
		for _, owner := range m.Codeowners {
			ownerSet[owner] = struct{}{}
		}
	}
	return slices.Collect(maps.Keys(ownerSet))
}

func (r *labelerYamlConfig) GetConfig() LabelerConfig {
	cfg := make(LabelerConfig, len(*r))
	for label, matches := range *r {
		matchers := []LabelerMatch{}
		for _, m := range matches {
			m.Normalize()
			if len(m.Any) != 0 || len(m.All) != 0 {
				matchers = append(matchers, LabelerMatch{
					Any: m.Any,
					All: m.All,
				})
			}
		}
		cfg[label] = LabelerLabelConfig{
			Matcher:     matchers,
			Color:       colorOfLabel(matches),
			Description: descriptionOfLabel(matches),
			Codeowners:  codeownersOfLabel(matches),
		}
	}
	return cfg
}

func (r *LabelerRule) Normalize() {
	// Integrate top-level all-files-to-any-glob into changed-files
	if len(r.AllFilesToAnyGlob) > 0 {
		r.ChangedFiles = append(r.ChangedFiles, ChangedFilesRule{
			AllFilesToAnyGlob: r.AllFilesToAnyGlob,
		})
		r.AllFilesToAnyGlob = nil // Clear to avoid duplication
	}
}

func (m *labelerYamlMatch) Normalize() {
	// Integrate top-level all-files-to-any-glob into changed-files
	if len(m.AllFilesToAnyGlob) > 0 {
		m.ChangedFiles = append(m.ChangedFiles, ChangedFilesRule{
			AllFilesToAnyGlob: m.AllFilesToAnyGlob,
		})
		m.AllFilesToAnyGlob = nil // Clear to avoid duplication
	}

	anyRules := make([]LabelerRule, 0)
	if m.BaseBranch != nil {
		anyRules = append(anyRules, LabelerRule{BaseBranch: m.GetBaseBranch()})
		m.BaseBranch = nil // Clear to avoid duplication
	}
	if m.HeadBranch != nil {
		anyRules = append(anyRules, LabelerRule{HeadBranch: m.GetHeadBranch()})
		m.HeadBranch = nil // Clear to avoid duplication
	}
	if m.Author != nil {
		anyRules = append(anyRules, LabelerRule{Author: m.GetAuthor()})
		m.Author = nil // Clear to avoid duplication
	}
	if len(m.ChangedFiles) > 0 {
		anyRules = append(anyRules, LabelerRule{ChangedFiles: m.ChangedFiles})
		m.ChangedFiles = nil // Clear to avoid duplication
	}

	if len(anyRules) > 0 {
		if m.Any == nil {
			m.Any = anyRules
		} else {
			m.Any = append(m.Any, anyRules...)
		}
	}

	for i := range m.Any {
		m.Any[i].Normalize()
	}
	for i := range m.All {
		m.All[i].Normalize()
	}
}
