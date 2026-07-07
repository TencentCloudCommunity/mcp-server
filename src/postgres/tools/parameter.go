package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterParameterTools 注册参数管理工具（5个：3现有迁移 + 2新增）
func RegisterParameterTools(s *server.MCPServer, cp security.CredentialProvider, g *security.Guard) {
	// ===== 现有迁移（3个）=====

	// DescribeDBInstanceParameters - 查询实例参数（只读）
	registerTool(s, cp, g, "DescribeDBInstanceParameters", "查询实例参数",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
			mcp.WithString("ParamName", mcp.Description("查询指定参数详情。为空返回全部参数")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeDBInstanceParametersRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeDBInstanceParameters(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeParameterTemplates - 查询参数模板列表（只读）
	registerTool(s, cp, g, "DescribeParameterTemplates", "查询参数模板列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithArray("Filters", mcp.Description("过滤条件: TemplateName|TemplateId|DBMajorVersion|DBEngine")),
			mcp.WithNumber("Limit", mcp.Description("每页显示数量[0,100]，默认20")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量")),
			mcp.WithString("OrderBy", mcp.Description("排序指标: CreateTime|TemplateName|DBMajorVersion")),
			mcp.WithString("OrderByType", mcp.Description("排序方式: asc|desc")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeParameterTemplatesRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeParameterTemplates(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeParameterTemplateAttributes - 查询参数模板详情（只读）
	registerTool(s, cp, g, "DescribeParameterTemplateAttributes", "查询参数模板详情",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("TemplateId", mcp.Required(), mcp.Description("参数模板ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeParameterTemplateAttributesRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeParameterTemplateAttributes(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeParamsEvent - 查询参数修改事件（只读，排障时追溯"谁在什么时候改了什么参数"）
	registerTool(s, cp, g, "DescribeParamsEvent", "查询参数修改事件",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeParamsEventRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeParamsEvent(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 新增（2个）=====

	// ModifyDBInstanceParameters - 修改实例参数（L3最高级确认，高危参数如max_connections）
	registerTool(s, cp, g, "ModifyDBInstanceParameters", "修改实例参数",
		security.LevelCritical,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithArray("ParamList", mcp.Required(), mcp.Description("参数列表，每项含 Name 和 ExpectedValue；也兼容更直观的 Value 写法")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			normalizeModifyDBInstanceParametersArgs(args)
			req := postgres.NewModifyDBInstanceParametersRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.ModifyDBInstanceParameters(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Parameter tools registered: 5")
}
