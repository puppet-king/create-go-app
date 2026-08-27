//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Default target to run when none is specified
var Default = Test

// Test 运行所有单元测试
func Test() error {
	fmt.Println("🧪 正在执行单元测试...")
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Build 编译 create-go-app 可执行文件
func Build() error {
	fmt.Println("🔨 正在构建 create-go-app 二进制文件...")
	_ = os.MkdirAll("bin", 0755)
	cmd := exec.Command("go", "build", "-o", "bin/create-go-app.exe", "./cmd/create-go-app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Install 将 create-go-app 安装到本地 GOPATH/bin
func Install() error {
	fmt.Println("📦 正在安装 create-go-app 到 GOPATH/bin...")
	cmd := exec.Command("go", "install", "./cmd/create-go-app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Tidy 同步并整理依赖
func Tidy() error {
	fmt.Println("🧹 正在执行 go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Clean 清理构建产物与临时文件
func Clean() error {
	fmt.Println("🗑️ 正在清理临时文件...")
	return os.RemoveAll("bin")
}
