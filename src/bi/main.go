package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-bi"
	mcpsvr := server.NewMCPServer(
		"腾讯云 BI MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	biCreateProject := mcp.NewTool(
		"bi-CreateProject",
		mcp.WithDescription(`创建项目`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Name",
			mcp.Description("项目名称"),
		),
		mcp.WithString(
			"ColorCode",
			mcp.Description("logo底色"),
		),
		mcp.WithString(
			"Logo",
			mcp.Description("项目Logo"),
		),
		mcp.WithString(
			"Mark",
			mcp.Description("备注"),
		),
		mcp.WithBoolean(
			"IsApply",
			mcp.Description("是否允许用户申请"),
		),
		mcp.WithNumber(
			"DefaultPanelType",
			mcp.Description("默认看板"),
		),
		mcp.WithString(
			"ManagePlatform",
			mcp.Description("管理平台"),
		),
	)
	mcpsvr.AddTool(biCreateProject, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewCreateProjectRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateProject(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biDeleteProject := mcp.NewTool(
		"bi-DeleteProject",
		mcp.WithDescription(`删除项目`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Id",
			mcp.Description("项目ID"),
		),
		mcp.WithString(
			"Seed",
			mcp.Description("随机数"),
		),
		mcp.WithNumber(
			"DefaultPanelType",
			mcp.Description("默认看板"),
		),
	)
	mcpsvr.AddTool(biDeleteProject, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewDeleteProjectRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DeleteProject(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biDescribeProjectInfo := mcp.NewTool(
		"bi-DescribeProjectInfo",
		mcp.WithDescription(`项目详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Id",
			mcp.Description("项目Id"),
		),
		mcp.WithNumber(
			"DefaultPanelType",
			mcp.Description("默认看板"),
		),
	)
	mcpsvr.AddTool(biDescribeProjectInfo, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewDescribeProjectInfoRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeProjectInfo(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biDescribeProjectList := mcp.NewTool(
		"bi-DescribeProjectList",
		mcp.WithDescription(`项目列表数据接口`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"PageSize",
			mcp.Description("页容，初版默认20，将来可能根据屏幕宽度动态变化"),
		),
		mcp.WithNumber(
			"PageNo",
			mcp.Description("页标"),
		),
		mcp.WithString(
			"Keyword",
			mcp.Description("检索模糊字段"),
		),
		mcp.WithBoolean(
			"AllPage",
			mcp.Description("是否全部展示，如果是ture，则忽略分页"),
		),
		mcp.WithString(
			"ModuleCollection",
			mcp.Description("角色信息"),
		),
		mcp.WithArray(
			"ModuleIdList",
			mcp.Description("moduleId集合"),
		),
	)
	mcpsvr.AddTool(biDescribeProjectList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewDescribeProjectListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeProjectList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biModifyProject := mcp.NewTool(
		"bi-ModifyProject",
		mcp.WithDescription(`修改项目`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Id",
			mcp.Description("项目Id"),
		),
		mcp.WithString(
			"Name",
			mcp.Description("名字"),
		),
		mcp.WithString(
			"ColorCode",
			mcp.Description("颜色值"),
		),
		mcp.WithString(
			"Logo",
			mcp.Description("图标"),
		),
		mcp.WithString(
			"Mark",
			mcp.Description("备注"),
		),
		mcp.WithBoolean(
			"IsApply",
			mcp.Description("可申请"),
		),
		mcp.WithString(
			"Seed",
			mcp.Description("种子"),
		),
		mcp.WithNumber(
			"DefaultPanelType",
			mcp.Description("默认看板"),
		),
		mcp.WithString(
			"PanelScope",
			mcp.Description("2"),
		),
		mcp.WithString(
			"ManagePlatform",
			mcp.Description("项目管理平台"),
		),
	)
	mcpsvr.AddTool(biModifyProject, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewModifyProjectRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ModifyProject(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biCreateDatasource := mcp.NewTool(
		"bi-CreateDatasource",
		mcp.WithDescription(`创建数据源`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DbHost",
			mcp.Description("HOST"),
		),
		mcp.WithNumber(
			"DbPort",
			mcp.Description("端口"),
		),
		mcp.WithString(
			"ServiceType",
			mcp.Description("后端提供字典：域类型，1、腾讯云，2、本地"),
		),
		mcp.WithString(
			"DbType",
			mcp.Description("驱动"),
		),
		mcp.WithString(
			"Charset",
			mcp.Description("数据库编码"),
		),
		mcp.WithString(
			"DbUser",
			mcp.Description("用户名"),
		),
		mcp.WithString(
			"DbPwd",
			mcp.Description("密码"),
		),
		mcp.WithString(
			"DbName",
			mcp.Description("数据库名称"),
		),
		mcp.WithString(
			"SourceName",
			mcp.Description("数据库别名"),
		),
		mcp.WithNumber(
			"ProjectId",
			mcp.Description("项目id"),
		),
		mcp.WithString(
			"Catalog",
			mcp.Description("catalog值"),
		),
		mcp.WithString(
			"DataOrigin",
			mcp.Description("第三方数据源标识"),
		),
		mcp.WithString(
			"DataOriginProjectId",
			mcp.Description("第三方项目id"),
		),
		mcp.WithString(
			"DataOriginDatasourceId",
			mcp.Description("第三方数据源id"),
		),
		mcp.WithString(
			"ExtraParam",
			mcp.Description("扩展参数"),
		),
		mcp.WithString(
			"UniqVpcId",
			mcp.Description("腾讯云私有网络统一标识"),
		),
		mcp.WithString(
			"Vip",
			mcp.Description("私有网络ip"),
		),
		mcp.WithString(
			"Vport",
			mcp.Description("私有网络端口"),
		),
		mcp.WithString(
			"VpcId",
			mcp.Description("腾讯云私有网络标识"),
		),
		mcp.WithArray(
			"OperationAuthLimit",
			mcp.Description("操作权限限制"),
		),
		mcp.WithBoolean(
			"UseVPC",
			mcp.Description("开启vpc"),
		),
		mcp.WithString(
			"RegionId",
			mcp.Description("地域"),
		),
	)
	mcpsvr.AddTool(biCreateDatasource, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewCreateDatasourceRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateDatasource(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biCreateDatasourceCloud := mcp.NewTool(
		"bi-CreateDatasourceCloud",
		mcp.WithDescription(`创建云数据库`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ServiceType",
			mcp.Description("后端提供字典：域类型，1、腾讯云，2、本地"),
		),
		mcp.WithString(
			"DbType",
			mcp.Description("驱动"),
		),
		mcp.WithString(
			"Charset",
			mcp.Description("数据库编码"),
		),
		mcp.WithString(
			"DbUser",
			mcp.Description("用户名"),
		),
		mcp.WithString(
			"DbPwd",
			mcp.Description("密码"),
		),
		mcp.WithString(
			"DbName",
			mcp.Description("数据库名称"),
		),
		mcp.WithString(
			"SourceName",
			mcp.Description("数据库别名"),
		),
		mcp.WithString(
			"ProjectId",
			mcp.Description("项目ID"),
		),
		mcp.WithString(
			"Vip",
			mcp.Description("公有云内网ip"),
		),
		mcp.WithString(
			"Vport",
			mcp.Description("公有云内网端口"),
		),
		mcp.WithString(
			"VpcId",
			mcp.Description("vpc标识"),
		),
		mcp.WithString(
			"UniqVpcId",
			mcp.Description("统一vpc标识"),
		),
		mcp.WithString(
			"RegionId",
			mcp.Description("区域标识（gz,bj)"),
		),
		mcp.WithString(
			"ExtraParam",
			mcp.Description("扩展参数"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例Id"),
		),
		mcp.WithString(
			"ProdDbName",
			mcp.Description("数据源产品名"),
		),
		mcp.WithString(
			"DataOrigin",
			mcp.Description("第三方数据源标识"),
		),
		mcp.WithString(
			"DataOriginProjectId",
			mcp.Description("第三方项目id"),
		),
		mcp.WithString(
			"DataOriginDatasourceId",
			mcp.Description("第三方数据源id"),
		),
		mcp.WithString(
			"ClusterId",
			mcp.Description("集群id"),
		),
	)
	mcpsvr.AddTool(biCreateDatasourceCloud, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewCreateDatasourceCloudRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateDatasourceCloud(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biDeleteDatasource := mcp.NewTool(
		"bi-DeleteDatasource",
		mcp.WithDescription(`删除数据源`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Id",
			mcp.Description("数据源id"),
		),
		mcp.WithNumber(
			"ProjectId",
			mcp.Description("项目id"),
		),
	)
	mcpsvr.AddTool(biDeleteDatasource, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewDeleteDatasourceRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DeleteDatasource(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biModifyDatasource := mcp.NewTool(
		"bi-ModifyDatasource",
		mcp.WithDescription(`更新数据源`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"DbHost",
			mcp.Description("HOST"),
		),
		mcp.WithNumber(
			"DbPort",
			mcp.Description("端口"),
		),
		mcp.WithString(
			"ServiceType",
			mcp.Description("后端提供字典：域类型，1、腾讯云，2、本地"),
		),
		mcp.WithString(
			"DbType",
			mcp.Description("驱动"),
		),
		mcp.WithString(
			"Charset",
			mcp.Description("数据库编码"),
		),
		mcp.WithString(
			"DbUser",
			mcp.Description("用户名"),
		),
		mcp.WithString(
			"DbPwd",
			mcp.Description("密码"),
		),
		mcp.WithString(
			"DbName",
			mcp.Description("数据库名称"),
		),
		mcp.WithString(
			"SourceName",
			mcp.Description("数据库别名"),
		),
		mcp.WithNumber(
			"Id",
			mcp.Description("数据源id"),
		),
		mcp.WithNumber(
			"ProjectId",
			mcp.Description("项目ID"),
		),
		mcp.WithString(
			"Catalog",
			mcp.Description("catalog值"),
		),
		mcp.WithString(
			"DataOrigin",
			mcp.Description("第三方数据源标识"),
		),
		mcp.WithString(
			"DataOriginProjectId",
			mcp.Description("第三方项目id"),
		),
		mcp.WithString(
			"DataOriginDatasourceId",
			mcp.Description("第三方数据源id"),
		),
		mcp.WithString(
			"ExtraParam",
			mcp.Description("扩展参数"),
		),
		mcp.WithString(
			"UniqVpcId",
			mcp.Description("腾讯云私有网络统一标识"),
		),
		mcp.WithString(
			"Vip",
			mcp.Description("私有网络ip"),
		),
		mcp.WithString(
			"Vport",
			mcp.Description("私有网络端口"),
		),
		mcp.WithString(
			"VpcId",
			mcp.Description("腾讯云私有网络标识"),
		),
		mcp.WithBoolean(
			"UseVPC",
			mcp.Description("开启vpc	"),
		),
		mcp.WithString(
			"RegionId",
			mcp.Description("地域"),
		),
	)
	mcpsvr.AddTool(biModifyDatasource, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewModifyDatasourceRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ModifyDatasource(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biModifyDatasourceCloud := mcp.NewTool(
		"bi-ModifyDatasourceCloud",
		mcp.WithDescription(`更新云数据库`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ServiceType",
			mcp.Description("后端提供字典：域类型，1、腾讯云，2、本地"),
		),
		mcp.WithString(
			"DbType",
			mcp.Description("驱动"),
		),
		mcp.WithString(
			"Charset",
			mcp.Description("数据库编码"),
		),
		mcp.WithString(
			"DbUser",
			mcp.Description("用户名"),
		),
		mcp.WithString(
			"DbPwd",
			mcp.Description("密码"),
		),
		mcp.WithString(
			"DbName",
			mcp.Description("数据库名称"),
		),
		mcp.WithString(
			"SourceName",
			mcp.Description("数据库别名"),
		),
		mcp.WithString(
			"ProjectId",
			mcp.Description("项目ID"),
		),
		mcp.WithNumber(
			"Id",
			mcp.Description("住键"),
		),
		mcp.WithString(
			"Vip",
			mcp.Description("公有云内网ip"),
		),
		mcp.WithString(
			"Vport",
			mcp.Description("公有云内网端口"),
		),
		mcp.WithString(
			"VpcId",
			mcp.Description("vpc标识"),
		),
		mcp.WithString(
			"UniqVpcId",
			mcp.Description("统一vpc标识"),
		),
		mcp.WithString(
			"RegionId",
			mcp.Description("区域标识（gz,bj)"),
		),
		mcp.WithString(
			"ExtraParam",
			mcp.Description("扩展参数"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("实例id"),
		),
		mcp.WithString(
			"ProdDbName",
			mcp.Description("数据源产品名"),
		),
		mcp.WithString(
			"DataOrigin",
			mcp.Description("第三方数据源标识"),
		),
		mcp.WithString(
			"DataOriginProjectId",
			mcp.Description("第三方项目id"),
		),
		mcp.WithString(
			"DataOriginDatasourceId",
			mcp.Description("第三方数据源id"),
		),
		mcp.WithString(
			"ClusterId",
			mcp.Description("集群id"),
		),
	)
	mcpsvr.AddTool(biModifyDatasourceCloud, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewModifyDatasourceCloudRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ModifyDatasourceCloud(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biDescribePageWidgetList := mcp.NewTool(
		"bi-DescribePageWidgetList",
		mcp.WithDescription(`查询页面组件信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ProjectId",
			mcp.Description("项目id"),
		),
		mcp.WithString(
			"PageId",
			mcp.Description("页面id"),
		),
	)
	mcpsvr.AddTool(biDescribePageWidgetList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewDescribePageWidgetListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribePageWidgetList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biExportScreenPage := mcp.NewTool(
		"bi-ExportScreenPage",
		mcp.WithDescription(`分享页截图导出`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"ProjectId",
			mcp.Description("项目id"),
		),
		mcp.WithString(
			"PageId",
			mcp.Description("页面id"),
		),
		mcp.WithString(
			"CanvasType",
			mcp.Description("画布类型。栅格画布：GRID；自由画布：FREE"),
		),
		mcp.WithString(
			"PicType",
			mcp.Description("图片导出类型。base64；url（有效期：1天）"),
		),
		mcp.WithArray(
			"WidgetIds",
			mcp.Description("组件Ids。为空时，导出整个页面"),
		),
		mcp.WithBoolean(
			"AsyncRequest",
			mcp.Description("是否是异步请求"),
		),
		mcp.WithString(
			"TranId",
			mcp.Description("事务id"),
		),
	)
	mcpsvr.AddTool(biExportScreenPage, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewExportScreenPageRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ExportScreenPage(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	biDescribeDatasourceList := mcp.NewTool(
		"bi-DescribeDatasourceList",
		mcp.WithDescription(`查询数据源列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"ProjectId",
			mcp.Description("无"),
		),
		mcp.WithBoolean(
			"AllPage",
			mcp.Description("返回所有页面，默认false"),
		),
		mcp.WithString(
			"DbName",
			mcp.Description("数据库名称检索"),
		),
		mcp.WithNumber(
			"PageNo",
			mcp.Description("无"),
		),
		mcp.WithNumber(
			"PageSize",
			mcp.Description("无"),
		),
		mcp.WithString(
			"Keyword",
			mcp.Description("搜索关键词"),
		),
		mcp.WithNumber(
			"PermissionType",
			mcp.Description("过滤无权限列表的参数（0 全量，1 使用权限，2 编辑权限）"),
		),
	)
	mcpsvr.AddTool(biDescribeDatasourceList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := bi.NewClient(credential, region_, cpf)
		req := bi.NewDescribeDatasourceListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDatasourceList(req)
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
