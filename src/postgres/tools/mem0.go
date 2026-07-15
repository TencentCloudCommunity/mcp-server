package tools

import (
	"encoding/json"

	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterMem0Tools 注册 Mem0 AI 服务相关工具（3个）
func RegisterMem0Tools(s *server.MCPServer, cp security.CredentialProvider, g *security.Guard) {
	// OpenMem0Service - 开启实例Mem0服务（L1费用确认）
	registerTool(s, cp, g, "OpenMem0Service", "开启实例Mem0服务",
		security.LevelFee,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("AgenticBaseId", mcp.Required(), mcp.Description("AgenticBaseID，开启Mem0服务前请先开通AgenticBase套餐")),
			mcp.WithString("LLMModel", mcp.Required(), mcp.Description(
				"Mem0服务使用的LLM模型。枚举值：auto/deepseek-v4-flash/deepseek-v4-pro/glm-5/glm-5-turbo/glm-5.1/kimi-k2.5/kimi-k2.6/minimax-m2.5/minimax-m2.7")),
			mcp.WithString("EmbeddingApiKey", mcp.Required(), mcp.Description("腾讯云TokenHub的Embedding ApiKey")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			return callMem0API(client, "OpenMem0Service", args)
		})

	// CloseMem0Service - 关闭实例Mem0服务（L2业务确认，删除容器+回收网络）
	registerTool(s, cp, g, "CloseMem0Service", "关闭实例Mem0服务",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			return callMem0API(client, "CloseMem0Service", args)
		})

	// DescribeMem0Service - 查询Mem0服务信息（只读）
	registerTool(s, cp, g, "DescribeMem0Service", "查询实例Mem0服务详情",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			return callMem0API(client, "DescribeMem0Service", args)
		})

	Log("Mem0 tools registered: 3")
}

// callMem0API 使用 CommonRequest 机制调用尚未进入 SDK 的 Mem0 API。
// 返回值已剥离外层 Response 包装，与 SDK 生成类型的 ToJsonString() 格式保持一致。
func callMem0API(client *postgres.Client, action string, args map[string]interface{}) (string, error) {
	req := tchttp.NewCommonRequest("postgres", "2017-03-12", action)
	if err := req.SetActionParameters(args); err != nil {
		return "", err
	}
	resp := tchttp.NewCommonResponse()
	if err := client.Send(req, resp); err != nil {
		return "", err
	}
	// CommonResponse.GetBody() 返回 {"Response":{...}}，
	// 需提取内层 Response 对象以与其他工具的返回格式对齐。
	raw := resp.GetBody()
	var wrapper struct {
		Response map[string]interface{} `json:"Response"`
	}
	if err := json.Unmarshal(raw, &wrapper); err != nil || wrapper.Response == nil {
		// 解包失败时回退返回原始内容
		return string(raw), nil
	}
	inner, err := json.Marshal(wrapper.Response)
	if err != nil {
		return string(raw), nil
	}
	return string(inner), nil
}
