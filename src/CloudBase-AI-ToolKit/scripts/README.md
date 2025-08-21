# 脚本使用说明

## sync-config.mjs - 配置同步脚本

这个脚本用于将 `config/` 目录下的配置文件同步到 [awsome-cloudbase-examples](https://github.com/TencentCloudBase/awsome-cloudbase-examples) 项目的各个模板中。

### 🚀 快速开始

1. **准备目标仓库**
   ```bash
   cd ..
   git clone https://github.com/TencentCloudBase/awsome-cloudbase-examples.git
   cd cloudbase-turbo-delploy
   ```

2. **运行同步脚本**
   ```bash
   # 同步所有模板
   node scripts/sync-config.mjs
   
   # 干运行模式（预览操作，不实际执行）
   node scripts/sync-config.mjs --dry-run
   ```

### 📋 命令选项

| 选项 | 说明 |
|-----|------|
| `--help`, `-h` | 显示帮助信息 |
| `--dry-run` | 干运行模式，只显示将要执行的操作 |
| `--skip-git` | 跳过 Git 提交和推送操作 |
| `--backup` | 创建备份（覆盖配置文件设置） |
| `--filter <关键词>` | 只同步路径包含指定关键词的模板 |

### 💡 使用示例

```bash
# 同步所有模板
node scripts/sync-config.mjs

# 只同步包含 "web" 的模板
node scripts/sync-config.mjs --filter web

# 只同步小程序模板
node scripts/sync-config.mjs --filter miniprogram

# 干运行模式，查看将要同步的内容
node scripts/sync-config.mjs --dry-run

# 同步但跳过 Git 操作
node scripts/sync-config.mjs --skip-git

# 创建备份并同步
node scripts/sync-config.mjs --backup
```

### ⚙️ 配置文件

#### template-config.json

脚本的主要配置文件，结构简单明了：

```json
{
  "description": "awsome-cloudbase-examples 项目模板路径配置",
  "version": "1.0.0",
  "templates": [
    "web/react-starter",
    "miniprogram/blog",
    "cloudbaserun/deepseek-agent",
    "..."
  ],
  "excludePatterns": ["node_modules", ".DS_Store", "*.log"],
  "syncConfig": {
    "createBackup": false,
    "overwrite": true,
    "autoCommit": true,
    "commitMessage": "chore: sync config and rules from cloudbase-turbo-deploy"
  }
}
```

**配置说明：**
- `templates`: 模板路径数组，每个路径对应 awsome-cloudbase-examples 中的一个模板目录
- `excludePatterns`: 同步时需要排除的文件/目录模式
- `syncConfig`: 同步相关配置
  - `createBackup`: 是否创建备份
  - `overwrite`: 是否覆盖已存在的文件
  - `autoCommit`: 是否自动提交和推送
  - `commitMessage`: 默认提交信息

### 🔧 工作流程

1. **读取配置**: 从 `template-config.json` 读取模板列表和配置
2. **过滤模板**: 根据 `--filter` 参数过滤要同步的模板
3. **检查环境**: 验证 `config/` 目录和目标仓库是否存在
4. **同步文件**: 将 `config/` 下的内容复制到各个模板目录
5. **Git 操作**: 自动提交更改并推送到远程仓库

### 📁 同步的内容

脚本会将 `config/` 目录下的所有内容同步到目标模板，包括：

- `.github/` - GitHub Actions 配置
- `.rules/` - 规则配置文件
- `.cursor/` - Cursor 编辑器配置
- `.windsurf/` - Windsurf 编辑器配置  
- `.vscode/` - VSCode 配置
- `.clinerules/` - CLI 规则配置
- `.trae/` - Trae 配置

### 🔍 过滤功能

使用 `--filter` 参数可以只同步特定类型的模板：

```bash
# 同步所有 web 相关模板
node scripts/sync-config.mjs --filter web

# 同步所有小程序模板  
node scripts/sync-config.mjs --filter miniprogram

# 同步所有云函数模板
node scripts/sync-config.mjs --filter cloudrunfunctions

# 同步所有 CloudBase Run 模板
node scripts/sync-config.mjs --filter cloudbaserun
```

### ⚠️ 注意事项

1. **权限要求**: 确保你有目标仓库的推送权限
2. **备份建议**: 首次使用建议启用 `--backup` 选项
3. **干运行**: 不确定操作结果时，先使用 `--dry-run` 预览
4. **路径结构**: 目标仓库必须位于 `../awsome-cloudbase-examples/`

### 🐛 故障排除

**问题**: 目标目录不存在
- **解决**: 确保已正确克隆目标仓库，且路径正确

**问题**: Git 操作失败  
- **解决**: 检查网络连接和仓库权限，可使用 `--skip-git` 跳过 Git 操作

**问题**: 配置文件读取失败
- **解决**: 检查 `template-config.json` 格式是否正确

### 📝 日志输出

脚本会输出详细的操作日志，包括：
- 📁 源目录信息
- 📋 要处理的模板数量  
- 🔍 过滤条件（如果使用）
- ✅ 成功同步的文件
- ⚠️ 跳过的文件和原因
- 📊 最终统计信息 