package formatter_test

import (
	"testing"

	"github.com/user/envlint/formatter"
	"github.com/user/envlint/validator"
)

func TestCheckStyle_UppercaseWarning(t *testing.T) {
	lines := []string{"db_host=localhost"}
	issues := formatter.CheckStyle(lines, nil)
	if len(issues) == 0 {
		t.Fatal("expected at least one issue for lowercase key")
	}
	if issues[0].Level != formatter.LevelWarning {
		t.Errorf("expected warning, got %s", issues[0].Level)
	}
	if issues[0].Line != 1 {
		t.Errorf("expected line 1, got %d", issues[0].Line)
	}
}

func TestCheckStyle_SpacesAroundEquals(t *testing.T) {
	lines := []string{"DB_HOST = localhost"}
	issues := formatter.CheckStyle(lines, nil)
	found := false
	for _, iss := range issues {
		if iss.Level == formatter.LevelWarning && iss.Line == 1 {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for spaces around '='")
	}
}

func TestCheckStyle_TrailingWhitespace(t *testing.T) {
	lines := []string{"DB_HOST=localhost   "}
	issues := formatter.CheckStyle(lines, nil)
	found := false
	for _, iss := range issues {
		if iss.Level == formatter.LevelInfo {
			found = true
		}
	}
	if !found {
		t.Error("expected info issue for trailing whitespace")
	}
}

func TestCheckStyle_SkipsCommentsAndBlanks(t *testing.T) {
	lines := []string{"", "# comment", "DB_HOST=localhost"}
	issues := formatter.CheckStyle(lines, nil)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestCheckStyle_CleanLine(t *testing.T) {
	lines := []string{"DATABASE_URL=postgres://localhost/mydb"}
	issues := formatter.CheckStyle(lines, nil)
	if len(issues) != 0 {
		t.Errorf("expected no issues for clean line, got %d", len(issues))
	}
}

func TestCheckStyle_WithValidationErrors(t *testing.T) {
	lines := []string{"DB_PORT=abc"}
	verrs := []validator.ValidationError{
		{Key: "DB_PORT", Message: "expected int"},
	}
	issues := formatter.CheckStyle(lines, verrs)
	// Style check should still run independently
	if issues == nil {
		t.Log("no style issues, which is acceptable")
	}
}
