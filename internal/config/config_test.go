package config

import (
	"bytes"
	"testing"
)

func TestParsePositionalArgs(t *testing.T) {
	args := []string{"user-service"}
	opts, err := Parse(args, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.NewModule != "user-service" {
		t.Errorf("expected user-service, got %s", opts.NewModule)
	}
	if opts.TargetDir != "user-service" {
		t.Errorf("expected targetDir user-service, got %s", opts.TargetDir)
	}
	if opts.DBName != "user_service" {
		t.Errorf("expected default DBName user_service, got %s", opts.DBName)
	}
	if opts.TemplateRepo != DefaultTemplateRepo {
		t.Errorf("expected default repo, got %s", opts.TemplateRepo)
	}
	if opts.OldModuleName != DefaultOldModuleName {
		t.Errorf("expected default old module, got %s", opts.OldModuleName)
	}
}

func TestParseCustomFlags(t *testing.T) {
	args := []string{
		"-repo", "https://github.com/test/custom-template.git",
		"-old-module", "custom-old",
		"-db-name", "custom_db",
		"-wx-appid", "wx123456",
		"-wx-secret", "secret654321",
		"-skip-wire",
		"-skip-tidy",
		"-skip-git",
		"github.com/my-org/my-app",
	}
	opts, err := Parse(args, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.NewModule != "github.com/my-org/my-app" {
		t.Errorf("expected github.com/my-org/my-app, got %s", opts.NewModule)
	}
	if opts.TargetDir != "my-app" {
		t.Errorf("expected targetDir my-app, got %s", opts.TargetDir)
	}
	if opts.DBName != "custom_db" {
		t.Errorf("expected custom_db, got %s", opts.DBName)
	}
	if opts.WechatAppID != "wx123456" {
		t.Errorf("expected wx123456, got %s", opts.WechatAppID)
	}
	if opts.WechatSecret != "secret654321" {
		t.Errorf("expected secret654321, got %s", opts.WechatSecret)
	}
	if !opts.SkipWire || !opts.SkipTidy || !opts.SkipGit {
		t.Errorf("expected all skip flags to be true")
	}
}

func TestParseInteractivePrompt(t *testing.T) {
	// simulate typing: project name -> db name (empty for default) -> appid -> secret
	input := bytes.NewBufferString("order-service\n\nwx999\nsec888\n")
	opts, err := Parse([]string{}, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.NewModule != "order-service" {
		t.Errorf("expected order-service, got %s", opts.NewModule)
	}
	if opts.DBName != "order_service" {
		t.Errorf("expected default DBName order_service, got %s", opts.DBName)
	}
	if opts.WechatAppID != "wx999" {
		t.Errorf("expected wx999, got %s", opts.WechatAppID)
	}
	if opts.WechatSecret != "sec888" {
		t.Errorf("expected sec888, got %s", opts.WechatSecret)
	}
}

func TestParseEmptyError(t *testing.T) {
	input := bytes.NewBufferString("\n")
	_, err := Parse([]string{}, input)
	if err == nil {
		t.Errorf("expected error when project name is empty, got nil")
	}
}

func TestParseInterleavedFlags(t *testing.T) {
	// Test positional arg first, flags after
	args := []string{
		"payment-service",
		"-skip-wire",
		"-skip-git",
		"-db-name", "pay_db",
		"-repo", "https://github.com/my-team/custom.git",
		"-old-module", "server-template",
	}
	opts, err := Parse(args, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.NewModule != "payment-service" {
		t.Errorf("expected payment-service, got %s", opts.NewModule)
	}
	if opts.DBName != "pay_db" {
		t.Errorf("expected pay_db, got %s", opts.DBName)
	}
	if !opts.SkipWire {
		t.Errorf("expected SkipWire to be true")
	}
	if !opts.SkipGit {
		t.Errorf("expected SkipGit to be true")
	}
	if opts.TemplateRepo != "https://github.com/my-team/custom.git" {
		t.Errorf("expected custom repo, got %s", opts.TemplateRepo)
	}
}
