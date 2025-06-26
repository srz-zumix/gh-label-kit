package labeler

import (
	"gopkg.in/yaml.v3"
)

// LabelerConfig represents the YAML config for labeler (v5 compatible, supports per-label color key)
type LabelerConfig map[string][]LabelerMatch

// LabelerMatch supports per-label color key (actions/labeler v5 style)
type LabelerMatch struct {
	Any          []LabelerRule      `yaml:"any,omitempty"`
	All          []LabelerRule      `yaml:"all,omitempty"`
	ChangedFiles []ChangedFilesRule `yaml:"changed-files,omitempty"`
	BaseBranch   StringOrSliceRaw   `yaml:"base-branch,omitempty"`
	HeadBranch   StringOrSliceRaw   `yaml:"head-branch,omitempty"`
	Color        string             `yaml:"color,omitempty"`
}

type LabelerRule struct {
	ChangedFiles []ChangedFilesRule `yaml:"changed-files,omitempty"`
	BaseBranch   StringOrSliceRaw   `yaml:"base-branch,omitempty"`
	HeadBranch   StringOrSliceRaw   `yaml:"head-branch,omitempty"`
}

type ChangedFilesRule struct {
	AnyGlobToAnyFile   StringOrSlice `yaml:"any-glob-to-any-file,omitempty"`
	AnyGlobToAllFiles  StringOrSlice `yaml:"any-glob-to-all-files,omitempty"`
	AllGlobsToAnyFile  StringOrSlice `yaml:"all-globs-to-any-file,omitempty"`
	AllGlobsToAllFiles StringOrSlice `yaml:"all-globs-to-all-files,omitempty"`
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

func (m *LabelerMatch) GetBaseBranch() []string {
	return flattenStringOrSliceRaw(m.BaseBranch)
}
func (m *LabelerMatch) GetHeadBranch() []string {
	return flattenStringOrSliceRaw(m.HeadBranch)
}
func (r *LabelerRule) GetBaseBranch() []string {
	return flattenStringOrSliceRaw(r.BaseBranch)
}
func (r *LabelerRule) GetHeadBranch() []string {
	return flattenStringOrSliceRaw(r.HeadBranch)
}

// ColorOfLabel returns the color string for a label (if any), allowing for color-only elements in the config.
func ColorOfLabel(matches []LabelerMatch) string {
	for _, m := range matches {
		if m.Color != "" {
			return m.Color
		}
	}
	return ""
}
