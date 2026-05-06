// Package linter provides a unified lint pipeline for .env files.
//
// It orchestrates the validator, formatter, auditor, and suggester packages
// into a single Run call that returns a consolidated Result.
//
// Usage:
//
//	env, _ := envparser.Parse("path/to/.env")
//	s, _   := schema.Load("path/to/schema.yaml")
//
//	result := linter.Run(env, s)
//	linter.WriteReport(os.Stdout, result)
//	if result.HasErrors {
//		os.Exit(1)
//	}
//
// The Result.HasErrors flag is true when there are validation errors or
// auditor issues at the LevelError severity, making it suitable for use
// in CI pipelines as an exit-code gate.
package linter
