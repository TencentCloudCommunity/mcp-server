package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-cdwch"
	mcpsvr := server.NewMCPServer(
		"腾讯云 CDWCH MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	cdwchDescribeInstanceShards := mcp.NewTool(
		"cdwch-DescribeInstanceShards",
		mcp.WithDescription(`获取实例shard信息列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群实例ID"),
		),
	)
	mcpsvr.AddTool(cdwchDescribeInstanceShards, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdwch.NewClient(credential, region_, cpf)
		req := cdwch.NewDescribeInstanceShardsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceShards(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cdwchDescribeInstancesNew := mcp.NewTool(
		"cdwch-DescribeInstancesNew",
		mcp.WithDescription(`获取实例简单信息列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"SearchInstanceId",
			mcp.Description("搜索的集群id名称"),
		),
		mcp.WithString(
			"SearchInstanceName",
			mcp.Description("搜索的集群name"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页参数，第一页为0，第二页为10"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页参数，分页步长，默认为10"),
		),
		mcp.WithArray(
			"SearchTags",
			mcp.Description("搜索标签列表"),
		),
		mcp.WithBoolean(
			"IsSimple",
			mcp.Description("信息详细与否"),
		),
		mcp.WithArray(
			"Vips",
			mcp.Description("vip列表"),
		),
	)
	mcpsvr.AddTool(cdwchDescribeInstancesNew, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdwch.NewClient(credential, region_, cpf)
		req := cdwch.NewDescribeInstancesNewRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstancesNew(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cdwchDescribeInstanceNodes := mcp.NewTool(
		"cdwch-DescribeInstanceNodes",
		mcp.WithDescription(`获取实例节点信息列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群实例ID"),
		),
		mcp.WithString(
			"NodeRole",
			mcp.Description("集群角色类型，“DATA” 为数据节点、“COMMON” 为 ZooKeeper 节点，默认为 &quot;DATA&quot; 数据节点。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页参数，第一页为0，第二页为10"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页参数，分页步长，默认为10"),
		),
		mcp.WithString(
			"DisplayPolicy",
			mcp.Description("展现策略，All时显示所有"),
		),
		mcp.WithBoolean(
			"ForceAll",
			mcp.Description("当true的时候返回所有节点，即Limit无限大"),
		),
	)
	mcpsvr.AddTool(cdwchDescribeInstanceNodes, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdwch.NewClient(credential, region_, cpf)
		req := cdwch.NewDescribeInstanceNodesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceNodes(req)
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
