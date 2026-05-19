// Package filter selects a subset of environment variables from a map based on
// configurable criteria such as key prefixes, inclusion patterns, and exclusion
// patterns.
//
// Example usage:
//
//	result, err := filter.Filter(env, filter.Options{
//		Prefixes: []string{"DB_", "REDIS_"},
//		Exclude:  ".*_TEST$",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	filter.WriteReport(os.Stdout, result)
//
// The original map is never mutated; Filter always returns a new map.
package filter
