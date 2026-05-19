// Package classifier categorises .env variables by their inferred purpose
// based on key naming conventions and value patterns.
package classifier

import (
	"regexp"
	"strings"
)

// Category represents a broad classification for an env variable.
type Category string

const (
	CategoryDatabase    Category = "database"
	CategoryAuth        Category = "auth"
	CategoryNetwork     Category = "network"
	CategoryFeatureFlag Category = "feature_flag"
	CategoryLogging     Category = "logging"
	CategoryStorage     Category = "storage"
	CategoryUnknown     Category = "unknown"
)

// Result holds the classification for a single env variable.
type Result struct {
	Key      string
	Value    string
	Category Category
	Reason   string
}

var (
	dbKeys      = regexp.MustCompile(`(?i)(DB_|DATABASE_|POSTGRES|MYSQL|MONGO|REDIS|DSN|SQL)`)
	authKeys    = regexp.MustCompile(`(?i)(AUTH|TOKEN|SECRET|PASSWORD|PASSWD|API_KEY|OAUTH|JWT)`)
	netKeys     = regexp.MustCompile(`(?i)(HOST|PORT|URL|ADDR|ENDPOINT|BASE_URL|PROXY)`)
	flagKeys    = regexp.MustCompile(`(?i)(FEATURE_|FLAG_|ENABLE_|DISABLE_|FF_)`)
	logKeys     = regexp.MustCompile(`(?i)(LOG|LOGGING|LOG_LEVEL|DEBUG|VERBOSE|TRACE)`)
	storageKeys = regexp.MustCompile(`(?i)(S3_|BUCKET|STORAGE|DISK|VOLUME|PATH|DIR)`)
)

// Classify assigns a Category to each key-value pair in the provided env map.
func Classify(env map[string]string) []Result {
	results := make([]Result, 0, len(env))
	for k, v := range env {
		results = append(results, classifyOne(k, v))
	}
	return results
}

func classifyOne(key, value string) Result {
	upper := strings.ToUpper(key)
	switch {
	case dbKeys.MatchString(upper):
		return Result{Key: key, Value: value, Category: CategoryDatabase, Reason: "key matches database pattern"}
	case authKeys.MatchString(upper):
		return Result{Key: key, Value: value, Category: CategoryAuth, Reason: "key matches auth/secret pattern"}
	case netKeys.MatchString(upper):
		return Result{Key: key, Value: value, Category: CategoryNetwork, Reason: "key matches network pattern"}
	case flagKeys.MatchString(upper):
		return Result{Key: key, Value: value, Category: CategoryFeatureFlag, Reason: "key matches feature-flag pattern"}
	case logKeys.MatchString(upper):
		return Result{Key: key, Value: value, Category: CategoryLogging, Reason: "key matches logging pattern"}
	case storageKeys.MatchString(upper):
		return Result{Key: key, Value: value, Category: CategoryStorage, Reason: "key matches storage pattern"}
	default:
		return Result{Key: key, Value: value, Category: CategoryUnknown, Reason: "no pattern matched"}
	}
}

// GroupByCategory returns a map of Category -> slice of Results.
func GroupByCategory(results []Result) map[Category][]Result {
	groups := make(map[Category][]Result)
	for _, r := range results {
		groups[r.Category] = append(groups[r.Category], r)
	}
	return groups
}
