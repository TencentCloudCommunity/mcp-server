import os
import requests
import json
from typing import Dict, Any, Optional
from base64 import b64encode
from mcp_server_tapd.app_config import AppConfig

class TAPDClient:
    def __init__(self):
        """
        初始化 TAPD API 客户端
        使用配置文件中的常量进行初始化
        """
        config = AppConfig()
        self.base_url = config.api_base_url
        self.bot_url = config.bot_url
        if config.access_token:  # 优先使用 token
            self.headers = {
                "Authorization": f"Bearer {config.access_token}",
                "Content-Type": "application/json",
                "Via": "mcp"
            }
            self.nick = self.get_user_info()  # 初始化时获取nick
        else:
            auth_str = f"{config.api_user}:{config.api_password}"
            self.headers = {
                "Authorization": f"Basic {b64encode(auth_str.encode()).decode()}",
                "Content-Type": "application/json",
                "Via": "mcp"
            }
    
    
    def _make_request(self, method: str, endpoint: str, params: Optional[Dict] = None, data: Optional[Dict] = None) -> Dict:
        """
        发送 API 请求的通用方法
        """
         # 获取环境变量
        if self.base_url is None:
            access_token = os.getenv("TAPD_ACCESS_TOKEN")
            api_user = os.getenv("TAPD_API_USER")
            api_password = os.getenv("TAPD_API_PASSWORD")
            self.base_url = os.getenv("TAPD_API_BASE_URL")
            self.bot_url = os.getenv("BOT_URL")
            if access_token:
                self.headers = {
                    "Authorization": f"Bearer {access_token}",
                    "Content-Type": "application/json",
                    "Via": "mcp"
                }
                self.nick = self.get_user_info()  # 初始化时获取nick
            else:
                auth_str = f"{api_user}:{api_password}"
                self.headers = {
                    "Authorization": f"Basic {b64encode(auth_str.encode()).decode()}",
                    "Content-Type": "application/json",
                    "Via": "mcp"
                }
        url = f"{self.base_url}/{endpoint}"
        separator = '&' if '?' in url else '?'
        url = f"{url}{separator}s=mcp"
        response = requests.request(
            method=method,
            url=url,
            headers=self.headers,
            params=params,
            json=data,
            timeout=30
        )
        response.raise_for_status()
        # ret_json = response.json()
        # if ret_json['status'] == 401:
        #     ret = "检查是否有传账号密码以及核实是否正确"
        # elif ret_json['status'] == 403 or ret_json['status'] == 404:
        #     ret = ret_json['info']
        # else:
        #     ret = ret_json['info']
            
        
        return response.json()

    def get_stories(self, params: Optional[Dict] = None) -> Dict:
        """
        获取需求或任务
        """
        entity_type = "stories"
        if 'entity_type' in params and params['entity_type'] == "tasks":
            entity_type = "tasks"  
        default_params = {
            "page": 1,
            "limit": 10
        }
        if params:
            default_params.update(params)
        # 新增：处理id为9位以下数字的情况，支持多个id（逗号分隔）
        if 'id' in default_params and 'workspace_id' in default_params:
            id_val = str(default_params['id'])
            workspace_id = str(default_params['workspace_id'])
            pre_id = "11" if self.is_cloud_env() else "10"
            def to_long_id(single_id):
                single_id = single_id.strip()
                if single_id.isdigit() and len(single_id) <= 9:
                    padded_id = single_id.zfill(9)
                    return f"{pre_id}{workspace_id}{padded_id}"
                return single_id
            if "," in id_val:
                id_list = id_val.split(",")
                long_id_list = [to_long_id(i) for i in id_list]
                default_params['id'] = ",".join(long_id_list)
            else:
                default_params['id'] = to_long_id(id_val)
        
        return self._make_request("GET", entity_type, params=default_params)

    def get_comments(self, params: Optional[Dict] = None) -> Dict:
        """
        获取评论
        """
        default_params = {
            "page": 1,
            "limit": 10
        }
        if params:
            default_params.update(params)
            
        return self._make_request("GET", "comments", params=default_params)

    def create_or_update_story(self, params: Dict[str, Any]) -> Dict:
        """
        创建/更新需求或任务
        """
        entity_type = "stories"
        if 'entity_type' in params:
            if params['entity_type'] == "tasks":
                entity_type = "tasks"  
            if 'id' in params:
                del params['entity_type']
        
         # 新增：处理id为9位以下数字的情况
        if 'id' in params and 'workspace_id' in params:
            id_val = str(params['id'])
            if id_val.isdigit() and len(id_val) <= 9:
                padded_id = id_val.zfill(9)
                workspace_id = str(params['workspace_id'])
                pre_id = "11" if self.is_cloud_env() else "10"
                long_id = f"{pre_id}{workspace_id}{padded_id}"
                params['id'] = long_id
        if getattr(self, "nick", None):
            if 'id' in params:
                params['current_user'] = self.nick
            else:
                params['creator'] = self.nick

        return self._make_request("POST", entity_type, data=params)

    def get_story_count(self, params: Optional[Dict] = None) -> Dict:
        """
        获取需求数量
        """
        entity_type = "stories"
        if 'entity_type' in params and params['entity_type'] == "tasks":
            entity_type = "tasks"
        return self._make_request("GET", entity_type + "/count", params=params)

    def get_entity_custom_fields(self, params: Optional[Dict] = None) -> Dict:
        """
        获取需求或者任务或者迭代或者测试用例的自定义字段配置
        """
        workspace_id = params['workspace_id']
        
        # 直接获取entity_type，默认为stories
        entity_type = params.get('entity_type', 'stories')
        
        return self._make_request(
            "GET", 
            f"{entity_type}/custom_fields_settings?workspace_id={workspace_id}"
        )

    def get_bug(self, params: Optional[Dict] = None) -> Dict:
        """
        获取缺陷
        """
        default_params = {
            "page": 1,
            "limit": 10
        }
        if params:
            default_params.update(params)
        
        if 'id' in default_params and 'workspace_id' in default_params:
            id_val = str(default_params['id'])
            workspace_id = str(default_params['workspace_id'])
            pre_id = "11" if self.is_cloud_env() else "10"
            def to_long_id(single_id):
                single_id = single_id.strip()
                if single_id.isdigit() and len(single_id) <= 9:
                    padded_id = single_id.zfill(9)
                    return f"{pre_id}{workspace_id}{padded_id}"
                return single_id
            if "," in id_val:
                id_list = id_val.split(",")
                long_id_list = [to_long_id(i) for i in id_list]
                default_params['id'] = ",".join(long_id_list)
            else:
                default_params['id'] = to_long_id(id_val)
            
        return self._make_request("GET", "bugs", params=default_params)

    def create_or_update_bug(self, data: Dict[str, Any]) -> Dict:
        """
        创建或更新缺陷
        """
        if 'id' in data and 'workspace_id' in data:
            id_val = str(data['id'])
            if id_val.isdigit() and len(id_val) <= 9:
                padded_id = id_val.zfill(9)
                workspace_id = str(data['workspace_id'])
                pre_id = "11" if self.is_cloud_env() else "10"
                long_id = f"{pre_id}{workspace_id}{padded_id}"
                data['id'] = long_id
        if getattr(self, "nick", None):
            if 'id' in data:
                data['current_user'] = self.nick
            else:
                data['reporter'] = self.nick

        return self._make_request("POST", "bugs", data=data)

    def get_bug_count(self, params: Optional[Dict] = None) -> Dict:
        """
        获取缺陷数量
        """
        return self._make_request("GET", "bugs/count", params=params)

    def get_bug_custom_fields(self, workspace_id: int) -> Dict:
        """
        获取缺陷自定义字段配置
        """
        params = {"workspace_id": workspace_id}
        return self._make_request("GET", f"bugs/custom_fields_settings?workspace_id={workspace_id}", params=params)
    
    def create_comments(self, data: Dict[str, Any]) -> Dict:
        """
        新建评论
        """
        # 新增：处理id为9位以下数字的情况，支持多个id（逗号分隔）
        if 'entry_id' in data and 'workspace_id' in data:
            id_val = str(data['entry_id'])
            workspace_id = str(data['workspace_id'])
            pre_id = "11" if self.is_cloud_env() else "10"
            def to_long_id(single_id):
                single_id = single_id.strip()
                if single_id.isdigit() and len(single_id) <= 9:
                    padded_id = single_id.zfill(9)
                    return f"{pre_id}{workspace_id}{padded_id}"
                return single_id
            data['entry_id'] = to_long_id(id_val)
        
        if getattr(self, "nick", None):
            if 'id' in data:
                data['change_creator'] = self.nick
            else:
                data['author'] = self.nick

        return self._make_request("POST", "comments", data=data)
    
    def create_wiki(self, data: Dict[str, Any]) -> Dict:
        """
        新建wiki
        """
        if getattr(self, "nick", None):
            if 'id' in data:
                data['modifier'] = self.nick
            else:
                data['creator'] = self.nick
        return self._make_request("POST", "tapd_wikis", data=data)
        
    def create_tcases(self, data: Dict[str, Any]) -> Dict:
        """
        新建测试用例
        """
        if getattr(self, "nick", None):
            if 'id' in data:
                data['modifier'] = self.nick
            else:
                data['creator'] = self.nick
        return self._make_request("POST", "tcases", data=data)

    def get_tcases(self, params: Optional[Dict] = None) -> Dict:
        """
        获取测试用例
        """
        default_params = {
            "page": 1,
            "limit": 30
        }
        if params:
            default_params.update(params)
            
        return self._make_request("GET", "tcases", params=default_params)

    def get_tcases_count(self, params: Optional[Dict] = None) -> Dict:
        """
        获取测试用例数量
        """
        return self._make_request("GET", "tcases/count", params=params)
    
    def get_tcases_custom_fields_settings(self, data: Dict[str, Any]) -> Dict:
        """
        获取测试用例自定义字段配置
        """
        workspace_id = data['workspace_id']
        return self._make_request("GET", f"tcases/custom_fields_settings?workspace_id={workspace_id}", data=data)

    def get_wiki(self, params: Optional[Dict] = None) -> Dict:
        """
        获取wiki
        """
        default_params = {
            "page": 1,
            "limit": 30
        }
        if params:
            default_params.update(params)
            
        return self._make_request("GET", "tapd_wikis", params=default_params)

    def get_wiki_count(self, params: Optional[Dict] = None) -> Dict:
        """
        获取测试用例数量
        """
        return self._make_request("GET", "tapd_wikis/count", params=params)
    
    def get_workflows_all_transitions(self, data: Dict[str, Any]) -> Dict:
        """
        获取工作流流转细则
        """
        system = data['system']
        workspace_id = data['workspace_id']
        params = f"?workspace_id={workspace_id}&system={system}"
        if 'workitem_type_id' in data:
            workitem_type_id = data['workitem_type_id']
            params += f"&workitem_type_id={workitem_type_id}"
        return self._make_request("GET", f"workflows/all_transitions{params}")
    
    def get_workflows_status_map(self, data: Dict[str, Any]) -> Dict:
        """
        获取工作流状态中英文名对应关系
        """
        system = data['system']
        workspace_id = data['workspace_id']
        params = f"?workspace_id={workspace_id}&system={system}"
        if 'workitem_type_id' in data:
            workitem_type_id = data['workitem_type_id']
            params += f"&workitem_type_id={workitem_type_id}"
        return self._make_request("GET", f"workflows/status_map{params}")
    
    def get_workitem_types(self, data: Dict[str, Any]) -> Dict:
        """
        返回符合查询条件的所有需求类别（分页显示，默认一页30条）
        """
        workspace_id = data['workspace_id']
        return self._make_request("GET", f"workitem_types?workspace_id={workspace_id}", data=data)
   
    def get_workflows_last_steps(self, data: Dict[str, Any]) -> Dict:
        """
        获取工作流结束状态
        """
        system = data['system']
        workspace_id = data['workspace_id']
        params = f"?workspace_id={workspace_id}&system={system}"
        if 'workitem_type_id' in data:
            workitem_type_id = data['workitem_type_id']
            params += f"&workitem_type_id={workitem_type_id}"
        if 'type' in data:
            type = data['type']
            params += f"&type={type}"
        return self._make_request("GET", f"workflows/last_steps{params}", data=data)
    
    def get_stories_custom_fields_settings(self, data: Dict[str, Any]) -> Dict:
        """
        获取需求自定义字段配置
        """
        workspace_id = data['workspace_id']
        return self._make_request("GET", f"stories/custom_fields_settings?workspace_id={workspace_id}", data=data)
   
    def get_stories_fields_lable(self, data: Dict[str, Any]) -> Dict:
        """
        获取需求需求所有字段的中英文
        """
        workspace_id = data['workspace_id']
        return self._make_request("GET", f"stories/get_fields_lable?workspace_id={workspace_id}")
    
    def get_stories_fields_info(self, data: Dict[str, Any]) -> Dict:
        """
        获取需求所有字段及候选值，返回符合查询条件的所有需求字段及候选值
        """
        workspace_id = data['workspace_id']
        return self._make_request("GET", f"stories/get_fields_info?workspace_id={workspace_id}")
    
    def get_workspace_info(self, data: Dict[str, Any]) -> Dict:
        """
        根据项目ID（workspace_id）获取项目信息，包含项目ID,项目名称,状态,创建时间,创建人等信息
        """
        workspace_id = data['workspace_id']
        return self._make_request("GET", f"workspaces/get_workspace_info?workspace_id={workspace_id}", data=data)
    
    def get_iterations(self, data: Dict[str, Any]) -> Dict:
        """
        符合查询条件的所有迭代（分页显示，默认一页30条）
        """
        workspace_id = data['workspace_id']
        params = f"?workspace_id={workspace_id}"
        if 'id' in data:
            id = data['id']
            params += f"&id={id}"
        if 'name' in data:
            name = data['name']
            params += f"&name={name}"
        return self._make_request("GET", f"iterations{params}", params=data)
    
    def get_related_bugs(self, data: Dict[str, Any]) -> Dict:
        """
        返回符合查询条件的所有需求关联的缺陷ID
        """
        return self._make_request("GET", f"stories/get_related_bugs", params=data)

    def create_or_update_iteration(self, data: Dict[str, Any]) -> Dict:
        """
        创建迭代，返回创建迭代的数据
        """
        if getattr(self, "nick", None):
            if 'id' in data:
                data['current_user'] = self.nick
            else:
                data['creator'] = self.nick
        return self._make_request("POST", "iterations", data=data)
    def add_entity_relations(self, data: Dict[str, Any]) -> Dict:
        """
        创建需求与缺陷关联关系
        """
        return self._make_request("POST", "relations", data=data)
   
    def get_todo(self, data: Dict[str, Any]) -> Dict:
        """
        获取待办
        """
        entity_type = data['entity_type']
        if getattr(self, "nick", None):
            user_nick = self.nick
        return self._make_request("GET", f"users/todo/{user_nick}/{entity_type}")
    
    def update_timesheets(self, data: Dict[str, Any]) -> Dict:
        """
        新建或更新花费工时
        """
        if getattr(self, "nick", None):
            data['owner'] = self.nick
        return self._make_request("POST", "timesheets", data=data)

    def get_timesheets(self, data: Dict[str, Any]) -> Dict:
        """
        获取花费工时
        """
        return self._make_request("GET", "timesheets", params=data)

    def get_scm_copy_keywords(self, data: Dict[str, Any]) -> Dict:
        """
        获取源码提交关键字
        """
        # 新增：处理id为9位以下数字的情况，支持多个id（逗号分隔）
        if 'object_id' in data and 'workspace_id' in data:
            id_val = str(data['object_id'])
            workspace_id = str(data['workspace_id'])
            pre_id = "11" if self.is_cloud_env() else "10"
            def to_long_id(single_id):
                single_id = single_id.strip()
                if single_id.isdigit() and len(single_id) <= 9:
                    padded_id = single_id.zfill(9)
                    return f"{pre_id}{workspace_id}{padded_id}"
                return single_id
            data['object_id'] = to_long_id(id_val)
        return self._make_request("GET", "svn_commits/get_scm_copy_keywords", params=data)
    
    def get_release_info(self, data: Dict[str, Any]) -> Dict:
        """
        返回符合查询条件的所有发布计划（分页显示，默认一页30条）
        """
        return self._make_request("GET", "releases", params=data)
    
    def get_category_id(self, data: Dict[str, Any]) -> Dict:
        """
        获取需求分类 ID
        """
        return self._make_request("GET", "story_categories", params=data)
    
    def get_user_participant_projects(self, data: Dict[str, Any]) -> Dict:
        """
        获取用户参与的项目列表
        """
        return self._make_request("GET", "workspaces/user_participant_projects", params=data)
    
    def check_mini_project(self, workspace_id: int) -> bool:
        """
        判断是否是轻协作项目
        """
        data = {
            "workspace_id": workspace_id,
        }
    
        ret = self.get_workspace_info(data)
        return ret.get('data', {}).get('Workspace', {}).get('category') == 'mini_project'
    
    def send_message(self, data: Dict[str, Any]):
        # 判断消息内容是否包含@字符，决定使用的消息类型
        if '@' in data['msg']:
            chat_data = {
                'msgtype': 'markdown',
                'markdown': {
                    'content': data['msg']
                }
            }
        else:
            chat_data = {
                'msgtype': 'markdown_v2',
                'markdown_v2': {
                    'content': data['msg']
                }
            }
        headers = {'Content-Type': 'application/json'}
        response = requests.post(
            url=self.bot_url,
            headers=headers,
            json=chat_data,
            timeout=500
        )
        
        return response.text
   
    def is_cloud_env(self) -> bool:
        """
        判断当前是否为CLOUD环境（base_url包含www.tapd.cn即为CLOUD环境）
        """
        if self.base_url is None:
            self.base_url = os.getenv("TAPD_API_BASE_URL")
        
        return self.base_url is not None and 'api.tapd.cn' in self.base_url
    
    def get_user_info(self) -> bool:
        """
        获取用户信息
        """
        response = self._make_request("GET", "users/info")
        nick = response.get("data", {}).get("nick") if response.get("data") else None
        return nick
    
   