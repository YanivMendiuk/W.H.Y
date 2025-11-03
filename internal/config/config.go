package config

import (
    "os"

    "gopkg.in/yaml.v2"
)

// Config holds all the configuration values for the webhook application.
// Each field corresponds to a key in your YAML configuration file (webhook-config.yaml).
type Config struct {
    ClientID        string `yaml:"clientId"`       
    ClientSecret    string `yaml:"clientSecret"`    
    EntityTypeID    string `yaml:"entityTypeId"`    
    ResourceType    string `yaml:"resourceType"`    
    PlainIDEndpoint string `yaml:"plainIdEndpoint"` 
    Action          string `yaml:"action"`          
}

// LoadConfig reads YAML file and returns a Config struct
func LoadConfig(path string) (*Config, error) {
    file, err := os.ReadFile(path) // Read file from path
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := yaml.Unmarshal(file, &cfg); err != nil { // Parse YAML into struct
        return nil, err
    }

    return &cfg, nil // Return populated config
}
