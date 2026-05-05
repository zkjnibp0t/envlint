package validator

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/user/envlint/schema"
)

// Severity indicates how serious a validation issue is.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// ValidationError holds details about a single failed validation.
type ValidationError struct {
	Key      string
	Message  string
	Severity Severity
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

// Validate checks env values against the schema and returns all errors found.
func Validate(s *schema.Schema, env map[string]string) []ValidationError {
	var errs []ValidationError

	for _, v := range s.Vars {
		val, present := env[v.Name]

		if !present || val == "" {
			if v.Required {
				errs = append(errs, ValidationError{
					Key:      v.Name,
					Message:  "is required but missing or empty",
					Severity: SeverityError,
				})
			}
			continue
		}

		if typeErr := validateValue(v.Name, val, v.Type); typeErr != nil {
			errs = append(errs, *typeErr)
		}

		if v.Pattern != "" {
			if matched, _ := regexp.MatchString("^"+v.Pattern+"$", val); !matched {
				errs = append(errs, ValidationError{
					Key:      v.Name,
					Message:  fmt.Sprintf("does not match pattern %q", v.Pattern),
					Severity: SeverityError,
				})
			}
		}
	}

	return errs
}

func validateValue(key, val, typ string) *ValidationError {
	switch typ {
	case "int":
		if _, err := strconv.Atoi(val); err != nil {
			return &ValidationError{
				Key:      key,
				Message:  fmt.Sprintf("expected int, got %q", val),
				Severity: SeverityError,
			}
		}
	case "bool":
		if _, err := strconv.ParseBool(val); err != nil {
			return &ValidationError{
				Key:      key,
				Message:  fmt.Sprintf("expected bool, got %q", val),
				Severity: SeverityError,
			}
		}
	case "url":
		u, err := url.ParseRequestURI(val)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return &ValidationError{
				Key:      key,
				Message:  fmt.Sprintf("expected valid URL, got %q", val),
				Severity: SeverityError,
			}
		}
	}
	return nil
}
