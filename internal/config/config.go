package config

import (
    "log"
    "os"

    "gopkg.in/yaml.v2"
)

type Config struct {
    ClientID       string `yaml:"clientId"`
    ClientSecret   string `yaml:"clientSecret"`
    EntityTypeID   string `yaml:"entityTypeId"`
    ResourceType   string `yaml:"resourceType"`
    PlainIDEndpoint string `yaml:"plainIdEndpoint"`
}

func LoadConfig(path string) *Config {
    file, err := os.ReadFile(path)
    if err != nil {
        log.Fatalf("Failed to read config: %v", err)
    }
    var cfg Config
    if err := yaml.Unmarshal(file, &cfg); err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    return &cfg
}

