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
	mcpServerName := "mcp-server-ocr-document-multimodal"
	mcpsvr := server.NewMCPServer(
		"腾讯云API MCP，用于将自然语言转换为API调用",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	ocrExtractDocMulti := mcp.NewTool(
		"ocr-ExtractDocMulti",
		mcp.WithDescription(`文档抽取（多模态版）`),
		mcp.WithString("region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ImageUrl",
			mcp.Description("图片的 Url 地址。支持的图片格式：PNG、JPG、JPEG，WORD，EXCEL，暂不支持 GIF 格式。支持的图片大小：所下载图片经 Base64 编码后不超过 10M。图片下载时间不超过 3 秒。支持的图片像素：需介于20-10000px之间。图片存储于腾讯云的 Url 可保障更高的下载速度和稳定性，建议图片存储于腾讯云。非腾讯云存储的 Url 速度和稳定性可能受一定影响。"),
		),
		mcp.WithString(
			"ImageBase64",
			mcp.Description("图片的 Base64 值。支持的图片格式：PNG、JPG、JPEG，WORD，EXCEL，暂不支持 GIF 格式。支持的图片大小：所下载图片经Base64编码后不超过 10M。图片下载时间不超过 3 秒。支持的图片像素：需介于20-10000px之间。图片的 ImageUrl、ImageBase64 必须提供一个，如果都提供，只使用 ImageUrl。"),
		),
		mcp.WithNumber(
			"PdfPageNumber",
			mcp.Description("需要识别的PDF页面的对应页码，仅支持PDF单页识别，当上传文件为PDF且IsPdf参数值为true时有效，默认值为前3页。"),
		),
		mcp.WithString(
			"ItemNames",
			mcp.Description("自定义结构化功能需返回的字段名称，例：若客户想新增返回姓名、性别两个字段的识别结果，则输入ItemNames=[\"姓名\",\"性别\"]"),
		),
		mcp.WithNumber(
			"ItemNamesShowMode",
			mcp.Description("true：仅输出自定义字段 false：输出默认字段+自定义字段 默认true"),
		),
		mcp.WithString(
			"ReturnFullText",
			mcp.Description("是否开启全文字段识别"),
		),
		mcp.WithNumber(
			"ConfigId",
			mcp.Description("配置id支持：General -- 通用场景 InvoiceEng -- 国际invoice模板 WayBillEng --海运订单模板 CustomsDeclaration -- 进出口报关单 WeightNote -- 磅单 MedicalMeter -- 血压仪表识别 BillOfLading -- 海运提单 EntrustmentBook -- 海运托书 Statement -- 对账单识别模板 BookingConfirmation -- 配舱通知书识别模板 AirWayBill -- 航空运单识别模板 Table -- 表格模板 SteelLabel -- 实物标签识别模板 CarInsurance -- 车辆保险单识别模板 MultiRealEstateCertificate -- 房产材料识别模板 MultiRealEstateMaterial -- 房产证明识别模板 HongKongUtilityBill -- 中国香港水电煤单识别模板 OverseasCheques -- 海外支票\nRegistrationCertificate -- 备案证 ​GridPhoto -- 电网系统照片 ​SignaturePage -- 签署页"),
		),
		mcp.WithString(
			"EnableCoord",
			mcp.Description("是否开启全文字段坐标值的识别"),
		),
		mcp.WithNumber(
			"OutputParentKey",
			mcp.Description("是否开启父子key识别，默认是"),
		),
		mcp.WithString(
			"ConfigAdvanced",
			mcp.Description("模板的单个属性配置"),
		),
		mcp.WithNumber(
			"OutputLanguage",
			mcp.Description("cn时，添加的key为中文 en时，添加的key为英语"),
		),
	)
	mcpsvr.AddTool(ocrExtractDocMulti, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ocr.NewClient(credential, region_, cpf)
		req := ocr.NewExtractDocMultiRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ExtractDocMulti(req)
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
