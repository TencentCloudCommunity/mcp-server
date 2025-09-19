package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func main() {
	mcpServerName := "mcp-server-ocr-wordrecognition-high"
	mcpsvr := server.NewMCPServer(
		"腾讯云OCR通用文字识别（高精度版）MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	ocrGeneralAccurateOCR := mcp.NewTool(
		"ocr-GeneralAccurateOCR",
		mcp.WithDescription(`通用印刷体识别（高精度版）`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ImageBase64",
			mcp.Description("图片的 Base64 值。"),
		),
		mcp.WithString(
			"ImageUrl",
			mcp.Description("图片的 Url 地址。"),
		),
		mcp.WithBoolean(
			"IsWords",
			mcp.Description("是否返回单字信息，默认关"),
		),
		mcp.WithBoolean(
			"EnableDetectSplit",
			mcp.Description("是否开启原图切图检测功能，开启后可提升“整图面积大，但单字符占比面积小”（例如：试卷）场景下的识别效果，默认关"),
		),
		mcp.WithBoolean(
			"IsPdf",
			mcp.Description("是否开启PDF识别，默认值为false，开启后可同时支持图片和PDF的识别。"),
		),
		mcp.WithNumber(
			"PdfPageNumber",
			mcp.Description("需要识别的PDF页面的对应页码，仅支持PDF单页识别，当上传文件为PDF且IsPdf参数值为true时有效，默认值为1。"),
		),
		mcp.WithBoolean(
			"EnableDetectText",
			mcp.Description("文本检测开关，默认为true。设置为false可直接进行单行识别，适用于仅包含正向单行文本的图片场景。"),
		),
		mcp.WithString(
			"ConfigID",
			mcp.Description("配置ID支持： OCR -- 通用场景 MulOCR--多语种场景，注：仅ConfigID配置为OCR时支持\n示例值：OCR"),
		),
	)
	mcpsvr.AddTool(ocrGeneralAccurateOCR, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ocr.NewClient(credential, region_, cpf)
		req := ocr.NewGeneralAccurateOCRRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.GeneralAccurateOCR(req)
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
