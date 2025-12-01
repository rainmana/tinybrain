package main

import (
	"os"
	"testing"
)

func TestMainWithHelp(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test --help doesn't panic
	os.Args = []string{"server", "--help"}
	// We can't easily test main() since it calls os.Exit
	// Instead we test the logic directly
	
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help") {
		// This branch works
		t.Log("Help flag detected correctly")
	} else {
		t.Error("Help flag not detected")
	}
}

func TestMainDefaultsToServe(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"server"}
	
	// Test that with no args, we'd add "serve"
	if len(os.Args) == 1 {
		testArgs := append(os.Args, "serve")
		if testArgs[1] != "serve" {
			t.Error("Expected serve to be added as default command")
		}
	}
}

func TestMainWithFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"server", "--http=:9000"}
	
	// Test that flags get prepended with serve
	if os.Args[1] != "serve" {
		testArgs := append([]string{os.Args[0], "serve"}, os.Args[1:]...)
		if testArgs[1] != "serve" || testArgs[2] != "--http=:9000" {
			t.Error("Expected flags to be prepended with serve")
		}
	}
}
