// Package differ compares two .env files and reports keys that are
// present in one but absent in the other, helping teams keep environments
// in sync (e.g. .env.example vs production .env).
package differ

import "sort"

// Diff holds the result of comparing two env maps.
type Diff struct {
	// OnlyInA are keys present in the first env but not the second.
	OnlyInA []string
	// OnlyInB are keys present in the second env but not the first.
	OnlyInB []string
	// Changed are keys present in both envs but with different values.
	Changed []string
}

// IsClean returns true when there are no differences.
func (d Diff) IsClean() bool {
	return len(d.OnlyInA) == 0 && len(d.OnlyInB) == 0 && len(d.Changed) == 0
}

// Compare computes the difference between two env maps a and b.
// Values are compared as plain strings; use masker.MaskEnv before
// passing sensitive maps if you do not want raw secrets in the result.
func Compare(a, b map[string]string) Diff {
	d := Diff{}

	for k, va := range a {
		vb, ok := b[k]
		if !ok {
			d.OnlyInA = append(d.OnlyInA, k)
		} else if va != vb {
			d.Changed = append(d.Changed, k)
		}
	}

	for k := range b {
		if _, ok := a[k]; !ok {
			d.OnlyInB = append(d.OnlyInB, k)
		}
	}

	sort.Strings(d.OnlyInA)
	sort.Strings(d.OnlyInB)
	sort.Strings(d.Changed)

	return d
}
