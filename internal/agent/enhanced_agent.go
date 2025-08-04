package agent

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/qiniu/codeagent/internal/code"
	"github.com/qiniu/codeagent/internal/config"
	"github.com/qiniu/codeagent/internal/events"
	ghclient "github.com/qiniu/codeagent/internal/github"
	"github.com/qiniu/codeagent/internal/interaction"
	"github.com/qiniu/codeagent/internal/mcp"
	"github.com/qiniu/codeagent/internal/mcp/servers"
	"github.com/qiniu/codeagent/internal/modes"
	"github.com/qiniu/codeagent/internal/workspace"
	"github.com/qiniu/codeagent/pkg/models"

	"github.com/google/go-github/v58/github"
	"github.com/qiniu/x/xlog"
)

// EnhancedAgent 增强版Agent，集成了新的组件架构
// 对应claude-code-action的完整智能化功能
type EnhancedAgent struct {
	// 原有组件
	config         *config.Config
	github         *ghclient.Client
	workspace      *workspace.Manager
	sessionManager *code.SessionManager
	
	// 新增组件
	eventParser    *events.Parser
	modeManager    *modes.Manager
	mcpManager     mcp.MCPManager
	mcpClient      mcp.MCPClient
	taskFactory    *interaction.TaskFactory
}

// NewEnhancedAgent 创建增强版Agent
func NewEnhancedAgent(cfg *config.Config, workspaceManager *workspace.Manager) (*EnhancedAgent, error) {
	xl := xlog.New("")
	
	// 1. 初始化GitHub客户端
	githubClient, err := ghclient.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}
	
	// 2. 初始化事件解析器
	eventParser := events.NewParser()
	
	// 3. 初始化MCP管理器和服务器
	mcpManager := mcp.NewManager()
	
	// 注册内置MCP服务器
	githubFiles := servers.NewGitHubFilesServer(githubClient)
	githubComments := servers.NewGitHubCommentsServer(githubClient)
	
	if err := mcpManager.RegisterServer("github-files", githubFiles); err != nil {
		return nil, fmt.Errorf("failed to register github-files server: %w", err)
	}
	
	if err := mcpManager.RegisterServer("github-comments", githubComments); err != nil {
		return nil, fmt.Errorf("failed to register github-comments server: %w", err)
	}
	
	// 4. 创建MCP客户端
	mcpClient := mcp.NewClient(mcpManager)
	
	// 5. 初始化SessionManager
	sessionManager := code.NewSessionManager(cfg)
	
	// 6. 初始化模式管理器
	modeManager := modes.NewManager()
	
	// 注册处理器（按优先级顺序）
	tagHandler := modes.NewTagHandler(githubClient, workspaceManager, mcpClient, sessionManager)
	agentHandler := modes.NewAgentHandler(githubClient, workspaceManager, mcpClient)
	reviewHandler := modes.NewReviewHandler(githubClient, workspaceManager, mcpClient)
	
	modeManager.RegisterHandler(tagHandler)
	modeManager.RegisterHandler(agentHandler)
	modeManager.RegisterHandler(reviewHandler)
	
	// 7. 创建任务工厂
	taskFactory := interaction.NewTaskFactory()
	
	agent := &EnhancedAgent{
		config:         cfg,
		github:         githubClient,
		workspace:      workspaceManager,
		sessionManager: sessionManager,
		eventParser:    eventParser,
		modeManager:    modeManager,
		mcpManager:     mcpManager,
		mcpClient:      mcpClient,
		taskFactory:    taskFactory,
	}
	
	xl.Infof("Enhanced Agent initialized with %d MCP servers and %d mode handlers", 
		len(mcpManager.GetServers()), modeManager.GetHandlerCount())
	
	return agent, nil
}

// ProcessGitHubEvent 处理GitHub事件的统一入口
// 替换原有的多个Process方法，使用新的事件系统
func (a *EnhancedAgent) ProcessGitHubEvent(ctx context.Context, eventType string, payload interface{}) error {
	xl := xlog.NewWith(ctx)
	
	startTime := time.Now()
	xl.Infof("Processing GitHub event: %s", eventType)
	
	// 1. 解析GitHub事件为类型安全的上下文
	githubCtx, err := a.eventParser.ParseEvent(ctx, eventType, payload)
	if err != nil {
		xl.Errorf("Failed to parse GitHub event: %v", err)
		return fmt.Errorf("failed to parse event: %w", err)
	}
	
	xl.Infof("Parsed event type: %s for repository: %s", 
		githubCtx.GetEventType(), githubCtx.GetRepository().GetFullName())
	
	// 2. 选择合适的处理器
	handler, err := a.modeManager.SelectHandler(ctx, githubCtx)
	if err != nil {
		xl.Errorf("No suitable handler found: %v", err)
		return fmt.Errorf("no handler available: %w", err)
	}
	
	xl.Infof("Selected handler with mode: %s (priority: %d)", 
		handler.GetMode(), handler.GetPriority())
	
	// 3. 执行处理
	err = handler.Execute(ctx, githubCtx)
	if err != nil {
		xl.Errorf("Handler execution failed: %v", err)
		return fmt.Errorf("handler execution failed: %w", err)
	}
	
	duration := time.Since(startTime)
	xl.Infof("GitHub event processed successfully in %v", duration)
	
	return nil
}

// ProcessIssueCommentEnhanced 增强版Issue评论处理
// 使用新的进度通信和MCP工具系统
func (a *EnhancedAgent) ProcessIssueCommentEnhanced(ctx context.Context, event *github.IssueCommentEvent) error {
	xl := xlog.NewWith(ctx)
	
	// 1. 解析为类型安全的上下文
	githubCtx, err := a.eventParser.ParseIssueCommentEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to parse issue comment event: %w", err)
	}
	
	issueCommentCtx, ok := githubCtx.(*models.IssueCommentContext)
	if !ok {
		return fmt.Errorf("invalid context type for issue comment")
	}
	
	// 2. 创建进度评论管理器
	pcm := interaction.NewProgressCommentManager(a.github, 
		issueCommentCtx.GetRepository(), issueCommentCtx.Issue.GetNumber())
	
	// 3. 创建任务列表
	tasks := a.taskFactory.CreateIssueProcessingTasks()
	
	// 4. 初始化进度跟踪
	if err := pcm.InitializeProgress(ctx, tasks); err != nil {
		xl.Errorf("Failed to initialize progress tracking: %v", err)
		return err
	}
	
	// 5. 创建MCP上下文
	mcpCtx := &models.MCPContext{
		Repository:  githubCtx,
		Issue:       issueCommentCtx.Issue,
		User:        issueCommentCtx.GetSender(),
		Permissions: []string{"github:read", "github:write"},
		Constraints: []string{}, // 根据需要添加约束
	}
	
	// 6. 执行处理流程
	result, err := a.executeIssueProcessingWithProgress(ctx, issueCommentCtx, mcpCtx, pcm)
	if err != nil {
		xl.Errorf("Issue processing failed: %v", err)
		
		// 最终化失败结果
		failureResult := &models.ProgressExecutionResult{
			Success: false,
			Error:   err.Error(),
			Duration: time.Since(pcm.GetTracker().StartTime),
		}
		
		if finalizeErr := pcm.FinalizeComment(ctx, failureResult); finalizeErr != nil {
			xl.Errorf("Failed to finalize failure comment: %v", finalizeErr)
		}
		
		return err
	}
	
	// 7. 最终化成功结果
	if err := pcm.FinalizeComment(ctx, result); err != nil {
		xl.Errorf("Failed to finalize success comment: %v", err)
		return err
	}
	
	xl.Infof("Issue comment processed successfully")
	return nil
}

// executeIssueProcessingWithProgress 执行Issue处理流程，带进度跟踪
func (a *EnhancedAgent) executeIssueProcessingWithProgress(
	ctx context.Context, 
	issueCtx *models.IssueCommentContext, 
	mcpCtx *models.MCPContext, 
	pcm *interaction.ProgressCommentManager,
) (*models.ProgressExecutionResult, error) {
	xl := xlog.NewWith(ctx)
	
	// 1. 收集上下文信息
	if err := pcm.UpdateTask(ctx, "gather-context", models.TaskStatusInProgress, "Analyzing issue and requirements"); err != nil {
		return nil, err
	}
	
	// 使用MCP工具收集更多上下文
	tools, err := a.mcpClient.PrepareTools(ctx, mcpCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare MCP tools: %w", err)
	}
	
	xl.Infof("Prepared %d MCP tools for issue processing", len(tools))
	
	if err := pcm.UpdateTask(ctx, "gather-context", models.TaskStatusCompleted); err != nil {
		return nil, err
	}
	
	// 2. 设置工作空间
	if err := pcm.UpdateTask(ctx, "setup-workspace", models.TaskStatusInProgress, "Creating workspace and branch"); err != nil {
		return nil, err
	}
	
	ws := a.workspace.CreateWorkspaceFromIssue(issueCtx.Issue)
	if ws == nil {
		return nil, fmt.Errorf("failed to create workspace")
	}
	
	// 更新MCP上下文
	mcpCtx.WorkspacePath = ws.Path
	mcpCtx.BranchName = ws.Branch
	
	if err := a.github.CreateBranch(ws); err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}
	
	if err := pcm.UpdateTask(ctx, "setup-workspace", models.TaskStatusCompleted); err != nil {
		return nil, err
	}
	
	// 3. 生成代码（使用MCP工具）
	if err := pcm.UpdateTask(ctx, "generate-code", models.TaskStatusInProgress, "Generating code implementation"); err != nil {
		return nil, err
	}
	
	// 创建AI会话进行代码生成
	session, err := a.sessionManager.GetSession(ws)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI session: %w", err)
	}
	
	// 构建代码生成提示
	codePrompt := fmt.Sprintf(`You are an AI coding assistant working on GitHub issue #%d.

Issue Title: %s
Issue Description: %s

Please implement the requested functionality by creating or modifying the necessary files.

Instructions:
1. Analyze the issue requirements carefully
2. Create well-structured, maintainable code
3. Follow best practices for the project's programming language
4. Include appropriate error handling
5. Add comments for complex logic
6. Ensure code is production-ready

Provide your implementation with clear explanations of the changes made.`,
		issueCtx.Issue.GetNumber(),
		issueCtx.Issue.GetTitle(),
		issueCtx.Issue.GetBody())
	
	xl.Infof("Executing AI code generation for issue #%d", issueCtx.Issue.GetNumber())
	
	// 执行AI代码生成
	resp, err := session.Prompt(codePrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate code with AI: %w", err)
	}
	
	// 读取AI生成的响应
	aiOutput, err := io.ReadAll(resp.Out)
	if err != nil {
		return nil, fmt.Errorf("failed to read AI response: %w", err)
	}
	
	aiOutputStr := string(aiOutput)
	xl.Infof("AI code generation completed, output length: %d", len(aiOutputStr))
	xl.Debugf("AI Output: %s", aiOutputStr)
	
	// 使用MCP工具获取和分析项目文件结构
	contextTools, err := a.mcpClient.PrepareTools(ctx, mcpCtx)
	if err != nil {
		xl.Warnf("Failed to prepare MCP tools for file operations: %v", err)
	} else {
		xl.Infof("Prepared %d MCP tools for file operations", len(contextTools))
	}
	
	var filesChanged []string
	
	// 如果有文件操作工具，使用它们来分析项目结构并协助AI代码生成
	if len(contextTools) > 0 {
		xl.Infof("Analyzing project structure with MCP tools")
		
		// 1. 列出项目文件以了解结构
		listFilesCall := &models.ToolCall{
			ID: "list_files_analysis",
			Function: models.ToolFunction{
				Name: "github-files_list_files",
				Arguments: map[string]interface{}{
					"path":      ".",
					"recursive": true,
				},
			},
		}
		
		results, err := a.mcpClient.ExecuteToolCalls(ctx, []*models.ToolCall{listFilesCall}, mcpCtx)
		if err != nil {
			xl.Warnf("Failed to list files via MCP: %v", err)
		} else if len(results) > 0 && results[0].Success {
			xl.Infof("Successfully analyzed project structure via MCP")
			xl.Debugf("Project files: %s", results[0].Content)
			
			// 2. 基于项目结构，向AI提供更详细的上下文
			structurePrompt := fmt.Sprintf(`Based on the project structure analysis:

Project Files: %s

Previous AI Response: %s

Now please provide specific file operations to implement the requested functionality:
1. Which files need to be created or modified?
2. What should be the content of each file?
3. Provide the exact file paths and content.

Format your response as a series of file operations, each clearly marked with the file path and content.`, 
				results[0].Content, aiOutputStr)
			
			// 执行结构化文件操作查询
			structureResp, err := session.Prompt(structurePrompt)
			if err != nil {
				xl.Warnf("Failed to get structured file operations: %v", err)
			} else {
				structureOutput, err := io.ReadAll(structureResp.Out)
				if err != nil {
					xl.Warnf("Failed to read structure response: %v", err)
				} else {
					xl.Infof("Got structured file operations from AI")
					xl.Debugf("Structure operations: %s", string(structureOutput))
					
					// 这里可以解析AI的响应并执行具体的文件操作
					// 为了演示，我们记录文件变更，实际项目中应该解析AI响应并执行文件操作
					filesChanged = a.extractFilesFromAIResponse(string(structureOutput))
				}
			}
		}
	}
	
	// 更新执行结果中的文件变更列表
	if len(filesChanged) > 0 {
		xl.Infof("AI identified %d files for modification: %v", len(filesChanged), filesChanged)
	}
	
	if err := pcm.UpdateTask(ctx, "generate-code", models.TaskStatusCompleted); err != nil {
		return nil, err
	}
	
	// 4. 提交变更
	if err := pcm.UpdateTask(ctx, "commit-changes", models.TaskStatusInProgress, "Committing changes"); err != nil {
		return nil, err
	}
	
	// 使用现有的提交逻辑
	execResult := &models.ExecutionResult{
		Success:      true,
		Output:       aiOutputStr,
		FilesChanged: filesChanged,
		Duration:     time.Since(pcm.GetTracker().StartTime),
	}
	
	if err := a.github.CommitAndPush(ws, execResult, nil); err != nil {
		return nil, fmt.Errorf("failed to commit and push: %w", err)
	}
	
	if err := pcm.UpdateTask(ctx, "commit-changes", models.TaskStatusCompleted); err != nil {
		return nil, err
	}
	
	// 5. 创建PR
	if err := pcm.UpdateTask(ctx, "create-pr", models.TaskStatusInProgress, "Creating pull request"); err != nil {
		return nil, err
	}
	
	pr, err := a.github.CreatePullRequest(ws)
	if err != nil {
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}
	
	if err := pcm.UpdateTask(ctx, "create-pr", models.TaskStatusCompleted); err != nil {
		return nil, err
	}
	
	// 构建结果
	result := &models.ProgressExecutionResult{
		Success:        true,
		Output:         execResult.Output,
		FilesChanged:   execResult.FilesChanged,
		Duration:       time.Since(pcm.GetTracker().StartTime),
		Summary:        fmt.Sprintf("Successfully implemented Issue #%d", issueCtx.Issue.GetNumber()),
		BranchName:     ws.Branch,
		PullRequestURL: pr.GetHTMLURL(),
		TaskResults:    pcm.GetTracker().Tasks,
	}
	
	return result, nil
}

// GetMCPManager 获取MCP管理器（用于外部扩展）
func (a *EnhancedAgent) GetMCPManager() mcp.MCPManager {
	return a.mcpManager
}

// GetModeManager 获取模式管理器（用于外部扩展）
func (a *EnhancedAgent) GetModeManager() *modes.Manager {
	return a.modeManager
}

// Shutdown 关闭增强版Agent
func (a *EnhancedAgent) Shutdown(ctx context.Context) error {
	xl := xlog.NewWith(ctx)
	
	// 关闭MCP管理器
	if err := a.mcpManager.Shutdown(ctx); err != nil {
		xl.Errorf("Failed to shutdown MCP manager: %v", err)
		return err
	}
	
	xl.Infof("Enhanced Agent shutdown completed")
	return nil
}

// extractFilesFromAIResponse 从AI响应中提取文件路径
// 这是一个简单的实现，实际项目中应该有更复杂的解析逻辑
func (a *EnhancedAgent) extractFilesFromAIResponse(response string) []string {
	var files []string
	
	// 简单的文件路径提取逻辑
	// 寻找常见的文件路径模式
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// 寻找文件路径模式（简化版本）
		if strings.Contains(line, ".go") || 
		   strings.Contains(line, ".py") || 
		   strings.Contains(line, ".js") ||
		   strings.Contains(line, ".ts") ||
		   strings.Contains(line, ".java") ||
		   strings.Contains(line, ".cpp") ||
		   strings.Contains(line, ".c") ||
		   strings.Contains(line, ".md") ||
		   strings.Contains(line, ".json") ||
		   strings.Contains(line, ".yaml") ||
		   strings.Contains(line, ".xml") {
			
			// 提取可能的文件路径
			words := strings.Fields(line)
			for _, word := range words {
				if strings.Contains(word, ".") && !strings.HasPrefix(word, "http") {
					// 简单验证是否看起来像文件路径
					if strings.Count(word, ".") == 1 || strings.Contains(word, "/") {
						files = append(files, word)
					}
				}
			}
		}
	}
	
	// 去重
	seen := make(map[string]bool)
	var uniqueFiles []string
	for _, file := range files {
		if !seen[file] {
			seen[file] = true
			uniqueFiles = append(uniqueFiles, file)
		}
	}
	
	return uniqueFiles
}