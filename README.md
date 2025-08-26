# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

> An AI-powered code generation system that processes GitHub Issues and Pull Requests through webhook-driven automation.

## 🚀 Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration with Docker/CLI options
- 🔄 **Webhook Automation**: Automatic processing of GitHub Issues and Pull Requests  
- 📁 **Smart Workspace**: Git Worktree-based temporary workspace management
- 🔐 **Security First**: Webhook signature verification and configurable secrets

## 🚀 Quick Start

### 1. Installation & Setup

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 2. Configuration

**Environment Variables (Required):**
```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"
```

**Create `config.yaml` (Optional):**
```yaml
code_provider: claude    # Options: claude, gemini
use_docker: false       # true for Docker, false for CLI
server:
  port: 8888
workspace:
  base_dir: "./codeagent"
  cleanup_after: "24h"
```

### 3. Run the Server

```bash
# Quick start with script (recommended)
./scripts/start.sh                    # Gemini + CLI (default)
./scripts/start.sh -p claude -d       # Claude + Docker

# Or run directly
go run ./cmd/server --port 8888
```

### 4. Setup GitHub Webhook

Configure in your repository settings:
- **URL**: `https://your-domain.com/hook`
- **Events**: Issue comments, Pull request reviews
- **Secret**: Same as `WEBHOOK_SECRET`

## 💡 Usage Examples

Trigger CodeAgent through GitHub comments:

```bash
# Generate code for an issue
/code Implement user authentication with JWT

# Continue development in PR
/continue Add unit tests for the login functionality

# Fix specific issues
/fix Resolve the validation logic bug
```

## 🔧 Development

### Project Structure

```
codeagent/
├── cmd/server/           # Main application entry
├── internal/             # Core business logic  
│   ├── agent/           # Orchestration logic
│   ├── webhook/         # GitHub webhook handling
│   ├── workspace/       # Git worktree management
│   └── code/            # AI provider implementations
├── pkg/models/          # Shared data structures
└── scripts/             # Utility scripts
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

# Debug mode
export LOG_LEVEL=debug
go run ./cmd/server
```

## 🛡️ Security

CodeAgent includes built-in security features:

- **Webhook Signature Verification**: SHA-256/SHA-1 signature validation
- **Secret Management**: Environment-based configuration for sensitive data  
- **HTTPS Support**: Secure webhook endpoints
- **Token Scope Limiting**: Minimal GitHub permissions required

**Security Best Practices:**
- Use strong webhook secrets (32+ characters)
- Always configure secrets in production
- Regularly rotate API keys and tokens
- Use HTTPS endpoints

## 🤝 Contributing

We welcome contributions! Check our [Contributing Guide](CONTRIBUTING.md) for details.

**Quick Links:**
- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Feature Requests](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)  
- 📝 [Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)

## 📄 License

Licensed under the [MIT License](LICENSE).
