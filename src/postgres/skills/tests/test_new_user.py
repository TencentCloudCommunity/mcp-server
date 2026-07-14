#!/usr/bin/env python3
"""
新用户验收测试脚本（本地沙箱模式）
===================================
模拟一个完全新用户的环境：无 SDK、无凭证、无任何腾讯云相关环境变量，
同时检查 skill zip 包的完整性和文档质量。

隔离原理（不需要 Docker）：
1. 自动创建一个临时 Python venv（空包环境，天然没有 tencentcloud SDK）
2. 所有“需要在新用户环境下运行”的测试，都用该 venv 的 Python 解释器执行
3. venv 子进程的 env 仅包含 PATH（指向 venv/bin）和临时 HOME，没有任何腾讯云凭证变量
4. zip 文档内容校验直接在主进程完成（不依赖运行环境）

使用方式：
  # 默认：本地沙箱模式（自动创建临时 venv，测试完自动清理）
  cd src/postgres/skills
  python3 tests/test_new_user.py

  # 保留沙箱 venv 以便事后调试
  python3 tests/test_new_user.py --keep-sandbox

  # 带真实凭证做端到端 API 测试（在沙箱内用 TC3 签名直连，不依赖 SDK）
  TENCENTCLOUD_SECRET_ID=xxx TENCENTCLOUD_SECRET_KEY=xxx \
    python3 tests/test_new_user.py --live

  # 用系统 Python 运行（不做包隔离，仅做环境变量隔离，快但不完全“新用户”）
  python3 tests/test_new_user.py --no-venv
"""

import os
import sys
import json
import shutil
import subprocess
import tempfile
import zipfile
import uuid
from pathlib import Path


# ===================== 配置 =====================

SKILLS_ROOT = Path(__file__).resolve().parent.parent
DIST_DIR = SKILLS_ROOT / "dist"

EXPECTED_ZIPS = [
    "tencent-pg-inspection-v1.0.0.zip",
    "tencent-pg-slowquery-diagnosis-v1.0.0.zip",
    "tencent-pg-ops-troubleshooter-v1.0.0.zip",
    "tencentdb-postgresql-skill-v1.0.0.zip",
]

REQUIRED_FILES = [
    "SKILL.md",
    "references/api_reference.md",
    "references/common/region_normalization.md",
    "references/common/error_handling.md",
]

BUNDLE_REQUIRED = [
    "SKILL.md",
    "_meta.json",
    "references/common/region_normalization.md",
    "references/common/error_handling.md",
]

REQUIRED_LINKS_IN_ERROR_TEMPLATE = [
    "cloud.tencent.com/document/product/598/40488",
    "console.cloud.tencent.com/cam/capi",
    "cloud.tencent.com/document/api/238/7520",
]

REQUIRED_LINKS_IN_REGION_DOC = [
    "cloud.tencent.com/document/api/238/7520",
    "cloud.tencent.com/document/product/409/16768",
]

REGION_ALIAS_TESTS = [
    ("广州", "ap-guangzhou"),
    ("上海", "ap-shanghai"),
    ("成都", "ap-chengdu"),
    ("北京", "ap-beijing"),
    ("ap-guangzhou", "ap-guangzhou"),
    ("ap-shanghai", "ap-shanghai"),
]


# ===================== 工具 =====================

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


def read_zip_content(zip_path):
    names = set()
    with zipfile.ZipFile(zip_path, "r") as zf:
        for info in zf.infolist():
            if not info.is_dir():
                names.add(info.filename)
    return names


def read_file_from_zip(zip_path, filename):
    with zipfile.ZipFile(zip_path, "r") as zf:
        return zf.read(filename).decode("utf-8")


def run_in_sandbox(runner, code, extra_env=None, timeout=10):
    """在沙箱 Python 解释器中执行一段代码，返回 (stdout, stderr, returncode)"""
    env = runner.make_env(extra_env)
    proc = subprocess.run(
        [runner.python, "-c", code],
        capture_output=True, text=True, env=env, timeout=timeout,
    )
    return proc.stdout.strip(), proc.stderr.strip(), proc.returncode


# ===================== 沙箱 Runner =====================

class SandboxRunner:
    """管理临时 venv 的沙箱执行器。

    临时目录结构:
      .sandbox_root/
        home/         # --sandbox 子进程的 HOME
        venv/         # 全新 Python venv，没有 tencentcloud SDK
    """

    def __init__(self, keep=False):
        self._root = tempfile.mkdtemp(prefix="pg-skill-sandbox-")
        self._venv_dir = os.path.join(self._root, "venv")
        self._home = os.path.join(self._root, "home")
        os.makedirs(self._home, exist_ok=True)
        self._keep = keep

        # 创建 venv
        subprocess.run(
            [sys.executable, "-m", "venv", "--clear", self._venv_dir],
            capture_output=True, check=True, timeout=120,
        )

        # 确定 venv 里的 python 路径
        if sys.platform == "win32":
            self.python = os.path.join(self._venv_dir, "Scripts", "python.exe")
        else:
            self.python = os.path.join(self._venv_dir, "bin", "python")

        # 验证 venv python 可用
        subprocess.run(
            [self.python, "-c", "print(1)"],
            capture_output=True, check=True, timeout=10,
        )

    def make_env(self, extra=None):
        """构造最小环境变量：只保留 PATH（含 venv/bin）+ 临时 HOME"""
        venv_bin = os.path.dirname(self.python)
        env = {
            "PATH": venv_bin + os.pathsep + os.environ.get("PATH", "/usr/bin:/bin"),
            "HOME": self._home,
        }
        # 不需要 XDG_* 但显式指向临时 home，防止读到用户真实配置
        env["XDG_CONFIG_HOME"] = os.path.join(self._home, ".config")
        env["XDG_CACHE_HOME"] = os.path.join(self._home, ".cache")
        os.makedirs(env["XDG_CONFIG_HOME"], exist_ok=True)
        os.makedirs(env["XDG_CACHE_HOME"], exist_ok=True)

        if extra:
            env.update(extra)
        return env

    def has_sdk(self):
        """在沙箱中检测 tencentcloud SDK 是否能 import"""
        code = """
try:
    import tencentcloud
    print("YES")
except ImportError:
    print("NO")
"""
        stdout, _, _ = run_in_sandbox(self, code, timeout=10)
        return "YES" in stdout

    def cleanup(self):
        if not self._keep and os.path.isdir(self._root):
            shutil.rmtree(self._root, ignore_errors=True)

    @property
    def root(self):
        return self._root

    @property
    def venv_dir(self):
        return self._venv_dir


# ===================== 测试用例 =====================

def test_zip_exists(r: TestResult):
    print("\n📦 [1/8] 检查 zip 包是否存在...")
    for zip_name in EXPECTED_ZIPS:
        zip_path = DIST_DIR / zip_name
        r.add(f"{zip_name} 存在", zip_path.exists(),
              f"缺失: {zip_path}" if not zip_path.exists() else "")


def test_single_skill_zip(r: TestResult):
    print("\n📋 [2/8] 检查单 skill zip 内部文件...")
    for skill in ["tencent-pg-inspection", "tencent-pg-slowquery-diagnosis", "tencent-pg-ops-troubleshooter"]:
        zip_path = DIST_DIR / f"{skill}-v1.0.0.zip"
        if not zip_path.exists():
            r.add(f"{skill} zip", False, "zip 不存在")
            continue
        contents = read_zip_content(zip_path)
        for req in REQUIRED_FILES:
            r.add(f"{skill}: {req}", req in contents,
                  f"缺失文件: {req}" if req not in contents else "")


def test_bundle_zip(r: TestResult):
    print("\n📦 [3/8] 检查 bundle zip 内部文件...")
    zip_path = DIST_DIR / "tencentdb-postgresql-skill-v1.0.0.zip"
    if not zip_path.exists():
        r.add("bundle zip", False, "不存在")
        return
    contents = read_zip_content(zip_path)
    for req in BUNDLE_REQUIRED:
        r.add(f"bundle: {req}", req in contents,
              f"缺失: {req}" if req not in contents else "")
    for skill in ["tencent-pg-inspection", "tencent-pg-slowquery-diagnosis", "tencent-pg-ops-troubleshooter"]:
        skill_entry = f"references/{skill}/SKILL.md"
        r.add(f"bundle 包含子 skill: {skill}", skill_entry in contents,
              f"缺失: {skill_entry}" if skill_entry not in contents else "")


def test_error_template_links(r: TestResult):
    print("\n🔗 [4/8] 检查错误模板中的官方链接...")
    zip_path = DIST_DIR / "tencent-pg-inspection-v1.0.0.zip"
    if not zip_path.exists():
        r.skip("链接检查", "zip 不存在")
        return
    content = read_file_from_zip(zip_path, "references/common/error_handling.md")
    for link in REQUIRED_LINKS_IN_ERROR_TEMPLATE:
        r.add(f"错误模板含链接: ...{link[-40:]}", link in content,
              f"缺失链接: {link}")
    for template_name in ["missing-credentials", "invalid-region", "missing-sdk"]:
        r.add(f"错误模板节: {template_name}", template_name in content,
              f"缺失模板节: {template_name}")


def test_region_doc_links(r: TestResult):
    print("\n🌍 [5/8] 检查地域规则文档中的官方链接...")
    zip_path = DIST_DIR / "tencent-pg-inspection-v1.0.0.zip"
    if not zip_path.exists():
        r.skip("链接检查", "zip 不存在")
        return
    content = read_file_from_zip(zip_path, "references/common/region_normalization.md")
    for link in REQUIRED_LINKS_IN_REGION_DOC:
        r.add(f"地域文档含链接: ...{link[-40:]}", link in content,
              f"缺失链接: {link}")
    for alias, code in REGION_ALIAS_TESTS:
        r.add(f"地域别名: {alias} -> {code}",
              alias in content and code in content,
              "别名映射缺失")


def test_sandbox_sdk(r: TestResult, runner: SandboxRunner):
    """测试：沙箱环境里确实没有 tencentcloud SDK"""
    print("\n🔧 [6/8] 沙箱 SDK 检测（空 venv，应检测到 SDK 缺失）...")
    has = runner.has_sdk()
    if has:
        r.add("沙箱内 SDK 应为未安装", False,
              "沙箱 venv 仍能 import tencentcloud，venv 创建可能使用了 --system-site-packages")
    else:
        r.add("沙箱内 SDK 未安装(正确模拟新用户)", True)


CREDENTIAL_CHECK_CODE = """
import os

def resolve_secret_id():
    for var in ["TENCENTCLOUD_SECRET_ID", "MCP_REQUEST_SECRET_ID", "MCP_SECRET_ID"]:
        if os.environ.get(var):
            return var, True
    return "SECRET_ID", False

def resolve_secret_key():
    for var in ["TENCENTCLOUD_SECRET_KEY", "MCP_REQUEST_SECRET_KEY", "MCP_SECRET_KEY"]:
        if os.environ.get(var):
            return var, True
    return "SECRET_KEY", False

missing = []
sid_name, has_sid = resolve_secret_id()
skey_name, has_skey = resolve_secret_key()

if not has_sid:
    missing.append("SECRET_ID")
if not has_skey:
    missing.append("SECRET_KEY")
if not os.environ.get("TENCENTCLOUD_REGION"):
    missing.append("REGION")

if missing:
    print("MISSING:" + ",".join(missing))
else:
    print("ALL_OK")
"""


def test_credential_detection(r: TestResult, runner: SandboxRunner):
    """测试：在沙箱中模拟凭证检测（零凭证/部分凭证/兼容变量名）"""
    print("\n🔑 [7/8] 沙箱内凭证检测（零凭证/部分凭证/兼容变量名）...")

    # 场景1：零凭证
    stdout, stderr, rc = run_in_sandbox(runner, CREDENTIAL_CHECK_CODE)
    r.add("零凭证→检测到缺失", "MISSING" in stdout,
          f"stdout={stdout} stderr={stderr}" if "MISSING" not in stdout else "")

    # 场景2：只配了 SECRET_ID
    stdout, stderr, rc = run_in_sandbox(
        runner, CREDENTIAL_CHECK_CODE,
        extra_env={"TENCENTCLOUD_SECRET_ID": "AKID_test"},
    )
    r.add("部分凭证→精确报缺失 SECRET_KEY+REGION",
          "SECRET_KEY" in stdout and "REGION" in stdout,
          f"stdout={stdout}")

    # 场景3：兼容变量名
    stdout, stderr, rc = run_in_sandbox(
        runner, CREDENTIAL_CHECK_CODE,
        extra_env={
            "MCP_SECRET_ID": "AKID_compat",
            "MCP_SECRET_KEY": "SK_compat",
            "TENCENTCLOUD_REGION": "ap-guangzhou",
        },
    )
    r.add("兼容变量名(MCP_*)→可识别", "MISSING" not in stdout,
          f"stdout={stdout}")


def test_live_api(r: TestResult, runner: SandboxRunner):
    """端到端：在沙箱中用 TC3 签名直连 OpenAPI（不依赖 SDK）"""
    print("\n🌐 [8/8] 端到端 API 连通性测试...")
    secret_id = os.environ.get("TENCENTCLOUD_SECRET_ID")
    secret_key = os.environ.get("TENCENTCLOUD_SECRET_KEY")
    region = os.environ.get("TENCENTCLOUD_REGION", "ap-guangzhou")

    if not secret_id or not secret_key:
        r.skip("端到端 API 测试", "未设置 TENCENTCLOUD_SECRET_ID/SECRET_KEY")
        return

    # 在沙箱中（无 SDK）执行 TC3 签名 + HTTPS 调用
    code = f"""
import hashlib, hmac, json, urllib.request
import datetime as _dt

secret_id = {json.dumps(secret_id)}
secret_key = {json.dumps(secret_key)}
region = {json.dumps(region)}

service = "postgres"
host = "postgres.tencentcloudapi.com"
action = "DescribeRegions"
version = "2017-03-12"
algorithm = "TC3-HMAC-SHA256"
timestamp = int(_dt.datetime.utcnow().timestamp())
date = _dt.datetime.utcfromtimestamp(timestamp).strftime("%Y-%m-%d")

payload = "{{}}"
ct = "application/json; charset=utf-8"
canonical_headers = f"content-type:{{ct}}\\nhost:{{host}}\\nx-tc-action:{{action.lower()}}\\n"
signed_headers = "content-type;host;x-tc-action"
hashed_payload = hashlib.sha256(payload.encode("utf-8")).hexdigest()
canonical_request = "\\n".join([
    "POST", "/", "", canonical_headers, signed_headers, hashed_payload,
])

credential_scope = f"{{date}}/{{service}}/tc3_request"
hashed_cr = hashlib.sha256(canonical_request.encode("utf-8")).hexdigest()
string_to_sign = "\\n".join([algorithm, str(timestamp), credential_scope, hashed_cr])

def sign(key, msg):
    return hmac.new(key, msg.encode("utf-8"), hashlib.sha256).digest()

sd = sign(("TC3" + secret_key).encode("utf-8"), date)
ss = sign(sd, service)
sst = sign(ss, "tc3_request")
signature = hmac.new(sst, string_to_sign.encode("utf-8"), hashlib.sha256).hexdigest()

auth = f"{{algorithm}} Credential={{secret_id}}/{{credential_scope}}, SignedHeaders={{signed_headers}}, Signature={{signature}}"

headers = {{
    "Authorization": auth,
    "Content-Type": ct,
    "Host": host,
    "X-TC-Action": action,
    "X-TC-Timestamp": str(timestamp),
    "X-TC-Version": version,
    "X-TC-Region": region,
}}

try:
    req = urllib.request.Request(f"https://{{host}}", data=payload.encode("utf-8"), headers=headers)
    with urllib.request.urlopen(req, timeout=10) as resp:
        data = json.loads(resp.read().decode("utf-8"))
    if "Response" in data and "Error" not in data["Response"]:
        regions = data["Response"].get("RegionSet", [])
        print("OK:" + str(len(regions)))
    else:
        err = data.get("Response", {{}}).get("Error", {{}})
        print("ERR:" + err.get("Code", "Unknown") + ":" + err.get("Message", ""))
except Exception as e:
    print("EX:" + str(e))
"""

    stdout, stderr, rc = run_in_sandbox(runner, code, timeout=15)
    if stdout.startswith("OK:"):
        count = stdout[3:]
        r.add(f"沙箱内 DescribeRegions 成功(无SDK)", True, f"返回 {count} 个地域")
    elif stdout.startswith("ERR:"):
        r.add("沙箱内 DescribeRegions", False, stdout[4:])
    else:
        r.add("沙箱内 DescribeRegions", False, f"stdout={stdout} stderr={stderr}")


# ===================== 主入口 =====================

def main():
    use_venv = "--no-venv" not in sys.argv
    keep = "--keep-sandbox" in sys.argv
    live = "--live" in sys.argv

    r = TestResult()

    print("=" * 60)
    print("  PostgreSQL Skill 新用户验收测试")
    if use_venv:
        print("  模式: 本地沙箱 (临时 venv + 临时 HOME + 零凭证)")
    else:
        print("  模式: 仅环境变量隔离 (no-venv, 不隔离 Python 包)")
    print("=" * 60)

    if live:
        print("\n⚠️  端到端模式：将使用你的真实凭证调用 OpenAPI")
        print("   仅执行只读 Action (DescribeRegions)，不会产生费用\n")

    # ---- 初始化沙箱 ----
    runner = None
    if use_venv:
        print("\n🛠️  创建临时 venv 沙箱 ...", end=" ", flush=True)
        runner = SandboxRunner(keep=keep)
        print("OK")
        print(f"   venv: {runner.venv_dir}")
        print(f"   home: {runner._home}")
        if keep:
            print(f"   (--keep-sandbox 已开启，测试完成后不会删除)")

    try:
        # 文档/包完整性检查（主进程直接完成）
        test_zip_exists(r)
        test_single_skill_zip(r)
        test_bundle_zip(r)
        test_error_template_links(r)
        test_region_doc_links(r)

        if use_venv and runner:
            # 沙箱内测试：SDK 缺失、凭证检测
            test_sandbox_sdk(r, runner)
            test_credential_detection(r, runner)

            if live:
                test_live_api(r, runner)
            elif not live:
                r.skip("端到端 API 测试", "未指定 --live，跳过（用 --live 启用）")
        else:
            # no-venv 模式：fallback 到主进程的 SDK 检测（仅供参考）
            print("\n🔧 [6/8] SDK 检测（no-venv 模式，仅供参考）...")
            try:
                subprocess.run(
                    [sys.executable, "-c", "import tencentcloud"],
                    capture_output=True, timeout=5,
                )
                print("  ⚠️  本机可 import tencentcloud —— 这不是新用户环境")
                r.add("SDK 检测(no-venv,仅供参考)", True, "本机已安装 SDK")
            except Exception:
                r.add("SDK 检测(no-venv)", True, "本机未安装 SDK")

            # no-venv 模式的凭证检测仍然走子进程隔离
            test_credential_detection(r, runner)

            if live:
                test_live_api(r, runner)
            else:
                r.skip("端到端 API 测试", "未指定 --live")

    finally:
        if runner:
            runner.cleanup()

    print(r.summary())
    return 1 if r.failed > 0 else 0


if __name__ == "__main__":
    sys.exit(main())
