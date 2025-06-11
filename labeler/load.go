package labeler

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfigFromReader(r io.Reader) (LabelerConfig, error) {
	var cfg LabelerConfig
	if err := yaml.NewDecoder(r).Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadConfig(path string) (LabelerConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadConfigFromReader(f)
}
