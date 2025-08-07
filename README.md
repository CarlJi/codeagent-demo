# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered code agent that automatically processes GitHub Issues and Pull Requests, generating code modification suggestions.

## 📋 Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Configuration](#configuration)
- [Development](#development)
- [Contributing](#contributing)
- [Security](#security)
- [License](#license)
- [Acknowledgments](#acknowledgments)

## ✨ Features

- 🤖 **Multi-Provider Support**: Integrates with various AI models like Claude and Gemini.
- 🔄 **Automated Workflow**: Automatically handles GitHub Issues and Pull Requests to generate code.
- 🐳 **Containerized Environment**: Utilizes Docker for a consistent and isolated execution environment.
- 📁 **Workspace Management**: Manages workspaces efficiently using Git Worktree.

## 🚀 Quick Start

### 1. Installation

```bash
git clone https://github.com/qiniu/codeagent.git
cd codeagent
go mod download
```

### 2. Configuration

CodeAgent can be configured via command-line arguments, environment variables, or a configuration file. For this quick start, we'll use environment variables.

**Note**: Sensitive information (like tokens and API keys) should be set via command-line arguments or environment variables.

```bash
# GitHub token with repo access
export GITHUB_TOKEN="your-github-token"

# Webhook secret for securing your webhook endpoint
export WEBHOOK_SECRET="your-webhook-secret"

# Google API Key for Gemini
export GOOGLE_API_KEY="your-google-api-key"
# or Claude API Key
# export CLAUDE_API_KEY="your-claude-api-key"
```

### 3. Running the Agent

We provide a convenient startup script that supports all configuration combinations.

```bash
# Start with Gemini in Local CLI mode (default)
./scripts/start.sh

# Start with Claude in Docker mode
# ./scripts/start.sh -p claude -d

# View help for more options
./scripts/start.sh --help
```

### 4. Configure GitHub Webhook

1.  In your GitHub repository settings, go to **Webhooks** and click **Add webhook**.
2.  **Payload URL**: `http://<your-server-address>:8888/hook`
3.  **Content type**: `application/json`
4.  **Secret**: Enter the same value as your `WEBHOOK_SECRET`.
5.  **Which events would you like to trigger this webhook?**: Select `Issue comments`, `Pull request review comments`, and `Pull requests`.
6.  Click **Add webhook**.

## 📖 Usage

You can interact with CodeAgent by commenting on GitHub Issues and Pull Requests.

-   **Start a task**:
    ```
    /code Implement user login functionality.
    ```
-   **Continue a task**:
    ```
    /continue Add unit tests for the login functionality.
    ```
-   **Fix code**:
    ```
    /fix The login validation is not working correctly.
    ```

## ⚙️ Configuration

CodeAgent offers flexible configuration options to suit your needs.

### Configuration Methods

1.  **Command-line arguments (highest priority)**:
    ```bash
    go run ./cmd/server --github-token "your-token" --claude-api-key "your-key"
    ```
2.  **Environment variables**:
    ```bash
    export GITHUB_TOKEN="your-token"
    export CLAUDE_API_KEY="your-key"
    go run ./cmd/server
    ```
3.  **Configuration file (lowest priority)**:
    Create a `config.yaml` file (see `config.example.yaml`) and run:
    ```bash
    go run ./cmd/server --config config.yaml
    ```

### Key Configuration Options

-   `code_provider`: `claude` or `gemini`.
-   `use_docker`: `true` to run in a Docker container (recommended for production), `false` to run locally (recommended for development).

### Configuration Examples

-   **Gemini + Local CLI (Default for `start.sh`)**:
    `CODE_PROVIDER=gemini`, `USE_DOCKER=false`
-   **Claude + Docker**:
    `CODE_PROVIDER=claude`, `USE_DOCKER=true`

## 🛠️ Development

### Project Structure
```
codeagent/
├── cmd/server/main.go      # Main program entry point
├── internal/               # Internal application logic
│   ├── agent/              # Agent core logic
│   ├── code/               # Code generation providers (Claude, Gemini)
│   ├── config/             # Configuration management
│   ├── github/             # GitHub API client
│   ├── webhook/            # Webhook handler
│   └── workspace/          # Workspace management
├── pkg/                    # Shared packages
│   └── models/             # Data models
├── scripts/                # Helper scripts
├── Dockerfile              # Docker build file
├── go.mod                  # Go module file
└── README.md               # This file
```

### Building

```bash
# Build for your local OS
go build -o bin/codeagent ./cmd/server

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o bin/codeagent-linux ./cmd/server
```

### Testing

```bash
# Run unit tests
go test ./...

# Run integration tests
# 1. Start the server with a test configuration
go run ./cmd/server --config test-config.yaml
# 2. Send a test webhook event
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

## 🤝 Contributing

We welcome all forms of contributions! Please check the [Contributing Guide](CONTRIBUTING.md) to learn how to participate in project development.

-   🐛 [Report Bugs](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
-   💡 [Suggest Features](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
-   📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)
-   🔧 [Submit Pull Requests](CONTRIBUTING.md#code-contributions)

## 🔒 Security

We take security seriously. Please review our security guidelines.

### Webhook Signature Verification

To protect against malicious requests, CodeAgent verifies webhook signatures from GitHub. Ensure you set a `WEBHOOK_SECRET` in your configuration. This feature uses SHA-256 (and falls back to SHA-1) signature verification and constant-time comparison to prevent timing attacks.

### Security Recommendations

-   Use a strong, unique `WEBHOOK_SECRET`.
-   Use HTTPS for your webhook endpoint.
-   Regularly rotate API keys and tokens.
-   Grant the `GITHUB_TOKEN` only the necessary permissions.

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Thank you to all developers and users who have contributed to this project!