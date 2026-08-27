package env

import (
	"fmt"
	"os/exec"
)

// CheckCommand checks if a given command is available in PATH.
func CheckCommand(name string) error {
	_, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("command '%s' not found in PATH. Please install it first", name)
	}
	return nil
}

// CheckGit verifies that git is available in system PATH.
func CheckGit() error {
	return CheckCommand("git")
}

// CheckGo verifies that go is available in system PATH.
func CheckGo() error {
	return CheckCommand("go")
}

// CheckAll verifies all required dependencies (git and go).
func CheckAll() error {
	if err := CheckGit(); err != nil {
		return err
	}
	if err := CheckGo(); err != nil {
		return err
	}
	return nil
}
