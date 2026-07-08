package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// ToolHandler 执行 API 调用并返回 JSON 响应字符串
//
//nolint:revive
type ToolHandler func(client *postgres.Client, args map[string]interface{}) (string, error)

// registerTool 创建并注册一个 MCP Tool，封装公共的 handler 样板代码。
func registerTool(
	s *server.MCPServer,
	cp security.CredentialProvider,
	g *security.Guard,
	name string,
	description string,
	guardLevel security.GuardLevel,
	params []mcp.ToolOption,
	handler ToolHandler,
) {
	if security.NeedsConfirm(guardLevel) {
		params = append(params, mcp.WithBoolean("confirm",
			mcp.Description("确认真执行此操作。设为 true 以确认。")))
	}

	toolOpts := []mcp.ToolOption{
		mcp.WithDescription(description),
		mcp.WithString("region", mcp.Required(), mcp.Description("地域")),
	}
	toolOpts = append(toolOpts, params...)

	tool := mcp.NewTool("postgres-"+name, toolOpts...)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		region := "ap-guangzhou"
		arguments := request.GetArguments()
		if v, ok := arguments["region"].(string); ok && v != "" {
			region = v
		}
		delete(arguments, "region")

		if err := security.AuthorizePrincipal(ctx, name, region, guardLevel); err != nil {
			return mcp.NewToolResultText(fmt.Sprintf(`{"code":403,"error":"%s"}`, err.Error())), nil
		}

		if g != nil {
			if err := g.Check(name, region, guardLevel); err != nil {
				return mcp.NewToolResultText(fmt.Sprintf(`{"code":403,"error":"%s"}`, err.Error())), nil
			}
		}

		confirmed := false
		if v, ok := arguments["confirm"].(bool); ok {
			confirmed = v
		}
		delete(arguments, "confirm")

		if security.NeedsConfirm(guardLevel) && !confirmed {
			warning := security.GetGuardWarning(guardLevel, name)
			return mcp.NewToolResultText(fmt.Sprintf(`{"code":403,"warning":"%s","require_confirm":true}`, warning)), nil
		}

		if cp == nil {
			return mcp.NewToolResultText(`{"code":500,"error":"credential provider is not configured"}`), nil
		}
		cred, err := cp.Resolve(ctx)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf(`{"code":500,"error":"resolve credential failed: %s"}`, err.Error())), nil
		}

		cpf := profile.NewClientProfile()
		cpf.Debug = getDebugFlag()
		client, err := postgres.NewClient(cred, region, cpf)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf(`{"code":500,"error":"create client failed: %s"}`, err.Error())), nil
		}

		result, err := handler(client, arguments)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf(`{"code":500,"error":"%s"}`, err.Error())), nil
		}
		return mcp.NewToolResultText(result), nil
	})
}

func marshalArgs(args map[string]interface{}) string {
	if len(args) == 0 {
		return "{}"
	}
	jsonstr, _ := json.Marshal(args)
	return string(jsonstr)
}

func getDebugFlag() bool {
	return os.Getenv("MCP_DEBUG") == "true"
}

func Log(format string, args ...any) {
	log.Printf(format, args...)
}
