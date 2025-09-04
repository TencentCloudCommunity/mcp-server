package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-ssl"
	mcpsvr := server.NewMCPServer(
		"腾讯云 SSL MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	sslCreateCertificate := mcp.NewTool(
		"ssl-CreateCertificate",
		mcp.WithDescription(`创建付费证书`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ProductId",
			mcp.Description("证书商品ID，3 = SecureSite 增强型企业版（EV Pro）， 4 = SecureSite 增强型（EV）， 5 = SecureSite 企业型专业版（OV Pro）， 6 = SecureSite 企业型（OV）， 7 = SecureSite 企业型（OV）通配符， 8 = Geotrust 增强型（EV）， 9 = Geotrust 企业型（OV）， 10 = Geotrust 企业型（OV）通配符， 11 = TrustAsia 域名型多域名 SSL 证书， 12 = TrustAsia 域名型（DV）通配符， 13 = TrustAsia 企业型通配符（OV）SSL 证书（D3）， 14 = TrustAsia 企业型（OV）SSL 证书（D3）， 15 = TrustAsia 企业型多域名 （OV）SSL 证书（D3）， 16 = TrustAsia 增强型 （EV）SSL 证书（D3）， 17 = TrustAsia 增强型多域名（EV）SSL 证书（D3）， 18 = GlobalSign 企业型（OV）SSL 证书， 19 = GlobalSign 企业型通配符 （OV）SSL 证书， 20 = GlobalSign 增强型 （EV）SSL 证书， 21 = TrustAsia 企业型通配符多域名（OV）SSL 证书（D3）， 22 = GlobalSign 企业型多域名（OV）SSL 证书， 23 = GlobalSign 企业型通配符多域名（OV）SSL 证书， 24 = GlobalSign 增强型多域名（EV）SSL 证书，25 = Wotrus 域名型证书，26 = Wotrus 域名型多域名证书，27 = Wotrus 域名型通配符证书，28 = Wotrus 企业型证书，29 = Wotrus 企业型多域名证书，30 = Wotrus 企业型通配符证书，31 = Wotrus 增强型证书，32 = Wotrus 增强型多域名证书，33 = Wotrus 国密域名型证书，34 = Wotrus 国密域名型多域名证书，35 = Wotrus 国密域名型通配符证书，37 = Wotrus 国密企业型证书，38 = Wotrus 国密企业型多域名证书，39 = Wotrus 国密企业型通配符证书，40 = Wotrus 国密增强型证书，41 = Wotrus 国密增强型多域名证书，42 = TrustAsia 域名型通配符多域名证书，43 = DNSPod-企业型(OV)SSL证书，44 = DNSPod-企业型(OV)通配符SSL证书，45 = DNSPod-企业型(OV)多域名SSL证书， 46 = DNSPod-增强型(EV)SSL证书，47 = DNSPod-增强型(EV)多域名SSL证书，48 = DNSPod-域名型(DV)SSL证书，49 = DNSPod-域名型(DV)通配符SSL证书，50 = DNSPod-域名型(DV)多域名SSL证书，51 = DNSPod（国密）-企业型(OV)SSL证书，52 = DNSPod（国密）-企业型(OV)通配符SSL证书，53 = DNSPod（国密）-企业型(OV)多域名SSL证书，54 = DNSPod（国密）-域名型(DV)SSL证书，55 = DNSPod（国密）-域名型(DV)通配符SSL证书， 56 = DNSPod（国密）-域名型(DV)多域名SSL证书，57 = SecureSite 企业型专业版多域名(OV Pro)，58 = SecureSite 企业型多域名(OV)，59 = SecureSite 增强型专业版多域名(EV Pro)，60 = SecureSite 增强型多域名(EV)，61 = Geotrust 增强型多域名(EV)"),
		),
		mcp.WithString(
			"DomainNum",
			mcp.Description("证书包含的域名数量"),
		),
		mcp.WithString(
			"TimeSpan",
			mcp.Description("证书年限"),
		),
		mcp.WithString(
			"AutoVoucher",
			mcp.Description("是否自动使用代金券：1是，0否；默认为1"),
		),
		mcp.WithString(
			"Tags",
			mcp.Description("标签， 生成证书打标签"),
		),
	)
	mcpsvr.AddTool(sslCreateCertificate, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ssl.NewClient(credential, region_, cpf)
		req := ssl.NewCreateCertificateRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateCertificate(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	sslCommitCertificateInformation := mcp.NewTool(
		"ssl-CommitCertificateInformation",
		mcp.WithDescription(`付费证书提交证书订单`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"CertificateId",
			mcp.Description("待提交资料的付费证书 ID。	"),
		),
		mcp.WithString(
			"VerifyType",
			mcp.Description("证书域名验证方式："),
		),
	)
	mcpsvr.AddTool(sslCommitCertificateInformation, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ssl.NewClient(credential, region_, cpf)
		req := ssl.NewCommitCertificateInformationRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CommitCertificateInformation(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	sslDescribeCertificateDetail := mcp.NewTool(
		"ssl-DescribeCertificateDetail",
		mcp.WithDescription(`获取证书详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"CertificateId",
			mcp.Description("证书 ID。"),
		),
	)
	mcpsvr.AddTool(sslDescribeCertificateDetail, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ssl.NewClient(credential, region_, cpf)
		req := ssl.NewDescribeCertificateDetailRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeCertificateDetail(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	sslUploadCertificate := mcp.NewTool(
		"ssl-UploadCertificate",
		mcp.WithDescription(`上传证书`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"CertificatePublicKey",
			mcp.Description("证书内容。"),
		),
		mcp.WithString(
			"CertificatePrivateKey",
			mcp.Description("私钥内容，证书类型为 SVR 时必填，为 CA 时可不填。"),
		),
		mcp.WithString(
			"CertificateType",
			mcp.Description("证书类型，默认 SVR。CA = CA证书，SVR = 服务器证书。"),
		),
		mcp.WithString(
			"Alias",
			mcp.Description("备注名称。"),
		),
		mcp.WithString(
			"ProjectId",
			mcp.Description("项目 ID。"),
		),
		mcp.WithString(
			"CertificateUse",
			mcp.Description("证书用途/证书来源。“CLB，CDN，WAF，LIVE，DDOS”"),
		),
		mcp.WithString(
			"Tags",
			mcp.Description("标签列表"),
		),
		mcp.WithString(
			"Repeatable",
			mcp.Description("相同的证书是否允许重复上传； true：允许上传相同指纹的证书；  false：不允许上传相同指纹的证书； 默认值：true"),
		),
	)
	mcpsvr.AddTool(sslUploadCertificate, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ssl.NewClient(credential, region_, cpf)
		req := ssl.NewUploadCertificateRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.UploadCertificate(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	sslDescribeCertificates := mcp.NewTool(
		"ssl-DescribeCertificates",
		mcp.WithDescription(`获取证书列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Offset",
			mcp.Description("分页偏移量，从0开始。 默认为0"),
		),
		mcp.WithString(
			"Limit",
			mcp.Description("每页数量，默认10。最大值1000，如超过1000按1000处理"),
		),
		mcp.WithString(
			"SearchKey",
			mcp.Description("搜索关键词，模糊匹配证书 ID、备注名称、证书域名"),
		),
		mcp.WithString(
			"CertificateType",
			mcp.Description("证书类型：CA = 客户端证书，SVR = 服务器证书。"),
		),
		mcp.WithString(
			"ProjectId",
			mcp.Description("项目 ID。"),
		),
		mcp.WithString(
			"ExpirationSort",
			mcp.Description("默认按照证书申请时间降序； 若传排序则按到期时间排序：DESC = 证书到期时间降序， ASC = 证书到期时间升序。"),
		),
		mcp.WithString(
			"CertificateStatus",
			mcp.Description("证书状态：0 = 审核中，1 = 已通过，2 = 审核失败，3 = 已过期，4 = 已添加DNS记录，5 = 企业证书，待提交，6 = 订单取消中，7 = 已取消，8 = 已提交资料， 待上传确认函，9 = 证书吊销中，10 = 已吊销，11 = 重颁发中，12 = 待上传吊销确认函，13 = 免费证书待提交资料。14 = 已退款。 15 = 证书迁移中"),
		),
		mcp.WithString(
			"Deployable",
			mcp.Description("是否可部署，可选值：1 = 可部署，0 =  不可部署。"),
		),
		mcp.WithString(
			"Upload",
			mcp.Description("是否筛选上传托管的 1筛选，0不筛选"),
		),
		mcp.WithString(
			"Renew",
			mcp.Description("是否筛选可续期证书 1筛选 0不筛选"),
		),
		mcp.WithString(
			"FilterSource",
			mcp.Description("筛选来源， upload：上传证书， buy：腾讯云证书， 不传默认全部"),
		),
		mcp.WithString(
			"IsSM",
			mcp.Description("是否筛选国密证书。1:筛选  0:不筛选"),
		),
		mcp.WithString(
			"FilterExpiring",
			mcp.Description("筛选证书是否即将过期，传1是筛选，0不筛选"),
		),
		mcp.WithString(
			"Hostable",
			mcp.Description("是否可托管，可选值：1 = 可托管，0 =  不可托管。"),
		),
		mcp.WithString(
			"Tags",
			mcp.Description("筛选指定标签的证书"),
		),
		mcp.WithString(
			"IsPendingIssue",
			mcp.Description("是否筛选等待签发的证书，传1是筛选，0和null不筛选"),
		),
		mcp.WithString(
			"CertIds",
			mcp.Description("筛选指定证书ID的证书，只支持有权限的证书ID"),
		),
	)
	mcpsvr.AddTool(sslDescribeCertificates, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ssl.NewClient(credential, region_, cpf)
		req := ssl.NewDescribeCertificatesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeCertificates(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	sslDescribeManagerDetail := mcp.NewTool(
		"ssl-DescribeManagerDetail",
		mcp.WithDescription(`查询管理人详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ManagerId",
			mcp.Description("管理人ID,可以从describeManagers接口获得"),
		),
		mcp.WithString(
			"Limit",
			mcp.Description("分页每页数量"),
		),
		mcp.WithString(
			"Offset",
			mcp.Description("分页偏移量"),
		),
	)
	mcpsvr.AddTool(sslDescribeManagerDetail, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ssl.NewClient(credential, region_, cpf)
		req := ssl.NewDescribeManagerDetailRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeManagerDetail(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	sslDescribeManagers := mcp.NewTool(
		"ssl-DescribeManagers",
		mcp.WithDescription(`查询管理人列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"CompanyId",
			mcp.Description("公司ID,可以从DescribeCompanies接口获取"),
		),
		mcp.WithString(
			"Offset",
			mcp.Description("分页偏移量，如果不传默认值为0"),
		),
		mcp.WithString(
			"Limit",
			mcp.Description("分页每页数量，如果不传默认值为10，最大值为1000"),
		),
		mcp.WithString(
			"ManagerName",
			mcp.Description("管理人姓名（将废弃），请使用SearchKey"),
		),
		mcp.WithString(
			"ManagerMail",
			mcp.Description("模糊查询管理人邮箱（将废弃），请使用SearchKey"),
		),
		mcp.WithString(
			"Status",
			mcp.Description("根据管理人状态进行筛选，取值有"),
		),
		mcp.WithString(
			"SearchKey",
			mcp.Description("根据这样的格式:管理人姓|管理人名|邮箱|部门 ,进行精准匹配"),
		),
	)
	mcpsvr.AddTool(sslDescribeManagers, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ssl.NewClient(credential, region_, cpf)
		req := ssl.NewDescribeManagersRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeManagers(req)
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
