// Package interpolator resolves variable references within .env values.
// It expands expressions like ${VAR_NAME} using already-parsed env values
// or OS environment variables as a fallback.
package interpolator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// ErrUnresolved is returned when a referenced variable cannot be found.
type ErrUnresolved struct {
	Key string
	Ref string
}

func (e *ErrUnresolved) Error() string {
	return fmt.Sprintf("variable %q references undefined variable %q", e.Key, e.Ref)
}

// Expand resolves all ${REF} expressions in the provided env map.
// It mutates and returns the map. References are resolved using other
// keys in the map first, then falling back to the OS environment.
// An ErrUnresolved slice is returned for any unresolvable references.
func Expand(env map[string]string) (map[string]string, []error) {
	var errs []error

	for key, value := range env {
		resolved, resolveErrs := expandValue(key, value, env)
		if len(resolveErrs) > 0 {
			errs = append(errs, resolveErrs...)
		}
		env[key] = resolved
	}

	return env, errs
}

func expandValue(key, value string, env map[string]string) (string, []error) {
	var errs []error

	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		ref := strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")

		if v, ok := env[ref]; ok {
			return v
		}
		if v, ok := os.LookupEnv(ref); ok {
			return v
		}

		errs = append(errs, &ErrUnresolved{Key: key, Ref: ref})
		return match
	})

	return result, errs
}
