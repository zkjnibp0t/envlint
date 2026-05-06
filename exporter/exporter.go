// Package exporter provides functionality to export validated env variables
// into different output formats such as shell export statements or Docker env-file format.
package exporter

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents the output format for exported variables.
type Format string

const (
	// FormatShell produces "export KEY=VALUE" lines suitable for shell sourcing.
	FormatShell Format = "shell"
	// FormatDocker produces "KEY=VALUE" lines suitable for Docker --env-file.
	FormatDocker Format = "docker"
	// FormatJSON produces a JSON object of key/value pairs.
	FormatJSON Format = "json"
)

// Export writes the given env map to w in the specified format.
// Keys are written in sorted order for deterministic output.
func Export(w io.Writer, env map[string]string, format Format) error {
	keys := sortedKeys(env)

	switch format {
	case FormatShell:
		return exportShell(w, env, keys)
	case FormatDocker:
		return exportDocker(w, env, keys)
	case FormatJSON:
		return exportJSON(w, env, keys)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

func exportShell(w io.Writer, env map[string]string, keys []string) error {
	for _, k := range keys {
		_, err := fmt.Fprintf(w, "export %s=%s\n", k, quoteIfNeeded(env[k]))
		if err != nil {
			return err
		}
	}
	return nil
}

func exportDocker(w io.Writer, env map[string]string, keys []string) error {
	for _, k := range keys {
		_, err := fmt.Fprintf(w, "%s=%s\n", k, env[k])
		if err != nil {
			return err
		}
	}
	return nil
}

func exportJSON(w io.Writer, env map[string]string, keys []string) error {
	_, err := fmt.Fprint(w, "{\n")
	if err != nil {
		return err
	}
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		_, err := fmt.Fprintf(w, "  %q: %q%s\n", k, env[k], comma)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprint(w, "}\n")
	return err
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n") {
		return fmt.Sprintf("%q", v)
	}
	return v
}
