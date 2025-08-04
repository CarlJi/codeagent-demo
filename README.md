# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

An AI-powered code agent that automatically processes GitHub Issues and Pull Requests through webhook integration, supporting multiple AI providers (Claude, Gemini) and execution environments.

## 🚀 Quick Start

### Prerequisites

- Go 1.19+
- GitHub Token with repository access
- API key for Claude or Gemini
- Docker (optional, for containerized execution)

### Installation & Basic Setup

```bash
# Clone and install
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download

# Set environment variables
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY
export WEBHOOK_SECRET="your-webhook-secret"

# Start with default settings (Gemini + CLI mode)
./scripts/start.sh

# Or use specific configuration
./scripts/start.sh -p claude -d  # Claude + Docker mode
```

### GitHub Webhook Configuration

1. Go to repository Settings → Webhooks → Add webhook
2. Set URL: `https://your-domain.com/hook`
3. Content type: `application/json`
4. Secret: Same as your `WEBHOOK_SECRET`
5. Events: `Issue comments`, `Pull request reviews`, `Pull requests`

### Usage Commands

**In GitHub Issues:**
```
/code Implement user authentication with JWT tokens
```

**In Pull Request comments:**
```
/continue Add error handling and validation
/fix Resolve the type casting issue in line 42
```

## ⚙️ Configuration

### Environment Variables (Recommended)

```bash
# Required
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
export WEBHOOK_SECRET="your-strong-secret"

# AI Provider (choose one)
export CLAUDE_API_KEY="sk-ant-xxxxxxxxxxxx"    # For Claude
export GOOGLE_API_KEY="AIxxxxxxxxxxxx"         # For Gemini

# Optional
export CODE_PROVIDER=claude    # claude or gemini (default: gemini)
export USE_DOCKER=true         # true or false (default: false)
export PORT=8888               # Server port (default: 8888)
```

### Configuration File

Create `config.yaml` for advanced settings:

```yaml
# Core settings
code_provider: claude          # claude or gemini
use_docker: false             # Use Docker containers or local CLI

server:
  port: 8888

workspace:
  base_dir: "./workspace"     # Supports relative paths
  cleanup_after: "24h"

# AI provider settings
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"
```

**Note:** Sensitive data (tokens, API keys) should only be set via environment variables, not in config files.

## 🏗️ Architecture & Features

- **🤖 Multi-AI Support**: Claude and Gemini integration
- **🔄 Webhook-Driven**: Automatic GitHub event processing
- **🐳 Flexible Execution**: Docker containers or local CLI
- **📁 Smart Workspace**: Git worktree-based isolation
- **🔒 Security**: Webhook signature verification
- **⚡ Performance**: Optimized for both development and production

### Project Structure

```
codeagent/
├── cmd/server/          # Application entry point
├── internal/
│   ├── agent/          # Core orchestration logic
│   ├── webhook/        # GitHub webhook handling
│   ├── workspace/      # Git worktree management
│   ├── code/           # AI provider implementations
│   └── github/         # GitHub API client
├── pkg/models/         # Shared data structures
├── scripts/            # Utility scripts
└── docs/              # Documentation
```

## 🔧 Development

### Local Development

```bash
# Development mode (fast iteration)
./scripts/start.sh -p gemini    # Local CLI mode

# Production-like testing
./scripts/start.sh -p claude -d # Docker mode
```

### Testing

```bash
# Health check
curl http://localhost:8888/health

# Test webhook
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json

# Build
make build
```

### Debugging

```bash
# Enable debug logging
export LOG_LEVEL=debug
go run ./cmd/server
```

## 🛡️ Security

- **Webhook Verification**: SHA-256 signature validation
- **Token Security**: Environment-only credential storage  
- **Workspace Isolation**: Temporary Git worktrees with automatic cleanup
- **HTTPS Recommended**: Use secure endpoints in production

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

- 🐛 [Report Issues](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Request Features](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📖 [Improve Docs](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.
