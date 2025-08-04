# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

一个基于 AI 的代码智能助手，自动处理 GitHub Issues 和 Pull Requests，生成代码修改建议。

## 特性

- 🤖 支持多种 AI 模型（Claude、Gemini）
- 🔄 自动处理 GitHub Issues 和 Pull Requests
- 🐳 Docker 容器化执行环境
- 📁 基于 Git Worktree 的工作空间管理
- 🔒 支持 Webhook 签名验证和安全配置
- 📝 支持相对路径配置

## 快速开始

### 安装

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 配置

创建配置文件 `config.yaml`：

```yaml
server:
  port: 8888

github:
  webhook_url: "http://localhost:8888/hook"

workspace:
  base_dir: "./codeagent"
  cleanup_after: "24h"

# 选择代码提供商: claude 或 gemini
code_provider: claude
# 选择执行方式: true(Docker) 或 false(本地CLI)
use_docker: true
```

设置必需的环境变量：

```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # 或 GOOGLE_API_KEY
export WEBHOOK_SECRET="your-webhook-secret"
```

### 运行

推荐使用启动脚本：

```bash
# Gemini + 本地CLI模式（默认，开发推荐）
./scripts/start.sh

# Claude + Docker模式（生产推荐）
./scripts/start.sh -p claude -d

# 其他组合
./scripts/start.sh -p gemini -d    # Gemini + Docker
./scripts/start.sh -p claude       # Claude + 本地CLI

# 查看帮助
./scripts/start.sh --help
```

或直接运行：

```bash
go run ./cmd/server --config config.yaml
```

### GitHub Webhook 配置

在 GitHub 仓库设置中添加 Webhook：

- **URL**: `https://your-domain.com/hook`
- **Content type**: `application/json`
- **Secret**: 与 `WEBHOOK_SECRET` 环境变量相同
- **Events**: 选择 `Issue comments`、`Pull request reviews`、`Pull requests`

### 使用示例

在 GitHub Issue 中触发代码生成：

```
/code 实现用户登录功能，包括用户名/密码验证和JWT令牌生成
```

在 PR 评论中继续开发：

```
/continue 添加单元测试
```

修复代码问题：

```
/fix 修复登录验证逻辑bug
```

## 配置选项

### 配置方式

支持三种配置方式，优先级：命令行参数 > 环境变量 > 配置文件

1. **配置文件**（推荐）- 创建 `config.yaml`
2. **环境变量** - 设置 `GITHUB_TOKEN`、API Keys 等
3. **命令行参数** - 使用 `--github-token`、`--claude-api-key` 等

### 代码提供商配置

- **Claude**: 设置 `CLAUDE_API_KEY`，配置 `code_provider: claude`
- **Gemini**: 设置 `GOOGLE_API_KEY`，配置 `code_provider: gemini`

### 执行模式

- **Docker 模式** (`use_docker: true`): 容器化执行，适合生产环境
- **本地 CLI 模式** (`use_docker: false`): 使用本地 CLI 工具，适合开发环境

### 安全配置

**Webhook 签名验证**：

```bash
export WEBHOOK_SECRET="your-strong-secret-here"
```

**安全建议**：
- 使用强密码作为 webhook 密钥（推荐32+字符）
- 生产环境务必配置 webhook 密钥
- 使用 HTTPS 保护 webhook 端点
- 定期轮换 API 密钥和 webhook 密钥

## 开发

### 项目结构

```
codeagent/
├── cmd/server/           # 主程序入口
├── internal/
│   ├── agent/           # 核心业务逻辑
│   ├── webhook/         # Webhook 处理
│   ├── workspace/       # 工作空间管理
│   ├── code/           # AI 提供商实现
│   └── github/         # GitHub API 客户端
├── pkg/models/         # 数据模型
├── scripts/           # 工具脚本
└── config.yaml        # 配置文件
```

### 构建和测试

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

## 🤝 Contributing

We welcome all forms of contributions! Please check the [Contributing Guide](CONTRIBUTING.md) to learn how to participate in project development.

### Ways to Contribute

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Feature Requests](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [Submit Code](CONTRIBUTING.md#code-contributions)

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thank you to all developers and users who have contributed to this project!
