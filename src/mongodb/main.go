package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-mongodb"
	mcpsvr := server.NewMCPServer(
		"腾讯云 MongoDB MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	mongodbDescribeDBInstances := mcp.NewTool(
		"mongodb-DescribeDBInstances",
		mcp.WithDescription(`查询云数据库实例列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("实例ID列表"),
		),
		mcp.WithNumber(
			"InstanceType",
			mcp.Description("实例类型（0:所有实例, 1:正式实例, 3:只读实例, 4:灾备实例）"),
		),
		mcp.WithNumber(
			"ClusterType",
			mcp.Description("集群类型（0:副本集实例, 1:分片实例, -1:副本集与分片实例）"),
		),
		mcp.WithArray(
			"Status",
			mcp.Description("实例状态（0:待初始化, 1:流程处理中, 2:正常运行, -2:已过期）"),
		),
		mcp.WithString(
			"VpcId",
			mcp.Description("私有网络ID"),
		),
		mcp.WithString(
			"SubnetId",
			mcp.Description("私有网络子网ID"),
		),
		mcp.WithNumber(
			"PayMode",
			mcp.Description("付费类型（0:按量计费, 1:包年包月, -1:所有）"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("单次返回数量，默认20，范围[1,100]"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认0"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序字段（ProjectId、InstanceName、CreateTime）"),
		),
		mcp.WithString(
			"OrderByType",
			mcp.Description("排序方式（ASC:升序, DESC:降序）"),
		),
		mcp.WithArray(
			"ProjectIds",
			mcp.Description("项目ID列表"),
		),
		mcp.WithString(
			"SearchKey",
			mcp.Description("搜索关键词（实例ID、实例名称或内网IP）"),
		),
		mcp.WithArray(
			"Tags",
			mcp.Description("标签信息"),
		),
	)
	mcpsvr.AddTool(mongodbDescribeDBInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := mongodb.NewClient(credential, region_, cpf)
		req := mongodb.NewDescribeDBInstancesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDBInstances(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	mongodbDescribeDBBackups := mcp.NewTool(
		"mongodb-DescribeDBBackups",
		mcp.WithDescription(`查询实例备份列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Required(),
			mcp.Description("实例ID"),
		),
		mcp.WithNumber(
			"BackupMethod",
			mcp.Description("备份方式（0:逻辑备份, 1:物理备份, 2:所有备份），默认逻辑备份"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页大小，最大100，默认查询所有"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页偏移量，最小0，默认0"),
		),
	)
	mcpsvr.AddTool(mongodbDescribeDBBackups, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := mongodb.NewClient(credential, region_, cpf)
		req := mongodb.NewDescribeDBBackupsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDBBackups(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	mongodbDescribeInstanceParams := mcp.NewTool(
		"mongodb-DescribeInstanceParams",
		mcp.WithDescription(`获取当前实例可修改的参数列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Required(),
			mcp.Description("实例ID"),
		),
	)
	mcpsvr.AddTool(mongodbDescribeInstanceParams, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := mongodb.NewClient(credential, region_, cpf)
		req := mongodb.NewDescribeInstanceParamsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceParams(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	// DescribeSlowLogs - 获取慢日志信息
	mongodbDescribeSlowLogs := mcp.NewTool(
		"mongodb-DescribeSlowLogs",
		mcp.WithDescription(`获取慢日志信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Required(),
			mcp.Description("实例ID"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Required(),
			mcp.Description("慢日志起始时间（格式：yyyy-mm-dd hh:mm:ss）"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Required(),
			mcp.Description("慢日志终止时间（格式：yyyy-mm-dd hh:mm:ss）"),
		),
		mcp.WithNumber(
			"SlowMS",
			mcp.Required(),
			mcp.Description("慢日志执行时间阈值，单位毫秒，最小100毫秒"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，最小0，最大10000，默认0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页大小，最小1，最大100，默认20"),
		),
		mcp.WithString(
			"Format",
			mcp.Description("慢日志返回格式，默认原生格式，4.4及以上版本可设置为json"),
		),
	)
	mcpsvr.AddTool(mongodbDescribeSlowLogs, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := mongodb.NewClient(credential, region_, cpf)
		req := mongodb.NewDescribeSlowLogsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSlowLogs(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	mongodbDescribeSpecInfo := mcp.NewTool(
		"mongodb-DescribeSpecInfo",
		mcp.WithDescription(`查询云数据库的售卖规格`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Zone",
			mcp.Description("待查询可用区"),
		),
	)
	mcpsvr.AddTool(mongodbDescribeSpecInfo, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := mongodb.NewClient(credential, region_, cpf)
		req := mongodb.NewDescribeSpecInfoRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSpecInfo(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	mongodbDescribeSecurityGroup := mcp.NewTool(
		"mongodb-DescribeSecurityGroup",
		mcp.WithDescription(`查询实例绑定的安全组`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Required(),
			mcp.Description("实例ID"),
		),
	)
	mcpsvr.AddTool(mongodbDescribeSecurityGroup, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := mongodb.NewClient(credential, region_, cpf)
		req := mongodb.NewDescribeSecurityGroupRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSecurityGroup(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
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
