package main

import (
	"testing"

	"godlv/debug"
)

func TestRun(t *testing.T) {
	debug.SetEnabled(true)
	if err := run(); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}
}
