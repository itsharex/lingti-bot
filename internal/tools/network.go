package tools

import (
	"context"
	"fmt"
	gonet "net"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	psnet "github.com/shirou/gopsutil/v4/net"
)

// NetworkInterfaces lists network interfaces
func NetworkInterfaces(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	interfaces, err := psnet.Interfaces()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get interfaces: %v", err)), nil
	}

	var result strings.Builder
	result.WriteString("=== Network Interfaces ===\n\n")

	for _, iface := range interfaces {
		result.WriteString(fmt.Sprintf("--- %s ---\n", iface.Name))
		result.WriteString(fmt.Sprintf("Index: %d\n", iface.Index))
		result.WriteString(fmt.Sprintf("MTU: %d\n", iface.MTU))
		result.WriteString(fmt.Sprintf("Hardware Addr: %s\n", iface.HardwareAddr))
		result.WriteString(fmt.Sprintf("Flags: %v\n", iface.Flags))

		if len(iface.Addrs) > 0 {
			result.WriteString("Addresses:\n")
			for _, addr := range iface.Addrs {
				result.WriteString(fmt.Sprintf("  - %s\n", addr.Addr))
			}
		}
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}

// NetworkConnections lists active network connections
func NetworkConnections(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	kind := "all"
	if k, ok := req.Params.Arguments["kind"].(string); ok && k != "" {
		kind = k // "tcp", "udp", "tcp4", "tcp6", "udp4", "udp6", "all"
	}

	conns, err := psnet.Connections(kind)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get connections: %v", err)), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("=== Network Connections (%s) ===\n\n", kind))
	result.WriteString(fmt.Sprintf("%-8s %-25s %-25s %-12s\n", "TYPE", "LOCAL", "REMOTE", "STATUS"))
	result.WriteString(strings.Repeat("-", 75) + "\n")

	count := 0
	for _, conn := range conns {
		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)
		remoteAddr := fmt.Sprintf("%s:%d", conn.Raddr.IP, conn.Raddr.Port)
		if conn.Raddr.IP == "" {
			remoteAddr = "-"
		}

		typeStr := "tcp"
		if conn.Type == 2 {
			typeStr = "udp"
		}

		result.WriteString(fmt.Sprintf("%-8s %-25s %-25s %-12s\n",
			typeStr, localAddr, remoteAddr, conn.Status))

		count++
		if count >= 50 {
			result.WriteString("\n... (truncated)")
			break
		}
	}

	return mcp.NewToolResultText(result.String()), nil
}

// NetworkPing pings a host
func NetworkPing(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	host, ok := req.Params.Arguments["host"].(string)
	if !ok {
		return mcp.NewToolResultError("host is required"), nil
	}

	timeout := 5.0
	if t, ok := req.Params.Arguments["timeout"].(float64); ok && t > 0 {
		timeout = t
	}

	// Simple TCP ping (connect test)
	port := "80"
	if p, ok := req.Params.Arguments["port"].(string); ok && p != "" {
		port = p
	}

	address := fmt.Sprintf("%s:%s", host, port)

	start := time.Now()
	conn, err := gonet.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	elapsed := time.Since(start)

	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Ping %s failed: %v (%.2fms)", address, err, float64(elapsed.Microseconds())/1000)), nil
	}
	conn.Close()

	return mcp.NewToolResultText(fmt.Sprintf("Ping %s succeeded in %.2fms", address, float64(elapsed.Microseconds())/1000)), nil
}

// NetworkDNSLookup performs DNS lookup
func NetworkDNSLookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hostname, ok := req.Params.Arguments["hostname"].(string)
	if !ok {
		return mcp.NewToolResultError("hostname is required"), nil
	}

	ips, err := gonet.LookupIP(hostname)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("DNS lookup failed: %v", err)), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("DNS lookup for %s:\n\n", hostname))

	for _, ip := range ips {
		ipType := "IPv4"
		if ip.To4() == nil {
			ipType = "IPv6"
		}
		result.WriteString(fmt.Sprintf("%s: %s\n", ipType, ip.String()))
	}

	return mcp.NewToolResultText(result.String()), nil
}
