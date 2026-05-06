// Package masker identifies and redacts sensitive environment variable
// values to prevent accidental exposure of secrets in logs, CLI output,
// or error reports.
//
// A key is considered sensitive if its name (case-insensitive) contains
// any of the well-known patterns such as SECRET, PASSWORD, TOKEN, API_KEY,
// PRIVATE_KEY, CREDENTIALS, or AUTH.
//
// Usage:
//
//	// Mask a single value before printing
//	safe := masker.MaskValue("DB_PASSWORD", os.Getenv("DB_PASSWORD"))
//	fmt.Println("DB_PASSWORD =", safe)
//
//	// Mask an entire env map for safe display
//	safeEnv := masker.MaskEnv(parsedEnv)
//	for k, v := range safeEnv {
//		fmt.Printf("%s=%s\n", k, v)
//	}
package masker
