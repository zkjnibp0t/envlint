package trimmer_test

import (
	"testing"

	"github.com/yourorg/envlint/trimmer"
)

func TestTrim_BasicClean(t *testing.T) {
	lines := []string{
		"APP_NAME=envlint",
		"PORT=8080",
	}
	res := trimmer.Trim(lines)
	if res.Env["APP_NAME"] != "envlint" {
		t.Errorf("expected envlint, got %q", res.Env["APP_NAME"])
	}
	if res.Env["PORT"] != "8080" {
		t.Errorf("expected 8080, got %q", res.Env["PORT"])
	}
}

func TestTrim_RemovesBlanks(t *testing.T) {
	lines := []string{
		"APP=foo",
		"",
		"   ",
		"DB=bar",
	}
	res := trimmer.Trim(lines)
	if res.BlankLinesRemoved != 2 {
		t.Errorf("expected 2 blank lines removed, got %d", res.BlankLinesRemoved)
	}
}

func TestTrim_CollatesDuplicates(t *testing.T) {
	lines := []string{
		"KEY=first",
		"KEY=second",
		"KEY=third",
	}
	res := trimmer.Trim(lines)
	if res.Env["KEY"] != "third" {
		t.Errorf("expected last value 'third', got %q", res.Env["KEY"])
	}
	if len(res.RemovedDuplicates) != 1 || res.RemovedDuplicates[0] != "KEY" {
		t.Errorf("expected KEY in duplicates, got %v", res.RemovedDuplicates)
	}
}

func TestTrim_TrimsValueWhitespace(t *testing.T) {
	lines := []string{
		"HOST=  localhost  ",
		"PORT=   9090",
	}
	res := trimmer.Trim(lines)
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", res.Env["HOST"])
	}
	if len(res.TrimmedValues) < 1 {
		t.Errorf("expected at least one trimmed value, got %v", res.TrimmedValues)
	}
}

func TestTrim_SkipsComments(t *testing.T) {
	lines := []string{
		"# this is a comment",
		"APP=hello",
		"# another comment",
	}
	res := trimmer.Trim(lines)
	if len(res.Env) != 1 {
		t.Errorf("expected 1 key, got %d", len(res.Env))
	}
	if res.BlankLinesRemoved != 0 {
		t.Errorf("comments should not count as blank lines")
	}
}

func TestTrim_SkipsInvalidLines(t *testing.T) {
	lines := []string{
		"NOTAKEYVALUE",
		"VALID=yes",
	}
	res := trimmer.Trim(lines)
	if _, ok := res.Env["NOTAKEYVALUE"]; ok {
		t.Error("invalid line should not appear in env map")
	}
	if res.Env["VALID"] != "yes" {
		t.Errorf("expected 'yes', got %q", res.Env["VALID"])
	}
}
