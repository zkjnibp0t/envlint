// Package redactor provides utilities for scrubbing sensitive values
// from .env maps before logging, displaying, or transmitting them.
package redactor

import (
	"strings"
)

// Rule defines a custom redaction rule based on key pattern and replacement.
type Rule struct {
	KeyPattern  string
	Replacement string
}

// Options controls redaction behaviour.
type Options struct {
	// ExtraRules are user-supplied patterns appended to the built-in list.
	ExtraRules []Rule
	// Placeholder overrides the default "[REDACTED]" string.
	Placeholder string
}

var defaultPatterns = []string{
	"password", "passwd", "secret", "token", "apikey", "api_key",
	"auth", "credential", "private", "cert", "key",
}

// Redact returns a copy of env with sensitive values replaced by placeholder.
func Redact(env map[string]string, opts Options) map[string]string {
	placeholder := "[REDACTED]"
	if opts.Placeholder != "" {
		placeholder = opts.Placeholder
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k, opts.ExtraRules) {
			result[k] = placeholder
		} else {
			result[k] = v
		}
	}
	return result
}

// isSensitive reports whether the key matches any built-in or custom pattern.
func isSensitive(key string, extra []Rule) bool {
	lower := strings.ToLower(key)
	for _, p := range defaultPatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	for _, r := range extra {
		if strings.Contains(lower, strings.ToLower(r.KeyPattern)) {
			return true
		}
	}
	return false
}

// Keys returns the list of keys that would be redacted for the given env.
func Keys(env map[string]string, opts Options) []string {
	var out []string
	for k := range env {
		if isSensitive(k, opts.ExtraRules) {
			out = append(out, k)
		}
	}
	return out
}
