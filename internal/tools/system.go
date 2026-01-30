package tools

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

// SystemInfo returns general system information
func SystemInfo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get host info: %v", err)), nil
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get memory info: %v", err)), nil
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get CPU info: %v", err)), nil
	}

	cpuPercent, _ := cpu.Percent(0, false)

	var result strings.Builder
	result.WriteString("=== System Information ===\n\n")

	result.WriteString("--- Host ---\n")
	result.WriteString(fmt.Sprintf("Hostname: %s\n", hostInfo.Hostname))
	result.WriteString(fmt.Sprintf("OS: %s\n", hostInfo.OS))
	result.WriteString(fmt.Sprintf("Platform: %s %s\n", hostInfo.Platform, hostInfo.PlatformVersion))
	result.WriteString(fmt.Sprintf("Kernel: %s\n", hostInfo.KernelVersion))
	result.WriteString(fmt.Sprintf("Architecture: %s\n", hostInfo.KernelArch))
	result.WriteString(fmt.Sprintf("Uptime: %d seconds\n", hostInfo.Uptime))

	result.WriteString("\n--- CPU ---\n")
	if len(cpuInfo) > 0 {
		result.WriteString(fmt.Sprintf("Model: %s\n", cpuInfo[0].ModelName))
		result.WriteString(fmt.Sprintf("Cores: %d physical, %d logical\n", cpuInfo[0].Cores, runtime.NumCPU()))
	}
	if len(cpuPercent) > 0 {
		result.WriteString(fmt.Sprintf("Usage: %.1f%%\n", cpuPercent[0]))
	}

	result.WriteString("\n--- Memory ---\n")
	result.WriteString(fmt.Sprintf("Total: %s\n", FormatBytes(memInfo.Total)))
	result.WriteString(fmt.Sprintf("Available: %s\n", FormatBytes(memInfo.Available)))
	result.WriteString(fmt.Sprintf("Used: %s (%.1f%%)\n", FormatBytes(memInfo.Used), memInfo.UsedPercent))

	return mcp.NewToolResultText(result.String()), nil
}

// DiskUsage returns disk usage information
func DiskUsage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := "/"
	if p, ok := req.Params.Arguments["path"].(string); ok && p != "" {
		path = p
	}

	usage, err := disk.Usage(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get disk usage: %v", err)), nil
	}

	partitions, _ := disk.Partitions(false)

	var result strings.Builder
	result.WriteString("=== Disk Usage ===\n\n")

	result.WriteString(fmt.Sprintf("--- %s ---\n", path))
	result.WriteString(fmt.Sprintf("Total: %s\n", FormatBytes(usage.Total)))
	result.WriteString(fmt.Sprintf("Free: %s\n", FormatBytes(usage.Free)))
	result.WriteString(fmt.Sprintf("Used: %s (%.1f%%)\n", FormatBytes(usage.Used), usage.UsedPercent))

	if len(partitions) > 0 {
		result.WriteString("\n--- Partitions ---\n")
		for _, p := range partitions {
			result.WriteString(fmt.Sprintf("%s -> %s (%s)\n", p.Device, p.Mountpoint, p.Fstype))
		}
	}

	return mcp.NewToolResultText(result.String()), nil
}

// EnvGet gets an environment variable
func EnvGet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := req.Params.Arguments["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name is required"), nil
	}

	value := os.Getenv(name)
	if value == "" {
		return mcp.NewToolResultText(fmt.Sprintf("%s is not set", name)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("%s=%s", name, value)), nil
}

// EnvList lists all environment variables
func EnvList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	env := os.Environ()

	var result strings.Builder
	for _, e := range env {
		result.WriteString(e)
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}
