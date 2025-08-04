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

3. **Setup GitHub Webhook**
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


## License

This project is licensed under the [MIT License](LICENSE).
