# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered code generation system that automatically processes GitHub Issues and Pull Requests, providing intelligent code modification suggestions through webhook integration.

## ✨ Features

- 🤖 **Multi-AI Support**: Works with Claude and Gemini models
- 🔄 **GitHub Integration**: Automatic processing of Issues and Pull Requests  
- 🐳 **Flexible Deployment**: Docker containers or local CLI execution
- 📁 **Smart Workspace**: Git worktree-based workspace management
- 🔒 **Security**: Webhook signature verification and secure token handling

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
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"
```

2. **Start the Server**
```bash
go run ./cmd/server --port 8888
```

3. **Configure GitHub Webhook**
   - URL: `https://your-domain.com/hook`
   - Events: `Issue comments`, `Pull request reviews`, `Pull requests`
   - Content type: `application/json`
   - Secret: Same as your `WEBHOOK_SECRET`

### Usage

Use these commands in GitHub Issues or Pull Request comments:

```bash
# Generate code for an Issue
/code Implement user authentication with JWT tokens

# Continue development in PR
/continue Add unit tests for the login function

# Fix issues in PR
/fix Handle edge case for empty username
```

## ⚙️ Configuration

### Configuration File

Create `config.yaml` for advanced configuration:

```yaml
# Basic settings
server:
  port: 8888

# AI Provider (claude or gemini)
code_provider: claude
use_docker: false  # true for Docker, false for local CLI

# Workspace settings
workspace:
  base_dir: "./workspace"  # Supports relative paths
  cleanup_after: "24h"

# Docker settings (when use_docker: true)
docker:
  socket: "unix:///var/run/docker.sock"
  network: "bridge"

# Provider-specific settings
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"
```

**Note**: Sensitive data (tokens, secrets) should be set via environment variables, not in config files.

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `code_provider` | AI provider: `claude` or `gemini` | `gemini` |
| `use_docker` | Use Docker containers vs local CLI | `false` |
| `workspace.base_dir` | Working directory for code generation | `./workspace` |
| `workspace.cleanup_after` | Cleanup interval for temporary files | `24h` |

### Security Configuration

CodeAgent supports GitHub webhook signature verification:

```bash
# Set a strong webhook secret (32+ characters recommended)
export WEBHOOK_SECRET="your-strong-secret-here"
```

**Security Recommendations:**
- Always use HTTPS in production
- Use strong webhook secrets (32+ characters)
- Regularly rotate API keys and secrets
- Limit GitHub token permissions to required scopes

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

### Build and Test

```bash
# Build binary
make build

# Run tests
make test

# Health check
curl http://localhost:8888/health

# Test webhook (with sample data)
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

### Development Modes

**Local CLI Mode (Recommended for Development)**
- Faster startup and execution
- Requires `claude` or `gemini` CLI tools installed
- Direct API communication

**Docker Mode (Recommended for Production)**
- Isolated execution environment
- Complete toolkit included
- Better security and reproducibility

### Debugging

```bash
# Enable debug logging
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

## 📖 Advanced Usage

### Supported Commands

| Command | Context | Description |
|---------|---------|-------------|
| `/code <description>` | Issue comments | Generate initial code and create PR |
| `/continue <instruction>` | PR comments | Continue development with custom instructions |
| `/fix <description>` | PR comments | Fix specific issues in the code |

### Workflow

1. **Issue Processing**: User comments `/code` in GitHub Issue → CodeAgent creates branch and generates code → Submits Pull Request
2. **PR Collaboration**: User comments `/continue` or `/fix` in PR → CodeAgent modifies code in existing branch → Updates Pull Request
3. **Review Integration**: CodeAgent processes batch review comments and responds with comprehensive code updates

## 🤝 Contributing

We welcome contributions! Here's how to get involved:

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Feature Requests](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [Submit Code](CONTRIBUTING.md#code-contributions)

Please check the [Contributing Guide](CONTRIBUTING.md) for detailed information.

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thank you to all developers and users who have contributed to this project!