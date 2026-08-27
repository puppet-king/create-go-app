package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/your-name/create-go-app/internal/config"
	"github.com/your-name/create-go-app/internal/env"
	"github.com/your-name/create-go-app/internal/ui"
)

// Run executes the end-to-end scaffolding process.
func Run(opts *config.Options) error {
	ui.Banner()

	// 0. Pre-check environment
	if err := env.CheckAll(); err != nil {
		return err
	}

	targetDir := opts.TargetDir
	// Check if directory already exists
	if entries, err := os.ReadDir(targetDir); err == nil && len(entries) > 0 {
		return fmt.Errorf("目标目录 ./%s 已存在且不为空，请指定新的目录或清空后再试", targetDir)
	}

	totalSteps := 7
	currentStep := 1

	// Step 1: Clone template repository
	ui.Step(currentStep, totalSteps, fmt.Sprintf("正在拉取模板到 ./%s ...", targetDir))
	cloneCmd := exec.Command("git", "clone", "--depth=1", opts.TemplateRepo, targetDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		// Clean up created empty directory on failure
		_ = os.RemoveAll(targetDir)
		return fmt.Errorf("Git 克隆失败: %w", err)
	}
	ui.Success("模板拉取成功")
	currentStep++

	// Step 2: Remove old git history
	ui.Step(currentStep, totalSteps, "正在清理原仓库 Git 版本历史...")
	if err := removeDirAll(filepath.Join(targetDir, ".git")); err != nil {
		ui.Warn("清理 .git 目录时出现警告: %v", err)
	} else {
		ui.Success("Git 历史清理完毕")
	}
	currentStep++

	// Step 3: Replace import paths and module names
	ui.Step(currentStep, totalSteps, fmt.Sprintf("正在替换模块路径与引用 [%s -> %s]...", opts.OldModuleName, opts.NewModule))
	if err := ReplaceModuleInDir(targetDir, opts.OldModuleName, opts.NewModule); err != nil {
		return fmt.Errorf("模块替换失败: %w", err)
	}
	ui.Success("模块与依赖引用替换完成")
	currentStep++

	// Step 4: Generate .env from .env.example
	ui.Step(currentStep, totalSteps, fmt.Sprintf("正在初始化 .env 配置文件 (数据库: %s)...", opts.DBName))
	if generated, err := GenerateDotEnv(targetDir, opts.OldModuleName, opts.DBName, opts.WechatAppID, opts.WechatSecret); err != nil {
		ui.Warn("初始化 .env 文件失败: %v", err)
	} else if generated {
		ui.Success("已根据 .env.example 自动生成并注入 .env 配置文件")
	} else {
		ui.Info("未检测到 .env.example，跳过 .env 文件创建")
	}
	currentStep++

	// Step 5: go mod tidy
	if !opts.SkipTidy {
		ui.Step(currentStep, totalSteps, "正在执行 go mod tidy 同步依赖...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = targetDir
		tidyCmd.Stdout = os.Stdout
		tidyCmd.Stderr = os.Stderr
		if err := tidyCmd.Run(); err != nil {
			ui.Warn("go mod tidy 执行出现警告 (可在项目生成后手动执行): %v", err)
		} else {
			ui.Success("依赖同步完成")
		}
	} else {
		ui.Info("已跳过 go mod tidy")
	}
	currentStep++

	// Step 6: wire code generation
	if !opts.SkipWire {
		wireTarget := filepath.Join(targetDir, opts.WirePath)
		if _, err := os.Stat(wireTarget); err == nil {
			ui.Step(currentStep, totalSteps, fmt.Sprintf("正在执行 Wire 依赖注入代码生成 (%s)...", opts.WirePath))
			wireCmd := exec.Command("go", "run", "github.com/google/wire/cmd/wire@latest", opts.WirePath)
			wireCmd.Dir = targetDir
			wireCmd.Stdout = os.Stdout
			wireCmd.Stderr = os.Stderr
			if err := wireCmd.Run(); err != nil {
				ui.Warn("Wire 代码生成告警 (若项目未引入 Wire 可忽略，或稍后手动执行 wire): %v", err)
			} else {
				ui.Success("Wire 代码生成完毕")
			}
		} else {
			ui.Info("未检测到 Wire 目录 %s，自动跳过 Wire 代码生成", opts.WirePath)
		}
	} else {
		ui.Info("已跳过 Wire 代码生成")
	}
	currentStep++

	// Step 7: Initialize new git repository
	if !opts.SkipGit {
		ui.Step(currentStep, totalSteps, "正在初始化全新 Git 仓库...")
		gitInitCmd := exec.Command("git", "init")
		gitInitCmd.Dir = targetDir
		_ = gitInitCmd.Run()
		ui.Success("Git 仓库初始化完毕")
	} else {
		ui.Info("已跳过 Git 初始化")
	}

	// Final success output
	printCompletion(opts)
	return nil
}

// removeDirAll safely removes a directory tree, stripping read-only permissions (needed for Windows .git directories).
func removeDirAll(dir string) error {
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err == nil {
			_ = os.Chmod(path, 0777)
		}
		return nil
	})
	return os.RemoveAll(dir)
}

func printCompletion(opts *config.Options) {
	fmt.Printf("\n")
	ui.Success("🎉 项目 %s 创建成功！", opts.NewModule)

	fmt.Println("\n🔑 默认管理员账号 (系统服务启动时自动初始化)：")
	fmt.Println("   后台账号: admin")
	fmt.Println("   初始密码: admin")

	fmt.Println("\n🗄️  数据库配置 (.env)：")
	fmt.Printf("   数据库名: %s\n", opts.DBName)
	fmt.Printf("   连接示例: DB_DSN=root:123456@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local\n", opts.DBName)
	fmt.Printf("   ⚠️ 提示: 首次启动前请确保在 MySQL 中创建数据库: CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;\n", opts.DBName)

	if opts.WechatAppID == "" || opts.WechatSecret == "" {
		fmt.Println("\n📱 微信小程序配置提醒 (.env)：")
		if opts.WechatAppID == "" {
			fmt.Println("   WECHAT_APPID   : 暂未配置 (可在 .env 中手动补充)")
		} else {
			fmt.Printf("   WECHAT_APPID   : %s\n", opts.WechatAppID)
		}
		if opts.WechatSecret == "" {
			fmt.Println("   WECHAT_SECRET  : 暂未配置 (可在 .env 中手动补充)")
		} else {
			fmt.Println("   WECHAT_SECRET  : [已配置]")
		}
	}

	fmt.Println("\n👉 快速上手体验：")
	fmt.Printf("   cd %s\n", opts.TargetDir)

	if _, err := os.Stat(filepath.Join(opts.TargetDir, "cmd", "server", "main.go")); err == nil {
		fmt.Println("   go run ./cmd/server/main.go         # 启动 Web / HTTP 后端服务")
	} else if _, err := os.Stat(filepath.Join(opts.TargetDir, "main.go")); err == nil {
		fmt.Println("   go run main.go")
	} else {
		fmt.Println("   go run .")
	}

	if _, err := os.Stat(filepath.Join(opts.TargetDir, "cmd", "mcp", "main.go")); err == nil {
		fmt.Println("   go run ./cmd/mcp/main.go            # 启动 MCP AI Gateway 工具服务")
	}

	if _, err := os.Stat(filepath.Join(opts.TargetDir, "magefile.go")); err == nil {
		fmt.Println("   mage -l                             # 查看项目可用的 Mage 自动化脚本任务")
	}

	if !opts.SkipWire {
		if _, err := os.Stat(filepath.Join(opts.TargetDir, opts.WirePath)); err == nil {
			fmt.Printf("\n💡 提示：如需重新生成 Wire 依赖注入代码，可运行: wire %s\n\n", opts.WirePath)
		} else {
			fmt.Println()
		}
	} else {
		fmt.Println()
	}
}
