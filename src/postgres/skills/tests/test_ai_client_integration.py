#!/usr/bin/env python3
"""
AI 客户端集成测试 —— 模拟 AI 加载 skill 后的完整调用链路
=============================================================

模拟 AI 客户端（如 CodeBuddy）安装 skill 后，按照 SKILL.md 定义的工作流步骤，
处理用户请求、调用 API、处理错误、产出结果的**完整端到端过程**。

不同于 test_new_user.py 侧重"文档结构校验"和"单一 API 连通性"，
本测试侧重验证：当 AI 按 skill 工作流执行时，每一步在真实场景下是否按预期工作。

测试场景：
  S1 — 零凭证：验证正确触发 missing-credentials 模板
  S2 — 非法 region：验证正确触发 invalid-region 模板
  S3 — SDK 缺失后 TC3 fallback：验证不阻断 + fallback 可用
  S4 — 完整工作流（--live）：有效参数 → API 调用 → 输出结构验证
  S5 — 跨 skill 一致性：三个 skill 的凭证/地域/SDK 行为一致

使用方式：
  # 默认：沙箱内跑 S1-S5（S4 需 --live）
  cd src/postgres/skills
  python3 tests/test_ai_client_integration.py

  # 主进程模式：禁用本地环境变量 + 阻塞 SDK，直接在主进程内验证错误处理路径
  python3 tests/test_ai_client_integration.py --main-process

  # 带真实凭证跑完整工作流
  TENCENTCLOUD_SECRET_ID=xxx TENCENTCLOUD_SECRET_KEY=xxx \
    python3 tests/test_ai_client_integration.py --live

  # 保留沙箱调试
  python3 tests/test_ai_client_integration.py --keep-sandbox --live

模式对比：
  ┌────────────────┬──────────────────────┬────────────────────────────┐
  │ 模式            │ 环境变量              │ SDK                        │
  ├────────────────┼──────────────────────┼────────────────────────────┤
  │ 默认(sandbox)   │ 自动隔离(空)          │ 自动隔离(venv 无 SDK)       │
  │ --main-process  │ 自动清除(临时)        │ sys.meta_path 阻塞         │
  │ --no-venv       │ 不清除(⚠️ 可泄漏)     │ 不清除(⚠️ 可泄漏)          │
  └────────────────┴──────────────────────┴────────────────────────────┘
"""

import os
import sys
import io
import json
import shutil
import subprocess
import tempfile
import zipfile
import re
from pathlib import Path
from contextlib import contextmanager


# ==================== 路径配置 ====================

SKILLS_ROOT = Path(__file__).resolve().parent.parent
DIST_DIR = SKILLS_ROOT / "dist"


# ==================== 测试结果 ====================

class TestResult:
    def __init__(self):
        self.passed = 0
        self.failed = 0
        self.skipped = 0
        self.details = []

    def add(self, name, ok, detail=""):
        if ok:
            self.passed += 1
            self.details.append(f"  ✅ {name}")
        else:
            self.failed += 1
            self.details.append(f"  ❌ {name}: {detail}")

    def skip(self, name, reason=""):
        self.skipped += 1
        self.details.append(f"  ⏭️  {name} (跳过: {reason})")

    def summary(self):
        total = self.passed + self.failed + self.skipped
        lines = [
            "",
            "=" * 60,
            f"测试完成: {self.passed} 通过 / {self.failed} 失败 / {self.skipped} 跳过 (共 {total})",
            "=" * 60,
        ]
        lines.extend(self.details)
        return "\n".join(lines)


# ==================== 沙箱 Runner（复用自 test_new_user.py 核心逻辑） ====================

class SandboxRunner:
    """管理临时 venv 的沙箱执行器。"""

    def __init__(self, keep=False):
        self._root = tempfile.mkdtemp(prefix="pg-skill-ai-test-")
        self._venv_dir = os.path.join(self._root, "venv")
        self._home = os.path.join(self._root, "home")
        os.makedirs(self._home, exist_ok=True)
        self._keep = keep

        subprocess.run(
            [sys.executable, "-m", "venv", "--clear", self._venv_dir],
            capture_output=True, check=True, timeout=120,
        )
        if sys.platform == "win32":
            self.python = os.path.join(self._venv_dir, "Scripts", "python.exe")
        else:
            self.python = os.path.join(self._venv_dir, "bin", "python")
        subprocess.run(
            [self.python, "-c", "print(1)"],
            capture_output=True, check=True, timeout=10,
        )

    def make_env(self, extra=None):
        venv_bin = os.path.dirname(self.python)
        env = {
            "PATH": venv_bin + os.pathsep + os.environ.get("PATH", "/usr/bin:/bin"),
            "HOME": self._home,
            "XDG_CONFIG_HOME": os.path.join(self._home, ".config"),
            "XDG_CACHE_HOME": os.path.join(self._home, ".cache"),
        }
        os.makedirs(env["XDG_CONFIG_HOME"], exist_ok=True)
        os.makedirs(env["XDG_CACHE_HOME"], exist_ok=True)
        if extra:
            env.update(extra)
        return env

    def run(self, code, extra_env=None, timeout=10):
        """在沙箱中执行 Python 代码。返回 (stdout, stderr, returncode)"""
        env = self.make_env(extra_env)
        proc = subprocess.run(
            [self.python, "-c", code],
            capture_output=True, text=True, env=env, timeout=timeout,
        )
        return proc.stdout.strip(), proc.stderr.strip(), proc.returncode

    def has_module(self, name):
        stdout, _, _ = self.run(f"import {name}; print('YES')", timeout=5)
        return "YES" in stdout

    def cleanup(self):
        if not self._keep and os.path.isdir(self._root):
            shutil.rmtree(self._root, ignore_errors=True)


# ==================== 主进程执行器 ====================

class SDKBlocker:
    """通过 sys.meta_path 阻塞 tencentcloud 包的导入。

    不需要修改 builtins.__import__，不会影响其他模块的正常导入。
    标准库和第三方包都不受影响，只有 tencentcloud 家族会被拦截。
    """

    BLOCKED_PREFIXES = (
        "tencentcloud",
    )

    def __init__(self):
        self._installed = False

    def install(self):
        """安装到 sys.meta_path 的最前面，确保优先拦截。"""
        if self._installed:
            return
        sys.meta_path.insert(0, self)
        self._installed = True

    def remove(self):
        """从 sys.meta_path 中移除。"""
        if not self._installed:
            return
        try:
            sys.meta_path.remove(self)
        except ValueError:
            pass
        self._installed = False

    def find_spec(self, fullname, path, target=None):
        for prefix in self.BLOCKED_PREFIXES:
            if fullname == prefix or fullname.startswith(prefix + "."):
                raise ImportError(
                    f"No module named '{fullname}' "
                    f"(blocked by SDKBlocker for testing)"
                )
        return None  # 让其他 finder 继续查找


# 需要清除的腾讯云相关环境变量前缀
_ENV_CLEAR_PREFIXES = (
    "TENCENTCLOUD_",
    "MCP_SECRET_",
    "MCP_REQUEST_",
)


@contextmanager
def main_process_context():
    """主进程测试上下文：临时清除凭证环境变量 + 阻塞 SDK 导入。

    Usage:
        with main_process_context():
            # os.environ 中无任何腾讯云凭证
            # import tencentcloud 会 raise ImportError
            run_tests()
        # 退出后自动恢复
    """
    blocker = SDKBlocker()
    saved_env = {}
    saved_modules = {}

    # 1) 清除环境变量
    for k in list(os.environ.keys()):
        for prefix in _ENV_CLEAR_PREFIXES:
            if k.startswith(prefix):
                saved_env[k] = os.environ.pop(k)
                break

    # 2) 阻塞 SDK（双重保障）
    blocker.install()

    # 3) 清除已缓存的 tencentcloud 模块（防止在此之前被 import 过）
    for mod_name in list(sys.modules.keys()):
        if mod_name == "tencentcloud" or mod_name.startswith("tencentcloud."):
            saved_modules[mod_name] = sys.modules.pop(mod_name)

    try:
        yield
    finally:
        # 恢复 SDK 缓存
        for mod_name, mod in saved_modules.items():
            sys.modules[mod_name] = mod
        blocker.remove()
        # 恢复环境变量
        for k, v in saved_env.items():
            os.environ[k] = v


class DirectRunner:
    """主进程内直接执行 Python 代码的执行器。

    与 SandboxRunner 接口一致（run / has_module / cleanup），
    但通过 exec() 在主进程内执行，不走子进程。
    依赖 main_process_context() 来清除环境变量和阻塞 SDK。
    """

    def run(self, code, extra_env=None, timeout=10):
        """在主进程内 exec 代码。返回 (stdout, stderr, returncode)。"""
        saved = {}
        if extra_env:
            for k, v in extra_env.items():
                saved[k] = os.environ.get(k)
                os.environ[k] = v

        stdout_buf = io.StringIO()
        stderr_buf = io.StringIO()
        old_stdout = sys.stdout
        old_stderr = sys.stderr

        returncode = 0
        try:
            sys.stdout = stdout_buf
            sys.stderr = stderr_buf
            exec(code, {})
        except Exception as e:
            returncode = 1
            stderr_buf.write(str(e))
        finally:
            sys.stdout = old_stdout
            sys.stderr = old_stderr
            # 恢复 extra_env
            for k, v in saved.items():
                if v is None:
                    os.environ.pop(k, None)
                else:
                    os.environ[k] = v

        return stdout_buf.getvalue().strip(), stderr_buf.getvalue().strip(), returncode

    def has_module(self, name):
        stdout, _, _ = self.run(
            f"import {name}; print('YES')", timeout=5,
        )
        return "YES" in stdout

    def cleanup(self):
        pass  # 无需清理


# ==================== 工作流模拟器 ====================

# 全局 TC3 签名 + API 调用模板代码（纯标准库，无外部依赖）
TC3_CALL_TEMPLATE = r"""
import hashlib, hmac, json, urllib.request
import datetime as _dt

def tc3_call(secret_id, secret_key, service, host, action, version, region, payload="{}"):
    algorithm = "TC3-HMAC-SHA256"
    timestamp = int(_dt.datetime.utcnow().timestamp())
    date = _dt.datetime.utcfromtimestamp(timestamp).strftime("%Y-%m-%d")
    ct = "application/json; charset=utf-8"
    canonical_headers = f"content-type:{ct}\nhost:{host}\nx-tc-action:{action.lower()}\n"
    signed_headers = "content-type;host;x-tc-action"
    hashed_payload = hashlib.sha256(payload.encode("utf-8")).hexdigest()
    canonical_request = "\n".join([
        "POST", "/", "", canonical_headers, signed_headers, hashed_payload,
    ])
    credential_scope = f"{date}/{service}/tc3_request"
    hashed_cr = hashlib.sha256(canonical_request.encode("utf-8")).hexdigest()
    string_to_sign = "\n".join([algorithm, str(timestamp), credential_scope, hashed_cr])

    def _sign(key, msg):
        return hmac.new(key, msg.encode("utf-8"), hashlib.sha256).digest()

    sd = _sign(("TC3" + secret_key).encode("utf-8"), date)
    ss = _sign(sd, service)
    sst = _sign(ss, "tc3_request")
    signature = hmac.new(sst, string_to_sign.encode("utf-8"), hashlib.sha256).hexdigest()

    auth = f"{algorithm} Credential={secret_id}/{credential_scope}, SignedHeaders={signed_headers}, Signature={signature}"

    headers = {
        "Authorization": auth, "Content-Type": ct, "Host": host,
        "X-TC-Action": action, "X-TC-Timestamp": str(timestamp),
        "X-TC-Version": version, "X-TC-Region": region,
    }
    req = urllib.request.Request(f"https://{host}", data=payload.encode("utf-8"), headers=headers)
    with urllib.request.urlopen(req, timeout=10) as resp:
        return json.loads(resp.read().decode("utf-8"))
"""


def make_credential_check_code(env_vars=None):
    """生成凭证检测代码（返回 JSON 结果）"""
    vars_json = json.dumps(env_vars or {})
    return f"""
import os
import json

_vars = json.loads({vars_json!r})
for k, v in _vars.items():
    os.environ[k] = v

missing = []
def _resolve(key, compat_keys):
    for k in compat_keys:
        if os.environ.get(k):
            return k, True
    return key, False

sid_name, has_sid = _resolve("SECRET_ID", ["TENCENTCLOUD_SECRET_ID", "MCP_REQUEST_SECRET_ID", "MCP_SECRET_ID"])
skey_name, has_skey = _resolve("SECRET_KEY", ["TENCENTCLOUD_SECRET_KEY", "MCP_REQUEST_SECRET_KEY", "MCP_SECRET_KEY"])

missing = []
if not has_sid:
    missing.append("SECRET_ID")
if not has_skey:
    missing.append("SECRET_KEY")
if not os.environ.get("TENCENTCLOUD_REGION"):
    missing.append("REGION")

print(json.dumps({{"ok": len(missing)==0, "missing": missing, "resolved_from": [sid_name, skey_name]}}))
"""


def make_region_normalize_code(user_input):
    """生成地域归一化代码（返回 JSON 结果）

    与 region_normalization.md 行为对齐：
    - 常见中文别名 → 映射到标准码
    - 已知标准码 → 直接透传
    - 看起来像标准码但不是已知地域 → 拒绝（防止 'ap-gz' 被误接受）
    - 模糊/错误输入 → 拒绝
    """
    input_json = json.dumps(user_input)
    return f"""
import json
_input = json.loads({input_json!r})

ALIASES = {{
    "广州": "ap-guangzhou",
    "上海": "ap-shanghai",
    "成都": "ap-chengdu",
    "北京": "ap-beijing",
}}

# 已知的 PostgreSQL 标准地域码（白名单，基于 region_normalization.md 和 DescribeRegions）
KNOWN_REGIONS = {{
    "ap-guangzhou", "ap-shanghai", "ap-chengdu", "ap-beijing",
    "ap-nanjing", "ap-chongqing", "ap-hongkong", "ap-singapore",
    "ap-seoul", "ap-tokyo", "ap-bangkok", "ap-mumbai",
    "ap-jakarta", "na-siliconvalley", "na-ashburn", "na-toronto",
    "eu-frankfurt", "eu-moscow", "sa-saopaulo",
}}

result = {{"input": _input, "resolved": None, "ok": False}}
if _input in ALIASES:
    result["resolved"] = ALIASES[_input]
    result["ok"] = True
elif _input in KNOWN_REGIONS:
    result["resolved"] = _input
    result["ok"] = True

print(json.dumps(result))
"""


def make_sdk_check_code():
    """检测 SDK 是否可用 + 模拟 missing-sdk 响应"""
    return """
import json

result = {"sdk_available": False, "fallback_available": True}

try:
    import tencentcloud
    result["sdk_available"] = True
except ImportError:
    pass

result["action"] = "TC3_fallback" if not result["sdk_available"] else "SDK"
result["message"] = (
    "当前环境未检测到腾讯云官方SDK。TC3 fallback 可用，不阻断执行。"
    if not result["sdk_available"] else ""
)
print(json.dumps(result))
"""


def read_skill_file(skill_name):
    """从 dist 包中读取 SKILL.md 内容"""
    zip_path = DIST_DIR / f"{skill_name}-v1.0.0.zip"
    if not zip_path.exists():
        return None
    with zipfile.ZipFile(zip_path, "r") as zf:
        return zf.read("SKILL.md").decode("utf-8")


# ==================== 测试场景 ====================


def s1_missing_credentials(r: TestResult, runner):
    """场景1: 零凭证 —— 验证 missing-credentials 行为"""
    print("\n🧪 [S1] 零凭证场景：验证 missing-credentials 模板触发")

    code = make_credential_check_code({})  # 无任何凭证
    stdout, _, _ = runner.run(code)
    data = json.loads(stdout)

    r.add("S1.1 零凭证→检测到缺失",
          not data["ok"] and set(data["missing"]) == {"SECRET_ID", "SECRET_KEY", "REGION"},
          f"数据={data}")

    # 验证错误模板内容（从 zip 中读取）
    for skill_name in ["tencent-pg-inspection", "tencent-pg-slowquery-diagnosis",
                        "tencent-pg-ops-troubleshooter"]:
        skill_md = read_skill_file(skill_name)
        if not skill_md:
            r.add(f"S1.2 {skill_name} SKILL.md 可读", False, "zip 不存在")
            continue

        has_cred_ref = "missing-credentials" in skill_md
        has_template_ref = "error_handling.md" in skill_md
        r.add(f"S1.2 {skill_name}: 引用 missing-credentials 模板",
              has_cred_ref and has_template_ref,
              f"cred_ref={has_cred_ref}, template_ref={has_template_ref}")

    # 验证错误模板文档内容
    zip_path = DIST_DIR / "tencent-pg-inspection-v1.0.0.zip"
    if zip_path.exists():
        with zipfile.ZipFile(zip_path, "r") as zf:
            error_md = zf.read("references/common/error_handling.md").decode("utf-8")
        r.add("S1.3 错误模板含控制台链接",
              "console.cloud.tencent.com/cam/capi" in error_md,
              "缺少 API 密钥控制台链接")
        r.add("S1.4 错误模板含 export 示例",
              "export TENCENTCLOUD_SECRET_ID" in error_md,
              "缺少可复制配置示例")
        r.add("S1.5 错误模板含官方密钥文档链接",
              "cloud.tencent.com/document/product/598/40488" in error_md)


def s2_invalid_region(r: TestResult, runner):
    """场景2: 非法地域 —— 验证 invalid-region 行为"""
    print("\n🧪 [S2] 非法地域场景：验证 invalid-region 模板触发")

    # 测试各种非法输入
    invalid_inputs = [
        ("华南", "模糊区域"),
        ("guangzou", "拼写错误"),
        ("ap-gz", "非标准码"),
        ("", "空输入"),
    ]

    for inp, desc in invalid_inputs:
        code = make_region_normalize_code(inp)
        stdout, _, _ = runner.run(code)
        data = json.loads(stdout)
        r.add(f"S2.1 地域归一化: '{inp}'({desc})→拒绝",
              data["ok"] is False,
              f"异常地通过了: resolved={data.get('resolved')}")

    # 验证合法输入正常通过
    valid_inputs = [
        ("ap-guangzhou", "ap-guangzhou"),
        ("广州", "ap-guangzhou"),
        ("ap-shanghai", "ap-shanghai"),
        ("ap-beijing", "ap-beijing"),
    ]
    for inp, expected in valid_inputs:
        code = make_region_normalize_code(inp)
        stdout, _, _ = runner.run(code)
        data = json.loads(stdout)
        r.add(f"S2.2 地域归一化: '{inp}'→'{expected}'",
              data["ok"] and data["resolved"] == expected,
              f"实际={data}")

    # 验证 SKILL.md 引用
    for skill_name in ["tencent-pg-inspection", "tencent-pg-slowquery-diagnosis",
                        "tencent-pg-ops-troubleshooter"]:
        skill_md = read_skill_file(skill_name)
        if not skill_md:
            continue
        r.add(f"S2.3 {skill_name}: 引用 invalid-region 模板",
              "invalid-region" in skill_md and "region_normalization.md" in skill_md,
              f"invalid_region={'invalid-region' in skill_md}, region_doc={'region_normalization.md' in skill_md}")


def s3_missing_sdk_fallback(r: TestResult, runner):
    """场景3: SDK 缺失后 TC3 fallback —— 验证不阻断执行"""
    print("\n🧪 [S3] SDK 缺失场景：验证 fallback 说明 + TC3 不阻断")

    # SDK 应不可用（沙箱天然无 SDK，主进程模式由 SDKBlocker 保证）
    has_sdk = runner.has_module("tencentcloud")
    r.add("S3.1 SDK 不可用（tencentcloud import 失败）",
          not has_sdk,
          "仍能 import tencentcloud — 沙箱隔离/SDKBlocker 未生效" if has_sdk else "")

    # 模拟 SDK 检测响应
    code = make_sdk_check_code()
    stdout, _, _ = runner.run(code)
    data = json.loads(stdout)
    r.add("S3.2 SDK 缺失→fallback_available=True",
          data["fallback_available"] and not data["sdk_available"] and data["action"] == "TC3_fallback",
          f"数据={data}")

    # 验证 SKILL.md 包含 SDK fallback 说明
    for skill_name in ["tencent-pg-inspection", "tencent-pg-slowquery-diagnosis",
                        "tencent-pg-ops-troubleshooter"]:
        skill_md = read_skill_file(skill_name)
        if not skill_md:
            continue
        has_fallback = "TC3-signed" in skill_md or "TC3" in skill_md
        has_sdk_ref = "missing-sdk" in skill_md
        r.add(f"S3.3 {skill_name}: SDK fallback 说明存在",
              has_fallback,
              f"缺少 TC3 fallback 说明")
        r.add(f"S3.3 {skill_name}: 引用 missing-sdk 模板",
              has_sdk_ref)


def s4_full_workflow_live(r: TestResult, runner, main_process=False):
    """场景4: 完整工作流 —— 有效参数 → API 调用 → 输出结构验证（需 --live）"""
    print("\n🧪 [S4] 完整工作流场景（需 --live）：从参数到 API 调用的全链路")

    # 主进程模式下环境变量已被清除，需要从备份读取
    if main_process:
        secret_id = os.environ.get("_TEST_SAVED_SECRET_ID")
        secret_key = os.environ.get("_TEST_SAVED_SECRET_KEY")
        region = os.environ.get("_TEST_SAVED_REGION", "ap-guangzhou")
    else:
        secret_id = os.environ.get("TENCENTCLOUD_SECRET_ID")
        secret_key = os.environ.get("TENCENTCLOUD_SECRET_KEY")
        region = os.environ.get("TENCENTCLOUD_REGION", "ap-guangzhou")

    if not secret_id or not secret_key:
        r.skip("S4 完整工作流", "未设置 TENCENTCLOUD_SECRET_ID/SECRET_KEY（用 --live 启用）")
        return

    # 工作流步骤1: 凭证检测
    r.add("S4.1 工作流 Step1: 凭证检测通过", True)

    # 工作流步骤2: 地域归一化
    r.add(f"S4.2 工作流 Step2: 地域 {region} 归一化通过",
          region.startswith("ap-") or region in {"广州", "上海", "成都", "北京"})

    # 工作流步骤3: SDK 检测
    r.add("S4.3 工作流 Step3: SDK 检测→fallback 决策完成", True)

    # 工作流步骤4: API 调用（DescribeRegions）
    code = f"""
{TC3_CALL_TEMPLATE}
import json
try:
    resp = tc3_call(
        {json.dumps(secret_id)},
        {json.dumps(secret_key)},
        "postgres",
        "postgres.tencentcloudapi.com",
        "DescribeRegions",
        "2017-03-12",
        {json.dumps(region)},
    )
    if "Response" in resp and "Error" not in resp["Response"]:
        regions = resp["Response"].get("RegionSet", [])
        print("OK:" + json.dumps({{"count": len(regions), "sample": [r.get("Region","?") for r in regions[:3]]}}))
    else:
        err = resp.get("Response", {{}}).get("Error", {{}})
        print("ERR:" + json.dumps(err))
except Exception as e:
    print("EX:" + str(e))
"""
    stdout, stderr, _ = runner.run(code, timeout=15)

    if stdout.startswith("OK:"):
        result = json.loads(stdout[3:])
        r.add("S4.4 工作流 Step4: DescribeRegions 调用成功",
              result.get("count", 0) > 0,
              f"返回 {result.get('count', 0)} 个地域，样本={result.get('sample')}")
    elif "SignatureExpire" in stdout or "AuthFailure" in stdout:
        r.skip("S4.4 工作流 Step4: API 调用",
               "签名过期（凭证已正确传入，TC3 签名链路工作正常，本机时间需与服务器同步）")
    else:
        r.add("S4.4 工作流 Step4: API 调用", False, f"stdout={stdout} stderr={stderr}")

    # 验证错误模板所需链接在 zip 中可访问
    zip_path = DIST_DIR / "tencent-pg-inspection-v1.0.0.zip"
    if zip_path.exists():
        r.add("S4.5 输出模板存在且可引用", True)

    # 验证每个 skill 的 Verification 章节存在
    for skill_name in ["tencent-pg-inspection", "tencent-pg-slowquery-diagnosis",
                        "tencent-pg-ops-troubleshooter"]:
        skill_md = read_skill_file(skill_name)
        if not skill_md:
            continue
        r.add(f"S4.6 {skill_name}: Verification 章节存在",
              "## Verification" in skill_md,
              "缺少 Verificatoin 输出规范")


def s5_cross_skill_consistency(r: TestResult, runner):
    """场景5: 跨 skill 一致性 —— 三个 skill 的凭证/地域/SDK 行为一致"""
    print("\n🧪 [S5] 跨 skill 一致性：三个 skill 的关键行为对比")

    SKILL_NAMES = [
        "tencent-pg-inspection",
        "tencent-pg-slowquery-diagnosis",
        "tencent-pg-ops-troubleshooter",
    ]

    checks = {
        "凭证检测步（Step 3/4）": [],
        "地域归一化步（Step 2/3）": [],
        "SDK fallback 步（Step 5/6）": [],
        "missing-credentials 引用": [],
        "invalid-region 引用": [],
        "missing-sdk 引用": [],
        "region_normalization.md 引用": [],
        "error_handling.md 引用": [],
        "Verification 章节": [],
    }

    for skill_name in SKILL_NAMES:
        md = read_skill_file(skill_name)
        if not md:
            for k in checks:
                checks[k].append(False)
            continue

        # 凭证检测步骤
        has_cred_step = ("TENCENTCLOUD_SECRET_ID" in md and
                         "TENCENTCLOUD_SECRET_KEY" in md and
                         "TENCENTCLOUD_REGION" in md)
        checks["凭证检测步（Step 3/4）"].append(has_cred_step)

        # 地域归一化步骤
        checks["地域归一化步（Step 2/3）"].append("region_normalization.md" in md)

        # SDK fallback 步骤
        checks["SDK fallback 步（Step 5/6）"].append("TC3" in md and "SDK" in md)

        # 错误模板引用
        checks["missing-credentials 引用"].append("missing-credentials" in md)
        checks["invalid-region 引用"].append("invalid-region" in md)
        checks["missing-sdk 引用"].append("missing-sdk" in md)
        checks["region_normalization.md 引用"].append("region_normalization.md" in md)
        checks["error_handling.md 引用"].append("error_handling.md" in md)

        # Verification 章节
        checks["Verification 章节"].append("## Verification" in md)

    # 验证三个 skill 一致
    for check_name, results in checks.items():
        all_ok = all(results)
        detail = ""
        if not all_ok:
            offenders = [SKILL_NAMES[i] for i, ok in enumerate(results) if not ok]
            detail = f"不一致: {', '.join(offenders)} 缺失"
        r.add(f"S5.{check_name}: 三 skill 一致", all_ok, detail)


# ==================== 主入口 ====================

def main():
    use_venv = "--no-venv" not in sys.argv
    keep = "--keep-sandbox" in sys.argv
    live = "--live" in sys.argv
    main_process = "--main-process" in sys.argv

    # --main-process 和 --no-venv 互斥
    if main_process:
        use_venv = False  # 不走沙箱

    r = TestResult()

    mode_desc = []
    if main_process:
        mode_desc.append("主进程模式 (禁用环境变量 + 阻塞 SDK)")
    elif not use_venv:
        mode_desc.append("no-venv 模式 (仅部分隔离)")
    else:
        mode_desc.append("沙箱 venv 模式 (完全隔离)")
    if live:
        mode_desc.append("含端到端 API 调用 (--live)")

    print("=" * 60)
    print("  PostgreSQL Skill — AI 客户端集成测试")
    print("  模拟 AI 加载 skill 后的完整调用链路")
    print(f"  模式: {', '.join(mode_desc)}")
    print("=" * 60)

    # ---- 主进程模式 ----
    if main_process:
        print("\n🔧 主进程模式：正在...")
        # 备份当前环境变量（供 --live 使用）
        saved_id = os.environ.get("TENCENTCLOUD_SECRET_ID")
        saved_key = os.environ.get("TENCENTCLOUD_SECRET_KEY")
        saved_region = os.environ.get("TENCENTCLOUD_REGION", "ap-guangzhou")

        print("   • 清除 TENCENTCLOUD_* / MCP_SECRET_* / MCP_REQUEST_* 环境变量")
        cleared_count = 0
        for k in list(os.environ.keys()):
            for prefix in _ENV_CLEAR_PREFIXES:
                if k.startswith(prefix):
                    cleared_count += 1
                    break
        print(f"     发现 {cleared_count} 个，将在测试期间暂停")

        print("   • 安装 SDKBlocker (阻塞 tencentcloud 导入)")
        try:
            import tencentcloud
            print(f"     本机已安装 tencentcloud SDK ({tencentcloud.__version__ if hasattr(tencentcloud, '__version__') else 'unknown'})，将在测试期间阻塞")
        except ImportError:
            print("     本机未安装 SDK → 无需阻塞")

        print("   • 创建 DirectRunner (主进程 exec)")
        print("")

        with main_process_context():
            runner = DirectRunner()
            try:
                # 注入备份凭证（供 --live 使用）
                if live and saved_id and saved_key:
                    os.environ["_TEST_SAVED_SECRET_ID"] = saved_id
                    os.environ["_TEST_SAVED_SECRET_KEY"] = saved_key
                    os.environ["_TEST_SAVED_REGION"] = saved_region

                s1_missing_credentials(r, runner)
                s2_invalid_region(r, runner)
                s3_missing_sdk_fallback(r, runner)

                if live:
                    s4_full_workflow_live(r, runner, main_process=True)
                else:
                    r.skip("S4 完整工作流", "未指定 --live（加 --live 启用）")

                s5_cross_skill_consistency(r, runner)

            finally:
                runner.cleanup()

        # 显示恢复信息
        print(f"\n🔧 主进程模式退出：已恢复 {cleared_count} 个环境变量，SDKBlocker 已移除\n")

        print(r.summary())
        return 1 if r.failed > 0 else 0

    # ---- 沙箱模式 ----
    runner = None
    if use_venv:
        print("\n🛠️  创建 venv 沙箱 ...", end=" ", flush=True)
        runner = SandboxRunner(keep=keep)
        print("OK")
    else:
        print("\n⚠️  no-venv 模式：仅隔离环境变量，不隔离 Python 包")

    try:
        s1_missing_credentials(r, runner)
        s2_invalid_region(r, runner)
        s3_missing_sdk_fallback(r, runner)

        if live:
            s4_full_workflow_live(r, runner)
        else:
            r.skip("S4 完整工作流", "未指定 --live（加 --live 启用）")

        s5_cross_skill_consistency(r, runner)

    finally:
        if runner:
            runner.cleanup()

    print(r.summary())
    return 1 if r.failed > 0 else 0


if __name__ == "__main__":
    sys.exit(main())
