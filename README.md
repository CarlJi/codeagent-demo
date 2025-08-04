# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered code agent that automatically processes GitHub Issues and Pull Requests, generating intelligent code modifications and suggestions through webhook-driven automation.

## 📋 Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage](#usage)
- [Development](#development)
- [Security](#security)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration with Docker/CLI modes
- 🔄 **Automated Workflow**: Seamless GitHub Issues and Pull Requests processing
- 🐳 **Flexible Deployment**: Docker containers or local CLI execution
- 📁 **Smart Workspace**: Git Worktree-based isolation and management
- 🔐 **Enterprise Security**: Webhook signature verification and secure token handling

## 🚀 Quick Start

### 1. Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 2. Environment Setup

Set up your environment variables:

```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"
```

### 3. Quick Launch

```bash
go run ./cmd/server --port 8888
```

### 4. Verify Installation

```bash
# Check health
curl http://localhost:8888/health

# Configure GitHub webhook
# URL: http://your-domain.com/hook
# Events: Issue comments, Pull request reviews
# Secret: Same as WEBHOOK_SECRET
```

## ⚙️ Configuration

### Environment Variables

Required environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `GITHUB_TOKEN` | GitHub personal access token | `ghp_xxxxxxxxxxxx` |
| `CLAUDE_API_KEY` | Claude API key (for Claude provider) | `sk-ant-xxxxxxxxxxxx` |
| `GOOGLE_API_KEY` | Google API key (for Gemini provider) | `AIzaxxxxxxxxxx` |
| `WEBHOOK_SECRET` | GitHub webhook secret | `your-strong-secret` |

### Configuration File (Optional)

Create `config.yaml` for advanced configuration:

```yaml
server:
  port: 8888

github:
  webhook_url: "http://localhost:8888/hook"

workspace:
  base_dir: "./codeagent"  # Supports relative paths
  cleanup_after: "24h"

# Provider selection
code_provider: claude  # Options: claude, gemini  
use_docker: false     # true for Docker, false for CLI

# Docker settings (when use_docker: true)
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"
```

### Provider & Mode Combinations

Configure your provider and mode through the configuration file:

| Provider | Docker Mode | CLI Mode | Recommended For |
|----------|-------------|----------|-----------------|
| Claude | `use_docker: true` | `use_docker: false` | Production / Development |
| Gemini | `use_docker: true` | `use_docker: false` | Development (default) |

## 📝 Usage

CodeAgent responds to GitHub comments with these commands:

### Issue Commands
```bash
/code <description>  # Generate code and create PR
```

### PR Commands  
```bash
/continue <instruction>  # Continue development
/fix <description>      # Fix code issues
```

### Examples
```bash
# In GitHub Issue
/code Implement user authentication with JWT tokens

# In PR Comment  
/continue Add comprehensive unit tests
/fix Handle edge case for empty input validation
```

## 🛠️ Development

### Project Structure

```
codeagent/
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── agent/agent.go          # Core orchestration logic
│   ├── webhook/handler.go      # GitHub webhook processing
│   ├── workspace/manager.go    # Git workspace management  
│   ├── code/                   # AI provider implementations
│   ├── github/client.go        # GitHub API integration
│   └── config/config.go        # Configuration management
├── pkg/models/                 # Shared data structures
└── docs/                       # Documentation
```

### Build & Test

```bash
# Build binary
go build -o bin/codeagent ./cmd/server

# Run tests  
go test ./...

# Integration testing
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json

# Debug mode
export LOG_LEVEL=debug
go run ./cmd/server
```

## 🔐 Security

### Webhook Security

CodeAgent implements GitHub webhook signature verification:

```bash
# Configure webhook secret
export WEBHOOK_SECRET="your-strong-secret-32-chars-minimum"
```

### GitHub Webhook Setup

1. Repository Settings → Webhooks → Add webhook
2. **Payload URL**: `https://your-domain.com/hook`
3. **Content type**: `application/json`  
4. **Secret**: Enter your `WEBHOOK_SECRET` value
5. **Events**: Select `Issue comments`, `Pull request reviews`, `Pull requests`

### Security Best Practices

- Use HTTPS for webhook endpoints in production
- Use strong webhook secrets (32+ characters)
- Rotate API keys and secrets regularly
- Limit GitHub token permissions to minimum required scope
- Never commit secrets to configuration files

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md).

**Quick Links:**
- [🐛 Report Bug](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- [💡 Request Feature](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)  
- [📝 Improve Docs](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)

## 📄 License

Licensed under the [MIT License](LICENSE).
