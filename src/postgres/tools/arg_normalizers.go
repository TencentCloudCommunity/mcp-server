package tools

func cloneArgs(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return map[string]interface{}{}
	}
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func stripControlArgs(src map[string]interface{}) map[string]interface{} {
	args := cloneArgs(src)
	delete(args, "region")
	delete(args, "confirm")
	return args
}

func normalizeDescribeBackupOverviewArgs(args map[string]interface{}) {
	delete(args, "DBInstanceId")
}

func normalizeDescribeDBVersionsArgs(args map[string]interface{}) {
	delete(args, "DBEngine")
}

func normalizeDBInstanceFilterArgs(args map[string]interface{}, filterName string) {
	if id, ok := args["DBInstanceId"].(string); ok && id != "" {
		args["Filters"] = []map[string]interface{}{
			{"Name": filterName, "Values": []string{id}},
		}
		delete(args, "DBInstanceId")
	}
}

func normalizeDescribeDBInstancesArgs(args map[string]interface{}) {
	normalizeDBInstanceFilterArgs(args, "db-instance-id")
}

func normalizeBackupListArgs(args map[string]interface{}) {
	if _, ok := args["MinFinishTime"]; !ok {
		if legacy, ok := args["StartTime"]; ok {
			args["MinFinishTime"] = legacy
		}
	}
	if _, ok := args["MaxFinishTime"]; !ok {
		if legacy, ok := args["EndTime"]; ok {
			args["MaxFinishTime"] = legacy
		}
	}
	delete(args, "StartTime")
	delete(args, "EndTime")
	normalizeDBInstanceFilterArgs(args, "db-instance-id")
}

func normalizeDescribeReadOnlyGroupsArgs(args map[string]interface{}) {
	if _, ok := args["Filters"]; !ok {
		filters := make([]map[string]interface{}, 0, 2)
		if id, ok := args["DBInstanceId"].(string); ok && id != "" {
			filters = append(filters, map[string]interface{}{"Name": "db-master-instance-id", "Values": []string{id}})
		}
		if groupID, ok := args["ReadOnlyGroupId"].(string); ok && groupID != "" {
			filters = append(filters, map[string]interface{}{"Name": "read-only-group-id", "Values": []string{groupID}})
		}
		if len(filters) > 0 {
			args["Filters"] = filters
		}
	}
	if _, ok := args["PageSize"]; !ok {
		if limit, ok := int64Value(args["Limit"]); ok && limit > 0 {
			args["PageSize"] = limit
		}
	}
	if _, ok := args["PageNumber"]; !ok {
		if offset, ok := int64Value(args["Offset"]); ok {
			pageNumber := int64(1)
			if pageSize, ok := int64Value(args["PageSize"]); ok && pageSize > 0 {
				pageNumber = offset/pageSize + 1
			}
			args["PageNumber"] = pageNumber
		}
	}
	delete(args, "DBInstanceId")
	delete(args, "ReadOnlyGroupId")
	delete(args, "Limit")
	delete(args, "Offset")
}

func normalizeCreateInstancesArgs(args map[string]interface{}) {
	if _, ok := args["SpecCode"]; !ok {
		if legacy, ok := args["InstanceSpec"]; ok {
			args["SpecCode"] = legacy
			delete(args, "InstanceSpec")
		}
	}
	if _, ok := args["Storage"]; !ok {
		if legacy, ok := args["Volume"]; ok {
			args["Storage"] = legacy
			delete(args, "Volume")
		}
	}
	if _, ok := args["Charset"]; !ok {
		if legacy, ok := args["DBCharset"]; ok {
			args["Charset"] = legacy
			delete(args, "DBCharset")
		}
	}
	if _, ok := args["Name"]; !ok {
		if legacy, ok := args["InstanceName"]; ok {
			args["Name"] = legacy
			delete(args, "InstanceName")
		}
	}
}

func normalizeModifyDBInstanceSpecArgs(args map[string]interface{}) {
	if _, ok := args["Storage"]; !ok {
		if legacy, ok := args["Volume"]; ok {
			args["Storage"] = legacy
			delete(args, "Volume")
		}
	}
	delete(args, "InstanceType")
}

func normalizeDBInstanceIdSetArgs(args map[string]interface{}) {
	if _, ok := args["DBInstanceIdSet"]; !ok {
		if id, ok := args["DBInstanceId"].(string); ok && id != "" {
			args["DBInstanceIdSet"] = []string{id}
		}
	}
	delete(args, "DBInstanceId")
}

func normalizeModifyDBInstanceParametersArgs(args map[string]interface{}) {
	switch rawList := args["ParamList"].(type) {
	case []interface{}:
		for _, item := range rawList {
			entry, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			normalizeParameterEntry(entry)
		}
	case []map[string]interface{}:
		for _, entry := range rawList {
			normalizeParameterEntry(entry)
		}
	}
}

func normalizeParameterEntry(entry map[string]interface{}) {
	if _, hasExpected := entry["ExpectedValue"]; !hasExpected {
		if v, hasValue := entry["Value"]; hasValue {
			entry["ExpectedValue"] = v
			delete(entry, "Value")
		}
	}
}

func normalizeModifyDBInstanceSecurityGroupsArgs(args map[string]interface{}) {
	if _, ok := args["SecurityGroupIdSet"]; !ok {
		if legacy, ok := args["SecurityGroupIds"]; ok {
			args["SecurityGroupIdSet"] = legacy
			delete(args, "SecurityGroupIds")
		}
	}
}

func normalizeCloneDBInstanceArgs(args map[string]interface{}) {
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
}

func normalizeCreateReadOnlyDBInstanceArgs(args map[string]interface{}) {
	if _, ok := args["MasterDBInstanceId"]; !ok {
		if legacy, ok := args["DBInstanceId"]; ok {
			args["MasterDBInstanceId"] = legacy
		}
	}
	if _, ok := args["SpecCode"]; !ok {
		if legacy, ok := args["SpecName"]; ok {
			args["SpecCode"] = legacy
		}
	}
	if _, ok := args["Name"]; !ok {
		if legacy, ok := args["InstanceName"]; ok {
			args["Name"] = legacy
		}
	}
	delete(args, "DBInstanceId")
	delete(args, "SpecName")
	delete(args, "InstanceName")
}

func int64Value(v interface{}) (int64, bool) {
	switch x := v.(type) {
	case int:
		return int64(x), true
	case int8:
		return int64(x), true
	case int16:
		return int64(x), true
	case int32:
		return int64(x), true
	case int64:
		return x, true
	case uint:
		return int64(x), true
	case uint8:
		return int64(x), true
	case uint16:
		return int64(x), true
	case uint32:
		return int64(x), true
	case uint64:
		return int64(x), true
	case float64:
		return int64(x), true
	default:
		return 0, false
	}
}
