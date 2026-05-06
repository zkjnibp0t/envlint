// Package merger provides utilities for merging multiple .env files
// with configurable precedence rules.
package merger

import "fmt"

// Strategy defines how conflicts are resolved when merging.
type Strategy int

const (
	// FirstWins keeps the value from the first file that defines a key.
	FirstWins Strategy = iota
	// LastWins keeps the value from the last file that defines a key.
	LastWins
)

// Source represents a named env map with an origin label.
type Source struct {
	Name string
	Env  map[string]string
}

// Result holds the merged environment and metadata about conflicts.
type Result struct {
	Env       map[string]string
	Conflicts []Conflict
}

// Conflict records a key that appeared in more than one source.
type Conflict struct {
	Key      string
	Winner   string
	Losers   []string
	FinalVal string
}

// Merge combines multiple Sources into a single Result using the given Strategy.
func Merge(sources []Source, strategy Strategy) (Result, error) {
	if len(sources) == 0 {
		return Result{}, fmt.Errorf("merger: no sources provided")
	}

	merged := make(map[string]string)
	// track which source last set each key
	origin := make(map[string]string)
	conflictMap := make(map[string]*Conflict)

	for _, src := range sources {
		for k, v := range src.Env {
			if existing, exists := merged[k]; exists {
				c, ok := conflictMap[k]
				if !ok {
					c = &Conflict{Key: k}
					conflictMap[k] = c
				}
				switch strategy {
				case FirstWins:
					c.Winner = origin[k]
					c.Losers = append(c.Losers, src.Name)
					c.FinalVal = existing
				case LastWins:
					c.Losers = append(c.Losers, origin[k])
					c.Winner = src.Name
					c.FinalVal = v
					merged[k] = v
					origin[k] = src.Name
				}
			} else {
				merged[k] = v
				origin[k] = src.Name
			}
		}
	}

	var conflicts []Conflict
	for _, c := range conflictMap {
		conflicts = append(conflicts, *c)
	}

	return Result{Env: merged, Conflicts: conflicts}, nil
}
