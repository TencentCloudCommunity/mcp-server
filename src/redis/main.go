package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-redis"
	mcpsvr := server.NewMCPServer(
		"腾讯云 Redis MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	redisDescribeProductInfo := mcp.NewTool(
		"redis-DescribeProductInfo",
		mcp.WithDescription(`查询产品售卖规格`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
	)
	mcpsvr.AddTool(redisDescribeProductInfo, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeProductInfoRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeProductInfo(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisDescribeInstanceZoneInfo := mcp.NewTool(
		"redis-DescribeInstanceZoneInfo",
		mcp.WithDescription(`查询Redis节点详细信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("指定实例 ID。例如：crs-xjhsdj****。请登录Redis控制台在实例列表复制实例 ID。"),
		),
	)
	mcpsvr.AddTool(redisDescribeInstanceZoneInfo, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeInstanceZoneInfoRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceZoneInfo(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisDescribeInstanceSecurityGroup := mcp.NewTool(
		"redis-DescribeInstanceSecurityGroup",
		mcp.WithDescription(`查询实例安全组信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("实例 ID 列表"),
		),
	)
	mcpsvr.AddTool(redisDescribeInstanceSecurityGroup, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeInstanceSecurityGroupRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceSecurityGroup(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisDescribeProjectSecurityGroup := mcp.NewTool(
		"redis-DescribeProjectSecurityGroup",
		mcp.WithDescription(`查询项目安全组信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"ProjectId",
			mcp.Description("0:默认项目；-1 所有项目; >0: 特定项目"),
		),
		mcp.WithString(
			"SecurityGroupId",
			mcp.Description("安全组Id"),
		),
	)
	mcpsvr.AddTool(redisDescribeProjectSecurityGroup, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeProjectSecurityGroupRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeProjectSecurityGroup(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	// DescribeInstances - 查询Redis实例列表
	redisDescribeInstances := mcp.NewTool(
		"redis-DescribeInstances",
		mcp.WithDescription(`查询Redis实例列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("实例列表的大小，参数默认值20"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，取Limit整数倍"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例Id，如：crs-6ubhgouj"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("枚举范围： projectId,createtime,instancename,type,curDeadline"),
		),
		mcp.WithNumber(
			"OrderType",
			mcp.Description("1倒序，0顺序，默认倒序"),
		),
		mcp.WithArray(
			"VpcIds",
			mcp.Description("私有网络ID数组，数组下标从0开始，如果不传则默认选择基础网络"),
		),
		mcp.WithArray(
			"SubnetIds",
			mcp.Description("子网ID数组，数组下标从0开始"),
		),
		mcp.WithString(
			"SearchKey",
			mcp.Description("模糊搜索关键词，支持实例Id、实例名称、完整IP"),
		),
		mcp.WithArray(
			"ProjectIds",
			mcp.Description("项目 ID 组成的数组。"),
		),
		mcp.WithString(
			"InstanceName",
			mcp.Description("实例名称"),
		),
		mcp.WithArray(
			"UniqVpcIds",
			mcp.Description("私有网络ID数组，数组下标从0开始，如果不传则默认选择基础网络，如：vpc-sad23jfdfk"),
		),
		mcp.WithArray(
			"UniqSubnetIds",
			mcp.Description("子网ID数组，数组下标从0开始，如：subnet-fdj24n34j2"),
		),
		mcp.WithArray(
			"RegionIds",
			mcp.Description("地域ID，已经弃用，可通过公共参数Region查询对应地域"),
		),
		mcp.WithArray(
			"Status",
			mcp.Description("实例状态：0-待初始化，1-流程中，2-运行中，-2-已隔离，-3-待删除"),
		),
		mcp.WithNumber(
			"TypeVersion",
			mcp.Description("类型版本：1-单机版,2-主从版,3-集群版"),
		),
		mcp.WithString(
			"EngineName",
			mcp.Description("引擎信息：Redis-2.8，Redis-4.0，Redis-5.0，Redis-6.0 或者 CKV"),
		),
		mcp.WithArray(
			"AutoRenew",
			mcp.Description("续费模式：0 - 默认状态（手动续费）；1 - 自动续费；2 - 明确不自动续费"),
		),
		mcp.WithString(
			"BillingMode",
			mcp.Description("计费模式：postpaid-按量计费；prepaid-包年包月"),
		),
		mcp.WithNumber(
			"Type",
			mcp.Description("实例类型：1 – Redis老集群版；2 – Redis 2.8主从版；3 – CKV主从版；4 – CKV集群版；5 – Redis 2.8单机版；6 – Redis 4.0主从版；7 – Redis 4.0集群版；8 – Redis 5.0主从版；9 – Redis 5.0集群版；"),
		),
		mcp.WithString(
			"SearchKeys",
			mcp.Description("搜索关键词：支持实例Id、实例名称、完整IP"),
		),
		mcp.WithArray(
			"TypeList",
			mcp.Description("内部参数，用户可以忽略"),
		),
		mcp.WithString(
			"MonitorVersion",
			mcp.Description("监控版本: 1m-分钟粒度监控，5s-5秒粒度监控"),
		),
		mcp.WithArray(
			"InstanceTags",
			mcp.Description("根据标签的Key和Value筛选资源"),
		),
		mcp.WithArray(
			"TagKeys",
			mcp.Description("根据标签的Key筛选资源"),
		),
		mcp.WithArray(
			"ProductVersions",
			mcp.Description("产品版本"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("实例ID"),
		),
		mcp.WithArray(
			"AzMode",
			mcp.Description("多AZ部署的实例"),
		),
	)
	mcpsvr.AddTool(redisDescribeInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeInstancesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstances(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisDescribeSlowLog := mcp.NewTool(
		"redis-DescribeSlowLog",
		mcp.WithDescription(`查询实例慢查询记录`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例Id"),
		),
		mcp.WithString(
			"BeginTime",
			mcp.Description("开始时间"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间"),
		),
		mcp.WithNumber(
			"MinQueryTime",
			mcp.Description("慢查询阈值(毫秒)"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页大小"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，取Limit整数倍"),
		),
		mcp.WithString(
			"Role",
			mcp.Description("节点角色"),
		),
	)
	mcpsvr.AddTool(redisDescribeSlowLog, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeSlowLogRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSlowLog(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisDescribeInstanceMonitorHotKey := mcp.NewTool(
		"redis-DescribeInstanceMonitorHotKey",
		mcp.WithDescription(`查询实例热Key`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例Id"),
		),
		mcp.WithNumber(
			"SpanType",
			mcp.Description("时间范围：1——实时，2——近30分钟，3——近6小时，4——近24小时"),
		),
	)
	mcpsvr.AddTool(redisDescribeInstanceMonitorHotKey, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeInstanceMonitorHotKeyRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceMonitorHotKey(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisDescribeInstanceBackups := mcp.NewTool(
		"redis-DescribeInstanceBackups",
		mcp.WithDescription(`查询Redis实例备份列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("待操作的实例ID，可通过 DescribeInstance 接口返回值中的 InstanceId 获取。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("实例列表的大小，参数默认值20"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，取Limit整数倍"),
		),
		mcp.WithString(
			"BeginTime",
			mcp.Description("开始时间，格式如：2017-02-08 16:46:34。查询实例在 [beginTime, endTime] 时间段内开始备份的备份列表。"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间，格式如：2017-02-08 19:09:26。查询实例在 [beginTime, endTime] 时间段内开始备份的备份列表。"),
		),
		mcp.WithArray(
			"Status",
			mcp.Description("1：备份在流程中，2：备份正常，3：备份转RDB文件处理中，4：已完成RDB转换，-1：备份已过期，-2：备份已删除。"),
		),
		mcp.WithString(
			"InstanceName",
			mcp.Description("实例名称，支持根据实例名称模糊搜索"),
		),
	)
	mcpsvr.AddTool(redisDescribeInstanceBackups, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeInstanceBackupsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceBackups(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisDescribeInstanceParams := mcp.NewTool(
		"redis-DescribeInstanceParams",
		mcp.WithDescription(`查询实例的参数列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例Id"),
		),
	)
	mcpsvr.AddTool(redisDescribeInstanceParams, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeInstanceParamsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceParams(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisDescribeParamTemplates := mcp.NewTool(
		"redis-DescribeParamTemplates",
		mcp.WithDescription(`查询参数模板列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"ProductTypes",
			mcp.Description("产品类型：1 – Redis2.8内存版（集群架构），2 – Redis2.8内存版（标准架构），3 – CKV 3.2内存版(标准架构)，4 – CKV 3.2内存版(集群架构)，5 – Redis2.8内存版（单机），6 – Redis4.0内存版（标准架构），7 – Redis4.0内存版（集群架构），8 – Redis5.0内存版（标准架构），9 – Redis5.0内存版（集群架构）"),
		),
		mcp.WithArray(
			"TemplateNames",
			mcp.Description("模板名称"),
		),
		mcp.WithArray(
			"TemplateIds",
			mcp.Description("模板ID"),
		),
	)
	mcpsvr.AddTool(redisDescribeParamTemplates, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeParamTemplatesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeParamTemplates(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisDescribeParamTemplateInfo := mcp.NewTool(
		"redis-DescribeParamTemplateInfo",
		mcp.WithDescription(`查询参数模板详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"TemplateId",
			mcp.Description("参数模板 ID"),
		),
	)
	mcpsvr.AddTool(redisDescribeParamTemplateInfo, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeParamTemplateInfoRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeParamTemplateInfo(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisDescribeInstanceAccount := mcp.NewTool(
		"redis-DescribeInstanceAccount",
		mcp.WithDescription(`查看实例子账号信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例ID"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页大小"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页偏移量"),
		),
	)
	mcpsvr.AddTool(redisDescribeInstanceAccount, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDescribeInstanceAccountRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceAccount(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	// CreateInstances - 创建Redis实例
	redisCreateInstances := mcp.NewTool(
		"redis-CreateInstances",
		mcp.WithDescription(`创建Redis实例`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"TypeId",
			mcp.Description("实例类型：2 – Redis2.8内存版(标准架构)，3 – CKV 3.2内存版(标准架构)，4 – CKV 3.2内存版(集群架构)，6 – Redis4.0内存版(标准架构)，7 – Redis4.0内存版(集群架构)，8 – Redis5.0内存版(标准架构)，9 – Redis5.0内存版(集群架构)。"),
		),
		mcp.WithNumber(
			"MemSize",
			mcp.Description("内存容量，单位为MB， 数值需为1024的整数倍，具体规格以 查询产品售卖规格 返回的规格为准。"),
		),
		mcp.WithNumber(
			"GoodsNum",
			mcp.Description("实例数量，单次购买实例数量以 查询产品售卖规格 返回的规格为准。"),
		),
		mcp.WithNumber(
			"Period",
			mcp.Description("购买时长，在创建包年包月实例的时候需要填写，按量计费实例填1即可，单位：月，取值范围 [1,2,3,4,5,6,7,8,9,10,11,12,24,36]。"),
		),
		mcp.WithNumber(
			"BillingMode",
			mcp.Description("计费方式:0-按量计费，1-包年包月。"),
		),
		mcp.WithString(
			"ZoneId",
			mcp.Description("可用区ID，可用区ID可通过 查询产品售卖规格 接口查询。"),
		),
		mcp.WithString(
			"Password",
			mcp.Description("实例密码"),
		),
		mcp.WithString(
			"VpcId",
			mcp.Description("私有网络ID，如果不传则默认选择基础网络，请使用私有网络列表 查询。"),
		),
		mcp.WithString(
			"SubnetId",
			mcp.Description("基础网络下， subnetId无效； vpc子网下，取值以查询查询VPC列表值为准。"),
		),
		mcp.WithNumber(
			"ProjectId",
			mcp.Description("项目id，取值以用户账户>用户账户相关接口查询>项目列表返回的projectId为准。"),
		),
		mcp.WithNumber(
			"AutoRenew",
			mcp.Description("自动续费标识。0 - 默认状态（手动续费）；1 - 自动续费；2 - 明确不自动续费。"),
		),
		mcp.WithArray(
			"SecurityGroupIdList",
			mcp.Description("安全组id数组。"),
		),
		mcp.WithString(
			"VPort",
			mcp.Description("用户自定义的端口 不填则默认为6379，范围[1024,65535]。"),
		),
		mcp.WithNumber(
			"RedisShardNum",
			mcp.Description("实例分片数量，购买标准版实例不需要填写，集群版分片数量范围[3,5,8,12,16,24,32,64,96,128]。"),
		),
		mcp.WithNumber(
			"RedisReplicasNum",
			mcp.Description("实例副本数量，Redis 2.8标准版、CKV标准版只支持1副本，4.0、5.0标准版可选择1-5副本，集群版可选择0-5副本。"),
		),
		mcp.WithNumber(
			"ReplicasReadonly",
			mcp.Description("是否支持副本只读，Redis 2.8标准版、CKV标准版不支持副本只读，开启副本只读，实例将自动读写分离，写请求路由到主节点，读请求路由到副本节点，如需开启副本只读建议副本数>=2。"),
		),
		mcp.WithString(
			"InstanceName",
			mcp.Description("实例名称"),
		),
		mcp.WithBoolean(
			"NoAuth",
			mcp.Description("是否支持免密，true-免密实例，false-非免密实例，默认为非免密实例，当前仅VPC网络的实例支持免密码访问。"),
		),
		mcp.WithArray(
			"NodeSet",
			mcp.Description("实例的节点信息，目前支持传入节点的类型（主节点或者副本节点），节点的可用区。单可用区部署不需要传递此参数。"),
		),
		mcp.WithArray(
			"ResourceTags",
			mcp.Description("购买实例绑定标签"),
		),
		mcp.WithString(
			"ZoneName",
			mcp.Description("指定实例所属的可用区名称。"),
		),
		mcp.WithString(
			"TemplateId",
			mcp.Description("指定实例所属的参数模板id。参数模板id可通过DescribeParamTemplates接口返回值获取。"),
		),
		mcp.WithBoolean(
			"DryRun",
			mcp.Description("内部参数，标识创建实例是否需要检查。"),
		),
		mcp.WithString(
			"ProductVersion",
			mcp.Description("指实例部署模式。- local：传统架构，默认为 local。- cdc：独享集群。- cloud：云原生，当前已暂停售卖。"),
		),
		mcp.WithString(
			"RedisClusterId",
			mcp.Description("独享集群 ID。当ProductVersion设置为cdc时，该参数必须设置。"),
		),
		mcp.WithArray(
			"AlarmPolicyList",
			mcp.Description("指定实例需要应用的告警策略id。告警策略id可通过DescribeAlarmPolicy接口返回值获取。请注意：HyperMemcached内存版实例和CKV实例暂不支持绑定告警策略。"),
		),
	)
	mcpsvr.AddTool(redisCreateInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewCreateInstancesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateInstances(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisModifyInstanceParams := mcp.NewTool(
		"redis-ModifyInstanceParams",
		mcp.WithDescription(`修改实例参数`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例Id"),
		),
		mcp.WithArray(
			"InstanceParams",
			mcp.Description("实例修改的参数列表"),
		),
	)
	mcpsvr.AddTool(redisModifyInstanceParams, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewModifyInstanceParamsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ModifyInstanceParams(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisCreateInstanceAccount := mcp.NewTool(
		"redis-CreateInstanceAccount",
		mcp.WithDescription(`创建实例子账号`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例ID"),
		),
		mcp.WithString(
			"AccountName",
			mcp.Description("子账号名称"),
		),
		mcp.WithString(
			"AccountPassword",
			mcp.Description("子账号密码"),
		),
		mcp.WithString(
			"ReadonlyPolicy",
			mcp.Description("指定账号的读请求路由分发至主节点或副本节点。未开启副本只读，不支持选择副本节点。- master：主节点 - replication：副本节点"),
		),
		mcp.WithString(
			"Privilege",
			mcp.Description("读写策略：填写r、w、rw，表示只读，只写，读写策略；不填默认为rw"),
		),
		mcp.WithString(
			"Remark",
			mcp.Description("子账号描述信息"),
		),
	)
	mcpsvr.AddTool(redisCreateInstanceAccount, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewCreateInstanceAccountRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateInstanceAccount(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisAssociateSecurityGroups := mcp.NewTool(
		"redis-AssociateSecurityGroups",
		mcp.WithDescription(`绑定安全组`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Product",
			mcp.Description("数据库引擎名称：mariadb,cdb,cynosdb,dcdb,redis,mongodb 等。"),
		),
		mcp.WithString(
			"SecurityGroupId",
			mcp.Description("要绑定的安全组ID，类似sg-efil73jd。"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("被绑定的实例ID，类似ins-lesecurk，支持指定多个实例。"),
		),
	)
	mcpsvr.AddTool(redisAssociateSecurityGroups, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewAssociateSecurityGroupsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.AssociateSecurityGroups(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	redisDisassociateSecurityGroups := mcp.NewTool(
		"redis-DisassociateSecurityGroups",
		mcp.WithDescription(`安全组批量解绑云资源`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Product",
			mcp.Description("数据库引擎名称：mariadb,cdb,cynosdb,dcdb,redis,mongodb 等。"),
		),
		mcp.WithString(
			"SecurityGroupId",
			mcp.Description("安全组Id。"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("实例ID列表，一个或者多个实例Id组成的数组。"),
		),
	)
	mcpsvr.AddTool(redisDisassociateSecurityGroups, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := redis.NewClient(credential, region_, cpf)
		req := redis.NewDisassociateSecurityGroupsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DisassociateSecurityGroups(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	sseEndpoint := getEnv("MCP_SERVER_SSE_ENDPOINT", "/sse")
	messageEndpoint := getEnv("MCP_SERVER_MESSAGE_ENDPOINT", "/message")
	ssePort := getEnv("MCP_SERVER_SSE_PORT", "9000")

	sseServer := server.NewSSEServer(mcpsvr,
		server.WithSSEEndpoint(sseEndpoint),
		server.WithMessageEndpoint(messageEndpoint),
		server.WithAppendQueryToMessageEndpoint())

	log.Printf("SSE server listening on :" + ssePort)
	serverURL := fmt.Sprintf("http://127.0.0.1:%s%s", ssePort, sseEndpoint)
	outputMCPServerConfig(mcpServerName, serverURL)
	if err := sseServer.Start(":" + ssePort); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func outputMCPServerConfig(mcpServerName, serverURL string) {
	config := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			mcpServerName: map[string]interface{}{
				"url":  serverURL,
				"type": "sse",
			},
		},
	}

	jsonOutput, _ := json.MarshalIndent(config, "", " ")
	fmt.Println("=== MCP Server Configuration ===")
	fmt.Println("Copy the following configuration to your MCP client:")
	fmt.Println()
	fmt.Println(string(jsonOutput))
	fmt.Println()
	fmt.Println("The server is now ready to accept connections.")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
