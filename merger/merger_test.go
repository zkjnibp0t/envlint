package merger

import (
	"testing"
)

func TestMerge_NoSources(t *testing.T) {
	_, err := Merge(nil, FirstWins)
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestMerge_SingleSource(t *testing.T) {
	src := Source{Name: "base", Env: map[string]string{"FOO": "bar", "BAZ": "qux"}}
	r, err := Merge([]Source{src}, FirstWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", r.Env["FOO"])
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(r.Conflicts))
	}
}

func TestMerge_FirstWins(t *testing.T) {
	sources := []Source{
		{Name: "base", Env: map[string]string{"KEY": "first"}},
		{Name: "override", Env: map[string]string{"KEY": "second"}},
	}
	r, err := Merge(sources, FirstWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["KEY"] != "first" {
		t.Errorf("FirstWins: expected 'first', got %s", r.Env["KEY"])
	}
	if len(r.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(r.Conflicts))
	}
	if r.Conflicts[0].Winner != "base" {
		t.Errorf("expected winner=base, got %s", r.Conflicts[0].Winner)
	}
}

func TestMerge_LastWins(t *testing.T) {
	sources := []Source{
		{Name: "base", Env: map[string]string{"KEY": "first"}},
		{Name: "override", Env: map[string]string{"KEY": "second"}},
	}
	r, err := Merge(sources, LastWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["KEY"] != "second" {
		t.Errorf("LastWins: expected 'second', got %s", r.Env["KEY"])
	}
	if r.Conflicts[0].Winner != "override" {
		t.Errorf("expected winner=override, got %s", r.Conflicts[0].Winner)
	}
}

func TestMerge_NoConflicts(t *testing.T) {
	sources := []Source{
		{Name: "a", Env: map[string]string{"A": "1"}},
		{Name: "b", Env: map[string]string{"B": "2"}},
	}
	r, err := Merge(sources, FirstWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(r.Conflicts))
	}
	if r.Env["A"] != "1" || r.Env["B"] != "2" {
		t.Errorf("unexpected merged values: %v", r.Env)
	}
}
