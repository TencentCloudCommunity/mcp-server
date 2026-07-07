package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterMonitoringTools 注册监控/日志工具（3个，全只读）
func RegisterMonitoringTools(s *server.MCPServer, cp security.CredentialProvider, g *security.Guard) {
	// DescribeSlowQueryList - 查询慢查询列表（只读）
	registerTool(s, cp, g, "DescribeSlowQueryList", "查询慢查询列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("StartTime", mcp.Description("开始时间")),
			mcp.WithString("EndTime", mcp.Description("结束时间")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeSlowQueryListRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeSlowQueryList(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeSlowQueryAnalysis - 慢查询分析（只读）
	registerTool(s, cp, g, "DescribeSlowQueryAnalysis", "慢查询分析",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("StartTime", mcp.Description("开始时间")),
			mcp.WithString("EndTime", mcp.Description("结束时间")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeSlowQueryAnalysisRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeSlowQueryAnalysis(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeDBErrlogs - 查询错误日志（只读）
	registerTool(s, cp, g, "DescribeDBErrlogs", "查询错误日志",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("StartTime", mcp.Description("开始时间")),
			mcp.WithString("EndTime", mcp.Description("结束时间")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeDBErrlogsRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeDBErrlogs(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Monitoring tools registered: 3")
}
