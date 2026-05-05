package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// VarType represents the expected type of an environment variable.
type VarType string

const (
	TypeString  VarType = "string"
	TypeInt     VarType = "int"
	TypeBool    VarType = "bool"
	TypeURL     VarType = "url"
	TypeEmail   VarType = "email"
)

// VarDefinition describes a single environment variable rule.
type VarDefinition struct {
	Required    bool    `yaml:"required"`
	Type        VarType `yaml:"type"`
	Default     string  `yaml:"default"`
	Description string  `yaml:"description"`
	Pattern     string  `yaml:"pattern"`
}

// Schema holds the full set of variable definitions loaded from a schema file.
type Schema struct {
	Vars map[string]VarDefinition `yaml:"vars"`
}

// Load reads and parses a YAML schema file from the given path.
func Load(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading schema file %q: %w", path, err)
	}

	var s Schema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing schema file %q: %w", path, err)
	}

	if s.Vars == nil {
		s.Vars = make(map[string]VarDefinition)
	}

	return &s, nil
}
