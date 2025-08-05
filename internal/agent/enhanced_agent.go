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
	
	// 使用MCP工具和AI进行代码生成
	codeResult, err := a.generateCodeWithMCP(ctx, issueCtx, mcpCtx, ws)
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
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
		Output:       codeResult.Output,
		FilesChanged: codeResult.FilesChanged,
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

// generateCodeWithMCP 使用MCP工具和AI生成代码
func (a *EnhancedAgent) generateCodeWithMCP(
	ctx context.Context,
	issueCtx *models.IssueCommentContext,
	mcpCtx *models.MCPContext,
	ws *models.Workspace,
) (*models.ExecutionResult, error) {
	xl := xlog.NewWith(ctx)
	
	// 1. 初始化code session
	codeClient, err := a.sessionManager.GetSession(ws)
	if err != nil {
		return nil, fmt.Errorf("failed to get code session: %w", err)
	}
	
	// 2. 使用MCP工具收集项目上下文 (简化实现)
	xl.Infof("Collecting project context using MCP tools")
	// TODO: 实现 ListFiles 方法
	projectFiles := []string{} // 暂时使用空列表
	xl.Infof("MCP ListFiles method not implemented, using empty file list")
	
	// 3. 分析Issue需求
	issue := issueCtx.Issue
	issueAnalysis := a.analyzeIssueRequirements(issue)
	
	// 4. 构建包含项目上下文的代码生成prompt
	codePrompt := a.buildCodeGenerationPrompt(issue, issueAnalysis, projectFiles)
	xl.Infof("Generated code generation prompt for Issue #%d", issue.GetNumber())
	
	// 5. 执行AI代码生成
	xl.Infof("Executing AI code generation")
	resp, err := codeClient.Prompt(codePrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to prompt AI for code generation: %w", err)
	}
	
	output, err := io.ReadAll(resp.Out)
	if err != nil {
		return nil, fmt.Errorf("failed to read AI output: %w", err)
	}
	
	aiOutput := string(output)
	xl.Infof("AI code generation completed, output length: %d", len(aiOutput))
	
	// 6. 使用MCP工具分析生成的文件变更 (简化实现)
	// TODO: 实现 GetChangedFiles 方法
	changedFiles := []string{} // 暂时使用空列表
	xl.Infof("MCP GetChangedFiles method not implemented, using empty list")
	
	result := &models.ExecutionResult{
		Success:      true,
		Output:       aiOutput,
		FilesChanged: changedFiles,
		Duration:     0, // 将在上层计算
	}
	
	xl.Infof("Code generation completed with %d files changed", len(changedFiles))
	return result, nil
}

// analyzeIssueRequirements 分析Issue需求
func (a *EnhancedAgent) analyzeIssueRequirements(issue *github.Issue) *models.IssueAnalysis {
	title := issue.GetTitle()
	body := issue.GetBody()
	
	analysis := &models.IssueAnalysis{
		Type:        "feature", // 默认类型
		Priority:    "medium",
		Complexity:  "medium",
		Keywords:    []string{},
		Suggestions: []string{},
	}
	
	titleLower := strings.ToLower(title)
	bodyLower := strings.ToLower(body)
	
	// 分析Issue类型
	if strings.Contains(titleLower, "bug") || strings.Contains(titleLower, "fix") || 
	   strings.Contains(bodyLower, "bug") || strings.Contains(bodyLower, "error") {
		analysis.Type = "bug"
	} else if strings.Contains(titleLower, "test") || strings.Contains(bodyLower, "test") {
		analysis.Type = "test"
	} else if strings.Contains(titleLower, "refactor") || strings.Contains(bodyLower, "refactor") {
		analysis.Type = "refactor"
	} else if strings.Contains(titleLower, "doc") || strings.Contains(bodyLower, "documentation") {
		analysis.Type = "documentation"
	}
	
	// 分析优先级
	if strings.Contains(titleLower, "urgent") || strings.Contains(titleLower, "critical") ||
	   strings.Contains(bodyLower, "urgent") || strings.Contains(bodyLower, "critical") {
		analysis.Priority = "high"
	} else if strings.Contains(titleLower, "minor") || strings.Contains(bodyLower, "minor") {
		analysis.Priority = "low"
	}
	
	// 提取关键词
	words := strings.Fields(titleLower + " " + bodyLower)
	techKeywords := []string{"api", "database", "frontend", "backend", "ui", "auth", "security", "performance"}
	for _, word := range words {
		for _, keyword := range techKeywords {
			if strings.Contains(word, keyword) {
				analysis.Keywords = append(analysis.Keywords, keyword)
			}
		}
	}
	
	return analysis
}

// buildCodeGenerationPrompt 构建代码生成的prompt
func (a *EnhancedAgent) buildCodeGenerationPrompt(issue *github.Issue, analysis *models.IssueAnalysis, projectFiles []string) string {
	var fileContext string
	if len(projectFiles) > 0 {
		fileContext = fmt.Sprintf("\n## 项目文件结构\n```\n%s\n```\n", strings.Join(projectFiles[:min(50, len(projectFiles))], "\n"))
	}
	
	prompt := fmt.Sprintf(`请根据以下Issue实现相应的代码：

## Issue信息
- **标题**: %s
- **描述**: %s
- **类型**: %s
- **优先级**: %s
- **关键词**: %s

%s

## 实现要求
1. **代码质量**: 遵循最佳实践，确保代码可读性和可维护性
2. **错误处理**: 添加适当的错误处理逻辑
3. **测试友好**: 编写易于测试的代码
4. **文档注释**: 为重要函数和复杂逻辑添加注释
5. **安全性**: 考虑安全相关的因素

## 输出格式要求
%s
简要说明实现的功能和主要变更

%s
- 列出修改或新增的文件
- 说明每个文件的主要变更

请开始实现代码，确保充分理解Issue需求并提供完整的解决方案。`,
		issue.GetTitle(),
		issue.GetBody(),
		analysis.Type,
		analysis.Priority,
		strings.Join(analysis.Keywords, ", "),
		fileContext,
		models.SectionSummary,
		models.SectionChanges)
	
	return prompt
}

// min 辅助函数：返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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