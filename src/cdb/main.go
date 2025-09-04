package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-cdb"
	mcpsvr := server.NewMCPServer(
		"腾讯云 CDB MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	cdbDescribeDBInstances := mcp.NewTool(
		"cdb-DescribeDBInstances",
		mcp.WithDescription(`查询实例列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"ProjectId",
			mcp.Description("项目 ID。"),
		),
		mcp.WithArray(
			"InstanceTypes",
			mcp.Description("实例类型，可取值：1 - 主实例，2 - 灾备实例，3 - 只读实例。"),
		),
		mcp.WithArray(
			"Vips",
			mcp.Description("实例的内网 IP 地址。"),
		),
		mcp.WithArray(
			"Status",
			mcp.Description("实例状态，可取值：<br>0 - 创建中<br>1 - 运行中<br>4 - 正在进行隔离操作<br>5 - 已隔离（可在回收站恢复开机）"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，默认值为 0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("单次请求返回的数量，默认值为 20，最大值为 2000。"),
		),
		mcp.WithNumber(
			"SecurityGroupId",
			mcp.Description("安全组 ID。当使用安全组 ID 为过滤条件时，需指定 WithSecurityGroup 参数为 1。"),
		),
		mcp.WithArray(
			"PayTypes",
			mcp.Description("付费类型，可取值：0 - 包年包月，1 - 小时计费。"),
		),
		mcp.WithArray(
			"InstanceNames",
			mcp.Description("实例名称。"),
		),
		mcp.WithArray(
			"TaskStatus",
			mcp.Description("实例任务状态，可能取值：<br>0 - 没有任务<br>1 - 升级中<br>2 - 数据导入中<br>3 - 开放Slave中<br>4 - 外网访问开通中<br>5 - 批量操作执行中<br>6 - 回档中<br>7 - 外网访问关闭中<br>8 - 密码修改中<br>9 - 实例名修改中<br>10 - 重启中<br>12 - 自建迁移中<br>13 - 删除库表中<br>14 - 灾备实例创建同步中<br>15 - 升级待切换<br>16 - 升级切换中<br>17 - 升级切换完成<br>19 - 参数设置待执行<br>34 - 原地升级待执行"),
		),
		mcp.WithArray(
			"EngineVersions",
			mcp.Description("实例数据库引擎版本，可能取值：5.1、5.5、5.6 和 5.7。"),
		),
		mcp.WithArray(
			"VpcIds",
			mcp.Description("私有网络的 ID。"),
		),
		mcp.WithArray(
			"ZoneIds",
			mcp.Description("可用区的 ID。"),
		),
		mcp.WithArray(
			"SubnetIds",
			mcp.Description("子网 ID。"),
		),
		mcp.WithArray(
			"CdbErrors",
			mcp.Description("是否锁定标记，可选值：0 - 不锁定，1 - 锁定，默认为0。"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("返回结果集排序的字段，目前支持：&quot;InstanceId&quot;，&quot;InstanceName&quot;，&quot;CreateTime&quot;，&quot;DeadlineTime&quot;。"),
		),
		mcp.WithString(
			"OrderDirection",
			mcp.Description("返回结果集排序方式，目前支持：&quot;ASC&quot; 或者 &quot;DESC&quot;。"),
		),
		mcp.WithString(
			"WithSecurityGroup",
			mcp.Description("是否以安全组 ID 为过滤条件。"),
		),
		mcp.WithNumber(
			"WithExCluster",
			mcp.Description("是否包含独享集群详细信息，可取值：0 - 不包含，1 - 包含。"),
		),
		mcp.WithString(
			"ExClusterId",
			mcp.Description("独享集群 ID。"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("实例 ID。"),
		),
		mcp.WithNumber(
			"InitFlag",
			mcp.Description("初始化标记，可取值：0 - 未初始化，1 - 初始化。"),
		),
		mcp.WithNumber(
			"WithDr",
			mcp.Description("是否包含灾备关系对应的实例，可取值：0 - 不包含，1 - 包含。默认取值为1。如果拉取主实例，则灾备关系的数据在DrInfo字段中， 如果拉取灾备实例， 则灾备关系的数据在MasterInfo字段中。灾备关系中只包含部分基本的数据，详细的数据需要自行调接口拉取。"),
		),
		mcp.WithNumber(
			"WithRo",
			mcp.Description("是否包含只读实例，可取值：0 - 不包含，1 - 包含。默认取值为1。"),
		),
		mcp.WithNumber(
			"WithMaster",
			mcp.Description("是否包含主实例，可取值：0 - 不包含，1 - 包含。默认取值为1。"),
		),
		mcp.WithArray(
			"DeployGroupIds",
			mcp.Description("置放群组ID列表。"),
		),
		mcp.WithArray(
			"TagKeysForSearch",
			mcp.Description("是否以标签键为过滤条件。"),
		),
		mcp.WithArray(
			"CageIds",
			mcp.Description("金融围拢 ID 。"),
		),
		mcp.WithArray(
			"TagValues",
			mcp.Description("标签值"),
		),
		mcp.WithArray(
			"UniqueVpcIds",
			mcp.Description("私有网络字符型vpcId"),
		),
		mcp.WithArray(
			"UniqSubnetIds",
			mcp.Description("私有网络字符型subnetId"),
		),
		mcp.WithArray(
			"Tags",
			mcp.Description("标签键值"),
		),
		mcp.WithArray(
			"ProxyVips",
			mcp.Description("数据库代理 IP 。"),
		),
		mcp.WithArray(
			"ProxyIds",
			mcp.Description("数据库代理 ID 。"),
		),
		mcp.WithArray(
			"EngineTypes",
			mcp.Description("数据库引擎类型。"),
		),
		mcp.WithString(
			"QueryClusterInfo",
			mcp.Description("是否获取集群版实例节点信息，可填：true或false"),
		),
	)
	mcpsvr.AddTool(cdbDescribeDBInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdb.NewClient(credential, region_, cpf)
		req := cdb.NewDescribeDBInstancesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDBInstances(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cdbDescribeDBInstanceInfo := mcp.NewTool(
		"cdb-DescribeDBInstanceInfo",
		mcp.WithDescription(`查询实例基本信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例 ID。"),
		),
	)
	mcpsvr.AddTool(cdbDescribeDBInstanceInfo, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdb.NewClient(credential, region_, cpf)
		req := cdb.NewDescribeDBInstanceInfoRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDBInstanceInfo(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cdbDescribeInstanceUpgradeType := mcp.NewTool(
		"cdb-DescribeInstanceUpgradeType",
		mcp.WithDescription(`查询数据库实例升级类型`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例 ID。"),
		),
		mcp.WithString(
			"DstCpu",
			mcp.Description("目标实例 CPU 的核数。"),
		),
		mcp.WithString(
			"DstMemory",
			mcp.Description("目标实例内存大小，单位：MB。"),
		),
		mcp.WithString(
			"DstDisk",
			mcp.Description("目标实例磁盘大小，单位：GB。"),
		),
		mcp.WithString(
			"DstVersion",
			mcp.Description("目标实例数据库版本。"),
		),
		mcp.WithString(
			"DstDeployMode",
			mcp.Description("目标实例部署模型。"),
		),
		mcp.WithString(
			"DstProtectMode",
			mcp.Description("目标实例复制类型。"),
		),
		mcp.WithString(
			"DstSlaveZone",
			mcp.Description("目标实例备机1可用区。"),
		),
		mcp.WithString(
			"DstBackupZone",
			mcp.Description("目标实例备机2可用区。"),
		),
		mcp.WithString(
			"DstCdbType",
			mcp.Description("目标实例类型。"),
		),
		mcp.WithString(
			"DstZoneId",
			mcp.Description("目标实例主可用区。"),
		),
		mcp.WithString(
			"NodeDistribution",
			mcp.Description("独享集群 CDB 实例的节点分布情况。"),
		),
		mcp.WithString(
			"ClusterTopology",
			mcp.Description("集群版的节点拓扑配置"),
		),
	)
	mcpsvr.AddTool(cdbDescribeInstanceUpgradeType, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdb.NewClient(credential, region_, cpf)
		req := cdb.NewDescribeInstanceUpgradeTypeRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceUpgradeType(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cdbDescribeAccounts := mcp.NewTool(
		"cdb-DescribeAccounts",
		mcp.WithDescription(`查询云数据库的所有账号信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例 ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例 ID 相同。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("记录偏移量，默认值为0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("单次请求返回的数量，默认值为20，最小值为1，最大值为100。"),
		),
		mcp.WithString(
			"AccountRegexp",
			mcp.Description("匹配账号名的正则表达式，规则同 MySQL 官网。"),
		),
	)
	mcpsvr.AddTool(cdbDescribeAccounts, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdb.NewClient(credential, region_, cpf)
		req := cdb.NewDescribeAccountsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeAccounts(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cdbDescribeDatabases := mcp.NewTool(
		"cdb-DescribeDatabases",
		mcp.WithDescription(`查询数据库`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例 ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例 ID 相同。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，最小值为0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("单次请求数量，默认值为20，最小值为1，最大值为100。"),
		),
		mcp.WithString(
			"DatabaseRegexp",
			mcp.Description("匹配数据库库名的正则表达式。"),
		),
	)
	mcpsvr.AddTool(cdbDescribeDatabases, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdb.NewClient(credential, region_, cpf)
		req := cdb.NewDescribeDatabasesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDatabases(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	cdbDescribeInstanceParams := mcp.NewTool(
		"cdb-DescribeInstanceParams",
		mcp.WithDescription(`查询实例的可设置参数列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例 ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例 ID 相同，可使用 [查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值。"),
		),
	)
	mcpsvr.AddTool(cdbDescribeInstanceParams, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdb.NewClient(credential, region_, cpf)
		req := cdb.NewDescribeInstanceParamsRequest()
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

	cdbDescribeParamTemplates := mcp.NewTool(
		"cdb-DescribeParamTemplates",
		mcp.WithDescription(`查询参数模板列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"EngineVersions",
			mcp.Description("引擎版本，缺省则查询所有"),
		),
		mcp.WithArray(
			"EngineTypes",
			mcp.Description("引擎类型，缺省则查询所有"),
		),
		mcp.WithArray(
			"TemplateNames",
			mcp.Description("模板名称，缺省则查询所有"),
		),
		mcp.WithArray(
			"TemplateIds",
			mcp.Description("模板id，缺省则查询所有"),
		),
	)
	mcpsvr.AddTool(cdbDescribeParamTemplates, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdb.NewClient(credential, region_, cpf)
		req := cdb.NewDescribeParamTemplatesRequest()
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

	cdbDescribeParamTemplateInfo := mcp.NewTool(
		"cdb-DescribeParamTemplateInfo",
		mcp.WithDescription(`查询参数模板详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"TemplateId",
			mcp.Description("参数模板 ID。"),
		),
	)
	mcpsvr.AddTool(cdbDescribeParamTemplateInfo, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdb.NewClient(credential, region_, cpf)
		req := cdb.NewDescribeParamTemplateInfoRequest()
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

	cdbCreateDatabase := mcp.NewTool(
		"cdb-CreateDatabase",
		mcp.WithDescription(`创建数据库`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例 ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例 ID 相同。"),
		),
		mcp.WithString(
			"DBName",
			mcp.Description("数据库名称，长度不超过64。"),
		),
		mcp.WithString(
			"CharacterSetName",
			mcp.Description("字符集，可选值：utf8，gbk，latin1，utf8mb4。"),
		),
	)
	mcpsvr.AddTool(cdbCreateDatabase, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := cdb.NewClient(credential, region_, cpf)
		req := cdb.NewCreateDatabaseRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateDatabase(req)
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
