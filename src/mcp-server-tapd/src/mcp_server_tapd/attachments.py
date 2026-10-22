"""Generic ticket attachments. Files remain files; only images become rich text."""
import os

MAX_BYTES = 150 * 1024 * 1024
IMAGE_EXTENSIONS = {'.png', '.jpg', '.jpeg', '.gif', '.webp', '.bmp'}


def prepare_attachments(data, resolve_path):
    raw = data.pop('attachments', [])
    if not isinstance(raw, list):
        raise ValueError('attachments 必须是数组，元素为 {path, filename}')
    if len(raw) > 20:
        raise ValueError('单次最多上传 20 个附件')
    pending = []
    images = []
    for item in raw:
        if not isinstance(item, dict) or not isinstance(item.get('path'), str):
            raise ValueError('附件必须提供本地 path 或 media:// 引用')
        p = resolve_path(item['path'])
        if '://' in p or not os.path.isfile(p):
            raise ValueError('附件必须是已下载的本地文件')
        name = item.get('filename') or os.path.basename(p)
        if not isinstance(name, str) or not name.strip() or any(c in name for c in '/\\\r\n\x00') or name in {'.', '..'}:
            raise ValueError('filename 必须是原始文件名，不能含目录或换行')
        if os.path.getsize(p) > MAX_BYTES:
            raise ValueError(f'附件 {name} 超过 150MB')
        if os.path.splitext(name)[1].lower() in IMAGE_EXTENSIONS:
            if os.path.splitext(p)[1].lower() not in IMAGE_EXTENSIONS:
                raise ValueError(f'图片 {name} 的本地路径没有图片扩展名，请先保存为正确图片类型后内嵌')
            images.append(p)
        else:
            pending.append({'path': p, 'filename': name})
    if raw and data.get('entity_type') == 'tasks':
        raise ValueError('Task 附件尚未验证，请使用需求 stories 或缺陷 bug；未创建/更新工作项')
    if images:
        existing = data.get('image_urls') or []
        if isinstance(existing, str): existing = [existing]
        data['image_urls'] = list(existing) + images
    return pending


def upload_files(client, workspace_id, entity, entry_id, files):
    if entity not in {'story', 'bug'}:
        raise ValueError('附件挂载只支持 story / bug')
    results = []
    for f in files:
        result = {'filename': f['filename'], 'status': 'failed'}
        try:
            if not entry_id:
                raise ValueError('missing entry id')
            ret = client.upload_attachment(workspace_id, f['path'], entity, entry_id,
                                           filename=f['filename'])
            attachment_id = ((ret.get('data') or {}).get('Attachment') or {}).get('id')
            if ret.get('status') == 1 and attachment_id:
                result.update(status='uploaded', attachment_id=str(attachment_id))
            else:
                result['error'] = 'TAPD 未确认上传成功或未返回附件 ID'
        except Exception as exc:
            # No credentials, file contents or server response dumps in error output.
            result['error'] = f'附件未确认上传成功（{type(exc).__name__}）；请先核对票内附件再重试'
        results.append(result)
    return results


def finalize_attachments(client, workspace_id, entity, response, files, warnings):
    if not files:
        return response
    key = 'Bug' if entity == 'bug' else 'Story'
    entry_id = ((response.get('data') or {}).get(key) or {}).get('id') if isinstance(response, dict) else None
    if not isinstance(response, dict):
        warnings.append('工作项响应格式不正确，附件未上传')
        return response
    # A failed create/update must not cause a secondary write.
    if response.get('status') != 1:
        entry_id = None
    results = upload_files(client, workspace_id, entity, entry_id, files)
    response['attachment_results'] = results
    for r in results:
        if r['status'] != 'uploaded': warnings.append(f"附件 {r['filename']}：{r['error']}")
    return response
