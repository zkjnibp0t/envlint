package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envlint/schema"
)

// Result holds the outcome of a single variable validation.
type Result struct {
	Key     string
	Passed  bool
	Message string
}

// Report aggregates all validation results.
type Report struct {
	Results []Result
}

// HasErrors returns true if any result failed.
func (r *Report) HasErrors() bool {
	for _, res := range r.Results {
		if !res.Passed {
			return true
		}
	}
	return false
}

// Validate checks the provided env map against the schema.
func Validate(s *schema.Schema, env map[string]string) *Report {
	report := &Report{}

	for _, v := range s.Vars {
		val, exists := env[v.Name]

		if !exists || val == "" {
			if v.Required {
				report.Results = append(report.Results, Result{
					Key:     v.Name,
					Passed:  false,
					Message: fmt.Sprintf("%s is required but missing or empty", v.Name),
				})
			}
			continue
		}

		if err := validateValue(v, val); err != nil {
			report.Results = append(report.Results, Result{
				Key:     v.Name,
				Passed:  false,
				Message: err.Error(),
			})
		} else {
			report.Results = append(report.Results, Result{
				Key:    v.Name,
				Passed: true,
			})
		}
	}

	return report
}

func validateValue(v schema.VarDef, val string) error {
	switch strings.ToLower(v.Type) {
	case "url":
		if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
			return fmt.Errorf("%s must be a valid URL (got %q)", v.Name, val)
		}
	case "int":
		if matched, _ := regexp.MatchString(`^-?\d+$`, val); !matched {
			return fmt.Errorf("%s must be an integer (got %q)", v.Name, val)
		}
	case "bool":
		lower := strings.ToLower(val)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return fmt.Errorf("%s must be a boolean (got %q)", v.Name, val)
		}
	}

	if v.Pattern != "" {
		matched, err := regexp.MatchString(v.Pattern, val)
		if err != nil {
			return fmt.Errorf("%s has invalid pattern %q: %w", v.Name, v.Pattern, err)
		}
		if !matched {
			return fmt.Errorf("%s does not match pattern %q (got %q)", v.Name, v.Pattern, val)
		}
	}

	return nil
}
