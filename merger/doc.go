// Package merger combines multiple .env files into a single environment map.
//
// It supports two merge strategies:
//
//   - FirstWins: the first source to define a key takes precedence.
//   - LastWins:  the last source to define a key takes precedence (typical
//     shell override behaviour).
//
// Conflicts — keys present in more than one source — are recorded in the
// Result so callers can surface them to the user via WriteReport.
//
// Example:
//
//	sources := []merger.Source{
//		{Name: ".env",       Env: base},
//		{Name: ".env.local", Env: local},
//	}
//	result, err := merger.Merge(sources, merger.LastWins)
package merger
