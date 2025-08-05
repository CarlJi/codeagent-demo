# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent 是一个基于 AI 的智能代码助手，通过 GitHub Webhooks 自动处理 Issues 和 Pull Requests，提供代码生成、修改和审查功能，支持多种 AI 模型和执行环境。

## ✨ 核心特性

- 🤖 **多模型支持**: 支持 Claude 和 Gemini 两种 AI 模型
- 🔄 **智能处理**: 自动处理 GitHub Issues 和 Pull Requests
- 🐳 **灵活部署**: 支持 Docker 容器化和本地 CLI 两种执行模式
- 📁 **工作空间管理**: 基于 Git Worktree 的临时工作空间，自动清理
- 🔒 **安全验证**: 支持 GitHub Webhook 签名验证
- 🎯 **指令驱动**: 通过简单的注释指令触发各种操作

## 🚀 快速开始

### 1. 安装配置

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 2. 环境配置

设置必要的环境变量：

```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # 或使用 GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"
```

### 3. 运行服务

**方式一：使用启动脚本（推荐）**

```bash
./scripts/start.sh                    # Gemini + CLI 模式（默认）
./scripts/start.sh -p claude -d       # Claude + Docker 模式
./scripts/start.sh -p claude          # Claude + CLI 模式
```

**方式二：直接运行**

```bash
go run ./cmd/server --port 8888
```

### 4. 配置 GitHub Webhook

在 GitHub 仓库设置中添加 Webhook：
- **URL**: `https://your-domain.com/hook`
- **Content type**: `application/json`
- **Secret**: 与 `WEBHOOK_SECRET` 相同
- **Events**: `Issue comments`, `Pull request reviews`, `Pull requests`

### 5. 开始使用

在 GitHub Issue 或 PR 中使用指令：

```bash
/code 实现用户登录功能，包括用户名密码验证和JWT生成
/continue 添加单元测试
/fix 修复登录验证逻辑的bug
```

## ⚙️ 配置说明

### 配置文件

创建 `config.yaml` 文件：

```yaml
# 服务配置
server:
  port: 8888

# 工作空间配置
workspace:
  base_dir: "./codeagent"
  cleanup_after: "24h"

# AI 模型选择
code_provider: claude  # claude 或 gemini
use_docker: false      # true=Docker模式，false=CLI模式

# Claude 配置
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

# Gemini 配置  
gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"
```

> **安全提示**: 敏感信息（如 Token、API Key、Webhook Secret）应通过环境变量设置，不要写入配置文件。

### 模式选择

| 模式 | 优势 | 适用场景 |
|------|------|----------|
| **CLI 模式** | 启动快，资源占用少 | 开发测试 |
| **Docker 模式** | 环境隔离，功能完整 | 生产部署 |

### AI 模型对比

| 模型 | 特点 | API Key |
|------|------|---------|
| **Claude** | 代码质量高，理解能力强 | `CLAUDE_API_KEY` |
| **Gemini** | 响应速度快，成本较低 | `GOOGLE_API_KEY` |

## 🔧 开发指南

### 项目架构

```
GitHub Events → Webhook → CodeAgent → 工作空间 → AI处理 → PR更新
```

**核心组件**：
- **Agent** (`internal/agent/`): 主要业务逻辑编排
- **Webhook Handler** (`internal/webhook/`): GitHub 事件处理
- **Workspace Manager** (`internal/workspace/`): Git 工作空间管理
- **Code Providers** (`internal/code/`): AI 模型接口实现

### 本地开发

**1. 构建项目**

```bash
# 构建二进制文件
make build

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o bin/codeagent-linux ./cmd/server
```

**2. 测试运行**

```bash
# 健康检查
curl http://localhost:8888/health

# 测试 Webhook
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

**3. 调试模式**

```bash
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

### 支持的指令

| 指令 | 使用场景 | 示例 |
|------|----------|------|
| `/code` | Issue 中生成代码 | `/code 实现用户认证功能` |
| `/continue` | PR 中继续开发 | `/continue 添加错误处理` |
| `/fix` | 修复代码问题 | `/fix 内存泄漏问题` |

## 🛡️ 安全配置

### Webhook 签名验证

为防止恶意请求，CodeAgent 支持 GitHub Webhook 签名验证：

```bash
# 设置强密码作为 webhook secret（建议32位以上）
export WEBHOOK_SECRET="your-strong-secret-32-chars-long"
```

**安全建议**：
- 生产环境必须配置 webhook secret
- 使用 HTTPS 保护 webhook 端点
- 定期轮换 API 密钥和 webhook secret
- 限制 GitHub Token 权限范围

## 🤔 常见问题

**Q: Docker 模式启动失败？**
A: 检查 Docker 服务是否运行，确保有足够的磁盘空间。

**Q: CLI 模式找不到命令？**
A: 确保已安装对应的 CLI 工具：`claude` 或 `gemini`。

**Q: Webhook 收不到事件？**
A: 检查 GitHub Webhook 配置和网络连接，确保端口可访问。

**Q: 工作空间清理失败？**
A: 检查磁盘权限和空间，默认24小时后自动清理。

## 🤝 贡献

欢迎各种形式的贡献！

- 🐛 [报告 Bug](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [功能建议](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [改进文档](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [提交代码](CONTRIBUTING.md)

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源。

---

感谢所有为此项目贡献的开发者和用户！
