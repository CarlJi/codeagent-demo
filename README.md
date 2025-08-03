# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered code generation system that automatically processes GitHub Issues and Pull Requests through webhooks, supporting Claude and Gemini models with Docker and CLI execution modes.

## Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration
- 🔄 **GitHub Integration**: Automatic Issue/PR processing via webhooks
- 🐳 **Flexible Execution**: Docker containers or local CLI
- 📁 **Workspace Management**: Git worktree-based isolation
- 🔒 **Security**: Webhook signature verification

## Quick Start

### 1. Installation & Setup

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 2. Configuration

Set required environment variables:

```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"
```

### 3. Run with Startup Script (Recommended)

```bash
# Gemini + CLI mode (fastest for development)
./scripts/start.sh

# Claude + Docker mode
./scripts/start.sh -p claude -d

# Other combinations
./scripts/start.sh -p gemini -d    # Gemini + Docker
./scripts/start.sh -p claude       # Claude + CLI
```

### 4. Manual Configuration (Alternative)

Create `config.yaml`:

```yaml
code_provider: claude    # or gemini
use_docker: false       # true for Docker, false for CLI

server:
  port: 8888

workspace:
  base_dir: "./workspace"
  cleanup_after: "24h"
```

Then run:
```bash
go run ./cmd/server --config config.yaml
```

### 5. GitHub Webhook Setup

Configure in your GitHub repository:
- **URL**: `https://your-domain.com/hook`
- **Content type**: `application/json`
- **Secret**: Same as your `WEBHOOK_SECRET`
- **Events**: Issue comments, Pull request reviews, Pull requests

## Usage

Trigger CodeAgent in GitHub comments:

```bash
# Generate code for an Issue
/code Implement user authentication with JWT

# Continue development in PR
/continue Add unit tests for login function

# Fix specific issues
/fix Handle edge case for empty username
```

## Development

### Project Structure

```
codeagent/
├── cmd/server/          # Application entry point
├── internal/
│   ├── agent/          # Core orchestration logic
│   ├── webhook/        # GitHub webhook handling
│   ├── workspace/      # Git worktree management
│   ├── code/          # AI provider implementations
│   └── github/        # GitHub API client
├── pkg/models/         # Shared data structures
└── scripts/           # Utility scripts
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
```

### Configuration Options

**Provider Selection:**
- `claude`: Anthropic Claude (requires `CLAUDE_API_KEY`)
- `gemini`: Google Gemini (requires `GOOGLE_API_KEY`)

**Execution Modes:**
- Docker: Full containerized environment (production)
- CLI: Local tools (faster development)

**Security Features:**
- SHA-256/SHA-1 webhook signature verification
- Constant-time comparison against timing attacks
- Configurable secrets via environment variables

## 🤝 Contributing

We welcome all forms of contributions! Please check the [Contributing Guide](CONTRIBUTING.md) to learn how to participate in project development.

### Ways to Contribute

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Feature Requests](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [Submit Code](CONTRIBUTING.md#code-contributions)

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thank you to all developers and users who have contributed to this project!
