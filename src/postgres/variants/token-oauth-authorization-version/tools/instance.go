package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterInstanceTools 注册实例相关工具（15个：2现有迁移 + 7查询 + 6管理）
func RegisterInstanceTools(s *server.MCPServer, cp security.CredentialProvider, g *security.Guard) {
	// ===== 现有迁移（2个）=====

	// DescribeDBInstanceAttribute - 查询实例详情（只读）
	registerTool(s, cp, g, "DescribeDBInstanceAttribute", "查询实例详情",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeDBInstanceAttributeRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeDBInstanceAttribute(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// UpgradeDBInstanceKernelVersion - 升级实例内核版本号（写，L2业务确认）
	registerTool(s, cp, g, "UpgradeDBInstanceKernelVersion", "升级实例内核版本号",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
			mcp.WithString("TargetDBKernelVersion", mcp.Description("升级的目标内核版本号")),
			mcp.WithNumber("SwitchTag", mcp.Description("指定切换时间。可选值: 0-立即切换, 1-维护时间切换, 2-指定时间切换")),
			mcp.WithString("SwitchStartTime", mcp.Description("切换开始时间，格式HH:MM:SS")),
			mcp.WithString("SwitchEndTime", mcp.Description("切换截止时间，格式HH:MM:SS")),
			mcp.WithBoolean("DryRun", mcp.Description("是否执行预检查")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewUpgradeDBInstanceKernelVersionRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.UpgradeDBInstanceKernelVersion(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 实例查询组（4个，只读）=====

	// DescribeDBInstances - 查询实例列表
	registerTool(s, cp, g, "DescribeDBInstances", "查询实例列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
			mcp.WithArray("Filters", mcp.Description("过滤条件")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目，默认20")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量，从0开始")),
			mcp.WithString("OrderBy", mcp.Description("排序字段")),
			mcp.WithString("OrderByType", mcp.Description("排序方式: asc|desc")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			normalizeDescribeDBInstancesArgs(args)
			req := postgres.NewDescribeDBInstancesRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeDBInstances(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeClasses - 查询可用规格
	registerTool(s, cp, g, "DescribeClasses", "查询可用规格列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("Zone", mcp.Description("可用区ID")),
			mcp.WithString("DBEngine", mcp.Description("数据库引擎，默认postgresql")),
			mcp.WithString("DBMajorVersion", mcp.Description("数据库主版本号")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeClassesRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeClasses(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeDBVersions - 查询可用数据库版本
	registerTool(s, cp, g, "DescribeDBVersions", "查询可用数据库版本",
		security.LevelNone,
		[]mcp.ToolOption{},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			normalizeDescribeDBVersionsArgs(args)
			req := postgres.NewDescribeDBVersionsRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeDBVersions(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeRegions - 查询售卖地域（只读，建实例前选地域）
	registerTool(s, cp, g, "DescribeRegions", "查询售卖地域",
		security.LevelNone,
		[]mcp.ToolOption{},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeRegionsRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeRegions(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeZones - 查询售卖可用区（只读，与DescribeRegions配对，先选地域再选可用区）
	registerTool(s, cp, g, "DescribeZones", "查询售卖可用区",
		security.LevelNone,
		[]mcp.ToolOption{},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeZonesRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeZones(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeProductConfig - 查询售卖规格配置（只读，一站式规格配置查询）
	registerTool(s, cp, g, "DescribeProductConfig", "查询售卖规格配置",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("Zone", mcp.Description("可用区名称")),
			mcp.WithString("DBEngine", mcp.Description("数据库引擎，默认postgresql")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeProductConfigRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeProductConfig(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeTasks - 查询异步任务状态
	registerTool(s, cp, g, "DescribeTasks", "查询异步任务状态",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeTasksRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeTasks(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 实例管理组（6个，写操作需guard）=====

	// CreateInstances - 创建实例（L1费用确认）
	// 兼容旧参数别名：InstanceSpec->SpecCode、Volume->Storage、DBCharset->Charset、InstanceName->Name。
	registerTool(s, cp, g, "CreateInstances", "创建实例",
		security.LevelFee,
		[]mcp.ToolOption{
			mcp.WithString("Zone", mcp.Required(), mcp.Description("可用区ID")),
			mcp.WithString("SpecCode", mcp.Description("售卖规格码，可由 DescribeClasses 获取")),
			mcp.WithNumber("Storage", mcp.Description("实例磁盘容量(GB)")),
			mcp.WithNumber("InstanceCount", mcp.Description("购买实例数量，默认1")),
			mcp.WithNumber("Period", mcp.Description("购买时长(月)")),
			mcp.WithString("Charset", mcp.Description("数据库字符集，如 UTF8")),
			mcp.WithString("AdminName", mcp.Description("实例管理员账号")),
			mcp.WithString("AdminPassword", mcp.Description("实例管理员密码")),
			mcp.WithString("DBMajorVersion", mcp.Description("PostgreSQL 大版本号，如 18")),
			mcp.WithString("DBVersion", mcp.Description("社区版本号，可选")),
			mcp.WithString("DBKernelVersion", mcp.Description("内核版本号，可选")),
			mcp.WithString("InstanceChargeType", mcp.Description("计费类型: POSTPAID_BY_HOUR|PREPAID")),
			mcp.WithString("VpcId", mcp.Description("私有网络ID")),
			mcp.WithString("SubnetId", mcp.Description("子网ID")),
			mcp.WithNumber("AutoRenewFlag", mcp.Description("续费标记：0手动续费，1自动续费")),
			mcp.WithString("Name", mcp.Description("实例名称")),
			mcp.WithArray("SecurityGroupIds", mcp.Description("安全组ID列表")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			normalizeCreateInstancesArgs(args)
			req := postgres.NewCreateInstancesRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.CreateInstances(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ModifyDBInstanceName - 修改实例名称（L4审计）
	registerTool(s, cp, g, "ModifyDBInstanceName", "修改实例名称",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("InstanceName", mcp.Required(), mcp.Description("新实例名称")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewModifyDBInstanceNameRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.ModifyDBInstanceName(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ModifyDBInstanceSpec - 变更实例规格（L1费用确认）
	// 兼容旧参数别名 Volume->Storage。
	registerTool(s, cp, g, "ModifyDBInstanceSpec", "变更实例规格(扩缩容)",
		security.LevelFee,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithNumber("Memory", mcp.Description("修改后的内存(GiB)")),
			mcp.WithNumber("Storage", mcp.Description("修改后的磁盘(GiB)")),
			mcp.WithNumber("Cpu", mcp.Description("修改后的 CPU 核数，可选")),
			mcp.WithNumber("SwitchTag", mcp.Description("切换时间选项：0立即切换，1指定时间，2维护窗口")),
			mcp.WithString("SwitchStartTime", mcp.Description("切换开始时间，HH:MM:SS")),
			mcp.WithString("SwitchEndTime", mcp.Description("切换截止时间，HH:MM:SS")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			if _, ok := args["Storage"]; !ok {
				if legacy, ok := args["Volume"]; ok {
					args["Storage"] = legacy
					delete(args, "Volume")
				}
			}
			delete(args, "InstanceType")
			req := postgres.NewModifyDBInstanceSpecRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.ModifyDBInstanceSpec(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// RestartDBInstance - 重启实例（L2业务确认）
	registerTool(s, cp, g, "RestartDBInstance", "重启实例",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewRestartDBInstanceRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.RestartDBInstance(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// IsolateDBInstances - 隔离实例（L2业务确认）
	// SDK 需要 DBInstanceIdSet；为兼容易用性保留 DBInstanceId 单实例别名。
	registerTool(s, cp, g, "IsolateDBInstances", "隔离实例",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("单实例ID，和 DBInstanceIdSet 二选一")),
			mcp.WithArray("DBInstanceIdSet", mcp.Description("实例ID数组，建议只传一个")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			if _, ok := args["DBInstanceIdSet"]; !ok {
				if id, ok := args["DBInstanceId"].(string); ok && id != "" {
					args["DBInstanceIdSet"] = []string{id}
				}
			}
			delete(args, "DBInstanceId")
			req := postgres.NewIsolateDBInstancesRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.IsolateDBInstances(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DisIsolateDBInstances - 解除隔离（L4审计）
	// SDK 需要 DBInstanceIdSet；为兼容易用性保留 DBInstanceId 单实例别名。
	registerTool(s, cp, g, "DisIsolateDBInstances", "解除隔离实例",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("单实例ID，和 DBInstanceIdSet 二选一")),
			mcp.WithArray("DBInstanceIdSet", mcp.Description("实例ID数组，建议只传一个")),
			mcp.WithNumber("Period", mcp.Description("购买时长(月)，预付费实例可用")),
			mcp.WithBoolean("AutoVoucher", mcp.Description("是否自动使用代金券")),
			mcp.WithArray("VoucherIds", mcp.Description("代金券ID列表")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			if _, ok := args["DBInstanceIdSet"]; !ok {
				if id, ok := args["DBInstanceId"].(string); ok && id != "" {
					args["DBInstanceIdSet"] = []string{id}
				}
			}
			delete(args, "DBInstanceId")
			req := postgres.NewDisIsolateDBInstancesRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DisIsolateDBInstances(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Instance tools registered: 15")
}
