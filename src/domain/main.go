package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	domain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/domain/v20180808"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-domain"
	mcpsvr := server.NewMCPServer(
		"腾讯云 Domain MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	domainDescribeDomainPriceList := mcp.NewTool(
		"domain-DescribeDomainPriceList",
		mcp.WithDescription(`域名价格列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"TldList",
			mcp.Description("查询指定域名后缀的域名价格列表，若不指定域名后缀，默认查询所有后缀的域名价格列表。"),
		),
		mcp.WithString(
			"Year",
			mcp.Description("查询购买的年份，默认会列出所有年份的价格，可选值【1,10】"),
		),
		mcp.WithString(
			"Operation",
			mcp.Description("域名的购买类型："),
		),
	)
	mcpsvr.AddTool(domainDescribeDomainPriceList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := domain.NewClient(credential, region_, cpf)
		req := domain.NewDescribeDomainPriceListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDomainPriceList(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	domainCheckDomain := mcp.NewTool(
		"domain-CheckDomain",
		mcp.WithDescription(`域名注册查询`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DomainName",
			mcp.Description("所查询域名名称"),
		),
		mcp.WithString(
			"Period",
			mcp.Description("年限。该参数为空时无法查询溢价词域名"),
		),
	)
	mcpsvr.AddTool(domainCheckDomain, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := domain.NewClient(credential, region_, cpf)
		req := domain.NewCheckDomainRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CheckDomain(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	domainDescribeDomainBaseInfo := mcp.NewTool(
		"domain-DescribeDomainBaseInfo",
		mcp.WithDescription(`域名基本信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Domain",
			mcp.Description("域名"),
		),
	)
	mcpsvr.AddTool(domainDescribeDomainBaseInfo, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := domain.NewClient(credential, region_, cpf)
		req := domain.NewDescribeDomainBaseInfoRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDomainBaseInfo(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	domainDescribeTemplate := mcp.NewTool(
		"domain-DescribeTemplate",
		mcp.WithDescription(`获取模板信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"TemplateId",
			mcp.Description("模板ID(模板列表接口可获取)"),
		),
	)
	mcpsvr.AddTool(domainDescribeTemplate, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := domain.NewClient(credential, region_, cpf)
		req := domain.NewDescribeTemplateRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTemplate(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	domainDescribeDomainNameList := mcp.NewTool(
		"domain-DescribeDomainNameList",
		mcp.WithDescription(`我的域名列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认为0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为20，取值范围[1,100]"),
		),
	)
	mcpsvr.AddTool(domainDescribeDomainNameList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := domain.NewClient(credential, region_, cpf)
		req := domain.NewDescribeDomainNameListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDomainNameList(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	domainDescribeTemplateList := mcp.NewTool(
		"domain-DescribeTemplateList",
		mcp.WithDescription(`信息模板列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认为0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100。"),
		),
		mcp.WithString(
			"Type",
			mcp.Description("用户注册类型，默认:all , 个人：I ,企业: E"),
		),
		mcp.WithString(
			"Status",
			mcp.Description("认证状态：未实名审核:NotUpload, 实名审核中:InAudit，已实名审核:Approved，实名审核失败:Reject，更新手机邮箱:NotVerified。"),
		),
		mcp.WithString(
			"Keyword",
			mcp.Description("关键字，用于域名所有者筛选"),
		),
	)
	mcpsvr.AddTool(domainDescribeTemplateList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := domain.NewClient(credential, region_, cpf)
		req := domain.NewDescribeTemplateListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTemplateList(req)
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
