package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-clb"
	mcpsvr := server.NewMCPServer(
		"腾讯云 CLB MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	clbDescribeCustomizedConfigAssociateList := mcp.NewTool(
		"clb-DescribeCustomizedConfigAssociateList",
		mcp.WithDescription(`拉取配置绑定的server或location`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"UconfigId",
			mcp.Description("配置ID"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("拉取绑定关系列表开始位置，默认值 0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("拉取绑定关系列表数目，默认值 20"),
		),
		mcp.WithString(
			"Domain",
			mcp.Description("搜索域名"),
		),
	)
	mcpsvr.AddTool(clbDescribeCustomizedConfigAssociateList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := clb.NewClient(credential, region_, cpf)
		req := clb.NewDescribeCustomizedConfigAssociateListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeCustomizedConfigAssociateList(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	clbDescribeCustomizedConfigList := mcp.NewTool(
		"clb-DescribeCustomizedConfigList",
		mcp.WithDescription(`拉取配置列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ConfigType",
			mcp.Description("配置类型:CLB 负载均衡维度。 SERVER 域名维度。 LOCATION 规则维度。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("拉取页偏移，默认值0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("拉取数目，默认值20"),
		),
		mcp.WithString(
			"ConfigName",
			mcp.Description("拉取指定配置名字，模糊匹配。"),
		),
		mcp.WithArray(
			"UconfigIds",
			mcp.Description("配置ID"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件如下："),
		),
	)
	mcpsvr.AddTool(clbDescribeCustomizedConfigList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := clb.NewClient(credential, region_, cpf)
		req := clb.NewDescribeCustomizedConfigListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeCustomizedConfigList(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	clbDescribeListeners := mcp.NewTool(
		"clb-DescribeListeners",
		mcp.WithDescription(`查询负载均衡的监听器列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"LoadBalancerId",
			mcp.Description("负载均衡实例 ID。"),
		),
		mcp.WithArray(
			"ListenerIds",
			mcp.Description("要查询的负载均衡监听器 ID 数组，最大为100个。"),
		),
		mcp.WithString(
			"Protocol",
			mcp.Description("要查询的监听器协议类型，取值 TCP | UDP | HTTP | HTTPS | TCP_SSL | QUIC。"),
		),
		mcp.WithNumber(
			"Port",
			mcp.Description("要查询的监听器的端口。"),
		),
	)
	mcpsvr.AddTool(clbDescribeListeners, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := clb.NewClient(credential, region_, cpf)
		req := clb.NewDescribeListenersRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeListeners(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	clbDescribeLoadBalancersDetail := mcp.NewTool(
		"clb-DescribeLoadBalancersDetail",
		mcp.WithDescription(`查询负载均衡详细信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回负载均衡列表数目，默认20，最大值100。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("返回负载均衡列表起始偏移量，默认0。"),
		),
		mcp.WithArray(
			"Fields",
			mcp.Description("选择返回的Fields列表，系统仅会返回Fileds中填写的字段，可填写的字段详情请参见<a href=&quot;https://cloud.tencent.com/document/api/214/30694#LoadBalancerDetail&quot;>LoadBalancerDetail</a>。若未在Fileds填写相关字段，则此字段返回null。Fileds中默认添加LoadBalancerId和LoadBalancerName字段。"),
		),
		mcp.WithString(
			"TargetType",
			mcp.Description("当Fields包含TargetId、TargetAddress、TargetPort、TargetWeight、ListenerId、Protocol、Port、LocationId、Domain、Url等Fields时，必选选择导出目标组的Target或者非目标组Target，取值范围NODE、GROUP。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("查询负载均衡详细信息列表条件，详细的过滤条件如下："),
		),
	)
	mcpsvr.AddTool(clbDescribeLoadBalancersDetail, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := clb.NewClient(credential, region_, cpf)
		req := clb.NewDescribeLoadBalancersDetailRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeLoadBalancersDetail(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	clbDescribeTargetHealth := mcp.NewTool(
		"clb-DescribeTargetHealth",
		mcp.WithDescription(`获取负载均衡后端服务的健康检查状态`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"LoadBalancerIds",
			mcp.Description("要查询的负载均衡实例ID列表。"),
		),
	)
	mcpsvr.AddTool(clbDescribeTargetHealth, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := clb.NewClient(credential, region_, cpf)
		req := clb.NewDescribeTargetHealthRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTargetHealth(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	clbDescribeTargetGroups := mcp.NewTool(
		"clb-DescribeTargetGroups",
		mcp.WithDescription(`查询目标组信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"TargetGroupIds",
			mcp.Description("目标组ID，与Filters互斥。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("显示条数限制，默认为20。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("显示的偏移起始量。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件数组，与TargetGroupIds互斥，支持TargetGroupVpcId和TargetGroupName。"),
		),
	)
	mcpsvr.AddTool(clbDescribeTargetGroups, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := clb.NewClient(credential, region_, cpf)
		req := clb.NewDescribeTargetGroupsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTargetGroups(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	clbDescribeResources := mcp.NewTool(
		"clb-DescribeResources",
		mcp.WithDescription(`查询用户在当前地域支持可用区列表和资源列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回可用区资源列表数目，默认20，最大值100。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("返回可用区资源列表起始偏移量，默认0。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("查询可用区资源列表条件，详细的过滤条件如下："),
		),
	)
	mcpsvr.AddTool(clbDescribeResources, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := clb.NewClient(credential, region_, cpf)
		req := clb.NewDescribeResourcesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeResources(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	clbDescribeTargetGroupList := mcp.NewTool(
		"clb-DescribeTargetGroupList",
		mcp.WithDescription(`获取目标组列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"TargetGroupIds",
			mcp.Description("目标组ID数组。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件数组，支持TargetGroupVpcId和TargetGroupName。与TargetGroupIds互斥，优先使用目标组ID。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("显示的偏移起始量。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("显示条数限制，默认为20。"),
		),
	)
	mcpsvr.AddTool(clbDescribeTargetGroupList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := clb.NewClient(credential, region_, cpf)
		req := clb.NewDescribeTargetGroupListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTargetGroupList(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	clbDescribeTargetGroupInstances := mcp.NewTool(
		"clb-DescribeTargetGroupInstances",
		mcp.WithDescription(`获取目标组绑定的服务器`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，当前仅支持TargetGroupId，BindIP，InstanceId过滤。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("显示数量限制，默认20。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("显示的偏移量，默认为0。"),
		),
	)
	mcpsvr.AddTool(clbDescribeTargetGroupInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := clb.NewClient(credential, region_, cpf)
		req := clb.NewDescribeTargetGroupInstancesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTargetGroupInstances(req)
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
