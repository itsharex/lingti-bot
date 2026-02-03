package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// GitHubPRList lists pull requests
func GitHubPRList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	state := "open"
	if s, ok := req.Params.Arguments["state"].(string); ok && s != "" {
		state = s
	}

	limit := 10
	if l, ok := req.Params.Arguments["limit"].(float64); ok {
		limit = int(l)
	}

	cmd := exec.CommandContext(ctx, "gh", "pr", "list", "--state", state, "--limit", fmt.Sprintf("%d", limit))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("gh pr list failed: %v\n%s", err, output)), nil
	}

	if len(output) == 0 {
		return mcp.NewToolResultText("No pull requests found"), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// GitHubPRView views a pull request
func GitHubPRView(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	number, ok := req.Params.Arguments["number"].(float64)
	if !ok {
		return mcp.NewToolResultError("PR number is required"), nil
	}

	cmd := exec.CommandContext(ctx, "gh", "pr", "view", fmt.Sprintf("%.0f", number))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("gh pr view failed: %v\n%s", err, output)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// GitHubIssueList lists issues
func GitHubIssueList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	state := "open"
	if s, ok := req.Params.Arguments["state"].(string); ok && s != "" {
		state = s
	}

	limit := 10
	if l, ok := req.Params.Arguments["limit"].(float64); ok {
		limit = int(l)
	}

	cmd := exec.CommandContext(ctx, "gh", "issue", "list", "--state", state, "--limit", fmt.Sprintf("%d", limit))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("gh issue list failed: %v\n%s", err, output)), nil
	}

	if len(output) == 0 {
		return mcp.NewToolResultText("No issues found"), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// GitHubIssueView views an issue
func GitHubIssueView(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	number, ok := req.Params.Arguments["number"].(float64)
	if !ok {
		return mcp.NewToolResultError("Issue number is required"), nil
	}

	cmd := exec.CommandContext(ctx, "gh", "issue", "view", fmt.Sprintf("%.0f", number))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("gh issue view failed: %v\n%s", err, output)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// GitHubIssueCreate creates an issue
func GitHubIssueCreate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, ok := req.Params.Arguments["title"].(string)
	if !ok || title == "" {
		return mcp.NewToolResultError("title is required"), nil
	}

	args := []string{"issue", "create", "--title", title}

	if body, ok := req.Params.Arguments["body"].(string); ok && body != "" {
		args = append(args, "--body", body)
	}

	if labels, ok := req.Params.Arguments["labels"].(string); ok && labels != "" {
		args = append(args, "--label", labels)
	}

	cmd := exec.CommandContext(ctx, "gh", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("gh issue create failed: %v\n%s", err, output)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// GitHubRepoView views repository info
func GitHubRepoView(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cmd := exec.CommandContext(ctx, "gh", "repo", "view")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("gh repo view failed: %v\n%s", err, output)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// GitStatus runs git status
func GitStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cmd := exec.CommandContext(ctx, "git", "status", "--short")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("git status failed: %v\n%s", err, output)), nil
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		return mcp.NewToolResultText("Working tree clean"), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// GitLog runs git log
func GitLog(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	limit := 10
	if l, ok := req.Params.Arguments["limit"].(float64); ok {
		limit = int(l)
	}

	cmd := exec.CommandContext(ctx, "git", "log", "--oneline", fmt.Sprintf("-n%d", limit))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("git log failed: %v\n%s", err, output)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// GitDiff runs git diff
func GitDiff(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"diff"}

	if staged, ok := req.Params.Arguments["staged"].(bool); ok && staged {
		args = append(args, "--staged")
	}

	if file, ok := req.Params.Arguments["file"].(string); ok && file != "" {
		args = append(args, file)
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("git diff failed: %v\n%s", err, output)), nil
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		return mcp.NewToolResultText("No changes"), nil
	}

	// Truncate if too long
	result := string(output)
	if len(result) > 5000 {
		result = result[:5000] + "\n... (truncated)"
	}

	return mcp.NewToolResultText(result), nil
}

// GitBranch lists or shows current branch
func GitBranch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cmd := exec.CommandContext(ctx, "git", "branch", "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("git branch failed: %v\n%s", err, output)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}
