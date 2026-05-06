// Package linter combines all validation, formatting, auditing, and suggestion
// passes into a single unified lint run over a .env file and schema.
package linter

import (
	"github.com/user/envlint/auditor"
	"github.com/user/envlint/formatter"
	"github.com/user/envlint/schema"
	"github.com/user/envlint/suggester"
	"github.com/user/envlint/validator"
)

// Result holds the aggregated output of a full lint run.
type Result struct {
	ValidationErrors []validator.Error
	FormatIssues     []formatter.Issue
	AuditIssues      []auditor.Issue
	Suggestions      []suggester.Suggestion
	HasErrors        bool
}

// Run executes all lint passes against the provided env map and schema.
func Run(env map[string]string, s *schema.Schema) Result {
	var r Result

	// Validation pass
	r.ValidationErrors = validator.Validate(env, s)

	// Format / style pass
	// formatter.CheckStyle works on raw lines; we re-derive them from the map
	// for the purposes of style checking using key=value pairs.
	lines := envLines(env)
	r.FormatIssues = formatter.CheckStyle(lines)

	// Audit pass (unused / deprecated keys)
	r.AuditIssues = auditor.Audit(env, s)

	// Suggestion pass (typo hints for validation errors)
	r.Suggestions = suggester.Suggest(r.ValidationErrors, s)

	r.HasErrors = len(r.ValidationErrors) > 0 || hasErrors(r.AuditIssues)
	return r
}

// envLines converts an env map to a slice of "KEY=VALUE" strings so that
// the formatter's line-based checks can be applied.
func envLines(env map[string]string) []string {
	lines := make([]string, 0, len(env))
	for k, v := range env {
		lines = append(lines, k+"="+v)
	}
	return lines
}

func hasErrors(issues []auditor.Issue) bool {
	for _, i := range issues {
		if i.Level == auditor.LevelError {
			return true
		}
	}
	return false
}
