package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterSSLTools 注册 SSL 相关工具（1个，现有迁移）
func RegisterSSLTools(s *server.MCPServer, cp security.CredentialProvider, g *security.Guard) {
	// DescribeDBInstanceSSLConfig - 查询实例SSL配置（只读）
	registerTool(s, cp, g, "DescribeDBInstanceSSLConfig", "查询实例SSL配置",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID，形如postgres-6bwgamo3")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeDBInstanceSSLConfigRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeDBInstanceSSLConfig(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("SSL tools registered: 1")
}
