package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func main() {
	mcpServerName := "mcp-server-cbs"
	mcpsvr := server.NewMCPServer(
		"腾讯云 CBS MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	cbsDescribeDisks := mcp.NewTool(
		"cbs-DescribeDisks",
		mcp.WithDescription(`查询云硬盘列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件。参数不支持同时指定DiskIds和Filters。详情见: https://cloud.tencent.com/document/api/362/16315"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。关于`Limit`的更进一步介绍请参考 API [简介](/document/product/362/15633)中的相关小节。"),
		),
		mcp.WithString(
			"OrderField",
			mcp.Description("云盘列表排序的依据字段。取值范围：<br><li>CREATE_TIME：依据云盘的创建时间排序<br></li><li>DEADLINE：依据云盘的到期时间排序<br>默认按云盘创建时间排序。</li>"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认为0。关于`Offset`的更进一步介绍请参考API[简介](/document/product/362/15633)中的相关小节。"),
		),
		mcp.WithString(
			"ReturnBindAutoSnapshotPolicy",
			mcp.Description("云盘详情中是否需要返回云盘绑定的定期快照策略ID，TRUE表示需要返回，FALSE表示不返回。"),
		),
		mcp.WithArray(
			"DiskIds",
			mcp.Description("按照一个或者多个云硬盘ID查询。云硬盘ID形如：`disk-11112222`，此参数的具体格式可参考API[简介](/document/product/362/15633)的ids.N一节）。参数不支持同时指定`DiskIds`和`Filters`。"),
		),
		mcp.WithString(
			"Order",
			mcp.Description("输出云盘列表的排列顺序。取值范围：<br><li>ASC：升序排列<br></li><li>DESC：降序排列。</li>"),
		),
	)
	mcpsvr.AddTool(cbsDescribeDisks, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cbs.NewClient(credential, region_, cpf)
		req := cbs.NewDescribeDisksRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDisks(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cbsDescribeSnapshots := mcp.NewTool(
		"cbs-DescribeSnapshots",
		mcp.WithDescription(`查询快照列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"SnapshotIds",
			mcp.Description("要查询快照的ID列表。参数不支持同时指定`SnapshotIds`和`Filters`。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件。参数不支持同时指定`SnapshotIds`和`Filters`。详情见: https://cloud.tencent.com/document/api/362/15647"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。关于`Limit`的更进一步介绍请参考 API [简介](/document/product/362/15633)中的相关小节。"),
		),
		mcp.WithString(
			"OrderField",
			mcp.Description("快照列表排序的依据字段。取值范围：<br><li>CREATE_TIME：依据快照的创建时间排序<br>默认按创建时间排序。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认为0。关于`Offset`的更进一步介绍请参考API[简介](/document/product/362/15633)中的相关小节。"),
		),
		mcp.WithString(
			"Order",
			mcp.Description("输出云盘列表的排列顺序。取值范围：<br><li>ASC：升序排列<br><li>DESC：降序排列。"),
		),
	)
	mcpsvr.AddTool(cbsDescribeSnapshots, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cbs.NewClient(credential, region_, cpf)
		req := cbs.NewDescribeSnapshotsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSnapshots(req)
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
