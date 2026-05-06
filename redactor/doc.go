// Package redactor scrubs sensitive environment variable values before they
// are written to logs, reports, or external systems.
//
// It uses a built-in list of common sensitive key patterns (password, token,
// secret, etc.) and supports user-supplied custom rules via Options.ExtraRules.
//
// Usage:
//
//	env := map[string]string{
//		"DB_PASSWORD": "s3cr3t",
//		"APP_PORT":    "8080",
//	}
//	clean := redactor.Redact(env, redactor.Options{})
//	// clean["DB_PASSWORD"] == "[REDACTED]"
//	// clean["APP_PORT"]    == "8080"
package redactor
