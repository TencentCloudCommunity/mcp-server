package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-cdn"
	mcpsvr := server.NewMCPServer(
		"腾讯云 CDN MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	cdnDescribeOriginData := mcp.NewTool(
		"cdn-DescribeOriginData",
		mcp.WithDescription(`回源数据查询`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("查询起始时间，如：2018-09-04 10:40:00，返回结果大于等于指定时间"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("查询结束时间，如：2018-09-04 10:40:00，返回结果小于等于指定时间"),
		),
		mcp.WithString(
			"Metric",
			mcp.Description("指定查询指标，支持的类型有："),
		),
		mcp.WithArray(
			"Domains",
			mcp.Description("指定查询域名列表，最多可一次性查询 30 个加速域名明细"),
		),
		mcp.WithNumber(
			"Project",
			mcp.Description("指定要查询的项目 ID，[前往查看项目 ID](https://console.cloud.tencent.com/project)"),
		),
		mcp.WithString(
			"Interval",
			mcp.Description("时间粒度，支持以下几种模式："),
		),
		mcp.WithBoolean(
			"Detail",
			mcp.Description("Domains 传入多个时，默认（false)返回多个域名的汇总数据"),
		),
		mcp.WithString(
			"Area",
			mcp.Description("指定服务地域查询，不填充表示查询中国境内 CDN 数据"),
		),
		mcp.WithString(
			"TimeZone",
			mcp.Description("指定查询时间的时区，默认UTC+08:00"),
		),
	)
	mcpsvr.AddTool(cdnDescribeOriginData, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdn.NewClient(credential, region_, cpf)
		req := cdn.NewDescribeOriginDataRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeOriginData(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cdnListTopData := mcp.NewTool(
		"cdn-ListTopData",
		mcp.WithDescription(`TOP 数据查询`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("查询起始时间：yyyy-MM-dd HH:mm:ss"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("查询结束时间：yyyy-MM-dd HH:mm:ss"),
		),
		mcp.WithString(
			"Metric",
			mcp.Description("排序对象，支持以下几种形式："),
		),
		mcp.WithString(
			"Filter",
			mcp.Description("排序使用的指标名称："),
		),
		mcp.WithArray(
			"Domains",
			mcp.Description("指定查询域名列表，最多可一次性查询 30 个加速域名明细"),
		),
		mcp.WithNumber(
			"Project",
			mcp.Description("指定要查询的项目 ID，[前往查看项目 ID](https://console.cloud.tencent.com/project)"),
		),
		mcp.WithBoolean(
			"Detail",
			mcp.Description("多域名查询时，默认（false)返回所有域名汇总排序结果"),
		),
		mcp.WithString(
			"Code",
			mcp.Description("Filter 为 statusCode、OriginStatusCode 时，填充指定状态码查询排序结果"),
		),
		mcp.WithString(
			"Area",
			mcp.Description("指定服务地域查询，不填充表示查询中国境内 CDN 数据"),
		),
		mcp.WithString(
			"AreaType",
			mcp.Description("查询中国境外CDN数据，且仅当 Metric 为 district 或 host 时，可指定地区类型查询，不填充表示查询服务地区数据（仅在 Area 为 overseas，且 Metric 是 district 或 host 时可用）"),
		),
		mcp.WithString(
			"Product",
			mcp.Description("指定查询的产品数据，可选为cdn或者ecdn，默认为cdn"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("只返回前N条数据，默认为最大值100，metric=url时默认为最大值1000"),
		),
	)
	mcpsvr.AddTool(cdnListTopData, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdn.NewClient(credential, region_, cpf)
		req := cdn.NewListTopDataRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ListTopData(req)
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
