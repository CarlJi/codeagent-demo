# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

An AI-powered code agent that automatically processes GitHub Issues and Pull Requests, generating intelligent code modifications through webhook-driven automation.

## ✨ Features

- 🤖 **Multiple AI Providers**: Support for Claude and Gemini models
- 🔄 **GitHub Integration**: Automatic processing of Issues and Pull Requests
- 🐳 **Flexible Execution**: Docker containers or local CLI modes
- 📁 **Smart Workspace**: Git worktree-based temporary workspace management
- 🔒 **Security First**: Webhook signature verification and secure token handling

## 🚀 Quick Start

### Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### Setup Environment

Set your required environment variables:

```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"
```

### Run with Scripts (Recommended)

Use the convenient startup script for different configurations:

```bash
./scripts/start.sh                    # Gemini + CLI mode (default)
./scripts/start.sh -p claude -d       # Claude + Docker mode
./scripts/start.sh -p gemini -d       # Gemini + Docker mode
./scripts/start.sh -p claude          # Claude + CLI mode
```

### Manual Configuration

Create `config.yaml`:

```yaml
server:
  port: 8888

workspace:
  base_dir: "./codeagent"
  cleanup_after: "24h"

# Choose your AI provider
code_provider: claude  # or gemini
use_docker: false      # true for containers, false for CLI

claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest" 
  timeout: "30m"
```

Then run:

```bash
go run ./cmd/server --config config.yaml
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `GITHUB_TOKEN` | GitHub personal access token | ✅ |
| `CLAUDE_API_KEY` | Anthropic Claude API key | ✅ (for Claude) |
| `GOOGLE_API_KEY` | Google Gemini API key | ✅ (for Gemini) |
| `WEBHOOK_SECRET` | GitHub webhook secret | ✅ |
| `CODE_PROVIDER` | AI provider: `claude` or `gemini` | ❌ |
| `USE_DOCKER` | Use Docker containers: `true`/`false` | ❌ |

### GitHub Webhook Setup

1. Go to your repository settings → Webhooks
2. Add webhook with:
   - **URL**: `https://your-domain.com/hook`
   - **Content type**: `application/json`
   - **Secret**: Same as your `WEBHOOK_SECRET`
   - **Events**: Issue comments, Pull request reviews, Pull requests

## 💡 Usage

Use these commands in GitHub Issues or Pull Request comments:

### Issue Commands
```bash
/code Implement user authentication with JWT tokens
```

### Pull Request Commands
```bash
/continue Add comprehensive unit tests
/fix Resolve the memory leak in the login handler
```

## 🏗️ Development

### Project Structure

```
codeagent/
├── cmd/server/           # Application entry point
├── internal/
│   ├── agent/           # Core orchestration logic
│   ├── webhook/         # GitHub webhook handling
│   ├── workspace/       # Git worktree management
│   ├── code/           # AI provider implementations
│   └── github/         # GitHub API client
├── pkg/models/         # Shared data structures
├── scripts/           # Utility scripts
└── docs/             # Documentation
```

### Build & Test

```bash
# Build binary
make build

# Run tests
make test

# Health check
curl http://localhost:8888/health

# Test webhook
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

### Debug Mode

```bash
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

## 🔒 Security

- **Webhook Verification**: SHA-256/SHA-1 signature validation with constant-time comparison
- **Token Security**: Environment-based credential management
- **HTTPS Required**: Use secure endpoints in production
- **Minimal Permissions**: Limit GitHub token scope to necessary permissions

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md).

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Feature Requests](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)

## 📄 License

This project is licensed under the [MIT License](LICENSE).
