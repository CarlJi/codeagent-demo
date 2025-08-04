# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

AI-powered code agent that automatically processes GitHub Issues and Pull Requests with intelligent code generation and modification capabilities.

## Features

- 🤖 Multi-AI support (Claude & Gemini)
- 🔄 Automated GitHub Issue/PR processing
- 🐳 Docker & CLI execution modes
- 📁 Git worktree workspace management
- 🔒 Webhook signature verification

## Quick Start

### Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### Configuration

Set required environment variables:

```bash
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"
```

Create `config.yaml`:

```yaml
server:
  port: 8888

workspace:
  base_dir: "./codeagent"
  cleanup_after: "24h"

code_provider: claude  # Options: claude, gemini
use_docker: false      # true for Docker, false for CLI
```

### Running

```bash
# Using startup script (recommended)
./scripts/start.sh                    # Gemini + CLI mode (default)
./scripts/start.sh -p claude -d       # Claude + Docker mode
./scripts/start.sh -p gemini -d       # Gemini + Docker mode
./scripts/start.sh -p claude          # Claude + CLI mode

# Or run directly
go run ./cmd/server
```

### GitHub Webhook Setup

1. Go to your repository **Settings** → **Webhooks** → **Add webhook**
2. Set **Payload URL**: `https://your-domain.com/hook`
3. Set **Content type**: `application/json`
4. Set **Secret**: Same as your `WEBHOOK_SECRET`
5. Select events: `Issue comments`, `Pull request reviews`, `Pull requests`

### Usage

In GitHub Issues or PRs, use these commands in comments:

```
/code Implement user authentication with JWT tokens
/continue Add error handling and validation
/fix Fix the login validation logic bug
```

## Development

### Build & Test

```bash
# Build
go build -o bin/codeagent ./cmd/server

# Test
go test ./...

# Health check
curl http://localhost:8888/health

# Test webhook
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

### Project Structure

```
├── cmd/server/           # Main application
├── internal/
│   ├── agent/           # Core orchestration logic
│   ├── webhook/         # GitHub webhook handler
│   ├── workspace/       # Git worktree management
│   ├── code/           # AI provider implementations
│   └── github/         # GitHub API client
├── pkg/models/         # Shared data structures
└── scripts/           # Utility scripts
```

### Debugging

```bash
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

## Contributing

We welcome contributions! Please check the [Contributing Guide](CONTRIBUTING.md).

- 🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Feature Requests](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
- 🔧 [Submit Code](CONTRIBUTING.md#code-contributions)

## License

This project is licensed under the [MIT License](LICENSE).
