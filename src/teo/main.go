package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-teo"
	mcpsvr := server.NewMCPServer(
		"腾讯云 TEO MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	teoDescribeIdentifications := mcp.NewTool(
		"teo-DescribeIdentifications",
		mcp.WithDescription(`查询站点的验证信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，Filters.Values的上限为20。详细的过滤条件如下："),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页查询偏移量。默认值：0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页查询限制数目。默认值：20，最大值：1000。"),
		),
	)
	mcpsvr.AddTool(teoDescribeIdentifications, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := teo.NewClient(credential, region_, cpf)
		req := teo.NewDescribeIdentificationsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeIdentifications(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	teoDescribeZones := mcp.NewTool(
		"teo-DescribeZones",
		mcp.WithDescription(`查询站点列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页查询偏移量。默认值：0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页查询限制数目。默认值：20，最大值：100。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，Filters.Values 的上限为 20。该参数不填写时，返回当前 appid 下有权限的所有站点信息。详细的过滤条件如下："),
		),
		mcp.WithString(
			"Order",
			mcp.Description("可根据该字段对返回结果进行排序，取值有："),
		),
		mcp.WithString(
			"Direction",
			mcp.Description("排序方向，如果是字段值为数字，则根据数字大小排序；如果字段值为文本，则根据 ascill 码的大小排序。取值有："),
		),
	)
	mcpsvr.AddTool(teoDescribeZones, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := teo.NewClient(credential, region_, cpf)
		req := teo.NewDescribeZonesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeZones(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	teoDescribeAccelerationDomains := mcp.NewTool(
		"teo-DescribeAccelerationDomains",
		mcp.WithDescription(`查询加速域名列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ZoneId",
			mcp.Description("加速域名所属站点 ID。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页查询偏移量，默认为 0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页查询限制数目，默认值：20，上限：200。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，Filters.Values 的上限为 20。该参数不填写时，返回当前 zone-id 下所有域名信息。详细的过滤条件如下："),
		),
		mcp.WithString(
			"Order",
			mcp.Description("可根据该字段对返回结果进行排序，取值有："),
		),
		mcp.WithString(
			"Direction",
			mcp.Description("排序方向，如果是字段值为数字，则根据数字大小排序；如果字段值为文本，则根据 ascill 码的大小排序。取值有："),
		),
		mcp.WithString(
			"Match",
			mcp.Description("匹配方式，取值有："),
		),
	)
	mcpsvr.AddTool(teoDescribeAccelerationDomains, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := teo.NewClient(credential, region_, cpf)
		req := teo.NewDescribeAccelerationDomainsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeAccelerationDomains(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	teoDescribeFunctions := mcp.NewTool(
		"teo-DescribeFunctions",
		mcp.WithDescription(`查询边缘函数列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ZoneId",
			mcp.Description("站点 ID。"),
		),
		mcp.WithArray(
			"FunctionIds",
			mcp.Description("按照函数 ID 列表过滤。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件列表，多个条件为且关系，Filters.Values 的上限为 20。详细的过滤条件如下："),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页查询偏移量。默认值：0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页查询限制数目。默认值：20，最大值：200。"),
		),
	)
	mcpsvr.AddTool(teoDescribeFunctions, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := teo.NewClient(credential, region_, cpf)
		req := teo.NewDescribeFunctionsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeFunctions(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	teoDescribeAliasDomains := mcp.NewTool(
		"teo-DescribeAliasDomains",
		mcp.WithDescription(`查询别称域名信息列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ZoneId",
			mcp.Description("站点 ID。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页查询偏移量。默认值：0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页查询限制数目。默认值：20，最大值：1000。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，Filters.Values的上限为20。详细的过滤条件如下："),
		),
	)
	mcpsvr.AddTool(teoDescribeAliasDomains, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := teo.NewClient(credential, region_, cpf)
		req := teo.NewDescribeAliasDomainsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeAliasDomains(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	teoDescribePrefetchTasks := mcp.NewTool(
		"teo-DescribePrefetchTasks",
		mcp.WithDescription(`查询预热任务状态`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ZoneId",
			mcp.Description("站点ID。该参数必填。"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("查询起始时间，时间与 job-id 必填一个。"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("查询结束时间，时间与 job-id 必填一个。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页查询偏移量，默认为 0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页查询限制数目，默认值：20，上限：1000。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，Filters.Values 的上限为 20。详细的过滤条件如下：<li>job-id：按照任务 ID 进行过滤。job-id 形如：1379afjk91u32h，暂不支持多值，不支持模糊查询；</li><li>target：按照目标资源信息进行过滤。target 形如：http://www.qq.com/1.txt，暂不支持多值，不支持模糊查询；</li><li>domains：按照域名行过滤。domains 形如：www.qq.com，不支持模糊查询；</li><li>statuses：按照任务状态进行过滤，不支持模糊查询。可选项：<br>   processing：处理中<br>   success：成功<br>   failed：失败<br>   timeout：超时<br>   invalid：无效。即源站响应非 2xx 状态码，请检查源站服务。</li>"),
		),
	)
	mcpsvr.AddTool(teoDescribePrefetchTasks, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := teo.NewClient(credential, region_, cpf)
		req := teo.NewDescribePrefetchTasksRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribePrefetchTasks(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	teoDescribeDDoSAttackEvent := mcp.NewTool(
		"teo-DescribeDDoSAttackEvent",
		mcp.WithDescription(`查询DDoS攻击事件列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("开始时间，时间范围为 30 天。"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间，时间范围为 30 天。"),
		),
		mcp.WithArray(
			"PolicyIds",
			mcp.Description("ddos策略组集合，不填默认选择全部策略。"),
		),
		mcp.WithArray(
			"ZoneIds",
			mcp.Description("站点集合，此参数必填。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页查询的限制数目，默认值为20，最大查询条目为1000。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页的偏移量，默认值为0。"),
		),
		mcp.WithBoolean(
			"ShowDetail",
			mcp.Description("展示攻击详情的参数，若填false，默认只返回攻击次数，不返回攻击详情；若填true，返回攻击详情。"),
		),
		mcp.WithString(
			"Area",
			mcp.Description("数据归属地区，取值有："),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序字段，取值有："),
		),
		mcp.WithString(
			"OrderType",
			mcp.Description("排序方式，取值有："),
		),
	)
	mcpsvr.AddTool(teoDescribeDDoSAttackEvent, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := teo.NewClient(credential, region_, cpf)
		req := teo.NewDescribeDDoSAttackEventRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDDoSAttackEvent(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	teoDescribeTimingL7AnalysisData := mcp.NewTool(
		"teo-DescribeTimingL7AnalysisData",
		mcp.WithDescription(`查询流量分析时序数据`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("开始时间。"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间。"),
		),
		mcp.WithArray(
			"MetricNames",
			mcp.Description("指标列表，取值有:"),
		),
		mcp.WithArray(
			"ZoneIds",
			mcp.Description("站点 ID 集合，此参数必填。"),
		),
		mcp.WithString(
			"Interval",
			mcp.Description("查询时间粒度，取值有："),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，详细的过滤条件 Key 值如下："),
		),
		mcp.WithString(
			"Area",
			mcp.Description("数据归属地区。该参数已废弃。请在 Filters.country 中按客户端地域过滤数据。"),
		),
	)
	mcpsvr.AddTool(teoDescribeTimingL7AnalysisData, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := teo.NewClient(credential, region_, cpf)
		req := teo.NewDescribeTimingL7AnalysisDataRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTimingL7AnalysisData(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	teoDescribeDefaultCertificates := mcp.NewTool(
		"teo-DescribeDefaultCertificates",
		mcp.WithDescription(`查询默认证书列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ZoneId",
			mcp.Description("站点 ID。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，Filters.Values的上限为5。详细的过滤条件如下："),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页查询偏移量。默认值：0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页查询限制数目。默认值：20，最大值：100。"),
		),
	)
	mcpsvr.AddTool(teoDescribeDefaultCertificates, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := teo.NewClient(credential, region_, cpf)
		req := teo.NewDescribeDefaultCertificatesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDefaultCertificates(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	teoDescribeDnsRecords := mcp.NewTool(
		"teo-DescribeDnsRecords",
		mcp.WithDescription(`查询 DNS 记录列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ZoneId",
			mcp.Description("站点 ID。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页查询偏移量，默认为 0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页查询限制数目，默认值：20，上限：1000。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤条件，Filters.Values 的上限为20。详细的过滤条件如下：<li>id： 按照 DNS 记录 ID 进行过滤，支持模糊查询；</li><li>name：按照 DNS 记录名称进行过滤，支持模糊查询；</li><li>content：按照 DNS 记录内容进行过滤，支持模糊查询；</li><li>type：按照 DNS 记录类型进行过滤，不支持模糊查询。可选项：<br>   A：将域名指向一个外网 IPv4 地址，如 8.8.8.8；<br>   AAAA：将域名指向一个外网 IPv6 地址；<br>   CNAME：将域名指向另一个域名，再由该域名解析出最终 IP 地址；<br>   TXT：对域名进行标识和说明，常用于域名验证和 SPF 记录（反垃圾邮件）；<br>   NS：如果需要将子域名交给其他 DNS 服务商解析，则需要添加 NS 记录。根域名无法添加 NS 记录；<br>   CAA：指定可为本站点颁发证书的 CA；<br>   SRV：标识某台服务器使用了某个服务，常见于微软系统的目录管理；<br>   MX：指定收件人邮件服务器。</li><li>ttl：按照解析生效时间进行过滤，不支持模糊查询。</li>"),
		),
		mcp.WithString(
			"SortBy",
			mcp.Description("排序依据，取值有：<li>content：DNS 记录内容；</li><li>created-on：DNS 记录创建时间；</li><li>name：DNS 记录名称；</li><li>ttl：缓存时间；</li><li>type：DNS 记录类型。</li>默认根据 type, name 属性组合排序。"),
		),
		mcp.WithString(
			"SortOrder",
			mcp.Description("列表排序方式，取值有：<li>asc：升序排列；</li><li>desc：降序排列。</li>默认值为 asc。"),
		),
		mcp.WithString(
			"Match",
			mcp.Description("匹配方式，取值有：<li>all：返回匹配所有查询条件的记录；</li><li>any：返回匹配任意一个查询条件的记录。</li>默认值为 all。"),
		),
	)
	mcpsvr.AddTool(teoDescribeDnsRecords, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := teo.NewClient(credential, region_, cpf)
		req := teo.NewDescribeDnsRecordsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDnsRecords(req)
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
