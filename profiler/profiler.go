// Package profiler analyzes .env files and produces a summary report
// of variable counts, types, and sensitive key statistics.
package profiler

import (
	"fmt"
	"sort"

	"github.com/user/envlint/masker"
	"github.com/user/envlint/schema"
)

// Profile holds aggregated statistics about an env file.
type Profile struct {
	TotalVars     int
	SensitiveVars int
	TypeBreakdown map[string]int
	DefinedInSchema int
	UndefinedInSchema int
	TopKeys       []string
}

// Analyze builds a Profile from a parsed env map and an optional schema.
// If s is nil, schema-related fields are left at zero.
func Analyze(env map[string]string, s *schema.Schema) Profile {
	p := Profile{
		TypeBreakdown: make(map[string]int),
	}

	p.TotalVars = len(env)

	schemaVars := map[string]schema.VarDef{}
	if s != nil {
		for _, v := range s.Vars {
			schemaVars[v.Name] = v
		}
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if masker.IsSensitive(k) {
			p.SensitiveVars++
		}

		if def, ok := schemaVars[k]; ok {
			p.DefinedInSchema++
			t := def.Type
			if t == "" {
				t = "string"
			}
			p.TypeBreakdown[t]++
		} else {
			p.UndefinedInSchema++
			p.TypeBreakdown["unknown"]++
		}
	}

	max := 10
	if len(keys) < max {
		max = len(keys)
	}
	p.TopKeys = keys[:max]

	return p
}

// Summary returns a human-readable summary string of the profile.
func Summary(p Profile) string {
	return fmt.Sprintf(
		"Total: %d | Sensitive: %d | In schema: %d | Undefined: %d",
		p.TotalVars, p.SensitiveVars, p.DefinedInSchema, p.UndefinedInSchema,
	)
}

// SensitiveRatio returns the fraction of variables considered sensitive,
// as a value between 0.0 and 1.0. Returns 0 if there are no variables.
func SensitiveRatio(p Profile) float64 {
	if p.TotalVars == 0 {
		return 0
	}
	return float64(p.SensitiveVars) / float64(p.TotalVars)
}
