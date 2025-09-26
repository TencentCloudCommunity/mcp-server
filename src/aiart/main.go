package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	aiart "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/aiart/v20221229"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func main() {
	mcpServerName := "mcp-server-aiart"
	mcpsvr := server.NewMCPServer(
		"腾讯云 AIART MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	aiartTextToImageLite := mcp.NewTool(
		"aiart-TextToImageLite",
		mcp.WithDescription(`混元生图（极速版）`),
		mcp.WithString("Prompt",
			mcp.Description("文本描述。将根据输入的文本智能生成与之相关的图像。"),
		),
		mcp.WithString("NegativePrompt",
			mcp.Description("反向提示词。 减少生成结果中出现描述内容。"),
		),
		mcp.WithString("Resolution",
			mcp.Description("生成图分辨率，默认1024:1024。支持的图像宽高比例: 1:1，3:4，4:3，9:16，16:9。支持的长边分辨率: 160，200，225，258，512，520，608，768，1024，1080，1280，1600，1620，1920，2048，2400，2560，2592，3440，3840，4096。"),
		),
		mcp.WithNumber("Seed",
			mcp.Description("随机种子，默认随机。0：随机种子生成。不传：随机种子生成。正数：固定种子生成。"),
		),
		mcp.WithNumber("LogoAdd",
			mcp.Description("为生成结果图添加标识的开关，默认为1。1：添加标识。0：不添加标识。"),
		),
		mcp.WithObject("LogoParam",
			mcp.Description("标识内容设置。"),
		),
		mcp.WithString("RspImgType",
			mcp.Description("返回图像方式（base64 或 url），二选一，默认为 base64。url 有效期为1小时。"),
		),
	)
	mcpsvr.AddTool(aiartTextToImageLite, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := aiart.NewClient(credential, region_, cpf)
		req := aiart.NewTextToImageLiteRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.TextToImageLite(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	aiartTextToImageRapid := mcp.NewTool(
		"aiart-TextToImageRapid",
		mcp.WithDescription(`混元生图（2.0）`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString("Prompt",
			mcp.Description("文本描述。算法将根据输入的文本智能生成与之相关的图像。建议详细描述画面主体、细节、场景等，文本描述越丰富，生成效果越精美。"),
		),
		mcp.WithString("Resolution",
			mcp.Description("生成图分辨率，默认1024:1024。支持的图像宽高比例: 1:1，3:4，4:3，9:16，16:9。支持的长边分辨率: 160，200，225，258，512，520，608，768，1024，1080，1280，1600，1620，1920，2048，2400，2560，2592，3440，3840，4096。"),
		),
		mcp.WithNumber("Seed",
			mcp.Description("随机种子，默认随机。0：随机种子生成。不传：随机种子生成。正数：固定种子生成。"),
		),
		mcp.WithObject("Image",
			mcp.Description("参考图。- Base64 和 Url 必须提供一个，如果都提供以 Url 为准。- 当传入Image参数时，Style和Resolution参数不生效，输出图分辨率将保持Image传入图分辨率。- 图片限制：单边分辨率大于128且小于2048；图片小于6M；格式支持 jpg、jpeg、png、bmp、tiff、webp。示例值：\"Image\":{\"Url\":\"https://cos.ap-guangzhou.myqcloud.com/image.jpg\", \"Base64\": \"\"}"),
		),
		mcp.WithString("Style",
			mcp.Description("生成的图片风格"),
		),
		mcp.WithNumber("LogoAdd",
			mcp.Description("为生成结果图添加标识的开关，默认为1。1：添加标识。0：不添加标识。"),
		),
		mcp.WithObject("LogoParam",
			mcp.Description("标识内容设置。默认在生成结果图右下角添加“图片由 AI 生成”字样，您可根据自身需要替换为其他的标识图片。"),
		),
		mcp.WithString("RspImgType",
			mcp.Description("返回图像方式（base64 或 url) ，二选一，默认为 base64。url 有效期为1小时。"),
		),
	)
	mcpsvr.AddTool(aiartTextToImageRapid, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := aiart.NewClient(credential, region_, cpf)
		req := aiart.NewTextToImageRapidRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.TextToImageRapid(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	aiartRefineImage := mcp.NewTool(
		"aiart-RefineImage",
		mcp.WithDescription(`图片变清晰`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InputUrl",
			mcp.Description("输入图 Url。"),
		),
		mcp.WithString(
			"InputImage",
			mcp.Description("输入图 Base64 数据。"),
		),
		mcp.WithString(
			"RspImgType",
			mcp.Description("返回图像方式（base64 或 url) ，二选一，默认为 base64。url 有效期为1小时。 示例值：url"),
		),
	)
	mcpsvr.AddTool(aiartRefineImage, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := aiart.NewClient(credential, region_, cpf)
		req := aiart.NewRefineImageRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.RefineImage(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	aiartImageToImage := mcp.NewTool(
		"aiart-ImageToImage",
		mcp.WithDescription(`图像风格化（图生图）`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InputImage",
			mcp.Description("输入图 Base64 数据。"),
		),
		mcp.WithString(
			"InputUrl",
			mcp.Description("输入图 Url。"),
		),
		mcp.WithString(
			"Prompt",
			mcp.Description("文本描述。"),
		),
		mcp.WithString(
			"NegativePrompt",
			mcp.Description("反向文本描述。"),
		),
		mcp.WithArray(
			"Styles",
			mcp.Description("绘画风格。"),
		),
		mcp.WithObject(
			"ResultConfig",
			mcp.Description("生成图结果的配置，包括输出图片分辨率和尺寸等。"),
		),
		mcp.WithNumber(
			"LogoAdd",
			mcp.Description("为生成结果图添加标识的开关，默认为1。"),
		),
		mcp.WithObject(
			"LogoParam",
			mcp.Description("标识内容设置。"),
		),
		mcp.WithNumber(
			"Strength",
			mcp.Description("生成自由度。"),
		),
		mcp.WithString(
			"RspImgType",
			mcp.Description("返回图像方式（base64 或 url) ，二选一，默认为 base64。url 有效期为1小时。"),
		),
		mcp.WithNumber(
			"EnhanceImage",
			mcp.Description("画质增强开关，默认关闭。"),
		),
		mcp.WithNumber(
			"RestoreFace",
			mcp.Description("细节优化的面部数量上限，支持0 ~ 6，默认为0。"),
		),
	)
	mcpsvr.AddTool(aiartImageToImage, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := aiart.NewClient(credential, region_, cpf)
		req := aiart.NewImageToImageRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ImageToImage(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	aiartImageOutpainting := mcp.NewTool(
		"aiart-ImageOutpainting",
		mcp.WithDescription(`扩图`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Ratio",
			mcp.Description("扩展后的比例（宽:高），需要不等于原图比例。"),
		),
		mcp.WithString(
			"InputImage",
			mcp.Description("输入图 Base64 数据。"),
		),
		mcp.WithString(
			"InputUrl",
			mcp.Description("输入图 Url。"),
		),
		mcp.WithString(
			"RspImgType",
			mcp.Description("返回图像方式（base64 或 url) ，二选一，默认为 base64。url 有效期为1小时。"),
		),
		mcp.WithNumber(
			"LogoAdd",
			mcp.Description("为生成结果图添加标识的开关，默认为1。"),
		),
		mcp.WithObject(
			"LogoParam",
			mcp.Description("标识内容设置。"),
		),
	)
	mcpsvr.AddTool(aiartImageOutpainting, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := aiart.NewClient(credential, region_, cpf)
		req := aiart.NewImageOutpaintingRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ImageOutpainting(req)
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
