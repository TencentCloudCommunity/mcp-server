package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterNetworkTools 注册网络和安全组工具（4个，全新增）
func RegisterNetworkTools(s *server.MCPServer, cp security.CredentialProvider, g *security.Guard) {
	// ===== 网络组（2个）=====

	// OpenDBExtranetAccess - 开启公网访问（L2业务确认，暴露公网）
	// 注意：SDK 只认 IsIpv6（1=开通Ipv6外网,0=否，默认0），此前误写为 WanPort，
	// 属于未知字段，会被 FromJsonString 拒绝解析导致整个请求体被吞掉，现已修正。
	registerTool(s, cp, g, "OpenDBExtranetAccess", "开启实例公网访问",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithNumber("IsIpv6", mcp.Description("是否开通Ipv6外网，1：是，0：否，默认0")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewOpenDBExtranetAccessRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.OpenDBExtranetAccess(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// CloseDBExtranetAccess - 关闭公网访问（L2业务确认，中断公网连接）
	registerTool(s, cp, g, "CloseDBExtranetAccess", "关闭实例公网访问",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithNumber("IsIpv6", mcp.Description("是否关闭Ipv6外网，1：是，0：否，默认0")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewCloseDBExtranetAccessRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.CloseDBExtranetAccess(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 安全组组（2个）=====

	// DescribeDBInstanceSecurityGroups - 查询实例安全组（只读）
	registerTool(s, cp, g, "DescribeDBInstanceSecurityGroups", "查询实例安全组",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeDBInstanceSecurityGroupsRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeDBInstanceSecurityGroups(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ModifyDBInstanceSecurityGroups - 修改实例安全组（L2业务确认，误改断连）
	// 注意：SDK 字段名是 SecurityGroupIdSet，此前误写为 SecurityGroupIds 会导致
	// FromJsonString 直接报 unknown keys。这里保留旧别名兼容，统一转换后再下发。
	registerTool(s, cp, g, "ModifyDBInstanceSecurityGroups", "修改实例安全组",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID，与 ReadOnlyGroupId 二选一")),
			mcp.WithString("ReadOnlyGroupId", mcp.Description("只读组ID，与 DBInstanceId 二选一")),
			mcp.WithArray("SecurityGroupIdSet", mcp.Required(), mcp.Description("安全组ID全量列表")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			normalizeModifyDBInstanceSecurityGroupsArgs(args)
			req := postgres.NewModifyDBInstanceSecurityGroupsRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.ModifyDBInstanceSecurityGroups(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Network tools registered: 4")
}
