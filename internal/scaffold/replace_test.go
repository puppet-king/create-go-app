package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReplaceModuleInDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "scaffold-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create dummy files
	goModContent := "module memory-server\n\ngo 1.21\n"
	if err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	mainGoContent := `package main

import (
	"fmt"
	"memory-server/internal/config"
)

func main() {
	fmt.Println("Starting memory-server")
}
`
	subDir := filepath.Join(tempDir, "cmd", "server")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create sub dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "main.go"), []byte(mainGoContent), 0644); err != nil {
		t.Fatalf("failed to write main.go: %v", err)
	}

	// Perform replacement
	err = ReplaceModuleInDir(tempDir, "memory-server", "user-service")
	if err != nil {
		t.Fatalf("ReplaceModuleInDir returned error: %v", err)
	}

	// Verify go.mod
	updatedMod, err := os.ReadFile(filepath.Join(tempDir, "go.mod"))
	if err != nil {
		t.Fatalf("failed to read go.mod: %v", err)
	}
	if !strings.Contains(string(updatedMod), "module user-service") {
		t.Errorf("expected module user-service in go.mod, got: %s", string(updatedMod))
	}

	// Verify main.go
	updatedMain, err := os.ReadFile(filepath.Join(subDir, "main.go"))
	if err != nil {
		t.Fatalf("failed to read main.go: %v", err)
	}
	if !strings.Contains(string(updatedMain), `"user-service/internal/config"`) {
		t.Errorf("expected user-service import in main.go, got: %s", string(updatedMain))
	}
	if !strings.Contains(string(updatedMain), `"Starting user-service"`) {
		t.Errorf("expected user-service string in main.go, got: %s", string(updatedMain))
	}
}
