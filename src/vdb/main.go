package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vdb/v20230616"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-vdb"
	mcpsvr := server.NewMCPServer(
		"腾讯云 TSE MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	vdbDescribeInstances := mcp.NewTool(
		"vdb-DescribeInstances",
		mcp.WithDescription(`查询实例列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("实例ID数组。"),
		),
		mcp.WithArray(
			"InstanceNames",
			mcp.Description("实例名称，支持模糊搜索。"),
		),
		mcp.WithArray(
			"InstanceKeys",
			mcp.Description("实例模糊搜索字段。"),
		),
		mcp.WithArray(
			"Status",
			mcp.Description("根据状态获取实例， 为空则获取全部非隔离和非下线的实例。"),
		),
		mcp.WithArray(
			"EngineVersions",
			mcp.Description("按照版本筛选实例。"),
		),
		mcp.WithArray(
			"ApiVersions",
			mcp.Description("按照api版本筛选实例"),
		),
		mcp.WithString(
			"CreateAt",
			mcp.Description("按照创建时间筛选实例。"),
		),
		mcp.WithString(
			"Zones",
			mcp.Description("按照可用区筛选实例。"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序字段。"),
		),
		mcp.WithString(
			"OrderDirection",
			mcp.Description("排序方式。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询开始位置。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("列表查询数量。"),
		),
		mcp.WithArray(
			"ResourceTags",
			mcp.Description("按照标签筛选实例"),
		),
		mcp.WithArray(
			"TaskStatus",
			mcp.Description("任务状态：1-待执行任务；2-密钥更新中；3-网络变更中；4-参数变更中；5-embedding变更中；6-ai套件变更中；7-滚动升级中；8-纵向扩容中；9-纵向缩容中；10-横向扩容中；11-横向缩容中"),
		),
		mcp.WithArray(
			"Networks",
			mcp.Description("根据实例vip搜索实例"),
		),
	)
	mcpsvr.AddTool(vdbDescribeInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vdb.NewClient(credential, region_, cpf)
		req := vdb.NewDescribeInstancesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstances(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vdbDescribeDBSecurityGroups := mcp.NewTool(
		"vdb-DescribeDBSecurityGroups",
		mcp.WithDescription(`查询实例安全组详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例ID，格式如：vdb-c1nl9***。"),
		),
	)
	mcpsvr.AddTool(vdbDescribeDBSecurityGroups, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vdb.NewClient(credential, region_, cpf)
		req := vdb.NewDescribeDBSecurityGroupsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDBSecurityGroups(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vdbDescribeInstanceMaintenanceWindow := mcp.NewTool(
		"vdb-DescribeInstanceMaintenanceWindow",
		mcp.WithDescription(`查询维护时间窗`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("指定查询维护时间窗的具体实例 ID。"),
		),
	)
	mcpsvr.AddTool(vdbDescribeInstanceMaintenanceWindow, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vdb.NewClient(credential, region_, cpf)
		req := vdb.NewDescribeInstanceMaintenanceWindowRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceMaintenanceWindow(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	sseEndpoint := getEnv("MCP_SERVER_SSE_ENDPOINT", "/sse")
	messageEndpoint := getEnv("MCP_SERVER_MESSAGE_ENDPOINT", "/message")
	ssePort := getEnv("MCP_SERVER_SSE_PORT", "9000")

	sseServer := server.NewSSEServer(mcpsvr,
		server.WithSSEEndpoint(sseEndpoint),
		server.WithMessageEndpoint(messageEndpoint),
		server.WithAppendQueryToMessageEndpoint())

	log.Printf("SSE server listening on :" + ssePort)
	serverURL := fmt.Sprintf("http://127.0.0.1:%s%s", ssePort, sseEndpoint)
	outputMCPServerConfig(mcpServerName, serverURL)
	if err := sseServer.Start(":" + ssePort); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func outputMCPServerConfig(mcpServerName, serverURL string) {
	config := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			mcpServerName: map[string]interface{}{
				"url":  serverURL,
				"type": "sse",
			},
		},
	}

	jsonOutput, _ := json.MarshalIndent(config, "", " ")
	fmt.Println("=== MCP Server Configuration ===")
	fmt.Println("Copy the following configuration to your MCP client:")
	fmt.Println()
	fmt.Println(string(jsonOutput))
	fmt.Println()
	fmt.Println("The server is now ready to accept connections.")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
