// Package scanner detects common security and configuration anti-patterns
// in .env files, such as hardcoded secrets, weak values, or dangerous defaults.
package scanner

import (
	"fmt"
	"regexp"
	"strings"
)

// Severity represents the risk level of a finding.
type Severity string

const (
	SeverityHigh   Severity = "HIGH"
	SeverityMedium Severity = "MEDIUM"
	SeverityLow    Severity = "LOW"
)

// Finding represents a single detected issue in the env map.
type Finding struct {
	Key      string
	Value    string
	Rule     string
	Message  string
	Severity Severity
}

var weakValues = map[string]bool{
	"password": true, "secret": true, "changeme": true,
	"example": true, "test": true, "1234": true, "admin": true,
}

var privateKeyPattern = regexp.MustCompile(`(?i)-----BEGIN .* PRIVATE KEY-----`)
var ipLocalhostPattern = regexp.MustCompile(`^(127\.0\.0\.1|localhost)$`)

// Scan inspects the provided env map and returns a list of findings.
func Scan(env map[string]string) []Finding {
	var findings []Finding

	for key, value := range env {
		lower := strings.ToLower(strings.TrimSpace(value))

		// Rule: weak/placeholder value for sensitive-looking keys
		if isSensitiveKey(key) && weakValues[lower] {
			findings = append(findings, Finding{
				Key:      key,
				Value:    value,
				Rule:     "weak-value",
				Message:  fmt.Sprintf("%q appears to be a placeholder or weak value", value),
				Severity: SeverityHigh,
			})
		}

		// Rule: private key material embedded in value
		if privateKeyPattern.MatchString(value) {
			findings = append(findings, Finding{
				Key:      key,
				Value:    "[REDACTED]",
				Rule:     "embedded-private-key",
				Message:  "value appears to contain a private key",
				Severity: SeverityHigh,
			})
		}

		// Rule: localhost used in a URL-like variable
		if isURLKey(key) && ipLocalhostPattern.MatchString(value) {
			findings = append(findings, Finding{
				Key:      key,
				Value:    value,
				Rule:     "localhost-url",
				Message:  "URL points to localhost — likely a development default",
				Severity: SeverityMedium,
			})
		}

		// Rule: empty value for non-optional-looking key
		if value == "" && isSensitiveKey(key) {
			findings = append(findings, Finding{
				Key:      key,
				Value:    "",
				Rule:     "empty-sensitive",
				Message:  "sensitive key has an empty value",
				Severity: SeverityMedium,
			})
		}
	}

	return findings
}

func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for _, kw := range []string{"password", "secret", "token", "api_key", "apikey", "auth", "private"} {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func isURLKey(key string) bool {
	lower := strings.ToLower(key)
	return strings.Contains(lower, "host") || strings.Contains(lower, "url") || strings.Contains(lower, "addr")
}
