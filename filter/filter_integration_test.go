package filter_test

import (
	"testing"

	"github.com/user/envlint/envparser"
	"github.com/user/envlint/filter"
)

func TestFilter_WithExampleEnv(t *testing.T) {
	env, err := envparser.Parse("../testdata/example.env")
	if err != nil {
		t.Fatalf("parse example.env: %v", err)
	}

	res, err := filter.Filter(env, filter.Options{
		Pattern: `.+`,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Kept) != len(env) {
		t.Errorf("expected all %d keys kept, got %d", len(env), len(res.Kept))
	}
	if len(res.Dropped) != 0 {
		t.Errorf("expected 0 dropped, got %d", len(res.Dropped))
	}
}

func TestFilter_ExcludeAllFromExampleEnv(t *testing.T) {
	env, err := envparser.Parse("../testdata/example.env")
	if err != nil {
		t.Fatalf("parse example.env: %v", err)
	}

	res, err := filter.Filter(env, filter.Options{
		Exclude: `.+`, // exclude everything
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Kept) != 0 {
		t.Errorf("expected 0 kept, got %d", len(res.Kept))
	}
	if len(res.Dropped) != len(env) {
		t.Errorf("expected all %d dropped, got %d", len(env), len(res.Dropped))
	}
}
