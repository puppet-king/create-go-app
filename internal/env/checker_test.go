package env

import (
	"testing"
)

func TestCheckCommands(t *testing.T) {
	// go and git should typically exist in the test environment
	if err := CheckGo(); err != nil {
		t.Logf("check go warning: %v", err)
	}

	if err := CheckGit(); err != nil {
		t.Logf("check git warning: %v", err)
	}

	err := CheckCommand("non-existent-command-xyz-12345")
	if err == nil {
		t.Errorf("expected error for non-existent command, got nil")
	}
}
