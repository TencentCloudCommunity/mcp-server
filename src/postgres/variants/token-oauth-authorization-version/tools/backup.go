package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterBackupTools 注册备份恢复工具（8个，全新增）
func RegisterBackupTools(s *server.MCPServer, cp security.CredentialProvider, g *security.Guard) {
	// ===== 只读（5个）=====

	// DescribeBackupOverview - 查询备份概览（只读）
	registerTool(s, cp, g, "DescribeBackupOverview", "查询备份概览",
		security.LevelNone,
		[]mcp.ToolOption{},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			normalizeDescribeBackupOverviewArgs(args)
			req := postgres.NewDescribeBackupOverviewRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeBackupOverview(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeBaseBackups - 查询基础备份列表（只读）
	registerTool(s, cp, g, "DescribeBaseBackups", "查询基础备份列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("MinFinishTime", mcp.Description("备份最小结束时间，形如 2018-01-01 00:00:00")),
			mcp.WithString("MaxFinishTime", mcp.Description("备份最大结束时间，形如 2018-01-01 00:00:00")),
			mcp.WithArray("Filters", mcp.Description("过滤条件，支持 db-instance-id|db-instance-name|db-instance-ip|base-backup-id|db-instance-status")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量")),
			mcp.WithString("OrderBy", mcp.Description("排序字段：StartTime|FinishTime|Size")),
			mcp.WithString("OrderByType", mcp.Description("排序方式：asc|desc")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			normalizeBackupListArgs(args)
			req := postgres.NewDescribeBaseBackupsRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeBaseBackups(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeLogBackups - 查询日志备份列表（只读）
	registerTool(s, cp, g, "DescribeLogBackups", "查询日志备份列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("MinFinishTime", mcp.Description("备份最小结束时间，形如 2018-01-01 00:00:00")),
			mcp.WithString("MaxFinishTime", mcp.Description("备份最大结束时间，形如 2018-01-01 00:00:00")),
			mcp.WithArray("Filters", mcp.Description("过滤条件，支持 db-instance-id|db-instance-name|db-instance-ip|db-instance-status")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量")),
			mcp.WithString("OrderBy", mcp.Description("排序字段：StartTime|FinishTime|Size")),
			mcp.WithString("OrderByType", mcp.Description("排序方式：asc|desc")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			normalizeBackupListArgs(args)
			req := postgres.NewDescribeLogBackupsRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeLogBackups(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeAvailableRecoveryTime - 查询可恢复时间范围（只读）
	registerTool(s, cp, g, "DescribeAvailableRecoveryTime", "查询可恢复时间范围",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeAvailableRecoveryTimeRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeAvailableRecoveryTime(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeCloneDBInstanceSpec - 查询克隆实例可购买的规格（只读，与CloneDBInstance配对使用）
	registerTool(s, cp, g, "DescribeCloneDBInstanceSpec", "查询克隆实例可购买的规格",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("源实例ID")),
			mcp.WithString("BackupSetId", mcp.Description("基础备份集ID，与RecoveryTargetTime二选一，同时设置时以此参数为准")),
			mcp.WithString("RecoveryTargetTime", mcp.Description("恢复目标时间，与BackupSetId二选一必传，东八区时间")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeCloneDBInstanceSpecRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeCloneDBInstanceSpec(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeBackupDownloadURL - 获取备份下载链接（L4审计，链接泄露风险）
	// 注意：此前缺少必传参数 BackupType（SDK 要求 LogBackup|BaseBackup），
	// 实测会返回 MissingParameter: 请求缺少必传参数 `BackupType`，现补上；
	// 同时补上缺失的 FromJsonString error 检查。
	registerTool(s, cp, g, "DescribeBackupDownloadURL", "获取备份下载链接",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("BackupType", mcp.Required(), mcp.Description("备份类型：LogBackup日志备份|BaseBackup基础备份")),
			mcp.WithString("BackupId", mcp.Required(), mcp.Description("备份ID")),
			mcp.WithNumber("URLExpireTime", mcp.Description("链接有效时间(小时)，取值[0,36]，默认12")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeBackupDownloadURLRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeBackupDownloadURL(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 写操作（3个）=====

	// CreateBaseBackup - 创建基础备份（L2业务确认，消耗IO/存储）
	registerTool(s, cp, g, "CreateBaseBackup", "创建基础备份",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewCreateBaseBackupRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.CreateBaseBackup(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// CloneDBInstance - 克隆实例（L1费用确认）
	// 兼容旧参数别名：SpecName->SpecCode、InstanceName->Name。
	registerTool(s, cp, g, "CloneDBInstance", "克隆实例",
		security.LevelFee,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("源实例ID")),
			mcp.WithString("SpecCode", mcp.Description("售卖规格码")),
			mcp.WithNumber("Storage", mcp.Description("实例容量大小(GB)")),
			mcp.WithNumber("Period", mcp.Description("购买时长(月)")),
			mcp.WithNumber("AutoRenewFlag", mcp.Description("续费标记：0手动，1自动")),
			mcp.WithString("VpcId", mcp.Description("私有网络ID")),
			mcp.WithString("SubnetId", mcp.Description("子网ID")),
			mcp.WithString("Name", mcp.Description("克隆实例名称")),
			mcp.WithString("InstanceChargeType", mcp.Description("计费类型: POSTPAID_BY_HOUR|PREPAID")),
			mcp.WithArray("SecurityGroupIds", mcp.Description("安全组ID列表")),
			mcp.WithArray("DBNodeSet", mcp.Description("实例节点部署信息，每项含 Role(Primary/Standby) 与 Zone")),
			mcp.WithNumber("ProjectId", mcp.Description("项目ID")),
			mcp.WithString("BackupSetId", mcp.Description("基础备份集ID，与 RecoveryTargetTime 二选一")),
			mcp.WithString("RecoveryTargetTime", mcp.Description("恢复时间点，与 BackupSetId 二选一")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			if _, ok := args["SpecCode"]; !ok {
				if legacy, ok := args["SpecName"]; ok {
					args["SpecCode"] = legacy
					delete(args, "SpecName")
				}
			}
			if _, ok := args["Name"]; !ok {
				if legacy, ok := args["InstanceName"]; ok {
					args["Name"] = legacy
					delete(args, "InstanceName")
				}
			}
			delete(args, "Zone")
			req := postgres.NewCloneDBInstanceRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.CloneDBInstance(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Backup tools registered: 8")
}
