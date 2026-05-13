// Package templater provides functionality to render .env files from
// a template with placeholder substitution using a provided values map.
package templater

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Result holds the rendered output and metadata about substitutions.
type Result struct {
	Rendered     string
	Substituted  []string
	Unresolved   []string
}

var placeholderRe = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

// Render takes a template string and a values map, replacing all
// {{ KEY }} placeholders with their corresponding values.
// Keys not found in the map are left as-is and recorded in Unresolved.
func Render(template string, values map[string]string) Result {
	substituted := []string{}
	unresolved := []string{}
	seenUnresolved := map[string]bool{}

	output := placeholderRe.ReplaceAllStringFunc(template, func(match string) string {
		key := strings.TrimSpace(placeholderRe.FindStringSubmatch(match)[1])
		if val, ok := values[key]; ok {
			substituted = append(substituted, key)
			return val
		}
		if !seenUnresolved[key] {
			unresolved = append(unresolved, key)
			seenUnresolved[key] = true
		}
		return match
	})

	return Result{
		Rendered:    output,
		Substituted: substituted,
		Unresolved:  unresolved,
	}
}

// RenderFile reads a template file from disk and renders it using values.
func RenderFile(path string, values map[string]string) (Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Result{}, fmt.Errorf("templater: cannot read file %q: %w", path, err)
	}
	return Render(string(data), values), nil
}
