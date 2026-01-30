package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shirou/gopsutil/v4/process"
)

// ProcessList lists running processes
func ProcessList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	procs, err := process.Processes()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get processes: %v", err)), nil
	}

	// Optional filter by name
	filter := ""
	if f, ok := req.Params.Arguments["filter"].(string); ok {
		filter = strings.ToLower(f)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("%-8s %-8s %-6s %-6s %s\n", "PID", "PPID", "CPU%", "MEM%", "NAME"))
	result.WriteString(strings.Repeat("-", 60) + "\n")

	count := 0
	for _, p := range procs {
		name, _ := p.Name()
		if filter != "" && !strings.Contains(strings.ToLower(name), filter) {
			continue
		}

		ppid, _ := p.Ppid()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()

		result.WriteString(fmt.Sprintf("%-8d %-8d %-6.1f %-6.1f %s\n",
			p.Pid, ppid, cpuPercent, memPercent, name))

		count++
		if count >= 50 {
			result.WriteString("\n... (truncated, use filter to narrow results)")
			break
		}
	}

	return mcp.NewToolResultText(result.String()), nil
}

// ProcessInfo gets detailed information about a specific process
func ProcessInfo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pidFloat, ok := req.Params.Arguments["pid"].(float64)
	if !ok {
		return mcp.NewToolResultError("pid is required"), nil
	}
	pid := int32(pidFloat)

	p, err := process.NewProcess(pid)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("process not found: %v", err)), nil
	}

	name, _ := p.Name()
	cmdline, _ := p.Cmdline()
	cwd, _ := p.Cwd()
	exe, _ := p.Exe()
	ppid, _ := p.Ppid()
	status, _ := p.Status()
	createTime, _ := p.CreateTime()
	cpuPercent, _ := p.CPUPercent()
	memPercent, _ := p.MemoryPercent()
	memInfo, _ := p.MemoryInfo()
	username, _ := p.Username()

	var result strings.Builder
	result.WriteString(fmt.Sprintf("=== Process %d ===\n\n", pid))
	result.WriteString(fmt.Sprintf("Name: %s\n", name))
	result.WriteString(fmt.Sprintf("Status: %v\n", status))
	result.WriteString(fmt.Sprintf("User: %s\n", username))
	result.WriteString(fmt.Sprintf("Parent PID: %d\n", ppid))
	result.WriteString(fmt.Sprintf("Executable: %s\n", exe))
	result.WriteString(fmt.Sprintf("Command: %s\n", cmdline))
	result.WriteString(fmt.Sprintf("Working Dir: %s\n", cwd))
	result.WriteString(fmt.Sprintf("Created: %d\n", createTime))
	result.WriteString(fmt.Sprintf("CPU: %.1f%%\n", cpuPercent))
	result.WriteString(fmt.Sprintf("Memory: %.1f%%\n", memPercent))
	if memInfo != nil {
		result.WriteString(fmt.Sprintf("RSS: %s\n", FormatBytes(memInfo.RSS)))
		result.WriteString(fmt.Sprintf("VMS: %s\n", FormatBytes(memInfo.VMS)))
	}

	return mcp.NewToolResultText(result.String()), nil
}

// ProcessKill terminates a process
func ProcessKill(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pidFloat, ok := req.Params.Arguments["pid"].(float64)
	if !ok {
		return mcp.NewToolResultError("pid is required"), nil
	}
	pid := int32(pidFloat)

	// Safety check - don't allow killing PID 1 or own process
	if pid == 1 {
		return mcp.NewToolResultError("cannot kill init process (PID 1)"), nil
	}

	p, err := process.NewProcess(pid)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("process not found: %v", err)), nil
	}

	name, _ := p.Name()

	if err := p.Kill(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to kill process: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully killed process %d (%s)", pid, name)), nil
}
