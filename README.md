# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

一个基于 AI 的自动化代码代理，通过 GitHub Webhook 自动处理 Issues 和 Pull Requests，生成代码修改建议。

## ✨ 核心特性

- 🤖 **多 AI 模型支持** - 支持 Claude 和 Gemini
- 🔄 **自动化工作流** - 自动处理 GitHub Issues 和 Pull Requests
- 🐳 **容器化执行** - Docker 容器化执行环境
- 📁 **智能工作区** - 基于 Git Worktree 的工作区管理
- 🔒 **安全可靠** - Webhook 签名验证，支持相对路径配置

## 🚀 快速开始

### 安装

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 配置

创建配置文件 `config.yaml`：

```yaml
# 服务器配置
server:
  port: 8888

# 工作区配置
workspace:
  base_dir: "./codeagent"  # 支持相对路径
  cleanup_after: "24h"

# AI 提供商配置
code_provider: claude      # 选项: claude, gemini
use_docker: false         # true=Docker容器, false=本地CLI

# AI 服务配置
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest" 
  timeout: "30m"

# Docker 配置
docker:
  socket: "unix:///var/run/docker.sock"
  network: "bridge"
```

### 环境变量

设置必需的环境变量：

```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # 或 GOOGLE_API_KEY
export WEBHOOK_SECRET="your-webhook-secret"
```

### 启动服务

#### 方式 1: 使用启动脚本（推荐）

```bash
# Gemini + CLI 模式（默认，开发推荐）
./scripts/start.sh

# Claude + Docker 模式（生产推荐）
./scripts/start.sh -p claude -d

# Claude + CLI 模式
./scripts/start.sh -p claude

# Gemini + Docker 模式
./scripts/start.sh -p gemini -d
```

#### 方式 2: 直接运行

```bash
go run ./cmd/server --config config.yaml
```

### GitHub Webhook 配置

在 GitHub 仓库设置中添加 Webhook：

- **URL**: `https://your-domain.com/hook`
- **Content type**: `application/json`
- **Secret**: 与 `WEBHOOK_SECRET` 相同
- **Events**: 选择 `Issue comments`, `Pull request reviews`, `Pull requests`

## 📖 使用指南

### 基本命令

在 GitHub Issue 或 PR 评论中使用以下命令：

| 命令 | 描述 | 示例 |
|------|------|------|
| `/code <描述>` | 在 Issue 中生成代码并创建 PR | `/code 实现用户登录功能` |
| `/continue <指令>` | 在 PR 中继续开发 | `/continue 添加单元测试` |
| `/fix <描述>` | 在 PR 中修复问题 | `/fix 修复登录验证逻辑错误` |

### 配置选项

#### AI 提供商选择

- **Claude**: Anthropic 的 Claude 模型，适合复杂代码生成
- **Gemini**: Google 的 Gemini 模型，快速响应

#### 执行模式选择

- **Docker 模式** (`use_docker: true`): 
  - 优点：隔离性好，适合生产环境
  - 缺点：启动稍慢，需要 Docker 环境

- **CLI 模式** (`use_docker: false`):
  - 优点：启动快速，适合开发环境
  - 缺点：需要本地安装 AI CLI 工具

## 🛠️ 开发指南

### 项目结构

```
codeagent/
├── cmd/server/          # 主程序入口
├── internal/            # 核心业务逻辑
│   ├── agent/          # 主要协调逻辑
│   ├── webhook/        # GitHub webhook 处理
│   ├── workspace/      # Git 工作区管理
│   ├── code/           # AI 提供商实现
│   ├── github/         # GitHub API 客户端
│   └── config/         # 配置管理
├── pkg/models/         # 共享数据结构
├── scripts/           # 工具脚本
└── docs/             # 文档
```

### 本地开发

```bash
# 构建
make build

# 运行测试
make test

# 健康检查
curl http://localhost:8888/health

# 测试 webhook
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

### 调试

```bash
# 设置详细日志
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

## 🔒 安全配置

### Webhook 签名验证

为防止恶意利用，CodeAgent 支持 GitHub Webhook 签名验证：

1. **配置 Webhook 密钥**：
   ```bash
   export WEBHOOK_SECRET="your-strong-secret-here"
   ```

2. **GitHub 设置**：在仓库 Webhook 设置中输入相同的密钥

3. **验证机制**：
   - 支持 SHA-256 签名验证（优先）
   - 向后兼容 SHA-1 签名验证
   - 使用恒定时间比较防止时序攻击

### 安全建议

- 使用强密码作为 webhook 密钥（建议 32+ 字符）
- 生产环境必须配置 webhook 密钥
- 使用 HTTPS 保护 webhook 端点
- 定期轮换 API 密钥和 webhook 密钥
- 限制 GitHub Token 权限范围

## 🤝 贡献

欢迎各种形式的贡献！请查看 [贡献指南](CONTRIBUTING.md) 了解如何参与项目开发。

### 贡献方式

- 🐛 [报告问题](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [功能建议](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [改进文档](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [提交代码](CONTRIBUTING.md#code-contributions)

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🙏 致谢

感谢所有为此项目做出贡献的开发者和用户！