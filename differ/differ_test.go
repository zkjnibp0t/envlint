package differ_test

import (
	"testing"

	"github.com/user/envlint/differ"
)

func TestCompare_Identical(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	d := differ.Compare(a, b)
	if !d.IsClean() {
		t.Errorf("expected clean diff, got %+v", d)
	}
}

func TestCompare_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "1", "EXTRA": "yes"}
	b := map[string]string{"FOO": "1"}
	d := differ.Compare(a, b)
	if len(d.OnlyInA) != 1 || d.OnlyInA[0] != "EXTRA" {
		t.Errorf("expected EXTRA in OnlyInA, got %v", d.OnlyInA)
	}
	if len(d.OnlyInB) != 0 || len(d.Changed) != 0 {
		t.Errorf("unexpected diff entries: %+v", d)
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"FOO": "1", "NEW_KEY": "hello"}
	d := differ.Compare(a, b)
	if len(d.OnlyInB) != 1 || d.OnlyInB[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY in OnlyInB, got %v", d.OnlyInB)
	}
}

func TestCompare_Changed(t *testing.T) {
	a := map[string]string{"PORT": "8080"}
	b := map[string]string{"PORT": "9090"}
	d := differ.Compare(a, b)
	if len(d.Changed) != 1 || d.Changed[0] != "PORT" {
		t.Errorf("expected PORT in Changed, got %v", d.Changed)
	}
	if len(d.OnlyInA) != 0 || len(d.OnlyInB) != 0 {
		t.Errorf("unexpected diff entries: %+v", d)
	}
}

func TestCompare_Mixed(t *testing.T) {
	a := map[string]string{"A": "1", "B": "old", "C": "same"}
	b := map[string]string{"B": "new", "C": "same", "D": "4"}
	d := differ.Compare(a, b)
	if len(d.OnlyInA) != 1 || d.OnlyInA[0] != "A" {
		t.Errorf("OnlyInA wrong: %v", d.OnlyInA)
	}
	if len(d.OnlyInB) != 1 || d.OnlyInB[0] != "D" {
		t.Errorf("OnlyInB wrong: %v", d.OnlyInB)
	}
	if len(d.Changed) != 1 || d.Changed[0] != "B" {
		t.Errorf("Changed wrong: %v", d.Changed)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	d := differ.Compare(map[string]string{}, map[string]string{})
	if !d.IsClean() {
		t.Errorf("expected clean diff for two empty maps")
	}
}
