// Package suggester provides suggestions for fixing validation errors
// found in .env files, such as typos in variable names or format hints.
package suggester

import (
	"fmt"
	"strings"

	"github.com/user/envlint/schema"
	"github.com/user/envlint/validator"
)

// Suggestion holds a human-readable fix hint for a validation error.
type Suggestion struct {
	Variable string
	Message  string
}

// Suggest returns a list of suggestions based on validation errors and the schema.
func Suggest(errs []validator.ValidationError, s *schema.Schema, envKeys []string) []Suggestion {
	var suggestions []Suggestion

	for _, e := range errs {
		switch e.Kind {
		case validator.ErrMissing:
			if close := closestMatch(e.Variable, envKeys); close != "" {
				suggestions = append(suggestions, Suggestion{
					Variable: e.Variable,
					Message:  fmt.Sprintf("did you mean '%s'?", close),
				})
			} else {
				suggestions = append(suggestions, Suggestion{
					Variable: e.Variable,
					Message:  fmt.Sprintf("add %s to your .env file", e.Variable),
				})
			}
		case validator.ErrInvalidType:
			if v, ok := s.Vars[e.Variable]; ok {
				suggestions = append(suggestions, Suggestion{
					Variable: e.Variable,
					Message:  fmt.Sprintf("expected type '%s' for %s", v.Type, e.Variable),
				})
			}
		}
	}

	return suggestions
}

// closestMatch finds the closest key in candidates to target using a simple
// case-insensitive prefix or substring heuristic.
func closestMatch(target string, candidates []string) string {
	lower := strings.ToLower(target)
	for _, c := range candidates {
		if strings.ToLower(c) == lower {
			return c
		}
	}
	for _, c := range candidates {
		if strings.HasPrefix(strings.ToLower(c), lower[:max(1, len(lower)-2)]) {
			return c
		}
	}
	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
