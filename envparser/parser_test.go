package envparser_test

import (
	"os"
	"testing"

	"github.com/user/envlint/envparser"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	env, err := envparser.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if env["APP_NAME"] != "myapp" || env["PORT"] != "8080" {
		t.Errorf("unexpected env: %v", env)
	}
}

func TestParse_CommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=value\n")
	env, err := envparser.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(env) != 1 || env["KEY"] != "value" {
		t.Errorf("unexpected env: %v", env)
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `DB_URL="postgres://localhost/db"` + "\n" + `TOKEN='abc123'` + "\n")
	env, err := envparser.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if env["DB_URL"] != "postgres://localhost/db" {
		t.Errorf("expected unquoted DB_URL, got %q", env["DB_URL"])
	}
	if env["TOKEN"] != "abc123" {
		t.Errorf("expected unquoted TOKEN, got %q", env["TOKEN"])
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := envparser.Parse(path)
	if err == nil {
		t.Error("expected error for invalid line")
	}
}

func TestParse_MissingFile(t *testing.T) {
	_, err := envparser.Parse("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestParse_ValueWithEquals(t *testing.T) {
	path := writeTempEnv(t, "CONN=host=localhost port=5432\n")
	env, err := envparser.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if env["CONN"] != "host=localhost port=5432" {
		t.Errorf("unexpected CONN value: %q", env["CONN"])
	}
}

func TestParse_EmptyValue(t *testing.T) {
	path := writeTempEnv(t, "EMPTY=\nALSO_EMPTY=\"\"\n")
	env, err := envparser.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if env["EMPTY"] != "" {
		t.Errorf("expected empty string for EMPTY, got %q", env["EMPTY"])
	}
	if env["ALSO_EMPTY"] != "" {
		t.Errorf("expected empty string for ALSO_EMPTY, got %q", env["ALSO_EMPTY"])
	}
}
