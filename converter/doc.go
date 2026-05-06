// Package converter provides utilities to transform .env variable maps
// into multiple output formats including JSON, YAML, TOML, shell export
// scripts, and standard dotenv notation.
//
// Supported formats:
//
//   - json   — JSON object with string key/value pairs
//   - yaml   — Simple YAML key: "value" pairs
//   - toml   — TOML key = "value" pairs
//   - shell  — Shell export statements (export KEY="value")
//   - dotenv — Standard KEY=value dotenv format
//
// Example usage:
//
//	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
//	out, err := converter.Convert(env, converter.FormatJSON)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(out)
package converter
