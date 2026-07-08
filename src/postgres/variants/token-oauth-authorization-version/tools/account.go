package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterAccountTools 注册账号管理工具（6个：1现有迁移 + 5新增）
func RegisterAccountTools(s *server.MCPServer, cp security.CredentialProvider, g *security.Guard) {
	// ===== 现有迁移（1个）=====

	// DescribeAccounts - 查询数据库账号列表（只读）
	registerTool(s, cp, g, "DescribeAccounts", "查询实例的数据库账号列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID，形如postgres-6fego161")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目，默认20，取值1-100")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量，从0开始")),
			mcp.WithString("OrderBy", mcp.Description("排序字段: createTime|name|updateTime")),
			mcp.WithString("OrderByType", mcp.Description("排序方式: desc|asc")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeAccountsRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeAccounts(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 新增（5个）=====

	// CreateAccount - 创建账号（L4审计）
	// 注意1：SDK 字段名是 Remark（无 s），此前误写为 Remarks，属于未知字段，
	// 会被 FromJsonString 拒绝解析导致整个请求体被吞掉，现已修正。
	// 注意2：模型注释未标注 Type 必填，但实测该接口会返回
	// MissingParameter: 请求缺少必传参数 `Type`，故将其标为必填参数，避免误导调用方。
	registerTool(s, cp, g, "CreateAccount", "创建数据库账号",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("UserName", mcp.Required(), mcp.Description("账号名")),
			mcp.WithString("Password", mcp.Required(), mcp.Description("账号密码")),
			mcp.WithString("Type", mcp.Required(), mcp.Description("账号类型：normal普通用户|tencentDBSuper超级用户")),
			mcp.WithString("Remark", mcp.Description("备注")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewCreateAccountRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.CreateAccount(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DeleteAccount - 删除账号（L2业务确认，误删导致应用断连）
	registerTool(s, cp, g, "DeleteAccount", "删除数据库账号",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("UserName", mcp.Required(), mcp.Description("账号名")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDeleteAccountRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DeleteAccount(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ModifyAccountPrivileges - 修改账号权限（L4审计，提权审计）
	// 注意：腾讯云 API 的 ModifyAccountPrivileges 只接受 DBInstanceId/UserName/ModifyPrivilegeSet
	// 三个顶层字段（SDK FromJsonString 对未知顶层字段会直接报错拒绝解析），不存在 DBName/Privileges
	// 这种扁平字符串参数。ModifyPrivilegeSet 是嵌套结构，每项形如：
	//   {
	//     "DatabasePrivilege": {
	//       "Object": {"ObjectType":"database|schema|table|...","ObjectName":"...","DatabaseName":"...","SchemaName":"...","TableName":"..."},
	//       "PrivilegeSet": ["SELECT","INSERT",...]
	//     },
	//     "ModifyType": "grantObject|revokeObject|alterRole",
	//     "IsCascade": false
	//   }
	// 调用方需按此嵌套结构直接传入 ModifyPrivilegeSet 数组（与 ModifyDBInstanceParameters 的
	// ParamList 用法一致，工具层不做字段名转换）。
	registerTool(s, cp, g, "ModifyAccountPrivileges", "修改账号权限（授权/收回/修改账号类型）",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("UserName", mcp.Required(), mcp.Description("账号名，可通过DescribeAccounts接口获取")),
			mcp.WithArray("ModifyPrivilegeSet", mcp.Required(), mcp.Description(
				"修改的权限信息数组，一次最高修改50条。每项结构："+
					"{DatabasePrivilege:{Object:{ObjectType(database|schema|table|...),ObjectName(必填,nullable为false),DatabaseName,SchemaName,TableName},PrivilegeSet:[权限字符串数组]},"+
					"ModifyType(grantObject授权|revokeObject收回|alterRole修改账号类型),IsCascade(仅revokeObject时可用，是否级联撤销，默认false)}")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewModifyAccountPrivilegesRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.ModifyAccountPrivileges(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ResetAccountPassword - 重置账号密码（L4审计，旧连接断开）
	registerTool(s, cp, g, "ResetAccountPassword", "重置账号密码",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("UserName", mcp.Required(), mcp.Description("账号名")),
			mcp.WithString("Password", mcp.Required(), mcp.Description("新密码")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewResetAccountPasswordRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.ResetAccountPassword(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeAccountPrivileges - 查询账号权限（只读，与ModifyAccountPrivileges配对使用）
	registerTool(s, cp, g, "DescribeAccountPrivileges", "查询数据库账号的权限信息",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("UserName", mcp.Description("账号名，可通过DescribeAccounts接口获取")),
			mcp.WithArray("DatabaseObjectSet", mcp.Description("要查询的数据库对象信息列表，每项含ObjectType/ObjectName/DatabaseName/SchemaName/TableName")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeAccountPrivilegesRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeAccountPrivileges(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Account tools registered: 6")
}
