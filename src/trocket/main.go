package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	trocket "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308"
	"log"
	"os"
)

func main() {
	mcpServerName := "mcp-server-trocket"
	mcpsvr := server.NewMCPServer(
		"腾讯云 TRocket MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	trocketDescribeFusionInstanceList := mcp.NewTool(
		"trocket-DescribeFusionInstanceList",
		mcp.WithDescription(`查询集群列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("查询结果限制数量"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("查询条件列表"),
		),
		mcp.WithArray(
			"TagFilters",
			mcp.Description("标签过滤器"),
		),
	)
	mcpsvr.AddTool(trocketDescribeFusionInstanceList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeFusionInstanceListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeFusionInstanceList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeInstance := mcp.NewTool(
		"trocket-DescribeInstance",
		mcp.WithDescription(`查询集群信息`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
	)
	mcpsvr.AddTool(trocketDescribeInstance, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeInstanceRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeInstance(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeTopicList := mcp.NewTool(
		"trocket-DescribeTopicList",
		mcp.WithDescription(`查询主题列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("查询条件列表"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("查询结果限制数量"),
		),
		mcp.WithString(
			"FromGroup",
			mcp.Description("按照消费组查询订阅的主题"),
		),
	)
	mcpsvr.AddTool(trocketDescribeTopicList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeTopicListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTopicList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeTopic := mcp.NewTool(
		"trocket-DescribeTopic",
		mcp.WithDescription(`查询主题详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
		mcp.WithString(
			"Topic",
			mcp.Description("主题名称"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("查询条件列表"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("查询结果限制数量"),
		),
	)
	mcpsvr.AddTool(trocketDescribeTopic, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeTopicRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTopic(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeTopicListByGroup := mcp.NewTool(
		"trocket-DescribeTopicListByGroup",
		mcp.WithDescription(`查询消费组订阅的主题列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("查询结果限制数量"),
		),
		mcp.WithString(
			"ConsumerGroup",
			mcp.Description("消费组名称"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("查询条件列表"),
		),
	)
	mcpsvr.AddTool(trocketDescribeTopicListByGroup, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeTopicListByGroupRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeTopicListByGroup(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeProducerList := mcp.NewTool(
		"trocket-DescribeProducerList",
		mcp.WithDescription("查询生产者信息列表"),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("RocketMQ 实例 ID"),
		),
		mcp.WithString(
			"Topic",
			mcp.Description("主题名称"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("查询结果限制数量"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("查询条件列表"),
		),
	)
	mcpsvr.AddTool(trocketDescribeProducerList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeProducerListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeProducerList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeConsumerGroup := mcp.NewTool(
		"trocket-DescribeConsumerGroup",
		mcp.WithDescription(`查询消费组详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
		mcp.WithString(
			"ConsumerGroup",
			mcp.Description("消费组名称"),
		),
	)
	mcpsvr.AddTool(trocketDescribeConsumerGroup, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeConsumerGroupRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeConsumerGroup(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeConsumerGroupList := mcp.NewTool(
		"trocket-DescribeConsumerGroupList",
		mcp.WithDescription(`查询消费组列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("查询条件列表"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("查询结果限制数量"),
		),
		mcp.WithString(
			"FromTopic",
			mcp.Description("查询指定主题下的消费组"),
		),
	)
	mcpsvr.AddTool(trocketDescribeConsumerGroupList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeConsumerGroupListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeConsumerGroupList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeConsumerClient := mcp.NewTool(
		"trocket-DescribeConsumerClient",
		mcp.WithDescription(`查询消费者客户端详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
		mcp.WithString(
			"ClientId",
			mcp.Description("客户端ID"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("查询条件列表"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("查询结果限制数量"),
		),
		mcp.WithString(
			"ConsumerGroup",
			mcp.Description("消费组名称"),
		),
	)
	mcpsvr.AddTool(trocketDescribeConsumerClient, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeConsumerClientRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeConsumerClient(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeConsumerClientList := mcp.NewTool(
		"trocket-DescribeConsumerClientList",
		mcp.WithDescription(`查询消费组下的客户端连接列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
		mcp.WithString(
			"ConsumerGroup",
			mcp.Description("消费组名称"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("查询条件列表"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("查询结果限制数量"),
		),
	)
	mcpsvr.AddTool(trocketDescribeConsumerClientList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeConsumerClientListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeConsumerClientList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeRoleList := mcp.NewTool(
		"trocket-DescribeRoleList",
		mcp.WithDescription(`查询角色列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("查询结果限制数量"),
		),
		mcp.WithArray(
			"Filters",
			mcp.Description("查询条件列表"),
		),
	)
	mcpsvr.AddTool(trocketDescribeRoleList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeRoleListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeRoleList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeMessageList := mcp.NewTool(
		"trocket-DescribeMessageList",
		mcp.WithDescription(`查询消息列表`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
		mcp.WithString(
			"Topic",
			mcp.Description("主题名称"),
		),
		mcp.WithNumber(
			"StartTime",
			mcp.Description("开始时间"),
		),
		mcp.WithNumber(
			"EndTime",
			mcp.Description("结束时间"),
		),
		mcp.WithString(
			"TaskRequestId",
			mcp.Description("一次查询标识"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("查询结果限制数量"),
		),
		mcp.WithString(
			"ConsumerGroup",
			mcp.Description("消费组名称"),
		),
		mcp.WithString(
			"MsgId",
			mcp.Description("消息 ID"),
		),
		mcp.WithString(
			"MsgKey",
			mcp.Description("消息 Key"),
		),
		mcp.WithNumber(
			"RecentMessageNum",
			mcp.Description("查询最近N条消息 最大不超过1024，默认-1为其他查询条件"),
		),
		mcp.WithBoolean(
			"QueryDeadLetterMessage",
			mcp.Description("是否查询死信消息"),
		),
		mcp.WithString(
			"Tag",
			mcp.Description("消息 Tag"),
		),
	)
	mcpsvr.AddTool(trocketDescribeMessageList, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeMessageListRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeMessageList(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeMessage := mcp.NewTool(
		"trocket-DescribeMessage",
		mcp.WithDescription(`查询消息详情`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
		mcp.WithString(
			"Topic",
			mcp.Description("主题名称"),
		),
		mcp.WithString(
			"MsgId",
			mcp.Description("消息ID"),
		),
		mcp.WithNumber(
			"Offset",
			mcp.Description("查询起始位置"),
		),
		mcp.WithNumber(
			"Limit",
			mcp.Description("查询结果限制数量"),
		),
		mcp.WithBoolean(
			"QueryDeadLetterMessage",
			mcp.Description("是否是死信消息"),
		),
		mcp.WithBoolean(
			"QueryDelayMessage",
			mcp.Description("是否是延时消息"),
		),
	)
	mcpsvr.AddTool(trocketDescribeMessage, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeMessageRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeMessage(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", rsp.ToJsonString())), nil
	})

	trocketDescribeMessageTrace := mcp.NewTool(
		"trocket-DescribeMessageTrace",
		mcp.WithDescription(`查询消息轨迹`),
		mcp.WithString(
			"region",
			mcp.Required(),
			mcp.Description("地域"),
		),
		mcp.WithString(
			"InstanceId",
			mcp.Description("集群ID"),
		),
		mcp.WithString(
			"Topic",
			mcp.Description("主题名称"),
		),
		mcp.WithString(
			"MsgId",
			mcp.Description("消息ID"),
		),
		mcp.WithBoolean(
			"QueryDeadLetterMessage",
			mcp.Description("是否是死信消息"),
		),
		mcp.WithBoolean(
			"QueryDelayMessage",
			mcp.Description("是否是延时消息"),
		),
	)
	mcpsvr.AddTool(trocketDescribeMessageTrace, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region_ := "ap-guangzhou"
		arguments := request.GetArguments()
		if nil != arguments["region"] {
			region_ = arguments["region"].(string)
		}
		delete(arguments, "region")
		cpf := profile.NewClientProfile()
		cpf.Debug = true
		client, _ := trocket.NewClient(credential, region_, cpf)
		req := trocket.NewDescribeMessageTraceRequest()
		if nil != arguments {
			jsonstr, _ := json.Marshal(arguments)
			req.FromJsonString(string(jsonstr))
		}
		rsp, err := client.DescribeMessageTrace(req)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("%s", err.Error())), err
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
