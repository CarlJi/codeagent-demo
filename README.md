# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

**CodeAgent** is an AI-powered automated code generation system that processes GitHub Issues and Pull Requests through webhooks, supporting Claude and Gemini models with both Docker and CLI execution modes.

## ✨ Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration
- 🔄 **GitHub Integration**: Automatic Issue/PR processing via webhooks  
- 🐳 **Flexible Execution**: Docker containers or local CLI
- 📁 **Smart Workspace**: Git worktree-based project management
- 🔐 **Security**: Webhook signature verification and secure token handling

## 🚀 Quick Start

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
code_provider: claude    # or gemini
use_docker: true        # or false for CLI mode
server:
  port: 8888
workspace:
  base_dir: "./codeagent"
  cleanup_after: "24h"
```

### Start the Service

```bash
go run ./cmd/server --config config.yaml
```

### GitHub Webhook Setup

1. Go to your repository → Settings → Webhooks
2. Add webhook:
   - **URL**: `https://your-domain.com/hook`
   - **Content type**: `application/json`
   - **Secret**: Same as your `WEBHOOK_SECRET`
   - **Events**: Issue comments, Pull request reviews, Pull requests

### Usage

Trigger AI actions with comments in GitHub Issues or PRs:

```bash
/code Implement user authentication with JWT tokens
/continue Add unit tests for the login function
/fix Resolve the null pointer exception in handleLogin
```

## 📖 Documentation

### Configuration Options

| Provider | Mode | Description |
|----------|------|-------------|
| Claude | Docker | Use Claude with Docker container |
| Claude | CLI | Use Claude with local CLI |
| Gemini | Docker | Use Gemini with Docker container |
| Gemini | CLI | Use Gemini with local CLI |

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GITHUB_TOKEN` | ✅ | GitHub personal access token |
| `CLAUDE_API_KEY` | ✅* | Claude API key (* if using Claude) |
| `GOOGLE_API_KEY` | ✅* | Google API key (* if using Gemini) |
| `WEBHOOK_SECRET` | ✅ | GitHub webhook secret |
| `CODE_PROVIDER` | ❌ | Override config file (claude/gemini) |
| `USE_DOCKER` | ❌ | Override config file (true/false) |

### Project Structure

```
codeagent/
├── cmd/server/           # Application entry point
├── internal/
│   ├── agent/           # Core orchestration logic
│   ├── webhook/         # GitHub webhook handling
│   ├── workspace/       # Git worktree management
│   ├── code/           # AI provider implementations
│   └── github/         # GitHub API client
├── pkg/models/         # Shared data structures
├── scripts/            # Utility scripts
└── config.yaml         # Configuration file
```

## 🛠️ Development

### Build

```bash
make build                           # Build binary
make test                           # Run tests
go build -o bin/codeagent ./cmd/server
```

### Testing

```bash
# Health check
curl http://localhost:8888/health

# Test webhook
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

### Debug Mode

```bash
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

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
