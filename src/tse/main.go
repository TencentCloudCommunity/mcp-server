package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-tse"
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

	tseDescribeSREInstances := mcp.NewTool(
		"tse-DescribeSREInstances",
		mcp.WithDescription(`查询引擎实例列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("请求过滤参数"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("翻页单页查询限制数量[0,1000], 默认值0"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("翻页单页偏移量，默认值0"),
		),
		mcp.WithString(
			"QueryType",
			mcp.Description("查询类型"),
		),
		mcp.WithString(
			"QuerySource",
			mcp.Description("调用方来源"),
		),
	)
	mcpsvr.AddTool(tseDescribeSREInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeSREInstancesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSREInstances(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tseDescribeSREInstanceAccessAddress := mcp.NewTool(
		"tse-DescribeSREInstanceAccessAddress",
		mcp.WithDescription(`查询引擎实例访问地址`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("注册引擎实例Id"),
		),
		mcp.WithString(
			"VpcId",
			mcp.Description("VPC ID"),
		),
		mcp.WithString(
			"SubnetId",
			mcp.Description("子网ID"),
		),
		mcp.WithString(
			"Workload",
			mcp.Description("引擎其他组件名称（pushgateway、polaris-limiter）"),
		),
		mcp.WithString(
			"EngineRegion",
			mcp.Description("部署地域"),
		),
	)
	mcpsvr.AddTool(tseDescribeSREInstanceAccessAddress, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeSREInstanceAccessAddressRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSREInstanceAccessAddress(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tseDescribeNacosReplicas := mcp.NewTool(
		"tse-DescribeNacosReplicas",
		mcp.WithDescription(`查询Nacos类型引擎实例副本信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("引擎实例ID"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("副本列表Limit"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("副本列表Offset"),
		),
	)
	mcpsvr.AddTool(tseDescribeNacosReplicas, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeNacosReplicasRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeNacosReplicas(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tseDescribeNacosServerInterfaces := mcp.NewTool(
		"tse-DescribeNacosServerInterfaces",
		mcp.WithDescription(`查询nacos服务接口列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例id"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回的列表个数"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("返回的列表起始偏移量"),
		),
	)
	mcpsvr.AddTool(tseDescribeNacosServerInterfaces, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeNacosServerInterfacesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeNacosServerInterfaces(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tseDescribeZookeeperReplicas := mcp.NewTool(
		"tse-DescribeZookeeperReplicas",
		mcp.WithDescription(`查询Zookeeper类型注册引擎实例副本信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("注册引擎实例ID"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("副本列表Limit"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("副本列表Offset"),
		),
	)
	mcpsvr.AddTool(tseDescribeZookeeperReplicas, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeZookeeperReplicasRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeZookeeperReplicas(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tseDescribeZookeeperServerInterfaces := mcp.NewTool(
		"tse-DescribeZookeeperServerInterfaces",
		mcp.WithDescription(`查询zookeeper服务接口列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例id"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回的列表个数"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("返回的列表起始偏移量"),
		),
	)
	mcpsvr.AddTool(tseDescribeZookeeperServerInterfaces, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeZookeeperServerInterfacesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeZookeeperServerInterfaces(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tseDescribeCloudNativeAPIGateway := mcp.NewTool(
		"tse-DescribeCloudNativeAPIGateway",
		mcp.WithDescription(`获取云原生API网关实例信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"GatewayId",
			mcp.Description("云原生API网关实例ID"),
		),
	)
	mcpsvr.AddTool(tseDescribeCloudNativeAPIGateway, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeCloudNativeAPIGatewayRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeCloudNativeAPIGateway(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tseDescribeCloudNativeAPIGatewayConfig := mcp.NewTool(
		"tse-DescribeCloudNativeAPIGatewayConfig",
		mcp.WithDescription(`获取云原生API网关实例网络配置信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"GatewayId",
			mcp.Description("云原生API网关实例ID。"),
		),
		mcp.WithString(
			"GroupId",
			mcp.Description("分组id，不填时为默认分组"),
		),
	)
	mcpsvr.AddTool(tseDescribeCloudNativeAPIGatewayConfig, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeCloudNativeAPIGatewayConfigRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeCloudNativeAPIGatewayConfig(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tseDescribeNativeGatewayServerGroups := mcp.NewTool(
		"tse-DescribeNativeGatewayServerGroups",
		mcp.WithDescription(`查询云原生网关分组信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"GatewayId",
			mcp.Description("云原生API网关实例ID。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认为 0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为 20。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤参数，支持按照分组名称、分组ID（Name、GroupId）筛选"),
		),
	)
	mcpsvr.AddTool(tseDescribeNativeGatewayServerGroups, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeNativeGatewayServerGroupsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeNativeGatewayServerGroups(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tseDescribeNativeGatewayServiceSources := mcp.NewTool(
		"tse-DescribeNativeGatewayServiceSources",
		mcp.WithDescription(`查询网关服务来源实例列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"GatewayID",
			mcp.Description("网关实例ID"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("单页条数，最大100"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页偏移量"),
		),
		mcp.WithString(
			"SourceName",
			mcp.Description("服务来源实例名称，模糊搜索"),
		),
		mcp.WithArray(
			"SourceTypes",
			mcp.Description("微服务引擎类型：TSE-Nacos｜TSE-Consul｜TSE-PolarisMesh｜Customer-Nacos｜Customer-Consul｜Customer-PolarisMesh"),
		),
		mcp.WithString(
			"OrderField",
			mcp.Description("排序字段类型，当前仅支持SourceName"),
		),
		mcp.WithString(
			"OrderType",
			mcp.Description("排序类型，AES/DESC"),
		),
	)
	mcpsvr.AddTool(tseDescribeNativeGatewayServiceSources, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeNativeGatewayServiceSourcesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeNativeGatewayServiceSources(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tseDescribeCloudNativeAPIGatewayServices := mcp.NewTool(
		"tse-DescribeCloudNativeAPIGatewayServices",
		mcp.WithDescription(`查询云原生网关服务列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"GatewayId",
			mcp.Description("网关ID"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("列表数量"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("列表 offset"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，多个过滤条件之间是与的关系，支持 name,upstreamType"),
		),
	)
	mcpsvr.AddTool(tseDescribeCloudNativeAPIGatewayServices, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeCloudNativeAPIGatewayServicesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeCloudNativeAPIGatewayServices(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tseDescribeCloudNativeAPIGatewayRoutes := mcp.NewTool(
		"tse-DescribeCloudNativeAPIGatewayRoutes",
		mcp.WithDescription(`查询云原生网关路由列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"GatewayId",
			mcp.Description("网关ID"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("翻页单页查询限制数量[0,1000], 默认值0"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("翻页单页偏移量，默认值0"),
		),
		mcp.WithString(
			"ServiceName",
			mcp.Description("服务的名字，精确匹配"),
		),
		mcp.WithString(
			"RouteName",
			mcp.Description("路由的名字，精确匹配"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，多个过滤条件之间是与的关系，支持 name, path, host, method, service, protocol"),
		),
	)
	mcpsvr.AddTool(tseDescribeCloudNativeAPIGatewayRoutes, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tse.NewClient(credential, region_, cpf)
		req := tse.NewDescribeCloudNativeAPIGatewayRoutesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeCloudNativeAPIGatewayRoutes(req)
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
