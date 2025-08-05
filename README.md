# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

An AI-powered code agent that automatically processes GitHub Issues and Pull Requests through webhook integration, generating intelligent code modifications and suggestions.

## ✨ Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration
- 🔄 **Automated Processing**: GitHub Issues and Pull Requests handling
- 🐳 **Flexible Execution**: Docker containers or local CLI
- 📁 **Smart Workspace**: Git Worktree-based management
- 🔒 **Secure**: Webhook signature verification and token protection

## 🚀 Quick Start

### Installation & Setup

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### Configuration

**Environment Variables (Recommended)**
```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"
```

**Configuration File** (`config.yaml`)
```yaml
server:
  port: 8888

workspace:
  base_dir: "./codeagent"  # Supports relative paths
  cleanup_after: "24h"

# Choose your AI provider
code_provider: claude  # Options: claude, gemini
use_docker: false      # true for Docker, false for local CLI

# Docker settings (if use_docker: true)
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"
```

### Running the Server

**Using Start Script (Recommended)**
```bash
./scripts/start.sh                    # Gemini + CLI (default)
./scripts/start.sh -p claude -d       # Claude + Docker
./scripts/start.sh -p gemini -d       # Gemini + Docker
./scripts/start.sh -p claude          # Claude + CLI
```

**Direct Command**
```bash
go run ./cmd/server --config config.yaml
# or
go run ./cmd/server --port 8888 --github-token "..." --claude-api-key "..."
```

### GitHub Webhook Setup

1. **Repository Settings** → **Webhooks** → **Add webhook**
2. **Payload URL**: `https://your-domain.com/hook`
3. **Content type**: `application/json`
4. **Secret**: Same as your `WEBHOOK_SECRET`
5. **Events**: Select `Issue comments`, `Pull request reviews`, `Pull requests`

### Usage

**Issue Commands**
```
/code Implement user authentication with JWT tokens
```

**PR Commands**
```
/continue Add comprehensive unit tests
/fix Resolve the memory leak in user session handling
```

## 🔧 Development

### Project Structure
```
codeagent/
├── cmd/server/          # Application entry point
├── internal/
│   ├── agent/           # Core orchestration logic
│   ├── webhook/         # GitHub webhook handling
│   ├── workspace/       # Git workspace management
│   ├── code/           # AI provider implementations
│   └── github/         # GitHub API client
├── pkg/models/         # Shared data structures
└── scripts/           # Utility scripts
```

### Build & Test

```bash
# Build
make build
# or
go build -o bin/codeagent ./cmd/server

# Test
make test
curl http://localhost:8888/health

# Integration test
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

### Debugging

```bash
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

## 🔒 Security

CodeAgent implements several security measures:

- **Webhook Signature Verification**: SHA-256/SHA-1 signature validation
- **Token Protection**: Sensitive data via environment variables only
- **Secure Defaults**: HTTPS endpoints and strong secret requirements

**Security Best Practices:**
- Use 32+ character webhook secrets
- Enable signature verification in production
- Regularly rotate API keys and secrets
- Limit GitHub token permissions to minimum required scope

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
