package sorter_test

import (
	"testing"

	"github.com/user/envlint/sorter"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"APP_PORT":    "8080",
		"APP_ENV":     "production",
		"SECRET_KEY":  "abc123",
		"DB_PASSWORD": "s3cr3t",
		"NOUNDERSCORE": "plain",
	}
}

func TestSort_KeysAreAlphabetical(t *testing.T) {
	result := sorter.Sort(sampleEnv())
	keys := sorter.SortedKeys(result.Sorted)

	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys out of order: %q before %q", keys[i-1], keys[i])
		}
	}
}

func TestSort_PreservesValues(t *testing.T) {
	env := sampleEnv()
	result := sorter.Sort(env)

	for k, v := range env {
		if result.Sorted[k] != v {
			t.Errorf("key %q: want %q, got %q", k, v, result.Sorted[k])
		}
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	env := sampleEnv()
	orig := make(map[string]string, len(env))
	for k, v := range env {
		orig[k] = v
	}

	sorter.Sort(env)

	for k, v := range orig {
		if env[k] != v {
			t.Errorf("original mutated at key %q", k)
		}
	}
}

func TestSort_GroupsByPrefix(t *testing.T) {
	result := sorter.Sort(sampleEnv())

	dbKeys, ok := result.Groups["DB"]
	if !ok {
		t.Fatal("expected group DB")
	}
	if len(dbKeys) != 2 {
		t.Errorf("DB group: want 2 keys, got %d", len(dbKeys))
	}

	appKeys := result.Groups["APP"]
	if len(appKeys) != 2 {
		t.Errorf("APP group: want 2 keys, got %d", len(appKeys))
	}
}

func TestSort_NoUnderscoreKeyInFallbackGroup(t *testing.T) {
	result := sorter.Sort(sampleEnv())

	fallback, ok := result.Groups["_"]
	if !ok {
		t.Fatal("expected fallback group _")
	}
	if len(fallback) != 1 || fallback[0] != "NOUNDERSCORE" {
		t.Errorf("unexpected fallback group contents: %v", fallback)
	}
}

func TestSort_EmptyEnv(t *testing.T) {
	result := sorter.Sort(map[string]string{})
	if len(result.Sorted) != 0 {
		t.Error("expected empty sorted map")
	}
	if len(result.Groups) != 0 {
		t.Error("expected empty groups")
	}
}
