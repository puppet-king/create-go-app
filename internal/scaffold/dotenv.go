package scaffold

import (
	"bufio"
	"bytes"
	"crypto/rand"
	mrand "math/rand"
	"os"
	"path/filepath"
	"strings"
)

// GenerateRandomKey generates a cryptographically secure random string of given length.
func GenerateRandomKey(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		for i := range b {
			b[i] = charset[mrand.Intn(len(charset))]
		}
		return string(b)
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

// GenerateDotEnv checks for .env.example in targetDir and creates a customized .env file.
// Returns true if .env was generated, false if .env.example was not found.
func GenerateDotEnv(targetDir, oldModuleName, dbName, wxAppID, wxSecret string) (bool, error) {
	examplePath := filepath.Join(targetDir, ".env.example")
	envPath := filepath.Join(targetDir, ".env")

	// If .env.example doesn't exist, nothing to do
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		return false, nil
	}

	content, err := os.ReadFile(examplePath)
	if err != nil {
		return false, err
	}

	parsed := GenerateCustomizedDotEnv(string(content), oldModuleName, dbName, wxAppID, wxSecret)
	if err := os.WriteFile(envPath, []byte(parsed), 0644); err != nil {
		return false, err
	}

	return true, nil
}

// GenerateCustomizedDotEnv converts .env.example into .env with DB DSN, WeChat credentials, and secure 32-char JWT keys.
func GenerateCustomizedDotEnv(example, oldModuleName, dbName, wxAppID, wxSecret string) string {
	scanner := bufio.NewScanner(strings.NewReader(example))
	var buf bytes.Buffer

	oldDBName := strings.ToLower(strings.ReplaceAll(oldModuleName, "-", "_"))
	if oldDBName == "" {
		oldDBName = "server_template"
	}

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Preserve comments and empty lines
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			buf.WriteString(line + "\n")
			continue
		}

		if strings.HasPrefix(trimmed, "DB_DSN=") {
			// Replace database name in DSN
			dsn := line
			if dbName != "" {
				if strings.Contains(dsn, "/"+oldDBName+"?") {
					dsn = strings.Replace(dsn, "/"+oldDBName+"?", "/"+dbName+"?", 1)
				} else if strings.Contains(dsn, "/server_template?") {
					dsn = strings.Replace(dsn, "/server_template?", "/"+dbName+"?", 1)
				} else if idx := strings.Index(dsn, "/"); idx != -1 {
					// Generic replacement after slash
					qIdx := strings.Index(dsn, "?")
					if qIdx != -1 && qIdx > idx {
						dsn = dsn[:idx+1] + dbName + dsn[qIdx:]
					}
				}
			}
			buf.WriteString(dsn + "\n")
			continue
		}

		if strings.HasPrefix(trimmed, "WECHAT_APPID=") {
			if wxAppID != "" {
				buf.WriteString("WECHAT_APPID=" + wxAppID + "\n")
			} else {
				buf.WriteString("WECHAT_APPID=\n")
			}
			continue
		}

		if strings.HasPrefix(trimmed, "WECHAT_SECRET=") {
			if wxSecret != "" {
				buf.WriteString("WECHAT_SECRET=" + wxSecret + "\n")
			} else {
				buf.WriteString("WECHAT_SECRET=\n")
			}
			continue
		}

		// Automatically generate secure 32-character random keys if these JWT fields exist in .env.example
		if strings.HasPrefix(trimmed, "JWT_TOKEN=") {
			buf.WriteString("JWT_TOKEN=" + GenerateRandomKey(32) + "\n")
			continue
		}
		if strings.HasPrefix(trimmed, "REFRESH_TOKEN=") {
			buf.WriteString("REFRESH_TOKEN=" + GenerateRandomKey(32) + "\n")
			continue
		}
		if strings.HasPrefix(trimmed, "ADMIN_JWT_TOKEN=") {
			buf.WriteString("ADMIN_JWT_TOKEN=" + GenerateRandomKey(32) + "\n")
			continue
		}
		if strings.HasPrefix(trimmed, "ADMIN_REFRESH_TOKEN=") {
			buf.WriteString("ADMIN_REFRESH_TOKEN=" + GenerateRandomKey(32) + "\n")
			continue
		}

		// Keep existing line intact (e.g. QINIU keys, SECRET_KEY, etc.)
		buf.WriteString(line + "\n")
	}

	return buf.String()
}
