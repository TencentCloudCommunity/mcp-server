package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-tsf"
	mcpsvr := server.NewMCPServer(
		"腾讯云 TSF MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	tsfDescribeApplications := mcp.NewTool(
		"tsf-DescribeApplications",
		mcp.WithDescription(`获取应用列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("搜索字段"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序字段"),
		),
		mcp.WithNumber(
			"OrderType",
			mcp.Description("排序类型"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页个数"),
		),
		mcp.WithString(
			"ApplicationType",
			mcp.Description("应用类型"),
		),
		mcp.WithString(
			"MicroserviceType",
			mcp.Description("应用的微服务类型"),
		),
		mcp.WithArray(
			"ApplicationResourceTypeList",
			mcp.Description("应用资源类型数组"),
		),
		mcp.WithArray(
			"ApplicationIdList",
			mcp.Description("IdList"),
		),
		mcp.WithArray(
			"MicroserviceTypeList",
			mcp.Description("查询多种微服务类型的应用"),
		),
	)
	mcpsvr.AddTool(tsfDescribeApplications, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribeApplicationsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeApplications(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfDescribeApplicationAttribute := mcp.NewTool(
		"tsf-DescribeApplicationAttribute",
		mcp.WithDescription(`获取应用列表其它字段`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ApplicationId",
			mcp.Description("应用ID"),
		),
	)
	mcpsvr.AddTool(tsfDescribeApplicationAttribute, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribeApplicationAttributeRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeApplicationAttribute(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfDescribeSimpleClusters := mcp.NewTool(
		"tsf-DescribeSimpleClusters",
		mcp.WithDescription(`查询简单集群列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"ClusterIdList",
			mcp.Description("需要查询的集群ID列表，不填或不传入时查询所有内容"),
		),
		mcp.WithString(
			"ClusterType",
			mcp.Description("需要查询的集群类型，不填或不传入时查询所有内容"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询偏移量，默认为0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页个数，默认为20， 取值应为1~50"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("对id和name进行关键词过滤"),
		),
		mcp.WithBoolean(
			"DisableProgramAuthCheck",
			mcp.Description("是否关闭鉴权"),
		),
	)
	mcpsvr.AddTool(tsfDescribeSimpleClusters, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribeSimpleClustersRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSimpleClusters(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfDescribeSimpleNamespaces := mcp.NewTool(
		"tsf-DescribeSimpleNamespaces",
		mcp.WithDescription(`查询简单命名空间列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithArray(
			"NamespaceIdList",
			mcp.Description("命名空间ID列表，不传入时查询全量"),
		),
		mcp.WithString(
			"ClusterId",
			mcp.Description("集群ID，不传入时查询全量"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("每页条数"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("起始偏移量"),
		),
		mcp.WithString(
			"NamespaceId",
			mcp.Description("命名空间ID，不传入时查询全量"),
		),
		mcp.WithArray(
			"NamespaceResourceTypeList",
			mcp.Description("查询资源类型列表"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("通过id和name进行过滤"),
		),
		mcp.WithArray(
			"NamespaceTypeList",
			mcp.Description("查询的命名空间类型列表"),
		),
		mcp.WithString(
			"NamespaceName",
			mcp.Description("通过命名空间名精确过滤"),
		),
		mcp.WithString(
			"IsDefault",
			mcp.Description("通过是否是默认命名空间过滤，不传表示拉取全部命名空间。0：默认命名空间。1：非默认命名空间"),
		),
		mcp.WithBoolean(
			"DisableProgramAuthCheck",
			mcp.Description("是否关闭鉴权查询"),
		),
	)
	mcpsvr.AddTool(tsfDescribeSimpleNamespaces, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribeSimpleNamespacesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeSimpleNamespaces(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfDescribeImageRepository := mcp.NewTool(
		"tsf-DescribeImageRepository",
		mcp.WithDescription(`查询镜像仓库列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("仓库名，搜索关键字,不带命名空间的"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，取值从0开始"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页个数，默认为20， 取值应为1~100"),
		),
		mcp.WithString(
			"RepoType",
			mcp.Description("企业: tcr ；个人: personal或者不填"),
		),
		mcp.WithString(
			"ApplicationId",
			mcp.Description("应用id"),
		),
		mcp.WithObject(
			"TcrRepoInfo",
			mcp.Description("TcrRepoInfo值"),
		),
		mcp.WithString(
			"RepoName",
			mcp.Description("镜像仓库名称"),
		),
	)
	mcpsvr.AddTool(tsfDescribeImageRepository, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribeImageRepositoryRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeImageRepository(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfDescribeImageTags := mcp.NewTool(
		"tsf-DescribeImageTags",
		mcp.WithDescription(`查询镜像版本列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ApplicationId",
			mcp.Description("应用Id"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，取值从0开始"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页个数，默认为20， 取值应为1~100"),
		),
		mcp.WithNumber(
			"QueryImageIdFlag",
			mcp.Description("不填和0:查询 1:不查询"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("可用于搜索的 tag 名字"),
		),
		mcp.WithString(
			"RepoType",
			mcp.Description("企业: tcr ；个人: personal或者不填"),
		),
		mcp.WithObject(
			"TcrRepoInfo",
			mcp.Description("TcrRepoInfo值"),
		),
		mcp.WithString(
			"RepoName",
			mcp.Description("仓库名"),
		),
	)
	mcpsvr.AddTool(tsfDescribeImageTags, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribeImageTagsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeImageTags(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfDescribeRepositories := mcp.NewTool(
		"tsf-DescribeRepositories",
		mcp.WithDescription(`查询仓库列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("查询关键字（按照仓库名称搜索）"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始偏移"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量限制"),
		),
		mcp.WithString(
			"RepositoryType",
			mcp.Description("仓库类型（默认仓库：default，私有仓库：private）"),
		),
	)
	mcpsvr.AddTool(tsfDescribeRepositories, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribeRepositoriesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeRepositories(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfDescribePkgs := mcp.NewTool(
		"tsf-DescribePkgs",
		mcp.WithDescription(`获取某个应用的程序包信息列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ApplicationId",
			mcp.Description("应用ID（只传入应用ID，返回该应用下所有软件包信息）"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("查询关键字（支持根据包ID，包名，包版本号搜索）"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序关键字（默认为&quot;UploadTime&quot;：上传时间）"),
		),
		mcp.WithString(
			"OrderType",
			mcp.Description("升序：0/降序：1（默认降序）"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始偏移"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量限制"),
		),
		mcp.WithString(
			"RepositoryType",
			mcp.Description("程序包仓库类型"),
		),
		mcp.WithString(
			"RepositoryId",
			mcp.Description("程序包仓库id"),
		),
		mcp.WithArray(
			"PackageTypeList",
			mcp.Description("程序包类型数组支持（fatjar jar war tar.gz zip）"),
		),
	)
	mcpsvr.AddTool(tsfDescribePkgs, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribePkgsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribePkgs(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfDescribeGroups := mcp.NewTool(
		"tsf-DescribeGroups",
		mcp.WithDescription(`获取虚拟机部署组列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("搜索字段"),
		),
		mcp.WithString(
			"ApplicationId",
			mcp.Description("应用ID"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序字段"),
		),
		mcp.WithNumber(
			"OrderType",
			mcp.Description("排序方式"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页个数"),
		),
		mcp.WithString(
			"NamespaceId",
			mcp.Description("命名空间ID"),
		),
		mcp.WithString(
			"ClusterId",
			mcp.Description("集群ID"),
		),
		mcp.WithArray(
			"GroupResourceTypeList",
			mcp.Description("部署组资源类型列表"),
		),
		mcp.WithString(
			"Status",
			mcp.Description("部署组状态过滤字段"),
		),
		mcp.WithArray(
			"GroupIdList",
			mcp.Description("无"),
		),
	)
	mcpsvr.AddTool(tsfDescribeGroups, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribeGroupsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeGroups(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfDescribeBusinessLogConfigs := mcp.NewTool(
		"tsf-DescribeBusinessLogConfigs",
		mcp.WithDescription(`查询日志配置项列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，取值范围大于等于0，默认值为0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("单页请求配置数量，取值范围[1, 50]，默认值为10"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("模糊匹配关键词"),
		),
		mcp.WithBoolean(
			"DisableProgramAuthCheck",
			mcp.Description("无"),
		),
		mcp.WithArray(
			"ConfigIdList",
			mcp.Description("无"),
		),
	)
	mcpsvr.AddTool(tsfDescribeBusinessLogConfigs, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribeBusinessLogConfigsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeBusinessLogConfigs(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfSearchBusinessLog := mcp.NewTool(
		"tsf-SearchBusinessLog",
		mcp.WithDescription(`业务日志搜索`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ConfigId",
			mcp.Description("日志配置项ID"),
		),
		mcp.WithArray(
			"InstanceIds",
			mcp.Description("机器实例ID，不传表示全部实例"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("开始时间"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("结束时间"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("请求偏移量，取值范围大于等于0，默认值为0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("单页请求配置数量，取值范围[1, 200]，默认值为50"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序规则，默认值&quot;time&quot;"),
		),
		mcp.WithString(
			"OrderType",
			mcp.Description("排序方式，取值&quot;asc&quot;或&quot;desc&quot;，默认值&quot;desc&quot;"),
		),
		mcp.WithArray(
			"SearchWords",
			mcp.Description("检索关键词"),
		),
		mcp.WithArray(
			"GroupIds",
			mcp.Description("部署组ID列表，不传表示全部部署组"),
		),
		mcp.WithString(
			"SearchWordType",
			mcp.Description("检索类型，取值&quot;LUCENE&quot;, &quot;REGEXP&quot;, &quot;NORMAL&quot;"),
		),
		mcp.WithString(
			"BatchType",
			mcp.Description("批量请求类型，取值&quot;page&quot;或&quot;scroll&quot;"),
		),
		mcp.WithString(
			"ScrollId",
			mcp.Description("游标ID"),
		),
		mcp.WithArray(
			"SearchAfter",
			mcp.Description("查询es使用searchAfter时，游标"),
		),
	)
	mcpsvr.AddTool(tsfSearchBusinessLog, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewSearchBusinessLogRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.SearchBusinessLog(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfSearchStdoutLog := mcp.NewTool(
		"tsf-SearchStdoutLog",
		mcp.WithDescription(`标准输出日志搜索`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("机器实例ID， 和  实例 ID 二者必选其一，不能同时为空"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("单页请求配置数量，取值范围[1, 500]，默认值为100"),
		),
		mcp.WithArray(
			"SearchWords",
			mcp.Description("检索关键词"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("查询起始时间"),
		),
		mcp.WithString(
			"GroupId",
			mcp.Description("部署组ID，和 InstanceId 二者必选其一，不能同时为空"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("查询结束时间"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("请求偏移量，取值范围大于等于0，默认值为"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("排序规则，默认值&quot;time&quot;"),
		),
		mcp.WithString(
			"OrderType",
			mcp.Description("排序方式，取值&quot;asc&quot;或&quot;desc&quot;，默认"),
		),
		mcp.WithString(
			"SearchWordType",
			mcp.Description("检索类型，取值&quot;LUCENE&quot;, &quot;REGEXP&quot;,"),
		),
		mcp.WithString(
			"BatchType",
			mcp.Description("批量请求类型，取值&quot;page&quot;或&quot;scroll&quot;，默认"),
		),
		mcp.WithString(
			"ScrollId",
			mcp.Description("游标ID"),
		),
		mcp.WithArray(
			"SearchAfter",
			mcp.Description("查询es使用searchAfter时，游标"),
		),
	)
	mcpsvr.AddTool(tsfSearchStdoutLog, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewSearchStdoutLogRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.SearchStdoutLog(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfDescribeInvocationMetricDataCurve := mcp.NewTool(
		"tsf-DescribeInvocationMetricDataCurve",
		mcp.WithDescription(`查询调用指标数据变化曲线`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"StartTime",
			mcp.Description("查询开始时间"),
		),
		mcp.WithString(
			"EndTime",
			mcp.Description("查询结束时间"),
		),
		mcp.WithNumber(
			"Period",
			mcp.Description("查询时间粒度，单位秒可选值：60、3600、86400"),
		),
		mcp.WithArray(
			"MetricDimensions",
			mcp.Description("查询指标维度，不能为空，支持 ServiceName, OperationName, PeerServiceName, PeerOperationName"),
		),
		mcp.WithArray(
			"Metrics",
			mcp.Description("查询指标名，不能为空."),
		),
		mcp.WithString(
			"Kind",
			mcp.Description("视图视角。可选值：SERVER, CLIENT。默认为SERVER"),
		),
		mcp.WithString(
			"Type",
			mcp.Description("类型。组件监控使用，可选值：SQL 或者 NoSQL"),
		),
	)
	mcpsvr.AddTool(tsfDescribeInvocationMetricDataCurve, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribeInvocationMetricDataCurveRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInvocationMetricDataCurve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	tsfDescribePrograms := mcp.NewTool(
		"tsf-DescribePrograms",
		mcp.WithDescription(`查询数据集列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("模糊查询数据集ID，数据集名称，不传入时查询全量"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("每页数量"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("起始偏移量"),
		),
	)
	mcpsvr.AddTool(tsfDescribePrograms, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := tsf.NewClient(credential, region_, cpf)
		req := tsf.NewDescribeProgramsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribePrograms(req)
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
