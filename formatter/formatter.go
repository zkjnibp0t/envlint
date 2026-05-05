package formatter

import (
	"fmt"
	"strings"

	"github.com/user/envlint/validator"
)

// Level represents the severity of a formatting style.
type Level string

const (
	LevelError   Level = "error"
	LevelWarning Level = "warning"
	LevelInfo    Level = "info"
)

// Issue represents a style or convention issue found in a .env file.
type Issue struct {
	Line    int
	Key     string
	Message string
	Level   Level
}

// CheckStyle inspects raw env lines for common style violations.
// It returns a slice of Issues alongside any validation errors.
func CheckStyle(lines []string, errs []validator.ValidationError) []Issue {
	var issues []Issue

	for i, line := range lines {
		lineNum := i + 1

		// Skip blank lines and comments
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		val := parts[1]

		// Check: key should be UPPER_SNAKE_CASE
		if key != strings.ToUpper(key) {
			issues = append(issues, Issue{
				Line:    lineNum,
				Key:     key,
				Message: fmt.Sprintf("key %q should be uppercase", key),
				Level:   LevelWarning,
			})
		}

		// Check: no spaces around the '=' sign
		if strings.Contains(line, " = ") {
			issues = append(issues, Issue{
				Line:    lineNum,
				Key:     key,
				Message: fmt.Sprintf("key %q has spaces around '='", key),
				Level:   LevelWarning,
			})
		}

		// Check: value should not have trailing whitespace
		if val != strings.TrimRight(val, " \t") {
			issues = append(issues, Issue{
				Line:    lineNum,
				Key:     key,
				Message: fmt.Sprintf("key %q has trailing whitespace in value", key),
				Level:   LevelInfo,
			})
		}
	}

	return issues
}
