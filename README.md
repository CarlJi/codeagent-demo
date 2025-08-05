# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered code agent that automatically processes GitHub Issues and Pull Requests, generating code modification suggestions through webhook-driven architecture.

## 📋 Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage](#usage)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

- 🤖 **Multi-AI Support** - Claude and Gemini integration with Docker/CLI modes
- 🔄 **Auto Processing** - Automatic GitHub Issues and Pull Requests handling
- 🐳 **Flexible Execution** - Docker containerized or local CLI execution
- 📁 **Smart Workspace** - Git Worktree-based workspace management with auto-cleanup
- 🔒 **Secure Webhooks** - GitHub webhook signature verification
- ⚡ **High Performance** - Efficient workspace isolation and parallel processing

## 🚀 Quick Start

### 1. Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 2. Environment Setup

```bash
# Required environment variables
export GITHUB_TOKEN="your-github-token"
export WEBHOOK_SECRET="your-webhook-secret"

# Choose AI provider (one of the following)
export CLAUDE_API_KEY="your-claude-api-key"     # For Claude
export GOOGLE_API_KEY="your-google-api-key"     # For Gemini
```

### 3. Quick Launch

```bash
go run ./cmd/server --port 8888
```

### 4. GitHub Webhook Setup

1. Go to your GitHub repository → Settings → Webhooks
2. Add webhook with URL: `https://your-domain.com/hook`
3. Content type: `application/json`
4. Secret: Same as your `WEBHOOK_SECRET`
5. Events: `Issue comments`, `Pull request reviews`, `Pull requests`

### 5. Test Installation

```bash
# Health check
curl http://localhost:8888/health

# Test webhook (optional)
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

## ⚙️ Configuration

### Configuration Methods

**Method 1: Configuration File (Recommended)**

Create `config.yaml`:

```yaml
# AI Provider Selection
code_provider: claude # Options: claude, gemini
use_docker: true      # true: Docker mode, false: CLI mode

# Server Configuration
server:
  port: 8888

# GitHub Configuration
github:
  webhook_url: "http://localhost:8888/hook"

# Workspace Management
workspace:
  base_dir: "./codeagent"  # Supports relative paths
  cleanup_after: "24h"

# AI Provider Settings
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"

```

**Method 2: Environment Variables**

```bash
export CODE_PROVIDER=claude    # or gemini
export USE_DOCKER=true         # or false
export PORT=8888
```

**Method 3: Command Line Arguments**

```bash
go run ./cmd/server \
  --github-token "your-token" \
  --claude-api-key "your-key" \
  --webhook-secret "your-secret" \
  --port 8888
```

### Security Configuration

**Webhook Security**
- Always set `WEBHOOK_SECRET` for production
- Use HTTPS endpoints
- Enable signature verification (SHA-256/SHA-1 supported)
- Regularly rotate API keys and secrets

**Best Practices**
- Use strong secrets (32+ characters)
- Limit GitHub token permissions
- Keep sensitive data in environment variables, not config files

### Provider Comparison

| Feature | CLI Mode | Docker Mode |
|---------|----------|-------------|
| **Speed** | ⚡ Faster | 🐢 Slower startup |  
| **Dependencies** | 📋 Requires local CLI | 🐳 Self-contained |
| **Security** | 🔓 Host access | 🔒 Isolated |
| **Use Case** | 🛠️ Development | 🚀 Production |

## 📖 Usage

### GitHub Commands

**In Issue Comments:**
```
/code Implement user authentication with JWT tokens
```

**In PR Comments:**
```
/continue Add comprehensive unit tests
```

**In PR Review Comments:**
```
/continue Optimize this database query
```

### Supported Scenarios

- ✅ Issue comment processing → Auto PR creation
- ✅ PR comment collaboration → Code updates  
- ✅ PR review comments → Targeted fixes
- ✅ Batch review processing → Multiple fixes
- ✅ Context-aware responses → Historical conversation

## 🛠️ Development

### Project Structure

```
codeagent/
├── cmd/server/           # Application entry point
├── internal/             # Core business logic
│   ├── agent/           # Main orchestration
│   ├── webhook/         # GitHub webhook handling  
│   ├── workspace/       # Git worktree management
│   ├── code/           # AI provider implementations
│   ├── github/         # GitHub API client
│   └── config/         # Configuration management
├── pkg/models/          # Shared data structures
├── scripts/            # Utility scripts
└── docs/              # Documentation
```

### Development Commands

```bash
# Build
make build
go build -o bin/codeagent ./cmd/server

# Test
make test

# Run with debugging
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml

# Cross-compilation
GOOS=linux GOARCH=amd64 go build -o bin/codeagent-linux ./cmd/server
```

### Development Workflow

1. **Local Setup**: Use CLI mode for faster iteration
2. **Testing**: Send test webhooks to verify functionality  
3. **Docker Testing**: Switch to Docker mode for production-like testing
4. **Workspace**: Temporary worktrees in `/tmp/codeagent` (auto-cleanup after 24h)

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Ways to Contribute

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Feature Requests](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)  
- 📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [Submit Code](CONTRIBUTING.md#code-contributions)

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thank you to all developers and users who have contributed to this project!