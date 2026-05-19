package filter_test

import (
	"strings"
	"testing"

	"github.com/user/envlint/filter"
)

var sampleEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"REDIS_URL":   "redis://localhost",
	"APP_NAME":    "envlint",
	"APP_VERSION": "1.0.0",
	"SECRET_KEY":  "s3cr3t",
	"DB_TEST":     "true",
}

func TestFilter_ByPrefix(t *testing.T) {
	res, err := filter.Filter(sampleEnv, filter.Options{Prefixes: []string{"DB_"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be kept")
	}
	if _, ok := res.Env["APP_NAME"]; ok {
		t.Error("expected APP_NAME to be dropped")
	}
	if len(res.Kept)+len(res.Dropped) != len(sampleEnv) {
		t.Error("kept+dropped should equal total keys")
	}
}

func TestFilter_ByPattern(t *testing.T) {
	res, err := filter.Filter(sampleEnv, filter.Options{Pattern: `^APP_`})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Kept) != 2 {
		t.Errorf("expected 2 kept, got %d", len(res.Kept))
	}
}

func TestFilter_Exclude(t *testing.T) {
	res, err := filter.Filter(sampleEnv, filter.Options{
		Prefixes: []string{"DB_"},
		Exclude:  `_TEST$`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.Env["DB_TEST"]; ok {
		t.Error("expected DB_TEST to be excluded")
	}
	if _, ok := res.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to remain")
	}
}

func TestFilter_CaseInsensitivePrefix(t *testing.T) {
	res, err := filter.Filter(sampleEnv, filter.Options{
		Prefixes:      []string{"db_"},
		CaseSensitive: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to match case-insensitive prefix db_")
	}
}

func TestFilter_NoOptions_KeepsAll(t *testing.T) {
	res, err := filter.Filter(sampleEnv, filter.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Kept) != len(sampleEnv) {
		t.Errorf("expected all %d keys kept, got %d", len(sampleEnv), len(res.Kept))
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := filter.Filter(sampleEnv, filter.Options{Pattern: `[invalid`})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestWriteReport_Output(t *testing.T) {
	res, _ := filter.Filter(sampleEnv, filter.Options{Prefixes: []string{"APP_"}})
	var sb strings.Builder
	filter.WriteReport(&sb, res)
	out := sb.String()
	if !strings.Contains(out, "Kept") {
		t.Error("expected 'Kept' section in report")
	}
	if !strings.Contains(out, "Dropped") {
		t.Error("expected 'Dropped' section in report")
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Error("expected APP_NAME in report")
	}
}
