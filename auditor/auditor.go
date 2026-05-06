// Package auditor provides functionality to audit .env files for
// unused variables (present in .env but not defined in schema) and
// deprecated variables flagged in the schema.
package auditor

import (
	"fmt"

	"github.com/user/envlint/schema"
)

// Issue represents a single audit finding.
type Issue struct {
	Key      string
	Severity string // "warning" or "info"
	Message  string
}

// Result holds the outcome of an audit run.
type Result struct {
	Unused     []Issue
	Deprecated []Issue
}

// HasIssues returns true if any issues were found.
func (r Result) HasIssues() bool {
	return len(r.Unused) > 0 || len(r.Deprecated) > 0
}

// Audit checks the parsed env map against the schema and returns
// unused and deprecated variable findings.
func Audit(env map[string]string, s *schema.Schema) Result {
	var result Result

	// Build a set of schema keys for quick lookup.
	schemaKeys := make(map[string]schema.VarDefinition, len(s.Vars))
	for _, v := range s.Vars {
		schemaKeys[v.Name] = v
	}

	// Find env keys not present in schema.
	for key := range env {
		if _, defined := schemaKeys[key]; !defined {
			result.Unused = append(result.Unused, Issue{
				Key:      key,
				Severity: "warning",
				Message:  fmt.Sprintf("%s is set in .env but not defined in schema", key),
			})
		}
	}

	// Find deprecated schema keys that are still in use.
	for _, v := range s.Vars {
		if v.Deprecated {
			if _, present := env[v.Name]; present {
				msg := fmt.Sprintf("%s is deprecated", v.Name)
				if v.DeprecationNote != "" {
					msg += ": " + v.DeprecationNote
				}
				result.Deprecated = append(result.Deprecated, Issue{
					Key:      v.Name,
					Severity: "warning",
					Message:  msg,
				})
			}
		}
	}

	return result
}
