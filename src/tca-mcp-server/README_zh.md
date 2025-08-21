# 腾讯云代码分析（TCA）MCP Server
> 官网地址(https://tca.tencent.com) 支持 MCP 协议的 MCP Server，用于快速启动分析代码，获取代码分析报告。

腾讯云代码分析（Tencent Code Analysis, TCA）起步于 2012 年（内部代号CodeDog），是集众多代码分析工具的云原生、分布式、高性能的代码综合分析跟踪管理平台，其主要功能是持续跟踪分析代码，观测项目代码质量，支撑团队传承代码文化。了解腾讯云代码助手更多信息，请访问官网使用指南: https://tca.tencent.com/document/zh/guide/


## TCA MCP Server 使用步骤

### 1，在TCA官网上创建相关资源

官网地址: https://tca.tencent.com/

- step1:【创建团队】 访问TCA官网，登录后选择创建一个团队，填写相关信息，等待申请通过：
![创建团队](https://cnb.cool/tca/plugins/tca-mcp-server/-/git/raw/master/docs/images/org.png)

- step2:【创建项目组】 创建团队后，点击选择团队，进入后创建一个项目组：
![创建项目组](https://cnb.cool/tca/plugins/tca-mcp-server/-/git/raw/master/docs/images/project_team.png)

- step3:【接入代码库】 创建项目组后，点击选择项目组，进入后选择接入需要进行分析的代码库：
![接入代码库](https://cnb.cool/tca/plugins/tca-mcp-server/-/git/raw/master/docs/images/repo.png)

- step4:【创建分析项目】 成功接入代码库后，创建一个分析项目(推荐先使用图中的官方体验方案进行使用体验)：
![创建分析项目](https://cnb.cool/tca/plugins/tca-mcp-server/-/git/raw/master/docs/images/project.png)

### 2，在代码库中创建 `tca-mcp.ini` 配置文件

在需要进行代码分析的代码库中创建 **tca-mcp.ini** 配置文件，配置文件存放在代码库**根目录**下，配置文件内容如下：

```ini
[config]
project_id=<project_id>
repo_id=<repo_id>
org_sid=<org_sid>
team_name=<team_name>
```

相关参数可以从对应页面的路由中获取，如下图所示：

![tca-mcp-ini参数](https://cnb.cool/tca/plugins/tca-mcp-server/-/git/raw/master/docs/images/tca-mcp-ini.png)
其中`4iYVpci9nAX` 对应 `org_sid`； `19485` 对应 `repo_id`; `234521` 对应 `project_id`; `first` 对应 `team_name`。按照实际情况填写即可。


### 3, 配置TCA MCP Server
```json
{
  "mcpServers": {
    "tca-mcp-server": {
      "command": "npx",
      "args": ["-y", "-p", "tca-mcp-server@latest", "tca-mcp-stdio"],
      "env": {
        "TCA_TOKEN": "<TCA_TOKEN>", 
        "TCA_USER_NAME": "<TCA_USER_NAME>"
      }
    }
  }
}
```
对应的 `TCA_TOKEN` 和 `TCA_USER_NAME` 从TCA官网, 【个人中心】 -> 【个人令牌
】 中获取，可访问 https://tca.tencent.com/user/token 获取。


## TCA MCP Server 开发步骤
> 要求： nodejs >= 22.0.0

1，npm run build   
2, 手动添加测试配置:
```bash
{
  "mcpServers": {
    "tca-mcp-server-test": {
      "command": "node",
      "args": ["/path/to/tca-mcp-server/dist/stdio.js"],
      "env": {
        "TCA_TOKEN": "<TCA_TOKEN>", 
        "TCA_USER_NAME": "<TCA_USER_NAME>",
      }
    }
  }
}
```
