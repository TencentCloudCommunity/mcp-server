package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
)

func main() {
	mcpServerName := "mcp-server-cynosdb"
	mcpsvr := server.NewMCPServer(
		"腾讯云 TDSQL-C MySQL MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	cynosdbAssociateSecurityGroups := mcp.NewTool(
		"cynosdb-AssociateSecurityGroups",
		mcp.WithDescription(`安全组批量绑定云资源`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("实例组 ID 数组，cynosdbmysql-grp-前缀开头或集群 ID。"),
		),
		mcp.WithArray(
			"SecurityGroupIds",
			mcp.Description("要修改的安全组ID列表，一个或者多个安全组Id组成的数组。"),
		),
		mcp.WithString(
			"Zone",
			mcp.Description("可用区"),
		),
	)
	mcpsvr.AddTool(cynosdbAssociateSecurityGroups, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cynosdb.NewClient(credential, region_, cpf)
		req := cynosdb.NewAssociateSecurityGroupsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.AssociateSecurityGroups(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cynosdbDescribeAccounts := mcp.NewTool(
		"cynosdb-DescribeAccounts",
		mcp.WithDescription(`查询数据库账号列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ClusterId",
			mcp.Description("集群ID"),
		),
		mcp.WithArray(
			"AccountNames",
			mcp.Description("需要过滤的账户列表"),
		),
		mcp.WithString(
			"DbType",
			mcp.Description("数据库类型，取值范围: "),
		),
		mcp.WithArray(
			"Hosts",
			mcp.Description("需要过滤的账户列表"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("限制量"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量"),
		),
		mcp.WithString(
			"AccountRegular",
			mcp.Description("模糊匹配关键字(同时匹配AccountName和AccountHost，返回并集结果，支持正则)"),
		),
	)
	mcpsvr.AddTool(cynosdbDescribeAccounts, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cynosdb.NewClient(credential, region_, cpf)
		req := cynosdb.NewDescribeAccountsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeAccounts(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cynosdbDescribeClusterDatabases := mcp.NewTool(
		"cynosdb-DescribeClusterDatabases",
		mcp.WithDescription(`获取集群数据库列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ClusterId",
			mcp.Description("集群id"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页偏移"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页限制数量"),
		),
	)
	mcpsvr.AddTool(cynosdbDescribeClusterDatabases, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cynosdb.NewClient(credential, region_, cpf)
		req := cynosdb.NewDescribeClusterDatabasesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeClusterDatabases(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cynosdbDescribeInstanceDetail := mcp.NewTool(
		"cynosdb-DescribeInstanceDetail",
		mcp.WithDescription(`查询实例详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例ID"),
		),
	)
	mcpsvr.AddTool(cynosdbDescribeInstanceDetail, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cynosdb.NewClient(credential, region_, cpf)
		req := cynosdb.NewDescribeInstanceDetailRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceDetail(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cynosdbDescribeInstances := mcp.NewTool(
		"cynosdb-DescribeInstances",
		mcp.WithDescription(`查询实例的列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为 20，取值范围为(0,100]"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("记录偏移量，默认值为0"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序字段，取值范围："),
		),
		mcp.WithString(
			"OrderByType",
			mcp.Description("排序类型，取值范围："),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("搜索条件，若存在多个Filter时，Filter间的关系为逻辑与（AND）关系。"),
		),
		mcp.WithString(
			"DbType",
			mcp.Description("引擎类型：目前支持“MYSQL”"),
		),
		mcp.WithString(
			"Status",
			mcp.Description("实例状态, 可选值:"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("实例id列表"),
		),
	)
	mcpsvr.AddTool(cynosdbDescribeInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cynosdb.NewClient(credential, region_, cpf)
		req := cynosdb.NewDescribeInstancesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstances(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cynosdbModifyInstanceName := mcp.NewTool(
		"cynosdb-ModifyInstanceName",
		mcp.WithDescription(`修改实例名称`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例ID"),
		),
		mcp.WithString(
			"InstanceName",
			mcp.Description("实例名称"),
		),
	)
	mcpsvr.AddTool(cynosdbModifyInstanceName, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cynosdb.NewClient(credential, region_, cpf)
		req := cynosdb.NewModifyInstanceNameRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ModifyInstanceName(req)
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
