package tools

import (
	"fmt"
	"sort"

	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

type requestParser interface {
	FromJsonString(string) error
}

type AlignmentResult struct {
	Name  string
	OK    bool
	Error string
}

type alignmentCase struct {
	Name       string
	GuardLevel security.GuardLevel
	Args       map[string]interface{}
	Validate   func(map[string]interface{}) error
}

func ValidateAllOpenAPIArgumentAlignment() ([]AlignmentResult, error) {
	cases := openAPIAlignmentCases()
	if len(cases) != 48 {
		return nil, fmt.Errorf("expected 48 tool alignment cases, got %d", len(cases))
	}

	seen := make(map[string]struct{}, len(cases))
	results := make([]AlignmentResult, 0, len(cases))
	failed := 0
	for _, tc := range cases {
		if _, ok := seen[tc.Name]; ok {
			return nil, fmt.Errorf("duplicate tool alignment case: %s", tc.Name)
		}
		seen[tc.Name] = struct{}{}

		args := cloneArgs(tc.Args)
		args["region"] = "ap-chengdu"
		if security.NeedsConfirm(tc.GuardLevel) {
			args["confirm"] = true
		}
		err := tc.Validate(args)
		result := AlignmentResult{Name: tc.Name, OK: err == nil}
		if err != nil {
			result.Error = err.Error()
			failed++
		}
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})
	if failed > 0 {
		return results, fmt.Errorf("%d/%d tools failed OpenAPI argument alignment", failed, len(results))
	}
	return results, nil
}

func parseRequest(args map[string]interface{}, req requestParser, normalizers ...func(map[string]interface{})) error {
	payload := stripControlArgs(args)
	for _, normalize := range normalizers {
		if normalize != nil {
			normalize(payload)
		}
	}
	return req.FromJsonString(marshalArgs(payload))
}

func openAPIAlignmentCases() []alignmentCase {
	return []alignmentCase{
		{Name: "DescribeDBInstanceAttribute", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeDBInstanceAttributeRequest())
		}},
		{Name: "UpgradeDBInstanceKernelVersion", GuardLevel: security.LevelBusiness, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "TargetDBKernelVersion": "v18.4_r1.9", "SwitchTag": 0, "DryRun": true}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewUpgradeDBInstanceKernelVersionRequest())
		}},
		{Name: "DescribeDBInstances", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "Limit": 20, "Offset": 0}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeDBInstancesRequest(), normalizeDescribeDBInstancesArgs)
		}},
		{Name: "DescribeClasses", GuardLevel: security.LevelNone, Args: map[string]interface{}{"Zone": "ap-chengdu-1", "DBEngine": "postgresql", "DBMajorVersion": "18"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeClassesRequest())
		}},
		{Name: "DescribeDBVersions", GuardLevel: security.LevelNone, Args: map[string]interface{}{}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeDBVersionsRequest(), normalizeDescribeDBVersionsArgs)
		}},
		{Name: "DescribeRegions", GuardLevel: security.LevelNone, Args: map[string]interface{}{}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeRegionsRequest())
		}},
		{Name: "DescribeZones", GuardLevel: security.LevelNone, Args: map[string]interface{}{}, Validate: func(args map[string]interface{}) error { return parseRequest(args, postgres.NewDescribeZonesRequest()) }},
		{Name: "DescribeProductConfig", GuardLevel: security.LevelNone, Args: map[string]interface{}{"Zone": "ap-chengdu-1", "DBEngine": "postgresql"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeProductConfigRequest())
		}},
		{Name: "DescribeTasks", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "Limit": 20, "Offset": 0}, Validate: func(args map[string]interface{}) error { return parseRequest(args, postgres.NewDescribeTasksRequest()) }},
		{Name: "CreateInstances", GuardLevel: security.LevelFee, Args: map[string]interface{}{"Zone": "ap-chengdu-1", "InstanceSpec": "pg.it.small2", "Volume": 100, "InstanceCount": 1, "Period": 1, "DBCharset": "UTF8", "AdminName": "mcpadmin", "AdminPassword": "McpCheck_2026!", "DBMajorVersion": "18", "InstanceChargeType": "POSTPAID_BY_HOUR", "VpcId": "vpc-test", "SubnetId": "subnet-test", "InstanceName": "mcp-create-check", "SecurityGroupIds": []interface{}{"sg-test"}}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewCreateInstancesRequest(), normalizeCreateInstancesArgs)
		}},
		{Name: "ModifyDBInstanceName", GuardLevel: security.LevelAudit, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "InstanceName": "mcp-renamed"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewModifyDBInstanceNameRequest())
		}},
		{Name: "ModifyDBInstanceSpec", GuardLevel: security.LevelFee, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "Memory": 4, "Volume": 120, "Cpu": 2, "SwitchTag": 0, "InstanceType": "legacy-ignore"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewModifyDBInstanceSpecRequest(), normalizeModifyDBInstanceSpecArgs)
		}},
		{Name: "RestartDBInstance", GuardLevel: security.LevelBusiness, Args: map[string]interface{}{"DBInstanceId": "postgres-test"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewRestartDBInstanceRequest())
		}},
		{Name: "IsolateDBInstances", GuardLevel: security.LevelBusiness, Args: map[string]interface{}{"DBInstanceId": "postgres-test"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewIsolateDBInstancesRequest(), normalizeDBInstanceIdSetArgs)
		}},
		{Name: "DisIsolateDBInstances", GuardLevel: security.LevelAudit, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "Period": 1, "AutoVoucher": true, "VoucherIds": []interface{}{"voucher-test"}}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDisIsolateDBInstancesRequest(), normalizeDBInstanceIdSetArgs)
		}},

		{Name: "DescribeAccounts", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "Limit": 20, "Offset": 0, "OrderBy": "createTime", "OrderByType": "desc"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeAccountsRequest())
		}},
		{Name: "CreateAccount", GuardLevel: security.LevelAudit, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "UserName": "mcp_user", "Password": "McpCheck_2026!", "Type": "normal", "Remark": "contract-check"}, Validate: func(args map[string]interface{}) error { return parseRequest(args, postgres.NewCreateAccountRequest()) }},
		{Name: "DeleteAccount", GuardLevel: security.LevelBusiness, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "UserName": "mcp_user"}, Validate: func(args map[string]interface{}) error { return parseRequest(args, postgres.NewDeleteAccountRequest()) }},
		{Name: "ModifyAccountPrivileges", GuardLevel: security.LevelAudit, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "UserName": "mcp_user", "ModifyPrivilegeSet": []interface{}{map[string]interface{}{"ModifyType": "grantObject", "IsCascade": false, "DatabasePrivilege": map[string]interface{}{"Object": map[string]interface{}{"ObjectType": "database", "ObjectName": "postgres", "DatabaseName": "postgres"}, "PrivilegeSet": []interface{}{"CREATE"}}}}}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewModifyAccountPrivilegesRequest())
		}},
		{Name: "ResetAccountPassword", GuardLevel: security.LevelAudit, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "UserName": "mcp_user", "Password": "McpCheck_2026_v2!"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewResetAccountPasswordRequest())
		}},
		{Name: "DescribeAccountPrivileges", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "UserName": "mcp_user", "DatabaseObjectSet": []interface{}{map[string]interface{}{"ObjectType": "database", "ObjectName": "postgres", "DatabaseName": "postgres"}}}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeAccountPrivilegesRequest())
		}},

		{Name: "DescribeDatabases", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "Filters": []interface{}{map[string]interface{}{"Name": "database-name", "Values": []interface{}{"postgres"}}}, "Limit": 20, "Offset": 0}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeDatabasesRequest())
		}},
		{Name: "CreateDatabase", GuardLevel: security.LevelBusiness, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "DatabaseName": "mcp_db", "DatabaseOwner": "postgres", "Encoding": "UTF8", "Collate": "C", "Ctype": "C"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewCreateDatabaseRequest())
		}},
		{Name: "ModifyDatabaseOwner", GuardLevel: security.LevelAudit, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "DatabaseName": "mcp_db", "DatabaseOwner": "owner_test"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewModifyDatabaseOwnerRequest())
		}},
		{Name: "DescribeDatabaseObjects", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "ObjectType": "schema", "DatabaseName": "postgres", "Limit": 20, "Offset": 0}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeDatabaseObjectsRequest())
		}},

		{Name: "DescribeDBInstanceParameters", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "ParamName": "max_connections"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeDBInstanceParametersRequest())
		}},
		{Name: "DescribeParameterTemplates", GuardLevel: security.LevelNone, Args: map[string]interface{}{"Filters": []interface{}{map[string]interface{}{"Name": "DBMajorVersion", "Values": []interface{}{"18"}}}, "Limit": 20, "Offset": 0, "OrderBy": "CreateTime", "OrderByType": "desc"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeParameterTemplatesRequest())
		}},
		{Name: "DescribeParameterTemplateAttributes", GuardLevel: security.LevelNone, Args: map[string]interface{}{"TemplateId": "tpl-test"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeParameterTemplateAttributesRequest())
		}},
		{Name: "DescribeParamsEvent", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeParamsEventRequest())
		}},
		{Name: "ModifyDBInstanceParameters", GuardLevel: security.LevelCritical, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "ParamList": []interface{}{map[string]interface{}{"Name": "log_min_duration_statement", "Value": "1000"}}}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewModifyDBInstanceParametersRequest(), normalizeModifyDBInstanceParametersArgs)
		}},

		{Name: "DescribeBackupOverview", GuardLevel: security.LevelNone, Args: map[string]interface{}{}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeBackupOverviewRequest(), normalizeDescribeBackupOverviewArgs)
		}},
		{Name: "DescribeBaseBackups", GuardLevel: security.LevelNone, Args: map[string]interface{}{"MinFinishTime": "2026-07-01 00:00:00", "MaxFinishTime": "2026-07-06 23:59:59", "Filters": []interface{}{map[string]interface{}{"Name": "db-instance-id", "Values": []interface{}{"postgres-test"}}}, "Limit": 20, "Offset": 0, "OrderBy": "FinishTime", "OrderByType": "desc"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeBaseBackupsRequest(), normalizeBackupListArgs)
		}},
		{Name: "DescribeLogBackups", GuardLevel: security.LevelNone, Args: map[string]interface{}{"MinFinishTime": "2026-07-01 00:00:00", "MaxFinishTime": "2026-07-06 23:59:59", "Filters": []interface{}{map[string]interface{}{"Name": "db-instance-id", "Values": []interface{}{"postgres-test"}}}, "Limit": 20, "Offset": 0, "OrderBy": "FinishTime", "OrderByType": "desc"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeLogBackupsRequest(), normalizeBackupListArgs)
		}},
		{Name: "DescribeAvailableRecoveryTime", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeAvailableRecoveryTimeRequest())
		}},
		{Name: "DescribeCloneDBInstanceSpec", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "BackupSetId": "backup-test"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeCloneDBInstanceSpecRequest())
		}},
		{Name: "DescribeBackupDownloadURL", GuardLevel: security.LevelAudit, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "BackupType": "BaseBackup", "BackupId": "backup-test", "URLExpireTime": 2}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeBackupDownloadURLRequest())
		}},
		{Name: "CreateBaseBackup", GuardLevel: security.LevelBusiness, Args: map[string]interface{}{"DBInstanceId": "postgres-test"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewCreateBaseBackupRequest())
		}},
		{Name: "CloneDBInstance", GuardLevel: security.LevelFee, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "SpecName": "pg.it.small2", "Storage": 100, "Period": 1, "AutoRenewFlag": 0, "VpcId": "vpc-test", "SubnetId": "subnet-test", "InstanceName": "mcp-clone-check", "InstanceChargeType": "POSTPAID_BY_HOUR", "SecurityGroupIds": []interface{}{"sg-test"}, "ProjectId": 0, "BackupSetId": "backup-test", "DBNodeSet": []interface{}{map[string]interface{}{"Role": "Primary", "Zone": "ap-chengdu-1"}, map[string]interface{}{"Role": "Standby", "Zone": "ap-chengdu-1"}}, "Zone": "ap-chengdu-1"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewCloneDBInstanceRequest(), normalizeCloneDBInstanceArgs)
		}},

		{Name: "DescribeSlowQueryList", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "StartTime": "2026-07-06 00:00:00", "EndTime": "2026-07-06 23:59:59", "Limit": 20, "Offset": 0}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeSlowQueryListRequest())
		}},
		{Name: "DescribeSlowQueryAnalysis", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "StartTime": "2026-07-06 00:00:00", "EndTime": "2026-07-06 23:59:59", "Limit": 20, "Offset": 0}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeSlowQueryAnalysisRequest())
		}},
		{Name: "DescribeDBErrlogs", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "StartTime": "2026-07-06 00:00:00", "EndTime": "2026-07-06 23:59:59", "Limit": 20, "Offset": 0}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeDBErrlogsRequest())
		}},

		{Name: "OpenDBExtranetAccess", GuardLevel: security.LevelBusiness, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "IsIpv6": 0}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewOpenDBExtranetAccessRequest())
		}},
		{Name: "CloseDBExtranetAccess", GuardLevel: security.LevelBusiness, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "IsIpv6": 0}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewCloseDBExtranetAccessRequest())
		}},
		{Name: "DescribeDBInstanceSecurityGroups", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeDBInstanceSecurityGroupsRequest())
		}},
		{Name: "ModifyDBInstanceSecurityGroups", GuardLevel: security.LevelBusiness, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "SecurityGroupIds": []interface{}{"sg-test"}}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewModifyDBInstanceSecurityGroupsRequest(), normalizeModifyDBInstanceSecurityGroupsArgs)
		}},

		{Name: "DescribeReadOnlyGroups", GuardLevel: security.LevelNone, Args: map[string]interface{}{"Filters": []interface{}{map[string]interface{}{"Name": "db-master-instance-id", "Values": []interface{}{"postgres-test"}}}, "PageSize": 20, "PageNumber": 1, "OrderBy": "CreateTime", "OrderByType": "asc"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeReadOnlyGroupsRequest(), normalizeDescribeReadOnlyGroupsArgs)
		}},
		{Name: "CreateReadOnlyDBInstance", GuardLevel: security.LevelFee, Args: map[string]interface{}{"DBInstanceId": "postgres-test", "SpecName": "pg.it.small2", "Storage": 100, "InstanceCount": 1, "Period": 1, "Zone": "ap-chengdu-1", "VpcId": "vpc-test", "SubnetId": "subnet-test", "InstanceChargeType": "POSTPAID_BY_HOUR", "AutoRenewFlag": 0, "InstanceName": "mcp-ro-check", "ReadOnlyGroupId": "rogrp-test", "SecurityGroupIds": []interface{}{"sg-test"}}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewCreateReadOnlyDBInstanceRequest(), normalizeCreateReadOnlyDBInstanceArgs)
		}},

		{Name: "DescribeDBInstanceSSLConfig", GuardLevel: security.LevelNone, Args: map[string]interface{}{"DBInstanceId": "postgres-test"}, Validate: func(args map[string]interface{}) error {
			return parseRequest(args, postgres.NewDescribeDBInstanceSSLConfigRequest())
		}},
	}
}
