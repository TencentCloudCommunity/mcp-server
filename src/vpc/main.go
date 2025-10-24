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
	mcpServerName := "mcp-server-vpc"
	mcpsvr := server.NewMCPServer(
		"腾讯云 VPC MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	vpcDescribeVpcs := mcp.NewTool(
		"vpc-DescribeVpcs",
		mcp.WithDescription(`查询VPC列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"VpcIds",
			mcp.Description("VPC实例ID。形如：vpc-f49l6u0z。每次请求的实例的上限为100。参数不支持同时指定VpcIds和Filters。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，不支持同时指定VpcIds和Filters参数。"),
		),
		mcp.WithString(
			"Offset",
			mcp.Description("偏移量，默认为0。"),
		),
		mcp.WithString(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeVpcs, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeVpcsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeVpcs(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeSubnets := mcp.NewTool(
		"vpc-DescribeSubnets",
		mcp.WithDescription(`查询子网列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"SubnetIds",
			mcp.Description("子网实例ID查询。形如：subnet-pxir56ns。每次请求的实例的上限为100。参数不支持同时指定SubnetIds和Filters。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，参数不支持同时指定SubnetIds和Filters。"),
		),
		mcp.WithString(
			"Offset",
			mcp.Description("偏移量，默认为0。"),
		),
		mcp.WithString(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeSubnets, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeSubnetsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSubnets(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeRouteTables := mcp.NewTool(
		"vpc-DescribeRouteTables",
		mcp.WithDescription(`查询路由表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，参数不支持同时指定RouteTableIds和Filters。"),
		),
		mcp.WithArray(
			"RouteTableIds",
			mcp.Description("路由表实例ID，例如：rtb-azd4dt1c。"),
		),
		mcp.WithString(
			"Offset",
			mcp.Description("偏移量。"),
		),
		mcp.WithString(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。"),
		),
		mcp.WithBoolean(
			"NeedRouterInfo",
			mcp.Description("是否需要获取路由策略信息，默认获取，当控制台不需要拉取路由策略信息时，改为False。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeRouteTables, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeRouteTablesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeRouteTables(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeNetworkInterfaces := mcp.NewTool(
		"vpc-DescribeNetworkInterfaces",
		mcp.WithDescription(`查询弹性网卡列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"NetworkInterfaceIds",
			mcp.Description("弹性网卡实例ID查询。形如：eni-pxir56ns。每次请求的实例的上限为100。参数不支持同时指定NetworkInterfaceIds和Filters。可通过[DescribeNetworkInterfaces](https://cloud.tencent.com/document/product/215/15817)接口获取。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，参数不支持同时指定NetworkInterfaceIds和Filters。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认为0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeNetworkInterfaces, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeNetworkInterfacesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeNetworkInterfaces(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeSecurityGroups := mcp.NewTool(
		"vpc-DescribeSecurityGroups",
		mcp.WithDescription(`查看安全组`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"SecurityGroupIds",
			mcp.Description("安全组实例ID，例如：sg-33ocnj9n。每次请求的实例的上限为100。参数不支持同时指定SecurityGroupIds和Filters。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，参数不支持同时指定SecurityGroupIds和Filters。"),
		),
		mcp.WithString(
			"Offset",
			mcp.Description("偏移量，默认为0。"),
		),
		mcp.WithString(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。"),
		),
		mcp.WithString(
			"OrderField",
			mcp.Description("排序字段。支持：`CreatedTime` `UpdateTime`。注意：该字段没有默认值。"),
		),
		mcp.WithString(
			"OrderDirection",
			mcp.Description("排序方法。升序：`ASC`，倒序：`DESC`。默认值：`ASC`"),
		),
	)
	mcpsvr.AddTool(vpcDescribeSecurityGroups, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeSecurityGroupsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSecurityGroups(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeHaVips := mcp.NewTool(
		"vpc-DescribeHaVips",
		mcp.WithDescription(`查询HAVIP列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"HaVipIds",
			mcp.Description("`HAVIP`唯一`ID`，形如：`havip-9o233uri`。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，参数不支持同时指定`HaVipIds`和`Filters`。<li>havip-id - String - `HAVIP`唯一`ID`，形如：`havip-9o233uri`。</li><li>havip-name - String - `HAVIP`名称。</li><li>vpc-id - String - `HAVIP`所在私有网络`ID`。</li><li>subnet-id - String - `HAVIP`所在子网`ID`。</li><li>vip - String - `HAVIP`的地址`VIP`。</li><li>address-ip - String - `HAVIP`绑定的弹性公网`IP`。</li><li>havip-association.instance-id - String - `HAVIP`绑定的子机或网卡。</li><li>havip-association.instance-type - String - `HAVIP`绑定的类型，取值:CVM, ENI。</li><li>check-associate - Bool - 是否开启HaVip飘移时校验绑定的子机或网卡。</li><li>cdc-id - String - CDC实例ID。</li>"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认为0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeHaVips, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeHaVipsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeHaVips(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeAddressTemplates := mcp.NewTool(
		"vpc-DescribeAddressTemplates",
		mcp.WithDescription(`查询IP地址模板`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Filters",
			mcp.Description("过滤条件。"),
		),
		mcp.WithString(
			"Offset",
			mcp.Description("偏移量，默认为0。"),
		),
		mcp.WithString(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。"),
		),
		mcp.WithString(
			"NeedMemberInfo",
			mcp.Description("是否获取IP地址模板成员标识。"),
		),
		mcp.WithString(
			"OrderField",
			mcp.Description("排序字段。支持：`AddressTemplateId` `CreatedTime` `UpdateTime`。注意：该字段没有默认值。"),
		),
		mcp.WithString(
			"OrderDirection",
			mcp.Description("排序方法。升序：`ASC`，倒序：`DESC`。注意：该字段没有默认值。"),
		),
		mcp.WithString(
			"MemberOrderField",
			mcp.Description("IP成员排序字段。支持：`Address` `UpdateTime`。注意：该字段没有默认值。"),
		),
		mcp.WithString(
			"MemberOrderDirection",
			mcp.Description("IP成员排序方法。升序：`ASC`，倒序：`DESC`。注意：该字段没有默认值。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeAddressTemplates, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeAddressTemplatesRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeAddressTemplates(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vpcDescribeNetworkAcls := mcp.NewTool(
		"vpc-DescribeNetworkAcls",
		mcp.WithDescription(`查询网络ACL列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，参数不支持同时指定NetworkAclIds和Filters。"),
		),
		mcp.WithArray(
			"NetworkAclIds",
			mcp.Description("网络ACL实例ID数组。形如：[acl-12345678]。每次请求的实例的上限为100。参数不支持同时指定NetworkAclIds和Filters。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认为0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为20，最小值为1，最大值为100。"),
		),
		mcp.WithString(
			"OrderField",
			mcp.Description("排序字段。支持：NetworkAclId,NetworkAclName,CreatedTime"),
		),
		mcp.WithString(
			"OrderDirection",
			mcp.Description("排序方法。顺序：ASC，倒序：DESC。"),
		),
	)
	mcpsvr.AddTool(vpcDescribeNetworkAcls, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vpc.NewClient(credential, region_, cpf)
		req := vpc.NewDescribeNetworkAclsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeNetworkAcls(req)
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
