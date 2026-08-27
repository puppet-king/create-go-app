package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// IsReplaceableFile checks if a file is suitable for string replacement.
func IsReplaceableFile(relPath string) bool {
	base := filepath.Base(relPath)
	ext := strings.ToLower(filepath.Ext(relPath))

	// Direct filename matches
	switch base {
	case "go.mod", "go.sum", "go.work", "Dockerfile", "Makefile", "README.md", "LICENSE", ".golangci.yml", ".env.example":
		return true
	}

	// Extension matches
	switch ext {
	case ".go", ".yaml", ".yml", ".json", ".md", ".toml", ".env", ".proto", ".sql", ".sh", ".bat", ".ps1":
		return true
	default:
		return false
	}
}

// ReplaceModuleInDir traverses targetDir and replaces occurrences of oldModule with newModule.
func ReplaceModuleInDir(targetDir, oldModule, newModule string) error {
	if oldModule == "" || newModule == "" || oldModule == newModule {
		return nil
	}

	return filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip .git and hidden directories except .agents
		if d.IsDir() {
			if d.Name() == ".agents" {
				return nil
			}
			if strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}

		if !IsReplaceableFile(path) {
			return nil
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(contentBytes)
		if strings.Contains(content, oldModule) {
			newContent := strings.ReplaceAll(content, oldModule, newModule)
			info, err := d.Info()
			perm := os.FileMode(0644)
			if err == nil {
				perm = info.Mode()
			}
			if err := os.WriteFile(path, []byte(newContent), perm); err != nil {
				return err
			}
		}

		return nil
	})
}
