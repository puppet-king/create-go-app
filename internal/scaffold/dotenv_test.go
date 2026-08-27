package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestGenerateCustomizedDotEnvContent(t *testing.T) {
	input := `# Configuration
DB_DSN=root:123456@tcp(127.0.0.1:3306)/server_template?charset=utf8mb4&parseTime=True&loc=Local

# WeChat
WECHAT_APPID=
WECHAT_SECRET=

# JWT
JWT_TOKEN=
REFRESH_TOKEN=
ADMIN_JWT_TOKEN=
ADMIN_REFRESH_TOKEN=
`

	got := GenerateCustomizedDotEnv(input, "server-template", "order_service", "wx12345678", "secret8888")

	if !strings.Contains(got, "DB_DSN=root:123456@tcp(127.0.0.1:3306)/order_service?charset=utf8mb4&parseTime=True&loc=Local") {
		t.Errorf("expected DB_DSN with order_service, got:\n%s", got)
	}

	if !strings.Contains(got, "WECHAT_APPID=wx12345678") {
		t.Errorf("expected WECHAT_APPID=wx12345678, got:\n%s", got)
	}

	if !strings.Contains(got, "WECHAT_SECRET=secret8888") {
		t.Errorf("expected WECHAT_SECRET=secret8888, got:\n%s", got)
	}

	// Verify 32-character keys for JWT
	keys := []string{"JWT_TOKEN", "REFRESH_TOKEN", "ADMIN_JWT_TOKEN", "ADMIN_REFRESH_TOKEN"}
	for _, key := range keys {
		re := regexp.MustCompile(key + `=([a-zA-Z0-9]{32})`)
		matches := re.FindStringSubmatch(got)
		if len(matches) < 2 {
			t.Errorf("expected %s to have a 32-character random alphanumeric key, but got matching failure in:\n%s", key, got)
		}
	}
}

func TestGenerateRandomKey(t *testing.T) {
	k1 := GenerateRandomKey(32)
	k2 := GenerateRandomKey(32)

	if len(k1) != 32 || len(k2) != 32 {
		t.Fatalf("expected key length 32, got %d and %d", len(k1), len(k2))
	}
	if k1 == k2 {
		t.Errorf("expected random keys to differ, got identical keys: %s", k1)
	}
}

func TestGenerateDotEnv(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cga-dotenv-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Case 1: No .env.example
	generated, err := GenerateDotEnv(tempDir, "server-template", "my_db", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if generated {
		t.Errorf("expected generated to be false when .env.example does not exist")
	}

	// Case 2: With .env.example
	exampleFile := filepath.Join(tempDir, ".env.example")
	exampleContent := "DB_DSN=root:123456@tcp(127.0.0.1:3306)/server_template?charset=utf8mb4\nWECHAT_APPID=\nJWT_TOKEN=\n"
	if err := os.WriteFile(exampleFile, []byte(exampleContent), 0644); err != nil {
		t.Fatalf("failed to write .env.example: %v", err)
	}

	generated, err = GenerateDotEnv(tempDir, "server-template", "my_app_db", "wx999", "secret999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !generated {
		t.Errorf("expected generated to be true")
	}

	envFile := filepath.Join(tempDir, ".env")
	contentBytes, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("failed to read created .env file: %v", err)
	}

	content := string(contentBytes)
	if !strings.Contains(content, "my_app_db") {
		t.Errorf("expected my_app_db in DB_DSN, got: %s", content)
	}
	if !strings.Contains(content, "WECHAT_APPID=wx999") {
		t.Errorf("expected WECHAT_APPID=wx999, got: %s", content)
	}
	if !regexp.MustCompile(`JWT_TOKEN=[a-zA-Z0-9]{32}`).MatchString(content) {
		t.Errorf("expected 32-char JWT_TOKEN in .env, got: %s", content)
	}
}
