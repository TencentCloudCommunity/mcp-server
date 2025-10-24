package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-vpn"
	mcpsvr := server.NewMCPServer(
		"腾讯云 VPN MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	vpcDescribeVpnGateways := mcp.NewTool(
		"vpc-DescribeVpnGateways",
		mcp.WithDescription(`查询VPN网关`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"VpnGatewayIds",
			mcp.Description("VPN网关实例ID。形如：vpngw-f49l6u0z。每次请求的实例的上限为100。参数不支持同时指定VpnGatewayIds和Filters。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，参数不支持同时指定VpnGatewayIds和Filters。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认值为0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("请求对象个数，默认值为20。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeVpnGateways, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeVpnGatewaysRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeVpnGateways(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeVpnConnections := mcp.NewTool(
		"vpc-DescribeVpnConnections",
		mcp.WithDescription(`查询VPN通道列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"VpnConnectionIds",
			mcp.Description("VPN通道实例ID。形如：vpnx-f49l6u0z。每次请求的实例的上限为100。参数不支持同时指定VpnConnectionIds和Filters。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件。每次请求的Filters的上限为10，Filter.Values的上限为5。参数不支持同时指定VpnConnectionIds和Filters。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认为0。关于Offset的更进一步介绍请参考 API 简介中的相关小节。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeVpnConnections, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeVpnConnectionsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeVpnConnections(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeVpnGatewayRoutes := mcp.NewTool(
		"vpc-DescribeVpnGatewayRoutes",
		mcp.WithDescription(`查询VPN网关路由`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"VpnGatewayId",
			mcp.Description("VPN网关实例ID。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件,  条件包括(DestinationCidr, InstanceId,InstanceType)。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量, 默认0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("单页个数, 默认20, 最大值100。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeVpnGatewayRoutes, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeVpnGatewayRoutesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeVpnGatewayRoutes(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeCustomerGateways := mcp.NewTool(
		"vpc-DescribeCustomerGateways",
		mcp.WithDescription(`查询对端网关`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"CustomerGatewayIds",
			mcp.Description("对端网关ID，例如：cgw-2wqq41m9。每次请求的实例的上限为100。参数不支持同时指定CustomerGatewayIds和Filters。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，详见下表：实例过滤条件表。每次请求的Filters的上限为10，Filter.Values的上限为5。参数不支持同时指定CustomerGatewayIds和Filters。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认为0。关于Offset的更进一步介绍请参考 API [简介](https://cloud.tencent.com/document/api/213/11646)中的相关小节。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeCustomerGateways, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeCustomerGatewaysRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeCustomerGateways(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeVpnGatewaySslServers := mcp.NewTool(
		"vpc-DescribeVpnGatewaySslServers",
		mcp.WithDescription(`查询SSL-VPN SERVER 列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("请求对象个数。"),
		),
		mcp.WithArray(
			"SslVpnServerIds",
			mcp.Description("SSL-VPN-SERVER实例ID。形如：vpns-0p4rj60。每次请求的实例的上限为100。参数不支持同时指定SslVpnServerIds和Filters。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，参数不支持同时指定SslVpnServerIds和Filters。"),
		),
		mcp.WithBoolean(
			"IsVpnPortal",
			mcp.Description("vpn门户使用。 默认Flase"),
		),
	)
	mcpsvr.AddTool(vpcDescribeVpnGatewaySslServers, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeVpnGatewaySslServersRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeVpnGatewaySslServers(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeVpnGatewaySslClients := mcp.NewTool(
		"vpc-DescribeVpnGatewaySslClients",
		mcp.WithDescription(`查询SSL-VPN-CLIENT 列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，参数不支持同时指定SslVpnClientIds和Filters。<li>vpc-id - String - （过滤条件）VPC实例ID形如：vpc-f49l6u0z。</li><li>vpn-gateway-id - String - （过滤条件）VPN实例ID形如：vpngw-5aluhh9t。</li><li>ssl-vpn-server-id - String - （过滤条件）SSL-VPN-SERVER实例ID形如：vpns-1j2w6xpx。</li><li>ssl-vpn-client-id - String - （过滤条件）SSL-VPN-CLIENT实例ID形如：vpnc-3rlxp4nd。</li><li>ssl-vpn-client-name - String - （过滤条件）SSL-VPN-CLIENT实例名称。</li><li>ssl-vpn-client-inner-ip - String - （过滤条件）SSL-VPN-CLIENT私网IP。</li>"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认值0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("请求对象个数，默认值20。"),
		),
		mcp.WithArray(
			"SslVpnClientIds",
			mcp.Description("SSL-VPN-CLIENT实例ID。形如：	"),
		),
		mcp.WithBoolean(
			"IsVpnPortal",
			mcp.Description("VPN门户网站使用。默认是False。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeVpnGatewaySslClients, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeVpnGatewaySslClientsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeVpnGatewaySslClients(req)
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
