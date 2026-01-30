package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pltanton/lingti-bot/internal/tools"
)

const (
	ServerName    = "lingti-bot"
	ServerVersion = "1.0.0"
)

// NewServer creates a new MCP server with all tools registered
func NewServer() *server.MCPServer {
	s := server.NewMCPServer(ServerName, ServerVersion,
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)

	registerFilesystemTools(s)
	registerShellTools(s)
	registerSystemTools(s)
	registerProcessTools(s)
	registerNetworkTools(s)

	return s
}

func registerFilesystemTools(s *server.MCPServer) {
	// file_read
	s.AddTool(mcp.NewTool("file_read",
		mcp.WithDescription("Read the contents of a file"),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the file to read")),
	), tools.FileRead)

	// file_write
	s.AddTool(mcp.NewTool("file_write",
		mcp.WithDescription("Write content to a file"),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the file to write")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Content to write to the file")),
	), tools.FileWrite)

	// file_list
	s.AddTool(mcp.NewTool("file_list",
		mcp.WithDescription("List contents of a directory"),
		mcp.WithString("path", mcp.Description("Path to the directory (default: current directory)")),
	), tools.FileList)

	// file_search
	s.AddTool(mcp.NewTool("file_search",
		mcp.WithDescription("Search for files matching a pattern"),
		mcp.WithString("pattern", mcp.Required(), mcp.Description("Glob pattern to match (e.g., *.go, *.txt)")),
		mcp.WithString("path", mcp.Description("Directory to search in (default: current directory)")),
	), tools.FileSearch)

	// file_info
	s.AddTool(mcp.NewTool("file_info",
		mcp.WithDescription("Get detailed information about a file"),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the file")),
	), tools.FileInfo)
}

func registerShellTools(s *server.MCPServer) {
	// shell_execute
	s.AddTool(mcp.NewTool("shell_execute",
		mcp.WithDescription("Execute a shell command"),
		mcp.WithString("command", mcp.Required(), mcp.Description("The command to execute")),
		mcp.WithNumber("timeout", mcp.Description("Timeout in seconds (default: 30)")),
		mcp.WithString("working_directory", mcp.Description("Working directory for the command")),
	), tools.ShellExecute)

	// shell_which
	s.AddTool(mcp.NewTool("shell_which",
		mcp.WithDescription("Find the path of an executable"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the executable to find")),
	), tools.ShellWhich)
}

func registerSystemTools(s *server.MCPServer) {
	// system_info
	s.AddTool(mcp.NewTool("system_info",
		mcp.WithDescription("Get system information (CPU, memory, OS)"),
	), tools.SystemInfo)

	// disk_usage
	s.AddTool(mcp.NewTool("disk_usage",
		mcp.WithDescription("Get disk usage information"),
		mcp.WithString("path", mcp.Description("Path to check (default: /)")),
	), tools.DiskUsage)

	// env_get
	s.AddTool(mcp.NewTool("env_get",
		mcp.WithDescription("Get an environment variable"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the environment variable")),
	), tools.EnvGet)

	// env_list
	s.AddTool(mcp.NewTool("env_list",
		mcp.WithDescription("List all environment variables"),
	), tools.EnvList)
}

func registerProcessTools(s *server.MCPServer) {
	// process_list
	s.AddTool(mcp.NewTool("process_list",
		mcp.WithDescription("List running processes"),
		mcp.WithString("filter", mcp.Description("Filter processes by name (optional)")),
	), tools.ProcessList)

	// process_info
	s.AddTool(mcp.NewTool("process_info",
		mcp.WithDescription("Get detailed information about a process"),
		mcp.WithNumber("pid", mcp.Required(), mcp.Description("Process ID")),
	), tools.ProcessInfo)

	// process_kill
	s.AddTool(mcp.NewTool("process_kill",
		mcp.WithDescription("Kill a process by PID"),
		mcp.WithNumber("pid", mcp.Required(), mcp.Description("Process ID to kill")),
	), tools.ProcessKill)
}

func registerNetworkTools(s *server.MCPServer) {
	// network_interfaces
	s.AddTool(mcp.NewTool("network_interfaces",
		mcp.WithDescription("List network interfaces"),
	), tools.NetworkInterfaces)

	// network_connections
	s.AddTool(mcp.NewTool("network_connections",
		mcp.WithDescription("List active network connections"),
		mcp.WithString("kind", mcp.Description("Connection type: tcp, udp, tcp4, tcp6, udp4, udp6, all (default: all)")),
	), tools.NetworkConnections)

	// network_ping
	s.AddTool(mcp.NewTool("network_ping",
		mcp.WithDescription("Ping a host (TCP connect test)"),
		mcp.WithString("host", mcp.Required(), mcp.Description("Host to ping")),
		mcp.WithString("port", mcp.Description("Port to connect to (default: 80)")),
		mcp.WithNumber("timeout", mcp.Description("Timeout in seconds (default: 5)")),
	), tools.NetworkPing)

	// network_dns_lookup
	s.AddTool(mcp.NewTool("network_dns_lookup",
		mcp.WithDescription("Perform DNS lookup for a hostname"),
		mcp.WithString("hostname", mcp.Required(), mcp.Description("Hostname to look up")),
	), tools.NetworkDNSLookup)
}
