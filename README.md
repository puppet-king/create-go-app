# 🚀 create-go-app

一个轻量、高效的 Go 项目脚手架 CLI 工具，类似前端生态的 `npx create-react-app` / `npm create vue`，支持一键拉取模板、自动替换模块名与 import 引用、依赖同步（`go mod tidy`）与 Wire 依赖注入代码生成。

仓库地址：`https://github.com/puppet-king/create-go-app.git`

---

## ⚡ 快速使用 (无需手动安装)

直接在终端执行：

```bash
# 格式: go run github.com/puppet-king/create-go-app@latest <新项目名/模块路径>
go run github.com/puppet-king/create-go-app@latest user-service
```

或者直接运行进入**交互式引导模式**：

```bash
go run github.com/puppet-king/create-go-app@latest
# 控制台将友好提示输入新项目名称、数据库名、微信凭证等
```

---

## 🛠️ CLI 参数与自定义选项

| 参数 / Flag | 说明 | 默认值 |
| :--- | :--- | :--- |
| `<project-name>` | 新项目名称或模块路径（如 `user-service` 或 `github.com/my-org/user-service`） | *交互式输入* |
| `-db-name` | 自定义 MySQL 数据库名（注入到 `.env` 的 `DB_DSN`） | 项目名的下划线形式 (如 `user_service`) |
| `-wx-appid` | 微信小程序 AppID（直接注入到 `.env`） | *留空或交互式输入* |
| `-wx-secret` | 微信小程序 AppSecret（直接注入到 `.env`） | *留空或交互式输入* |
| `-repo` | 自定义模板仓库 Git 地址 | `git@github.com:puppet-king/server-template.git` |
| `-old-module` | 原模板中的模块名称标识 | `server-template` |
| `-wire-path` | Wire 依赖注入代码生成的目录路径 | `./cmd/server` |
| `-skip-wire` | 是否跳过 Wire 代码生成 | `false` |
| `-skip-tidy` | 是否跳过 `go mod tidy` 依赖同步 | `false` |
| `-skip-git` | 是否跳过新项目 `git init` 仓库初始化 | `false` |

### 高级使用示例

```bash
# 指定自定义数据库名与微信凭证
go run github.com/puppet-king/create-go-app@latest order-service \
  -db-name order_db \
  -wx-appid wx1234567890abcdef \
  -wx-secret 0123456789abcdef0123456789abcdef
```

---

## 📦 工作流流水线

```text
  go run github.com/puppet-king/create-go-app@latest user-service
                        │
                        ▼
            【仓库 B: create-go-app】
                        │ 1. git clone --depth=1
                        ▼
         【仓库 A: server-template】
                        │ 2. 清理原 .git 历史
                        │ 3. 递归替换模块名与 import 引用 (含 cmd/mcp、magefile.go 等)
                        │ 4. 自动生成并注入定制化 .env 配置文件:
                        │    - DB_DSN 自动更新为当前项目库名: /user_service?
                        │    - 自动生成 32 位高强度安全 JWT 密钥
                        │    - 注入用户输入的 WECHAT_APPID 与 WECHAT_SECRET
                        │ 5. 执行 go mod tidy
                        │ 6. 执行 Wire 依赖注入代码生成
                        │ 7. 执行 git init
                        ▼
            生成独立的全新项目: ./user-service
```

---

## 🔑 默认后台账号与开箱即用

当新项目首次启动（`go run ./cmd/server/main.go`）时，系统自动执行数据库迁移并创建默认固定的超级管理员账号：

- **后台登录账号**：`admin`
- **后台初始密码**：`admin`

---

## 🤖 生成项目的运行指南

生成新项目后，进入项目目录：

```bash
cd user-service

# 1. 启动 Web / HTTP 后端服务
go run ./cmd/server/main.go

# 2. 启动 MCP AI Gateway 工具服务 (支持 Claude Desktop / Cursor / Antigravity 等 AI 工具直连)
go run ./cmd/mcp/main.go

# 3. 使用 Mage 管理与执行自动化脚本
mage -l                # 查看所有可用脚本任务
mage wire              # 重新生成 Wire 代码
mage build             # 编译 server、task、mcp 全量二进制
```
