package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-live"
	mcpsvr := server.NewMCPServer(
		"腾讯云 Live MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	liveDescribeLiveDomains := mcp.NewTool(
		"live-DescribeLiveDomains",
		mcp.WithDescription(`查询域名列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DomainStatus",
			mcp.Description("域名状态过滤。0-停用，1-启用。"),
		),
		mcp.WithString(
			"DomainType",
			mcp.Description("域名类型过滤。0-推流，1-播放。"),
		),
		mcp.WithString(
			"PageSize",
			mcp.Description("分页大小，范围：10~100。默认10。"),
		),
		mcp.WithString(
			"PageNum",
			mcp.Description("取第几页，范围：1~100000。默认1。"),
		),
		mcp.WithString(
			"IsDelayLive",
			mcp.Description("0 普通直播 1慢直播 默认0。"),
		),
		mcp.WithString(
			"DomainPrefix",
			mcp.Description("域名前缀。"),
		),
		mcp.WithString(
			"PlayType",
			mcp.Description("播放区域，只在 DomainType=1 时该参数有意义。"),
		),
	)
	mcpsvr.AddTool(liveDescribeLiveDomains, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := live.NewClient(credential, region_, cpf)
		req := live.NewDescribeLiveDomainsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeLiveDomains(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	liveDescribeLiveDomain := mcp.NewTool(
		"live-DescribeLiveDomain",
		mcp.WithDescription(`查询域名信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DomainName",
			mcp.Description("域名。"),
		),
	)
	mcpsvr.AddTool(liveDescribeLiveDomain, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := live.NewClient(credential, region_, cpf)
		req := live.NewDescribeLiveDomainRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeLiveDomain(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	liveDescribeLiveStreamOnlineList := mcp.NewTool(
		"live-DescribeLiveStreamOnlineList",
		mcp.WithDescription(`查询直播中的流`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DomainName",
			mcp.Description("推流域名。多域名用户需要填写 DomainName。"),
		),
		mcp.WithString(
			"AppName",
			mcp.Description("推流路径，与推流和播放地址中的 AppName 保持一致，默认为 live。多路径用户需要填写 AppName。"),
		),
		mcp.WithString(
			"PageNum",
			mcp.Description("取得第几页，默认1。"),
		),
		mcp.WithString(
			"PageSize",
			mcp.Description("每页大小，最大100。 "),
		),
		mcp.WithString(
			"StreamName",
			mcp.Description("流名称，用于精确查询。"),
		),
	)
	mcpsvr.AddTool(liveDescribeLiveStreamOnlineList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := live.NewClient(credential, region_, cpf)
		req := live.NewDescribeLiveStreamOnlineListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeLiveStreamOnlineList(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	liveDescribePushBandwidthAndFluxList := mcp.NewTool(
		"live-DescribePushBandwidthAndFluxList",
		mcp.WithDescription(`直播推流带宽和流量数据查询`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("查询开始时间点，精确到分钟粒度，接口查询支持两种时间格式："),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("查询结束时间点，精确到分钟粒度，接口查询支持两种时间格式："),
		),
		mcp.WithString(
			"PushDomains",
			mcp.Description("域名，可以填多个，若不填，表示总体数据。"),
		),
		mcp.WithString(
			"MainlandOrOversea",
			mcp.Description("可选值："),
		),
		mcp.WithString(
			"Granularity",
			mcp.Description("数据粒度，支持如下粒度："),
		),
		mcp.WithString(
			"RegionNames",
			mcp.Description("大区，映射表如下："),
		),
		mcp.WithString(
			"CountryNames",
			mcp.Description("国家，映射表参照如下文档："),
		),
	)
	mcpsvr.AddTool(liveDescribePushBandwidthAndFluxList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := live.NewClient(credential, region_, cpf)
		req := live.NewDescribePushBandwidthAndFluxListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribePushBandwidthAndFluxList(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	liveDescribeBillBandwidthAndFluxList := mcp.NewTool(
		"live-DescribeBillBandwidthAndFluxList",
		mcp.WithDescription(`直播播放带宽和流量数据查询`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("起始时间点，接口查询支持两种时间格式："),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间点，接口查询支持两种时间格式："),
		),
		mcp.WithString(
			"PlayDomains",
			mcp.Description("直播播放域名，若不填，表示总体数据。"),
		),
		mcp.WithString(
			"MainlandOrOversea",
			mcp.Description("可选值："),
		),
		mcp.WithString(
			"Granularity",
			mcp.Description("数据粒度，支持如下粒度："),
		),
		mcp.WithString(
			"ServiceName",
			mcp.Description("服务名称，可选值包括LVB(标准直播)，LEB(快直播)，不填则查LVB+LEB总值。"),
		),
		mcp.WithString(
			"RegionNames",
			mcp.Description("大区，映射表如下："),
		),
	)
	mcpsvr.AddTool(liveDescribeBillBandwidthAndFluxList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := live.NewClient(credential, region_, cpf)
		req := live.NewDescribeBillBandwidthAndFluxListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeBillBandwidthAndFluxList(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	liveDescribeStreamPlayInfoList := mcp.NewTool(
		"live-DescribeStreamPlayInfoList",
		mcp.WithDescription(`查询流的播放信息列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("起始时间点，接口查询支持两种时间格式："),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间点，接口查询支持两种时间格式："),
		),
		mcp.WithString(
			"PlayDomain",
			mcp.Description("播放域名，"),
		),
		mcp.WithString(
			"StreamName",
			mcp.Description("流名称，精确匹配。"),
		),
		mcp.WithString(
			"AppName",
			mcp.Description("推流路径，与播放地址中的AppName保持一致，会精确匹配，在同时传递了StreamName时生效。"),
		),
		mcp.WithString(
			"ServiceName",
			mcp.Description("服务名称，可选值包括LVB(标准直播)，LEB(快直播)，不填则查LVB+LEB总值。"),
		),
	)
	mcpsvr.AddTool(liveDescribeStreamPlayInfoList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := live.NewClient(credential, region_, cpf)
		req := live.NewDescribeStreamPlayInfoListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeStreamPlayInfoList(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	liveDescribeStreamPushInfoList := mcp.NewTool(
		"live-DescribeStreamPushInfoList",
		mcp.WithDescription(`查询某条流上行推流质量数据`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"StreamName",
			mcp.Description("流名称。"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("起始时间点，"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间点，"),
		),
		mcp.WithString(
			"PushDomain",
			mcp.Description("推流域名。"),
		),
		mcp.WithString(
			"AppName",
			mcp.Description("推流路径，与推流和播放地址中的AppName保持一致，默认为 live。"),
		),
	)
	mcpsvr.AddTool(liveDescribeStreamPushInfoList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := live.NewClient(credential, region_, cpf)
		req := live.NewDescribeStreamPushInfoListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeStreamPushInfoList(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	liveDescribeLiveStreamEventList := mcp.NewTool(
		"live-DescribeLiveStreamEventList",
		mcp.WithDescription(`查询推断流事件`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("起始时间。 UTC 格式，例如：2018-12-29T19:00:00Z。支持查询2个月内的历史记录。"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间。UTC 格式，例如：2018-12-29T20:00:00Z。不超过当前时间，且和起始时间相差不得超过1个月。"),
		),
		mcp.WithString(
			"AppName",
			mcp.Description("推流路径，与推流和播放地址中的AppName保持一致，默认为 live。"),
		),
		mcp.WithString(
			"DomainName",
			mcp.Description("推流域名。"),
		),
		mcp.WithString(
			"StreamName",
			mcp.Description("流名称，不支持通配符（*）查询，默认模糊匹配。"),
		),
		mcp.WithString(
			"PageNum",
			mcp.Description("取得第几页。"),
		),
		mcp.WithString(
			"PageSize",
			mcp.Description("分页大小。"),
		),
		mcp.WithString(
			"IsFilter",
			mcp.Description("是否过滤，默认不过滤。"),
		),
		mcp.WithString(
			"IsStrict",
			mcp.Description("是否精确查询，默认模糊匹配。"),
		),
		mcp.WithString(
			"IsAsc",
			mcp.Description("是否按结束时间正序显示，默认逆序。"),
		),
	)
	mcpsvr.AddTool(liveDescribeLiveStreamEventList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := live.NewClient(credential, region_, cpf)
		req := live.NewDescribeLiveStreamEventListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeLiveStreamEventList(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	liveDescribeLiveTranscodeDetailInfo := mcp.NewTool(
		"live-DescribeLiveTranscodeDetailInfo",
		mcp.WithDescription(`查询直播转码统计信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"PushDomain",
			mcp.Description("推流域名。"),
		),
		mcp.WithString(
			"StreamName",
			mcp.Description("流名称。"),
		),
		mcp.WithString(
			"DayTime",
			mcp.Description("查询时间，北京时间，"),
		),
		mcp.WithString(
			"PageNum",
			mcp.Description("页数，默认1，"),
		),
		mcp.WithString(
			"PageSize",
			mcp.Description("每页个数，默认20，"),
		),
		mcp.WithString(
			"StartDayTime",
			mcp.Description("起始天时间，北京时间，"),
		),
		mcp.WithString(
			"EndDayTime",
			mcp.Description("结束天时间，北京时间，"),
		),
	)
	mcpsvr.AddTool(liveDescribeLiveTranscodeDetailInfo, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := live.NewClient(credential, region_, cpf)
		req := live.NewDescribeLiveTranscodeDetailInfoRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeLiveTranscodeDetailInfo(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	liveDescribeVisitTopSumInfoList := mcp.NewTool(
		"live-DescribeVisitTopSumInfoList",
		mcp.WithDescription(`查询某时间段top n的域名或流id信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("起始时间点，格式为yyyy-mm-dd HH:MM:SS。"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间点，格式为yyyy-mm-dd HH:MM:SS"),
		),
		mcp.WithString(
			"TopIndex",
			mcp.Description("峰值指标，可选值包括”Domain”，”StreamId”。"),
		),
		mcp.WithString(
			"PlayDomains",
			mcp.Description("播放域名，默认为不填，表示求总体数据。"),
		),
		mcp.WithString(
			"PageNum",
			mcp.Description("页号，"),
		),
		mcp.WithString(
			"PageSize",
			mcp.Description("每页个数，范围是[1,1000]，"),
		),
		mcp.WithString(
			"OrderParam",
			mcp.Description("排序指标，可选值包括” AvgFluxPerSecond”，”TotalRequest”（默认）,“TotalFlux”。"),
		),
	)
	mcpsvr.AddTool(liveDescribeVisitTopSumInfoList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := live.NewClient(credential, region_, cpf)
		req := live.NewDescribeVisitTopSumInfoListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeVisitTopSumInfoList(req)
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
