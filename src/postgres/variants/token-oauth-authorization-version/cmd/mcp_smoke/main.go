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
	url := flag.String("url", envOrDefault("SMOKE_SSE_URL", "http://127.0.0.1:9000/sse"), "MCP SSE URL")
	region := flag.String("region", envOrDefault("SMOKE_REGION", "ap-guangzhou"), "region for readonly tool calls")
	instanceID := flag.String("instance-id", envOrDefault("SMOKE_INSTANCE_ID", ""), "instance id for instance-scoped readonly tool calls")
	listLimit := flag.Int("list-limit", envOrDefaultInt("SMOKE_LIST_LIMIT", 12), "max tool names to print from tools/list")
	flag.Parse()

	fmt.Println("== MCP smoke test ==")
	fmt.Printf("SSE URL: %s\n", *url)
	fmt.Printf("Region: %s\n", *region)
	if *instanceID == "" {
		fmt.Println("InstanceID: <not set>")
		fmt.Println("Note: instance-scoped readonly calls will be skipped.")
	} else {
		fmt.Printf("InstanceID: %s\n", *instanceID)
	}
	fmt.Println()

	c, err := client.NewSSEMCPClient(*url, security.MCPClientOptionsFromEnv()...)
	must("create SSE client", err)

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
