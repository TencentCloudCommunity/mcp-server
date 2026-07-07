package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"postgres_server/security"
)

func main() {
	transportMode := flag.String("transport", normalizeTransport(envOrDefault("SMOKE_TRANSPORT", envOrDefault("MCP_TRANSPORT", "streamable-http"))), "transport: streamable-http | sse | stdio")
	url := flag.String("url", defaultURL(*transportMode, "SMOKE_SERVER_URL", "SMOKE_SSE_URL"), "MCP URL for HTTP/SSE transports")
	stdioCommand := flag.String("command", envOrDefault("SMOKE_STDIO_COMMAND", ""), "stdio command path")
	region := flag.String("region", envOrDefault("SMOKE_REGION", "ap-guangzhou"), "region for readonly tool calls")
	instanceID := flag.String("instance-id", envOrDefault("SMOKE_INSTANCE_ID", ""), "instance id for instance-scoped readonly tool calls")
	listLimit := flag.Int("list-limit", envOrDefaultInt("SMOKE_LIST_LIMIT", 12), "max tool names to print from tools/list")
	flag.Parse()

	mode := normalizeTransport(*transportMode)
	*transportMode = mode

	fmt.Println("== MCP smoke test ==")
	fmt.Printf("Transport: %s\n", mode)
	switch mode {
	case "stdio":
		fmt.Printf("Command: %s\n", *stdioCommand)
	default:
		fmt.Printf("Server URL: %s\n", *url)
	}
	fmt.Printf("Region: %s\n", *region)
	if *instanceID == "" {
		fmt.Println("InstanceID: <not set>")
		fmt.Println("Note: instance-scoped readonly calls will be skipped.")
	} else {
		fmt.Printf("InstanceID: %s\n", *instanceID)
	}
	fmt.Println()

	c, err := newMCPClient(mode, *url, *stdioCommand)
	must("create MCP client", err)

	ctx := context.Background()
	must("start client", c.Start(ctx))
	defer c.Close()

	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "pg-mcp-smoke-client",
		Version: "1.0.0",
	}
	initReq.Params.Capabilities = mcp.ClientCapabilities{}

	initCtx, cancelInit := context.WithTimeout(ctx, 15*time.Second)
	defer cancelInit()
	initRes, err := c.Initialize(initCtx, initReq)
	must("initialize", err)

	capsJSON, _ := json.Marshal(initRes.Capabilities)
	fmt.Println("== initialize ==")
	fmt.Printf("server: %s %s\n", initRes.ServerInfo.Name, initRes.ServerInfo.Version)
	fmt.Printf("capabilities: %s\n", string(capsJSON))
	fmt.Println()

	pingCtx, cancelPing := context.WithTimeout(ctx, 10*time.Second)
	defer cancelPing()
	must("ping", c.Ping(pingCtx))
	fmt.Println("== ping ==")
	fmt.Println("status: ok")
	fmt.Println()

	listCtx, cancelList := context.WithTimeout(ctx, 15*time.Second)
	defer cancelList()
	toolList, err := c.ListTools(listCtx, mcp.ListToolsRequest{})
	must("tools/list", err)

	fmt.Println("== tools/list ==")
	fmt.Printf("tool_count: %d\n", len(toolList.Tools))
	toolNames := make([]string, 0, len(toolList.Tools))
	toolByName := make(map[string]mcp.Tool, len(toolList.Tools))
	for _, tool := range toolList.Tools {
		toolNames = append(toolNames, tool.Name)
		toolByName[tool.Name] = tool
	}
	sort.Strings(toolNames)
	limit := *listLimit
	if limit <= 0 || limit > len(toolNames) {
		limit = len(toolNames)
	}
	fmt.Printf("sample_tools(first_%d_sorted):\n", limit)
	for i := 0; i < limit; i++ {
		fmt.Printf("  - %s\n", toolNames[i])
	}
	fmt.Println()

	required := []string{
		"postgres-DescribeRegions",
		"postgres-DescribeDBVersions",
		"postgres-DescribeDBInstances",
		"postgres-DescribeDBInstanceAttribute",
		"postgres-CreateInstances",
		"postgres-CreateReadOnlyDBInstance",
	}
	fmt.Println("== schema spot check ==")
	for _, name := range required {
		tool, ok := toolByName[name]
		if !ok {
			fmt.Printf("%s -> found=false\n", name)
			continue
		}
		_, hasConfirm := tool.InputSchema.Properties["confirm"]
		fmt.Printf("%s -> found=true required=%v has_confirm=%v\n", name, tool.InputSchema.Required, hasConfirm)
	}
	fmt.Println()

	fmt.Println("== tools/call ==")
	callAndPrint(ctx, c, "postgres-DescribeRegions", map[string]any{"region": *region})
	callAndPrint(ctx, c, "postgres-DescribeDBVersions", map[string]any{"region": *region})
	callAndPrint(ctx, c, "postgres-DescribeDBInstances", map[string]any{"region": *region, "Limit": 2, "Offset": 0})
	if *instanceID != "" {
		callAndPrint(ctx, c, "postgres-DescribeDBInstanceAttribute", map[string]any{"region": *region, "DBInstanceId": *instanceID})
	} else {
		fmt.Println("-- postgres-DescribeDBInstanceAttribute --")
		fmt.Println("skipped: missing instance id")
		fmt.Println()
	}
	callAndPrint(ctx, c, "postgres-CreateInstances", map[string]any{"region": *region, "confirm": false})
}

func newMCPClient(mode, url, stdioCommand string) (*client.Client, error) {
	switch mode {
	case "sse":
		return client.NewSSEMCPClient(url, security.MCPClientOptionsFromEnv()...)
	case "stdio":
		if strings.TrimSpace(stdioCommand) == "" {
			return nil, fmt.Errorf("missing stdio command: set SMOKE_STDIO_COMMAND or pass --command")
		}
		return client.NewStdioMCPClient(strings.TrimSpace(stdioCommand), ensureTransportEnv(os.Environ(), "stdio"))
	default:
		return client.NewStreamableHttpClient(url, security.MCPStreamableHTTPClientOptionsFromEnv()...)
	}
}

func defaultURL(mode, primaryKey, legacyKey string) string {
	if v := envOrDefault(primaryKey, ""); v != "" {
		return v
	}
	if v := envOrDefault(legacyKey, ""); v != "" {
		return v
	}
	if mode == "sse" {
		return "http://127.0.0.1:9000/sse"
	}
	return "http://127.0.0.1:9000/mcp"
}

func normalizeTransport(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "sse":
		return "sse"
	case "stdio":
		return "stdio"
	default:
		return "streamable-http"
	}
}

func ensureTransportEnv(env []string, mode string) []string {
	result := append([]string(nil), env...)
	prefix := "MCP_TRANSPORT="
	for i, item := range result {
		if strings.HasPrefix(item, prefix) {
			result[i] = prefix + mode
			return result
		}
	}
	return append(result, prefix+mode)
}

func callAndPrint(ctx context.Context, c *client.Client, toolName string, args map[string]any) {
	fmt.Printf("-- %s --\n", toolName)
	argsJSON, _ := json.Marshal(args)
	fmt.Printf("request: %s\n", string(argsJSON))

	callCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	req := mcp.CallToolRequest{}
	req.Params.Name = toolName
	req.Params.Arguments = args
	res, err := c.CallTool(callCtx, req)
	if err != nil {
		fmt.Printf("transport_error: %v\n\n", err)
		return
	}

	fmt.Printf("is_error: %v\n", res.IsError)
	fmt.Printf("response: %s\n\n", truncate(renderToolResult(res), 700))
}

func renderToolResult(res *mcp.CallToolResult) string {
	parts := make([]string, 0, len(res.Content)+1)
	for _, content := range res.Content {
		switch v := content.(type) {
		case mcp.TextContent:
			parts = append(parts, v.Text)
		default:
			raw, err := json.Marshal(v)
			if err == nil {
				parts = append(parts, string(raw))
			}
		}
	}
	if len(parts) == 0 && res.StructuredContent != nil {
		raw, err := json.Marshal(res.StructuredContent)
		if err == nil {
			parts = append(parts, string(raw))
		}
	}
	if len(parts) == 0 {
		return "<empty>"
	}
	return strings.Join(parts, "\n")
}

func truncate(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max] + "...(truncated)"
}

func must(step string, err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s failed: %v\n", step, err)
	os.Exit(1)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envOrDefaultInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		var parsed int
		if _, err := fmt.Sscanf(v, "%d", &parsed); err == nil {
			return parsed
		}
	}
	return fallback
}
