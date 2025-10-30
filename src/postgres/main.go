package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-postgres"
	mcpsvr := server.NewMCPServer(
		"腾讯云 Postgres MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	postgresDescribeDBInstanceSSLConfig := mcp.NewTool(
		"postgres-DescribeDBInstanceSSLConfig",
		mcp.WithDescription(`查询实例SSL配置`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DBInstanceId",
			mcp.Description("实例ID，形如postgres-6bwgamo3"),
		),
	)
	mcpsvr.AddTool(postgresDescribeDBInstanceSSLConfig, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := postgres.NewClient(credential, region_, cpf)
		req := postgres.NewDescribeDBInstanceSSLConfigRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDBInstanceSSLConfig(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	postgresDescribeDBInstanceAttribute := mcp.NewTool(
		"postgres-DescribeDBInstanceAttribute",
		mcp.WithDescription(`查询实例详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DBInstanceId",
			mcp.Description("实例ID"),
		),
	)
	mcpsvr.AddTool(postgresDescribeDBInstanceAttribute, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := postgres.NewClient(credential, region_, cpf)
		req := postgres.NewDescribeDBInstanceAttributeRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDBInstanceAttribute(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	postgresUpgradeDBInstanceKernelVersion := mcp.NewTool(
		"postgres-UpgradeDBInstanceKernelVersion",
		mcp.WithDescription(`升级实例内核版本号`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DBInstanceId",
			mcp.Description("实例ID。"),
		),
		mcp.WithString(
			"TargetDBKernelVersion",
			mcp.Description("升级的目标内核版本号。可以通过接口[DescribeDBVersions](https://cloud.tencent.com/document/api/409/89018)的返回字段AvailableUpgradeTarget获取。"),
		),
		mcp.WithNumber(
			"SwitchTag",
			mcp.Description("指定实例升级内核版本号完成后的切换时间。可选值:"),
		),
		mcp.WithString(
			"SwitchStartTime",
			mcp.Description("切换开始时间，时间格式：HH:MM:SS，例如：01:00:00。当SwitchTag为0或2时，该参数失效。"),
		),
		mcp.WithString(
			"SwitchEndTime",
			mcp.Description("切换截止时间，时间格式：HH:MM:SS，例如：01:30:00。当SwitchTag为0或2时，该参数失效。SwitchStartTime和SwitchEndTime时间窗口不能小于30分钟。"),
		),
		mcp.WithBoolean(
			"DryRun",
			mcp.Description("是否对本次升级实例内核版本号操作执行预检查。"),
		),
	)
	mcpsvr.AddTool(postgresUpgradeDBInstanceKernelVersion, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := postgres.NewClient(credential, region_, cpf)
		req := postgres.NewUpgradeDBInstanceKernelVersionRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.UpgradeDBInstanceKernelVersion(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	postgresDescribeAccounts := mcp.NewTool(
		"postgres-DescribeAccounts",
		mcp.WithDescription(`查询实例的数据库账号列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DBInstanceId",
			mcp.Description("实例ID，形如postgres-6fego161"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页返回，每页最大返回数目，默认20，取值范围为1-100"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("数据偏移量，从0开始。"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("返回数据按照创建时间或者用户名排序。取值支持createTime、name、updateTime。createTime-按照创建时间排序；name-按照用户名排序; updateTime-按照更新时间排序。"),
		),
		mcp.WithString(
			"OrderByType",
			mcp.Description("返回结果是升序还是降序。取值只能为desc或者asc。desc-降序；asc-升序"),
		),
	)
	mcpsvr.AddTool(postgresDescribeAccounts, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := postgres.NewClient(credential, region_, cpf)
		req := postgres.NewDescribeAccountsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeAccounts(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	postgresDescribeDatabases := mcp.NewTool(
		"postgres-DescribeDatabases",
		mcp.WithDescription(`查询实例的数据库列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DBInstanceId",
			mcp.Description("实例ID"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("按照一个或者多个过滤条件进行查询，目前支持的过滤条件有：database-name：按照数据库名称过滤，类型为string。此处使用模糊匹配搜索符合条件的数据库。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("数据偏移量，从0开始。	"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("单次显示数量"),
		),
	)
	mcpsvr.AddTool(postgresDescribeDatabases, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := postgres.NewClient(credential, region_, cpf)
		req := postgres.NewDescribeDatabasesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDatabases(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	postgresDescribeDBInstanceParameters := mcp.NewTool(
		"postgres-DescribeDBInstanceParameters",
		mcp.WithDescription(`查询实例参数`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DBInstanceId",
			mcp.Description("实例ID"),
		),
		mcp.WithString(
			"ParamName",
			mcp.Description("查询指定参数详情。ParamName为空或不传，默认返回全部参数列表"),
		),
	)
	mcpsvr.AddTool(postgresDescribeDBInstanceParameters, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := postgres.NewClient(credential, region_, cpf)
		req := postgres.NewDescribeDBInstanceParametersRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDBInstanceParameters(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	postgresDescribeParameterTemplates := mcp.NewTool(
		"postgres-DescribeParameterTemplates",
		mcp.WithDescription(`查询参数模板列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，目前支持的过滤条件有：TemplateName, TemplateId，DBMajorVersion，DBEngine"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("每页显示数量，[0，100]，默认 20"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("数据偏移量"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序指标，枚举值，支持：CreateTime，TemplateName，DBMajorVersion"),
		),
		mcp.WithString(
			"OrderByType",
			mcp.Description("排序方式，枚举值，支持：asc（升序） ，desc（降序）"),
		),
	)
	mcpsvr.AddTool(postgresDescribeParameterTemplates, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := postgres.NewClient(credential, region_, cpf)
		req := postgres.NewDescribeParameterTemplatesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeParameterTemplates(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	postgresDescribeParameterTemplateAttributes := mcp.NewTool(
		"postgres-DescribeParameterTemplateAttributes",
		mcp.WithDescription(`查询参数模板详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"TemplateId",
			mcp.Description("参数模板ID"),
		),
	)
	mcpsvr.AddTool(postgresDescribeParameterTemplateAttributes, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := postgres.NewClient(credential, region_, cpf)
		req := postgres.NewDescribeParameterTemplateAttributesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeParameterTemplateAttributes(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	postgresCreateDatabase := mcp.NewTool(
		"postgres-CreateDatabase",
		mcp.WithDescription(`创建数据库`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DBInstanceId",
			mcp.Description("实例ID，形如postgres-6fego161"),
		),
		mcp.WithString(
			"DatabaseName",
			mcp.Description("创建的数据库名"),
		),
		mcp.WithString(
			"DatabaseOwner",
			mcp.Description("数据库的所有者"),
		),
		mcp.WithString(
			"Encoding",
			mcp.Description("数据库的字符编码"),
		),
		mcp.WithString(
			"Collate",
			mcp.Description("数据库的排序规则"),
		),
		mcp.WithString(
			"Ctype",
			mcp.Description("数据库的字符分类"),
		),
	)
	mcpsvr.AddTool(postgresCreateDatabase, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := postgres.NewClient(credential, region_, cpf)
		req := postgres.NewCreateDatabaseRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateDatabase(req)
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
