// Package schema handles loading and parsing the envlint schema definition.
package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// VarDefinition describes a single environment variable in the schema.
type VarDefinition struct {
	Name            string `yaml:"name"`
	Type            string `yaml:"type"`
	Required        bool   `yaml:"required"`
	Default         string `yaml:"default"`
	Pattern         string `yaml:"pattern"`
	Deprecated      bool   `yaml:"deprecated"`
	DeprecationNote string `yaml:"deprecation_note"`
}

// Schema is the top-level schema definition loaded from a YAML file.
type Schema struct {
	Version string          `yaml:"version"`
	Vars    []VarDefinition `yaml:"vars"`
}

// Load reads and parses a schema YAML file from the given path.
func Load(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading schema file: %w", err)
	}

	var s Schema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing schema YAML: %w", err)
	}

	if len(s.Vars) == 0 {
		return nil, fmt.Errorf("schema defines no variables")
	}

	return &s, nil
}
