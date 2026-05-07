// Package snapshoter captures and restores .env state snapshots,
// useful for comparing environments across deployments or CI runs.
package snapshoter

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// Snapshot represents a captured state of an environment.
type Snapshot struct {
	CapturedAt time.Time         `json:"captured_at"`
	Label      string            `json:"label"`
	Env        map[string]string `json:"env"`
}

// Capture creates a new Snapshot from the given env map.
func Capture(env map[string]string, label string) Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return Snapshot{
		CapturedAt: time.Now().UTC(),
		Label:      label,
		Env:        copy,
	}
}

// Save writes a snapshot to a JSON file at the given path.
func Save(s Snapshot, path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshoter: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshoter: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshoter: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshoter: unmarshal: %w", err)
	}
	return s, nil
}

// Keys returns the sorted list of keys in the snapshot.
func (s Snapshot) Keys() []string {
	keys := make([]string, 0, len(s.Env))
	for k := range s.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
