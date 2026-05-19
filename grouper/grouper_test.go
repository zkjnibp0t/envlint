package grouper_test

import (
	"bytes"
	"strings"
	"testing"

	"envlint/grouper"
)

var sampleEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"AWS_KEY":     "AKIA123",
	"AWS_SECRET":  "secret",
	"APP_ENV":     "production",
	"LOG_LEVEL":   "info",
	"UNRELATED":   "value",
}

func TestByPrefix_GroupsCorrectly(t *testing.T) {
	r := grouper.ByPrefix(sampleEnv, []string{"DB_", "AWS_"})
	if len(r.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(r.Groups))
	}
	if _, ok := r.Ungrouped["APP_ENV"]; !ok {
		t.Error("expected APP_ENV in ungrouped")
	}
	if _, ok := r.Ungrouped["UNRELATED"]; !ok {
		t.Error("expected UNRELATED in ungrouped")
	}
}

func TestByPrefix_NoPrefixes(t *testing.T) {
	r := grouper.ByPrefix(sampleEnv, nil)
	if len(r.Groups) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(r.Groups))
	}
	if len(r.Ungrouped) != len(sampleEnv) {
		t.Errorf("expected all keys ungrouped, got %d", len(r.Ungrouped))
	}
}

func TestByPrefix_EmptyEnv(t *testing.T) {
	r := grouper.ByPrefix(map[string]string{}, []string{"DB_"})
	if len(r.Groups) != 0 {
		t.Error("expected no groups for empty env")
	}
	if len(r.Ungrouped) != 0 {
		t.Error("expected no ungrouped for empty env")
	}
}

func TestByPrefix_GroupsAreSorted(t *testing.T) {
	r := grouper.ByPrefix(sampleEnv, []string{"LOG_", "APP_", "DB_"})
	previous := ""
	for _, g := range r.Groups {
		if g.Prefix < previous {
			t.Errorf("groups not sorted: %s before %s", previous, g.Prefix)
		}
		previous = g.Prefix
	}
}

func TestWriteReport_ContainsGroups(t *testing.T) {
	r := grouper.ByPrefix(sampleEnv, []string{"DB_", "AWS_"})
	var buf bytes.Buffer
	grouper.WriteReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "[AWS_]") {
		t.Error("expected AWS_ group in report")
	}
	if !strings.Contains(out, "[DB_]") {
		t.Error("expected DB_ group in report")
	}
	if !strings.Contains(out, "[UNGROUPED]") {
		t.Error("expected UNGROUPED section in report")
	}
	if !strings.Contains(out, "total:") {
		t.Error("expected total line in report")
	}
}

func TestWriteReport_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	grouper.WriteReport(&buf, grouper.Result{})
	if !strings.Contains(buf.String(), "no variables") {
		t.Error("expected empty message for empty result")
	}
}
