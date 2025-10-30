package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-es"
	mcpsvr := server.NewMCPServer(
		"腾讯云 ES MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	esDescribeInstanceLogs := mcp.NewTool(
		"es-DescribeInstanceLogs",
		mcp.WithDescription(`查询ES集群日志`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群实例ID"),
		),
		mcp.WithNumber(
			"LogType",
			mcp.Description("日志类型，默认值为1"),
		),
		mcp.WithString(
			"SearchKey",
			mcp.Description("搜索词，支持LUCENE语法，如 level:WARN、ip:1.1.1.1、message:test-index等"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("日志开始时间，格式为YYYY-MM-DD HH:MM:SS, 如2019-01-22 20:15:53"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("日志结束时间，格式为YYYY-MM-DD HH:MM:SS, 如2019-01-22 20:15:53"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页起始值, 默认值为0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页大小，默认值为100，最大值100"),
		),
		mcp.WithNumber(
			"OrderByType",
			mcp.Description("时间排序方式，默认值为0"),
		),
		mcp.WithArray(
			"LogLevels",
			mcp.Description("日志级别"),
		),
		mcp.WithArray(
			"NodeIds",
			mcp.Description("节点ID"),
		),
		mcp.WithString(
			"IndexName",
			mcp.Description("慢日志索引名"),
		),
		mcp.WithNumber(
			"QueryCost",
			mcp.Description("慢日志查询耗时"),
		),
	)
	mcpsvr.AddTool(esDescribeInstanceLogs, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeInstanceLogsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceLogs(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeInstancePluginList := mcp.NewTool(
		"es-DescribeInstancePluginList",
		mcp.WithDescription(`查询实例插件列表`),
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
			"Offset",
			mcp.Description("分页起始值, 默认值0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页大小，默认值10"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序字段<li>1：插件名 pluginName</li>"),
		),
		mcp.WithString(
			"OrderByType",
			mcp.Description("排序方式<li>0：升序 asc</li><li>1：降序 desc</li>"),
		),
		mcp.WithNumber(
			"PluginType",
			mcp.Description("0：系统插件"),
		),
	)
	mcpsvr.AddTool(esDescribeInstancePluginList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeInstancePluginListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstancePluginList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esGetRequestTargetNodeTypes := mcp.NewTool(
		"es-GetRequestTargetNodeTypes",
		mcp.WithDescription(`获取接收客户端请求的节点类型`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例ID"),
		),
	)
	mcpsvr.AddTool(esGetRequestTargetNodeTypes, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewGetRequestTargetNodeTypesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.GetRequestTargetNodeTypes(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esInquirePriceRenewInstance := mcp.NewTool(
		"es-InquirePriceRenewInstance",
		mcp.WithDescription(`集群续费询价`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群实例Id"),
		),
	)
	mcpsvr.AddTool(esInquirePriceRenewInstance, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewInquirePriceRenewInstanceRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.InquirePriceRenewInstance(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeIndexList := mcp.NewTool(
		"es-DescribeIndexList",
		mcp.WithDescription(`获取索引列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"IndexType",
			mcp.Description("索引类型。auto：自治索引；normal：普通索引"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ES集群ID"),
		),
		mcp.WithString(
			"IndexName",
			mcp.Description("索引名，若填空则获取所有索引"),
		),
		mcp.WithString(
			"Username",
			mcp.Description("集群访问用户名"),
		),
		mcp.WithString(
			"Password",
			mcp.Description("集群访问密码"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("一页展示数量"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序字段，支持索引名：IndexName、索引存储量：IndexStorage、索引创建时间：IndexCreateTime"),
		),
		mcp.WithArray(
			"IndexStatusList",
			mcp.Description("过滤索引状态"),
		),
		mcp.WithString(
			"Order",
			mcp.Description("排序顺序，支持asc、desc，默认为desc 数据格式&quot;asc&quot;,&quot;desc&quot;"),
		),
	)
	mcpsvr.AddTool(esDescribeIndexList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeIndexListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeIndexList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeIndexMeta := mcp.NewTool(
		"es-DescribeIndexMeta",
		mcp.WithDescription(`获取索引元数据`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ES集群ID"),
		),
		mcp.WithString(
			"IndexType",
			mcp.Description("索引类型。auto：自治索引；normal：普通索引"),
		),
		mcp.WithString(
			"IndexName",
			mcp.Description("索引名，若填空则获取所有索引"),
		),
		mcp.WithString(
			"Username",
			mcp.Description("集群访问用户名"),
		),
		mcp.WithString(
			"Password",
			mcp.Description("集群访问密码"),
		),
	)
	mcpsvr.AddTool(esDescribeIndexMeta, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeIndexMetaRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeIndexMeta(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeInstanceOperations := mcp.NewTool(
		"es-DescribeInstanceOperations",
		mcp.WithDescription(`查询实例操作记录`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群实例ID"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("起始时间, e.g. &quot;2019-03-07 16:30:39&quot;"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间, e.g. &quot;2019-03-30 20:18:03&quot;"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页起始值"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页大小"),
		),
	)
	mcpsvr.AddTool(esDescribeInstanceOperations, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeInstanceOperationsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceOperations(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeInstances := mcp.NewTool(
		"es-DescribeInstances",
		mcp.WithDescription(`查询ES集群实例`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Zone",
			mcp.Description("集群实例所属可用区，不传则默认所有可用区"),
		),
		mcp.WithString(
			"InstanceIds",
			mcp.Description("集群实例ID列表"),
		),
		mcp.WithString(
			"InstanceNames",
			mcp.Description("集群实例名称列表"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页起始值, 默认值0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页大小，默认值20"),
		),
		mcp.WithNumber(
			"OrderByKey",
			mcp.Description("排序字段<li>1：实例ID</li><li>2：实例名称</li><li>3：可用区</li><li>4：创建时间</li>若orderByKey未传递则按创建时间降序排序"),
		),
		mcp.WithNumber(
			"OrderByType",
			mcp.Description("排序方式<li>0：升序</li><li>1：降序</li>若传递了orderByKey未传递orderByType, 则默认升序"),
		),
		mcp.WithArray(
			"TagList",
			mcp.Description("节点标签信息列表"),
		),
		mcp.WithArray(
			"IpList",
			mcp.Description("私有网络vip列表"),
		),
		mcp.WithArray(
			"ZoneList",
			mcp.Description("可用区列表"),
		),
		mcp.WithArray(
			"HealthStatus",
			mcp.Description("健康状态筛列表:0表示绿色，1表示黄色，2表示红色,-1表示未知"),
		),
		mcp.WithArray(
			"VpcIds",
			mcp.Description("Vpc列表 筛选项"),
		),
		mcp.WithString(
			"CdcId",
			mcp.Description("cdc集群id"),
		),
	)
	mcpsvr.AddTool(esDescribeInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeInstancesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstances(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeViews := mcp.NewTool(
		"es-DescribeViews",
		mcp.WithDescription(`查询集群视图`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群实例ID"),
		),
	)
	mcpsvr.AddTool(esDescribeViews, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeViewsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeViews(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeClusterSnapshot := mcp.NewTool(
		"es-DescribeClusterSnapshot",
		mcp.WithDescription(`获取快照备份列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群实例Id，格式：es-xxxx"),
		),
		mcp.WithString(
			"RepositoryName",
			mcp.Description("快照仓库名称"),
		),
		mcp.WithString(
			"SnapshotName",
			mcp.Description("集群快照名称"),
		),
	)
	mcpsvr.AddTool(esDescribeClusterSnapshot, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeClusterSnapshotRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeClusterSnapshot(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeServerlessSpaceUser := mcp.NewTool(
		"es-DescribeServerlessSpaceUser",
		mcp.WithDescription(`查看Serverless空间子用户`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"SpaceId",
			mcp.Description("空间的ID"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("游标"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("页条数"),
		),
		mcp.WithArray(
			"UserNames",
			mcp.Description("用户名列表过滤"),
		),
		mcp.WithArray(
			"UserTypes",
			mcp.Description("用户类型"),
		),
		mcp.WithArray(
			"PrivilegeTypes",
			mcp.Description("权限类型"),
		),
	)
	mcpsvr.AddTool(esDescribeServerlessSpaceUser, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeServerlessSpaceUserRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeServerlessSpaceUser(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeServerlessInstances := mcp.NewTool(
		"es-DescribeServerlessInstances",
		mcp.WithDescription(`Serverless获取索引列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("索引集群ID"),
		),
		mcp.WithArray(
			"IndexNames",
			mcp.Description("索引名"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("一页展示数量"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序字段，支持索引名：IndexName、索引存储量：IndexStorage、索引创建时间：IndexCreateTime"),
		),
		mcp.WithArray(
			"IndexStatusList",
			mcp.Description("过滤索引状态"),
		),
		mcp.WithString(
			"Order",
			mcp.Description("排序顺序，支持asc、desc，默认为desc"),
		),
		mcp.WithArray(
			"SpaceIds",
			mcp.Description("索引空间ID列表"),
		),
		mcp.WithArray(
			"DiSourceTypes",
			mcp.Description("数据链路数据源类型"),
		),
		mcp.WithArray(
			"TagList",
			mcp.Description("标签信息"),
		),
	)
	mcpsvr.AddTool(esDescribeServerlessInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeServerlessInstancesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeServerlessInstances(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeServerlessSpaces := mcp.NewTool(
		"es-DescribeServerlessSpaces",
		mcp.WithDescription(`获取Serverless索引空间列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"SpaceIds",
			mcp.Description("过滤的空间ID"),
		),
		mcp.WithArray(
			"SpaceNames",
			mcp.Description("过滤的空间名"),
		),
		mcp.WithString(
			"Order",
			mcp.Description("排序顺序，支持升序asc、降序desc"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序字段，支持空间创建时间SpaceCreateTime"),
		),
		mcp.WithArray(
			"VpcIds",
			mcp.Description("vpcId信息数组"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页起始"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页条数"),
		),
		mcp.WithArray(
			"TagList",
			mcp.Description("标签信息"),
		),
	)
	mcpsvr.AddTool(esDescribeServerlessSpaces, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeServerlessSpacesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeServerlessSpaces(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeServerlessMetrics := mcp.NewTool(
		"es-DescribeServerlessMetrics",
		mcp.WithDescription(`获取实例对应的监控指标`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"SpaceId",
			mcp.Description("space空间id"),
		),
		mcp.WithString(
			"IndexId",
			mcp.Description("index索引id"),
		),
		mcp.WithArray(
			"MetricType",
			mcp.Description("指标类型，暂时只支持Storage(存储大小),AllMetric(所有存储指标：索引流量、存储大小、文档数量、读请求和写请求)"),
		),
		mcp.WithNumber(
			"DurationType",
			mcp.Description("时间长度类型DurationType(1: 3小时, 2: 昨天1天,3: 今日0点到现在)"),
		),
		mcp.WithArray(
			"BatchIndexList",
			mcp.Description("索引数据"),
		),
	)
	mcpsvr.AddTool(esDescribeServerlessMetrics, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeServerlessMetricsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeServerlessMetrics(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeUserCosSnapshotList := mcp.NewTool(
		"es-DescribeUserCosSnapshotList",
		mcp.WithDescription(`查询快照信息接口`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"CosBucket",
			mcp.Description("cos桶名"),
		),
		mcp.WithString(
			"BasePath",
			mcp.Description("bucket 桶下的备份路径"),
		),
		mcp.WithString(
			"ClusterInstanceId",
			mcp.Description("云上集群迁移集群名"),
		),
	)
	mcpsvr.AddTool(esDescribeUserCosSnapshotList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeUserCosSnapshotListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeUserCosSnapshotList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeLogstashInstanceLogs := mcp.NewTool(
		"es-DescribeLogstashInstanceLogs",
		mcp.WithDescription(`查询Logstash实例日志`),
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
			"LogType",
			mcp.Description("日志类型，默认值为1"),
		),
		mcp.WithString(
			"SearchKey",
			mcp.Description("搜索词，支持LUCENE语法，如 level:WARN、ip:1.1.1.1、message:test-index等"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("日志开始时间，格式为YYYY-MM-DD HH:MM:SS, 如2019-01-22 20:15:53"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("日志结束时间，格式为YYYY-MM-DD HH:MM:SS, 如2019-01-22 20:15:53"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页起始值, 默认值为0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页大小，默认值为100，最大值100"),
		),
		mcp.WithNumber(
			"OrderByType",
			mcp.Description("时间排序方式，默认值为0"),
		),
	)
	mcpsvr.AddTool(esDescribeLogstashInstanceLogs, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeLogstashInstanceLogsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeLogstashInstanceLogs(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeLogstashInstances := mcp.NewTool(
		"es-DescribeLogstashInstances",
		mcp.WithDescription(`获取Logstash实例列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Zone",
			mcp.Description("实例所属可用区，不传则默认所有可用区"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("实例ID列表"),
		),
		mcp.WithString(
			"InstanceNames",
			mcp.Description("实例名称列表"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页起始值, 默认值0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页大小，默认值20"),
		),
		mcp.WithNumber(
			"OrderByKey",
			mcp.Description("排序字段<li>1：实例ID</li><li>2：实例名称</li><li>3：可用区</li><li>4：创建时间</li>若orderKey未传递则按创建时间降序排序"),
		),
		mcp.WithNumber(
			"OrderByType",
			mcp.Description("排序方式<li>0：升序</li><li>1：降序</li>若传递了orderByKey未传递orderByType, 则默认升序"),
		),
		mcp.WithArray(
			"VpcIds",
			mcp.Description("VpcId 筛选项"),
		),
		mcp.WithArray(
			"TagList",
			mcp.Description("标签信息列表"),
		),
	)
	mcpsvr.AddTool(esDescribeLogstashInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeLogstashInstancesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeLogstashInstances(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeLogstashInstanceOperations := mcp.NewTool(
		"es-DescribeLogstashInstanceOperations",
		mcp.WithDescription(`查询Logstash实例操作记录`),
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
			"StartTime",
			mcp.Description("起始时间, e.g. &quot;2019-03-07 16:30:39&quot;"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间, e.g. &quot;2019-03-30 20:18:03&quot;"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页起始值"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页大小"),
		),
	)
	mcpsvr.AddTool(esDescribeLogstashInstanceOperations, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeLogstashInstanceOperationsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeLogstashInstanceOperations(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	esDescribeLogstashPipelines := mcp.NewTool(
		"es-DescribeLogstashPipelines",
		mcp.WithDescription(`获取Logstash实例管道列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例ID"),
		),
	)
	mcpsvr.AddTool(esDescribeLogstashPipelines, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := es.NewClient(credential, region_, cpf)
		req := es.NewDescribeLogstashPipelinesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeLogstashPipelines(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
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
