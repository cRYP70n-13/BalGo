package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

// NOTE: Now i need a way to start marking the servers as healthy or not healthy
// and in that case we will need to wait for them again or just remove them from our list
// in case they will never come back again.
// Also we can do something more fancy which is to make a specific logic for each service
// everything will be defined in the config file and it will be our source of truth
// which algorithm to use which each one of them and which endpoints to hit
// or what sort of health check we should use, gRPC HTTP TCP conn ...
// Log our events what we are doing and how is everything going.
// More importantly use best practices and enough sleep you sucker

type Service struct {
	Name     string    `yaml:"name"`
	Address  string    `yaml:"address"`
	Strategy string    `yaml:"strategy"`
	Replicas []Replica `yaml:"replicas"`
}

type Replica struct {
	Address string `yaml:"address"`
}

type Config struct {
	Services []Service `yaml:"services"`
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
