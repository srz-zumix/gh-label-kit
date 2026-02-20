package labeler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"

	"gopkg.in/yaml.v3"
)

func LoadConfigFromReader(r io.Reader, strictMode bool) (LabelerConfig, error) {
	// Read all content into a buffer so we can decode it twice if needed
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// First pass: try strict decoding to detect unknown fields
	var cfgStrict labelerYamlConfig
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfgStrict); err != nil {
		// Check if it's an unknown field error
		if strings.Contains(err.Error(), "field") && strings.Contains(err.Error(), "not found") {
			if strictMode {
				// In strict mode, return the error
				return nil, fmt.Errorf("config validation failed: %w", err)
			}
			logger.Warn("Config contains unknown or unsupported fields", "details", err.Error())
			logger.Warn("Unknown fields will be ignored. Please check the labeler configuration documentation")

			// Second pass: decode normally (allowing unknown fields)
			var cfg labelerYamlConfig
			if err := yaml.NewDecoder(bytes.NewReader(data)).Decode(&cfg); err != nil {
				return nil, err
			}
			logger.Debug("Config loaded successfully", "labels", len(cfg))
			return cfg.GetConfig(), nil
		}
		// If it's not an unknown field error, return it as actual error
		return nil, err
	}

	// Successfully loaded with strict validation
	logger.Debug("Config loaded successfully", "labels", len(cfgStrict))
	return cfgStrict.GetConfig(), nil
}

// ConfigFileExists checks if the config file exists at the given path.
func ConfigFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func LoadConfig(path string, strictMode bool) (LabelerConfig, error) {
	logger.Debug("Loading config from local file", "path", path, "strictMode", strictMode)
	f, err := os.Open(path)
	if err != nil {
		logger.Debug("Failed to open local config file", "path", path, "error", err)
		return nil, err
	}
	defer f.Close() // nolint
	cfg, err := LoadConfigFromReader(f, strictMode)
	if err != nil {
		logger.Debug("Failed to parse config file", "path", path, "error", err)
		return nil, err
	}
	logger.Debug("Successfully loaded config from local file", "path", path, "labels", len(cfg))
	return cfg, nil
}

// LoadConfigFromRepo loads a labeler config YAML from a GitHub repository using go-github's Contents API.
func LoadConfigFromRepo(ctx context.Context, g *gh.GitHubClient, repo repository.Repository, path string, ref *string, strictMode bool) (LabelerConfig, error) {
	refStr := "default"
	if ref != nil {
		refStr = *ref
	}
	logger.Debug("Loading config from repository", "owner", repo.Owner, "repo", repo.Name, "path", path, "ref", refStr, "strictMode", strictMode)
	fileContent, err := gh.GetRepositoryFileContent(ctx, g, repo, path, ref)
	if err != nil {
		logger.Debug("Failed to get file content from repository", "owner", repo.Owner, "repo", repo.Name, "path", path, "ref", refStr, "error", err)
		return nil, err
	}
	if fileContent == nil {
		logger.Debug("File not found in repository", "owner", repo.Owner, "repo", repo.Name, "path", path, "ref", refStr)
		return nil, os.ErrNotExist
	}
	content, err := fileContent.GetContent()
	if err != nil {
		logger.Debug("Failed to decode file content", "owner", repo.Owner, "repo", repo.Name, "path", path, "ref", refStr, "error", err)
		return nil, err
	}
	reader := strings.NewReader(content)
	cfg, err := LoadConfigFromReader(reader, strictMode)
	if err != nil {
		logger.Debug("Failed to parse config from repository", "owner", repo.Owner, "repo", repo.Name, "path", path, "ref", refStr, "error", err)
		return nil, err
	}
	logger.Debug("Successfully loaded config from repository", "owner", repo.Owner, "repo", repo.Name, "path", path, "ref", refStr, "labels", len(cfg))
	return cfg, nil
}
