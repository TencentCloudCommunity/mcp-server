package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterReadonlyTools 注册只读实例工具（2个，全新增）
func RegisterReadonlyTools(s *server.MCPServer, cp security.CredentialProvider, g *security.Guard) {
	// DescribeReadOnlyGroups - 查询只读组列表（只读）
	registerTool(s, cp, g, "DescribeReadOnlyGroups", "查询只读组列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithArray("Filters", mcp.Description("过滤条件，支持 db-master-instance-id|read-only-group-id；其中 db-master-instance-id 为必填项")),
			mcp.WithNumber("PageSize", mcp.Description("每页返回数目，默认10，最大99")),
			mcp.WithNumber("PageNumber", mcp.Description("页码，默认1")),
			mcp.WithString("OrderBy", mcp.Description("排序字段：ROGroupId|CreateTime|Name")),
			mcp.WithString("OrderByType", mcp.Description("排序方式：asc|desc")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			normalizeDescribeReadOnlyGroupsArgs(args)
			req := postgres.NewDescribeReadOnlyGroupsRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeReadOnlyGroups(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// CreateReadOnlyDBInstance - 创建只读实例（L1费用确认）
	// 兼容旧参数别名：DBInstanceId->MasterDBInstanceId、SpecName->SpecCode、InstanceName->Name。
	registerTool(s, cp, g, "CreateReadOnlyDBInstance", "创建只读实例",
		security.LevelFee,
		[]mcp.ToolOption{
			mcp.WithString("MasterDBInstanceId", mcp.Description("主实例ID")),
			mcp.WithString("DBInstanceId", mcp.Description("主实例ID旧别名，兼容保留")),
			mcp.WithString("SpecCode", mcp.Description("售卖规格码")),
			mcp.WithString("SpecName", mcp.Description("售卖规格码旧别名，兼容保留")),
			mcp.WithNumber("Storage", mcp.Description("实例硬盘容量(GB)")),
			mcp.WithNumber("InstanceCount", mcp.Description("购买数量，默认1")),
			mcp.WithNumber("Period", mcp.Description("购买时长(月)")),
			mcp.WithString("Zone", mcp.Description("可用区")),
			mcp.WithString("VpcId", mcp.Description("私有网络ID")),
			mcp.WithString("SubnetId", mcp.Description("子网ID")),
			mcp.WithString("InstanceChargeType", mcp.Description("计费类型: POSTPAID_BY_HOUR|PREPAID")),
			mcp.WithNumber("AutoRenewFlag", mcp.Description("续费标记：0手动续费，1自动续费")),
			mcp.WithString("Name", mcp.Description("只读实例名称")),
			mcp.WithString("InstanceName", mcp.Description("只读实例名称旧别名，兼容保留")),
			mcp.WithString("ReadOnlyGroupId", mcp.Description("只读组ID，可选")),
			mcp.WithArray("SecurityGroupIds", mcp.Description("安全组ID列表")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			normalizeCreateReadOnlyDBInstanceArgs(args)
			req := postgres.NewCreateReadOnlyDBInstanceRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.CreateReadOnlyDBInstance(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Readonly tools registered: 2")
}
