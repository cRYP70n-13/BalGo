package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Services []string `yaml:"services"`
}

func Load(reader io.Reader) (*Config, error) {
	config := &Config{}

	yamlFile, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
