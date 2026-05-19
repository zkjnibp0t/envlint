package classifier_test

import (
	"testing"

	"github.com/envlint/envlint/classifier"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"DB_HOST":          "localhost",
		"DATABASE_URL":     "postgres://localhost/mydb",
		"JWT_SECRET":       "supersecret",
		"API_KEY":          "abc123",
		"SERVER_PORT":      "8080",
		"BASE_URL":         "https://example.com",
		"FEATURE_DARK_MODE": "true",
		"LOG_LEVEL":        "info",
		"S3_BUCKET":        "my-bucket",
		"APP_NAME":         "envlint",
	}
}

func TestClassify_DatabaseKeys(t *testing.T) {
	results := classifier.Classify(map[string]string{
		"DB_HOST":      "localhost",
		"DATABASE_URL": "postgres://localhost/db",
	})
	for _, r := range results {
		if r.Category != classifier.CategoryDatabase {
			t.Errorf("expected database for %s, got %s", r.Key, r.Category)
		}
	}
}

func TestClassify_AuthKeys(t *testing.T) {
	results := classifier.Classify(map[string]string{
		"JWT_SECRET": "s3cr3t",
		"API_KEY":    "key123",
	})
	for _, r := range results {
		if r.Category != classifier.CategoryAuth {
			t.Errorf("expected auth for %s, got %s", r.Key, r.Category)
		}
	}
}

func TestClassify_NetworkKeys(t *testing.T) {
	results := classifier.Classify(map[string]string{
		"SERVER_PORT": "8080",
		"BASE_URL":    "https://example.com",
	})
	for _, r := range results {
		if r.Category != classifier.CategoryNetwork {
			t.Errorf("expected network for %s, got %s", r.Key, r.Category)
		}
	}
}

func TestClassify_FeatureFlagKeys(t *testing.T) {
	results := classifier.Classify(map[string]string{
		"FEATURE_DARK_MODE": "true",
		"ENABLE_SIGNUP":     "false",
	})
	for _, r := range results {
		if r.Category != classifier.CategoryFeatureFlag {
			t.Errorf("expected feature_flag for %s, got %s", r.Key, r.Category)
		}
	}
}

func TestClassify_UnknownKey(t *testing.T) {
	results := classifier.Classify(map[string]string{"APP_NAME": "envlint"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Category != classifier.CategoryUnknown {
		t.Errorf("expected unknown, got %s", results[0].Category)
	}
}

func TestGroupByCategory_CorrectBuckets(t *testing.T) {
	results := classifier.Classify(sampleEnv())
	groups := classifier.GroupByCategory(results)

	if len(groups[classifier.CategoryDatabase]) == 0 {
		t.Error("expected database group to be non-empty")
	}
	if len(groups[classifier.CategoryAuth]) == 0 {
		t.Error("expected auth group to be non-empty")
	}
	if len(groups[classifier.CategoryNetwork]) == 0 {
		t.Error("expected network group to be non-empty")
	}
	if len(groups[classifier.CategoryUnknown]) == 0 {
		t.Error("expected unknown group to contain APP_NAME")
	}
}

func TestGroupByCategory_TotalCount(t *testing.T) {
	env := sampleEnv()
	results := classifier.Classify(env)
	groups := classifier.GroupByCategory(results)

	total := 0
	for _, v := range groups {
		total += len(v)
	}
	if total != len(env) {
		t.Errorf("expected total %d, got %d", len(env), total)
	}
}
