package security

import (
	"fmt"
	"os"
	"strings"
)

// GuardLevel 定义 guard 防护级别
type GuardLevel int

const (
	LevelNone     GuardLevel = iota // L0: 无 guard，只读操作直接执行
	LevelFee                        // L1: 费用确认（建实例/扩缩容等花钱操作）
	LevelBusiness                   // L2: 业务确认（重启/隔离等可能中断业务）
	LevelCritical                   // L3: 最高级确认（数据覆盖/恢复等不可逆操作）
	LevelAudit                      // L4: 审计（执行但记录日志）
)

// Guard 安全中间件
//
// 支持"画像（Profile）"机制，方便在开发/测试/生产等环境间快速切换限制策略：
// 同一份 .env 中可以并存多套画像的配置（如 GUARD_DEV_* / GUARD_TEST_*），
// 只需切换 GUARD_PROFILE 这一个变量并重启进程，即可整体切换限制策略，
// 不需要逐项修改/删除 SCOPE_ENABLED、REGION_SCOPE 等值。
// 未设置 GUARD_PROFILE 时，行为与旧版本完全一致（直接读不带前缀的变量）。
type Guard struct {
	Profile      string // 当前生效的画像名（GUARD_PROFILE），空表示未启用画像机制
	ReadOnly     bool   // 默认 true，写操作被拒
	ScopeEnabled bool   // 默认 true，是否启用 region scoping 限制
	RegionScope  string // 限定 region（空=不限）
}

// NewGuard 从环境变量创建 Guard
//
// 画像变量优先级（以 READ_ONLY 为例）：
//  1. GUARD_<PROFILE>_READ_ONLY（画像专属值，PROFILE 会转为大写）
//  2. READ_ONLY（未加前缀的全局值，向后兼容旧配置）
//  3. 内置默认值
func NewGuard() *Guard {
	profile := strings.TrimSpace(os.Getenv("GUARD_PROFILE"))
	return &Guard{
		Profile:      profile,
		ReadOnly:     getProfileBool(profile, "READ_ONLY", true),
		ScopeEnabled: getProfileBool(profile, "SCOPE_ENABLED", true),
		RegionScope:  getProfileEnv(profile, "REGION_SCOPE", ""),
	}
}

// Check 执行 guard 校验
func (g *Guard) Check(toolName string, region string, level GuardLevel) error {
	// 1. 只读模式检查：写操作在 ReadOnly=true 时拒绝
	if g.ReadOnly && level != LevelNone {
		return fmt.Errorf("write operation '%s' is blocked by READ_ONLY mode (profile=%s). Switch to a profile with READ_ONLY=false to enable write operations", toolName, g.profileLabel())
	}

	// 资源 scoping 总开关：关闭时跳过 region 限制检查
	if !g.ScopeEnabled {
		return nil
	}

	// 2. 资源 scoping 检查：region 是否在允许范围
	if g.RegionScope != "" && region != g.RegionScope {
		return fmt.Errorf("region '%s' is outside allowed scope '%s' (profile=%s)", region, g.RegionScope, g.profileLabel())
	}

	return nil
}

// GetGuardWarning 返回对应级别的确认提示
func GetGuardWarning(level GuardLevel, toolName string) string {
	switch level {
	case LevelFee:
		return fmt.Sprintf("Tool '%s' will incur costs. Set confirm=true to proceed.", toolName)
	case LevelBusiness:
		return fmt.Sprintf("Tool '%s' may cause service interruption. Set confirm=true to proceed.", toolName)
	case LevelCritical:
		return fmt.Sprintf("Tool '%s' will overwrite data and is irreversible. Set confirm=true to proceed.", toolName)
	default:
		return fmt.Sprintf("Tool '%s' requires confirmation. Set confirm=true to proceed.", toolName)
	}
}

// NeedsConfirm 判断该级别是否需要确认参数
func NeedsConfirm(level GuardLevel) bool {
	return level == LevelFee || level == LevelBusiness || level == LevelCritical
}

// RegionScopeActive 判断 region scoping 限制当前是否生效（总开关开启且已配置限定 region）
func (g *Guard) RegionScopeActive() bool {
	return g != nil && g.ScopeEnabled && g.RegionScope != ""
}

func (g *Guard) profileLabel() string {
	if g == nil || g.Profile == "" {
		return "default"
	}
	return g.Profile
}

// getProfileEnv 按"画像专属值 -> 全局值 -> 默认值"的优先级读取字符串型配置。
// 画像专属变量名格式为 GUARD_<PROFILE大写>_<KEY>，例如 GUARD_TEST_REGION_SCOPE。
func getProfileEnv(profile, key, fallback string) string {
	if profile != "" {
		if v := os.Getenv(fmt.Sprintf("GUARD_%s_%s", strings.ToUpper(profile), key)); v != "" {
			return v
		}
	}
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getProfileBool 按同样的优先级读取布尔型配置。
func getProfileBool(profile, key string, fallback bool) bool {
	v := getProfileEnv(profile, key, "")
	if v == "" {
		return fallback
	}
	return v == "true" || v == "1"
}
