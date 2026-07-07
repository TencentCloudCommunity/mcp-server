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
	url := flag.String("url", envOrDefault("VERIFY_SSE_URL", "http://127.0.0.1:9000/sse"), "MCP SSE URL")
	region := flag.String("region", envOrDefault("VERIFY_REGION", "ap-guangzhou"), "region for readonly tool calls")
	instanceID := flag.String("instance-id", envOrDefault("VERIFY_INSTANCE_ID", ""), "instance id for instance-scoped readonly tool calls")
	flag.Parse()

	if strings.TrimSpace(*instanceID) == "" {
		fmt.Println("missing instance id: set VERIFY_INSTANCE_ID or pass --instance-id")
		os.Exit(1)
	}

	c, err := client.NewSSEMCPClient(*url, security.MCPClientOptionsFromEnv()...)
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

	// 1-5 instance group
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

	// parameter group
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

	// ssl
	out, err = call("DescribeDBInstanceSSLConfig", base)
	record("DescribeDBInstanceSSLConfig", out, err)

	// account
	out, err = call("DescribeAccounts", base)
	record("DescribeAccounts", out, err)
	userName := extractFirst(out, "UserName")
	if userName == "" {
		userName = "postgres"
	}

	// UserName + DatabaseObjectSet 均为查询权限的必填参数
	out, err = call("DescribeAccountPrivileges", map[string]interface{}{
		"region":       *region,
		"DBInstanceId": *instanceID,
		"UserName":     userName,
		"DatabaseObjectSet": []map[string]interface{}{
			{"ObjectType": "database", "ObjectName": "postgres"},
		},
	})
	record("DescribeAccountPrivileges", out, err)

	// network
	out, err = call("DescribeDBInstanceSecurityGroups", base)
	record("DescribeDBInstanceSecurityGroups", out, err)

	// monitoring
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

	// database
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

	// backup
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

	// BackupSetId 与 RecoveryTargetTime 必须二选一传入，这里优先用真实的基础备份集ID
	cloneArgs := map[string]interface{}{"region": *region, "DBInstanceId": *instanceID}
	if backupSetID != "" {
		cloneArgs["BackupSetId"] = backupSetID
	} else {
		cloneArgs["RecoveryTargetTime"] = time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")
	}
	out, err = call("DescribeCloneDBInstanceSpec", cloneArgs)
	record("DescribeCloneDBInstanceSpec", out, err)

	// readonly
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

// extractFirst 从 JSON 文本里粗略提取第一个字段值（仅用于验证脚本内部串联参数，非生产代码）
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
	instances, ok := resp["DBInstanceSet"].([]interface{})
	if !ok {
		return ""
	}
	for _, item := range instances {
		instance, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if instanceID != "" && stringValue(instance["DBInstanceId"]) != instanceID {
			continue
		}
		return stringValue(instance[field])
	}
	return ""
}

func stringValue(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	case float64:
		if x == float64(int64(x)) {
			return fmt.Sprintf("%d", int64(x))
		}
		return fmt.Sprintf("%v", x)
	default:
		return ""
	}
}

func envOrDefault(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
