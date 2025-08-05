# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered automation system that processes GitHub Issues and Pull Requests, generating intelligent code modifications through AI models like Claude and Gemini.

## 🚀 Quick Start

### Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### Basic Setup

1. **Set Environment Variables**
   ```bash
   export GITHUB_TOKEN="your-github-token"
   export CLAUDE_API_KEY="your-claude-api-key"    # or GOOGLE_API_KEY for Gemini
   export WEBHOOK_SECRET="your-webhook-secret"
   ```

2. **Start the Server**
   ```bash
   go run ./cmd/server
   ```

3. **Configure GitHub Webhook**
   - URL: `https://your-domain.com/hook`
   - Events: `Issue comments`, `Pull request reviews`, `Pull requests`
   - Secret: Same as your `WEBHOOK_SECRET`

4. **Test Usage**
   ```bash
   # Health check
   curl http://localhost:8888/health
   
   # In GitHub Issues/PRs, use these commands:
   /code Implement user authentication
   /continue Add unit tests
   ```

## 📋 Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration
- 🔄 **Automated Processing**: Handles GitHub Issues and Pull Requests automatically
- 🐳 **Flexible Deployment**: Docker containers or local CLI execution
- 📁 **Smart Workspace**: Git worktree-based isolated environments
- 🔒 **Security First**: Webhook signature verification and secure token handling
- ⚡ **Developer Friendly**: Simple setup with comprehensive configuration options

## ⚙️ Configuration

### Configuration File (Recommended)

Create `config.yaml`:

```yaml
# Provider settings
code_provider: claude    # Options: claude, gemini
use_docker: false        # true for Docker, false for CLI

# Server configuration
server:
  port: 8888

# Workspace settings
workspace:
  base_dir: "./codeagent"
  cleanup_after: "24h"

# AI provider configurations
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"
```

**Security Note**: Always set sensitive values via environment variables:
- `GITHUB_TOKEN`
- `CLAUDE_API_KEY` or `GOOGLE_API_KEY`
- `WEBHOOK_SECRET`

### Configuration Options

| Provider | Mode | Best For | Prerequisites |
|----------|------|----------|---------------|
| Claude + Docker | Production | Isolated execution | Docker installed |
| Claude + CLI | Development | Fast iteration | Claude CLI installed |
| Gemini + Docker | Production | Isolated execution | Docker installed |
| Gemini + CLI | Development | Fast iteration | Gemini CLI installed |

### Alternative Configuration Methods

**Environment Variables:**
```bash
export CODE_PROVIDER=claude
export USE_DOCKER=false
export PORT=8888
go run ./cmd/server
```

**Command Line Arguments:**
```bash
go run ./cmd/server \
  --github-token "token" \
  --claude-api-key "key" \
  --webhook-secret "secret" \
  --port 8888
```

## 🏗️ Development

### Building

```bash
# Development build
go build -o bin/codeagent ./cmd/server

# Production build
GOOS=linux GOARCH=amd64 go build -o bin/codeagent-linux ./cmd/server
```

### Testing

```bash
# Integration test
go run ./cmd/server --config test-config.yaml

# Test webhook
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

### Debugging

```bash
# Enable detailed logging
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

## 🔒 Security

### Webhook Security

1. **Enable Signature Verification**
   ```bash
   export WEBHOOK_SECRET="your-strong-secret-32-chars-min"
   ```

2. **GitHub Webhook Configuration**
   - Use HTTPS endpoints in production
   - Set the same secret in GitHub webhook settings
   - Enable signature verification (SHA-256 supported)

### Security Best Practices

- Use strong webhook secrets (32+ characters)
- Regularly rotate API keys and secrets
- Limit GitHub token permissions to minimum required
- Deploy with HTTPS in production
- Monitor webhook endpoint access logs

## 🤝 Contributing

We welcome contributions! Here's how to get involved:

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Request Features](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [Improve Docs](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [Submit Code](CONTRIBUTING.md)

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

**Questions?** Check our [documentation](docs/) or [open an issue](https://github.com/qiniu/codeagent/issues/new).
