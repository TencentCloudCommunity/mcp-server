# TAPD MCP Server

TAPD 是腾讯敏捷研发管理平台，覆盖需求、计划、研发、测试、发布研发全生命周期。支持用自然语言与 TAPD 对话，实现需求、缺陷、任务、迭代等管理。

* 与 TAPD API 无缝集成，提升开发效率

## 系统要求

* uv
* TAPD Access Token（推荐）或 TAPD API 账号密码

## 快速开始
### Install uv
```
brew install uv
# OR
curl -LsSf https://astral.sh/uv/install.sh | sh
```

### 参数
- TAPD_ACCESS_TOKEN: 个人令牌（推荐）
- TAPD_API_USER: API账号 
- TAPD_API_PASSWORD: API密钥 
- BOT_URL: 企业微信机器人 webhook 地址，如果需要发送消息到企业微信群才需要填（选填）

注意 TAPD_API_USER/TAPD_API_PASSWORD（TAPD API 账号密码） 与 TAPD_ACCESS_TOKEN（TAPD 个人授权令牌）为两种调用 MCP 的凭据，选一种即可。


### 获取 TAPD Access Token（推荐）
TAPD_ACCESS_TOKEN 为 TAPD 个人令牌，在 [我的设置-个人访问令牌](https://www.tapd.cn/personal_settings/index?tab=personal_token)  点“创建个人访问令牌”，令牌只显示一次，注意保存下来。

### 获取 API 账号密码（兼容）

1. 未注册，请前往 [注册](https://www.tapd.cn?from_partner=copilot&source=tapd_operation_copilot)
2. 已注册但未授权API，请前往API配置：登录TAPD，点击进入 [公司管理-API账号管理](https://www.tapd.cn/open_platform/open_api_redirect?from_partner=copilot&source=tapd_operation_copilot)，复制API账号和API密钥


## IDE 配置
### CodeBuddy / Claude Desktop / Cursor
#### 本地 stdio 运行，环境变量配置示例（推荐个人 Token）
```json
{
  "mcpServers": {
    "mcp-server-tapd": {
      "command": "uvx",
      "args": ["mcp-server-tapd"],
      "env": {
        "TAPD_ACCESS_TOKEN": "您的个人Token",
        // 兼容账号密码方式（如未配置 Token 时生效）
        "TAPD_API_USER": "",
        "TAPD_API_PASSWORD": "",
        "TAPD_API_BASE_URL": "https://api.tapd.cn",
        "TAPD_BASE_URL": "https://www.tapd.cn",
        "BOT_URL": ""
      }
    }
  }
}
```


### Streamable HTTP 配置
使用 Streamable HTTP 替代 stdio
1. 在终端手动启动服务
```
git clone https://cnb.cool/tapd_mcp/mcp-server-tapd.git

cd mcp-server-tapd/src/mcp_server_tapd

uv venv && source .venv/bin/activate

uv pip install requests markdown mcp mcp_server_tapd

python server.py --mode=streamable-http --host="0.0.0.0" --port=8000 --api-user=your_api_user --api-password=your_api_password --api-base-url=https://api.tapd.cn --tapd-base-url=https://www.tapd.cn  --bot-url="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=XXX"

# 如果使用个人 Token，则将 --api-user=your_api_user --api-password=your_api_password 替换为 --access-token=您的个人Token

```

2. 通过如下配置连接到已启动的服务
```json
{
  "mcpServers": {
    "tapd_mcp_http": {
      "url": "http://localhost:8000/mcp/"
    }
  }
}
```
