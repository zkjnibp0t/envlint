// Package snapshoter provides functionality to capture, persist, and reload
// snapshots of .env environments as JSON files.
//
// A snapshot records the full key-value state of an environment at a point in
// time along with a human-readable label and a UTC timestamp. Snapshots can be
// saved to disk and later loaded for comparison with differ.Compare or for
// audit purposes.
//
// Typical usage:
//
//	env, _ := envparser.Parse(".env")
//	snap := snapshoter.Capture(env, "pre-deploy")
//	_ = snapshoter.Save(snap, "snapshots/pre-deploy.json")
//
//	// Later...
//	prev, _ := snapshoter.Load("snapshots/pre-deploy.json")
//	diff := differ.Compare(prev.Env, env)
package snapshoter
