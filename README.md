# CodeAgent

[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/codeagent)](https://goreportcard.com/report/github.com/qiniu/codeagent)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qiniu/codeagent)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/qiniu/codeagent/workflows/CI/badge.svg)](https://github.com/qiniu/codeagent/actions)

CodeAgent is an AI-powered agent that automates the handling of GitHub Issues and Pull Requests by generating code modification suggestions.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running the Agent](#running-the-agent)
- [Usage](#usage)
- [Development](#development)
  - [Project Structure](#project-structure)
  - [Building](#building)
  - [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Features

- 🤖 **Multi-Provider Support**: Integrates with various AI models, including Claude and Gemini.
- 🔄 **GitHub Automation**: Automatically processes GitHub Issues and Pull Requests based on comments.
- 🐳 **Containerized Environment**: Supports running AI models in a sandboxed Docker environment for security and consistency.
- 📁 **Isolated Workspaces**: Manages different tasks in separate Git worktrees to avoid conflicts.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) 1.18+
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Docker](https://docs.docker.com/get-docker/) (Optional, for running models in containers)

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/qiniu/codeagent.git
    cd codeagent
    ```

2.  Download Go modules:
    ```bash
    go mod download
    ```

### Configuration

CodeAgent can be configured via a YAML file, environment variables, or command-line flags.

1.  **Create a configuration file** by copying the example:
    ```bash
    cp config.example.yaml config.yaml
    ```

2.  **Set required credentials**. The agent needs a GitHub token and an API key for the desired AI provider (e.g., Claude or Gemini). For security, it's best to provide these using environment variables.

    ```bash
    export GITHUB_TOKEN="your-github-token"
    export CLAUDE_API_KEY="your-claude-api-key" # or GOOGLE_API_KEY for Gemini
    export WEBHOOK_SECRET="a-very-strong-secret"
    ```

3.  **Configure your `config.yaml`**. At a minimum, set your preferred `code_provider` (`claude` or `gemini`) and `use_docker` (`true` or `false`).

    ```yaml
    # config.yaml
    code_provider: claude # Or gemini
    use_docker: true      # Or false to use local CLI tools
    
    github:
      webhook_url: "http://your-public-url:8888/hook"
    
    workspace:
      base_dir: "./codeagent_work"
    ```

    **Note**: Relative paths like `./codeagent_work` are supported for `base_dir`.

### Running the Agent

You can run the server directly or use the provided start script.

**Using the start script (Recommended for local development):**

The script simplifies running the agent with different configurations.

```bash
# Ensure your environment variables are set (GITHUB_TOKEN, etc.)

# Start with Gemini in local CLI mode (default)
./scripts/start.sh

# Start with Claude in Docker mode
./scripts/start.sh -p claude -d

# View all options
./scripts/start.sh --help
```

**Running directly:**

```bash
go run ./cmd/server --config config.yaml
```

Once running, the agent will listen for webhook events from GitHub.

## Usage

To use the agent, simply comment on a GitHub Issue or Pull Request.

-   **Start a task in an Issue**:
    ```
    /code Implement user login functionality.
    ```

-   **Continue a task in a Pull Request**:
    ```
    /continue Add unit tests for the login endpoint.
    ```

-   **Fix code in a Pull Request**:
    ```
    /fix The login validation is incorrect, please fix it.
    ```

## Development

### Project Structure

```
.
├── cmd/server/main.go    # Main application entrypoint
├── internal/             # Internal application logic
│   ├── agent/            # Core agent logic
│   ├── code/             # AI provider clients (Claude, Gemini)
│   ├── config/           # Configuration management
│   ├── github/           # GitHub API client
│   ├── mcp/              # Multi-Cloud-Provider abstractions
│   ├── modes/            # Handlers for different modes (agent, review)
│   ├── webhook/          # GitHub webhook handler
│   └── workspace/        # Workspace management
├── pkg/                  # Shared libraries and data models
├── scripts/              # Helper scripts
└── config.example.yaml   # Example configuration
```

### Building

To build the binary from source:

```bash
go build -o bin/codeagent ./cmd/server
```

To cross-compile for Linux:

```bash
GOOS=linux GOARCH=amd64 go build -o bin/codeagent-linux ./cmd/server
```

### Testing

To run integration tests, you can send a test webhook payload to a running server:

```bash
# Start the server with a test configuration
go run ./cmd/server --config test-config.yaml &

# Send a test event
curl -X POST http://localhost:8888/hook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: issue_comment" \
  -d @test-data/issue-comment.json
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) to learn how you can get involved.

-   🐛 [Report a Bug](https://github.com/qiniu/codeagent/issues/new?template=bug_report.md)
-   💡 [Suggest a Feature](https://github.com/qiniu/codeagent/issues/new?template=feature_request.md)
-   📝 [Improve Documentation](https://github.com/qiniu/codeagent/issues/new?template=documentation.md)

## License

This project is licensed under the [MIT License](LICENSE).