// Package interpolator provides variable interpolation for .env files.
//
// It resolves ${VAR_NAME} style references within env values, enabling
// composed values such as:
//
//	BASE_URL=http://localhost:8080
//	API_URL=${BASE_URL}/api/v1
//
// Resolution order:
//  1. Other keys present in the parsed env map.
//  2. OS environment variables (via os.LookupEnv).
//
// If a reference cannot be resolved from either source, an ErrUnresolved
// error is collected and returned alongside the partially expanded map.
// The original placeholder (e.g. ${MISSING}) is preserved in the value
// so downstream validators can still inspect the raw reference.
package interpolator
