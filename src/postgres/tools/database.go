package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterDatabaseTools 注册数据库管理工具（4个：2现有迁移 + 2新增）
func RegisterDatabaseTools(s *server.MCPServer, cp security.CredentialProvider, g *security.Guard) {
	// ===== 现有迁移（2个）=====

	// DescribeDatabases - 查询数据库列表（只读）
	registerTool(s, cp, g, "DescribeDatabases", "查询实例的数据库列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
			mcp.WithArray("Filters", mcp.Description("过滤条件: database-name")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量，从0开始")),
			mcp.WithNumber("Limit", mcp.Description("单次显示数量")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeDatabasesRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeDatabases(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// CreateDatabase - 创建数据库（L2业务确认）
	registerTool(s, cp, g, "CreateDatabase", "创建数据库",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID，形如postgres-6fego161")),
			mcp.WithString("DatabaseName", mcp.Description("创建的数据库名")),
			mcp.WithString("DatabaseOwner", mcp.Description("数据库的所有者")),
			mcp.WithString("Encoding", mcp.Description("数据库的字符编码")),
			mcp.WithString("Collate", mcp.Description("数据库的排序规则")),
			mcp.WithString("Ctype", mcp.Description("数据库的字符分类")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewCreateDatabaseRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.CreateDatabase(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 新增（2个）=====

	// ModifyDatabaseOwner - 修改数据库属主（L4审计，属主变更审计）
	registerTool(s, cp, g, "ModifyDatabaseOwner", "修改数据库属主",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("DatabaseName", mcp.Required(), mcp.Description("数据库名")),
			mcp.WithString("DatabaseOwner", mcp.Required(), mcp.Description("新属主用户名")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewModifyDatabaseOwnerRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.ModifyDatabaseOwner(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeDatabaseObjects - 查询数据库对象列表（只读）
	registerTool(s, cp, g, "DescribeDatabaseObjects", "查询数据库对象列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("ObjectType", mcp.Required(), mcp.Description("查询的对象类型: database|schema|sequence|procedure|type|function|table|view|matview|column")),
			mcp.WithString("DatabaseName", mcp.Description("查询对象所属的数据库。当查询对象类型不为database时必填")),
			mcp.WithString("SchemaName", mcp.Description("查询对象所属的模式。当查询对象类型不为database、schema时必填")),
			mcp.WithString("TableName", mcp.Description("查询对象所属的表。当查询对象类型为column时必填")),
			mcp.WithNumber("Limit", mcp.Description("单次显示数量，默认20，可选范围[0,100]")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量，从0开始")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeDatabaseObjectsRequest()
			if err := req.FromJsonString(marshalArgs(args)); err != nil {
				return "", err
			}
			rsp, err := client.DescribeDatabaseObjects(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Database tools registered: 4")
}
