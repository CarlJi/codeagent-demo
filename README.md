# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered code agent that automatically processes GitHub Issues and Pull Requests, generating intelligent code modifications and suggestions through webhook-driven automation.

## 🚀 Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration with Docker/CLI execution modes
- 🔄 **Automated Workflow**: Real-time processing of GitHub Issues and Pull Requests
- 🐳 **Flexible Deployment**: Docker containerized or local CLI execution
- 📁 **Smart Workspace**: Git Worktree-based isolated workspace management
- 🔒 **Enterprise Security**: Webhook signature verification and secure token handling

## 📋 Table of Contents

- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage](#usage)
- [Development](#development)
- [Security](#security)
- [Contributing](#contributing)
- [License](#license)

## 🚀 Quick Start

### Prerequisites

- Go 1.19+ installed
- GitHub personal access token
- API key for Claude or Gemini
- Docker (for container mode) or CLI tools (for local mode)

### Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### Quick Run

The fastest way to get started:

```bash
# Set required environment variables
export GITHUB_TOKEN="your-github-token"
export GOOGLE_API_KEY="your-google-api-key"  # or CLAUDE_API_KEY
export WEBHOOK_SECRET="your-webhook-secret"

# Start with default settings (Gemini + CLI mode)
./scripts/start.sh

# Or choose your preferred combination
./scripts/start.sh -p claude -d       # Claude + Docker mode
./scripts/start.sh -p gemini -d       # Gemini + Docker mode  
./scripts/start.sh -p claude          # Claude + CLI mode
```

### Health Check

```bash
curl http://localhost:8888/health
```

## ⚙️ Configuration

### Environment Variables

**Required:**
- `GITHUB_TOKEN` - GitHub personal access token with repo permissions
- `WEBHOOK_SECRET` - Secret for webhook signature verification
- `CLAUDE_API_KEY` or `GOOGLE_API_KEY` - API key for your chosen AI provider

**Optional:**
- `CODE_PROVIDER` - AI provider: `claude` or `gemini` (default: `gemini`)
- `USE_DOCKER` - Execution mode: `true` for Docker, `false` for CLI (default: `false`)
- `PORT` - Server port (default: `8888`)

### Configuration File

Create `config.yaml` for advanced configuration:

```yaml
# Core settings
code_provider: claude  # claude or gemini
use_docker: true      # true for Docker, false for CLI

server:
  port: 8888

workspace:
  base_dir: "./codeagent"  # Supports relative paths
  cleanup_after: "24h"

# AI provider settings
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"

docker:
  socket: "unix:///var/run/docker.sock"
  network: "bridge"
```

**Security Note**: Sensitive values (tokens, API keys, secrets) should be set via environment variables or command line flags, not in configuration files.

### Command Line Usage

```bash
# Using config file
go run ./cmd/server --config config.yaml

# Direct command line
go run ./cmd/server \
  --github-token "your-token" \
  --claude-api-key "your-key" \
  --webhook-secret "your-secret" \
  --port 8888
```

## 📖 Usage

### GitHub Webhook Setup

1. **Configure webhook in your GitHub repository:**
   - URL: `https://your-domain.com/hook`
   - Content type: `application/json`
   - Secret: Same as your `WEBHOOK_SECRET`
   - Events: `Issue comments`, `Pull request reviews`, `Pull requests`

### Supported Commands

**In GitHub Issues:**
```
/code Implement user login functionality with JWT authentication
```

**In Pull Request Comments:**
```
/continue Add comprehensive unit tests for the login module
/fix Resolve the null pointer exception in user validation
```

**In Pull Request Reviews:**
```
/continue Optimize database queries for better performance
```

### Execution Modes

**Docker Mode** (Recommended for Production):
- Isolated execution environment
- Complete AI tooling included
- Better security and consistency

**CLI Mode** (Recommended for Development):
- Faster execution
- Requires local AI CLI installation
- Direct system access

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
├── pkg/models/         # Shared data structures
├── scripts/           # Utility scripts
└── docs/             # Documentation
```

### Building

```bash
# Build for current platform
make build

# Cross-compilation example
GOOS=linux GOARCH=amd64 go build -o bin/codeagent-linux ./cmd/server
```

### Testing

```bash
# Run tests
make test

# Integration testing
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

### Debugging

```bash
# Enable debug logging
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

## 🔒 Security

### Webhook Security

CodeAgent implements robust webhook security:

- **SHA-256 signature verification** (with SHA-1 fallback)
- **Constant-time comparison** to prevent timing attacks
- **Required webhook secrets** in production environments

### Security Best Practices

- Use strong webhook secrets (32+ characters recommended)
- Deploy with HTTPS endpoints
- Regularly rotate API keys and secrets
- Limit GitHub token permissions to minimum required scope
- Never commit sensitive information to repositories

### Relative Path Support

CodeAgent safely handles relative paths in configuration:

```yaml
workspace:
  base_dir: "./codeagent"     # Relative to config file
  base_dir: "../workspace"    # Relative to parent directory
  base_dir: "/tmp/codeagent"  # Absolute path
```

## 🤝 Contributing

We welcome contributions! Here's how you can help:

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Request Features](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [Submit Code](CONTRIBUTING.md)

Please read our [Contributing Guide](CONTRIBUTING.md) for development setup and guidelines.

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thanks to all the developers and users who have contributed to making CodeAgent better!

---

For detailed documentation, visit our [docs](docs/) directory or check out specific topics:
- [Relative Path Support](docs/relative-path-support.md)
- [Architecture Design](docs/xgo-agent.md)