import os
import json
import requests
import markdown
from mcp.server.fastmcp import FastMCP
from mcp_server_tapd.tapd import TAPDClient
from mcp_server_tapd.app_config import AppConfig

config = AppConfig()
if config.mode == "http" or config.mode == "sse" or config.mode == "streamable-http":
    mcp = FastMCP("mcp-tapd", host=config.host, port=config.port)
else:
    mcp = FastMCP("mcp-tapd")
client = TAPDClient()

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
        options: 可选参数字典，支持以下字段：
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
        options: 可选参数字典，支持以下字段：
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
        options: 可选参数字典，支持以下字段：
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
    """更新 TAPD 需求或任务，任务的状态只有三个，open: 未开始，progressing: 进行中，done: 已完成，如果是流转需求的状态，需要调用工具先获取状态信息get_workflows_all_transitions
    TAPD API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/story/update_story.html
    Args:
        workspace_id: 项目ID（必填）
        options: 可选参数字典，支持以下字段：
            - entity_type: 类型，需求/工作项 stories 或 任务 tasks（必填）
            - id: 需求或任务ID（必填）
            - name: 标题
            - v_status: 状态别名(支持传入中文状态名称)，默认用这个字段
            - status: 状态，原名
            - priority_label: 优先级，需要先检查get_stories_fields_info是否配置了候选值，如果不存在则使用默认的候选值 High => 高、Middle => 中、Low => 低、Nice To Have => 无关紧要，不使用 priority 字段
            - description: 描述
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
    if 'description' in new_story:
        html = markdown.markdown(new_story["description"], extensions=['extra', 'codehilite', 'toc'])
        if '<p>' in html or '<h' in html or '<ul>' in html or '<ol>' in html:
            new_story["description"] = html
    created_story = client.create_or_update_story(new_story)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    
    entity_type = options['entity_type']
    url_template = get_story_or_task_url_template(workspace_id, entity_type, tapd_base_url)
    return json.dumps({
        "url_template": url_template,
        "data": json.dumps(created_story, indent=2, ensure_ascii=False)
    }, ensure_ascii=False, indent=2)

@mcp.tool()
def create_story_or_task(workspace_id: int, name: str, options: dict = None) -> dict:
    """创建 TAPD 需求或任务
    Args:
        workspace_id: 项目ID（必填）
        name: 需求或任务的标题（必填）
        options: 可选参数字典，支持以下字段：
            - entity_type: 类型，需求/工作项 stories 或 任务 tasks（必填）
            - priority_label: 优先级，需要先检查get_stories_fields_info是否配置了候选值，如果不存在则使用默认的候选值 High => 高、Middle => 中、Low => 低、Nice To Have => 无关紧要，不使用 priority 字段
            - description: 描述
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
    
    if 'description' in data:
        html = markdown.markdown(data["description"], extensions=['extra', 'codehilite', 'toc'])
        if '<p>' in html or '<h' in html or '<ul>' in html or '<ol>' in html:
            data["description"] = html

    created_story = client.create_or_update_story(data)
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
    """更新 TAPD 缺陷
    TAPD API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/update_bug.html
    Args:
        workspace_id: 项目ID（必填）
        options: 可选参数字典，支持以下字段：
            - id: 缺陷ID（必填）
            - title: 标题
            - description: 描述
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

    if 'description' in new_bug:
        html = markdown.markdown(new_bug["description"], extensions=['extra', 'codehilite', 'toc'])
        if '<p>' in html or '<h' in html or '<ul>' in html or '<ol>' in html:
            new_bug["description"] = html
    user_nick = os.getenv("CURRENT_USER_NICK")
    if 'lastmodify' not in new_bug and user_nick:
        new_bug['lastmodify'] = user_nick
    created_bug = client.create_or_update_bug(new_bug)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    return {
        "base_url": tapd_base_url, # 返回给用户时，链接要可点击
        "data": json.dumps(created_bug, indent=2, ensure_ascii=False)
    }

@mcp.tool()
def create_bug(workspace_id: int, title: str, options: dict = None) -> dict:
    """创建 TAPD 缺陷
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
            "data": <str>,  # 所有字段数据的 json 格式
            "url": <str>  #  url，返回给用户时，链接要可点击
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
    
    if 'description' in new_bug:
        html = markdown.markdown(new_bug["description"], extensions=['extra', 'codehilite', 'toc'])
        if '<p>' in html or '<h' in html or '<ul>' in html or '<ol>' in html:
            new_bug["description"] = html
    user_nick = os.getenv("CURRENT_USER_NICK")
    if 'reporter' not in new_bug and user_nick:
        new_bug['reporter'] = user_nick
    ret = client.create_or_update_bug(new_bug)
    config = AppConfig()
    if config.tapd_base_url is None:
        tapd_base_url = os.getenv("TAPD_BASE_URL")
    else:
        tapd_base_url = config.tapd_base_url
    return {
        "url": f'{tapd_base_url}/{workspace_id}/bugtrace/bugs/view/{ret["data"]["Bug"]["id"]}', # 返回给用户时，链接要可点击
        "data": json.dumps(ret, indent=2, ensure_ascii=False),
    }

@mcp.tool()
def create_comments(workspace_id: int, options: dict = None) -> dict:
    """添加评论
    TAPD API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/comment/add_comment.html
    Args:
        workspace_id: 项目ID（必填）
        options: 可选参数字典，支持以下字段：
            - entry_id: 评论所依附的业务对象实体id（必填）
            - entry_type: 评论类型（取值： bug、 bug_remark （流转缺陷时候的评论）、 stories、 tasks 。）（必填）
            - author: 评论人（必填）
            - description: 内容（必填）
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
    
    if 'description' in data:
        html = markdown.markdown(data["description"], extensions=['extra', 'codehilite', 'toc'])
        if '<p>' in html or '<h' in html or '<ul>' in html or '<ol>' in html:
            data["description"] = html
    created_story = client.create_comments(data)
    return json.dumps(created_story, indent=2, ensure_ascii=False)

@mcp.tool()
def update_comments(workspace_id: int, options: dict = None) -> dict:
    """更新评论，返回更新评论的数据
    Args:
        workspace_id: 项目ID（必填）
        options: 可选参数字典，支持以下字段：
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
    
    if 'description' in data:
        html = markdown.markdown(data["description"], extensions=['extra', 'codehilite', 'toc'])
        if '<p>' in html or '<h' in html or '<ul>' in html or '<ol>' in html:
            data["description"] = html
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
        options: 可选参数字典，支持以下字段：
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
        options: 可选参数字典，支持以下字段：
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