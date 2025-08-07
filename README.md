# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

🤖 **CodeAgent** is an AI-powered code generation system that automatically processes GitHub Issues and Pull Requests, providing intelligent code modifications through multiple AI providers (Claude, Gemini).

## ✨ Features

- 🤖 **Multi-AI Support**: Claude and Gemini integration with Docker/CLI execution modes
- 🔄 **GitHub Integration**: Automatic processing of Issues and Pull Requests via webhooks  
- 📁 **Smart Workspace**: Git worktree-based isolated workspace management
- 🔒 **Security First**: Webhook signature verification and secure token handling
- 🚀 **Easy Setup**: Multiple configuration methods with intelligent defaults

## 🚀 Quick Start

### Prerequisites
- Go 1.19+
- Git
- Docker (optional, for containerized execution)
- GitHub token with repository access
- AI provider API key (Claude or Gemini)

### Installation & Setup

1. **Clone and install dependencies**:
   ```bash
   git clone https://github.com/qiniu/codeagent.git
   cd codeagent
   go mod download
   ```

2. **Set environment variables**:
   ```bash
   export GITHUB_TOKEN="your-github-token"
   export CLAUDE_API_KEY="your-claude-api-key"  # or GOOGLE_API_KEY for Gemini
   export WEBHOOK_SECRET="your-webhook-secret"
   ```

3. **Start the server**:
   ```bash
   go run ./cmd/server --port 8888
   ```

4. **Configure GitHub webhook**:
   - URL: `https://your-domain.com/hook`
   - Events: `Issue comments`, `Pull request reviews`, `Pull requests`
   - Secret: Same value as your `WEBHOOK_SECRET`
   - Content type: `application/json`

5. **Test with GitHub commands**:
   ```
   /code Implement user authentication with JWT
   /continue Add error handling
   /fix Fix the validation logic bug
   ```

## ⚙️ Configuration

CodeAgent supports flexible configuration through environment variables, command line arguments, or YAML files.

### Environment Variables
```bash
# Required
export GITHUB_TOKEN="your-github-token"
export WEBHOOK_SECRET="your-webhook-secret"

# AI Provider (choose one)
export CLAUDE_API_KEY="your-claude-api-key"
export GOOGLE_API_KEY="your-google-api-key"

# Optional
export CODE_PROVIDER="claude"  # or "gemini"
export USE_DOCKER="false"      # or "true"
export PORT="8888"
```

### Configuration File (config.yaml)
```yaml
server:
  port: 8888

workspace:
  base_dir: "./codeagent"  # Supports relative paths
  cleanup_after: "24h"

# AI Provider Configuration
code_provider: claude  # Options: claude, gemini
use_docker: false      # true for Docker, false for CLI

# Docker settings (when use_docker: true)
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

**Security Note**: Never store sensitive tokens in configuration files. Use environment variables or command line arguments.

### Execution Modes

| Mode | Provider | Method | Use Case |
|------|----------|---------|----------|
| `claude + cli` | Claude | Local CLI | Development (fastest) |
| `claude + docker` | Claude | Docker | Production (isolated) |
| `gemini + cli` | Gemini | Local CLI | Development (recommended) |
| `gemini + docker` | Gemini | Docker | Production (isolated) |

## 🔧 Development

### Project Structure
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
└── docs/              # Documentation
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

### Debugging
```bash
# Enable debug logging
export LOG_LEVEL=debug
go run ./cmd/server --config config.yaml
```

## 🔒 Security

### Webhook Security
- **Signature Verification**: SHA-256/SHA-1 signature verification with constant-time comparison
- **Strong Secrets**: Use 32+ character webhook secrets
- **HTTPS Only**: Always use HTTPS endpoints in production
- **Token Rotation**: Regularly rotate API keys and webhook secrets

### GitHub Permissions
Minimum required GitHub token permissions:
- Repository: Read/Write access
- Issues: Read/Write access  
- Pull requests: Read/Write access
- Contents: Read/Write access

## 📖 Usage Examples

### Issue Commands
```bash
# Generate code for new feature
/code Implement REST API for user management with CRUD operations

# Generate code with specific requirements  
/code Create a login system with bcrypt password hashing and JWT tokens
```

### Pull Request Commands
```bash
# Continue development
/continue Add input validation and error handling

# Fix specific issues
/fix Fix the memory leak in the connection pool

# Add tests
/continue Write comprehensive unit tests for the auth module
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### How to Contribute
- 🐛 [Report Issues](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
- 💡 [Request Features](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)  
- 📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
- 🔧 [Submit Code](CONTRIBUTING.md#code-contributions)

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

**Built with ❤️ by the Qiniu Team** | [Documentation](docs/) | [Issues](https://github.com/qiniu/codeagent/issues) | [Discussions](https://github.com/qiniu/codeagent/discussions)