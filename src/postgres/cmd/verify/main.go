package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"postgres_server/security"
)

func main() {
	transportMode := flag.String("transport", normalizeTransport(envOrDefault("VERIFY_TRANSPORT", envOrDefault("MCP_TRANSPORT", "streamable-http"))), "transport: streamable-http | sse | stdio")
	url := flag.String("url", defaultURL(*transportMode, "VERIFY_SERVER_URL", "VERIFY_SSE_URL"), "MCP URL for HTTP/SSE transports")
	stdioCommand := flag.String("command", envOrDefault("VERIFY_STDIO_COMMAND", ""), "stdio command path")
	region := flag.String("region", envOrDefault("VERIFY_REGION", "ap-guangzhou"), "region for readonly tool calls")
	instanceID := flag.String("instance-id", envOrDefault("VERIFY_INSTANCE_ID", ""), "instance id for instance-scoped readonly tool calls")
	flag.Parse()

	mode := normalizeTransport(*transportMode)
	*transportMode = mode

	if strings.TrimSpace(*instanceID) == "" {
		fmt.Println("missing instance id: set VERIFY_INSTANCE_ID or pass --instance-id")
		os.Exit(1)
	}

	c, err := newMCPClient(mode, *url, *stdioCommand)
	if err != nil {
		fmt.Println("new client error:", err)
		os.Exit(1)
	}
	ctx := context.Background()
	if err := c.Start(ctx); err != nil {
		fmt.Println("start error:", err)
		os.Exit(1)
	}
	defer c.Close()

	_, err = c.Initialize(ctx, mcp.InitializeRequest{})
	if err != nil {
		fmt.Println("initialize error:", err)
		os.Exit(1)
	}

	call := func(name string, args map[string]interface{}) (string, error) {
		req := mcp.CallToolRequest{}
		req.Params.Name = "postgres-" + name
		req.Params.Arguments = args
		cctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		res, err := c.CallTool(cctx, req)
		if err != nil {
			return "", err
		}
		var sb strings.Builder
		for _, content := range res.Content {
			if tc, ok := content.(mcp.TextContent); ok {
				sb.WriteString(tc.Text)
			}
		}
		return sb.String(), nil
	}

	results := map[string]string{}
	order := []string{}
	record := func(name, out string, err error) {
		order = append(order, name)
		if err != nil {
			results[name] = "TRANSPORT_ERROR: " + err.Error()
			return
		}
		results[name] = out
	}

	base := map[string]interface{}{"region": *region, "DBInstanceId": *instanceID}
	noID := map[string]interface{}{"region": *region}

	out, err := call("DescribeDBInstanceAttribute", base)
	record("DescribeDBInstanceAttribute", out, err)

	out, err = call("DescribeDBInstances", noID)
	record("DescribeDBInstances", out, err)

	zone := extractInstanceField(out, *instanceID, "Zone")
	majorVersion := extractInstanceField(out, *instanceID, "DBMajorVersion")
	if zone != "" && majorVersion != "" {
		out, err = call("DescribeClasses", map[string]interface{}{"region": *region, "Zone": zone, "DBEngine": "postgresql", "DBMajorVersion": majorVersion})
		record("DescribeClasses", out, err)
	} else {
		record("DescribeClasses", fmt.Sprintf("SKIPPED: missing Zone/DBMajorVersion for instance %s", *instanceID), nil)
	}

	out, err = call("DescribeDBVersions", noID)
	record("DescribeDBVersions", out, err)

	out, err = call("DescribeTasks", base)
	record("DescribeTasks", out, err)

	out, err = call("DescribeRegions", noID)
	record("DescribeRegions", out, err)

	out, err = call("DescribeZones", noID)
	record("DescribeZones", out, err)

	out, err = call("DescribeProductConfig", map[string]interface{}{"region": *region, "DBEngine": "postgresql"})
	record("DescribeProductConfig", out, err)

	out, err = call("DescribeDBInstanceParameters", base)
	record("DescribeDBInstanceParameters", out, err)

	out, err = call("DescribeParamsEvent", base)
	record("DescribeParamsEvent", out, err)

	out, err = call("DescribeParameterTemplates", noID)
	record("DescribeParameterTemplates", out, err)
	templateID := extractFirst(out, "TemplateId")

	if templateID != "" {
		out, err = call("DescribeParameterTemplateAttributes", map[string]interface{}{"region": *region, "TemplateId": templateID})
	} else {
		out, err = call("DescribeParameterTemplateAttributes", map[string]interface{}{"region": *region, "TemplateId": "notfound"})
	}
	record("DescribeParameterTemplateAttributes", out, err)

	out, err = call("DescribeDBInstanceSSLConfig", base)
	record("DescribeDBInstanceSSLConfig", out, err)

	out, err = call("DescribeAccounts", base)
	record("DescribeAccounts", out, err)
	userName := extractFirst(out, "UserName")
	if userName == "" {
		userName = "postgres"
	}

	out, err = call("DescribeAccountPrivileges", map[string]interface{}{
		"region":       *region,
		"DBInstanceId": *instanceID,
		"UserName":     userName,
		"DatabaseObjectSet": []map[string]interface{}{
			{"ObjectType": "database", "ObjectName": "postgres"},
		},
	})
	record("DescribeAccountPrivileges", out, err)

	out, err = call("DescribeDBInstanceSecurityGroups", base)
	record("DescribeDBInstanceSecurityGroups", out, err)

	now := time.Now()
	monArgs := map[string]interface{}{
		"region":       *region,
		"DBInstanceId": *instanceID,
		"StartTime":    now.Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
		"EndTime":      now.Format("2006-01-02 15:04:05"),
	}
	out, err = call("DescribeSlowQueryList", monArgs)
	record("DescribeSlowQueryList", out, err)

	out, err = call("DescribeSlowQueryAnalysis", monArgs)
	record("DescribeSlowQueryAnalysis", out, err)

	out, err = call("DescribeDBErrlogs", monArgs)
	record("DescribeDBErrlogs", out, err)

	out, err = call("DescribeDatabases", base)
	record("DescribeDatabases", out, err)
	dbName := extractFirst(out, "DBName")
	if dbName == "" {
		dbName = extractFirst(out, "DatabaseName")
	}

	if dbName != "" {
		out, err = call("DescribeDatabaseObjects", map[string]interface{}{"region": *region, "DBInstanceId": *instanceID, "DatabaseName": dbName, "ObjectType": "schema"})
	} else {
		out, err = call("DescribeDatabaseObjects", map[string]interface{}{"region": *region, "DBInstanceId": *instanceID, "DatabaseName": "postgres", "ObjectType": "schema"})
	}
	record("DescribeDatabaseObjects", out, err)

	out, err = call("DescribeBackupOverview", noID)
	record("DescribeBackupOverview", out, err)

	backupQuery := map[string]interface{}{
		"region": *region,
		"Filters": []map[string]interface{}{
			{"Name": "db-instance-id", "Values": []string{*instanceID}},
		},
		"Limit": 20,
	}
	out, err = call("DescribeBaseBackups", backupQuery)
	record("DescribeBaseBackups", out, err)
	backupSetID := extractFirst(out, "Id")

	out, err = call("DescribeLogBackups", backupQuery)
	record("DescribeLogBackups", out, err)

	out, err = call("DescribeAvailableRecoveryTime", base)
	record("DescribeAvailableRecoveryTime", out, err)

	cloneArgs := map[string]interface{}{"region": *region, "DBInstanceId": *instanceID}
	if backupSetID != "" {
		cloneArgs["BackupSetId"] = backupSetID
	} else {
		cloneArgs["RecoveryTargetTime"] = time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")
	}
	out, err = call("DescribeCloneDBInstanceSpec", cloneArgs)
	record("DescribeCloneDBInstanceSpec", out, err)

	out, err = call("DescribeReadOnlyGroups", base)
	record("DescribeReadOnlyGroups", out, err)

	fmt.Println("========== RESULTS ==========")
	for _, name := range order {
		out := results[name]
		if len(out) > 500 {
			out = out[:500] + "...(truncated)"
		}
		fmt.Printf("\n--- %s ---\n%s\n", name, out)
	}
}

func newMCPClient(mode, url, stdioCommand string) (*client.Client, error) {
	switch mode {
	case "sse":
		return client.NewSSEMCPClient(url, security.MCPClientOptionsFromEnv()...)
	case "stdio":
		if strings.TrimSpace(stdioCommand) == "" {
			return nil, fmt.Errorf("missing stdio command: set VERIFY_STDIO_COMMAND or pass --command")
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

func extractFirst(jsonText string, field string) string {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonText), &raw); err != nil {
		return ""
	}
	resp, ok := raw["Response"].(map[string]interface{})
	if !ok {
		return ""
	}
	for _, v := range resp {
		arr, ok := v.([]interface{})
		if !ok || len(arr) == 0 {
			continue
		}
		first, ok := arr[0].(map[string]interface{})
		if !ok {
			continue
		}
		if val, ok := first[field].(string); ok {
			return val
		}
	}
	return ""
}

func extractInstanceField(jsonText, instanceID, field string) string {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonText), &raw); err != nil {
		return ""
	}
	resp, ok := raw["Response"].(map[string]interface{})
	if !ok {
		return ""
	}
	set, ok := resp["DBInstanceSet"].([]interface{})
	if !ok {
		return ""
	}
	for _, item := range set {
		instance, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if id, _ := instance["DBInstanceId"].(string); id != instanceID {
			continue
		}
		if value, _ := instance[field].(string); value != "" {
			return value
		}
		if nested, ok := instance[field].(map[string]interface{}); ok {
			if name, _ := nested["Zone"].(string); name != "" {
				return name
			}
		}
	}
	return ""
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
