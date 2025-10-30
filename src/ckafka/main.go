package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-ckafka"
	mcpsvr := server.NewMCPServer(
		"腾讯云 Ckafka MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	ckafkaDescribeTopic := mcp.NewTool(
		"ckafka-DescribeTopic",
		mcp.WithDescription(`获取主题列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("过滤条件，按照 topicName 过滤，支持模糊查询"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，不填默认为0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，不填则默认为20，最大值为50"),
		),
		mcp.WithString(
			"AclRuleName",
			mcp.Description("Acl预设策略名称"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeTopic, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeTopicRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTopic(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeTopicAttributes := mcp.NewTool(
		"ckafka-DescribeTopicAttributes",
		mcp.WithDescription(`获取主题属性`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"TopicName",
			mcp.Description("主题名称"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeTopicAttributes, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeTopicAttributesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTopicAttributes(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeTopicDetail := mcp.NewTool(
		"ckafka-DescribeTopicDetail",
		mcp.WithDescription(`获取主题列表详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("（过滤条件）按照topicName过滤，支持模糊查询"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，不填默认为0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，不填则默认 10，最大值20，取值要大于0"),
		),
		mcp.WithString(
			"AclRuleName",
			mcp.Description("Acl预设策略名称"),
		),
		mcp.WithString(
			"OrderBy",
			mcp.Description("根据特定的属性排序(目前支持PartitionNum/CreateTime)"),
		),
		mcp.WithNumber(
			"OrderType",
			mcp.Description("0-顺序、1-倒序"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("目前支持 ReplicaNum （副本数）筛选"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeTopicDetail, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeTopicDetailRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTopicDetail(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeTopicFlowRanking := mcp.NewTool(
		"ckafka-DescribeTopicFlowRanking",
		mcp.WithDescription(`Topic 流量排行`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"RankingType",
			mcp.Description("排行类别，PRO：Topic生产流量；CON：Topic消费流量"),
		),
		mcp.WithString(
			"BeginDate",
			mcp.Description("排行起始日期"),
		),
		mcp.WithString(
			"EndDate",
			mcp.Description("排行结束日期"),
		),
		mcp.WithString(
			"BrokerIp",
			mcp.Description("Broker IP 地址"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeTopicFlowRanking, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeTopicFlowRankingRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTopicFlowRanking(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeTopicProduceConnection := mcp.NewTool(
		"ckafka-DescribeTopicProduceConnection",
		mcp.WithDescription(`查询topic生产端连接信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"TopicName",
			mcp.Description("主题名"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeTopicProduceConnection, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeTopicProduceConnectionRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTopicProduceConnection(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeTopicSubscribeGroup := mcp.NewTool(
		"ckafka-DescribeTopicSubscribeGroup",
		mcp.WithDescription(`查询订阅某主题消息分组信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"TopicName",
			mcp.Description("主题名"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页时的起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("分页时的个数"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeTopicSubscribeGroup, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeTopicSubscribeGroupRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTopicSubscribeGroup(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaFetchMessageByOffset := mcp.NewTool(
		"ckafka-FetchMessageByOffset",
		mcp.WithDescription(`查询消息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"Topic",
			mcp.Description("主题名"),
		),
		mcp.WithNumber(
			"Partition",
			mcp.Description("分区id"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("位点信息，必填"),
		),
	)
	mcpsvr.AddTool(ckafkaFetchMessageByOffset, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewFetchMessageByOffsetRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.FetchMessageByOffset(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaFetchMessageListByOffset := mcp.NewTool(
		"ckafka-FetchMessageListByOffset",
		mcp.WithDescription(`根据位点查询消息列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"Topic",
			mcp.Description("主题名"),
		),
		mcp.WithNumber(
			"Partition",
			mcp.Description("分区id"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("位点信息"),
		),
		mcp.WithNumber(
			"SinglePartitionRecordNumber",
			mcp.Description("最大查询条数，默认20，最大20"),
		),
	)
	mcpsvr.AddTool(ckafkaFetchMessageListByOffset, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewFetchMessageListByOffsetRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.FetchMessageListByOffset(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaFetchMessageListByTimestamp := mcp.NewTool(
		"ckafka-FetchMessageListByTimestamp",
		mcp.WithDescription(`根据时间戳查询消息列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"Topic",
			mcp.Description("主题名"),
		),
		mcp.WithNumber(
			"Partition",
			mcp.Description("分区id"),
		),
		mcp.WithNumber(
			"StartTime",
			mcp.Description("查询开始时间，13位时间戳"),
		),
		mcp.WithNumber(
			"SinglePartitionRecordNumber",
			mcp.Description("最大查询条数，默认20，最大20, 最小1"),
		),
	)
	mcpsvr.AddTool(ckafkaFetchMessageListByTimestamp, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewFetchMessageListByTimestampRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.FetchMessageListByTimestamp(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeGroupInfo := mcp.NewTool(
		"ckafka-DescribeGroupInfo",
		mcp.WithDescription(`获取消费分组信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithArray(
			"GroupList",
			mcp.Description("Kafka 消费分组列表"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeGroupInfo, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeGroupInfoRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeGroupInfo(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeGroupOffsets := mcp.NewTool(
		"ckafka-DescribeGroupOffsets",
		mcp.WithDescription(`获取消费分组offset`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"Group",
			mcp.Description("Kafka 消费分组"),
		),
		mcp.WithArray(
			"Topics",
			mcp.Description("group 订阅的主题名称数组，如果没有该数组，则表示指定的 group 下所有 topic 信息"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("模糊匹配 topicName"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("本次查询的偏移位置，默认为0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("本次返回结果的最大个数，默认为50，最大值为50"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeGroupOffsets, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeGroupOffsetsRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeGroupOffsets(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeConsumerGroup := mcp.NewTool(
		"ckafka-DescribeConsumerGroup",
		mcp.WithDescription(`查询消费分组信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"GroupName",
			mcp.Description("可选，用户需要查询的group名称。"),
		),
		mcp.WithString(
			"TopicName",
			mcp.Description("可选，用户需要查询的group中的对应的topic名称，如果指定了该参数，而group又未指定则忽略该参数。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("本次返回个数限制，最大支持50"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移位置"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeConsumerGroup, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeConsumerGroupRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeConsumerGroup(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeGroup := mcp.NewTool(
		"ckafka-DescribeGroup",
		mcp.WithDescription(`枚举消费分组(精简版)`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("搜索关键字"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("最大返回数量"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("仅支持 GroupState 筛选,   支持的筛选状态有 Empty/Stable  注意：该参数只能在2.8/3.2 版本生效"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeGroup, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeGroupRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeGroup(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeACL := mcp.NewTool(
		"ckafka-DescribeACL",
		mcp.WithDescription(`枚举ACL`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithNumber(
			"ResourceType",
			mcp.Description("Acl资源类型，(2:TOPIC，3:GROUP，4:CLUSTER)"),
		),
		mcp.WithString(
			"ResourceName",
			mcp.Description("资源名称，和resourceType相关，如当resourceType为TOPIC时，则该字段表示topic名称，当resourceType为GROUP时，该字段表示group名称，当resourceType为CLUSTER时，该字段可为空。"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("个数限制"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("关键字匹配"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeACL, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeACLRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeACL(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeAclRule := mcp.NewTool(
		"ckafka-DescribeAclRule",
		mcp.WithDescription(`查询ACL规则列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
		mcp.WithString(
			"RuleName",
			mcp.Description("ACL规则名"),
		),
		mcp.WithString(
			"PatternType",
			mcp.Description("ACL规则匹配类型"),
		),
		mcp.WithBoolean(
			"IsSimplified",
			mcp.Description("是否读取简略的ACL规则"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeAclRule, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeAclRuleRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeAclRule(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeInstanceAttributes := mcp.NewTool(
		"ckafka-DescribeInstanceAttributes",
		mcp.WithDescription(`获取实例属性`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("ckafka集群实例Id"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeInstanceAttributes, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeInstanceAttributesRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstanceAttributes(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeInstances := mcp.NewTool(
		"ckafka-DescribeInstances",
		mcp.WithDescription(`获取实例列表信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("（查询条件）按照ckafka集群实例Id过滤"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("搜索词   ex:（查询条件）按照实例名称过滤，支持模糊查询"),
		),
		mcp.WithArray(
			"Status",
			mcp.Description("（查询条件）实例的状态  0：创建中，1：运行中，2：删除中，5: 隔离中,  7:升级中 不填默认返回全部"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，不填默认为0"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，不填则默认10，最大值100"),
		),
		mcp.WithString(
			"TagKey",
			mcp.Description("已废弃。匹配标签key值。"),
		),
		mcp.WithString(
			"VpcId",
			mcp.Description("（查询条件）私有网络Id"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeInstances, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeInstancesRequest()
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

	ckafkaDescribeInstancesDetail := mcp.NewTool(
		"ckafka-DescribeInstancesDetail",
		mcp.WithDescription(`获取实例集群列表详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("（过滤条件）按照实例ID过滤"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("（过滤条件）按照实例名,实例Id,可用区,私有网络id,子网id 过滤，支持模糊查询"),
		),
		mcp.WithArray(
			"Status",
			mcp.Description("（过滤条件）实例的状态。0：创建中，1：运行中，2：删除中，不填默认返回全部"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("偏移量，不填默认为0。"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，不填则默认10，最大值20。"),
		),
		mcp.WithString(
			"TagKey",
			mcp.Description("匹配标签key值。"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("过滤器。filter.Name 支持('Ip', 'VpcId', 'SubNetId', 'InstanceType','InstanceId') ,filter.Values最多传递10个值."),
		),
		mcp.WithArray(
			"InstanceIdList",
			mcp.Description("按照实例ID过滤"),
		),
		mcp.WithArray(
			"TagList",
			mcp.Description("根据标签列表过滤实例（取交集）"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeInstancesDetail, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeInstancesDetailRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstancesDetail(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaInquireCkafkaPrice := mcp.NewTool(
		"ckafka-InquireCkafkaPrice",
		mcp.WithDescription(`Ckafka询价`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceType",
			mcp.Description("国内站标准版填写standards2, 国际站标准版填写standard,专业版填写profession,高级版填写premium"),
		),
		mcp.WithObject(
			"InstanceChargeParam",
			mcp.Description("购买/续费付费类型(购买时不填的话, 默认获取购买包年包月一个月的费用)"),
		),
		mcp.WithNumber(
			"InstanceNum",
			mcp.Description("购买/续费时购买的实例数量(不填时, 默认为1个)"),
		),
		mcp.WithNumber(
			"Bandwidth",
			mcp.Description("实例内网带宽大小, 单位MB/s (购买时必填，专业版/高级版询价时带宽信息必填)"),
		),
		mcp.WithObject(
			"InquiryDiskParam",
			mcp.Description("实例的硬盘购买类型以及大小 (购买时必填，专业版/高级版询价时磁盘信息必填)"),
		),
		mcp.WithNumber(
			"MessageRetention",
			mcp.Description("实例消息保留时间大小, 单位小时 (购买时必填)"),
		),
		mcp.WithNumber(
			"Topic",
			mcp.Description("购买实例topic数, 单位个 (购买时必填)"),
		),
		mcp.WithNumber(
			"Partition",
			mcp.Description("购买实例分区数, 单位个 (购买时必填，专业版/高级版询价时带宽信息必填)"),
		),
		mcp.WithArray(
			"ZoneIds",
			mcp.Description("购买地域, 可通过查看DescribeCkafkaZone这个接口获取ZoneId"),
		),
		mcp.WithString(
			"CategoryAction",
			mcp.Description("标记操作, 新购填写purchase, 续费填写renew, (不填时, 默认为purchase)"),
		),
		mcp.WithString(
			"BillType",
			mcp.Description("国内站购买的版本, sv_ckafka_instance_s2_1(入门型), sv_ckafka_instance_s2_2(标准版), sv_ckafka_instance_s2_3(进阶型), 如果instanceType为standards2, 但该参数为空, 则默认值为sv_ckafka_instance_s2_1"),
		),
		mcp.WithObject(
			"PublicNetworkParam",
			mcp.Description("公网带宽计费模式, 目前只有专业版支持公网带宽 (购买公网带宽时必填)"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("续费时的实例id, 续费时填写"),
		),
	)
	mcpsvr.AddTool(ckafkaInquireCkafkaPrice, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewInquireCkafkaPriceRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.InquireCkafkaPrice(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeCkafkaZone := mcp.NewTool(
		"ckafka-DescribeCkafkaZone",
		mcp.WithDescription(`查看可用区列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"CdcId",
			mcp.Description("cdc集群Id"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeCkafkaZone, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeCkafkaZoneRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeCkafkaZone(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeDatahubTask := mcp.NewTool(
		"ckafka-DescribeDatahubTask",
		mcp.WithDescription(`查询Datahub任务信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"TaskId",
			mcp.Description("任务id"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeDatahubTask, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeDatahubTaskRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDatahubTask(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	ckafkaDescribeDatahubTasks := mcp.NewTool(
		"ckafka-DescribeDatahubTasks",
		mcp.WithDescription(`查询Datahub任务列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("返回数量，默认为20，最大值为100"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("分页偏移量，默认为0"),
		),
		mcp.WithString(
			"SearchWord",
			mcp.Description("过滤条件，按照 TaskName 过滤，支持模糊查询"),
		),
		mcp.WithString(
			"TargetType",
			mcp.Description("转储的目标类型"),
		),
		mcp.WithString(
			"TaskType",
			mcp.Description("任务类型，SOURCE数据接入，SINK数据流出"),
		),
		mcp.WithString(
			"SourceType",
			mcp.Description("转储的源类型"),
		),
		mcp.WithString(
			"Resource",
			mcp.Description("转储的资源"),
		),
	)
	mcpsvr.AddTool(ckafkaDescribeDatahubTasks, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := ckafka.NewClient(credential, region_, cpf)
		req := ckafka.NewDescribeDatahubTasksRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeDatahubTasks(req)
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
