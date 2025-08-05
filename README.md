# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

**CodeAgent** is an AI-powered automation system that processes GitHub Issues and Pull Requests through webhook events, automatically generating code modifications and suggestions using Claude or Gemini AI models.

## ✨ Features

- 🤖 **Multi-AI Support** - Compatible with Anthropic Claude and Google Gemini
- 🔄 **GitHub Integration** - Automatic processing of Issues and Pull Requests via webhooks
- 🐳 **Flexible Execution** - Docker containers or local CLI tools
- 📁 **Smart Workspace** - Git worktree-based temporary workspace management
- 🔒 **Secure** - GitHub webhook signature verification and secure API handling

## 🚀 Quick Start

### 1. Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 2. Configuration

Create a `config.yaml` file:

```yaml
server:
  port: 8888

workspace:
  base_dir: "./codeagent"
  cleanup_after: "24h"

# Choose your AI provider
code_provider: claude  # Options: claude, gemini
use_docker: false      # true: Docker containers, false: local CLI

# Provider-specific settings
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"
```

Set up environment variables:

```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"    # or GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"
```

### 3. Run the Server

```bash
go run ./cmd/server --config config.yaml
```

### 4. Setup GitHub Webhook

In your GitHub repository settings:
- **URL**: `https://your-domain.com/hook`
- **Content type**: `application/json`
- **Secret**: Same value as your `WEBHOOK_SECRET`
- **Events**: Select `Issue comments`, `Pull request reviews`, `Pull requests`

### 5. Test the Setup

```bash
# Health check
curl http://localhost:8888/health

# Test webhook (optional)
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

## 📖 Usage

### Issue Commands

Create a new PR with generated code:
```
/code Implement user authentication with JWT tokens
```

### Pull Request Commands

Continue development with additional instructions:
```
/continue Add comprehensive unit tests for the authentication module
```

Fix specific issues:
```
/fix Handle edge case when user session expires
```

## ⚙️ Configuration

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GITHUB_TOKEN` | ✅ | GitHub personal access token |
| `WEBHOOK_SECRET` | ✅ | GitHub webhook secret for signature verification |
| `CLAUDE_API_KEY` | ⚠️ | Required when using Claude |
| `GOOGLE_API_KEY` | ⚠️ | Required when using Gemini |
| `CODE_PROVIDER` | ❌ | Override config file setting (claude/gemini) |
| `USE_DOCKER` | ❌ | Override config file setting (true/false) |
| `LOG_LEVEL` | ❌ | Set logging level (debug/info/warn/error) |

### Configuration File Options

#### Server Settings
```yaml
server:
  port: 8888
```

#### Workspace Management
```yaml
workspace:
  base_dir: "./codeagent"     # Supports relative paths
  cleanup_after: "24h"       # Auto-cleanup interval
```

#### Docker Configuration
```yaml
docker:
  socket: "unix:///var/run/docker.sock"
  network: "bridge"
```

### Execution Modes

#### Docker Mode (Production Recommended)
- **Pros**: Isolated environment, consistent across systems
- **Cons**: Slightly slower startup
- **Setup**: Ensure Docker is running

#### CLI Mode (Development Recommended)  
- **Pros**: Faster execution, easier debugging
- **Cons**: Requires local AI CLI installation
- **Setup**: Install `claude` or `gemini` CLI tools

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
│   ├── github/         # GitHub API client
│   └── config/         # Configuration management
├── pkg/models/         # Shared data structures
├── scripts/           # Utility scripts
└── docs/             # Documentation
```

### Building

```bash
# Build binary
make build

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o bin/codeagent-linux ./cmd/server
```

### Testing

```bash
# Run tests
make test

# Debug mode
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

## 🔒 Security

### Webhook Security
- Always use HTTPS in production
- Configure strong webhook secrets (32+ characters recommended)
- Enable signature verification for all webhook endpoints
- Regularly rotate API keys and secrets

### GitHub Token Permissions
Limit your GitHub token to the minimum required scopes:
- `repo` - Repository access
- `pull_requests:write` - PR management
- `issues:write` - Issue management

## 🤝 Contributing

We welcome contributions! Here's how you can help:

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Request Features](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [Submit Code](CONTRIBUTING.md)

Please read our [Contributing Guide](CONTRIBUTING.md) for development setup and guidelines.

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

**Questions?** Check our [documentation](docs/) or [open an issue](https://github.com/qiniu/codeagent/issues/new).