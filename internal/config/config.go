package config

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

const (
	// DefaultTemplateRepo is the default Git repository URL for the template.
	DefaultTemplateRepo = "git@github.com:puppet-king/server-template.git"
	// DefaultOldModuleName is the original Go module name in the template.
	DefaultOldModuleName = "server-template"
	// DefaultWirePath is the default directory to execute wire.
	DefaultWirePath = "./cmd/server"
)

// Options holds the resolved configuration for scaffolding.
type Options struct {
	NewModule     string
	TargetDir     string
	DBName        string
	WechatAppID   string
	WechatSecret  string
	TemplateRepo  string
	OldModuleName string
	WirePath      string
	SkipWire      bool
	SkipTidy      bool
	SkipGit       bool
}

// Parse parses CLI arguments and flags, falling back to interactive prompt if new module is omitted.
func Parse(args []string, stdinReader io.Reader) (*Options, error) {
	fs := flag.NewFlagSet("create-go-app", flag.ContinueOnError)

	repo := fs.String("repo", DefaultTemplateRepo, "Template repository Git URL")
	oldMod := fs.String("old-module", DefaultOldModuleName, "Original module name in template")
	wirePath := fs.String("wire-path", DefaultWirePath, "Target directory for wire code generation")
	dbName := fs.String("db-name", "", "Database name (defaults to snake_case of target directory)")
	wxAppID := fs.String("wx-appid", "", "WeChat Mini Program AppID")
	wxSecret := fs.String("wx-secret", "", "WeChat Mini Program AppSecret")
	skipWire := fs.Bool("skip-wire", false, "Skip wire code generation")
	skipTidy := fs.Bool("skip-tidy", false, "Skip go mod tidy")
	skipGit := fs.Bool("skip-git", false, "Skip git init")

	reordered := reorderArgs(args)
	if err := fs.Parse(reordered); err != nil {
		return nil, err
	}

	var newModule string
	isInteractive := false
	positional := fs.Args()
	if len(positional) > 0 {
		newModule = strings.TrimSpace(positional[0])
	}

	var scanner *bufio.Scanner
	if stdinReader != nil {
		scanner = bufio.NewScanner(stdinReader)
	}

	// Interactive fallback if not specified
	if newModule == "" && scanner != nil {
		isInteractive = true
		fmt.Print("👉 请输入新项目名称或模块路径 (例如 user-service 或 github.com/user/service): ")
		if scanner.Scan() {
			newModule = strings.TrimSpace(scanner.Text())
		}
	}

	if newModule == "" {
		return nil, fmt.Errorf("项目名称/模块路径不能为空。用法: go run github.com/your-name/create-go-app@latest <新项目名>")
	}

	// Target directory is the last element of module path or direct folder name
	targetDir := filepath.Base(newModule)
	if targetDir == "." || targetDir == "/" || targetDir == "\\" {
		targetDir = newModule
	}

	resolvedDBName := *dbName
	if resolvedDBName == "" {
		resolvedDBName = strings.ToLower(strings.ReplaceAll(targetDir, "-", "_"))
	}

	resolvedWxAppID := *wxAppID
	resolvedWxSecret := *wxSecret

	// If interactive mode, prompt for DB name and WeChat credentials
	if isInteractive && scanner != nil {
		fmt.Printf("👉 请确认数据库名称 (直接回车默认: %s): ", resolvedDBName)
		if scanner.Scan() {
			val := strings.TrimSpace(scanner.Text())
			if val != "" {
				resolvedDBName = val
			}
		}

		if resolvedWxAppID == "" {
			fmt.Print("👉 请输入/粘贴 微信 WECHAT_APPID (若暂无请直接回车): ")
			if scanner.Scan() {
				resolvedWxAppID = strings.TrimSpace(scanner.Text())
			}
		}

		if resolvedWxSecret == "" {
			fmt.Print("👉 请输入/粘贴 微信 WECHAT_SECRET (若暂无请直接回车): ")
			if scanner.Scan() {
				resolvedWxSecret = strings.TrimSpace(scanner.Text())
			}
		}
	}

	return &Options{
		NewModule:     newModule,
		TargetDir:     targetDir,
		DBName:        resolvedDBName,
		WechatAppID:   resolvedWxAppID,
		WechatSecret:  resolvedWxSecret,
		TemplateRepo:  *repo,
		OldModuleName: *oldMod,
		WirePath:      *wirePath,
		SkipWire:      *skipWire,
		SkipTidy:      *skipTidy,
		SkipGit:       *skipGit,
	}, nil
}

// reorderArgs moves flags before positional arguments to support flexible CLI flag ordering.
func reorderArgs(args []string) []string {
	var flagArgs []string
	var posArgs []string

	valFlags := map[string]bool{
		"repo":       true,
		"old-module": true,
		"wire-path":  true,
		"db-name":    true,
		"wx-appid":   true,
		"wx-secret":  true,
	}

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "--" {
			posArgs = append(posArgs, args[i+1:]...)
			break
		}
		if strings.HasPrefix(arg, "-") {
			flagArgs = append(flagArgs, arg)
			name := strings.TrimLeft(arg, "-")
			if !strings.Contains(arg, "=") && valFlags[name] {
				if i+1 < len(args) {
					i++
					flagArgs = append(flagArgs, args[i])
				}
			}
		} else {
			posArgs = append(posArgs, arg)
		}
		i++
	}

	return append(flagArgs, posArgs...)
}
