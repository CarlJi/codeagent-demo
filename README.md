# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

An AI-powered code agent that automatically processes GitHub Issues and Pull Requests, generating intelligent code modifications through webhook-driven automation.

## 📋 Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage](#usage)
- [Development](#development)
- [Security](#security)
- [Contributing](#contributing)
- [License](#license)

## Features

- 🤖 Support for multiple AI models (Claude, Gemini)
- 🔄 Automatic processing of GitHub Issues and Pull Requests
- 🐳 Docker containerized execution environment
- 📁 Git Worktree-based workspace management

## Quick Start

### 1. Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 2. Environment Setup

```bash
# Required environment variables
export GITHUB_TOKEN="your-github-token"
export WEBHOOK_SECRET="your-webhook-secret"
export CLAUDE_API_KEY="your-claude-api-key"  # OR GOOGLE_API_KEY for Gemini
```

### 3. Start the Server

```bash
# Quick start with defaults (Gemini + CLI mode)
./scripts/start.sh

# Or run directly
go run ./cmd/server --port 8888
```

### 4. Configure GitHub Webhook

- **URL**: `https://your-domain.com/hook`
- **Events**: `Issue comments`, `Pull request reviews`, `Pull requests`
- **Secret**: Use your `WEBHOOK_SECRET` value
- **Content type**: `application/json`

### 5. Test the Setup

```bash
# Health check
curl http://localhost:8888/health

# Test with sample webhook
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

## Configuration

### Configuration Methods

Choose one of the following configuration methods:

#### Environment Variables (Recommended for Development)

```bash
export GITHUB_TOKEN="your-github-token"
export WEBHOOK_SECRET="your-webhook-secret"
export CLAUDE_API_KEY="your-claude-api-key"  # OR GOOGLE_API_KEY
export CODE_PROVIDER=claude  # or gemini
export USE_DOCKER=false     # true for Docker mode

go run ./cmd/server --port 8888
```

#### Configuration File (Recommended for Production)

Create `config.yaml`:

```yaml
server:
  port: 8888

workspace:
  base_dir: "./codeagent"  # Supports relative paths
  cleanup_after: "24h"

# Provider selection
code_provider: claude  # Options: claude, gemini
use_docker: true      # Use Docker containers vs local CLI

# AI Provider settings
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

**Note**: Store sensitive data (tokens, API keys, secrets) in environment variables, not config files.

#### Startup Script (Simplest)

```bash
# Set required environment variables first
export GITHUB_TOKEN="your-github-token"
export WEBHOOK_SECRET="your-webhook-secret"
export CLAUDE_API_KEY="your-claude-api-key"  # OR GOOGLE_API_KEY

# Use startup script with different modes
./scripts/start.sh                    # Gemini + CLI (default)
./scripts/start.sh -p claude -d       # Claude + Docker
./scripts/start.sh -p gemini -d       # Gemini + Docker
./scripts/start.sh -p claude          # Claude + CLI
./scripts/start.sh --help             # View all options
```

### Provider & Execution Modes

| Provider | CLI Mode | Docker Mode | Best For |
|----------|----------|-------------|----------|
| Claude   | ✅ Fast | ✅ Isolated | Development / Production |
| Gemini   | ✅ Fast | ✅ Isolated | Development / Production |

- **CLI Mode**: Uses locally installed tools, faster startup
- **Docker Mode**: Containerized execution, better isolation

## Usage

### GitHub Commands

Use these commands in GitHub Issues and Pull Request comments:

#### Issue Commands
```
/code <description>
```
Generate initial code implementation and create a Pull Request.

**Example:**
```
/code Implement user authentication with JWT tokens and password hashing
```

#### Pull Request Commands
```
/continue <instruction>
/fix <description>
```

**Examples:**
```
/continue Add comprehensive unit tests for the authentication module
/fix Handle edge case when user email is null
```

### Webhook Events

CodeAgent responds to these GitHub webhook events:
- **Issue Comments**: Processes `/code` commands
- **PR Comments**: Processes `/continue` and `/fix` commands  
- **PR Reviews**: Processes batch review comments with instructions

## Development

### Project Structure

```
codeagent/
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── agent/agent.go           # Core orchestration logic
│   ├── webhook/handler.go       # GitHub webhook handler
│   ├── workspace/manager.go     # Git worktree management
│   ├── code/                    # AI provider implementations
│   ├── github/client.go         # GitHub API client
│   └── config/config.go         # Configuration management
├── pkg/models/                  # Shared data structures
├── scripts/start.sh             # Development startup script
└── docs/                        # Documentation
```

### Building & Testing

```bash
# Build binary
make build
# or
go build -o bin/codeagent ./cmd/server

# Run tests
make test

# Integration testing
go run ./cmd/server --config test-config.yaml

# Send test webhook
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

# Monitor workspace cleanup
ls -la /tmp/codeagent/  # Default workspace location
```

### Architecture

```
GitHub Events → Webhook → Agent → Workspace → AI Provider → Code Generation → PR Updates
```

**Key Components:**
- **Agent**: Orchestrates the entire workflow
- **Webhook Handler**: Processes GitHub events (Issues/PRs)
- **Workspace Manager**: Manages temporary Git worktrees
- **Code Providers**: Claude/Gemini integration (Docker/CLI)
- **GitHub Client**: Handles API interactions

## Security

### Webhook Signature Verification

CodeAgent supports GitHub webhook signature verification to prevent malicious requests:

```bash
# Configure webhook secret (required for production)
export WEBHOOK_SECRET="your-strong-secret-here"  # 32+ characters recommended
```

**GitHub Webhook Settings:**
- URL: `https://your-domain.com/hook`
- Content type: `application/json`
- Secret: Use your `WEBHOOK_SECRET` value
- Events: `Issue comments`, `Pull request reviews`, `Pull requests`

**Security Features:**
- SHA-256 signature verification (primary)
- SHA-1 backward compatibility
- Constant-time comparison (prevents timing attacks)
- Development mode: signature verification optional

### Security Best Practices

- ✅ Use strong webhook secrets (32+ characters)
- ✅ Always configure secrets in production
- ✅ Use HTTPS for webhook endpoints
- ✅ Regularly rotate API keys and secrets
- ✅ Limit GitHub token permissions to minimum required scope
- ✅ Store sensitive data in environment variables, not config files

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

**Quick Links:**
- 🐛 [Report Bug](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Request Feature](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [Improve Docs](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)

## License

This project is licensed under the [MIT License](LICENSE).
