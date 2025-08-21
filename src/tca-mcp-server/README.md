# Tencent Cloud Code Analysis (TCA) MCP Server​
> Official website (https://tca.tencent.com) MCP Server supporting MCP protocol for quickly starting code analysis and obtaining code analysis reports.

Tencent Cloud Code Analysis (TCA), which started in 2012 (internal code name: CodeDog), is a cloud-native, distributed, and high-performance comprehensive code analysis and tracking management platform integrating numerous code analysis tools. Its main functions are to continuously track and analyze code, observe project code quality, and support teams in inheriting code culture. For more information about Tencent Cloud Code Assistant, please visit the official website usage guide: https://tca.tencent.com/document/zh/guide/.

## TCA MCP Server Usage Steps​

### 1，Create relevant resources on the TCA official website

Official website: https://tca.tencent.com/​

- step1: [Create a team] Visit the TCA official website, log in, select to create a team, fill in relevant information, and wait for the application to be approved:
![create_team](https://cnb.cool/tca/plugins/tca-mcp-server/-/git/raw/master/docs/images/org.png)

- step2: [Create a project team] After creating the team, click to select the team, and create a project team after entering:
![create a project team](https://cnb.cool/tca/plugins/tca-mcp-server/-/git/raw/master/docs/images/project_team.png)

- step3: [Access the code repository] After creating the project team, click to select the project team, and select to access the code repository that needs to be analyzed after entering:
![repo](https://cnb.cool/tca/plugins/tca-mcp-server/-/git/raw/master/docs/images/repo.png)

- step4: [Create an analysis project] After successfully accessing the code repository, create an analysis project (it is recommended to first use the official experience plan in the figure for usage experience):
![create a project](https://cnb.cool/tca/plugins/tca-mcp-server/-/git/raw/master/docs/images/project.png)

### 2, Create a `tca-mcp.ini` configuration file in the code repository

Create a tca-mcp.ini configuration file in the code repository that needs code analysis. The configuration file is stored in the root directory of the code repository, and the content of the configuration file is as follows:

```ini
[config]
project_id=<project_id>
repo_id=<repo_id>
org_sid=<org_sid>
team_name=<team_name>
```

Relevant parameters can be obtained from the route of the corresponding page, as shown in the following figure:

![tca-mcp-ini参数](https://cnb.cool/tca/plugins/tca-mcp-server/-/git/raw/master/docs/images/tca-mcp-ini.png)

Where `4iYVpci9nAX` corresponds to `org_sid`; `19485` corresponds to `repo_id`; `234521` corresponds to `project_id`; `first` corresponds to `team_name`. Fill in according to the actual situation.


### 3, Configure TCA MCP Server
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
The corresponding TCA_TOKEN and TCA_USER_NAME are obtained from the TCA official website, [Personal Center] -> [Personal Token], and can be accessed at https://tca.tencent.com/user/token.

## TCA MCP Server Development Steps​
> Requirements: nodejs >= 22.0.0

1，npm run build
2, Manually add test configuration:
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
