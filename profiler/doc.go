// Package profiler provides statistical analysis of parsed .env files.
//
// It computes variable counts, identifies sensitive keys using the masker
// package, cross-references variables against a schema definition, and
// breaks down variables by their declared types.
//
// Example usage:
//
//	env, _ := envparser.Parse("app.env")
//	s, _   := schema.Load("app.schema.yaml")
//	p      := profiler.Analyze(env, s)
//	fmt.Println(profiler.Summary(p))
package profiler
