package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-vod"
	mcpsvr := server.NewMCPServer(
		"腾讯云 Vod MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	vodDescribeSubAppIds := mcp.NewTool(
		"vod-DescribeSubAppIds",
		mcp.WithDescription(`获取子应用列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Name",
			mcp.Description("应用名称"),
		),
		mcp.WithString(
			"Tags",
			mcp.Description("标签信息，查询指定标签的子应用列表"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页拉取的起始偏移量。默认值：0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页拉取的最大返回结果数。默认值：200；最大值：200"),
		),
	)
	mcpsvr.AddTool(vodDescribeSubAppIds, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vod.NewClient(credential, region_, cpf)
		req := vod.NewDescribeSubAppIdsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSubAppIds(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vodCreateSubAppId := mcp.NewTool(
		"vod-CreateSubAppId",
		mcp.WithDescription(`创建子应用`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Name",
			mcp.Description("应用名称，长度限制：40个字符。"),
		),
		mcp.WithString(
			"Description",
			mcp.Description("应用简介，长度限制：300个字符。"),
		),
		mcp.WithString(
			"Type",
			mcp.Description("应用类型， 取值有：AllInOne：一体化；Professional：专业版。默认值为 AllInOne。"),
		),
	)
	mcpsvr.AddTool(vodCreateSubAppId, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vod.NewClient(credential, region_, cpf)
		req := vod.NewCreateSubAppIdRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateSubAppId(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vodApplyUpload := mcp.NewTool(
		"vod-ApplyUpload",
		mcp.WithDescription(`申请上传`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"MediaType",
			mcp.Description("媒体类型，可选值请参考 上传能力综述。"),
		),
		mcp.WithString(
			"MediaName",
			mcp.Description("媒体名称。"),
		),
		mcp.WithString(
			"CoverType",
			mcp.Description("封面类型，可选值请参考 上传能力综述。"),
		),
		mcp.WithString(
			"Procedure",
			mcp.Description("任务流上下文，任务流模板名。"),
		),
		mcp.WithString(
			"ExpireTime",
			mcp.Description("媒体文件过期时间，格式按照 ISO 8601标准表示，详见 ISO 日期格式说明。"),
		),
		mcp.WithString(
			"StorageRegion",
			mcp.Description("指定上传园区，仅适用于对上传地域有特殊需求的用户。"),
		),
		mcp.WithString(
			"ClassId",
			mcp.Description("分类ID，用于对媒体进行分类管理，可通过 创建分类 接口，创建分类，获得分类 ID。"),
		),
		mcp.WithString(
			"SourceContext",
			mcp.Description("来源上下文，用于透传用户请求信息，上传回调接口将返回该字段值，最长 250 个字符。"),
		),
		mcp.WithString(
			"SessionContext",
			mcp.Description("会话上下文，用于透传用户请求信息，当指定 Procedure 参数后，任务流状态变更回调将返回该字段值，最长 1000 个字符。"),
		),
		mcp.WithString(
			"SubAppId",
			mcp.Description("点播 子应用 ID。如果要访问子应用中的资源，则将该字段填写为子应用 ID；否则无需填写该字段。"),
		),
	)
	mcpsvr.AddTool(vodApplyUpload, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vod.NewClient(credential, region_, cpf)
		req := vod.NewApplyUploadRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ApplyUpload(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vodCommitUpload := mcp.NewTool(
		"vod-CommitUpload",
		mcp.WithDescription(`确认上传`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"VodSessionKey",
			mcp.Description("点播会话，取申请上传接口的返回值 VodSessionKey。"),
		),
		mcp.WithString(
			"SubAppId",
			mcp.Description("点播 子应用 ID。如果要访问子应用中的资源，则将该字段填写为子应用 ID；否则无需填写该字段。"),
		),
	)
	mcpsvr.AddTool(vodCommitUpload, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vod.NewClient(credential, region_, cpf)
		req := vod.NewCommitUploadRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CommitUpload(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vodDescribeTaskDetail := mcp.NewTool(
		"vod-DescribeTaskDetail",
		mcp.WithDescription(`查询任务详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"TaskId",
			mcp.Description("视频处理任务的任务 ID。"),
		),
		mcp.WithString(
			"SubAppId",
			mcp.Description("点播 子应用 ID。如果要访问子应用中的资源，则将该字段填写为子应用 ID；否则无需填写该字段。"),
		),
	)
	mcpsvr.AddTool(vodDescribeTaskDetail, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vod.NewClient(credential, region_, cpf)
		req := vod.NewDescribeTaskDetailRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTaskDetail(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vodDescribeMediaInfos := mcp.NewTool(
		"vod-DescribeMediaInfos",
		mcp.WithDescription(`获取媒体详细信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FileIds",
			mcp.Description("媒体文件 ID 列表，N 从 0 开始取值，最大 19。"),
		),
		mcp.WithString(
			"Filters",
			mcp.Description("指定所有媒体文件需要返回的信息，可同时指定多个信息，N 从 0 开始递增。如果未填写该字段，默认返回所有信息。选项有："),
		),
		mcp.WithString(
			"SubAppId",
			mcp.Description("点播 子应用 ID。如果要访问子应用中的资源，则将该字段填写为子应用 ID；否则无需填写该字段。"),
		),
	)
	mcpsvr.AddTool(vodDescribeMediaInfos, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vod.NewClient(credential, region_, cpf)
		req := vod.NewDescribeMediaInfosRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeMediaInfos(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	vodDescribeVodDomains := mcp.NewTool(
		"vod-DescribeVodDomains",
		mcp.WithDescription(`查询点播域名信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Domains",
			mcp.Description("域名列表。当该字段不填时，则默认列出所有域名信息。本字段限制如下："),
		),
		mcp.WithString(
			"Limit",
			mcp.Description("分页拉取的最大返回结果数。默认值：20；最大值：20。"),
		),
		mcp.WithString(
			"Offset",
			mcp.Description("分页拉取的起始偏移量。默认值：0。"),
		),
		mcp.WithString(
			"SubAppId",
			mcp.Description("点播 子应用 ID。如果要访问子应用中的资源，则将该字段填写为子应用 ID；否则无需填写该字段。"),
		),
	)
	mcpsvr.AddTool(vodDescribeVodDomains, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := vod.NewClient(credential, region_, cpf)
		req := vod.NewDescribeVodDomainsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeVodDomains(req)
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
