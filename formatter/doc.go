// Package formatter provides style-checking utilities for .env files.
//
// It inspects raw .env lines for common convention violations such as:
//   - Lowercase key names (keys should be UPPER_SNAKE_CASE)
//   - Spaces surrounding the '=' assignment operator
//   - Trailing whitespace in values
//
// Usage:
//
//	lines := []string{"db_host=localhost", "PORT = 8080"}
//	issues := formatter.CheckStyle(lines, nil)
//	formatter.WriteIssues(os.Stdout, issues)
//	fmt.Println(formatter.Summary(issues))
//
// Style issues are categorised by Level: error, warning, or info.
// These checks are independent of schema validation and run as a
// complementary linting pass.
package formatter
