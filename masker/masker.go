// Package masker provides utilities for redacting sensitive values
// in .env files before logging or displaying them.
package masker

import "strings"

// SensitivePatterns is a list of key substrings that indicate a value
// should be masked in output.
var SensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

// Masked is the replacement string used for sensitive values.
const Masked = "***REDACTED***"

// IsSensitive reports whether the given key name matches any known
// sensitive pattern (case-insensitive).
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range SensitivePatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// MaskEnv returns a copy of the provided env map with sensitive values
// replaced by the Masked constant.
func MaskEnv(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		if IsSensitive(k) {
			result[k] = Masked
		} else {
			result[k] = v
		}
	}
	return result
}

// MaskValue returns the value itself if the key is not sensitive,
// or Masked if it is.
func MaskValue(key, value string) string {
	if IsSensitive(key) {
		return Masked
	}
	return value
}
