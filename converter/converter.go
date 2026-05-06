// Package converter transforms .env files between different formats
// such as YAML, TOML, JSON, and shell export scripts.
package converter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Format represents a supported output format.
type Format string

const (
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
	FormatTOML  Format = "toml"
	FormatShell Format = "shell"
	FormatDotEnv Format = "dotenv"
)

// Convert transforms a map of env vars into the target format string.
func Convert(env map[string]string, format Format) (string, error) {
	switch format {
	case FormatJSON:
		return toJSON(env)
	case FormatYAML:
		return toYAML(env)
	case FormatTOML:
		return toTOML(env)
	case FormatShell:
		return toShell(env)
	case FormatDotEnv:
		return toDotEnv(env)
	default:
		return "", fmt.Errorf("unsupported format: %q", format)
	}
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func toJSON(env map[string]string) (string, error) {
	b, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func toYAML(env map[string]string) (string, error) {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		fmt.Fprintf(&sb, "%s: %q\n", k, env[k])
	}
	return sb.String(), nil
}

func toTOML(env map[string]string) (string, error) {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		fmt.Fprintf(&sb, "%s = %q\n", k, env[k])
	}
	return sb.String(), nil
}

func toShell(env map[string]string) (string, error) {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String(), nil
}

func toDotEnv(env map[string]string) (string, error) {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}
	return sb.String(), nil
}
