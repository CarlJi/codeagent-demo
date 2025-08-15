# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

**CodeAgent** is an AI-powered automated code generation and collaboration system that processes GitHub Issues and Pull Requests through webhooks, providing intelligent code modifications and suggestions.

## ✨ Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration
- 🔄 **GitHub Integration**: Automatic Issue and PR processing
- 🐳 **Flexible Deployment**: Docker containers or local CLI
- 📁 **Smart Workspace**: Git worktree-based management
- 🔒 **Security**: Webhook signature verification
- ⚡ **High Performance**: Efficient workspace cleanup and management

## 🚀 Quick Start

### Prerequisites

- Go 1.19+
- GitHub Personal Access Token
- Claude API Key or Google API Key
- Docker (optional, for container mode)

### Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### Configuration

Create `config.yaml`:

```yaml
# Core Configuration
code_provider: claude  # Options: claude, gemini
use_docker: false      # true: Docker mode, false: CLI mode

server:
  port: 8888

workspace:
  base_dir: "./codeagent"  # Supports relative paths
  cleanup_after: "24h"

# Provider Settings (API keys set via environment variables)
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"
```

### Environment Variables

```bash
# Required
export GITHUB_TOKEN="your-github-token"
export WEBHOOK_SECRET="your-webhook-secret"

# Choose one based on your provider
export CLAUDE_API_KEY="your-claude-api-key"      # For Claude
export GOOGLE_API_KEY="your-google-api-key"      # For Gemini
```

### Run

```bash
# Method 1: Direct execution
go run ./cmd/server --config config.yaml

# Method 2: Using convenience script
./scripts/start.sh                    # Gemini + CLI (default)
./scripts/start.sh -p claude -d       # Claude + Docker
./scripts/start.sh -p gemini -d       # Gemini + Docker
./scripts/start.sh -p claude          # Claude + CLI
```

### GitHub Webhook Setup

1. Go to your repository → Settings → Webhooks
2. Add webhook:
   - **URL**: `https://your-domain.com/hook`
   - **Content type**: `application/json`
   - **Secret**: Your `WEBHOOK_SECRET` value
   - **Events**: `Issue comments`, `Pull request reviews`, `Pull requests`

## 💬 Usage

### Issue Commands

```bash
# Generate code for an issue
/code Implement user authentication with JWT tokens
```

### Pull Request Commands

```bash
# Continue development
/continue Add unit tests for the authentication module

# Fix issues
/fix Resolve the login validation bug
```

## 🛠️ Development

### Project Structure

```
codeagent/
├── cmd/server/           # Application entry point
├── internal/
│   ├── agent/           # Core orchestration logic
│   ├── webhook/         # GitHub webhook handling
│   ├── workspace/       # Git worktree management
│   ├── code/           # AI provider implementations
│   ├── github/         # GitHub API client
│   └── config/         # Configuration management
├── pkg/models/          # Shared data structures
├── scripts/            # Utility scripts
└── docs/              # Documentation
```

### Local Development

```bash
# Build
make build

# Test
make test

# Health check
curl http://localhost:8888/health

# Test webhook
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

### Configuration Options

| Option | Description | Values |
|--------|-------------|---------|
| `code_provider` | AI service provider | `claude`, `gemini` |
| `use_docker` | Execution environment | `true` (containers), `false` (local CLI) |
| `workspace.base_dir` | Workspace location | Absolute or relative path |
| `workspace.cleanup_after` | Auto-cleanup interval | Duration (e.g., "24h") |

### Security Best Practices

- ✅ Use strong webhook secrets (32+ characters)
- ✅ Always configure secrets in production
- ✅ Use HTTPS for webhook endpoints
- ✅ Regularly rotate API keys
- ✅ Limit GitHub token permissions

## 📖 Documentation

- [Architecture Design](docs/xgo-agent.md)
- [Relative Path Support](docs/relative-path-support.md)
- [Contributing Guide](CONTRIBUTING.md)

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Feature Requests](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)

## 📄 License

Licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thanks to all contributors who make this project possible!