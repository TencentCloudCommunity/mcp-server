package config

import (
	"os"
	"strings"
)

// Features 定义功能组开关
var features = map[string]bool{
	"instance":   true,  // 实例查询+管理
	"account":    true,  // 账号管理
	"database":   true,  // 数据库管理
	"parameter":  true,  // 参数管理
	"backup":     true,  // 备份恢复
	"monitoring": true,  // 监控/日志
	"network":    true,  // 网络+安全组
	"readonly":   true,  // 只读实例
	"lifecycle":  false, // 高危生命周期操作（默认关）
}

// Init 从环境变量 FEATURES 初始化功能组开关
// 格式: FEATURES=instance,account,database 或 FEATURES=all
func Init() {
	envFeatures := os.Getenv("FEATURES")
	if envFeatures == "" {
		return
	}

	if envFeatures == "all" {
		for k := range features {
			features[k] = true
		}
		return
	}

	// 先全部关闭，再按列表开启
	for k := range features {
		features[k] = false
	}

	for _, f := range strings.Split(envFeatures, ",") {
		f = strings.TrimSpace(f)
		if _, ok := features[f]; ok {
			features[f] = true
		}
	}
}

// IsEnabled 检查功能组是否启用
func IsEnabled(group string) bool {
	enabled, ok := features[group]
	return ok && enabled
}
