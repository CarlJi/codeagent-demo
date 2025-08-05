# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered automation system that processes GitHub Issues and Pull Requests, automatically generating and modifying code through intelligent webhook-driven workflows.

## ✨ Features

- 🤖 **Multi-AI Support**: Works with Claude and Gemini models
- 🔄 **Automated Workflows**: Processes GitHub Issues and PRs automatically
- 🐳 **Flexible Deployment**: Docker containers or local CLI execution
- 📁 **Smart Workspace**: Git worktree-based isolation and management
- 🔒 **Secure**: GitHub webhook signature verification
- ⚡ **Fast Setup**: Quick configuration with multiple options

## 🚀 Quick Start

### 1. Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 2. Configuration

Create `config.yaml` with your preferred AI provider:

```yaml
# Choose your AI provider
code_provider: claude  # or gemini
use_docker: false      # true for Docker, false for CLI

server:
  port: 8888

workspace:
  base_dir: "./codeagent"
  cleanup_after: "24h"
```

### 3. Set Environment Variables

```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"
```

### 4. Start the Server

```bash
# Using the convenient startup script
./scripts/start.sh -p claude          # Claude + CLI (recommended for development)
./scripts/start.sh -p gemini -d      # Gemini + Docker

# Or directly with Go
go run ./cmd/server --config config.yaml
```

### 5. Configure GitHub Webhook

In your GitHub repository settings:
- **URL**: `https://your-domain.com/hook`
- **Content type**: `application/json`
- **Secret**: Your `WEBHOOK_SECRET` value
- **Events**: `Issue comments`, `Pull request reviews`, `Pull requests`

## 💬 Usage

Interact with CodeAgent using simple commands in GitHub:

### In Issues
```
/code Implement user authentication with JWT tokens
```

### In Pull Requests
```
/continue Add comprehensive unit tests
/fix Resolve memory leak in connection pool
```

## ⚙️ Configuration Options

### AI Providers

| Provider | Docker Image | CLI Tool | Notes |
|----------|-------------|----------|-------|
| Claude | `anthropic/claude-code:latest` | `claude` | Recommended for complex tasks |
| Gemini | `google-gemini/gemini-cli:latest` | `gemini` | Fast development cycles |

### Execution Modes

- **Docker Mode** (`use_docker: true`): Containerized execution, better for production
- **CLI Mode** (`use_docker: false`): Direct CLI execution, faster for development

### Security Configuration

Enable webhook signature verification for production:

```yaml
# In config.yaml - secrets should be in environment variables
webhook_secret: # Set via WEBHOOK_SECRET env var
```

**Security Best Practices:**
- Use strong webhook secrets (32+ characters)
- Always use HTTPS in production
- Regularly rotate API keys
- Limit GitHub token permissions

## 🏗️ Architecture

```
GitHub Events → Webhook → CodeAgent → Workspace → AI Provider → Code Generation → PR Updates
```

### Core Components

- **Agent** (`internal/agent/`): Orchestrates workflows
- **Webhook Handler** (`internal/webhook/`): Processes GitHub events
- **Workspace Manager** (`internal/workspace/`): Manages Git worktrees
- **Code Providers** (`internal/code/`): Interfaces with AI services
- **GitHub Client** (`internal/github/`): Handles GitHub API interactions

## 🛠️ Development

### Project Structure

```
codeagent/
├── cmd/server/          # Application entry point
├── internal/            # Core business logic
│   ├── agent/          # Workflow orchestration
│   ├── webhook/        # GitHub webhook handling
│   ├── workspace/      # Git workspace management
│   ├── code/           # AI provider implementations
│   └── github/         # GitHub API client
├── pkg/models/         # Shared data structures
├── scripts/            # Utility scripts
└── docs/               # Documentation
```

### Development Commands

```bash
# Build
make build

# Test
make test

# Run with different configurations
./scripts/start.sh --help

# Health check
curl http://localhost:8888/health
```

### Testing

```bash
# Integration test
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

## 📖 Documentation

- [Contributing Guide](CONTRIBUTING.md)
- [Architecture Details](docs/)
- [Workspace Management](internal/workspace/README.md)

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Feature Requests](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)  
- 📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thank you to all developers and contributors who have made this project possible!