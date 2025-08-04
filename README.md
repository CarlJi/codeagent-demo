# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered code agent that automatically processes GitHub Issues and Pull Requests, generating and modifying code through multiple AI providers.

## Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration
- 🔄 **Auto Processing**: GitHub Issues and Pull Requests automation
- 🐳 **Flexible Deployment**: Docker containers or local CLI
- 📁 **Smart Workspace**: Git Worktree-based management

## Quick Start

1. **Install**
   ```bash
   git clone https://github.com/qiniu/codeagent.git
   cd codeagent
   go mod download
   ```

2. **Configure**
   ```bash
   export GITHUB_TOKEN="your-github-token"
   export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY
   export WEBHOOK_SECRET="your-webhook-secret"
   ```

3. **Run**
   ```bash
   # Quick start with script
   ./scripts/start.sh                # Gemini + CLI (default)
   ./scripts/start.sh -p claude      # Claude + CLI
   ./scripts/start.sh -p claude -d   # Claude + Docker
   
   # Or run directly
   go run ./cmd/server --port 8888
   ```

4. **Setup GitHub Webhook**
   - URL: `https://your-domain.com/hook`
   - Events: `Issue comments`, `Pull request reviews`, `Pull requests`
   - Secret: Same as `WEBHOOK_SECRET`

## Usage

Comment in GitHub Issues or PRs with:

- `/code <description>` - Generate code for an Issue
- `/continue <instruction>` - Continue development in PR
- `/fix <description>` - Fix code issues

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `GITHUB_TOKEN` | GitHub personal access token | ✅ |
| `WEBHOOK_SECRET` | GitHub webhook secret | ✅ |
| `CLAUDE_API_KEY` | Claude API key | ✅* |
| `GOOGLE_API_KEY` | Google API key | ✅* |
| `CODE_PROVIDER` | AI provider: `claude` or `gemini` | ❌ |
| `USE_DOCKER` | Use Docker: `true` or `false` | ❌ |

*One of the AI provider keys is required

### Configuration File (config.yaml)

```yaml
server:
  port: 8888

code_provider: claude  # or gemini
use_docker: false      # true for Docker, false for CLI

workspace:
  base_dir: "./codeagent"
  cleanup_after: "24h"

claude:
  container_image: "anthropic/claude-code:latest"
  timeout: "30m"

gemini:
  container_image: "google-gemini/gemini-cli:latest"
  timeout: "30m"
```

## Development

### Project Structure

```
codeagent/
├── cmd/server/           # Main entry point
├── internal/
│   ├── agent/           # Core orchestration
│   ├── webhook/         # GitHub webhook handling
│   ├── workspace/       # Git workspace management
│   ├── code/           # AI provider implementations
│   └── github/         # GitHub API client
├── pkg/models/         # Data structures
└── scripts/           # Utility scripts
```

### Build & Test

```bash
# Build
make build

# Test
make test
curl http://localhost:8888/health

# Debug
export LOG_LEVEL=debug
go run ./cmd/server
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
