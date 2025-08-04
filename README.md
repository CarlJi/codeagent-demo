# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered code agent that automatically processes GitHub Issues and Pull Requests, generating code modification suggestions with support for Claude and Gemini models.

## ✨ Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration with Docker/CLI modes
- 🔄 **Auto Processing**: GitHub Issues and Pull Requests automation
- 🐳 **Flexible Deployment**: Docker containers or local CLI execution
- 📁 **Smart Workspace**: Git Worktree-based management with auto-cleanup
- 🔐 **Secure**: Webhook signature verification and HTTPS support

## 🚀 Quick Start

### Installation & Setup

```bash
# Clone and install
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download

# Set required environment variables
export GITHUB_TOKEN="your-github-token"
export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY for Gemini
export WEBHOOK_SECRET="your-webhook-secret"

# Start with default settings (Gemini + CLI mode)
./scripts/start.sh

# Or use other configurations
./scripts/start.sh -p claude -d    # Claude + Docker mode
./scripts/start.sh -p gemini -d    # Gemini + Docker mode
./scripts/start.sh -p claude       # Claude + CLI mode
```

### Configuration

Create `config.yaml` for custom settings:

```yaml
# Core configuration
code_provider: claude    # Options: claude, gemini
use_docker: true        # true: Docker mode, false: CLI mode

server:
  port: 8888

workspace:
  base_dir: "./codeagent"    # Supports relative paths
  cleanup_after: "24h"

# AI provider settings
claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"
  
gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"
```

**Security Note**: Keep sensitive data (tokens, API keys, webhook secrets) in environment variables, not config files.

### GitHub Integration

1. **Add Webhook to your repository**:
   - URL: `https://your-domain.com/hook`
   - Content type: `application/json`
   - Secret: Same as your `WEBHOOK_SECRET`
   - Events: `Issue comments`, `Pull request reviews`, `Pull requests`

2. **Verify setup**:
   ```bash
   curl http://localhost:8888/health
   ```

## 💬 Usage

Use these commands in GitHub Issues and Pull Requests:

### Issue Commands
- `/code <description>` - Generate code and create PR
  ```
  /code Implement user authentication with JWT tokens
  ```

### Pull Request Commands  
- `/continue <instruction>` - Continue development
  ```
  /continue Add comprehensive unit tests
  ```
- `/fix <description>` - Fix issues
  ```
  /fix Handle edge case for empty user input
  ```

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

### Build & Test

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
