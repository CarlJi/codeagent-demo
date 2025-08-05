# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered code agent that automatically processes GitHub Issues and Pull Requests, generating code modification suggestions using Claude or Gemini models.

## ✨ Features

- 🤖 **Multi-AI Support**: Claude and Gemini models with Docker/CLI execution modes
- 🔄 **GitHub Integration**: Automatic processing of Issues and Pull Requests via webhooks
- 🐳 **Flexible Deployment**: Docker containerized or local CLI execution
- 📁 **Smart Workspace**: Git Worktree-based workspace management with automatic cleanup
- 🔒 **Secure**: Webhook signature verification and configurable security settings

## 🚀 Quick Start

### 1. Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 2. Basic Setup

Create a `config.yaml` file:

```yaml
server:
  port: 8888

code_provider: gemini  # Options: claude, gemini
use_docker: false      # true for Docker, false for local CLI

workspace:
  base_dir: "./codeagent"
  cleanup_after: "24h"
```

### 3. Environment Variables

```bash
export GITHUB_TOKEN="your-github-token"
export GOOGLE_API_KEY="your-google-api-key"  # or CLAUDE_API_KEY for Claude
export WEBHOOK_SECRET="your-webhook-secret"
```

### 4. Start the Server

```bash
# Using the convenient start script
./scripts/start.sh                    # Gemini + CLI (default)
./scripts/start.sh -p claude -d       # Claude + Docker
./scripts/start.sh -p gemini -d       # Gemini + Docker

# Or run directly
go run ./cmd/server --config config.yaml
```

### 5. Configure GitHub Webhook

In your GitHub repository settings:
- **URL**: `https://your-domain.com/hook`
- **Content type**: `application/json`
- **Secret**: Same as your `WEBHOOK_SECRET`
- **Events**: `Issue comments`, `Pull request reviews`, `Pull requests`

### 6. Usage

Comment on GitHub Issues or PRs with these commands:

```bash
# Generate code for an issue
/code Implement user authentication with JWT tokens

# Continue development in a PR
/continue Add unit tests for the login functionality

# Fix issues in a PR
/fix Handle edge case for empty username
```

## ⚙️ Configuration

### Complete Configuration Example

```yaml
server:
  port: 8888

github:
  webhook_url: "http://localhost:8888/hook"

workspace:
  base_dir: "./codeagent"     # Supports relative paths
  cleanup_after: "24h"

claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"

docker:
  socket: "unix:///var/run/docker.sock"
  network: "bridge"

code_provider: gemini  # claude or gemini
use_docker: false      # true for Docker, false for CLI
```

### Configuration Methods

**Priority order**: Command line args > Environment variables > Config file

1. **Command Line Arguments**:
   ```bash
   go run ./cmd/server \
     --github-token "token" \
     --claude-api-key "key" \
     --webhook-secret "secret" \
     --port 8888
   ```

2. **Environment Variables**:
   ```bash
   export GITHUB_TOKEN="token"
   export CLAUDE_API_KEY="key"    # or GOOGLE_API_KEY
   export WEBHOOK_SECRET="secret"
   export CODE_PROVIDER="claude"  # or gemini
   export USE_DOCKER="true"       # or false
   ```

3. **Configuration File**: See example above

**Security Note**: Never store sensitive tokens in config files. Use environment variables or command line arguments.

### Execution Modes

| Mode | Description | Use Case |
|------|-------------|----------|
| **Claude + Docker** | Full containerized Claude environment | Production, isolated execution |
| **Claude + CLI** | Local Claude CLI | Development, faster startup |
| **Gemini + Docker** | Containerized Gemini environment | Production with Google AI |
| **Gemini + CLI** | Local Gemini CLI | Development, recommended for quick testing |

## 🛡️ Security

### Webhook Security

CodeAgent supports GitHub webhook signature verification to prevent unauthorized access:

```bash
# Set a strong webhook secret (32+ characters recommended)
export WEBHOOK_SECRET="your-very-strong-secret-here"
```

**Security Features**:
- SHA-256 signature verification (with SHA-1 fallback)
- Constant-time comparison to prevent timing attacks
- Automatic signature validation for all webhook requests

### Security Best Practices

- ✅ Use strong webhook secrets (32+ characters)
- ✅ Always configure webhook secrets in production
- ✅ Use HTTPS for webhook endpoints
- ✅ Regularly rotate API keys and secrets
- ✅ Limit GitHub token permissions to minimum required scope

## 🔧 Development

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

### Testing

```bash
# Health check
curl http://localhost:8888/health

# Test webhook processing
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json

# Build
make build

# Run tests
make test
```

### Debugging

```bash
# Enable debug logging
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

## 🤝 Contributing

We welcome contributions! Here's how you can help:

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Request Features](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [Submit Code](CONTRIBUTING.md)

Please read our [Contributing Guide](CONTRIBUTING.md) before making contributions.

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thank you to all developers and contributors who have made this project possible!