package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func main() {
	mcpServerName := "mcp-server-scf"
	mcpsvr := server.NewMCPServer(
		"腾讯云 SCF MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	scfCreateFunction := mcp.NewTool(
		"scf-CreateFunction",
		mcp.WithDescription(`创建函数`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("创建的函数名称，函数名称支持26个英文字母大小写、数字、连接符和下划线，第一个字符只能以字母开头，最后一个字符不能为连接符或者下划线，名称长度2-60"),
		),
		mcp.WithObject(
			"Code",
			mcp.Description("函数代码. 注意：不能同时指定Cos、ZipFile或 DemoId。详情请参考 https://cloud.tencent.com/document/api/583/17244#Code"),
		),
		mcp.WithString(
			"Handler",
			mcp.Description("函数处理方法名称，名称格式支持 &quot;文件名称.方法名称&quot; 形式（java 名称格式 包名.类名::方法名），文件名称和函数名称之间以&quot;.&quot;隔开，文件名称和函数名称要求以字母开始和结尾，中间允许插入字母、数字、下划线和连接符，文件名称和函数名字的长度要求是 2-60 个字符"),
		),
		mcp.WithString(
			"Description",
			mcp.Description("函数描述,最大支持 1000 个英文字母、数字、空格、逗号、换行符和英文句号，支持中文"),
		),
		mcp.WithNumber(
			"MemorySize",
			mcp.Description("函数运行时内存大小，默认为 128M，可选范围 64、128MB-3072MB，并且以 128MB 为阶梯"),
		),
		mcp.WithNumber(
			"Timeout",
			mcp.Description("函数最长执行时间，单位为秒，可选值范围 1-900 秒，默认为 3 秒"),
		),
		mcp.WithObject(
			"Environment",
			mcp.Description("函数的环境变量"),
		),
		mcp.WithString(
			"Runtime",
			mcp.Description("函数运行环境，默认Python2.7"),
		),
		mcp.WithObject(
			"VpcConfig",
			mcp.Description("函数的私有网络配置"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("函数所属命名空间"),
		),
		mcp.WithString(
			"Role",
			mcp.Description("函数绑定的角色"),
		),
		mcp.WithString(
			"InstallDependency",
			mcp.Description("[在线依赖安装](https://cloud.tencent.com/document/product/583/37920)，TRUE 表示安装，默认值为 FALSE。仅支持 Node.js 函数。"),
		),
		mcp.WithString(
			"ClsLogsetId",
			mcp.Description("函数日志投递到的CLS LogsetID"),
		),
		mcp.WithString(
			"ClsTopicId",
			mcp.Description("函数日志投递到的CLS TopicID"),
		),
		mcp.WithString(
			"Type",
			mcp.Description("函数类型，默认值为Event，创建触发器函数请填写Event，创建HTTP函数级服务请填写HTTP"),
		),
		mcp.WithArray(
			"Layers",
			mcp.Description("函数要关联的Layer版本列表，Layer会按照在列表中顺序依次覆盖。"),
		),
		mcp.WithObject(
			"DeadLetterConfig",
			mcp.Description("死信队列参数"),
		),
		mcp.WithObject(
			"PublicNetConfig",
			mcp.Description("公网访问配置"),
		),
		mcp.WithObject(
			"CfsConfig",
			mcp.Description("文件系统配置参数，用于云函数挂载文件系统"),
		),
		mcp.WithNumber(
			"InitTimeout",
			mcp.Description("函数初始化超时时间，默认 65s，镜像部署函数默认 90s。"),
		),
		mcp.WithArray(
			"Tags",
			mcp.Description("函数 Tag 参数，以键值对数组形式传入"),
		),
		mcp.WithString(
			"AsyncRunEnable",
			mcp.Description("是否开启异步属性，TRUE 为开启，FALSE为关闭"),
		),
		mcp.WithString(
			"TraceEnable",
			mcp.Description("是否开启事件追踪，TRUE 为开启，FALSE为关闭"),
		),
		mcp.WithString(
			"AutoDeployClsTopicIndex",
			mcp.Description("是否自动创建cls索引，TRUE 为开启，FALSE为关闭"),
		),
		mcp.WithString(
			"AutoCreateClsTopic",
			mcp.Description("是否自动创建cls主题，TRUE 为开启，FALSE为关闭"),
		),
		mcp.WithString(
			"ProtocolType",
			mcp.Description("HTTP函数支持的访问协议。当前支持WebSockets协议，值为WS"),
		),
		mcp.WithObject(
			"ProtocolParams",
			mcp.Description("HTTP函数配置ProtocolType访问协议，当前协议可配置的参数"),
		),
		mcp.WithObject(
			"InstanceConcurrencyConfig",
			mcp.Description("单实例多并发配置。只支持Web函数。"),
		),
		mcp.WithString(
			"DnsCache",
			mcp.Description("是否开启Dns缓存能力。只支持EVENT函数。默认为FALSE，TRUE 为开启，FALSE为关闭"),
		),
		mcp.WithObject(
			"IntranetConfig",
			mcp.Description("内网访问配置"),
		),
	)
	mcpsvr.AddTool(scfCreateFunction, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewCreateFunctionRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateFunction(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfUpdateFunctionCode := mcp.NewTool(
		"scf-UpdateFunctionCode",
		mcp.WithDescription(`更新函数代码`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("要修改的函数名称"),
		),
		mcp.WithString(
			"Handler",
			mcp.Description("函数处理方法名称。名称格式支持“文件名称.函数名称”形式（java 名称格式 包名.类名::方法名），文件名称和函数名称之间以&quot;.&quot;隔开，文件名称和函数名称要求以字母开始和结尾，中间允许插入字母、数字、下划线和连接符，文件名称和函数名字的长度要求 2-60 个字符"),
		),
		mcp.WithString(
			"CosBucketName",
			mcp.Description("对象存储桶名称"),
		),
		mcp.WithString(
			"CosObjectName",
			mcp.Description("对象存储对象路径"),
		),
		mcp.WithString(
			"ZipFile",
			mcp.Description("包含函数代码文件及其依赖项的 zip 格式文件，使用该接口时要求将 zip 文件的内容转成 base64 编码，最大支持20M"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("函数所属命名空间"),
		),
		mcp.WithString(
			"CosBucketRegion",
			mcp.Description("对象存储的地域，注：北京分为ap-beijing和ap-beijing-1"),
		),
		mcp.WithString(
			"InstallDependency",
			mcp.Description("是否自动安装依赖"),
		),
		mcp.WithString(
			"EnvId",
			mcp.Description("函数所属环境"),
		),
		mcp.WithString(
			"Publish",
			mcp.Description("在更新时是否同步发布新版本，默认为：FALSE，不发布"),
		),
		mcp.WithObject(
			"Code",
			mcp.Description("函数代码，详情请参考 https://cloud.tencent.com/document/api/583/17244#Code"),
		),
	)
	mcpsvr.AddTool(scfUpdateFunctionCode, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewUpdateFunctionCodeRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.UpdateFunctionCode(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfUpdateFunctionConfiguration := mcp.NewTool(
		"scf-UpdateFunctionConfiguration",
		mcp.WithDescription(`更新函数配置`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("要修改的函数名称"),
		),
		mcp.WithString(
			"Description",
			mcp.Description("函数描述。最大支持 1000 个英文字母、数字、空格、逗号和英文句号，支持中文"),
		),
		mcp.WithString(
			"MemorySize",
			mcp.Description("函数运行时内存大小，默认为 128 M，可选范64M、128 M-3072 M，以 128MB 为阶梯。"),
		),
		mcp.WithString(
			"Timeout",
			mcp.Description("函数最长执行时间，单位为秒，可选值范 1-900 秒，默认为 3 秒"),
		),
		mcp.WithString(
			"Runtime",
			mcp.Description("函数运行环境，创建时指定，目前不支持修改。"),
		),
		mcp.WithObject(
			"Environment",
			mcp.Description("函数的环境变量"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("函数所属命名空间"),
		),
		mcp.WithObject(
			"VpcConfig",
			mcp.Description("函数的私有网络配置"),
		),
		mcp.WithString(
			"Role",
			mcp.Description("函数绑定的角色"),
		),
		mcp.WithString(
			"InstallDependency",
			mcp.Description("[在线依赖安装](https://cloud.tencent.com/document/product/583/37920)，TRUE 表示安装，默认值为 FALSE。仅支持 Node.js 函数。"),
		),
		mcp.WithString(
			"ClsLogsetId",
			mcp.Description("日志投递到的cls日志集ID"),
		),
		mcp.WithString(
			"ClsTopicId",
			mcp.Description("日志投递到的cls Topic ID"),
		),
		mcp.WithString(
			"Publish",
			mcp.Description("在更新时是否同步发布新版本，默认为：FALSE，不发布新版本"),
		),
		mcp.WithString(
			"L5Enable",
			mcp.Description("是否开启L5访问能力，TRUE 为开启，FALSE为关闭"),
		),
		mcp.WithArray(
			"Layers",
			mcp.Description("函数要关联的层版本列表，层的版本会按照在列表中顺序依次覆盖。"),
		),
		mcp.WithObject(
			"DeadLetterConfig",
			mcp.Description("函数关联的死信队列信息"),
		),
		mcp.WithObject(
			"PublicNetConfig",
			mcp.Description("公网访问配置"),
		),
		mcp.WithObject(
			"CfsConfig",
			mcp.Description("文件系统配置入参，用于云函数绑定CFS文件系统"),
		),
		mcp.WithNumber(
			"InitTimeout",
			mcp.Description("函数初始化执行超时时间"),
		),
		mcp.WithObject(
			"ProtocolParams",
			mcp.Description("HTTP函数配置ProtocolType访问协议，当前协议可配置的参数"),
		),
		mcp.WithObject(
			"InstanceConcurrencyConfig",
			mcp.Description("单实例多并发配置。只支持Web函数。"),
		),
		mcp.WithString(
			"DnsCache",
			mcp.Description("是否开启Dns缓存能力。只支持EVENT函数。默认为FALSE，TRUE 为开启，FALSE为关闭"),
		),
		mcp.WithObject(
			"IntranetConfig",
			mcp.Description("内网访问配置"),
		),
		mcp.WithBoolean(
			"IgnoreSysLog",
			mcp.Description("忽略系统日志上报"),
		),
	)
	mcpsvr.AddTool(scfUpdateFunctionConfiguration, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewUpdateFunctionConfigurationRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.UpdateFunctionConfiguration(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfPublishVersion := mcp.NewTool(
		"scf-PublishVersion",
		mcp.WithDescription(`发布新版本`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("发布函数的名称"),
		),
		mcp.WithString(
			"Description",
			mcp.Description("函数的描述"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("函数的命名空间"),
		),
	)
	mcpsvr.AddTool(scfPublishVersion, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewPublishVersionRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.PublishVersion(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfDeleteFunction := mcp.NewTool(
		"scf-DeleteFunction",
		mcp.WithDescription(`删除函数`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("要删除的函数名称"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("函数所属命名空间"),
		),
		mcp.WithString(
			"Qualifier",
			mcp.Description("填写需要删除的版本号，不填默认删除函数下全部版本。"),
		),
	)
	mcpsvr.AddTool(scfDeleteFunction, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewDeleteFunctionRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DeleteFunction(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfGetFunction := mcp.NewTool(
		"scf-GetFunction",
		mcp.WithDescription(`获取函数详细信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("需要获取详情的函数名称，ResourceId和FunctionName只能传一个"),
		),
		mcp.WithString(
			"Qualifier",
			mcp.Description("函数的版本号"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("函数所属命名空间"),
		),
		mcp.WithString(
			"ShowCode",
			mcp.Description("是否显示代码, TRUE表示显示代码，FALSE表示不显示代码,大于1M的入口文件不会显示"),
		),
	)
	mcpsvr.AddTool(scfGetFunction, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewGetFunctionRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.GetFunction(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfGetAlias := mcp.NewTool(
		"scf-GetAlias",
		mcp.WithDescription(`获取别名详细信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("函数名称"),
		),
		mcp.WithString(
			"Name",
			mcp.Description("别名的名称"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("函数所在的命名空间"),
		),
	)
	mcpsvr.AddTool(scfGetAlias, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewGetAliasRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.GetAlias(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfListFunctions := mcp.NewTool(
		"scf-ListFunctions",
		mcp.WithDescription(`获取函数列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Order",
			mcp.Description("以升序还是降序的方式返回结果，可选值 ASC 和 DESC"),
		),
		mcp.WithString(
			"Orderby",
			mcp.Description("根据哪个字段进行返回结果排序,支持以下字段：AddTime, ModTime, FunctionName"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("数据偏移量，默认值为 0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数据长度，默认值为 20"),
		),
		mcp.WithString(
			"SearchKey",
			mcp.Description("支持FunctionName模糊匹配"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("命名空间"),
		),
		mcp.WithString(
			"Description",
			mcp.Description("函数描述，支持模糊搜索"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("`过滤特定属性或者有特定标签的函数。`- 传值方式key-value 进行传值  例如：&quot;Filters&quot;: [{ &quot;Name&quot;: &quot;Status&quot;, &quot;Values&quot;: [&quot;CreateFailed&quot;,&quot;Creating&quot;]}, {&quot;Name&quot;: &quot;Type&quot;,&quot;Values&quot;: [&quot;HTTP&quot;]}]上述条件的函数是，函数状态为创建失败或者创建中，且函数类型为 HTTP 函数`如果通过标签进行过滤：`- tag:tag-key - String - 是否必填：否 - （过滤条件）按照标签键值对进行过滤。 tag-key使用具体的标签键进行替换。示例值：&quot;Filters&quot;: [{&quot;Name&quot;:&quot;tag-dmtest&quot;,&quot;Values&quot;:[&quot;dmtest&quot;]}]`入参限制：`1.每次请求的Filters的上限为10，Filter.Values的上限为5。2.[VpcId', 'SubnetId', 'ClsTopicId', 'ClsLogsetId', 'Role', 'CfsId', 'CfsMountInsId', 'Eip'] 过滤的Name 为这些属性时， values 只能传一个值3.['Status', 'Runtime', 'Type', 'PublicNetStatus', 'AsyncRunEnable', 'TraceEnable', 'Stamp'] 过滤的Name 为这些属性时 ，values 可以传多个值"),
		),
	)
	mcpsvr.AddTool(scfListFunctions, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewListFunctionsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ListFunctions(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfCreateTrigger := mcp.NewTool(
		"scf-CreateTrigger",
		mcp.WithDescription(`设置函数触发方式`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("新建触发器绑定的函数名称"),
		),
		mcp.WithString(
			"TriggerName",
			mcp.Description("新建触发器名称。如果是定时触发器，名称支持英文字母、数字、连接符和下划线，最长100个字符；如果是cos触发器，需要是对应cos存储桶适用于XML API的访问域名(例如:5401-5ff414-12345.cos.ap-shanghai.myqcloud.com);如果是其他触发器，见具体触发器绑定参数的说明"),
		),
		mcp.WithString(
			"Type",
			mcp.Description("触发器类型，目前支持 cos 、cls 、 timer、 ckafka、http类型。创建函数 URL 请使用 http 类型，参考[创建函数 URL ](https://cloud.tencent.com/document/product/583/100227#33bbbda4-9131-48a6-ac37-ac62ffe01424)。创建 cls 触发器请参考[CLS 创建投递 SCF 任务](https://cloud.tencent.com/document/product/614/61096)。"),
		),
		mcp.WithString(
			"TriggerDesc",
			mcp.Description("触发器对应的参数，可见具体[触发器描述说明](https://cloud.tencent.com/document/product/583/39901)"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("函数的命名空间"),
		),
		mcp.WithString(
			"Qualifier",
			mcp.Description("触发器所生效的版本或别名，建议填写 [$DEFAULT](https://cloud.tencent.com/document/product/583/36149#.E9.BB.98.E8.AE.A4.E5.88.AB.E5.90.8D)方便后续进行版本的灰度发布，默认为 $LATEST。"),
		),
		mcp.WithString(
			"Enable",
			mcp.Description("触发器的初始是能状态 OPEN表示开启 CLOSE表示关闭"),
		),
		mcp.WithString(
			"CustomArgument",
			mcp.Description("用户自定义参数，仅支持timer触发器"),
		),
		mcp.WithString(
			"Description",
			mcp.Description("触发器描述"),
		),
	)
	mcpsvr.AddTool(scfCreateTrigger, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewCreateTriggerRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateTrigger(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfDeleteTrigger := mcp.NewTool(
		"scf-DeleteTrigger",
		mcp.WithDescription(`删除触发器`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("函数的名称"),
		),
		mcp.WithString(
			"TriggerName",
			mcp.Description("要删除的触发器名称"),
		),
		mcp.WithString(
			"Type",
			mcp.Description("要删除的触发器类型，目前只支持  timer、ckafka 、apigw 、cls 、cos 、cmq 、http 类型"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("函数所属命名空间"),
		),
		mcp.WithString(
			"TriggerDesc",
			mcp.Description("如果删除的触发器类型为 COS 触发器，该字段为必填值，存放 JSON 格式的数据 {&quot;event&quot;:&quot;cos:ObjectCreated:*&quot;}，数据内容和 SetTrigger 接口中该字段的格式相同；如果删除的触发器类型为定时触发器或 CMQ 触发器，可以不指定该字段"),
		),
		mcp.WithString(
			"Qualifier",
			mcp.Description("要删除的触发器实际所指向的版本或别名，默认值为 $LATEST"),
		),
	)
	mcpsvr.AddTool(scfDeleteTrigger, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewDeleteTriggerRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DeleteTrigger(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfCreateCustomDomain := mcp.NewTool(
		"scf-CreateCustomDomain",
		mcp.WithDescription(`创建云函数自定义域名`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Domain",
			mcp.Description("域名，不支持泛域名"),
		),
		mcp.WithString(
			"Protocol",
			mcp.Description("协议，取值范围：HTTP, HTTPS, HTTP&HTTPS"),
		),
		mcp.WithArray(
			"EndpointsConfig",
			mcp.Description("路由配置"),
		),
		mcp.WithObject(
			"CertConfig",
			mcp.Description("证书配置信息，有使用HTTPS协议时候必须传"),
		),
		mcp.WithObject(
			"WafConfig",
			mcp.Description("web 应用防火墙配置"),
		),
		mcp.WithArray(
			"Tags",
			mcp.Description("标签"),
		),
	)
	mcpsvr.AddTool(scfCreateCustomDomain, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewCreateCustomDomainRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.CreateCustomDomain(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfDeleteCustomDomain := mcp.NewTool(
		"scf-DeleteCustomDomain",
		mcp.WithDescription(`删除云函数自定义域名`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Domain",
			mcp.Description("域名"),
		),
	)
	mcpsvr.AddTool(scfDeleteCustomDomain, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewDeleteCustomDomainRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DeleteCustomDomain(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfGetCustomDomain := mcp.NewTool(
		"scf-GetCustomDomain",
		mcp.WithDescription(`查看云函数自定义域名详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"Domain",
			mcp.Description("域名"),
		),
	)
	mcpsvr.AddTool(scfGetCustomDomain, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewGetCustomDomainRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.GetCustomDomain(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfInvoke := mcp.NewTool(
		"scf-Invoke",
		mcp.WithDescription(`运行函数`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("函数名称"),
		),
		mcp.WithString(
			"InvocationType",
			mcp.Description("同步调用请使用[同步 Invoke 调用接口](https://cloud.tencent.com/document/product/583/58400) 或填写同步调用参数 RequestResponse ，建议使用同步调用接口以获取最佳性能；异步调用填写 Event；默认为同步。接口超时时间为 300s，更长超时时间请使用异步调用。"),
		),
		mcp.WithString(
			"Qualifier",
			mcp.Description("触发函数的版本号或别名，默认值为 $LATEST"),
		),
		mcp.WithString(
			"ClientContext",
			mcp.Description("运行函数时的参数，以json格式传入，同步调用最大支持 6MB，异步调用最大支持 128 KB。该字段信息对应函数 [event 入参](https://cloud.tencent.com/document/product/583/9210#.E5.87.BD.E6.95.B0.E5.85.A5.E5.8F.82.3Ca-id.3D.22input.22.3E.3C.2Fa.3E)。"),
		),
		mcp.WithString(
			"LogType",
			mcp.Description("异步调用该字段返回为空。"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("命名空间"),
		),
		mcp.WithString(
			"RoutingKey",
			mcp.Description("函数灰度流量控制调用，以json格式传入，例如{&quot;k&quot;:&quot;v&quot;}，注意kv都需要是字符串类型，最大支持的参数长度是1024字节"),
		),
	)
	mcpsvr.AddTool(scfInvoke, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewInvokeRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.Invoke(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfGetAsyncEventStatus := mcp.NewTool(
		"scf-GetAsyncEventStatus",
		mcp.WithDescription(`获取函数异步事件状态`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InvokeRequestId",
			mcp.Description("异步执行请求 id"),
		),
	)
	mcpsvr.AddTool(scfGetAsyncEventStatus, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewGetAsyncEventStatusRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.GetAsyncEventStatus(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfTerminateAsyncEvent := mcp.NewTool(
		"scf-TerminateAsyncEvent",
		mcp.WithDescription(`终止函数异步事件`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("函数名称"),
		),
		mcp.WithString(
			"InvokeRequestId",
			mcp.Description("终止的调用请求id"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("命名空间"),
		),
		mcp.WithBoolean(
			"GraceShutdown",
			mcp.Description("true，向指定请求[发送 SIGTERM 终止信号](https://cloud.tencent.com/document/product/583/63969#.E5.8F.91.E9.80.81.E7.BB.88.E6.AD.A2.E4.BF.A1.E5.8F.B7]， ，默认值为 false。"),
		),
	)
	mcpsvr.AddTool(scfTerminateAsyncEvent, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewTerminateAsyncEventRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.TerminateAsyncEvent(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfListAsyncEvents := mcp.NewTool(
		"scf-ListAsyncEvents",
		mcp.WithDescription(`拉取函数异步事件列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("函数名称"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("命名空间"),
		),
		mcp.WithString(
			"Qualifier",
			mcp.Description("过滤条件，函数版本"),
		),
		mcp.WithArray(
			"InvokeType",
			mcp.Description("过滤条件，调用类型列表"),
		),
		mcp.WithArray(
			"Status",
			mcp.Description("过滤条件，事件状态列表"),
		),
		mcp.WithObject(
			"StartTimeInterval",
			mcp.Description("过滤条件，开始执行时间左闭右开区间"),
		),
		mcp.WithObject(
			"EndTimeInterval",
			mcp.Description("过滤条件，结束执行时间左闭右开区间"),
		),
		mcp.WithString(
			"Order",
			mcp.Description("可选值 ASC 和 DESC，默认 DESC"),
		),
		mcp.WithString(
			"Orderby",
			mcp.Description("可选值 StartTime 和 EndTime，默认值 StartTime"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("数据偏移量，默认值为 0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数据长度，默认值为 20，最大值 100"),
		),
		mcp.WithString(
			"InvokeRequestId",
			mcp.Description("过滤条件，事件调用请求id"),
		),
	)
	mcpsvr.AddTool(scfListAsyncEvents, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewListAsyncEventsRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.ListAsyncEvents(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfGetLayerVersion := mcp.NewTool(
		"scf-GetLayerVersion",
		mcp.WithDescription(`获取层版本详细信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"LayerName",
			mcp.Description("层名称"),
		),
		mcp.WithNumber(
			"LayerVersion",
			mcp.Description("版本号"),
		),
	)
	mcpsvr.AddTool(scfGetLayerVersion, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewGetLayerVersionRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.GetLayerVersion(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfDeleteLayerVersion := mcp.NewTool(
		"scf-DeleteLayerVersion",
		mcp.WithDescription(`删除层版本`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"LayerName",
			mcp.Description("层名称"),
		),
		mcp.WithNumber(
			"LayerVersion",
			mcp.Description("版本号"),
		),
	)
	mcpsvr.AddTool(scfDeleteLayerVersion, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewDeleteLayerVersionRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DeleteLayerVersion(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfPublishLayerVersion := mcp.NewTool(
		"scf-PublishLayerVersion",
		mcp.WithDescription(`发布层版本`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"LayerName",
			mcp.Description("层名称，支持26个英文字母大小写、数字、连接符和下划线，第一个字符只能以字母开头，最后一个字符不能为连接符或者下划线，名称长度1-64"),
		),
		mcp.WithArray(
			"CompatibleRuntimes",
			mcp.Description("层适用的运行时，可多选，可选的值对应函数的 Runtime 可选值。"),
		),
		mcp.WithObject(
			"Content",
			mcp.Description("层的文件来源或文件内容"),
		),
		mcp.WithString(
			"Description",
			mcp.Description("层的版本的描述"),
		),
		mcp.WithString(
			"LicenseInfo",
			mcp.Description("层的软件许可证"),
		),
		mcp.WithArray(
			"Tags",
			mcp.Description("层Tag 参数，以键值对数组形式传入"),
		),
	)
	mcpsvr.AddTool(scfPublishLayerVersion, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewPublishLayerVersionRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.PublishLayerVersion(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	scfUpdateAlias := mcp.NewTool(
		"scf-UpdateAlias",
		mcp.WithDescription(`更新别名的配置`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"FunctionName",
			mcp.Description("函数名称"),
		),
		mcp.WithString(
			"Name",
			mcp.Description("别名的名称"),
		),
		mcp.WithString(
			"FunctionVersion",
			mcp.Description("别名指向的主版本"),
		),
		mcp.WithString(
			"Namespace",
			mcp.Description("函数所在的命名空间"),
		),
		mcp.WithObject(
			"RoutingConfig",
			mcp.Description("别名的路由信息，需要为别名指定附加版本时，必须提供此参数；	  附加版本指的是：除主版本 FunctionVersion 外，为此别名再指定一个函数可正常使用的版本；   这里附加版本中的 Version 值 不能是别名指向的主版本；  要注意的是：如果想要某个版本的流量全部指向这个别名，不需配置此参数； 目前一个别名最多只能指定一个附加版本"),
		),
		mcp.WithString(
			"Description",
			mcp.Description("别名的描述"),
		),
	)
	mcpsvr.AddTool(scfUpdateAlias, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := scf.NewClient(credential, region_, cpf)
		req := scf.NewUpdateAliasRequest()
		if nil != request.Params.Arguments {
			jsonstr, _ := json.Marshal(request.Params.Arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.UpdateAlias(req)
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
