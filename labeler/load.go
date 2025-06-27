package labeler

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"

	"gopkg.in/yaml.v3"
)

func LoadConfigFromReader(r io.Reader) (LabelerConfig, error) {
	var cfg labelerYamlConfig
	if err := yaml.NewDecoder(r).Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg.GetConfig(), nil
}

func LoadConfig(path string) (LabelerConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close() // nolint
	return LoadConfigFromReader(f)
}

// LoadConfigFromRepo loads a labeler config YAML from a GitHub repository using go-github's Contents API.
func LoadConfigFromRepo(ctx context.Context, g *gh.GitHubClient, repo repository.Repository, path string, ref *string) (LabelerConfig, error) {
	fileContent, err := gh.GetRepositoryFileContent(ctx, g, repo, path, ref)
	if err != nil {
		return nil, err
	}
	if fileContent == nil {
		return nil, os.ErrNotExist
	}
	content, err := fileContent.GetContent()
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(content)
	return LoadConfigFromReader(reader)
}
