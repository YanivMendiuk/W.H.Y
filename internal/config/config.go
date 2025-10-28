package config

import (
    "os"

    "gopkg.in/yaml.v2"
)

type Config struct {
    ClientID        string `yaml:"clientId"`
    ClientSecret    string `yaml:"clientSecret"`
    EntityTypeID    string `yaml:"entityTypeId"`
    ResourceType    string `yaml:"resourceType"`
    PlainIDEndpoint string `yaml:"plainIdEndpoint"`
}

// LoadConfig now returns (*Config, error)
func LoadConfig(path string) (*Config, error) {
    file, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := yaml.Unmarshal(file, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

