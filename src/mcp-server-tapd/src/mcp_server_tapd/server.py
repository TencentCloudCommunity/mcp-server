import os
import json
import html
import requests
import markdown
from mcp.server.fastmcp import FastMCP
from mcp_server_tapd.tapd import TAPDClient
from mcp_server_tapd.app_config import AppConfig
from mcp_server_tapd.attachments import prepare_attachments, finalize_attachments, upload_files, VIDEO_EXTENSIONS

config = AppConfig()
if config.mode == "http" or config.mode == "sse" or config.mode == "streamable-http":
    mcp = FastMCP("mcp-tapd", host=config.host, port=config.port)
else:
    mcp = FastMCP("mcp-tapd")
client = TAPDClient()

# OpenClaw 把入站文件落到 <媒体目录>/inbound/<id>，但给模型看的引用是
# `media://inbound/<id>`（core 在消息里追加 "[media attached: media://inbound/…]"）。
# 模型拿到的就是这个形态，原样传进 image_url 是它唯一能做的事——2026-09-10 那条
# "补充截图"评论就是这么坏的：media:// 既不是本地路径也不是 http，被原样写进
# <img src>，TAPD 渲染成灰色占位，而 bot 还回了句"截图也带上了"。
# 媒体目录与 MCP 同机（stdio 模式的前提，和 _is_local_image_path 的假设一致）。
_MEDIA_REF_PREFIX = "media://"
_OPENCLAW_MEDIA_DIR = os.path.expanduser(os.getenv("OPENCLAW_MEDIA_DIR", "~/.openclaw/media"))


def _resolve_media_ref(url: str) -> str:
    """把 `media://<子目录>/<文件>` 还原成本机绝对路径；不是这种引用就原样返回。

    找不到文件直接抛：这是「静默丢图」的源头，宁可评论/票不建，也不要贴一张
    打不开的图然后告诉用户已经带上了。`..` 一律拒绝，别让一个引用逃出媒体目录。
    """
    if not isinstance(url, str) or not url.startswith(_MEDIA_REF_PREFIX):
        return url
    rel = url[len(_MEDIA_REF_PREFIX):].strip().lstrip("/")
    parts = [seg for seg in rel.split("/") if seg]
    if not parts or any(seg in ("..", ".") for seg in parts):
        raise ValueError(f"非法的媒体引用：{url}")
    path = os.path.realpath(os.path.join(_OPENCLAW_MEDIA_DIR, *parts))
    root = os.path.realpath(_OPENCLAW_MEDIA_DIR)
    if not (path == root or path.startswith(root + os.sep)):
        raise ValueError(f"媒体引用越出媒体目录：{url}")
    if not os.path.isfile(path):
        raise ValueError(
            f"媒体引用 {url} 在本机找不到对应文件（查找路径 {path}）。"
            "可能已被清理，或 MCP 与文件不在同一台机器。"
        )
    return path


def _normalize_media_items(raw_media) -> list[dict]:
    """Normalize supported media parameters into a single list."""
    if not raw_media:
        return []
    if isinstance(raw_media, dict):
        raw_media = [raw_media]
    normalized = []
    for item in raw_media:
        if not isinstance(item, dict):
            continue
        media_type = str(item.get("type", "")).strip().lower()
        url = str(item.get("url", "")).strip()
        if media_type not in {"image", "video"} or not url:
            continue
        normalized.append({
            "type": media_type,
            "url": url,
            "alt": str(item.get("alt", "")).strip(),
            "poster": str(item.get("poster", "")).strip(),
        })
    return normalized

def _extract_media_items(data: dict) -> list[dict]:
    """Pop media-related fields from the payload before sending to TAPD."""
    media_items = _normalize_media_items(data.pop("media", None))

    image_url = data.pop("image_url", None)
    if image_url:
        media_items.append({"type": "image", "url": str(image_url).strip(), "alt": ""})

    video_url = data.pop("video_url", None)
    if video_url:
        media_items.append({"type": "video", "url": str(video_url).strip(), "alt": "", "poster": ""})

    for key, media_type in (("image_urls", "image"), ("video_urls", "video")):
        urls = data.pop(key, None)
        if not urls:
            continue
        if isinstance(urls, str):
            urls = [urls]
        if isinstance(urls, list):
            for url in urls:
                if isinstance(url, str) and url.strip():
                    media_items.append({
                        "type": media_type,
                        "url": url.strip(),
                        "alt": "",
                        "poster": "",
                    })

    return [
        {**item, "url": _resolve_media_ref(item["url"])}
        for item in media_items if item.get("url")
    ]

_UPLOADABLE_IMAGE_EXT = {".png", ".gif", ".jpg", ".jpeg", ".bmp"}
_MAX_UPLOAD_BYTES = 5 * 1024 * 1024  # TAPD upload_image 接口硬限制

def _human_size(num_bytes: int) -> str:
    """给人看的大小。报"8388608 字节"没人算得动。"""
    mb = num_bytes / (1024 * 1024)
    return f"{mb:.1f}MB" if mb >= 1 else f"{num_bytes / 1024:.0f}KB"


_UPLOADABLE_VIDEO_EXT = VIDEO_EXTENSIONS
_MAX_ATTACHMENT_BYTES = 150 * 1024 * 1024  # TAPD upload_attachment 接口硬限制

# 挂载点。文档没枚举 type 的取值，实测（2026-09-09，真实空间）**上传接口只认
# 实体名本身**：story / bug 可以，story_description / story_description_attachment /
# story_attachment / stories / workitem 全部 422 "type is invalid"。
#
# 别被 GET /attachments 的返回误导：那里确实存在 bug_description 这种 type，但那是
# 网页端走另一套内部接口写进去的形态，上传 API 不接受。我第一版就是从那里抄的值，
# 结果每次都 422，而 422 又被吞掉了（见下面 warnings 的修复），表现成"没报错但也没传上"。
_VIDEO_ATTACH_TYPE = {
    "bug": "bug",
    "story": "story",
}


def _is_local_video_path(url: str) -> bool:
    """本地视频路径，需要上传；已经是 URL 的一律放行。"""
    if not url or "://" in url or url.startswith("//"):
        return False
    if os.path.splitext(url)[1].lower() not in _UPLOADABLE_VIDEO_EXT:
        return False
    return os.path.isfile(url)


def _split_local_videos(media_items: list[dict]) -> tuple[list[dict], list[str]]:
    """把待上传的本地视频挑出来。

    图片能在建票前传（upload_image 不需要 entry_id），视频不行：附件接口必须
    指定挂在哪个工作项上，所以只能建票之后再传，描述也就得回填。这个函数就是
    那道分水岭。
    """
    remaining, local_videos = [], []
    for item in media_items:
        if item.get("type") == "video" and _is_local_video_path(item.get("url", "")):
            local_videos.append(item["url"])
        else:
            remaining.append(item)
    return remaining, local_videos


def _render_tapd_video_html(workspace_id, attachment_id, attach_type: str) -> str:
    """TAPD 自己内嵌视频用的写法，从网页端建的票里逆出来的。

    src 只依赖 attachment_id 和 type，不是带签名的临时下载链接，所以不会过期。
    poster 是 TAPD 自动生成的首帧，我们不用管。preload=none 让票打开时不拉视频。
    """
    config = AppConfig()
    base = (config.tapd_base_url or os.getenv("TAPD_BASE_URL") or "https://www.tapd.cn").rstrip("/")
    src = f"{base}/{workspace_id}/attachments/preview_attachments/{attachment_id}/{attach_type}"
    poster = f"{base}/{workspace_id}/attachments/preview_video_poster/{attachment_id}"
    return (
        f'<p><video controls="controls" preload="none" '
        f'video_original_width="300" video_original_height="150" '
        f'src="{html.escape(src, quote=True)}" poster="{html.escape(poster, quote=True)}">'
        f'<source src="{html.escape(src, quote=True)}"></video></p>'
    )


def _attach_local_videos(workspace_id, entity: str, entry_id, video_paths: list[str],
                         warnings: list[str]) -> str:
    """建票之后把本地视频传上去，返回要追加进描述的 HTML。

    **这里一律不抛。**票此刻已经存在了，抛出去的结果是"票建了、调用方收到失败"
    ——最糟的组合，调用方多半会重试，于是多出一张重复的票。所以每一个失败都
    降级成警告：票留着，描述里少这一个视频，警告如实回给调用方去转达。
    """
    attach_type = _VIDEO_ATTACH_TYPE[entity]
    blocks = []
    for path in video_paths:
        name = os.path.basename(path)
        try:
            size = os.path.getsize(path)
            if size > _MAX_ATTACHMENT_BYTES:
                warnings.append(
                    f"视频 {name} 为 {_human_size(size)}，超过 TAPD 附件 "
                    f"{_human_size(_MAX_ATTACHMENT_BYTES)} 上限，未附到票上"
                )
                continue

            ret = client.upload_attachment(workspace_id, path, attach_type, entry_id)
            if ret.get("status") != 1:
                warnings.append(f"视频 {name} 上传失败（{ret.get('info')}），未附到票上")
                continue

            attachment_id = ((ret.get("data") or {}).get("Attachment") or {}).get("id")
            if not attachment_id:
                warnings.append(f"视频 {name} 上传后没拿到附件 id，未附到票上")
                continue

            blocks.append(_render_tapd_video_html(workspace_id, attachment_id, attach_type))
        except Exception as exc:  # noqa: BLE001 - 票已存在，任何异常都不能往上抛
            warnings.append(f"视频 {name} 处理失败（{type(exc).__name__}: {exc}），未附到票上")
    return "\n".join(blocks)

def _is_local_image_path(url: str) -> bool:
    """判断是不是一个待上传的本地图片路径（而非已经可用的 URL）。

    要求 MCP server 与文件所在机器是同一台；streamable-http 模式跨机部署时不成立。
    """
    if not url or "://" in url or url.startswith("//"):
        return False
    if url.startswith("/tfl/"):
        return False  # 已经在 TAPD 图床上了
    if os.path.splitext(url)[1].lower() not in _UPLOADABLE_IMAGE_EXT:
        return False
    return os.path.isfile(url)

def _upload_local_images(media_items: list[dict], workspace_id,
                         warnings: list[str] | None = None) -> list[dict]:
    """把本地图片路径上传到 TAPD 图床，替换成图床返回的 image_src。

    非本地路径、非图片一律原样放行。

    两类失败分开处理：
    - **超限**是事先就知道、重试也没用的确定性情况 —— 跳过这一张、记一条警告，
      票照建。为一张超大图把整张票挡下来，比缺一张图更糟。
    - **上传失败**（网络、接口报错）是意外，仍然直接抛：这一步在建票之前，
      抛出去不会留下孤票，而"悄悄少一张图"的票没人会发现。
    """
    if not workspace_id:
        return media_items

    resolved = []
    for item in media_items:
        url = item.get("url", "")
        if item.get("type") != "image" or not _is_local_image_path(url):
            resolved.append(item)
            continue

        size = os.path.getsize(url)
        if size > _MAX_UPLOAD_BYTES:
            msg = (
                f"图片 {os.path.basename(url)} 为 {_human_size(size)}，"
                f"超过 TAPD 图床 {_human_size(_MAX_UPLOAD_BYTES)} 上限，未附到票上"
            )
            if warnings is None:
                raise ValueError(msg)
            warnings.append(msg)
            continue

        ret = client.upload_image(workspace_id, url)
        if ret.get("status") != 1:
            raise ValueError(f"上传图片 {url} 失败：{ret.get('info')}")

        payload = ret.get("data") or {}
        image_src = payload.get("image_src")
        if not image_src:
            raise ValueError(f"上传图片 {url} 后返回中没有 image_src：{ret}")

        resolved.append({**item, "url": image_src})

    return resolved

# 嵌入宽度上限。`max-width: 100%` 只防溢出、不限尺寸——它的意思是「最大到容器
# 那么宽」，而 TAPD 的描述区很宽、手机截图本身又是 1000+ px，结果就是铺满整屏。
# 上限写在外层容器而不是 <img> 上：这样窄容器下 <p> 会自己缩，图片跟着缩，
# 不需要 CSS min() 也能两头都对。
MEDIA_MAX_WIDTH_PX = 680

# 高度上限。只限宽是不够的：手机截图是竖版的，1260x2844 限宽到 680 之后仍然是
# 一面 1536px 高的墙，票读起来照样被它占满——第一版就漏了这一维。
#
# 内嵌图在这里的角色是「指路的缩略图」，不是阅读界面：TAPD 自己给图挂了
# 点击放大（img.class 含 gallery，cursor: zoom-in），细节始终一点即看。
# 所以宁可缩到看不清字，也要让票本身能当文档读。
MEDIA_MAX_HEIGHT_PX = 480


# 宽高都给上限、宽高都留 auto：浏览器取更严的那个约束，等比缩放，不裁剪。
# 横版整页截图因此受宽度约束，竖版手机截图受高度约束，各自用合适的规则。
_MEDIA_FIT = (
    f"max-width: 100%; max-height: {MEDIA_MAX_HEIGHT_PX}px; width: auto; height: auto;"
)


def _render_media_html(media_items: list[dict]) -> str:
    blocks = []
    wrap = f'<p style="max-width: {MEDIA_MAX_WIDTH_PX}px;">'
    for item in media_items:
        url = html.escape(item["url"], quote=True)
        if item["type"] == "image":
            alt = html.escape(item.get("alt", ""), quote=True)
            blocks.append(
                f'{wrap}<img src="{url}" alt="{alt}" style="{_MEDIA_FIT}" /></p>'
            )
            continue

        poster = item.get("poster", "")
        poster_attr = f' poster="{html.escape(poster, quote=True)}"' if poster else ""
        blocks.append(
            f'{wrap}<video controls src="{url}"{poster_attr} style="{_MEDIA_FIT}"></video></p>'
        )
    return "\n".join(blocks)

def _render_rich_description(data: dict, field_name: str = "description",
                             supports_local_video: bool = False,
                             warnings: list[str] | None = None) -> list[str]:
    """渲染描述，并把需要建票后才能传的本地视频路径回给调用方。

    `supports_local_video=False` 的调用点（比如评论）拿到本地视频会直接抛错，
    而不是把它悄悄丢掉。静默丢弃最坏：调用方会以为录屏已经贴上去了。
    """
    content = data.get(field_name, "")
    if content is None:
        content = ""
    elif not isinstance(content, str):
        content = str(content)

    media_items, pending_videos = _split_local_videos(_extract_media_items(data))
    media_items = _upload_local_images(media_items, data.get("workspace_id"), warnings)
    media_html = _render_media_html(media_items)
    if field_name not in data and not media_html and not pending_videos:
        return []

    html_content = markdown.markdown(content, extensions=['extra', 'codehilite', 'toc'])
    if '<p>' in html_content or '<h' in html_content or '<ul>' in html_content or '<ol>' in html_content:
        content = html_content

    if media_html:
        content = f"{content}\n{media_html}" if content else media_html

    data[field_name] = content
    if pending_videos and not supports_local_video:
        raise ValueError(
            "这个接口不支持内嵌本地视频（附件接口需要 entry_id，只有建票/改票路径才拿得到）："
            + "、".join(pending_videos)
            + "。请改用 create_bug / update_bug / create_story_or_task，或先自行上传后传公网直链。"
        )
    return pending_videos

def _finalize_local_videos(entity: str, workspace_id, created: dict, payload: dict,
                           pending_videos: list[str], warnings: list[str]) -> dict:
    """Legacy video inputs use the same upload/player path as attachments display=auto."""
    files = [{"path": p, "filename": os.path.basename(p), "display": "auto"}
             for p in pending_videos]
    return finalize_attachments(client, workspace_id, entity, created, files, warnings,
                                render_video=_render_tapd_video_html)


def _filter_fields(data_list, fields_param=None):
    """
    过滤 custom_field_ 开头的自定义字段以及 description：如果值为空且不在 fields_param 中，则移除。
    同时过滤 custom_plan_field_ 开头的字段：如果值为 '0'，也移除。
    data_list: list of dict（需求或缺陷或任务列表，支持嵌套{"Story": {...}}或{"Bug": {...}}或{"TASK": {...}}）
    fields_param: str 或 list，fields参数，多个字段用逗号分隔
    """
    if not data_list:
        return data_list
    if isinstance(fields_param, str):
        fields = [f.strip() for f in fields_param.split(',') if f.strip()]
    elif isinstance(fields_param, list):
        fields = fields_param
    else:
        fields = []
    filtered = []
    for item in data_list:
        # 兼容 {"Story": {...}} 或 {"Bug": {...}} 或 {"TASK": {...}} 结构
        if isinstance(item, dict):
            if 'Story' in item and isinstance(item['Story'], dict):
                obj = item['Story']
            elif 'Bug' in item and isinstance(item['Bug'], dict):
                obj = item['Bug']
            elif 'Task' in item and isinstance(item['Task'], dict):
                obj = item['Task']
            elif 'Iteration' in item and isinstance(item['Iteration'], dict):
                obj = item['Iteration']
            else:
                obj = item
            
            # status_map = {}
            # if 'Story' in item or 'Bug' in item:
            #     system = "story" if 'Story' in item else "bug"
            #     data = {
            #         "workspace_id": workspace_id,
            #         "system": system
            #     }
            #     ret = client.get_workflows_status_map(data)
            #     if ret and isinstance(ret, dict) and 'data' in ret and isinstance(ret['data'], dict):
            #         status_map = ret['data']
            
            new_obj = {}
            for k, v in obj.items():
                if k.startswith('custom_field_') and (v is None or v == '') and (not fields_param or k not in fields):
                    continue
                if k.startswith('description') and 'Iteration' not in item and (not fields_param or k not in fields):
                    continue
                if k.startswith('custom_plan_field_') and v == '0':
                    continue
                # 如果是 status 字段，检查 status_map 中是否有对应的中文状态名
                # if k == 'status' and status_map and v in status_map:
                #     new_obj['v_status'] = status_map[v]
                
                new_obj[k] = v
            # 保持原有嵌套结构
            if 'Story' in item:
                filtered.append({'Story': new_obj})
            elif 'Bug' in item:
                filtered.append({'Bug': new_obj})
            elif 'Task' in item:
                filtered.append({'Task': new_obj})
            elif 'Iteration' in item:
                filtered.append({'Iteration': new_obj})
            else:
                filtered.append(new_obj)
        else:
            filtered.append(item)
    return filtered

def get_story_or_task_url_template(workspace_id: int, entity_type: str, tapd_base_url: str) -> str:
    """
    根据项目类型和实体类型返回 url 模板
    """
    is_mini_project = client.check_mini_project(workspace_id)
    if entity_type == 'tasks':
        return f'{tapd_base_url}/{workspace_id}/prong/tasks/view/{{id}}'
    else:
        if is_mini_project:
            return f'{tapd_base_url}/tapd_fe/t/index/{workspace_id}?workitemId={{id}}'
        else:
            return f'{tapd_base_url}/{workspace_id}/prong/stories/view/{{id}}'

@mcp.tool()
def get_user_participant_projects(nick: str = '') -> str:
    """
    nick 非必填，如果传入，则使用传入的参数，否则使用client.nick，最后使用环境变量CURRENT_USER_NICK
    获取用户参与的项目列表，如果用户进行 TAPD 操作时，没有指定 workspace_id，可以调用这个工具获取
    一次返回所有符合条件的值,只能传用户 nick 参数，一次只能查一个用户
    """
    # 获取用户昵称，优先使用传入参数，其次使用client.nick，最后使用环境变量
    nick = nick.strip() if nick else client.nick if client.nick else os.getenv("CURRENT_USER_NICK", "").strip()
    # 参数验证
    if not nick:
        return json.dumps({
            "status": 0,
            "message": "用户昵称不能为空，请提供nick参数或设置CURRENT_USER_NICK环境变量",
            "data": []
        }, indent=2, ensure_ascii=False)
    
    try:
        condition = {"nick": nick}
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    except Exception as e:
        return json.dumps({
            "status": 0,
            "message": f"获取用户参与项目失败: {str(e)}",
            "data": []
        }, indent=2, ensure_ascii=False)

@mcp.tool()
def get_stories_or_tasks(workspace_id: int, options: dict = None) -> str:
    """获取 TAPD 需求或任务，如果没有 limit 参数，则需要同时调用 get_story_count 工具获取需求数量
    
    !!! IMPORTANT !!!
    使用自定义字段(custom_field_*)前必须先调用 get_entity_custom_fields 获取字段配置
    示例流程：
    1. 调用 get_entity_custom_fields(workspace_id, {"entity_type": "stories"}) 获取需求的自定义字段配置，如果是任务，则调用 get_entity_custom_fields(workspace_id, {"entity_type": "tasks"}) 获取任务的自定义字段配置
    2. 使用返回的 custom_field_* 作为 options 参数的键名进行查询

    Args:
        workspace_id: 项目ID（必填）
        options: 参数字典（包含必填字段，不可省略），支持以下字段：
            - entity_type: 类型，需求/工作项 stories 或 任务 tasks（必填）
            - id: ID，支持多ID查询，格式为 id1,id2,id3
            - name: 标题，支持模糊匹配，例如："%搜索词%"
            - v_status: 状态别名(支持传入中文状态名称)，默认用这个字段查询状态，支持枚举查询
            - status: 状态，支持枚举查询，例如：status=status1|status2|status3
            - description: 需求内容或详细描述
            - category_name: 需求分类名称, 支持模糊匹配，例如："%搜索词%"
            - category_id: 需求分类ID
            - priority_label: 优先级，需要先检查get_stories_fields_info是否配置了候选值，如果不存在则使用默认的候选值 High => 高、Middle => 中、Low => 低、Nice To Have => 无关紧要，不使用 priority 字段
            - release_id: 发布计划ID
            - owner: 处理人
            - cc: 抄送人
            - developer: 开发人员
            - begin: 预计开始
            - due: 预计结束
            - iteration_id: 迭代ID
            - iteration_name: 迭代名称
            - parent_id: 父需求ID
            - workitem_type_id: 需求类别ID
            - workitem_type_name: 需求类别名称
            - version: 版本
            - module: 模块
            - size: 规模点，整数类型
            - ancestor_id: 祖先需求，查询指定需求下所有子需求
            - parent_id: 父需求，为0表示是根需求
            - children_id: 子需求，为空查询传：|，可以用来查询子需求或叶子需求，如果要获取所有的父需求，则传：!= |
            - creator: 创建人
            - custom_field_*: 自定义字段参数，，必须先调用 get_entity_custom_fields 获取字段配置
                示例：{"custom_field_1": "值"}
            - fields: 设置获取的字段，多个字段间以','逗号隔开，如果需要获取需求的内容/详细描述，必须带上这个字段 description
            等等...
    Returns:  <str>  # 需求或者任务所有字段数据的 json 格式，返回链接给用户，链接要可点击。而且需要返回剩余的需求数量，让用户确定是否需要继续获取剩余的需求
    Note: 如果没有给limit 参数，则需要提醒用户剩余的数量
    Note: 任务的状态status只有三个，open: 未开始，progressing: 进行中，done: 已完成

    Example:
        >>> # 第一步：获取需求的自定义字段配置
        >>> fields_config = get_entity_custom_fields(
        ...     workspace_id=123,
        ...     options={"entity_type": "stories"}
        ... )
        >>> # 第二步：使用配置中的字段进行查询需求
        >>> stories = get_stories_or_tasks(
        ...     workspace_id=123,
        ...     options={"entity_type": "stories", "custom_field_1": "值"}
        ... )

    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    story_condition = {
        "workspace_id": workspace_id,
    }
    
    if options is None:
        options = {}
    if 'entity_type' not in options:
        options['entity_type'] = 'stories'
    if 'category_name' in options:
        categroy_condition = {
            "workspace_id": workspace_id,
            "name": options['category_name']
        }
        category_result = client.get_category_id(categroy_condition)
        # 判断返回的分类数量
        if category_result and isinstance(category_result, dict) and 'data' in category_result and isinstance(category_result['data'], list):
            categories = category_result['data']
            if len(categories) == 0:
                return "没有找到需求分类，请检查需求分类是否存在"
            elif len(categories) == 1:
                # 只返回一个分类，赋值 category_id
                category_id = categories[0]['Category']['id']
                options['category_id'] = category_id
                del options['category_name']
            else:
                # 多个分类，返回分类id和名列表和提示
                names = [f"{c['Category']['id']}:{c['Category']['name']}" for c in categories]
                msg = f"有多个分类，请选择：{', '.join(names)}"
                return msg
        else:
            return "没有找到需求分类，请检查需求分类是否存在"
    
    if 'iteration_name' in options:
        iteration_condition = {
            "workspace_id": workspace_id,
            "name": options['iteration_name']
        }
        iteration_result = client.get_iterations(iteration_condition)
        if iteration_result and isinstance(iteration_result, dict) and 'data' in iteration_result and isinstance(iteration_result['data'], list):
            iterations = iteration_result['data']
            if len(iterations) == 0:
                return "没有找到迭代，请检查迭代是否存在"
            elif len(iterations) == 1:
                # 只返回一个迭代，赋值 iteration_id
                iteration_id = iterations[0]['Iteration']['id']
                options['iteration_id'] = iteration_id
                del options['iteration_name']
            else:
                # 多个迭代，返回迭代id和名列表和提示
                names = [f"{c['Iteration']['id']}:{c['Iteration']['name']}" for c in iterations]
                msg = f"有多个迭代，请选择：{', '.join(names)}"
                return msg
        else:
            return "没有找到迭代，请检查迭代是否存在"
    if 'workitem_type_name' in options:
        workitem_type_condition = {
            "workspace_id": workspace_id,
            "name": options['workitem_type_name']
        }
        workitem_type_result = client.get_workitem_types(workitem_type_condition)
        if workitem_type_result and isinstance(workitem_type_result, dict) and 'data' in workitem_type_result and isinstance(workitem_type_result['data'], list):
            workitem_types = workitem_type_result['data']
            if len(workitem_types) == 0:
                return "没有找到需求类别，请检查需求类别是否存在"
            elif len(workitem_types) == 1:
                # 只返回一个需求类别，赋值 workitem_type_id
                workitem_type_id = workitem_types[0]['WorkitemType']['id']
                options['workitem_type_id'] = workitem_type_id
                del options['workitem_type_name']
            else:
                # 多个需求类别，返回需求类别id和名列表和提示
                names = [f"{c['WorkitemType']['id']}:{c['WorkitemType']['name']}" for c in workitem_types]
                msg = f"有多个需求类别，请选择：{', '.join(names)}"
                return msg
        else:
            return "没有找到需求类别，请检查需求类别是否存在"
    
    if options:
        story_condition.update(options)
    
    ret = client.get_stories(story_condition)
    count_ret = client.get_story_count(story_condition)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    
    # 处理 custom_field_ 字段过滤
    fields_param = options.get('fields') if options else None
    # ret 可能是 {"data": [{...}, ...]} 或 list
    if isinstance(ret, dict) and 'data' in ret and isinstance(ret['data'], list):
        ret['data'] = _filter_fields(ret['data'], fields_param)
    elif isinstance(ret, list):
        ret = _filter_fields(ret, fields_param)

    entity_type = options['entity_type']
    url_template = get_story_or_task_url_template(workspace_id, entity_type, tapd_base_url)

    return json.dumps({
        "url_template": url_template,
        "data": ret['data'] if isinstance(ret, dict) and 'data' in ret else ret,
        "count": count_ret
    }, ensure_ascii=False, indent=2)

@mcp.tool()
def get_entity_custom_fields(workspace_id: int, options: dict = None) -> str:
    """获取 TAPD 需求或任务或迭代或测试用例的自定义字段配置
    Args:
        workspace_id: 项目ID（必填）
        options: 参数字典（包含必填字段，不可省略），支持以下字段：
            - entity_type: 类型，需求 stories 或任务 tasks 或迭代 iterations 或测试用例 tcases（必填）
    Returns:  <str>  # 需求或者任务或迭代或测试用例所有自定义字段配置数据的 json 格式
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    if options:
        data.update(options)
    ret = client.get_entity_custom_fields(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def get_story_or_task_count(workspace_id: int, options: dict = None) -> str:
    """获取 TAPD 需求或者任务的数量
    TAPD API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/story/get_stories_count.html
    Args:
        workspace_id: 项目ID（必填）
        options: 参数字典（包含必填字段，不可省略），支持以下字段：
            - entity_type: 类型，需求stories或任务tasks（必填）
            - id: ID
            - name: 标题
            - workitem_type_id: 需求类别ID，支持枚举查询，例如：workitem_type_id=workitem_type_id_1|workitem_type_id_2|workitem_type_id_3
            - status: 状态，支持枚举查询，例如：status=status1|status2|status3
            等等...
    Returns:  <str>  # 需求所有字段数据的 json 格式
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    new_story = {
        "workspace_id": workspace_id,
    }
    
    if options:
        new_story.update(options)
    
    ret = client.get_story_count(new_story)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def update_story_or_task(workspace_id: int, options: dict = None) -> str:
    """
    options.attachments: [{"path": "/local/path or media://inbound/...", "filename": "original.eml"}].
    attachments 默认 display="auto"：图片内嵌、视频附件加正文播放器、其他文件仅附件；display="attachment" 强制仅附件。
    不解压 ZIP，不解析邮件。视频上传与 preview_status 分别检查，预览失败不能重传已成功附件。
    必须检查返回 data 内 attachment_results 和 warnings，不能把建票成功当成附件成功。
    更新 TAPD 需求或任务，任务的状态只有三个，open: 未开始，progressing: 进行中，done: 已完成，如果是流转需求的状态，需要调用工具先获取状态信息get_workflows_all_transitions
    TAPD API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/story/update_story.html
    Args:
        workspace_id: 项目ID（必填）
        options: 参数字典（包含必填字段，不可省略），支持以下字段：
            - entity_type: 类型，需求/工作项 stories 或 任务 tasks（必填）
            - id: 需求或任务ID（必填）
            - name: 标题
            - v_status: 状态别名(支持传入中文状态名称)，默认用这个字段
            - status: 状态，原名
            - priority_label: 优先级，需要先检查get_stories_fields_info是否配置了候选值，如果不存在则使用默认的候选值 High => 高、Middle => 中、Low => 低、Nice To Have => 无关紧要，不使用 priority 字段
            - description: 描述
            - media: 富媒体列表，支持 [{"type": "image"|"video", "url": "...", "alt": "...", "poster": "..."}]
            - image_url: 单张图片，会以内嵌图片形式写入描述。可以是公网直链，也可以是本机图片的绝对路径（如 /root/.openclaw/media/inbound/xxx.png），本地路径会自动上传到 TAPD 图床后再内嵌，无需自己找图床；消息里给出的 media://inbound/xxx 引用也可以直接传
            - image_urls: 多张图片，规则同 image_url，传数组
            - video_url: 单个视频，可以是公网直链，也可以是本机视频的绝对路径（如 /root/.openclaw/media/inbound/xxx.mp4）。本地路径会在票创建后自动上传为 TAPD 附件并内嵌进描述，可在票里直接播放，不需要任何外部图床或直链。单文件限 150MB
            - video_urls: 多个视频，规则同 video_url，传数组
            - owner: 处理人
            - cc: 抄送人
            - developer: 开发人员
            - begin: 预计开始
            - due: 预计结束
            - iteration_id: 迭代ID
            - iteration_name: 迭代名称
            - templated_id: 模板ID
            - parent_id: 父需求ID
            - workitem_type_id: 需求类别ID
            - workitem_type_name: 需求类别名称
            - category_id: 需求分类ID
            - category_name: 需求分类名称
            - release_id: 发布计划ID
            - version: 版本
            - module: 模块
            - size: 规模点，整数类型
            - story_id: 如果 entity_type 为 tasks，则表示任务关联的需求ID
            - cus_{$自定义字段别名}: 自定义字段值，参数名会由后台自动转义为custom_field_*，如：cus_自定义字段的名称
            - custom_field_*: 自定义字段参数
            等等...
    Returns: <str>,  # 需求或者任务所有字段数据的 json 格式
    Note: 需求链接格式为 {tapd_base_url}/{workspace_id}/prong/stories/view/{story_id}
    Note: 任务链接格式为 {tapd_base_url}/{workspace_id}/prong/tasks/view/{id}
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    new_story = {
        "workspace_id": workspace_id,
    }
    if options is None:
        options = {}
    if 'entity_type' not in options:
        options['entity_type'] = 'stories'
    if options:
        new_story.update(options)
    pending_attachments = prepare_attachments(new_story, _resolve_media_ref)
    warnings: list[str] = []
    pending_videos = _render_rich_description(
        new_story, supports_local_video=True, warnings=warnings
    )
    created_story = client.create_or_update_story(new_story)
    created_story = _finalize_local_videos(
        "story", workspace_id, created_story, new_story, pending_videos, warnings
    )
    created_story = finalize_attachments(client, workspace_id, "story", created_story, pending_attachments, warnings, render_video=_render_tapd_video_html)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    
    entity_type = options['entity_type']
    url_template = get_story_or_task_url_template(workspace_id, entity_type, tapd_base_url)
    return json.dumps({
        "url_template": url_template,
        "data": json.dumps(created_story, indent=2, ensure_ascii=False),
        **({"warnings": warnings} if warnings else {}),
    }, ensure_ascii=False, indent=2)

@mcp.tool()
def create_story_or_task(workspace_id: int, name: str, options: dict = None) -> dict:
    """
    options.attachments: [{"path": "/local/path or media://inbound/...", "filename": "original.eml"}].
    attachments 默认 display="auto"：图片内嵌、视频附件加正文播放器、其他文件仅附件；display="attachment" 强制仅附件。
    不解压 ZIP，不解析邮件。视频上传与 preview_status 分别检查，预览失败不能重传已成功附件。
    必须检查返回 data 内 attachment_results 和 warnings，不能把建票成功当成附件成功。
    创建 TAPD 需求或任务
    Args:
        workspace_id: 项目ID（必填）
        name: 需求或任务的标题（必填）
        options: 参数字典（包含必填字段，不可省略），支持以下字段：
            - entity_type: 类型，需求/工作项 stories 或 任务 tasks（必填）
            - priority_label: 优先级，需要先检查get_stories_fields_info是否配置了候选值，如果不存在则使用默认的候选值 High => 高、Middle => 中、Low => 低、Nice To Have => 无关紧要，不使用 priority 字段
            - description: 描述
            - media: 富媒体列表，支持 [{"type": "image"|"video", "url": "...", "alt": "...", "poster": "..."}]
            - image_url: 单张图片，会以内嵌图片形式写入描述。可以是公网直链，也可以是本机图片的绝对路径（如 /root/.openclaw/media/inbound/xxx.png），本地路径会自动上传到 TAPD 图床后再内嵌，无需自己找图床；消息里给出的 media://inbound/xxx 引用也可以直接传
            - image_urls: 多张图片，规则同 image_url，传数组
            - video_url: 单个视频，可以是公网直链，也可以是本机视频的绝对路径（如 /root/.openclaw/media/inbound/xxx.mp4）。本地路径会在票创建后自动上传为 TAPD 附件并内嵌进描述，可在票里直接播放，不需要任何外部图床或直链。单文件限 150MB
            - video_urls: 多个视频，规则同 video_url，传数组
            - owner: 处理人
            - cc: 抄送人
            - developer: 开发人员
            - begin: 预计开始
            - due: 预计结束
            - iteration_id: 迭代ID
            - iteration_name: 迭代名称
            - templated_id: 模板ID
            - parent_id: 父需求ID
            - workitem_type_id: 需求类别ID
            - workitem_type_name: 需求类别名称
            - category_id: 需求分类ID
            - category_name: 需求分类名称
            - release_id: 发布计划ID
            - version: 版本
            - module: 模块
            - size: 规模点，整数类型
            - story_id 如果是 entity_type 为 tasks 则表示任务关联的需求ID
            - creator: 创建人
            - cus_{$自定义字段别名}: 自定义字段值，参数名会由后台自动转义为custom_field_*，如：cus_自定义字段的名称
            - custom_field_*: 自定义字段参数
            等等...
    Returns:
        {
            "data": <str>,  # 需求或任务所有字段数据的 json 格式
            "url": <str>  # 需求或者任务 url，返回给用户时，链接要可点击
        }
    """
    user_nick = os.getenv("CURRENT_USER_NICK")
    if not workspace_id and user_nick:
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
        "name": name
    }
    if 'creator' not in data and user_nick:
        data['creator'] = user_nick
    if options is None:
        options = {}
    if 'entity_type' not in options:
        options['entity_type'] = 'stories'
    if 'category_name' in options:
        categroy_condition = {
            "workspace_id": workspace_id,
            "name": options['category_name']
        }
        category_result = client.get_category_id(categroy_condition)
        # 判断返回的分类数量
        if category_result and isinstance(category_result, dict) and 'data' in category_result and isinstance(category_result['data'], list):
            categories = category_result['data']
            if len(categories) == 0:
                return "没有找到需求分类，请检查需求分类是否存在"
            elif len(categories) == 1:
                # 只返回一个分类，赋值 category_id
                category_id = categories[0]['Category']['id']
                options['category_id'] = category_id
                del options['category_name']
            else:
                # 多个分类，返回分类id和名列表和提示
                names = [f"{c['Category']['id']}:{c['Category']['name']}" for c in categories]
                msg = f"有多个分类，请选择：{', '.join(names)}"
                return msg
        else:
            return "没有找到需求分类，请检查需求分类是否存在"
    if 'iteration_name' in options:
        iteration_condition = {
            "workspace_id": workspace_id,
            "name": options['iteration_name']
        }
        iteration_result = client.get_iterations(iteration_condition)
        if iteration_result and isinstance(iteration_result, dict) and 'data' in iteration_result and isinstance(iteration_result['data'], list):
            iterations = iteration_result['data']
            if len(iterations) == 0:
                return "没有找到迭代，请检查迭代是否存在"
            elif len(iterations) == 1:
                # 只返回一个迭代，赋值 iteration_id
                iteration_id = iterations[0]['Iteration']['id']
                options['iteration_id'] = iteration_id
                del options['iteration_name']
            else:
                # 多个迭代，返回迭代id和名列表和提示
                names = [f"{c['Iteration']['id']}:{c['Iteration']['name']}" for c in iterations]
                msg = f"有多个迭代，请选择：{', '.join(names)}"
                return msg
        else:
            return "没有找到迭代，请检查迭代是否存在"
    if 'workitem_type_name' in options:
        workitem_type_condition = {
            "workspace_id": workspace_id,
            "name": options['workitem_type_name']
        }
        workitem_type_result = client.get_workitem_types(workitem_type_condition)
        if workitem_type_result and isinstance(workitem_type_result, dict) and 'data' in workitem_type_result and isinstance(workitem_type_result['data'], list):
            workitem_types = workitem_type_result['data']
            if len(workitem_types) == 0:
                return "没有找到需求类别，请检查需求类别是否存在"
            elif len(workitem_types) == 1:
                # 只返回一个需求类别，赋值 workitem_type_id
                workitem_type_id = workitem_types[0]['WorkitemType']['id']
                options['workitem_type_id'] = workitem_type_id
                del options['workitem_type_name']
            else:
                # 多个需求类别，返回需求类别id和名列表和提示
                names = [f"{c['WorkitemType']['id']}:{c['WorkitemType']['name']}" for c in workitem_types]
                msg = f"有多个需求类别，请选择：{', '.join(names)}"
                return msg
        else:
            return "没有找到需求类别，请检查需求类别是否存在"

    if options:
        data.update(options)
    
    pending_attachments = prepare_attachments(data, _resolve_media_ref)
    warnings: list[str] = []
    pending_videos = _render_rich_description(data, supports_local_video=True, warnings=warnings)

    created_story = client.create_or_update_story(data)
    created_story = _finalize_local_videos(
        "story", workspace_id, created_story, data, pending_videos, warnings
    )
    created_story = finalize_attachments(client, workspace_id, "story", created_story, pending_attachments, warnings, render_video=_render_tapd_video_html)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    entity_type = data['entity_type']
    if entity_type == 'tasks':
        entity_id = created_story['data']['Task']['id']
    else:
        entity_id = created_story['data']['Story']['id']
    
    entity_type = options['entity_type']
    url_template = get_story_or_task_url_template(workspace_id, entity_type, tapd_base_url)
    return {
        "url_template": url_template,
        "data": json.dumps(created_story, indent=2, ensure_ascii=False),
        # 漏掉这一行的代价：附件上传失败被完全吞掉，调用方看到一切正常。
        **({"warnings": warnings} if warnings else {}),
    }

@mcp.tool()
def get_bug(workspace_id: int, options: dict = None) -> str:
    """获取 TAPD 缺陷
    
    !!! IMPORTANT !!!
    使用自定义字段(custom_field_*)前必须先调用 get_entity_custom_fields 获取字段配置
    示例流程：
    1. 调用 get_entity_custom_fields(workspace_id, {"entity_type": "bugs"}) 获取需求的自定义字段配置，如果是任务，则调用 get_entity_custom_fields(workspace_id, {"entity_type": "bugs"}) 获取缺陷的自定义字段配置
    2. 使用返回的 custom_field_* 作为 options 参数的键名进行查询
    
    Args:
        workspace_id: 项目ID（必填）
        options: 可选参数字典，支持以下字段：
            - id: ID
            - title: 标题
            - status: 状态，支持枚举查询，例如：status=status1|status2|status3
            - priority_label: 优先级，需要先检查get_entity_custom_fields是否配置了自定义候选值，没有的话用默认的候选值 urgent=> 紧急、high=> 高、medium=> 中、low=> 低、insignificant=> 无关紧要，不使用 priority
            - severity: 严重程度，取值为 fatal => 致命、 serious => 严重、 normal => 一般、 prompt => 提示、advice => 建议
            - custom_field_{number}: 自定义字段参数，必须先调用 get_entity_custom_fields 获取字段配置
                示例：{"custom_field_1": "值"}
            - fields: 设置获取的字段，多个字段间以','逗号隔开，如果需要获取需求的内容/详细描述，必须带上这个字段 description
            等等...
    Returns:  <str>  # 缺陷所有字段数据的 json 格式
    Note: 缺陷链接格式为 {tapd_base_url}/{workspace_id}/bugtrace/bugs/view/{id}
    Example:
        >>> # 第一步：获取自定义字段配置
        >>> fields_config = get_entity_custom_fields(
        ...     workspace_id=123,
        ...     options={"entity_type": "bugs"}
        ... )
        >>> # 第二步：使用配置中的字段进行查询
        >>> bugs = get_bug(
        ...     workspace_id=123,
        ...     options={"custom_field_1": "值"}
        ... )

    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    bug_condition = {
        "workspace_id": workspace_id,
    }
    
    if options:
        bug_condition.update(options)
    
    ret = client.get_bug(bug_condition)
    count_ret = client.get_bug_count(bug_condition)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url

    # 处理 custom_field_ 和 custom_plan_field_ 字段过滤
    fields_param = options.get('fields') if options else None
    if isinstance(ret, dict) and 'data' in ret and isinstance(ret['data'], list):
        ret['data'] = _filter_fields(ret['data'], fields_param)
    elif isinstance(ret, list):
        ret = _filter_fields(ret, fields_param)

    return json.dumps({
        "base_url": tapd_base_url, # 返回给用户时，链接要可点击
        "data": ret['data'] if isinstance(ret, dict) and 'data' in ret else ret,
        "count": count_ret
    }, ensure_ascii=False, indent=2)

@mcp.tool()
def get_bug_count(workspace_id: int, options: dict = None) -> str:
    """获取 TAPD 缺陷
    TAPD API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/get_bugs_count.html
    Args:
        workspace_id: 项目ID（必填）
        options: 可选参数字典，支持以下字段：
            - id: ID
            - name: 标题
            等等...
    Returns:  <str>  # 需求所有字段数据的 json 格式
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    new_bug = {
        "workspace_id": workspace_id,
    }
    
    if options:
        new_bug.update(options)
    
    ret = client.get_bug_count(new_bug)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def update_bug(workspace_id: int, options: dict = None) -> dict:
    """
    options.attachments: [{"path": "/local/path or media://inbound/...", "filename": "original.eml"}].
    attachments 默认 display="auto"：图片内嵌、视频附件加正文播放器、其他文件仅附件；display="attachment" 强制仅附件。
    不解压 ZIP，不解析邮件。视频上传与 preview_status 分别检查，预览失败不能重传已成功附件。
    必须检查返回 data 内 attachment_results 和 warnings，不能把建票成功当成附件成功。
    更新 TAPD 缺陷
    TAPD API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/update_bug.html
    Args:
        workspace_id: 项目ID（必填）
        options: 参数字典（包含必填字段，不可省略），支持以下字段：
            - id: 缺陷ID（必填）
            - title: 标题
            - description: 描述
            - media: 富媒体列表，支持 [{"type": "image"|"video", "url": "...", "alt": "...", "poster": "..."}]
            - image_url: 单张图片，会以内嵌图片形式写入描述。可以是公网直链，也可以是本机图片的绝对路径（如 /root/.openclaw/media/inbound/xxx.png），本地路径会自动上传到 TAPD 图床后再内嵌，无需自己找图床；消息里给出的 media://inbound/xxx 引用也可以直接传
            - image_urls: 多张图片，规则同 image_url，传数组
            - video_url: 单个视频，可以是公网直链，也可以是本机视频的绝对路径（如 /root/.openclaw/media/inbound/xxx.mp4）。本地路径会在票创建后自动上传为 TAPD 附件并内嵌进描述，可在票里直接播放，不需要任何外部图床或直链。单文件限 150MB
            - video_urls: 多个视频，规则同 video_url，传数组
            - v_status: 状态别名(支持传入中文状态名称)，默认用这个字段
            - status: 状态，原名
            - priority_label: 优先级，需要先检查get_entity_custom_fields是否配置了自定义候选值，没有的话用默认的候选值 urgent=> 紧急、high=> 高、medium=> 中、low=> 低、insignificant=> 无关紧要，不使用 priority
            - severity: 严重程度，取值为 fatal => 致命、 serious => 严重、 normal => 一般、 prompt => 提示、advice => 建议
            等等...
    Returns: <str>  # 缺陷所有字段数据的 json 格式
    Note: 缺陷链接格式为 {tapd_base_url}/{workspace_id}/bugtrace/bugs/view/{id}
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    new_bug = {
        "workspace_id": workspace_id,
    }
    
    if options:
        new_bug.update(options)

    pending_attachments = prepare_attachments(new_bug, _resolve_media_ref)
    warnings: list[str] = []
    pending_videos = _render_rich_description(
        new_bug, supports_local_video=True, warnings=warnings
    )
    user_nick = os.getenv("CURRENT_USER_NICK")
    if 'lastmodify' not in new_bug and user_nick:
        new_bug['lastmodify'] = user_nick
    created_bug = client.create_or_update_bug(new_bug)
    created_bug = _finalize_local_videos(
        "bug", workspace_id, created_bug, new_bug, pending_videos, warnings
    )
    created_bug = finalize_attachments(client, workspace_id, "bug", created_bug, pending_attachments, warnings, render_video=_render_tapd_video_html)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    return {
        "base_url": tapd_base_url, # 返回给用户时，链接要可点击
        "data": json.dumps(created_bug, indent=2, ensure_ascii=False),
        **({"warnings": warnings} if warnings else {}),
    }

@mcp.tool()
def create_bug(workspace_id: int, title: str, options: dict = None) -> dict:
    """
    options.attachments: [{"path": "/local/path or media://inbound/...", "filename": "original.eml"}].
    attachments 默认 display="auto"：图片内嵌、视频附件加正文播放器、其他文件仅附件；display="attachment" 强制仅附件。
    不解压 ZIP，不解析邮件。视频上传与 preview_status 分别检查，预览失败不能重传已成功附件。
    必须检查返回 data 内 attachment_results 和 warnings，不能把建票成功当成附件成功。
    创建 TAPD 缺陷
    TAPD API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/add_bug.html
    Args:
        workspace_id: 项目ID（必填）
        title: 缺陷标题（必填）
        options: 可选参数字典，支持以下字段：
            - priority_label: 优先级，默认的候选值 urgent=> 紧急、high=> 高、medium=> 中、low=> 低、insignificant=> 无关紧要，不使用 priority
            - current_owner: 处理人
            - cc: 抄送人
            - reporter: 创建人
            - participator: 参与人
            - te: 测试人员
            - de: 开发人员
            - auditer: 审核人员
            - confirmer: 验证人员
            - description: 描述
            - media: 富媒体列表，支持 [{"type": "image"|"video", "url": "...", "alt": "...", "poster": "..."}]
            - image_url: 单张图片，会以内嵌图片形式写入描述。可以是公网直链，也可以是本机图片的绝对路径（如 /root/.openclaw/media/inbound/xxx.png），本地路径会自动上传到 TAPD 图床后再内嵌，无需自己找图床；消息里给出的 media://inbound/xxx 引用也可以直接传
            - image_urls: 多张图片，规则同 image_url，传数组
            - video_url: 单个视频，可以是公网直链，也可以是本机视频的绝对路径（如 /root/.openclaw/media/inbound/xxx.mp4）。本地路径会在票创建后自动上传为 TAPD 附件并内嵌进描述，可在票里直接播放，不需要任何外部图床或直链。单文件限 150MB
            - video_urls: 多个视频，规则同 video_url，传数组
            - priority_label: 优先级，需要先检查get_entity_custom_fields是否配置了自定义候选值，没有的话用默认的候选值 urgent=> 紧急、high=> 高、medium=> 中、low=> 低、insignificant=> 无关紧要，不使用 priority
            - severity: 严重程度，取值为 fatal => 致命、 serious => 严重、 normal => 一般、 prompt => 提示、advice => 建议
            - module: 模块
            - feature: 特性
            - release_id: 发布计划 ID
            - template_id: 模板 ID
            - iteration_id: 迭代 ID
            - size: 规模点，整数类型
             - cus_{$自定义字段别名}: 自定义字段值，参数名会由后台自动转义为custom_field_*，如：cus_自定义字段的名称
            - custom_field_*: 自定义字段参数
            等等...
    Returns:
        {
            "data": <str>,      # 所有字段数据的 json 格式
            "url": <str>,       #  url，返回给用户时，链接要可点击
            "warnings": [<str>] # 可选。票已建好，但有附件没贴上（超限/上传失败）。
                                # 有这个字段就如实转达给用户，不要重试建票
        }
    Note: 缺陷链接格式为 {tapd_base_url}/{workspace_id}/bugtrace/bugs/view/{id}
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    new_bug = {
        "workspace_id": workspace_id,
        "title": title
    }
    
    if options:
        new_bug.update(options)
    
    pending_attachments = prepare_attachments(new_bug, _resolve_media_ref)
    warnings: list[str] = []
    pending_videos = _render_rich_description(
        new_bug, supports_local_video=True, warnings=warnings
    )
    user_nick = os.getenv("CURRENT_USER_NICK")
    if 'reporter' not in new_bug and user_nick:
        new_bug['reporter'] = user_nick
    ret = client.create_or_update_bug(new_bug)
    ret = _finalize_local_videos("bug", workspace_id, ret, new_bug, pending_videos, warnings)
    ret = finalize_attachments(client, workspace_id, "bug", ret, pending_attachments, warnings, render_video=_render_tapd_video_html)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    return {
        "url": f'{tapd_base_url}/{workspace_id}/bugtrace/bugs/view/{ret["data"]["Bug"]["id"]}', # 返回给用户时，链接要可点击
        "data": json.dumps(ret, indent=2, ensure_ascii=False),
        # 票建成了但有附件没贴上：如实转达，别让用户以为录屏在票里
        **({"warnings": warnings} if warnings else {}),
    }

@mcp.tool()
def create_comments(workspace_id: int, options: dict = None) -> dict:
    """添加评论
    TAPD API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/comment/add_comment.html
    Args:
        workspace_id: 项目ID（必填）
        options: 参数字典（包含必填字段，不可省略），支持以下字段：
            - entry_id: 评论所依附的业务对象实体id（必填）
            - entry_type: 评论类型（取值： bug、 bug_remark （流转缺陷时候的评论）、 stories、 tasks 。）（必填）
            - author: 评论人（必填）
            - description: 内容（必填）
            - media: 富媒体列表，支持 [{"type": "image"|"video", "url": "...", "alt": "...", "poster": "..."}]
            - image_url: 单张图片，会以内嵌图片形式写入评论。可以是公网直链，也可以是本机图片的绝对路径（如 /root/.openclaw/media/inbound/xxx.png），本地路径会自动上传到 TAPD 图床后再内嵌，无需自己找图床；消息里给出的 media://inbound/xxx 引用也可以直接传
            - image_urls: 多张图片，规则同 image_url，传数组
            - video_url: 单个视频直链，会以内嵌视频形式写入评论
            - video_urls: 多个视频直链，会以内嵌视频形式写入评论
            - root_id: 根评论ID
            - reply_id: 需求评论回复的ID
    Returns: <str>  # 新建评论的数据
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    if options:
        data.update(options)

    # warnings 必须回给调用方。之前这里把它丢了，超限/失败的图就"没报错但也没
    # 传上"，bot 会对用户说"截图也带上了"——最坏的失败方式。
    warnings: list[str] = []
    _render_rich_description(data, warnings=warnings)
    created = client.create_comments(data)
    if warnings and isinstance(created, dict):
        created = {**created, "warnings": warnings}
    return json.dumps(created, indent=2, ensure_ascii=False)

@mcp.tool()
def update_comments(workspace_id: int, options: dict = None) -> dict:
    """更新评论，返回更新评论的数据
    Args:
        workspace_id: 项目ID（必填）
        options: 参数字典（包含必填字段，不可省略），支持以下字段：
            - id: 评论id（必填）
            - description: 内容（必填）
            - change_creator: 变更人（必填）
    Returns: <str>  # 更新后的评论数据
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    if options:
        data.update(options)
    
    _render_rich_description(data)
    created_story = client.create_comments(data)
    return json.dumps(created_story, indent=2, ensure_ascii=False)

@mcp.tool()
def get_comments(workspace_id: int, options: dict = None) -> dict:
    """获取符合查询条件的所有评论（分页显示，默认一页30条）
    Args:
        workspace_id: 项目ID（必填）
        options: 可选参数字典，支持以下字段：
            - id: 评论ID（格式如：1159271484001002933）
            - title: 标题
            - description: 内容
            - author: 评论人
            - entry_type: 评论类型（取值： bug、 bug_remark （流转缺陷时候的评论）、 stories、 tasks 。多个类型间以竖线隔开）
            - entry_id: 评论所依附的业务对象实体id
            - created: 创建时间，支持时间查询
            - limit: 设置返回数量限制，默认为30
            - page: 返回当前数量限制下第N页的数据，默认为1（第一页）
            - order: 排序规则，规则：字段名 ASC或者DESC，然后 urlencode，如按创建时间逆序：order=created%20desc
            - fields: 设置获取的字段，多个字段间以','逗号隔开
    Returns: <str>  # 符合条件的评论数据
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    if options:
        data.update(options)
    
    ret = client.get_comments(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def get_workflows_all_transitions(workspace_id: int, options: dict = None) -> dict:
    """获取项目下的工作流流转细则，例如要流转到"实现中"，则需要调用这个工具查看当前状态能流转到的状态
    Args:
        workspace_id: 项目ID（必填）
        options:
            - system: 系统名。取 bug （缺陷的）或者 story（需求的）（必填）
            - workitem_type_id: 需求类别ID（必填）
    Returns: <str>  状态流转细则
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    if options:
        data.update(options)
    
    ret = client.get_workflows_all_transitions(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def get_workflows_status_map(workspace_id: int, options: dict = None) -> dict:
    """获取工作流状态中英文名对应关系，例如想要获取所有"实现中"的需求/缺陷等，可以调用这个接口先获取状态的英文名，再根据这个英文名作为 status 字段的值去查询
    Args:
        workspace_id: 项目ID（必填）
        options:
            - system: 系统名。取 bug （缺陷的）或者 story（需求的）（必填）
            - workitem_type_id: 需求类别ID（必填）
    Returns: <str>  状态流转细则
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    if options:
        data.update(options)
    
    ret = client.get_workflows_status_map(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def get_workitem_types(workspace_id: int, options: dict = None) -> dict:
    """获取符合查询条件的所有需求类别（分页显示，默认一页30条）
    TAPD API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/story/get_workitem_types.html
    Args:
        workspace_id: 项目ID（必填）
        options:
            - id: id，支持多ID查询
            - name: 需求类别名称
            等等...
    Returns: <str>  项目下所有需求类别字段数据
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    if options:
        data.update(options)
    
    ret = client.get_workitem_types(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def get_workflows_last_steps(workspace_id: int, options: dict = None) -> dict:
    """获取工作流结束状态，一次只能获取一个项目的工作流结束状态
    Args:
        workspace_id: 项目ID（必填）
        options:
            - system: 系统名。取 bug （缺陷的）或者 story（需求的）（必填）
            - workitem_type_id: 需求类别id（不传则取项目下所有工作流的结束状态）
            - type: 节点类型，仅并行工作流需区分。status 状态，step 并行工作流节点。默认只返回结束状态。若需要同时返回结束状态和结束节点，支持数组type[]=status&type[]=step
    Returns: <str> 工作流结束状态
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    if options:
        data.update(options)
    
    ret = client.get_workflows_last_steps(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def get_stories_fields_lable(workspace_id: int) -> dict:
    """获取需求所有字段的中英文
    Args:
        workspace_id: 项目ID（必填）
    Returns: <str> 
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    ret = client.get_stories_fields_lable(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def get_stories_fields_info(workspace_id: int) -> dict:
    """获取需求所有字段及候选值，返回符合查询条件的所有需求字段及候选值。 部分字段为静态候选值，建议参考下方 "可选值说明"部分。其余动态字段（如：status(状态)/iteration_id(迭代)/categories(需求分类)），需要通过该接口获取对应的候选值（中英文映射）
    Args:
        workspace_id: 项目ID（必填）
    Returns: <str> 
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    ret = client.get_stories_fields_info(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def get_workspace_info(workspace_id: int) -> dict:
    """根据项目ID（workspace_id）获取项目信息，包含项目ID,项目名称,状态,创建时间,创建人等信息
    Args:
        workspace_id: 项目ID（必填）
    Returns: <str> 
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    ret = client.get_workspace_info(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def get_iterations(workspace_id: int, options: dict = None) -> dict:
    """根据项目ID（workspace_id）获取符合查询条件的所有迭代，例如可以通过名称获取迭代 id，然后作为iteration_id字段的值获取到迭代下的需求或缺陷

    !!! IMPORTANT !!!
    使用自定义字段(custom_field_*)前必须先调用 get_entity_custom_fields 获取字段配置
    示例流程：
    1. 调用 get_entity_custom_fields(workspace_id, {"entity_type": "iterations"}) 获取配置
    2. 使用返回的 custom_field_* 作为 options 参数的键名进行查询

    Args:
        workspace_id: 项目ID（必填）
        options:
            - id: ID
            - name: 标题
            - description: 描述
            - startdate: 开始时间
            - enddate: 结束时间
            - workitem_type_id: 迭代类别
            - completed: 完成时间
            - status: 状态（系统状态open/done，自定义状态可传中文）
            - custom_field_{number}: 自定义字段参数，必须先调用 get_entity_custom_fields 获取字段配置
                示例：{"custom_field_1": "值"}
            ...
    Returns: <str> 

    Example:
        >>> # 第一步：获取自定义字段配置
        >>> fields_config = get_entity_custom_fields(
        ...     workspace_id=123,
        ...     options={"entity_type": "iterations"}
        ... )
        >>> # 第二步：使用配置中的字段进行查询
        >>> iterations = get_iterations(
        ...     workspace_id=123,
        ...     options={"custom_field_1": "值"}
        ... )
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    if options:
        data.update(options)
    ret = client.get_iterations(data)
    fields_param = options.get('fields') if options else None
    if isinstance(ret, dict) and 'data' in ret and isinstance(ret['data'], list):
        ret['data'] = _filter_fields(ret['data'], fields_param)
    elif isinstance(ret, list):
        ret = _filter_fields(ret, fields_param)

    return json.dumps(ret, indent=2, ensure_ascii=False)
    
@mcp.tool()
def update_iteration(workspace_id: int, options: dict = None) -> dict:
    """更新迭代，返回更新迭代的数据
    Args:
        workspace_id: 项目ID（必填）
        options:
            - id: 迭代ID（必填）
            - current_user: 变更人（必填）
            - name: 标题
            - startdate: 开始时间
            - enddate: 结束时间
            - creator: 创建人
            - status: 状态
            ...
    Returns: <str> # 迭代所有字段数据的 json 格式
    Note: 迭代字段说明 status: 迭代状态，open: 开启，done: 已关闭
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    if options:
        data.update(options)
    ret = client.create_or_update_iteration(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def create_iteration(workspace_id: int, options: dict = None) -> dict:
    """创建迭代，返回创建迭代的数据
    Args:
        workspace_id: 项目ID（必填）
        options:
            - name: 迭代名称（必填）
            - startdate: 开始时间（必填）
            - enddate: 结束时间（必填）
            - creator: 创建人（必填）
            - status: 状态
            - description: 详细描述
            - parent_id: 默认为0. 指定当前迭代的父迭代，值为迭代ID，19位纯数字
            - label: 标签，标签不存在时将自动创建，多个以英文坚线分格
            ...
    Returns: <str> # 迭代所有字段数据的 json 格式，需要返回链接给用户，链接要可点击。
    Note: 迭代字段说明 status: 迭代状态，open: 开启，done: 已关闭
    Note: 迭代链接格式为 {tapd_base_url}/{workspace_id}/prong/iterations/card_view/{id}
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    if options:
        data.update(options)
    ret = client.create_or_update_iteration(data)

    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    return {
        "url": f'{tapd_base_url}/{workspace_id}/prong/iterations/card_view/{ret["data"]["Iteration"]["id"]}', # 返回给用户时，链接要可点击
        "data": json.dumps(ret, indent=2, ensure_ascii=False),
    }

@mcp.tool()
def get_related_bugs(workspace_id: int, options: dict = None) -> dict:
    """获取符合查询条件的所有需求关联的缺陷ID
    Args:
        workspace_id: 项目ID（必填）
        options:
            - story_id: ID（必填），支持多ID查询
    Returns: <str> 
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    if options:
        data.update(options)
    ret = client.get_related_bugs(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def entity_relations(workspace_id: int, options: dict = None) -> dict:
    """创建需求与缺陷关联关系
    Args:
        workspace_id: 所属TAPD项目ID（必填）
        options:
            - source_type: 关联关系源对象类型（story、bug）（必填）
            - target_type: 关联关系目标对象类型（story、bug）（必填）
            - source_id: 关联关系源对象id（必填）
            - target_id: 关联关系目标对象id（必填）
    Returns: <str> # 返回关联关系
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    if options:
        data.update(options)
    ret = client.add_entity_relations(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def create_or_update_tcases(workspace_id: int, options: dict = None) -> dict:
    """新建测试用例，返回新建测试用例的数据
    Args:
        workspace_id(int): 项目ID（必填）
        options: 可选参数字典，支持以下字段：
            - id(int): 测试用例 ID
            - name(str): 用例名称
            - category_id: 用例目录 enum('updating','abandon','normal')
            - status: 前置条件
            - precondition: 前置条件
            - steps: 用例步骤
            - expectation: 预期结果
            - type: 用例类型
            - priority: 用例等级
            - creator: 创建人
            - custom_field_*: 自定义字段参数，具体字段名通过接口 获取测试用例自定义字段配置获取
    Returns: <str>  # 新建测试用例的数据
    Note: 测试用例链接格式为 {tapd_base_url}/{workspace_id}/sparrow/tcase/view/{id}
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    if options:
        data.update(options)
    
    if 'precondition' in data:
        html = markdown.markdown(data["precondition"], extensions=['extra', 'codehilite', 'toc'])
        if '<p>' in html or '<h' in html or '<ul>' in html or '<ol>' in html:
            data["precondition"] = html
    if 'steps' in data:
        html = markdown.markdown(data["steps"], extensions=['extra', 'codehilite', 'toc'])
        if '<p>' in html or '<h' in html or '<ul>' in html or '<ol>' in html:
            data["steps"] = html
    if 'expectation' in data:
        html = markdown.markdown(data["expectation"], extensions=['extra', 'codehilite', 'toc'])
        if '<p>' in html or '<h' in html or '<ul>' in html or '<ol>' in html:
            data["expectation"] = html
    created_tcases = client.create_tcases(data)
    return json.dumps(created_tcases, indent=2, ensure_ascii=False)

@mcp.tool()
def get_tcases(workspace_id: int, options: dict = None) -> dict:
    """获取测试用例，返回符合查询条件的所有测试用例（分页显示，默认一页30条），默认返回 30 条。可通过传 limit 参数设置，最大取 200。也可以传 page 参数翻页
    Args:
        workspace_id(int): 项目ID（必填）
        options: 可选参数，支持以下字段：
            - id(str): ID
            - name: 用例名称
            - steps: 用例步骤
            - category_id(int): 用例目录
            - created: 创建时间
            - modifier: 最后修改人
            - modified: 最后修改时间
            - creator: 创建人
            - status: 用例状态 enum('updating','abandon','normal')	
            - precondition: 前置条件
            - expectation: 预期结果
            - type: 用例类型
            - priority: 用例等级
            - is_automated: 是否实现自动化
            - automation_type: 自动化测试类型
            - automation_platform: 自动化测试平台
            - is_serving: 是否上架
            - limit: 设置返回数量限制，默认为30
            - page: 返回当前数量限制下第N页的数据，默认为1（第一页）
            - order: 排序规则，规则：字段名 ASC或者DESC，然后 urlencode，如按创建时间逆序：order=created%20desc
            - fields: 设置获取的字段，多个字段间以','逗号隔开
            - custom_field_*: 自定义字段参数，具体字段名通过接口 获取
    Returns:  <str>  # 符合查询条件的测试用例所有字段数据的 json 格式
    Note: 测试用例链接格式为 {tapd_base_url}/{workspace_id}/sparrow/tcase/view/{id}
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    condition = {
        "workspace_id": workspace_id,
    }
    
    if options:
        condition.update(options)
    
    ret = client.get_tcases(condition)
    count_ret = client.get_tcases_count(condition)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    return {
        "base_url": {tapd_base_url}, # 返回给用户时，链接要可点击
        "data": json.dumps(ret, indent=2, ensure_ascii=False),
        "count": json.dumps(count_ret, indent=2, ensure_ascii=False)
    }

@mcp.tool()
def create_wiki(workspace_id: int, options: dict = None) -> dict:
    """新建Wiki，返回新建Wiki的数据
    Args:
        workspace_id(int): 项目ID（必填）
        options: 参数字典（包含必填字段，不可省略），支持以下字段：
            - name(str): 标题（必填）
            - markdown_description: wiki 内容，Markdown 格式
            - creator: 创建人（必填）
            - note: 备注
            - parent_wiki_id(str): 父wiki名
    Returns: <str>  # 新建wiki的数据
    Note: wiki链接格式为 {tapd_base_url}/{workspace_id}/markdown_wikis/show/#{id}
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    
    if options:
        data.update(options)
    created_wiki = client.create_wiki(data)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    return {
        "base_url": {tapd_base_url}, # 返回给用户时，链接要可点击
        "data": json.dumps(created_wiki, indent=2, ensure_ascii=False)
    }

@mcp.tool()
def update_wiki(workspace_id: int, options: dict = None) -> dict:
    """更新 TAPD wiki
    Args:
        workspace_id: 项目ID（必填）
        options: 参数字典（包含必填字段，不可省略），支持以下字段：
            - id: wiki ID（必填）
            - name(str): 标题
            - markdown_description: wiki 内容，Markdown 格式
            - note: 备注
            - parent_wiki_id(str): 父wiki名
    Returns: <str>  # wiki 所有字段数据的 json 格式
    Note: wiki链接格式为 {tapd_base_url}/{workspace_id}/markdown_wikis/show/#{id}
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    new_wiki = {
        "workspace_id": workspace_id,
    }
    
    if options:
        new_wiki.update(options)

    updated_wiki = client.create_wiki(new_wiki)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    return {
        "base_url": {tapd_base_url}, # 返回给用户时，链接要可点击
        "data": json.dumps(updated_wiki, indent=2, ensure_ascii=False)
    }

@mcp.tool()
def get_wiki(workspace_id: int, options: dict = None) -> dict:
    """获取 Wiki，返回符合查询条件的所有Wiki（分页显示，默认一页30条）
    Args:
        workspace_id(int): 项目ID（必填）
        options: 可选参数，支持以下字段：
            - id(int): id
            - name(str): 标题
            - modifier(str): 最后修改人
            - creator: 创建人
            - note: 备注
            - view_count: 浏览量
            - created: 创建时间
            - modified: 最后修改时间	
            - limit: 设置返回数量限制，默认为30	
            - page: 返回当前数量限制下第N页的数据，默认为1（第一页）	
            - order(str): 排序规则，规则：字段名 ASC或者DESC，然后 urlencode，如按创建时间逆序：order=created%20desc
            - fields: 设置获取的字段，多个字段间以','逗号隔开	
    Returns:  <str>  # 符合查询条件的wiki所有字段数据的 json 格式
    Note: wiki链接格式为 {tapd_base_url}/{workspace_id}/markdown_wikis/show/#{id}
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    condition = {
        "workspace_id": workspace_id,
    }
    
    if options:
        condition.update(options)
    
    ret = client.get_wiki(condition)
    count_ret = client.get_wiki_count(condition)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    return {
        "base_url": {tapd_base_url}, # 返回给用户时，链接要可点击
        "data": json.dumps(ret, indent=2, ensure_ascii=False),
        "count": json.dumps(count_ret, indent=2, ensure_ascii=False)
    }


@mcp.tool()
def get_todo(entity_type: str, user_nick: str = None) -> dict:
    """获取用户的待办，返回给用户的时候需要说明这是待办
    Args:
        entity_type: 业务对象类型，如需求story、缺陷bug、任务task（必填）
        user_nick: 用户昵称，如果为空则从环境变量CURRENT_USER_NICK获取
    Returns: 
        dict: 待办数据
    Raises:
        ValueError: 当user_nick为空且环境变量CURRENT_USER_NICK也为空时抛出
    """
    if not entity_type:
        raise ValueError("entity_type is required")
        
    if not user_nick:
        user_nick = os.getenv("CURRENT_USER_NICK")
    
    data = {
        "entity_type": entity_type,
        "user_nick": user_nick,
    }
    
    ret = client.get_todo(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def add_timesheets(workspace_id: int, options: dict = None) -> dict:
    """填写花费工时，返回新建花费工时的数据，可以先用get_timesheets查下owner是否在spentdate已经填写过花费，如果spentdate日期已存在花费记录，则需要调用update_timesheets进行更新
    Args:
        workspace_id: TAPD 项目 ID（必填）
        options:
            - entity_type: 对象类型，如story、task、bug（必填）
            - entity_id: 对象ID（必填）
            - timespent: 花费工时（必填）
            - owner: 花费创建人
            - spentdate: 花费日期，格式为 YYYY-MM-DD
            - memo: 花费描述
            - timeremain: 剩余工时
    Returns: <str> # 返回新建花费工时的数据
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    if options:
        data.update(options)
    user_nick = os.getenv("CURRENT_USER_NICK")
    if 'owner' not in data and user_nick:
        data['owner'] = user_nick
    else:
        if 'owner' not in data:
            raise ValueError("owner is required")
    ret = client.update_timesheets(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def update_timesheets(workspace_id: int, options: dict = None) -> dict:
    """更新花费工时，返回花费工时更新后的数据，每次只允许更新一条数据
    Args:
        workspace_id: TAPD 项目 ID（必填）
        options:
            - id: 花费id
            - timespent: 花费工时（必填）
            - timeremain: 剩余工时
            - memo: 花费描述
    Returns: <str> # 返回花费工时更新后的数据
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    if options:
        data.update(options)
    ret = client.update_timesheets(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def get_timesheets(workspace_id: int, options: dict = None) -> dict:
    """获取花费工时，返回符合查询条件的所有花费工时（分页显示，默认一页30条）默认返回 30 条。可通过传 limit 参数设置，最大取 200。也可以传 page 参数翻页
    Args:
        workspace_id: TAPD 项目 ID（必填）
        options:
            - id: 花费id
            - entity_type: 对象类型，如story、task、bug
            - entity_id: 对象ID
            - timespent: 花费工时
            - spentdate: 花费日期，格式为 YYYY-MM-DD
            - modified: 最后修改时间
            - owner: 花费创建人
            - include_parent_story_timesheet: 值=0不返回父需求的花费
            - created: 创建时间
            - memo: 花费描述
            - is_delete: 是否已删除。默认取 0，不返回已删除的工时记录。取 1 可以返回已删除的记录
            - limit: 设置返回数量限制，默认为30
            - page: 返回当前数量限制下第N页的数据，默认为1（第一页）
            - order: 排序规则，规则：字段名 ASC或者DESC，然后 urlencode，如按创建时间逆序：order=created%20desc
            - fields: 设置获取的字段，多个字段间以','逗号隔开
    Returns: <str> # 返回符合查询条件的所有花费工时（分页显示，默认一页30条
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    if options:
        data.update(options)
    ret = client.get_timesheets(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)


@mcp.tool()
def get_commit_msg(workspace_id: int, options: dict = None) -> dict:
    """获取需求/缺陷/任务的源码提交关键字，将 commit 和需求/缺陷/任务关联
    Args:
        workspace_id: TAPD 项目 ID（必填）
        options:
            - object_id: 对象ID（必填）
            - type: 对象类型，如story、task、bug（必填）
    Returns: <str> # 返回源码提交关键字，作为 commit message
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    if options:
        data.update(options)
    ret = client.get_scm_copy_keywords(data)

    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def get_release_info(workspace_id: int, options: dict = None) -> dict:
    """获取发布计划信息，返回符合查询条件的所有发布计划（分页显示），默认返回 30 条。可通过传 limit 参数设置，最大取 200。也可以传 page 参数翻页
    如果要获取某个日期发布的需求，可以先根据startdate和enddate的日期范围，然后调用 get_release_info 获取发布计划ID，然后调用 get_stories 获取需求
    Args:
        workspace_id: TAPD 项目 ID（必填）
        options:
            - id: 发布计划ID，支持多 ID 查询
            - name: 标题，支持模糊匹配
            - description: 详细描述
            - startdate: 开始时间
            - enddate: 结束时间
            - creator: 创建人
            - created: 创建时间
            - modified: 最后修改时间
            - status: 状态，done 或者 open
            - order: 排序规则，规则：字段名 ASC或者DESC，然后 urlencode，如按创建时间逆序：order=created%20desc
            - fields: 设置获取的字段，多个字段间以','逗号隔开
    Returns: <str> # 返回发布计划信息
    """
    if not workspace_id:
        user_nick = os.getenv("CURRENT_USER_NICK")
        condition = {
            "nick": user_nick
        }
        ret = client.get_user_participant_projects(condition)
        return json.dumps(ret, indent=2, ensure_ascii=False)
    data = {
        "workspace_id": workspace_id,
    }
    if options:
        data.update(options)
   
    ret = client.get_release_info(data)
    return json.dumps(ret, indent=2, ensure_ascii=False)

@mcp.tool()
def send_qiwei_message(msg: str) -> dict:
    """发送信息到企业微信群
    Args:
        msg: 推送的企业微信的信息，Markdown 格式（必填）
    Returns: <str> 
    """
    data = {
        "msg": msg,
    }
    return client.send_message(data)


@mcp.tool()
def upload_image(workspace_id: int, image_path: str) -> str:
    """上传本机图片到 TAPD 图床，返回可直接内嵌到描述/评论里的图片地址。

    多数情况下不需要单独调用本工具：create_bug / create_story_or_task /
    create_comments 等工具的 image_url、image_urls 参数已经可以直接接收本机图片
    绝对路径并自动完成上传。只有在需要先拿到地址、自己拼 description 时才用它。

    Args:
        workspace_id: 项目ID（必填）
        image_path: 本机图片绝对路径（必填）。限 5MB 以内，
                    支持 PNG/GIF/JPG/JPEG/BMP，一次一张。
                    要求本 MCP server 与图片在同一台机器上。
    Returns:
        {
            "image_src": <str>,   # 形如 /tfl/pictures/202608/api_xxx.png
            "html_code": <str>    # 形如 <img src="/tfl/pictures/..."/>，可直接写入 description
        }
    """
    if not os.path.isfile(image_path):
        return json.dumps(
            {"status": 0, "info": f"找不到文件：{image_path}"},
            indent=2, ensure_ascii=False
        )

    size = os.path.getsize(image_path)
    if size > _MAX_UPLOAD_BYTES:
        return json.dumps(
            {"status": 0, "info": f"图片为 {size} 字节，超过 TAPD 的 5MB 限制，请先压缩"},
            indent=2, ensure_ascii=False
        )

    ret = client.upload_image(workspace_id, image_path)
    return json.dumps(ret, indent=2, ensure_ascii=False)


@mcp.tool()
def attach_files_to_ticket(workspace_id: int, entry_id: str, entity_type: str, attachments: list[dict]) -> dict:
    """给已有需求/缺陷补普通附件。entity_type 为 story 或 bug。
    attachments=[{"path":"本地路径或 media://inbound/...","filename":"原始文件名.zip"}]。
    用微信文件索引中的 localPath 和原始文件名，不要把 UUID 名称当原始名称。
    默认 display="auto"：视频会附加正文播放器，其他文件仅附件；display="attachment" 强制仅附件（图片也可）。
    不解压 ZIP 或解析 EML。自动内嵌图片请通过 update_story_or_task/update_bug 的 image_urls。
    返回每个附件的状态/ID；部分失败时只处理失败项，不重新建票或重传成功项。
    """
    if entity_type not in {"story", "bug"} or not str(entry_id).isdigit():
        raise ValueError("需要有效的 entry_id 和 story/bug 类型")
    data = {"attachments": attachments}
    files = prepare_attachments(data, _resolve_media_ref)
    if data.get("image_urls"):
        raise ValueError("图片应通过 update_story_or_task/update_bug 的 image_urls 内嵌；本次未上传")
    results = upload_files(client, workspace_id, entity_type, entry_id, files, render_video=_render_tapd_video_html)
    return {"workspace_id": workspace_id, "entry_id": entry_id, "attachment_results": results,
            "ok": all(r["status"] == "uploaded" and r.get("preview_status") != "failed" for r in results)}


def start_mcp_server():
    config = AppConfig()
    if config.mode == "http" or config.mode == "sse":
        mcp.run(transport="sse")
    elif config.mode == "streamable-http":
        mcp.run(transport="streamable-http")
    else:
        mcp.run()

def main():
    start_mcp_server()

if __name__ == "__main__":
    start_mcp_server()
