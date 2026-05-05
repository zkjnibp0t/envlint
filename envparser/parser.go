// Package envparser reads .env files into a key-value map.
package envparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Parse reads a .env file and returns a map of key-value pairs.
// Lines starting with '#' are treated as comments and ignored.
// Blank lines are skipped. Values may optionally be quoted.
func Parse(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envparser: cannot open file %q: %w", path, err)
	}
	defer f.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("envparser: invalid syntax on line %d: %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = stripQuotes(val)

		if key == "" {
			return nil, fmt.Errorf("envparser: empty key on line %d", lineNum)
		}

		env[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envparser: scan error: %w", err)
	}

	return env, nil
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
