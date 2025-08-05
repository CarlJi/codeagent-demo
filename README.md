# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered code agent that automatically processes GitHub Issues and Pull Requests, generating intelligent code modification suggestions through webhook integration.

## 📋 Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage](#usage)
- [Development](#development)
- [Security](#security)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration
- 🔄 **Automatic Processing**: GitHub Issues and Pull Requests automation
- 🐳 **Flexible Deployment**: Docker containers or local CLI execution
- 📁 **Smart Workspace**: Git Worktree-based isolated environments
- 🔒 **Security First**: Webhook signature verification and secure configuration
- ⚡ **High Performance**: Efficient workspace management with automatic cleanup

## 🚀 Quick Start

### Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### Basic Setup

1. **Set Environment Variables**:
```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"
```

2. **Start the Server**:
```bash
go run ./cmd/server --port 8888
```

3. **Verify Installation**:
```bash
curl http://localhost:8888/health
```

### GitHub Webhook Setup

Configure webhook in your repository settings:
- **URL**: `https://your-domain.com/hook`
- **Content Type**: `application/json`
- **Secret**: Same as your `WEBHOOK_SECRET`
- **Events**: `Issue comments`, `Pull request reviews`, `Pull requests`

## ⚙️ Configuration

### Configuration Methods

CodeAgent supports three configuration approaches:

#### 1. Environment Variables (Simplest)
```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY
export WEBHOOK_SECRET="your-webhook-secret"
export CODE_PROVIDER=claude  # or gemini
export USE_DOCKER=false      # or true
export PORT=8888

go run ./cmd/server
```

#### 2. Command Line Arguments
```bash
go run ./cmd/server \
  --github-token "your-github-token" \
  --claude-api-key "your-claude-api-key" \
  --webhook-secret "your-webhook-secret" \
  --port 8888
```

#### 3. Configuration File (Production)
Create `config.yaml`:
```yaml
server:
  port: 8888

workspace:
  base_dir: "./codeagent"  # Supports relative paths
  cleanup_after: "24h"

# Provider configuration
code_provider: claude  # Options: claude, gemini
use_docker: true       # false for local CLI


```

Then run:
```bash
go run ./cmd/server --config config.yaml
```

### Configuration Options

| Option | Values | Description |
|--------|--------|-------------|
| `code_provider` | `claude`, `gemini` | AI service provider |
| `use_docker` | `true`, `false` | Execution environment |
| `workspace.base_dir` | Path string | Workspace directory (supports relative paths) |
| `workspace.cleanup_after` | Duration | Auto-cleanup interval |

**Security Note**: Never store sensitive information (tokens, API keys, secrets) in configuration files. Use environment variables or command-line arguments.

## 📖 Usage

### AI Commands

Trigger CodeAgent with these commands in GitHub:

#### Issue Commands
```bash
/code Implement user authentication with JWT tokens
```

#### Pull Request Commands
```bash
/continue Add comprehensive unit tests
```

### Execution Modes

#### CLI Mode (Development)
- Faster execution
- Requires local CLI installation (`claude` or `gemini`)
- Direct system integration

#### Docker Mode (Production)
- Isolated execution environment
- Consistent runtime across systems
- Enhanced security

## 🛠️ Development

### Project Structure
```
codeagent/
├── cmd/server/           # Application entry point
├── internal/
│   ├── agent/           # Core orchestration logic
│   ├── webhook/         # GitHub webhook handler
│   ├── workspace/       # Git worktree management
│   ├── code/           # AI provider implementations
│   ├── github/         # GitHub API client
│   └── config/         # Configuration management
├── pkg/models/         # Shared data structures
├── scripts/           # Utility scripts
└── docs/             # Documentation
```

### Build Commands
```bash
# Development build
make build

# Cross-platform build
GOOS=linux GOARCH=amd64 go build -o bin/codeagent-linux ./cmd/server

# Run tests
make test
```

### Development Workflow

1. **Local Testing**:
```bash
# Start development server
go run ./cmd/server --port 8888

# Test webhook endpoint
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

2. **Debug Mode**:
```bash
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

### Prerequisites

**For CLI Mode**:
- Claude CLI: Install from [Anthropic CLI](https://docs.anthropic.com/claude/docs/cli)
- Gemini CLI: Install from [Google AI CLI](https://ai.google.dev/gemini-api)

**For Docker Mode**:
- Docker Engine running
- Appropriate container images pulled

## 🔒 Security

### Webhook Security

CodeAgent implements GitHub webhook signature verification:

- **SHA-256 signature verification** (primary)
- **SHA-1 backward compatibility**
- **Constant-time comparison** (prevents timing attacks)
- **Environment-based secret management**

### Security Best Practices

✅ **Recommended**:
- Use strong webhook secrets (32+ characters)
- Always configure secrets in production
- Use HTTPS for webhook endpoints
- Regularly rotate API keys and secrets
- Limit GitHub token permissions to minimum required

❌ **Avoid**:
- Storing secrets in configuration files
- Using weak or default passwords
- Exposing webhook endpoints without verification
- Overly broad GitHub token permissions

### Token Permissions

Minimum required GitHub token scopes:
- `repo` - Repository access for code operations
- `pull_requests:write` - PR creation and updates
- `issues:write` - Issue comment responses

## 🤝 Contributing

We welcome contributions! Here's how to get involved:

### Quick Contribution
- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Request Features](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [Submit Code](CONTRIBUTING.md#code-contributions)

### Development Setup
1. Fork and clone the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

See [Contributing Guide](CONTRIBUTING.md) for detailed instructions.

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thank you to all contributors and the open-source community that makes this project possible!

---

**Need Help?** Check our [documentation](docs/) or [open an issue](https://github.com/qiniu/codeagent/issues/new).