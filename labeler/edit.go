package labeler

import (
	"context"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/google/go-github/v79/github"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
)

// EditLabelsByConfig edits the given labels according to the config (color, etc). Returns the edited labels.
func EditLabelsByConfig(ctx context.Context, g *gh.GitHubClient, repo repository.Repository, labels []*github.Label, config LabelerConfig) ([]*github.Label, error) {
	logger.Debug("Editing labels by config", "labelsCount", len(labels))
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
				logger.Debug("Updating label", "name", name, "color", color, "description", description)
				result, err := gh.EditLabel(ctx, g, repo, *l.Name, l)
				if err != nil {
					logger.Debug("Failed to update label", "name", name, "error", err)
					return nil, err
				}
				edited = append(edited, result)
			}
		}
	}
	logger.Debug("Finished editing labels", "editedCount", len(edited))
	return edited, nil
}
