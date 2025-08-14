# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered code generation and collaboration system that automatically processes GitHub Issues and Pull Requests through webhooks, providing intelligent code generation, modification, and review capabilities.

## ✨ Features

- 🤖 **Multi-AI Support**: Supports both Claude and Gemini AI models
- 🔄 **Webhook-Driven**: Automatic processing of GitHub Issues and Pull Requests
- 🐳 **Flexible Deployment**: Both Docker containerized and local CLI execution modes
- 📁 **Workspace Management**: Git Worktree-based temporary workspace management with auto-cleanup
- 🔒 **Security**: Built-in webhook signature verification and secure token management
- 🚀 **Easy Setup**: Simple configuration with multiple setup options

## 🚀 Quick Start

### 1. Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 2. Configuration

Create a `config.yaml` file with your preferred settings:

```yaml
# Basic server configuration
server:
  port: 8888

# Workspace settings
workspace:
  base_dir: "./codeagent"  # Supports both relative and absolute paths
  cleanup_after: "24h"

# AI provider settings
code_provider: claude  # Options: claude, gemini
use_docker: false      # true for Docker, false for local CLI

# Docker settings (when use_docker: true)
docker:
  socket: "unix:///var/run/docker.sock"
  network: "bridge"

# Claude configuration
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

# Gemini configuration
gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"
```

### 3. Set Environment Variables

```bash
# Required for all configurations
export GITHUB_TOKEN="your-github-token"
export WEBHOOK_SECRET="your-webhook-secret"

# Choose one based on your AI provider
export CLAUDE_API_KEY="your-claude-api-key"    # For Claude
export GOOGLE_API_KEY="your-google-api-key"    # For Gemini
```

### 4. Run CodeAgent

**Option A: Using the convenient startup script (Recommended)**
```bash
./scripts/start.sh                    # Gemini + CLI (default)
./scripts/start.sh -p claude -d       # Claude + Docker
./scripts/start.sh -p gemini -d       # Gemini + Docker  
./scripts/start.sh -p claude          # Claude + CLI
```

**Option B: Direct Go execution**
```bash
go run ./cmd/server --config config.yaml
```

**Option C: Build and run binary**
```bash
make build
./bin/codeagent --config config.yaml
```

### 5. Configure GitHub Webhook

1. Go to your GitHub repository settings
2. Add a new webhook with:
   - **URL**: `https://your-domain.com/hook`
   - **Content type**: `application/json`
   - **Secret**: Same value as your `WEBHOOK_SECRET`
   - **Events**: Select `Issue comments`, `Pull request reviews`, and `Pull requests`

### 6. Test Your Setup

```bash
# Health check
curl http://localhost:8888/health

# Test webhook (optional)
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

## 💡 Usage

CodeAgent responds to specific commands in GitHub Issues and Pull Requests:

### Issue Commands
```bash
/code <description>  # Generate initial code implementation and create a PR
```

### Pull Request Commands  
```bash
/continue <instruction>  # Continue development with custom instructions
/fix <description>       # Fix specific issues in the code
```

### Examples
```bash
# In a GitHub Issue
/code Implement user authentication with JWT tokens and password hashing

# In a PR comment
/continue Add comprehensive unit tests for all functions

# In a PR review
/fix Handle edge case when username is empty
```

## 🔧 Configuration Options

### AI Providers

**Claude (Anthropic)**
- **Docker mode**: Uses containerized Claude Code with full toolkit
- **CLI mode**: Uses locally installed Claude CLI (faster for development)

**Gemini (Google)**
- **Docker mode**: Uses containerized Gemini CLI 
- **CLI mode**: Uses locally installed Gemini CLI with optimized single-prompt approach

### Execution Modes

| Mode | Pros | Cons | Best For |
|------|------|------|----------|
| **CLI** | Faster startup, lower resource usage | Requires local AI CLI installation | Development |
| **Docker** | Isolated environment, no local dependencies | Slower startup, higher resource usage | Production |

### Security Configuration

CodeAgent includes robust security features:

- **Webhook Signature Verification**: SHA-256/SHA-1 signature validation with constant-time comparison
- **Environment-based Secrets**: API keys and tokens via environment variables (never in config files)
- **Token Scope Limiting**: Minimize GitHub token permissions
- **HTTPS Enforcement**: Recommended for production deployments

## 📁 Project Structure

```
codeagent/
├── cmd/server/          # Application entry point
├── internal/
│   ├── agent/          # Core orchestration logic  
│   ├── webhook/        # GitHub webhook handling
│   ├── workspace/      # Git worktree management
│   ├── code/           # AI provider implementations
│   ├── github/         # GitHub API client
│   └── config/         # Configuration management
├── pkg/models/         # Shared data structures
├── scripts/            # Utility scripts
├── docs/              # Documentation
└── config.yaml        # Configuration file
```

## 🛠️ Development

### Build Commands
```bash
make build                                    # Build binary
make test                                     # Run tests
GOOS=linux GOARCH=amd64 make build          # Cross-compile for Linux
```

### Debug Mode
```bash
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

### Testing
```bash
# Run all tests
go test ./...

# Integration testing
go run ./cmd/server --config test-config.yaml
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Ways to Contribute
- 🐛 [Report Issues](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Request Features](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)  
- 📝 [Improve Docs](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [Submit PRs](CONTRIBUTING.md#code-contributions)

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thanks to all contributors and users who make this project possible!