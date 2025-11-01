package labeler

import (
	"context"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/google/go-github/v73/github"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
)

// EditLabelsByConfig edits the given labels according to the config (color, etc). Returns the edited labels.
func EditLabelsByConfig(ctx context.Context, g *gh.GitHubClient, repo repository.Repository, labels []*github.Label, config LabelerConfig) ([]*github.Label, error) {
	labelMap := make(map[string]*github.Label)
	var edited []*github.Label
	for _, l := range labels {
		if l == nil || l.Name == nil {
			continue
		}
		labelMap[*l.Name] = l
	}
	for name, cfg := range config {
		color := cfg.Color
		description := cfg.Description
		if color == "" && description == "" {
			continue
		}
		if color != "" && color[0] == '#' {
			color = color[1:]
		}
		if l, ok := labelMap[name]; ok {
			needsUpdate := false
			if color != "" && (l.Color == nil || *l.Color != color) {
				l.Color = github.Ptr(color)
				needsUpdate = true
			}
			if description != "" && (l.Description == nil || *l.Description != description) {
				l.Description = github.Ptr(description)
				needsUpdate = true
			}
			if needsUpdate {
				result, err := gh.EditLabel(ctx, g, repo, *l.Name, l)
				if err != nil {
					return nil, err
				}
				edited = append(edited, result)
			}
		}
	}
	return edited, nil
}
